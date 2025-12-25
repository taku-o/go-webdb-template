#!/bin/bash
# マイグレーション適用スクリプト
# 使用方法: ./scripts/migrate.sh [master|sharding|all]

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
DB_DIR="$PROJECT_ROOT/db/migrations"
DATA_DIR="$PROJECT_ROOT/server/data"

# データディレクトリの作成
mkdir -p "$DATA_DIR"

# マスターグループのマイグレーション
migrate_master() {
    echo "Migrating master group..."
    local master_db="$DATA_DIR/master.db"

    for sql_file in "$DB_DIR/master"/*.sql; do
        if [ -f "$sql_file" ]; then
            echo "  Applying: $(basename "$sql_file")"
            sqlite3 "$master_db" < "$sql_file"
        fi
    done

    echo "Master group migration completed."
}

# シャーディンググループのマイグレーション
migrate_sharding() {
    echo "Migrating sharding group..."

    # まずテンプレートからSQLを生成
    echo "  Generating SQL from templates..."
    cd "$PROJECT_ROOT/server/cmd/migrate-gen"
    go run main.go
    cd "$PROJECT_ROOT"

    # 各データベースにマイグレーションを適用
    local db_mapping=(
        "1:0:7"   # DB1: テーブル 000-007
        "2:8:15"  # DB2: テーブル 008-015
        "3:16:23" # DB3: テーブル 016-023
        "4:24:31" # DB4: テーブル 024-031
    )

    for mapping in "${db_mapping[@]}"; do
        IFS=':' read -r db_id start end <<< "$mapping"
        local sharding_db="$DATA_DIR/sharding_db_${db_id}.db"

        echo "  Migrating sharding_db_${db_id} (tables ${start}-${end})..."

        # テーブル定義を適用
        for table_type in users posts; do
            for i in $(seq $start $end); do
                suffix=$(printf "%03d" $i)
                sql_file="$DB_DIR/sharding/generated/${table_type}_${suffix}.sql"

                if [ -f "$sql_file" ]; then
                    sqlite3 "$sharding_db" < "$sql_file"
                fi
            done
        done
    done

    echo "Sharding group migration completed."
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

echo "All migrations completed successfully!"
