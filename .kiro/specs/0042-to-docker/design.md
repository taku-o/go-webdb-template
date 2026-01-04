# Docker化設計書

## 1. 概要

### 1.1 設計の目的
要件定義書に基づき、APIサーバー、クライアントサーバー、AdminサーバーをDockerコンテナ上で動作させるための詳細設計を定義する。開発環境の統一と本番環境（AWSまたはTencent Cloud）へのデプロイを容易にするため、Dockerイメージの作成とDocker Composeによる管理を実現する。

### 1.2 設計の範囲
- APIサーバー（Go）のDockerfile作成
- Adminサーバー（Go）のDockerfile作成
- クライアントサーバー（Next.js）のDockerfile作成
- 環境別docker-composeファイルの作成（9ファイル）
- 既存のPostgreSQL、Redisコンテナとの統合
- Dockerイメージの最適化（マルチステージビルド）
- 環境別設定（develop/staging/production）のサポート
- 本番環境へのデプロイ準備

### 1.3 設計方針
- **既存システムとの統合**: 既存のPostgreSQL、Redisコンテナと同一ネットワークで通信
- **環境別分離**: サービス別・環境別にdocker-composeファイルを分離
- **マルチステージビルド**: イメージサイズを最小化
- **非rootユーザー実行**: セキュリティベストプラクティスに準拠
- **既存設定ファイルの維持**: `config/{env}/`ディレクトリ構造を維持
- **データ永続化**: ボリュームマウントにより設定ファイルとデータベースファイルを永続化

## 2. アーキテクチャ設計

### 2.1 ディレクトリ構造

#### 2.1.1 変更前の構造
```
.
├── server/
│   ├── cmd/
│   │   ├── server/
│   │   │   └── main.go
│   │   └── admin/
│   │       └── main.go
│   ├── data/                    # SQLiteデータベースファイル
│   └── ...
├── client/
│   ├── src/
│   ├── package.json
│   └── ...
├── config/
│   ├── develop/
│   ├── staging/
│   └── production/
└── docker-compose.postgres.yml
```

#### 2.1.2 変更後の構造
```
.
├── server/
│   ├── Dockerfile                # APIサーバー用（新規）
│   ├── Dockerfile.admin          # Adminサーバー用（新規）
│   ├── .dockerignore             # APIサーバー用（新規）
│   ├── cmd/
│   │   ├── server/
│   │   │   └── main.go
│   │   └── admin/
│   │       └── main.go
│   ├── data/                     # 既存（維持）
│   └── ...
├── client/
│   ├── Dockerfile                # クライアントサーバー用（新規）
│   ├── .dockerignore             # クライアントサーバー用（新規）
│   ├── src/
│   ├── package.json
│   └── ...
├── config/
│   ├── develop/                  # 既存（維持）
│   ├── staging/                  # 既存（維持）
│   └── production/               # 既存（維持）
├── docker-compose.api.develop.yml        # 新規
├── docker-compose.api.staging.yml        # 新規
├── docker-compose.api.production.yml     # 新規
├── docker-compose.client.develop.yml     # 新規
├── docker-compose.client.staging.yml     # 新規
├── docker-compose.client.production.yml  # 新規
├── docker-compose.admin.develop.yml      # 新規
├── docker-compose.admin.staging.yml      # 新規
├── docker-compose.admin.production.yml   # 新規
├── docker-compose.postgres.yml           # 既存（維持）
├── docker-compose.redis.yml              # 既存（維持）
└── docs/
    └── Docker.md                 # 新規
```

### 2.2 システム構成図

```
┌─────────────────────────────────────────────────────────────┐
│                    開発者/運用者                              │
└──────────────────┬────────────────────────────────────┘
                    │
                    │ docker-compose -f docker-compose.{service}.{env}.yml up
                    │
                    ▼
┌─────────────────────────────────────────────────────────────┐
│              Docker Compose (環境別ファイル)                  │
│  - docker-compose.api.{env}.yml                             │
│  - docker-compose.client.{env}.yml                           │
│  - docker-compose.admin.{env}.yml                            │
└──────────────────┬────────────────────────────────────┘
                    │
                    │ Dockerネットワーク
                    │
        ┌───────────┴───────────┐
        │                       │
        ▼                       ▼
┌──────────────┐      ┌──────────────┐
│  API Server │      │ Admin Server │
│  (Port 8080)│      │ (Port 8081)  │
│             │      │              │
│  ┌────────┐ │      │  ┌────────┐ │
│  │ Go App │ │      │  │ Go App │ │
│  └────┬───┘ │      │  └────┬───┘ │
└───────┼─────┘      └───────┼─────┘
        │                    │
        │                    │
        ▼                    ▼
┌──────────────┐      ┌──────────────┐
│  PostgreSQL │      │    Redis     │
│  Container  │      │  Container   │
│  (postgres) │      │   (redis)    │
└──────────────┘      └──────────────┘
        │                    │
        │                    │
        └──────────┬─────────┘
                   │
                   ▼
        ┌──────────────────┐
        │  Client Server  │
        │   (Port 3000)   │
        │                 │
        │  ┌────────────┐ │
        │  │ Next.js   │ │
        │  └─────┬─────┘ │
        └────────┼───────┘
                 │
                 │ HTTP Request
                 │
                 ▼
        ┌──────────────┐
        │  API Server  │
        │  (Port 8080) │
        └──────────────┘
```

### 2.3 データフロー

#### 2.3.1 コンテナ起動フロー
```
開発者が docker-compose -f docker-compose.{service}.{env}.yml up を実行
    ↓
Docker ComposeがDockerfileからイメージをビルド（初回のみ）
    ↓
コンテナが起動
    ↓
環境変数APP_ENVに基づいて設定ファイルを読み込み
    ↓
ボリュームマウントにより設定ファイルとデータディレクトリにアクセス
    ↓
既存のPostgreSQL、Redisコンテナと同一ネットワークで通信
    ↓
サービスが正常に起動
```

#### 2.3.2 クライアント→APIサーバー通信フロー
```
クライアント（Next.js）がHTTPリクエストを送信
    ↓
Dockerネットワーク経由でAPIサーバー（api:8080）に接続
    ↓
APIサーバーがリクエストを処理
    ↓
必要に応じてPostgreSQL、Redisコンテナに接続
    ↓
レスポンスをクライアントに返却
```

## 3. コンポーネント設計

### 3.1 APIサーバーDockerfile設計

#### 3.1.1 Dockerfile構造
```dockerfile
# ビルドステージ
FROM golang:1.21-alpine AS builder

WORKDIR /build

# 依存関係をコピー
COPY go.mod go.sum ./
RUN go mod download

# ソースコードをコピー
COPY . .

# ビルド
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server ./cmd/server/main.go

# 実行ステージ
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# 非rootユーザーを作成
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

# ビルド成果物をコピー
COPY --from=builder /build/server .

# 設定ファイル、データディレクトリ、ログディレクトリ用のディレクトリを作成
RUN mkdir -p /app/config /app/data /app/logs && \
    chown -R appuser:appuser /app

USER appuser

EXPOSE 8080

CMD ["./server"]
```

#### 3.1.2 .dockerignore構造
```
# server/.dockerignore
.git
.gitignore
*.md
*.test.go
*.sum
vendor/
.env
.env.local
*.log
coverage/
.idea/
.vscode/
```

### 3.2 AdminサーバーDockerfile設計

#### 3.2.1 Dockerfile.admin構造
```dockerfile
# ビルドステージ
FROM golang:1.21-alpine AS builder

WORKDIR /build

# 依存関係をコピー
COPY go.mod go.sum ./
RUN go mod download

# ソースコードをコピー
COPY . .

# ビルド
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o admin ./cmd/admin/main.go

# 実行ステージ
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# 非rootユーザーを作成
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

# ビルド成果物をコピー
COPY --from=builder /build/admin .

# 設定ファイル、データディレクトリ、ログディレクトリ用のディレクトリを作成
RUN mkdir -p /app/config /app/data /app/logs && \
    chown -R appuser:appuser /app

USER appuser

EXPOSE 8081

CMD ["./admin"]
```

### 3.3 クライアントサーバーDockerfile設計

#### 3.3.1 Dockerfile構造（開発モード対応）
```dockerfile
# 依存関係インストールステージ
FROM node:22-alpine AS deps

WORKDIR /app

COPY package.json package-lock.json ./
RUN npm ci

# ビルドステージ
FROM node:22-alpine AS builder

WORKDIR /app

COPY --from=deps /app/node_modules ./node_modules
COPY . .

# 環境変数を設定（ビルド時に必要）
# NEXT_PUBLIC_で始まる環境変数はビルド時にクライアント側のJavaScriptに埋め込まれる
ARG NEXT_PUBLIC_API_URL
ARG NEXT_PUBLIC_API_KEY
ENV NEXT_PUBLIC_API_URL=${NEXT_PUBLIC_API_URL}
ENV NEXT_PUBLIC_API_KEY=${NEXT_PUBLIC_API_KEY}

RUN npm run build

# 実行ステージ（開発モード）
FROM node:22-alpine AS dev

WORKDIR /app

COPY --from=deps /app/node_modules ./node_modules
COPY . .

ENV NODE_ENV=development
ENV NEXT_PUBLIC_API_URL=${NEXT_PUBLIC_API_URL:-http://api:8080}

EXPOSE 3000

CMD ["npm", "run", "dev"]

# 実行ステージ（本番モード）
FROM node:22-alpine AS production

WORKDIR /app

ENV NODE_ENV=production

# ビルド成果物をコピー
COPY --from=builder /app/.next ./.next
COPY --from=builder /app/public ./public
COPY --from=builder /app/package.json ./package.json
COPY --from=deps /app/node_modules ./node_modules

# 非rootユーザーを作成
RUN addgroup -g 1000 nodeuser && \
    adduser -D -u 1000 -G nodeuser nodeuser && \
    chown -R nodeuser:nodeuser /app

USER nodeuser

EXPOSE 3000

CMD ["npm", "start"]
```

#### 3.3.2 .dockerignore構造
```
# client/.dockerignore
.git
.gitignore
*.md
node_modules
.next
.env
.env.local
.env.development
.env.staging
.env.production
*.log
coverage/
.idea/
.vscode/
```

**注意**: `.env*`ファイルは全て`.dockerignore`で除外します。これにより：
- 機密情報（APIキー等）がイメージに含まれない
- 環境変数はDockerfileの`ARG`やdocker-compose.ymlの`environment`で明示的に指定する必要がある
- ビルド時に必要な`NEXT_PUBLIC_*`環境変数は、docker-compose.ymlの`build.args`で指定する

### 3.4 Docker Compose設定設計

#### 3.4.1 docker-compose.api.develop.yml構造
```yaml
version: '3.8'

services:
  api:
    build:
      context: ./server
      dockerfile: Dockerfile
    container_name: api-develop
    ports:
      - "8080:8080"
    environment:
      - APP_ENV=develop
    volumes:
      - ./config/develop:/app/config/develop:ro
      - ./server/data:/app/data
      - ./logs:/app/logs
    networks:
      - postgres-network
      - redis-network
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "wget", "--quiet", "--tries=1", "--spider", "http://localhost:8080/health"]
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

#### 3.4.2 docker-compose.client.develop.yml構造
```yaml
version: '3.8'

services:
  client:
    build:
      context: ./client
      dockerfile: Dockerfile
      target: dev
      args:
        - NEXT_PUBLIC_API_URL=http://api:8080
        - NEXT_PUBLIC_API_KEY=${NEXT_PUBLIC_API_KEY:-}
    container_name: client-develop
    ports:
      - "3000:3000"
    # .env.localファイルを読み込む（Not Docker版と挙動を統一するため）
    env_file:
      - ./client/.env.local
    environment:
      - NODE_ENV=development
      # 開発モードでは実行時にも環境変数を設定可能（ホットリロード時に反映）
      # env_fileで読み込んだ環境変数は、environmentセクションで上書き可能
      - NEXT_PUBLIC_API_URL=${NEXT_PUBLIC_API_URL:-http://api:8080}
      - NEXT_PUBLIC_API_KEY=${NEXT_PUBLIC_API_KEY:-}
    volumes:
      - ./client:/app
      - /app/node_modules
      - /app/.next
    networks:
      - postgres-network
    restart: unless-stopped
    depends_on:
      - api

networks:
  postgres-network:
    external: true
    name: postgres-network
```

#### 3.4.3 docker-compose.admin.develop.yml構造
```yaml
version: '3.8'

services:
  admin:
    build:
      context: ./server
      dockerfile: Dockerfile.admin
    container_name: admin-develop
    ports:
      - "8081:8081"
    environment:
      - APP_ENV=develop
    volumes:
      - ./config/develop:/app/config/develop:ro
      - ./server/data:/app/data
      - ./logs:/app/logs
    networks:
      - postgres-network
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "wget", "--quiet", "--tries=1", "--spider", "http://localhost:8081/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s

networks:
  postgres-network:
    external: true
    name: postgres-network
```

#### 3.4.4 環境別設定の違い

**staging/production環境**:
- 環境変数`APP_ENV`を`staging`または`production`に変更
- 設定ファイルのパスを`config/staging`または`config/production`に変更
- 本番環境では開発モードのボリュームマウントを削除（クライアント）
- リソース制限を追加（必要に応じて）

### 3.5 ネットワーク統合設計

#### 3.5.1 既存ネットワークの参照
既存の`docker-compose.postgres.yml`と`docker-compose.redis.yml`で定義されているネットワークを外部ネットワークとして参照する。

**postgres-network**:
- 既存の`docker-compose.postgres.yml`で定義
- ネットワーク名: `postgres-network`
- サービス名: `postgres`

**redis-network**:
- 既存の`docker-compose.redis.yml`で定義
- ネットワーク名: `redis-network`（または既存のネットワーク名）
- サービス名: `redis`

#### 3.5.2 サービス間通信
- APIサーバー → PostgreSQL: `postgres:5432`
- APIサーバー → Redis: `redis:6379`
- Adminサーバー → PostgreSQL: `postgres:5432`
- クライアント → APIサーバー: `api:8080`

## 4. 環境変数管理設計

### 4.1 環境変数の定義

#### 4.1.1 APIサーバー・Adminサーバー
- **APP_ENV**: 環境を指定（`develop`/`staging`/`production`）
  - デフォルト値: `develop`
  - 用途: 設定ファイルの読み込み先を決定
  - 設定方法: docker-compose.ymlの`environment`セクションで指定

#### 4.1.2 クライアントサーバー（Next.js）

**Next.jsの環境変数の特徴**:
- `NEXT_PUBLIC_`で始まる環境変数は、**ビルド時に**クライアント側のJavaScriptに埋め込まれます
- ビルド後に変更しても反映されません（本番モードの場合）
- 開発モードでは実行時にも環境変数を設定可能（ホットリロード時に反映）

**環境変数の定義**:
- **NODE_ENV**: Node.js環境（`development`/`production`）
  - デフォルト値: `development`（開発環境）、`production`（本番環境）
  - 設定方法: docker-compose.ymlの`environment`セクションで指定

- **NEXT_PUBLIC_API_URL**: APIサーバーのURL
  - デフォルト値: `http://api:8080`
  - 用途: クライアントからAPIサーバーへの接続
  - **ビルド時に必要**: Dockerfileの`ARG`で受け取り、`ENV`で設定
  - **実行時にも設定可能**: docker-compose.ymlの`environment`セクションで指定（開発モード）

- **NEXT_PUBLIC_API_KEY**: APIキー
  - デフォルト値: なし（空文字列）
  - 用途: API認証用のキー
  - **ビルド時に必要**: Dockerfileの`ARG`で受け取り、`ENV`で設定
  - **実行時にも設定可能**: docker-compose.ymlの`environment`セクションで指定（開発モード）

**重要な注意事項**:
- `.env.local`、`.env.development`などのファイルは`.dockerignore`で除外されます
- これらのファイルで定義されている環境変数は**自動的に取り込まれません**
- ビルド時に必要な環境変数は、docker-compose.ymlの`build.args`で明示的に指定する必要があります
- 実行時に変更可能な環境変数は、docker-compose.ymlの`environment`セクションで指定します

### 4.2 環境変数の設定方法

#### 4.2.1 docker-compose.ymlでの設定（APIサーバー・Adminサーバー）
```yaml
environment:
  - APP_ENV=develop
```

#### 4.2.2 docker-compose.ymlでの設定（クライアントサーバー）

**ビルド時の環境変数（build.args）**:
```yaml
build:
  context: ./client
  dockerfile: Dockerfile
  target: dev
  args:
    - NEXT_PUBLIC_API_URL=http://api:8080
    - NEXT_PUBLIC_API_KEY=${NEXT_PUBLIC_API_KEY:-}
```

**実行時の環境変数（env_file + environment）**:
```yaml
# .env.localファイルを読み込む（Not Docker版と挙動を統一）
env_file:
  - ./client/.env.local

environment:
  - NODE_ENV=development
  # env_fileで読み込んだ環境変数は、environmentセクションで上書き可能
  - NEXT_PUBLIC_API_URL=${NEXT_PUBLIC_API_URL:-http://api:8080}
  - NEXT_PUBLIC_API_KEY=${NEXT_PUBLIC_API_KEY:-}
```

**環境変数の受け渡しフロー**:
1. **ビルド時**: `build.args`で指定した環境変数がDockerfileの`ARG`に渡される
   - `${NEXT_PUBLIC_API_KEY:-}`は、ホスト側の環境変数または`.env`ファイルから解決される
   - `.env.local`の値は`env_file`で読み込まれるが、ビルド時には`build.args`で明示的に渡す必要がある
2. **ビルド時**: `ARG`から`ENV`に設定され、`npm run build`時にクライアント側のJavaScriptに埋め込まれる
3. **実行時**: 
   - `env_file`で`.env.local`を読み込み
   - `environment`セクションで指定した環境変数がコンテナ内で使用可能（`env_file`の値を上書き可能）
   - 開発モードでは実行時の環境変数がホットリロード時に反映される

#### 4.2.3 ホスト側の環境変数ファイル（.env）の扱い

**Docker Composeの環境変数解決**:
- Docker Composeは、`docker-compose.yml`と同じディレクトリにある`.env`ファイルを**自動的に読み込みます**
- `${変数名:-デフォルト値}`形式は、以下の順序で解決されます:
  1. ホスト側の環境変数
  2. `.env`ファイルの変数（docker-compose.ymlと同じディレクトリ）
  3. デフォルト値（`:-`の後に指定した値）

**`.env.local`、`.env.development`などのファイル**:
- これらのファイルは`.dockerignore`で除外されます（イメージに含まれない）
- Docker Composeは`.env.local`を**自動的には読み込みません**
- `.env.local`を読み込むには、`env_file`オプションで明示的に指定する必要があります:
  ```yaml
  env_file:
    - ./client/.env.local
  ```

**推奨される方法（採用）**:
**方法1: `env_file`オプションで`.env.local`を明示的に指定**
```yaml
env_file:
  - ./client/.env.local
```

**採用理由**:
- Not Docker版（通常の開発環境）と挙動を統一できる
- `.env.local`ファイルがそのまま使用できる（追加の設定ファイル不要）
- 環境変数の解決順序: `env_file`で読み込んだ値 → `environment`セクションの値（上書き可能）
- 既存の`.env.local`ファイルをそのまま活用できる

**実装方法**:
- `docker-compose.client.develop.yml`に`env_file`オプションを追加
- `build.args`では、ホスト側の環境変数または`.env`ファイルから値を取得（`${NEXT_PUBLIC_API_KEY:-}`）
- 実行時には`env_file`で`.env.local`を読み込み、`environment`セクションで必要に応じて上書き

**注意事項**:
- `.env.local`は`.dockerignore`で除外されるため、イメージには含まれない（セキュリティ上問題なし）
- `env_file`は実行時に読み込まれるため、ビルド時の`build.args`には影響しない
- ビルド時に`.env.local`の値を使う場合は、ホスト側の環境変数として設定するか、`.env`ファイルに値を定義する必要がある

**注意事項**:
- `.env.local`は通常、Gitにコミットしないファイルです（`.gitignore`に含まれる）
- 本番環境では、環境変数管理サービス（AWS Secrets Manager等）の使用を推奨

## 5. ボリュームマウント設計

### 5.1 設定ファイルのマウント
- **パス**: `./config/{env}:/app/config/{env}:ro`
- **モード**: 読み取り専用（`:ro`）
- **用途**: 環境別設定ファイルの読み込み

### 5.2 データディレクトリのマウント
- **パス**: `./server/data:/app/data`
- **モード**: 読み書き可能
- **用途**: SQLiteデータベースファイルの永続化

### 5.3 ログディレクトリのマウント
- **パス**: `./logs:/app/logs`
- **モード**: 読み書き可能
- **用途**: ログファイルの永続化
- **注意**: デフォルト設定でログファイルは`logs/`ディレクトリに出力されるため、ボリュームマウントが必要

### 5.3 クライアント開発モードのマウント
- **ソースコード**: `./client:/app`
- **node_modules**: `/app/node_modules`（ボリュームとして分離）
- **.next**: `/app/.next`（ボリュームとして分離）
- **用途**: ホットリロード対応

## 6. エラーハンドリング設計

### 6.1 ビルドエラー
- **問題**: Dockerfileの構文エラー、依存関係の不足
- **対処**: エラーメッセージを確認し、Dockerfileを修正

### 6.2 起動エラー
- **問題**: ポート競合、ネットワーク接続エラー、設定ファイルの読み込みエラー
- **対処**: 
  - ポートが使用されていないか確認
  - ネットワークが正しく設定されているか確認
  - 設定ファイルのパスが正しいか確認

### 6.3 実行時エラー
- **問題**: データベース接続エラー、Redis接続エラー
- **対処**: 
  - 既存のPostgreSQL、Redisコンテナが起動しているか確認
  - ネットワーク接続を確認
  - ログを確認（`docker-compose logs`）

## 7. セキュリティ設計

### 7.1 非rootユーザー実行
- すべてのコンテナで非rootユーザーで実行
- Goアプリケーション: `appuser`（UID 1000）
- Node.jsアプリケーション: `nodeuser`（UID 1000）

### 7.2 最小権限の原則
- 必要最小限のファイルのみをコピー
- 設定ファイルは読み取り専用でマウント
- 不要なツールやパッケージをインストールしない

### 7.3 シークレット管理
- 機密情報は環境変数で管理
- `.env`ファイルはGitにコミットしない
- 本番環境ではシークレット管理サービス（AWS Secrets Manager等）を使用

## 8. パフォーマンス最適化設計

### 8.1 マルチステージビルド
- ビルドステージと実行ステージを分離
- 最終イメージサイズを最小化
- ビルドツールを含まない軽量な実行イメージ

### 8.2 ビルドキャッシュの活用
- 依存関係のインストールを先に実行
- ソースコードの変更時のみ再ビルド
- `.dockerignore`で不要なファイルを除外

### 8.3 イメージサイズの最適化
- Alpine Linuxベースイメージを使用
- 不要なパッケージをインストールしない
- マルチステージビルドでビルドツールを除外

## 9. テスト戦略

### 9.1 Dockerfileのテスト
- Dockerイメージが正常にビルドされるか確認
- コンテナが正常に起動するか確認
- ヘルスチェックが正常に動作するか確認

### 9.2 Docker Composeのテスト
- 各docker-composeファイルが正常に動作するか確認
- サービス間の通信が正常に行われるか確認
- 環境変数が正しく設定されるか確認

### 9.3 統合テスト
- APIサーバーがPostgreSQL、Redisに接続できるか確認
- クライアントがAPIサーバーに接続できるか確認
- AdminサーバーがPostgreSQLに接続できるか確認

## 10. 運用・保守

### 10.1 ビルドコマンド
```bash
# APIサーバー
docker-compose -f docker-compose.api.develop.yml build

# Adminサーバー
docker-compose -f docker-compose.admin.develop.yml build

# クライアントサーバー
docker-compose -f docker-compose.client.develop.yml build
```

### 10.2 起動コマンド
```bash
# APIサーバー
docker-compose -f docker-compose.api.develop.yml up -d

# Adminサーバー
docker-compose -f docker-compose.admin.develop.yml up -d

# クライアントサーバー
docker-compose -f docker-compose.client.develop.yml up -d
```

### 10.3 停止コマンド
```bash
# APIサーバー
docker-compose -f docker-compose.api.develop.yml down

# Adminサーバー
docker-compose -f docker-compose.admin.develop.yml down

# クライアントサーバー
docker-compose -f docker-compose.client.develop.yml down
```

### 10.4 ログ確認
```bash
# APIサーバー
docker-compose -f docker-compose.api.develop.yml logs -f

# Adminサーバー
docker-compose -f docker-compose.admin.develop.yml logs -f

# クライアントサーバー
docker-compose -f docker-compose.client.develop.yml logs -f
```

## 11. 本番環境へのデプロイ準備

### 11.1 イメージのタグ付け
```bash
# APIサーバー
docker tag api:latest your-registry/api:v1.0.0

# Adminサーバー
docker tag admin:latest your-registry/admin:v1.0.0

# クライアントサーバー
docker tag client:latest your-registry/client:v1.0.0
```

### 11.2 コンテナレジストリへのプッシュ
```bash
# AWS ECR
aws ecr get-login-password --region ap-northeast-1 | docker login --username AWS --password-stdin your-account.dkr.ecr.ap-northeast-1.amazonaws.com
docker push your-account.dkr.ecr.ap-northeast-1.amazonaws.com/api:v1.0.0

# Tencent Cloud TCR
docker login your-registry.tencentcloudcr.com
docker push your-registry.tencentcloudcr.com/api:v1.0.0
```

### 11.3 本番環境用設定
- 環境変数を本番環境用に設定
- リソース制限（CPU・メモリ）を設定
- ヘルスチェックを実装
- ログ集約の設定（将来の拡張項目）

## 12. 実装上の注意事項

### 12.1 Dockerfileの作成
- マルチステージビルドを使用してイメージサイズを最小化
- 非rootユーザーで実行するように設定
- CGO_ENABLED=0で静的リンクビルド
- 適切な`.dockerignore`ファイルを作成

### 12.2 Docker Compose設定
- 環境別にdocker-composeファイルを分離
- 既存のネットワークを外部ネットワークとして参照
- ボリュームマウントを適切に設定
- 環境変数を各ファイルで適切に設定

### 12.3 既存サービスとの統合
- 既存のPostgreSQL、Redisコンテナと同一ネットワークで通信
- サービス名（`postgres`、`redis`）を正しく参照
- ネットワーク名を正しく参照

### 12.4 環境別制御
- `APP_ENV`環境変数で環境を切り替え
- 設定ファイルのパスを環境に応じて変更
- デフォルト値を適切に設定

### 12.6 ログファイルの永続化
- **ログディレクトリのマウント**: `./logs:/app/logs`をボリュームマウント
- **デフォルト設定**: サーバーアプリケーションのデフォルト設定でログファイルは`logs/`ディレクトリに出力される
- **永続化の必要性**: コンテナが削除されるとログファイルも消えるため、ボリュームマウントが必要
- **Dockerfileでの準備**: `/app/logs`ディレクトリを作成し、適切な権限を設定
- **適用範囲**: APIサーバーとAdminサーバーの両方に適用

### 12.5 Next.js環境変数の扱い（重要）
- **`.env*`ファイルは`.dockerignore`で除外**: `.env.local`、`.env.development`などのファイルは自動的に取り込まれません
- **ビルド時に必要な環境変数**: `NEXT_PUBLIC_*`で始まる環境変数は、docker-compose.ymlの`build.args`で明示的に指定する必要があります
- **実行時の環境変数**: 開発モードでは実行時にも環境変数を設定可能ですが、本番モードではビルド時に埋め込まれた値が使用されます
- **環境変数の受け渡し**:
  1. docker-compose.ymlの`build.args`でビルド時の環境変数を指定
  2. Dockerfileの`ARG`で受け取り、`ENV`で設定
  3. `npm run build`時にクライアント側のJavaScriptに埋め込まれる
  4. 実行時にはdocker-compose.ymlの`environment`セクションで指定（開発モードの場合）

## 13. 参考情報

### 13.1 関連ドキュメント
- 要件定義書: `.kiro/specs/0042-to-docker/requirements.md`
- プロジェクトREADME: `README.md`
- システムアーキテクチャ: `docs/Architecture.md`

### 13.2 技術スタック
- **Go**: 1.21+
- **Node.js**: 22+
- **Docker**: Docker Compose
- **ベースイメージ**: Alpine Linux

### 13.3 参考リンク
- Docker公式ドキュメント: https://docs.docker.com/
- Docker Compose公式ドキュメント: https://docs.docker.com/compose/
- Go公式Dockerイメージ: https://hub.docker.com/_/golang
- Node.js公式Dockerイメージ: https://hub.docker.com/_/node
- Alpine Linux公式サイト: https://alpinelinux.org/
