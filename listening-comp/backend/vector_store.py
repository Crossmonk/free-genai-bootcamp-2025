import chromadb
from chromadb.utils import embedding_functions
import json
import os
from typing import Dict, List, Optional

# Import necessary libraries for Ollama embeddings
# from langchain_community.embeddings import OllamaEmbeddings # Removed: Deprecated
from langchain_ollama import OllamaEmbeddings # Updated: Import from langchain_ollama

# Use a smaller model for faster processing. You can change this if needed.
DEFAULT_MODEL_NAME = "all-MiniLM-L6-v2"  # Keeping this for comparison or if Ollama fails

class OllamaEmbeddingFunction(embedding_functions.EmbeddingFunction):
    def __init__(self, model_name: str = "llama2", base_url: str = "http://localhost:11434", device: str = "cpu"):
        """
        Initializes the embedding function with an Ollama model.

        Args:
            model_name (str): The name of the Ollama model. Defaults to "llama2"
            base_url (str): The base URL of the Ollama API. Defaults to "http://localhost:11434"
            device (str): Not used for Ollama, but kept for consistency.
        """
        self.model = OllamaEmbeddings(model=model_name, base_url=base_url)

    def __call__(self, texts: List[str]) -> List[List[float]]:
        """
        Generates embeddings for a list of texts using Ollama.

        Args:
            texts (List[str]): The list of texts to embed.

        Returns:
            List[List[float]]: The list of embeddings.
        """
        try:
            embeddings = self.model.embed_documents(texts)
            return embeddings
        except Exception as e:
            print(f"Error generating embeddings with Ollama: {str(e)}")
            # Fallback to a zero vector if Ollama fails
            return [[0.0] * 384] * len(texts)  # Fallback: 384-dimensional zero vector


class SpacySentenceTransformerEmbeddingFunction(embedding_functions.EmbeddingFunction): #Added as fallback
    def __init__(self, model_name: str = DEFAULT_MODEL_NAME, device: str = "cpu"):
        """
        Initializes the embedding function with a sentence-transformers model.

        Args:
            model_name (str): The name of the sentence-transformers model.
            device (str): The device to use for the model ("cpu" or "cuda").
        """
        from sentence_transformers import SentenceTransformer
        self.model = SentenceTransformer(model_name, device=device)

    def __call__(self, texts: List[str]) -> List[List[float]]:
        """
        Generates embeddings for a list of texts.

        Args:
            texts (List[str]): The list of texts to embed.

        Returns:
            List[List[float]]: The list of embeddings.
        """
        try:
            embeddings = self.model.encode(texts).tolist()
            return embeddings
        except Exception as e:
            print(f"Error generating embeddings: {str(e)}")
            return [[0.0] * 384] * len(texts)  # Fallback: 384-dimensional zero vector (adjust if you use a different model)


class QuestionVectorStore:
    def __init__(self, persist_directory: str = "backend/data/vectorstore", embedding_model_name: str = "llama2", embedding_device: str = "cpu", use_ollama:bool = True, ollama_base_url: str = "http://localhost:11434"):
        self.persist_directory = persist_directory
        self.client = chromadb.PersistentClient(path=persist_directory)
        
        if use_ollama:
            self.embedding_fn = OllamaEmbeddingFunction(embedding_model_name, ollama_base_url, embedding_device)
        else:
            self.embedding_fn = SpacySentenceTransformerEmbeddingFunction(DEFAULT_MODEL_NAME, embedding_device) #Fallback

        self.collections = {
            "section2": self.client.get_or_create_collection(
                name="section2_questions",
                embedding_function=self.embedding_fn,
                metadata={"description": "JLPT listening comprehension questions - Section 2"}
            ),
            "section3": self.client.get_or_create_collection(
                name="section3_questions",
                embedding_function=self.embedding_fn,
                metadata={"description": "JLPT phrase matching questions - Section 3"}
            )
        }

    def add_questions(self, section_num: int, questions: List[Dict], video_id: str):
        """Add questions to the vector store"""
        if section_num not in [2, 3]:
            raise ValueError("Only sections 2 and 3 are currently supported")

        collection = self.collections[f"section{section_num}"]

        ids = []
        documents = []
        metadatas = []

        for idx, question in enumerate(questions):
            # Create a unique ID for each question
            question_id = f"{video_id}_{section_num}_{idx}"
            ids.append(question_id)

            # Store the full question structure as metadata
            metadatas.append({
                "video_id": video_id,
                "section": section_num,
                "question_index": idx,
                "full_structure": json.dumps(question)
            })

            # Create a searchable document from the question content
            if section_num == 2:
                document = f"""
                Introduction: {question['Introduction']}
                Dialogue: {question['Conversation']}
                Question: {question['Question']}
                """
            else:  # section 3
                document = f"""
                Situation: {question['Situation']}
                Question: {question['Question']}
                """
            documents.append(document)

        # Add to collection
        collection.add(
            ids=ids,
            documents=documents,
            metadatas=metadatas
        )

    def search_similar_questions(
        self,
        section_num: int,
        query: str,
        n_results: int = 5
    ) -> List[Dict]:
        """Search for similar questions in the vector store"""
        if section_num not in [2, 3]:
            raise ValueError("Only sections 2 and 3 are currently supported")

        collection = self.collections[f"section{section_num}"]

        results = collection.query(
            query_texts=[query],
            n_results=n_results
        )

        # Convert results to more usable format
        questions = []
        for idx, metadata in enumerate(results['metadatas'][0]):
            question_data = json.loads(metadata['full_structure'])
            question_data['similarity_score'] = results['distances'][0][idx]
            questions.append(question_data)
            #Add section_num key to be used in main.py
            question_data['section_num'] = section_num

        return questions

    def get_question_by_id(self, section_num: int, question_id: str) -> Optional[Dict]:
        """Retrieve a specific question by its ID"""
        if section_num not in [2, 3]:
            raise ValueError("Only sections 2 and 3 are currently supported")

        collection = self.collections[f"section{section_num}"]

        result = collection.get(
            ids=[question_id],
            include=['metadatas']
        )

        if result['metadatas']:
            question_data = json.loads(result['metadatas'][0]['full_structure'])
            #Add section_num key to be used in main.py
            question_data['section_num'] = section_num
            return question_data
        return None

    def parse_questions_from_file(self, filename: str) -> List[Dict]:
        """Parse questions from a structured text file"""
        questions = []
        current_question = {}
        section_num = int(os.path.basename(filename).split('_section')[1].split('.')[0]) # Extract section number from filename

        try:
            with open(filename, 'r', encoding='utf-8') as f:
                lines = f.readlines()

            i = 0
            while i < len(lines):
                line = lines[i].strip()

                if line.startswith('<question>'):
                    current_question = {}
                elif line.startswith('Introduction:'):
                    i += 1
                    if i < len(lines):
                        current_question['Introduction'] = lines[i].strip()
                elif line.startswith('Conversation:'):
                    i += 1
                    if i < len(lines):
                        current_question['Conversation'] = lines[i].strip()
                elif line.startswith('Situation:'):
                    i += 1
                    if i < len(lines):
                        current_question['Situation'] = lines[i].strip()
                elif line.startswith('Question:'):
                    i += 1
                    if i < len(lines):
                        current_question['Question'] = lines[i].strip()
                elif line.startswith('Options:'):
                    options = []
                    for _ in range(4):
                        i += 1
                        if i < len(lines):
                            option = lines[i].strip()
                            if option.startswith('1.') or option.startswith('2.') or option.startswith('3.') or option.startswith('4.'):
                                options.append(option[2:].strip())
                    current_question['Options'] = options
                elif line.startswith('</question>'):
                    if current_question:
                        questions.append(current_question)
                        current_question = {}
                i += 1
            return questions
        except Exception as e:
            print(f"Error parsing questions from {filename}: {str(e)}")
            return []

    def index_questions_file(self, filename: str, section_num: int):
        """Index all questions from a file into the vector store"""
        # Extract video ID from filename
        video_id = os.path.basename(filename).split('_section')[0]

        # Parse questions from file
        questions = self.parse_questions_from_file(filename)

        # Add to vector store
        if questions:
            self.add_questions(section_num, questions, video_id)
            print(f"Indexed {len(questions)} questions from {filename}")

if __name__ == "__main__":
    # Example usage
    store = QuestionVectorStore()

    # Index questions from files
    question_files = [
        ("backend/data/questions/sY7L5cfCWno_section2.txt", 2),
        ("backend/data/questions/sY7L5cfCWno_section3.txt", 3)
    ]

    for filename, section_num in question_files:
        if os.path.exists(filename):
            store.index_questions_file(filename, section_num)

    # Search for similar questions
    similar = store.search_similar_questions(2, "誕生日について質問", n_results=1)
    print("Similar questions",similar)
    # Example of get question by id
    question_id = "sY7L5cfCWno_2_0"
    question = store.get_question_by_id(2, question_id)
    print("question by id", question)

