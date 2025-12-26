# Atlas導入設計書

## 1. 概要

### 1.1 設計の目的
要件定義書に基づき、Atlas CLIを使用したデータベース管理システムの詳細設計を定義する。既存の手動マイグレーション管理から、宣言的なスキーマ定義（HCL形式）による管理に移行し、バージョン管理型のマイグレーション運用を実現する。

### 1.2 設計の範囲
- Atlas CLIの導入とセットアップ
- スキーマ定義ファイル（HCL形式）の作成
- 環境別設定ファイルの作成
- マイグレーション生成と適用のワークフロー構築
- 運用ドキュメントの作成
- 既存システムとの統合確認
- 運用実験の実施

### 1.3 設計方針
- **宣言的なスキーマ定義**: HCL形式によるスキーマ定義により、現在のスキーマ状態を明確に管理
- **バージョン管理型の運用**: Atlasによるマイグレーション履歴の自動管理
- **環境別設定の分離**: 各環境（develop, staging, production）での設定を分離
- **既存システムとの互換性**: Atlasで構築したデータベースが既存のAPIサーバー、クライアント、管理画面と互換性を持つ
- **段階的な移行**: 既存のマイグレーションファイルを参考にしつつ、新規にスキーマを定義

## 2. アーキテクチャ設計

### 2.1 ディレクトリ構造

#### 2.1.1 変更前の構造
```
.
├── db/
│   └── migrations/
│       ├── master/
│       │   ├── 001_init.sql
│       │   ├── 002_goadmin.sql
│       │   ├── 003_menu.sql
│       │   └── 004_api_key_menu.sql
│       └── sharding/
│           ├── templates/
│           │   ├── users.sql.template
│           │   └── posts.sql.template
│           └── generated/
│               ├── users_000.sql ～ users_031.sql
│               └── posts_000.sql ～ posts_031.sql
├── config/
│   ├── develop/
│   │   └── database.yaml
│   ├── staging/
│   │   └── database.yaml
│   └── production/
│       └── database.yaml
└── scripts/
    └── migrate.sh
```

#### 2.1.2 変更後の構造
```
.
├── db/
│   ├── schema/                    # 新規: スキーマ定義ファイル
│   │   ├── master.hcl            # マスターデータベーススキーマ
│   │   └── sharding.hcl          # シャーディングデータベーススキーマ
│   └── migrations/                # Atlasマイグレーションファイル
│       ├── master/
│       │   └── (Atlas生成ファイル)
│       └── sharding/
│           └── (Atlas生成ファイル)
├── config/
│   ├── develop/
│   │   ├── database.yaml         # 既存（維持）
│   │   └── atlas.hcl             # 新規: Atlas設定
│   ├── staging/
│   │   ├── database.yaml         # 既存（維持）
│   │   └── atlas.hcl             # 新規: Atlas設定
│   └── production/
│       ├── database.yaml         # 既存（維持）
│       └── atlas.hcl             # 新規: Atlas設定
├── scripts/
│   └── migrate.sh                # Atlas対応に更新、または新規スクリプト
└── docs/
    └── Atlas.md                   # 新規: Atlas運用ドキュメント
```

### 2.2 ファイル構成

#### 2.2.1 スキーマ定義ファイル
- **`db/schema/master.hcl`**: マスターデータベースの全テーブル定義
  - `news`テーブル
  - GoAdmin関連テーブル（goadmin_menu, goadmin_operation_log, goadmin_site, goadmin_permissions, goadmin_roles, goadmin_role_menu, goadmin_role_permissions, goadmin_role_users, goadmin_session, goadmin_user_permissions, goadmin_users）
  
- **`db/schema/sharding.hcl`**: シャーディングデータベースのテーブル定義
  - `users`テーブル（32分割: users_000 ～ users_031）
  - `posts`テーブル（32分割: posts_000 ～ posts_031）

#### 2.2.2 環境別設定ファイル
- **`config/develop/atlas.hcl`**: 開発環境用Atlas設定
  - データソース: SQLite（`./data/master.db`, `./data/sharding_db_*.db`）
  - マイグレーションディレクトリ設定
  
- **`config/staging/atlas.hcl`**: ステージング環境用Atlas設定
  - データソース: PostgreSQL/MySQL
  - マイグレーションディレクトリ設定
  
- **`config/production/atlas.hcl`**: 本番環境用Atlas設定
  - データソース: PostgreSQL/MySQL
  - マイグレーションディレクトリ設定

### 2.3 データフロー

#### 2.3.1 マイグレーション生成フロー
```
スキーマ定義ファイル（.hcl）
    ↓
atlas migrate diff
    ↓
マイグレーションファイル（.sql）生成
    ↓
Gitでバージョン管理
```

#### 2.3.2 マイグレーション適用フロー（開発環境）
```
マイグレーションファイル（.sql）
    ↓
atlas migrate apply
    ↓
データベースに適用
    ↓
マイグレーション履歴を記録（atlas_schema_migrationsテーブル）
```

#### 2.3.3 マイグレーション適用フロー（本番環境）
```
マイグレーションファイル（.sql）
    ↓
atlas migrate apply --dry-run
    ↓
SQL生成
    ↓
SQL確認
    ↓
本番環境に適用（手動または自動化スクリプト）
```

## 3. コンポーネント設計

### 3.1 Atlas CLIの設定

#### 3.1.1 Atlas CLIのインストール
- **公式サイト**: https://atlasgo.io/
- **インストール場所**: システムのPATHにインストールされる（プロジェクト内には配置しない）
- **インストール方法**: 
  - **macOS**: 
    - Homebrew: `brew install ariga/tap/atlas`
    - または公式インストールスクリプト: `curl -sSf https://atlasgo.sh | sh`
  - **Linux**: 
    - 公式インストールスクリプト: `curl -sSf https://atlasgo.sh | sh`
    - または手動でバイナリをダウンロードして`/usr/local/bin/`に配置
  - **Windows**: 
    - Scoop: `scoop install atlas`
    - または手動でバイナリをダウンロードしてPATHに追加
- **インストール後の確認**: 
  - `atlas version`コマンドでバージョンを確認
  - どのディレクトリからでも`atlas`コマンドが実行できることを確認
- **注意**: Atlas CLIはシステム全体で使用するツールとしてPATHにインストールされるため、プロジェクトディレクトリ内には配置しない

#### 3.1.2 Atlasの初期化
```bash
# プロジェクトルートで実行
atlas migrate init --dir db/migrations/master
atlas migrate init --dir db/migrations/sharding
```

### 3.2 スキーマ定義ファイルの設計

#### 3.2.1 マスターデータベーススキーマ（master.hcl）

**newsテーブル**:
```hcl
table "news" {
  schema = schema.main
  column "id" {
    null = false
    type = integer
    auto_increment = true
  }
  column "title" {
    null = false
    type = text
  }
  column "content" {
    null = false
    type = text
  }
  column "author_id" {
    null = true
    type = integer
  }
  column "published_at" {
    null = true
    type = datetime
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
  index "idx_news_published_at" {
    columns = [column.published_at]
  }
  index "idx_news_author_id" {
    columns = [column.author_id]
  }
}
```

**GoAdmin関連テーブル**:
- `goadmin_menu`: メニューテーブル
- `goadmin_operation_log`: 操作ログテーブル
- `goadmin_site`: サイト設定テーブル
- `goadmin_permissions`: 権限テーブル
- `goadmin_roles`: ロールテーブル
- `goadmin_role_menu`: ロール-メニュー関連テーブル
- `goadmin_role_permissions`: ロール-権限関連テーブル
- `goadmin_role_users`: ロール-ユーザー関連テーブル
- `goadmin_session`: セッションテーブル
- `goadmin_user_permissions`: ユーザー-権限関連テーブル
- `goadmin_users`: 管理者ユーザーテーブル

#### 3.2.2 シャーディングデータベーススキーマ（sharding.hcl）

**usersテーブル（32分割）**:
```hcl
# users_000 ～ users_031 の32テーブルを定義
# 各テーブルは同じ構造を持つ

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
    columns = [column.email]
    unique = true
  }
}

# users_001 ～ users_031 も同様に定義
# （実際の実装では、テンプレートやループを使用して定義する可能性がある）
```

**postsテーブル（32分割）**:
```hcl
# posts_000 ～ posts_031 の32テーブルを定義
# 各テーブルは同じ構造を持つ

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
  foreign_key "fk_posts_000_user_id" {
    columns = [column.user_id]
    ref_columns = [table.users_000.column.id]
    on_delete = CASCADE
  }
  index "idx_posts_000_user_id" {
    columns = [column.user_id]
  }
  index "idx_posts_000_created_at" {
    columns = [column.created_at]
  }
}

# posts_001 ～ posts_031 も同様に定義
```

### 3.3 環境別設定ファイルの設計

#### 3.3.1 開発環境設定（config/develop/atlas.hcl）

```hcl
# 開発環境用Atlas設定

env "develop" {
  # マスターデータベース
  src = "file://db/schema/master.hcl"
  url = "sqlite://./server/data/master.db"
  dev = "sqlite://./server/data/master.db"
  
  # マイグレーションディレクトリ
  migration {
    dir = "file://db/migrations/master"
  }
}

# シャーディングデータベース用の設定
# 各データベースに対して個別に設定するか、スクリプトで一括適用する
```

#### 3.3.2 ステージング環境設定（config/staging/atlas.hcl）

```hcl
# ステージング環境用Atlas設定

env "staging" {
  # マスターデータベース
  src = "file://db/schema/master.hcl"
  url = "postgres://user:password@host:5432/master_db?sslmode=disable"
  dev = "postgres://user:password@host:5432/master_db_dev?sslmode=disable"
  
  # マイグレーションディレクトリ
  migration {
    dir = "file://db/migrations/master"
  }
}
```

#### 3.3.3 本番環境設定（config/production/atlas.hcl）

```hcl
# 本番環境用Atlas設定

env "production" {
  # マスターデータベース
  src = "file://db/schema/master.hcl"
  url = "postgres://user:password@host:5432/master_db?sslmode=require"
  dev = "postgres://user:password@host:5432/master_db_dev?sslmode=require"
  
  # マイグレーションディレクトリ
  migration {
    dir = "file://db/migrations/master"
  }
}
```

### 3.4 マイグレーション管理の設計

#### 3.4.1 マイグレーション生成コマンド

**開発環境でのマイグレーション生成**:
```bash
# マスターデータベース
atlas migrate diff \
  --dir file://db/migrations/master \
  --to file://db/schema/master.hcl \
  --dev-url sqlite://./server/data/master.db

# シャーディングデータベース
atlas migrate diff \
  --dir file://db/migrations/sharding \
  --to file://db/schema/sharding.hcl \
  --dev-url sqlite://./server/data/sharding_db_1.db
```

#### 3.4.2 マイグレーション適用コマンド

**開発環境でのマイグレーション適用**:
```bash
# マスターデータベース
atlas migrate apply \
  --dir file://db/migrations/master \
  --url sqlite://./server/data/master.db

# シャーディングデータベース（4つ）
for i in {1..4}; do
  atlas migrate apply \
    --dir file://db/migrations/sharding \
    --url sqlite://./server/data/sharding_db_${i}.db
done
```

**本番環境でのSQL生成**:
```bash
# マスターデータベース
atlas migrate apply \
  --dir file://db/migrations/master \
  --url postgres://user:password@host:5432/master_db \
  --dry-run

# 生成されたSQLを確認してから適用
```

### 3.5 スクリプトの設計

#### 3.5.1 マイグレーション適用スクリプト（scripts/migrate.sh）

```bash
#!/bin/bash
# Atlas対応マイグレーション適用スクリプト

set -e

ENV=${1:-develop}
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# 環境別設定の読み込み
ATLAS_CONFIG="$PROJECT_ROOT/config/${ENV}/atlas.hcl"

# マスターデータベースのマイグレーション適用
atlas migrate apply \
  --config "$ATLAS_CONFIG" \
  --env master

# シャーディングデータベースのマイグレーション適用
for i in {1..4}; do
  atlas migrate apply \
    --config "$ATLAS_CONFIG" \
    --env "sharding_${i}"
done
```

## 4. データモデル設計

### 4.1 マスターデータベースのテーブル構造

#### 4.1.1 newsテーブル
- **目的**: ニュース記事の管理
- **カラム**:
  - `id`: INTEGER PRIMARY KEY AUTOINCREMENT
  - `title`: TEXT NOT NULL
  - `content`: TEXT NOT NULL
  - `author_id`: INTEGER
  - `published_at`: DATETIME
  - `created_at`: DATETIME NOT NULL
  - `updated_at`: DATETIME NOT NULL
- **インデックス**:
  - `idx_news_published_at`: published_at
  - `idx_news_author_id`: author_id

#### 4.1.2 GoAdmin関連テーブル
- **goadmin_menu**: メニュー管理
- **goadmin_operation_log**: 操作ログ
- **goadmin_site**: サイト設定
- **goadmin_permissions**: 権限管理
- **goadmin_roles**: ロール管理
- **goadmin_role_menu**: ロール-メニュー関連
- **goadmin_role_permissions**: ロール-権限関連
- **goadmin_role_users**: ロール-ユーザー関連
- **goadmin_session**: セッション管理
- **goadmin_user_permissions**: ユーザー-権限関連
- **goadmin_users**: 管理者ユーザー

### 4.2 シャーディングデータベースのテーブル構造

#### 4.2.1 usersテーブル（32分割）
- **目的**: ユーザー情報の管理（シャーディング）
- **テーブル名**: `users_000` ～ `users_031`
- **カラム**:
  - `id`: INTEGER PRIMARY KEY
  - `name`: TEXT NOT NULL
  - `email`: TEXT NOT NULL UNIQUE
  - `created_at`: DATETIME NOT NULL
  - `updated_at`: DATETIME NOT NULL
- **インデックス**:
  - `idx_users_{suffix}_email`: email (UNIQUE)

#### 4.2.2 postsテーブル（32分割）
- **目的**: 投稿情報の管理（シャーディング）
- **テーブル名**: `posts_000` ～ `posts_031`
- **カラム**:
  - `id`: INTEGER PRIMARY KEY
  - `user_id`: INTEGER NOT NULL
  - `title`: TEXT NOT NULL
  - `content`: TEXT NOT NULL
  - `created_at`: DATETIME NOT NULL
  - `updated_at`: DATETIME NOT NULL
- **外部キー**:
  - `fk_posts_{suffix}_user_id`: user_id → users_{suffix}.id (ON DELETE CASCADE)
- **インデックス**:
  - `idx_posts_{suffix}_user_id`: user_id
  - `idx_posts_{suffix}_created_at`: created_at

## 5. ワークフロー設計

### 5.1 初期セットアップワークフロー

```
1. Atlas CLIのインストール
   ↓
2. スキーマ定義ファイルの作成（db/schema/master.hcl, db/schema/sharding.hcl）
   ↓
3. 環境別設定ファイルの作成（config/{env}/atlas.hcl）
   ↓
4. 初期マイグレーションの生成（atlas migrate diff）
   ↓
5. マイグレーションの適用（atlas migrate apply）
   ↓
6. 既存システムとの統合確認
```

### 5.2 スキーマ変更ワークフロー

```
1. スキーマ定義ファイル（.hcl）を編集
   ↓
2. マイグレーションファイルを生成（atlas migrate diff）
   ↓
3. 生成されたマイグレーションファイルを確認
   ↓
4. Gitでコミット
   ↓
5. 開発環境でマイグレーション適用（atlas migrate apply）
   ↓
6. テスト実行
   ↓
7. 本番環境でSQL生成（atlas migrate apply --dry-run）
   ↓
8. SQL確認後、本番環境に適用
```

### 5.3 イレギュラーケースの対処ワークフロー

```
1. データベースに直接SQLを適用
   ↓
2. マイグレーション履歴の確認（atlas migrate status）
   ↓
3. マイグレーションハッシュの確認（atlas migrate hash）
   ↓
4. ベースラインの設定（atlas migrate baseline）
   ↓
5. マイグレーション履歴とスキーマの整合性を確認
   ↓
6. 必要に応じてマイグレーション履歴を修正
```

## 6. エラーハンドリング

### 6.1 マイグレーション適用時のエラー

#### 6.1.1 スキーマ不一致エラー
- **原因**: マイグレーション履歴と実際のスキーマが不一致
- **対処**: `atlas migrate baseline`でベースラインを設定

#### 6.1.2 マイグレーションハッシュ不一致エラー
- **原因**: マイグレーションファイルが変更された
- **対処**: `atlas migrate hash`でハッシュを確認し、必要に応じて修正

### 6.2 スキーマ定義ファイルのエラー

#### 6.2.1 HCL構文エラー
- **原因**: HCL形式の構文エラー
- **対処**: `atlas schema validate`でスキーマを検証

#### 6.2.2 スキーマ検証エラー
- **原因**: スキーマ定義がデータベースの制約に違反
- **対処**: エラーメッセージを確認し、スキーマ定義を修正

## 7. テスト戦略

### 7.1 スキーマ定義の検証

#### 7.1.1 スキーマ検証コマンド
```bash
# スキーマ定義ファイルの検証
atlas schema validate --schema file://db/schema/master.hcl
atlas schema validate --schema file://db/schema/sharding.hcl
```

### 7.2 マイグレーションのテスト

#### 7.2.1 開発環境でのテスト
- マイグレーション生成のテスト
- マイグレーション適用のテスト
- ロールバックのテスト（必要に応じて）

#### 7.2.2 既存システムとの統合テスト
- APIサーバーの動作確認
- クライアント（Next.js）の動作確認
- 管理画面（GoAdmin）の動作確認

### 7.3 運用実験

#### 7.3.1 実験シナリオ
1. **0からのデータベースの初期化**
2. **master側の操作**
   - テーブルの追加
   - テーブルにカラムを追加
   - テーブルのデータを更新
   - テーブルを削除
3. **sharding側の操作**
   - テーブルの追加
   - テーブルにカラムを追加
   - テーブルのデータを更新
   - テーブルを削除
4. **イレギュラーケースのシナリオ**
   - 直接SQLを適用した後の作業
   - マイグレーションハッシュの確認
   - マイグレーション履歴とスキーマの整合性を取る方法

## 8. 運用ドキュメント

### 8.1 ドキュメント構成

#### 8.1.1 docs/Atlas.mdの構成
1. **Atlasの基本的な使用方法**
   - Atlas CLIのインストール
   - スキーマ定義ファイルの作成
   - マイグレーション生成と適用

2. **環境別のマイグレーション適用手順**
   - 開発環境での運用
   - ステージング環境での運用
   - 本番環境での運用

3. **ケース別の運用方法**
   - スキーマ変更の手順
   - ロールバックの手順
   - データ移行の手順

4. **イレギュラーケースの対処方法**
   - 直接SQL適用後の作業
   - マイグレーションハッシュの確認
   - マイグレーション履歴とスキーマの整合性を取る方法

5. **トラブルシューティング**
   - よくあるエラーと対処方法
   - マイグレーション履歴の修復方法

## 9. 既存システムとの統合

### 9.1 既存システムとの互換性

#### 9.1.1 APIサーバーとの統合
- Atlasで構築したデータベースが既存のAPIサーバーと互換性を持つことを確認
- 既存のAPIエンドポイントが正常に動作することを確認

#### 9.1.2 クライアント（Next.js）との統合
- 既存のクライアントが正常に動作することを確認
- APIエンドポイントとの通信が正常に動作することを確認

#### 9.1.3 管理画面（GoAdmin）との統合
- 既存の管理画面が正常に動作することを確認
- GoAdmin関連テーブルが正常に動作することを確認

### 9.2 統合確認の手順

#### 9.2.1 確認項目
1. Atlasでデータベースを構築
2. APIサーバーを起動
3. クライアント（Next.js）を起動
4. 管理画面（GoAdmin）を起動
5. 各システムの動作確認
6. 統合確認結果の文書化

## 10. 参考情報

### 10.1 Atlas公式ドキュメント
- Atlas公式サイト: https://atlasgo.io/
- Atlas GitHub: https://github.com/ariga/atlas
- Atlas CLI リファレンス: https://atlasgo.io/cli-reference

### 10.2 既存ドキュメント
- `README.md`: プロジェクト概要とセットアップ手順
- `docs/Sharding.md`: シャーディングの詳細仕様
- `db/migrations/master/*.sql`: 既存のマスターデータベースマイグレーション（参考）
- `db/migrations/sharding/templates/*.sql.template`: 既存のシャーディングテンプレート（参考）

