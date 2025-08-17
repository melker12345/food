#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
PID_DIR="$ROOT_DIR/.pids"

kill_pid_file() {
	local file="$1"
	if [ -f "$file" ]; then
		local pid
		pid=$(cat "$file" || true)
		if [ -n "$pid" ] && kill -0 "$pid" >/dev/null 2>&1; then
			echo "[dev-down] Killing PID $pid from $(basename "$file")"
			kill "$pid" >/dev/null 2>&1 || true
			sleep 1
			kill -9 "$pid" >/dev/null 2>&1 || true
		fi
		rm -f "$file"
	fi
}

kill_port() {
	local port="$1"
	local pids
	pids=$(lsof -ti:"$port" 2>/dev/null || true)
	if [ -n "$pids" ]; then
		echo "[dev-down] Killing processes on port $port: $pids"
		kill $pids >/dev/null 2>&1 || true
		sleep 1
		kill -9 $pids >/dev/null 2>&1 || true
	fi
}

kill_pid_file "$PID_DIR/frontend.pid"
kill_pid_file "$PID_DIR/backend.pid"

# Fallback: ensure ports are free even if PID files are missing/stale
kill_port 3000  # React frontend
kill_port 4000  # API

echo "[dev-down] Stopping docker compose..."
docker compose -f "$ROOT_DIR/docker-compose.yml" down || true

echo "[dev-down] Done."


