"""Typesetting service for rendering translated text"""

from typing import List
import numpy as np
import cv2
from PIL import Image, ImageDraw, ImageFont

from config.settings import (
    FONT_PATH,
    FONT_SIZE_START,
    FONT_SIZE_MIN,
    FONT_SIZE_STEP,
    BOX_PADDING,
    BOX_BORDER_RADIUS,
    TEXT_PADDING_X_PCT,
    TEXT_PADDING_X_MIN_PX,
    TEXT_PADDING_Y_PCT,
    TEXT_PADDING_Y_MIN_PX,
    LINE_SPACING,
    TEXT_STROKE_WIDTH,
    TEXT_STROKE_FILL,
    INPAINT_RADIUS,
    INPAINT_DILATE_ITERATIONS,
    INPAINT_DILATE_KERNEL_SIZE,
    INPAINT_TEXT_THRESHOLD,
    INPAINT_USE_OTSU,
    INPAINT_ALGORITHM,
    INPAINT_BACKEND,
    LAMA_MODEL_NAME,
    LAMA_DEVICE,
)


class Typesetter:
    """Handles text rendering and box cleaning for manga pages."""

    def __init__(self, font_path: str = FONT_PATH):
        self.font_path = font_path
        self._lama_model = None
        self._InpaintRequest = None

        if INPAINT_BACKEND == "lama":
            self._load_lama_model()

    def _load_lama_model(self):
        """Load the IOPaint manga inpainting model once at startup, downloading if needed.

        Loads the raw Manga model directly (not via ModelManager) so we can call
        forward() on individual small crops, avoiding full-image HD strategy passes
        that would exhaust VRAM when running alongside Qwen and other models.
        """
        try:
            import torch
            from iopaint.model.manga import Manga
            from iopaint.model import models as iopaint_models
            from iopaint.schema import InpaintRequest, ModelInfo, ModelType

            # Download weights on first use (no-op if already cached)
            if not iopaint_models[LAMA_MODEL_NAME].is_downloaded():
                print(f"Downloading IOPaint '{LAMA_MODEL_NAME}' model weights (first run only)...")
                iopaint_models[LAMA_MODEL_NAME].download()

            print(f"Loading IOPaint '{LAMA_MODEL_NAME}' inpainting model...")
            device = torch.device(LAMA_DEVICE)

            # Instantiate the Manga model directly, bypassing ModelManager and
            # HD strategy — we call forward() on small per-box crops instead
            model_info = ModelInfo(name=LAMA_MODEL_NAME, path=LAMA_MODEL_NAME, model_type=ModelType.INPAINT)
            self._lama_model = Manga(device=device, model_info=model_info)
            self._InpaintRequest = InpaintRequest
            print(f"IOPaint '{LAMA_MODEL_NAME}' model loaded.")
        except ImportError:
            print("WARNING: iopaint not installed. Falling back to cv2 inpainting.")
            self._lama_model = None
        except Exception as e:
            print(f"WARNING: Could not load IOPaint model ({e}). Falling back to cv2.")
            self._lama_model = None

    def get_font(self, size: int) -> ImageFont.FreeTypeFont:
        try:
            return ImageFont.truetype(self.font_path, size)
        except Exception:
            return ImageFont.load_default()

    def _make_text_mask(self, gray_roi: np.ndarray) -> np.ndarray:
        """Generate a binary mask marking dark text pixels in a grayscale ROI."""
        if INPAINT_USE_OTSU:
            otsu_thr, text_mask = cv2.threshold(
                gray_roi, 0, 255, cv2.THRESH_BINARY_INV | cv2.THRESH_OTSU
            )
            # Safety guard: Otsu on a mostly-white blank bubble can choose a
            # near-maximum threshold — fall back to fixed value in that case.
            if otsu_thr > 230:
                _, text_mask = cv2.threshold(
                    gray_roi, INPAINT_TEXT_THRESHOLD, 255, cv2.THRESH_BINARY_INV
                )
        else:
            _, text_mask = cv2.threshold(
                gray_roi, INPAINT_TEXT_THRESHOLD, 255, cv2.THRESH_BINARY_INV
            )
        return text_mask

    def _dilate_mask(self, mask: np.ndarray) -> np.ndarray:
        """Dilate a binary mask to cover text anti-aliasing and compression artifacts."""
        kernel = np.ones(
            (INPAINT_DILATE_KERNEL_SIZE, INPAINT_DILATE_KERNEL_SIZE), np.uint8
        )
        return cv2.dilate(mask, kernel, iterations=INPAINT_DILATE_ITERATIONS)

    # ------------------------------------------------------------------
    # Public entry points
    # ------------------------------------------------------------------

    def clean_box(self, image: Image.Image, box: List[int]) -> Image.Image:
        """
        Remove original text from a bounding box region.

        Dispatches to the LaMa neural backend (if loaded) or the cv2 fallback.
        For bulk processing of multiple boxes on one image use clean_all_boxes()
        instead — it runs a single LaMa pass which is significantly faster.
        """
        if INPAINT_BACKEND == "lama" and self._lama_model is not None:
            return self._clean_box_lama(image, box)
        return self._clean_box_cv2(image, box)

    def clean_all_boxes(self, image: Image.Image, boxes: List[List[int]]) -> Image.Image:
        """
        Remove text from all boxes.

        With the manga backend, runs forward() directly on each box crop to keep
        VRAM usage minimal (avoids full-image HD strategy passes that cause OOM
        when Qwen and other models are loaded simultaneously).
        """
        if INPAINT_BACKEND != "lama" or self._lama_model is None:
            for box in boxes:
                self._clean_box_cv2(image, box)
            return image

        import torch
        for box in boxes:
            self._clean_box_lama(image, box)
        torch.cuda.empty_cache()
        return image

    # ------------------------------------------------------------------
    # Private implementations
    # ------------------------------------------------------------------

    def _clean_box_cv2(self, image: Image.Image, box: List[int]) -> Image.Image:
        """
        Clean a text box using OpenCV masked inpainting.

        1. Extract ROI from bounding box
        2. Create binary mask of dark text pixels (Otsu or fixed threshold)
        3. Dilate mask to cover anti-aliasing edges
        4. Inpaint only the masked pixels (Telea or Navier-Stokes algorithm)
        """
        x1, y1, x2, y2 = box
        w = x2 - x1
        h = y2 - y1

        if w <= 0 or h <= 0:
            return image

        img_array = np.array(image)
        img_bgr = cv2.cvtColor(img_array, cv2.COLOR_RGB2BGR)

        roi = img_bgr[y1:y1+h, x1:x1+w]
        if roi.size == 0:
            return image

        gray_roi = cv2.cvtColor(roi, cv2.COLOR_BGR2GRAY)
        text_mask = self._make_text_mask(gray_roi)
        dilated_mask = self._dilate_mask(text_mask)

        algo = cv2.INPAINT_NS if INPAINT_ALGORITHM == "ns" else cv2.INPAINT_TELEA
        cleaned_roi = cv2.inpaint(roi, dilated_mask, INPAINT_RADIUS, algo)

        img_bgr[y1:y1+h, x1:x1+w] = cleaned_roi
        img_rgb = cv2.cvtColor(img_bgr, cv2.COLOR_BGR2RGB)
        result = Image.fromarray(img_rgb)
        image.paste(result, (0, 0))
        return image

    def _clean_box_lama(self, image: Image.Image, box: List[int]) -> Image.Image:
        """
        Clean a single text box using the manga inpainting model.

        Calls forward() directly on the box crop (RGB [H,W,C] input, [H,W,1] mask).
        Pastes only the inpainted crop back at (x1,y1) — no full-page array ops.
        """
        import torch

        x1, y1, x2, y2 = box
        w, h = x2 - x1, y2 - y1
        if w <= 0 or h <= 0:
            return image

        # Crop directly from PIL — avoids converting the full page to numpy
        roi_pil = image.crop((x1, y1, x2, y2))
        roi_rgb = np.ascontiguousarray(np.array(roi_pil))
        if roi_rgb.size == 0:
            return image

        # Pad to multiple of pad_mod=16
        pad_mod = 16
        ph = ((h + pad_mod - 1) // pad_mod) * pad_mod
        pw = ((w + pad_mod - 1) // pad_mod) * pad_mod
        padded = np.zeros((ph, pw, 3), dtype=np.uint8)
        padded[:h, :w] = roi_rgb

        # Build mask [H, W, 1] for this crop
        gray = cv2.cvtColor(roi_rgb, cv2.COLOR_RGB2GRAY)
        text_mask = self._make_text_mask(gray)
        text_mask = self._dilate_mask(text_mask)
        padded_mask = np.zeros((ph, pw), dtype=np.uint8)
        padded_mask[:h, :w] = text_mask
        padded_mask_3d = padded_mask[:, :, np.newaxis]  # [H, W, 1] as forward() expects

        # Run model forward on the small crop
        config = self._InpaintRequest()
        with torch.no_grad():
            result_bgr = self._lama_model.forward(padded, padded_mask_3d, config)

        # Crop back to original box size and convert BGR -> RGB
        # iopaint forward() may return float32 [0,255]; cast to uint8 before cv2
        result_arr = np.clip(result_bgr, 0, 255).astype(np.uint8) if result_bgr.dtype != np.uint8 else result_bgr
        result_crop = np.ascontiguousarray(result_arr[:h, :w])
        result_rgb = cv2.cvtColor(result_crop, cv2.COLOR_BGR2RGB)

        # Blend: keep original pixels where mask is 0, use inpainted where mask is 255
        mask_bool = text_mask > 0
        roi_out = roi_rgb.copy()
        roi_out[mask_bool] = result_rgb[mask_bool]

        # Paste only the small crop back into the PIL image at the correct position
        image.paste(Image.fromarray(roi_out), (x1, y1))
        return image

    # ------------------------------------------------------------------
    # Text rendering
    # ------------------------------------------------------------------

    def pixel_wrap(self, text: str, font: ImageFont.FreeTypeFont, max_width: int) -> List[str]:
        """Wrap text to fit within a given pixel width."""
        text = text.replace('\n', ' ').replace('\r', '')

        words = text.split(' ')
        lines = []
        current_line = []

        for word in words:
            if not word:
                continue

            test_line = ' '.join(current_line + [word])

            try:
                w = font.getlength(test_line)
            except AttributeError:
                # Pillow < 9.2 fallback
                w = font.getsize(test_line)[0]
            except Exception:
                w = 999999

            if w <= max_width:
                current_line.append(word)
            else:
                if current_line:
                    lines.append(' '.join(current_line))
                    current_line = [word]
                else:
                    lines.append(word)
                    current_line = []

        if current_line:
            lines.append(' '.join(current_line))
        return lines

    def draw_text(self, image: Image.Image, text: str, box: List[int]) -> Image.Image:
        """
        Draw text into a bounding box with auto-sizing, wrapping, and stroke.

        - Padding uses max(pixel_floor, percentage) to handle both tiny and
          large bubbles correctly.
        - Font size steps down by FONT_SIZE_STEP (1pt) for better fit precision.
        - Text is rendered with a white stroke outline for readability on
          complex backgrounds (screentones, gradients).
        """
        x1, y1, x2, y2 = box
        w_box = x2 - x1
        h_box = y2 - y1
        draw = ImageDraw.Draw(image)

        padx = max(TEXT_PADDING_X_MIN_PX, int(w_box * TEXT_PADDING_X_PCT))
        pady = max(TEXT_PADDING_Y_MIN_PX, int(h_box * TEXT_PADDING_Y_PCT))

        w_usable = w_box - (2 * padx)
        h_usable = h_box - (2 * pady)
        start_x_inner = x1 + padx
        start_y_inner = y1 + pady

        fontsize = FONT_SIZE_START
        final_lines = []
        final_font = None
        final_line_height = 0

        while fontsize >= FONT_SIZE_MIN:
            font = self.get_font(fontsize)
            lines = self.pixel_wrap(text, font, w_usable)
            ascent, descent = font.getmetrics()
            base_height = ascent + descent
            total_height = base_height * len(lines) * LINE_SPACING

            if total_height <= h_usable:
                final_lines = lines
                final_font = font
                final_line_height = base_height * LINE_SPACING
                break
            fontsize -= FONT_SIZE_STEP

        if final_font is None:
            final_font = self.get_font(FONT_SIZE_MIN)
            final_lines = self.pixel_wrap(text, final_font, w_usable)
            ascent, descent = final_font.getmetrics()
            final_line_height = (ascent + descent) * LINE_SPACING

        total_block_height = final_line_height * len(final_lines)
        current_y = start_y_inner + (h_usable - total_block_height) / 2

        for line in final_lines:
            try:
                line_w = final_font.getlength(line)
            except AttributeError:
                line_w = final_font.getsize(line)[0]
            line_x = start_x_inner + (w_usable - line_w) / 2
            draw.text(
                (line_x, current_y),
                line,
                font=final_font,
                fill="black",
                anchor="la",
                stroke_width=TEXT_STROKE_WIDTH,
                stroke_fill=TEXT_STROKE_FILL,
            )
            current_y += final_line_height

        return image
