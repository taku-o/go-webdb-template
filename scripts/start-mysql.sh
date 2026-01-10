#!/bin/bash
set -e

SCRIPT_DIR=$(cd $(dirname $0); pwd)
PROJECT_DIR=$(cd $SCRIPT_DIR/..; pwd)
COMPOSE_FILE="$PROJECT_DIR/docker-compose.mysql.yml"

usage() {
    echo "Usage: $0 {start|stop|status|health}"
    echo ""
    echo "Commands:"
    echo "  start   Start MySQL containers"
    echo "  stop    Stop MySQL containers"
    echo "  status  Show container status"
    echo "  health  Show health check status"
    exit 1
}

start() {
    echo "Starting MySQL containers..."
    docker-compose -f "$COMPOSE_FILE" up -d
    echo "MySQL containers started successfully."
    echo ""
    echo "Connection URLs:"
    echo "  Master:     mysql://webdb:webdb@tcp(localhost:3306)/webdb_master"
    echo "  Sharding 1: mysql://webdb:webdb@tcp(localhost:3307)/webdb_sharding_1"
    echo "  Sharding 2: mysql://webdb:webdb@tcp(localhost:3308)/webdb_sharding_2"
    echo "  Sharding 3: mysql://webdb:webdb@tcp(localhost:3309)/webdb_sharding_3"
    echo "  Sharding 4: mysql://webdb:webdb@tcp(localhost:3310)/webdb_sharding_4"
}

stop() {
    echo "Stopping MySQL containers..."
    docker-compose -f "$COMPOSE_FILE" down
    echo "MySQL containers stopped successfully."
}

status() {
    echo "MySQL container status:"
    docker-compose -f "$COMPOSE_FILE" ps
}

health() {
    echo "MySQL health check status:"
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
