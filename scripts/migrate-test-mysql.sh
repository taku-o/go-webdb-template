#!/bin/bash
# MySQL用マイグレーションスクリプト（テスト環境用）
# 使用方法: ./scripts/migrate-test-mysql.sh [master|sharding|all]
# テスト用データベースは事前に作成しておく必要がある

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# テスト環境を固定で設定
APP_ENV="test"
CONFIG_FILE="$PROJECT_ROOT/config/$APP_ENV/database.mysql.yaml"

# テスト環境用MySQL接続情報
MASTER_HOST="localhost"
MASTER_PORT="3306"
MASTER_USER="webdb"
MASTER_PASSWORD="webdb"
MASTER_DB="webdb_master_test"

SHARDING_1_HOST="localhost"
SHARDING_1_PORT="3307"
SHARDING_1_USER="webdb"
SHARDING_1_PASSWORD="webdb"
SHARDING_1_DB="webdb_sharding_1_test"

SHARDING_2_HOST="localhost"
SHARDING_2_PORT="3308"
SHARDING_2_USER="webdb"
SHARDING_2_PASSWORD="webdb"
SHARDING_2_DB="webdb_sharding_2_test"

SHARDING_3_HOST="localhost"
SHARDING_3_PORT="3309"
SHARDING_3_USER="webdb"
SHARDING_3_PASSWORD="webdb"
SHARDING_3_DB="webdb_sharding_3_test"

SHARDING_4_HOST="localhost"
SHARDING_4_PORT="3310"
SHARDING_4_USER="webdb"
SHARDING_4_PASSWORD="webdb"
SHARDING_4_DB="webdb_sharding_4_test"

# MySQL URL形式を構築
build_mysql_url() {
    local host=$1
    local port=$2
    local user=$3
    local password=$4
    local dbname=$5
    echo "mysql://${user}:${password}@${host}:${port}/${dbname}"
}

# 使用方法を表示
usage() {
    echo "Usage: $0 [master|sharding|all]"
    echo ""
    echo "Commands:"
    echo "  master    Apply migrations to master database only"
    echo "  sharding  Apply migrations to sharding databases only"
    echo "  all       Apply migrations to all databases (default)"
    echo ""
    echo "This script uses APP_ENV=test and config/test/atlas.mysql.hcl"
    exit 1
}

# マスターグループのマイグレーション
migrate_master() {
    echo "Migrating master database (test environment - MySQL)..."
    local url=$(build_mysql_url "$MASTER_HOST" "$MASTER_PORT" "$MASTER_USER" "$MASTER_PASSWORD" "$MASTER_DB")

    # Atlasマイグレーション適用
    echo "  Applying Atlas migrations..."
    atlas migrate apply \
        --dir "file://$PROJECT_ROOT/db/migrations/master-mysql" \
        --url "$url"

    # Viewマイグレーション適用（生SQL）
    echo "  Applying View migrations..."
    local view_dir="$PROJECT_ROOT/db/migrations/view_master-mysql"
    if [ -d "$view_dir" ]; then
        for sql_file in $(ls "$view_dir"/*.sql 2>/dev/null | sort); do
            echo "    Applying $(basename "$sql_file")..."
            docker exec -i mysql-master mysql -u"$MASTER_USER" -p"$MASTER_PASSWORD" "$MASTER_DB" < "$sql_file"
        done
    fi

    echo "Master database migration applied."
}

# シャーディンググループのマイグレーション
migrate_sharding() {
    echo "Migrating sharding databases (test environment - MySQL)..."

    # Sharding 1
    echo "  Migrating sharding_1..."
    local url1=$(build_mysql_url "$SHARDING_1_HOST" "$SHARDING_1_PORT" "$SHARDING_1_USER" "$SHARDING_1_PASSWORD" "$SHARDING_1_DB")
    atlas migrate apply \
        --dir "file://$PROJECT_ROOT/db/migrations/sharding_1-mysql" \
        --url "$url1"

    # Sharding 2
    echo "  Migrating sharding_2..."
    local url2=$(build_mysql_url "$SHARDING_2_HOST" "$SHARDING_2_PORT" "$SHARDING_2_USER" "$SHARDING_2_PASSWORD" "$SHARDING_2_DB")
    atlas migrate apply \
        --dir "file://$PROJECT_ROOT/db/migrations/sharding_2-mysql" \
        --url "$url2"

    # Sharding 3
    echo "  Migrating sharding_3..."
    local url3=$(build_mysql_url "$SHARDING_3_HOST" "$SHARDING_3_PORT" "$SHARDING_3_USER" "$SHARDING_3_PASSWORD" "$SHARDING_3_DB")
    atlas migrate apply \
        --dir "file://$PROJECT_ROOT/db/migrations/sharding_3-mysql" \
        --url "$url3"

    # Sharding 4
    echo "  Migrating sharding_4..."
    local url4=$(build_mysql_url "$SHARDING_4_HOST" "$SHARDING_4_PORT" "$SHARDING_4_USER" "$SHARDING_4_PASSWORD" "$SHARDING_4_DB")
    atlas migrate apply \
        --dir "file://$PROJECT_ROOT/db/migrations/sharding_4-mysql" \
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
