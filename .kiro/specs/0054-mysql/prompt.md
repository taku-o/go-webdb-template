MySQLでも動作するように修正する

PostgreSQLが主のデータベースだが、
MySQLでも動作するように修正する。

その場合、どのような作業が発生するか。
またどの設定ファイルをMySQL用に分けなければいけないか。

まずは修正する箇所、発生する作業を特定してください。


/kiro:spec-requirements "MySQL対応します。
cc-sddのfeature名は0054-mysqlとしてください。

規模が大きそうなので、
要件、設計、タスクは多めに。

事前分析
.kiro/specs/0054-mysql/MySQL-Support-Analysis.md

### 1. データベース接続設定（database.yaml）
MySQL用のdatabase.yamlを用意する。
- `config/develop/database.yaml` → `config/develop/database.mysql.yaml`
- `config/staging/database.yaml` → `config/staging/database.mysql.yaml`
- `config/production/database.yaml` → `config/production/database.mysql.yaml`
- `config/test/database.yaml` → `config/test/database.mysql.yaml`

### 2. マイグレーションファイル
MySQL用のマイグレーションディレクトリを用意。
- `db/migrations/master-mysql/` → MySQL用
- `db/migrations/sharding_1-mysql/` → MySQL用

これらのSQLはMySQL版を作って移植。
- db/migrations/master/20260108145415_seed_data.sql
- db/migrations/view_master/20260103030225_create_dm_news_view.sql

### 3. テストコード（testutil/db.go）
MySQL用の処理を用意。
- `InitMySQLMasterSchema()`関数
- `InitMySQLShardingSchema()`関数
- `clearMySQLDatabaseTables()`関数

### 4. Atlas設定ファイル（atlas.hcl）
MySQL用のatlas.hclを用意。
- `config/develop/atlas.hcl` → `config/develop/atlas.mysql.hcl`
- `config/staging/atlas.hcl` → `config/staging/atlas.mysql.hcl`
- `config/production/atlas.hcl` → `config/production/atlas.mysql.hcl`
- `config/test/atlas.hcl` → `config/test/atlas.mysql.hcl`

### 5. Docker Compose設定
MySQL用の起動スクリプトを用意する。
- docker-compose.mysql.yml

### 6. スクリプト（start-postgres.sh, migrate-test.sh）
MySQL用のスクリプトを用意する。
- script/start-mysql.sh
- scripts/migrate-test-mysql.sh

### 7. DSN生成ロジック（config.ShardConfig）
- ⚠️ 推奨: MySQLのDSNに`charset=utf8mb4&loc=Local`を追加

### 8. 環境情報の取得
- config/{env}/config.yaml にデータベースの種類を定義してよい。
- DB_TYPE=postgresql
- DB_TYPE=mysql

" think.


この5ファイルはatlasコマンドで生成する
>#### 3.2.3 移植が必要なファイル
>- `db/migrations/master/20260108145414_initial_schema.sql`
>- `db/migrations/sharding_1/20260108145537_initial_schema.sql`
>- `db/migrations/sharding_2/20260108145546_initial_schema.sql`
>- `db/migrations/sharding_3/20260108145548_initial_schema.sql`
>- `db/migrations/sharding_4/20260108145549_initial_schema.sql`

要件定義書を承認します。

/kiro:spec-design 0054-mysql
think.

これは可能？
"既存のHCLスキーマ（`db/schema/master.hcl`など）をそのまま使用"
可能ならOK。可能で無いならならmaster.mysql.hclとか分けてね。
> #### 3.2.2 Atlasコマンドによる自動生成の設計
> **スキーマファイル（initial_schema.sql）の生成**:
> - 既存のHCLスキーマ（`db/schema/master.hcl`など）をそのまま使用

設計書を承認します。

/kiro:spec-tasks 0054-mysql
タスクの数は通常よりも多めに用意してください。
think.

途中でコンテキストが尽きて作業が途切れると思われるので、
各タスクに完了チェック用のチェックボックスをつけてください。

次の作業のために会話をまとめてください。

タスクリストを承認します。

/sdd-fix-plan

_serena_indexing

/serena-initialize

/kiro:spec-impl 0053-parallel-dbtest






