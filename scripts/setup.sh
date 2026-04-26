#!/usr/bin/env bash
set -euo pipefail

echo "Starting local infrastructure..."
docker compose up -d
echo "PROMETHEUS local dependencies are up."
