#!/usr/bin/env bash
set -e

CONTAINER_NAME="goblog-postgres"

echo "=== GoBlog - Stop Script ==="

if podman container exists "$CONTAINER_NAME" 2>/dev/null; then
    echo "Stopping and removing PostgreSQL container..."
    podman stop "$CONTAINER_NAME" 2>/dev/null || true
    podman rm "$CONTAINER_NAME" 2>/dev/null || true
    echo "Done. Container removed."
else
    echo "Container '$CONTAINER_NAME' does not exist."
fi
