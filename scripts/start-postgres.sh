#!/bin/bash

# PostgreSQL起動スクリプト
# 使用方法: ./scripts/start-postgres.sh {start|stop}

SCRIPT_DIR=$(cd "$(dirname "$0")/.." && pwd)
COMPOSE_FILE="$SCRIPT_DIR/docker-compose.postgres.yml"

case "$1" in
  start)
    echo "Starting PostgreSQL..."
    docker-compose -f "$COMPOSE_FILE" up -d
    echo "PostgreSQL started. Port: 5432"
    echo "Connection: postgresql://webdb:webdb@localhost:5432/webdb"
    ;;
  stop)
    echo "Stopping PostgreSQL..."
    docker-compose -f "$COMPOSE_FILE" down
    echo "PostgreSQL stopped."
    ;;
  *)
    echo "Usage: $0 {start|stop}"
    exit 1
    ;;
esac
