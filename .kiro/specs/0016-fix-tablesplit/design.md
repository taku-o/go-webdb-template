# テーブル分割修正設計書

## 1. 概要

### 1.1 設計の目的
要件定義書に基づき、シャーディングデータベースごとに適切なテーブル範囲のみが作成されるよう、スキーマ定義ファイルをテーブルごとに分割し、マイグレーション管理システムを修正する詳細設計を定義する。既存のAtlas設定とマイグレーションスクリプトを修正し、各データベースに必要なテーブルのみが作成されるようにする。

### 1.2 設計の範囲
- スキーマ定義ディレクトリの分割（データベースごとに個別のディレクトリを作成）
- スキーマ定義ファイルの分割（テーブルごとにファイルを分割：`_schema.hcl`, `users.hcl`, `posts.hcl`）
- マイグレーションディレクトリの分割（データベースごとに個別のディレクトリを作成）
- マイグレーション適用スクリプトの修正（データベースごとのマイグレーションディレクトリを参照するように修正）
- Atlas設定ファイルの修正（データベースごとのスキーマディレクトリとマイグレーションディレクトリを参照するように修正）
- 既存データベースのクリーンアップ（削除して再作成）

### 1.3 設計方針
- **テーブルごとのファイル分割**: スキーマファイルをテーブルタイプごとに分割し、将来の拡張を容易にする
- **データベースごとの分離**: 各データベース用のスキーマディレクトリとマイグレーションディレクトリを分離し、管理を明確化
- **Atlasの複数ファイル読み込み機能を活用**: ディレクトリ内のすべての`.hcl`ファイルを自動的に読み込む機能を利用
- **既存システムとの互換性**: 既存のAtlas設定ファイルとマイグレーションスクリプトの構造を維持しつつ、参照先を変更
- **データベースのリセット**: 既存データの移行は不要。データベースを削除して再作成する

## 2. アーキテクチャ設計

### 2.1 ディレクトリ構造

#### 2.1.1 変更前の構造
```
db/
├── schema/
│   ├── master.hcl
│   └── sharding.hcl              # 全32テーブルの定義を含む単一ファイル
└── migrations/
    ├── master/
    │   ├── 20251226074846_initial.sql
    │   └── atlas.sum
    └── sharding/                  # 全データベース共通
        ├── 20251226074934_initial.sql  # 全32テーブルのCREATE文を含む
        └── atlas.sum

config/
├── develop/
│   └── atlas.hcl                 # 全データベースが同じスキーマファイルとマイグレーションディレクトリを参照
├── staging/
│   └── atlas.hcl
└── production/
    └── atlas.hcl

scripts/
└── migrate.sh                    # 全データベースに同じマイグレーションファイルを適用
```

#### 2.1.2 変更後の構造
```
db/
├── schema/
│   ├── master.hcl                # 既存（維持）
│   └── sharding_1/               # 新規: シャード1用のディレクトリ
│       ├── _schema.hcl           # スキーマ定義（schema "main" {} のみ）
│       ├── users.hcl              # users_000 〜 users_007
│       └── posts.hcl              # posts_000 〜 posts_007
│   └── sharding_2/               # 新規: シャード2用のディレクトリ
│       ├── _schema.hcl
│       ├── users.hcl              # users_008 〜 users_015
│       └── posts.hcl              # posts_008 〜 posts_015
│   └── sharding_3/               # 新規: シャード3用のディレクトリ
│       ├── _schema.hcl
│       ├── users.hcl              # users_016 〜 users_023
│       └── posts.hcl              # posts_016 〜 posts_023
│   └── sharding_4/               # 新規: シャード4用のディレクトリ
│       ├── _schema.hcl
│       ├── users.hcl              # users_024 〜 users_031
│       └── posts.hcl              # posts_024 〜 posts_031
└── migrations/
    ├── master/                   # 既存（維持）
    │   ├── 20251226074846_initial.sql
    │   └── atlas.sum
    ├── sharding_1/               # 新規: sharding_db_1.db用マイグレーション
    │   ├── 20251226074934_initial.sql  # テーブル000-007のみを含む
    │   └── atlas.sum
    ├── sharding_2/               # 新規: sharding_db_2.db用マイグレーション
    │   ├── 20251226074934_initial.sql  # テーブル008-015のみを含む
    │   └── atlas.sum
    ├── sharding_3/               # 新規: sharding_db_3.db用マイグレーション
    │   ├── 20251226074934_initial.sql  # テーブル016-023のみを含む
    │   └── atlas.sum
    └── sharding_4/               # 新規: sharding_db_4.db用マイグレーション
        ├── 20251226074934_initial.sql  # テーブル024-031のみを含む
        └── atlas.sum

config/
├── develop/
│   └── atlas.hcl                 # 修正: 各データベースごとのスキーマディレクトリとマイグレーションディレクトリを参照
├── staging/
│   └── atlas.hcl                 # 修正: 各データベースごとのスキーマディレクトリとマイグレーションディレクトリを参照
└── production/
    └── atlas.hcl                 # 修正: 各データベースごとのスキーマディレクトリとマイグレーションディレクトリを参照

scripts/
└── migrate.sh                    # 修正: 各データベースごとのマイグレーションディレクトリを参照
```

### 2.2 ファイル構成

#### 2.2.1 スキーマ定義ファイル

**`_schema.hcl`**: スキーマ定義のみ
```hcl
schema "main" {
}
```

**`users.hcl`**: usersテーブルの定義（該当するテーブル範囲のみ）
- `db/schema/sharding_1/users.hcl`: users_000 〜 users_007の定義
- `db/schema/sharding_2/users.hcl`: users_008 〜 users_015の定義
- `db/schema/sharding_3/users.hcl`: users_016 〜 users_023の定義
- `db/schema/sharding_4/users.hcl`: users_024 〜 users_031の定義

**`posts.hcl`**: postsテーブルの定義（該当するテーブル範囲のみ）
- `db/schema/sharding_1/posts.hcl`: posts_000 〜 posts_007の定義
- `db/schema/sharding_2/posts.hcl`: posts_008 〜 posts_015の定義
- `db/schema/sharding_3/posts.hcl`: posts_016 〜 posts_023の定義
- `db/schema/sharding_4/posts.hcl`: posts_024 〜 posts_031の定義

#### 2.2.2 Atlas設定ファイル

**`config/{env}/atlas.hcl`**: 環境別Atlas設定ファイル
- 各データベース用の環境設定で、個別のスキーマディレクトリとマイグレーションディレクトリを参照
- `env "sharding_1"`, `env "sharding_2"`, `env "sharding_3"`, `env "sharding_4"`の各設定を修正

#### 2.2.3 マイグレーション適用スクリプト

**`scripts/migrate.sh`**: マイグレーション適用スクリプト
- 各データベースごとに個別のマイグレーションディレクトリを参照するように修正

### 2.3 システム構成図

```
┌─────────────────────────────────────────────────────────┐
│                    開発者                                │
└──────────────────┬────────────────────────────────────┘
                    │
                    │ ./scripts/migrate.sh sharding
                    │
                    ▼
┌─────────────────────────────────────────────────────────┐
│              scripts/migrate.sh                          │
│  - 各データベースごとのマイグレーションディレクトリを参照 │
│  - sharding_db_1.db → db/migrations/sharding_1/        │
│  - sharding_db_2.db → db/migrations/sharding_2/        │
│  - sharding_db_3.db → db/migrations/sharding_3/        │
│  - sharding_db_4.db → db/migrations/sharding_4/        │
└──────────────────┬────────────────────────────────────┘
                    │
                    │ atlas migrate apply
                    │
                    ▼
┌─────────────────────────────────────────────────────────┐
│         config/{env}/atlas.hcl                           │
│                                                          │
│  env "sharding_1" {                                      │
│    src = "file://db/schema/sharding_1"                  │
│    dir = "file://db/migrations/sharding_1"              │
│  }                                                       │
│  env "sharding_2" {                                      │
│    src = "file://db/schema/sharding_2"                  │
│    dir = "file://db/migrations/sharding_2"              │
│  }                                                       │
│  env "sharding_3" {                                      │
│    src = "file://db/schema/sharding_3"                  │
│    dir = "file://db/migrations/sharding_3"              │
│  }                                                       │
│  env "sharding_4" {                                      │
│    src = "file://db/schema/sharding_4"                  │
│    dir = "file://db/migrations/sharding_4"              │
│  }                                                       │
└──────────────────┬────────────────────────────────────┘
                    │
                    │ Atlasがスキーマディレクトリ内の
                    │ すべての.hclファイルを自動的に読み込み
                    │
                    ▼
┌─────────────────────────────────────────────────────────┐
│         db/schema/sharding_1/                            │
│  ├── _schema.hcl  (schema "main" {})                    │
│  ├── users.hcl   (users_000 〜 users_007)               │
│  └── posts.hcl   (posts_000 〜 posts_007)               │
│                                                          │
│  Atlasが自動的に結合してスキーマとして解釈               │
└──────────────────┬────────────────────────────────────┘
                    │
                    │ マイグレーション生成・適用
                    │
                    ▼
┌─────────────────────────────────────────────────────────┐
│         server/data/                                     │
│  ├── sharding_db_1.db  (テーブル000-007のみ)            │
│  ├── sharding_db_2.db  (テーブル008-015のみ)            │
│  ├── sharding_db_3.db  (テーブル016-023のみ)            │
│  └── sharding_db_4.db  (テーブル024-031のみ)            │
└─────────────────────────────────────────────────────────┘
```

### 2.4 データフロー

#### 2.4.1 マイグレーション生成フロー
```
既存のdb/schema/sharding.hclを分析
    ↓
各データベース用のスキーマディレクトリを作成
    ↓
db/schema/sharding_1/, sharding_2/, sharding_3/, sharding_4/
    ↓
各ディレクトリに_schema.hcl, users.hcl, posts.hclを作成
    ↓
Atlasでマイグレーションを生成
    ↓
atlas migrate diff initial \
  --dir file://db/migrations/sharding_1 \
  --to file://db/schema/sharding_1 \
  --dev-url "sqlite://file?mode=memory"
    ↓
各マイグレーションディレクトリにマイグレーションファイルが生成される
```

#### 2.4.2 マイグレーション適用フロー
```
./scripts/migrate.sh sharding を実行
    ↓
各データベースごとにループ処理
    ↓
atlas migrate apply \
  --dir file://db/migrations/sharding_${db_id} \
  --url "sqlite://server/data/sharding_db_${db_id}.db"
    ↓
各データベースに適切なテーブル範囲のみが作成される
```

## 3. コンポーネント設計

### 3.1 スキーマ定義ファイルの分割

#### 3.1.1 _schema.hclの構造
各データベース用のスキーマディレクトリに、スキーマ定義のみを含むファイルを作成：

```hcl
// db/schema/sharding_1/_schema.hcl
schema "main" {
}
```

#### 3.1.2 users.hclの構造
各データベース用のスキーマディレクトリに、該当するテーブル範囲のみを含むファイルを作成：

**例: `db/schema/sharding_1/users.hcl`**
```hcl
// Users テーブル（sharding_db_1.db用: users_000 〜 users_007）

table "users_000" {
  schema = schema.main
  column "id" {
    null = false
    type = integer
  }
  column "name" {
    null = false
    type = text
  }
  column "email" {
    null = false
    type = text
  }
  column "created_at" {
    null = false
    type = datetime
  }
  column "updated_at" {
    null = false
    type = datetime
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_users_000_email" {
    unique  = true
    columns = [column.email]
  }
}

// users_001 〜 users_007 も同様に定義
```

#### 3.1.3 posts.hclの構造
各データベース用のスキーマディレクトリに、該当するテーブル範囲のみを含むファイルを作成：

**例: `db/schema/sharding_1/posts.hcl`**
```hcl
// Posts テーブル（sharding_db_1.db用: posts_000 〜 posts_007）

table "posts_000" {
  schema = schema.main
  column "id" {
    null = false
    type = integer
  }
  column "user_id" {
    null = false
    type = integer
  }
  column "title" {
    null = false
    type = text
  }
  column "content" {
    null = false
    type = text
  }
  column "created_at" {
    null = false
    type = datetime
  }
  column "updated_at" {
    null = false
    type = datetime
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_posts_000_user_id" {
    columns = [column.user_id]
  }
  index "idx_posts_000_created_at" {
    columns = [column.created_at]
  }
}

// posts_001 〜 posts_007 も同様に定義
```

#### 3.1.4 Atlasの複数ファイル読み込み
- Atlasは、ディレクトリを指定すると、そのディレクトリ内のすべての`.hcl`ファイルを自動的に読み込む
- `src = "file://db/schema/sharding_1"`のようにディレクトリを指定することで、`_schema.hcl`, `users.hcl`, `posts.hcl`が自動的に結合される
- ファイル名の順序は関係なく、すべてのファイルが結合される

### 3.2 Atlas設定ファイルの修正

#### 3.2.1 config/develop/atlas.hclの修正
```hcl
// 開発環境用Atlas設定ファイル

// マスターデータベース用環境（既存、維持）
env "master" {
  src = "file://db/schema/master.hcl"
  url = "sqlite://server/data/master.db"
  dev = "sqlite://file?mode=memory"

  migration {
    dir = "file://db/migrations/master"
  }
}

// シャーディングDB 1（修正）
env "sharding_1" {
  src = "file://db/schema/sharding_1"  // ディレクトリを指定
  url = "sqlite://server/data/sharding_db_1.db"
  dev = "sqlite://file?mode=memory"

  migration {
    dir = "file://db/migrations/sharding_1"  // 個別のマイグレーションディレクトリ
  }
}

// シャーディングDB 2（修正）
env "sharding_2" {
  src = "file://db/schema/sharding_2"
  url = "sqlite://server/data/sharding_db_2.db"
  dev = "sqlite://file?mode=memory"

  migration {
    dir = "file://db/migrations/sharding_2"
  }
}

// シャーディングDB 3（修正）
env "sharding_3" {
  src = "file://db/schema/sharding_3"
  url = "sqlite://server/data/sharding_db_3.db"
  dev = "sqlite://file?mode=memory"

  migration {
    dir = "file://db/migrations/sharding_3"
  }
}

// シャーディングDB 4（修正）
env "sharding_4" {
  src = "file://db/schema/sharding_4"
  url = "sqlite://server/data/sharding_db_4.db"
  dev = "sqlite://file?mode=memory"

  migration {
    dir = "file://db/migrations/sharding_4"
  }
}
```

#### 3.2.2 config/staging/atlas.hclとconfig/production/atlas.hclの修正
- 同様の修正を`config/staging/atlas.hcl`と`config/production/atlas.hcl`にも適用
- データベースURLは環境に応じて変更（PostgreSQL/MySQL想定）

### 3.3 マイグレーション適用スクリプトの修正

#### 3.3.1 scripts/migrate.shの修正
```bash
#!/bin/bash
# マイグレーション適用スクリプト (Atlas版)
# 使用方法: ./scripts/migrate.sh [master|sharding|all]

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
DATA_DIR="$PROJECT_ROOT/server/data"

# データディレクトリの作成
mkdir -p "$DATA_DIR"

# マスターグループのマイグレーション（既存、維持）
migrate_master() {
    echo "Migrating master group..."
    local master_db="$DATA_DIR/master.db"

    atlas migrate apply \
        --dir "file://$PROJECT_ROOT/db/migrations/master" \
        --url "sqlite://$master_db"

    echo "Master group migration applied."
}

# シャーディンググループのマイグレーション（修正）
migrate_sharding() {
    echo "Migrating sharding group..."

    # 各シャーディングDBにマイグレーションを適用
    for db_id in 1 2 3 4; do
        local sharding_db="$DATA_DIR/sharding_db_${db_id}.db"
        echo "  Migrating sharding_db_${db_id}..."

        atlas migrate apply \
            --dir "file://$PROJECT_ROOT/db/migrations/sharding_${db_id}" \
            --url "sqlite://$sharding_db"
    done

    echo "Sharding group migration applied."
}

# メイン処理
case "${1:-all}" in
    master)
        migrate_master
        ;;
    sharding)
        migrate_sharding
        ;;
    all)
        migrate_master
        migrate_sharding
        ;;
    *)
        echo "Usage: $0 [master|sharding|all]"
        exit 1
        ;;
esac

echo "All migrations applied successfully!"
```

### 3.4 マイグレーションファイルの生成

#### 3.4.1 マイグレーション生成コマンド
各データベース用のスキーマディレクトリから、Atlasでマイグレーションを生成：

```bash
# sharding_db_1.db用のマイグレーション生成
atlas migrate diff initial \
  --dir file://db/migrations/sharding_1 \
  --to file://db/schema/sharding_1 \
  --dev-url "sqlite://file?mode=memory"

# sharding_db_2.db用のマイグレーション生成
atlas migrate diff initial \
  --dir file://db/migrations/sharding_2 \
  --to file://db/schema/sharding_2 \
  --dev-url "sqlite://file?mode=memory"

# sharding_db_3.db用のマイグレーション生成
atlas migrate diff initial \
  --dir file://db/migrations/sharding_3 \
  --to file://db/schema/sharding_3 \
  --dev-url "sqlite://file?mode=memory"

# sharding_db_4.db用のマイグレーション生成
atlas migrate diff initial \
  --dir file://db/migrations/sharding_4 \
  --to file://db/schema/sharding_4 \
  --dev-url "sqlite://file?mode=memory"
```

#### 3.4.2 マイグレーションファイルの内容
各マイグレーションファイルには、該当するテーブル範囲のみのCREATE文が含まれる：

**例: `db/migrations/sharding_1/20251226074934_initial.sql`**
```sql
-- Create "users_000" table
CREATE TABLE `users_000` (
  `id` integer NOT NULL,
  `name` text NOT NULL,
  `email` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
CREATE UNIQUE INDEX `idx_users_000_email` ON `users_000` (`email`);

-- users_001 〜 users_007 も同様に定義

-- Create "posts_000" table
CREATE TABLE `posts_000` (
  `id` integer NOT NULL,
  `user_id` integer NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
);
CREATE INDEX `idx_posts_000_user_id` ON `posts_000` (`user_id`);
CREATE INDEX `idx_posts_000_created_at` ON `posts_000` (`created_at`);

-- posts_001 〜 posts_007 も同様に定義
```

### 3.5 既存データベースのクリーンアップ

#### 3.5.1 開発環境でのデータベース削除
開発環境（develop）では、既存のシャーディングデータベースファイルを削除してからマイグレーションを適用する：

```bash
# 既存データベースファイルを削除（develop環境のみ）
rm server/data/sharding_db_1.db
rm server/data/sharding_db_2.db
rm server/data/sharding_db_3.db
rm server/data/sharding_db_4.db

# 修正後のマイグレーションを適用
./scripts/migrate.sh sharding
```

#### 3.5.2 ステージング・本番環境での対応
ステージング・本番環境では、データを安易に削除できないため、適切な手順に従って対応する必要がある。

## 4. エラーハンドリング

### 4.1 スキーマファイルの分割時のエラー
- **問題**: 既存の`db/schema/sharding.hcl`からテーブルを抽出する際に、テーブル範囲を間違える可能性
- **対策**: テーブル範囲を正確に確認し、各スキーマファイルに必要なテーブルのみが含まれていることを検証

### 4.2 マイグレーション生成時のエラー
- **問題**: Atlasでマイグレーションを生成する際に、スキーマディレクトリの参照が正しくない可能性
- **対策**: スキーマディレクトリのパスを正確に指定し、ディレクトリ内のすべての`.hcl`ファイルが読み込まれることを確認

### 4.3 マイグレーション適用時のエラー
- **問題**: マイグレーション適用時に、既存のデータベースにテーブルが存在する場合のエラー
- **対策**: 開発環境では既存データベースファイルを削除してからマイグレーションを適用する。ステージング・本番環境では適切な手順に従う

### 4.4 Atlas設定ファイルの参照エラー
- **問題**: Atlas設定ファイルでスキーマディレクトリを参照する際に、パスが正しくない可能性
- **対策**: 相対パスを正確に指定し、ディレクトリが存在することを確認

## 5. テスト戦略

### 5.1 スキーマファイルの検証
- 各スキーマディレクトリに`_schema.hcl`, `users.hcl`, `posts.hcl`が存在することを確認
- 各スキーマファイルに必要なテーブルのみが含まれていることを確認
- テーブル範囲が正しいことを確認（sharding_1: 000-007, sharding_2: 008-015, sharding_3: 016-023, sharding_4: 024-031）

### 5.2 マイグレーションファイルの検証
- 各マイグレーションディレクトリにマイグレーションファイルが存在することを確認
- 各マイグレーションファイルに必要なテーブルのみが含まれていることを確認
- `atlas.sum`が各マイグレーションディレクトリに生成されていることを確認

### 5.3 マイグレーション適用の検証
- 修正後のマイグレーション適用スクリプトが正常に動作することを確認
- 各データベースに適切なテーブルのみが作成されることを確認
- 不要なテーブル（範囲外のテーブル）が作成されていないことを確認

### 5.4 データベース構造の検証
- `sharding_db_1.db`にテーブル000-007のみが作成されていることを確認
- `sharding_db_2.db`にテーブル008-015のみが作成されていることを確認
- `sharding_db_3.db`にテーブル016-023のみが作成されていることを確認
- `sharding_db_4.db`にテーブル024-031のみが作成されていることを確認

### 5.5 Atlas設定ファイルの検証
- 各環境設定ファイル（develop, staging, production）が正しく修正されていることを確認
- 各データベース用の環境設定で、個別のスキーマディレクトリとマイグレーションディレクトリを参照していることを確認

## 6. 実装手順

### 6.1 スキーマファイルの分割
1. 既存の`db/schema/sharding.hcl`を分析し、各データベース用のテーブル範囲を特定
2. 各データベース用のスキーマディレクトリを作成（`db/schema/sharding_1/` 〜 `sharding_4/`）
3. 各ディレクトリに`_schema.hcl`を作成（`schema "main" {}`のみ）
4. 各ディレクトリに`users.hcl`を作成（該当するテーブル範囲のみ）
5. 各ディレクトリに`posts.hcl`を作成（該当するテーブル範囲のみ）
6. 既存の`db/schema/sharding.hcl`を削除

### 6.2 マイグレーションディレクトリの分割
1. 各データベース用のマイグレーションディレクトリを作成（`db/migrations/sharding_1/` 〜 `sharding_4/`）
2. 各データベース用のスキーマディレクトリから、Atlasでマイグレーションを生成
3. 各マイグレーションディレクトリに`atlas.sum`が生成されることを確認
4. 既存の`db/migrations/sharding/`ディレクトリを削除

### 6.3 マイグレーション適用スクリプトの修正
1. `scripts/migrate.sh`を修正し、各データベースごとのマイグレーションディレクトリを参照するように変更
2. 修正後のスクリプトが正常に動作することを確認

### 6.4 Atlas設定ファイルの修正
1. `config/develop/atlas.hcl`を修正し、各データベースごとのスキーマディレクトリとマイグレーションディレクトリを参照するように変更
2. `config/staging/atlas.hcl`を修正
3. `config/production/atlas.hcl`を修正

### 6.5 既存データベースのクリーンアップ（開発環境のみ）
1. 開発環境では、既存のシャーディングデータベースファイル（sharding_db_1.db ～ sharding_db_4.db）を削除
2. 修正後のマイグレーションを適用して、正しいテーブル構造で再作成
3. 各データベースに適切なテーブルのみが作成されることを確認
4. ステージング・本番環境では、適切な手順に従って対応する

### 6.6 ドキュメントの更新
1. `docs/atlas-operations.md`に新しいディレクトリ構造とマイグレーション適用方法を記載
2. 変更内容が適切にドキュメント化されていることを確認

## 7. 参考情報

### 7.1 Atlasの複数ファイル読み込み機能
- Atlasは、ディレクトリを指定すると、そのディレクトリ内のすべての`.hcl`ファイルを自動的に読み込む
- ファイル名の順序は関係なく、すべてのファイルが結合される
- スキーマ定義、テーブル定義、インデックス定義などが複数のファイルに分かれていても、自動的に結合される

### 7.2 既存ドキュメント
- `docs/atlas-operations.md`: Atlas運用ガイド
- `docs/Sharding.md`: シャーディングの詳細仕様
- `.kiro/specs/0012-sharding/`: シャーディング実装の仕様書
- `.kiro/specs/0014-db-atlas/`: Atlas導入の仕様書

### 7.3 技術スタック
- **Atlas**: データベーススキーマ管理ツール（複数HCLファイルの自動読み込みをサポート）
- **SQLite**: 開発環境のデータベース
- **Bash**: マイグレーション適用スクリプト

