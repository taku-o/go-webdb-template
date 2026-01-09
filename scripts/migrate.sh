#!/bin/bash
# PostgreSQL用マイグレーションスクリプト
# 使用方法: APP_ENV=develop ./scripts/migrate.sh [master|sharding|all]

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# 環境変数からAPP_ENVを取得（デフォルト: develop）
APP_ENV="${APP_ENV:-develop}"
CONFIG_FILE="$PROJECT_ROOT/config/$APP_ENV/database.yaml"

# デフォルトのPostgreSQL接続情報
# config/{env}/database.yamlのPostgreSQL設定が有効になった場合はそちらを使用
MASTER_HOST="localhost"
MASTER_PORT="5432"
MASTER_USER="webdb"
MASTER_PASSWORD="webdb"
MASTER_DB="webdb_master"

SHARDING_1_HOST="localhost"
SHARDING_1_PORT="5433"
SHARDING_1_USER="webdb"
SHARDING_1_PASSWORD="webdb"
SHARDING_1_DB="webdb_sharding_1"

SHARDING_2_HOST="localhost"
SHARDING_2_PORT="5434"
SHARDING_2_USER="webdb"
SHARDING_2_PASSWORD="webdb"
SHARDING_2_DB="webdb_sharding_2"

SHARDING_3_HOST="localhost"
SHARDING_3_PORT="5435"
SHARDING_3_USER="webdb"
SHARDING_3_PASSWORD="webdb"
SHARDING_3_DB="webdb_sharding_3"

SHARDING_4_HOST="localhost"
SHARDING_4_PORT="5436"
SHARDING_4_USER="webdb"
SHARDING_4_PASSWORD="webdb"
SHARDING_4_DB="webdb_sharding_4"

# PostgreSQL URL形式を構築
build_postgres_url() {
    local host=$1
    local port=$2
    local user=$3
    local password=$4
    local dbname=$5
    echo "postgres://${user}:${password}@${host}:${port}/${dbname}?sslmode=disable"
}

# 使用方法を表示
usage() {
    echo "Usage: APP_ENV=develop $0 [master|sharding|all]"
    echo ""
    echo "Commands:"
    echo "  master    Apply migrations to master database only"
    echo "  sharding  Apply migrations to sharding databases only"
    echo "  all       Apply migrations to all databases (default)"
    echo ""
    echo "Environment variables:"
    echo "  APP_ENV   Environment name (develop/staging/production, default: develop)"
    exit 1
}

# マスターグループのマイグレーション
migrate_master() {
    echo "Migrating master database..."
    local url=$(build_postgres_url "$MASTER_HOST" "$MASTER_PORT" "$MASTER_USER" "$MASTER_PASSWORD" "$MASTER_DB")

    # Atlasマイグレーション適用
    echo "  Applying Atlas migrations..."
    atlas migrate apply \
        --dir "file://$PROJECT_ROOT/db/migrations/master" \
        --url "$url"

    # Viewマイグレーション適用（生SQL）
    echo "  Applying View migrations..."
    local view_dir="$PROJECT_ROOT/db/migrations/view_master"
    if [ -d "$view_dir" ]; then
        for sql_file in $(ls "$view_dir"/*.sql 2>/dev/null | sort); do
            echo "    Applying $(basename "$sql_file")..."
            docker exec -i postgres-master psql -U "$MASTER_USER" -d "$MASTER_DB" < "$sql_file"
        done
    fi

    echo "Master database migration applied."
}

# シャーディンググループのマイグレーション
migrate_sharding() {
    echo "Migrating sharding databases..."

    # Sharding 1
    echo "  Migrating sharding_1..."
    local url1=$(build_postgres_url "$SHARDING_1_HOST" "$SHARDING_1_PORT" "$SHARDING_1_USER" "$SHARDING_1_PASSWORD" "$SHARDING_1_DB")
    atlas migrate apply \
        --dir "file://$PROJECT_ROOT/db/migrations/sharding_1" \
        --url "$url1"

    # Sharding 2
    echo "  Migrating sharding_2..."
    local url2=$(build_postgres_url "$SHARDING_2_HOST" "$SHARDING_2_PORT" "$SHARDING_2_USER" "$SHARDING_2_PASSWORD" "$SHARDING_2_DB")
    atlas migrate apply \
        --dir "file://$PROJECT_ROOT/db/migrations/sharding_2" \
        --url "$url2"

    # Sharding 3
    echo "  Migrating sharding_3..."
    local url3=$(build_postgres_url "$SHARDING_3_HOST" "$SHARDING_3_PORT" "$SHARDING_3_USER" "$SHARDING_3_PASSWORD" "$SHARDING_3_DB")
    atlas migrate apply \
        --dir "file://$PROJECT_ROOT/db/migrations/sharding_3" \
        --url "$url3"

    # Sharding 4
    echo "  Migrating sharding_4..."
    local url4=$(build_postgres_url "$SHARDING_4_HOST" "$SHARDING_4_PORT" "$SHARDING_4_USER" "$SHARDING_4_PASSWORD" "$SHARDING_4_DB")
    atlas migrate apply \
        --dir "file://$PROJECT_ROOT/db/migrations/sharding_4" \
        --url "$url4"

    echo "Sharding databases migration applied."
}

# メイン処理
case "${1:-all}" in
    master)
        migrate_master
        ;;
    sharding)
        migrate_sharding
        ;;
    all)
        migrate_master
        migrate_sharding
        ;;
    -h|--help)
        usage
        ;;
    *)
        usage
        ;;
esac

echo "All migrations applied successfully!"
