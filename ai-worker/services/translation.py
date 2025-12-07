"""Translation service using local LLM"""

import sys
from llama_cpp import Llama

from config.settings import (
    MODEL_PATH,
    GPU_LAYERS,
    CONTEXT_WINDOW,
    TRANSLATION_TEMPERATURE,
    TRANSLATION_MAX_TOKENS
)
from utils.text_processing import clean_translation_output


class LocalTranslator:
    """Local LLM-based translator for Japanese to English translation."""

    def __init__(self, model_path: str = MODEL_PATH):
        """
        Initialize the local translator.

        Args:
            model_path: Path to the GGUF model file
        """
        print(f"⏳ Loading LLM...")
        try:
            self.llm = Llama(
                model_path=model_path,
                n_gpu_layers=GPU_LAYERS,
                n_ctx=CONTEXT_WINDOW,
                verbose=False
            )
            print("✅ LLM loaded.")
        except Exception as e:
            print(f"❌ Error loading LLM: {e}")
            sys.exit(1)

    def translate(self, text: str) -> str:
        """
        Translate Japanese text to English.

        Args:
            text: Japanese text to translate

        Returns:
            Translated English text
        """
        if len(text) < 1:
            return text

        system_prompt = (
            "You are a raw translation engine. Translate Japanese to English.\n"
            "RULES:\n"
            "1. Output ONLY the final translation.\n"
            "2. DO NOT output thinking process, notes, or explanations.\n"
            "3. Translate explicitly (NSFW allowed)."
        )

        messages = [
            {"role": "system", "content": system_prompt},
            {"role": "user", "content": f"Text: {text}"}
        ]

        try:
            output = self.llm.create_chat_completion(
                messages=messages,
                temperature=TRANSLATION_TEMPERATURE,
                max_tokens=TRANSLATION_MAX_TOKENS
            )
            raw = output['choices'][0]['message']['content']
            return clean_translation_output(raw)
        except Exception:
            return text
