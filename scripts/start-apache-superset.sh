#!/bin/bash
# Apache Superset起動スクリプト

# Docker ComposeでApache Supersetを起動
docker-compose -f docker-compose.apache-superset.yml up -d

echo "Waiting for container to be ready..."
sleep 5

# PostgreSQLドライバをインストール（leanイメージにはドライバが含まれていないため）
echo "Installing PostgreSQL driver..."
docker exec -u root apache-superset pip install --target /app/.venv/lib/python3.10/site-packages psycopg2-binary --quiet 2>/dev/null || true

# 初期化が必要かどうかをチェック（superset.dbが存在しない場合は初期化が必要）
if [ ! -f "apache-superset/data/superset.db" ]; then
    echo "Initializing Apache Superset..."

    # データベースマイグレーションを実行
    echo "Running database migrations..."
    docker exec apache-superset superset db upgrade

    # 管理者ユーザーを作成（admin/admin）
    echo "Creating admin user..."
    docker exec apache-superset superset fab create-admin \
        --username admin \
        --firstname Admin \
        --lastname User \
        --email admin@example.com \
        --password admin

    # デフォルトのロールと権限を設定
    echo "Initializing roles and permissions..."
    docker exec apache-superset superset init

    echo "Initialization completed."
fi

# 起動確認メッセージ
echo "Apache Superset started."
echo "Access URL: http://localhost:8088"
echo "Default credentials: admin/admin"
