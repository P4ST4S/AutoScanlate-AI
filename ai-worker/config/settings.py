"""Configuration settings for manga translator"""

# Model Configuration
MODEL_PATH = "./models/Qwen2.5-7B-Instruct-abliterated-v2.Q4_K_M.gguf"
GPU_LAYERS = -1
CONTEXT_WINDOW = 4096

# YOLO Configuration
YOLO_MODEL_NAME = "./models/manga-text-detector.pt"
YOLO_CONFIDENCE_THRESHOLD = 0.20

# Font Configuration
FONT_PATH = "./fonts/animeace2_reg.ttf"
FONT_SIZE_START = 20
FONT_SIZE_MIN = 2

# Translation Configuration
TRANSLATION_TEMPERATURE = 0.1
TRANSLATION_MAX_TOKENS = 200

# Typesetting Configuration
BOX_PADDING = 6
BOX_BORDER_RADIUS = 65
TEXT_PADDING_X_PCT = 0.15
TEXT_PADDING_Y_PCT = 0.02
LINE_SPACING = 0.9

# Inpainting Configuration
INPAINT_RADIUS = 3  # Radius for cv2.inpaint algorithm
INPAINT_DILATE_ITERATIONS = 1  # Dilation iterations for text mask
INPAINT_DILATE_KERNEL_SIZE = 3  # Kernel size for dilation (3x3)
INPAINT_TEXT_THRESHOLD = 180  # Threshold for detecting dark text (0-255, lower = more aggressive)

# Box Consolidation Configuration
BOX_DISTANCE_THRESHOLD = 25

# Output Configuration
OUTPUT_QUALITY = 95
TEMP_DIR = "temp_process"
