# JobQueueサーバーのDocker対応設計書

## 1. 概要

### 1.1 設計の目的
要件定義書に基づき、JobQueueサーバーをDocker上で動作させるための詳細設計を定義する。既存のAPIサーバーやAdminサーバーと同様のDocker環境を構築し、開発環境での一貫性を保つ。また、既存のdocker-composeファイルのコンテナ名を統一的な命名規則に従うように修正する。

### 1.2 設計の範囲
- `server/Dockerfile.jobqueue`の新規作成
- `docker-compose.jobqueue.yml`の新規作成
- `docker-compose.api.yml`のコンテナ名修正（`api-develop` → `api`）
- `docker-compose.admin.yml`のコンテナ名修正（`admin-develop` → `admin`）
- `docker-compose.client.yml`のコンテナ名修正（`client-develop` → `client`）
- ドキュメントの更新（`docs/ja/Docker.md`、`docs/en/Docker.md`、`README.md`、`README.ja.md`）

### 1.3 設計方針
- **既存パターンの遵守**: 既存の`server/Dockerfile`と`server/Dockerfile.admin`と同様の構造を維持
- **docker-compose設定の統一**: 既存の`docker-compose.api.yml`と`docker-compose.admin.yml`と同様の構造を維持
- **コンテナ名の統一**: 環境名（`-develop`）を含めない統一的な命名規則に従う
- **既存機能の維持**: 既存のJobQueueサーバーの機能を完全に維持
- **ドキュメントの整合性**: すべてのドキュメントで一貫した情報を提供

## 2. Dockerfile設計

### 2.1 server/Dockerfile.jobqueueの設計

#### 2.1.1 基本構造
既存の`server/Dockerfile`と`server/Dockerfile.admin`と同様のマルチステージビルド構造を採用する。

```
ビルドステージ（golang:1.24-alpine）
  ↓
実行ステージ（alpine:latest）
```

#### 2.1.2 ビルドステージの詳細

```dockerfile
# ビルドステージ
FROM golang:1.24-alpine AS builder

WORKDIR /build

# 依存関係をコピー
COPY go.mod go.sum ./
RUN go mod download

# ソースコードをコピー
COPY . .

# ビルド
RUN CGO_ENABLED=0 go build -o jobqueue ./cmd/jobqueue/main.go
```

**設計ポイント**:
- ベースイメージ: `golang:1.24-alpine`（既存のDockerfileと同一）
- CGO設定: `CGO_ENABLED=0`（既存のDockerfileと同一）
- ビルドコマンド: `./cmd/jobqueue/main.go`をビルド
- バイナリ名: `jobqueue`（既存の`server`や`admin`と同様の命名規則）

#### 2.1.3 実行ステージの詳細

```dockerfile
# 実行ステージ
FROM alpine:latest

# 必要なパッケージをインストール
RUN apk add --no-cache ca-certificates tzdata wget

# 非rootユーザーを作成（Alpine用コマンド）
RUN addgroup -g 1000 appuser && \
    adduser -u 1000 -G appuser -D appuser

# ディレクトリ構造を作成（server/から起動される想定）
# /app/server/ - バイナリ配置・作業ディレクトリ
# /app/config/ - 設定ファイル（../config/でアクセス）
# /app/logs/   - ログファイル（../logs/でアクセス）
RUN mkdir -p /app/server/data /app/config /app/logs && \
    chown -R appuser:appuser /app

WORKDIR /app/server

# ビルド成果物をコピー
COPY --from=builder /build/jobqueue .

USER appuser

EXPOSE 8082

CMD ["./jobqueue"]
```

**設計ポイント**:
- ベースイメージ: `alpine:latest`（既存のDockerfileと同一）
- パッケージ: `ca-certificates`、`tzdata`、`wget`（既存のDockerfileと同一、`wget`はヘルスチェック用）
- 非rootユーザー: `appuser`（UID 1000、GID 1000、既存のDockerfileと同一）
- ディレクトリ構造: 既存のDockerfileと同一
- ポート: 8082（JobQueueサーバーのポート）
- CMD: `["./jobqueue"]`（既存の`["./server"]`や`["./admin"]`と同様）

#### 2.1.4 既存Dockerfileとの比較

| 項目 | Dockerfile | Dockerfile.admin | Dockerfile.jobqueue |
|------|-----------|-----------------|-------------------|
| ビルドコマンド | `./cmd/server/main.go` | `./cmd/admin/main.go` | `./cmd/jobqueue/main.go` |
| バイナリ名 | `server` | `admin` | `jobqueue` |
| ポート | 8080 | 8081 | 8082 |
| CMD | `["./server"]` | `["./admin"]` | `["./jobqueue"]` |
| その他 | 同一 | 同一 | 同一 |

### 2.2 ディレクトリ構造

```
/app/
├── server/
│   ├── jobqueue          # バイナリ
│   └── data/             # データディレクトリ
├── config/
│   └── develop/          # 設定ファイル（マウント）
└── logs/                 # ログファイル（マウント）
```

**設計ポイント**:
- 既存のDockerfileと同一のディレクトリ構造を維持
- 設定ファイルとログファイルはボリュームマウントで提供

## 3. docker-compose設定設計

### 3.1 docker-compose.jobqueue.ymlの設計

#### 3.1.1 基本構造

```yaml
services:
  jobqueue:
    build:
      context: ./server
      dockerfile: Dockerfile.jobqueue
    container_name: jobqueue
    ports:
      - "8082:8082"
    environment:
      - APP_ENV=develop
      - REDIS_JOBQUEUE_ADDR=redis:6379
      # Database DSNs for Docker environment
      - DB_MASTER_WRITER_DSN=postgres://webdb:webdb@postgres-master:5432/webdb_master?sslmode=disable
      - DB_MASTER_READER_DSN=postgres://webdb:webdb@postgres-master:5432/webdb_master?sslmode=disable
      - DB_SHARD1_WRITER_DSN=postgres://webdb:webdb@postgres-sharding-1:5432/webdb_sharding_1?sslmode=disable
      - DB_SHARD1_READER_DSN=postgres://webdb:webdb@postgres-sharding-1:5432/webdb_sharding_1?sslmode=disable
      - DB_SHARD2_WRITER_DSN=postgres://webdb:webdb@postgres-sharding-1:5432/webdb_sharding_1?sslmode=disable
      - DB_SHARD2_READER_DSN=postgres://webdb:webdb@postgres-sharding-1:5432/webdb_sharding_1?sslmode=disable
      - DB_SHARD3_WRITER_DSN=postgres://webdb:webdb@postgres-sharding-2:5432/webdb_sharding_2?sslmode=disable
      - DB_SHARD3_READER_DSN=postgres://webdb:webdb@postgres-sharding-2:5432/webdb_sharding_2?sslmode=disable
      - DB_SHARD4_WRITER_DSN=postgres://webdb:webdb@postgres-sharding-2:5432/webdb_sharding_2?sslmode=disable
      - DB_SHARD4_READER_DSN=postgres://webdb:webdb@postgres-sharding-2:5432/webdb_sharding_2?sslmode=disable
      - DB_SHARD5_WRITER_DSN=postgres://webdb:webdb@postgres-sharding-3:5432/webdb_sharding_3?sslmode=disable
      - DB_SHARD5_READER_DSN=postgres://webdb:webdb@postgres-sharding-3:5432/webdb_sharding_3?sslmode=disable
      - DB_SHARD6_WRITER_DSN=postgres://webdb:webdb@postgres-sharding-3:5432/webdb_sharding_3?sslmode=disable
      - DB_SHARD6_READER_DSN=postgres://webdb:webdb@postgres-sharding-3:5432/webdb_sharding_3?sslmode=disable
      - DB_SHARD7_WRITER_DSN=postgres://webdb:webdb@postgres-sharding-4:5432/webdb_sharding_4?sslmode=disable
      - DB_SHARD7_READER_DSN=postgres://webdb:webdb@postgres-sharding-4:5432/webdb_sharding_4?sslmode=disable
      - DB_SHARD8_WRITER_DSN=postgres://webdb:webdb@postgres-sharding-4:5432/webdb_sharding_4?sslmode=disable
      - DB_SHARD8_READER_DSN=postgres://webdb:webdb@postgres-sharding-4:5432/webdb_sharding_4?sslmode=disable
    volumes:
      - ./config/develop:/app/config/develop:ro
      - ./server/data:/app/server/data
      - ./logs:/app/logs
    networks:
      - postgres-network
      - redis-network
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "wget", "--quiet", "--tries=1", "--spider", "http://localhost:8082/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s

networks:
  postgres-network:
    external: true
    name: postgres-network
  redis-network:
    external: true
    name: redis-network
```

#### 3.1.2 設計ポイント

**サービス名とコンテナ名**:
- サービス名: `jobqueue`
- コンテナ名: `jobqueue`（環境名を含めない）

**ビルド設定**:
- `context: ./server`: 既存のdocker-composeファイルと同一
- `dockerfile: Dockerfile.jobqueue`: 新規作成するDockerfile

**ポート設定**:
- `"8082:8082"`: JobQueueサーバーのポート

**環境変数**:
- `APP_ENV=develop`: 開発環境を指定
- `REDIS_JOBQUEUE_ADDR=redis:6379`: Redis接続先（Docker環境では`redis`ホスト名を使用）
- データベースDSN: 既存のAPI/Adminサーバーと同様の設定（現在のジョブ処理ではPostgreSQLを使用していないが、将来的にジョブ処理でPostgreSQLを使用する可能性が高いため、拡張性を考慮して設定を残す）

**ボリュームマウント**:
- `./config/develop:/app/config/develop:ro`: 設定ファイル（読み取り専用、既存と同一）
- `./server/data:/app/server/data`: データディレクトリ（読み書き可、既存と同一）
- `./logs:/app/logs`: ログファイル（読み書き可、既存と同一）

**ネットワーク**:
- `postgres-network`: PostgreSQL接続用（external: true、既存と同一）
- `redis-network`: Redis接続用（external: true、既存と同一）

**ヘルスチェック**:
- `test: ["CMD", "wget", "--quiet", "--tries=1", "--spider", "http://localhost:8082/health"]`: `/health`エンドポイントを確認
- 既存のAPI/Adminサーバーと同様の設定（ポートのみ変更）

#### 3.1.3 既存docker-composeファイルとの比較

| 項目 | docker-compose.api.yml | docker-compose.admin.yml | docker-compose.jobqueue.yml |
|------|----------------------|------------------------|---------------------------|
| サービス名 | `api` | `admin` | `jobqueue` |
| コンテナ名 | `api-develop` → `api` | `admin-develop` → `admin` | `jobqueue` |
| ポート | `8080:8080` | `8081:8081` | `8082:8082` |
| ネットワーク | `postgres-network`, `redis-network` | `postgres-network` | `postgres-network`, `redis-network` |
| ヘルスチェック | `/health` (8080) | `/health` (8081) | `/health` (8082) |
| その他 | 同一 | 同一 | 同一 |

### 3.2 既存docker-composeファイルのコンテナ名修正

#### 3.2.1 docker-compose.api.ymlの修正

**変更前**:
```yaml
container_name: api-develop
```

**変更後**:
```yaml
container_name: api
```

#### 3.2.2 docker-compose.admin.ymlの修正

**変更前**:
```yaml
container_name: admin-develop
```

**変更後**:
```yaml
container_name: admin
```

#### 3.2.3 docker-compose.client.ymlの修正

**変更前**:
```yaml
container_name: client-develop
```

**変更後**:
```yaml
container_name: client
```

#### 3.2.4 修正の理由
- 環境名（`-develop`）を含めない統一的な命名規則に従う
- コンテナ名を簡潔にし、管理を容易にする
- すべてのdocker-composeファイルで一貫した命名規則を適用

## 4. ドキュメント更新設計

### 4.1 docs/ja/Docker.mdの更新

#### 4.1.1 概要セクションの更新

**変更前**:
```markdown
APIサーバー、Adminサーバー、クライアントサーバーをDockerコンテナ上で動作させることができます。
```

**変更後**:
```markdown
APIサーバー、Adminサーバー、JobQueueサーバー、クライアントサーバーをDockerコンテナ上で動作させることができます。
```

#### 4.1.2 アーキテクチャ図の更新

**変更前**:
```
┌──────────────┐    ┌──────────────┐    ┌──────────────┐
│  Client      │    │  API Server  │    │ Admin Server │
│  (Port 3000) │───▶│  (Port 8080) │    │ (Port 8081)  │
└──────────────┘    └───────┬──────┘    └───────┬──────┘
                            │                    │
                    ┌───────┴────────────────────┘
                    ▼
            ┌──────────────┐    ┌──────────────┐
            │  PostgreSQL  │    │    Redis     │
            │  (postgres)  │    │   (redis)    │
            └──────────────┘    └──────────────┘
```

**変更後**:
```
┌──────────────┐    ┌──────────────┐    ┌──────────────┐    ┌──────────────┐
│  Client      │    │  API Server  │    │ Admin Server │    │ JobQueue     │
│  (Port 3000) │───▶│  (Port 8080) │    │ (Port 8081)  │    │ (Port 8082)  │
└──────────────┘    └───────┬──────┘    └───────┬──────┘    └───────┬──────┘
                            │                    │                    │
                    ┌───────┴────────────────────┴────────────────────┘
                    ▼
            ┌──────────────┐    ┌──────────────┐
            │  PostgreSQL  │    │    Redis     │
            │  (postgres)  │    │   (redis)    │
            └──────────────┘    └──────────────┘
```

#### 4.1.3 前提条件セクションの更新

**変更前**:
```markdown
- 8080（APIサーバー）
- 8081（Adminサーバー）
- 3000（クライアントサーバー）
```

**変更後**:
```markdown
- 8080（APIサーバー）
- 8081（Adminサーバー）
- 8082（JobQueueサーバー）
- 3000（クライアントサーバー）
```

#### 4.1.4 Dockerfile構成表の更新

**変更前**:
| ファイル | 用途 | CGO | ベースイメージ |
|---------|------|-----|---------------|
| `server/Dockerfile` | 全環境（develop/staging/production） | 0 | golang:1.24-alpine → alpine:latest |
| `server/Dockerfile.admin` | 全環境（develop/staging/production） | 0 | golang:1.24-alpine → alpine:latest |

**変更後**:
| ファイル | 用途 | CGO | ベースイメージ |
|---------|------|-----|---------------|
| `server/Dockerfile` | 全環境（develop/staging/production） | 0 | golang:1.24-alpine → alpine:latest |
| `server/Dockerfile.admin` | 全環境（develop/staging/production） | 0 | golang:1.24-alpine → alpine:latest |
| `server/Dockerfile.jobqueue` | 全環境（develop/staging/production） | 0 | golang:1.24-alpine → alpine:latest |

#### 4.1.5 Docker Compose設定ファイル一覧の更新

**変更前**:
| ファイル | サービス | 環境 |
|---------|---------|------|
| `docker-compose.api.yml` | APIサーバー | develop |
| `docker-compose.admin.yml` | Adminサーバー | develop |
| `docker-compose.client.yml` | クライアント | develop |

**変更後**:
| ファイル | サービス | 環境 |
|---------|---------|------|
| `docker-compose.api.yml` | APIサーバー | develop |
| `docker-compose.admin.yml` | Adminサーバー | develop |
| `docker-compose.jobqueue.yml` | JobQueueサーバー | develop |
| `docker-compose.client.yml` | クライアント | develop |

#### 4.1.6 ビルド・起動・停止コマンドの更新

**変更前**:
```bash
# ビルド
docker-compose -f docker-compose.api.yml build
docker-compose -f docker-compose.admin.yml build
docker-compose -f docker-compose.client.yml build

# 起動
docker-compose -f docker-compose.api.yml up -d
docker-compose -f docker-compose.admin.yml up -d
docker-compose -f docker-compose.client.yml up -d

# 停止
docker-compose -f docker-compose.api.yml down
docker-compose -f docker-compose.admin.yml down
docker-compose -f docker-compose.client.yml down

# ログ確認
docker-compose -f docker-compose.api.yml logs -f
docker-compose -f docker-compose.admin.yml logs -f
docker-compose -f docker-compose.client.yml logs -f
```

**変更後**:
```bash
# ビルド
docker-compose -f docker-compose.api.yml build
docker-compose -f docker-compose.admin.yml build
docker-compose -f docker-compose.jobqueue.yml build
docker-compose -f docker-compose.client.yml build

# 起動
docker-compose -f docker-compose.api.yml up -d
docker-compose -f docker-compose.admin.yml up -d
docker-compose -f docker-compose.jobqueue.yml up -d
docker-compose -f docker-compose.client.yml up -d

# 停止
docker-compose -f docker-compose.api.yml down
docker-compose -f docker-compose.admin.yml down
docker-compose -f docker-compose.jobqueue.yml down
docker-compose -f docker-compose.client.yml down

# ログ確認
docker-compose -f docker-compose.api.yml logs -f
docker-compose -f docker-compose.admin.yml logs -f
docker-compose -f docker-compose.jobqueue.yml logs -f
docker-compose -f docker-compose.client.yml logs -f
```

#### 4.1.7 起動順序の更新

**変更前**:
1. PostgreSQL、Redisコンテナを起動
2. APIサーバーを起動
3. Adminサーバーを起動
4. クライアントサーバーを起動

**変更後**:
1. PostgreSQL、Redisコンテナを起動
2. APIサーバーを起動
3. Adminサーバーを起動
4. JobQueueサーバーを起動
5. クライアントサーバーを起動

#### 4.1.8 サービス間通信の更新

**変更前**:
| 接続元 | 接続先 | ホスト名 |
|--------|--------|---------|
| APIサーバー | PostgreSQL | postgres:5432 |
| APIサーバー | Redis | redis:6379 |
| Adminサーバー | PostgreSQL | postgres:5432 |
| クライアント | APIサーバー | api:8080 |

**変更後**:
| 接続元 | 接続先 | ホスト名 |
|--------|--------|---------|
| APIサーバー | PostgreSQL | postgres:5432 |
| APIサーバー | Redis | redis:6379 |
| Adminサーバー | PostgreSQL | postgres:5432 |
| JobQueueサーバー | Redis | redis:6379 |
| クライアント | APIサーバー | api:8080 |

### 4.2 docs/en/Docker.mdの更新

`docs/ja/Docker.md`と同様の更新を英語版にも適用する。

### 4.3 README.mdの更新

#### 4.3.1 Docker環境の説明セクションの追加

JobQueueサーバーのDocker起動方法を追加する。

```markdown
#### JobQueue Server

```bash
# ビルド
docker-compose -f docker-compose.jobqueue.yml build

# 起動
docker-compose -f docker-compose.jobqueue.yml up -d

# 停止
docker-compose -f docker-compose.jobqueue.yml down

# ログ確認
docker-compose -f docker-compose.jobqueue.yml logs -f
```

The JobQueue server runs an HTTP server (port 8082) and an Asynq server in parallel. The HTTP server provides a `/health` endpoint for health monitoring. The Asynq server processes jobs from Redis. Make sure Redis is running before starting the JobQueue server.
```

### 4.4 README.ja.mdの更新

`README.md`と同様の更新を日本語版にも適用する。

## 5. 実装手順

### 5.1 Dockerfile.jobqueueの作成

1. `server/Dockerfile`をコピーして`server/Dockerfile.jobqueue`を作成
2. ビルドコマンドを`./cmd/jobqueue/main.go`に変更
3. バイナリ名を`jobqueue`に変更
4. ポートを`8082`に変更
5. CMDを`["./jobqueue"]`に変更

### 5.2 docker-compose.jobqueue.ymlの作成

1. `docker-compose.api.yml`をコピーして`docker-compose.jobqueue.yml`を作成
2. サービス名を`jobqueue`に変更
3. コンテナ名を`jobqueue`に変更
4. ビルド設定の`dockerfile`を`Dockerfile.jobqueue`に変更
5. ポートを`8082:8082`に変更
6. ヘルスチェックのポートを`8082`に変更
7. 環境変数は既存のAPI/Adminサーバーと同様に設定（データベースDSNも含む）

### 5.3 既存docker-composeファイルの修正

1. `docker-compose.api.yml`の`container_name`を`api-develop`から`api`に変更
2. `docker-compose.admin.yml`の`container_name`を`admin-develop`から`admin`に変更
3. `docker-compose.client.yml`の`container_name`を`client-develop`から`client`に変更

### 5.4 ドキュメントの更新

1. `docs/ja/Docker.md`を更新
2. `docs/en/Docker.md`を更新
3. `README.md`を更新
4. `README.ja.md`を更新

## 6. テスト計画

### 6.1 Dockerfile.jobqueueのテスト

1. **ビルドテスト**: `docker build -f server/Dockerfile.jobqueue -t jobqueue-test ./server`でビルドが成功することを確認
2. **イメージサイズ確認**: 既存のAPI/Adminサーバーと同等のイメージサイズであることを確認
3. **バイナリ確認**: ビルドしたイメージ内で`jobqueue`バイナリが存在することを確認

### 6.2 docker-compose.jobqueue.ymlのテスト

1. **ビルドテスト**: `docker-compose -f docker-compose.jobqueue.yml build`でビルドが成功することを確認
2. **起動テスト**: `docker-compose -f docker-compose.jobqueue.yml up -d`でコンテナが正常に起動することを確認
3. **ヘルスチェックテスト**: `/health`エンドポイントが正常に応答することを確認
4. **Redis接続テスト**: Redisへの接続が正常に動作することを確認
5. **ログ確認**: ログにエラーが出力されていないことを確認

### 6.3 既存docker-composeファイルの修正テスト

1. **コンテナ名確認**: 修正後のコンテナ名が正しく設定されていることを確認
2. **起動テスト**: 修正後のdocker-composeファイルでコンテナが正常に起動することを確認
3. **既存機能確認**: 既存のAPI/Admin/Clientサーバーの機能が正常に動作することを確認

### 6.4 統合テスト

1. **全サーバー起動テスト**: PostgreSQL、Redis、API、Admin、JobQueue、Clientサーバーをすべて起動し、正常に動作することを確認
2. **サービス間通信テスト**: 各サーバー間の通信が正常に動作することを確認
3. **ジョブ処理テスト**: JobQueueサーバーが正常にジョブを処理できることを確認

## 7. 注意事項

### 7.1 コンテナ名変更の影響

- 既存のコンテナ（`api-develop`、`admin-develop`、`client-develop`）が起動している場合、停止して再起動する必要がある
- コンテナ名の変更により、既存のコンテナと新しいコンテナが競合する可能性がある
- 修正後は既存のコンテナを停止してから新しいコンテナを起動することを推奨

### 7.2 ネットワーク設定

- `postgres-network`と`redis-network`は外部ネットワーク（external: true）として定義されている
- これらのネットワークが存在しない場合、事前に作成する必要がある
- PostgreSQLとRedisコンテナを起動することで、これらのネットワークが作成される

### 7.3 環境変数

- 現在のジョブ処理ではPostgreSQLを使用していないが、将来的にジョブ処理でPostgreSQLを使用する可能性が高いため、拡張性を考慮してデータベースDSN環境変数を設定
- `REDIS_JOBQUEUE_ADDR=redis:6379`はDocker環境でのRedis接続先を指定

### 7.4 ヘルスチェック

- ヘルスチェックは`wget`を使用して`/health`エンドポイントを確認
- `wget`はAlpineイメージにインストールする必要がある（Dockerfileでインストール済み）

## 8. 参考情報

### 8.1 既存ファイル

- `server/Dockerfile`: APIサーバー用Dockerfile
- `server/Dockerfile.admin`: Adminサーバー用Dockerfile
- `docker-compose.api.yml`: APIサーバー用docker-compose設定
- `docker-compose.admin.yml`: Adminサーバー用docker-compose設定
- `docker-compose.client.yml`: クライアント用docker-compose設定

### 8.2 関連ドキュメント

- `docs/ja/Docker.md`: Docker環境の説明（日本語）
- `docs/en/Docker.md`: Docker環境の説明（英語）
- `README.md`: プロジェクトの説明（英語）
- `README.ja.md`: プロジェクトの説明（日本語）

### 8.3 技術スタック

- **Go**: 1.24
- **ベースイメージ**: `golang:1.24-alpine`（ビルドステージ）、`alpine:latest`（実行ステージ）
- **CGO**: CGO_ENABLED=0
- **Redis**: AsynqサーバーがRedisからジョブを取得
- **HTTPサーバー**: ポート8082、`/health`エンドポイント提供
