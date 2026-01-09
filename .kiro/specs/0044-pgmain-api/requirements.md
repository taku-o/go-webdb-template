# APIサーバーPostgreSQL対応要件定義書

## 1. 概要

### 1.1 プロジェクト情報
- **プロジェクト名**: go-webdb-template
- **Issue番号**: #87
- **親Issue番号**: #85
- **Issueタイトル**: APIサーバーの修正
- **Feature名**: 0044-pgmain-api
- **作成日**: 2025-01-27

### 1.2 目的
APIサーバーでPostgreSQLを利用するように修正する。既存のSQLite設定をPostgreSQL設定に切り替え、開発環境、staging環境、production環境でPostgreSQLを利用できるようにする。

### 1.3 スコープ
- `config/develop/database.yaml`のPostgreSQL設定を有効化（SQLite設定を削除）
- `config/staging/database.yaml`のPostgreSQL設定を確認・修正（SQLite設定があれば削除）
- `config/production/database.yaml`のPostgreSQL設定を確認・修正（存在する場合、SQLite設定があれば削除）
- `config/production/database.yaml.example`のPostgreSQL設定を確認・修正（SQLite設定があれば削除）
- SQLite用ライブラリのインポートを削除
- ソースコード中のSQLite用処理分岐を削除
- テストコードのSQLite依存を確認し、必要に応じて修正
- ドキュメントの更新

**本実装の範囲外**:
- PostgreSQLの起動スクリプト・マイグレーションスクリプトの修正（Issue #86で対応）
- GoAdminサーバーのPostgreSQL接続設定変更（別Issue対応）
- `server/cmd/generate-sample-data`コマンドのPostgreSQL対応（別Issue対応）
- 本番環境へのデプロイ（準備のみ）

## 2. 背景・現状分析

### 2.1 現在の実装
- **設定ファイル**: `config/develop/database.yaml`でSQLite設定が有効、PostgreSQL設定がコメントアウトされている（または未定義）
- **データベース接続**: `server/internal/db/connection.go`で既にPostgreSQLドライバー（`gorm.io/driver/postgres`）がサポートされている
- **APIサーバー**: `server/cmd/server/main.go`で`config.Load()`で設定を読み込み、`db.NewGroupManager(cfg)`でデータベース接続を初期化
- **データベース構成**:
  - マスターデータベース: 1台（`master.db`）
  - シャーディングデータベース: 4台（`sharding_db_1.db` ～ `sharding_db_4.db`）
  - **論理シャーディング数: 8**（物理DB 4台 × 2論理シャード）
  - **バグ**: 現在のSQLite版の`config/develop/database.yaml`ではshardingグループに4つのデータベース設定（id: 1-4）しかないが、正しくは8つ（id: 1-8）必要

### 2.2 課題点
1. **設定ファイルのSQLite依存**: `config/develop/database.yaml`でSQLite設定が有効になっており、PostgreSQL設定がコメントアウトされている（または未定義）
2. **設定ファイルのバグ**: 現在のSQLite版の`config/develop/database.yaml`ではshardingグループに4つのデータベース設定（id: 1-4）しかないが、論理シャーディング数は8であるため、正しくは8つ（id: 1-8）必要
3. **環境別設定の未整備**: staging環境、production環境でのPostgreSQL設定が未確認・未整備
4. **SQLite用ライブラリの残存**: ソースコード中にSQLite用ライブラリのインポートが残っている
5. **SQLite用処理分岐の残存**: ソースコード中にSQLite用の処理分岐（`case "sqlite3":`など）が残っている
6. **テストコードのSQLite依存**: テストコードがSQLiteに依存している可能性がある
7. **ドキュメントの未更新**: PostgreSQL利用に関するドキュメントが未整備

### 2.3 本実装による改善点
1. **PostgreSQL環境への移行**: 開発環境、staging環境、production環境でPostgreSQLを利用できるようにする
2. **設定ファイルの整備**: 各環境の設定ファイルでPostgreSQL設定を有効化
3. **設定ファイルのバグ修正**: 論理シャーディング数8に合わせて、shardingグループに8つのデータベース設定を定義
4. **コードの簡素化**: SQLite用ライブラリと処理分岐を削除し、コードを簡素化
5. **テストコードの修正**: テストコードをPostgreSQL対応に修正
6. **ドキュメントの整備**: PostgreSQL利用に関するドキュメントを整備

## 3. 機能要件

### 3.1 設定ファイルの修正

#### 3.1.1 config/develop/database.yamlの修正
- **ファイル**: `config/develop/database.yaml`
- **変更内容**:
  - SQLite設定を削除（現在4つのshardingデータベース設定しかないが、これはバグ。正しくは8つ必要）
  - PostgreSQL設定のコメントアウトを解除（または新規追加）
  - PostgreSQL接続情報をIssue #86で定義された構成に合わせる
    - masterグループ: 1台（`postgres-master`、ポート5432、データベース名`webdb_master`）
    - shardingグループ: **論理シャーディング数8**（物理DB 4台、各物理DBに2つの論理シャード）
      - 物理データベース: 4台（`postgres-sharding-1` ～ `postgres-sharding-4`、ポート5433-5436、データベース名`webdb_sharding_1` ～ `webdb_sharding_4`）
      - 論理シャーディング設定: 8つ（id: 1-8）
        - id: 1, table_range: [0, 3] → postgres-sharding-1 (webdb_sharding_1)
        - id: 2, table_range: [4, 7] → postgres-sharding-1 (webdb_sharding_1)
        - id: 3, table_range: [8, 11] → postgres-sharding-2 (webdb_sharding_2)
        - id: 4, table_range: [12, 15] → postgres-sharding-2 (webdb_sharding_2)
        - id: 5, table_range: [16, 19] → postgres-sharding-3 (webdb_sharding_3)
        - id: 6, table_range: [20, 23] → postgres-sharding-3 (webdb_sharding_3)
        - id: 7, table_range: [24, 27] → postgres-sharding-4 (webdb_sharding_4)
        - id: 8, table_range: [28, 31] → postgres-sharding-4 (webdb_sharding_4)
  - 接続情報:
    - `host`: `localhost`（Docker環境の場合）
    - `port`: masterは`5432`、shardingは`5433`、`5434`、`5435`、`5436`
    - `user`: `webdb`
    - `password`: `webdb`
    - `name`: masterは`webdb_master`、shardingは`webdb_sharding_1` ～ `webdb_sharding_4`（物理DB名、論理シャードは同じ物理DBを参照）
    - `sslmode`: `disable`（開発環境）

#### 3.1.2 config/staging/database.yamlの確認・修正
- **ファイル**: `config/staging/database.yaml`
- **確認内容**:
  - PostgreSQL設定が正しく定義されているか確認
  - シャーディング構成が正しいか確認（物理DB 4台、論理シャーディング8）
  - shardingグループに8つのデータベース設定（id: 1-8）が定義されているか確認
  - 各論理シャードのtable_rangeが正しく設定されているか確認
  - 接続情報が適切か確認
  - SQLite設定が存在する場合は削除
- **修正内容**:
  - 不備があれば修正
  - SQLite設定があれば削除
  - 接続情報をIssue #86で定義された構成に合わせる
  - 論理シャーディング数が8つになっていることを確認

#### 3.1.3 config/production/database.yamlの確認・修正
- **ファイル**: `config/production/database.yaml`（存在する場合）
- **確認内容**:
  - PostgreSQL設定が正しく定義されているか確認
  - シャーディング構成が正しいか確認（物理DB 4台、論理シャーディング8）
  - shardingグループに8つのデータベース設定（id: 1-8）が定義されているか確認
  - 各論理シャードのtable_rangeが正しく設定されているか確認
  - 接続情報が適切か確認
  - SQLite設定が存在する場合は削除
- **修正内容**:
  - 不備があれば修正
  - SQLite設定があれば削除
  - 接続情報をIssue #86で定義された構成に合わせる
  - 論理シャーディング数が8つになっていることを確認

#### 3.1.4 config/production/database.yaml.exampleの確認・修正
- **ファイル**: `config/production/database.yaml.example`
- **確認内容**:
  - PostgreSQL設定が正しく定義されているか確認
  - シャーディング構成が正しいか確認（物理DB 4台、論理シャーディング8）
  - shardingグループに8つのデータベース設定（id: 1-8）が定義されているか確認
  - 各論理シャードのtable_rangeが正しく設定されているか確認
  - 接続情報が適切か確認
  - SQLite設定が存在する場合は削除
- **修正内容**:
  - 不備があれば修正
  - SQLite設定があれば削除
  - 接続情報をIssue #86で定義された構成に合わせる
  - 論理シャーディング数が8つになっていることを確認

### 3.2 SQLite用ライブラリと処理分岐の削除

#### 3.2.1 SQLite用ライブラリのインポート削除
- **対象ファイル**: 
  - `server/internal/db/connection.go`: `_ "github.com/mattn/go-sqlite3"`、`"gorm.io/driver/sqlite"`のインポートを削除
  - `server/cmd/admin/main.go`: `_ "github.com/GoAdminGroup/go-admin/modules/db/drivers/sqlite"`のインポートを削除
- **依存関係**: `server/go.mod`から`gorm.io/driver/sqlite`の依存関係を削除

#### 3.2.2 SQLite用処理分岐の削除
- **対象ファイル**: `server/internal/db/connection.go`
- **削除内容**:
  - `NewConnection`関数内の`driver = "sqlite3"`デフォルト値設定を削除
  - `createGORMConnection`関数内の`case "sqlite3":`分岐を削除
  - `createGORMConnectionFromDSN`関数内の`case "sqlite3":`分岐を削除
  - `NewGORMConnection`関数内の`case "sqlite3":`分岐を削除（Reader接続作成部分）

#### 3.2.3 エラーハンドリングの確認
- **確認内容**:
  - SQLite用処理分岐削除後、未サポートドライバー指定時のエラーハンドリングが適切か確認
  - エラーメッセージが適切か確認

### 3.3 テストコードの修正

#### 3.3.1 テストコードのSQLite依存確認
- **対象ファイル**: `server/test/`配下のテストファイル
- **確認内容**:
  - SQLite固有の設定やコードが含まれているか確認
  - テストデータベースの初期化方法を確認
  - テスト実行時のデータベース接続方法を確認

#### 3.3.2 テストコードのPostgreSQL対応
- **修正内容**:
  - SQLite固有の設定をPostgreSQL設定に変更
  - テストデータベースの初期化方法をPostgreSQL対応に変更
  - テスト実行時のデータベース接続方法をPostgreSQL対応に変更
  - テストユーティリティ（`server/test/testutil/db.go`など）をPostgreSQL対応に修正

### 3.4 ドキュメントの更新

#### 3.4.1 README.mdの更新
- **ファイル**: `README.md`
- **更新内容**:
  - PostgreSQL利用に関する記述を追加
  - 設定ファイルの変更方法を記載
  - 開発環境でのPostgreSQL起動手順を記載

#### 3.4.2 その他のドキュメントの更新
- **対象ファイル**: `docs/`配下の関連ドキュメント
- **更新内容**:
  - SQLiteに関する記述をPostgreSQLに変更
  - PostgreSQL利用に関する記述を追加

## 4. 非機能要件

### 4.1 PostgreSQL環境の前提条件
- Issue #86でPostgreSQLの起動スクリプトとマイグレーションスクリプトが修正されていること
- PostgreSQLコンテナが起動していること（開発環境の場合）
- マイグレーションが適用されていること

### 4.2 パフォーマンス
- **接続時間**: データベース接続の確立時間を最小化
- **クエリ実行時間**: PostgreSQLでのクエリ実行時間を最適化

### 4.3 セキュリティ
- **パスワード管理**: 開発環境では設定ファイルに固定パスワード（`webdb`）を記載する
- **SSL/TLS**: 開発環境では`sslmode=disable`、本番環境では適切なSSL設定を推奨

### 4.4 環境別対応
- **開発環境**: ローカルDocker環境でのPostgreSQL接続
- **staging環境**: staging環境用のPostgreSQL接続設定
- **production環境**: production環境用のPostgreSQL接続設定

## 5. 制約事項

### 5.1 既存システムとの関係
- **既存のデータベース接続コード**: `server/internal/db/connection.go`は既にPostgreSQLドライバーをサポートしているが、SQLite用の処理分岐を削除する必要がある
- **既存の設定ファイル構造**: `config/{env}/database.yaml`の構造は維持（SQLite設定の削除とPostgreSQL設定の追加のみ）
- **依存関係**: `server/go.mod`からSQLite関連の依存関係を削除する必要がある

### 5.2 技術スタック
- **PostgreSQL**: PostgreSQL 15-alpine（Dockerイメージ）
- **GORM**: 既存のGORMライブラリを使用
- **PostgreSQLドライバー**: `gorm.io/driver/postgres`（既存）

### 5.3 シャーディング構成
- **物理データベース数**: master 1台 + sharding 4台 = 合計5台
- **論理シャーディング数**: **8**（各物理DBに2つの論理シャードを割り当て）
- **論理シャーディング設定**: config/develop/database.yamlには8つのshardingデータベース設定（id: 1-8）が必要
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
- **起動順序**: PostgreSQLコンテナを起動してからAPIサーバーを起動する必要がある
- **マイグレーション**: マイグレーション適用後にAPIサーバーを起動する必要がある

## 6. 受け入れ基準

### 6.1 設定ファイルの修正
- [ ] `config/develop/database.yaml`でPostgreSQL設定が有効になっている（SQLite設定が削除されている）
- [ ] `config/develop/database.yaml`のPostgreSQL接続情報が正しく設定されている（master 1台、sharding論理シャーディング8つ、物理DB 4台）
- [ ] `config/develop/database.yaml`のshardingグループに8つのデータベース設定（id: 1-8）が定義されている
- [ ] `config/develop/database.yaml`の各論理シャードのtable_rangeが正しく設定されている（[0,3], [4,7], [8,11], [12,15], [16,19], [20,23], [24,27], [28,31]）
- [ ] `config/staging/database.yaml`のPostgreSQL設定が正しく定義されている（SQLite設定が削除されている、論理シャーディング8つ）
- [ ] `config/production/database.yaml`のPostgreSQL設定が正しく定義されている（存在する場合、SQLite設定が削除されている、論理シャーディング8つ）
- [ ] `config/production/database.yaml.example`のPostgreSQL設定が正しく定義されている（SQLite設定が削除されている、論理シャーディング8つ）

### 6.2 SQLite用ライブラリと処理分岐の削除
- [ ] `server/internal/db/connection.go`からSQLite用ライブラリのインポートが削除されている
- [ ] `server/cmd/admin/main.go`からSQLite用ライブラリのインポートが削除されている
- [ ] `server/go.mod`から`gorm.io/driver/sqlite`の依存関係が削除されている
- [ ] `server/internal/db/connection.go`からSQLite用処理分岐（`case "sqlite3":`など）が削除されている
- [ ] `server/internal/db/connection.go`から`driver = "sqlite3"`のデフォルト値設定が削除されている
- [ ] 未サポートドライバー指定時のエラーハンドリングが適切に実装されている

### 6.3 テストコードの修正
- [ ] テストコードのSQLite依存が確認され、必要に応じて修正されている
- [ ] テストユーティリティ（`server/test/testutil/db.go`など）がPostgreSQL対応に修正されている
- [ ] テストがPostgreSQL環境で正常に実行できる

### 6.4 APIサーバーの動作確認
- [ ] APIサーバーがPostgreSQLに正常に接続できる
- [ ] APIサーバーがmasterデータベースに正常に接続できる
- [ ] APIサーバーがshardingデータベースに正常に接続できる（4台すべて）
- [ ] APIサーバーが正常に起動し、リクエストを処理できる

### 6.5 ドキュメント
- [ ] `README.md`にPostgreSQL利用に関する記述が追加されている
- [ ] 設定ファイルの変更方法が記載されている
- [ ] 開発環境でのPostgreSQL起動手順が記載されている
- [ ] その他の関連ドキュメントが更新されている

## 7. 影響範囲

### 7.1 変更が必要なファイル

#### 設定ファイル
- `config/develop/database.yaml`: SQLite設定を削除、PostgreSQL設定を有効化
- `config/staging/database.yaml`: SQLite設定を削除（存在する場合）、PostgreSQL設定を確認・修正
- `config/production/database.yaml`: SQLite設定を削除（存在する場合）、PostgreSQL設定を確認・修正（存在する場合）
- `config/production/database.yaml.example`: SQLite設定を削除（存在する場合）、PostgreSQL設定を確認・修正

#### ソースコード
- `server/internal/db/connection.go`: SQLite用ライブラリのインポートと処理分岐を削除
- `server/cmd/admin/main.go`: SQLite用ライブラリのインポートを削除
- `server/go.mod`: SQLite関連の依存関係を削除

#### テストコード
- `server/test/testutil/db.go`: PostgreSQL対応に修正
- `server/test/`配下のその他のテストファイル: SQLite依存を確認し、必要に応じて修正

#### ドキュメント
- `README.md`: PostgreSQL利用に関する記述を追加
- `docs/`配下の関連ドキュメント: SQLiteに関する記述をPostgreSQLに変更

### 7.2 既存ファイルの扱い
- `server/internal/db/connection.go`: SQLite用ライブラリのインポートと処理分岐を削除（PostgreSQLドライバーは既にサポート）
- `server/cmd/admin/main.go`: SQLite用ライブラリのインポートを削除
- `server/go.mod`: SQLite関連の依存関係を削除
- `server/cmd/server/main.go`: 変更なし（設定ファイルから読み込むため）
- `server/internal/config/config.go`: 変更なし（設定ファイル構造は維持）

## 8. 実装上の注意事項

### 8.1 設定ファイルの管理
- **設定ファイル**: `config/{env}/database.yaml`から接続情報を読み込む
- **環境変数**: `APP_ENV`環境変数で環境を指定（develop/staging/production）
- **設定構造**: 既存の`config/{env}/database.yaml`の構造を維持
- **PostgreSQL設定**: Issue #86で定義されたPostgreSQL構成に合わせる

### 8.2 SQLite用ライブラリと処理分岐の削除
- **インポート削除**: `server/internal/db/connection.go`と`server/cmd/admin/main.go`からSQLite用ライブラリのインポートを削除
- **依存関係削除**: `server/go.mod`から`gorm.io/driver/sqlite`の依存関係を削除（`go mod tidy`を実行）
- **処理分岐削除**: `server/internal/db/connection.go`内のSQLite用処理分岐（`case "sqlite3":`など）を削除
- **デフォルト値削除**: `driver = "sqlite3"`のデフォルト値設定を削除
- **エラーハンドリング**: 未サポートドライバー指定時のエラーメッセージを確認

### 8.3 テストコードの修正
- **テストユーティリティ**: `server/test/testutil/db.go`をPostgreSQL対応に修正
- **テストデータベース**: テスト実行時にPostgreSQLデータベースを使用
- **テストデータ**: テストデータの投入方法をPostgreSQL対応に変更

### 8.4 ドキュメント整備
- **起動手順**: PostgreSQLコンテナの起動・マイグレーション適用・APIサーバー起動の手順を記載
- **設定ファイル**: `config/{env}/database.yaml`の設定方法を記載
- **トラブルシューティング**: よくある問題と解決方法を記載

### 8.5 動作確認
- **接続確認**: APIサーバー起動時にPostgreSQLへの接続を確認
- **クエリ実行**: APIサーバーが正常にクエリを実行できることを確認
- **エラーハンドリング**: 接続エラー時のエラーハンドリングを確認

## 9. 参考情報

### 9.1 関連Issue
- GitHub Issue #85: 開発環境はPostgreSQLを利用する前提とする
- GitHub Issue #86: PostgreSQLの起動スクリプトと、Atlasマイグレーションスクリプトの修正
- GitHub Issue #87: APIサーバーの修正

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
