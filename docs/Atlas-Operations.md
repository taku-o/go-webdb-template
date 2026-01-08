# Atlas マイグレーション運用ガイド

## 概要

本プロジェクトでは、データベーススキーマ管理に [Atlas](https://atlasgo.io/) を使用しています。
Atlas は宣言的スキーマ管理ツールで、HCL ファイルでスキーマを定義し、
自動的にマイグレーションを生成・適用します。

## ディレクトリ構成

```
db/
├── schema/                    # スキーマ定義ファイル（HCL）
│   ├── master.hcl            # マスターDBのスキーマ定義
│   ├── sharding_1/           # シャード1用スキーマ（テーブル000-007）
│   │   ├── _schema.hcl       # スキーマ定義（schema "public" {}）
│   │   ├── dm_users.hcl      # dm_users_000 〜 dm_users_007
│   │   └── dm_posts.hcl      # dm_posts_000 〜 dm_posts_007
│   ├── sharding_2/           # シャード2用スキーマ（テーブル008-015）
│   │   ├── _schema.hcl
│   │   ├── dm_users.hcl      # dm_users_008 〜 dm_users_015
│   │   └── dm_posts.hcl      # dm_posts_008 〜 dm_posts_015
│   ├── sharding_3/           # シャード3用スキーマ（テーブル016-023）
│   │   ├── _schema.hcl
│   │   ├── dm_users.hcl      # dm_users_016 〜 dm_users_023
│   │   └── dm_posts.hcl      # dm_posts_016 〜 dm_posts_023
│   └── sharding_4/           # シャード4用スキーマ（テーブル024-031）
│       ├── _schema.hcl
│       ├── dm_users.hcl      # dm_users_024 〜 dm_users_031
│       └── dm_posts.hcl      # dm_posts_024 〜 dm_posts_031
└── migrations/               # マイグレーションファイル（初期データ含む）
    ├── master/               # マスターDB用マイグレーション
    │   ├── YYYYMMDD_*.sql
    │   └── atlas.sum
    ├── sharding_1/           # webdb_sharding_1用マイグレーション（テーブル000-007）
    │   ├── YYYYMMDD_*.sql
    │   └── atlas.sum
    ├── sharding_2/           # webdb_sharding_2用マイグレーション（テーブル008-015）
    │   ├── YYYYMMDD_*.sql
    │   └── atlas.sum
    ├── sharding_3/           # webdb_sharding_3用マイグレーション（テーブル016-023）
    │   ├── YYYYMMDD_*.sql
    │   └── atlas.sum
    ├── sharding_4/           # webdb_sharding_4用マイグレーション（テーブル024-031）
    │   ├── YYYYMMDD_*.sql
    │   └── atlas.sum
    └── view_master/          # マスターDB用ビューマイグレーション
        ├── YYYYMMDD_*.sql
        └── atlas.sum

config/
├── develop/atlas.hcl         # 開発環境用Atlas設定
├── staging/atlas.hcl         # ステージング環境用Atlas設定
└── production/atlas.hcl      # 本番環境用Atlas設定
```

### テーブル分割ルール

シャーディングDBは4つのデータベースに分割され、各データベースには8つのテーブル分割が含まれます。

| データベース | テーブル範囲 | スキーマディレクトリ | マイグレーションディレクトリ |
|------------|-----------|------------------|----------------------|
| webdb_sharding_1 | dm_users_000-007, dm_posts_000-007 | db/schema/sharding_1/ | db/migrations/sharding_1/ |
| webdb_sharding_2 | dm_users_008-015, dm_posts_008-015 | db/schema/sharding_2/ | db/migrations/sharding_2/ |
| webdb_sharding_3 | dm_users_016-023, dm_posts_016-023 | db/schema/sharding_3/ | db/migrations/sharding_3/ |
| webdb_sharding_4 | dm_users_024-031, dm_posts_024-031 | db/schema/sharding_4/ | db/migrations/sharding_4/ |

## PostgreSQLコンテナ構成

| コンテナ名 | データベース名 | ホストポート |
|-----------|--------------|-------------|
| postgres-master | webdb_master | 5432 |
| postgres-sharding-1 | webdb_sharding_1 | 5433 |
| postgres-sharding-2 | webdb_sharding_2 | 5434 |
| postgres-sharding-3 | webdb_sharding_3 | 5435 |
| postgres-sharding-4 | webdb_sharding_4 | 5436 |

## 基本コマンド

### マイグレーションスクリプト

```bash
# すべてのマイグレーションを適用（デフォルト）
APP_ENV=develop ./scripts/migrate.sh all

# masterデータベースのみ
APP_ENV=develop ./scripts/migrate.sh master

# shardingデータベースのみ
APP_ENV=develop ./scripts/migrate.sh sharding
```

### マイグレーションの生成

スキーマ定義を変更した後、差分マイグレーションを生成します。

```bash
# マスターDBのマイグレーション生成
atlas migrate diff <migration_name> \
    --dir file://db/migrations/master \
    --to file://db/schema/master.hcl \
    --dev-url "postgres://webdb:webdb@localhost:5432/webdb_master?sslmode=disable"

# シャーディングDBのマイグレーション生成（各シャードごとに実行）
# シャード1
atlas migrate diff <migration_name> \
    --dir file://db/migrations/sharding_1 \
    --to file://db/schema/sharding_1 \
    --dev-url "postgres://webdb:webdb@localhost:5433/webdb_sharding_1?sslmode=disable"

# シャード2
atlas migrate diff <migration_name> \
    --dir file://db/migrations/sharding_2 \
    --to file://db/schema/sharding_2 \
    --dev-url "postgres://webdb:webdb@localhost:5434/webdb_sharding_2?sslmode=disable"

# シャード3
atlas migrate diff <migration_name> \
    --dir file://db/migrations/sharding_3 \
    --to file://db/schema/sharding_3 \
    --dev-url "postgres://webdb:webdb@localhost:5435/webdb_sharding_3?sslmode=disable"

# シャード4
atlas migrate diff <migration_name> \
    --dir file://db/migrations/sharding_4 \
    --to file://db/schema/sharding_4 \
    --dev-url "postgres://webdb:webdb@localhost:5436/webdb_sharding_4?sslmode=disable"
```

### マイグレーションの適用

```bash
# マスターDBへのマイグレーション適用
atlas migrate apply \
    --dir file://db/migrations/master \
    --url "postgres://webdb:webdb@localhost:5432/webdb_master?sslmode=disable"

# シャーディングDBへのマイグレーション適用（各シャードごとに対応するマイグレーションディレクトリを使用）
for i in 1 2 3 4; do
    port=$((5432 + i))
    atlas migrate apply \
        --dir file://db/migrations/sharding_${i} \
        --url "postgres://webdb:webdb@localhost:${port}/webdb_sharding_${i}?sslmode=disable"
done
```

### マイグレーション状態の確認

```bash
# マスターDBのマイグレーション状態
atlas migrate status \
    --dir file://db/migrations/master \
    --url "postgres://webdb:webdb@localhost:5432/webdb_master?sslmode=disable"

# シャーディングDBのマイグレーション状態（各シャードごとに確認）
for i in 1 2 3 4; do
    echo "=== webdb_sharding_${i} ==="
    port=$((5432 + i))
    atlas migrate status \
        --dir file://db/migrations/sharding_${i} \
        --url "postgres://webdb:webdb@localhost:${port}/webdb_sharding_${i}?sslmode=disable"
done
```

## 環境別適用手順

### 開発環境

```bash
# 設定ファイルを使用したマイグレーション
atlas migrate apply \
    --config file://config/develop/atlas.hcl \
    --env master

atlas migrate apply \
    --config file://config/develop/atlas.hcl \
    --env sharding_1

# または簡易スクリプト
APP_ENV=develop ./scripts/migrate.sh all
```

### ステージング環境

```bash
# 環境変数でDB URLを設定するか、設定ファイルを編集
atlas migrate apply \
    --config file://config/staging/atlas.hcl \
    --env master

# 各シャードにも適用
for env in sharding_1 sharding_2 sharding_3 sharding_4; do
    atlas migrate apply \
        --config file://config/staging/atlas.hcl \
        --env $env
done
```

### 本番環境

```bash
# 本番環境では dry-run で確認後に適用
atlas migrate apply \
    --config file://config/production/atlas.hcl \
    --env master \
    --dry-run

# 問題なければ適用
atlas migrate apply \
    --config file://config/production/atlas.hcl \
    --env master
```

## ケース別運用方法

### テーブルの追加

1. スキーマファイル (`db/schema/master.hcl` または `db/schema/sharding_*/`) にテーブル定義を追加

```hcl
table "new_table" {
  schema = schema.public
  column "id" {
    null = false
    type = serial
  }
  column "name" {
    null = false
    type = text
  }
  primary_key {
    columns = [column.id]
  }
}
```

2. マイグレーションを生成

```bash
atlas migrate diff add_new_table \
    --dir file://db/migrations/master \
    --to file://db/schema/master.hcl \
    --dev-url "postgres://webdb:webdb@localhost:5432/webdb_master?sslmode=disable"
```

3. 生成されたSQLを確認し、適用

### カラムの追加

1. スキーマファイルのテーブル定義にカラムを追加

```hcl
table "existing_table" {
  # ... 既存のカラム ...

  column "new_column" {
    null = true
    type = text
  }
}
```

2. マイグレーション生成・適用

### テーブルの削除

1. スキーマファイルからテーブル定義を削除
2. マイグレーション生成・適用

**注意**: テーブル削除はデータも削除されるため、本番環境では慎重に実行してください。

### インデックスの追加

```hcl
table "users" {
  # ... カラム定義 ...

  index "idx_users_email" {
    columns = [column.email]
  }

  # ユニークインデックス
  index "idx_users_unique_email" {
    unique  = true
    columns = [column.email]
  }
}
```

## イレギュラーケース対応

### 直接SQLでスキーマ変更した場合

Atlas管理外でSQLを直接実行した場合、スキーマとマイグレーション履歴に不整合が生じます。

**対処法 1: ベースライン設定**

```bash
# 現在のDBスキーマをベースラインとして設定
atlas migrate hash --dir file://db/migrations/master

# atlas_schema_revisions テーブルを手動で更新
docker exec -i postgres-master psql -U webdb -d webdb_master -c "DELETE FROM atlas_schema_revisions"
atlas migrate apply \
    --dir file://db/migrations/master \
    --url "postgres://webdb:webdb@localhost:5432/webdb_master?sslmode=disable" \
    --baseline <version>
```

**対処法 2: スキーマの同期**

```bash
# 現在のDBからスキーマをインスペクト
atlas schema inspect \
    --url "postgres://webdb:webdb@localhost:5432/webdb_master?sslmode=disable" \
    --format hcl > db/schema/master_current.hcl

# 差分を確認して master.hcl を更新
diff db/schema/master.hcl db/schema/master_current.hcl
```

### マイグレーション適用に失敗した場合

```bash
# マイグレーション状態を確認
atlas migrate status \
    --dir file://db/migrations/master \
    --url "postgres://webdb:webdb@localhost:5432/webdb_master?sslmode=disable"

# 部分適用されたマイグレーションがある場合は手動で修正
# 必要に応じてロールバック用SQLを実行
```

### データベースの初期化（ゼロから再構築）

```bash
# PostgreSQLコンテナを停止し、データディレクトリを削除
./scripts/start-postgres.sh stop
rm -rf postgres/data/master postgres/data/sharding_*

# PostgreSQLコンテナを再起動
./scripts/start-postgres.sh start

# マイグレーションを適用（初期データも含む）
APP_ENV=develop ./scripts/migrate.sh all
```

## トラブルシューティング

### 接続エラー

**症状**: `connection refused` または `no such host`

**解決策**:
```bash
# PostgreSQLコンテナの状態確認
./scripts/start-postgres.sh status
./scripts/start-postgres.sh health
```

### チェックサムエラー

**症状**: `checksum mismatch`

**解決策**:
```bash
atlas migrate hash --dir "file://db/migrations/master"
atlas migrate hash --dir "file://db/migrations/sharding_1"
atlas migrate hash --dir "file://db/migrations/sharding_2"
atlas migrate hash --dir "file://db/migrations/sharding_3"
atlas migrate hash --dir "file://db/migrations/sharding_4"
```

### Atlas CLIエラー

**症状**: `atlas: command not found`

**解決策**:
```bash
# macOS
brew install ariga/tap/atlas

# その他のOS
curl -sSf https://atlasgo.sh | sh
```

## 注意事項

- マイグレーションファイルは一度生成したら編集しないでください
- `atlas.sum` ファイルはマイグレーションの整合性チェックに使用されます
- 本番環境では必ずバックアップを取得してからマイグレーションを実行してください
- シャーディングDBは各シャードに対応するマイグレーションディレクトリを使用して適用する必要があります
- シャーディングスキーマを変更する場合は、4つのスキーマディレクトリ（sharding_1〜4）すべてを更新してください

## データ更新用マイグレーションの作成

Atlasはスキーマ変更のみを検出するため、データ更新用のマイグレーションは手動で作成する必要があります。

```bash
# 空のマイグレーションファイルを作成
atlas migrate new insert_data --dir file://db/migrations/master

# 生成されたファイルにSQLを追加
# 例: INSERT INTO table_name (col1, col2) VALUES ('value1', 'value2');

# atlas.sumを更新
atlas migrate hash --dir file://db/migrations/master

# マイグレーションを適用
atlas migrate apply \
    --dir file://db/migrations/master \
    --url "postgres://webdb:webdb@localhost:5432/webdb_master?sslmode=disable"
```

## ビュー（VIEW）の管理

### 概要

ビューはSQLファイルを作り、それをAtlasで適用する形で運用します。

### ディレクトリ構成

ビュー用のマイグレーションはテーブル用とは別ディレクトリで管理します。

```
db/
├── schema/
│   └── master.hcl              # テーブル定義のみ（ビュー定義は含めない）
└── migrations/
    ├── master/                 # テーブル用マイグレーション
    │   ├── YYYYMMDD_*.sql
    │   └── atlas.sum
    └── view_master/            # マスターDB用ビューマイグレーション
        ├── YYYYMMDD_*.sql
        └── atlas.sum
```

### ビューの作成手順

#### 1. SQLファイルを手動作成

`db/migrations/view_master/`にSQLファイルを作成します。

```bash
# ディレクトリと名前を指定して、空のファイルを生成
atlas migrate new create_view_name \
    --dir "file://db/migrations/view_master"
```

```sql
-- Create dm_news_view
CREATE VIEW dm_news_view AS SELECT id, title, content, published_at FROM dm_news;
```

#### 2. チェックサムを更新

```bash
atlas migrate hash --dir "file://db/migrations/view_master"
```

#### 3. マイグレーションを適用

ビューのマイグレーションは`scripts/migrate.sh`で自動的に適用されます。

### ビューの削除

```sql
-- Drop dm_news_view
DROP VIEW IF EXISTS dm_news_view;
```

### 注意事項

- ビューはベーステーブルに依存するため、テーブル用マイグレーションを先に適用する必要があります
- 履歴テーブル（`atlas_schema_revisions`）は共通です。ファイル名（バージョン番号）で管理されます
- HCLファイルにビュー定義を追加すると、`atlas migrate diff`が使用できなくなるため、ビュー定義はHCLには追加しないでください
