# generate-sample-dataコマンドPostgreSQL対応要件定義書

## 1. 概要

### 1.1 プロジェクト情報
- **プロジェクト名**: go-webdb-template
- **Issue番号**: #89
- **親Issue番号**: #85
- **Issueタイトル**: server/cmd/generate-sample-data コマンド
- **Feature名**: 0046-pgmain-sample
- **作成日**: 2025-01-27

### 1.2 目的
`server/cmd/generate-sample-data`コマンドで新しく用意したPostgreSQLを利用するように修正する。既存のSQLite設定をPostgreSQL設定に切り替え、開発環境、staging環境、production環境でPostgreSQLを利用できるようにする。

### 1.3 スコープ
- `server/cmd/generate-sample-data/main.go`のPostgreSQL対応確認・修正
- 設定ファイル（`config/{env}/database.yaml`）のPostgreSQL設定確認
- SQLite用ライブラリのインポート削除（存在する場合）
- SQLite用処理分岐の削除（存在する場合）
- テストコードのPostgreSQL対応確認・修正（存在する場合）
- ドキュメントの更新

**本実装の範囲外**:
- PostgreSQLの起動スクリプト・マイグレーションスクリプトの修正（Issue #86で対応済み）
- APIサーバーのPostgreSQL接続設定変更（Issue #87で対応済み）
- GoAdminサーバーのPostgreSQL接続設定変更（別Issue対応）
- 本番環境へのデプロイ（準備のみ）

## 2. 背景・現状分析

### 2.1 現在の実装
- **コマンド実装**: `server/cmd/generate-sample-data/main.go`で`config.Load()`と`db.NewGroupManager(cfg)`を使用してデータベース接続を行っている
- **データベース接続**: `server/internal/db/connection.go`で既にPostgreSQLドライバー（`gorm.io/driver/postgres`）がサポートされている
- **設定ファイル**: `config/{env}/database.yaml`から接続情報を読み込む
- **データベース構成**:
  - マスターデータベース: 1台（`master.db` → `webdb_master`）
  - シャーディングデータベース: 4台（`sharding_db_1.db` ～ `sharding_db_4.db` → `webdb_sharding_1` ～ `webdb_sharding_4`）
  - **論理シャーディング数: 8**（物理DB 4台 × 2論理シャード）
- **生成対象データ**:
  - `dm_users`: shardingグループの各テーブル（`dm_users_000` ～ `dm_users_031`）にデータを生成
  - `dm_posts`: shardingグループの各テーブル（`dm_posts_000` ～ `dm_posts_031`）にデータを生成
  - `dm_news`: masterグループの`dm_news`テーブルにデータを生成

### 2.2 課題点
1. **設定ファイルのSQLite依存**: `config/{env}/database.yaml`でSQLite設定が有効になっており、PostgreSQL設定がコメントアウトされている（または未定義）可能性がある
2. **SQLite用ライブラリの残存**: ソースコード中にSQLite用ライブラリのインポートが残っている可能性がある
3. **SQLite用処理分岐の残存**: ソースコード中にSQLite用の処理分岐が残っている可能性がある
4. **テストコードのSQLite依存**: テストコードがSQLiteに依存している可能性がある
5. **ドキュメントの未更新**: PostgreSQL利用に関するドキュメントが未整備

### 2.3 本実装による改善点
1. **PostgreSQL環境への移行**: 開発環境、staging環境、production環境でPostgreSQLを利用できるようにする
2. **コードの簡素化**: SQLite用ライブラリと処理分岐を削除し、コードを簡素化
3. **テストコードの修正**: テストコードをPostgreSQL対応に修正（存在する場合）
4. **ドキュメントの整備**: PostgreSQL利用に関するドキュメントを整備

## 3. 機能要件

### 3.1 コマンド実装の確認・修正

#### 3.1.1 server/cmd/generate-sample-data/main.goの確認
- **ファイル**: `server/cmd/generate-sample-data/main.go`
- **確認内容**:
  - `config.Load()`で設定ファイルから接続情報を読み込んでいることを確認
  - `db.NewGroupManager(cfg)`でPostgreSQL接続が確立されることを確認
  - SQLite用ライブラリのインポートが存在しないことを確認
  - SQLite用処理分岐が存在しないことを確認
  - データ生成処理がPostgreSQLで正常に動作することを確認

#### 3.1.2 設定ファイルの確認
- **ファイル**: `config/{env}/database.yaml`
- **確認内容**:
  - PostgreSQL設定が有効になっていることを確認（Issue #87で対応済みの想定）
  - masterグループのPostgreSQL接続情報が正しく設定されていることを確認
  - shardingグループのPostgreSQL接続情報が正しく設定されていることを確認（論理シャーディング数8）
  - SQLite設定が削除されていることを確認

#### 3.1.3 SQLite用ライブラリと処理分岐の削除
- **対象ファイル**: `server/cmd/generate-sample-data/main.go`
- **確認内容**:
  - SQLite用ライブラリのインポートが存在しないことを確認
  - SQLite用処理分岐が存在しないことを確認
- **削除内容**（存在する場合）:
  - `_ "github.com/mattn/go-sqlite3"`のインポートを削除
  - `"gorm.io/driver/sqlite"`のインポートを削除
  - SQLite用の処理分岐を削除

### 3.2 データ生成処理の確認

#### 3.2.1 dm_usersテーブルへのデータ生成
- **対象**: shardingグループの各テーブル（`dm_users_000` ～ `dm_users_031`）
- **確認内容**:
  - `generateDmUsers`関数がPostgreSQLで正常に動作することを確認
  - UUIDv7の生成が正常に動作することを確認
  - テーブル番号の計算が正しく動作することを確認
  - バッチ挿入がPostgreSQLで正常に動作することを確認
  - `insertDmUsersBatch`関数のSQL構文がPostgreSQL対応であることを確認

#### 3.2.2 dm_postsテーブルへのデータ生成
- **対象**: shardingグループの各テーブル（`dm_posts_000` ～ `dm_posts_031`）
- **確認内容**:
  - `generateDmPosts`関数がPostgreSQLで正常に動作することを確認
  - UUIDv7の生成が正常に動作することを確認
  - テーブル番号の計算が正しく動作することを確認（user_idベース）
  - バッチ挿入がPostgreSQLで正常に動作することを確認
  - `insertDmPostsBatch`関数のSQL構文がPostgreSQL対応であることを確認

#### 3.2.3 dm_newsテーブルへのデータ生成
- **対象**: masterグループの`dm_news`テーブル
- **確認内容**:
  - `generateDmNews`関数がPostgreSQLで正常に動作することを確認
  - バッチ挿入がPostgreSQLで正常に動作することを確認
  - `insertDmNewsBatch`関数がPostgreSQLで正常に動作することを確認

### 3.3 SQL構文の確認

#### 3.3.1 INSERT文の確認
- **対象**: `insertDmUsersBatch`, `insertDmPostsBatch`, `insertDmNewsBatch`関数
- **確認内容**:
  - PostgreSQLのプレースホルダー構文（`$1`, `$2`, ...）またはGORMのプレースホルダー構文（`?`）が正しく使用されていることを確認
  - 動的テーブル名の使用がPostgreSQLで正常に動作することを確認
  - バッチ挿入のSQL構文がPostgreSQL対応であることを確認

#### 3.3.2 データ型の確認
- **確認内容**:
  - UUIDv7のデータ型が`varchar(32)`であることを確認（Issue #86で対応済み）
  - タイムスタンプのデータ型がPostgreSQL対応であることを確認
  - その他のデータ型がPostgreSQL対応であることを確認

### 3.4 テストコードの確認・修正

#### 3.4.1 テストコードの存在確認
- **対象ファイル**: `server/cmd/generate-sample-data/`配下のテストファイル
- **確認内容**:
  - テストファイルが存在するか確認
  - テストコードがSQLiteに依存しているか確認

#### 3.4.2 テストコードのPostgreSQL対応
- **修正内容**（テストコードが存在する場合）:
  - SQLite固有の設定をPostgreSQL設定に変更
  - テストデータベースの初期化方法をPostgreSQL対応に変更
  - テスト実行時のデータベース接続方法をPostgreSQL対応に変更

### 3.5 ドキュメントの更新

#### 3.5.1 README.mdの更新
- **ファイル**: `README.md`
- **更新内容**:
  - `generate-sample-data`コマンドのPostgreSQL利用に関する記述を追加
  - コマンド実行前のPostgreSQL起動手順を記載
  - コマンド実行方法を記載

#### 3.5.2 その他のドキュメントの更新
- **対象ファイル**: `docs/`配下の関連ドキュメント
- **更新内容**:
  - SQLiteに関する記述をPostgreSQLに変更
  - PostgreSQL利用に関する記述を追加

## 4. 非機能要件

### 4.1 PostgreSQL環境の前提条件
- Issue #86でPostgreSQLの起動スクリプトとマイグレーションスクリプトが修正されていること
- Issue #87でAPIサーバーのPostgreSQL接続設定が修正されていること
- PostgreSQLコンテナが起動していること（開発環境の場合）
- マイグレーションが適用されていること

### 4.2 パフォーマンス
- **データ生成時間**: サンプルデータ生成の実行時間を最小化
- **バッチサイズ**: 既存のバッチサイズ（500件）を維持
- **接続時間**: データベース接続の確立時間を最小化

### 4.3 セキュリティ
- **パスワード管理**: 開発環境では設定ファイルに固定パスワード（`webdb`）を記載する
- **SSL/TLS**: 開発環境では`sslmode=disable`、本番環境では適切なSSL設定を推奨

### 4.4 環境別対応
- **開発環境**: ローカルDocker環境でのPostgreSQL接続
- **staging環境**: staging環境用のPostgreSQL接続設定
- **production環境**: production環境用のPostgreSQL接続設定

## 5. 制約事項

### 5.1 既存システムとの関係
- **既存のデータベース接続コード**: `server/internal/db/connection.go`は既にPostgreSQLドライバーをサポートしている
- **既存の設定ファイル構造**: `config/{env}/database.yaml`の構造は維持（Issue #87でPostgreSQL設定が有効化されている想定）
- **依存関係**: `server/go.mod`からSQLite関連の依存関係が削除されている想定（Issue #87で対応済み）

### 5.2 技術スタック
- **PostgreSQL**: PostgreSQL 15-alpine（Dockerイメージ）
- **GORM**: 既存のGORMライブラリを使用
- **PostgreSQLドライバー**: `gorm.io/driver/postgres`（既存）

### 5.3 シャーディング構成
- **物理データベース数**: master 1台 + sharding 4台 = 合計5台
- **論理シャーディング数**: **8**（各物理DBに2つの論理シャードを割り当て）
- **table_range**: 
  - id: 1, table_range: [0, 3] → postgres-sharding-1 (webdb_sharding_1)
  - id: 2, table_range: [4, 7] → postgres-sharding-1 (webdb_sharding_1)
  - id: 3, table_range: [8, 11] → postgres-sharding-2 (webdb_sharding_2)
  - id: 4, table_range: [12, 15] → postgres-sharding-2 (webdb_sharding_2)
  - id: 5, table_range: [16, 19] → postgres-sharding-3 (webdb_sharding_3)
  - id: 6, table_range: [20, 23] → postgres-sharding-3 (webdb_sharding_3)
  - id: 7, table_range: [24, 27] → postgres-sharding-4 (webdb_sharding_4)
  - id: 8, table_range: [28, 31] → postgres-sharding-4 (webdb_sharding_4)

### 5.4 運用上の制約
- **起動順序**: PostgreSQLコンテナを起動してからコマンドを実行する必要がある
- **マイグレーション**: マイグレーション適用後にコマンドを実行する必要がある

## 6. 受け入れ基準

### 6.1 コマンド実装の確認・修正
- [ ] `server/cmd/generate-sample-data/main.go`でSQLite用ライブラリのインポートが存在しない
- [ ] `server/cmd/generate-sample-data/main.go`でSQLite用処理分岐が存在しない
- [ ] `server/cmd/generate-sample-data/main.go`がPostgreSQLで正常に動作する

### 6.2 設定ファイルの確認
- [ ] `config/{env}/database.yaml`でPostgreSQL設定が有効になっている（Issue #87で対応済みの想定）
- [ ] `config/{env}/database.yaml`のPostgreSQL接続情報が正しく設定されている（master 1台、sharding論理シャーディング8つ、物理DB 4台）

### 6.3 データ生成処理の確認
- [ ] `generateDmUsers`関数がPostgreSQLで正常に動作する
- [ ] `generateDmPosts`関数がPostgreSQLで正常に動作する
- [ ] `generateDmNews`関数がPostgreSQLで正常に動作する
- [ ] `insertDmUsersBatch`関数のSQL構文がPostgreSQL対応である
- [ ] `insertDmPostsBatch`関数のSQL構文がPostgreSQL対応である
- [ ] `insertDmNewsBatch`関数がPostgreSQLで正常に動作する

### 6.4 コマンドの動作確認
- [ ] `generate-sample-data`コマンドがPostgreSQLに正常に接続できる
- [ ] `generate-sample-data`コマンドがmasterデータベースに正常に接続できる
- [ ] `generate-sample-data`コマンドがshardingデータベースに正常に接続できる（4台すべて）
- [ ] `generate-sample-data`コマンドが正常に実行され、サンプルデータが生成される
- [ ] `dm_users`テーブルにデータが正常に生成される（shardingグループの各テーブル）
- [ ] `dm_posts`テーブルにデータが正常に生成される（shardingグループの各テーブル）
- [ ] `dm_news`テーブルにデータが正常に生成される（masterグループ）

### 6.5 テストコードの確認・修正
- [ ] テストコードのSQLite依存が確認され、必要に応じて修正されている（テストコードが存在する場合）
- [ ] テストがPostgreSQL環境で正常に実行できる（テストコードが存在する場合）

### 6.6 ドキュメント
- [ ] `README.md`に`generate-sample-data`コマンドのPostgreSQL利用に関する記述が追加されている
- [ ] コマンド実行前のPostgreSQL起動手順が記載されている
- [ ] コマンド実行方法が記載されている
- [ ] その他の関連ドキュメントが更新されている

## 7. 影響範囲

### 7.1 変更が必要なファイル

#### ソースコード
- `server/cmd/generate-sample-data/main.go`: SQLite用ライブラリのインポートと処理分岐を削除（存在する場合）

#### テストコード
- `server/cmd/generate-sample-data/`配下のテストファイル: SQLite依存を確認し、必要に応じて修正（存在する場合）

#### ドキュメント
- `README.md`: `generate-sample-data`コマンドのPostgreSQL利用に関する記述を追加
- `docs/`配下の関連ドキュメント: SQLiteに関する記述をPostgreSQLに変更

### 7.2 既存ファイルの扱い
- `server/cmd/generate-sample-data/main.go`: SQLite用ライブラリのインポートと処理分岐を削除（存在する場合、PostgreSQLドライバーは既にサポート）
- `server/internal/db/connection.go`: 変更なし（既にPostgreSQLドライバーをサポート）
- `server/internal/config/config.go`: 変更なし（設定ファイル構造は維持）
- `config/{env}/database.yaml`: 変更なし（Issue #87でPostgreSQL設定が有効化されている想定）

## 8. 実装上の注意事項

### 8.1 設定ファイルの管理
- **設定ファイル**: `config/{env}/database.yaml`から接続情報を読み込む
- **環境変数**: `APP_ENV`環境変数で環境を指定（develop/staging/production）
- **設定構造**: 既存の`config/{env}/database.yaml`の構造を維持
- **PostgreSQL設定**: Issue #87で定義されたPostgreSQL構成に合わせる

### 8.2 SQLite用ライブラリと処理分岐の削除
- **インポート削除**: `server/cmd/generate-sample-data/main.go`からSQLite用ライブラリのインポートを削除（存在する場合）
- **処理分岐削除**: `server/cmd/generate-sample-data/main.go`内のSQLite用処理分岐を削除（存在する場合）

### 8.3 SQL構文の確認
- **プレースホルダー**: GORMのプレースホルダー構文（`?`）がPostgreSQLで正常に動作することを確認
- **動的テーブル名**: 動的テーブル名の使用がPostgreSQLで正常に動作することを確認
- **バッチ挿入**: バッチ挿入のSQL構文がPostgreSQL対応であることを確認

### 8.4 データ生成処理の確認
- **UUIDv7**: UUIDv7の生成とデータ型（`varchar(32)`）が正しく動作することを確認
- **テーブル番号計算**: テーブル番号の計算が正しく動作することを確認
- **バッチサイズ**: 既存のバッチサイズ（500件）を維持

### 8.5 動作確認
- **接続確認**: コマンド実行時にPostgreSQLへの接続を確認
- **データ生成確認**: 各テーブルにデータが正常に生成されることを確認
- **エラーハンドリング**: 接続エラー時のエラーハンドリングを確認

### 8.6 ドキュメント整備
- **起動手順**: PostgreSQLコンテナの起動・マイグレーション適用・コマンド実行の手順を記載
- **コマンド実行**: `generate-sample-data`コマンドの実行方法を記載
- **トラブルシューティング**: よくある問題と解決方法を記載

## 9. 参考情報

### 9.1 関連Issue
- GitHub Issue #85: 開発環境はPostgreSQLを利用する前提とする
- GitHub Issue #86: PostgreSQLの起動スクリプトと、Atlasマイグレーションスクリプトの修正
- GitHub Issue #87: APIサーバーの修正
- GitHub Issue #89: server/cmd/generate-sample-data コマンド

### 9.2 既存ドキュメント
- `README.md`: プロジェクト概要とセットアップ手順
- `docs/Architecture.md`: システムアーキテクチャ
- `docs/Initial-Setup.md`: 初期セットアップ手順
- `config/{env}/database.yaml`: 環境別データベース設定

### 9.3 技術スタック
- **PostgreSQL**: 15-alpine（Dockerイメージ）
- **GORM**: 既存のGORMライブラリ
- **PostgreSQLドライバー**: `gorm.io/driver/postgres`（既存）

### 9.4 参考リンク
- PostgreSQL公式ドキュメント: https://www.postgresql.org/docs/
- GORM公式ドキュメント: https://gorm.io/docs/
- GORM PostgreSQLドライバー: https://github.com/go-gorm/postgres
