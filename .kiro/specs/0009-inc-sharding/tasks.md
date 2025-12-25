# シャーディング数増加実装タスク一覧

## 概要
シャーディング数を2から4に増やす実装を段階的に進めるためのタスク一覧。要件定義書と設計書に基づき、実装可能な粒度でタスクを分解。

## 実装フェーズ

### Phase 1: マイグレーションファイルの作成

#### - [ ] タスク 1.1: shard3ディレクトリの作成
**目的**: shard3のマイグレーションファイル用ディレクトリを作成

**作業内容**:
- `db/migrations/shard3/`ディレクトリを作成

**受け入れ基準**:
- `db/migrations/shard3/`ディレクトリが存在する

---

#### - [ ] タスク 1.2: shard3のマイグレーションファイル作成
**目的**: shard3の初期化スクリプトを作成

**作業内容**:
- `db/migrations/shard3/001_init.sql`を作成
- shard1の`001_init.sql`をベースに、コメントを「Shard 3」に変更
- 以下の内容を含む:
  - `users`テーブル
  - `posts`テーブル
  - インデックス（`idx_users_email`, `idx_posts_user_id`, `idx_posts_created_at`）

**受け入れ基準**:
- `db/migrations/shard3/001_init.sql`が存在する
- shard1の`001_init.sql`と同じスキーマを使用している
- コメントが「Shard 3」に更新されている

---

#### - [ ] タスク 1.3: shard4ディレクトリの作成
**目的**: shard4のマイグレーションファイル用ディレクトリを作成

**作業内容**:
- `db/migrations/shard4/`ディレクトリを作成

**受け入れ基準**:
- `db/migrations/shard4/`ディレクトリが存在する

---

#### - [ ] タスク 1.4: shard4のマイグレーションファイル作成
**目的**: shard4の初期化スクリプトを作成

**作業内容**:
- `db/migrations/shard4/001_init.sql`を作成
- shard1の`001_init.sql`をベースに、コメントを「Shard 4」に変更
- 以下の内容を含む:
  - `users`テーブル
  - `posts`テーブル
  - インデックス（`idx_users_email`, `idx_posts_user_id`, `idx_posts_created_at`）

**受け入れ基準**:
- `db/migrations/shard4/001_init.sql`が存在する
- shard1の`001_init.sql`と同じスキーマを使用している
- コメントが「Shard 4」に更新されている

---

### Phase 2: 設定ファイルの更新

#### - [ ] タスク 2.1: 開発環境設定ファイルの更新
**目的**: develop環境の設定ファイルにshard3とshard4を追加

**作業内容**:
- `config/develop/database.yaml`を開く
- 既存のshard1とshard2の設定は変更しない
- shard3とshard4を追加:
  ```yaml
  - id: 3
    driver: sqlite3
    dsn: ./data/shard3.db
    writer_dsn: ./data/shard3.db
    reader_dsns:
      - ./data/shard3.db
    reader_policy: random
    max_connections: 10
    max_idle_connections: 5
    connection_max_lifetime: 300s
  - id: 4
    driver: sqlite3
    dsn: ./data/shard4.db
    writer_dsn: ./data/shard4.db
    reader_dsns:
      - ./data/shard4.db
    reader_policy: random
    max_connections: 10
    max_idle_connections: 5
    connection_max_lifetime: 300s
  ```

**受け入れ基準**:
- `config/develop/database.yaml`にshard3とshard4が追加されている
- 既存のshard1とshard2の設定が変更されていない
- YAML形式が正しい
- 既存のパターンに従っている

---

#### - [ ] タスク 2.2: ステージング環境設定ファイルの更新
**目的**: staging環境の設定ファイルにshard3とshard4を追加

**作業内容**:
- `config/staging/database.yaml`を開く
- 既存のshard1とshard2の設定は変更しない
- shard3とshard4を追加:
  ```yaml
  - id: 3
    driver: postgres
    host: staging-db-shard3.example.com
    port: 5432
    name: app_db_shard3
    user: staging_user
    password: ${DB_PASSWORD_SHARD3}
    writer_dsn: host=staging-db-shard3-writer.example.com port=5432 user=staging_user password=${DB_PASSWORD_SHARD3} dbname=app_db_shard3 sslmode=require
    reader_dsns:
      - host=staging-db-shard3-reader1.example.com port=5432 user=staging_user password=${DB_PASSWORD_SHARD3} dbname=app_db_shard3 sslmode=require
    reader_policy: random
    max_connections: 25
    max_idle_connections: 10
    connection_max_lifetime: 600s
  - id: 4
    driver: postgres
    host: staging-db-shard4.example.com
    port: 5432
    name: app_db_shard4
    user: staging_user
    password: ${DB_PASSWORD_SHARD4}
    writer_dsn: host=staging-db-shard4-writer.example.com port=5432 user=staging_user password=${DB_PASSWORD_SHARD4} dbname=app_db_shard4 sslmode=require
    reader_dsns:
      - host=staging-db-shard4-reader1.example.com port=5432 user=staging_user password=${DB_PASSWORD_SHARD4} dbname=app_db_shard4 sslmode=require
    reader_policy: random
    max_connections: 25
    max_idle_connections: 10
    connection_max_lifetime: 600s
  ```

**受け入れ基準**:
- `config/staging/database.yaml`にshard3とshard4が追加されている
- 既存のshard1とshard2の設定が変更されていない
- 環境変数（`DB_PASSWORD_SHARD3`、`DB_PASSWORD_SHARD4`）が使用されている
- YAML形式が正しい
- 既存のパターンに従っている

---

#### - [ ] タスク 2.3: 本番環境設定ファイル例の更新
**目的**: production環境の設定ファイル例にshard3とshard4を追加

**作業内容**:
- `config/production/database.yaml.example`を開く
- 既存のshard1とshard2の設定は変更しない
- shard3とshard4を追加（ステージング環境と同様の構造で、ホスト名を`prod-db-shard*`に変更）:
  ```yaml
  - id: 3
    driver: postgres
    host: prod-db-shard3.example.com
    port: 5432
    name: app_db_shard3
    user: prod_user
    password: ${DB_PASSWORD_SHARD3}
    writer_dsn: host=prod-db-shard3-writer.example.com port=5432 user=prod_user password=${DB_PASSWORD_SHARD3} dbname=app_db_shard3 sslmode=require
    reader_dsns:
      - host=prod-db-shard3-reader1.example.com port=5432 user=prod_user password=${DB_PASSWORD_SHARD3} dbname=app_db_shard3 sslmode=require
      - host=prod-db-shard3-reader2.example.com port=5432 user=prod_user password=${DB_PASSWORD_SHARD3} dbname=app_db_shard3 sslmode=require
    reader_policy: round_robin
    max_connections: 50
    max_idle_connections: 20
    connection_max_lifetime: 600s
  - id: 4
    driver: postgres
    host: prod-db-shard4.example.com
    port: 5432
    name: app_db_shard4
    user: prod_user
    password: ${DB_PASSWORD_SHARD4}
    writer_dsn: host=prod-db-shard4-writer.example.com port=5432 user=prod_user password=${DB_PASSWORD_SHARD4} dbname=app_db_shard4 sslmode=require
    reader_dsns:
      - host=prod-db-shard4-reader1.example.com port=5432 user=prod_user password=${DB_PASSWORD_SHARD4} dbname=app_db_shard4 sslmode=require
      - host=prod-db-shard4-reader2.example.com port=5432 user=prod_user password=${DB_PASSWORD_SHARD4} dbname=app_db_shard4 sslmode=require
    reader_policy: round_robin
    max_connections: 50
    max_idle_connections: 20
    connection_max_lifetime: 600s
  ```

**受け入れ基準**:
- `config/production/database.yaml.example`にshard3とshard4が追加されている
- 既存のshard1とshard2の設定が変更されていない
- 環境変数（`DB_PASSWORD_SHARD3`、`DB_PASSWORD_SHARD4`）が使用されている
- 複数のReaderが設定されている（本番環境のパターンに従う）
- YAML形式が正しい
- 既存のパターンに従っている

---

### Phase 3: ドキュメントの更新

#### - [ ] タスク 3.1: Sharding.mdの更新（シャーディング数の説明）
**目的**: シャーディング数の説明を2シャードから4シャードに更新

**作業内容**:
- `docs/Sharding.md`を開く
- 「2シャード」という記述を「4シャード」に更新
- シャーディング数の説明を更新

**受け入れ基準**:
- `docs/Sharding.md`のシャーディング数の説明が4シャードに更新されている

---

#### - [ ] タスク 3.2: Sharding.mdの更新（設定例）
**目的**: 設定例を4シャード構成に更新

**作業内容**:
- `docs/Sharding.md`の設定例セクションを更新
- 開発環境の設定例にshard3とshard4を追加
- 本番環境の設定例にshard3とshard4を追加

**受け入れ基準**:
- `docs/Sharding.md`の設定例が4シャード構成を反映している
- 開発環境と本番環境の両方の設定例が更新されている

---

#### - [ ] タスク 3.3: Sharding.mdの更新（データ分散の例）
**目的**: データ分散の例を4シャード構成に更新

**作業内容**:
- `docs/Sharding.md`のデータ分散セクションを更新
- 4シャードでのデータ分散例を追加
- 例: `User ID 1 → hash(1) % 4 = Shard 1, 2, 3, または4`

**受け入れ基準**:
- `docs/Sharding.md`のデータ分散の例が4シャード構成を反映している

---

### Phase 4: 動作確認

#### - [ ] タスク 4.1: 設定ファイルの読み込み確認
**目的**: 4シャードの設定が正しく読み込まれることを確認

**作業内容**:
- アプリケーションを起動
- 設定ファイルが正しく読み込まれることを確認
- 4つのシャードが検出されることを確認

**受け入れ基準**:
- アプリケーションが正常に起動する
- 設定ファイルから4つのシャードが読み込まれる
- エラーが発生しない

---

#### - [ ] タスク 4.2: データベース接続確認
**目的**: 4つのシャードに正常に接続できることを確認

**作業内容**:
- アプリケーションを起動
- 4つのシャードすべてに接続できることを確認
- 接続エラーが発生しないことを確認

**受け入れ基準**:
- 4つのシャードすべてに正常に接続できる
- 接続エラーが発生しない

---

#### - [ ] タスク 4.3: シャーディングロジックの動作確認
**目的**: 4シャードでのシャーディングロジックが正常に動作することを確認

**作業内容**:
- 新しいデータを作成
- データが4つのシャードに適切に分散されることを確認
- シャードIDが1, 2, 3, 4のいずれかになることを確認

**受け入れ基準**:
- 新しいデータが4つのシャードに適切に分散される
- シャードIDが1, 2, 3, 4のいずれかになる
- 既存のシャーディングロジックが正常に動作する

---

#### - [ ] タスク 4.4: クロスシャードクエリの動作確認
**目的**: 4シャードでのクロスシャードクエリが正常に動作することを確認

**作業内容**:
- 全ユーザー取得などのクロスシャードクエリを実行
- 4つのシャードすべてからデータを取得できることを確認
- 結果が正しくマージされることを確認

**受け入れ基準**:
- クロスシャードクエリが正常に動作する
- 4つのシャードすべてからデータを取得できる
- 結果が正しくマージされる

---

#### - [ ] タスク 4.5: 既存テストの動作確認
**目的**: 既存のテストが4シャードでも正常に動作することを確認

**作業内容**:
- 既存のテストを実行
- 4シャードでも正常に動作することを確認
- テストが失敗しないことを確認

**受け入れ基準**:
- 既存のテストが4シャードでも正常に動作する
- テストが失敗しない

---

## 実装上の注意事項

### 既存設定の保持
- 既存のshard1とshard2の設定は変更しない
- 既存のマイグレーションファイルは変更しない

### 一貫性の維持
- shard3とshard4の設定は既存のパターンに従う
- shard3とshard4のマイグレーションファイルはshard1とshard2と同じスキーマを使用

### データ損失の許容
- 既存データの移行は行わない
- データ損失を許容する
- 必要に応じてデータベースをリセットしても良い

### コード変更の不要性
- 既存のシャーディングロジックは変更不要（動的にシャード数を検出するため）
- 既存のデータベース接続処理は変更不要
- 既存のテストコードは変更不要

## 参考情報

### 既存ファイル
- `db/migrations/shard1/001_init.sql`: shard1の初期化スクリプト（ベースとして使用）
- `config/develop/database.yaml`: 開発環境設定ファイル（既存のパターンを参照）
- `config/staging/database.yaml`: ステージング環境設定ファイル（既存のパターンを参照）
- `config/production/database.yaml.example`: 本番環境設定ファイル例（既存のパターンを参照）
- `docs/Sharding.md`: シャーディング戦略の詳細（更新対象）

### 既存実装
- `server/internal/db/sharding.go`: `HashBasedSharding`の実装（変更不要）
- `server/internal/config/config.go`: 設定構造体の定義（変更不要）
- `server/internal/db/connection.go`: データベース接続処理（変更不要）

