#!/usr/bin/env bash
set -e

CONTAINER_NAME="goblog-postgres"
DB_USER="bloguser"
DB_PASS="blogpass"
DB_NAME="blogdb"
DB_PORT="5432"

echo "=== GoBlog - Start Script ==="

# Check if podman is available
if ! command -v podman &> /dev/null; then
    echo "Error: podman is not installed."
    exit 1
fi

# Check if container already exists
if podman container exists "$CONTAINER_NAME" 2>/dev/null; then
    STATE=$(podman inspect --format '{{.State.Status}}' "$CONTAINER_NAME" 2>/dev/null)
    if [ "$STATE" == "running" ]; then
        echo "PostgreSQL container '$CONTAINER_NAME' is already running."
    else
        echo "Starting existing PostgreSQL container..."
        podman start "$CONTAINER_NAME"
    fi
else
    echo "Creating PostgreSQL container with podman..."
    podman run -d \
        --name "$CONTAINER_NAME" \
        -e POSTGRES_USER="$DB_USER" \
        -e POSTGRES_PASSWORD="$DB_PASS" \
        -e POSTGRES_DB="$DB_NAME" \
        -p "$DB_PORT":5432 \
        docker.io/library/postgres:16-alpine
fi

# Wait for PostgreSQL to be ready
echo "Waiting for PostgreSQL to be ready..."
for i in $(seq 1 30); do
    if podman exec "$CONTAINER_NAME" pg_isready -U "$DB_USER" -d "$DB_NAME" &>/dev/null; then
        echo "PostgreSQL is ready!"
        break
    fi
    sleep 1
done

# Download Go dependencies
echo "Downloading Go dependencies..."
export PATH="$HOME/.local/go/bin:$PATH"
go mod tidy

# Build and run the Go app
echo "Building and starting GoBlog..."
go run main.go
