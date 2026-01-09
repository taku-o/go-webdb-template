# Docker

本プロジェクトのDocker環境について説明します。

## 概要

APIサーバー、Adminサーバー、クライアントサーバーをDockerコンテナ上で動作させることができます。

### 対応環境

| 環境 | 用途 | データベース |
|------|------|-------------|
| develop | 開発環境 | SQLite |
| staging | ステージング環境 | PostgreSQL/MySQL |
| production | 本番環境 | PostgreSQL/MySQL |

### アーキテクチャ

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

---

## 前提条件

- Docker Desktop（macOS/Windows）または Docker Engine（Linux）
- Docker Compose v2以上
- 以下のポートが使用可能であること：
  - 8080（APIサーバー）
  - 8081（Adminサーバー）
  - 3000（クライアントサーバー）
  - 5432（PostgreSQL）
  - 6379（Redis）

---

## Dockerfileの構成

### サーバー側（Go）

| ファイル | 用途 | CGO | ベースイメージ |
|---------|------|-----|---------------|
| `server/Dockerfile` | staging/production | 0 | golang:1.24-bookworm → debian:bookworm-slim |
| `server/Dockerfile.develop` | develop | 1 | golang:1.24-bookworm → debian:bookworm-slim |
| `server/Dockerfile.admin` | staging/production | 0 | golang:1.24-bookworm → debian:bookworm-slim |
| `server/Dockerfile.admin.develop` | develop | 1 | golang:1.24-bookworm → debian:bookworm-slim |

**注意**: 開発環境（develop）ではSQLiteを使用するため、CGO_ENABLED=1でビルドする必要があります。

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
| `docker-compose.api.develop.yml` | APIサーバー | develop |
| `docker-compose.api.staging.yml` | APIサーバー | staging |
| `docker-compose.api.production.yml` | APIサーバー | production |
| `docker-compose.admin.develop.yml` | Adminサーバー | develop |
| `docker-compose.admin.staging.yml` | Adminサーバー | staging |
| `docker-compose.admin.production.yml` | Adminサーバー | production |
| `docker-compose.client.develop.yml` | クライアント | develop |
| `docker-compose.client.staging.yml` | クライアント | staging |
| `docker-compose.client.production.yml` | クライアント | production |

### ボリュームマウント

| パス | 用途 | モード |
|------|------|--------|
| `./config/{env}:/app/config/{env}` | 設定ファイル | 読み取り専用 |
| `./server/data:/app/server/data` | SQLiteデータベース | 読み書き可 |
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
docker-compose -f docker-compose.api.develop.yml build
docker-compose -f docker-compose.admin.develop.yml build
docker-compose -f docker-compose.client.develop.yml build

# 起動
docker-compose -f docker-compose.api.develop.yml up -d
docker-compose -f docker-compose.admin.develop.yml up -d
docker-compose -f docker-compose.client.develop.yml up -d

# 停止
docker-compose -f docker-compose.api.develop.yml down
docker-compose -f docker-compose.admin.develop.yml down
docker-compose -f docker-compose.client.develop.yml down

# ログ確認
docker-compose -f docker-compose.api.develop.yml logs -f
docker-compose -f docker-compose.admin.develop.yml logs -f
docker-compose -f docker-compose.client.develop.yml logs -f
```

### ステージング環境（staging）

```bash
# ビルド
docker-compose -f docker-compose.api.staging.yml build
docker-compose -f docker-compose.admin.staging.yml build
docker-compose -f docker-compose.client.staging.yml build

# 起動
docker-compose -f docker-compose.api.staging.yml up -d
docker-compose -f docker-compose.admin.staging.yml up -d
docker-compose -f docker-compose.client.staging.yml up -d

# 停止
docker-compose -f docker-compose.api.staging.yml down
docker-compose -f docker-compose.admin.staging.yml down
docker-compose -f docker-compose.client.staging.yml down
```

### 本番環境（production）

```bash
# ビルド
docker-compose -f docker-compose.api.production.yml build
docker-compose -f docker-compose.admin.production.yml build
docker-compose -f docker-compose.client.production.yml build

# 起動
docker-compose -f docker-compose.api.production.yml up -d
docker-compose -f docker-compose.admin.production.yml up -d
docker-compose -f docker-compose.client.production.yml up -d

# 停止
docker-compose -f docker-compose.api.production.yml down
docker-compose -f docker-compose.admin.production.yml down
docker-compose -f docker-compose.client.production.yml down
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
4. クライアントサーバーを起動

### サービス間通信

| 接続元 | 接続先 | ホスト名 |
|--------|--------|---------|
| APIサーバー | PostgreSQL | postgres:5432 |
| APIサーバー | Redis | redis:6379 |
| Adminサーバー | PostgreSQL | postgres:5432 |
| クライアント | APIサーバー | api:8080 |

---

## トラブルシューティング

### CGO関連エラー

**症状**: `Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work`

**原因**: 開発環境でSQLiteを使用する場合、CGO_ENABLED=1でビルドする必要があります。

**解決策**: `Dockerfile.develop`を使用してビルドしてください。

```bash
# 正しいDockerfileを使用
docker-compose -f docker-compose.api.develop.yml build
```

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
docker-compose -f docker-compose.api.develop.yml build --no-cache
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
docker-compose -f docker-compose.api.production.yml build

# Adminサーバー
docker-compose -f docker-compose.admin.production.yml build

# クライアントサーバー
docker-compose -f docker-compose.client.production.yml build
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
