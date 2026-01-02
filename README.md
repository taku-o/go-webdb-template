# Go DB Project Sample

Go + Next.js + Database Sharding対応のサンプルプロジェクトです。大規模プロジェクト向けの構成を採用しています。

## プロジェクト概要

- **サーバー**: Go言語、レイヤードアーキテクチャ、Database Sharding対応
- **クライアント**: Next.js 14 (App Router)、TypeScript
- **データベース**: SQLite (開発環境)、PostgreSQL/MySQL (本番想定)
- **テスト**: Go testing、Jest、Playwright

## 特徴

- ✅ **Sharding対応**: テーブルベースシャーディング（32分割）で複数DBへデータ分散
- ✅ **GORM対応**: Writer/Reader分離をサポート (GORM v1.25.12)
- ✅ **GoAdmin管理画面**: Webベースの管理画面でデータ管理
- ✅ **レイヤー分離**: API層、Service層、Repository層、DB層で責務を明確化
- ✅ **環境別設定**: develop/staging/production環境で設定切り替え
- ✅ **型安全**: TypeScriptによる型定義
- ✅ **テスト**: ユニット/統合/E2Eテスト対応
- ✅ **レートリミット**: IPアドレス単位でのAPI呼び出し制限（ulule/limiter使用）
- ✅ **ジョブキュー**: Redis + Asynqを使用したバックグラウンドジョブ処理
- ✅ **メール送信**: 標準出力、Mailpit、AWS SES対応のメール送信機能
- ✅ **ファイルアップロード**: TUSプロトコルによる大容量ファイルアップロード（ローカル/S3ストレージ対応）
- ✅ **ログ機能**: アクセスログ、メール送信ログ、SQLログの出力

## セットアップ

### 前提条件

- Go 1.21+
- Node.js 18+
- SQLite3
- Atlas CLI（データベースマイグレーション管理用）
  - インストール方法: `brew install ariga/tap/atlas`（macOS）
  - インストール確認: `atlas version`
  - 詳細: https://atlasgo.io/
- Redis（ジョブキュー機能を使用する場合、オプション）
  - Dockerを使用して起動可能（`./scripts/start-redis.sh`）

### 1. 依存関係のインストール

#### サーバー側
```bash
cd server
go mod download
```

#### クライアント側
```bash
cd client
npm install
```

### 2. データベースのセットアップ

本プロジェクトでは [Atlas](https://atlasgo.io/) を使用してデータベースマイグレーションを管理しています。

#### マイグレーションスクリプトを使用

```bash
# 全データベースにマイグレーションを適用（初期データも含む）
./scripts/migrate.sh all
```

#### 手動でAtlasコマンドを使用

```bash
mkdir -p server/data

# マスターDBにマイグレーションを適用（初期データも含む）
atlas migrate apply \
    --dir file://db/migrations/master \
    --url "sqlite://server/data/master.db"

# シャーディングDBにマイグレーションを適用
for i in 1 2 3 4; do
    atlas migrate apply \
        --dir file://db/migrations/sharding \
        --url "sqlite://server/data/sharding_db_${i}.db"
done
```

#### スキーマ変更時のマイグレーション生成

```bash
# master.hclを変更した後
atlas migrate diff <migration_name> \
    --dir file://db/migrations/master \
    --to file://db/schema/master.hcl \
    --dev-url "sqlite://file?mode=memory"

# sharding.hclを変更した後
atlas migrate diff <migration_name> \
    --dir file://db/migrations/sharding \
    --to file://db/schema/sharding.hcl \
    --dev-url "sqlite://file?mode=memory"
```

詳細は [docs/Atlas-Operations.md](docs/Atlas-Operations.md) を参照してください。

### PostgreSQL環境（動作確認用）

開発環境でPostgreSQLを使用して遅延接続・自動再接続機能を確認できます。

#### PostgreSQLの起動・停止

```bash
# PostgreSQLを起動
./scripts/start-postgres.sh start

# PostgreSQLを停止
./scripts/start-postgres.sh stop
```

PostgreSQLは http://localhost:5432 で起動します。

**接続情報**（開発環境）:
- ホスト: `localhost`
- ポート: `5432`
- ユーザー名: `webdb`
- パスワード: `webdb`
- データベース: `webdb`

#### 設定ファイルの切り替え

`config/develop/database.yaml`でデータベースドライバを切り替えます：

```yaml
# PostgreSQL設定
database:
  groups:
    master:
      - id: 1
        driver: postgres
        host: localhost
        port: 5432
        user: webdb
        password: webdb
        name: webdb
        max_connections: 25
        max_idle_connections: 5
        connection_max_lifetime: 1h
```

#### SQLite環境との併用

- **開発時（通常）**: SQLite設定を使用（設定不要）
- **動作確認時**: PostgreSQL設定に切り替えて使用
- 設定ファイル内でSQLite/PostgreSQLの設定をコメントアウトで切り替え可能

#### 遅延接続・自動再接続機能

本プロジェクトでは以下の機能を実装しています：

- **遅延接続**: サーバー起動時にDB接続を行わず、最初のクエリ実行時に接続を確立
- **自動再接続**: データベースが復旧した際に自動的に再接続
- **リトライ機能**: 接続エラー時に最大3回、1秒間隔でリトライ

これらの機能はPostgreSQL環境で動作確認できます。

### 3. サーバー起動

```bash
cd server
APP_ENV=develop go run cmd/server/main.go
```

サーバーは http://localhost:8080 で起動します。

### 4. 管理画面起動

```bash
cd server
APP_ENV=develop go run cmd/admin/main.go
```

管理画面は http://localhost:8081/admin で起動します。

**認証情報**（開発環境）:
- ユーザー名: `admin`
- パスワード: `admin123`

詳細は [Admin.md](docs/Admin.md) を参照してください。

### 5. データベースビューア起動（CloudBeaver）

```bash
# 開発環境（デフォルト）
npm run cloudbeaver:start

# 環境を指定して起動
APP_ENV=develop npm run cloudbeaver:start
APP_ENV=staging npm run cloudbeaver:start
APP_ENV=production npm run cloudbeaver:start

# 停止
npm run cloudbeaver:stop
```

データベースビューアは http://localhost:8978 で起動します。

**認証情報**（開発環境）:
- ユーザー名: `cbadmin`
- パスワード: `Admin123`

**主な機能**:
- Webブラウザからデータベースを操作
- テーブル構造の確認・データ閲覧
- SQLクエリの実行
- SQLスクリプトの保存・管理（Resource Manager）

詳細は [Database-Viewer.md](docs/Database-Viewer.md) を参照してください。

### 6. データ可視化ツール起動（Metabase）

```bash
# 開発環境（デフォルト）
npm run metabase:start

# 環境を指定して起動
APP_ENV=develop npm run metabase:start
APP_ENV=staging npm run metabase:start
APP_ENV=production npm run metabase:start

# 停止
npm run metabase:stop
```

Metabaseは http://localhost:8970 で起動します。

**主な機能**:
- データの可視化・グラフ作成
- ダッシュボードの作成・共有
- 非エンジニア向けのデータ分析

**CloudBeaverとMetabaseの使い分け**:
- **CloudBeaver**: データの直接編集・操作、テーブル構造の確認
- **Metabase**: データの可視化・分析、ダッシュボード作成

**注意**: CloudBeaverとMetabaseはメモリ使用量が大きいため、開発環境では片方ずつしか起動しない運用を推奨します。

詳細は [Metabase.md](docs/Metabase.md) を参照してください。

### 7. Redisの起動（ジョブキュー機能用）

```bash
# Redisを起動
./scripts/start-redis.sh start

# Redis Insightを起動（オプション、データビューワ）
./scripts/start-redis-insight.sh start
```

Redisは http://localhost:6379 で起動します。
Redis Insightは http://localhost:8001 で起動します。

詳細は [Queue-Job.md](docs/Queue-Job.md) を参照してください。

### 8. Mailpitの起動（メール送信機能用、オプション）

```bash
./scripts/start-mailpit.sh start
```

Mailpitは http://localhost:8025 で起動します。

詳細は [Send-Mail.md](docs/Send-Mail.md) を参照してください。

### 9. クライアント起動

```bash
cd client
npm run dev
```

クライアントは http://localhost:3000 で起動します。

## API エンドポイント

### 基本エンドポイント

#### ユーザー関連

- `GET /api/dm-users` - ユーザー一覧取得
- `GET /api/dm-users/{id}` - ユーザー取得
- `POST /api/dm-users` - ユーザー作成
- `PUT /api/dm-users/{id}` - ユーザー更新
- `DELETE /api/dm-users/{id}` - ユーザー削除
- `GET /api/export/dm-users/csv` - ユーザー情報をCSV形式でダウンロード

#### 投稿関連

- `GET /api/dm-posts` - 投稿一覧取得
- `GET /api/dm-posts/{id}` - 投稿取得
- `POST /api/dm-posts` - 投稿作成
- `PUT /api/dm-posts/{id}` - 投稿更新
- `DELETE /api/dm-posts/{id}` - 投稿削除
- `GET /api/dm-user-posts` - ユーザーと投稿をJOIN（クロスシャードクエリ）

#### その他

- `GET /api/today` - 今日の日付取得（private API、Auth0 JWT必須）
- `GET /health` - ヘルスチェック（認証不要）

### 機能別エンドポイント

- `POST /api/email/send` - メール送信
- `POST /api/dm-jobqueue/register` - ジョブ登録
- `POST /api/upload/dm_movie` - ファイルアップロード（TUSプロトコル）

### OpenAPI仕様

- `GET /docs` - API Documentation UI (Stoplight Elements)
- `GET /openapi.json` - OpenAPI 3.1 (JSON)
- `GET /openapi.yaml` - OpenAPI 3.1 (YAML)
- `GET /openapi-3.0.json` - OpenAPI 3.0.3 (JSON)

※ OpenAPIドキュメントエンドポイントは認証不要でアクセス可能です。

詳細は [API.md](docs/API.md) を参照してください。

## 機能別ドキュメント

以下の機能の詳細な利用手順は、各ドキュメントを参照してください：

- [ジョブキュー機能](docs/Queue-Job.md) - Redis + Asynqを使用したバックグラウンドジョブ処理
- [メール送信機能](docs/Send-Mail.md) - 標準出力、Mailpit、AWS SES対応のメール送信
- [ファイルアップロード機能](docs/File-Upload.md) - TUSプロトコルによる大容量ファイルアップロード
- [ログ機能](docs/Logging.md) - アクセスログ、メール送信ログ、SQLログ
- [レートリミット機能](docs/Rate-Limit.md) - APIレートリミットの詳細設定

## APIレートリミット

APIエンドポイントへのリクエストはIPアドレス単位でレート制限されています。

### レスポンスヘッダー

すべてのAPIレスポンスに以下のヘッダーが付与されます：

**分制限（常に付与）:**

| ヘッダー | 説明 | 例 |
|---------|------|-----|
| `X-RateLimit-Limit` | 1分あたりの制限値 | `60` |
| `X-RateLimit-Remaining` | 残りリクエスト数 | `45` |
| `X-RateLimit-Reset` | リセット時刻（Unix timestamp） | `1706342400` |

**時間制限（`requests_per_hour`が設定されている場合のみ）:**

| ヘッダー | 説明 | 例 |
|---------|------|-----|
| `X-RateLimit-Hour-Limit` | 1時間あたりの制限値 | `1000` |
| `X-RateLimit-Hour-Remaining` | 残りリクエスト数 | `950` |
| `X-RateLimit-Hour-Reset` | リセット時刻（Unix timestamp） | `1706346000` |

### レートリミット超過時

制限を超過した場合、HTTP 429ステータスコードが返されます：

```json
{
  "code": 429,
  "message": "Too Many Requests"
}
```

### 設定

レートリミットの設定は`config/{env}/config.yaml`で管理します：

```yaml
api:
  rate_limit:
    enabled: true
    requests_per_minute: 60
    requests_per_hour: 1000
```

### ストレージ

レートリミットのカウンターは環境に応じて異なるストレージを使用します：

| 環境 | ストレージ | 設定ファイル |
|------|----------|-------------|
| develop | In-Memory | `config/develop/cacheserver.yaml` |
| staging | Redis Cluster | `config/staging/cacheserver.yaml` |
| production | Redis Cluster | `config/production/cacheserver.yaml` |

Redis Clusterを使用する場合は`cacheserver.yaml`でアドレスを設定します：

```yaml
redis:
  cluster:
    addrs:
      - host1:6379
      - host2:6379
      - host3:6379
```

`addrs`が空または未設定の場合はIn-Memoryストレージが使用されます。

### 動作確認

```bash
# レートリミットヘッダーの確認
curl -i -H "Authorization: Bearer <YOUR_API_KEY>" http://localhost:8080/api/users

# レスポンスヘッダー例
# X-RateLimit-Limit: 60
# X-RateLimit-Remaining: 59
# X-RateLimit-Reset: 1706342460
```

詳細は [Rate-Limit.md](docs/Rate-Limit.md) を参照してください。

## API認証

APIエンドポイント（`/api/*`）へのアクセスにはJWTベースのPublic APIキーが必要です。

### APIキーの発行

GoAdmin管理画面からAPIキーを発行できます。

1. 管理画面（http://localhost:8081/admin）にログイン
2. サイドメニューから「カスタムページ」→「APIキー発行」を選択
3. 「APIキーを発行」ボタンをクリック
4. 生成されたJWTトークンをダウンロードまたはコピー

### 秘密鍵の生成

APIキーの署名に使用する秘密鍵を生成するツールが用意されています。

```bash
cd server
go run cmd/generate-secret/main.go
```

生成された秘密鍵を`config/{env}/config.yaml`の`api.secret_key`に設定してください。

### APIリクエストの認証

APIリクエストには`Authorization`ヘッダーでJWTトークンを送信します。

```bash
curl -H "Authorization: Bearer <YOUR_API_KEY>" http://localhost:8080/api/users
```

### クライアント側の設定

Next.jsクライアントでは、環境変数`NEXT_PUBLIC_API_KEY`にAPIキーを設定します。

```bash
# client/.env.local
NEXT_PUBLIC_API_KEY=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

### エラーレスポンス

- `401 Unauthorized` - APIキーが無効または未設定
- `403 Forbidden` - スコープ不足（readスコープなしでGET、writeスコープなしでPOST/PUT/DELETE）

エラーレスポンス形式:
```json
{
  "code": 401,
  "message": "Invalid API key"
}
```

## CLIツール

バッチ処理用のCLIツールが利用できます。

### 秘密鍵生成（generate-secret）

APIキー署名用の秘密鍵を生成します。

#### 実行

```bash
cd server
go run cmd/generate-secret/main.go
```

#### 出力

Base64エンコードされた32バイト（256ビット）のランダムな秘密鍵が標準出力に表示されます。

### ユーザー一覧出力（list-users）

ユーザー一覧をTSV形式で出力します。

#### ビルド

```bash
cd server
go build -o bin/list-users ./cmd/list-users
```

#### 実行

```bash
# デフォルト（20件）
APP_ENV=develop ./bin/list-users

# 件数を指定（最大100件）
APP_ENV=develop ./bin/list-users --limit 50
```

#### オプション

| オプション | 説明 | デフォルト | 範囲 |
|-----------|------|----------|------|
| `--limit` | 出力件数 | 20 | 1-100 |

#### 出力形式

TSV（タブ区切り）形式で、以下の項目を出力します：
- ID, Name, Email, CreatedAt, UpdatedAt

## Sharding戦略

テーブルベースシャーディング（32分割）を採用しています。

```
table_number = id % 32      # 0-31
table_name = "users_" + sprintf("%03d", table_number)  # users_000 ~ users_031
db_id = (table_number / 8) + 1  # 1-4
```

**データベースグループ**:
- **Master グループ**: 共有テーブル（news）を格納
- **Sharding グループ**: 32分割されたテーブル（users_000〜031, posts_000〜031）を4つのDBに分散

| DB | テーブル範囲 |
|----|------------|
| DB1 | _000 〜 _007 |
| DB2 | _008 〜 _015 |
| DB3 | _016 〜 _023 |
| DB4 | _024 〜 _031 |

詳細は [Sharding.md](docs/Sharding.md) を参照してください。

## 設定ファイル構造

設定ファイルは環境別ディレクトリに分割されています。

### ディレクトリ構造

```
config/
├── develop/                  # 開発環境設定ディレクトリ
│   ├── config.yaml           # メイン設定（server, admin, logging, cors, api）
│   ├── database.yaml         # データベース設定（groups構造）
│   └── cacheserver.yaml      # キャッシュサーバー設定（Redis Cluster）
├── production/               # 本番環境設定ディレクトリ
│   ├── config.yaml.example   # メイン設定テンプレート
│   ├── database.yaml.example # データベース設定テンプレート
│   └── cacheserver.yaml.example # キャッシュサーバー設定テンプレート
└── staging/                  # ステージング環境設定ディレクトリ
    ├── config.yaml           # メイン設定
    ├── database.yaml         # データベース設定
    └── cacheserver.yaml      # キャッシュサーバー設定

db/
└── migrations/
    ├── master/               # Masterグループ用マイグレーション
    │   └── 001_init.sql      # newsテーブル
    └── sharding/             # Shardingグループ用マイグレーション
        ├── templates/        # テンプレートファイル
        │   ├── users.sql.template
        │   └── posts.sql.template
        └── generated/        # 生成されたマイグレーション
```

### 設定ファイルの読み込み順序

1. メイン設定ファイル（`config/{env}/config.yaml`）を読み込み
2. データベース設定ファイル（`config/{env}/database.yaml`）をマージ
3. キャッシュサーバー設定ファイル（`config/{env}/cacheserver.yaml`）をマージ（オプション）
4. 統合された設定を`Config`構造体にマッピング
5. 環境変数（`DB_PASSWORD_SHARD*`）でパスワードを上書き

### 環境切り替え

環境変数`APP_ENV`で環境を切り替えます：

```bash
APP_ENV=develop go run cmd/server/main.go    # 開発環境
APP_ENV=staging go run cmd/server/main.go    # ステージング環境
APP_ENV=production go run cmd/server/main.go # 本番環境
```

## GORM対応

Writer/Reader分離をサポートするGORM版のRepositoryを実装しています。

### Writer/Reader分離の設定例

`config/production/database.yaml`:
```yaml
database:
  shards:
    - id: 1
      driver: postgres
      writer_dsn: host=writer.example.com port=5432 user=app password=xxx dbname=db sslmode=require
      reader_dsns:
        - host=reader1.example.com port=5432 user=app password=xxx dbname=db sslmode=require
        - host=reader2.example.com port=5432 user=app password=xxx dbname=db sslmode=require
      reader_policy: round_robin
```

### 主要な依存パッケージ

- `gorm.io/gorm` v1.25.12
- `gorm.io/driver/sqlite`
- `gorm.io/driver/postgres`
- `gorm.io/plugin/dbresolver` (Writer/Reader分離)
- `gorm.io/sharding` (将来使用予定)
- `github.com/labstack/echo/v4` v4.13.3 (HTTPルーター)
- `github.com/danielgtaylor/huma/v2` v2.34.1 (OpenAPI仕様自動生成)
- `github.com/ulule/limiter/v3` v3.11.2 (レートリミット)
- `github.com/redis/go-redis/v9` v9.17.2 (Redis Cluster接続)

詳細は [Architecture.md](docs/Architecture.md) を参照してください。

## ライセンス

MIT License
