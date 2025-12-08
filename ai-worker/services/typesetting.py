"""Typesetting service for rendering translated text"""

from typing import List, Tuple
import numpy as np
import cv2
from PIL import Image, ImageDraw, ImageFont

from config.settings import (
    FONT_PATH,
    FONT_SIZE_START,
    FONT_SIZE_MIN,
    BOX_PADDING,
    BOX_BORDER_RADIUS,
    TEXT_PADDING_X_PCT,
    TEXT_PADDING_Y_PCT,
    LINE_SPACING,
    INPAINT_RADIUS,
    INPAINT_DILATE_ITERATIONS,
    INPAINT_DILATE_KERNEL_SIZE,
    INPAINT_TEXT_THRESHOLD
)


class Typesetter:
    """Handles text rendering and box cleaning for manga pages."""

    def __init__(self, font_path: str = FONT_PATH):
        """
        Initialize the typesetter.

        Args:
            font_path: Path to the TrueType font file
        """
        self.font_path = font_path
        self.dummy_draw = ImageDraw.Draw(Image.new("RGB", (1, 1)))

    def get_font(self, size: int) -> ImageFont.FreeTypeFont:
        """
        Load a font at the specified size.

        Args:
            size: Font size in points

        Returns:
            Font object
        """
        try:
            return ImageFont.truetype(self.font_path, size)
        except Exception:
            return ImageFont.load_default()

    def clean_box(self, image: Image.Image, box: List[int]) -> Image.Image:
        """
        Intelligently clean a text box using masked inpainting.

        This method targets ONLY the dark text pixels for removal:
        1. Extracts the ROI (region of interest) from the bounding box
        2. Creates a binary mask using fixed threshold to detect dark pixels (text)
        3. Dilates the mask to cover text edges and compression artifacts
        4. Uses cv2.inpaint to fill only the masked text pixels

        This preserves artwork and backgrounds, even when boxes overlap.

        Args:
            image: PIL Image to modify (will be modified in-place)
            box: Bounding box [x1, y1, x2, y2]

        Returns:
            Modified image
        """
        x1, y1, x2, y2 = box

        # Convert bounding box to width/height format
        w = x2 - x1
        h = y2 - y1

        if w <= 0 or h <= 0:
            return image  # Skip invalid boxes

        # Convert PIL image to OpenCV format (RGB -> BGR)
        img_array = np.array(image)
        img_bgr = cv2.cvtColor(img_array, cv2.COLOR_RGB2BGR)

        # 1. Extract ROI
        roi = img_bgr[y1:y1+h, x1:x1+w]

        if roi.size == 0:
            return image  # Skip if ROI is empty

        # 2. Create Text Mask using fixed threshold
        # Convert ROI to grayscale
        gray_roi = cv2.cvtColor(roi, cv2.COLOR_BGR2GRAY)

        # Apply binary threshold: pixels darker than threshold become white (255) in mask
        # THRESH_BINARY_INV inverts: dark pixels (text) -> white, light pixels -> black
        _, text_mask = cv2.threshold(
            gray_roi,
            INPAINT_TEXT_THRESHOLD,
            255,
            cv2.THRESH_BINARY_INV
        )

        # 3. Dilate Mask to ensure complete coverage of text anti-aliasing
        kernel = np.ones((INPAINT_DILATE_KERNEL_SIZE, INPAINT_DILATE_KERNEL_SIZE), np.uint8)
        dilated_mask = cv2.dilate(text_mask, kernel, iterations=INPAINT_DILATE_ITERATIONS)

        # 4. Inpaint only the white pixels in the dilated mask
        # Telea algorithm is fast and produces good results for text removal
        cleaned_roi = cv2.inpaint(roi, dilated_mask, INPAINT_RADIUS, cv2.INPAINT_TELEA)

        # 5. Paste the cleaned ROI back into the original image
        img_bgr[y1:y1+h, x1:x1+w] = cleaned_roi

        # Convert back to PIL Image (BGR -> RGB) and update the original
        img_rgb = cv2.cvtColor(img_bgr, cv2.COLOR_BGR2RGB)
        result = Image.fromarray(img_rgb)

        # Update the original image in-place
        image.paste(result, (0, 0))

        return image

    def pixel_wrap(self, text: str, font: ImageFont.FreeTypeFont, max_width: int) -> List[str]:
        """
        Wrap text to fit within a given pixel width.

        Args:
            text: Text to wrap
            font: Font to use for measurement
            max_width: Maximum width in pixels

        Returns:
            List of wrapped text lines
        """
        # CRASH FIX: Ensure there are no \n characters
        text = text.replace('\n', ' ').replace('\r', '')

        words = text.split(' ')
        lines = []
        current_line = []

        for word in words:
            if not word:
                continue  # Skip empty words

            test_line = ' '.join(current_line + [word])

            # This is where it crashed if test_line had a \n
            try:
                w = self.dummy_draw.textlength(test_line, font=font)
            except ValueError:
                # If it still crashes somehow, force calculation without the word
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
        Draw text into a bounding box with auto-sizing and wrapping.

        Args:
            image: PIL Image to modify
            text: Text to render
            box: Bounding box [x1, y1, x2, y2]

        Returns:
            Modified image
        """
        x1, y1, x2, y2 = box
        w_box = x2 - x1
        h_box = y2 - y1
        draw = ImageDraw.Draw(image)

        padx = int(w_box * TEXT_PADDING_X_PCT)
        pady = int(h_box * TEXT_PADDING_Y_PCT)

        w_usable = w_box - (2 * padx)
        h_usable = h_box - (2 * pady)
        start_x_inner = x1 + padx
        start_y_inner = y1 + pady

        fontsize = FONT_SIZE_START
        final_lines = []
        final_font = None
        final_line_height = 0

        # Try to fit text with decreasing font sizes
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
            fontsize -= 2

        # Fallback to minimum font size
        if final_font is None:
            final_font = self.get_font(FONT_SIZE_MIN)
            final_lines = self.pixel_wrap(text, final_font, w_usable)
            ascent, descent = final_font.getmetrics()
            final_line_height = (ascent + descent) * LINE_SPACING

        # Center text vertically
        total_block_height = final_line_height * len(final_lines)
        current_y = start_y_inner + (h_usable - total_block_height) / 2

        # Draw each line centered horizontally
        for line in final_lines:
            line_w = draw.textlength(line, font=final_font)
            line_x = start_x_inner + (w_usable - line_w) / 2
            draw.text(
                (line_x, current_y),
                line,
                font=final_font,
                fill="black",
                anchor="la"
            )
            current_y += final_line_height

        return image
