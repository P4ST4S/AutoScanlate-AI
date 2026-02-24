# Changelog

All notable changes to the Manga Translator Backend API will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [2.1.0] - 2026-02-24

### Added

- **`run.bat` launcher**: One-command startup script for Windows — starts Docker services, creates storage directories, launches the Go worker in a separate window, and opens the browser automatically
- **`run.sh` launcher**: Equivalent startup script for Linux/Mac
- **Auto storage directory creation**: `run.bat`/`run.sh` now create `storage/uploads`, `storage/originals`, `storage/translated`, and `storage/temp` if missing, preventing Docker bind-mount failures
- **`PYTHONIOENCODING=utf-8` env var**: Added to the Python subprocess environment in `executor.go` to prevent `UnicodeEncodeError` crashes when emoji characters (e.g. `⏳`, `✅`) are printed to piped stdout on Windows

### Fixed

- **Python worker crash (`exit status 0xc0000005`)**: Worker subprocess was crashing with an access violation on Windows when spawned by the Go worker. Root cause: Python fell back to `cp1252` encoding for piped stdout, and emoji characters in `translation.py` caused a fatal `UnicodeEncodeError`. Fixed by setting `PYTHONIOENCODING=utf-8` in the subprocess environment.
- **MangaOCR GPU conflict**: `MangaOcr()` now initializes with `force_cpu=True` to avoid a CUDA context conflict between llama-cpp-python and the transformers ViT model when both are loaded in the same process on Windows
- **Docker storage bind-mount failure** (`mkdir /app/storage: file exists`): Caused by a missing `./storage` directory on the host — Docker auto-created it as a file instead of a directory. Fixed by pre-creating the directory structure in `run.bat`/`run.sh` and via `docker-compose.yml` volume mount ordering

### Changed

- **AI pipeline loading order**: LLM is now loaded before YOLO and MangaOCR to avoid GPU context initialization conflicts on Windows
- **`manga-ocr` pinned to `0.1.13`** in `requirements.txt` for stability

### Technical Details

- `executor.go`: Added `PYTHONIOENCODING=utf-8` alongside existing `PYTHONUNBUFFERED=1`
- `pipeline.py`: Changed `MangaOcr()` → `MangaOcr(force_cpu=True)` and reordered model loading (LLM → YOLO → OCR)
- `run.bat`: Added storage directory pre-creation block before Docker Compose startup

## [2.0.0] - 2026-02-01

### Added

- **Real-time SSE Progress Tracking**: Live page-by-page translation progress with Server-Sent Events
- **ZIP File Extraction**: Automatic extraction of uploaded ZIP archives to originals/ and translated/ directories
- **Subdirectory Support**: Full support for nested folder structures within ZIP archives
- **Instant Page Counting**: Display total pages immediately on upload via `countImagesInZip()`
- **Redis Pub/Sub Integration**: Reliable progress broadcasting from worker to SSE clients
- **Progress Pattern Parsing**: Regex-based parsing of `PROGRESS: XX% - message` format from Python worker
- **Wildcard File Routing**: Changed from `:filename` to `*` parameter for nested file paths (`/api/files/:id/:type/*`)
- **Enhanced Logging**: Detailed debug output for progress tracking and file operations

### Changed

- **SSE Cleanup Logic**: Moved `defer subscriber.Close()` and `cancel()` inside `SetBodyStreamWriter` goroutine for proper connection lifecycle
- **File Path Handling**: Use `filepath.Rel()` instead of `filepath.Base()` to preserve subdirectory structure
- **Python Output Buffering**: Set `PYTHONUNBUFFERED=1` environment variable for real-time stdout streaming
- **Log Level**: Changed worker output from `Debug` to `Info` level for better visibility
- **Database Schema**: Added relative path support in Result entries for ZIP-extracted images

### Fixed

- **SSE Connection Stability**: Connections now stay open properly instead of closing immediately
- **Frontend JSON Parsing**: SSE events no longer crash on malformed data
- **Python Progress Messages**: Real-time progress now appears immediately instead of being buffered
- **404 Errors**: Files in ZIP subdirectories now load correctly
- **Page Count Display**: Shows correct count during translation instead of "0 pages"
- **ZIP Processing**: Worker now properly extracts archives instead of copying them as-is

### Technical Details

- Improved error handling throughout SSE pipeline
- Type-safe progress callbacks in worker executor
- Proper cleanup of goroutines and Redis subscriptions
- Enhanced file validation with recursive image collection
- Sorted file ordering in `collectImageFiles()` for consistent page numbering

## [1.0.0] - 2025-12-15

### Added

- **Core Architecture**: Hexagonal/Clean Architecture with Fiber v3 framework
- **Database Layer**: PostgreSQL 16+ integration with pgx driver
- **Migration System**: Up/down SQL migrations for Request and Result tables
- **Job Queue**: Asynq + Redis for asynchronous job processing
- **Upload Handler**: File validation and storage for images and ZIP archives
- **Python Worker Integration**: Subprocess executor for AI worker communication
- **Basic SSE Support**: Initial Server-Sent Events implementation for progress tracking
- **File Serving**: Static file endpoints for originals, translated, and upload files
- **Docker Support**: Multi-stage Dockerfile and docker-compose.yml
- **CORS Configuration**: Cross-origin support for frontend integration

### Features

- Request/Result domain entities with repository pattern
- HTTP handlers with proper error handling
- Worker mode (separate process) for job processing
- Progress parsing from Python stdout/stderr
- Storage organization (originals/, translated/, temp/, uploads/)
- Environment-based configuration
- Graceful shutdown handling

---

## Version History

- **v2.1.0** (2026-02-24): Windows subprocess fix, storage directory auto-creation, `run.bat`/`run.sh` launchers
- **v2.0.0** (2026-02-01): Real-time SSE progress, ZIP extraction, subdirectory support
- **v1.0.0** (2025-12-15): Initial release with core features
