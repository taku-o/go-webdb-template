#!/bin/bash

# MySQL起動スクリプト
# 使用方法: ./scripts/start-mysql.sh {start|stop}

SCRIPT_DIR=$(cd "$(dirname "$0")/.." && pwd)
COMPOSE_FILE="$SCRIPT_DIR/docker-compose.mysql.yml"

case "$1" in
  start)
    echo "Starting MySQL..."
    docker-compose -f "$COMPOSE_FILE" up -d
    echo "MySQL started."
    echo "Master: mysql://webdb:webdb@tcp(localhost:3306)/master_db"
    echo "Sharding 1: mysql://webdb:webdb@tcp(localhost:3307)/sharding_db_1"
    echo "Sharding 2: mysql://webdb:webdb@tcp(localhost:3308)/sharding_db_2"
    echo "Sharding 3: mysql://webdb:webdb@tcp(localhost:3309)/sharding_db_3"
    echo "Sharding 4: mysql://webdb:webdb@tcp(localhost:3310)/sharding_db_4"
    ;;
  stop)
    echo "Stopping MySQL..."
    docker-compose -f "$COMPOSE_FILE" down
    echo "MySQL stopped."
    ;;
  *)
    echo "Usage: $0 {start|stop}"
    exit 1
    ;;
esac
