#!/bin/bash

# Create logs and pids directories if they don't exist
mkdir -p logs .pids

# Stop any existing processes first
./scripts/dev-down.sh 2>/dev/null

echo "Starting Food Planning App development environment..."

# Check if Docker Compose is available for PostgreSQL
if command -v docker-compose &> /dev/null; then
    echo "Starting PostgreSQL database with Docker..."
    docker-compose up -d database
    # Wait for database to be ready
    echo "Waiting for database to be ready..."
    sleep 10
    export DB_TYPE=postgres
else
    echo "Docker Compose not found. Using SQLite database for development..."
    export DB_TYPE=sqlite
    export DB_PATH="./backend/food_app.db"
fi

# Start backend
echo "Starting Go backend..."
cd backend
go mod tidy
go run . > ../logs/backend.log 2>&1 &
echo $! > ../.pids/backend.pid
cd ..

# Wait for backend to start
echo "Waiting for backend to start..."
sleep 5

# Start frontend
echo "Starting React frontend..."
cd frontend
npm install
npm run dev > ../logs/frontend.log 2>&1 &
echo $! > ../.pids/frontend.pid
cd ..

echo "Development environment started!"
echo "Frontend: http://localhost:3000"
echo "Backend: http://localhost:8080"
echo "Database: localhost:5432"
echo ""
echo "To view logs:"
echo "  Backend: tail -f logs/backend.log"
echo "  Frontend: tail -f logs/frontend.log"
echo ""
echo "To stop: ./scripts/dev-down.sh"
