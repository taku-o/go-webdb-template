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
│   │   ├── _schema.hcl       # スキーマ定義（schema "main" {}）
│   │   ├── users.hcl         # users_000 〜 users_007
│   │   └── posts.hcl         # posts_000 〜 posts_007
│   ├── sharding_2/           # シャード2用スキーマ（テーブル008-015）
│   │   ├── _schema.hcl
│   │   ├── users.hcl         # users_008 〜 users_015
│   │   └── posts.hcl         # posts_008 〜 posts_015
│   ├── sharding_3/           # シャード3用スキーマ（テーブル016-023）
│   │   ├── _schema.hcl
│   │   ├── users.hcl         # users_016 〜 users_023
│   │   └── posts.hcl         # posts_016 〜 posts_023
│   └── sharding_4/           # シャード4用スキーマ（テーブル024-031）
│       ├── _schema.hcl
│       ├── users.hcl         # users_024 〜 users_031
│       └── posts.hcl         # posts_024 〜 posts_031
└── migrations/               # マイグレーションファイル（初期データ含む）
    ├── master/               # マスターDB用マイグレーション
    │   ├── 20251226_initial.sql
    │   └── atlas.sum
    ├── sharding_1/           # sharding_db_1.db用マイグレーション（テーブル000-007）
    │   ├── YYYYMMDD_initial.sql
    │   └── atlas.sum
    ├── sharding_2/           # sharding_db_2.db用マイグレーション（テーブル008-015）
    │   ├── YYYYMMDD_initial.sql
    │   └── atlas.sum
    ├── sharding_3/           # sharding_db_3.db用マイグレーション（テーブル016-023）
    │   ├── YYYYMMDD_initial.sql
    │   └── atlas.sum
    └── sharding_4/           # sharding_db_4.db用マイグレーション（テーブル024-031）
        ├── YYYYMMDD_initial.sql
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
| sharding_db_1.db | users_000-007, posts_000-007 | db/schema/sharding_1/ | db/migrations/sharding_1/ |
| sharding_db_2.db | users_008-015, posts_008-015 | db/schema/sharding_2/ | db/migrations/sharding_2/ |
| sharding_db_3.db | users_016-023, posts_016-023 | db/schema/sharding_3/ | db/migrations/sharding_3/ |
| sharding_db_4.db | users_024-031, posts_024-031 | db/schema/sharding_4/ | db/migrations/sharding_4/ |

## 基本コマンド

### マイグレーションの生成

スキーマ定義を変更した後、差分マイグレーションを生成します。

```bash
# マスターDBのマイグレーション生成
atlas migrate diff <migration_name> \
    --dir file://db/migrations/master \
    --to file://db/schema/master.hcl \
    --dev-url "sqlite://file?mode=memory"

# シャーディングDBのマイグレーション生成（各シャードごとに実行）
# シャード1
atlas migrate diff <migration_name> \
    --dir file://db/migrations/sharding_1 \
    --to file://db/schema/sharding_1 \
    --dev-url "sqlite://file?mode=memory"

# シャード2
atlas migrate diff <migration_name> \
    --dir file://db/migrations/sharding_2 \
    --to file://db/schema/sharding_2 \
    --dev-url "sqlite://file?mode=memory"

# シャード3
atlas migrate diff <migration_name> \
    --dir file://db/migrations/sharding_3 \
    --to file://db/schema/sharding_3 \
    --dev-url "sqlite://file?mode=memory"

# シャード4
atlas migrate diff <migration_name> \
    --dir file://db/migrations/sharding_4 \
    --to file://db/schema/sharding_4 \
    --dev-url "sqlite://file?mode=memory"
```

### マイグレーションの適用

```bash
# マスターDBへのマイグレーション適用
atlas migrate apply \
    --dir file://db/migrations/master \
    --url "sqlite://server/data/master.db"

# シャーディングDBへのマイグレーション適用（各シャードごとに対応するマイグレーションディレクトリを使用）
for i in 1 2 3 4; do
    atlas migrate apply \
        --dir file://db/migrations/sharding_${i} \
        --url "sqlite://server/data/sharding_db_${i}.db"
done
```

### マイグレーション状態の確認

```bash
# マスターDBのマイグレーション状態
atlas migrate status \
    --dir file://db/migrations/master \
    --url "sqlite://server/data/master.db"

# シャーディングDBのマイグレーション状態（各シャードごとに確認）
for i in 1 2 3 4; do
    echo "=== sharding_db_${i} ==="
    atlas migrate status \
        --dir file://db/migrations/sharding_${i} \
        --url "sqlite://server/data/sharding_db_${i}.db"
done
```

### スクリプトを使用したマイグレーション

```bash
# 全データベースにマイグレーションを適用
./scripts/migrate.sh all

# マスターDBのみ
./scripts/migrate.sh master

# シャーディングDBのみ
./scripts/migrate.sh sharding
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
./scripts/migrate.sh all
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

1. スキーマファイル (`db/schema/master.hcl` または `db/schema/sharding.hcl`) にテーブル定義を追加

```hcl
table "new_table" {
  schema = schema.main
  column "id" {
    null           = false
    type           = integer
    auto_increment = true
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
    --dev-url "sqlite://file?mode=memory"
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
sqlite3 server/data/master.db "DELETE FROM atlas_schema_revisions"
atlas migrate apply \
    --dir file://db/migrations/master \
    --url "sqlite://server/data/master.db" \
    --baseline <version>
```

**対処法 2: スキーマの同期**

```bash
# 現在のDBからスキーマをインスペクト
atlas schema inspect \
    --url "sqlite://server/data/master.db" \
    --format hcl > db/schema/master_current.hcl

# 差分を確認して master.hcl を更新
diff db/schema/master.hcl db/schema/master_current.hcl
```

### マイグレーション適用に失敗した場合

```bash
# マイグレーション状態を確認
atlas migrate status \
    --dir file://db/migrations/master \
    --url "sqlite://server/data/master.db"

# 部分適用されたマイグレーションがある場合は手動で修正
# 必要に応じてロールバック用SQLを実行
```

### データベースの初期化（ゼロから再構築）

```bash
# 既存のDBを削除
rm -f server/data/master.db server/data/sharding_db_*.db

# マイグレーションを適用（初期データも含む）
./scripts/migrate.sh all
```

## 注意事項

- マイグレーションファイルは一度生成したら編集しないでください
- `atlas.sum` ファイルはマイグレーションの整合性チェックに使用されます
- 本番環境では必ずバックアップを取得してからマイグレーションを実行してください
- シャーディングDBは各シャードに対応するマイグレーションディレクトリを使用して適用する必要があります
- シャーディングスキーマを変更する場合は、4つのスキーマディレクトリ（sharding_1〜4）すべてを更新してください

## 運用実験結果

本プロジェクトで実施したAtlas運用実験の結果を記録します。

### 実験結果サマリ

| 操作 | Master | Sharding | 結果 |
|------|--------|----------|------|
| テーブル追加 | OK | OK | `atlas migrate diff` で自動生成 |
| カラム追加 | OK | OK | `ALTER TABLE ADD COLUMN` が自動生成 |
| データ更新 | OK | OK | `atlas migrate new` で空ファイル作成後、手動でSQL追加 |
| テーブル削除 | OK | OK | `DROP TABLE` が自動生成（PRAGMA foreign_keys含む） |

### データ更新用マイグレーションの作成

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
    --url "sqlite://server/data/master.db"
```

### 直接SQL実行後の同期方法

Atlas管理外でSQLを直接実行した場合の対処方法：

1. **現在のDBスキーマをインスペクト**
   ```bash
   atlas schema inspect --url "sqlite://server/data/master.db" --format hcl
   ```

2. **HCLスキーマファイルを更新**
   インスペクト結果を参考に `db/schema/master.hcl` を更新

3. **マイグレーションを生成**
   ```bash
   atlas migrate diff sync_changes \
       --dir file://db/migrations/master \
       --to file://db/schema/master.hcl \
       --dev-url "sqlite://file?mode=memory"
   ```

4. **マイグレーション履歴を手動で更新**（DBにテーブルが既に存在する場合）
   ```bash
   # atlas_schema_revisionsテーブルに直接レコードを追加
   sqlite3 server/data/master.db "INSERT INTO atlas_schema_revisions (...) VALUES (...)"
   ```

5. **ステータスを確認**
   ```bash
   atlas migrate status \
       --dir file://db/migrations/master \
       --url "sqlite://server/data/master.db"
   ```

### シャーディングDBへの一括適用

シャーディングDBは各シャードに対応するマイグレーションディレクトリを使用して適用します。

```bash
# 全シャードにマイグレーションを適用（各シャードに対応するマイグレーションディレクトリを使用）
for i in 1 2 3 4; do
    echo "=== Applying to sharding_db_${i}.db ==="
    atlas migrate apply \
        --dir file://db/migrations/sharding_${i} \
        --url "sqlite://server/data/sharding_db_${i}.db"
done
```

### 実験で生成されたマイグレーションファイル

#### Master
- `20251226074846_initial.sql` - 初期スキーマ
- `20251226130113_add_experiment_table.sql` - テーブル追加実験
- `20251226130242_add_description_column.sql` - カラム追加実験
- `20251226130340_insert_experiment_data.sql` - データ挿入実験
- `20251226130500_drop_experiment_table.sql` - テーブル削除実験
- `20251226131107_sync_direct_sql_table.sql` - 直接SQL同期実験
- `20251226131150_drop_direct_sql_table.sql` - クリーンアップ

#### Sharding（現在の構成）
各シャーディングDBに対応するマイグレーションディレクトリが存在します：

- `db/migrations/sharding_1/` - sharding_db_1.db用（テーブル000-007）
- `db/migrations/sharding_2/` - sharding_db_2.db用（テーブル008-015）
- `db/migrations/sharding_3/` - sharding_db_3.db用（テーブル016-023）
- `db/migrations/sharding_4/` - sharding_db_4.db用（テーブル024-031）

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
    ├── view_master/            # マスターDB用ビューマイグレーション
    │   ├── YYYYMMDD_*.sql
    │   └── atlas.sum
    ├── view_sharding_1/        # シャード1用ビューマイグレーション（将来用）
    ├── view_sharding_2/        # シャード2用ビューマイグレーション（将来用）
    ├── view_sharding_3/        # シャード3用ビューマイグレーション（将来用）
    └── view_sharding_4/        # シャード4用ビューマイグレーション（将来用）
```

### ビューの作成手順

#### 1. SQLファイルを手動作成

`db/migrations/view_master/`にSQLファイルを作成します。

```bash
# ディレクトリと名前を指定して、空のファイルを生成
# db/migrations/view_master/20260103030226_create_view_name.sqlが生成される
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

```bash
# テーブル用を先に適用（ベーステーブルが必要）
atlas migrate apply \
    --dir "file://db/migrations/master" \
    --url "sqlite://server/data/master.db"

# ビュー用を後に適用
atlas migrate apply \
    --dir "file://db/migrations/view_master" \
    --url "sqlite://server/data/master.db"
```

### ビューの削除

```sql
-- Drop dm_news_view
DROP VIEW IF EXISTS dm_news_view;
```

### 注意事項

- ビューはベーステーブルに依存するため、テーブル用マイグレーションを先に適用する必要があります
- 履歴テーブル（`atlas_schema_revisions`）は共通です。ファイル名（バージョン番号）で管理されます
- HCLファイルにビュー定義を追加すると、`atlas migrate diff`が使用できなくなるため、ビュー定義はHCLには追加しないでください
