"""Text processing utilities"""

import re


def sanitize_for_font(text: str) -> str:
    """
    Clean forbidden characters to avoid squares or crashes in font rendering.

    Args:
        text: The text to sanitize

    Returns:
        Sanitized text safe for font rendering
    """
    # 1. Smart replacements
    replacements = {
        '…': '...', ''': "'", ''': "'", '"': '"', '"': '"',
        '–': '-', '—': '-', 'œ': 'oe', 'Œ': 'OE'
    }
    for old, new in replacements.items():
        text = text.replace(old, new)

    # 2. Remove line breaks (CRITICAL for Pillow)
    text = text.replace('\n', ' ').replace('\r', ' ')

    # 3. Whitelist (Letters, Numbers, Basic punctuation only)
    # Remove everything that's not standard (emojis, remaining kanjis...)
    allowed_chars_pattern = r"[^a-zA-Z0-9àâäéèêëîïôöùûüçÀÂÄÉÈÊËÎÏÔÖÙÛÜÇ.,?!:;'\-\"()\[\]\s]"
    clean_text = re.sub(allowed_chars_pattern, '', text)

    # Reduce multiple spaces (e.g., "Hello   World" -> "Hello World")
    clean_text = re.sub(r'\s+', ' ', clean_text)

    return clean_text.strip()


def clean_translation_output(raw: str) -> str:
    """
    Clean LLM translation output from thinking tags and prefixes.

    Args:
        raw: Raw LLM output

    Returns:
        Cleaned translation text
    """
    # 1. Clean <think> tags
    clean = re.sub(r'<think>.*?</think>', '', raw, flags=re.DOTALL)
    clean = clean.strip()

    # 2. Clean prefixes
    prefixes = r'(?i)^(text|translation|english|en|output|response)\s*[:\-]?\s*'
    clean = re.sub(prefixes, '', clean)

    # 3. Clean quotes
    clean = clean.strip().strip('"').strip("'")

    # 4. Final sanitization (Font + Crash fix)
    clean = sanitize_for_font(clean)

    return clean
