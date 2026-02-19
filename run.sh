#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo "=== Manga Translator - Starting all services ==="
echo

# Check Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "[ERROR] Docker is not running. Please start Docker first."
    exit 1
fi

# Start Docker services (infra + api + frontend)
echo "[1/2] Starting Docker services..."
docker compose -f "$SCRIPT_DIR/docker-compose.yml" up -d --build
echo "[OK] Docker services started."
echo

# Wait for services to be ready
echo "Waiting for services to be ready..."
sleep 5

# Start the Go worker in a new terminal (background)
echo "[2/2] Starting Go worker (host)..."
cd "$SCRIPT_DIR/backend-api"

# Detect terminal emulator and open new window, fallback to background process
if command -v gnome-terminal > /dev/null 2>&1; then
    gnome-terminal -- bash -c "cd '$SCRIPT_DIR/backend-api' && go run ./cmd/api --mode=worker; exec bash"
elif command -v xterm > /dev/null 2>&1; then
    xterm -title "Manga Translator - Worker" -e "cd '$SCRIPT_DIR/backend-api' && go run ./cmd/api --mode=worker; exec bash" &
elif command -v konsole > /dev/null 2>&1; then
    konsole --new-tab -e bash -c "cd '$SCRIPT_DIR/backend-api' && go run ./cmd/api --mode=worker; exec bash" &
else
    # No GUI terminal — run worker in background and log to file
    echo "[INFO] No GUI terminal found. Running worker in background (logs: $SCRIPT_DIR/worker.log)"
    cd "$SCRIPT_DIR/backend-api"
    go run ./cmd/api --mode=worker > "$SCRIPT_DIR/worker.log" 2>&1 &
    WORKER_PID=$!
    echo "[OK] Worker started (PID $WORKER_PID). Tail logs with: tail -f $SCRIPT_DIR/worker.log"
fi

# Open default browser
sleep 3
if command -v xdg-open > /dev/null 2>&1; then
    xdg-open "http://localhost:3000" > /dev/null 2>&1 &
elif command -v open > /dev/null 2>&1; then
    open "http://localhost:3000"
fi

echo
echo "=== All services started ==="
echo "  Frontend:      http://localhost:3000"
echo "  Backend API:   http://localhost:8080"
echo "  Asynq Monitor: http://localhost:8081"
echo
echo "Press Ctrl+C to stop Docker services."
echo

cd "$SCRIPT_DIR"
docker compose logs -f --tail=20
