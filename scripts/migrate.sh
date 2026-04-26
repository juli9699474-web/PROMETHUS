#!/usr/bin/env bash
set -euo pipefail

DB_URL="${DATABASE_URL:-postgres://prometheus:devpassword@localhost:5432/prometheus}"
MIGRATION_FILE="infra/sql/001_init.sql"

if command -v psql >/dev/null 2>&1; then
  echo "Applying migrations using local psql..."
  psql "$DB_URL" -f "$MIGRATION_FILE"
  echo "Migrations applied."
  exit 0
fi

echo "psql not found; trying postgres container..."
docker compose exec -T postgres psql -U prometheus -d prometheus < "$MIGRATION_FILE" || {
  echo "Failed to run migrations. Install psql or run manually."
  exit 1
}
