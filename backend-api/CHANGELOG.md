# Changelog

All notable changes to the Manga Translator Backend API will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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

- **v2.0.0** (2026-02-01): Real-time SSE progress, ZIP extraction, subdirectory support
- **v1.0.0** (2025-12-15): Initial release with core features
