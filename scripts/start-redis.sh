#!/bin/bash

# Redis起動スクリプト
# 使用方法: ./scripts/start-redis.sh {start|stop}

SCRIPT_DIR=$(cd "$(dirname "$0")/.." && pwd)
COMPOSE_FILE="$SCRIPT_DIR/docker-compose.redis.yml"

case "$1" in
  start)
    echo "Starting Redis..."
    docker-compose -f "$COMPOSE_FILE" up -d
    echo "Redis started. Port: 6379"
    ;;
  stop)
    echo "Stopping Redis..."
    docker-compose -f "$COMPOSE_FILE" down
    echo "Redis stopped."
    ;;
  *)
    echo "Usage: $0 {start|stop}"
    exit 1
    ;;
esac
