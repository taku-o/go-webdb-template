# generate-sample-dataコマンドPostgreSQL対応実装タスク一覧

## 概要
`server/cmd/generate-sample-data`コマンドでPostgreSQLを利用するように修正を段階的に進めるためのタスク一覧。要件定義書と設計書に基づき、実装可能な粒度でタスクを分解。

## 実装フェーズ

### Phase 1: コマンド実装の確認

#### タスク 1.1: server/cmd/generate-sample-data/main.goの確認
**目的**: SQLite用ライブラリのインポートと処理分岐が存在しないことを確認。PostgreSQL対応のコードパスを使用していることを確認。

**作業内容**:
- `server/cmd/generate-sample-data/main.go`を開く
- インポートセクションを確認:
  - `_ "github.com/mattn/go-sqlite3"`のインポートが存在しないことを確認
  - `"gorm.io/driver/sqlite"`のインポートが存在しないことを確認
  - PostgreSQLドライバー（`gorm.io/driver/postgres`）が使用されていないことを確認（`db.NewGroupManager(cfg)`で自動的に使用される）
- 処理分岐を確認:
  - `case "sqlite3":`分岐が存在しないことを確認
  - `driver = "sqlite3"`のデフォルト値設定が存在しないことを確認
- データベース接続の確認:
  - `config.Load()`で設定ファイルから接続情報を読み込んでいることを確認
  - `db.NewGroupManager(cfg)`でPostgreSQL接続が確立されることを確認
- 必要に応じて`grep`コマンドで確認:
  - `grep -n "sqlite" server/cmd/generate-sample-data/main.go`でSQLite関連のコードを検索
  - `grep -n "postgres" server/cmd/generate-sample-data/main.go`でPostgreSQL関連のコードを検索

**受け入れ基準**:
- SQLite用ライブラリ（`_ "github.com/mattn/go-sqlite3"`, `"gorm.io/driver/sqlite"`）のインポートが存在しない
- SQLite用処理分岐（`case "sqlite3":`, `driver = "sqlite3"`）が存在しない
- `config.Load()`で設定ファイルから接続情報を読み込んでいる
- `db.NewGroupManager(cfg)`でPostgreSQL接続が確立される

- _Requirements: 3.1.1, 3.1.3_
- _Design: 3.1.1_

---

### Phase 2: 設定ファイルの確認

#### タスク 2.1: config/develop/database.yamlの確認
**目的**: PostgreSQL設定が有効になっていることを確認。論理シャーディング数8を確認。SQLite設定が削除されていることを確認。

**作業内容**:
- `config/develop/database.yaml`を開く
- PostgreSQL設定の確認:
  - PostgreSQL設定が有効になっていることを確認（Issue #87で対応済みの想定）
  - masterグループのPostgreSQL接続情報が正しく設定されていることを確認
  - shardingグループのPostgreSQL接続情報が正しく設定されていることを確認（論理シャーディング数8）
- SQLite設定の確認:
  - SQLite設定が削除されていることを確認（コメントアウトではない）
- シャーディング構成の確認:
  - 物理データベース数: master 1台 + sharding 4台 = 合計5台
  - 論理シャーディング数: 8（各物理DBに2つの論理シャードを割り当て）
  - shardingグループに8つのデータベース設定（id: 1-8）が定義されていることを確認
  - 各論理シャードのtable_rangeが正しく設定されていることを確認:
    - id: 1, table_range: [0, 3] → postgres-sharding-1 (webdb_sharding_1)
    - id: 2, table_range: [4, 7] → postgres-sharding-1 (webdb_sharding_1)
    - id: 3, table_range: [8, 11] → postgres-sharding-2 (webdb_sharding_2)
    - id: 4, table_range: [12, 15] → postgres-sharding-2 (webdb_sharding_2)
    - id: 5, table_range: [16, 19] → postgres-sharding-3 (webdb_sharding_3)
    - id: 6, table_range: [20, 23] → postgres-sharding-3 (webdb_sharding_3)
    - id: 7, table_range: [24, 27] → postgres-sharding-4 (webdb_sharding_4)
    - id: 8, table_range: [28, 31] → postgres-sharding-4 (webdb_sharding_4)

**受け入れ基準**:
- PostgreSQL設定が有効になっている
- masterグループのPostgreSQL接続情報が正しく設定されている
- shardingグループのPostgreSQL接続情報が正しく設定されている（論理シャーディング数8）
- SQLite設定が削除されている
- shardingグループに8つのデータベース設定（id: 1-8）が定義されている
- 各論理シャードのtable_rangeが正しく設定されている

- _Requirements: 3.1.2_
- _Design: 3.1.2_

---

#### タスク 2.2: config/staging/database.yamlの確認
**目的**: PostgreSQL設定を確認し、論理シャーディング数8を確認。SQLite設定があれば削除。

**作業内容**:
- `config/staging/database.yaml`を開く
- PostgreSQL設定が正しく定義されているか確認
- シャーディング構成が正しいか確認（物理DB 4台、論理シャーディング8）
- shardingグループに8つのデータベース設定（id: 1-8）が定義されているか確認
- 各論理シャードのtable_rangeが正しく設定されているか確認
- SQLite設定が存在する場合は削除

**受け入れ基準**:
- PostgreSQL設定が正しく定義されている
- シャーディング構成が正しい（物理DB 4台、論理シャーディング8）
- shardingグループに8つのデータベース設定（id: 1-8）が定義されている
- 各論理シャードのtable_rangeが正しく設定されている
- SQLite設定が削除されている（存在する場合）

- _Requirements: 3.1.2_
- _Design: 3.1.2_

---

#### タスク 2.3: config/production/database.yamlの確認
**目的**: PostgreSQL設定を確認し、論理シャーディング数8を確認。SQLite設定があれば削除（存在する場合）。

**作業内容**:
- `config/production/database.yaml`が存在するか確認
- 存在する場合:
  - ファイルを開く
  - PostgreSQL設定が正しく定義されているか確認
  - シャーディング構成が正しいか確認（物理DB 4台、論理シャーディング8）
  - shardingグループに8つのデータベース設定（id: 1-8）が定義されているか確認
  - 各論理シャードのtable_rangeが正しく設定されているか確認
  - SQLite設定が存在する場合は削除

**受け入れ基準**:
- ファイルが存在しない場合はスキップ
- ファイルが存在する場合:
  - PostgreSQL設定が正しく定義されている
  - シャーディング構成が正しい（物理DB 4台、論理シャーディング8）
  - shardingグループに8つのデータベース設定（id: 1-8）が定義されている
  - 各論理シャードのtable_rangeが正しく設定されている
  - SQLite設定が削除されている（存在する場合）

- _Requirements: 3.1.2_
- _Design: 3.1.2_

---

### Phase 3: SQL構文の確認

#### タスク 3.1: プレースホルダー構文の確認
**目的**: GORMのプレースホルダー構文（`?`）がPostgreSQLで正常に動作することを確認。

**作業内容**:
- `server/cmd/generate-sample-data/main.go`を開く
- `insertDmUsersBatch`関数を確認:
  - GORMのプレースホルダー構文（`?`）が使用されていることを確認
  - GORMがPostgreSQLドライバーで自動的に`$1`, `$2`に変換することを確認
- `insertDmPostsBatch`関数を確認:
  - GORMのプレースホルダー構文（`?`）が使用されていることを確認
  - GORMがPostgreSQLドライバーで自動的に`$1`, `$2`に変換することを確認
- 必要に応じてGORMのドキュメントを確認:
  - GORMのプレースホルダー構文がPostgreSQLドライバーで自動的に変換されることを確認

**受け入れ基準**:
- `insertDmUsersBatch`関数でGORMのプレースホルダー構文（`?`）が使用されている
- `insertDmPostsBatch`関数でGORMのプレースホルダー構文（`?`）が使用されている
- GORMがPostgreSQLドライバーで自動的に`$1`, `$2`に変換することを確認

- _Requirements: 3.3.1_
- _Design: 3.2.1_

---

#### タスク 3.2: 動的テーブル名の確認
**目的**: 動的テーブル名の使用がPostgreSQLで正常に動作することを確認。

**作業内容**:
- `server/cmd/generate-sample-data/main.go`を開く
- `insertDmUsersBatch`関数を確認:
  - 動的テーブル名（`fmt.Sprintf("dm_users_%03d", tableNumber)`）が使用されていることを確認
  - SQLインジェクション対策が適切に行われていることを確認（テーブル名は数値から生成されるため安全）
- `insertDmPostsBatch`関数を確認:
  - 動的テーブル名（`fmt.Sprintf("dm_posts_%03d", tableNumber)`）が使用されていることを確認
  - SQLインジェクション対策が適切に行われていることを確認（テーブル名は数値から生成されるため安全）

**受け入れ基準**:
- `insertDmUsersBatch`関数で動的テーブル名が使用されている
- `insertDmPostsBatch`関数で動的テーブル名が使用されている
- SQLインジェクション対策が適切に行われている（テーブル名は数値から生成されるため安全）

- _Requirements: 3.3.1_
- _Design: 3.2.2_

---

#### タスク 3.3: バッチ挿入の確認
**目的**: バッチ挿入のSQL構文がPostgreSQL対応であることを確認。

**作業内容**:
- `server/cmd/generate-sample-data/main.go`を開く
- `insertDmUsersBatch`関数を確認:
  - バッチ挿入のSQL構文がPostgreSQL対応であることを確認
  - バッチサイズ（500件）が適切であることを確認
  - エラーハンドリングが適切であることを確認
- `insertDmPostsBatch`関数を確認:
  - バッチ挿入のSQL構文がPostgreSQL対応であることを確認
  - バッチサイズ（500件）が適切であることを確認
  - エラーハンドリングが適切であることを確認
- `insertDmNewsBatch`関数を確認:
  - GORMの`CreateInBatches`がPostgreSQLで正常に動作することを確認
  - バッチサイズ（500件）が適切であることを確認
  - エラーハンドリングが適切であることを確認

**受け入れ基準**:
- `insertDmUsersBatch`関数のバッチ挿入SQL構文がPostgreSQL対応である
- `insertDmPostsBatch`関数のバッチ挿入SQL構文がPostgreSQL対応である
- `insertDmNewsBatch`関数のGORMの`CreateInBatches`がPostgreSQL対応である
- バッチサイズ（500件）が適切である
- エラーハンドリングが適切である

- _Requirements: 3.3.1_
- _Design: 3.2.3_

---

### Phase 4: データ生成処理の確認

#### タスク 4.1: dm_usersテーブルへのデータ生成の確認
**目的**: `generateDmUsers`関数がPostgreSQLで正常に動作することを確認。

**作業内容**:
- `server/cmd/generate-sample-data/main.go`を開く
- `generateDmUsers`関数を確認:
  - UUIDv7の生成が正常に動作することを確認
  - テーブル番号の計算が正しく動作することを確認
    - `tableSelector.GetTableNumberFromUUID(id)`が正しく動作することを確認
  - シャーディング接続が正しく動作することを確認
    - `groupManager.GetShardingConnection(tableNumber)`が正しく動作することを確認
  - バッチ挿入がPostgreSQLで正常に動作することを確認
    - `insertDmUsersBatch`関数がPostgreSQLで正常に動作することを確認

**受け入れ基準**:
- UUIDv7の生成が正常に動作する
- テーブル番号の計算が正しく動作する
- シャーディング接続が正しく動作する
- バッチ挿入がPostgreSQLで正常に動作する

- _Requirements: 3.2.1_
- _Design: 3.3.1_

---

#### タスク 4.2: dm_postsテーブルへのデータ生成の確認
**目的**: `generateDmPosts`関数がPostgreSQLで正常に動作することを確認。

**作業内容**:
- `server/cmd/generate-sample-data/main.go`を開く
- `generateDmPosts`関数を確認:
  - UUIDv7の生成が正常に動作することを確認
  - テーブル番号の計算が正しく動作することを確認（user_idベース）
    - `tableSelector.GetTableNumberFromUUID(dmUserID)`が正しく動作することを確認
  - シャーディング接続が正しく動作することを確認
    - `groupManager.GetShardingConnection(tableNumber)`が正しく動作することを確認
  - バッチ挿入がPostgreSQLで正常に動作することを確認
    - `insertDmPostsBatch`関数がPostgreSQLで正常に動作することを確認

**受け入れ基準**:
- UUIDv7の生成が正常に動作する
- テーブル番号の計算が正しく動作する（user_idベース）
- シャーディング接続が正しく動作する
- バッチ挿入がPostgreSQLで正常に動作する

- _Requirements: 3.2.2_
- _Design: 3.3.2_

---

#### タスク 4.3: dm_newsテーブルへのデータ生成の確認
**目的**: `generateDmNews`関数がPostgreSQLで正常に動作することを確認。

**作業内容**:
- `server/cmd/generate-sample-data/main.go`を開く
- `generateDmNews`関数を確認:
  - master接続が正しく動作することを確認
    - `groupManager.GetMasterConnection()`が正しく動作することを確認
  - バッチ挿入がPostgreSQLで正常に動作することを確認
    - `insertDmNewsBatch`関数がPostgreSQLで正常に動作することを確認
    - GORMの`CreateInBatches`がPostgreSQLで正常に動作することを確認

**受け入れ基準**:
- master接続が正しく動作する
- バッチ挿入がPostgreSQLで正常に動作する
- GORMの`CreateInBatches`がPostgreSQLで正常に動作する

- _Requirements: 3.2.3_
- _Design: 3.3.3_

---

### Phase 5: テストコードの確認・修正

#### タスク 5.1: テストコードの存在確認
**目的**: テストファイルが存在するか確認。SQLite依存がないことを確認。

**作業内容**:
- `server/cmd/generate-sample-data/`配下のテストファイルを確認
- テストファイルが存在するか確認:
  - `server/cmd/generate-sample-data/*_test.go`ファイルを検索
  - テストファイルが存在しない場合は、このタスクをスキップ
- テストコードがSQLiteに依存しているか確認:
  - `grep -n "sqlite" server/cmd/generate-sample-data/*_test.go`でSQLite関連のコードを検索
  - SQLite固有の設定やコードが含まれているか確認

**受け入れ基準**:
- テストファイルが存在するか確認された
- テストコードがSQLiteに依存していないことが確認された（存在する場合）

- _Requirements: 3.4.1_
- _Design: 3.4.1_

---

#### タスク 5.2: テストコードのPostgreSQL対応
**目的**: テストコードをPostgreSQL対応に修正（存在する場合）。

**作業内容**:
- テストファイルが存在する場合のみ実行
- SQLite固有の設定をPostgreSQL設定に変更:
  - テストデータベースの接続情報をPostgreSQL設定に変更
  - SQLite固有の設定を削除
- テストデータベースの初期化方法をPostgreSQL対応に変更:
  - テストデータベースの初期化方法をPostgreSQL対応に変更
  - テスト実行時のデータベース接続方法をPostgreSQL対応に変更

**受け入れ基準**:
- テストファイルが存在しない場合はスキップ
- テストファイルが存在する場合:
  - SQLite固有の設定がPostgreSQL設定に変更されている
  - テストデータベースの初期化方法がPostgreSQL対応に変更されている
  - テスト実行時のデータベース接続方法がPostgreSQL対応に変更されている

- _Requirements: 3.4.2_
- _Design: 3.4.2_

---

### Phase 6: 動作確認

#### タスク 6.1: PostgreSQL環境の準備
**目的**: PostgreSQLコンテナを起動し、マイグレーションを適用する。

**作業内容**:
- PostgreSQLコンテナを起動:
  - `docker-compose -f docker-compose.postgres.yml up -d`を実行
  - 各PostgreSQLコンテナが正常に起動していることを確認
- マイグレーションを適用:
  - `scripts/migrate.sh`を実行（Issue #86で対応済み）
  - マイグレーションが正常に適用されたことを確認
- データベース接続を確認:
  - `psql`コマンドまたは`docker exec`コマンドで各データベースに接続できることを確認

**受け入れ基準**:
- PostgreSQLコンテナが正常に起動している
- マイグレーションが正常に適用されている
- 各データベースに接続できる

- _Requirements: 4.1_
- _Design: 9.4_

---

#### タスク 6.2: コマンドの動作確認
**目的**: `generate-sample-data`コマンドがPostgreSQLで正常に動作することを確認。

**作業内容**:
- 環境変数を設定:
  - `APP_ENV=develop`を設定
- コマンドを実行:
  - `cd server && go run cmd/generate-sample-data/main.go`を実行
  - または、ビルド済みのバイナリを実行
- 実行結果を確認:
  - コマンドが正常に実行されることを確認
  - エラーメッセージが出力されないことを確認
  - ログメッセージを確認:
    - "Starting sample data generation..."が出力されることを確認
    - "Generated X dm_users in dm_users_XXX"が出力されることを確認
    - "Generated X dm_posts in dm_posts_XXX"が出力されることを確認
    - "Generated X dm_news articles"が出力されることを確認
    - "Sample data generation completed successfully"が出力されることを確認

**受け入れ基準**:
- コマンドが正常に実行される
- エラーメッセージが出力されない
- ログメッセージが正しく出力される

- _Requirements: 6.4_
- _Design: 9.4_

---

#### タスク 6.3: データ生成の確認
**目的**: 各テーブルにデータが正常に生成されることを確認。

**作業内容**:
- masterデータベースの確認:
  - `psql`コマンドまたは`docker exec`コマンドで`webdb_master`データベースに接続
  - `dm_news`テーブルにデータが生成されていることを確認:
    - `SELECT COUNT(*) FROM dm_news;`を実行してデータ件数を確認
    - 100件のデータが生成されていることを確認
- shardingデータベースの確認:
  - 各shardingデータベース（`webdb_sharding_1` ～ `webdb_sharding_4`）に接続
  - `dm_users`テーブルにデータが生成されていることを確認:
    - 各テーブル（`dm_users_000` ～ `dm_users_031`）にデータが生成されていることを確認
    - データが正しいテーブルに振り分けられていることを確認
  - `dm_posts`テーブルにデータが生成されていることを確認:
    - 各テーブル（`dm_posts_000` ～ `dm_posts_031`）にデータが生成されていることを確認
    - データが正しいテーブルに振り分けられていることを確認

**受け入れ基準**:
- `dm_news`テーブルに100件のデータが生成されている
- `dm_users`テーブルにデータが生成されている（shardingグループの各テーブル）
- `dm_posts`テーブルにデータが生成されている（shardingグループの各テーブル）
- データが正しいテーブルに振り分けられている

- _Requirements: 6.4_
- _Design: 9.4_

---

### Phase 7: ドキュメントの更新

#### タスク 7.1: README.mdの更新
**目的**: `generate-sample-data`コマンドのPostgreSQL利用に関する記述を追加。

**作業内容**:
- `README.md`を開く
- `generate-sample-data`コマンドのセクションを追加または更新:
  - コマンドの説明を追加:
    - `generate-sample-data`コマンドがPostgreSQLを利用することを記載
    - コマンドの目的を記載（サンプルデータの生成）
  - 前提条件を記載:
    - PostgreSQLコンテナの起動方法を記載
    - マイグレーションの適用方法を記載
  - コマンド実行方法を記載:
    - コマンドの実行コマンドを記載
    - 環境変数の設定方法を記載（`APP_ENV=develop`）
    - 実行結果の確認方法を記載
  - トラブルシューティングを記載:
    - よくある問題と解決方法を記載

**受け入れ基準**:
- `generate-sample-data`コマンドのPostgreSQL利用に関する記述が追加されている
- コマンド実行前のPostgreSQL起動手順が記載されている
- コマンド実行方法が記載されている
- 実行結果の確認方法が記載されている
- トラブルシューティングが記載されている

- _Requirements: 3.5.1_
- _Design: 3.5.1_

---

#### タスク 7.2: その他のドキュメントの更新
**目的**: `docs/`配下の関連ドキュメントを更新。

**作業内容**:
- `docs/`配下の関連ドキュメントを確認:
  - `docs/Architecture.md`を確認
  - `docs/Initial-Setup.md`を確認
  - その他の関連ドキュメントを確認
- SQLiteに関する記述をPostgreSQLに変更:
  - `generate-sample-data`コマンドに関するSQLiteの記述をPostgreSQLに変更
- PostgreSQL利用に関する記述を追加:
  - `generate-sample-data`コマンドのPostgreSQL利用に関する記述を追加

**受け入れ基準**:
- SQLiteに関する記述がPostgreSQLに変更されている
- PostgreSQL利用に関する記述が追加されている

- _Requirements: 3.5.2_
- _Design: 3.5.2_

---

## 実装順序

1. Phase 1: コマンド実装の確認
2. Phase 2: 設定ファイルの確認
3. Phase 3: SQL構文の確認
4. Phase 4: データ生成処理の確認
5. Phase 5: テストコードの確認・修正（存在する場合）
6. Phase 6: 動作確認
7. Phase 7: ドキュメントの更新

## 注意事項

- 既存のコードは既にPostgreSQL対応のコードパスを使用しているため、主に確認作業とドキュメント更新が中心になります
- 設定ファイルはIssue #87でPostgreSQL設定が有効化されている想定です
- テストコードは存在しない可能性が高いため、存在しない場合は該当タスクをスキップします
- 動作確認は実際のPostgreSQL環境で実施する必要があります
