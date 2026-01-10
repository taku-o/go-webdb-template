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

/kiro:spec-impl 0054-mysql 1

タスクリストの完了した受け入れ基準にチェックをつけてください。

/kiro:spec-impl 0054-mysql 2
/kiro:spec-impl 0054-mysql 3

/compact

/kiro:spec-impl 0054-mysql 4
/kiro:spec-impl 0054-mysql 5
/kiro:spec-impl 0054-mysql 6

>  注意: タスク6.2の受け入れ基準のうち、マイグレーション実行確認（残り3項目）はPhase 8-9でMySQL用マイグレーションディレクトリとファイルが作成された後に確認できます。

/compact

/kiro:spec-impl 0054-mysql 7
/kiro:spec-impl 0054-mysql 8
/kiro:spec-impl 0054-mysql 9

タスク9 OK。
/kiro:spec-impl 0054-mysql 10

これらのファイルのMySQL版が作られていない。
db/migrations/master/20260108145415_seed_data.sql
db/migrations/view_master/20260103030225_create_dm_news_view.sql

タスクリストの完了した受け入れ基準にチェックをつけてください。
タスク6.2のチェックも。

MySQLにSQLを流し込むのはまだやっていない？

はい、MySQLを起動して、SQLを実行してください。

MySQLに接続してみたい。
ホストと、DB名、ユーザー、パスワードの情報ください

  │ ホスト     │ localhost         │
  │ ポート     │ 3306              │
  │ DB名       │ webdb_master_test │
  │ ユーザー   │ webdb             │
  │ パスワード │ webdb             │

  Sharding Databases
  │  シャード  │ ポート │         DB名          │
  │ Sharding 1 │ 3307   │ webdb_sharding_1_test │
  │ Sharding 2 │ 3308   │ webdb_sharding_2_test │
  │ Sharding 3 │ 3309   │ webdb_sharding_3_test │
  │ Sharding 4 │ 3310   │ webdb_sharding_4_test │

docker-compose.cloudbeaver.yml でMySQLに接続したい。
けど、おそらくネットワーク設定がいる？

  ┌────────────┬───────────────────┬───────────────────────┬───────────────────────┬───────────────────────┬───────────────────────┐
  │    項目    │      Master       │      Sharding 1       │      Sharding 2       │      Sharding 3       │      Sharding 4       │
  ├────────────┼───────────────────┼───────────────────────┼───────────────────────┼───────────────────────┼───────────────────────┤
  │ ホスト     │ mysql-master      │ mysql-sharding-1      │ mysql-sharding-2      │ mysql-sharding-3      │ mysql-sharding-4      │
  ├────────────┼───────────────────┼───────────────────────┼───────────────────────┼───────────────────────┼───────────────────────┤
  │ ポート     │ 3306              │ 3306                  │ 3306                  │ 3306                  │ 3306                  │
  ├────────────┼───────────────────┼───────────────────────┼───────────────────────┼───────────────────────┼───────────────────────┤
  │ DB名       │ webdb_master_test │ webdb_sharding_1_test │ webdb_sharding_2_test │ webdb_sharding_3_test │ webdb_sharding_4_test │
  ├────────────┼───────────────────┼───────────────────────┼───────────────────────┼───────────────────────┼───────────────────────┤
  │ ユーザー   │ webdb             │ webdb                 │ webdb                 │ webdb                 │ webdb                 │
  ├────────────┼───────────────────┼───────────────────────┼───────────────────────┼───────────────────────┼───────────────────────┤
  │ パスワード │ webdb             │ webdb                 │ webdb                 │ webdb                 │ webdb                 │
  └────────────┴───────────────────┴───────────────────────┴───────────────────────┴───────────────────────┴───────────────────────┘

DB接続成功した。
/kiro:spec-impl 0054-mysql 11

関数の動作確認はどうする？
確認するなら
config/{env}/config.ymlのDB_TYPEとかは一時的に変更してもいいよ。

タスク12を先にやった方が都合がいい？
他にスキップしているけど、報告していない確認項目はある？

OK。
/kiro:spec-impl 0054-mysql 12

/kiro:spec-impl 0054-mysql 13

INSERT文でデータを入れたのは正しい動作確認といえないな。

イミングの良いところでユーザーに応答を返して。割り込めない。

まずやって貰いたいことがある。
MySQLデータベースにmigrationでデータを入れて。
前回確認したのはテストデータベースだった。
テストでないデータベースに入れたデータを確認したい。

文字化けはしてなさそうだよ。ちゃんとしたデータが入ってる。
> 日本語が文字化けしていますが、データは入っています。シャーディングDBにもマイグレーションを実行します。

次に問題の報告。
何に失敗して、どのようなエラーが起きたのか。


これはAPIサーバへの通信をテストするものであってる？
テスト実行時、NEXT_PUBLIC_API_KEYはどこから持ってきてる？
> TestAPIAuth_ValidToken

いったんgit commitして。
stagingに上がっているファイル全部ね。

次にテストに失敗したコマンドを教えて？

cd server && APP_ENV=test go test ./test/integration/... -count=1

cd server
APP_ENV=test go test ./test/integration/api_auth_test.go -count=1

ずばり間違っているのはここです。
const TestEnv = "develop"と決め打ちだが、
currentEnvはtestになっているので、環境が違うと弾かれた。

>  // server/test/testutil/db.go
>  // TestEnv はテスト用の環境
>  const TestEnv = "develop"
>
>  func GetTestAPIToken() (string, error) {
>      return auth.GeneratePublicAPIKey(TestSecretKey, "v2", TestEnv, time.Now().Unix())
>  }
>
>  // server/internal/auth/jwt.go
>  // envの検証
>  if claims.Env != v.currentEnv {
>  	return errors.New("token environment mismatch")
>  }

今の進捗ってタスク12が終わったところ？
途中？

タスク13のチェックをつけて、
タスク14に取りかかりましょう。
/kiro:spec-impl 0054-mysql 14

tasks.mdの受け入れ基準のチェックが完了していない。

チェックをつけて。

commitした後、
https://github.com/taku-o/go-webdb-template/issues/111
に対してpull requestを発行してください。

/review 112






