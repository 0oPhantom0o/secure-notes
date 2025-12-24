#!/usr/bin/env bash
set -euo pipefail

# Project DB initializer for secure-notes
# - Creates the database if missing
# - Creates required tables
# - Creates required global indexes (parent table indexes are kept minimal; child
#   indexes are created by the app's partition maintenance if you enabled it)

# ------------------------------
# Configuration
# ------------------------------
: "${DB_HOST:=localhost}"
: "${DB_PORT:=5432}"
: "${DB_USER:=postgres}"
: "${DB_PASSWORD:=postgres}"
: "${DB_NAME:=notepad}"

export PGPASSWORD="${DB_PASSWORD}"

# Resolve docker compose command if needed
get_compose_cmd() {
  if command -v docker >/dev/null 2>&1 && docker compose version >/dev/null 2>&1; then
    echo "docker compose"
  elif command -v docker-compose >/dev/null 2>&1; then
    echo "docker-compose"
  else
    echo ""
  fi
}

have_psql() { command -v psql >/dev/null 2>&1; }

run_psql_cmd() {
  local db="$1"; shift
  local sql="$1"; shift || true
  if have_psql; then
    psql -h "${DB_HOST}" -p "${DB_PORT}" -U "${DB_USER}" -d "$db" -v ON_ERROR_STOP=1 -tAc "$sql"
  else
    local dcmd
    dcmd=$(get_compose_cmd)
    if [[ -z "$dcmd" ]]; then
      echo "Error: psql is not installed and docker compose is unavailable." >&2
      return 127
    fi
    # Use psql inside the postgres container
    ${dcmd} exec -T -e PGPASSWORD="${DB_PASSWORD}" postgres \
      psql -h localhost -p 5432 -U "${DB_USER}" -d "$db" -v ON_ERROR_STOP=1 -tAc "$sql"
  fi
}

echo "Checking database '${DB_NAME}' on ${DB_HOST}:${DB_PORT} as ${DB_USER}..."
DB_EXISTS=$(run_psql_cmd postgres "SELECT 1 FROM pg_database WHERE datname='${DB_NAME}';" || echo "")
if [[ "$DB_EXISTS" != "1" ]]; then
  echo "Creating database '${DB_NAME}'..."
  if have_psql; then
    psql -h "${DB_HOST}" -p "${DB_PORT}" -U "${DB_USER}" -v ON_ERROR_STOP=1 -c "CREATE DATABASE ${DB_NAME};"
  else
    dcmd=$(get_compose_cmd)
    if [[ -z "$dcmd" ]]; then
      echo "Error: cannot create DB — psql missing and docker compose not found." >&2
      exit 1
    fi
    ${dcmd} exec -T -e PGPASSWORD="${DB_PASSWORD}" postgres \
      psql -h localhost -p 5432 -U "${DB_USER}" -v ON_ERROR_STOP=1 -c "CREATE DATABASE ${DB_NAME};"
  fi
else
  echo "Database '${DB_NAME}' already exists."
fi
