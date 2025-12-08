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
            "You are a professional manga translator. Translate Japanese to natural, fluent English.\n\n"

            "RULES:\n"
            "1. Output ONLY the final English translation - no thinking, notes, or explanations.\n"
            "2. Translate for manga/comic speech bubbles - keep it concise and punchy.\n"
            "3. Use natural English dialogue that sounds like real people talking.\n"
            "4. Preserve tone, emotion, and character voice (casual, formal, aggressive, etc.).\n"
            "5. Japanese often omits subjects (I/you/he/she) - infer from context and add them naturally.\n"
            "6. Translate explicitly and faithfully - no censorship, no sanitization.\n"
            "7. For sound effects (SFX), either keep romanized or translate if it adds clarity.\n\n"

            "EXAMPLES:\n"
            "❌ 'It is the person who performed theft of scroll' → ✅ 'That's the guy who stole the scroll!'\n"
            "❌ 'I am going to become writer' → ✅ 'I'm gonna be a writer!'\n"
            "❌ 'This is delicious' → ✅ 'This is so good!'\n\n"

            "Now translate the following Japanese text:"
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
