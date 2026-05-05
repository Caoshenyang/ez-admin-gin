#!/bin/bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
SERVER_DIR="$ROOT_DIR/server"

POSTGRES_IMAGE="${EZ_REAL_DB_POSTGRES_IMAGE:-postgres:18-alpine}"
MYSQL_IMAGE="${EZ_REAL_DB_MYSQL_IMAGE:-mysql:8.4}"

POSTGRES_CONTAINER="ez-admin-migration-postgres"
MYSQL_CONTAINER="ez-admin-migration-mysql"

POSTGRES_PORT="${EZ_REAL_DB_POSTGRES_PORT:-55432}"
MYSQL_PORT="${EZ_REAL_DB_MYSQL_PORT:-53306}"

POSTGRES_USER="${EZ_REAL_DB_POSTGRES_USER:-ez_admin}"
POSTGRES_PASSWORD="${EZ_REAL_DB_POSTGRES_PASSWORD:-ez_admin_123456}"
POSTGRES_DB="${EZ_REAL_DB_POSTGRES_NAME:-ez_admin_migration}"

MYSQL_USER="${EZ_REAL_DB_MYSQL_USER:-ez_admin}"
MYSQL_PASSWORD="${EZ_REAL_DB_MYSQL_PASSWORD:-ez_admin_123456}"
MYSQL_DB="${EZ_REAL_DB_MYSQL_NAME:-ez_admin_migration}"
MYSQL_ROOT_PASSWORD="${EZ_REAL_DB_MYSQL_ROOT_PASSWORD:-root_123456}"

cleanup() {
  docker rm -f "$POSTGRES_CONTAINER" >/dev/null 2>&1 || true
  docker rm -f "$MYSQL_CONTAINER" >/dev/null 2>&1 || true
}

wait_for_postgres() {
  for _ in $(seq 1 60); do
    if docker exec "$POSTGRES_CONTAINER" pg_isready -U "$POSTGRES_USER" -d "$POSTGRES_DB" >/dev/null 2>&1; then
      return 0
    fi
    sleep 1
  done
  echo "PostgreSQL did not become ready in time" >&2
  return 1
}

wait_for_mysql() {
  for _ in $(seq 1 90); do
    if docker exec "$MYSQL_CONTAINER" mysqladmin ping -h 127.0.0.1 -uroot "-p$MYSQL_ROOT_PASSWORD" --silent >/dev/null 2>&1; then
      return 0
    fi
    sleep 1
  done
  echo "MySQL did not become ready in time" >&2
  return 1
}

trap cleanup EXIT

echo ">>> Starting disposable PostgreSQL container..."
docker rm -f "$POSTGRES_CONTAINER" >/dev/null 2>&1 || true
docker run -d \
  --name "$POSTGRES_CONTAINER" \
  -e POSTGRES_USER="$POSTGRES_USER" \
  -e POSTGRES_PASSWORD="$POSTGRES_PASSWORD" \
  -e POSTGRES_DB="$POSTGRES_DB" \
  -p "127.0.0.1:${POSTGRES_PORT}:5432" \
  "$POSTGRES_IMAGE" >/dev/null
wait_for_postgres

echo ">>> Starting disposable MySQL container..."
docker rm -f "$MYSQL_CONTAINER" >/dev/null 2>&1 || true
docker run -d \
  --name "$MYSQL_CONTAINER" \
  -e MYSQL_ROOT_PASSWORD="$MYSQL_ROOT_PASSWORD" \
  -e MYSQL_DATABASE="$MYSQL_DB" \
  -e MYSQL_USER="$MYSQL_USER" \
  -e MYSQL_PASSWORD="$MYSQL_PASSWORD" \
  -p "127.0.0.1:${MYSQL_PORT}:3306" \
  "$MYSQL_IMAGE" >/dev/null
wait_for_mysql

echo ">>> Running embedded migration smoke tests on PostgreSQL and MySQL..."
cd "$SERVER_DIR"
EZ_REAL_DB_MIGRATION=1 \
EZ_REAL_DB_DRIVERS=postgres,mysql \
EZ_REAL_DB_POSTGRES_HOST=127.0.0.1 \
EZ_REAL_DB_POSTGRES_PORT="$POSTGRES_PORT" \
EZ_REAL_DB_POSTGRES_USER="$POSTGRES_USER" \
EZ_REAL_DB_POSTGRES_PASSWORD="$POSTGRES_PASSWORD" \
EZ_REAL_DB_POSTGRES_NAME="$POSTGRES_DB" \
EZ_REAL_DB_MYSQL_HOST=127.0.0.1 \
EZ_REAL_DB_MYSQL_PORT="$MYSQL_PORT" \
EZ_REAL_DB_MYSQL_USER="$MYSQL_USER" \
EZ_REAL_DB_MYSQL_PASSWORD="$MYSQL_PASSWORD" \
EZ_REAL_DB_MYSQL_NAME="$MYSQL_DB" \
go test -run TestEmbeddedMigrationsApplyOnRealDatabases .

echo ">>> Real database migration smoke tests passed."
