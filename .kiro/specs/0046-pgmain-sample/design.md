# generate-sample-dataコマンドPostgreSQL対応設計書

## 1. 概要

### 1.1 設計の目的
要件定義書に基づき、`server/cmd/generate-sample-data`コマンドでPostgreSQLを利用するように修正するための詳細設計を定義する。既存のコードは既にPostgreSQL対応のコードパスを使用しているため、主に確認作業と必要に応じた修正、ドキュメント更新を行う。

### 1.2 設計の範囲
- `server/cmd/generate-sample-data/main.go`のPostgreSQL対応確認
- SQLite用ライブラリのインポート確認（存在しないことを確認）
- SQLite用処理分岐の確認（存在しないことを確認）
- SQL構文の確認（GORMのプレースホルダー構文がPostgreSQLで動作することを確認）
- テストコードの確認・修正（存在する場合）
- ドキュメントの更新

### 1.3 設計方針
- **既存システムとの統合**: 既存のコードは既に`config.Load()`と`db.NewGroupManager(cfg)`を使用しており、PostgreSQL対応のコードパスを使用している
- **設定ファイルベース**: 環境変数ではなく設定ファイル（`config/{env}/database.yaml`）から接続情報を読み込む（Issue #87でPostgreSQL設定が有効化されている想定）
- **環境別対応**: `APP_ENV`環境変数で環境を指定（develop/staging/production）
- **GORMのプレースホルダー**: GORMのプレースホルダー構文（`?`）はPostgreSQLドライバーで自動的に`$1`, `$2`に変換されるため、既存のコードで問題なく動作する
- **既存コードの維持**: 既存のコード構造を維持し、確認と必要に応じた修正のみを行う

## 2. アーキテクチャ設計

### 2.1 既存アーキテクチャの分析

#### 2.1.1 現在の構成
- **コマンド実装**: `server/cmd/generate-sample-data/main.go`で`config.Load()`と`db.NewGroupManager(cfg)`を使用してデータベース接続を行っている
- **データベース接続**: `server/internal/db/connection.go`で既にPostgreSQLドライバー（`gorm.io/driver/postgres`）がサポートされている
- **設定ファイル**: `config/{env}/database.yaml`から接続情報を読み込む（Issue #87でPostgreSQL設定が有効化されている想定）
- **SQLite用ライブラリ**: `server/cmd/generate-sample-data/main.go`にはSQLite用ライブラリのインポートが存在しない
- **SQLite用処理分岐**: `server/cmd/generate-sample-data/main.go`にはSQLite用処理分岐が存在しない
- **SQL構文**: GORMのプレースホルダー構文（`?`）を使用しており、GORMがPostgreSQLドライバーで自動的に`$1`, `$2`に変換する
- **データベース構成**:
  - マスターデータベース: 1台（`webdb_master`）
  - シャーディングデータベース: 4台（`webdb_sharding_1` ～ `webdb_sharding_4`）
  - **論理シャーディング数: 8**（物理DB 4台 × 2論理シャード）

#### 2.1.2 既存パターンの維持
- **設定ファイル構造**: `config/{env}/database.yaml`の構造は維持（Issue #87でPostgreSQL設定が有効化されている想定）
- **データベース接続コード**: `server/internal/db/connection.go`は既にPostgreSQLドライバーをサポートしているため、変更不要
- **コマンドの初期化**: `server/cmd/generate-sample-data/main.go`の初期化処理は変更不要
- **データ生成処理**: 既存のデータ生成処理（`generateDmUsers`, `generateDmPosts`, `generateDmNews`）は変更不要

### 2.2 システム構成図

```
┌─────────────────────────────────────────────────────────────┐
│      generate-sample-dataコマンド (server/cmd/generate-sample-data/main.go) │
│                                                              │
│  ┌────────────────────────────────────────────────────┐   │
│  │ config.Load()                                       │   │
│  │   ↓                                                 │   │
│  │ config/{env}/database.yaml                         │   │
│  │   - master: 1台 (webdb_master)                      │   │
│  │   - sharding: 8つの論理シャード (物理DB 4台)        │   │
│  └────────────────────────────────────────────────────┘   │
│                          │                                  │
│                          ▼                                  │
│  ┌────────────────────────────────────────────────────┐   │
│  │ db.NewGroupManager(cfg)                            │   │
│  │   ↓                                                 │   │
│  │ server/internal/db/connection.go                   │   │
│  │   - PostgreSQLドライバー: サポート済み              │   │
│  │   - SQLiteドライバー: 使用していない                │   │
│  └────────────────────────────────────────────────────┘   │
│                          │                                  │
│                          ▼                                  │
│  ┌────────────────────────────────────────────────────┐   │
│  │ データ生成処理                                       │   │
│  │   - generateDmUsers(): shardingグループ              │   │
│  │   - generateDmPosts(): shardingグループ              │   │
│  │   - generateDmNews(): masterグループ                 │   │
│  └────────────────────────────────────────────────────┘   │
│                          │                                  │
│                          ▼                                  │
│  ┌────────────────────────────────────────────────────┐   │
│  │ PostgreSQL接続                                      │   │
│  │   - master: postgres-master:5432                    │   │
│  │   - sharding_1: postgres-sharding-1:5433            │   │
│  │   - sharding_2: postgres-sharding-2:5434            │   │
│  │   - sharding_3: postgres-sharding-3:5435            │   │
│  │   - sharding_4: postgres-sharding-4:5436            │   │
│  └────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────────┐
│              Docker Compose (docker-compose.postgres.yml)    │
│                                                              │
│  ┌──────────────────┐                                      │
│  │ postgres-master   │                                      │
│  │ (ポート: 5432)    │                                      │
│  │ DB: webdb_master  │                                      │
│  └──────────────────┘                                      │
│                                                              │
│  ┌──────────────────┐  ┌──────────────────┐              │
│  │ postgres-sharding │  │ postgres-sharding │              │
│  │ -1 (ポート: 5433) │  │ -2 (ポート: 5434) │              │
│  │ DB: webdb_sharding_1│ │ DB: webdb_sharding_2│             │
│  └──────────────────┘  └──────────────────┘              │
│                                                              │
│  ┌──────────────────┐  ┌──────────────────┐              │
│  │ postgres-sharding │  │ postgres-sharding │              │
│  │ -3 (ポート: 5435) │  │ -4 (ポート: 5436) │              │
│  │ DB: webdb_sharding_3│ │ DB: webdb_sharding_4│             │
│  └──────────────────┘  └──────────────────┘              │
└─────────────────────────────────────────────────────────────┘
```

### 2.3 ディレクトリ構造

#### 2.3.1 変更前の構造
```
.
├── config/
│   ├── develop/
│   │   └── database.yaml          # PostgreSQL設定が有効（Issue #87で対応済み）
│   ├── staging/
│   │   └── database.yaml          # PostgreSQL設定（Issue #87で対応済み）
│   └── production/
│       ├── database.yaml          # PostgreSQL設定（存在する場合、Issue #87で対応済み）
│       └── database.yaml.example  # PostgreSQL設定（Issue #87で対応済み）
├── server/
│   ├── cmd/
│   │   └── generate-sample-data/
│   │       └── main.go            # PostgreSQL対応のコードパスを使用（確認が必要）
│   └── internal/
│       └── db/
│           └── connection.go      # PostgreSQLドライバーをサポート（変更不要）
└── README.md                      # generate-sample-dataコマンドのPostgreSQL利用に関する記述が未整備
```

#### 2.3.2 変更後の構造
```
.
├── config/
│   ├── develop/
│   │   └── database.yaml          # 変更なし（Issue #87でPostgreSQL設定が有効）
│   ├── staging/
│   │   └── database.yaml          # 変更なし（Issue #87でPostgreSQL設定が有効）
│   └── production/
│       ├── database.yaml          # 変更なし（Issue #87でPostgreSQL設定が有効）
│       └── database.yaml.example  # 変更なし（Issue #87でPostgreSQL設定が有効）
├── server/
│   ├── cmd/
│   │   └── generate-sample-data/
│   │       └── main.go            # 変更なし（既にPostgreSQL対応のコードパスを使用）
│   └── internal/
│       └── db/
│           └── connection.go      # 変更なし（既にPostgreSQLドライバーをサポート）
└── README.md                      # generate-sample-dataコマンドのPostgreSQL利用に関する記述を追加
```

## 3. 詳細設計

### 3.1 コマンド実装の確認

#### 3.1.1 server/cmd/generate-sample-data/main.goの確認
- **ファイル**: `server/cmd/generate-sample-data/main.go`
- **確認項目**:
  1. **インポート文の確認**: SQLite用ライブラリのインポートが存在しないことを確認
     - `_ "github.com/mattn/go-sqlite3"`が存在しないことを確認
     - `"gorm.io/driver/sqlite"`が存在しないことを確認
  2. **処理分岐の確認**: SQLite用処理分岐が存在しないことを確認
     - `case "sqlite3":`分岐が存在しないことを確認
     - `driver = "sqlite3"`のデフォルト値設定が存在しないことを確認
  3. **データベース接続の確認**: PostgreSQL接続が確立されることを確認
     - `config.Load()`で設定ファイルから接続情報を読み込んでいることを確認
     - `db.NewGroupManager(cfg)`でPostgreSQL接続が確立されることを確認
  4. **データ生成処理の確認**: データ生成処理がPostgreSQLで正常に動作することを確認
     - `generateDmUsers`関数がPostgreSQLで正常に動作することを確認
     - `generateDmPosts`関数がPostgreSQLで正常に動作することを確認
     - `generateDmNews`関数がPostgreSQLで正常に動作することを確認

#### 3.1.2 設定ファイルの確認
- **ファイル**: `config/{env}/database.yaml`
- **確認項目**:
  1. **PostgreSQL設定の有効化**: PostgreSQL設定が有効になっていることを確認（Issue #87で対応済みの想定）
  2. **masterグループの設定**: masterグループのPostgreSQL接続情報が正しく設定されていることを確認
  3. **shardingグループの設定**: shardingグループのPostgreSQL接続情報が正しく設定されていることを確認（論理シャーディング数8）
  4. **SQLite設定の削除**: SQLite設定が削除されていることを確認

### 3.2 SQL構文の確認

#### 3.2.1 プレースホルダー構文の確認
- **対象関数**: `insertDmUsersBatch`, `insertDmPostsBatch`
- **確認内容**:
  - GORMのプレースホルダー構文（`?`）が使用されていることを確認
  - GORMがPostgreSQLドライバーで自動的に`$1`, `$2`に変換することを確認
  - 既存のコードで問題なく動作することを確認

#### 3.2.2 動的テーブル名の確認
- **対象関数**: `insertDmUsersBatch`, `insertDmPostsBatch`
- **確認内容**:
  - 動的テーブル名（`fmt.Sprintf("dm_users_%03d", tableNumber)`）がPostgreSQLで正常に動作することを確認
  - SQLインジェクション対策が適切に行われていることを確認（テーブル名は数値から生成されるため安全）

#### 3.2.3 バッチ挿入の確認
- **対象関数**: `insertDmUsersBatch`, `insertDmPostsBatch`, `insertDmNewsBatch`
- **確認内容**:
  - バッチ挿入のSQL構文がPostgreSQL対応であることを確認
  - バッチサイズ（500件）が適切であることを確認
  - エラーハンドリングが適切であることを確認

### 3.3 データ生成処理の確認

#### 3.3.1 dm_usersテーブルへのデータ生成
- **対象関数**: `generateDmUsers`
- **確認内容**:
  1. **UUIDv7の生成**: UUIDv7の生成が正常に動作することを確認
  2. **テーブル番号の計算**: テーブル番号の計算が正しく動作することを確認
     - `tableSelector.GetTableNumberFromUUID(id)`が正しく動作することを確認
  3. **シャーディング接続**: 正しいシャーディングデータベースに接続できることを確認
     - `groupManager.GetShardingConnection(tableNumber)`が正しく動作することを確認
  4. **バッチ挿入**: バッチ挿入がPostgreSQLで正常に動作することを確認
     - `insertDmUsersBatch`関数がPostgreSQLで正常に動作することを確認

#### 3.3.2 dm_postsテーブルへのデータ生成
- **対象関数**: `generateDmPosts`
- **確認内容**:
  1. **UUIDv7の生成**: UUIDv7の生成が正常に動作することを確認
  2. **テーブル番号の計算**: テーブル番号の計算が正しく動作することを確認（user_idベース）
     - `tableSelector.GetTableNumberFromUUID(dmUserID)`が正しく動作することを確認
  3. **シャーディング接続**: 正しいシャーディングデータベースに接続できることを確認
     - `groupManager.GetShardingConnection(tableNumber)`が正しく動作することを確認
  4. **バッチ挿入**: バッチ挿入がPostgreSQLで正常に動作することを確認
     - `insertDmPostsBatch`関数がPostgreSQLで正常に動作することを確認

#### 3.3.3 dm_newsテーブルへのデータ生成
- **対象関数**: `generateDmNews`
- **確認内容**:
  1. **master接続**: masterデータベースに接続できることを確認
     - `groupManager.GetMasterConnection()`が正しく動作することを確認
  2. **バッチ挿入**: バッチ挿入がPostgreSQLで正常に動作することを確認
     - `insertDmNewsBatch`関数がPostgreSQLで正常に動作することを確認
     - GORMの`CreateInBatches`がPostgreSQLで正常に動作することを確認

### 3.4 テストコードの確認・修正

#### 3.4.1 テストコードの存在確認
- **対象ファイル**: `server/cmd/generate-sample-data/`配下のテストファイル
- **確認内容**:
  - テストファイルが存在するか確認
  - テストコードがSQLiteに依存しているか確認

#### 3.4.2 テストコードのPostgreSQL対応
- **修正内容**（テストコードが存在する場合）:
  - SQLite固有の設定をPostgreSQL設定に変更
  - テストデータベースの初期化方法をPostgreSQL対応に変更
  - テスト実行時のデータベース接続方法をPostgreSQL対応に変更

### 3.5 ドキュメントの更新

#### 3.5.1 README.mdの更新
- **ファイル**: `README.md`
- **更新内容**:
  1. **コマンドの説明**: `generate-sample-data`コマンドのPostgreSQL利用に関する記述を追加
  2. **前提条件**: コマンド実行前のPostgreSQL起動手順を記載
     - PostgreSQLコンテナの起動方法
     - マイグレーションの適用方法
  3. **コマンド実行方法**: コマンドの実行方法を記載
     - コマンドの実行コマンド
     - 環境変数の設定方法
     - 実行結果の確認方法

#### 3.5.2 その他のドキュメントの更新
- **対象ファイル**: `docs/`配下の関連ドキュメント
- **更新内容**:
  - SQLiteに関する記述をPostgreSQLに変更
  - PostgreSQL利用に関する記述を追加

## 4. データモデル設計

### 4.1 データ生成対象テーブル

#### 4.1.1 dm_usersテーブル
- **データベース**: shardingグループ（論理シャーディング数8、物理DB 4台）
- **テーブル名**: `dm_users_000` ～ `dm_users_031`（32テーブル）
- **データ型**:
  - `id`: `varchar(32)`（UUIDv7、ハイフン抜き小文字）
  - `name`: `varchar(255)`
  - `email`: `varchar(255)`
  - `created_at`: `timestamp`
  - `updated_at`: `timestamp`
- **生成数**: 100件（デフォルト）
- **シャーディングキー**: `id`（UUIDv7）

#### 4.1.2 dm_postsテーブル
- **データベース**: shardingグループ（論理シャーディング数8、物理DB 4台）
- **テーブル名**: `dm_posts_000` ～ `dm_posts_031`（32テーブル）
- **データ型**:
  - `id`: `varchar(32)`（UUIDv7、ハイフン抜き小文字）
  - `user_id`: `varchar(32)`（UUIDv7、ハイフン抜き小文字）
  - `title`: `varchar(255)`
  - `content`: `text`
  - `created_at`: `timestamp`
  - `updated_at`: `timestamp`
- **生成数**: 100件（デフォルト）
- **シャーディングキー**: `user_id`（UUIDv7）

#### 4.1.3 dm_newsテーブル
- **データベース**: masterグループ
- **テーブル名**: `dm_news`（固定テーブル名）
- **データ型**:
  - `id`: `bigint`（自動インクリメント）
  - `title`: `varchar(255)`
  - `content`: `text`
  - `author_id`: `bigint`（NULL許可）
  - `published_at`: `timestamp`（NULL許可）
  - `created_at`: `timestamp`
  - `updated_at`: `timestamp`
- **生成数**: 100件（デフォルト）
- **シャーディング**: なし（masterデータベースのみ）

### 4.2 データ生成フロー

```
1. 設定ファイルの読み込み
   ↓
2. GroupManagerの初期化
   ↓
3. データベース接続確認（PingAll）
   ↓
4. dm_usersテーブルへのデータ生成
   ├─ UUIDv7の生成
   ├─ テーブル番号の計算（idベース）
   ├─ シャーディング接続の取得
   └─ バッチ挿入
   ↓
5. dm_postsテーブルへのデータ生成
   ├─ UUIDv7の生成
   ├─ dm_user_idの選択
   ├─ テーブル番号の計算（user_idベース）
   ├─ シャーディング接続の取得
   └─ バッチ挿入
   ↓
6. dm_newsテーブルへのデータ生成
   ├─ master接続の取得
   └─ バッチ挿入
   ↓
7. 生成完了
```

## 5. エラーハンドリング設計

### 5.1 エラー処理の方針
- **設定ファイル読み込みエラー**: `log.Fatalf`でエラーメッセージを出力して終了
- **データベース接続エラー**: `log.Fatalf`でエラーメッセージを出力して終了
- **データ生成エラー**: `log.Fatalf`でエラーメッセージを出力して終了
- **バッチ挿入エラー**: `fmt.Errorf`でエラーを返し、上位で処理

### 5.2 エラーメッセージ
- **設定ファイル読み込みエラー**: `"Failed to load config: %v"`
- **GroupManager初期化エラー**: `"Failed to create group manager: %v"`
- **データベース接続確認エラー**: `"Failed to ping databases: %v"`
- **dm_users生成エラー**: `"Failed to generate users: %v"`
- **dm_posts生成エラー**: `"Failed to generate posts: %v"`
- **dm_news生成エラー**: `"Failed to generate news: %v"`
- **バッチ挿入エラー**: `"failed to insert batch: %w"`（各関数内）

## 6. パフォーマンス設計

### 6.1 バッチサイズ
- **デフォルトバッチサイズ**: 500件
- **理由**: 既存のコードで使用されているバッチサイズを維持

### 6.2 データ生成時間
- **目標**: 100件のデータ生成を数秒以内に完了
- **最適化**: バッチ挿入を使用してパフォーマンスを最適化

### 6.3 接続プール
- **既存の設定**: `server/internal/db/connection.go`で接続プールが設定されている
- **変更不要**: 既存の接続プール設定を維持

## 7. セキュリティ設計

### 7.1 SQLインジェクション対策
- **動的テーブル名**: テーブル名は数値から生成されるため、SQLインジェクションのリスクは低い
- **プレースホルダー**: GORMのプレースホルダー構文を使用してSQLインジェクションを防止

### 7.2 パスワード管理
- **開発環境**: 設定ファイルに固定パスワード（`webdb`）を記載
- **本番環境**: 適切なパスワード管理を推奨

## 8. テスト戦略

### 8.1 単体テスト
- **対象**: 各データ生成関数（`generateDmUsers`, `generateDmPosts`, `generateDmNews`）
- **テスト内容**:
  - PostgreSQL接続の確認
  - データ生成の確認
  - エラーハンドリングの確認

### 8.2 統合テスト
- **対象**: コマンド全体の動作
- **テスト内容**:
  - コマンドの実行
  - データベースへのデータ投入確認
  - エラーハンドリングの確認

### 8.3 動作確認
- **対象**: 実際のPostgreSQL環境での動作確認
- **確認内容**:
  - コマンドの実行
  - データベースへのデータ投入確認
  - 各テーブルへのデータ投入確認

## 9. 実装上の注意事項

### 9.1 設定ファイルの管理
- **設定ファイル**: `config/{env}/database.yaml`から接続情報を読み込む
- **環境変数**: `APP_ENV`環境変数で環境を指定（develop/staging/production）
- **設定構造**: 既存の`config/{env}/database.yaml`の構造を維持
- **PostgreSQL設定**: Issue #87で定義されたPostgreSQL構成に合わせる

### 9.2 SQL構文の確認
- **プレースホルダー**: GORMのプレースホルダー構文（`?`）がPostgreSQLで正常に動作することを確認
- **動的テーブル名**: 動的テーブル名の使用がPostgreSQLで正常に動作することを確認
- **バッチ挿入**: バッチ挿入のSQL構文がPostgreSQL対応であることを確認

### 9.3 データ生成処理の確認
- **UUIDv7**: UUIDv7の生成とデータ型（`varchar(32)`）が正しく動作することを確認
- **テーブル番号計算**: テーブル番号の計算が正しく動作することを確認
- **バッチサイズ**: 既存のバッチサイズ（500件）を維持

### 9.4 動作確認
- **接続確認**: コマンド実行時にPostgreSQLへの接続を確認
- **データ生成確認**: 各テーブルにデータが正常に生成されることを確認
- **エラーハンドリング**: 接続エラー時のエラーハンドリングを確認

### 9.5 ドキュメント整備
- **起動手順**: PostgreSQLコンテナの起動・マイグレーション適用・コマンド実行の手順を記載
- **コマンド実行**: `generate-sample-data`コマンドの実行方法を記載
- **トラブルシューティング**: よくある問題と解決方法を記載

## 10. 参考情報

### 10.1 関連Issue
- GitHub Issue #85: 開発環境はPostgreSQLを利用する前提とする
- GitHub Issue #86: PostgreSQLの起動スクリプトと、Atlasマイグレーションスクリプトの修正
- GitHub Issue #87: APIサーバーの修正
- GitHub Issue #89: server/cmd/generate-sample-data コマンド

### 10.2 既存ドキュメント
- `README.md`: プロジェクト概要とセットアップ手順
- `docs/Architecture.md`: システムアーキテクチャ
- `docs/Initial-Setup.md`: 初期セットアップ手順
- `config/{env}/database.yaml`: 環境別データベース設定

### 10.3 技術スタック
- **PostgreSQL**: 15-alpine（Dockerイメージ）
- **GORM**: 既存のGORMライブラリ
- **PostgreSQLドライバー**: `gorm.io/driver/postgres`（既存）

### 10.4 参考リンク
- PostgreSQL公式ドキュメント: https://www.postgresql.org/docs/
- GORM公式ドキュメント: https://gorm.io/docs/
- GORM PostgreSQLドライバー: https://github.com/go-gorm/postgres
