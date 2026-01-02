#!/bin/bash

# Redis Cluster起動スクリプト
# 使用方法: ./scripts/start-redis-cluster.sh {start|stop}

SCRIPT_DIR=$(cd "$(dirname "$0")/.." && pwd)
COMPOSE_FILE="$SCRIPT_DIR/docker-compose.redis-cluster.yml"

case "$1" in
  start)
    echo "Starting Redis Cluster..."
    docker-compose -f "$COMPOSE_FILE" up -d
    echo "Redis Cluster started. Ports: 7100-7105"
    ;;
  stop)
    echo "Stopping Redis Cluster..."
    docker-compose -f "$COMPOSE_FILE" down
    echo "Redis Cluster stopped."
    ;;
  *)
    echo "Usage: $0 {start|stop}"
    exit 1
    ;;
esac
