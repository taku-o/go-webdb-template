# MySQL対応の要件定義書

## 1. 概要

### 1.1 プロジェクト情報
- **プロジェクト名**: go-webdb-template
- **Feature名**: 0054-mysql
- **作成日**: 2026-01-10

### 1.2 目的
PostgreSQLが主のデータベースだが、MySQLでも動作するように修正する。PostgreSQLとMySQLの両方に対応することで、データベース選択の柔軟性を提供する。

### 1.3 スコープ
- MySQL用の設定ファイル（database.yaml, atlas.hcl）の作成
- MySQL用のマイグレーションファイルの作成
- テストコードのMySQL対応
- MySQL用のDocker Compose設定の作成
- MySQL用のスクリプトの作成
- DSN生成ロジックの改善
- 環境情報の取得機能（config.yamlにDB_TYPEを追加）

**本実装の範囲外**:
- PostgreSQLの既存機能への影響（PostgreSQLは引き続き動作する）
- データベース間の自動マイグレーション（手動でのマイグレーション実行を想定）
- 本番環境での自動切り替え機能（設定ファイルで手動切り替え）

## 2. 背景・現状分析

### 2.1 現在の状況
- **主データベース**: PostgreSQL
- **データベースドライバー**: `gorm.io/driver/postgres`（既存）、`gorm.io/driver/mysql`（既に依存関係に含まれている）
- **接続設定**: `config/{env}/database.yaml`でPostgreSQL接続を定義
- **マイグレーション**: Atlasを使用してPostgreSQL用のSQLを生成
- **テスト環境**: PostgreSQL用のテストデータベースを使用

### 2.2 課題点
1. **設定ファイルのPostgreSQL依存**: 各環境の設定ファイルがPostgreSQL専用
2. **マイグレーションファイルのPostgreSQL構文**: SQLファイルがPostgreSQL固有の構文を使用
3. **テストコードのPostgreSQL依存**: テストユーティリティがPostgreSQL固有のSQLを使用
4. **Docker ComposeのPostgreSQL専用**: PostgreSQLコンテナのみ定義
5. **スクリプトのPostgreSQL専用**: 起動・マイグレーションスクリプトがPostgreSQL専用
6. **DSN生成の改善余地**: MySQLのDSNに`charset=utf8mb4&loc=Local`が未追加

### 2.3 本実装による改善点
1. **データベース選択の柔軟性**: PostgreSQLとMySQLの両方に対応
2. **設定ファイルの分離**: 環境ごとにPostgreSQL用とMySQL用の設定ファイルを用意
3. **マイグレーションの対応**: MySQL用のマイグレーションファイルを作成
4. **テスト環境の対応**: MySQL用のテスト環境を構築可能
5. **開発環境の選択**: 開発者がPostgreSQLまたはMySQLを選択可能

## 3. 機能要件

### 3.1 データベース接続設定（database.yaml）

#### 3.1.1 MySQL用設定ファイルの作成
- **目的**: 各環境でMySQL接続設定を定義できるようにする
- **作成対象**:
  - `config/develop/database.mysql.yaml`
  - `config/staging/database.mysql.yaml`
  - `config/production/database.mysql.yaml.example`
  - `config/test/database.mysql.yaml`
- **設定内容**:
  - `driver: mysql`を指定
  - MySQL接続情報（host, port, user, password, name）を設定
  - 接続プール設定（max_connections, max_idle_connections, connection_max_lifetime）を設定
  - シャーディング設定も同様にMySQL用に設定

#### 3.1.2 環境情報の取得機能
- **目的**: 実行時に使用するデータベースタイプを判定できるようにする
- **実装方法**: `config/{env}/config.yaml`に`DB_TYPE`フィールドを追加
  - `DB_TYPE: postgresql`（デフォルト）
  - `DB_TYPE: mysql`
- **使用箇所**: 設定ファイルの読み込み時に、`DB_TYPE`に応じて適切な`database.yaml`を読み込む

### 3.2 マイグレーションファイル（SQL構文の違い）

#### 3.2.1 MySQL用マイグレーションディレクトリの作成
- **目的**: MySQL用のSQLファイルを管理する
- **作成対象**:
  - `db/migrations/master-mysql/`
  - `db/migrations/sharding_1-mysql/`
  - `db/migrations/sharding_2-mysql/`
  - `db/migrations/sharding_3-mysql/`
  - `db/migrations/sharding_4-mysql/`
  - `db/migrations/view_master-mysql/`

#### 3.2.2 PostgreSQL構文からMySQL構文への変換
- **主な変換内容**:
  - `SERIAL` → `INT AUTO_INCREMENT`
  - `character varying(n)` → `VARCHAR(n)`
  - `ON CONFLICT DO NOTHING` → `INSERT IGNORE`
  - `TRUNCATE TABLE ... RESTART IDENTITY CASCADE` → `TRUNCATE TABLE`
  - ダブルクォート `"` → バッククォート `` ` ``
  - `pg_tables` → `INFORMATION_SCHEMA.TABLES`

#### 3.2.3 マイグレーションファイルの生成・移植

**Atlasコマンドで自動生成するファイル**:
以下のスキーマファイルは、Atlasコマンドを使用してMySQL用に自動生成する：
- `db/migrations/master/20260108145414_initial_schema.sql` → AtlasでMySQL版を生成
- `db/migrations/sharding_1/20260108145537_initial_schema.sql` → AtlasでMySQL版を生成
- `db/migrations/sharding_2/20260108145546_initial_schema.sql` → AtlasでMySQL版を生成
- `db/migrations/sharding_3/20260108145548_initial_schema.sql` → AtlasでMySQL版を生成
- `db/migrations/sharding_4/20260108145549_initial_schema.sql` → AtlasでMySQL版を生成

**手動移植が必要なファイル**:
以下のデータファイルとビューファイルは、手動でMySQL構文に変換して移植する：
- `db/migrations/master/20260108145415_seed_data.sql` → MySQL版を作成（`ON CONFLICT DO NOTHING` → `INSERT IGNORE`など）
- `db/migrations/view_master/20260103030225_create_dm_news_view.sql` → MySQL版を作成（必要に応じて）

### 3.3 テストコード（testutil/db.go）

#### 3.3.1 MySQL用スキーマ初期化関数の作成
- **目的**: テスト環境でMySQLデータベースのスキーマを初期化する
- **作成関数**:
  - `InitMySQLMasterSchema(t *testing.T, database *gorm.DB)`: マスターデータベースのスキーマを初期化
  - `InitMySQLShardingSchema(t *testing.T, database *gorm.DB, startTable, endTable int)`: シャーディングデータベースのスキーマを初期化
- **実装内容**:
  - MySQL用のSQL構文を使用（`INT AUTO_INCREMENT`, `VARCHAR`など）
  - 既存の`InitMasterSchema()`、`InitShardingSchema()`を参考に実装

#### 3.3.2 MySQL用データクリア関数の作成
- **目的**: テスト環境でMySQLデータベースのデータをクリアする
- **作成関数**:
  - `clearMySQLDatabaseTables(t *testing.T, database *gorm.DB)`: MySQLデータベースの全テーブルのデータをクリア
- **実装内容**:
  - `INFORMATION_SCHEMA.TABLES`を使用してテーブル一覧を取得
  - `TRUNCATE TABLE`を使用してデータをクリア（AUTO_INCREMENTは自動リセット）

#### 3.3.3 データベース判定機能の追加
- **目的**: 実行時にデータベースタイプを判定して適切な関数を呼び出す
- **実装方法**: `SetupTestGroupManager()`内で、データベースドライバーを判定
  - `driver == "postgres"` → 既存の関数を使用
  - `driver == "mysql"` → MySQL用の関数を使用

### 3.4 Atlas設定ファイル（atlas.hcl）

#### 3.4.1 MySQL用Atlas設定ファイルの作成
- **目的**: AtlasでMySQL用のマイグレーションを管理できるようにする
- **作成対象**:
  - `config/develop/atlas.mysql.hcl`
  - `config/staging/atlas.mysql.hcl`
  - `config/production/atlas.mysql.hcl`
  - `config/test/atlas.mysql.hcl`
- **設定内容**:
  - `url = "user:pass@tcp(host:port)/dbname?charset=utf8mb4&parseTime=true&loc=Local"`
  - `dev = "docker://mysql/8/dev"`
  - マイグレーションディレクトリをMySQL用に設定

### 3.5 Docker Compose設定

#### 3.5.1 MySQL用Docker Composeファイルの作成
- **目的**: MySQLコンテナを起動できるようにする
- **作成ファイル**: `docker-compose.mysql.yml`
- **設定内容**:
  - MySQL 8のコンテナを定義
  - マスターデータベース: port 3306
  - シャーディングデータベース: port 3307, 3308, 3309, 3310
  - 環境変数: `MYSQL_USER`, `MYSQL_PASSWORD`, `MYSQL_DATABASE`
  - ボリュームマウント: `./mysql/data/{database_name}:/var/lib/mysql/data`
  - ヘルスチェック: `mysqladmin ping`

### 3.6 スクリプト

#### 3.6.1 MySQL用起動スクリプトの作成
- **目的**: MySQLコンテナを起動・停止できるようにする
- **作成ファイル**: `scripts/start-mysql.sh`
- **機能**:
  - `start`: MySQLコンテナを起動
  - `stop`: MySQLコンテナを停止
  - `status`: コンテナの状態を表示
  - `health`: ヘルスチェック状態を表示
- **実装内容**: `scripts/start-postgres.sh`を参考に実装

#### 3.6.2 MySQL用マイグレーションスクリプトの作成
- **目的**: MySQL用のマイグレーションを実行できるようにする
- **作成ファイル**: `scripts/migrate-test-mysql.sh`
- **機能**:
  - マスターデータベースへのマイグレーション実行
  - シャーディングデータベースへのマイグレーション実行
  - 各データベースへのSQLファイルの適用
- **実装内容**: `scripts/migrate-test.sh`を参考に実装

### 3.7 DSN生成ロジックの改善

#### 3.7.1 MySQL DSNの改善
- **目的**: MySQL接続時の文字セットとタイムゾーンを適切に設定する
- **修正対象**: `server/internal/config/config.go`の`GetDSN()`メソッド
- **改善内容**:
  - MySQLのDSNに`charset=utf8mb4`を追加
  - MySQLのDSNに`loc=Local`を追加
  - 既存の`parseTime=true`は維持
- **変更後のDSN形式**: `user:pass@tcp(host:port)/dbname?charset=utf8mb4&parseTime=true&loc=Local`

## 4. 非機能要件

### 4.1 パフォーマンス
- **接続プール**: PostgreSQLと同様の接続プール設定を適用
- **クエリ性能**: MySQLでもPostgreSQLと同等の性能を維持（可能な限り）

### 4.2 信頼性
- **データ整合性**: MySQLでもPostgreSQLと同等のデータ整合性を保証
- **トランザクション**: GORMが抽象化しているため、両方で正常に動作
- **エラーハンドリング**: データベース固有のエラーを適切に処理

### 4.3 保守性
- **コードの可読性**: データベース固有の処理を明確に分離
- **一貫性**: 既存のコードスタイルと一貫性を保つ
- **テスト容易性**: PostgreSQLとMySQLの両方でテスト可能

### 4.4 互換性
- **既存機能**: PostgreSQLの既存機能に影響を与えない
- **後方互換性**: 既存のPostgreSQL設定ファイルは引き続き動作
- **設定ファイル**: 環境変数や設定ファイルの読み込み方法を変更しない（追加のみ）

## 5. 制約事項

### 5.1 技術的制約
- **データベースドライバー**: `gorm.io/driver/mysql`を使用（既に依存関係に含まれている）
- **Atlas**: AtlasがMySQLに対応していることを確認
- **Docker**: MySQL 8のDockerイメージを使用

### 5.2 実装上の制約
- **設定ファイルの分離**: PostgreSQL用とMySQL用の設定ファイルを分離（環境ごと）
- **マイグレーションディレクトリの分離**: PostgreSQL用とMySQL用のマイグレーションディレクトリを分離
- **テストコードの分離**: データベースタイプに応じて適切な関数を呼び出す

### 5.3 動作環境
- **ローカル環境**: ローカル環境でPostgreSQLとMySQLの両方が動作することを確認
- **CI環境**: CI環境でもPostgreSQLとMySQLの両方が動作することを確認（該当する場合）

## 6. 受け入れ基準

### 6.1 データベース接続設定
- [ ] `config/develop/database.mysql.yaml`が作成されている
- [ ] `config/staging/database.mysql.yaml`が作成されている
- [ ] `config/production/database.mysql.yaml.example`が作成されている
- [ ] `config/test/database.mysql.yaml`が作成されている
- [ ] 各設定ファイルで`driver: mysql`が指定されている
- [ ] MySQL接続情報が正しく設定されている
- [ ] `config/{env}/config.yaml`に`DB_TYPE`フィールドが追加されている

### 6.2 マイグレーションファイル
- [ ] `db/migrations/master-mysql/`ディレクトリが作成されている
- [ ] `db/migrations/sharding_1-mysql/`ディレクトリが作成されている
- [ ] `db/migrations/sharding_2-mysql/`ディレクトリが作成されている
- [ ] `db/migrations/sharding_3-mysql/`ディレクトリが作成されている
- [ ] `db/migrations/sharding_4-mysql/`ディレクトリが作成されている
- [ ] Atlasコマンドでスキーマファイル（initial_schema.sql）がMySQL用に生成されている
- [ ] 手動移植が必要なファイル（seed_data.sql等）がMySQL構文に変換されている
- [ ] マイグレーションが正常に実行できる

### 6.3 テストコード
- [ ] `InitMySQLMasterSchema()`関数が実装されている
- [ ] `InitMySQLShardingSchema()`関数が実装されている
- [ ] `clearMySQLDatabaseTables()`関数が実装されている
- [ ] `SetupTestGroupManager()`でデータベースタイプを判定して適切な関数を呼び出している
- [ ] MySQL環境でテストが正常に実行できる

### 6.4 Atlas設定ファイル
- [ ] `config/develop/atlas.mysql.hcl`が作成されている
- [ ] `config/staging/atlas.mysql.hcl`が作成されている
- [ ] `config/production/atlas.mysql.hcl`が作成されている
- [ ] `config/test/atlas.mysql.hcl`が作成されている
- [ ] 各設定ファイルでMySQL用のURLが設定されている
- [ ] AtlasでMySQL用のマイグレーションが正常に実行できる

### 6.5 Docker Compose設定
- [ ] `docker-compose.mysql.yml`が作成されている
- [ ] マスターデータベースコンテナが定義されている（port 3306）
- [ ] シャーディングデータベースコンテナが定義されている（port 3307-3310）
- [ ] コンテナが正常に起動できる
- [ ] ヘルスチェックが正常に動作する

### 6.6 スクリプト
- [ ] `scripts/start-mysql.sh`が作成されている
- [ ] `scripts/migrate-test-mysql.sh`が作成されている
- [ ] `start-mysql.sh`でコンテナの起動・停止が正常に動作する
- [ ] `migrate-test-mysql.sh`でマイグレーションが正常に実行できる

### 6.7 DSN生成ロジック
- [ ] `GetDSN()`メソッドでMySQLのDSNに`charset=utf8mb4`が追加されている
- [ ] `GetDSN()`メソッドでMySQLのDSNに`loc=Local`が追加されている
- [ ] MySQL接続が正常に確立できる

### 6.8 動作確認
- [ ] ローカル環境でMySQLコンテナが正常に起動できる
- [ ] ローカル環境でMySQL用のマイグレーションが正常に実行できる
- [ ] ローカル環境でMySQL環境でテストが正常に実行できる
- [ ] PostgreSQLの既存機能が正常に動作することを確認
- [ ] CI環境でMySQL環境が正常に動作することを確認（該当する場合）

## 7. 影響範囲

### 7.1 変更が必要なファイル

#### 新規作成するファイル
- `config/develop/database.mysql.yaml`
- `config/staging/database.mysql.yaml`
- `config/production/database.mysql.yaml.example`
- `config/test/database.mysql.yaml`
- `config/develop/atlas.mysql.hcl`
- `config/staging/atlas.mysql.hcl`
- `config/production/atlas.mysql.hcl`
- `config/test/atlas.mysql.hcl`
- `docker-compose.mysql.yml`
- `scripts/start-mysql.sh`
- `scripts/migrate-test-mysql.sh`
- `db/migrations/master-mysql/`（ディレクトリとファイル）
- `db/migrations/sharding_1-mysql/`（ディレクトリとファイル）
- `db/migrations/sharding_2-mysql/`（ディレクトリとファイル）
- `db/migrations/sharding_3-mysql/`（ディレクトリとファイル）
- `db/migrations/sharding_4-mysql/`（ディレクトリとファイル）

#### 修正が必要なファイル
- `server/internal/config/config.go`: `GetDSN()`メソッドの改善
- `server/test/testutil/db.go`: MySQL用関数の追加、データベース判定機能の追加
- `config/{env}/config.yaml`: `DB_TYPE`フィールドの追加（各環境）

#### 確認が必要なファイル
- 既存のPostgreSQL設定ファイル: 正常に動作することを確認
- 既存のテストファイル: 正常に動作することを確認

### 7.2 既存機能への影響
- **PostgreSQL機能**: 既存のPostgreSQL機能に影響を与えない（追加のみ）
- **既存のテスト**: 既存のPostgreSQLテストが正常に動作することを確認

### 7.3 将来の拡張への影響
- **データベース選択の柔軟性**: 将来的に他のデータベース（SQLite等）にも対応可能な構造
- **設定ファイルの拡張**: 環境変数や設定ファイルの読み込み方法を拡張可能

## 8. 実装上の注意事項

### 8.1 設定ファイルの管理
- **分離方針**: PostgreSQL用とMySQL用の設定ファイルを分離（環境ごと）
- **デフォルト**: PostgreSQLをデフォルトとして維持
- **環境変数**: `DB_TYPE`環境変数でデータベースタイプを切り替え可能

### 8.2 マイグレーションファイルの管理
- **ディレクトリ分離**: PostgreSQL用とMySQL用のマイグレーションディレクトリを分離
- **Atlasによる自動生成**: スキーマファイル（initial_schema.sql）はAtlasコマンドでMySQL用に自動生成
- **手動移植**: データファイル（seed_data.sql）やビューファイルは手動でMySQL構文に変換して移植
- **構文変換**: PostgreSQL構文からMySQL構文への変換を正確に実施
- **バージョン管理**: マイグレーションファイルのバージョンを一致させる

### 8.3 テストコードの実装
- **関数分離**: PostgreSQL用とMySQL用の関数を分離
- **データベース判定**: 実行時にデータベースタイプを判定して適切な関数を呼び出す
- **エラーハンドリング**: データベース固有のエラーを適切に処理

### 8.4 DSN生成の実装
- **文字セット**: MySQLでは`utf8mb4`を明示的に指定
- **タイムゾーン**: MySQLでは`loc=Local`を指定
- **後方互換性**: 既存のPostgreSQL DSN生成に影響を与えない

### 8.5 Docker Composeの実装
- **ポート番号**: PostgreSQLと重複しないポート番号を使用（3306-3310）
- **ボリューム**: データ永続化のためのボリュームマウントを設定
- **ヘルスチェック**: コンテナの起動確認のためのヘルスチェックを設定

## 9. 参考情報

### 9.1 関連ドキュメント
- `docs/MySQL-Support-Analysis.md`: MySQL対応の分析結果
- `docs/Architecture.md`: アーキテクチャドキュメント
- `docs/Project-Structure.md`: プロジェクト構造ドキュメント

### 9.2 既存実装の参考
- `config/develop/database.yaml`: PostgreSQL用設定ファイル
- `config/develop/atlas.hcl`: PostgreSQL用Atlas設定
- `docker-compose.postgres.yml`: PostgreSQL用Docker Compose設定
- `scripts/start-postgres.sh`: PostgreSQL用起動スクリプト
- `scripts/migrate-test.sh`: PostgreSQL用マイグレーションスクリプト
- `server/test/testutil/db.go`: テストユーティリティ

### 9.3 技術スタック
- **言語**: Go
- **データベース**: PostgreSQL, MySQL
- **ORM**: GORM（`gorm.io/driver/postgres`, `gorm.io/driver/mysql`）
- **マイグレーション**: Atlas
- **コンテナ**: Docker, Docker Compose

### 9.4 主な構文の違い

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
