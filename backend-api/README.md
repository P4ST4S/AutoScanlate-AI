# Manga Translator Backend API

**Version 2.0** - High-performance Go backend API for the Manga Translator project. Orchestrates file uploads, job queuing, and translation processing between the Next.js frontend and Python AI worker.

## Architecture

- **Framework**: Fiber v3 (high-performance HTTP framework)
- **Database**: PostgreSQL 16+ (via pgx driver)
- **Queue**: Asynq + Redis (for async job processing)
- **Pattern**: Hexagonal/Clean Architecture
- **Real-time**: Server-Sent Events (SSE) for live progress tracking

## What's New in v2.0

### Real-time Progress Tracking

- ✅ **Live SSE streaming** with page-by-page progress updates
- ✅ **Python stdout unbuffering** (`PYTHONUNBUFFERED=1`) for instant progress reporting
- ✅ **Redis pub/sub integration** for reliable progress broadcasting
- ✅ **Connection stability** with proper error handling and cleanup

### Enhanced ZIP Support

- ✅ **Automatic ZIP extraction** to originals/ and translated/ directories
- ✅ **Subdirectory support** with proper relative path handling
- ✅ **Instant page counting** - displays total pages immediately on upload
- ✅ **Wildcard routing** (`/api/files/:id/:type/*`) for nested file structures
- ✅ **Database cataloging** - all pages indexed with proper paths

### Improved Architecture

- ✅ **Proper resource cleanup** in SSE streams (defer placement fixes)
- ✅ **Enhanced error logging** with detailed debugging information
- ✅ **Progress parsing** with regex patterns for structured updates
- ✅ **Type-safe progress callbacks** throughout the pipeline

## Project Status

### Phase 1: Foundation ✅ Complete

- [x] Project structure and configuration
- [x] PostgreSQL database with migrations
- [x] Domain entities (Request, Result)
- [x] Repository layer (PostgreSQL)
- [x] HTTP server with Fiber + CORS
- [x] Basic API endpoints
- [x] Docker Compose setup

### Phase 2: Core Features ✅ Complete

- [x] Upload handler with file validation
- [x] Asynq job queue integration (client + server)
- [x] Python worker subprocess executor
- [x] Progress parsing from stdout/stderr
- [x] Worker mode (separate process)
- [x] End-to-end translation flow

### Phase 3: Real-time & Advanced ✅ Complete (v2.0)

- [x] SSE progress tracking with live updates
- [x] Redis pub/sub integration
- [x] Real-time progress broadcasting
- [x] Event streaming (connected/progress/complete/error)
- [x] ZIP extraction and cataloging
- [x] Subdirectory support for complex archives
- [ ] Cleanup service (Phase 4)

### Phase 4: Production ✅ Complete

- [x] Docker multi-stage build
- [x] Docker Compose orchestration
- [x] Production deployment guide
- [ ] Unit + integration tests (future)

## Prerequisites

- **Go**: 1.23+ (automatically upgraded to 1.24 by dependencies)
- **Docker Desktop**: For PostgreSQL and Redis
- **Python**: 3.10+ (for AI worker integration)

## Quick Start

### 0. Prerequisites - AI Worker Setup

**Before starting the backend**, ensure the AI worker is set up with its venv:

```bash
cd ../ai-worker

# Create virtual environment if not exists
python -m venv venv

# Activate venv
# Windows:
venv\Scripts\activate
# Linux/Mac:
source venv/bin/activate

# Install dependencies
pip install -r requirements.txt

# Verify installation
python -c "import torch; print('PyTorch OK')"

cd ../backend-api
```

### 1. Clone and Navigate

```bash
cd backend-api
```

### 2. Install Dependencies

```bash
go mod download
```

### 3. Configure Environment

```bash
# Copy example environment file
cp .env.example .env

# IMPORTANT: Verify PYTHON_PATH points to the ai-worker venv
# Windows (default):
# PYTHON_PATH=../ai-worker/venv/Scripts/python.exe
#
# Linux/Mac:
# PYTHON_PATH=../ai-worker/venv/bin/python

# Edit .env with your settings if needed
```

**⚠️ Critical**: The `PYTHON_PATH` must point to the Python executable inside the `ai-worker/venv` directory, not the system Python. The ai-worker dependencies (PyTorch, llama-cpp-python, etc.) are installed in this virtual environment.

### 4. Start Database and Redis

```bash
# Start PostgreSQL and Redis with Docker Compose
docker-compose up -d postgres redis

# Verify services are running
docker-compose ps
```

### 5. Run Database Migrations

```bash
# Install golang-migrate tool (one-time setup)
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Run migrations
migrate -path ./migrations -database "postgres://manga_user:secure_pass@localhost:5432/manga_translator?sslmode=disable" up

# To rollback migrations
# migrate -path ./migrations -database "..." down
```

### 6. Build and Run

```bash
# Build
go build -o api ./cmd/api

# Run API server (Windows)
./api.exe

# Run API server (Linux/Mac)
./api

# Run in worker mode (processes translation jobs)
./api --mode=worker

# Or run both simultaneously (recommended)
# Terminal 1:
./api

# Terminal 2:
./api --mode=worker
```

The API will start on `http://localhost:8080`.

**Note**: You need BOTH the API server and worker running for full functionality. The API server handles HTTP requests and enqueues jobs, while the worker processes the translation tasks.

## API Endpoints

### Health Check

```
GET /health
```

Returns server health status.

### Upload Translation Request

```
POST /api/translate
Content-Type: multipart/form-data

Body:
  files: File[] (max 10 files, .zip/.png/.jpg/.jpeg/.webp)

Response 201:
{
  "id": "uuid",
  "filename": "chapter.zip",
  "status": "queued",
  "progress": 0,
  "pageCount": 0,
  "createdAt": "2026-02-01T10:00:00Z"
}
```

### List Requests

```
GET /api/requests?status=processing&limit=20&offset=0

Response 200:
{
  "requests": [...],
  "total": 42,
  "limit": 20,
  "offset": 0
}
```

### Get Request Status

```
GET /api/requests/:id

Response 200:
{
  "id": "uuid",
  "filename": "chapter.zip",
  "status": "completed",
  "progress": 100,
  "pageCount": 18,
  "createdAt": "...",
  "completedAt": "..."
}
```

### Get Translation Results

```
GET /api/results/:id

Response 200:
{
  "requestId": "uuid",
  "pages": [
    {
      "pageNumber": 1,
      "original": "/api/files/uuid/originals/page_001.jpg",
      "translated": "/api/files/uuid/translated/page_001.jpg"
    }
  ]
}
```

### Real-time Progress Updates (SSE)

```
GET /api/requests/:id/events
Accept: text/event-stream

Event Stream:
event: connected
data: {"status":"queued","progress":0,"message":"Connected to progress stream"}

event: progress
data: {"status":"processing","progress":25,"message":"Processing page 5/18"}

event: progress
data: {"status":"processing","progress":50,"message":"Processing page 9/18"}

event: complete
data: {"status":"completed","progress":100,"message":"Translation completed successfully"}

event: error
data: {"status":"failed","progress":0,"message":"Translation failed: ..."}
```

**Usage with JavaScript:**

```javascript
const eventSource = new EventSource(
  `http://localhost:8080/api/requests/${requestId}/events`,
);

eventSource.addEventListener("connected", (e) => {
  const data = JSON.parse(e.data);
  console.log("Connected:", data);
});

eventSource.addEventListener("progress", (e) => {
  const data = JSON.parse(e.data);
  console.log("Progress:", data.progress, "%", data.message);
});

eventSource.addEventListener("complete", (e) => {
  const data = JSON.parse(e.data);
  console.log("Completed!");
  eventSource.close();
});

eventSource.addEventListener("error", (e) => {
  const data = JSON.parse(e.data);
  console.error("Error:", data.message);
  eventSource.close();
});
```

### Serve Files

```
GET /api/files/:requestId/:type/:filename

Parameters:
  - type: uploads | originals | translated
```

## Development

### Project Structure

```
backend-api/
├── cmd/
│   └── api/
│       └── main.go              # Application entry point
├── internal/
│   ├── domain/                  # Business entities
│   │   ├── request.go
│   │   ├── result.go
│   │   └── errors.go
│   ├── ports/                   # Interfaces (dependency inversion)
│   │   ├── repository.go
│   │   ├── queue.go
│   │   ├── worker.go
│   │   └── storage.go
│   ├── adapters/                # External implementations
│   │   ├── http/
│   │   │   ├── router.go
│   │   │   ├── handlers/
│   │   │   └── middleware/
│   │   ├── repository/postgres/
│   │   ├── queue/asynq/
│   │   ├── worker/python/
│   │   └── storage/local/
│   ├── application/             # Use cases
│   └── infrastructure/
│       ├── config/              # Configuration loader
│       ├── logger/              # Zap logger
│       └── database/            # PostgreSQL setup
├── migrations/                  # SQL migrations
├── storage/                     # File storage
└── docker-compose.yml
```

### Configuration

All configuration is managed through environment variables (see `.env.example`):

| Variable             | Description                                  | Default                              |
| -------------------- | -------------------------------------------- | ------------------------------------ |
| `PORT`               | Server port                                  | 8080                                 |
| `DB_HOST`            | PostgreSQL host                              | localhost                            |
| `DB_PORT`            | PostgreSQL port                              | 5432                                 |
| `REDIS_ADDR`         | Redis address                                | localhost:6379                       |
| `PYTHON_PATH`        | Python executable path (⚠️ **must be venv**) | ../ai-worker/venv/Scripts/python.exe |
| `WORKER_PATH`        | AI worker directory                          | ../ai-worker                         |
| `WORKER_CONCURRENCY` | Max concurrent jobs                          | 1                                    |
| `MAX_UPLOAD_SIZE`    | Max file size (bytes)                        | 104857600 (100MB)                    |
| `CORS_ORIGINS`       | Allowed CORS origins                         | http://localhost:3000                |

### Running Tests

```bash
# Run unit tests
go test ./internal/domain/...
go test ./internal/application/...

# Run integration tests (requires Docker)
go test ./internal/adapters/repository/...

# Run all tests with coverage
go test -cover ./...
```

### Database Migrations

Create a new migration:

```bash
migrate create -ext sql -dir migrations -seq migration_name
```

This creates two files:

- `XXX_migration_name.up.sql` - Forward migration
- `XXX_migration_name.down.sql` - Rollback migration

### Logging

The application uses structured logging with Zap:

```go
logger.Info("message",
    zap.String("key", "value"),
    zap.Int("count", 42),
)
```

Log levels: `debug`, `info`, `warn`, `error`

Configure via `LOG_LEVEL` environment variable.

## Docker Deployment

### Option 1: Full Stack with Root docker-compose.yml (Recommended)

From the project root, you can deploy the entire stack (frontend + backend + services):

```bash
cd ..  # Navigate to project root

# Build and start all services (frontend, backend, postgres, redis, worker, asynqmon)
docker-compose up -d

# View logs for all services
docker-compose logs -f

# View logs for specific service
docker-compose logs -f api
docker-compose logs -f worker
docker-compose logs -f frontend

# Stop all services
docker-compose down

# Stop and remove volumes (will delete database data)
docker-compose down -v
```

**Services included:**

- `postgres` - PostgreSQL 16 database (port 5432)
- `redis` - Redis cache & pub/sub (port 6379)
- `api` - Backend API server (port 8080)
- `worker` - Asynq worker process (background jobs)
- `frontend` - Next.js frontend (port 3000)
- `asynqmon` - Asynq monitoring UI (port 8081)

**Access the application:**

- Frontend: http://localhost:3000
- Backend API: http://localhost:8080
- Asynq Monitor: http://localhost:8081

### Option 2: Backend Services Only (Local Development)

If you want to run the backend services with Docker but develop the API locally:

```bash
# From backend-api directory
docker-compose up -d postgres redis

# Run API locally
go run ./cmd/api

# Run worker locally (in another terminal)
go run ./cmd/api --mode=worker
```

### Build Docker Image Manually

```bash
docker build -t manga-translator-api .
```

### Environment Variables for Docker

When using Docker Compose from the root, environment variables are pre-configured. For custom deployments, ensure these are set:

```yaml
environment:
  DB_HOST: postgres # Use service name, not localhost
  REDIS_ADDR: redis:6379 # Use service name, not localhost
  PYTHON_PATH: /app/ai-worker/venv/bin/python # Linux path in container
  WORKER_PATH: /app/ai-worker
  CORS_ORIGINS: http://localhost:3000
```

**⚠️ Important Notes:**

- The Python venv must be created on the host before building Docker images
- Storage volumes are mounted from the host to persist uploaded/translated files
- Database migrations run automatically on first startup

## Troubleshooting

### Python Worker Fails / ModuleNotFoundError

**Symptom**: Worker logs show errors like `ModuleNotFoundError: No module named 'torch'` or similar.

**Cause**: The `PYTHON_PATH` is pointing to the system Python instead of the ai-worker venv.

**Solution**:

```bash
# Verify your .env file has the correct path
# Windows:
PYTHON_PATH=../ai-worker/venv/Scripts/python.exe

# Linux/Mac:
PYTHON_PATH=../ai-worker/venv/bin/python

# Test the Python path manually
../ai-worker/venv/Scripts/python.exe -c "import torch; print('OK')"
```

If you haven't created the venv yet:

```bash
cd ../ai-worker
python -m venv venv
venv\Scripts\activate  # Windows
# or
source venv/bin/activate  # Linux/Mac
pip install -r requirements.txt
```

### Database Connection Failed

```bash
# Check if PostgreSQL is running
docker-compose ps postgres

# Check logs
docker-compose logs postgres

# Restart PostgreSQL
docker-compose restart postgres
```

### Migration Errors

```bash
# Check current migration version
migrate -path ./migrations -database "postgres://..." version

# Force version (if stuck)
migrate -path ./migrations -database "postgres://..." force VERSION
```

### Port Already in Use

```bash
# Change PORT in .env file
PORT=8081

# Or stop conflicting process
# Windows:
netstat -ano | findstr :8080
taskkill /PID <PID> /F

# Linux/Mac:
lsof -ti:8080 | xargs kill -9
```

## Next Steps

1. **Phase 2**: Implement Asynq job queue and Python worker integration
2. **Phase 3**: Add SSE real-time progress tracking
3. **Phase 4**: Production deployment with monitoring

## Contributing

This backend follows Clean Architecture principles:

- **Domain layer** contains pure business logic (no external dependencies)
- **Ports** define interfaces
- **Adapters** implement interfaces
- **Dependencies flow inward** (adapters → ports → domain)

When adding new features:

1. Define domain entities in `internal/domain/`
2. Define interfaces in `internal/ports/`
3. Implement adapters in `internal/adapters/`
4. Wire dependencies in `cmd/api/main.go`

## License

Part of the Manga Translator project.
