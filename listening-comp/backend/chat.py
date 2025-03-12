# groq_chat.py
import requests
import streamlit as st
from typing import Optional, Dict, Any
import os

# Model ID (replace with the Groq model you want to use)
MODEL_ID = os.environ.get("GROQ_MODEL_ID")  # Example: "mixtral-8x7b-3277" or "llama3-8b-8192"

class GroqChat:
    def __init__(self, model_id: str = MODEL_ID, api_key: Optional[str] = None):
        """Initialize Groq chat client"""
        self.model_id = model_id
        self.api_key = api_key or os.environ.get("GROQ_API_KEY")

        if not self.api_key:
            raise ValueError("Groq API key not found. Set it as the environment variable GROQ_API_KEY or pass it to the constructor.")

        self.base_url = "https://api.groq.com/openai" #correct
        self.headers = {
            "Authorization": f"Bearer {self.api_key}",
            "Content-Type": "application/json",
        }

    def generate_response(self, message: str, inference_config: Optional[Dict[str, Any]] = None) -> Optional[str]:
        """Generate a response using Groq"""
        if inference_config is None:
            inference_config = {"temperature": 0.7, "max_tokens": 512}  # Add max_tokens

        messages = [{"role": "user", "content": message}]

        payload = {
            "model": self.model_id,
            "messages": messages,
            "temperature": inference_config.get("temperature", 0.7),
            "max_tokens" : inference_config.get("max_tokens", 512)
        }

        try:
            response = requests.post(
                url=f"{self.base_url}/chat/completions", #correct
                headers=self.headers,
                json=payload,
            )
            response.raise_for_status()  # Raise an exception for bad status codes
            response_json = response.json()

            if response_json['choices']:
                return response_json['choices'][0]['message']['content']
            else:
                st.error(f"Error: No content returned from Groq")
                return None

        except requests.exceptions.RequestException as e:
            st.error(f"Error generating response: {str(e)}")
            return None
        except (KeyError, IndexError) as e:
            st.error(f"Error parsing response: {str(e)}")
            return None


if __name__ == "__main__":
    chat = GroqChat()  # Uses environment variable for API key
    while True:
        user_input = input("You: ")
        if user_input.lower() == '/exit':
            break
        response = chat.generate_response(user_input)
        print("Bot:", response)
