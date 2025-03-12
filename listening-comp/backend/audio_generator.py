import requests
import json
import os
from typing import Dict, List, Tuple
import tempfile
import subprocess
from datetime import datetime

class AudioGenerator:
    def __init__(self):
        # VOICEVOX server URL - default runs on localhost:50021
        self.voicevox_url = os.getenv('VOICEVOX_URL', 'http://localhost:50021')
        
        # Define Japanese voices by gender
        # VOICEVOX speaker IDs (these are examples, check actual available voices)
        self.voices = {
            'male': [1, 6, 2],     # Example male voice IDs: 1:ずんだもん(ノーマル), 6:四国めたん(あまあま), 2:冥鳴ひまり(ノーマル)
            'female': [3, 4, 7],    # Example female voice IDs: 3:春日部つむぎ, 4:雨晴はう(ノーマル), 7:九州そら(あまあま)
            'announcer': 1          # Default announcer voice ID: ずんだもん(ノーマル)
        }
        
        # Create audio output directory
        self.audio_dir = os.path.join(
            os.path.dirname(os.path.dirname(os.path.abspath(__file__))),
            "frontend/static/audio"
        )
        os.makedirs(self.audio_dir, exist_ok=True)

    def get_voice_for_gender(self, gender: str) -> int:
        """Get appropriate VOICEVOX speaker ID for the given gender"""
        if gender == 'male':
          #return self.voices['male'][0]  # Use first male voice
            return self.voices['male'][0]
        elif gender == 'female':
            #return self.voices['female'][0]  # Use first female voice
            return self.voices['female'][0]
        else:
            return self.voices['announcer']

    def generate_audio_part(self, text: str, gender:str) -> str:
        print("""Generate audio for a single part using VOICEVOX""")
        speaker_id = self.get_voice_for_gender(gender)
        try:
            # Generate audio query from text
            query_response = requests.post(
                f'{self.voicevox_url}/audio_query',
                params={
                    'text': text,
                    'speaker': speaker_id
                }
            )
            query_response.raise_for_status()
            query_data = query_response.json()

            print("Synthesize speech")
            # Synthesize speech from query
            synthesis_response = requests.post(
                f'{self.voicevox_url}/synthesis',
                params={
                    'speaker': speaker_id
                },
                json=query_data
            )
            synthesis_response.raise_for_status()

            print("Save to temporary file")
            # Save to temporary file
            with tempfile.NamedTemporaryFile(suffix='.wav', delete=False) as temp_file:
                temp_file.write(synthesis_response.content)
                wav_file = temp_file.name

            print("Convert to MP3")
            # Convert to MP3 (since VOICEVOX outputs WAV)
            mp3_file = wav_file.replace('.wav', '.mp3')
            subprocess.run([
                'ffmpeg', '-i', wav_file,
                '-codec:a', 'libmp3lame', '-qscale:a', '2',
                mp3_file
            ], check=True)

            # Clean up WAV file
            os.unlink(wav_file)
            return mp3_file

        except Exception as e:
            print(f"Error generating audio: {str(e)}")
            raise e

    def combine_audio_files(self, audio_files: List[str], output_file: str):
        """Combine multiple audio files using ffmpeg"""
        file_list = None
        try:
            # Create file list for ffmpeg
            with tempfile.NamedTemporaryFile('w', suffix='.txt', delete=False) as f:
                for audio_file in audio_files:
                    f.write(f"file '{audio_file}'\n")
                file_list = f.name

            # Combine audio files
            subprocess.run([
                'ffmpeg', '-f', 'concat', '-safe', '0',
                '-i', file_list,
                '-c', 'copy',
                output_file
            ], check=True)

            return True
        except Exception as e:
            print(f"Error combining audio files: {str(e)}")
            if os.path.exists(output_file):
                os.unlink(output_file)
            return False
        finally:
            # Clean up temporary files
            if file_list and os.path.exists(file_list):
                os.unlink(file_list)
            for audio_file in audio_files:
                if os.path.exists(audio_file):
                    try:
                        os.unlink(audio_file)
                    except Exception as e:
                        print(f"Error cleaning up {audio_file}: {str(e)}")
                        
    def generate_silence(self, duration_ms: int) -> str:
        """Generate a silent audio file of specified duration"""
        output_file = os.path.join(self.audio_dir, f'silence_{duration_ms}ms.mp3')
        if not os.path.exists(output_file):
            subprocess.run([
                'ffmpeg', '-f', 'lavfi', '-i',
                f'anullsrc=r=24000:cl=mono:d={duration_ms/1000}',
                '-c:a', 'libmp3lame', '-b:a', '48k',
                output_file
            ])
        return output_file
                        
    def generate_audio(self, conversation_parts: List[Tuple[str, str, str]]) -> str:
      """
      Generates the complete audio file for a conversation.
      Args:
          conversation_parts: A list of tuples, where each tuple is:
                              (speaker_name, text, gender)
      Returns:
          The path to the generated audio file.
      """
      timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
      output_file = os.path.join(self.audio_dir, f"conversation_{timestamp}.mp3")

      audio_files = []
      try:
          long_pause = self.generate_silence(2000)
          short_pause = self.generate_silence(500)
          current_section = None

          for speaker, text, gender in conversation_parts:
            # Detect section changes and add appropriate pauses
            if speaker.lower() == 'announcer':
                if '次の会話' in text:  # Introduction
                    if current_section is not None:
                        audio_files.append(long_pause)
                    current_section = 'intro'
                elif '質問' in text or '選択肢' in text:  # Question or options
                    audio_files.append(long_pause)
                    current_section = 'question'
            elif current_section == 'intro':
                audio_files.append(long_pause)
                current_section = 'conversation'
            
            # Generate audio for this part
            audio_file = self.generate_audio_part(text, gender)
            audio_files.append(audio_file)
            audio_files.append(short_pause)

          if not self.combine_audio_files(audio_files, output_file):
            raise Exception("Failed to combine audio files")
          
      except Exception as e:
          print(f"Error in generate_audio: {str(e)}")
          raise e
      finally:
        #Clean up the files.
          for file in audio_files:
            try:
              if os.path.exists(file):
                os.unlink(file)
            except Exception as e:
              print(f"Error deleting temporary audio files {e}")
      return output_file
