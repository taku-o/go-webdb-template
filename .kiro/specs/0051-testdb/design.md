# テスト用データベース導入の設計書

## 1. 概要

### 1.1 設計の目的
要件定義書に基づき、テスト実行時に使用する専用のデータベースを導入するための詳細設計を定義する。開発環境のデータベースと分離し、テスト開始時にデータベースをクリアする機能を実装する。

### 1.2 設計の範囲
- テスト用データベース設定ファイル（`config/test/database.yaml`）の作成
- テスト用Atlas設定ファイル（`config/test/atlas.hcl`）の作成
- テスト実行時の設定ファイル読み込み機能の実装
- テスト開始時のデータベースクリア機能の実装
- 既存テストコードの修正

### 1.3 設計方針
- **設定ファイルベース**: ハードコードされた設定を削除し、設定ファイルから読み込む
- **環境変数による切り替え**: `APP_ENV=test`でテスト環境の設定を読み込む
- **データベース分離**: テスト用データベース名に`_test`サフィックスを付与
- **高速クリア**: TRUNCATEを使用してデータベースを高速にクリア
- **後方互換性**: 既存のテストコードへの影響を最小限に抑える

## 2. アーキテクチャ設計

### 2.1 設定ファイル構成

```
config/
├── develop/
│   ├── database.yaml      # 開発環境用データベース設定
│   └── atlas.hcl          # 開発環境用Atlas設定
└── test/
    ├── database.yaml      # テスト環境用データベース設定（新規作成）
    └── atlas.hcl          # テスト環境用Atlas設定（新規作成）
```

### 2.2 設定読み込みフロー

```
┌─────────────────────────────────────────────────────────────┐
│              テスト実行時の設定読み込みフロー                  │
└─────────────────────────────────────────────────────────────┘
                          │
                          ▼
        ┌─────────────────────────────────┐
        │  APP_ENV=test を設定            │
        │  (環境変数またはテストコード内)  │
        └─────────────────────────────────┘
                          │
                          ▼
        ┌─────────────────────────────────┐
        │  config.Load() を呼び出し        │
        │  - config/test/config.yaml      │
        │  - config/test/database.yaml    │
        └─────────────────────────────────┘
                          │
                          ▼
        ┌─────────────────────────────────┐
        │  Config構造体に設定を読み込み    │
        │  - データベース名: *_test       │
        │  - その他の設定: 開発環境と同じ  │
        └─────────────────────────────────┘
                          │
                          ▼
        ┌─────────────────────────────────┐
        │  db.NewGroupManager(cfg)         │
        │  テスト用データベースに接続        │
        └─────────────────────────────────┘
```

### 2.3 データベースクリアフロー

```
┌─────────────────────────────────────────────────────────────┐
│              テスト開始時のデータベースクリアフロー             │
└─────────────────────────────────────────────────────────────┘
                          │
                          ▼
        ┌─────────────────────────────────┐
        │  SetupTestGroupManager()        │
        │  または                          │
        │  SetupTestGroupManager8Sharding()│
        └─────────────────────────────────┘
                          │
                          ▼
        ┌─────────────────────────────────┐
        │  ClearTestDatabase() を呼び出し │
        │  (自動実行または手動実行)        │
        └─────────────────────────────────┘
                          │
        ┌─────────────────┴─────────────────┐
        │                                   │
        ▼                                   ▼
┌──────────────────┐            ┌──────────────────┐
│  マスタDB      │            │  シャーディングDB │
│  クリア処理       │            │  クリア処理        │
│                  │            │                   │
│  1. テーブル一覧  │            │  1. 各DBの接続取得 │
│     取得          │            │  2. テーブル一覧  │
│  2. TRUNCATE実行 │            │     取得          │
│                  │            │  3. TRUNCATE実行  │
└──────────────────┘            └──────────────────┘
```

## 3. 実装設計

### 3.1 設定ファイルの作成

#### 3.1.1 `config/test/database.yaml`の作成

**ファイル構成**: `config/develop/database.yaml`をベースに、データベース名に`_test`サフィックスを付与

**主要な変更点**:
- マスターデータベース名: `webdb_master` → `webdb_master_test`
- シャーディングデータベース名: `webdb_sharding_1` → `webdb_sharding_1_test`（同様に2-4も）
- その他の設定（ホスト、ポート、ユーザー、パスワード等）は開発環境と同じ

**実装例**:
```yaml
# PostgreSQL設定（テスト環境用）
# テスト用データベースは事前に作成し、マイグレーションを実行しておく必要がある
database:
  groups:
    master:
      - id: 1
        driver: postgres
        host: localhost
        port: 5432
        user: webdb
        password: webdb
        name: webdb_master_test
        max_connections: 25
        max_idle_connections: 5
        connection_max_lifetime: 1h

    sharding:
      databases:
        - id: 1
          driver: postgres
          host: localhost
          port: 5433
          user: webdb
          password: webdb
          name: webdb_sharding_1_test
          max_connections: 25
          max_idle_connections: 5
          connection_max_lifetime: 1h
          table_range: [0, 3]
        # ... 同様に2-4も定義
```

#### 3.1.2 `config/test/atlas.hcl`の作成

**ファイル構成**: `config/develop/atlas.hcl`をベースに、データベース名に`_test`サフィックスを付与

**主要な変更点**:
- 各環境のURLにテスト用データベース名を使用
- マイグレーションディレクトリは開発環境と同じ設定を使用

**実装例**:
```hcl
// テスト環境用Atlas設定ファイル (PostgreSQL)

// マスターデータベース用環境
env "master" {
  src = "file://db/schema/master.hcl"
  url = "postgres://webdb:webdb@localhost:5432/webdb_master_test?sslmode=disable"
  dev = "docker://postgres/15/dev?search_path=public"

  migration {
    dir = "file://db/migrations/master"
  }
}

// シャーディングDB 1
env "sharding_1" {
  src = "file://db/schema/sharding_1"
  url = "postgres://webdb:webdb@localhost:5433/webdb_sharding_1_test?sslmode=disable"
  dev = "docker://postgres/15/dev?search_path=public"

  migration {
    dir = "file://db/migrations/sharding_1"
  }
}
// ... 同様に2-4も定義
```

### 3.2 設定読み込み機能の実装

#### 3.2.1 既存の`config.Load()`関数の活用

**実装方針**: 既存の`config.Load()`関数は環境変数`APP_ENV`を使用して環境別の設定を読み込むため、テスト実行時に`APP_ENV=test`を設定することで`config/test/`ディレクトリから設定を読み込む。

**実装場所**: `server/test/testutil/db.go`

**実装方法**:
1. テスト実行前に`APP_ENV=test`を設定（環境変数または`os.Setenv`）
2. `config.Load()`を呼び出して設定を読み込む
3. 読み込んだ設定を使用して`db.NewGroupManager`を作成

**実装コード例**:
```go
// LoadTestConfig はテスト環境の設定を読み込む
func LoadTestConfig() (*config.Config, error) {
    // 既存のAPP_ENVを保存
    oldEnv := os.Getenv("APP_ENV")
    defer func() {
        if oldEnv != "" {
            os.Setenv("APP_ENV", oldEnv)
        } else {
            os.Unsetenv("APP_ENV")
        }
    }()
    
    // テスト環境を設定
    os.Setenv("APP_ENV", "test")
    
    // 設定を読み込む
    return config.Load()
}
```

#### 3.2.2 `SetupTestGroupManager`関数の修正

**変更前**: ハードコードされた設定でデータベース接続を作成

**変更後**: 設定ファイルから読み込んだ設定を使用

**実装コード例**:
```go
func SetupTestGroupManager(t *testing.T, dbCount int, tablesPerDB int) *db.GroupManager {
    // 設定ファイルから読み込む
    cfg, err := LoadTestConfig()
    require.NoError(t, err)
    
    // 設定からGroupManagerを作成
    manager, err := db.NewGroupManager(cfg)
    require.NoError(t, err)
    
    // データベースをクリア
    ClearTestDatabase(t, manager)
    
    // スキーマを初期化（既存の実装を維持）
    // ...
    
    return manager
}
```

#### 3.2.3 `SetupTestGroupManager8Sharding`関数の修正

**変更内容**: `SetupTestGroupManager`と同様に、設定ファイルから読み込むように修正

### 3.3 データベースクリア機能の実装

#### 3.3.1 `ClearTestDatabase`関数の実装

**実装場所**: `server/test/testutil/db.go`

**機能**:
- マスターデータベースの全テーブルのデータを削除
- シャーディングデータベースの全テーブルのデータを削除
- テーブル構造は維持（TRUNCATEを使用）

**実装コード例**:
```go
// ClearTestDatabase はテスト用データベースの全テーブルのデータをクリアする
func ClearTestDatabase(t *testing.T, manager *db.GroupManager) {
    // マスターデータベースのクリア
    masterConn, err := manager.GetMasterConnection()
    require.NoError(t, err)
    clearDatabaseTables(t, masterConn.DB)
    
    // シャーディングデータベースのクリア
    connections := manager.GetAllShardingConnections()
    for _, conn := range connections {
        clearDatabaseTables(t, conn.DB)
    }
}

// clearDatabaseTables は指定されたデータベースの全テーブルのデータをクリアする
func clearDatabaseTables(t *testing.T, database *gorm.DB) {
    // テーブル一覧を取得
    var tables []string
    err := database.Raw(`
        SELECT tablename 
        FROM pg_tables 
        WHERE schemaname = 'public'
    `).Scan(&tables).Error
    require.NoError(t, err)
    
    // 各テーブルをTRUNCATE
    for _, table := range tables {
        err := database.Exec(fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", table)).Error
        require.NoError(t, err)
    }
}
```

**注意事項**:
- `RESTART IDENTITY`: シーケンス（SERIAL等）をリセット
- `CASCADE`: 依存するオブジェクトもクリア（このシステムでは外部キー制約を使用しないため、実質的な影響はない）

#### 3.3.2 クリア処理の実行タイミング

**自動実行**: `SetupTestGroupManager`関数内で自動実行

**手動実行**: 必要に応じて、テストコード内で明示的に呼び出し可能

**実装方針**: 
- デフォルトでは自動実行
- テストの独立性を保証するため、各テストの開始時にクリア

### 3.4 エラーハンドリング

#### 3.4.1 設定ファイル読み込みエラー

**エラーケース**:
- `config/test/database.yaml`が存在しない
- 設定ファイルの形式が不正

**対応**:
- 明確なエラーメッセージを表示
- テスト用データベースの設定ファイルが必要であることを示す

#### 3.4.2 データベース接続エラー

**エラーケース**:
- テスト用データベースが存在しない
- データベースへの接続に失敗

**対応**:
- 明確なエラーメッセージを表示
- テスト用データベースの事前作成が必要であることを示す

#### 3.4.3 データベースクリアエラー

**エラーケース**:
- テーブル一覧の取得に失敗
- TRUNCATEの実行に失敗

**対応**:
- エラーをログに記録
- テストを失敗させる（`require.NoError`を使用）

## 4. テスト設計

### 4.1 設定ファイル読み込みのテスト

#### 4.1.1 テスト内容
- `LoadTestConfig()`関数が正常に設定を読み込むことを確認
- データベース名に`_test`サフィックスが付与されていることを確認
- その他の設定が開発環境と同じであることを確認

#### 4.1.2 テストコード例
```go
func TestLoadTestConfig(t *testing.T) {
    cfg, err := LoadTestConfig()
    require.NoError(t, err)
    
    // マスターデータベース名を確認
    assert.Equal(t, "webdb_master_test", cfg.Database.Groups.Master[0].Name)
    
    // シャーディングデータベース名を確認
    assert.Equal(t, "webdb_sharding_1_test", cfg.Database.Groups.Sharding.Databases[0].Name)
    
    // その他の設定を確認（ホスト、ポート等）
    assert.Equal(t, "localhost", cfg.Database.Groups.Master[0].Host)
    assert.Equal(t, 5432, cfg.Database.Groups.Master[0].Port)
}
```

### 4.2 データベースクリア機能のテスト

#### 4.2.1 テスト内容
- `ClearTestDatabase()`関数が正常にデータベースをクリアすることを確認
- マスターデータベースの全テーブルがクリアされることを確認
- シャーディングデータベースの全テーブルがクリアされることを確認
- テーブル構造が維持されることを確認

#### 4.2.2 テストコード例
```go
func TestClearTestDatabase(t *testing.T) {
    manager := SetupTestGroupManager(t, 4, 8)
    defer CleanupTestGroupManager(manager)
    
    // テストデータを挿入
    masterConn, err := manager.GetMasterConnection()
    require.NoError(t, err)
    masterConn.DB.Exec("INSERT INTO dm_news (title, content) VALUES ('test', 'test')")
    
    // データが存在することを確認
    var count int64
    masterConn.DB.Model(&dm_news{}).Count(&count)
    assert.Greater(t, count, int64(0))
    
    // データベースをクリア
    ClearTestDatabase(t, manager)
    
    // データがクリアされたことを確認
    masterConn.DB.Model(&dm_news{}).Count(&count)
    assert.Equal(t, int64(0), count)
    
    // テーブルが存在することを確認（構造が維持される）
    var exists bool
    masterConn.DB.Raw(`
        SELECT EXISTS (
            SELECT FROM information_schema.tables 
            WHERE table_schema = 'public' 
            AND table_name = 'dm_news'
        )
    `).Scan(&exists)
    assert.True(t, exists)
}
```

### 4.3 統合テスト

#### 4.3.1 テスト内容
- 既存のテストが正常に動作することを確認
- テスト用データベースが使用されることを確認
- 開発環境のデータベースに影響がないことを確認

#### 4.3.2 テスト実行方法
```bash
# テスト用データベースを事前に作成
createdb -U webdb webdb_master_test
createdb -U webdb webdb_sharding_1_test
# ... 同様に2-4も作成

# マイグレーションを実行
APP_ENV=test atlas migrate apply --env master
APP_ENV=test atlas migrate apply --env sharding_1
# ... 同様に2-4も実行

# テストを実行
APP_ENV=test go test ./server/test/...
```

## 5. 既存機能への影響

### 5.1 既存テストへの影響

**影響内容**:
- `SetupTestGroupManager`関数と`SetupTestGroupManager8Sharding`関数の実装が変更される
- テスト実行時に`APP_ENV=test`を設定する必要がある（または関数内で自動設定）

**対応**:
- 既存のテストコードは変更不要（関数のシグネチャは維持）
- テスト実行時に環境変数を設定するか、関数内で自動設定

### 5.2 開発環境への影響

**影響**: なし
- テスト用データベースは別データベースのため、開発環境のデータベースに影響しない

### 5.3 本番環境への影響

**影響**: なし
- テスト用データベースはローカル環境のみで使用

## 6. 実装上の注意事項

### 6.1 設定ファイル作成の注意事項
- **データベース名の一貫性**: すべてのデータベース名に`_test`サフィックスを付与する
- **設定の完全性**: 開発環境設定のすべての項目を含める
- **コメント**: 設定ファイルに適切なコメントを追加する（テスト用データベースの事前作成が必要であることを明記）

### 6.2 テスト実装の注意事項
- **環境変数の管理**: `APP_ENV`の設定と復元を適切に行う
- **エラーハンドリング**: 設定ファイルが存在しない場合やデータベース接続に失敗した場合、明確なエラーメッセージを表示
- **データベースクリアのタイミング**: テスト開始時に必ずクリアする

### 6.3 動作確認の注意事項
- **データベースの事前準備**: テスト実行前にテスト用データベースを作成し、マイグレーションを実行する
- **開発環境との分離**: テスト実行後、開発環境のデータベースに影響がないことを確認
- **パフォーマンス**: データベースクリア処理のパフォーマンスを確認（10秒以内を目安）

## 7. 実装チェックリスト

### 7.1 設定ファイル作成
- [ ] `config/test/database.yaml`を作成
- [ ] `config/test/atlas.hcl`を作成
- [ ] データベース名に`_test`サフィックスが付与されている
- [ ] その他の設定が開発環境と同じである

### 7.2 実装
- [ ] `LoadTestConfig()`関数を実装
- [ ] `SetupTestGroupManager()`関数を修正
- [ ] `SetupTestGroupManager8Sharding()`関数を修正
- [ ] `ClearTestDatabase()`関数を実装
- [ ] `clearDatabaseTables()`関数を実装

### 7.3 テスト
- [ ] 設定ファイル読み込みのテストを実装
- [ ] データベースクリア機能のテストを実装
- [ ] 既存のテストが正常に動作することを確認

### 7.4 動作確認
- [ ] テスト用データベースが事前に作成されている
- [ ] テスト用データベースに対してマイグレーションが実行されている
- [ ] テストを実行して、テスト用データベースが使用されることを確認
- [ ] テスト実行後、開発環境のデータベースに影響がないことを確認

## 8. 参考情報

### 8.1 関連Issue
- GitHub Issue #105: テスト用データベースの導入

### 8.2 既存実装の参考
- **開発環境設定**: `config/develop/database.yaml`
- **開発環境Atlas設定**: `config/develop/atlas.hcl`
- **設定読み込み**: `server/internal/config/config.go`の`Load()`関数
- **テストユーティリティ**: `server/test/testutil/db.go`

### 8.3 技術スタック
- **言語**: Go
- **データベース**: PostgreSQL
- **ORM**: GORM
- **マイグレーションツール**: Atlas
- **設定管理**: Viper
- **設定ファイル形式**: YAML (database.yaml), HCL (atlas.hcl)

### 8.4 関連ドキュメント
- `config/develop/database.yaml`: 開発環境のデータベース設定
- `config/develop/atlas.hcl`: 開発環境のAtlas設定
- `server/test/testutil/db.go`: テスト用のデータベース接続設定
- `server/internal/config/config.go`: 設定読み込み関数
