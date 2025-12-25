# シャーディング数増加設計書

## 1. 概要

### 1.1 設計の目的
要件定義書に基づき、データベースのシャーディング数を2から4に増やすための詳細設計を定義する。既存のシャーディングロジックは動的にシャード数を検出するため、コード変更は不要で、設定ファイルとマイグレーションファイルの追加のみで実現する。

### 1.2 設計の範囲
- 設定ファイル（develop/staging/production環境）へのshard3とshard4の追加設計
- マイグレーションファイル（shard3とshard4）の作成設計
- ドキュメント（`docs/Sharding.md`）の更新設計
- 既存のシャーディングロジックとの統合確認
- テスト戦略

### 1.3 設計方針
- **コード変更なし**: 既存のシャーディングロジックは動的にシャード数を検出するため、コード変更は不要
- **設定ファイル追加**: 各環境の設定ファイルにshard3とshard4を追加
- **マイグレーションファイル作成**: shard1と同じスキーマでshard3とshard4のマイグレーションファイルを作成
- **一貫性の維持**: 既存のshard1とshard2の設定パターンに従う

## 2. アーキテクチャ設計

### 2.1 ディレクトリ構造

#### 2.1.1 変更前の構造
```
config/
├── develop/
│   └── database.yaml      # shard1, shard2のみ
├── staging/
│   └── database.yaml      # shard1, shard2のみ
└── production/
    └── database.yaml.example  # shard1, shard2のみ

db/
└── migrations/
    ├── shard1/
    │   ├── 001_init.sql
    │   ├── 002_goadmin.sql
    │   └── 003_menu.sql
    └── shard2/
        └── 001_init.sql
```

#### 2.1.2 変更後の構造
```
config/
├── develop/
│   └── database.yaml      # shard1, shard2, shard3, shard4
├── staging/
│   └── database.yaml      # shard1, shard2, shard3, shard4
└── production/
    └── database.yaml.example  # shard1, shard2, shard3, shard4

db/
└── migrations/
    ├── shard1/
    │   ├── 001_init.sql
    │   ├── 002_goadmin.sql
    │   └── 003_menu.sql
    ├── shard2/
    │   └── 001_init.sql
    ├── shard3/              # 新規追加
    │   └── 001_init.sql     # 新規追加
    └── shard4/              # 新規追加
        └── 001_init.sql     # 新規追加
```

### 2.2 シャーディングロジックの動作

#### 2.2.1 シャード数の自動検出
既存の`HashBasedSharding`は、設定ファイルから読み込まれた`shards`配列の長さから自動的にシャード数を決定する：

```go
// server/internal/db/manager.go（既存コード）
func NewGORMManager(cfg *config.Config) (*GORMManager, error) {
    shardCount := len(cfg.Database.Shards)  // 設定ファイルから自動検出
    sharding := db.NewHashBasedSharding(shardCount)
    // ...
}
```

#### 2.2.2 シャードIDの計算
`HashBasedSharding.GetShardID()`は、ハッシュ値とシャード数からシャードIDを計算する：

```go
func (h *HashBasedSharding) GetShardID(key int64) int {
    hash := fnv.New32a()
    hash.Write([]byte(fmt.Sprintf("%d", key)))
    hashValue := hash.Sum32()
    shardID := int(hashValue%uint32(h.shardCount)) + 1  // 1からNの範囲
    return shardID
}
```

**動作例**:
- 2シャードの場合: `hashValue % 2 + 1` → 1または2
- 4シャードの場合: `hashValue % 4 + 1` → 1, 2, 3, または4

#### 2.2.3 データ分散の変化
シャード数が2から4に増えると、既存データのシャード割り当てが変わる可能性がある：

```
user_id=1: hash(1) % 2 + 1 = 2  →  hash(1) % 4 + 1 = 2（同じ）
user_id=2: hash(2) % 2 + 1 = 1  →  hash(2) % 4 + 1 = 3（変更）
user_id=3: hash(3) % 2 + 1 = 2  →  hash(3) % 4 + 1 = 4（変更）
user_id=4: hash(4) % 2 + 1 = 1  →  hash(4) % 4 + 1 = 1（同じ）
```

**注意**: データ損失を許容するため、既存データの移行は行わない。

### 2.3 設定ファイルの読み込みフロー

```
┌─────────────────────────────────────────────────────────────┐
│              1. アプリケーション起動                           │
│              server/cmd/server/main.go                        │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              2. 設定ファイル読み込み                           │
│              config.Load()                                  │
│              - config/{env}/database.yaml を読み込み         │
│              - shards配列を解析                              │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              3. シャード数の自動検出                           │
│              len(cfg.Database.Shards)                       │
│              - 2シャード → 4シャードに自動的に変更            │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              4. GORMManager初期化                            │
│              NewGORMManager(cfg)                            │
│              - HashBasedSharding(shardCount=4) を作成       │
│              - 4つのシャードに接続                           │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              5. シャーディングロジック動作                     │
│              GetShardID(user_id)                            │
│              - hash(user_id) % 4 + 1 でシャードIDを計算      │
│              - 1, 2, 3, 4のいずれかを返す                    │
└─────────────────────────────────────────────────────────────┘
```

### 2.4 既存アーキテクチャとの統合

```
┌─────────────────────────────────────────────────────────────┐
│              HTTPリクエスト                                  │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              Handler (internal/api/handler)                 │
│              - UserHandler                                  │
│              - PostHandler                                  │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              Service (internal/service)                     │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              Repository (internal/repository)               │
│              - GetConnectionByKey(user_id)                  │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              GORMManager (internal/db)                      │
│              - HashBasedSharding(shardCount=4)              │
│              - GetShardID(user_id) → 1, 2, 3, または4       │
│              - GetConnection(shardID)                       │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              GORMConnection (internal/db)                   │
│              - Shard 1, 2, 3, または4への接続                │
└─────────────────────────────────────────────────────────────┘
```

**重要なポイント**:
- 既存のコードは変更不要
- 設定ファイルの`shards`配列にshard3とshard4を追加するだけで、自動的に4シャード構成になる
- シャーディングロジックは動的にシャード数を検出するため、コード変更は不要

## 3. コンポーネント設計

### 3.1 設定ファイル設計

#### 3.1.1 開発環境（develop）
**ファイル**: `config/develop/database.yaml`

```yaml
database:
  shards:
    # 既存のshard1（変更なし）
    - id: 1
      driver: sqlite3
      dsn: ./data/shard1.db
      writer_dsn: ./data/shard1.db
      reader_dsns:
        - ./data/shard1.db
      reader_policy: random
      max_connections: 10
      max_idle_connections: 5
      connection_max_lifetime: 300s
    
    # 既存のshard2（変更なし）
    - id: 2
      driver: sqlite3
      dsn: ./data/shard2.db
      writer_dsn: ./data/shard2.db
      reader_dsns:
        - ./data/shard2.db
      reader_policy: random
      max_connections: 10
      max_idle_connections: 5
      connection_max_lifetime: 300s
    
    # 新規追加: shard3
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
    
    # 新規追加: shard4
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

#### 3.1.2 ステージング環境（staging）
**ファイル**: `config/staging/database.yaml`

```yaml
database:
  shards:
    # 既存のshard1（変更なし）
    - id: 1
      driver: postgres
      host: staging-db-shard1.example.com
      port: 5432
      name: app_db_shard1
      user: staging_user
      password: ${DB_PASSWORD_SHARD1}
      writer_dsn: host=staging-db-shard1-writer.example.com port=5432 user=staging_user password=${DB_PASSWORD_SHARD1} dbname=app_db_shard1 sslmode=require
      reader_dsns:
        - host=staging-db-shard1-reader1.example.com port=5432 user=staging_user password=${DB_PASSWORD_SHARD1} dbname=app_db_shard1 sslmode=require
      reader_policy: random
      max_connections: 25
      max_idle_connections: 10
      connection_max_lifetime: 600s
    
    # 既存のshard2（変更なし）
    - id: 2
      driver: postgres
      host: staging-db-shard2.example.com
      port: 5432
      name: app_db_shard2
      user: staging_user
      password: ${DB_PASSWORD_SHARD2}
      writer_dsn: host=staging-db-shard2-writer.example.com port=5432 user=staging_user password=${DB_PASSWORD_SHARD2} dbname=app_db_shard2 sslmode=require
      reader_dsns:
        - host=staging-db-shard2-reader1.example.com port=5432 user=staging_user password=${DB_PASSWORD_SHARD2} dbname=app_db_shard2 sslmode=require
      reader_policy: random
      max_connections: 25
      max_idle_connections: 10
      connection_max_lifetime: 600s
    
    # 新規追加: shard3
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
    
    # 新規追加: shard4
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

#### 3.1.3 本番環境（production）
**ファイル**: `config/production/database.yaml.example`

ステージング環境と同様の構造で、ホスト名を`prod-db-shard*`に変更し、複数のReaderを設定可能にする。

### 3.2 マイグレーションファイル設計

#### 3.2.1 shard3のマイグレーションファイル
**ファイル**: `db/migrations/shard3/001_init.sql`

```sql
-- Shard 3 初期化スクリプト

-- Users テーブル
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);

-- Posts テーブル
CREATE TABLE IF NOT EXISTS posts (
    id INTEGER PRIMARY KEY,
    user_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- インデックス
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_posts_user_id ON posts(user_id);
CREATE INDEX IF NOT EXISTS idx_posts_created_at ON posts(created_at);
```

#### 3.2.2 shard4のマイグレーションファイル
**ファイル**: `db/migrations/shard4/001_init.sql`

shard3と同じ内容で、コメントのみ「Shard 4」に変更。

### 3.3 ドキュメント更新設計

#### 3.3.1 `docs/Sharding.md`の更新内容
1. **シャーディング数の説明**: 2シャードから4シャードに変更
2. **設定例の更新**: 4シャード構成の設定例を追加
3. **データ分散の例**: 4シャードでのデータ分散例を追加

## 4. データフロー設計

### 4.1 シャード選択のフロー

```
┌─────────────────────────────────────────────────────────────┐
│              1. ユーザーID取得                                │
│              user_id = 123                                  │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              2. ハッシュ値計算                               │
│              hash = FNV-1a(user_id)                        │
│              hashValue = 0x12345678                        │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              3. シャードID計算                               │
│              shardID = hashValue % 4 + 1                    │
│              shardID = 0x12345678 % 4 + 1 = 2              │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              4. シャード接続取得                              │
│              conn = GetConnection(shardID=2)                │
│              → shard2への接続を返す                          │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              5. データベース操作                              │
│              SELECT/INSERT/UPDATE/DELETE                   │
│              → shard2で実行                                 │
└─────────────────────────────────────────────────────────────┘
```

### 4.2 クロスシャードクエリのフロー

```
┌─────────────────────────────────────────────────────────────┐
│              1. 全ユーザー取得要求                            │
│              GetAllUsers()                                  │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              2. 全シャードをループ                            │
│              for shardID := 1; shardID <= 4; shardID++    │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              3. 各シャードから並列取得                         │
│              - shard1から取得                                │
│              - shard2から取得                                │
│              - shard3から取得                                │
│              - shard4から取得                                │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              4. 結果をマージ                                  │
│              allUsers = append(shard1, shard2, shard3, shard4)
└─────────────────────────────────────────────────────────────┘
```

## 5. 実装詳細

### 5.1 設定ファイルの実装

#### 5.1.1 開発環境設定ファイル
- **ファイル**: `config/develop/database.yaml`
- **変更内容**: shard3とshard4を追加
- **既存設定**: shard1とshard2は変更しない

#### 5.1.2 ステージング環境設定ファイル
- **ファイル**: `config/staging/database.yaml`
- **変更内容**: shard3とshard4を追加
- **環境変数**: `DB_PASSWORD_SHARD3`、`DB_PASSWORD_SHARD4`を使用

#### 5.1.3 本番環境設定ファイル例
- **ファイル**: `config/production/database.yaml.example`
- **変更内容**: shard3とshard4を追加
- **環境変数**: `DB_PASSWORD_SHARD3`、`DB_PASSWORD_SHARD4`を使用

### 5.2 マイグレーションファイルの実装

#### 5.2.1 shard3のマイグレーションファイル
- **ディレクトリ作成**: `db/migrations/shard3/`
- **ファイル作成**: `db/migrations/shard3/001_init.sql`
- **内容**: shard1の`001_init.sql`をベースに、コメントを「Shard 3」に変更

#### 5.2.2 shard4のマイグレーションファイル
- **ディレクトリ作成**: `db/migrations/shard4/`
- **ファイル作成**: `db/migrations/shard4/001_init.sql`
- **内容**: shard1の`001_init.sql`をベースに、コメントを「Shard 4」に変更

### 5.3 ドキュメントの実装

#### 5.3.1 `docs/Sharding.md`の更新
- **シャーディング数の説明**: 「2シャード」→「4シャード」に更新
- **設定例の更新**: 4シャード構成の設定例を追加
- **データ分散の例**: 4シャードでのデータ分散例を追加

## 6. エラーハンドリング

### 6.1 既存のエラーハンドリング
既存のシャーディングロジックは変更不要のため、既存のエラーハンドリングがそのまま適用される：

- **設定ファイル読み込みエラー**: `config.Load()`でエラーハンドリング済み
- **データベース接続エラー**: `NewGORMConnection()`でエラーハンドリング済み
- **シャード選択エラー**: `GetShardID()`でエラーハンドリング済み

### 6.2 新規追加時のエラーハンドリング
- **設定ファイルの構文エラー**: YAMLパーサーが自動的に検出
- **マイグレーションファイルのエラー**: SQL実行時にエラーが発生

## 7. テスト戦略

### 7.1 既存テストの動作確認
- **既存のテストコード**: 変更不要
- **4シャードでの動作確認**: 既存のテストが4シャードでも正常に動作することを確認

### 7.2 新規テスト（オプション）
- **設定ファイル読み込みテスト**: 4シャードの設定が正しく読み込まれることを確認
- **シャード選択テスト**: 4シャードでのシャード選択が正しく動作することを確認
- **データ分散テスト**: 新しいデータが4つのシャードに適切に分散されることを確認

### 7.3 手動テスト
- **アプリケーション起動**: 4つのシャードに正常に接続できることを確認
- **データ作成**: 新しいデータが4つのシャードに適切に分散されることを確認
- **クロスシャードクエリ**: 全シャードからデータを取得できることを確認

## 8. 実装上の注意事項

### 8.1 設定ファイルの更新
- **既存設定の保持**: shard1とshard2の設定は変更しない
- **一貫性の維持**: shard3とshard4の設定は既存のパターンに従う
- **環境変数の使用**: ステージング/本番環境では環境変数を使用

### 8.2 マイグレーションファイルの作成
- **スキーマの一貫性**: shard1とshard2と同じスキーマを使用
- **コメントの更新**: 「Shard 1」を「Shard 3」「Shard 4」に変更

### 8.3 データベース接続
- **自動検出**: 設定ファイルから自動的に4シャードが検出される
- **接続確認**: 4つのシャードすべてに接続できることを確認

### 8.4 データ損失の許容
- **既存データの移行**: 行わない
- **データ損失**: 許容する
- **データリセット**: 必要に応じて実行可能

## 9. 参考情報

### 9.1 既存実装
- `server/internal/db/sharding.go`: `HashBasedSharding`の実装
- `server/internal/config/config.go`: 設定構造体の定義
- `server/internal/db/connection.go`: データベース接続処理
- `db/migrations/shard1/001_init.sql`: shard1の初期化スクリプト

### 9.2 既存ドキュメント
- `docs/Sharding.md`: シャーディング戦略の詳細
- `config/develop/database.yaml`: 開発環境設定ファイル
- `config/staging/database.yaml`: ステージング環境設定ファイル

### 9.3 関連Issue
- GitHub Issue #17: データベースのシャーディング数を増やす

