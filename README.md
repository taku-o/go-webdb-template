# Go DB Project Sample

大量ユーザー、大量アクセスの運用に耐えうるGo APIサーバー、データベース構成の
大規模プロジェクト向けのテンプレートプロジェクトです。

ナンバー付きDBテーブルでデータベース分割への配慮。
APIサーバー、Webサーバーを増やして負荷に耐えられるサーバー構成。

保守しやすいソースコード構成に、
レートリミット、各種ログ、データベースのスキーマ管理など各種運用向け機能。

**GitHub Pages**: [https://taku-o.github.io/go-webdb-template/pages/ja/](https://taku-o.github.io/go-webdb-template/pages/ja/)

## プロジェクト概要

- **サーバー**: Go言語、レイヤードアーキテクチャ、Database Sharding対応
- **クライアント**: Next.js 14 (App Router)、TypeScript
- **データベース**: PostgreSQL/MySQL（全環境）
- **テスト**: Go testing、Jest、Playwright

## 特徴

- ✅ **Sharding対応**: テーブルベースシャーディング（32分割）で複数DBへデータ分散
- ✅ **GORM対応**: Writer/Reader分離をサポート (GORM v1.25.12)
- ✅ **GoAdmin管理画面**: Webベースの管理画面でデータ管理
- ✅ **レイヤー分離**: API層、Usecase層、Service層、Repository層、DB層で責務を明確化
- ✅ **環境別設定**: develop/staging/production環境で設定切り替え
- ✅ **型安全**: TypeScriptによる型定義
- ✅ **テスト**: ユニット/統合/E2Eテスト対応
- ✅ **レートリミット**: IPアドレス単位でのAPI呼び出し制限（ulule/limiter使用）
- ✅ **ジョブキュー**: Redis + Asynqを使用したバックグラウンドジョブ処理
- ✅ **メール送信**: 標準出力、Mailpit、AWS SES対応のメール送信機能
- ✅ **ファイルアップロード**: TUSプロトコルによる大容量ファイルアップロード（ローカル/S3ストレージ対応）
- ✅ **ログ機能**: アクセスログ、メール送信ログ、SQLログの出力
- ✅ **Docker対応**: APIサーバー、Adminサーバー、クライアントサーバーのDocker化

## セットアップ

### 前提条件

- Go 1.21+
- Node.js 18+
- Docker（PostgreSQLコンテナ用）
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

本プロジェクトではPostgreSQLを使用し、[Atlas](https://atlasgo.io/) でマイグレーションを管理しています。

#### PostgreSQLの起動

```bash
# PostgreSQLコンテナを起動（master + sharding 4台）
./scripts/start-postgres.sh start
```

**接続情報**（開発環境）:

| データベース | ホスト | ポート | ユーザー | パスワード | データベース名 |
|------------|--------|--------|---------|-----------|--------------|
| Master | localhost | 5432 | webdb | webdb | webdb_master |
| Sharding 1 | localhost | 5433 | webdb | webdb | webdb_sharding_1 |
| Sharding 2 | localhost | 5434 | webdb | webdb | webdb_sharding_2 |
| Sharding 3 | localhost | 5435 | webdb | webdb | webdb_sharding_3 |
| Sharding 4 | localhost | 5436 | webdb | webdb | webdb_sharding_4 |

#### マイグレーションの適用

```bash
# 全データベースにマイグレーションを適用（初期データも含む）
./scripts/migrate.sh all
```

#### PostgreSQLの停止

```bash
./scripts/start-postgres.sh stop
```

#### MySQLの起動（オプション）

PostgreSQLの代わりにMySQLを使用することもできます。

```bash
# MySQLコンテナを起動（master + sharding 4台）
./scripts/start-mysql.sh start
```

**接続情報**（開発環境）:

| データベース | ホスト | ポート | ユーザー | パスワード | データベース名 |
|------------|--------|--------|---------|-----------|--------------|
| Master | localhost | 3306 | webdb | webdb | webdb_master |
| Sharding 1 | localhost | 3307 | webdb | webdb | webdb_sharding_1 |
| Sharding 2 | localhost | 3308 | webdb | webdb | webdb_sharding_2 |
| Sharding 3 | localhost | 3309 | webdb | webdb | webdb_sharding_3 |
| Sharding 4 | localhost | 3310 | webdb | webdb | webdb_sharding_4 |

#### MySQLマイグレーションの適用

```bash
# 開発環境用マイグレーション
atlas migrate apply --dir "file://db/migrations/master-mysql" --url "mysql://webdb:webdb@localhost:3306/webdb_master"
atlas migrate apply --dir "file://db/migrations/sharding_1-mysql" --url "mysql://webdb:webdb@localhost:3307/webdb_sharding_1"
atlas migrate apply --dir "file://db/migrations/sharding_2-mysql" --url "mysql://webdb:webdb@localhost:3308/webdb_sharding_2"
atlas migrate apply --dir "file://db/migrations/sharding_3-mysql" --url "mysql://webdb:webdb@localhost:3309/webdb_sharding_3"
atlas migrate apply --dir "file://db/migrations/sharding_4-mysql" --url "mysql://webdb:webdb@localhost:3310/webdb_sharding_4"

# テスト環境用マイグレーション
./scripts/migrate-test-mysql.sh
```

#### MySQLの停止

```bash
./scripts/start-mysql.sh stop
```

#### データベースタイプの切り替え

`config/{env}/config.yaml`の`DB_TYPE`でデータベースを切り替えます：

```yaml
# PostgreSQLを使用（デフォルト）
DB_TYPE: postgresql

# MySQLを使用
DB_TYPE: mysql
```

MySQLを使用する場合は`database.mysql.yaml`から設定が読み込まれます。

#### シャーディング構成

本プロジェクトでは**8つの論理シャード**を**4つの物理データベース**に分散配置しています：

| 論理シャードID | テーブル範囲 | 物理データベース |
|--------------|-------------|----------------|
| 1 | _000 〜 _003 | webdb_sharding_1 (port 5433) |
| 2 | _004 〜 _007 | webdb_sharding_1 (port 5433) |
| 3 | _008 〜 _011 | webdb_sharding_2 (port 5434) |
| 4 | _012 〜 _015 | webdb_sharding_2 (port 5434) |
| 5 | _016 〜 _019 | webdb_sharding_3 (port 5435) |
| 6 | _020 〜 _023 | webdb_sharding_3 (port 5435) |
| 7 | _024 〜 _027 | webdb_sharding_4 (port 5436) |
| 8 | _028 〜 _031 | webdb_sharding_4 (port 5436) |

#### スキーマ変更時のマイグレーション生成

```bash
# master.hclを変更した後
atlas migrate diff <migration_name> \
    --dir file://db/migrations/master \
    --to file://db/schema/master.hcl \
    --dev-url "postgres://webdb:webdb@localhost:5432/webdb_master?sslmode=disable"

# sharding.hclを変更した後
atlas migrate diff <migration_name> \
    --dir file://db/migrations/sharding \
    --to file://db/schema/sharding.hcl \
    --dev-url "postgres://webdb:webdb@localhost:5433/webdb_sharding_1?sslmode=disable"
```

詳細は [docs/Atlas-Operations.md](docs/Atlas-Operations.md) を参照してください。

#### 遅延接続・自動再接続機能

本プロジェクトでは以下の機能を実装しています：

- **遅延接続**: サーバー起動時にDB接続を行わず、最初のクエリ実行時に接続を確立
- **自動再接続**: データベースが復旧した際に自動的に再接続
- **リトライ機能**: 接続エラー時に最大3回、1秒間隔でリトライ

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

**データベース接続**:
GoAdminサーバーはPostgreSQLを利用します。設定ファイル（`config/{env}/database.yaml`）から接続情報を読み込み、masterデータベース（`webdb_master`）に接続します。

起動前に以下を確認してください：
1. PostgreSQLコンテナが起動していること（`./scripts/start-postgres.sh start`）
2. マイグレーションが適用されていること（`./scripts/migrate.sh all`）

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

#### 依存関係のインストール

```bash
cd client
npm install --legacy-peer-deps
```

**注意**: peer dependencyの競合がある場合は`--legacy-peer-deps`フラグを使用してください。

#### 環境変数の設定

**AUTH_SECRETの生成**:
```bash
# プロジェクトルートで実行
npm run cli:generate-secret
```
このコマンドで生成された秘密鍵をコピーします。

`.env.local`を作成して以下の環境変数を設定：
```
# NextAuth (Auth.js)
AUTH_SECRET=<npm run cli:generate-secretで生成した秘密鍵>
AUTH_URL=http://localhost:3000

# Auth0設定
AUTH0_ISSUER=https://your-tenant.auth0.com
AUTH0_CLIENT_ID=your-client-id
AUTH0_CLIENT_SECRET=your-client-secret
AUTH0_AUDIENCE=https://your-api-audience

# API設定
NEXT_PUBLIC_API_KEY=your-api-key
NEXT_PUBLIC_API_BASE_URL=http://localhost:8080

# テスト環境用（テスト実行時に必要）
APP_ENV=test
```

**注意**: 
- `AUTH_SECRET`は`npm run cli:generate-secret`コマンドで生成します（`server/cmd/generate-secret`を使用）。
- `APP_ENV=test`はテスト実行時に必要です（`npm test`、`npm run e2e`実行時）。

#### Auth0アプリケーション設定

Auth0ダッシュボード（`Applications > [対象アプリ] > Settings`）で以下のURLを設定：

**Allowed Callback URLs:**
```
http://localhost:3000/api/auth/callback/auth0
```

**Allowed Logout URLs:**
```
http://localhost:3000
```

**Allowed Web Origins:**
```
http://localhost:3000
```

#### 開発サーバーの起動

```bash
cd client
npm run dev
```

クライアントは http://localhost:3000 で起動します。

#### 技術スタック

- **フレームワーク**: Next.js 14+ (App Router)
- **言語**: TypeScript 5+
- **UIコンポーネント**: shadcn/ui
- **認証**: NextAuth (Auth.js) v5
- **スタイリング**: Tailwind CSS
- **フォーム管理**: react-hook-form
- **バリデーション**: zod
- **ファイルアップロード**: Uppy (TUSプロトコル)
- **テスト**: Playwright (E2E), Jest (単体・統合), MSW (APIモック)

#### 利用可能なスクリプト

**開発**:
- `npm run dev` - 開発サーバーを起動（ポート3000）
- `npm run build` - プロダクションビルドを実行
- `npm run start` - プロダクションビルドを起動（ポート3000）
- `npm run lint` - ESLintを実行
- `npm run type-check` - TypeScript型チェックを実行
- `npm run format` - Prettierでフォーマットを確認
- `npm run format:write` - Prettierでフォーマットを適用

**テスト**:
- `npm test` - Jestテストを実行（単体・統合テスト）
- `npm run test:watch` - Jestテストをウォッチモードで実行
- `npm run test:coverage` - Jestテストのカバレッジを取得
- `npm run e2e` - Playwright E2Eテストを実行
- `npm run e2e:ui` - Playwright E2EテストをUIモードで実行
- `npm run e2e:headed` - Playwright E2Eテストをヘッドモードで実行

**注意**: テスト実行時は`APP_ENV=test`が自動的に設定されます（`package.json`のスクリプトに含まれています）。

### 10. Docker環境での起動（オプション）

Docker環境でサーバーを起動することもできます。

```bash
# PostgreSQL、Redisコンテナの起動（先に起動が必要）
docker-compose -f docker-compose.postgres.yml up -d
docker-compose -f docker-compose.redis.yml up -d

# APIサーバーのビルドと起動
docker-compose -f docker-compose.api.yml build
docker-compose -f docker-compose.api.yml up -d

# Adminサーバーのビルドと起動
docker-compose -f docker-compose.admin.yml build
docker-compose -f docker-compose.admin.yml up -d

# クライアントサーバーのビルドと起動
docker-compose -f docker-compose.client.yml build
docker-compose -f docker-compose.client.yml up -d
```

**起動後のアクセス先**:
- APIサーバー: http://localhost:8080
- Adminサーバー: http://localhost:8081/admin
- クライアント: http://localhost:3000

詳細は [Docker.md](docs/Docker.md) を参照してください。

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
- [Docker](docs/Docker.md) - Docker環境での起動・デプロイ

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

詳細なクライアントアプリのセットアップ手順は、セクション9「クライアント起動」を参照してください。

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

### サンプルデータ生成（generate-sample-data）

開発用のサンプルデータを生成します。PostgreSQLを使用してmaster/shardingデータベースにデータを投入します。

#### 前提条件

1. PostgreSQLコンテナが起動していること
2. マイグレーションが適用されていること

```bash
# PostgreSQL起動
./scripts/start-postgres.sh start

# マイグレーション適用
./scripts/migrate.sh
```

#### 実行

```bash
cd server
APP_ENV=develop go run cmd/generate-sample-data/main.go
```

#### 生成されるデータ

| テーブル | データベース | 件数 |
|---------|------------|------|
| dm_users_000〜031 | sharding (4台に分散) | 100件 |
| dm_posts_000〜031 | sharding (4台に分散) | 100件 |
| dm_news | master | 100件 |

詳細は [Generate-Sample-Data.md](docs/Generate-Sample-Data.md) を参照してください。

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

テーブルベースシャーディング（32分割、8論理シャード）を採用しています。

```
table_number = id % 32      # 0-31
table_name = "dm_users_" + sprintf("%03d", table_number)  # dm_users_000 ~ dm_users_031
logical_shard_id = (table_number / 4) + 1  # 1-8
```

**データベースグループ**:
- **Master グループ**: 共有テーブル（dm_news）を格納（PostgreSQL、port 5432）
- **Sharding グループ**: 32分割されたテーブル（dm_users_000〜031, dm_posts_000〜031）を8論理シャード→4物理DBに分散

| 論理シャード | テーブル範囲 | 物理DB（PostgreSQL） |
|------------|------------|---------------------|
| 1 | _000 〜 _003 | webdb_sharding_1 (port 5433) |
| 2 | _004 〜 _007 | webdb_sharding_1 (port 5433) |
| 3 | _008 〜 _011 | webdb_sharding_2 (port 5434) |
| 4 | _012 〜 _015 | webdb_sharding_2 (port 5434) |
| 5 | _016 〜 _019 | webdb_sharding_3 (port 5435) |
| 6 | _020 〜 _023 | webdb_sharding_3 (port 5435) |
| 7 | _024 〜 _027 | webdb_sharding_4 (port 5436) |
| 8 | _028 〜 _031 | webdb_sharding_4 (port 5436) |

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
