# PostgreSQL起動スクリプト・マイグレーションスクリプト修正設計書

## 1. 概要

### 1.1 設計の目的
要件定義書に基づき、PostgreSQLの起動スクリプトとAtlasマイグレーションスクリプトを修正するための詳細設計を定義する。開発環境、staging環境、production環境でPostgreSQLを利用する前提に移行するため、Docker ComposeによるPostgreSQL環境構築と、設定ファイルベースのマイグレーションスクリプトを実現する。

### 1.2 設計の範囲
- `docker-compose.postgres.yml`の修正（master 1台 + sharding 4台の構成）
- `scripts/start-postgres.sh`の作成（PostgreSQL起動スクリプト）
- `scripts/migrate.sh`の新規作成（既存スクリプトは破棄、PostgreSQL対応、設定ファイルからの読み込み、初期データとViewの適用）
- HCLファイルのUUIDv7カラム定義の修正（`bigint`から`varchar(32)`に変更）
- データ投入確認手順の実装

### 1.3 設計方針
- **既存システムとの統合**: 既存の`postgres-network`を使用し、既存のマイグレーションファイルをそのまま利用
- **設定ファイルベース**: 環境変数ではなく設定ファイル（`config/{env}/database.yaml`）から接続情報を読み込む
- **環境別対応**: `APP_ENV`環境変数で環境を指定（develop/staging/production）
- **データ永続化**: ボリュームマウントにより各PostgreSQLコンテナのデータを永続化
- **既存設定ファイルの維持**: `config/{env}/database.yaml`の構造を維持

## 2. アーキテクチャ設計

### 2.1 既存アーキテクチャの分析

#### 2.1.1 現在の構成
- **PostgreSQL**: `docker-compose.postgres.yml`で1台のPostgreSQLコンテナを定義
- **マイグレーション**: `scripts/migrate.sh`がSQLite用のAtlasマイグレーションを実行（本実装で破棄）
- **設定ファイル**: `config/{env}/database.yaml`にPostgreSQL設定がコメントアウトされている
- **初期データ**: `db/migrations/master/`内の初期データSQLファイル（現在: `20251230045548_seed_data.sql`）にGoAdmin関連の初期データが含まれている（SQLite構文）
- **View**: `db/migrations/view_master/`内のView定義SQLファイル（現在: `20260103030225_create_dm_news_view.sql`）にView定義が含まれている（生SQL）
- **適用順序**: Atlasはファイル名順にマイグレーションを適用するため、適用タイミングを調整するためにファイル名を変更しても良い

#### 2.1.2 既存パターンの維持
- **ネットワーク**: 既存の`postgres-network`を使用
- **マイグレーションファイル**: `db/migrations/master/`, `db/migrations/sharding_*/`は変更しない
- **設定ファイル構造**: `config/{env}/database.yaml`の構造を維持

### 2.2 システム構成図

```
┌─────────────────────────────────────────────────────────────┐
│                    開発者/運用者                              │
└──────────────────┬────────────────────────────────────┘
                    │
                    │ ./scripts/start-postgres.sh start
                    │
                    ▼
┌─────────────────────────────────────────────────────────────┐
│              Docker Compose (docker-compose.postgres.yml)    │
│                                                              │
│  ┌──────────────────┐  ┌──────────────────┐              │
│  │ postgres-master   │  │ postgres-sharding │              │
│  │ (ポート: 5432)    │  │ -1 (ポート: 5433) │              │
│  │ DB: webdb_master  │  │ DB: webdb_sharding_1│             │
│  └────────┬─────────┘  └────────┬─────────┘              │
│           │                      │                          │
│           │                      ├── postgres-sharding-2    │
│           │                      │   (ポート: 5434)        │
│           │                      │   DB: webdb_sharding_2   │
│           │                      │                          │
│           │                      ├── postgres-sharding-3    │
│           │                      │   (ポート: 5435)        │
│           │                      │   DB: webdb_sharding_3  │
│           │                      │                          │
│           │                      └── postgres-sharding-4    │
│           │                          (ポート: 5436)        │
│           │                          DB: webdb_sharding_4   │
│           │                                                  │
│           └──────────────────────────────────────────────┘
│                          │
│                          │ postgres-network
│                          │
└──────────────────────────┼──────────────────────────────────┘
                           │
                           │ ./scripts/migrate.sh
                           │ (config/{env}/database.yamlから読み込み)
                           │
                           ▼
┌─────────────────────────────────────────────────────────────┐
│                    Atlas CLI                                  │
│                                                              │
│  ┌──────────────────┐  ┌──────────────────┐              │
│  │ db/migrations/   │  │ db/migrations/   │              │
│  │ master/          │  │ sharding_1/      │              │
│  │                  │  │ sharding_2/      │              │
│  │                  │  │ sharding_3/      │              │
│  │                  │  │ sharding_4/      │              │
│  └──────────────────┘  └──────────────────┘              │
└─────────────────────────────────────────────────────────────┘
```

### 2.3 ディレクトリ構造

#### 2.3.1 変更前の構造
```
.
├── docker-compose.postgres.yml    # 1台のPostgreSQLコンテナ
├── scripts/
│   └── migrate.sh                 # SQLite用マイグレーション
├── config/
│   └── {env}/
│       └── database.yaml          # PostgreSQL設定がコメントアウト
└── postgres/
    └── data/                      # 1台分のデータ
```

#### 2.3.2 変更後の構造
```
.
├── docker-compose.postgres.yml    # master 1台 + sharding 4台の構成
├── scripts/
│   ├── start-postgres.sh          # PostgreSQL起動スクリプト（新規）
│   └── migrate.sh                 # PostgreSQL対応の新規スクリプト（既存は破棄）
├── config/
│   └── {env}/
│       └── database.yaml          # PostgreSQL設定のコメントアウト解除（別Issue）
├── db/
│   └── migrations/
│       ├── master/                 # Atlasマイグレーション + 初期データ
│       │   ├── 20251230045547_initial_schema.sql
│       │   └── {timestamp}_seed_data.sql  # 初期データ（ファイル名は適用順序に応じて変更可能）
│       ├── sharding_1/             # Atlasマイグレーション
│       ├── sharding_2/             # Atlasマイグレーション
│       ├── sharding_3/             # Atlasマイグレーション
│       ├── sharding_4/             # Atlasマイグレーション
│       └── view_master/            # Viewマイグレーション（生SQL）
│           └── {timestamp}_create_dm_news_view.sql  # ファイル名は適用順序に応じて変更可能
└── postgres/
    └── data/
        ├── master/                 # master PostgreSQLのデータ（新規）
        ├── sharding_1/             # sharding_1 PostgreSQLのデータ（新規）
        ├── sharding_2/             # sharding_2 PostgreSQLのデータ（新規）
        ├── sharding_3/             # sharding_3 PostgreSQLのデータ（新規）
        └── sharding_4/             # sharding_4 PostgreSQLのデータ（新規）
```

## 3. コンポーネント設計

### 3.1 Docker Compose設定

#### 3.1.1 docker-compose.postgres.yml

| フィールド | 詳細 |
|-----------|------|
| Intent | master 1台 + sharding 4台のPostgreSQLコンテナを定義 |
| Requirements | 3.1.1 |

**構成**:
- **masterグループ**: `postgres-master`サービス（1台）
- **shardingグループ**: `postgres-sharding-1` ～ `postgres-sharding-4`サービス（4台）

**サービス定義**:

```yaml
services:
  postgres-master:
    image: postgres:15-alpine
    container_name: postgres-master
    ports:
      - "5432:5432"
    environment:
      POSTGRES_USER: webdb
      POSTGRES_PASSWORD: webdb
      POSTGRES_DB: webdb_master
    volumes:
      - ./postgres/data/master:/var/lib/postgresql/data
    restart: unless-stopped
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U webdb"]
      interval: 10s
      timeout: 5s
      retries: 5
    networks:
      - postgres-network

  postgres-sharding-1:
    image: postgres:15-alpine
    container_name: postgres-sharding-1
    ports:
      - "5433:5432"
    environment:
      POSTGRES_USER: webdb
      POSTGRES_PASSWORD: webdb
      POSTGRES_DB: webdb_sharding_1
    volumes:
      - ./postgres/data/sharding_1:/var/lib/postgresql/data
    restart: unless-stopped
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U webdb"]
      interval: 10s
      timeout: 5s
      retries: 5
    networks:
      - postgres-network

  postgres-sharding-2:
    # postgres-sharding-1と同様（ポート: 5434, DB: webdb_sharding_2）

  postgres-sharding-3:
    # postgres-sharding-1と同様（ポート: 5435, DB: webdb_sharding_3）

  postgres-sharding-4:
    # postgres-sharding-1と同様（ポート: 5436, DB: webdb_sharding_4）

networks:
  postgres-network:
    name: postgres-network
    driver: bridge
```

**実装上の注意事項**:
- 各コンテナのデータディレクトリを個別にマウント（`./postgres/data/{service-name}/`）
- ホスト側ポートを競合しないように設定（5432, 5433, 5434, 5435, 5436）
- 既存の`postgres-network`を使用
- ヘルスチェックを設定してコンテナの正常性を確認

### 3.2 起動スクリプト

#### 3.2.1 scripts/start-postgres.sh

| フィールド | 詳細 |
|-----------|------|
| Intent | PostgreSQLコンテナの起動・停止・状態確認を行うスクリプト |
| Requirements | 3.1.2 |

**機能**:
- `start`: PostgreSQLコンテナを起動（`docker-compose -f docker-compose.postgres.yml up -d`）
- `stop`: PostgreSQLコンテナを停止（`docker-compose -f docker-compose.postgres.yml down`）
- `status`: PostgreSQLコンテナの状態を確認（`docker-compose -f docker-compose.postgres.yml ps`）
- `health`: ヘルスチェックを確認（`docker-compose -f docker-compose.postgres.yml ps`でhealthcheck状態を表示）

**実装上の注意事項**:
- Bashスクリプトとして実装
- `set -e`でエラー時に即座に終了
- 各コマンドのエラーハンドリングを実装
- 使用方法のヘルプメッセージを表示

### 3.4 HCLファイルのUUIDv7カラム定義修正

#### 3.4.1 db/schema/sharding_*/dm_posts.hcl

| フィールド | 詳細 |
|-----------|------|
| Intent | UUIDv7カラム（`id`, `user_id`）の型定義を`bigint`から`varchar(32)`に修正 |
| Requirements | 3.2.1（UUIDv7対応） |

**修正内容**:
- `dm_posts_XXX.id`: `bigint` → `varchar(32)`
- `dm_posts_XXX.user_id`: `bigint` → `varchar(32)`
- `unsigned`属性を削除（PostgreSQLには存在しない）

**対象ファイル**:
- `db/schema/sharding_1/dm_posts.hcl`
- `db/schema/sharding_2/dm_posts.hcl`
- `db/schema/sharding_3/dm_posts.hcl`
- `db/schema/sharding_4/dm_posts.hcl`

**修正例**:
```hcl
column "id" {
  null = false
  type = varchar(32)  # bigint → varchar(32)
  # unsigned = true を削除
  # auto_increment = false を削除
}
column "user_id" {
  null = false
  type = varchar(32)  # bigint → varchar(32)
  # unsigned = true を削除
}
```

#### 3.4.2 db/schema/sharding_*/dm_users.hcl

| フィールド | 詳細 |
|-----------|------|
| Intent | UUIDv7カラム（`id`）の型定義を`bigint`から`varchar(32)`に修正 |
| Requirements | 3.2.1（UUIDv7対応） |

**修正内容**:
- `dm_users_XXX.id`: `bigint` → `varchar(32)`
- `unsigned`属性を削除（PostgreSQLには存在しない）

**対象ファイル**:
- `db/schema/sharding_1/dm_users.hcl`
- `db/schema/sharding_2/dm_users.hcl`
- `db/schema/sharding_3/dm_users.hcl`
- `db/schema/sharding_4/dm_users.hcl`

**修正例**:
```hcl
column "id" {
  null = false
  type = varchar(32)  # bigint → varchar(32)
  # unsigned = true を削除
  # auto_increment = false を削除
}
```

**実装上の注意事項**:
- UUIDv7はハイフン抜き小文字で32文字の文字列として保存される
- PostgreSQLでは`varchar(32)`を使用
- `unsigned`属性はPostgreSQLには存在しないため削除
- `auto_increment`属性も削除（UUIDv7はアプリケーション側で生成）

### 3.3 マイグレーションスクリプト

#### 3.3.1 scripts/migrate.sh

| フィールド | 詳細 |
|-----------|------|
| Intent | 設定ファイルからPostgreSQL接続情報を読み込み、Atlasマイグレーションと生SQLを適用するスクリプト |
| Requirements | 3.2.1, 3.2.2, 3.2.3 |

**新規作成**:
- 既存の`scripts/migrate.sh`は破棄し、PostgreSQL対応の新規スクリプトを作成する

**機能**:
1. 設定ファイル（`config/{env}/database.yaml`）から接続情報を読み込む
2. `APP_ENV`環境変数で環境を指定（develop/staging/production）
3. Atlasマイグレーションを適用（`db/migrations/master/`, `db/migrations/sharding_*/`）
4. 初期データ（seed_data.sql）を適用
5. Viewのマイグレーションを生SQLで適用（`db/migrations/view_master/`）

**設定ファイルからの読み込み**:
- **masterグループ**: `database.groups.master[0]`から接続情報を取得
  - `host`, `port`, `user`, `password`, `name`（データベース名）
- **shardingグループ**: `database.groups.sharding.databases[]`から各shardingデータベースの接続情報を取得
  - `host`, `port`, `user`, `password`, `name`（データベース名）

**PostgreSQL URL形式**:
```
postgres://{user}:{password}@{host}:{port}/{dbname}?sslmode=disable
```

**マイグレーション適用順序**:
1. **Atlasマイグレーション**: `db/migrations/master/`, `db/migrations/sharding_*/`のAtlasマイグレーションを適用
   - Atlasはファイル名順にマイグレーションを適用するため、適用順序を調整する場合はファイル名を変更する
2. **初期データ**: `db/migrations/master/`内の初期データSQLファイルを適用（PostgreSQL用に修正が必要な場合あり）
   - 現在のファイル名: `20251230045548_seed_data.sql`
   - 適用タイミングを調整するためにファイル名を変更しても良い
3. **Viewマイグレーション**: `db/migrations/view_master/`の生SQLを適用（Atlas PROライセンスが必要なため、生SQLで実行）
   - 現在のファイル名: `20260103030225_create_dm_news_view.sql`
   - 適用タイミングを調整するためにファイル名を変更しても良い

**初期データの扱い**:
- `db/migrations/master/`内の初期データSQLファイル（現在: `20251230045548_seed_data.sql`）にGoAdmin関連の初期データが含まれている
- SQLite用の構文（`INSERT OR IGNORE`, `datetime('now')`）をPostgreSQL用に変換する必要がある
  - `INSERT OR IGNORE` → `INSERT ... ON CONFLICT DO NOTHING`
  - `datetime('now')` → `NOW()`
- Atlasはファイル名順にマイグレーションを適用するため、適用タイミングを調整するためにファイル名を変更しても良い

**Viewの扱い**:
- `db/migrations/view_master/`内のView定義SQLファイル（現在: `20260103030225_create_dm_news_view.sql`）にView定義が含まれている
- AtlasでViewを使うにはPROライセンスが必要なため、生SQLで実行する
- `psql`コマンドまたはPostgreSQLクライアントライブラリを使用して生SQLを実行
- Atlasはファイル名順にマイグレーションを適用するため、適用タイミングを調整するためにファイル名を変更しても良い

**実装上の注意事項**:
- YAMLファイルのパースには`yq`コマンドまたはGoスクリプトを使用（既存の設定読み込みロジックを参考）
- 設定ファイルが存在しない場合のエラーハンドリング
- 接続情報が不足している場合のエラーハンドリング
- Atlasマイグレーション適用時のエラーハンドリング
- 初期データ適用時のエラーハンドリング（SQLite構文の変換エラーなど）
- Viewマイグレーション適用時のエラーハンドリング

## 4. データモデル

### 4.1 設定ファイル構造

#### 4.1.1 config/{env}/database.yaml

既存の設定ファイル構造を維持し、PostgreSQL設定のコメントアウトを解除して使用する。

**masterグループ設定**:
```yaml
database:
  groups:
    master:
      - id: 1
        driver: postgres
        host: localhost
        port: 5432
        user: webdb
        password: webdb
        name: webdb_master
        max_connections: 25
        max_idle_connections: 5
        connection_max_lifetime: 1h
```

**shardingグループ設定**:
```yaml
    sharding:
      databases:
        - id: 1
          driver: postgres
          host: localhost
          port: 5433
          user: webdb
          password: webdb
          name: webdb_sharding_1
          max_connections: 25
          max_idle_connections: 5
          connection_max_lifetime: 1h
          table_range: [0, 7]
        - id: 2
          driver: postgres
          host: localhost
          port: 5434
          user: webdb
          password: webdb
          name: webdb_sharding_2
          max_connections: 25
          max_idle_connections: 5
          connection_max_lifetime: 1h
          table_range: [8, 15]
        - id: 3
          driver: postgres
          host: localhost
          port: 5435
          user: webdb
          password: webdb
          name: webdb_sharding_3
          max_connections: 25
          max_idle_connections: 5
          connection_max_lifetime: 1h
          table_range: [16, 23]
        - id: 4
          driver: postgres
          host: localhost
          port: 5436
          user: webdb
          password: webdb
          name: webdb_sharding_4
          max_connections: 25
          max_idle_connections: 5
          connection_max_lifetime: 1h
          table_range: [24, 31]

      tables:
        - name: dm_users
          suffix_count: 32
        - name: dm_posts
          suffix_count: 32
```

**実装上の注意事項**:
- 設定ファイルの構造は既存の`config/{env}/database.yaml`を維持
- PostgreSQL設定のコメントアウト解除は別Issue対応（本実装では設定ファイルからの読み込みロジックのみ実装）

## 5. エラーハンドリング

### 5.1 エラー戦略

#### 5.1.1 起動スクリプト（scripts/start-postgres.sh）
- **Docker Composeコマンドエラー**: エラーメッセージを表示して終了
- **コンテナ起動失敗**: ログを表示して原因を特定できるようにする
- **ポート競合**: エラーメッセージでポート番号を明示

#### 5.1.2 マイグレーションスクリプト（scripts/migrate.sh）
- **設定ファイル不存在**: `config/{env}/database.yaml`が存在しない場合、エラーメッセージを表示して終了
- **設定ファイルパースエラー**: YAMLのパースに失敗した場合、エラーメッセージを表示して終了
- **接続情報不足**: 必要な接続情報（host, port, user, password, name）が不足している場合、エラーメッセージを表示して終了
- **PostgreSQL接続エラー**: データベースに接続できない場合、エラーメッセージを表示して終了
- **マイグレーション適用エラー**: Atlasマイグレーションの適用に失敗した場合、エラーメッセージを表示して終了

### 5.2 エラーカテゴリと対応

| エラーカテゴリ | 原因 | 対応 |
|--------------|------|------|
| 設定ファイルエラー | 設定ファイルが存在しない、パースエラー | エラーメッセージを表示して終了 |
| 接続情報エラー | 接続情報が不足している | エラーメッセージを表示して終了 |
| データベース接続エラー | PostgreSQLに接続できない | エラーメッセージを表示して終了（PostgreSQLコンテナが起動しているか確認） |
| マイグレーションエラー | Atlasマイグレーションの適用に失敗 | エラーメッセージを表示して終了（マイグレーションファイルの確認） |

### 5.3 モニタリング

- **起動スクリプト**: 各コマンドの実行結果を標準出力に表示
- **マイグレーションスクリプト**: マイグレーション適用の進捗を標準出力に表示
- **エラーログ**: エラー発生時に詳細なエラーメッセージを標準エラー出力に表示

## 6. テスト戦略

### 6.1 ユニットテスト

#### 6.1.1 起動スクリプト（scripts/start-postgres.sh）
- スクリプトの構文チェック（`bash -n scripts/start-postgres.sh`）
- 各コマンド（start, stop, status, health）の動作確認

#### 6.1.2 マイグレーションスクリプト（scripts/migrate.sh）
- スクリプトの構文チェック（`bash -n scripts/migrate.sh`）
- 設定ファイルからの読み込みロジックのテスト（モック設定ファイルを使用）

### 6.2 統合テスト

#### 6.2.0 HCLファイル修正テスト
- HCLファイルのUUIDv7カラム定義が正しく修正されていることを確認
  - `dm_posts_XXX.id`と`dm_posts_XXX.user_id`が`varchar(32)`になっている
  - `dm_users_XXX.id`が`varchar(32)`になっている
  - `unsigned`属性が削除されている
  - `auto_increment`属性が削除されている（`id`カラム）

#### 6.2.1 PostgreSQL起動テスト
- `./scripts/start-postgres.sh start`でPostgreSQLコンテナが正常に起動することを確認
- 各PostgreSQLコンテナが正常に動作していることを確認（ヘルスチェック）
- `./scripts/start-postgres.sh stop`でPostgreSQLコンテナが正常に停止することを確認

#### 6.2.2 マイグレーションテスト
- `./scripts/migrate.sh master`でmasterグループのマイグレーションが正常に適用されることを確認
  - Atlasマイグレーション（`db/migrations/master/`）の適用（ファイル名順）
  - 初期データ（`db/migrations/master/`内の初期データSQLファイル）の適用（ファイル名順）
  - Viewマイグレーション（`db/migrations/view_master/`）の適用（ファイル名順）
- `./scripts/migrate.sh sharding`でshardingグループのマイグレーションが正常に適用されることを確認（4台すべて）
- `./scripts/migrate.sh all`で全マイグレーションが正常に適用されることを確認

### 6.3 データ投入確認テスト

#### 6.3.1 テーブル存在確認
- 各PostgreSQLデータベースにテーブルが作成されていることを確認（`psql`コマンドまたはAtlas CLI）

#### 6.3.2 データ投入確認
- テストデータを投入してデータベースにデータが入ることを確認
- 各shardingデータベースにデータが正しく分散されていることを確認

## 7. 実装上の注意事項

### 7.1 Docker Compose設定
- **サービス名**: `postgres-master`, `postgres-sharding-1` ～ `postgres-sharding-4`
- **ネットワーク**: 既存の`postgres-network`を使用
- **ボリューム**: 各PostgreSQLコンテナのデータを個別のボリュームにマウント
- **環境変数**: 各コンテナで適切な環境変数を設定
- **ポート**: ホスト側ポートを競合しないように設定（5432, 5433, 5434, 5435, 5436）

### 7.2 マイグレーションスクリプト
- **新規作成**: 既存の`scripts/migrate.sh`は破棄し、PostgreSQL対応の新規スクリプトを作成
- **URL形式**: PostgreSQL用のURL形式（`postgres://user:password@host:port/dbname?sslmode=disable`）を使用
- **設定ファイル**: `config/{env}/database.yaml`から接続情報を読み込む
- **環境変数**: `APP_ENV`環境変数で環境を指定（develop/staging/production）
- **Atlasマイグレーション**: `db/migrations/master/`, `db/migrations/sharding_*/`のAtlasマイグレーションを適用
- **初期データ**: `db/migrations/master/`内の初期データSQLファイルを生SQLで適用（SQLite構文をPostgreSQL構文に変換）
  - ファイル名は適用順序に応じて変更可能（Atlasはファイル名順に適用）
- **Viewマイグレーション**: `db/migrations/view_master/`の生SQLを適用（Atlas PROライセンスが必要なため、生SQLで実行）
  - ファイル名は適用順序に応じて変更可能（Atlasはファイル名順に適用）
- **エラーハンドリング**: マイグレーション適用時のエラーを適切に処理

### 7.3 設定ファイルの管理
- **設定ファイル**: `config/{env}/database.yaml`から接続情報を読み込む
- **環境変数**: `APP_ENV`環境変数で環境を指定（develop/staging/production）
- **設定構造**: 既存の`config/{env}/database.yaml`の構造を維持
- **PostgreSQL設定**: `config/{env}/database.yaml`のPostgreSQL設定のコメントアウト解除は別Issue対応

### 7.4 データ永続化
- **ボリュームマウント**: 各PostgreSQLコンテナのデータをボリュームマウントで永続化
- **データディレクトリ**: `postgres/data/master/`, `postgres/data/sharding_1/` ～ `postgres/data/sharding_4/`にデータを保存

## 8. 参考情報

### 8.1 関連Issue
- GitHub Issue #85: 開発環境はPostgreSQLを利用する前提とする
- GitHub Issue #86: PostgreSQLの起動スクリプトと、Atlasマイグレーションスクリプトの修正

### 8.2 既存ドキュメント
- `README.md`: プロジェクト概要とセットアップ手順
- `docs/Architecture.md`: システムアーキテクチャ
- `docs/Initial-Setup.md`: 初期セットアップ手順
- `config/{env}/database.yaml`: 環境別データベース設定

### 8.3 技術スタック
- **PostgreSQL**: 15-alpine（Dockerイメージ）
- **Atlas**: Atlas CLI（既存）
- **Docker**: Docker Compose（既存）
- **シェルスクリプト**: Bashスクリプトを使用

### 8.4 参考リンク
- PostgreSQL公式ドキュメント: https://www.postgresql.org/docs/
- Atlas公式ドキュメント: https://atlasgo.io/
- Docker Compose公式ドキュメント: https://docs.docker.com/compose/
