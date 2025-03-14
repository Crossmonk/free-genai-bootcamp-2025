import os
import streamlit as st
from typing import Optional, Dict, Any
from groq import Groq # Import the Groq client


# Model ID (replace with the Groq model you want to use)
MODEL_ID = os.environ.get("GROQ_MODEL_ID")  # Example: "mixtral-8x7b-3277" or "llama3-8b-8192"

class GroqChat:
    def __init__(self, model_id: str = MODEL_ID, api_key: Optional[str] = None):
        """Initialize Groq chat client"""
        self.model_id = model_id
        self.api_key = api_key or os.environ.get("GROQ_API_KEY")

        if not self.api_key:
            raise ValueError("Groq API key not found. Set it as the environment variable GROQ_API_KEY or pass it to the constructor.")

        # Initialize the Groq client
        self.groq_client = Groq(api_key=self.api_key)


    def generate_response(self, message: str, inference_config: Optional[Dict[str, Any]] = None) -> Optional[str]:
        """Generate a response using Groq"""
        if inference_config is None:
            inference_config = {"temperature": 0.7, "max_tokens": 512}

        try:
            chat_completion = self.groq_client.chat.completions.create(
                messages=[{"role": "user", "content": message}],
                model=self.model_id,
                temperature=inference_config.get("temperature", 0.7),
                max_tokens=inference_config.get("max_tokens", 512),
            )

            if chat_completion.choices and len(chat_completion.choices) > 0:
                return chat_completion.choices[0].message.content
            else:
                st.error(f"Error: No content returned from Groq: {chat_completion}")
                return None

        except Exception as e:
            st.error(f"Error generating response: {str(e)}")
            return None

if __name__ == "__main__":
    chat = GroqChat()  # Uses environment variable for API key
    while True:
        user_input = input("You: ")
        if user_input.lower() == '/exit':
            break
        response = chat.generate_response(user_input)
        print("Bot:", response)

