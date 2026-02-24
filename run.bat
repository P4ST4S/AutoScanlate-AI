@echo off
setlocal

echo === Manga Translator - Starting all services ===
echo.

:: Check Docker is running
docker info >nul 2>&1
if errorlevel 1 (
    echo [ERROR] Docker is not running. Please start Docker Desktop first.
    pause
    exit /b 1
)

:: Create storage directories if not present
echo Creating storage directories if needed...
if not exist "storage\uploads"    mkdir "storage\uploads"
if not exist "storage\originals"  mkdir "storage\originals"
if not exist "storage\translated" mkdir "storage\translated"
if not exist "storage\temp"       mkdir "storage\temp"
echo [OK] Storage directories ready.
echo.

:: Start Docker services (infra + api + frontend)
echo [1/2] Starting Docker services...
docker compose up -d --build
if errorlevel 1 (
    echo [ERROR] Failed to start Docker services.
    pause
    exit /b 1
)
echo [OK] Docker services started.
echo.

:: Wait a moment for services to be healthy
echo Waiting for services to be ready...
timeout /t 5 /nobreak >nul

:: Start the Go worker in a new window
echo [2/2] Starting Go worker (host)...
start "Manga Translator - Worker" cmd /k "cd /d "%~dp0backend-api" && go run ./cmd/api --mode=worker"

:: Open browser after a short delay to let the frontend finish starting
timeout /t 3 /nobreak >nul
start "" "http://localhost:3000"

echo.
echo === All services started ===
echo   Frontend:      http://localhost:3000
echo   Backend API:   http://localhost:8080
echo   Asynq Monitor: http://localhost:8081
echo.
echo The worker is running in a separate window.
echo Close this window or press Ctrl+C to stop Docker services.
echo.

:: Keep window open and show docker logs
docker compose logs -f --tail=20

endlocal
