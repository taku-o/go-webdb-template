#!/bin/bash

# Redis Insight起動スクリプト
# 使用方法: ./scripts/start-redis-insight.sh {start|stop}
# docker-compose.redis.ymlで起動した1台のRedisサーバーと接続

SCRIPT_DIR=$(cd "$(dirname "$0")/.." && pwd)
COMPOSE_FILE="$SCRIPT_DIR/docker-compose.redis-insight.yml"

case "$1" in
  start)
    echo "Starting Redis Insight..."
    docker-compose -f "$COMPOSE_FILE" up -d
    echo "Redis Insight started. Web UI: http://localhost:8001"
    ;;
  stop)
    echo "Stopping Redis Insight..."
    docker-compose -f "$COMPOSE_FILE" down
    echo "Redis Insight stopped."
    ;;
  *)
    echo "Usage: $0 {start|stop}"
    exit 1
    ;;
esac
