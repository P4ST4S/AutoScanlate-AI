"""Manga Translator - Main Entry Point"""

import os
import argparse

from core.pipeline import MangaPipeline


def main():
    """Main entry point for the manga translator."""
    parser = argparse.ArgumentParser(
        description="Translate manga images from Japanese to English"
    )
    parser.add_argument(
        "input",
        help="Path to image or ZIP file to process"
    )
    args = parser.parse_args()

    if not os.path.exists(args.input):
        print("❌ File not found.")
        return

    pipeline = MangaPipeline()
    pipeline.run(args.input)


if __name__ == "__main__":
    main()
