# ログ出力機能実装タスク一覧

## 概要
ログ出力機能の実装を段階的に進めるためのタスク一覧。要件定義書と設計書に基づき、実装可能な粒度でタスクを分解。

## 実装フェーズ

### Phase 1: 依存関係とディレクトリ構造の準備

#### - [x] タスク 1.1: ログライブラリの依存関係追加
**目的**: `logrus`と`lumberjack`ライブラリの依存関係を追加

**作業内容**:
- `server/go.mod`に以下の依存関係を追加:
  - `github.com/sirupsen/logrus`
  - `gopkg.in/natefinch/lumberjack.v2`
- `go mod tidy`を実行して依存関係を更新

**受け入れ基準**:
- `server/go.mod`に`logrus`と`lumberjack`の依存関係が追加されている
- `go mod tidy`が正常に実行される
- コンパイルエラーがない

---

#### - [x] タスク 1.2: ログディレクトリの作成
**目的**: ログファイルを出力するディレクトリを作成

**作業内容**:
- `logs/`ディレクトリを作成（プロジェクトルート）
- `logs/.gitkeep`ファイルを作成（ディレクトリを保持）

**受け入れ基準**:
- `logs/`ディレクトリが存在する
- `logs/.gitkeep`ファイルが存在する

---

#### - [x] タスク 1.3: .gitignoreの更新
**目的**: ログファイルをGitの管理対象外にする

**作業内容**:
- `.gitignore`に`logs/*`と`logs/**`を追加
- `logs`ディレクトリ自体は除外しない（`logs/`は追加しない）

**受け入れ基準**:
- `.gitignore`に`logs/*`と`logs/**`が追加されている
- `logs`ディレクトリ自体は除外されていない

---

### Phase 2: 設定の拡張

#### - [x] タスク 2.1: LoggingConfig構造体の拡張
**目的**: ログ出力先ディレクトリの設定項目を追加

**作業内容**:
- `server/internal/config/config.go`の`LoggingConfig`構造体に`OutputDir`フィールドを追加:
  ```go
  type LoggingConfig struct {
      Level     string `mapstructure:"level"`
      Format    string `mapstructure:"format"`
      Output    string `mapstructure:"output"`
      OutputDir string `mapstructure:"output_dir"` // 新規追加
  }
  ```

**受け入れ基準**:
- `LoggingConfig`構造体に`OutputDir`フィールドが追加されている
- コンパイルエラーがない

---

#### - [x] タスク 2.2: デフォルト値の設定
**目的**: `OutputDir`のデフォルト値を設定

**作業内容**:
- `server/internal/config/config.go`の`Load()`関数で、`cfg.Logging.OutputDir`が空の場合は`"logs"`を設定

**受け入れ基準**:
- `OutputDir`が空の場合、デフォルト値`"logs"`が設定される
- コンパイルエラーがない

---

#### - [x] タスク 2.3: 設定ファイルの更新（develop環境）
**目的**: develop環境の設定ファイルに`output_dir`項目を追加

**作業内容**:
- `config/develop/config.yaml`の`logging`セクションに`output_dir: logs`を追加

**受け入れ基準**:
- `config/develop/config.yaml`に`output_dir: logs`が追加されている
- YAML形式が正しい

---

#### - [x] タスク 2.4: 設定ファイルの更新（staging環境）
**目的**: staging環境の設定ファイルに`output_dir`項目を追加

**作業内容**:
- `config/staging/config.yaml`の`logging`セクションに`output_dir: logs`を追加

**受け入れ基準**:
- `config/staging/config.yaml`に`output_dir: logs`が追加されている
- YAML形式が正しい

---

#### - [x] タスク 2.5: 設定ファイルの更新（production環境）
**目的**: production環境の設定ファイルに`output_dir`項目を追加（例示ファイル）

**作業内容**:
- `config/production/config.yaml.example`の`logging`セクションに`output_dir: /var/log/go-webdb-template`を追加（コメント付きで説明）

**受け入れ基準**:
- `config/production/config.yaml.example`に`output_dir`が追加されている
- コメントで説明が記載されている

---

### Phase 3: ログ出力機能の実装

#### - [x] タスク 3.1: loggingパッケージディレクトリの作成
**目的**: ログ出力機能を実装するパッケージのディレクトリを作成

**作業内容**:
- `server/internal/logging/`ディレクトリを作成

**受け入れ基準**:
- `server/internal/logging/`ディレクトリが存在する

---

#### - [x] タスク 3.2: AccessLoggerの基本構造の作成
**目的**: AccessLoggerの基本構造を作成

**作業内容**:
- `server/internal/logging/access_logger.go`を作成
- パッケージ宣言とインポート文を追加:
  - `package logging`
  - 必要なパッケージのインポート（`fmt`, `os`, `path/filepath`）
  - `github.com/sirupsen/logrus`
  - `gopkg.in/natefinch/lumberjack.v2`
- `AccessLogger`構造体を定義:
  ```go
  type AccessLogger struct {
      logger *logrus.Logger
      writer *lumberjack.Logger
  }
  ```

**受け入れ基準**:
- `server/internal/logging/access_logger.go`が存在する
- コンパイルエラーがない

---

#### - [x] タスク 3.3: NewAccessLogger関数の実装
**目的**: AccessLoggerを作成する関数を実装

**作業内容**:
- `NewAccessLogger(logType string, outputDir string) (*AccessLogger, error)`関数を実装
- 出力ディレクトリの作成（`os.MkdirAll`）:
  - `outputDir`は絶対パスと相対パスの両方をサポート
  - 相対パスの場合は、サーバーの実行ディレクトリからの相対パスとして解釈
  - `filepath.Join(outputDir, ...)`でログファイルパスを結合
- `lumberjack.Logger`の設定:
  - `Filename`: `{logType}-access-2006-01-02.log`形式（`filepath.Join(outputDir, ...)`で結合）
  - `LocalTime: true`
  - `MaxSize: 0`, `MaxBackups: 0`, `MaxAge: 0`（日付別分割のみ使用）
- `logrus.Logger`の設定:
  - `SetOutput(writer)`
  - `SetFormatter(&logrus.TextFormatter{...})`
  - `SetLevel(logrus.InfoLevel)`

**受け入れ基準**:
- `NewAccessLogger`関数が実装されている
- 絶対パスと相対パスの両方が正しく処理される
- ディレクトリが存在しない場合は自動作成される
- エラーハンドリングが実装されている
- コンパイルエラーがない

---

#### - [x] タスク 3.4: カスタムテキストフォーマッターの実装
**目的**: 要件定義書の形式に合わせたカスタムフォーマッターを実装

**作業内容**:
- `CustomTextFormatter`構造体を定義（`logrus.TextFormatter`を埋め込み）
- `Format(entry *logrus.Entry) ([]byte, error)`メソッドを実装
- ログエントリを以下の形式で出力:
  ```
  [YYYY-MM-DD HH:MM:SS] METHOD /path HTTP/1.1 STATUS_CODE RESPONSE_TIME_MS REMOTE_IP USER_AGENT
  ```

**受け入れ基準**:
- `CustomTextFormatter`が実装されている
- ログエントリが要件定義書の形式で出力される
- コンパイルエラーがない

---

#### - [x] タスク 3.5: LogAccessメソッドの実装
**目的**: アクセスログを出力するメソッドを実装

**作業内容**:
- `LogAccess(method, path, protocol string, statusCode int, responseTimeMs float64, remoteIP, userAgent string)`メソッドを実装
- `logrus.WithFields()`を使用してログエントリを作成
- `Info("access")`でログを出力

**受け入れ基準**:
- `LogAccess`メソッドが実装されている
- すべてのパラメータがログに記録される
- コンパイルエラーがない

---

#### - [x] タスク 3.6: Closeメソッドの実装
**目的**: AccessLoggerをクローズするメソッドを実装

**作業内容**:
- `Close() error`メソッドを実装
- `writer.Close()`を呼び出す

**受け入れ基準**:
- `Close`メソッドが実装されている
- コンパイルエラーがない

---

#### - [x] タスク 3.7: responseWriter構造体の実装
**目的**: ステータスコードを記録するResponseWriterラッパーを実装

**作業内容**:
- `server/internal/logging/middleware.go`を作成
- `responseWriter`構造体を定義:
  ```go
  type responseWriter struct {
      http.ResponseWriter
      statusCode int
  }
  ```
- `WriteHeader(code int)`メソッドを実装

**受け入れ基準**:
- `responseWriter`構造体が実装されている
- ステータスコードが正しく記録される
- コンパイルエラーがない

---

#### - [x] タスク 3.8: AccessLogMiddleware構造体の実装
**目的**: アクセスログを記録するHTTPミドルウェアを実装

**作業内容**:
- `AccessLogMiddleware`構造体を定義:
  ```go
  type AccessLogMiddleware struct {
      accessLogger *AccessLogger
  }
  ```
- `NewAccessLogMiddleware(accessLogger *AccessLogger) *AccessLogMiddleware`関数を実装

**受け入れ基準**:
- `AccessLogMiddleware`構造体が実装されている
- `NewAccessLogMiddleware`関数が実装されている
- コンパイルエラーがない

---

#### - [x] タスク 3.9: Middlewareメソッドの実装
**目的**: HTTPミドルウェア関数を実装

**作業内容**:
- `Middleware(next http.Handler) http.Handler`メソッドを実装
- リクエスト開始時刻を記録
- `responseWriter`でラップしてステータスコードを取得
- 次のハンドラーを実行
- レスポンス時間を計算
- リモートIPアドレスを取得（`X-Forwarded-For`ヘッダーを考慮）
- User-Agentを取得
- `accessLogger.LogAccess()`を呼び出してログを出力

**受け入れ基準**:
- `Middleware`メソッドが実装されている
- リクエスト情報が正しく取得される
- レスポンス時間が正しく計算される
- アクセスログが出力される
- コンパイルエラーがない

---

### Phase 4: サーバーへの統合

#### - [x] タスク 4.1: APIサーバーへの統合
**目的**: APIサーバーにアクセスログ機能を統合

**作業内容**:
- `server/cmd/server/main.go`を修正
- `logging`パッケージをインポート
- `NewAccessLogger("api", cfg.Logging.OutputDir)`を呼び出してアクセスログを初期化
- エラーハンドリングを実装（エラー時は警告を出力して続行）
- `defer accessLogger.Close()`を追加
- `NewAccessLogMiddleware(accessLogger)`でミドルウェアを作成
- ルーターにミドルウェアを適用（`accessLogMiddleware.Middleware(r)`）

**受け入れ基準**:
- APIサーバー起動時にアクセスログが初期化される
- エラー時は警告が出力され、サーバーは正常に起動する
- すべてのHTTPリクエストがログに記録される
- コンパイルエラーがない

---

#### - [x] タスク 4.2: 管理画面サーバーへの統合（環境判定含む）
**目的**: 管理画面サーバーにアクセスログ機能を統合（production環境以外）

**作業内容**:
- `server/cmd/admin/main.go`を修正
- `logging`パッケージをインポート
- `os.Getenv("APP_ENV")`で環境を取得（デフォルトは`"develop"`）
- `env != "production"`の場合のみアクセスログを初期化
- `NewAccessLogger("admin", cfg.Logging.OutputDir)`を呼び出してアクセスログを初期化
- エラーハンドリングを実装（エラー時は警告を出力して続行）
- `defer accessLogger.Close()`を追加
- `NewAccessLogMiddleware(accessLogger)`でミドルウェアを作成
- ルーターにミドルウェアを適用

**受け入れ基準**:
- develop/staging環境では管理画面サーバーのアクセスログが初期化される
- production環境では管理画面サーバーのアクセスログが初期化されない
- すべてのHTTPリクエストがログに記録される（production環境以外）
- コンパイルエラーがない

---

#### - [ ] タスク 4.3: Routerへの統合（オプション）【スキップ】
**目的**: Routerでアクセスログミドルウェアを統合（オプション実装）

**作業内容**:
- `server/internal/api/router/router.go`を修正
- `NewRouter`関数のシグネチャに`accessLogger *logging.AccessLogger`パラメータを追加（オプション）
- `accessLogger`が`nil`でない場合、ミドルウェアを適用
- `server/cmd/server/main.go`と`server/cmd/admin/main.go`で`NewRouter`の呼び出しを更新

**注意**: このタスクはオプションです。タスク4.1と4.2で直接ミドルウェアを適用する方法でも実装可能です。

**受け入れ基準**:
- Routerでアクセスログミドルウェアが統合されている（オプション実装の場合）
- 既存の機能が正常に動作する
- コンパイルエラーがない

---

### Phase 5: テスト

#### - [x] タスク 5.1: AccessLoggerのユニットテスト
**目的**: AccessLoggerの機能をテスト

**作業内容**:
- `server/internal/logging/access_logger_test.go`を作成
- `TestNewAccessLogger`テストを実装:
  - 正常系: AccessLoggerが正常に作成される
  - 異常系: ディレクトリ作成に失敗した場合、エラーを返す
- `TestAccessLogger_LogAccess`テストを実装:
  - 正常系: アクセスログが正常に出力される
  - 正常系: ログエントリのフォーマットが正しい

**受け入れ基準**:
- すべてのテストがパスする
- テストカバレッジが適切である

---

#### - [x] タスク 5.2: AccessLogMiddlewareのユニットテスト
**目的**: AccessLogMiddlewareの機能をテスト

**作業内容**:
- `server/internal/logging/middleware_test.go`を作成
- `TestAccessLogMiddleware_Middleware`テストを実装:
  - 正常系: ミドルウェアが正常に動作する
  - 正常系: レスポンス時間が正しく記録される
  - 正常系: ステータスコードが正しく記録される
  - 正常系: リモートIPアドレスが正しく取得される

**受け入れ基準**:
- すべてのテストがパスする
- テストカバレッジが適切である

---

#### - [ ] タスク 5.3: 統合テストの実装【スキップ：オプション】
**目的**: アクセスログ機能の統合テストを実装

**作業内容**:
- `server/test/integration/access_log_test.go`を作成
- `TestAPIServerAccessLog`テストを実装:
  - APIサーバーへのリクエストがログファイルに記録される
  - ログファイル名に日付が含まれる
- `TestAdminServerAccessLog`テストを実装:
  - 管理画面サーバーへのリクエストがログファイルに記録される（production環境以外）
  - production環境ではログが出力されない

**受け入れ基準**:
- すべてのテストがパスする
- ログファイルが正しく作成される
- ログファイルの内容が正しい形式である

---

#### - [ ] タスク 5.4: 日付別ファイル分割のテスト【スキップ：オプション】
**目的**: 日付別ファイル分割機能をテスト

**作業内容**:
- 日付が変わったら新しいファイルに切り替わることを確認するテストを実装
- ログファイル名に日付（YYYY-MM-DD形式）が含まれることを確認

**受け入れ基準**:
- 日付別ファイル分割が正常に動作する
- ログファイル名が正しい形式である

---

#### - [x] タスク 5.5: 既存テストの確認
**目的**: 既存のテストが正常に動作することを確認

**作業内容**:
- 既存のテストスイートを実行
- テストが失敗しないことを確認

**受け入れ基準**:
- 既存のテストがすべてパスする
- 新規実装による影響がない

---

### Phase 6: 動作確認とドキュメント更新

#### - [x] タスク 6.1: 動作確認（APIサーバー）
**目的**: APIサーバーのアクセスログ機能を動作確認

**作業内容**:
- APIサーバーを起動
- いくつかのHTTPリクエストを送信
- `logs/api-access-YYYY-MM-DD.log`ファイルが作成されることを確認
- ログファイルの内容が正しい形式であることを確認
- ログエントリに必要な情報が含まれていることを確認

**受け入れ基準**:
- ログファイルが正常に作成される
- ログファイルの内容が正しい形式である
- すべてのHTTPリクエストがログに記録される

---

#### - [x] タスク 6.2: 動作確認（管理画面サーバー）
**目的**: 管理画面サーバーのアクセスログ機能を動作確認

**作業内容**:
- develop環境で管理画面サーバーを起動
- いくつかのHTTPリクエストを送信
- `logs/admin-access-YYYY-MM-DD.log`ファイルが作成されることを確認
- production環境で管理画面サーバーを起動
- ログファイルが作成されないことを確認

**受け入れ基準**:
- develop環境ではログファイルが正常に作成される
- production環境ではログファイルが作成されない
- ログファイルの内容が正しい形式である

---

#### - [ ] タスク 6.3: 日付別ファイル分割の動作確認【スキップ：動作確認困難】
**目的**: 日付別ファイル分割機能を動作確認

**作業内容**:
- サーバーを起動してログを出力
- 日付が変わったら（または手動で日付を変更して）新しいファイルに切り替わることを確認
- ログファイル名に日付（YYYY-MM-DD形式）が含まれることを確認

**受け入れ基準**:
- 日付が変わったら新しいファイルに切り替わる
- ログファイル名が正しい形式である

---

#### - [ ] タスク 6.4: 設定変更の動作確認【スキップ：動作確認困難】
**目的**: ログ出力先の設定変更を動作確認

**作業内容**:
- 設定ファイルで`output_dir`を変更（相対パスと絶対パスの両方をテスト）
- サーバーを起動
- 指定したディレクトリにログファイルが作成されることを確認
- 相対パスと絶対パスの両方が正しく動作することを確認

**受け入れ基準**:
- 設定ファイルでログ出力先を変更できる
- 相対パス（例: `logs`）が正しく動作する
- 絶対パス（例: `/var/log/go-webdb-template`）が正しく動作する
- デフォルト値（`logs`）が正しく動作する

---

#### - [ ] タスク 6.5: エラーハンドリングの動作確認【スキップ：動作確認困難】
**目的**: エラーハンドリングが正常に動作することを確認

**作業内容**:
- ログディレクトリの作成に失敗するケースをテスト（権限不足等）
- ログファイルへの書き込みに失敗するケースをテスト
- エラーが発生してもHTTPリクエスト処理が継続することを確認
- エラーメッセージが標準エラー出力に出力されることを確認

**受け入れ基準**:
- エラーが発生してもHTTPリクエスト処理は継続する
- エラーメッセージが適切に出力される

---

#### - [ ] タスク 6.6: ドキュメント更新（オプション）【スキップ】
**目的**: README.mdまたは関連ドキュメントにログ出力機能の説明を追加

**作業内容**:
- README.mdまたは関連ドキュメントにログ出力機能の説明を追加
- ログファイルの場所、形式、設定方法を記載

**注意**: このタスクはオプションです。必要に応じて実装してください。

**受け入れ基準**:
- ドキュメントにログ出力機能の説明が追加されている（オプション）

---

## 実装順序の推奨

1. **Phase 1**: 依存関係とディレクトリ構造の準備（タスク1.1〜1.3）
2. **Phase 2**: 設定の拡張（タスク2.1〜2.5）
3. **Phase 3**: ログ出力機能の実装（タスク3.1〜3.9）
4. **Phase 4**: サーバーへの統合（タスク4.1〜4.2、4.3はオプション）
5. **Phase 5**: テスト（タスク5.1〜5.5）
6. **Phase 6**: 動作確認とドキュメント更新（タスク6.1〜6.6）

## 注意事項

- 各タスクは独立して実装可能な粒度に分解されています
- タスク4.3と6.6はオプションです
- 実装中は既存の機能に影響を与えないよう注意してください
- エラーハンドリングは適切に実装し、HTTPリクエスト処理に影響を与えないようにしてください

