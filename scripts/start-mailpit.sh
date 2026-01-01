#!/bin/bash

# Mailpit起動スクリプト
# 使用方法: ./scripts/start-mailpit.sh {start|stop}

SCRIPT_DIR=$(cd "$(dirname "$0")/.." && pwd)
COMPOSE_FILE="$SCRIPT_DIR/docker-compose.mailpit.yml"

case "$1" in
  start)
    echo "Starting Mailpit..."
    docker-compose -f "$COMPOSE_FILE" up -d
    echo "Mailpit started. Web UI: http://localhost:8025"
    ;;
  stop)
    echo "Stopping Mailpit..."
    docker-compose -f "$COMPOSE_FILE" down
    echo "Mailpit stopped."
    ;;
  *)
    echo "Usage: $0 {start|stop}"
    exit 1
    ;;
esac
