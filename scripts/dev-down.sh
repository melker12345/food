#!/bin/bash

echo "Stopping Food Planning App development environment..."

# Stop frontend
if [ -f .pids/frontend.pid ]; then
    echo "Stopping React frontend..."
    kill $(cat .pids/frontend.pid) 2>/dev/null
    rm .pids/frontend.pid
fi

# Stop backend
if [ -f .pids/backend.pid ]; then
    echo "Stopping Go backend..."
    kill $(cat .pids/backend.pid) 2>/dev/null
    rm .pids/backend.pid
fi

# Stop database
echo "Stopping PostgreSQL database..."
docker-compose down

echo "Development environment stopped!"
