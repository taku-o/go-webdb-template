# MySQL対応の実装タスク一覧

## 概要
PostgreSQLが主のデータベースだが、MySQLでも動作するように修正するためのタスク一覧。要件定義書と設計書に基づき、実装可能な粒度でタスクを分解。

## 実装フェーズ

### Phase 1: DSN生成ロジックの改善

#### タスク 1.1: GetDSN()メソッドの改善
**目的**: MySQL接続時の文字セットとタイムゾーンを適切に設定するため、`GetDSN()`メソッドを改善する。

**作業内容**:
- `server/internal/config/config.go`の`GetDSN()`メソッドを修正
- MySQLのDSNに`charset=utf8mb4`を追加
- MySQLのDSNに`loc=Local`を追加
- 既存の`parseTime=true`は維持

**実装内容**:
- 修正対象: `server/internal/config/config.go`の`GetDSN()`メソッド（362-364行目）
- 修正前:
  ```go
  case "mysql":
      return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
          s.User, s.Password, s.Host, s.Port, s.Name)
  ```
- 修正後:
  ```go
  case "mysql":
      return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local",
          s.User, s.Password, s.Host, s.Port, s.Name)
  ```

**受け入れ基準**:
- [ ] `GetDSN()`メソッドでMySQLのDSNに`charset=utf8mb4`が追加されている
- [ ] `GetDSN()`メソッドでMySQLのDSNに`loc=Local`が追加されている
- [ ] 既存の`parseTime=true`が維持されている
- [ ] PostgreSQLのDSN生成に影響がない（後方互換性を維持）

- _Requirements: 3.7.1, 6.7_
- _Design: 3.7.1_

---

### Phase 2: config.yamlへのDB_TYPE追加

#### タスク 2.1: develop環境のconfig.yamlにDB_TYPEを追加
**目的**: 開発環境でデータベースタイプを指定できるようにする。

**作業内容**:
- `config/develop/config.yaml`に`DB_TYPE`フィールドを追加
- デフォルト値は`postgresql`（既存の動作を維持）

**実装内容**:
- 追加場所: `config/develop/config.yaml`の先頭（`server:`セクションの前）
- 追加内容:
  ```yaml
  # データベースタイプの指定
  # postgresql: PostgreSQLを使用（デフォルト）
  # mysql: MySQLを使用
  DB_TYPE: postgresql
  ```

**受け入れ基準**:
- [ ] `config/develop/config.yaml`に`DB_TYPE`フィールドが追加されている
- [ ] デフォルト値が`postgresql`に設定されている

- _Requirements: 3.1.2, 6.1_
- _Design: 3.1.2_

---

#### タスク 2.2: staging環境のconfig.yamlにDB_TYPEを追加
**目的**: ステージング環境でデータベースタイプを指定できるようにする。

**作業内容**:
- `config/staging/config.yaml`に`DB_TYPE`フィールドを追加
- デフォルト値は`postgresql`

**実装内容**:
- 追加場所: `config/staging/config.yaml`の先頭
- 追加内容: タスク2.1と同様

**受け入れ基準**:
- [ ] `config/staging/config.yaml`に`DB_TYPE`フィールドが追加されている

- _Requirements: 3.1.2, 6.1_
- _Design: 3.1.2_

---

#### タスク 2.3: production環境のconfig.yaml.exampleにDB_TYPEを追加
**目的**: 本番環境でデータベースタイプを指定できるようにする。

**作業内容**:
- `config/production/config.yaml.example`に`DB_TYPE`フィールドを追加
- デフォルト値は`postgresql`

**実装内容**:
- 追加場所: `config/production/config.yaml.example`の先頭
- 追加内容: タスク2.1と同様

**受け入れ基準**:
- [ ] `config/production/config.yaml.example`に`DB_TYPE`フィールドが追加されている

- _Requirements: 3.1.2, 6.1_
- _Design: 3.1.2_

---

#### タスク 2.4: test環境のconfig.yamlにDB_TYPEを追加
**目的**: テスト環境でデータベースタイプを指定できるようにする。

**作業内容**:
- `config/test/config.yaml`に`DB_TYPE`フィールドを追加
- デフォルト値は`postgresql`

**実装内容**:
- 追加場所: `config/test/config.yaml`の先頭
- 追加内容: タスク2.1と同様

**受け入れ基準**:
- [ ] `config/test/config.yaml`に`DB_TYPE`フィールドが追加されている

- _Requirements: 3.1.2, 6.1_
- _Design: 3.1.2_

---

### Phase 3: 設定ファイル読み込みロジックの修正

#### タスク 3.1: Load()関数の修正（DB_TYPE判定機能の追加）
**目的**: `DB_TYPE`に応じて適切な`database.yaml`を読み込むように修正する。

**作業内容**:
- `server/internal/config/config.go`の`Load()`関数を修正
- `config.yaml`から`DB_TYPE`を読み込む
- `DB_TYPE`が`mysql`の場合、`database.mysql.yaml`を読み込む
- `DB_TYPE`が`postgresql`または未指定の場合、`database.yaml`を読み込む（既存の動作）

**実装内容**:
- 修正対象: `server/internal/config/config.go`の`Load()`関数（238-348行目）
- 修正箇所: データベース設定ファイルの読み込み部分（264-268行目）
- 修正前:
  ```go
  // データベース設定ファイルのマージ
  viper.SetConfigName("database")
  if err := viper.MergeInConfig(); err != nil {
      return nil, fmt.Errorf("failed to read database config file: %w", err)
  }
  ```
- 修正後:
  ```go
  // DB_TYPEを読み込む（config.yamlから）
  dbType := viper.GetString("DB_TYPE")
  if dbType == "" {
      dbType = "postgresql" // デフォルト値
  }

  // データベース設定ファイルのマージ
  var databaseFileName string
  if dbType == "mysql" {
      databaseFileName = "database.mysql"
  } else {
      databaseFileName = "database"
  }
  viper.SetConfigName(databaseFileName)
  if err := viper.MergeInConfig(); err != nil {
      return nil, fmt.Errorf("failed to read database config file: %w", err)
  }
  ```

**受け入れ基準**:
- [ ] `Load()`関数で`DB_TYPE`が正しく読み込まれる
- [ ] `DB_TYPE`が`mysql`の場合、`database.mysql.yaml`が読み込まれる
- [ ] `DB_TYPE`が`postgresql`または未指定の場合、`database.yaml`が読み込まれる（既存の動作）
- [ ] 既存のPostgreSQL設定ファイルの読み込みに影響がない

- _Requirements: 3.1.2, 6.1_
- _Design: 3.1.2_

---

### Phase 4: MySQL用設定ファイルの作成

#### タスク 4.1: develop環境のdatabase.mysql.yamlの作成
**目的**: 開発環境でMySQL接続設定を定義できるようにする。

**作業内容**:
- `config/develop/database.mysql.yaml`を作成
- PostgreSQL用の`config/develop/database.yaml`をベースに作成
- `driver: mysql`に変更
- ポート番号をMySQL用に変更（3306, 3307, 3308, 3309, 3310）

**実装内容**:
- 参考ファイル: `config/develop/database.yaml`
- 作成ファイル: `config/develop/database.mysql.yaml`
- 主な変更点:
  - `driver: postgres` → `driver: mysql`
  - ポート番号: `5432-5436` → `3306-3310`
  - その他の設定はPostgreSQLと同じ値を維持

**受け入れ基準**:
- [ ] `config/develop/database.mysql.yaml`が作成されている
- [ ] 各設定で`driver: mysql`が指定されている
- [ ] MySQL接続情報が正しく設定されている
- [ ] ポート番号がMySQL用に変更されている（3306, 3307, 3308, 3309, 3310）

- _Requirements: 3.1.1, 6.1_
- _Design: 3.1.1_

---

#### タスク 4.2: staging環境のdatabase.mysql.yamlの作成
**目的**: ステージング環境でMySQL接続設定を定義できるようにする。

**作業内容**:
- `config/staging/database.mysql.yaml`を作成
- PostgreSQL用の`config/staging/database.yaml`をベースに作成
- タスク4.1と同様の変更を適用

**実装内容**:
- 参考ファイル: `config/staging/database.yaml`
- 作成ファイル: `config/staging/database.mysql.yaml`
- 変更内容: タスク4.1と同様

**受け入れ基準**:
- [ ] `config/staging/database.mysql.yaml`が作成されている
- [ ] 各設定で`driver: mysql`が指定されている
- [ ] MySQL接続情報が正しく設定されている

- _Requirements: 3.1.1, 6.1_
- _Design: 3.1.1_

---

#### タスク 4.3: production環境のdatabase.mysql.yaml.exampleの作成
**目的**: 本番環境でMySQL接続設定を定義できるようにする。

**作業内容**:
- `config/production/database.mysql.yaml.example`を作成
- PostgreSQL用の`config/production/database.yaml.example`をベースに作成
- タスク4.1と同様の変更を適用

**実装内容**:
- 参考ファイル: `config/production/database.yaml.example`
- 作成ファイル: `config/production/database.mysql.yaml.example`
- 変更内容: タスク4.1と同様

**受け入れ基準**:
- [ ] `config/production/database.mysql.yaml.example`が作成されている
- [ ] 各設定で`driver: mysql`が指定されている
- [ ] MySQL接続情報が正しく設定されている

- _Requirements: 3.1.1, 6.1_
- _Design: 3.1.1_

---

#### タスク 4.4: test環境のdatabase.mysql.yamlの作成
**目的**: テスト環境でMySQL接続設定を定義できるようにする。

**作業内容**:
- `config/test/database.mysql.yaml`を作成
- PostgreSQL用の`config/test/database.yaml`をベースに作成
- タスク4.1と同様の変更を適用

**実装内容**:
- 参考ファイル: `config/test/database.yaml`
- 作成ファイル: `config/test/database.mysql.yaml`
- 変更内容: タスク4.1と同様

**受け入れ基準**:
- [ ] `config/test/database.mysql.yaml`が作成されている
- [ ] 各設定で`driver: mysql`が指定されている
- [ ] MySQL接続情報が正しく設定されている

- _Requirements: 3.1.1, 6.1_
- _Design: 3.1.1_

---

### Phase 5: MySQL用Docker Compose設定の作成

#### タスク 5.1: docker-compose.mysql.ymlの作成
**目的**: MySQLコンテナを起動できるようにする。

**作業内容**:
- `docker-compose.mysql.yml`を作成
- PostgreSQL用の`docker-compose.postgres.yml`をベースに作成
- MySQL 8のDockerイメージを使用
- ポート番号をMySQL用に設定（3306, 3307, 3308, 3309, 3310）

**実装内容**:
- 参考ファイル: `docker-compose.postgres.yml`
- 作成ファイル: `docker-compose.mysql.yml`
- 主な変更点:
  - `image: postgres:15-alpine` → `image: mysql:8`
  - ポート番号: `5432-5436` → `3306-3310`
  - 環境変数: `POSTGRES_*` → `MYSQL_*`
  - ボリュームマウント: `./postgres/data/` → `./mysql/data/`
  - ヘルスチェック: `pg_isready` → `mysqladmin ping`
  - コマンド: `--default-authentication-plugin=mysql_native_password`を追加

**受け入れ基準**:
- [ ] `docker-compose.mysql.yml`が作成されている
- [ ] マスターデータベースコンテナが定義されている（port 3306）
- [ ] シャーディングデータベースコンテナが定義されている（port 3307-3310）
- [ ] 各コンテナでMySQL 8のイメージが使用されている
- [ ] 環境変数がMySQL用に設定されている
- [ ] ボリュームマウントがMySQL用に設定されている
- [ ] ヘルスチェックがMySQL用に設定されている

- _Requirements: 3.5.1, 6.5_
- _Design: 3.5.1_

---

### Phase 6: MySQL用スクリプトの作成

#### タスク 6.1: start-mysql.shの作成
**目的**: MySQLコンテナを起動・停止できるようにする。

**作業内容**:
- `scripts/start-mysql.sh`を作成
- PostgreSQL用の`scripts/start-postgres.sh`をベースに作成
- `docker-compose.mysql.yml`を使用
- 接続URLをMySQL用に変更

**実装内容**:
- 参考ファイル: `scripts/start-postgres.sh`
- 作成ファイル: `scripts/start-mysql.sh`
- 主な変更点:
  - `COMPOSE_FILE`を`docker-compose.mysql.yml`に変更
  - 接続URLをMySQL用に変更（`mysql://webdb:webdb@tcp(localhost:3306)/webdb_master`など）
  - メッセージをMySQL用に変更

**受け入れ基準**:
- [ ] `scripts/start-mysql.sh`が作成されている
- [ ] 実行権限が設定されている（`chmod +x scripts/start-mysql.sh`）
- [ ] `start`コマンドでMySQLコンテナが起動できる
- [ ] `stop`コマンドでMySQLコンテナが停止できる
- [ ] `status`コマンドでコンテナの状態が表示できる
- [ ] `health`コマンドでヘルスチェック状態が表示できる

- _Requirements: 3.6.1, 6.6_
- _Design: 3.6.1_

---

#### タスク 6.2: migrate-test-mysql.shの作成
**目的**: MySQL用のマイグレーションを実行できるようにする。

**作業内容**:
- `scripts/migrate-test-mysql.sh`を作成
- PostgreSQL用の`scripts/migrate-test.sh`をベースに作成
- MySQL接続情報を使用
- MySQL URL形式を構築
- `docker exec`でMySQLコンテナに接続してSQLファイルを適用

**実装内容**:
- 参考ファイル: `scripts/migrate-test.sh`
- 作成ファイル: `scripts/migrate-test-mysql.sh`
- 主な変更点:
  - ポート番号をMySQL用に変更（3306-3310）
  - `build_postgres_url()` → `build_mysql_url()`に変更
  - MySQL URL形式: `user:pass@tcp(host:port)/dbname?charset=utf8mb4&parseTime=true&loc=Local`
  - `docker exec`コマンド: `psql` → `mysql`
  - マイグレーションディレクトリ: `master` → `master-mysql`など

**受け入れ基準**:
- [ ] `scripts/migrate-test-mysql.sh`が作成されている
- [ ] 実行権限が設定されている（`chmod +x scripts/migrate-test-mysql.sh`）
- [ ] マスターデータベースへのマイグレーションが正常に実行できる
- [ ] シャーディングデータベースへのマイグレーションが正常に実行できる
- [ ] 各データベースへのSQLファイルが正常に適用できる

- _Requirements: 3.6.2, 6.6_
- _Design: 3.6.2_

---

### Phase 7: Atlas設定ファイルの作成

#### タスク 7.1: develop環境のatlas.mysql.hclの作成
**目的**: 開発環境でAtlasでMySQL用のマイグレーションを管理できるようにする。

**作業内容**:
- `config/develop/atlas.mysql.hcl`を作成
- PostgreSQL用の`config/develop/atlas.hcl`をベースに作成
- MySQL用のURLとdev環境を設定
- マイグレーションディレクトリをMySQL用に設定

**実装内容**:
- 参考ファイル: `config/develop/atlas.hcl`
- 作成ファイル: `config/develop/atlas.mysql.hcl`
- 主な変更点:
  - `url = "postgres://..."` → `url = "webdb:webdb@tcp(localhost:3306)/webdb_master?charset=utf8mb4&parseTime=true&loc=Local"`
  - `dev = "docker://postgres/15/dev?search_path=public"` → `dev = "docker://mysql/8/dev"`
  - `dir = "file://db/migrations/master"` → `dir = "file://db/migrations/master-mysql"`

**受け入れ基準**:
- [ ] `config/develop/atlas.mysql.hcl`が作成されている
- [ ] 各環境でMySQL用のURLが設定されている
- [ ] 各環境で`dev = "docker://mysql/8/dev"`が設定されている
- [ ] マイグレーションディレクトリがMySQL用に設定されている

- _Requirements: 3.4.1, 6.4_
- _Design: 3.4.1_

---

#### タスク 7.2: staging環境のatlas.mysql.hclの作成
**目的**: ステージング環境でAtlasでMySQL用のマイグレーションを管理できるようにする。

**作業内容**:
- `config/staging/atlas.mysql.hcl`を作成
- タスク7.1と同様の内容で作成

**実装内容**:
- 参考ファイル: `config/staging/atlas.hcl`
- 作成ファイル: `config/staging/atlas.mysql.hcl`
- 変更内容: タスク7.1と同様（接続情報はstaging環境用に調整）

**受け入れ基準**:
- [ ] `config/staging/atlas.mysql.hcl`が作成されている
- [ ] 各環境でMySQL用のURLが設定されている

- _Requirements: 3.4.1, 6.4_
- _Design: 3.4.1_

---

#### タスク 7.3: production環境のatlas.mysql.hclの作成
**目的**: 本番環境でAtlasでMySQL用のマイグレーションを管理できるようにする。

**作業内容**:
- `config/production/atlas.mysql.hcl`を作成
- タスク7.1と同様の内容で作成

**実装内容**:
- 参考ファイル: `config/production/atlas.hcl`
- 作成ファイル: `config/production/atlas.mysql.hcl`
- 変更内容: タスク7.1と同様（接続情報はproduction環境用に調整）

**受け入れ基準**:
- [ ] `config/production/atlas.mysql.hcl`が作成されている
- [ ] 各環境でMySQL用のURLが設定されている

- _Requirements: 3.4.1, 6.4_
- _Design: 3.4.1_

---

#### タスク 7.4: test環境のatlas.mysql.hclの作成
**目的**: テスト環境でAtlasでMySQL用のマイグレーションを管理できるようにする。

**作業内容**:
- `config/test/atlas.mysql.hcl`を作成
- タスク7.1と同様の内容で作成

**実装内容**:
- 参考ファイル: `config/test/atlas.hcl`
- 作成ファイル: `config/test/atlas.mysql.hcl`
- 変更内容: タスク7.1と同様（接続情報はtest環境用に調整）

**受け入れ基準**:
- [ ] `config/test/atlas.mysql.hcl`が作成されている
- [ ] 各環境でMySQL用のURLが設定されている

- _Requirements: 3.4.1, 6.4_
- _Design: 3.4.1_

---

### Phase 8: MySQL用マイグレーションディレクトリの作成

#### タスク 8.1: master-mysqlディレクトリの作成
**目的**: マスターデータベース用のMySQLマイグレーションファイルを管理するディレクトリを作成する。

**作業内容**:
- `db/migrations/master-mysql/`ディレクトリを作成

**実装内容**:
- ディレクトリ作成: `mkdir -p db/migrations/master-mysql`

**受け入れ基準**:
- [ ] `db/migrations/master-mysql/`ディレクトリが作成されている

- _Requirements: 3.2.1, 6.2_
- _Design: 3.2.1_

---

#### タスク 8.2: sharding_1-mysqlディレクトリの作成
**目的**: シャーディング1用のMySQLマイグレーションファイルを管理するディレクトリを作成する。

**作業内容**:
- `db/migrations/sharding_1-mysql/`ディレクトリを作成

**実装内容**:
- ディレクトリ作成: `mkdir -p db/migrations/sharding_1-mysql`

**受け入れ基準**:
- [ ] `db/migrations/sharding_1-mysql/`ディレクトリが作成されている

- _Requirements: 3.2.1, 6.2_
- _Design: 3.2.1_

---

#### タスク 8.3: sharding_2-mysqlディレクトリの作成
**目的**: シャーディング2用のMySQLマイグレーションファイルを管理するディレクトリを作成する。

**作業内容**:
- `db/migrations/sharding_2-mysql/`ディレクトリを作成

**実装内容**:
- ディレクトリ作成: `mkdir -p db/migrations/sharding_2-mysql`

**受け入れ基準**:
- [ ] `db/migrations/sharding_2-mysql/`ディレクトリが作成されている

- _Requirements: 3.2.1, 6.2_
- _Design: 3.2.1_

---

#### タスク 8.4: sharding_3-mysqlディレクトリの作成
**目的**: シャーディング3用のMySQLマイグレーションファイルを管理するディレクトリを作成する。

**作業内容**:
- `db/migrations/sharding_3-mysql/`ディレクトリを作成

**実装内容**:
- ディレクトリ作成: `mkdir -p db/migrations/sharding_3-mysql`

**受け入れ基準**:
- [ ] `db/migrations/sharding_3-mysql/`ディレクトリが作成されている

- _Requirements: 3.2.1, 6.2_
- _Design: 3.2.1_

---

#### タスク 8.5: sharding_4-mysqlディレクトリの作成
**目的**: シャーディング4用のMySQLマイグレーションファイルを管理するディレクトリを作成する。

**作業内容**:
- `db/migrations/sharding_4-mysql/`ディレクトリを作成

**実装内容**:
- ディレクトリ作成: `mkdir -p db/migrations/sharding_4-mysql`

**受け入れ基準**:
- [ ] `db/migrations/sharding_4-mysql/`ディレクトリが作成されている

- _Requirements: 3.2.1, 6.2_
- _Design: 3.2.1_

---

#### タスク 8.6: view_master-mysqlディレクトリの作成
**目的**: ビューマスター用のMySQLマイグレーションファイルを管理するディレクトリを作成する。

**作業内容**:
- `db/migrations/view_master-mysql/`ディレクトリを作成

**実装内容**:
- ディレクトリ作成: `mkdir -p db/migrations/view_master-mysql`

**受け入れ基準**:
- [ ] `db/migrations/view_master-mysql/`ディレクトリが作成されている

- _Requirements: 3.2.1, 6.2_
- _Design: 3.2.1_

---

### Phase 9: Atlasコマンドによるマイグレーションファイルの生成

#### タスク 9.1: 既存HCLスキーマでMySQL用SQL生成の確認
**目的**: 既存のHCLスキーマ（`db/schema/master.hcl`など）をそのまま使用して、MySQL用のSQLが生成できるか確認する。

**作業内容**:
- `config/develop/atlas.mysql.hcl`を使用してAtlasコマンドを実行
- 既存のHCLスキーマ（`db/schema/master.hcl`）を指定
- MySQL用のSQLが生成できるか確認
- 生成されたSQLに問題がないか確認

**実装内容**:
- 実行コマンド:
  ```bash
  atlas migrate diff \
    --env master \
    --config config/develop/atlas.mysql.hcl \
    --to file://db/schema/master.hcl \
    --dir file://db/migrations/master-mysql
  ```
- 確認内容:
  - SQLが正常に生成されるか
  - `schema.public`の参照が適切に処理されるか
  - `type = serial`が`INT AUTO_INCREMENT`に変換されるか
  - 生成されたSQLに構文エラーがないか

**受け入れ基準**:
- [ ] Atlasコマンドが正常に実行できる
- [ ] MySQL用のSQLが生成される
- [ ] 生成されたSQLに構文エラーがない
- [ ] 問題が発生した場合は、MySQL用HCLスキーマの作成が必要であることを記録

- _Requirements: 3.2.2, 6.2_
- _Design: 3.2.2, 7.2_

---

#### タスク 9.2: マスターデータベース用マイグレーションファイルの生成
**目的**: マスターデータベース用のMySQLマイグレーションファイルを生成する。

**作業内容**:
- タスク9.1で問題がなかった場合: 既存のHCLスキーマを使用
- タスク9.1で問題があった場合: MySQL用HCLスキーマを作成してから生成
- AtlasコマンドでMySQL用のSQLを生成

**実装内容**:
- 既存HCLスキーマを使用する場合:
  ```bash
  atlas migrate diff \
    --env master \
    --config config/develop/atlas.mysql.hcl \
    --to file://db/schema/master.hcl \
    --dir file://db/migrations/master-mysql
  ```
- MySQL用HCLスキーマを作成する場合:
  - `db/schema/master-mysql.hcl`を作成（`schema.public`の参照を削除など）
  - 上記コマンドの`--to`を`file://db/schema/master-mysql.hcl`に変更

**受け入れ基準**:
- [ ] `db/migrations/master-mysql/`ディレクトリにMySQL用のSQLファイルが生成されている
- [ ] 生成されたSQLファイルがMySQL構文になっている
- [ ] 生成されたSQLファイルに構文エラーがない

- _Requirements: 3.2.2, 3.2.3, 6.2_
- _Design: 3.2.2, 8.5_

---

#### タスク 9.3: シャーディング1用マイグレーションファイルの生成
**目的**: シャーディング1用のMySQLマイグレーションファイルを生成する。

**作業内容**:
- タスク9.1と同様の方針で、シャーディング1用のマイグレーションファイルを生成

**実装内容**:
- 既存HCLスキーマを使用する場合:
  ```bash
  atlas migrate diff \
    --env sharding_1 \
    --config config/develop/atlas.mysql.hcl \
    --to file://db/schema/sharding_1 \
    --dir file://db/migrations/sharding_1-mysql
  ```
- MySQL用HCLスキーマを作成する場合:
  - `db/schema/sharding_1-mysql/`ディレクトリを作成
  - 各HCLファイルをMySQL用に修正

**受け入れ基準**:
- [ ] `db/migrations/sharding_1-mysql/`ディレクトリにMySQL用のSQLファイルが生成されている
- [ ] 生成されたSQLファイルがMySQL構文になっている

- _Requirements: 3.2.2, 3.2.3, 6.2_
- _Design: 3.2.2, 8.5_

---

#### タスク 9.4: シャーディング2用マイグレーションファイルの生成
**目的**: シャーディング2用のMySQLマイグレーションファイルを生成する。

**作業内容**:
- タスク9.3と同様の方針で、シャーディング2用のマイグレーションファイルを生成

**実装内容**:
- タスク9.3と同様（`sharding_2`に変更）

**受け入れ基準**:
- [ ] `db/migrations/sharding_2-mysql/`ディレクトリにMySQL用のSQLファイルが生成されている
- [ ] 生成されたSQLファイルがMySQL構文になっている

- _Requirements: 3.2.2, 3.2.3, 6.2_
- _Design: 3.2.2, 8.5_

---

#### タスク 9.5: シャーディング3用マイグレーションファイルの生成
**目的**: シャーディング3用のMySQLマイグレーションファイルを生成する。

**作業内容**:
- タスク9.3と同様の方針で、シャーディング3用のマイグレーションファイルを生成

**実装内容**:
- タスク9.3と同様（`sharding_3`に変更）

**受け入れ基準**:
- [ ] `db/migrations/sharding_3-mysql/`ディレクトリにMySQL用のSQLファイルが生成されている
- [ ] 生成されたSQLファイルがMySQL構文になっている

- _Requirements: 3.2.2, 3.2.3, 6.2_
- _Design: 3.2.2, 8.5_

---

#### タスク 9.6: シャーディング4用マイグレーションファイルの生成
**目的**: シャーディング4用のMySQLマイグレーションファイルを生成する。

**作業内容**:
- タスク9.3と同様の方針で、シャーディング4用のマイグレーションファイルを生成

**実装内容**:
- タスク9.3と同様（`sharding_4`に変更）

**受け入れ基準**:
- [ ] `db/migrations/sharding_4-mysql/`ディレクトリにMySQL用のSQLファイルが生成されている
- [ ] 生成されたSQLファイルがMySQL構文になっている

- _Requirements: 3.2.2, 3.2.3, 6.2_
- _Design: 3.2.2, 8.5_

---

### Phase 10: 手動移植が必要なファイルの移植

#### タスク 10.1: seed_data.sqlのMySQL版作成
**目的**: マスターデータベース用のseed_data.sqlをMySQL構文に変換して移植する。

**作業内容**:
- `db/migrations/master/20260108145415_seed_data.sql`をMySQL構文に変換
- `db/migrations/master-mysql/20260108145415_seed_data.sql`を作成

**実装内容**:
- 参考ファイル: `db/migrations/master/20260108145415_seed_data.sql`
- 作成ファイル: `db/migrations/master-mysql/20260108145415_seed_data.sql`
- 主な変換内容:
  - `ON CONFLICT DO NOTHING` → `INSERT IGNORE`
  - ダブルクォート `"` → バッククォート `` ` ``
  - その他のPostgreSQL固有の構文をMySQL構文に変換

**変換例**:
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

**受け入れ基準**:
- [ ] `db/migrations/master-mysql/20260108145415_seed_data.sql`が作成されている
- [ ] PostgreSQL固有の構文がMySQL構文に変換されている
- [ ] 生成されたSQLファイルがMySQLで正常に実行できる

- _Requirements: 3.2.3, 6.2_
- _Design: 3.2.3_

---

#### タスク 10.2: view_masterのSQLのMySQL版作成
**目的**: ビューマスター用のSQLをMySQL構文に変換して移植する。

**作業内容**:
- `db/migrations/view_master/20260103030225_create_dm_news_view.sql`をMySQL構文に変換
- `db/migrations/view_master-mysql/20260103030225_create_dm_news_view.sql`を作成

**実装内容**:
- 参考ファイル: `db/migrations/view_master/20260103030225_create_dm_news_view.sql`
- 作成ファイル: `db/migrations/view_master-mysql/20260103030225_create_dm_news_view.sql`
- 主な変換内容:
  - PostgreSQL固有のビュー構文をMySQL構文に変換
  - 必要に応じてビューの定義を調整

**注意事項**:
- ビューの定義がPostgreSQLとMySQLで大きく異なる場合は、設計を確認する

**受け入れ基準**:
- [ ] `db/migrations/view_master-mysql/20260103030225_create_dm_news_view.sql`が作成されている（必要に応じて）
- [ ] PostgreSQL固有の構文がMySQL構文に変換されている
- [ ] 生成されたSQLファイルがMySQLで正常に実行できる

- _Requirements: 3.2.3, 6.2_
- _Design: 3.2.3_

---

### Phase 11: MySQL用テスト関数の実装

#### タスク 11.1: InitMySQLMasterSchema()関数の実装
**目的**: テスト環境でMySQLデータベースのマスターデータベーススキーマを初期化する関数を実装する。

**作業内容**:
- `server/test/testutil/db.go`に`InitMySQLMasterSchema()`関数を追加
- PostgreSQL用の`InitMasterSchema()`を参考に実装
- MySQL用のSQL構文を使用

**実装内容**:
- 関数名: `InitMySQLMasterSchema(t *testing.T, database *gorm.DB)`
- 実装内容:
  ```go
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
- 主な変換:
  - `SERIAL PRIMARY KEY` → `INT AUTO_INCREMENT PRIMARY KEY`
  - その他の構文は両方で動作する

**受け入れ基準**:
- [ ] `InitMySQLMasterSchema()`関数が実装されている
- [ ] MySQL用のSQL構文が使用されている
- [ ] 関数が正常に動作する

- _Requirements: 3.3.1, 6.3_
- _Design: 3.3.1_

---

#### タスク 11.2: InitMySQLShardingSchema()関数の実装
**目的**: テスト環境でMySQLデータベースのシャーディングデータベーススキーマを初期化する関数を実装する。

**作業内容**:
- `server/test/testutil/db.go`に`InitMySQLShardingSchema()`関数を追加
- PostgreSQL用の`InitShardingSchema()`を参考に実装
- MySQL用のSQL構文を使用

**実装内容**:
- 関数名: `InitMySQLShardingSchema(t *testing.T, database *gorm.DB, startTable, endTable int)`
- 実装内容:
  ```go
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
- 主な変換:
  - `TEXT PRIMARY KEY` → `VARCHAR(32) PRIMARY KEY`（より適切）

**受け入れ基準**:
- [ ] `InitMySQLShardingSchema()`関数が実装されている
- [ ] MySQL用のSQL構文が使用されている
- [ ] 関数が正常に動作する

- _Requirements: 3.3.1, 6.3_
- _Design: 3.3.1_

---

#### タスク 11.3: clearMySQLDatabaseTables()関数の実装
**目的**: テスト環境でMySQLデータベースの全テーブルのデータをクリアする関数を実装する。

**作業内容**:
- `server/test/testutil/db.go`に`clearMySQLDatabaseTables()`関数を追加
- PostgreSQL用の`clearDatabaseTables()`を参考に実装
- MySQL用のSQL構文を使用

**実装内容**:
- 関数名: `clearMySQLDatabaseTables(t *testing.T, database *gorm.DB)`
- 実装内容:
  ```go
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
- 主な変換:
  - `pg_tables` → `INFORMATION_SCHEMA.TABLES`
  - `TRUNCATE TABLE ... RESTART IDENTITY CASCADE` → `TRUNCATE TABLE`（AUTO_INCREMENTは自動リセット）

**受け入れ基準**:
- [ ] `clearMySQLDatabaseTables()`関数が実装されている
- [ ] MySQL用のSQL構文が使用されている
- [ ] 関数が正常に動作する
- [ ] テーブルのデータが正常にクリアされる

- _Requirements: 3.3.2, 6.3_
- _Design: 3.3.2_

---

### Phase 12: データベース判定機能の追加

#### タスク 12.1: SetupTestGroupManager()の修正（データベース判定機能の追加）
**目的**: `SetupTestGroupManager()`関数でデータベースタイプを判定して、適切なスキーマ初期化関数を呼び出すように修正する。

**作業内容**:
- `server/test/testutil/db.go`の`SetupTestGroupManager()`関数を修正
- データベースドライバーを判定
- `driver == "postgres"`の場合、既存の関数を使用
- `driver == "mysql"`の場合、MySQL用の関数を使用

**実装内容**:
- 修正対象: `server/test/testutil/db.go`の`SetupTestGroupManager()`関数
- 修正箇所: スキーマ初期化部分（104-124行目）
- 修正内容:
  ```go
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
  ```

**受け入れ基準**:
- [ ] `SetupTestGroupManager()`関数でデータベースドライバーが判定されている
- [ ] `driver == "postgres"`の場合、既存の関数が呼び出される
- [ ] `driver == "mysql"`の場合、MySQL用の関数が呼び出される
- [ ] 既存のPostgreSQLテストが正常に動作する

- _Requirements: 3.3.3, 6.3_
- _Design: 3.3.3_

---

#### タスク 12.2: ClearTestDatabase()の修正（データベース判定機能の追加）
**目的**: `ClearTestDatabase()`関数でデータベースタイプを判定して、適切なデータクリア関数を呼び出すように修正する。

**作業内容**:
- `server/test/testutil/db.go`の`ClearTestDatabase()`関数を修正
- データベースドライバーを判定
- `driver == "postgres"`の場合、既存の関数を使用
- `driver == "mysql"`の場合、MySQL用の関数を使用

**実装内容**:
- 修正対象: `server/test/testutil/db.go`の`ClearTestDatabase()`関数（265-276行目）
- 修正内容:
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

**受け入れ基準**:
- [ ] `ClearTestDatabase()`関数でデータベースドライバーが判定されている
- [ ] `driver == "postgres"`の場合、既存の関数が呼び出される
- [ ] `driver == "mysql"`の場合、MySQL用の関数が呼び出される
- [ ] 既存のPostgreSQLテストが正常に動作する

- _Requirements: 3.3.3, 6.3_
- _Design: 3.3.3_

---

### Phase 13: 動作確認

#### タスク 13.1: MySQLコンテナの起動確認
**目的**: MySQLコンテナが正常に起動できることを確認する。

**作業内容**:
- `./scripts/start-mysql.sh start`を実行
- コンテナが正常に起動することを確認
- ヘルスチェックが正常に動作することを確認

**実装内容**:
- 実行コマンド: `./scripts/start-mysql.sh start`
- 確認コマンド: `./scripts/start-mysql.sh status`
- 確認コマンド: `./scripts/start-mysql.sh health`

**受け入れ基準**:
- [ ] MySQLコンテナが正常に起動できる
- [ ] ヘルスチェックが正常に動作する
- [ ] コンテナの状態が正常に表示される

- _Requirements: 6.5, 6.8_
- _Design: 3.5.1_

---

#### タスク 13.2: MySQL接続の確認
**目的**: MySQL接続が正常に確立できることを確認する。

**作業内容**:
- `DB_TYPE=mysql`を設定してアプリケーションを起動
- MySQL接続が正常に確立されることを確認
- 接続エラーが発生しないことを確認

**実装内容**:
- 環境変数設定: `export DB_TYPE=mysql`または`config/develop/config.yaml`で`DB_TYPE: mysql`に設定
- アプリケーション起動: `cd server && go run cmd/server/main.go`
- 接続ログを確認

**受け入れ基準**:
- [ ] MySQL接続が正常に確立できる
- [ ] 接続エラーが発生しない
- [ ] DSNが正しく生成されている（`charset=utf8mb4&parseTime=true&loc=Local`が含まれている）

- _Requirements: 6.7, 6.8_
- _Design: 3.7.1_

---

#### タスク 13.3: MySQL用マイグレーションの実行確認
**目的**: MySQL用のマイグレーションが正常に実行できることを確認する。

**作業内容**:
- `./scripts/migrate-test-mysql.sh`を実行
- マイグレーションが正常に実行されることを確認
- エラーが発生しないことを確認

**実装内容**:
- 実行コマンド: `./scripts/migrate-test-mysql.sh`
- 確認内容:
  - マスターデータベースへのマイグレーションが正常に実行される
  - シャーディングデータベースへのマイグレーションが正常に実行される
  - エラーが発生しない

**受け入れ基準**:
- [ ] マイグレーションが正常に実行できる
- [ ] エラーが発生しない
- [ ] テーブルが正常に作成される

- _Requirements: 6.2, 6.8_
- _Design: 3.6.2_

---

#### タスク 13.4: MySQL環境でのテスト実行確認
**目的**: MySQL環境でテストが正常に実行できることを確認する。

**作業内容**:
- `config/test/config.yaml`で`DB_TYPE: mysql`に設定
- `go test ./test/integration/...`を実行
- テストが正常に実行されることを確認

**実装内容**:
- 設定変更: `config/test/config.yaml`で`DB_TYPE: mysql`に設定
- テスト実行: `cd server && go test ./test/integration/...`
- 確認内容:
  - テストが正常に実行される
  - エラーが発生しない
  - データベーススキーマが正常に初期化される

**受け入れ基準**:
- [ ] MySQL環境でテストが正常に実行できる
- [ ] エラーが発生しない
- [ ] データベーススキーマが正常に初期化される

- _Requirements: 6.3, 6.8_
- _Design: 3.3.3_

---

#### タスク 13.5: PostgreSQL環境でのテスト実行確認（既存動作確認）
**目的**: PostgreSQL環境でテストが正常に実行できることを確認する（既存の動作確認）。

**作業内容**:
- `config/test/config.yaml`で`DB_TYPE: postgresql`に設定（または未設定）
- `go test ./test/integration/...`を実行
- テストが正常に実行されることを確認

**実装内容**:
- 設定確認: `config/test/config.yaml`で`DB_TYPE: postgresql`に設定（または未設定）
- テスト実行: `cd server && go test ./test/integration/...`
- 確認内容:
  - テストが正常に実行される
  - エラーが発生しない
  - 既存のPostgreSQLテストが正常に動作する

**受け入れ基準**:
- [ ] PostgreSQL環境でテストが正常に実行できる
- [ ] エラーが発生しない
- [ ] 既存のPostgreSQLテストが正常に動作する

- _Requirements: 6.8, 7.2_
- _Design: 3.3.3_

---

### Phase 14: ドキュメントの更新

#### タスク 14.1: READMEの更新（MySQL対応の手順追加）
**目的**: READMEにMySQL対応の手順を追加する。

**作業内容**:
- `README.md`にMySQL対応の手順を追加
- MySQL用のコマンドを追加
- MySQL用の設定方法を追加

**実装内容**:
- 追加場所: `README.md`の適切なセクション
- 追加内容:
  - MySQLコンテナの起動方法
  - MySQL用のマイグレーション実行方法
  - `DB_TYPE`の設定方法
  - MySQL用の設定ファイルの説明

**受け入れ基準**:
- [ ] `README.md`にMySQL対応の手順が追加されている
- [ ] MySQL用のコマンドが記載されている
- [ ] MySQL用の設定方法が記載されている

- _Requirements: 6.8_
- _Design: 6.4_

---

## 実装順序の推奨

1. **Phase 1**: DSN生成ロジックの改善
2. **Phase 2**: config.yamlへのDB_TYPE追加（各環境）
3. **Phase 3**: 設定ファイル読み込みロジックの修正
4. **Phase 4**: MySQL用設定ファイルの作成（各環境）
5. **Phase 5**: MySQL用Docker Compose設定の作成
6. **Phase 6**: MySQL用スクリプトの作成
7. **Phase 7**: Atlas設定ファイルの作成（各環境）
8. **Phase 8**: MySQL用マイグレーションディレクトリの作成
9. **Phase 9**: Atlasコマンドによるマイグレーションファイルの生成
10. **Phase 10**: 手動移植が必要なファイルの移植
11. **Phase 11**: MySQL用テスト関数の実装
12. **Phase 12**: データベース判定機能の追加
13. **Phase 13**: 動作確認
14. **Phase 14**: ドキュメントの更新

## 注意事項

### 実装時の注意点

1. **後方互換性**: 既存のPostgreSQL機能に影響を与えない（追加のみ）
2. **設定ファイルの分離**: PostgreSQL用とMySQL用の設定ファイルを分離（環境ごと）
3. **デフォルト値**: PostgreSQLをデフォルトとして維持
4. **構文変換**: PostgreSQL構文からMySQL構文への変換を正確に実施
5. **テストの確認**: 各フェーズで既存のPostgreSQLテストが正常に動作することを確認

### テスト時の注意点

1. **データベースタイプの切り替え**: `DB_TYPE`環境変数または`config.yaml`で切り替え
2. **コンテナの起動**: MySQLコンテナとPostgreSQLコンテナは同時に起動可能
3. **マイグレーションの確認**: マイグレーションが正常に実行できることを確認
4. **テストの実行**: PostgreSQLとMySQLの両方でテストが正常に実行できることを確認

## 参考情報

- 要件定義書: `requirements.md`
- 設計書: `design.md`
- 既存のPostgreSQL設定ファイル: `config/develop/database.yaml`
- 既存のPostgreSQLAtlas設定: `config/develop/atlas.hcl`
- 既存のPostgreSQLDocker Compose: `docker-compose.postgres.yml`
- 既存のPostgreSQL起動スクリプト: `scripts/start-postgres.sh`
- 既存のPostgreSQLマイグレーションスクリプト: `scripts/migrate-test.sh`
- 既存のテストユーティリティ: `server/test/testutil/db.go`
- DSN生成ロジック: `server/internal/config/config.go`
