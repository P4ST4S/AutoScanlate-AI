"""Configuration settings for manga translator"""

# Model Configuration
MODEL_PATH = "./models/Qwen2.5-7B-Instruct-abliterated-v2.Q4_K_M.gguf"
GPU_LAYERS = 0
CONTEXT_WINDOW = 4096

# YOLO Configuration
YOLO_MODEL_NAME = "./models/manga-text-detector.pt"
YOLO_CONFIDENCE_THRESHOLD = 0.20

# Font Configuration
FONT_PATH = "./fonts/animeace2_reg.ttf"
FONT_SIZE_START = 24          # Increased from 20 — more headroom before shrinking
FONT_SIZE_MIN = 2
FONT_SIZE_STEP = 1            # Step size when reducing font (was hardcoded 2)

# Translation Configuration
TRANSLATION_TEMPERATURE = 0.1
TRANSLATION_MAX_TOKENS = 200

# Typesetting Configuration
BOX_PADDING = 6
BOX_BORDER_RADIUS = 65
TEXT_PADDING_X_PCT = 0.10     # Reduced from 0.15 (pixel floor below prevents over-padding)
TEXT_PADDING_X_MIN_PX = 4     # Hard pixel floor for X padding on small bubbles
TEXT_PADDING_Y_PCT = 0.02
TEXT_PADDING_Y_MIN_PX = 2     # Hard pixel floor for Y padding
LINE_SPACING = 0.9
TEXT_STROKE_WIDTH = 1         # White outline around black text (0 to disable)
TEXT_STROKE_FILL = "white"    # Outline colour

# Inpainting Configuration
INPAINT_RADIUS = 7            # Increased from 3 — better fill for larger text strokes
INPAINT_DILATE_ITERATIONS = 2 # Increased from 1 — covers anti-aliasing & JPEG artifacts
INPAINT_DILATE_KERNEL_SIZE = 5 # Increased from 3 — 5x5 covers more edge pixels
INPAINT_TEXT_THRESHOLD = 180  # Fallback threshold when Otsu is disabled
INPAINT_USE_OTSU = True       # Adaptive threshold per ROI (True) vs fixed threshold (False)
INPAINT_ALGORITHM = "telea"   # "telea" (Telea, fast) or "ns" (Navier-Stokes, better large areas)

# Inpainting Backend (Phase 2)
INPAINT_BACKEND = "lama"       # "cv2" = OpenCV inpainting, "lama" = IOPaint manga model
LAMA_MODEL_NAME = "manga"      # IOPaint model name — "manga" is the manga-specific inpainter
LAMA_DEVICE = "cpu"            # cpu avoids CUDA conflicts with llama.cpp (Qwen)

# Box Consolidation Configuration
BOX_DISTANCE_THRESHOLD = 25

# Output Configuration
OUTPUT_QUALITY = 95
TEMP_DIR = "temp_process"
