# MySQL対応のための修正箇所と作業一覧

## 概要

PostgreSQLが主のデータベースだが、MySQLでも動作するように修正するための分析結果。

## 修正が必要な箇所

### 1. データベース接続設定（database.yaml）

**現状**: 各環境の`config/{env}/database.yaml`で`driver: postgres`を指定

**対応**: 
- `driver: mysql`に変更するだけで動作する（`server/internal/db/connection.go`で既にMySQLドライバーがサポート済み）
- `config.ShardConfig.GetDSN()`は既にMySQLに対応済み（`server/internal/config/config.go`）
- ただし、MySQLのDSNに`charset=utf8mb4&loc=Local`を追加することを推奨（現状は`parseTime=true`のみ）

**分離が必要な設定ファイル**:
- `config/develop/database.yaml` → `config/develop/database.mysql.yaml`（オプション）
- `config/staging/database.yaml` → `config/staging/database.mysql.yaml`（オプション）
- `config/production/database.yaml` → `config/production/database.mysql.yaml`（オプション）
- `config/test/database.yaml` → `config/test/database.mysql.yaml`（オプション）

**注意**: 環境変数や設定ファイルの読み込み方法によっては、1つのファイルで両方の設定を保持することも可能

---

### 2. マイグレーションファイル（SQL構文の違い）

**現状**: `db/migrations/`配下のSQLファイルがPostgreSQL構文で記述されている

**PostgreSQL固有の構文**:
- `SERIAL` → MySQLでは`AUTO_INCREMENT`（INTEGER型と組み合わせ）
- `character varying(32)` → MySQLでは`VARCHAR(32)`
- `ON CONFLICT DO NOTHING` → MySQLでは`INSERT IGNORE`または`INSERT ... ON DUPLICATE KEY UPDATE`
- `TRUNCATE TABLE ... RESTART IDENTITY CASCADE` → MySQLでは`TRUNCATE TABLE`（AUTO_INCREMENTは自動リセット）

**対応方法**:
1. **方法A**: マイグレーションディレクトリを分離
   - `db/migrations/master/` → PostgreSQL用
   - `db/migrations/master-mysql/` → MySQL用
   - `db/migrations/sharding_1/` → PostgreSQL用
   - `db/migrations/sharding_1-mysql/` → MySQL用

2. **方法B**: Atlasの機能を活用して、HCLスキーマから自動生成（推奨）
   - `db/schema/master.hcl`からPostgreSQL/MySQL両方のSQLを生成
   - Atlasがデータベース固有の構文を自動変換

**修正が必要なファイル**:
- `db/migrations/master/20260108145414_initial_schema.sql`
- `db/migrations/master/20260108145415_seed_data.sql`
- `db/migrations/sharding_1/20260108145537_initial_schema.sql`
- `db/migrations/sharding_2/20260108145546_initial_schema.sql`
- `db/migrations/sharding_3/20260108145548_initial_schema.sql`
- `db/migrations/sharding_4/20260108145549_initial_schema.sql`

**主な変換例**:
```sql
-- PostgreSQL
CREATE TABLE "dm_news" (
  "id" serial NOT NULL,
  "title" text NOT NULL,
  ...
);

-- MySQL
CREATE TABLE `dm_news` (
  `id` INT AUTO_INCREMENT NOT NULL,
  `title` TEXT NOT NULL,
  ...
);
```

```sql
-- PostgreSQL
INSERT INTO goadmin_roles (...) VALUES (...)
ON CONFLICT DO NOTHING;

-- MySQL
INSERT IGNORE INTO goadmin_roles (...) VALUES (...);
```

---

### 3. テストコード（testutil/db.go）

**現状**: PostgreSQL固有の構文を使用

**修正が必要な箇所**:

#### 3.1. `InitMasterSchema()`関数
- `SERIAL PRIMARY KEY` → `INT AUTO_INCREMENT PRIMARY KEY`
- `TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP` → 両方で動作するが、MySQLでは`TIMESTAMP`の挙動が異なる場合がある

#### 3.2. `InitShardingSchema()`関数
- `TEXT PRIMARY KEY` → MySQLでも動作するが、`VARCHAR`の方が適切な場合がある
- `TIMESTAMP DEFAULT CURRENT_TIMESTAMP` → 両方で動作

#### 3.3. `clearDatabaseTables()`関数
- `pg_tables`システムカタログ → MySQLでは`INFORMATION_SCHEMA.TABLES`
- `TRUNCATE TABLE ... RESTART IDENTITY CASCADE` → MySQLでは`TRUNCATE TABLE`（AUTO_INCREMENTは自動リセット）

**修正方法**:
- データベースドライバーを判定して、適切なSQLを実行するように修正
- または、GORMの機能を活用してデータベース固有の処理を抽象化

**修正が必要なファイル**:
- `server/test/testutil/db.go`

---

### 4. Atlas設定ファイル（atlas.hcl）

**現状**: PostgreSQL用のURLとdev環境が設定されている

**修正が必要な箇所**:
- `url = "postgres://..."` → `url = "mysql://..."`または`url = "mysql://user:pass@tcp(host:port)/dbname"`
- `dev = "docker://postgres/15/dev?search_path=public"` → `dev = "docker://mysql/8/dev"`

**分離が必要な設定ファイル**:
- `config/develop/atlas.hcl` → `config/develop/atlas.mysql.hcl`（オプション）
- `config/staging/atlas.hcl` → `config/staging/atlas.mysql.hcl`（オプション）
- `config/production/atlas.hcl` → `config/production/atlas.mysql.hcl`（オプション）
- `config/test/atlas.hcl` → `config/test/atlas.mysql.hcl`（オプション）

**修正が必要なファイル**:
- `config/develop/atlas.hcl`
- `config/staging/atlas.hcl`
- `config/production/atlas.hcl`
- `config/test/atlas.hcl`

---

### 5. Docker Compose設定

**現状**: `docker-compose.postgres.yml`でPostgreSQLコンテナを定義

**対応**:
- `docker-compose.mysql.yml`を作成してMySQLコンテナを定義
- ポート番号はPostgreSQLと重複しないように設定（例: 3306, 3307, 3308, 3309, 3310）

**作成が必要なファイル**:
- `docker-compose.mysql.yml`

---

### 6. スクリプト（start-postgres.sh, migrate-test.sh）

**現状**: PostgreSQL専用のスクリプト

**対応**:
- `scripts/start-mysql.sh`を作成
- `scripts/migrate-test.sh`をMySQL対応版に修正、または`scripts/migrate-test-mysql.sh`を作成

**修正が必要なファイル**:
- `scripts/start-postgres.sh` → 参考にして`scripts/start-mysql.sh`を作成
- `scripts/migrate-test.sh` → MySQL対応版を作成

---

### 7. DSN生成ロジック（config.ShardConfig）

**現状**: `config.ShardConfig.GetDSN()`でPostgreSQL/MySQL用のDSNを生成（既に対応済み）

**実装状況**:
- ✅ PostgreSQL: `host=... port=... user=... password=... dbname=... sslmode=disable`
- ✅ MySQL: `user:pass@tcp(host:port)/dbname?parseTime=true`
- ⚠️ 推奨: MySQLのDSNに`charset=utf8mb4&loc=Local`を追加

**改善提案**:
- MySQLのDSN生成時に`charset=utf8mb4&loc=Local`を追加することを推奨

**確認済みファイル**:
- `server/internal/config/config.go`（`GetDSN()`メソッド、351-368行目）

---

### 8. その他のPostgreSQL固有の機能

**確認が必要な箇所**:
- GORMの機能で自動的に変換される部分（型マッピングなど）
- トランザクション処理（両方で動作するはず）
- 接続プール設定（両方で動作するはず）

---

## 設定ファイルの分離方針

### オプション1: 環境ごとに分離（推奨）

各環境（develop, staging, production, test）ごとにPostgreSQL用とMySQL用の設定ファイルを分離：

```
config/
  develop/
    database.yaml          # PostgreSQL用（デフォルト）
    database.mysql.yaml    # MySQL用
    atlas.hcl              # PostgreSQL用（デフォルト）
    atlas.mysql.hcl         # MySQL用
  staging/
    database.yaml
    database.mysql.yaml
    atlas.hcl
    atlas.mysql.hcl
  production/
    database.yaml.example
    database.mysql.yaml.example
    atlas.hcl
    atlas.mysql.hcl
  test/
    database.yaml
    database.mysql.yaml
    atlas.hcl
    atlas.mysql.hcl
```

### オプション2: 環境変数で切り替え

1つの設定ファイルで、環境変数（例: `DB_DRIVER=mysql`）で切り替え：

```
config/
  develop/
    database.yaml    # driver: ${DB_DRIVER:-postgres}
    atlas.hcl         # 環境変数でURLを切り替え
```

### オプション3: 別ディレクトリに分離

```
config/
  develop/
    postgres/
      database.yaml
      atlas.hcl
    mysql/
      database.yaml
      atlas.hcl
```

---

## 作業の優先順位

### Phase 1: 基本対応（必須）
1. **DSN生成ロジックの改善**（`config.ShardConfig.GetDSN()`に`charset=utf8mb4&loc=Local`を追加）
2. **テストコードの修正**（`server/test/testutil/db.go`）
3. **MySQL用設定ファイルの作成**（`config/{env}/database.mysql.yaml`）
4. **MySQL用Docker Composeの作成**（`docker-compose.mysql.yml`）
5. **MySQL用スクリプトの作成**（`scripts/start-mysql.sh`）

### Phase 2: マイグレーション対応
1. **Atlas設定のMySQL対応**（`config/{env}/atlas.mysql.hcl`）
2. **マイグレーションファイルのMySQL版作成**（またはAtlasで自動生成）

### Phase 3: ドキュメント・CI/CD対応
1. **READMEの更新**（MySQL対応の手順追加）
2. **CI/CDパイプラインの更新**（MySQLテストの追加）

---

## 主な構文の違いまとめ

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

---

## 注意事項

1. **外部キー制約**: データ分散環境では外部キー制約を使わない設計になっているため、問題なし
2. **トランザクション**: GORMが抽象化しているため、基本的に問題なし
3. **接続プール**: 両方で動作するが、設定値の最適化が必要な場合がある
4. **文字セット**: MySQLでは`utf8mb4`を明示的に指定することを推奨
5. **タイムゾーン**: MySQLでは`parseTime=True&loc=Local`をDSNに含めることを推奨
