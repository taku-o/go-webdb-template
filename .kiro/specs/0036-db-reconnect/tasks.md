# データベース遅延接続・自動再接続機能実装タスク一覧

## 概要
データベース遅延接続・自動再接続機能の実装を段階的に進めるためのタスク一覧。要件定義書と設計書に基づき、実装可能な粒度でタスクを分解。

## 実装フェーズ

### Phase 1: 起動時のDB接続確認処理の削除

#### - [x] タスク 1.1: server/cmd/server/main.goのPingAll()呼び出し削除
**目的**: APIサーバー起動時にDB接続確認を行わないようにする

**作業内容**:
- `server/cmd/server/main.go`を開く
- 38-40行目の`groupManager.PingAll()`呼び出しを削除
- 削除箇所に警告ログを追加:
  ```go
  // 起動時のDB接続確認は削除（遅延接続のため）
  // 最初のクエリ実行時に接続が確立される
  log.Println("Database connections will be established on first query execution (lazy connection)")
  ```
- 既存の`log.Println("Successfully connected to all database groups")`を削除

**受け入れ基準**:
- `PingAll()`呼び出しが削除されている
- 警告ログが追加されている
- DB接続できない場合でもサーバーが起動できること
- 既存のコードスタイルに従っている

---

### Phase 2: 接続作成時のPing削除

#### - [x] タスク 2.1: server/internal/db/connection.goのNewConnection関数のPing削除
**目的**: 接続作成時にPingを実行せず、遅延接続を実現する

**作業内容**:
- `server/internal/db/connection.go`を開く
- `NewConnection`関数内の55行目付近の`db.Ping()`呼び出しを削除
- 削除箇所にコメントを追加:
  ```go
  // 接続確認は削除（遅延接続のため）
  // sql.Open()は接続を確立せず、接続オブジェクトを作成するのみ
  // 実際のクエリ実行時に接続が確立される
  ```
- `db.Close()`の呼び出しも削除（Ping失敗時の処理のため）

**受け入れ基準**:
- `db.Ping()`呼び出しが削除されている
- 適切なコメントが追加されている
- `sql.Open()`はエラーを返さない（接続を確立しないため）
- 既存のコードスタイルに従っている

---

### Phase 3: 接続プール設定の確認と追加

#### - [x] タスク 3.1: 接続プール設定の実装状況確認
**目的**: 既存の接続プール設定が適切に実装されているか確認する

**作業内容**:
- `server/internal/db/connection.go`の`NewConnection`関数を確認
- `server/internal/db/connection.go`の`createGORMConnection`関数を確認
- 以下の設定が実装されているか確認:
  - `SetMaxIdleConns()`の呼び出し
  - `SetMaxOpenConns()`の呼び出し
  - `SetConnMaxLifetime()`の呼び出し
- 設定値が設定ファイルから読み込まれているか確認

**受け入れ基準**:
- 接続プール設定が実装されていることを確認
- 設定値が設定ファイルから読み込まれていることを確認
- 確認結果を記録

---

#### - [x] タスク 3.2: 接続プール設定のデフォルト値追加
**目的**: 設定値が0以下の場合に適切なデフォルト値を設定する

**作業内容**:
- `server/internal/db/connection.go`を開く
- デフォルト値の定数を追加:
  ```go
  const (
      DefaultMaxConnections        = 25
      DefaultMaxIdleConnections    = 5
      DefaultConnectionMaxLifetime = 1 * time.Hour
  )
  ```
- `NewConnection`関数で設定値が0以下の場合にデフォルト値を使用するように修正:
  ```go
  maxConnections := cfg.MaxConnections
  if maxConnections <= 0 {
      maxConnections = DefaultMaxConnections
  }
  db.SetMaxOpenConns(maxConnections)
  
  maxIdleConnections := cfg.MaxIdleConnections
  if maxIdleConnections <= 0 {
      maxIdleConnections = DefaultMaxIdleConnections
  }
  db.SetMaxIdleConns(maxIdleConnections)
  
  connectionMaxLifetime := cfg.ConnectionMaxLifetime
  if connectionMaxLifetime <= 0 {
      connectionMaxLifetime = DefaultConnectionMaxLifetime
  }
  db.SetConnMaxLifetime(connectionMaxLifetime)
  ```
- `createGORMConnection`関数にも同様の処理を追加

**受け入れ基準**:
- デフォルト値の定数が定義されている
- 設定値が0以下の場合にデフォルト値が使用される
- `NewConnection`と`createGORMConnection`の両方で実装されている
- 既存のコードスタイルに従っている

---

### Phase 4: PostgreSQL環境の構築

#### - [x] タスク 4.1: docker-compose.postgres.ymlの作成
**目的**: PostgreSQLをDocker Composeで起動するための設定ファイルを作成する

**作業内容**:
- `docker-compose.postgres.yml`ファイルを新規作成
- PostgreSQL 15 Alpineイメージを使用
- ポート5432を公開
- 環境変数を設定:
  - `POSTGRES_USER: webdb`
  - `POSTGRES_PASSWORD: webdb`
  - `POSTGRES_DB: webdb`
- ボリューム設定:
  - `./postgres/data:/var/lib/postgresql/data`
- ヘルスチェックを追加:
  ```yaml
  healthcheck:
    test: ["CMD-SHELL", "pg_isready -U webdb"]
    interval: 10s
    timeout: 5s
    retries: 5
  ```
- `restart: unless-stopped`を設定

**受け入れ基準**:
- `docker-compose.postgres.yml`ファイルが作成されている
- PostgreSQLコンテナが起動できること
- データが`postgres/data`ディレクトリに保存されること
- コンテナ再起動後もデータが保持されること
- ヘルスチェックが正常に動作すること

---

#### - [x] タスク 4.2: scripts/start-postgres.shの作成
**目的**: PostgreSQLを起動・停止するスクリプトを作成する

**作業内容**:
- `scripts/start-postgres.sh`ファイルを新規作成
- 実行権限を付与（`chmod +x scripts/start-postgres.sh`）
- 既存の`start-mailpit.sh`と同じパターンで実装
- `start`コマンド: Docker ComposeでPostgreSQLを起動
- `stop`コマンド: Docker ComposeでPostgreSQLを停止
- 適切なフィードバックメッセージを出力（接続情報を含む）

**受け入れ基準**:
- `scripts/start-postgres.sh`ファイルが作成されている
- 実行権限が付与されている
- `./scripts/start-postgres.sh start`でPostgreSQLが起動できること
- `./scripts/start-postgres.sh stop`でPostgreSQLが停止できること
- 既存の起動スクリプトと同じパターンに従っている

---

#### - [x] タスク 4.3: PostgreSQL設定ファイルの追加
**目的**: PostgreSQL用の設定ファイルを追加する

**作業内容**:
- `config/develop/database.yaml`を開く（または新規作成）
- PostgreSQL用の設定を追加:
  ```yaml
  database:
    groups:
      master:
        - id: 1
          driver: "postgres"
          host: "localhost"
          port: 5432
          user: "webdb"
          password: "webdb"
          name: "webdb"
          max_connections: 25
          max_idle_connections: 5
          connection_max_lifetime: 1h
  ```
- 既存のSQLite設定と併用できるようにする（コメントで説明）
- `config/staging/database.yaml`にも同様の設定を追加（オプション）

**受け入れ基準**:
- PostgreSQL用の設定が追加されている
- 設定値が適切に設定されている
- YAML形式が正しい
- 既存のSQLite設定と併用できる

---

### Phase 5: リトライ機能の実装

#### - [x] タスク 5.1: avast/retry-go/v4ライブラリの導入
**目的**: リトライ機能を実装するためのライブラリを導入する

**作業内容**:
- `go.mod`ファイルを開く
- `github.com/avast/retry-go/v4`を依存関係に追加
- `go mod tidy`を実行して依存関係を解決

**受け入れ基準**:
- `go.mod`に`github.com/avast/retry-go/v4`が追加されている
- `go mod tidy`がエラーなく実行できること
- 依存関係が正しく解決されている

---

#### - [x] タスク 5.2: リトライ設定定数の追加
**目的**: リトライ設定の定数を定義する

**作業内容**:
- `server/internal/db/connection.go`を開く
- リトライ設定の定数を追加:
  ```go
  const (
      MaxRetryAttempts = 3  // 最大3回（初回 + 2回のリトライ）
      RetryDelay       = 1 * time.Second
  )
  ```
- 適切なコメントを追加

**受け入れ基準**:
- リトライ設定の定数が定義されている
- 適切なコメントが追加されている
- 既存のコードスタイルに従っている

---

#### - [x] タスク 5.3: NewConnection関数へのリトライ機能追加
**目的**: 接続作成時にリトライ機能を追加する

**作業内容**:
- `server/internal/db/connection.go`を開く
- `github.com/avast/retry-go/v4`をインポート
- `NewConnection`関数内でリトライ機能を実装:
  ```go
  var db *sql.DB
  var err error
  
  err = retry.Do(
      func() error {
          db, err = sql.Open(driver, dsn)
          if err != nil {
              return err
          }
          
          // 接続プール設定
          // ... 設定処理 ...
          
          // 接続確認（リトライ対象）
          // 注意: 遅延接続のため、実際の接続確認は最初のクエリ実行時に行われる
          // ここでは接続オブジェクトの作成のみを確認
          return nil
      },
      retry.Attempts(MaxRetryAttempts),
      retry.Delay(RetryDelay),
      retry.OnRetry(func(n uint, err error) {
          log.Printf("Retrying database connection (attempt %d/%d): %v", n+1, MaxRetryAttempts, err)
      }),
  )
  ```
- エラーハンドリングを追加（すべてのリトライが失敗した場合）

**受け入れ基準**:
- リトライ機能が実装されている
- リトライ回数と間隔が正しく設定されている
- リトライ時にログが出力される
- エラーハンドリングが適切に実装されている
- 既存のコードスタイルに従っている

**注意**: 遅延接続のため、`sql.Open()`自体はエラーを返さない。リトライは実際のクエリ実行時の接続エラーに対して行う（タスク5.5参照）。

---

#### - [x] タスク 5.4: NewGORMConnection関数へのリトライ機能追加
**目的**: GORM接続作成時にリトライ機能を追加する

**作業内容**:
- `server/internal/db/connection.go`を開く
- `createGORMConnection`関数内でリトライ機能を実装
- `gorm.Open`呼び出しをリトライでラップ
- エラーハンドリングを追加

**受け入れ基準**:
- リトライ機能が実装されている
- リトライ回数と間隔が正しく設定されている
- リトライ時にログが出力される
- エラーハンドリングが適切に実装されている
- 既存のコードスタイルに従っている

---

#### - [x] タスク 5.5: クエリ実行時のリトライ機能追加（オプション）
**目的**: クエリ実行時の接続エラーに対してリトライ機能を追加する

**作業内容**:
- `server/internal/db/connection.go`を開く
- 接続エラー判定関数を追加:
  ```go
  func isConnectionError(err error) bool {
      if err == nil {
          return false
      }
      // sql.ErrConnDone などの接続エラーを判定
      return errors.Is(err, sql.ErrConnDone) || 
             strings.Contains(err.Error(), "connection") ||
             strings.Contains(err.Error(), "network")
  }
  ```
- 必要に応じて、Repository層でクエリ実行時のリトライ機能を実装

**受け入れ基準**:
- 接続エラー判定関数が実装されている
- クエリ実行時のリトライ機能が実装されている（必要に応じて）
- 接続エラー以外のエラーはリトライしない
- 既存のコードスタイルに従っている

**注意**: このタスクはオプションです。`database/sql`の標準機能により、接続プール設定により自動的に再接続されるため、明示的なリトライが不要な場合があります。

---

### Phase 6: 動作確認

#### - [x] タスク 6.1: 遅延接続の動作確認
**目的**: 遅延接続が正常に動作することを確認する

**作業内容**:
- PostgreSQLを起動（`./scripts/start-postgres.sh start`）
- PostgreSQLを停止（`./scripts/start-postgres.sh stop`）
- APIサーバーを起動（DB接続なしで起動可能であることを確認）
- 最初のクエリを実行（例: ユーザー一覧取得API）
- ログで接続確立のタイミングを確認
- PostgreSQLを起動してから再度クエリを実行
- 接続が確立されることを確認

**受け入れ基準**:
- DB接続なしでAPIサーバーが起動できること
- 最初のクエリ実行時に接続が確立されること
- ログで接続確立のタイミングが確認できること
- 接続が正常に確立されること

---

#### - [x] タスク 6.2: 自動再接続の動作確認
**目的**: 自動再接続が正常に動作することを確認する

**作業内容**:
- PostgreSQLを起動
- APIサーバーを起動
- クエリを実行して正常に動作することを確認
- PostgreSQLを停止
- クエリを実行（エラーが発生することを確認）
- PostgreSQLを再起動
- 次のクエリを実行（自動的に再接続されることを確認）
- ログで再接続のタイミングを確認

**受け入れ基準**:
- PostgreSQL停止時にクエリがエラーになること
- PostgreSQL再起動後に自動的に再接続されること
- ログで再接続のタイミングが確認できること
- 再接続後、クエリが正常に実行されること

---

#### - [x] タスク 6.3: リトライ機能の動作確認
**目的**: リトライ機能が正常に動作することを確認する

**作業内容**:
- PostgreSQLを停止
- APIサーバーを起動
- クエリを実行（リトライが実行されることを確認）
- ログでリトライの実行を確認
- 最大3回までリトライされることを確認
- すべてのリトライが失敗した場合にエラーが返されることを確認

**受け入れ基準**:
- リトライが実行されること
- リトライ回数が最大3回であること
- リトライ間隔が1秒であること
- ログでリトライの実行が確認できること
- すべてのリトライが失敗した場合にエラーが返されること

---

#### - [x] タスク 6.4: 接続プール設定の動作確認
**目的**: 接続プール設定が正常に動作することを確認する

**作業内容**:
- 設定ファイルで接続プール設定を確認
- デフォルト値が適用されることを確認（設定値が0以下の場合）
- 接続の再利用が正常に動作することを確認
- 複数のクエリを実行して接続が再利用されることを確認

**受け入れ基準**:
- 接続プール設定が正しく読み込まれること
- デフォルト値が適用されること
- 接続の再利用が正常に動作すること
- 接続プール設定により、パフォーマンスが向上すること

---

#### - [x] タスク 6.5: 既存機能の回帰テスト
**目的**: 既存のDB接続機能が正常に動作することを確認する

**作業内容**:
- SQLite環境で既存のテストを実行
- 既存のAPIエンドポイントが正常に動作することを確認
- 既存のRepository層が正常に動作することを確認
- 既存のService層が正常に動作することを確認

**受け入れ基準**:
- 既存のテストがすべて成功すること
- 既存のAPIエンドポイントが正常に動作すること
- 既存の機能が壊れていないこと
- 後方互換性が保たれていること

---

### Phase 7: ドキュメント更新

#### - [x] タスク 7.1: README.mdの更新
**目的**: PostgreSQL環境の構築方法をREADMEに追加する

**作業内容**:
- `README.md`を開く
- PostgreSQL環境の構築方法を追加:
  - Docker Composeでの起動方法
  - 起動スクリプトの使用方法
  - 設定ファイルの設定方法
- 既存のSQLite設定との併用方法を説明

**受け入れ基準**:
- PostgreSQL環境の構築方法が記載されている
- 起動スクリプトの使用方法が記載されている
- 設定ファイルの設定方法が記載されている
- 既存のSQLite設定との併用方法が説明されている

---

#### - [x] タスク 7.2: 設計ドキュメントの更新
**目的**: 実装内容を設計ドキュメントに反映する

**作業内容**:
- 実装完了後の設計ドキュメントを確認
- 実装内容と設計内容の差異を確認
- 必要に応じて設計ドキュメントを更新

**受け入れ基準**:
- 実装内容が設計ドキュメントに反映されている
- 設計ドキュメントが最新の状態であること

---

#### - [x] タスク 7.3: 開発環境のデータベース設定をSQLite版に戻す
**目的**: 動作確認完了後、開発環境のデータベース設定をSQLite版に戻す

**作業内容**:
- **注意**: このタスクは最後の最後、ユーザーの確認作業が終わった後に実施する
- `config/develop/database.yaml`を開く
- PostgreSQL用の設定をコメントアウトまたは削除
- SQLite用の設定が有効になっていることを確認
- 既存のSQLite設定が正常に動作することを確認
- 必要に応じて、PostgreSQL設定をコメントアウトして残す（将来の参照用）

**受け入れ基準**:
- PostgreSQL用の設定がコメントアウトまたは削除されている
- SQLite用の設定が有効になっていること
- 開発環境でSQLiteが正常に動作すること
- 既存のSQLite設定が壊れていないこと

**注意事項**:
- このタスクは、すべての動作確認が完了し、ユーザーの確認作業が終わった後に実施する
- PostgreSQL設定は削除せず、コメントアウトして残すことを推奨（将来の参照用）

---

## 実装順序の推奨

1. **Phase 1**: 起動時のDB接続確認処理の削除（最も影響が大きいため最初に実施）
2. **Phase 2**: 接続作成時のPing削除（遅延接続の実現）
3. **Phase 3**: 接続プール設定の確認と追加（既存機能の改善）
4. **Phase 4**: PostgreSQL環境の構築（動作確認環境の整備）
5. **Phase 5**: リトライ機能の実装（エラーハンドリングの強化）
6. **Phase 6**: 動作確認（全機能の検証）
7. **Phase 7**: ドキュメント更新（実装内容の記録）
8. **Phase 7.3**: 開発環境のデータベース設定をSQLite版に戻す（**最後の最後、ユーザーの確認作業が終わった後に実施**）

## 注意事項

### 遅延接続の注意点
- `sql.Open()`は接続を確立しないため、エラーが発生しない
- 実際のクエリ実行時に接続が確立される
- 接続エラーは最初のクエリ実行時に検出される

### 自動再接続の注意点
- 接続プール設定により、古い接続は破棄され、新しい接続が作成される
- `SetConnMaxLifetime`により、古い接続は定期的に破棄される
- データベースが復旧した際に、次のクエリ実行時に自動的に再接続される

### リトライ機能の注意点
- リトライは接続エラーに対してのみ実行される
- その他のエラー（SQL構文エラーなど）はリトライしない
- リトライ回数と間隔は適切に設定する

### PostgreSQL環境の注意点
- PostgreSQL環境は動作確認用（開発用途）
- 本番・staging環境では別の方法でPostgreSQLを起動
- SQLite環境との併用が可能
- **動作確認完了後、開発環境のデータベース設定はSQLite版に戻す（タスク7.3）**
