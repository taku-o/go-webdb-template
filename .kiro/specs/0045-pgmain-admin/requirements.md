# GoAdminサーバーPostgreSQL対応要件定義書

## 1. 概要

### 1.1 プロジェクト情報
- **プロジェクト名**: go-webdb-template
- **Issue番号**: #88
- **親Issue番号**: #85
- **Issueタイトル**: GoAdminサーバーの修正
- **Feature名**: 0045-pgmain-admin
- **作成日**: 2025-01-27

### 1.2 目的
GoAdminサーバーでPostgreSQLを利用するように修正する。既存のSQLite設定をPostgreSQL設定に切り替え、開発環境、staging環境、production環境でPostgreSQLを利用できるようにする。

### 1.3 スコープ
- `server/internal/admin/config.go`の`getDatabaseConfig()`関数をPostgreSQL対応に修正（SQLite設定を削除）
- GoAdminのデータベース接続設定をPostgreSQLドライバーに変更
- 設定ファイル（`config/{env}/database.yaml`）からPostgreSQL接続情報を読み込む
- SQLite用ライブラリのインポートを削除（既にPostgreSQLドライバーはインポート済み）
- ドキュメントの更新

**本実装の範囲外**:
- PostgreSQLの起動スクリプト・マイグレーションスクリプトの修正（Issue #86で対応済み）
- APIサーバーのPostgreSQL接続設定変更（Issue #87で対応済み）
- `server/cmd/generate-sample-data`コマンドのPostgreSQL対応（別Issue対応）
- 本番環境へのデプロイ（準備のみ）

## 2. 背景・現状分析

### 2.1 現在の実装
- **設定ファイル**: `config/develop/database.yaml`でPostgreSQL設定が有効になっている（Issue #87で対応済み）
- **データベース接続**: `server/cmd/admin/main.go`で既にPostgreSQLドライバー（`github.com/GoAdminGroup/go-admin/modules/db/drivers/postgres`）がインポートされている
- **GoAdmin設定**: `server/internal/admin/config.go`の`getDatabaseConfig()`関数でSQLite設定（`Driver: "sqlite"`）がハードコードされている
- **設定読み込み**: `server/cmd/admin/main.go`で`config.Load()`で設定を読み込み、`appdb.NewGroupManager(cfg)`でデータベース接続を初期化しているが、GoAdminの設定では使用されていない
- **データベース構成**:
  - マスターデータベース: 1台（`webdb_master`）
  - シャーディングデータベース: 4台（`webdb_sharding_1` ～ `webdb_sharding_4`）
  - **論理シャーディング数: 8**（物理DB 4台 × 2論理シャード）
  - GoAdminはmasterデータベースのみを使用

### 2.2 課題点
1. **GoAdmin設定のSQLite依存**: `server/internal/admin/config.go`の`getDatabaseConfig()`関数でSQLite設定（`Driver: "sqlite"`）がハードコードされている
2. **設定ファイルの未活用**: `config.Load()`で読み込んだPostgreSQL設定がGoAdminの設定に反映されていない
3. **DSN形式の不一致**: GoAdminの設定で使用するDSN形式がPostgreSQL用の形式になっていない
4. **ドキュメントの未更新**: GoAdminサーバーのPostgreSQL利用に関するドキュメントが未整備

### 2.3 本実装による改善点
1. **PostgreSQL環境への移行**: GoAdminサーバーがPostgreSQLを利用できるようにする
2. **設定ファイルの活用**: `config/{env}/database.yaml`からPostgreSQL接続情報を読み込んでGoAdminの設定に反映
3. **コードの簡素化**: SQLite用設定を削除し、PostgreSQL設定のみに統一
4. **ドキュメントの整備**: GoAdminサーバーのPostgreSQL利用に関するドキュメントを整備

## 3. 機能要件

### 3.1 GoAdmin設定の修正

#### 3.1.1 server/internal/admin/config.goの修正
- **ファイル**: `server/internal/admin/config.go`
- **変更内容**:
  - `getDatabaseConfig()`関数をPostgreSQL対応に修正
  - SQLite設定（`Driver: "sqlite"`）を削除
  - PostgreSQL設定（`Driver: "postgres"`）を追加
  - `config.Load()`で読み込んだ設定（`c.appConfig.Database.Groups.Master[0]`）からPostgreSQL接続情報を取得
  - PostgreSQL用のDSN形式（`postgres://user:password@host:port/dbname?sslmode=disable`）を構築
  - 接続情報:
    - `host`: `c.appConfig.Database.Groups.Master[0].Host`
    - `port`: `c.appConfig.Database.Groups.Master[0].Port`
    - `user`: `c.appConfig.Database.Groups.Master[0].User`
    - `password`: `c.appConfig.Database.Groups.Master[0].Password`
    - `name`: `c.appConfig.Database.Groups.Master[0].Name`
    - `sslmode`: 開発環境では`disable`、本番環境では適切なSSL設定を推奨

#### 3.1.2 エラーハンドリングの追加
- **確認内容**:
  - masterグループのデータベース設定が存在しない場合のエラーハンドリング
  - 接続情報が不完全な場合のエラーハンドリング
  - エラーメッセージが適切か確認

### 3.2 SQLite用ライブラリの削除確認

#### 3.2.1 SQLite用ライブラリのインポート確認
- **対象ファイル**: 
  - `server/cmd/admin/main.go`: SQLite用ライブラリのインポートが存在しないことを確認（既にPostgreSQLドライバーのみインポートされている）
- **確認内容**:
  - `_ "github.com/GoAdminGroup/go-admin/modules/db/drivers/sqlite"`のインポートが存在しないことを確認
  - PostgreSQLドライバー（`_ "github.com/GoAdminGroup/go-admin/modules/db/drivers/postgres"`）がインポートされていることを確認

### 3.3 ドキュメントの更新

#### 3.3.1 README.mdの更新
- **ファイル**: `README.md`
- **更新内容**:
  - GoAdminサーバーのPostgreSQL利用に関する記述を追加
  - 設定ファイルの変更方法を記載
  - 開発環境でのPostgreSQL起動手順を記載

#### 3.3.2 その他のドキュメントの更新
- **対象ファイル**: `docs/`配下の関連ドキュメント
- **更新内容**:
  - GoAdminサーバーのSQLiteに関する記述をPostgreSQLに変更
  - GoAdminサーバーのPostgreSQL利用に関する記述を追加

## 4. 非機能要件

### 4.1 PostgreSQL環境の前提条件
- Issue #86でPostgreSQLの起動スクリプトとマイグレーションスクリプトが修正されていること
- Issue #87でAPIサーバーのPostgreSQL接続設定が修正されていること
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
- **既存のデータベース接続コード**: `server/cmd/admin/main.go`で既にPostgreSQLドライバーがインポートされている
- **既存の設定ファイル構造**: `config/{env}/database.yaml`の構造は維持（Issue #87でPostgreSQL設定が有効になっている）
- **GoAdminの設定**: `server/internal/admin/config.go`の`getDatabaseConfig()`関数のみを修正

### 5.2 技術スタック
- **PostgreSQL**: PostgreSQL 15-alpine（Dockerイメージ）
- **GoAdmin**: 既存のGoAdminライブラリを使用
- **PostgreSQLドライバー**: `github.com/GoAdminGroup/go-admin/modules/db/drivers/postgres`（既存）

### 5.3 シャーディング構成
- **物理データベース数**: master 1台 + sharding 4台 = 合計5台
- **論理シャーディング数**: **8**（各物理DBに2つの論理シャードを割り当て）
- **GoAdminの使用データベース**: masterデータベースのみ（`webdb_master`）

### 5.4 運用上の制約
- **起動順序**: PostgreSQLコンテナを起動してからGoAdminサーバーを起動する必要がある
- **マイグレーション**: マイグレーション適用後にGoAdminサーバーを起動する必要がある

## 6. 受け入れ基準

### 6.1 GoAdmin設定の修正
- [ ] `server/internal/admin/config.go`の`getDatabaseConfig()`関数でPostgreSQL設定（`Driver: "postgres"`）が使用されている（SQLite設定が削除されている）
- [ ] `server/internal/admin/config.go`の`getDatabaseConfig()`関数で`config.Load()`で読み込んだPostgreSQL接続情報が使用されている
- [ ] `server/internal/admin/config.go`の`getDatabaseConfig()`関数でPostgreSQL用のDSN形式が正しく構築されている
- [ ] masterグループのデータベース設定が存在しない場合のエラーハンドリングが適切に実装されている
- [ ] 接続情報が不完全な場合のエラーハンドリングが適切に実装されている

### 6.2 SQLite用ライブラリの削除確認
- [ ] `server/cmd/admin/main.go`にSQLite用ライブラリのインポートが存在しないことを確認
- [ ] `server/cmd/admin/main.go`にPostgreSQLドライバー（`_ "github.com/GoAdminGroup/go-admin/modules/db/drivers/postgres"`）がインポートされていることを確認

### 6.3 GoAdminサーバーの動作確認
- [ ] GoAdminサーバーがPostgreSQLに正常に接続できる
- [ ] GoAdminサーバーがmasterデータベース（`webdb_master`）に正常に接続できる
- [ ] GoAdminサーバーが正常に起動し、管理画面にアクセスできる
- [ ] GoAdmin管理画面でデータベースのテーブルが正常に表示できる

### 6.4 ドキュメント
- [ ] `README.md`にGoAdminサーバーのPostgreSQL利用に関する記述が追加されている
- [ ] 設定ファイルの変更方法が記載されている
- [ ] 開発環境でのPostgreSQL起動手順が記載されている
- [ ] その他の関連ドキュメントが更新されている

## 7. 影響範囲

### 7.1 変更が必要なファイル

#### ソースコード
- `server/internal/admin/config.go`: `getDatabaseConfig()`関数をPostgreSQL対応に修正

#### ドキュメント
- `README.md`: GoAdminサーバーのPostgreSQL利用に関する記述を追加
- `docs/`配下の関連ドキュメント: GoAdminサーバーのSQLiteに関する記述をPostgreSQLに変更

### 7.2 既存ファイルの扱い
- `server/internal/admin/config.go`: `getDatabaseConfig()`関数のみを修正（他の関数は変更なし）
- `server/cmd/admin/main.go`: 変更なし（既にPostgreSQLドライバーがインポートされている）
- `server/internal/config/config.go`: 変更なし（設定ファイル構造は維持）
- `config/{env}/database.yaml`: 変更なし（Issue #87でPostgreSQL設定が有効になっている）

## 8. 実装上の注意事項

### 8.1 GoAdmin設定の修正
- **設定ファイル**: `config/{env}/database.yaml`から接続情報を読み込む
- **環境変数**: `APP_ENV`環境変数で環境を指定（develop/staging/production）
- **設定構造**: 既存の`config/{env}/database.yaml`の構造を維持（Issue #87でPostgreSQL設定が有効になっている）
- **PostgreSQL設定**: Issue #86で定義されたPostgreSQL構成に合わせる
- **DSN形式**: PostgreSQL用のDSN形式（`postgres://user:password@host:port/dbname?sslmode=disable`）を使用

### 8.2 エラーハンドリング
- **masterグループの存在確認**: `c.appConfig.Database.Groups.Master`が空でないことを確認
- **接続情報の完全性確認**: 必要な接続情報（host, port, user, password, name）がすべて存在することを確認
- **エラーメッセージ**: 適切なエラーメッセージを返す

### 8.3 ドキュメント整備
- **起動手順**: PostgreSQLコンテナの起動・マイグレーション適用・GoAdminサーバー起動の手順を記載
- **設定ファイル**: `config/{env}/database.yaml`の設定方法を記載
- **トラブルシューティング**: よくある問題と解決方法を記載

### 8.4 動作確認
- **接続確認**: GoAdminサーバー起動時にPostgreSQLへの接続を確認
- **クエリ実行**: GoAdmin管理画面でデータベースのテーブルが正常に表示できることを確認
- **エラーハンドリング**: 接続エラー時のエラーハンドリングを確認

## 9. 参考情報

### 9.1 関連Issue
- GitHub Issue #85: 開発環境はPostgreSQLを利用する前提とする
- GitHub Issue #86: PostgreSQLの起動スクリプトと、Atlasマイグレーションスクリプトの修正
- GitHub Issue #87: APIサーバーの修正
- GitHub Issue #88: GoAdminサーバーの修正

### 9.2 既存ドキュメント
- `README.md`: プロジェクト概要とセットアップ手順
- `docs/Architecture.md`: システムアーキテクチャ
- `docs/Initial-Setup.md`: 初期セットアップ手順
- `config/{env}/database.yaml`: 環境別データベース設定

### 9.3 技術スタック
- **PostgreSQL**: 15-alpine（Dockerイメージ）
- **GoAdmin**: 既存のGoAdminライブラリ
- **PostgreSQLドライバー**: `github.com/GoAdminGroup/go-admin/modules/db/drivers/postgres`

### 9.4 参考リンク
- PostgreSQL公式ドキュメント: https://www.postgresql.org/docs/
- GoAdmin公式ドキュメント: https://github.com/GoAdminGroup/go-admin
- GoAdmin PostgreSQLドライバー: https://github.com/GoAdminGroup/go-admin/tree/master/modules/db/drivers/postgres
