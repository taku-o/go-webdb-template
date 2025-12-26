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

## セットアップ

### 前提条件

- Go 1.21+
- Node.js 18+
- SQLite3
- Atlas CLI（データベースマイグレーション管理用）
  - インストール方法: `brew install ariga/tap/atlas`（macOS）
  - インストール確認: `atlas version`
  - 詳細: https://atlasgo.io/

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

#### 方法A: マイグレーションスクリプトを使用（推奨）

```bash
./scripts/migrate.sh
```

#### 方法B: 手動セットアップ

```bash
mkdir -p server/data

# Master データベース（news テーブル）
sqlite3 server/data/master.db < db/migrations/master/001_init.sql

# Sharding データベース（マイグレーション生成＆適用）
cd server
go run cmd/migrate-gen/main.go \
    -template ../db/migrations/sharding/templates/users.sql.template \
    -output ../db/migrations/sharding/generated/
go run cmd/migrate-gen/main.go \
    -template ../db/migrations/sharding/templates/posts.sql.template \
    -output ../db/migrations/sharding/generated/
cd ..

# DB1: テーブル _000-007
for i in {0..7}; do
    sqlite3 server/data/sharding_db_1.db < db/migrations/sharding/generated/users_$(printf "%03d" $i).sql
    sqlite3 server/data/sharding_db_1.db < db/migrations/sharding/generated/posts_$(printf "%03d" $i).sql
done

# DB2: テーブル _008-015
for i in {8..15}; do
    sqlite3 server/data/sharding_db_2.db < db/migrations/sharding/generated/users_$(printf "%03d" $i).sql
    sqlite3 server/data/sharding_db_2.db < db/migrations/sharding/generated/posts_$(printf "%03d" $i).sql
done

# DB3: テーブル _016-023
for i in {16..23}; do
    sqlite3 server/data/sharding_db_3.db < db/migrations/sharding/generated/users_$(printf "%03d" $i).sql
    sqlite3 server/data/sharding_db_3.db < db/migrations/sharding/generated/posts_$(printf "%03d" $i).sql
done

# DB4: テーブル _024-031
for i in {24..31}; do
    sqlite3 server/data/sharding_db_4.db < db/migrations/sharding/generated/users_$(printf "%03d" $i).sql
    sqlite3 server/data/sharding_db_4.db < db/migrations/sharding/generated/posts_$(printf "%03d" $i).sql
done
```

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

### 5. クライアント起動

```bash
cd client
npm run dev
```

クライアントは http://localhost:3000 で起動します。

## API エンドポイント

- `GET /api/users` - ユーザー一覧
- `POST /api/users` - ユーザー作成
- `GET /api/posts` - 投稿一覧
- `POST /api/posts` - 投稿作成
- `GET /api/user-posts` - ユーザーと投稿をJOIN（クロスシャードクエリ）
- `GET /health` - ヘルスチェック（認証不要）

### OpenAPI仕様

- `GET /docs` - API Documentation UI (Stoplight Elements)
- `GET /openapi.json` - OpenAPI 3.1 (JSON)
- `GET /openapi.yaml` - OpenAPI 3.1 (YAML)
- `GET /openapi-3.0.json` - OpenAPI 3.0.3 (JSON)

※ OpenAPIドキュメントエンドポイントは認証不要でアクセス可能です。

詳細は [API.md](docs/API.md) を参照してください。

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
│   ├── config.yaml           # メイン設定（server, admin, logging, cors）
│   └── database.yaml         # データベース設定（groups構造）
├── production/               # 本番環境設定ディレクトリ
│   ├── config.yaml.example   # メイン設定テンプレート
│   └── database.yaml.example # データベース設定テンプレート
└── staging/                  # ステージング環境設定ディレクトリ
    ├── config.yaml           # メイン設定
    └── database.yaml         # データベース設定

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
3. 統合された設定を`Config`構造体にマッピング
4. 環境変数（`DB_PASSWORD_SHARD*`）でパスワードを上書き

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

詳細は [Architecture.md](docs/Architecture.md) を参照してください。

## 静的ファイル（CSS・画像）の配置

クライアント側（Next.js）の静的ファイルは`client/public/`ディレクトリに配置します。

### ディレクトリ構造

```
client/
├── public/
│   ├── css/              # CSSファイル用ディレクトリ
│   │   └── style.css     # サンプルCSSファイル
│   └── images/           # 画像ファイル用ディレクトリ
│       ├── logo.svg      # サンプルSVG画像
│       ├── logo.png      # サンプルPNG画像
│       └── icon.jpg      # サンプルJPG画像
```

### 参照方法

Next.jsの`public/`ディレクトリ配下のファイルは、ルート（`/`）から直接参照できます。

#### CSSファイルの参照

`client/src/app/layout.tsx`でCSSファイルを参照:

```tsx
<html lang="en">
  <head>
    <link rel="stylesheet" href="/css/style.css" />
  </head>
  <body>{children}</body>
</html>
```

#### 画像ファイルの参照

`client/src/app/page.tsx`などで画像ファイルを参照:

```tsx
// <img>タグを使用する場合
<img src="/images/logo.png" alt="Logo" />

// next/imageコンポーネントを使用する場合（推奨）
import Image from 'next/image'
<Image src="/images/logo.png" alt="Logo" width={100} height={100} />
```

**重要なポイント**:
- パスに`public/`を含めない（`/css/style.css`が正しい）
- `public/`ディレクトリがルート（`/`）として扱われる
- 開発環境と本番環境の両方で同じ動作

## ライセンス

MIT License
