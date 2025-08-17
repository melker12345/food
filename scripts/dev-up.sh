#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
LOG_DIR="$ROOT_DIR/logs"
PID_DIR="$ROOT_DIR/.pids"
mkdir -p "$LOG_DIR" "$PID_DIR"

echo "[dev-up] Checking Docker daemon..."
if ! docker info >/dev/null 2>&1; then
	echo "[dev-up] Docker daemon is not running. Start Docker and re-run."
	exit 1
fi

echo "[dev-up] Starting MongoDB (docker compose up -d)..."
docker compose -f "$ROOT_DIR/docker-compose.yml" up -d

echo "[dev-up] Building Go backend..."
( cd "$ROOT_DIR/backend" && go build -o food-api . )

echo "[dev-up] Starting backend API..."
( cd "$ROOT_DIR/backend" && nohup ./food-api >"$LOG_DIR/backend.log" 2>&1 & echo $! > "$PID_DIR/backend.pid" )

echo "[dev-up] Starting React frontend..."
( cd "$ROOT_DIR/frontend" && nohup npm start >"$LOG_DIR/frontend.log" 2>&1 & echo $! > "$PID_DIR/frontend.pid" )

echo "[dev-up] Done."
echo "API:   http://localhost:4000"
echo "Web:   http://localhost:3000"
echo "Logs:  $LOG_DIR"


