# 🏯 AI Manga Translator (Local GPU Version)

A high-performance, fully local pipeline for automatically translating manga/comics from Japanese to English (or other languages).

This project uses a combination of state-of-the-art AI models to **Detect**, **Read**, **Translate**, and **Typeset** manga pages locally on your GPU. No external APIs, no costs, and complete privacy.

## 🚀 Key Features (V10)

- **⚡ 100% Local & GPU Accelerated**: Powered by llama.cpp and CUDA. Optimized for NVIDIA RTX cards (runs entirely in VRAM).
- **👁️ Smart Detection**: Uses YOLOv8 (fine-tuned on Manga109) to detect speech bubbles.
  - Smart Box Merging algorithm to consolidate fragmented vertical text bubbles into single coherent blocks.
- **📖 Robust OCR**: Utilizes MangaOCR to accurately read vertical and handwritten Japanese text.
- **🧠 Uncensored Translation**: Integrated with Qwen 2.5 7B (Abliterated) for high-quality, unfiltered translations (supports NSFW, slang, and honorifics).
  - Custom "Anti-Thinking" prompt engineering to prevent LLM hallucinations or internal monologues appearing in the final text.
- **🎨 Advanced Typesetting**:
  - **NEW (V10)**: **Intelligent Masked Inpainting** - Uses OpenCV threshold detection and cv2.inpaint to remove ONLY dark text pixels, preserving artwork and backgrounds even when bounding boxes overlap.
  - **Pixel-Perfect Wrapping**: Custom algorithm that measures text width in pixels (not characters) to prevent words from being cut off or overlapping borders.
  - **Sanitization**: Automatically filters unsupported characters (emojis, complex symbols) to prevent font glitches.
- **📦 Batch Processing**:
  - Supports single images (.jpg, .png, .webp).
  - **Native ZIP Support**: Automatically extracts chapters, translates all images, and re-packages them into a `_translated.zip`.
  - **Format Normalization**: Automatically converts all outputs to high-quality JPG.
- **🏗️ Modular Architecture**: Clean, maintainable codebase with separation of concerns for easy customization and extension.

## 🛠️ Prerequisites

- **OS**: Windows / Linux
- **Python**: 3.10 or higher
- **GPU**: NVIDIA RTX series recommended (Min 6GB VRAM for 7B models, 8GB+ preferred).
- **System Tools**: NVIDIA CUDA Toolkit 12.x installed.

## 📥 Installation

### 1. Clone & Setup Environment

```bash
# Create a virtual environment
python -m venv venv
.\venv\Scripts\activate

# Upgrade pip
pip install --upgrade pip
```

### 2. Install PyTorch (CUDA Version)

You must install the version compatible with your CUDA drivers.

```bash
pip install torch torchvision torchaudio --index-url https://download.pytorch.org/whl/cu121
```

### 3. Install Dependencies

```bash
pip install -r requirements.txt
```

Dependencies include: `manga-ocr`, `ultralytics`, `Pillow`, `numpy`.

### 4. Compile Llama-cpp-python (GPU)

This is the most critical step to enable hardware acceleration for the translator.

```bash
# Set flags to force CUDA build
$env:CMAKE_ARGS="-DGGML_CUDA=on"

# Install/Reinstall forcing compilation
pip install llama-cpp-python --no-cache-dir --force-reinstall
```

## 🤖 Models Setup

Create a `models/` directory in the root folder and download the following:

### Translation Model (LLM):

- **Model**: `Qwen2.5-7B-Instruct-abliterated-Q4_K_M.gguf`
- **Why**: Best balance of speed/quality, fits in 12GB VRAM, and does not refuse NSFW/Contextual translations.
- **Source**: [HuggingFace Link](https://huggingface.co/QuantFactory/Qwen2.5-7B-Instruct-abliterated-v2-GGUF)

### Detection Model (YOLO):

- **Model**: `manga-text-detector.pt`
- **Source**: [HuggingFace Link](https://huggingface.co/ogkalu/manga-text-detector-yolov8s/blob/main/manga-text-detector.pt)

### Fonts:

Place a `.ttf` font file in the root directory (e.g., `arial.ttf` or `animeace2_reg.ttf`). The font `animeace2_reg.ttf` is already included in the `fonts/` folder, and used by default.

## 📸 Examples

See the intelligent masked inpainting in action! These examples showcase V10's ability to preserve artwork while cleanly removing text.

### Example 1: Naruto

<table>
<tr>
<td width="50%">
<img src="exemples/exemple_naruto.png" alt="Original Naruto page" />
<p align="center"><b>Original (Japanese)</b></p>
</td>
<td width="50%">
<img src="exemples/translated_exemple_naruto.jpg" alt="Translated Naruto page" />
<p align="center"><b>Translated (English)</b></p>
</td>
</tr>
</table>

### Example 2: One Piece

<table>
<tr>
<td width="50%">
<img src="exemples/exemple_one_piece.png" alt="Original One Piece page" />
<p align="center"><b>Original (Japanese)</b></p>
</td>
<td width="50%">
<img src="exemples/translated_exemple_one_piece.jpg" alt="Translated One Piece page" />
<p align="center"><b>Translated (English)</b></p>
</td>
</tr>
</table>

**V10 Improvements Demonstrated:**
- Clean text removal without damaging background artwork
- Preserved bubble borders and shading
- Accurate text positioning and sizing
- No artifacts in overlapping bubble regions

---

## 💻 Usage

### Translate a Single Page

```bash
python main.py .\input\image_01.jpg
```

**Output**: `translated_image_01.jpg`

### Translate a Full Chapter (ZIP)

```bash
python main.py .\input\OnePiece_Chapter_1050.zip
```

- Extracts the zip.
- Translates every image inside.
- Deletes original files to save space.

**Output**: `OnePiece_Chapter_1050_translated.zip`

## ⚙️ Configuration

All configuration settings are centralized in `config/settings.py`. Key settings you can adjust:

```python
# Model Configuration
MODEL_PATH = "./models/7b/Qwen2.5-7B-Instruct-abliterated-v2.Q4_K_M.gguf"
GPU_LAYERS = -1  # -1 = offload all layers to GPU
CONTEXT_WINDOW = 4096  # Lower to 2048 to save VRAM on smaller cards

# YOLO Configuration
YOLO_MODEL_NAME = "./models/manga-text-detector.pt"
YOLO_CONFIDENCE_THRESHOLD = 0.20

# Font Configuration
FONT_PATH = "./fonts/animeace2_reg.ttf"
FONT_SIZE_START = 20  # Starting font size (auto-reduces if needed)
FONT_SIZE_MIN = 14

# Translation Configuration
TRANSLATION_TEMPERATURE = 0.1
TRANSLATION_MAX_TOKENS = 200

# Typesetting Configuration
BOX_PADDING = 6
LINE_SPACING = 0.9
TEXT_PADDING_X_PCT = 0.15
TEXT_PADDING_Y_PCT = 0.02

# Inpainting Configuration (V10+)
INPAINT_RADIUS = 3  # Radius for cv2.inpaint algorithm
INPAINT_DILATE_ITERATIONS = 1  # Dilation iterations for text mask
INPAINT_DILATE_KERNEL_SIZE = 3  # Kernel size for dilation (3x3)
INPAINT_TEXT_THRESHOLD = 180  # Threshold for detecting dark text (0-255, lower = more aggressive)

# Output Configuration
OUTPUT_QUALITY = 95  # JPEG quality
```

## 🧩 Architecture

The codebase follows a modular architecture with clear separation of concerns:

```
ai-worker/
├── main.py                    # Entry point & CLI argument parsing
├── config/
│   └── settings.py           # Centralized configuration
├── core/
│   └── pipeline.py           # Main processing pipeline orchestration
├── services/
│   ├── translation.py        # LLM translation service
│   └── typesetting.py        # Text rendering & box cleaning
└── utils/
    ├── text_processing.py    # Text sanitization utilities
    └── box_processing.py     # Box consolidation algorithms
```

### Processing Pipeline

1. **Input**: Image load + conversion to RGB.
2. **YOLO Detection**: Scans the page for text bubbles (via `core/pipeline.py`).
3. **Post-Processing**: Merges overlapping/nearby boxes to handle split vertical text (`utils/box_processing.py`).
4. **MangaOCR**: Crops the merged boxes and extracts Japanese text.
5. **LLM Translator**: Sends text to Qwen 2.5 (running on GPU via llama.cpp) (`services/translation.py`).
   - **Prompting**: "Raw translation engine" system prompt + Few-shot examples + Regex cleaning to remove `<think>` tags and prefixes.
6. **Typesetter** (`services/typesetting.py`):
   - **V10**: Uses intelligent masked inpainting to remove only text pixels while preserving backgrounds and artwork.
   - Calculates optimal font size using `pixel_wrap` (dynamic width measurement).
   - Centers text vertically and horizontally.
7. **Output**: Saves as High-Quality JPEG.

## 📝 Credits

- **MangaOCR**: [kha-white](https://github.com/kha-white)
- **YOLOv8**: [Ultralytics](https://github.com/ultralytics/ultralytics)
- **Llama.cpp**: [ggerganov](https://github.com/ggerganov/llama.cpp)
- **Qwen**: Alibaba Cloud

---

## 📋 Changelog

### V10 (Stable) - 2025-12-08
**Major Enhancement: Intelligent Masked Inpainting**

- **Feature**: Completely refactored text cleaning pipeline to use OpenCV-based masked inpainting
  - Replaced simple rectangle erasure with intelligent text-only detection
  - Uses `cv2.threshold` with configurable threshold (default: 180) to identify dark text pixels
  - Applies `cv2.dilate` to expand mask and cover text anti-aliasing
  - Uses `cv2.inpaint` (TELEA algorithm) to fill only detected text regions
- **Benefits**:
  - Preserves artwork and backgrounds even when bounding boxes overlap
  - Eliminates artifacts from merged/overlapping detection boxes
  - No damage to surrounding art in complex bubble arrangements
- **New Configuration Options** (`config/settings.py`):
  - `INPAINT_RADIUS` - Controls inpainting algorithm radius (default: 3)
  - `INPAINT_DILATE_ITERATIONS` - Dilation passes for mask expansion (default: 1)
  - `INPAINT_DILATE_KERNEL_SIZE` - Kernel size for dilation (default: 3x3)
  - `INPAINT_TEXT_THRESHOLD` - Threshold for dark text detection 0-255 (default: 180)
- **Files Modified**:
  - `services/typesetting.py` - Complete rewrite of `clean_box()` method
  - `config/settings.py` - Added inpainting configuration section

### V9 (Stable)
- Smart Box Merging algorithm for vertical text consolidation
- Anti-Thinking prompt engineering to prevent LLM hallucinations
- Enhanced text sanitization and pixel-perfect wrapping
- Improved batch processing with ZIP support

---

**Current Version**: V10 (Stable)
