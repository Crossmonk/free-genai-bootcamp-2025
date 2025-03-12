import json
import os  # Import the os module
from typing import Dict, List, Optional
from backend.vector_store import QuestionVectorStore
import requests
from dotenv import load_dotenv # Import load_dotenv

# Load environment variables from .env file
load_dotenv()

class QuestionGenerator:
    def __init__(self):
        """Initialize Groq client and vector store"""
        self.vector_store = QuestionVectorStore()
        self.model_id = os.environ.get("GROQ_MODEL_ID")  # Replace with the Groq model ID you want to use
        self.groq_api_url = "https://api.groq.com/openai/chat/completions" #correct
        # Read the API key from the environment variable
        self.groq_api_key = os.getenv("GROQ_API_KEY")
        if not self.groq_api_key:
            raise ValueError("GROQ_API_KEY environment variable not set.")

    def _invoke_groq(self, prompt: str) -> Optional[str]:
        """Invoke Groq with the given prompt"""
        try:
            headers = {
                "Authorization": f"Bearer {self.groq_api_key}",
                "Content-Type": "application/json"
            }

            messages = [{"role": "user", "content": prompt}]

            data = {
                "model": self.model_id,
                "messages": messages,
                "temperature": 0.7,
                "max_tokens": 2048  # adjust if needed
            }

            response = requests.post(self.groq_api_url, headers=headers, json=data)
            response.raise_for_status()  # Raise an exception for bad status codes

            response_json = response.json()

            if 'choices' in response_json and len(response_json['choices']) > 0 and 'message' in response_json['choices'][0]:
                return response_json['choices'][0]['message']['content']
            else:
                print("Unexpected response format from Groq:", response_json)
                return None

        except requests.exceptions.RequestException as e:
            print(f"Error invoking Groq: {str(e)}")
            return None
        except json.JSONDecodeError as e:
            print(f"Error decoding JSON response from Groq: {str(e)}")
            return None
        except Exception as e:
            print(f"An unexpected error occurred: {e}")
            return None

    def generate_similar_question(self, section_num: int, topic: str) -> Dict:
        """Generate a new question similar to existing ones on a given topic"""
        # Get similar questions for context
        similar_questions = self.vector_store.search_similar_questions(section_num, topic, n_results=3)

        if not similar_questions:
            return None

        # Create context from similar questions
        context = "Here are some example JLPT listening questions:\n\n"
        for idx, q in enumerate(similar_questions, 1):
            if section_num == 2:
                context += f"Example {idx}:\n"
                context += f"Introduction: {q.get('Introduction', '')}\n"
                context += f"Conversation: {q.get('Conversation', '')}\n"
                context += f"Question: {q.get('Question', '')}\n"
                if 'Options' in q:
                    context += "Options:\n"
                    for i, opt in enumerate(q['Options'], 1):
                        context += f"{i}. {opt}\n"
            else:  # section 3
                context += f"Example {idx}:\n"
                context += f"Situation: {q.get('Situation', '')}\n"
                context += f"Question: {q.get('Question', '')}\n"
                if 'Options' in q:
                    context += "Options:\n"
                    for i, opt in enumerate(q['Options'], 1):
                        context += f"{i}. {opt}\n"
            context += "\n"

        # Create prompt for generating new question
        prompt = f"""Based on the following example JLPT listening questions, create a new question about {topic}.
        The question should follow the same format but be different from the examples.
        Make sure the question tests listening comprehension and has a clear correct answer.

        {context}

        Generate a new question following the exact same format as above. Include all components (Introduction/Situation, 
        Conversation/Question, and Options). Make sure the question is challenging but fair, and the options are plausible 
        but with only one clearly correct answer. Return ONLY the question without any additional text.

        New Question:
        """

        # Generate new question
        response = self._invoke_groq(prompt)
        if not response:
            return None

        # Parse the generated question
        try:
            lines = response.strip().split('\n')
            question = {}
            current_key = None
            current_value = []

            for line in lines:
                line = line.strip()
                if not line:
                    continue

                if line.startswith("Introduction:"):
                    if current_key:
                        question[current_key] = ' '.join(current_value)
                    current_key = 'Introduction'
                    current_value = [line.replace("Introduction:", "").strip()]
                elif line.startswith("Conversation:"):
                    if current_key:
                        question[current_key] = ' '.join(current_value)
                    current_key = 'Conversation'
                    current_value = [line.replace("Conversation:", "").strip()]
                elif line.startswith("Situation:"):
                    if current_key:
                        question[current_key] = ' '.join(current_value)
                    current_key = 'Situation'
                    current_value = [line.replace("Situation:", "").strip()]
                elif line.startswith("Question:"):
                    if current_key:
                        question[current_key] = ' '.join(current_value)
                    current_key = 'Question'
                    current_value = [line.replace("Question:", "").strip()]
                elif line.startswith("Options:"):
                    if current_key:
                        question[current_key] = ' '.join(current_value)
                    current_key = 'Options'
                    current_value = []
                elif line[0].isdigit() and line[1] == "." and current_key == 'Options':
                    current_value.append(line[2:].strip())
                elif current_key:
                    current_value.append(line)

            if current_key:
                if current_key == 'Options':
                    question[current_key] = current_value
                else:
                    question[current_key] = ' '.join(current_value)

            # Ensure we have exactly 4 options
            if 'Options' not in question or len(question.get('Options', [])) != 4:
                # Use default options if we don't have exactly 4
                question['Options'] = [
                    "ピザを食べる",
                    "ハンバーガーを食べる",
                    "サラダを食べる",
                    "パスタを食べる"
                ]

            return question
        except Exception as e:
            print(f"Error parsing generated question: {str(e)}")
            return None

    def get_feedback(self, question: Dict, selected_answer: int) -> Dict:
        """Generate feedback for the selected answer"""
        if not question or 'Options' not in question:
            return None

        # Create prompt for generating feedback
        prompt = f"""Given this JLPT listening question and the selected answer, provide feedback explaining if it's correct 
        and why. Keep the explanation clear and concise.

        """
        if 'Introduction' in question:
            prompt += f"Introduction: {question['Introduction']}\n"
            prompt += f"Conversation: {question['Conversation']}\n"
        else:
            prompt += f"Situation: {question['Situation']}\n"

        prompt += f"Question: {question['Question']}\n"
        prompt += "Options:\n"
        for i, opt in enumerate(question['Options'], 1):
            prompt += f"{i}. {opt}\n"

        prompt += f"\nSelected Answer: {selected_answer}\n"
        prompt += "\nProvide feedback in JSON format with these fields:\n"
        prompt += "- correct: true/false\n"
        prompt += "- explanation: brief explanation of why the answer is correct/incorrect\n"
        prompt += "- correct_answer: the number of the correct option (1-4)\n"

        # Get feedback
        response = self._invoke_groq(prompt)
        if not response:
            return None

        try:
            # Parse the JSON response
            feedback = json.loads(response.strip())
            return feedback
        except:
            # If JSON parsing fails, return a basic response with a default correct answer
            return {
                "correct": False,
                "explanation": "Unable to generate detailed feedback. Please try again.",
                "correct_answer": 1  # Default to first option
            }
