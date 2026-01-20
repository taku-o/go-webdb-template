**[日本語]** | [English](../en/Docker.md)

# Docker

本プロジェクトのDocker環境について説明します。

## 概要

APIサーバー、Adminサーバー、JobQueueサーバー、クライアントサーバーをDockerコンテナ上で動作させることができます。

### 対応環境

| 環境 | 用途 | データベース |
|------|------|-------------|
| develop | 開発環境 | PostgreSQL/MySQL |

### アーキテクチャ

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

---

## 前提条件

- Docker Desktop（macOS/Windows）または Docker Engine（Linux）
- Docker Compose v2以上
- 以下のポートが使用可能であること：
  - 8080（APIサーバー）
  - 8081（Adminサーバー）
  - 8082（JobQueueサーバー）
  - 3000（クライアントサーバー）
  - 5432（PostgreSQL）
  - 6379（Redis）

---

## Dockerfileの構成

### サーバー側（Go）

| ファイル | 用途 | CGO | ベースイメージ |
|---------|------|-----|---------------|
| `server/Dockerfile` | 全環境（develop/staging/production） | 0 | golang:1.24-alpine → alpine:latest |
| `server/Dockerfile.admin` | 全環境（develop/staging/production） | 0 | golang:1.24-alpine → alpine:latest |
| `server/Dockerfile.jobqueue` | 全環境（develop/staging/production） | 0 | golang:1.24-alpine → alpine:latest |

### クライアント側（Next.js）

| ファイル | ターゲット | 用途 |
|---------|-----------|------|
| `client/Dockerfile` | dev | 開発環境（ホットリロード対応） |
| `client/Dockerfile` | production | staging/production環境 |

---

## Docker Compose設定ファイル

### ファイル一覧

| ファイル | サービス | 環境 |
|---------|---------|------|
| `docker-compose.api.yml` | APIサーバー | develop |
| `docker-compose.admin.yml` | Adminサーバー | develop |
| `docker-compose.jobqueue.yml` | JobQueueサーバー | develop |
| `docker-compose.client.yml` | クライアント | develop |

### ボリュームマウント

| パス | 用途 | モード |
|------|------|--------|
| `./config/{env}:/app/config/{env}` | 設定ファイル | 読み取り専用 |
| `./server/data:/app/server/data` | データディレクトリ | 読み書き可 |
| `./logs:/app/logs` | ログファイル | 読み書き可 |

### 環境変数

| 変数名 | 説明 | 例 |
|--------|------|-----|
| `APP_ENV` | 環境指定 | develop/staging/production |
| `REDIS_JOBQUEUE_ADDR` | Redis接続先 | redis:6379 |
| `NEXT_PUBLIC_API_URL` | APIサーバーURL | http://api:8080 |
| `NEXT_PUBLIC_API_KEY` | APIキー | （.env.localから読み込み） |

---

## ビルド・起動・停止コマンド

### 開発環境（develop）

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

---

## 既存サービスとの統合

### 必要な外部ネットワーク

Docker化されたサーバーは、既存のPostgreSQL・Redisコンテナと同じネットワークで通信します。

```bash
# PostgreSQLコンテナの起動（postgres-networkを作成）
docker-compose -f docker-compose.postgres.yml up -d

# Redisコンテナの起動（redis-networkを作成）
docker-compose -f docker-compose.redis.yml up -d
```

### 起動順序

1. PostgreSQL、Redisコンテナを起動
2. APIサーバーを起動
3. Adminサーバーを起動
4. JobQueueサーバーを起動
5. クライアントサーバーを起動

### サービス間通信

| 接続元 | 接続先 | ホスト名 |
|--------|--------|---------|
| APIサーバー | PostgreSQL | postgres:5432 |
| APIサーバー | Redis | redis:6379 |
| Adminサーバー | PostgreSQL | postgres:5432 |
| JobQueueサーバー | Redis | redis:6379 |
| クライアント | APIサーバー | api:8080 |

---

## PostgreSQLコンテナの起動・管理

### 構成概要

本プロジェクトではmaster 1台 + sharding 4台のPostgreSQLコンテナを使用します。

| コンテナ名 | データベース名 | ホストポート |
|-----------|--------------|-------------|
| postgres-master | webdb_master | 5432 |
| postgres-sharding-1 | webdb_sharding_1 | 5433 |
| postgres-sharding-2 | webdb_sharding_2 | 5434 |
| postgres-sharding-3 | webdb_sharding_3 | 5435 |
| postgres-sharding-4 | webdb_sharding_4 | 5436 |

### 起動・停止コマンド

```bash
# 起動
./scripts/start-postgres.sh start

# 停止
./scripts/start-postgres.sh stop

# 状態確認
./scripts/start-postgres.sh status

# ヘルスチェック
./scripts/start-postgres.sh health
```

### マイグレーション

```bash
# PostgreSQL起動後にマイグレーションを適用
APP_ENV=develop ./scripts/migrate.sh all
```

マイグレーションの詳細は [Atlas-Operations.md](./Atlas-Operations.md) を参照してください。

---

## トラブルシューティング

### ネットワーク接続エラー

**症状**: `network postgres-network not found`

**原因**: 外部ネットワークが作成されていません。

**解決策**: PostgreSQL/Redisコンテナを先に起動してください。

```bash
docker-compose -f docker-compose.postgres.yml up -d
docker-compose -f docker-compose.redis.yml up -d
```

### ポート競合

**症状**: `Bind for 0.0.0.0:8080 failed: port is already allocated`

**原因**: 指定したポートが既に使用されています。

**解決策**:
1. 使用中のプロセスを停止する
2. または、docker-compose.ymlのポートマッピングを変更する

```bash
# 使用中のプロセスを確認
lsof -i :8080

# プロセスを停止
kill <PID>
```

### イメージの再ビルド

キャッシュを使用せずに再ビルドする場合：

```bash
docker-compose -f docker-compose.api.yml build --no-cache
```

### 全コンテナ・イメージの削除

```bash
# 全コンテナ停止・削除
docker stop $(docker ps -aq)
docker rm $(docker ps -aq)

# 未使用イメージ削除
docker image prune -a -f

# ビルドキャッシュ削除
docker builder prune -a -f
```

---

## 本番環境へのデプロイ

### デプロイフロー

1. 本番用イメージをビルド
2. イメージにタグを付与
3. コンテナレジストリにプッシュ
4. 本番環境でイメージをプル・起動

---

## コンテナレジストリへのプッシュ

本番環境へのデプロイに向けて、Dockerイメージをコンテナレジストリにプッシュする手順を説明します。

### イメージのビルド

プッシュ前に、本番用イメージをビルドします。

```bash
# APIサーバー
docker-compose -f docker-compose.api.yml build

# Adminサーバー
docker-compose -f docker-compose.admin.yml build

# クライアントサーバー
docker-compose -f docker-compose.client.yml build
```

### イメージのタグ付け

レジストリにプッシュするためにタグを付与します。

```bash
# APIサーバー
docker tag go-webdb-template-api:latest <registry>/api:<version>

# Adminサーバー
docker tag go-webdb-template-admin:latest <registry>/admin:<version>

# クライアントサーバー
docker tag go-webdb-template-client:latest <registry>/client:<version>
```

**タグ形式の例**:
- `api:v1.0.0` - バージョン指定
- `api:latest` - 最新版
- `api:staging` - ステージング環境用

---

### AWS ECR（Elastic Container Registry）

#### 認証

```bash
# AWS CLIでECRにログイン
aws ecr get-login-password --region <region> | docker login --username AWS --password-stdin <account-id>.dkr.ecr.<region>.amazonaws.com
```

**パラメータ**:
- `<region>`: AWSリージョン（例: `ap-northeast-1`）
- `<account-id>`: AWSアカウントID（12桁の数字）

#### リポジトリ作成（初回のみ）

```bash
# APIサーバー用リポジトリ作成
aws ecr create-repository --repository-name api --region <region>

# Adminサーバー用リポジトリ作成
aws ecr create-repository --repository-name admin --region <region>

# クライアントサーバー用リポジトリ作成
aws ecr create-repository --repository-name client --region <region>
```

#### タグ付けとプッシュ

```bash
# タグ付け
docker tag go-webdb-template-api:latest <account-id>.dkr.ecr.<region>.amazonaws.com/api:v1.0.0
docker tag go-webdb-template-admin:latest <account-id>.dkr.ecr.<region>.amazonaws.com/admin:v1.0.0
docker tag go-webdb-template-client:latest <account-id>.dkr.ecr.<region>.amazonaws.com/client:v1.0.0

# プッシュ
docker push <account-id>.dkr.ecr.<region>.amazonaws.com/api:v1.0.0
docker push <account-id>.dkr.ecr.<region>.amazonaws.com/admin:v1.0.0
docker push <account-id>.dkr.ecr.<region>.amazonaws.com/client:v1.0.0
```

---

### Tencent Cloud TCR（Tencent Container Registry）

#### 認証

```bash
# TCRにログイン
docker login <registry-name>.tencentcloudcr.com --username <username>
```

**パラメータ**:
- `<registry-name>`: TCRインスタンス名
- `<username>`: Tencent Cloudアカウントのユーザー名またはアクセスキーID

パスワードはTencent Cloudコンソールから取得したアクセスキーシークレットを使用します。

#### 名前空間作成（初回のみ）

Tencent Cloudコンソールから名前空間を作成するか、CLIを使用します。

#### タグ付けとプッシュ

```bash
# タグ付け
docker tag go-webdb-template-api:latest <registry-name>.tencentcloudcr.com/<namespace>/api:v1.0.0
docker tag go-webdb-template-admin:latest <registry-name>.tencentcloudcr.com/<namespace>/admin:v1.0.0
docker tag go-webdb-template-client:latest <registry-name>.tencentcloudcr.com/<namespace>/client:v1.0.0

# プッシュ
docker push <registry-name>.tencentcloudcr.com/<namespace>/api:v1.0.0
docker push <registry-name>.tencentcloudcr.com/<namespace>/admin:v1.0.0
docker push <registry-name>.tencentcloudcr.com/<namespace>/client:v1.0.0
```

---

### Docker Hub（オプション）

#### 認証

```bash
# Docker Hubにログイン
docker login --username <username>
```

パスワードまたはアクセストークンを入力します。

#### タグ付けとプッシュ

```bash
# タグ付け
docker tag go-webdb-template-api:latest <username>/api:v1.0.0
docker tag go-webdb-template-admin:latest <username>/admin:v1.0.0
docker tag go-webdb-template-client:latest <username>/client:v1.0.0

# プッシュ
docker push <username>/api:v1.0.0
docker push <username>/admin:v1.0.0
docker push <username>/client:v1.0.0
```

---

### セキュリティに関する注意事項

1. **認証情報の管理**
   - 認証情報はCI/CD環境の秘密変数として管理する
   - 認証情報をソースコードにコミットしない
   - IAMロールやサービスアカウントを使用して認証を自動化する

2. **イメージスキャン**
   - プッシュ前にセキュリティスキャンを実行することを推奨
   - AWS ECRやTCRの脆弱性スキャン機能を活用する

3. **タグの運用**
   - 本番環境には特定バージョンのタグを使用する（`latest`は避ける）
   - セマンティックバージョニング（例: `v1.2.3`）を推奨
