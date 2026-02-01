# <img src="https://raw.githubusercontent.com/lucide-icons/lucide/main/icons/castle.svg" width="24" height="24" /> Manga AI Translator

An automated, privacy-focused, GPU-accelerated pipeline to translate manga and comics locally.

This project aims to provide a full-stack solution (Frontend, Backend, and AI Worker) to detect text bubbles, perform OCR, translate contextually using LLMs, and typeset the result back into the original image—all without external APIs or recurring costs.

## <img src="https://raw.githubusercontent.com/lucide-icons/lucide/main/icons/video.svg" width="24" height="24" /> Demo

![Demo](docs/demo.gif)

## <img src="https://raw.githubusercontent.com/lucide-icons/lucide/main/icons/layers.svg" width="24" height="24" /> Architecture

The project follows a Microservices architecture to ensure the heavy AI processing doesn't block the web server.

![Architecture Diagram](docs/manga-translation-architecture.drawio.png)

## <img src="https://raw.githubusercontent.com/lucide-icons/lucide/main/icons/puzzle.svg" width="24" height="24" /> Project Structure

| Module         | Status                                                                                                                          | Description                                                                                                      |
| -------------- | ------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------- |
| `/ai-worker`   | <img src="https://raw.githubusercontent.com/lucide-icons/lucide/main/icons/circle-check.svg" width="24" height="24" /> v10.0    | The core Python engine. Handles Computer Vision, OCR, and LLM Inference on GPU.                                  |
| `/backend-api` | <img src="https://raw.githubusercontent.com/lucide-icons/lucide/main/icons/circle-check.svg" width="24" height="24" /> **v2.0** | High-performance Go API with **real-time SSE progress**, Redis pub/sub, ZIP extraction, and nested file support. |
| `/frontend`    | <img src="https://raw.githubusercontent.com/lucide-icons/lucide/main/icons/circle-check.svg" width="24" height="24" /> v1.0     | Modern Web UI (Next.js 16) for drag-and-drop uploads and reading translated chapters.                            |

## <img src="https://raw.githubusercontent.com/lucide-icons/lucide/main/icons/rocket.svg" width="24" height="24" /> What's New in Backend v2.0

The backend API has been significantly enhanced with production-ready features:

### 🔴 Live Progress Streaming

- **Real-time SSE updates** showing page-by-page translation progress
- **Instant feedback** with proper Python stdout unbuffering
- **Reliable broadcasting** via Redis pub/sub architecture
- **Connection stability** with proper resource cleanup and error handling

### 📦 Enhanced ZIP Support

- **Automatic extraction** of original and translated archives
- **Subdirectory preservation** - maintains complex folder structures
- **Instant page counting** - displays total pages immediately on upload
- **Smart path handling** - supports nested directories and Unicode filenames

### 🏗️ Architecture Improvements

- **Proper SSE lifecycle** with deferred cleanup in goroutines
- **Wildcard routing** for flexible file serving
- **Enhanced logging** with detailed progress tracking
- **Type-safe callbacks** throughout the translation pipeline

## <img src="https://raw.githubusercontent.com/lucide-icons/lucide/main/icons/sparkles.svg" width="24" height="24" /> Key Features (AI Worker V10)

The core engine is currently fully operational.

**<img src="https://raw.githubusercontent.com/lucide-icons/lucide/main/icons/chart-bar.svg" width="24" height="24" /> Perfs (RTX 2060 12GB)**:

- 29 pages/minute
- ~1,700 pages/hour
- Batch processing (.zip native)

- **<img src="https://raw.githubusercontent.com/lucide-icons/lucide/main/icons/zap.svg" width="24" height="24" /> 100% Local & Uncensored**: Powered by llama.cpp and Abliterated models. No moralizing, just translation.
- **<img src="https://raw.githubusercontent.com/lucide-icons/lucide/main/icons/eye.svg" width="24" height="24" /> Smart Detection**: Uses YOLOv8 fine-tuned on Manga109 to detect speech bubbles.
  - Smart Box Merging automatically consolidates fragmented vertical text bubbles.
- **<img src="https://raw.githubusercontent.com/lucide-icons/lucide/main/icons/book-open.svg" width="24" height="24" /> Specialized OCR**: Uses MangaOCR to handle vertical Japanese text and handwritten fonts.
- **<img src="https://raw.githubusercontent.com/lucide-icons/lucide/main/icons/brain.svg" width="24" height="24" /> Context-Aware Translation**:
  - Uses Qwen 2.5 7B (Instruction tuned).
  - Custom prompt engineering to handle "Subject-less" Japanese sentences.
  - "Anti-Thinking" regex filters to remove internal LLM monologues.
- **<img src="https://raw.githubusercontent.com/lucide-icons/lucide/main/icons/palette.svg" width="24" height="24" /> Advanced Typesetting**:
  - **NEW (V10)**: **Intelligent Masked Inpainting** - Uses OpenCV threshold detection and cv2.inpaint to remove ONLY dark text pixels, preserving artwork and backgrounds even when bounding boxes overlap.
  - **Pixel-Perfect Wrapping**: Custom algorithm measuring exact pixel width of words to avoid overflow.
  - **Sanitization**: Filters out unsupported characters (emojis, math symbols) to prevent font rendering glitches.
- **<img src="https://raw.githubusercontent.com/lucide-icons/lucide/main/icons/package.svg" width="24" height="24" /> Batch Processing**: Native support for .zip archives (extract → translate → repack).
- **<img src="https://raw.githubusercontent.com/lucide-icons/lucide/main/icons/blocks.svg" width="24" height="24" /> Modular Architecture**: Clean, maintainable codebase with separation of concerns for easy customization and extension.

## <img src="https://raw.githubusercontent.com/lucide-icons/lucide/main/icons/camera.svg" width="24" height="24" /> Examples

See the V10 intelligent masked inpainting in action! These examples showcase the ability to preserve artwork while cleanly removing text.

### Example 1: Naruto

<table>
<tr>
<td width="50%">
<img src="ai-worker/exemples/exemple_naruto.png" alt="Original Naruto page" />
<p align="center"><b>Original (Japanese)</b></p>
</td>
<td width="50%">
<img src="ai-worker/exemples/translated_exemple_naruto.jpg" alt="Translated Naruto page" />
<p align="center"><b>Translated (English)</b></p>
</td>
</tr>
</table>

### Example 2: One Piece

<table>
<tr>
<td width="50%">
<img src="ai-worker/exemples/exemple_one_piece.png" alt="Original One Piece page" />
<p align="center"><b>Original (Japanese)</b></p>
</td>
<td width="50%">
<img src="ai-worker/exemples/translated_exemple_one_piece.jpg" alt="Translated One Piece page" />
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

## <img src="https://raw.githubusercontent.com/lucide-icons/lucide/main/icons/download.svg" width="24" height="24" /> Download Models

Before starting, download the required AI models:

**📦 [Download Models (Google Drive)](https://drive.google.com/drive/folders/18nlj90zpwe57XLsK2slb9OxJwcgNMqQM?usp=sharing)**

Required files:

- `Qwen2.5-7B-Instruct-abliterated-v2.Q4_K_M.gguf` (~4.6 GB) - LLM for translation
- `manga-text-detector.pt` - YOLO model for text bubble detection

Place these files in the `ai-worker/models/` directory.

---

## <img src="https://raw.githubusercontent.com/lucide-icons/lucide/main/icons/rocket.svg" width="24" height="24" /> Quick Start

### Option 1: Full Stack with Docker (Recommended)

Deploy the entire application (frontend + backend + AI worker) with one command:

```bash
# Prerequisites: Docker Desktop installed
# Clone the repository
git clone <repository-url>
cd manga-translator

# Set up AI worker Python environment
cd ai-worker
python -m venv venv
venv\Scripts\activate  # Windows
# or: source venv/bin/activate  # Linux/Mac
pip install -r requirements.txt
cd ..

# Start all services
docker-compose up -d

# Access the application
# Frontend: http://localhost:3000
# Backend API: http://localhost:8080
# Asynq Monitor: http://localhost:8081
```

**Services included:**

- PostgreSQL database
- Redis cache & pub/sub
- Backend Go API
- Asynq worker (background translation jobs)
- Next.js frontend
- Asynq monitoring UI

### Option 2: Local Development

Run each component separately for development:

#### 1. Start Database Services

```bash
cd backend-api
docker-compose up -d postgres redis
```

#### 2. Set Up AI Worker

```bash
cd ../ai-worker
python -m venv venv
venv\Scripts\activate  # Windows
# or: source venv/bin/activate  # Linux/Mac
pip install -r requirements.txt
```

#### 3. Run Backend

```bash
cd ../backend-api
cp .env.example .env
# Edit .env to configure paths (especially PYTHON_PATH)

# Run migrations
migrate -path ./migrations -database "postgres://manga_user:secure_pass@localhost:5432/manga_translator?sslmode=disable" up

# Start API server
go run ./cmd/api

# In another terminal, start worker
go run ./cmd/api --mode=worker
```

#### 4. Run Frontend

```bash
cd ../frontend
npm install  # or: pnpm install
cp .env.local.example .env.local
npm run dev  # or: pnpm dev
```

### Option 3: AI Worker Only (CLI)

Use just the AI worker for command-line batch translation:

```bash
cd ai-worker
python -m venv venv
venv\Scripts\activate
pip install -r requirements.txt

# Translate a single image or ZIP file
python main.py path/to/manga_chapter.zip
```

### System Requirements

- **GPU**: NVIDIA GPU with 6GB+ VRAM (Recommended: 8GB+)
- **CUDA**: CUDA Toolkit 12.x
- **Python**: 3.10+
- **Go**: 1.23+ (for backend development)
- **Node.js**: 20+ (for frontend development)
- **Docker**: Docker Desktop (for containerized deployment)

## <img src="https://raw.githubusercontent.com/lucide-icons/lucide/main/icons/map.svg" width="24" height="24" /> Roadmap

### AI Worker

- [x] Core AI Pipeline (Detection, OCR, Translation, Inpainting)
- [x] GPU Optimization (VRAM management, 4-bit quantization)
- [x] Smart Typesetting (Pixel wrapping, box merging)
- [x] Modular Code Architecture (Config, Services, Utils separation)

### Backend API (v1.0 - Complete ✅)

- [x] Go/Fiber HTTP server with hexagonal architecture
- [x] PostgreSQL database with migrations
- [x] Asynq + Redis job queue
- [x] Python worker subprocess integration
- [x] File upload and validation
- [x] SSE real-time progress tracking
- [x] Redis pub/sub for event broadcasting
- [x] Docker multi-stage build
- [x] Production Docker Compose orchestration
- [ ] Unit & integration tests (future)

### Frontend (v0.1 - Complete ✅)

- [x] Modern UI with Next.js 16 and Tailwind CSS
- [x] Drag-and-drop file upload zone
- [x] API integration with backend
- [x] Real-time SSE progress tracking
- [x] Translation status dashboard
- [x] Interactive result viewer (original/translated toggle)
- [ ] Thumbnail generation (future)
- [ ] User authentication (future)

### Infrastructure (Complete ✅)

- [x] Docker Compose (one-command full stack deployment)
- [x] PostgreSQL + Redis services
- [x] Multi-container orchestration (API + Worker + Frontend)
- [x] Asynq monitoring UI
- [ ] CI/CD pipeline (future)
- [ ] Prometheus/Grafana monitoring (future)

## <img src="https://raw.githubusercontent.com/lucide-icons/lucide/main/icons/users.svg" width="24" height="24" /> Contributing

We welcome contributions from the community! Whether you want to fix bugs, add features, improve documentation, or optimize performance, your help is appreciated.

**Before contributing:**

- 📖 Read our [CONTRIBUTING.md](CONTRIBUTING.md) guide
- 💬 Open an Issue to discuss significant changes (especially for `/ai-worker` modifications)
- ✅ Follow code standards: Ruff (Python), golangci-lint (Go), ESLint (Frontend)
- 🧪 Include tests and documentation with your changes

**Languages:** Contributions can be made in French or English.

## <img src="https://raw.githubusercontent.com/lucide-icons/lucide/main/icons/scale.svg" width="24" height="24" /> License

This project is licensed under a **Custom Non-Commercial Open Source License**.

**You are free to:**

- ✅ Use, modify, and distribute for personal, educational, or research purposes
- ✅ Fork and create derivative works (non-commercial)
- ✅ Contribute back to the project

**Restrictions:**

- ❌ Commercial use requires explicit permission from the author
- 📧 For commercial licensing: contact@antoinerospars.dev

See [LICENSE](LICENSE) for full terms.

Copyright (c) 2026 P4ST4S / Antoine Rospars

## <img src="https://raw.githubusercontent.com/lucide-icons/lucide/main/icons/handshake.svg" width="24" height="24" /> Credits

- **Models**: Qwen (Alibaba Cloud), YOLOv8 (Ultralytics), MangaOCR (kha-white).
- **Tech**: Llama.cpp, PyTorch, Pillow.

---

**Current Version**: V10 (Stable) - Intelligent Masked Inpainting

See [CHANGELOG](ai-worker/CHANGELOG.md) for detailed version history.
