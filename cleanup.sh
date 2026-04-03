#!/usr/bin/env bash
set -e

if podman pod exists goblog; then
    echo "Stopping and removing pod goblog..."
    podman pod rm -f goblog
else
    echo "Pod goblog does not exist."
fi
