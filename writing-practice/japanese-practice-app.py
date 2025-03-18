import streamlit as st
import requests
import json
import base64
from PIL import Image
import io
import manga_ocr  # MangaOCR for Japanese text recognition
import openai  # For LLM interactions
from dotenv import load_dotenv
import os

# Load environment variables from .env file
load_dotenv()

# Get the API key from the environment variables
openai.api_key = os.getenv("OPENAI_API_KEY")

# Check if API key is loaded correctly
if not openai.api_key:
    st.error("OPENAI_API_KEY not found in environment variables. Please check your .env file.")
    st.stop()

# Initialize session state variables if they don't exist
if 'state' not in st.session_state:
    st.session_state.state = "setup"
if 'current_sentence' not in st.session_state:
    st.session_state.current_sentence = ""
if 'english_sentence' not in st.session_state:
    st.session_state.english_sentence = ""
if 'vocabulary' not in st.session_state:
    st.session_state.vocabulary = []
if 'uploaded_image' not in st.session_state:
    st.session_state.uploaded_image = None
if 'review_data' not in st.session_state:
    st.session_state.review_data = None

# Initialize MangaOCR
@st.cache_resource
def load_ocr():
    return manga_ocr.MangaOcr()

mocr = load_ocr()

# Functions
def fetch_vocabulary(group_id):
    """Fetch vocabulary from the API"""
    try:
        response = requests.get(f"http://localhost:5000/api/words")
        if response.status_code == 200:
            return response.json()
        else:
            st.error(f"Failed to fetch vocabulary: {response.status_code}")
            return []
    except Exception as e:
        st.error(f"Error fetching vocabulary: {e}")
        return []

def generate_sentence(word):
    """Generate a simple Japanese practice sentence using the provided word"""
    prompt = f"""Generate a simple sentence using the following word: {word}
    The grammar should be scoped to JLPT N5 grammar.
    You can use the following vocabulary to construct a simple sentence:
    - simple objects e.g. book, car, ramen, sushi
    - simple verbs, to drink, to eat, to meet
    - simple times e.g. tomorrow, today, yesterday
    
    Return only the English sentence.
    """
    
    try:
        response = openai.ChatCompletion.create(
            model="gpt-4",
            messages=[{"role": "system", "content": "You are a Japanese language tutor."},
                      {"role": "user", "content": prompt}]
        )
        return response.choices[0].message["content"].strip()
    except Exception as e:
        st.error(f"Error generating sentence: {e}")
        return "I eat sushi."  # Fallback sentence

def transcribe_image(image):
    """Transcribe Japanese text from the uploaded image using MangaOCR"""
    try:
        text = mocr(image)
        return text
    except Exception as e:
        st.error(f"Error transcribing image: {e}")
        return ""

def translate_text(japanese_text):
    """Translate Japanese text to English"""
    prompt = f"Translate the following Japanese text to English literally: {japanese_text}"
    
    try:
        response = openai.ChatCompletion.create(
            model="gpt-4",
            messages=[{"role": "system", "content": "You are a Japanese translator."},
                      {"role": "user", "content": prompt}]
        )
        return response.choices[0].message["content"].strip()
    except Exception as e:
        st.error(f"Error translating text: {e}")
        return ""

def grade_attempt(original_english, japanese_attempt, english_translation):
    """Grade the user's Japanese writing attempt"""
    prompt = f"""Grade this Japanese language practice attempt:
    
    Original English sentence: "{original_english}"
    User's Japanese writing (transcribed): "{japanese_attempt}"
    Literal translation of user's writing: "{english_translation}"
    
    Provide:
    1. A letter grade using S, A, B, C, D, F ranking
    2. A brief explanation of whether the attempt accurately conveyed the English sentence
    3. 1-2 specific suggestions for improvement
    
    Format your response as a JSON object with keys: "grade", "explanation", "suggestions"
    """
    
    try:
        response = openai.ChatCompletion.create(
            model="gpt-4",
            messages=[{"role": "system", "content": "You are a Japanese language teacher."},
                      {"role": "user", "content": prompt}]
        )
        
        # Parse the JSON response
        feedback = json.loads(response.choices[0].message["content"])
        return feedback
    except Exception as e:
        st.error(f"Error grading attempt: {e}")
        return {
            "grade": "C",
            "explanation": "Unable to properly assess due to system error.",
            "suggestions": ["Please try again."]
        }

def generate_new_question():
    """Generate a new practice sentence and update the state"""
    if st.session_state.vocabulary:
        # Randomly select a word from vocabulary
        import random
        word = random.choice(st.session_state.vocabulary)
        
        # Generate English sentence
        english_sentence = generate_sentence(word["japanese"])
        
        # Update session state
        st.session_state.english_sentence = english_sentence
        st.session_state.state = "practice"
        st.session_state.uploaded_image = None
        st.session_state.review_data = None
    else:
        st.error("No vocabulary available. Please check the API connection.")

def handle_submit_for_review():
    """Process the uploaded image and generate review feedback"""
    if st.session_state.uploaded_image is not None:
        # Transcribe image
        japanese_text = transcribe_image(st.session_state.uploaded_image)
        
        # Translate transcription
        english_translation = translate_text(japanese_text)
        
        # Grade the attempt
        feedback = grade_attempt(
            st.session_state.english_sentence,
            japanese_text,
            english_translation
        )
        
        # Update session state
        st.session_state.review_data = {
            "transcription": japanese_text,
            "translation": english_translation,
            "feedback": feedback
        }
        st.session_state.state = "review"
    else:
        st.error("Please upload an image before submitting for review.")

# Main app logic
def main():
    st.title("Japanese Writing Practice App")
    
    # On first load, fetch vocabulary
    if len(st.session_state.vocabulary) == 0:
        group_id = st.text_input("Enter vocabulary group ID:", "1")
        if st.button("Load Vocabulary"):
            st.session_state.vocabulary = fetch_vocabulary(group_id)
            st.success(f"Loaded {len(st.session_state.vocabulary)} vocabulary items!")
    
    # Render current state
    if st.session_state.state == "setup":
        if len(st.session_state.vocabulary) > 0:
            if st.button("Generate Sentence"):
                generate_new_question()
    
    elif st.session_state.state == "practice":
        st.subheader("Translate this sentence to Japanese:")
        st.write(st.session_state.english_sentence)
        
        uploaded_file = st.file_uploader("Upload your handwritten Japanese", type=["jpg", "jpeg", "png"])
        if uploaded_file is not None:
            image = Image.open(uploaded_file)
            st.session_state.uploaded_image = image
            st.image(image, caption="Uploaded Image", use_column_width=True)
        
        if st.button("Submit for Review"):
            handle_submit_for_review()
    
    elif st.session_state.state == "review":
        st.subheader("Original Sentence:")
        st.write(st.session_state.english_sentence)
        
        review_data = st.session_state.review_data
        if review_data:
            st.subheader("Your Japanese Writing:")
            st.write(review_data["transcription"])
            
            st.subheader("Translation:")
            st.write(review_data["translation"])
            
            st.subheader("Feedback:")
            st.markdown(f"**Grade:** {review_data['feedback']['grade']}")
            st.markdown(f"**Assessment:** {review_data['feedback']['explanation']}")
            
            st.subheader("Suggestions:")
            for suggestion in review_data['feedback']['suggestions']:
                st.write(f"- {suggestion}")
        
        if st.button("Next Question"):
            generate_new_question()

if __name__ == "__main__":
    main()
