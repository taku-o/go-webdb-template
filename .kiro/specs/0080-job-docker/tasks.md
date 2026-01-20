# JobQueueサーバーのDocker対応実装タスク一覧

## 概要
JobQueueサーバーをDocker上で動作させるための実装タスク一覧。要件定義書と設計書に基づき、実装可能な粒度でタスクを分解。

## 実装フェーズ

### Phase 1: Dockerfile.jobqueueの作成

#### - [ ] タスク 1.1: `server/Dockerfile.jobqueue`の作成
**目的**: JobQueueサーバー用のDockerfileを作成する。

**作業内容**:
- `server/Dockerfile`をコピーして`server/Dockerfile.jobqueue`を作成
- ビルドコマンドを`./cmd/jobqueue/main.go`に変更
- バイナリ名を`jobqueue`に変更
- ポートを`8082`に変更
- CMDを`["./jobqueue"]`に変更

**実装コード**:
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

**確認事項**:
- 既存の`server/Dockerfile`と`server/Dockerfile.admin`と同様の構造になっている
- ビルドコマンドが`./cmd/jobqueue/main.go`になっている
- バイナリ名が`jobqueue`になっている
- ポートが`8082`になっている
- CMDが`["./jobqueue"]`になっている
- `wget`パッケージがインストールされている（ヘルスチェック用）

### Phase 2: docker-compose.jobqueue.ymlの作成

#### - [ ] タスク 2.1: `docker-compose.jobqueue.yml`の作成
**目的**: JobQueueサーバー用のdocker-compose設定ファイルを作成する。

**作業内容**:
- `docker-compose.api.yml`をコピーして`docker-compose.jobqueue.yml`を作成
- サービス名を`jobqueue`に変更
- コンテナ名を`jobqueue`に変更
- ビルド設定の`dockerfile`を`Dockerfile.jobqueue`に変更
- ポートを`8082:8082`に変更
- ヘルスチェックのポートを`8082`に変更
- 環境変数は既存のAPI/Adminサーバーと同様に設定（データベースDSNも含む）

**実装コード**:
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

**確認事項**:
- サービス名が`jobqueue`になっている
- コンテナ名が`jobqueue`になっている（環境名を含めない）
- ビルド設定が`context: ./server`と`dockerfile: Dockerfile.jobqueue`になっている
- ポート設定が`"8082:8082"`になっている
- 環境変数が適切に設定されている（`APP_ENV=develop`、`REDIS_JOBQUEUE_ADDR=redis:6379`、データベースDSN等）
- ボリュームマウントが適切に設定されている
- ネットワーク設定が`postgres-network`と`redis-network`になっている
- ヘルスチェックが適切に設定されている（ポート8082）

### Phase 3: 既存docker-composeファイルのコンテナ名修正

#### - [ ] タスク 3.1: `docker-compose.api.yml`のコンテナ名修正
**目的**: APIサーバーのコンテナ名を`api-develop`から`api`に変更する。

**作業内容**:
- `docker-compose.api.yml`を開く
- `container_name: api-develop`を`container_name: api`に変更

**変更内容**:
```yaml
# 変更前
container_name: api-develop

# 変更後
container_name: api
```

**確認事項**:
- コンテナ名が`api`になっている
- 他の設定に影響がない

#### - [ ] タスク 3.2: `docker-compose.admin.yml`のコンテナ名修正
**目的**: Adminサーバーのコンテナ名を`admin-develop`から`admin`に変更する。

**作業内容**:
- `docker-compose.admin.yml`を開く
- `container_name: admin-develop`を`container_name: admin`に変更

**変更内容**:
```yaml
# 変更前
container_name: admin-develop

# 変更後
container_name: admin
```

**確認事項**:
- コンテナ名が`admin`になっている
- 他の設定に影響がない

#### - [ ] タスク 3.3: `docker-compose.client.yml`のコンテナ名修正
**目的**: クライアントサーバーのコンテナ名を`client-develop`から`client`に変更する。

**作業内容**:
- `docker-compose.client.yml`を開く
- `container_name: client-develop`を`container_name: client`に変更

**変更内容**:
```yaml
# 変更前
container_name: client-develop

# 変更後
container_name: client
```

**確認事項**:
- コンテナ名が`client`になっている
- 他の設定に影響がない

### Phase 4: ドキュメントの更新

#### - [ ] タスク 4.1: `docs/ja/Docker.md`の更新
**目的**: 日本語版のDockerドキュメントにJobQueueサーバーの情報を追加する。

**作業内容**:
- 概要セクションにJobQueueサーバーを追加
- アーキテクチャ図にJobQueueサーバーを追加
- 前提条件セクションにポート8082を追加
- Dockerfile構成表に`server/Dockerfile.jobqueue`を追加
- Docker Compose設定ファイル一覧に`docker-compose.jobqueue.yml`を追加
- ビルド・起動・停止コマンドにJobQueueサーバー用のコマンドを追加
- 起動順序にJobQueueサーバーを追加
- サービス間通信にJobQueueサーバーとRedisの接続を追加

**確認事項**:
- すべてのセクションでJobQueueサーバーの情報が追加されている
- 既存の情報に影響がない

#### - [ ] タスク 4.2: `docs/en/Docker.md`の更新
**目的**: 英語版のDockerドキュメントにJobQueueサーバーの情報を追加する。

**作業内容**:
- `docs/ja/Docker.md`と同様の更新を英語版にも適用
- 概要セクションにJobQueueサーバーを追加
- アーキテクチャ図にJobQueueサーバーを追加
- 前提条件セクションにポート8082を追加
- Dockerfile構成表に`server/Dockerfile.jobqueue`を追加
- Docker Compose設定ファイル一覧に`docker-compose.jobqueue.yml`を追加
- ビルド・起動・停止コマンドにJobQueueサーバー用のコマンドを追加
- 起動順序にJobQueueサーバーを追加
- サービス間通信にJobQueueサーバーとRedisの接続を追加

**確認事項**:
- すべてのセクションでJobQueueサーバーの情報が追加されている
- 既存の情報に影響がない
- 日本語版と内容が一致している

#### - [ ] タスク 4.3: `README.md`の更新
**目的**: 英語版のREADMEにJobQueueサーバーのDocker起動方法を追加する。

**作業内容**:
- JobQueue Serverセクションを追加
- Docker環境でのビルド・起動・停止コマンドを追加
- JobQueueサーバーの説明を追加

**追加内容**:
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

**確認事項**:
- JobQueue Serverセクションが追加されている
- コマンドが正しく記載されている
- 説明が適切である
- 既存の情報に影響がない

#### - [ ] タスク 4.4: `README.ja.md`の更新
**目的**: 日本語版のREADMEにJobQueueサーバーのDocker起動方法を追加する。

**作業内容**:
- `README.md`と同様の更新を日本語版にも適用
- JobQueueサーバーセクションを追加
- Docker環境でのビルド・起動・停止コマンドを追加
- JobQueueサーバーの説明を追加

**追加内容**:
```markdown
#### JobQueueサーバー

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

JobQueueサーバーはHTTPサーバー（ポート8082）とAsynqサーバーを並行して起動します。HTTPサーバーは死活監視用の`/health`エンドポイントを提供します。Asynqサーバーは Redis からジョブを取得して処理します。JobQueueサーバーを起動する前に、Redisが起動していることを確認してください。
```

**確認事項**:
- JobQueueサーバーセクションが追加されている
- コマンドが正しく記載されている
- 説明が適切である
- 既存の情報に影響がない
- 英語版と内容が一致している

### Phase 5: ビルドと動作確認

#### - [ ] タスク 5.1: Dockerfile.jobqueueのビルド確認
**目的**: Dockerfile.jobqueueが正常にビルドできることを確認する。

**作業内容**:
- `docker build -f server/Dockerfile.jobqueue -t jobqueue-test ./server`を実行
- ビルドが成功することを確認
- イメージサイズを確認（既存のAPI/Adminサーバーと同等であることを確認）

**確認事項**:
- ビルドが正常に完了する
- エラーが発生しない
- イメージサイズが既存のAPI/Adminサーバーと同等である

#### - [ ] タスク 5.2: docker-compose.jobqueue.ymlのビルド確認
**目的**: docker-compose.jobqueue.ymlでイメージが正常にビルドできることを確認する。

**作業内容**:
- `docker-compose -f docker-compose.jobqueue.yml build`を実行
- ビルドが成功することを確認

**確認事項**:
- ビルドが正常に完了する
- エラーが発生しない

#### - [ ] タスク 5.3: JobQueueサーバーの起動確認
**目的**: JobQueueサーバーが正常に起動することを確認する。

**作業内容**:
- PostgreSQLとRedisコンテナが起動していることを確認
- `docker-compose -f docker-compose.jobqueue.yml up -d`を実行
- コンテナが正常に起動することを確認
- ログを確認してエラーが出力されていないことを確認

**確認事項**:
- コンテナが正常に起動する
- ログにエラーが出力されていない
- コンテナ名が`jobqueue`になっている

#### - [ ] タスク 5.4: ヘルスチェックの確認
**目的**: `/health`エンドポイントが正常に動作することを確認する。

**作業内容**:
- `curl http://localhost:8082/health`を実行
- `OK`が返されることを確認
- Dockerのヘルスチェックが正常に動作することを確認（`docker ps`で`healthy`ステータスを確認）

**確認事項**:
- `/health`エンドポイントが正常に応答する
- `OK`が返される
- Dockerのヘルスチェックが正常に動作する

#### - [ ] タスク 5.5: Redis接続の確認
**目的**: Redisへの接続が正常に動作することを確認する。

**作業内容**:
- JobQueueサーバーのログを確認
- Redis接続エラーが出力されていないことを確認
- Asynqサーバーが正常に起動していることを確認

**確認事項**:
- Redis接続エラーが出力されていない
- Asynqサーバーが正常に起動している

#### - [ ] タスク 5.6: 既存docker-composeファイルの修正確認
**目的**: 既存のdocker-composeファイルのコンテナ名修正が正しく反映されていることを確認する。

**作業内容**:
- `docker-compose.api.yml`のコンテナ名が`api`になっていることを確認
- `docker-compose.admin.yml`のコンテナ名が`admin`になっていることを確認
- `docker-compose.client.yml`のコンテナ名が`client`になっていることを確認
- 既存のコンテナを停止して再起動し、新しいコンテナ名で起動することを確認

**確認事項**:
- すべてのdocker-composeファイルのコンテナ名が正しく修正されている
- 新しいコンテナ名でコンテナが起動する
- 既存の機能に影響がない

#### - [ ] タスク 5.7: 統合テスト
**目的**: すべてのサーバーが正常に動作することを確認する。

**作業内容**:
- PostgreSQL、Redis、API、Admin、JobQueue、Clientサーバーをすべて起動
- 各サーバーが正常に動作することを確認
- サービス間通信が正常に動作することを確認

**確認事項**:
- すべてのサーバーが正常に起動する
- サービス間通信が正常に動作する
- エラーが発生しない

## 受け入れ基準の確認

### Dockerfileの作成
- [ ] `server/Dockerfile.jobqueue`が作成されている
- [ ] ベースイメージが`golang:1.24-alpine`（ビルドステージ）と`alpine:latest`（実行ステージ）である
- [ ] CGO_ENABLED=0でビルドされている
- [ ] 必要なパッケージ（`ca-certificates`, `tzdata`, `wget`）がインストールされている
- [ ] 非rootユーザー（appuser）が作成されている
- [ ] ディレクトリ構造が既存のAPI/Adminサーバーと同様である
- [ ] ポート8082がEXPOSEされている
- [ ] CMDが`["./jobqueue"]`である

### docker-compose設定ファイルの作成
- [ ] `docker-compose.jobqueue.yml`が作成されている
- [ ] サービス名が`jobqueue`である
- [ ] コンテナ名が`jobqueue`である
- [ ] `docker-compose.api.yml`のコンテナ名が`api`に修正されている
- [ ] `docker-compose.admin.yml`のコンテナ名が`admin`に修正されている
- [ ] `docker-compose.client.yml`のコンテナ名が`client`に修正されている
- [ ] ビルド設定が`context: ./server`と`dockerfile: Dockerfile.jobqueue`である
- [ ] ポート設定が`"8082:8082"`である
- [ ] 環境変数が適切に設定されている（`APP_ENV=develop`、`REDIS_JOBQUEUE_ADDR=redis:6379`等）
- [ ] ボリュームマウントが適切に設定されている
- [ ] ネットワーク設定が`postgres-network`と`redis-network`である
- [ ] ヘルスチェックが適切に設定されている

### ビルドと動作確認
- [ ] `docker-compose.jobqueue.yml`でイメージが正常にビルドできる
- [ ] ビルドしたイメージからコンテナが正常に起動できる
- [ ] JobQueueサーバーが正常に動作する（HTTPサーバーとAsynqサーバー）
- [ ] `/health`エンドポイントが正常に動作する
- [ ] Redis接続が正常に動作する
- [ ] ヘルスチェックが正常に動作する

### ドキュメント
- [ ] `docs/ja/Docker.md`が更新されている
- [ ] `docs/en/Docker.md`が更新されている
- [ ] `README.md`が更新されている
- [ ] `README.ja.md`が更新されている
- [ ] すべてのドキュメントでJobQueueサーバーのDocker対応が記載されている

### 既存機能との整合性
- [ ] 既存のAPI/AdminサーバーのDocker設定に影響がない
- [ ] 既存のdocker-composeファイルに影響がない
- [ ] 既存のDockerfileに影響がない
- [ ] 既存のドキュメントの他の部分に影響がない
