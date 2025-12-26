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
│   └── sharding.hcl          # シャーディングDBのスキーマ定義
└── migrations/               # マイグレーションファイル（初期データ含む）
    ├── master/               # マスターDB用マイグレーション
    │   ├── 20251226_initial.sql
    │   └── atlas.sum
    └── sharding/             # シャーディングDB用マイグレーション
        ├── 20251226_initial.sql
        └── atlas.sum

config/
├── develop/atlas.hcl         # 開発環境用Atlas設定
├── staging/atlas.hcl         # ステージング環境用Atlas設定
└── production/atlas.hcl      # 本番環境用Atlas設定
```

## 基本コマンド

### マイグレーションの生成

スキーマ定義を変更した後、差分マイグレーションを生成します。

```bash
# マスターDBのマイグレーション生成
atlas migrate diff <migration_name> \
    --dir file://db/migrations/master \
    --to file://db/schema/master.hcl \
    --dev-url "sqlite://file?mode=memory"

# シャーディングDBのマイグレーション生成
atlas migrate diff <migration_name> \
    --dir file://db/migrations/sharding \
    --to file://db/schema/sharding.hcl \
    --dev-url "sqlite://file?mode=memory"
```

### マイグレーションの適用

```bash
# マスターDBへのマイグレーション適用
atlas migrate apply \
    --dir file://db/migrations/master \
    --url "sqlite://server/data/master.db"

# シャーディングDBへのマイグレーション適用（全シャード）
for i in 1 2 3 4; do
    atlas migrate apply \
        --dir file://db/migrations/sharding \
        --url "sqlite://server/data/sharding_db_${i}.db"
done
```

### マイグレーション状態の確認

```bash
# マスターDBのマイグレーション状態
atlas migrate status \
    --dir file://db/migrations/master \
    --url "sqlite://server/data/master.db"

# シャーディングDBのマイグレーション状態
atlas migrate status \
    --dir file://db/migrations/sharding \
    --url "sqlite://server/data/sharding_db_1.db"
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
- シャーディングDBは全シャードに同じマイグレーションを適用する必要があります

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

シャーディングDBは全シャードに同じマイグレーションを適用する必要があります。

```bash
# 全シャードにマイグレーションを適用
for i in 1 2 3 4; do
    echo "=== Applying to sharding_db_${i}.db ==="
    atlas migrate apply \
        --dir file://db/migrations/sharding \
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

#### Sharding
- `20251226074934_initial.sql` - 初期スキーマ
- `20251226130618_add_sharding_experiment.sql` - テーブル追加実験
- `20251226130721_add_sharding_experiment_desc.sql` - カラム追加実験
- `20251226130814_insert_sharding_experiment_data.sql` - データ挿入実験
- `20251226130931_drop_sharding_experiment.sql` - テーブル削除実験
