# AdminサーバーのMySQL対応の要件定義書

## 1. 概要

### 1.1 プロジェクト情報
- **プロジェクト名**: go-webdb-template
- **Feature名**: 0055-admin-mysql
- **作成日**: 2026-01-10
- **関連Issue**: https://github.com/taku-o/go-webdb-template/issues/113

### 1.2 目的
Adminサーバー（GoAdmin管理画面）をPostgreSQLに加えて、MySQLでも動作するように修正する。デフォルトのデータベース自体はPostgreSQLとする。これにより、Adminサーバーでもデータベース選択の柔軟性を提供する。

### 1.3 スコープ
- Adminサーバーのデータベース接続設定のMySQL対応
- GoAdminのMySQLドライバーのインポート追加
- DSN生成ロジックのMySQL対応

**本実装の範囲外**:
- メインサーバー（APIサーバー）のMySQL対応（0054-mysqlで対応済み）
- PostgreSQLの既存機能への影響（PostgreSQLは引き続き動作する）
- データベース間の自動マイグレーション（手動でのマイグレーション実行を想定）
- 本番環境での自動切り替え機能（設定ファイルで手動切り替え）

## 2. 背景・現状分析

### 2.1 現在の状況
- **Adminサーバー**: GoAdminを使用した管理画面サーバー
- **データベース接続**: `server/internal/admin/config.go`の`getDatabaseConfig()`メソッドでPostgreSQL用のDSNを構築
- **ドライバー**: `github.com/GoAdminGroup/go-admin/modules/db/drivers/postgres`をインポート
- **DSN形式**: PostgreSQL用の`host=... port=... user=... password=... dbname=... sslmode=disable`形式をハードコード
- **ドライバー指定**: `Driver: "postgresql"`をハードコード
- **設定ファイル**: `config/{env}/config.yaml`の`DB_TYPE`設定を参照していない

### 2.2 課題点
1. **PostgreSQL専用のDSN生成**: `getDatabaseConfig()`メソッドでPostgreSQL用のDSN形式のみを生成
2. **MySQLドライバーの未インポート**: `main.go`でMySQLドライバーをインポートしていない
3. **ドライバー名のハードコード**: `Driver: "postgresql"`をハードコードしており、MySQLに対応していない
4. **設定ファイルの未参照**: `config.yaml`の`DB_TYPE`設定を参照していない

### 2.3 本実装による改善点
1. **データベース選択の柔軟性**: AdminサーバーでもPostgreSQLとMySQLの両方に対応
2. **設定ファイルの活用**: `config.yaml`の`DB_TYPE`設定に基づいて適切なデータベース接続を確立
3. **DSN生成の改善**: データベースタイプに応じて適切なDSN形式を生成

## 3. 機能要件

### 3.1 データベース接続設定（config.go）

#### 3.1.1 DSN生成ロジックの改善
- **目的**: データベースタイプに応じて適切なDSN形式を生成する
- **修正対象**: `server/internal/admin/config.go`の`getDatabaseConfig()`メソッド
- **実装内容**:
  - `c.appConfig.Database.Groups.Master[0].Driver`を参照してデータベースタイプを判定
  - PostgreSQLの場合: `host=... port=... user=... password=... dbname=... sslmode=disable`形式
  - MySQLの場合: `user:pass@tcp(host:port)/dbname?charset=utf8mb4&parseTime=true&loc=Local`形式
  - ドライバー名もデータベースタイプに応じて設定（`postgresql`または`mysql`）

#### 3.1.2 データベースタイプの取得
- **目的**: 設定ファイルからデータベースタイプを取得する
- **実装方法**: `c.appConfig.Database.Groups.Master[0].Driver`を参照
  - `postgres` → PostgreSQL
  - `mysql` → MySQL
- **エラーハンドリング**: ドライバーが指定されていない場合、または`postgres`/`mysql`以外の値の場合はエラーを返す

### 3.2 ドライバーのインポート（main.go）

#### 3.2.1 MySQLドライバーのインポート追加
- **目的**: GoAdminでMySQLを使用できるようにする
- **修正対象**: `server/cmd/admin/main.go`
- **実装内容**:
  - `_ "github.com/GoAdminGroup/go-admin/modules/db/drivers/mysql"`をインポート追加
  - 既存のPostgreSQLドライバーのインポートは維持

## 4. 非機能要件

### 4.1 パフォーマンス
- **接続プール**: 既存の接続プール設定を維持
- **クエリ性能**: MySQLでもPostgreSQLと同等の性能を維持（可能な限り）

### 4.2 信頼性
- **データ整合性**: MySQLでもPostgreSQLと同等のデータ整合性を保証
- **エラーハンドリング**: データベース固有のエラーを適切に処理
- **接続エラー**: データベース接続エラーを適切にハンドリング

### 4.3 保守性
- **コードの可読性**: データベース固有の処理を明確に分離
- **一貫性**: 既存のコードスタイルと一貫性を保つ
- **設定の一元化**: `config.yaml`の`DB_TYPE`設定を活用

### 4.4 互換性
- **既存機能**: PostgreSQLの既存機能に影響を与えない
- **後方互換性**: 既存のPostgreSQL設定ファイルは引き続き動作
- **設定ファイルの検証**: ドライバーが指定されていない場合はエラーとする

## 5. 制約事項

### 5.1 技術的制約
- **GoAdminドライバー**: `github.com/GoAdminGroup/go-admin/modules/db/drivers/mysql`を使用
- **データベース設定**: `config.yaml`の`DB_TYPE`設定と`database.yaml`（または`database.mysql.yaml`）の`driver`設定が一致している必要がある
- **DSN形式**: GoAdminが期待するDSN形式に準拠する必要がある

### 5.2 実装上の制約
- **設定ファイルの整合性**: `config.yaml`の`DB_TYPE`と`database.yaml`の`driver`が一致している必要がある
- **データベース接続**: メインサーバーと同じデータベース設定を使用する想定

### 5.3 動作環境
- **ローカル環境**: ローカル環境でPostgreSQLとMySQLの両方が動作することを確認
- **CI環境**: CI環境でもPostgreSQLとMySQLの両方が動作することを確認（該当する場合）

## 6. 受け入れ基準

### 6.1 データベース接続設定
- [ ] `server/internal/admin/config.go`の`getDatabaseConfig()`メソッドでデータベースタイプを判定している
- [ ] PostgreSQLの場合、PostgreSQL用のDSN形式を生成している
- [ ] MySQLの場合、MySQL用のDSN形式を生成している（`charset=utf8mb4&parseTime=true&loc=Local`を含む）
- [ ] ドライバー名がデータベースタイプに応じて設定されている（`postgresql`または`mysql`）
- [ ] ドライバーが指定されていない場合、または`postgres`/`mysql`以外の値の場合はエラーを返している

### 6.2 ドライバーのインポート
- [ ] `server/cmd/admin/main.go`にMySQLドライバーのインポートが追加されている
- [ ] 既存のPostgreSQLドライバーのインポートが維持されている

### 6.3 動作確認
- [ ] ローカル環境でPostgreSQL設定でAdminサーバーが正常に起動できる
- [ ] ローカル環境でMySQL設定でAdminサーバーが正常に起動できる
- [ ] PostgreSQL環境で管理画面が正常に表示される
- [ ] MySQL環境で管理画面が正常に表示される
- [ ] 既存のPostgreSQL機能が正常に動作することを確認
- [ ] CI環境でMySQL環境が正常に動作することを確認（該当する場合）

## 7. 影響範囲

### 7.1 変更が必要なファイル

#### 修正が必要なファイル
- `server/internal/admin/config.go`: `getDatabaseConfig()`メソッドの改善
- `server/cmd/admin/main.go`: MySQLドライバーのインポート追加

#### 確認が必要なファイル
- 既存の設定ファイル: 正常に動作することを確認

### 7.2 既存機能への影響
- **PostgreSQL機能**: 既存のPostgreSQL機能に影響を与えない（追加のみ）
- **既存の設定**: 既存のPostgreSQL設定ファイルは引き続き動作

## 8. 実装上の注意事項

### 8.1 DSN生成の実装
- **PostgreSQL DSN**: `host=... port=... user=... password=... dbname=... sslmode=disable`形式
- **MySQL DSN**: `user:pass@tcp(host:port)/dbname?charset=utf8mb4&parseTime=true&loc=Local`形式
- **文字セット**: MySQLでは`utf8mb4`を明示的に指定
- **タイムゾーン**: MySQLでは`loc=Local`を指定
- **後方互換性**: 既存のPostgreSQL DSN生成に影響を与えない

### 8.2 データベースタイプの判定
- **設定ファイルの参照**: `c.appConfig.Database.Groups.Master[0].Driver`を参照
- **エラーハンドリング**: 
  - ドライバーが指定されていない場合はエラーを返す
  - `postgres`/`mysql`以外の値の場合はエラーを返す
  - エラーメッセージは明確に設定ファイルの不備を指摘する

### 8.3 ドライバーのインポート
- **インポート順序**: 既存のインポートの後に追加
- **ブランクインポート**: `_`を使用してブランクインポートとする

## 9. 参考情報

### 9.1 関連ドキュメント
- `docs/Architecture.md`: アーキテクチャドキュメント
- `docs/Project-Structure.md`: プロジェクト構造ドキュメント
- `.kiro/specs/0054-mysql/requirements.md`: メインサーバーのMySQL対応要件定義書

### 9.2 既存実装の参考
- `server/internal/admin/config.go`: 既存のPostgreSQL用設定
- `server/cmd/admin/main.go`: 既存のエントリーポイント
- `server/internal/config/config.go`: メインサーバーのDSN生成ロジック（`GetDSN()`メソッド）

### 9.3 技術スタック
- **言語**: Go
- **データベース**: PostgreSQL, MySQL
- **管理画面**: GoAdmin（`github.com/GoAdminGroup/go-admin`）
- **ドライバー**: 
  - `github.com/GoAdminGroup/go-admin/modules/db/drivers/postgres`
  - `github.com/GoAdminGroup/go-admin/modules/db/drivers/mysql`

### 9.4 主なDSN形式の違い

| 項目 | PostgreSQL | MySQL |
|------|-----------|-------|
| DSN形式 | `host=... port=... user=... password=... dbname=... sslmode=disable` | `user:pass@tcp(host:port)/dbname?charset=utf8mb4&parseTime=true&loc=Local` |
| ドライバー名 | `postgresql` | `mysql` |
| 文字セット | 不要 | `charset=utf8mb4` |
| タイムゾーン | 不要 | `loc=Local` |
| 時刻解析 | 不要 | `parseTime=true` |
