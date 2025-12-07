"""Services module for manga translator"""

from .translation import LocalTranslator
from .typesetting import Typesetter

__all__ = ["LocalTranslator", "Typesetter"]
