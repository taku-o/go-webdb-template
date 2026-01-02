# データベース遅延接続・自動再接続機能要件定義書

## 1. 概要

### 1.1 プロジェクト情報
- **プロジェクト名**: go-webdb-template
- **Issue番号**: #72
- **Issueタイトル**: データベースに遅延接続＋データベースダウン後の再開時に再接続するように修正したい
- **Feature名**: 0036-db-reconnect
- **作成日**: 2025-01-27

### 1.2 目的
`database/sql`（GORMのベース）の標準機能を活用して、データベースへの遅延接続とダウン後の自動再接続を実現する。
これにより、APIサーバー起動時にデータベースが利用できない場合でもサーバーを起動でき、データベースが復旧した際に自動的に接続が確立されるようにする。

### 1.3 スコープ
- APIサーバー起動時のDB接続確認処理の削除（遅延接続の実現）
- 接続プール設定の確認と追加（`SetMaxIdleConns`, `SetMaxOpenConns`, `SetConnMaxLifetime`）
- 遅延接続と自動再接続の動作確認（PostgreSQLを使用）
- DB接続エラー時のリトライ機能の実装（`avast/retry-go`を使用）
- PostgreSQL環境の構築（Dockerを使用）

**本実装の範囲外**:
- 本番環境での本格的な接続管理機能
- 接続監視・アラート機能
- 接続プールの動的な調整機能
- 複数データベースへの接続管理の最適化

## 2. 背景・現状分析

### 2.1 現在の実装
- **起動時の接続確認**: `server/cmd/server/main.go`の38-40行目で`groupManager.PingAll()`を実行し、失敗時にサーバーを終了
- **接続作成時のPing**: `server/internal/db/connection.go`の55行目で`NewConnection`内で`db.Ping()`を実行し、失敗時にエラーを返す
- **接続プール設定**: `connection.go`の50-52行目で接続プール設定が実装されているが、設定値が適切か確認が必要
- **GORM接続**: `NewGORMConnection`では接続プール設定は実装されているが、Pingは実行されていない
- **データベース**: 現在はSQLiteを主に使用（PostgreSQLの動作確認が困難）

### 2.2 課題点
1. **起動時の接続必須**: APIサーバー起動時にDB接続できない場合、サーバーが終了してしまう
2. **遅延接続の未活用**: `database/sql`の標準機能である遅延接続が活用されていない
3. **自動再接続の未確認**: データベースダウン後の再開時の自動再接続が動作するか確認されていない
4. **接続プール設定の不備**: 接続プール設定が適切に設定されていない可能性
5. **接続リトライ機能の不足**: DB接続エラー時のリトライ機能が実装されていない
6. **動作確認環境の不足**: SQLiteでは遅延接続・自動再接続の動作確認が困難

### 2.3 本実装による改善点
1. **起動時の柔軟性向上**: DB接続できない場合でもサーバーを起動できる
2. **遅延接続の実現**: 実際にクエリを実行する際に接続が確立される
3. **自動再接続の実現**: データベースが復旧した際に自動的に接続が確立される
4. **接続プールの最適化**: 適切な接続プール設定により、パフォーマンスと安定性が向上
5. **接続リトライ機能**: DB接続エラー時に自動的にリトライする機能を実装
6. **動作確認環境の整備**: PostgreSQL環境を構築し、動作確認が可能になる

## 3. 機能要件

### 3.1 起動時のDB接続確認処理の削除

#### 3.1.1 起動時のPing削除
- `server/cmd/server/main.go`の38-40行目の`groupManager.PingAll()`呼び出しを削除
- 起動時にDB接続できない場合でもサーバーを起動できるようにする
- 警告ログを出力して、DB接続が利用できない状態であることを通知

#### 3.1.2 接続作成時のPing削除
- `server/internal/db/connection.go`の`NewConnection`関数内の`db.Ping()`呼び出し（55行目）を削除
- `sql.Open()`は接続を確立せず、接続オブジェクトを作成するのみ
- 実際のクエリ実行時に接続が確立される（遅延接続）

### 3.2 接続プール設定の確認と追加

#### 3.2.1 接続プール設定の確認
- `server/internal/db/connection.go`の`NewConnection`関数で接続プール設定が実装されていることを確認
- `server/internal/db/connection.go`の`createGORMConnection`関数で接続プール設定が実装されていることを確認
- 設定値が適切か確認し、必要に応じて調整

#### 3.2.2 接続プール設定の追加
以下の設定が実装されていることを確認し、未実装の場合は追加する：
- `sqlDB.SetMaxIdleConns()`: アイドル状態のコネクション数
- `sqlDB.SetMaxOpenConns()`: 最大同時接続数
- `sqlDB.SetConnMaxLifetime()`: 古くなった接続を安全に破棄して再作成させる

#### 3.2.3 設定値の確認
- 設定ファイル（`config/{env}/database.yaml`）から接続プール設定を読み込む
- デフォルト値が適切に設定されていることを確認
- 設定値が0以下の場合は適切なデフォルト値を設定

### 3.3 PostgreSQL環境の構築

#### 3.3.1 Docker Compose設定
- `docker-compose.postgres.yml`ファイルを作成（開発用途）
- PostgreSQLコンテナの定義
- ネットワーク設定
- ボリューム設定（データ永続化用）
- デフォルトポート: 5432

#### 3.3.2 起動スクリプトの実装
- `scripts/start-postgres.sh`を作成（開発用途）
- Docker Composeを使用してPostgreSQLを起動
- 既存の起動スクリプト（`start-mailpit.sh`など）と同じパターンで実装

#### 3.3.3 設定ファイルの追加
- PostgreSQL用の設定ファイル（`config/{env}/database.yaml`）を追加
- PostgreSQL接続情報（DSN、ユーザー名、パスワードなど）を設定
- 既存のSQLite設定と併用できるようにする

### 3.4 遅延接続と自動再接続の動作確認

#### 3.4.1 遅延接続の確認
- PostgreSQL環境でAPIサーバーを起動（DB接続なしで起動可能であることを確認）
- 最初のクエリ実行時に接続が確立されることを確認
- ログで接続確立のタイミングを確認

#### 3.4.2 自動再接続の確認
- PostgreSQLを停止して、クエリ実行時にエラーが発生することを確認
- PostgreSQLを再起動して、次のクエリ実行時に自動的に再接続されることを確認
- ログで再接続のタイミングを確認

### 3.5 DB接続エラー時のリトライ機能

#### 3.5.1 リトライライブラリの導入
- `github.com/avast/retry-go/v4`ライブラリを導入
- 接続エラー時にリトライを実行する機能を実装

#### 3.5.2 リトライ設定
- リトライ回数: 最大3回（初回 + 2回のリトライ）
- リトライ間隔: 1秒待機
- リトライ対象: DB接続エラー（`sql.ErrConnDone`など）

#### 3.5.3 リトライ実装箇所
- `server/internal/db/connection.go`の接続作成処理
- `server/internal/db/manager.go`の接続作成処理
- クエリ実行時の接続エラー処理（必要に応じて）

## 4. 非機能要件

### 4.1 パフォーマンス
- 接続プール設定により、接続の再利用が適切に行われること
- リトライ処理が過度にパフォーマンスに影響を与えないこと

### 4.2 可用性
- DB接続できない場合でもサーバーを起動できること
- データベースが復旧した際に自動的に再接続されること
- 接続エラー時に適切なリトライが実行されること

### 4.3 保守性
- 既存のDocker Compose設定パターンに従うこと
- 既存の起動スクリプトパターンに従うこと
- 接続プール設定が設定ファイルから読み込めること

### 4.4 拡張性
- 将来的に接続監視・アラート機能を追加できる設計であること
- 接続プール設定を動的に調整できる設計であること（将来の拡張項目）

## 5. 技術仕様

### 5.1 サーバー側技術スタック
- **言語**: Go 1.21+
- **データベース**: SQLite（既存）、PostgreSQL（動作確認用）
- **リトライライブラリ**: `github.com/avast/retry-go/v4`
- **Webフレームワーク**: 既存のEcho + Huma API

### 5.2 インフラストラクチャ
- **PostgreSQL**: Docker Composeで起動
- **Docker Compose**: 既存のパターンに従う

### 5.3 ファイル構造
- **Docker Compose設定**: `docker-compose.postgres.yml`（新規作成）
- **起動スクリプト**: `scripts/start-postgres.sh`（新規作成）
- **DB接続処理**: `server/internal/db/connection.go`（修正）
- **DB Manager**: `server/internal/db/manager.go`（修正）
- **サーバー起動**: `server/cmd/server/main.go`（修正）
- **設定**: `server/internal/config/config.go`（確認・修正）

## 6. 受け入れ基準

### 6.1 機能要件
1. **起動時のDB接続確認処理の削除**: DB接続できない場合でもサーバーを起動できること
2. **接続プール設定の確認と追加**: 接続プール設定が適切に実装されていること
3. **PostgreSQL環境の構築**: PostgreSQLがDocker Composeで起動できること
4. **遅延接続の確認**: 最初のクエリ実行時に接続が確立されること
5. **自動再接続の確認**: データベースが復旧した際に自動的に再接続されること
6. **接続リトライ機能**: DB接続エラー時に1秒待機して最大3回までリトライすること

### 6.2 非機能要件
1. **パフォーマンス**: 接続プール設定により、接続の再利用が適切に行われること
2. **可用性**: DB接続できない場合でもサーバーを起動できること
3. **保守性**: 既存のパターンに従った実装であること
4. **拡張性**: 将来的に接続監視・アラート機能を追加できる設計であること

## 7. 制約事項

1. **動作確認環境**: SQLiteでは遅延接続・自動再接続の動作確認が困難なため、PostgreSQLを使用する
2. **開発環境での利用**: 本実装は主に開発環境での利用を想定（本番環境での利用は将来の拡張項目）
3. **既存機能の保持**: 既存のDB接続機能を壊さないこと
4. **後方互換性**: 既存の設定ファイルとの互換性を保つこと

## 8. 将来の拡張項目（現時点では未実装）

以下の機能は将来の拡張として検討されていますが、現時点では実装対象外です：

- 接続監視・アラート機能
- 接続プールの動的な調整機能
- 複数データベースへの接続管理の最適化
- 本番環境での本格的な接続管理機能

## Project Description (Input)

https://github.com/taku-o/go-webdb-template/issues/72 に対応するための要件を作成してください。

database/sql（GORMのベース）は実は正しく設定すると、標準で「遅延接続」と「DBダウン後の再開時の自動再接続」の機能を持っているらしい。

1. APIサーバー起動時にDBに接続できない場合、サーバーが終了する実装になっている。この処理を止める。
2. 遅延接続と、ダウン後の再開時の自動接続が動作するか確認したい。
   * SQLiteでは確認が難しいか？PostgreSQLを一時的に使うことが考えたい。PostgreSQLを使う場合はDockerで導入する。
3. これらの設定が入っていなければ入れる。
   ```go
   sqlDB.SetMaxIdleConns()    // アイドル状態のコネクション数
   sqlDB.SetMaxOpenConns()    // 最大同時接続数
   sqlDB.SetConnMaxLifetime() // 古くなった接続を安全に破棄して再作成させる
   ```
4. DB接続エラー時に接続リトライしたい。
   * avast/retry-go 辺りを利用して。
   * 1秒待機して、追加で2回までDB接続に挑戦する。

## Requirements

### Requirement 1: 起動時のDB接続確認処理の削除
**Objective:** As a system administrator, I want the API server to start even when the database is unavailable, so that the server can be started before the database is ready and automatically connect when the database becomes available.

#### Acceptance Criteria
1. WHEN the API server starts THEN it SHALL NOT require database connection to be available
2. IF database connection is unavailable at startup THEN the server SHALL start successfully with a warning log
3. WHERE database connection fails at startup THEN the server SHALL log a warning message indicating database connection is unavailable
4. WHEN the first database query is executed THEN the connection SHALL be established automatically (lazy connection)
5. IF database connection is unavailable at startup THEN the server SHALL continue to run and wait for database to become available

### Requirement 2: 接続作成時のPing削除
**Objective:** As a system, I want database connections to be created without immediate ping, so that lazy connection can be properly implemented.

#### Acceptance Criteria
1. WHEN NewConnection is called THEN it SHALL NOT execute db.Ping() immediately
2. IF sql.Open() is called THEN it SHALL create a connection object without establishing actual connection
3. WHERE connection is created THEN the actual connection SHALL be established when the first query is executed
4. WHEN NewGORMConnection is called THEN it SHALL NOT execute ping immediately
5. IF GORM connection is created THEN the actual connection SHALL be established when the first query is executed

### Requirement 3: 接続プール設定の確認と追加
**Objective:** As a system, I want proper connection pool settings to be configured, so that database connections can be efficiently managed and reused.

#### Acceptance Criteria
1. WHEN connection pool settings are checked THEN they SHALL be implemented in NewConnection function
2. IF connection pool settings are checked THEN they SHALL be implemented in createGORMConnection function
3. WHERE SetMaxIdleConns is called THEN it SHALL set the maximum number of idle connections
4. WHEN SetMaxOpenConns is called THEN it SHALL set the maximum number of open connections
5. IF SetConnMaxLifetime is called THEN it SHALL set the maximum lifetime of connections
6. WHERE connection pool settings are read THEN they SHALL be read from config file (config/{env}/database.yaml)
7. WHEN connection pool settings are not configured THEN default values SHALL be used
8. IF connection pool settings are 0 or negative THEN appropriate default values SHALL be set

### Requirement 4: PostgreSQL環境の構築
**Objective:** As a developer, I want PostgreSQL environment to be set up for testing lazy connection and auto-reconnection, so that I can verify these features work correctly.

#### Acceptance Criteria
1. WHEN docker-compose.postgres.yml is created THEN it SHALL define PostgreSQL service
2. IF PostgreSQL service is defined THEN it SHALL use appropriate Docker image
3. WHERE PostgreSQL is configured THEN it SHALL listen on port 5432 by default
4. WHEN PostgreSQL volumes are defined THEN they SHALL persist data across container restarts
5. IF start-postgres.sh is created THEN it SHALL start PostgreSQL using Docker Compose
6. WHERE start-postgres.sh is executed THEN it SHALL follow existing script patterns
7. WHEN PostgreSQL configuration is added THEN it SHALL be added to config/{env}/database.yaml
8. IF PostgreSQL configuration is added THEN it SHALL coexist with existing SQLite configuration

### Requirement 5: 遅延接続の動作確認
**Objective:** As a developer, I want to verify that lazy connection works correctly, so that I can confirm connections are established only when needed.

#### Acceptance Criteria
1. WHEN API server starts without database connection THEN it SHALL start successfully
2. IF first database query is executed THEN the connection SHALL be established automatically
3. WHERE connection is established THEN it SHALL be logged appropriately
4. WHEN lazy connection works THEN the connection timing SHALL be verified in logs

### Requirement 6: 自動再接続の動作確認
**Objective:** As a developer, I want to verify that auto-reconnection works correctly, so that I can confirm connections are automatically re-established when database becomes available.

#### Acceptance Criteria
1. WHEN database is stopped THEN queries SHALL fail with appropriate errors
2. IF database is restarted THEN the next query SHALL automatically re-establish connection
3. WHERE auto-reconnection works THEN the reconnection timing SHALL be verified in logs
4. WHEN database becomes available THEN connections SHALL be automatically re-established

### Requirement 7: DB接続エラー時のリトライ機能
**Objective:** As a system, I want database connection errors to be retried automatically, so that temporary connection failures can be recovered without manual intervention.

#### Acceptance Criteria
1. WHEN retry library is integrated THEN it SHALL use github.com/avast/retry-go/v4
2. IF database connection fails THEN it SHALL retry up to 3 times (initial attempt + 2 retries)
3. WHERE retry is executed THEN it SHALL wait 1 second between retries
4. WHEN retry is implemented THEN it SHALL be implemented in connection creation process
5. IF retry is implemented THEN it SHALL be implemented in manager connection creation process
6. WHERE retry is executed THEN it SHALL handle connection errors (sql.ErrConnDone, etc.)
7. WHEN all retries fail THEN it SHALL return appropriate error message
