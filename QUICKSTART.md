# 🚀 Quick Start Guide - Manga Translator

This guide will help you get the full stack running in under 5 minutes.

## Prerequisites

- ✅ Docker Desktop installed and running
- ✅ NVIDIA GPU with 6GB+ VRAM
- ✅ Python 3.10+
- ✅ Git

## Step-by-Step Setup

### 1. Clone the Repository

```bash
git clone <repository-url>
cd manga-translator
```

### 2. Set Up AI Worker Python Environment

The AI worker dependencies must be installed in a Python virtual environment **before** building Docker containers:

```bash
cd ai-worker

# Create virtual environment
python -m venv venv

# Activate the environment
# Windows:
venv\Scripts\activate
# Linux/Mac:
source venv/bin/activate

# Install dependencies
pip install -r requirements.txt

# Verify installation
python -c "import torch; print('PyTorch:', torch.__version__)"
python -c "from manga_ocr import MangaOcr; print('MangaOCR: OK')"

# Deactivate and return to root
deactivate
cd ..
```

**Important**: The Docker containers will mount this `venv` directory, so all dependencies must be installed on the host machine first.

### 3. Launch the Full Stack

```bash
# Build and start all services
docker-compose up -d

# This will start:
# - PostgreSQL database (port 5432)
# - Redis cache (port 6379)
# - Backend API (port 8080)
# - Background worker (translation processing)
# - Frontend UI (port 3000)
# - Asynq Monitor (port 8081)
```

### 4. Access the Application

Open your browser and navigate to:

- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:8080/health
- **Asynq Monitor**: http://localhost:8081

### 5. Upload and Translate

1. Go to http://localhost:3000
2. Drag & drop a manga image or ZIP file
3. Click "Upload & Translate"
4. Watch real-time progress in the status dashboard
5. View results when complete

## Monitoring

### View Logs

```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f api
docker-compose logs -f worker
docker-compose logs -f frontend
```

### Check Service Status

```bash
docker-compose ps
```

### Asynq Job Queue Monitor

Visit http://localhost:8081 to see:
- Active jobs
- Pending queue
- Completed jobs
- Failed jobs with retry logic

## Stopping the Stack

```bash
# Stop all services (keeps data)
docker-compose down

# Stop and remove all data (fresh start)
docker-compose down -v
```

## Troubleshooting

### "ModuleNotFoundError" in Worker Logs

**Problem**: Python dependencies not found in container.

**Solution**: Ensure you ran `pip install -r requirements.txt` in the `ai-worker/venv` directory on your host machine before starting Docker.

```bash
cd ai-worker
venv\Scripts\activate  # Windows
pip install -r requirements.txt
deactivate
cd ..
docker-compose restart worker
```

### Port Already in Use

**Problem**: Error like `bind: address already in use`

**Solution**: Stop the conflicting service or change the port in `docker-compose.yml`:

```yaml
services:
  frontend:
    ports:
      - "3001:3000"  # Changed from 3000:3000
```

### GPU Not Detected

**Problem**: Translation is very slow or fails.

**Solution**: Ensure NVIDIA drivers and CUDA Toolkit are installed:

```bash
# Check NVIDIA driver
nvidia-smi

# Check CUDA
nvcc --version
```

If using Docker with GPU support, ensure Docker is configured for NVIDIA runtime (requires `nvidia-docker2`).

### Database Connection Failed

**Problem**: Backend can't connect to PostgreSQL.

**Solution**: Wait for PostgreSQL to fully start:

```bash
# Check PostgreSQL logs
docker-compose logs postgres

# Restart the API
docker-compose restart api
```

## Next Steps

- 📚 Read the full [README.md](README.md) for architecture details
- 🔧 Check [backend-api/README.md](backend-api/README.md) for API documentation
- 🤖 See [ai-worker/README.md](ai-worker/README.md) for AI pipeline details
- 🎨 Explore [frontend/README.md](frontend/README.md) for UI customization

## Development Mode

If you want to develop locally without Docker:

```bash
# Terminal 1: Start databases only
cd backend-api
docker-compose up -d postgres redis

# Terminal 2: Run backend API
cd backend-api
go run ./cmd/api

# Terminal 3: Run background worker
cd backend-api
go run ./cmd/api --mode=worker

# Terminal 4: Run frontend
cd frontend
npm install
npm run dev
```

Access at:
- Frontend: http://localhost:3000
- Backend: http://localhost:8080

---

**Need help?** Check the individual README files in each module directory or open an issue.
