# Manga Translator Backend API

High-performance Go backend API for the Manga Translator project. Orchestrates file uploads, job queuing, and translation processing between the Next.js frontend and Python AI worker.

## Architecture

- **Framework**: Fiber v3 (high-performance HTTP framework)
- **Database**: PostgreSQL 16+ (via pgx driver)
- **Queue**: Asynq + Redis (for async job processing)
- **Pattern**: Hexagonal/Clean Architecture
- **Real-time**: Server-Sent Events (SSE) for progress tracking

## Project Status

### Phase 1: Foundation ✅ Complete
- [x] Project structure and configuration
- [x] PostgreSQL database with migrations
- [x] Domain entities (Request, Result)
- [x] Repository layer (PostgreSQL)
- [x] HTTP server with Fiber + CORS
- [x] Basic API endpoints
- [x] Docker Compose setup

### Phase 2: Core Features (In Progress)
- [ ] Upload handler with file validation
- [ ] Asynq job queue integration
- [ ] Python worker subprocess executor
- [ ] Request/Result endpoints

### Phase 3: Real-time & Advanced
- [ ] SSE progress tracking
- [ ] Redis pub/sub integration
- [ ] File serving endpoint
- [ ] Cleanup service

### Phase 4: Production
- [ ] Docker multi-stage build
- [ ] Unit + integration tests
- [ ] Production deployment guide

## Prerequisites

- **Go**: 1.23+ (automatically upgraded to 1.24 by dependencies)
- **Docker Desktop**: For PostgreSQL and Redis
- **Python**: 3.10+ (for AI worker integration)

## Quick Start

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

# Edit .env with your settings (defaults work for local development)
```

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

# Run (Windows)
./api.exe

# Run (Linux/Mac)
./api
```

The API will start on `http://localhost:8080`.

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

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | 8080 |
| `DB_HOST` | PostgreSQL host | localhost |
| `DB_PORT` | PostgreSQL port | 5432 |
| `REDIS_ADDR` | Redis address | localhost:6379 |
| `PYTHON_PATH` | Python executable path | python |
| `WORKER_PATH` | AI worker directory | ../ai-worker |
| `WORKER_CONCURRENCY` | Max concurrent jobs | 1 |
| `MAX_UPLOAD_SIZE` | Max file size (bytes) | 104857600 (100MB) |
| `CORS_ORIGINS` | Allowed CORS origins | http://localhost:3000 |

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

### Build Docker Image

```bash
docker build -t manga-translator-api .
```

### Run with Docker Compose

```bash
# Start all services (PostgreSQL + Redis + API)
docker-compose up -d

# View logs
docker-compose logs -f api

# Stop services
docker-compose down

# Stop and remove volumes
docker-compose down -v
```

## Troubleshooting

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
