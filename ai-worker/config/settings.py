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
FONT_SIZE_MIN = 14

# Translation Configuration
TRANSLATION_TEMPERATURE = 0.1
TRANSLATION_MAX_TOKENS = 200

# Typesetting Configuration
BOX_PADDING = 6
BOX_BORDER_RADIUS = 65
TEXT_PADDING_X_PCT = 0.15
TEXT_PADDING_Y_PCT = 0.02
LINE_SPACING = 0.9

# Box Consolidation Configuration
BOX_DISTANCE_THRESHOLD = 25

# Output Configuration
OUTPUT_QUALITY = 95
TEMP_DIR = "temp_process"
