#!/bin/bash
# マイグレーション適用スクリプト (Atlas版)
# 使用方法: ./scripts/migrate.sh [master|sharding|all]

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
DATA_DIR="$PROJECT_ROOT/server/data"

# データディレクトリの作成
mkdir -p "$DATA_DIR"

# マスターグループのマイグレーション
migrate_master() {
    echo "Migrating master group..."
    local master_db="$DATA_DIR/master.db"

    atlas migrate apply \
        --dir "file://$PROJECT_ROOT/db/migrations/master" \
        --url "sqlite://$master_db"

    echo "Master group migration applied."
}

# シャーディンググループのマイグレーション
migrate_sharding() {
    echo "Migrating sharding group..."

    # 各シャーディングDBにマイグレーションを適用
    for db_id in 1 2 3 4; do
        local sharding_db="$DATA_DIR/sharding_db_${db_id}.db"
        echo "  Migrating sharding_db_${db_id}..."

        atlas migrate apply \
            --dir "file://$PROJECT_ROOT/db/migrations/sharding" \
            --url "sqlite://$sharding_db"
    done

    echo "Sharding group migration applied."
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
    *)
        echo "Usage: $0 [master|sharding|all]"
        exit 1
        ;;
esac

echo "All migrations applied successfully!"
