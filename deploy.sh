#!/usr/bin/env bash
set -e

echo "Building GoBlog app..."
podman build -t goblog-app .

# Check if pod exists and remove
if podman pod exists goblog; then
    echo "Removing existing pod goblog..."
    podman pod rm -f goblog
fi

echo "Creating podman pod goblog..."
podman pod create --name goblog -p 8080:8080

echo "Starting Postgres container..."
podman run -d --pod goblog --name goblog-postgres \
  -e POSTGRES_USER=bloguser \
  -e POSTGRES_PASSWORD=blogpass \
  -e POSTGRES_DB=blogdb \
  docker.io/library/postgres:16-alpine

echo "Waiting for postgres..."
for i in $(seq 1 10); do
    if podman exec goblog-postgres pg_isready -U bloguser -d blogdb &>/dev/null; then
        echo "PostgreSQL is ready!"
        break
    fi
    sleep 1
done

echo "Starting GoBlog app..."
podman run -d --pod goblog --name goblog-app \
  -e DB_HOST=127.0.0.1 \
  -e DB_PORT=5432 \
  -e DB_USER=bloguser \
  -e DB_PASSWORD=blogpass \
  -e DB_NAME=blogdb \
  -e PORT=8080 \
  goblog-app

echo "GoBlog is running at http://localhost:8080"
