# MySQL対応の設計書

## 1. 概要

### 1.1 設計の目的
要件定義書に基づき、PostgreSQLが主のデータベースだが、MySQLでも動作するように修正するための詳細設計を定義する。PostgreSQLとMySQLの両方に対応することで、データベース選択の柔軟性を提供する。

### 1.2 設計の範囲
- MySQL用の設定ファイル（database.yaml, atlas.hcl）の設計
- MySQL用のマイグレーションファイルの生成・移植方法の設計
- テストコードのMySQL対応設計
- MySQL用のDocker Compose設定の設計
- MySQL用のスクリプトの設計
- DSN生成ロジックの改善設計
- 環境情報の取得機能（config.yamlにDB_TYPEを追加）の設計

### 1.3 設計方針
- **設定ファイルの分離**: PostgreSQL用とMySQL用の設定ファイルを分離（環境ごと）
- **後方互換性**: 既存のPostgreSQL機能に影響を与えない（追加のみ）
- **一貫性**: 既存のコードスタイルと一貫性を保つ
- **段階的な実装**: 各機能ごとに実装し、動作確認を行う
- **既存コードの最小限の変更**: 既存のPostgreSQLコードは可能な限り変更せず、MySQL対応を追加する

## 2. アーキテクチャ設計

### 2.1 全体構成

```
設定ファイル読み込み
  ↓
config.yaml から DB_TYPE を取得
  ↓
DB_TYPE に応じて適切な database.yaml を読み込む
  ├── DB_TYPE=postgresql → database.yaml (既存)
  └── DB_TYPE=mysql → database.mysql.yaml (新規)
  ↓
GroupManager の初期化
  ↓
データベース接続（PostgreSQL または MySQL）
  ├── PostgreSQL: gorm.io/driver/postgres
  └── MySQL: gorm.io/driver/mysql
```

### 2.2 ディレクトリ構造

```
config/
├── develop/
│   ├── config.yaml              # DB_TYPE フィールドを追加
│   ├── database.yaml            # PostgreSQL用（既存）
│   ├── database.mysql.yaml      # MySQL用（新規）
│   ├── atlas.hcl                # PostgreSQL用（既存）
│   └── atlas.mysql.hcl          # MySQL用（新規）
├── staging/
│   ├── config.yaml              # DB_TYPE フィールドを追加
│   ├── database.yaml            # PostgreSQL用（既存）
│   ├── database.mysql.yaml      # MySQL用（新規）
│   ├── atlas.hcl                # PostgreSQL用（既存）
│   └── atlas.mysql.hcl          # MySQL用（新規）
├── production/
│   ├── config.yaml.example      # DB_TYPE フィールドを追加
│   ├── database.yaml.example    # PostgreSQL用（既存）
│   ├── database.mysql.yaml.example  # MySQL用（新規）
│   ├── atlas.hcl                # PostgreSQL用（既存）
│   └── atlas.mysql.hcl         # MySQL用（新規）
└── test/
    ├── config.yaml              # DB_TYPE フィールドを追加
    ├── database.yaml            # PostgreSQL用（既存）
    ├── database.mysql.yaml      # MySQL用（新規）
    ├── atlas.hcl                # PostgreSQL用（既存）
    └── atlas.mysql.hcl          # MySQL用（新規）

db/
└── migrations/
    ├── master/                  # PostgreSQL用（既存）
    ├── master-mysql/            # MySQL用（新規）
    ├── sharding_1/              # PostgreSQL用（既存）
    ├── sharding_1-mysql/       # MySQL用（新規）
    ├── sharding_2/              # PostgreSQL用（既存）
    ├── sharding_2-mysql/       # MySQL用（新規）
    ├── sharding_3/              # PostgreSQL用（既存）
    ├── sharding_3-mysql/        # MySQL用（新規）
    ├── sharding_4/              # PostgreSQL用（既存）
    ├── sharding_4-mysql/       # MySQL用（新規）
    ├── view_master/             # PostgreSQL用（既存）
    └── view_master-mysql/       # MySQL用（新規）

server/
├── internal/
│   ├── config/
│   │   └── config.go            # GetDSN() メソッドの改善
│   └── test/
│       └── testutil/
│           └── db.go            # MySQL用関数の追加
├── docker-compose.mysql.yml     # MySQL用（新規）
└── scripts/
    ├── start-mysql.sh           # MySQL用（新規）
    └── migrate-test-mysql.sh    # MySQL用（新規）
```

## 3. 詳細設計

### 3.1 データベース接続設定（database.yaml）

#### 3.1.1 MySQL用設定ファイルの設計

**ファイル**: `config/{env}/database.mysql.yaml`

**設計内容**:
- PostgreSQL用の`database.yaml`をベースに、`driver: mysql`に変更
- ポート番号をMySQL用に変更（3306, 3307, 3308, 3309, 3310）
- データベース名はPostgreSQLと同じ（`webdb_master`, `webdb_sharding_1`など）
- 接続プール設定はPostgreSQLと同じ値を維持

**設計例** (`config/develop/database.mysql.yaml`):
```yaml
# MySQL設定（開発環境用）
# MySQLの起動: ./scripts/start-mysql.sh start
# マイグレーション: ./scripts/migrate.sh
database:
  groups:
    master:
      - id: 1
        driver: mysql
        host: localhost
        port: 3306
        user: webdb
        password: webdb
        name: webdb_master
        max_connections: 25
        max_idle_connections: 5
        connection_max_lifetime: 1h

    sharding:
      databases:
        # 論理シャード 1: テーブル _000-003 → mysql-sharding-1
        - id: 1
          driver: mysql
          host: localhost
          port: 3307
          user: webdb
          password: webdb
          name: webdb_sharding_1
          max_connections: 25
          max_idle_connections: 5
          connection_max_lifetime: 1h
          table_range: [0, 3]
        # 論理シャード 2: テーブル _004-007 → mysql-sharding-1
        - id: 2
          driver: mysql
          host: localhost
          port: 3307
          user: webdb
          password: webdb
          name: webdb_sharding_1
          max_connections: 25
          max_idle_connections: 5
          connection_max_lifetime: 1h
          table_range: [4, 7]
        # 論理シャード 3: テーブル _008-011 → mysql-sharding-2
        - id: 3
          driver: mysql
          host: localhost
          port: 3308
          user: webdb
          password: webdb
          name: webdb_sharding_2
          max_connections: 25
          max_idle_connections: 5
          connection_max_lifetime: 1h
          table_range: [8, 11]
        # 論理シャード 4: テーブル _012-015 → mysql-sharding-2
        - id: 4
          driver: mysql
          host: localhost
          port: 3308
          user: webdb
          password: webdb
          name: webdb_sharding_2
          max_connections: 25
          max_idle_connections: 5
          connection_max_lifetime: 1h
          table_range: [12, 15]
        # 論理シャード 5: テーブル _016-019 → mysql-sharding-3
        - id: 5
          driver: mysql
          host: localhost
          port: 3309
          user: webdb
          password: webdb
          name: webdb_sharding_3
          max_connections: 25
          max_idle_connections: 5
          connection_max_lifetime: 1h
          table_range: [16, 19]
        # 論理シャード 6: テーブル _020-023 → mysql-sharding-3
        - id: 6
          driver: mysql
          host: localhost
          port: 3309
          user: webdb
          password: webdb
          name: webdb_sharding_3
          max_connections: 25
          max_idle_connections: 5
          connection_max_lifetime: 1h
          table_range: [20, 23]
        # 論理シャード 7: テーブル _024-027 → mysql-sharding-4
        - id: 7
          driver: mysql
          host: localhost
          port: 3310
          user: webdb
          password: webdb
          name: webdb_sharding_4
          max_connections: 25
          max_idle_connections: 5
          connection_max_lifetime: 1h
          table_range: [24, 27]
        # 論理シャード 8: テーブル _028-031 → mysql-sharding-4
        - id: 8
          driver: mysql
          host: localhost
          port: 3310
          user: webdb
          password: webdb
          name: webdb_sharding_4
          max_connections: 25
          max_idle_connections: 5
          connection_max_lifetime: 1h
          table_range: [28, 31]

      tables:
        - name: dm_users
          suffix_count: 32
        - name: dm_posts
          suffix_count: 32
```

#### 3.1.2 環境情報の取得機能の設計

**ファイル**: `config/{env}/config.yaml`

**設計内容**:
- `DB_TYPE`フィールドを追加
- デフォルト値は`postgresql`（既存の動作を維持）
- `mysql`を指定した場合、`database.mysql.yaml`を読み込む

**設計例** (`config/develop/config.yaml`):
```yaml
# データベースタイプの指定
# postgresql: PostgreSQLを使用（デフォルト）
# mysql: MySQLを使用
DB_TYPE: postgresql  # または mysql
```

**設定ファイル読み込みロジックの設計**:
- `config.Load()`関数内で、`config.yaml`から`DB_TYPE`を読み込む
- `DB_TYPE`が`mysql`の場合、`database.mysql.yaml`を読み込む
- `DB_TYPE`が`postgresql`または未指定の場合、`database.yaml`を読み込む（既存の動作）

**実装イメージ**:
```go
// config.Load() 内での処理
dbType := cfg.DBType // config.yaml から読み込んだ値
if dbType == "" {
    dbType = "postgresql" // デフォルト値
}

var databaseFile string
if dbType == "mysql" {
    databaseFile = "database.mysql.yaml"
} else {
    databaseFile = "database.yaml"
}

// databaseFile を読み込む
```

### 3.2 マイグレーションファイル（SQL構文の違い）

#### 3.2.1 MySQL用マイグレーションディレクトリの設計

**ディレクトリ構造**:
- `db/migrations/master-mysql/`: マスターデータベース用
- `db/migrations/sharding_1-mysql/`: シャーディング1用
- `db/migrations/sharding_2-mysql/`: シャーディング2用
- `db/migrations/sharding_3-mysql/`: シャーディング3用
- `db/migrations/sharding_4-mysql/`: シャーディング4用
- `db/migrations/view_master-mysql/`: ビューマスター用（必要に応じて）

#### 3.2.2 Atlasコマンドによる自動生成の設計

**スキーマファイル（initial_schema.sql）の生成**:
- Atlasコマンドを使用して、HCLスキーマからMySQL用のSQLを生成
- **方針**: まず既存のHCLスキーマ（`db/schema/master.hcl`など）をそのまま使用して試行
  - Atlasは`dev`環境で指定されたデータベースタイプに応じて、型を自動変換する（例: `serial` → `INT AUTO_INCREMENT`）
  - ただし、`schema.public`の参照など、PostgreSQL固有の構文が含まれている場合は、MySQL用に別のHCLスキーマが必要になる可能性がある
- **問題が発生した場合**: MySQL用のHCLスキーマ（`db/schema/master-mysql.hcl`など）を作成して分離

**生成コマンドの例（既存HCLスキーマを使用する場合）**:
```bash
# マスターデータベース用
atlas migrate diff \
  --env mysql_master \
  --to file://db/schema/master.hcl \
  --dir file://db/migrations/master-mysql

# シャーディング1用
atlas migrate diff \
  --env mysql_sharding_1 \
  --to file://db/schema/sharding_1 \
  --dir file://db/migrations/sharding_1-mysql
```

**生成コマンドの例（MySQL用HCLスキーマを分離する場合）**:
```bash
# マスターデータベース用
atlas migrate diff \
  --env mysql_master \
  --to file://db/schema/master-mysql.hcl \
  --dir file://db/migrations/master-mysql

# シャーディング1用
atlas migrate diff \
  --env mysql_sharding_1 \
  --to file://db/schema/sharding_1-mysql \
  --dir file://db/migrations/sharding_1-mysql
```

**注意事項**:
- `schema.public`の参照: MySQLにはスキーマの概念がないため、Atlasが適切に処理するか確認が必要
- `type = serial`: AtlasがMySQL用に`INT AUTO_INCREMENT`に自動変換するか確認が必要
- 実際にAtlasコマンドを実行して、既存のHCLスキーマでMySQL用のSQLが生成できるか確認する
- 生成できない場合は、MySQL用のHCLスキーマを作成して分離する

**Atlas設定ファイルの設計**:
- `config/{env}/atlas.mysql.hcl`でMySQL用の環境を定義
- `url`をMySQL用に設定
- `dev`をMySQL用に設定（`docker://mysql/8/dev`）

#### 3.2.3 手動移植が必要なファイルの設計

**seed_data.sqlの移植**:
- `db/migrations/master/20260108145415_seed_data.sql`をMySQL構文に変換
- 主な変換内容:
  - `ON CONFLICT DO NOTHING` → `INSERT IGNORE`
  - ダブルクォート `"` → バッククォート `` ` ``
  - その他のPostgreSQL固有の構文をMySQL構文に変換

**設計例**:
```sql
-- PostgreSQL版
INSERT INTO goadmin_roles (id, name, slug, created_at, updated_at) VALUES
    (1, 'Administrator', 'administrator', NOW(), NOW()),
    (2, 'Operator', 'operator', NOW(), NOW())
ON CONFLICT DO NOTHING;

-- MySQL版
INSERT IGNORE INTO `goadmin_roles` (`id`, `name`, `slug`, `created_at`, `updated_at`) VALUES
    (1, 'Administrator', 'administrator', NOW(), NOW()),
    (2, 'Operator', 'operator', NOW(), NOW());
```

**view_masterのSQLの移植**:
- `db/migrations/view_master/20260103030225_create_dm_news_view.sql`をMySQL構文に変換
- MySQLのビュー構文に合わせて変換（必要に応じて）

### 3.3 テストコード（testutil/db.go）

#### 3.3.1 MySQL用スキーマ初期化関数の設計

**関数**: `InitMySQLMasterSchema(t *testing.T, database *gorm.DB)`

**設計内容**:
- PostgreSQL用の`InitMasterSchema()`を参考に実装
- MySQL用のSQL構文を使用
- 主な変換:
  - `SERIAL PRIMARY KEY` → `INT AUTO_INCREMENT PRIMARY KEY`
  - `TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP` → `TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP`（両方で動作）

**設計例**:
```go
// InitMySQLMasterSchema initializes the master database schema for MySQL
func InitMySQLMasterSchema(t *testing.T, database *gorm.DB) {
	schema := `
		CREATE TABLE IF NOT EXISTS dm_news (
			id INT AUTO_INCREMENT PRIMARY KEY,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			author_id INT,
			published_at TIMESTAMP,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
	`
	err := database.Exec(schema).Error
	require.NoError(t, err)
}
```

**関数**: `InitMySQLShardingSchema(t *testing.T, database *gorm.DB, startTable, endTable int)`

**設計内容**:
- PostgreSQL用の`InitShardingSchema()`を参考に実装
- MySQL用のSQL構文を使用
- 主な変換:
  - `TEXT PRIMARY KEY` → `VARCHAR(32) PRIMARY KEY`（より適切）
  - `TIMESTAMP DEFAULT CURRENT_TIMESTAMP` → `TIMESTAMP DEFAULT CURRENT_TIMESTAMP`（両方で動作）

**設計例**:
```go
// InitMySQLShardingSchema initializes the sharding database schema for MySQL
func InitMySQLShardingSchema(t *testing.T, database *gorm.DB, startTable, endTable int) {
	for i := startTable; i <= endTable; i++ {
		suffix := fmt.Sprintf("%03d", i)

		usersSchema := fmt.Sprintf(`
			CREATE TABLE IF NOT EXISTS dm_users_%s (
				id VARCHAR(32) PRIMARY KEY,
				name TEXT NOT NULL,
				email TEXT NOT NULL UNIQUE,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			);
		`, suffix)
		err := database.Exec(usersSchema).Error
		require.NoError(t, err)

		postsSchema := fmt.Sprintf(`
			CREATE TABLE IF NOT EXISTS dm_posts_%s (
				id VARCHAR(32) PRIMARY KEY,
				user_id VARCHAR(32) NOT NULL,
				title TEXT NOT NULL,
				content TEXT NOT NULL,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			);
		`, suffix)
		err = database.Exec(postsSchema).Error
		require.NoError(t, err)
	}
}
```

#### 3.3.2 MySQL用データクリア関数の設計

**関数**: `clearMySQLDatabaseTables(t *testing.T, database *gorm.DB)`

**設計内容**:
- PostgreSQL用の`clearDatabaseTables()`を参考に実装
- MySQL用のSQL構文を使用
- 主な変換:
  - `pg_tables` → `INFORMATION_SCHEMA.TABLES`
  - `TRUNCATE TABLE ... RESTART IDENTITY CASCADE` → `TRUNCATE TABLE`（AUTO_INCREMENTは自動リセット）

**設計例**:
```go
// clearMySQLDatabaseTables clears all tables in a MySQL database
func clearMySQLDatabaseTables(t *testing.T, database *gorm.DB) {
	// テーブル一覧を取得
	var tables []string
	err := database.Raw(`
		SELECT table_name
		FROM INFORMATION_SCHEMA.TABLES
		WHERE table_schema = DATABASE()
	`).Scan(&tables).Error
	require.NoError(t, err)

	// 各テーブルをTRUNCATE
	for _, tableName := range tables {
		err := database.Exec(fmt.Sprintf("TRUNCATE TABLE `%s`", tableName)).Error
		require.NoError(t, err)
	}
}
```

#### 3.3.3 データベース判定機能の設計

**修正対象**: `SetupTestGroupManager()`関数

**設計内容**:
- 設定ファイルから読み込んだデータベースドライバーを判定
- `driver == "postgres"`の場合、既存の関数を使用
- `driver == "mysql"`の場合、MySQL用の関数を使用

**設計例**:
```go
func SetupTestGroupManager(t *testing.T, dbCount int, tablesPerDB int) *db.GroupManager {
	// ロックを取得
	fileLock, err := AcquireTestLock(t)
	if err != nil {
		t.Fatalf("Failed to acquire test lock: %v", err)
	}
	defer func() {
		if err := fileLock.Unlock(); err != nil {
			t.Logf("Warning: failed to unlock test lock: %v", err)
		}
	}()

	// 設定ファイルから読み込む
	cfg, err := LoadTestConfig()
	require.NoError(t, err)

	// 設定からGroupManagerを作成
	manager, err := db.NewGroupManager(cfg)
	require.NoError(t, err)

	// データベースをクリア
	ClearTestDatabase(t, manager)

	// データベースドライバーを判定
	masterConn, err := manager.GetMasterConnection()
	require.NoError(t, err)
	
	driver := masterConn.Driver // または設定から取得

	// Initialize master database schema
	if driver == "mysql" {
		InitMySQLMasterSchema(t, masterConn.DB)
	} else {
		InitMasterSchema(t, masterConn.DB)
	}

	// Initialize sharding database schemas
	tableRanges := map[int][2]int{
		1: {0, 7},   // Entries 1,2 -> tables 0-7
		3: {8, 15},  // Entries 3,4 -> tables 8-15
		5: {16, 23}, // Entries 5,6 -> tables 16-23
		7: {24, 31}, // Entries 7,8 -> tables 24-31
	}

	connections := manager.GetAllShardingConnections()
	for _, conn := range connections {
		tableRange, ok := tableRanges[conn.ShardID]
		if ok {
			if driver == "mysql" {
				InitMySQLShardingSchema(t, conn.DB, tableRange[0], tableRange[1])
			} else {
				InitShardingSchema(t, conn.DB, tableRange[0], tableRange[1])
			}
		}
	}

	return manager
}
```

**修正対象**: `ClearTestDatabase()`関数

**設計内容**:
- データベースドライバーを判定して、適切なクリア関数を呼び出す

**設計例**:
```go
func ClearTestDatabase(t *testing.T, manager *db.GroupManager) {
	// マスターデータベースのクリア
	masterConn, err := manager.GetMasterConnection()
	require.NoError(t, err)
	
	driver := masterConn.Driver
	if driver == "mysql" {
		clearMySQLDatabaseTables(t, masterConn.DB)
	} else {
		clearDatabaseTables(t, masterConn.DB)
	}

	// シャーディングデータベースのクリア
	connections := manager.GetAllShardingConnections()
	for _, conn := range connections {
		if driver == "mysql" {
			clearMySQLDatabaseTables(t, conn.DB)
		} else {
			clearDatabaseTables(t, conn.DB)
		}
	}
}
```

### 3.4 Atlas設定ファイル（atlas.hcl）

#### 3.4.1 MySQL用Atlas設定ファイルの設計

**ファイル**: `config/{env}/atlas.mysql.hcl`

**設計内容**:
- PostgreSQL用の`atlas.hcl`をベースに、MySQL用に変更
- `url`をMySQL用に設定（`user:pass@tcp(host:port)/dbname?charset=utf8mb4&parseTime=true&loc=Local`）
- `dev`をMySQL用に設定（`docker://mysql/8/dev`）
- マイグレーションディレクトリをMySQL用に設定（`db/migrations/{database}-mysql`）

**設計例** (`config/develop/atlas.mysql.hcl`):
```hcl
// 開発環境用Atlas設定ファイル (MySQL)

// マスターデータベース用環境
env "master" {
  src = "file://db/schema/master.hcl"
  url = "webdb:webdb@tcp(localhost:3306)/webdb_master?charset=utf8mb4&parseTime=true&loc=Local"
  dev = "docker://mysql/8/dev"

  migration {
    dir = "file://db/migrations/master-mysql"
  }
}

// シャーディングDB 1
env "sharding_1" {
  src = "file://db/schema/sharding_1"
  url = "webdb:webdb@tcp(localhost:3307)/webdb_sharding_1?charset=utf8mb4&parseTime=true&loc=Local"
  dev = "docker://mysql/8/dev"

  migration {
    dir = "file://db/migrations/sharding_1-mysql"
  }
}

// シャーディングDB 2
env "sharding_2" {
  src = "file://db/schema/sharding_2"
  url = "webdb:webdb@tcp(localhost:3308)/webdb_sharding_2?charset=utf8mb4&parseTime=true&loc=Local"
  dev = "docker://mysql/8/dev"

  migration {
    dir = "file://db/migrations/sharding_2-mysql"
  }
}

// シャーディングDB 3
env "sharding_3" {
  src = "file://db/schema/sharding_3"
  url = "webdb:webdb@tcp(localhost:3309)/webdb_sharding_3?charset=utf8mb4&parseTime=true&loc=Local"
  dev = "docker://mysql/8/dev"

  migration {
    dir = "file://db/migrations/sharding_3-mysql"
  }
}

// シャーディングDB 4
env "sharding_4" {
  src = "file://db/schema/sharding_4"
  url = "webdb:webdb@tcp(localhost:3310)/webdb_sharding_4?charset=utf8mb4&parseTime=true&loc=Local"
  dev = "docker://mysql/8/dev"

  migration {
    dir = "file://db/migrations/sharding_4-mysql"
  }
}
```

### 3.5 Docker Compose設定

#### 3.5.1 MySQL用Docker Composeファイルの設計

**ファイル**: `docker-compose.mysql.yml`

**設計内容**:
- PostgreSQL用の`docker-compose.postgres.yml`をベースに、MySQL用に変更
- MySQL 8のDockerイメージを使用
- ポート番号をMySQL用に設定（3306, 3307, 3308, 3309, 3310）
- 環境変数をMySQL用に設定（`MYSQL_USER`, `MYSQL_PASSWORD`, `MYSQL_DATABASE`）
- ボリュームマウントをMySQL用に設定（`./mysql/data/{database_name}:/var/lib/mysql`）
- ヘルスチェックをMySQL用に設定（`mysqladmin ping`）

**設計例**:
```yaml
services:
  mysql-master:
    image: mysql:8
    container_name: mysql-master
    ports:
      - "3306:3306"
    environment:
      MYSQL_USER: webdb
      MYSQL_PASSWORD: webdb
      MYSQL_DATABASE: webdb_master
      MYSQL_ROOT_PASSWORD: rootpassword
    volumes:
      - ./mysql/data/master:/var/lib/mysql
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "mysqladmin", "ping", "-h", "localhost", "-u", "webdb", "-pwebdb"]
      interval: 10s
      timeout: 5s
      retries: 5
    networks:
      - mysql-network
    command: --default-authentication-plugin=mysql_native_password

  mysql-sharding-1:
    image: mysql:8
    container_name: mysql-sharding-1
    ports:
      - "3307:3306"
    environment:
      MYSQL_USER: webdb
      MYSQL_PASSWORD: webdb
      MYSQL_DATABASE: webdb_sharding_1
      MYSQL_ROOT_PASSWORD: rootpassword
    volumes:
      - ./mysql/data/sharding_1:/var/lib/mysql
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "mysqladmin", "ping", "-h", "localhost", "-u", "webdb", "-pwebdb"]
      interval: 10s
      timeout: 5s
      retries: 5
    networks:
      - mysql-network
    command: --default-authentication-plugin=mysql_native_password

  mysql-sharding-2:
    image: mysql:8
    container_name: mysql-sharding-2
    ports:
      - "3308:3306"
    environment:
      MYSQL_USER: webdb
      MYSQL_PASSWORD: webdb
      MYSQL_DATABASE: webdb_sharding_2
      MYSQL_ROOT_PASSWORD: rootpassword
    volumes:
      - ./mysql/data/sharding_2:/var/lib/mysql
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "mysqladmin", "ping", "-h", "localhost", "-u", "webdb", "-pwebdb"]
      interval: 10s
      timeout: 5s
      retries: 5
    networks:
      - mysql-network
    command: --default-authentication-plugin=mysql_native_password

  mysql-sharding-3:
    image: mysql:8
    container_name: mysql-sharding-3
    ports:
      - "3309:3306"
    environment:
      MYSQL_USER: webdb
      MYSQL_PASSWORD: webdb
      MYSQL_DATABASE: webdb_sharding_3
      MYSQL_ROOT_PASSWORD: rootpassword
    volumes:
      - ./mysql/data/sharding_3:/var/lib/mysql
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "mysqladmin", "ping", "-h", "localhost", "-u", "webdb", "-pwebdb"]
      interval: 10s
      timeout: 5s
      retries: 5
    networks:
      - mysql-network
    command: --default-authentication-plugin=mysql_native_password

  mysql-sharding-4:
    image: mysql:8
    container_name: mysql-sharding-4
    ports:
      - "3310:3306"
    environment:
      MYSQL_USER: webdb
      MYSQL_PASSWORD: webdb
      MYSQL_DATABASE: webdb_sharding_4
      MYSQL_ROOT_PASSWORD: rootpassword
    volumes:
      - ./mysql/data/sharding_4:/var/lib/mysql
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "mysqladmin", "ping", "-h", "localhost", "-u", "webdb", "-pwebdb"]
      interval: 10s
      timeout: 5s
      retries: 5
    networks:
      - mysql-network
    command: --default-authentication-plugin=mysql_native_password

networks:
  mysql-network:
    name: mysql-network
    driver: bridge
```

**注意事項**:
- `--default-authentication-plugin=mysql_native_password`を指定して、古い認証方式との互換性を確保
- ボリュームマウントのパスは`./mysql/data/{database_name}`を使用
- ヘルスチェックは`mysqladmin ping`を使用

### 3.6 スクリプト

#### 3.6.1 MySQL用起動スクリプトの設計

**ファイル**: `scripts/start-mysql.sh`

**設計内容**:
- PostgreSQL用の`scripts/start-postgres.sh`をベースに、MySQL用に変更
- `docker-compose.mysql.yml`を使用
- 接続URLをMySQL用に変更

**設計例**:
```bash
#!/bin/bash
set -e

SCRIPT_DIR=$(cd $(dirname $0); pwd)
PROJECT_DIR=$(cd $SCRIPT_DIR/..; pwd)
COMPOSE_FILE="$PROJECT_DIR/docker-compose.mysql.yml"

usage() {
    echo "Usage: $0 {start|stop|status|health}"
    echo ""
    echo "Commands:"
    echo "  start   Start MySQL containers"
    echo "  stop    Stop MySQL containers"
    echo "  status  Show container status"
    echo "  health  Show health check status"
    exit 1
}

start() {
    echo "Starting MySQL containers..."
    docker-compose -f "$COMPOSE_FILE" up -d
    echo "MySQL containers started successfully."
    echo ""
    echo "Connection URLs:"
    echo "  Master:     mysql://webdb:webdb@tcp(localhost:3306)/webdb_master"
    echo "  Sharding 1: mysql://webdb:webdb@tcp(localhost:3307)/webdb_sharding_1"
    echo "  Sharding 2: mysql://webdb:webdb@tcp(localhost:3308)/webdb_sharding_2"
    echo "  Sharding 3: mysql://webdb:webdb@tcp(localhost:3309)/webdb_sharding_3"
    echo "  Sharding 4: mysql://webdb:webdb@tcp(localhost:3310)/webdb_sharding_4"
}

stop() {
    echo "Stopping MySQL containers..."
    docker-compose -f "$COMPOSE_FILE" down
    echo "MySQL containers stopped successfully."
}

status() {
    echo "MySQL container status:"
    docker-compose -f "$COMPOSE_FILE" ps
}

health() {
    echo "MySQL health check status:"
    docker-compose -f "$COMPOSE_FILE" ps --format "table {{.Name}}\t{{.Status}}"
}

if [ $# -eq 0 ]; then
    usage
fi

case "$1" in
    start)
        start
        ;;
    stop)
        stop
        ;;
    status)
        status
        ;;
    health)
        health
        ;;
    *)
        usage
        ;;
esac
```

#### 3.6.2 MySQL用マイグレーションスクリプトの設計

**ファイル**: `scripts/migrate-test-mysql.sh`

**設計内容**:
- PostgreSQL用の`scripts/migrate-test.sh`をベースに、MySQL用に変更
- MySQL接続情報を使用
- MySQL URL形式を構築
- `docker exec`でMySQLコンテナに接続してSQLファイルを適用

**設計例**:
```bash
#!/bin/bash
set -e

# MySQL用マイグレーションスクリプト（テスト環境用）

SCRIPT_DIR=$(cd $(dirname $0); pwd)
PROJECT_DIR=$(cd $SCRIPT_DIR/..; pwd)

# テスト環境用MySQL接続情報
MASTER_HOST="localhost"
MASTER_PORT="3306"
MASTER_USER="webdb"
MASTER_PASSWORD="webdb"
MASTER_DB="webdb_master_test"

SHARDING_1_HOST="localhost"
SHARDING_1_PORT="3307"
SHARDING_1_USER="webdb"
SHARDING_1_PASSWORD="webdb"
SHARDING_1_DB="webdb_sharding_1_test"

SHARDING_2_HOST="localhost"
SHARDING_2_PORT="3308"
SHARDING_2_USER="webdb"
SHARDING_2_PASSWORD="webdb"
SHARDING_2_DB="webdb_sharding_2_test"

SHARDING_3_HOST="localhost"
SHARDING_3_PORT="3309"
SHARDING_3_USER="webdb"
SHARDING_3_PASSWORD="webdb"
SHARDING_3_DB="webdb_sharding_3_test"

SHARDING_4_HOST="localhost"
SHARDING_4_PORT="3310"
SHARDING_4_USER="webdb"
SHARDING_4_PASSWORD="webdb"
SHARDING_4_DB="webdb_sharding_4_test"

# MySQL URL形式を構築
build_mysql_url() {
    local host=$1
    local port=$2
    local user=$3
    local password=$4
    local dbname=$5
    echo "${user}:${password}@tcp(${host}:${port})/${dbname}?charset=utf8mb4&parseTime=true&loc=Local"
}

# マスターデータベースへのマイグレーション
migrate_master() {
    echo "Migrating master database..."
    local url=$(build_mysql_url "$MASTER_HOST" "$MASTER_PORT" "$MASTER_USER" "$MASTER_PASSWORD" "$MASTER_DB")
    
    # マイグレーションディレクトリ内のSQLファイルを実行
    local migration_dir="$PROJECT_DIR/db/migrations/master-mysql"
    if [ -d "$migration_dir" ]; then
        for sql_file in "$migration_dir"/*.sql; do
            if [ -f "$sql_file" ]; then
                echo "Executing: $sql_file"
                docker exec -i mysql-master mysql -u"$MASTER_USER" -p"$MASTER_PASSWORD" "$MASTER_DB" < "$sql_file"
            fi
        done
    fi
}

# シャーディングデータベースへのマイグレーション
migrate_sharding() {
    echo "Migrating sharding databases..."
    
    # Sharding 1
    local url1=$(build_mysql_url "$SHARDING_1_HOST" "$SHARDING_1_PORT" "$SHARDING_1_USER" "$SHARDING_1_PASSWORD" "$SHARDING_1_DB")
    local migration_dir1="$PROJECT_DIR/db/migrations/sharding_1-mysql"
    if [ -d "$migration_dir1" ]; then
        for sql_file in "$migration_dir1"/*.sql; do
            if [ -f "$sql_file" ]; then
                echo "Executing: $sql_file (sharding_1)"
                docker exec -i mysql-sharding-1 mysql -u"$SHARDING_1_USER" -p"$SHARDING_1_PASSWORD" "$SHARDING_1_DB" < "$sql_file"
            fi
        done
    fi
    
    # Sharding 2, 3, 4 も同様に処理
    # ...
}

# メイン処理
main() {
    migrate_master
    migrate_sharding
    echo "Migration completed."
}

main
```

### 3.7 DSN生成ロジックの改善設計

#### 3.7.1 GetDSN()メソッドの改善設計

**ファイル**: `server/internal/config/config.go`

**修正箇所**: `GetDSN()`メソッド（351-368行目）

**改善内容**:
- MySQLのDSNに`charset=utf8mb4`を追加
- MySQLのDSNに`loc=Local`を追加
- 既存の`parseTime=true`は維持
- PostgreSQLのDSN生成は変更しない（後方互換性を維持）

**修正前**:
```go
case "mysql":
    return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
        s.User, s.Password, s.Host, s.Port, s.Name)
```

**修正後**:
```go
case "mysql":
    return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local",
        s.User, s.Password, s.Host, s.Port, s.Name)
```

**変更後のDSN形式**: `user:pass@tcp(host:port)/dbname?charset=utf8mb4&parseTime=true&loc=Local`

**注意事項**:
- `charset=utf8mb4`: UTF-8の4バイト文字（絵文字など）をサポート
- `parseTime=true`: 時刻型を`time.Time`に自動変換
- `loc=Local`: タイムゾーンをローカル時間に設定

## 4. エラーハンドリング設計

### 4.1 データベース接続エラー

**エラーケース**:
- MySQLコンテナが起動していない
- 接続情報が間違っている
- データベースが存在しない

**対応**:
- 接続エラーを適切にログに出力
- エラーメッセージに接続情報を含める（パスワードは除く）

### 4.2 マイグレーションエラー

**エラーケース**:
- SQL構文エラー
- テーブルが既に存在する
- 外部キー制約エラー

**対応**:
- マイグレーションエラーを適切にログに出力
- エラーメッセージにSQLファイル名を含める
- ロールバック処理を実装（可能な限り）

### 4.3 テストエラー

**エラーケース**:
- スキーマ初期化エラー
- データクリアエラー
- データベース判定エラー

**対応**:
- テストエラーを適切にログに出力
- `t.Fatalf()`でテストを失敗させる
- エラーメッセージにデータベースタイプを含める

## 5. テスト設計

### 5.1 設定ファイルのテスト

**テスト内容**:
- `database.mysql.yaml`が正しく読み込める
- `DB_TYPE`が正しく判定される
- 適切な設定ファイルが読み込まれる

### 5.2 マイグレーションのテスト

**テスト内容**:
- AtlasコマンドでMySQL用のSQLが生成される
- 生成されたSQLがMySQLで正常に実行できる
- 手動移植したSQLがMySQLで正常に実行できる

### 5.3 テストコードのテスト

**テスト内容**:
- `InitMySQLMasterSchema()`が正常に動作する
- `InitMySQLShardingSchema()`が正常に動作する
- `clearMySQLDatabaseTables()`が正常に動作する
- `SetupTestGroupManager()`でデータベース判定が正常に動作する

### 5.4 統合テスト

**テスト内容**:
- MySQL環境でテストが正常に実行できる
- PostgreSQL環境でテストが正常に実行できる（既存の動作確認）
- 両方の環境で同じテストが正常に実行できる

## 6. 実装順序

### 6.1 Phase 1: 基本対応（必須）

1. **DSN生成ロジックの改善**
   - `server/internal/config/config.go`の`GetDSN()`メソッドを修正
   - MySQLのDSNに`charset=utf8mb4&loc=Local`を追加

2. **MySQL用設定ファイルの作成**
   - `config/{env}/database.mysql.yaml`を作成（各環境）
   - `config/{env}/config.yaml`に`DB_TYPE`フィールドを追加

3. **MySQL用Docker Composeの作成**
   - `docker-compose.mysql.yml`を作成
   - MySQLコンテナを定義

4. **MySQL用スクリプトの作成**
   - `scripts/start-mysql.sh`を作成
   - `scripts/migrate-test-mysql.sh`を作成

5. **動作確認**
   - MySQLコンテナが正常に起動できる
   - MySQL接続が正常に確立できる

### 6.2 Phase 2: マイグレーション対応

1. **Atlas設定のMySQL対応**
   - `config/{env}/atlas.mysql.hcl`を作成（各環境）

2. **Atlasコマンドでスキーマファイルを生成**
   - マスターデータベース用
   - シャーディングデータベース用（4つ）

3. **手動移植が必要なファイルの移植**
   - `seed_data.sql`をMySQL構文に変換
   - `view_master`のSQLをMySQL構文に変換（必要に応じて）

4. **動作確認**
   - マイグレーションが正常に実行できる

### 6.3 Phase 3: テストコード対応

1. **MySQL用関数の実装**
   - `InitMySQLMasterSchema()`を実装
   - `InitMySQLShardingSchema()`を実装
   - `clearMySQLDatabaseTables()`を実装

2. **データベース判定機能の追加**
   - `SetupTestGroupManager()`を修正
   - `ClearTestDatabase()`を修正

3. **動作確認**
   - MySQL環境でテストが正常に実行できる
   - PostgreSQL環境でテストが正常に実行できる（既存の動作確認）

### 6.4 Phase 4: ドキュメント・CI/CD対応

1. **READMEの更新**
   - MySQL対応の手順を追加
   - MySQL用のコマンドを追加

2. **CI/CDパイプラインの更新**
   - MySQLテストを追加（該当する場合）

## 7. 実装上の注意事項

### 7.1 設定ファイルの管理
- **分離方針**: PostgreSQL用とMySQL用の設定ファイルを分離（環境ごと）
- **デフォルト**: PostgreSQLをデフォルトとして維持
- **環境変数**: `DB_TYPE`環境変数でデータベースタイプを切り替え可能

### 7.2 マイグレーションファイルの管理
- **ディレクトリ分離**: PostgreSQL用とMySQL用のマイグレーションディレクトリを分離
- **Atlasによる自動生成**: スキーマファイル（initial_schema.sql）はAtlasコマンドでMySQL用に自動生成
  - **方針**: まず既存のHCLスキーマ（`db/schema/master.hcl`など）をそのまま使用して試行
  - **問題が発生した場合**: MySQL用のHCLスキーマ（`db/schema/master-mysql.hcl`など）を作成して分離
  - **注意点**: `schema.public`の参照は、MySQLにはスキーマの概念がないため、Atlasが適切に処理するか確認が必要
- **手動移植**: データファイル（seed_data.sql）やビューファイルは手動でMySQL構文に変換して移植
- **構文変換**: PostgreSQL構文からMySQL構文への変換を正確に実施
- **バージョン管理**: マイグレーションファイルのバージョンを一致させる

### 7.3 テストコードの実装
- **関数分離**: PostgreSQL用とMySQL用の関数を分離
- **データベース判定**: 実行時にデータベースタイプを判定して適切な関数を呼び出す
- **エラーハンドリング**: データベース固有のエラーを適切に処理

### 7.4 DSN生成の実装
- **文字セット**: MySQLでは`utf8mb4`を明示的に指定
- **タイムゾーン**: MySQLでは`loc=Local`を指定
- **後方互換性**: 既存のPostgreSQL DSN生成に影響を与えない

### 7.5 Docker Composeの実装
- **ポート番号**: PostgreSQLと重複しないポート番号を使用（3306-3310）
- **ボリューム**: データ永続化のためのボリュームマウントを設定
- **ヘルスチェック**: コンテナの起動確認のためのヘルスチェックを設定
- **認証方式**: `--default-authentication-plugin=mysql_native_password`を指定

## 8. 参考情報

### 8.1 関連ドキュメント
- 要件定義書: `.kiro/specs/0054-mysql/requirements.md`
- 分析結果: `.kiro/specs/0054-mysql/MySQL-Support-Analysis.md`
- アーキテクチャドキュメント: `docs/Architecture.md`
- プロジェクト構造ドキュメント: `docs/Project-Structure.md`

### 8.2 既存実装の参考
- PostgreSQL用設定ファイル: `config/develop/database.yaml`
- PostgreSQL用Atlas設定: `config/develop/atlas.hcl`
- PostgreSQL用Docker Compose: `docker-compose.postgres.yml`
- PostgreSQL用起動スクリプト: `scripts/start-postgres.sh`
- PostgreSQL用マイグレーションスクリプト: `scripts/migrate-test.sh`
- テストユーティリティ: `server/test/testutil/db.go`
- DSN生成ロジック: `server/internal/config/config.go`

### 8.3 技術スタック
- **言語**: Go
- **データベース**: PostgreSQL, MySQL
- **ORM**: GORM（`gorm.io/driver/postgres`, `gorm.io/driver/mysql`）
- **マイグレーション**: Atlas
- **コンテナ**: Docker, Docker Compose

### 8.4 主な構文の違い

| 項目 | PostgreSQL | MySQL |
|------|-----------|-------|
| 自動増分ID | `SERIAL` | `INT AUTO_INCREMENT` |
| 可変長文字列 | `character varying(n)` | `VARCHAR(n)` |
| 固定長文字列 | `CHAR(n)` | `CHAR(n)` |
| テキスト型 | `TEXT` | `TEXT` |
| タイムスタンプ | `TIMESTAMP` | `TIMESTAMP` / `DATETIME` |
| デフォルト値（現在時刻） | `DEFAULT CURRENT_TIMESTAMP` | `DEFAULT CURRENT_TIMESTAMP` |
| 重複時の無視 | `ON CONFLICT DO NOTHING` | `INSERT IGNORE` |
| テーブル一覧取得 | `pg_tables` | `INFORMATION_SCHEMA.TABLES` |
| TRUNCATE（IDリセット） | `TRUNCATE ... RESTART IDENTITY` | `TRUNCATE TABLE`（自動リセット） |
| 引用符 | ダブルクォート `"` | バッククォート `` ` `` |
| DSN形式 | `postgres://user:pass@host:port/dbname` | `user:pass@tcp(host:port)/dbname` |

### 8.5 Atlasコマンドの使用例

**スキーマファイルの生成（既存HCLスキーマを使用する場合）**:
```bash
# マスターデータベース用
atlas migrate diff \
  --env mysql_master \
  --to file://db/schema/master.hcl \
  --dir file://db/migrations/master-mysql

# シャーディング1用
atlas migrate diff \
  --env mysql_sharding_1 \
  --to file://db/schema/sharding_1 \
  --dir file://db/migrations/sharding_1-mysql
```

**スキーマファイルの生成（MySQL用HCLスキーマを分離する場合）**:
```bash
# マスターデータベース用
atlas migrate diff \
  --env mysql_master \
  --to file://db/schema/master-mysql.hcl \
  --dir file://db/migrations/master-mysql

# シャーディング1用
atlas migrate diff \
  --env mysql_sharding_1 \
  --to file://db/schema/sharding_1-mysql \
  --dir file://db/migrations/sharding_1-mysql
```

**マイグレーションの適用**:
```bash
# マスターデータベース用
atlas migrate apply \
  --env mysql_master \
  --dir file://db/migrations/master-mysql

# シャーディング1用
atlas migrate apply \
  --env mysql_sharding_1 \
  --dir file://db/migrations/sharding_1-mysql
```

**注意事項**:
- まず既存のHCLスキーマ（`db/schema/master.hcl`など）をそのまま使用してAtlasコマンドを実行し、MySQL用のSQLが生成できるか確認する
- `schema.public`の参照や`type = serial`など、PostgreSQL固有の構文が含まれている場合、Atlasが適切に変換するか確認が必要
- 生成できない場合や、生成されたSQLに問題がある場合は、MySQL用のHCLスキーマ（`db/schema/master-mysql.hcl`など）を作成して分離する
