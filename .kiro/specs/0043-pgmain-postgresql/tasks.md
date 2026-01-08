# PostgreSQL起動スクリプト・マイグレーションスクリプト修正実装タスク一覧

## 概要
PostgreSQL起動スクリプトとAtlasマイグレーションスクリプトの修正を段階的に進めるためのタスク一覧。要件定義書と設計書に基づき、実装可能な粒度でタスクを分解。

## 実装フェーズ

### Phase 1: HCLファイルのUUIDv7カラム定義修正

#### タスク 1.1: db/schema/sharding_1/dm_posts.hclの修正
**目的**: UUIDv7カラム（`id`, `user_id`）の型定義を`bigint`から`varchar(32)`に修正

**作業内容**:
- `db/schema/sharding_1/dm_posts.hcl`を開く
- すべての`dm_posts_XXX`テーブル（000～007）の`id`カラムを修正:
  - `type = bigint` → `type = varchar(32)`
  - `unsigned = true`を削除
  - `auto_increment = false`を削除
- すべての`dm_posts_XXX`テーブル（000～007）の`user_id`カラムを修正:
  - `type = bigint` → `type = varchar(32)`
  - `unsigned = true`を削除

**受け入れ基準**:
- `dm_posts_000`～`dm_posts_007`のすべての`id`カラムが`varchar(32)`になっている
- `dm_posts_000`～`dm_posts_007`のすべての`user_id`カラムが`varchar(32)`になっている
- `unsigned`属性がすべて削除されている
- `auto_increment`属性がすべて削除されている

- _Requirements: 3.4.1_

---

#### タスク 1.2: db/schema/sharding_1/dm_users.hclの修正
**目的**: UUIDv7カラム（`id`）の型定義を`bigint`から`varchar(32)`に修正

**作業内容**:
- `db/schema/sharding_1/dm_users.hcl`を開く
- すべての`dm_users_XXX`テーブル（000～007）の`id`カラムを修正:
  - `type = bigint` → `type = varchar(32)`
  - `unsigned = true`を削除
  - `auto_increment = false`を削除

**受け入れ基準**:
- `dm_users_000`～`dm_users_007`のすべての`id`カラムが`varchar(32)`になっている
- `unsigned`属性がすべて削除されている
- `auto_increment`属性がすべて削除されている

- _Requirements: 3.4.2_

---

#### タスク 1.3: db/schema/sharding_2/dm_posts.hclの修正 (P)
**目的**: UUIDv7カラム（`id`, `user_id`）の型定義を`bigint`から`varchar(32)`に修正

**作業内容**:
- タスク1.1と同様の作業を`db/schema/sharding_2/dm_posts.hcl`に対して実施
- すべての`dm_posts_XXX`テーブル（008～015）を修正

**受け入れ基準**:
- タスク1.1と同様

- _Requirements: 3.4.1_

---

#### タスク 1.4: db/schema/sharding_2/dm_users.hclの修正 (P)
**目的**: UUIDv7カラム（`id`）の型定義を`bigint`から`varchar(32)`に修正

**作業内容**:
- タスク1.2と同様の作業を`db/schema/sharding_2/dm_users.hcl`に対して実施
- すべての`dm_users_XXX`テーブル（008～015）を修正

**受け入れ基準**:
- タスク1.2と同様

- _Requirements: 3.4.2_

---

#### タスク 1.5: db/schema/sharding_3/dm_posts.hclの修正 (P)
**目的**: UUIDv7カラム（`id`, `user_id`）の型定義を`bigint`から`varchar(32)`に修正

**作業内容**:
- タスク1.1と同様の作業を`db/schema/sharding_3/dm_posts.hcl`に対して実施
- すべての`dm_posts_XXX`テーブル（016～023）を修正

**受け入れ基準**:
- タスク1.1と同様

- _Requirements: 3.4.1_

---

#### タスク 1.6: db/schema/sharding_3/dm_users.hclの修正 (P)
**目的**: UUIDv7カラム（`id`）の型定義を`bigint`から`varchar(32)`に修正

**作業内容**:
- タスク1.2と同様の作業を`db/schema/sharding_3/dm_users.hcl`に対して実施
- すべての`dm_users_XXX`テーブル（016～023）を修正

**受け入れ基準**:
- タスク1.2と同様

- _Requirements: 3.4.2_

---

#### タスク 1.7: db/schema/sharding_4/dm_posts.hclの修正 (P)
**目的**: UUIDv7カラム（`id`, `user_id`）の型定義を`bigint`から`varchar(32)`に修正

**作業内容**:
- タスク1.1と同様の作業を`db/schema/sharding_4/dm_posts.hcl`に対して実施
- すべての`dm_posts_XXX`テーブル（024～031）を修正

**受け入れ基準**:
- タスク1.1と同様

- _Requirements: 3.4.1_

---

#### タスク 1.8: db/schema/sharding_4/dm_users.hclの修正 (P)
**目的**: UUIDv7カラム（`id`）の型定義を`bigint`から`varchar(32)`に修正

**作業内容**:
- タスク1.2と同様の作業を`db/schema/sharding_4/dm_users.hcl`に対して実施
- すべての`dm_users_XXX`テーブル（024～031）を修正

**受け入れ基準**:
- タスク1.2と同様

- _Requirements: 3.4.2_

---

### Phase 2: PostgreSQL起動スクリプトの作成

#### タスク 2.1: 既存PostgreSQLデータのクリーンアップ
**目的**: 過去の1台構成で使用していたPostgreSQLデータを破棄

**作業内容**:
- 既存の`postgres/data/`ディレクトリを確認
- 過去の1台構成で使用していたデータが存在する場合は削除
- 以下のディレクトリをクリーンアップ:
  - `postgres/data/`（既存のデータディレクトリ）
- データディレクトリを削除する前に、必要に応じてバックアップを確認（通常は不要）

**受け入れ基準**:
- 既存のPostgreSQLデータが削除されている
- `postgres/data/`ディレクトリがクリーンな状態になっている

- _Requirements: 3.1.1_

---

#### タスク 2.2: docker-compose.postgres.ymlの修正
**目的**: master 1台 + sharding 4台のPostgreSQLコンテナ構成に修正

**作業内容**:
- `docker-compose.postgres.yml`を開く
- 既存の`postgres`サービスを`postgres-master`に変更:
  - サービス名: `postgres-master`
  - データベース名: `webdb_master`
  - ポート: `5432:5432`
  - ボリューム: `./postgres/data/master:/var/lib/postgresql/data`
- `postgres-sharding-1`サービスを追加:
  - サービス名: `postgres-sharding-1`
  - データベース名: `webdb_sharding_1`
  - ポート: `5433:5432`
  - ボリューム: `./postgres/data/sharding_1:/var/lib/postgresql/data`
- `postgres-sharding-2`サービスを追加:
  - サービス名: `postgres-sharding-2`
  - データベース名: `webdb_sharding_2`
  - ポート: `5434:5432`
  - ボリューム: `./postgres/data/sharding_2:/var/lib/postgresql/data`
- `postgres-sharding-3`サービスを追加:
  - サービス名: `postgres-sharding-3`
  - データベース名: `webdb_sharding_3`
  - ポート: `5435:5432`
  - ボリューム: `./postgres/data/sharding_3:/var/lib/postgresql/data`
- `postgres-sharding-4`サービスを追加:
  - サービス名: `postgres-sharding-4`
  - データベース名: `webdb_sharding_4`
  - ポート: `5436:5432`
  - ボリューム: `./postgres/data/sharding_4:/var/lib/postgresql/data`
- すべてのサービスで既存の`postgres-network`を使用
- すべてのサービスでヘルスチェックを設定

**受け入れ基準**:
- `docker-compose.postgres.yml`に5つのPostgreSQLサービスが定義されている
- 各サービスのポートが競合していない（5432, 5433, 5434, 5435, 5436）
- 各サービスのデータディレクトリが個別にマウントされている
- すべてのサービスが`postgres-network`を使用している

- _Requirements: 3.1.1_

---

#### タスク 2.3: scripts/start-postgres.shの作成
**目的**: PostgreSQLコンテナの起動・停止・状態確認を行うスクリプトを作成

**作業内容**:
- `scripts/start-postgres.sh`を新規作成
- 以下の機能を実装:
  - `start`: `docker-compose -f docker-compose.postgres.yml up -d`を実行
  - `stop`: `docker-compose -f docker-compose.postgres.yml down`を実行
  - `status`: `docker-compose -f docker-compose.postgres.yml ps`を実行
  - `health`: ヘルスチェック状態を確認（`docker-compose -f docker-compose.postgres.yml ps`でhealthcheck状態を表示）
- エラーハンドリングを実装（`set -e`を使用）
- 使用方法のヘルプメッセージを表示

**受け入れ基準**:
- `scripts/start-postgres.sh`が作成されている
- `./scripts/start-postgres.sh start`でPostgreSQLコンテナが起動する
- `./scripts/start-postgres.sh stop`でPostgreSQLコンテナが停止する
- `./scripts/start-postgres.sh status`でコンテナの状態が確認できる
- `./scripts/start-postgres.sh health`でヘルスチェック状態が確認できる

- _Requirements: 3.1.2_

---

#### タスク 2.4: postgres/dataディレクトリの作成
**目的**: PostgreSQLコンテナのデータ永続化用ディレクトリを作成

**作業内容**:
- `postgres/data/master/`ディレクトリを作成
- `postgres/data/sharding_1/`ディレクトリを作成
- `postgres/data/sharding_2/`ディレクトリを作成
- `postgres/data/sharding_3/`ディレクトリを作成
- `postgres/data/sharding_4/`ディレクトリを作成
- `.gitkeep`ファイルを各ディレクトリに追加（空ディレクトリをGit管理下に置くため）

**受け入れ基準**:
- 5つのデータディレクトリが作成されている
- 各ディレクトリに`.gitkeep`ファイルが存在する

- _Requirements: 3.1.1_

---

### Phase 3: マイグレーションスクリプトの新規作成

#### タスク 3.1: 既存scripts/migrate.shのバックアップと削除
**目的**: 既存のSQLite用マイグレーションスクリプトを破棄

**作業内容**:
- 既存の`scripts/migrate.sh`を確認
- 必要に応じてバックアップを作成（`scripts/migrate.sh.backup`など）
- 既存の`scripts/migrate.sh`を削除

**受け入れ基準**:
- 既存の`scripts/migrate.sh`が削除されている（またはバックアップが作成されている）

- _Requirements: 3.2.1_

---

#### タスク 3.2: 設定ファイル読み込み機能の実装
**目的**: `config/{env}/database.yaml`からPostgreSQL接続情報を読み込む機能を実装

**作業内容**:
- `scripts/migrate.sh`の新規作成を開始
- `APP_ENV`環境変数から環境を取得（デフォルト: `develop`）
- `config/{env}/database.yaml`ファイルの存在確認
- YAMLファイルのパース機能を実装:
  - `yq`コマンドを使用するか、Goスクリプトを使用
  - masterグループの接続情報を取得: `database.groups.master[0]`から`host`, `port`, `user`, `password`, `name`
  - shardingグループの接続情報を取得: `database.groups.sharding.databases[]`から各データベースの`host`, `port`, `user`, `password`, `name`
- エラーハンドリングを実装:
  - 設定ファイルが存在しない場合
  - 接続情報が不足している場合

**受け入れ基準**:
- `APP_ENV`環境変数で環境を指定できる
- `config/{env}/database.yaml`から接続情報を読み込める
- masterグループとshardingグループの接続情報を取得できる
- エラーハンドリングが実装されている

- _Requirements: 3.2.3_

---

#### タスク 3.3: Atlasマイグレーション適用機能の実装
**目的**: AtlasマイグレーションをPostgreSQLに適用する機能を実装

**作業内容**:
- PostgreSQL用のURL形式（`postgres://user:password@host:port/dbname?sslmode=disable`）を構築
- masterグループのAtlasマイグレーション適用機能を実装:
  - `atlas migrate apply --dir "file://db/migrations/master" --url "postgres://..."`を実行
- shardingグループのAtlasマイグレーション適用機能を実装:
  - 各shardingデータベース（1～4）に対して`atlas migrate apply`を実行
  - `db/migrations/sharding_1/` ～ `db/migrations/sharding_4/`を適用
- エラーハンドリングを実装:
  - Atlasコマンドの実行エラー
  - データベース接続エラー

**受け入れ基準**:
- masterグループのAtlasマイグレーションが適用できる
- shardingグループのAtlasマイグレーションが適用できる（4台すべて）
- エラーハンドリングが実装されている

- _Requirements: 3.2.1, 3.2.2_

---

#### タスク 3.4: 初期データ適用機能の実装
**目的**: 初期データSQLファイルをPostgreSQL構文に変換して適用する機能を実装

**作業内容**:
- `db/migrations/master/`内の初期データSQLファイルを検出（現在: `20251230045548_seed_data.sql`）
- SQLite構文をPostgreSQL構文に変換:
  - `INSERT OR IGNORE` → `INSERT ... ON CONFLICT DO NOTHING`
  - `datetime('now')` → `NOW()`
- 変換後のSQLをPostgreSQLに適用（`psql`コマンドまたはPostgreSQLクライアントライブラリを使用）
- エラーハンドリングを実装:
  - SQLファイルが存在しない場合
  - SQL構文変換エラー
  - SQL適用エラー

**受け入れ基準**:
- 初期データSQLファイルが検出できる
- SQLite構文がPostgreSQL構文に正しく変換される
- 変換後のSQLがPostgreSQLに適用できる
- エラーハンドリングが実装されている

- _Requirements: 3.2.1_

---

#### タスク 3.5: Viewマイグレーション適用機能の実装
**目的**: View定義SQLファイルを生SQLで適用する機能を実装

**作業内容**:
- `db/migrations/view_master/`内のView定義SQLファイルを検出（現在: `20260103030225_create_dm_news_view.sql`）
- ファイル名順にソートして適用（Atlasはファイル名順に適用するため）
- 生SQLをPostgreSQLに適用（`psql`コマンドまたはPostgreSQLクライアントライブラリを使用）
- エラーハンドリングを実装:
  - SQLファイルが存在しない場合
  - SQL適用エラー

**受け入れ基準**:
- View定義SQLファイルが検出できる
- ファイル名順にソートして適用できる
- 生SQLがPostgreSQLに適用できる
- エラーハンドリングが実装されている

- _Requirements: 3.2.1_

---

#### タスク 3.6: マイグレーションスクリプトの統合とテスト
**目的**: すべての機能を統合し、動作確認を行う

**作業内容**:
- タスク3.2～3.5で実装した機能を統合
- コマンドライン引数の処理を実装:
  - `master`: masterグループのみマイグレーション適用
  - `sharding`: shardingグループのみマイグレーション適用
  - `all`: すべてのマイグレーション適用（デフォルト）
- 使用方法のヘルプメッセージを表示
- 動作確認:
  - `./scripts/migrate.sh master`でmasterグループのマイグレーションが適用できることを確認
  - `./scripts/migrate.sh sharding`でshardingグループのマイグレーションが適用できることを確認
  - `./scripts/migrate.sh all`で全マイグレーションが適用できることを確認

**受け入れ基準**:
- すべての機能が統合されている
- コマンドライン引数でマイグレーション対象を指定できる
- 各コマンドが正常に動作する
- エラーハンドリングが適切に機能する

- _Requirements: 3.2.1_

---

### Phase 4: データ投入確認

#### タスク 4.1: PostgreSQL起動とマイグレーション適用の確認
**目的**: PostgreSQLコンテナが正常に起動し、マイグレーションが適用できることを確認

**作業内容**:
- `./scripts/start-postgres.sh start`でPostgreSQLコンテナを起動
- 各PostgreSQLコンテナが正常に起動していることを確認（`./scripts/start-postgres.sh status`）
- ヘルスチェックが正常であることを確認（`./scripts/start-postgres.sh health`）
- `./scripts/migrate.sh all`で全マイグレーションを適用
- マイグレーション適用時にエラーが発生しないことを確認

**受け入れ基準**:
- PostgreSQLコンテナが正常に起動する
- すべてのマイグレーションが正常に適用される
- エラーが発生しない

- _Requirements: 3.3.1_

---

#### タスク 4.2: テーブル存在確認
**目的**: 各データベースにテーブルが作成されていることを確認

**作業内容**:
- masterデータベースに接続してテーブル一覧を確認（`psql`コマンドまたはAtlas CLI）
- 各shardingデータベース（1～4）に接続してテーブル一覧を確認
- 以下のテーブルが存在することを確認:
  - master: `dm_news`, GoAdmin関連テーブル
  - sharding: `dm_users_XXX`, `dm_posts_XXX`（各shardingデータベースごと）

**受け入れ基準**:
- すべてのデータベースにテーブルが作成されている
- 期待されるテーブルが存在する

- _Requirements: 3.3.1_

---

#### タスク 4.3: データ投入確認
**目的**: テストデータを投入してデータベースにデータが入ることを確認

**作業内容**:
- テストデータを各データベースに投入
- データが正しく保存されていることを確認（`psql`コマンドでSELECT文を実行）
- UUIDv7カラム（`id`, `user_id`）が`varchar(32)`型で正しく保存されていることを確認

**受け入れ基準**:
- テストデータが正しく投入できる
- データが正しく保存されている
- UUIDv7カラムが正しい型で保存されている

- _Requirements: 3.3.1_

---

### Phase 5: ドキュメント整備

#### タスク 5.1: PostgreSQL起動手順のドキュメント作成
**目的**: PostgreSQL起動手順をドキュメント化

**作業内容**:
- `README.md`または`docs/`にPostgreSQL起動手順を追記
- 以下の内容を含める:
  - PostgreSQLコンテナの起動方法（`./scripts/start-postgres.sh start`）
  - PostgreSQLコンテナの停止方法（`./scripts/start-postgres.sh stop`）
  - PostgreSQLコンテナの状態確認方法（`./scripts/start-postgres.sh status`）
  - ヘルスチェック確認方法（`./scripts/start-postgres.sh health`）

**受け入れ基準**:
- PostgreSQL起動手順がドキュメント化されている
- 各コマンドの使用方法が記載されている

- _Requirements: 6.5_

---

#### タスク 5.2: マイグレーション手順のドキュメント作成
**目的**: PostgreSQLマイグレーション手順をドキュメント化

**作業内容**:
- `README.md`または`docs/`にPostgreSQLマイグレーション手順を追記
- 以下の内容を含める:
  - マイグレーション適用方法（`./scripts/migrate.sh [master|sharding|all]`）
  - 設定ファイル（`config/{env}/database.yaml`）の設定方法
  - 環境変数（`APP_ENV`）の使用方法
  - トラブルシューティング情報

**受け入れ基準**:
- PostgreSQLマイグレーション手順がドキュメント化されている
- 設定ファイルの設定方法が記載されている
- トラブルシューティング情報が記載されている

- _Requirements: 6.5_

---

## 実装順序

1. **Phase 1**: HCLファイルのUUIDv7カラム定義修正（タスク1.1～1.8）
2. **Phase 2**: PostgreSQL起動スクリプトの作成（タスク2.1～2.4）
3. **Phase 3**: マイグレーションスクリプトの新規作成（タスク3.1～3.6）
4. **Phase 4**: データ投入確認（タスク4.1～4.3）
5. **Phase 5**: ドキュメント整備（タスク5.1～5.2）

## 注意事項

- タスク1.3～1.8は並列実行可能（Pマーク付き）
- タスク2.1はタスク2.2の前に実行（既存データのクリーンアップ）
- タスク3.2～3.5は順次実行が必要（依存関係あり）
- タスク4.1はタスク2.1～2.3とタスク3.1～3.6の完了後に実行
- タスク4.2～4.3はタスク4.1の完了後に実行
