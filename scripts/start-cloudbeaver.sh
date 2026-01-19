#!/bin/bash

# CloudBeaver起動スクリプト
# 使用方法: ./scripts/start-cloudbeaver.sh {start|stop}

SCRIPT_DIR=$(cd "$(dirname "$0")/.." && pwd)
COMPOSE_FILE="$SCRIPT_DIR/docker-compose.cloudbeaver.yml"

case "$1" in
  start)
    echo "Starting CloudBeaver..."
    docker-compose -f "$COMPOSE_FILE" up -d
    echo "CloudBeaver started. Port: 8978"
    ;;
  stop)
    echo "Stopping CloudBeaver..."
    docker-compose -f "$COMPOSE_FILE" down
    echo "CloudBeaver stopped."
    ;;
  *)
    echo "Usage: $0 {start|stop}"
    exit 1
    ;;
esac
