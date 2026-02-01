"""Main pipeline for manga translation processing"""

import os
import sys
import shutil
import zipfile
import torch
from typing import Optional
from PIL import Image
from manga_ocr import MangaOcr
from ultralytics import YOLO

from config.settings import (
    YOLO_MODEL_NAME,
    YOLO_CONFIDENCE_THRESHOLD,
    OUTPUT_QUALITY,
    TEMP_DIR,
    MODEL_PATH,
    FONT_PATH
)
from services.translation import LocalTranslator
from services.typesetting import Typesetter
from utils.box_processing import consolidate_boxes


class MangaPipeline:
    """Main pipeline for processing manga images and ZIP files."""

    def __init__(self):
        """Initialize the manga translation pipeline."""
        self.device = 'cuda' if torch.cuda.is_available() else 'cpu'
        print(f"⚙️ Device: {self.device.upper()}")

        if not os.path.exists(YOLO_MODEL_NAME):
            print(f"❌ Missing YOLO model: {YOLO_MODEL_NAME}")
            sys.exit(1)

        self.detector = YOLO(YOLO_MODEL_NAME)
        self.mocr = MangaOcr()
        self.translator = LocalTranslator(MODEL_PATH)
        self.typesetter = Typesetter(FONT_PATH)
        print("✅ Pipeline Ready (V10 - Stable | Masked Inpainting).", flush=True)

    def process_image(self, image_path: str, output_path: Optional[str] = None) -> Optional[str]:
        """
        Process a single manga image.

        Args:
            image_path: Path to input image
            output_path: Optional path for output image

        Returns:
            Path to saved output image, or None if processing failed
        """
        try:
            with Image.open(image_path) as img_src:
                img_src.load()
                original_img = img_src.convert("RGB")
        except Exception as e:
            print(f"⚠️ Skipped invalid image: {image_path} ({e})")
            return None

        print(f"   Processing: {os.path.basename(image_path)}")

        # Detect text boxes
        results = self.detector(
            original_img, conf=YOLO_CONFIDENCE_THRESHOLD, verbose=False)
        boxes = []
        for r in results:
            if r.boxes:
                for box in r.boxes.xyxy.cpu().numpy():
                    boxes.append(list(map(int, box)))

        boxes = consolidate_boxes(boxes)

        # Process each text box
        if boxes:
            for box in boxes:
                x1, y1, x2, y2 = box
                crop = original_img.crop((x1, y1, x2, y2))

                try:
                    jap_text = self.mocr(crop)
                except Exception:
                    continue

                if not jap_text.strip():
                    continue

                fr_text = self.translator.translate(jap_text)
                self.typesetter.clean_box(original_img, box)
                self.typesetter.draw_text(original_img, fr_text, box)

        # Save output
        if output_path:
            base, _ = os.path.splitext(output_path)
            save_path = base + ".jpg"
        else:
            base = os.path.splitext(os.path.basename(image_path))[0]
            save_path = f"translated_{base}.jpg"

        original_img.save(save_path, "JPEG", quality=OUTPUT_QUALITY)
        return save_path

    def process_zip(self, zip_path: str) -> None:
        """
        Process a ZIP file containing manga images.

        Args:
            zip_path: Path to input ZIP file
        """
        print(f"\n📦 ZIP Detected: {zip_path}")
        temp_dir = TEMP_DIR

        if os.path.exists(temp_dir):
            try:
                shutil.rmtree(temp_dir)
            except Exception:
                pass
        os.makedirs(temp_dir)

        # Extract ZIP
        with zipfile.ZipFile(zip_path, 'r') as zip_ref:
            zip_ref.extractall(temp_dir)

        # Collect all image files first
        image_files = []
        for root, _, files in os.walk(temp_dir):
            for file in files:
                if file.lower().endswith(('.png', '.jpg', '.jpeg', '.webp', '.bmp')):
                    image_files.append(os.path.join(root, file))

        total_images = len(image_files)
        print(f"📄 Found {total_images} images to process", flush=True)

        # Process all images with progress tracking
        for idx, full_input_path in enumerate(image_files, 1):
            print(
                f"   Processing: {os.path.basename(full_input_path)} ({idx}/{total_images})", flush=True)
            new_jpg_path = self.process_image(
                full_input_path, output_path=full_input_path)

            if new_jpg_path and os.path.normpath(new_jpg_path) != os.path.normpath(full_input_path):
                try:
                    os.remove(full_input_path)
                except Exception:
                    pass

            # Report progress after each image
            progress = int((idx / total_images) * 90) + 5  # 5-95% range
            print(
                f"PROGRESS: {progress}% - Translated {idx}/{total_images} pages", flush=True)

        print(f" -> Processed {total_images} images.", flush=True)
        # Create output ZIP in the same directory as input
        zip_dir = os.path.dirname(zip_path)
        zip_basename = os.path.splitext(os.path.basename(zip_path))[0]
        output_zip_base = os.path.join(zip_dir, f"{zip_basename}_translated")
        shutil.make_archive(output_zip_base, 'zip', temp_dir)
        output_zip = output_zip_base

        try:
            shutil.rmtree(temp_dir)
        except Exception:
            pass

        print(f"✅ Created: {output_zip}.zip", flush=True)

    def run(self, input_path: str) -> None:
        """
        Run the pipeline on an input file (image or ZIP).

        Args:
            input_path: Path to input file
        """
        if input_path.lower().endswith('.zip'):
            self.process_zip(input_path)
        else:
            self.process_image(input_path)
