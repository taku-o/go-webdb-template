#!/bin/bash
set -e

SCRIPT_DIR=$(cd $(dirname $0); pwd)
PROJECT_DIR=$(cd $SCRIPT_DIR/..; pwd)
COMPOSE_FILE="$PROJECT_DIR/docker-compose.postgres.yml"

usage() {
    echo "Usage: $0 {start|stop|status|health}"
    echo ""
    echo "Commands:"
    echo "  start   Start PostgreSQL containers"
    echo "  stop    Stop PostgreSQL containers"
    echo "  status  Show container status"
    echo "  health  Show health check status"
    exit 1
}

start() {
    echo "Starting PostgreSQL containers..."
    docker-compose -f "$COMPOSE_FILE" up -d
    echo "PostgreSQL containers started successfully."
    echo ""
    echo "Connection URLs:"
    echo "  Master:     postgresql://webdb:webdb@localhost:5432/webdb_master"
    echo "  Sharding 1: postgresql://webdb:webdb@localhost:5433/webdb_sharding_1"
    echo "  Sharding 2: postgresql://webdb:webdb@localhost:5434/webdb_sharding_2"
    echo "  Sharding 3: postgresql://webdb:webdb@localhost:5435/webdb_sharding_3"
    echo "  Sharding 4: postgresql://webdb:webdb@localhost:5436/webdb_sharding_4"
}

stop() {
    echo "Stopping PostgreSQL containers..."
    docker-compose -f "$COMPOSE_FILE" down
    echo "PostgreSQL containers stopped successfully."
}

status() {
    echo "PostgreSQL container status:"
    docker-compose -f "$COMPOSE_FILE" ps
}

health() {
    echo "PostgreSQL health check status:"
    docker-compose -f "$COMPOSE_FILE" ps --format "table {{.Name}}\t{{.Status}}"
}

if [ $# -eq 0 ]; then
    usage
fi

case "$1" in
    start)
        start
        ;;
    stop)
        stop
        ;;
    status)
        status
        ;;
    health)
        health
        ;;
    *)
        usage
        ;;
esac
