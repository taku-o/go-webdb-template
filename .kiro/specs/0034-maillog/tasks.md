# メール送信ログ機能実装タスク一覧

## 概要
メール送信ログ機能の実装を段階的に進めるためのタスク一覧。要件定義書と設計書に基づき、実装可能な粒度でタスクを分解。

## 実装フェーズ

### Phase 1: 設定とインフラの準備

#### - [ ] タスク 1.1: LoggingConfig構造体の拡張
**目的**: メール送信ログの設定を管理するフィールドを追加する

**作業内容**:
- `server/internal/config/config.go`の`LoggingConfig`構造体に`MailLogOutputDir string`フィールドを追加
- `mapstructure:"mail_log_output_dir"`タグを追加
- 設定ファイル読み込み処理で`MailLogOutputDir`が空の場合は`OutputDir`を使用するロジックを追加
- `config.Load()`関数内でデフォルト値の設定処理を追加

**受け入れ基準**:
- `LoggingConfig`構造体に`MailLogOutputDir`フィールドが追加されている
- 設定ファイル読み込み時にデフォルト値が正しく設定される
- 既存のコードスタイルに従っている

---

#### - [ ] タスク 1.2: 設定ファイルにメール送信ログ設定を追加
**目的**: 各環境の設定ファイルにメール送信ログ設定を追加する

**作業内容**:
- `config/develop/config.yaml`の`logging`セクションに`mail_log_output_dir: "logs"`を追加（オプション）
- `config/staging/config.yaml`の`logging`セクションに`mail_log_output_dir: "logs"`を追加（オプション）
- `config/production/config.yaml.example`の`logging`セクションに`mail_log_output_dir`のコメントを追加（本番環境では使用しない旨を記載）

**受け入れ基準**:
- 3つの設定ファイルに適切な設定が追加されている
- 設定値が適切に設定されている
- YAML形式が正しい

---

### Phase 2: MailLoggerの実装

#### - [ ] タスク 2.1: MailLogEntry構造体の定義
**目的**: メール送信ログのJSON構造体を定義する

**作業内容**:
- `server/internal/logging/mail_logger.go`を新規作成
- `MailLogEntry`構造体を定義:
  - `Timestamp string` - 送信時刻（YYYY-MM-DD HH:MM:SS形式）
  - `To []string` - 送信先メールアドレスリスト
  - `Subject string` - メール件名
  - `Body string` - メール本文（200文字制限）
  - `BodyTruncated bool` - 本文が切り捨てられたかどうか
  - `SenderType string` - 送信実装の種類（mock/mailpit/ses）
  - `Success bool` - 送信成功/失敗
  - `Error string` - エラーメッセージ（エラー時のみ、omitemptyタグ付き）
- JSONタグを適切に設定

**受け入れ基準**:
- `MailLogEntry`構造体が定義されている
- すべてのフィールドにJSONタグが設定されている
- 既存のコードスタイルに従っている

---

#### - [ ] タスク 2.2: MailLogger構造体の定義
**目的**: メール送信ログを出力するロガー構造体を定義する

**作業内容**:
- `MailLogger`構造体を定義:
  - `logger *logrus.Logger` - logrusロガー
  - `writer io.WriteCloser` - ログファイルへの書き込み用
  - `enabled bool` - ログ出力が有効かどうか
- パッケージのインポートを追加:
  - `encoding/json`
  - `fmt`
  - `io`
  - `os`
  - `path/filepath`
  - `time`
  - `github.com/sirupsen/logrus`
  - `gopkg.in/natefinch/lumberjack.v2`

**受け入れ基準**:
- `MailLogger`構造体が定義されている
- 必要なパッケージがインポートされている
- 既存のコードスタイルに従っている

---

#### - [ ] タスク 2.3: NewMailLogger関数の実装
**目的**: MailLoggerを作成する関数を実装する

**作業内容**:
- `NewMailLogger(outputDir string, enabled bool) (*MailLogger, error)`関数を実装
- `enabled`が`false`の場合は`nil`を返す
- 出力ディレクトリの作成処理:
  - `os.MkdirAll(outputDir, 0755)`を使用
  - エラーが発生した場合はエラーを返す
- lumberjackの設定:
  - ファイル名: `mail-YYYY-MM-DD.log`形式
  - `MaxSize: 0`（サイズ制限なし）
  - `MaxBackups: 0`（バックアップ保持数なし）
  - `MaxAge: 0`（保持日数なし）
  - `LocalTime: true`（ローカルタイムゾーンを使用）
  - `Compress: false`（圧縮なし）
- logrusの設定:
  - `logrus.New()`でロガーを作成
  - `SetOutput(writer)`で出力先を設定
  - `SetFormatter(&MailLogFormatter{})`でフォーマッターを設定
  - `SetLevel(logrus.InfoLevel)`でログレベルを設定
- `MailLogger`インスタンスを返す

**受け入れ基準**:
- `NewMailLogger`関数が実装されている
- 環境判定に基づいて適切にロガーを作成または`nil`を返す
- ディレクトリ作成エラーが適切に処理されている
- lumberjackとlogrusが正しく設定されている

---

#### - [ ] タスク 2.4: MailLogFormatterの実装
**目的**: logrusのフォーマッターを実装する（実際には使用しないが、インターフェースに適合させるため）

**作業内容**:
- `MailLogFormatter`構造体を定義
- `Format(entry *logrus.Entry) ([]byte, error)`メソッドを実装
- 実際のJSON出力は`LogMail`メソッド内で直接行うため、このフォーマッターは空の実装でOK

**受け入れ基準**:
- `MailLogFormatter`構造体が定義されている
- `Format`メソッドが実装されている
- logrus.Formatterインターフェースを実装している

---

#### - [ ] タスク 2.5: LogMailメソッドの実装
**目的**: メール送信ログを出力するメソッドを実装する

**作業内容**:
- `LogMail(to []string, subject, body, senderType string, success bool, err error)`メソッドを実装
- `MailLogger`が`nil`または`enabled`が`false`の場合は何もしない
- メール本文の切り捨て処理:
  - `len(body) > 200`の場合、`body[:200] + "..."`に切り捨て
  - `bodyTruncated`フラグを`true`に設定
- エラーメッセージの取得:
  - `err != nil`の場合、`err.Error()`でエラーメッセージを取得
  - それ以外は空文字列
- `MailLogEntry`構造体の作成:
  - `Timestamp`: `time.Now().Format("2006-01-02 15:04:05")`
  - その他のフィールドを設定
- JSON形式で出力:
  - `json.Marshal(entry)`でJSONエンコード
  - エンコードエラーが発生した場合は標準エラー出力に記録して`return`
  - `writer.Write(append(jsonBytes, '\n'))`でログファイルに書き込み
  - 書き込みエラーが発生した場合は標準エラー出力に記録（メール送信処理は継続）

**受け入れ基準**:
- `LogMail`メソッドが実装されている
- メール本文の200文字制限が正しく実装されている
- JSON形式で正しく出力される
- エラーハンドリングが適切に実装されている
- ログ出力に失敗してもメール送信処理は継続する

---

#### - [ ] タスク 2.6: Closeメソッドの実装
**目的**: MailLoggerをクローズするメソッドを実装する

**作業内容**:
- `Close() error`メソッドを実装
- `MailLogger`が`nil`または`writer`が`nil`の場合は`nil`を返す
- `writer.Close()`を呼び出してエラーを返す

**受け入れ基準**:
- `Close`メソッドが実装されている
- 適切にリソースが解放される
- エラーハンドリングが適切に実装されている

---

### Phase 3: EmailServiceへの統合

#### - [ ] タスク 3.1: EmailService構造体の拡張
**目的**: EmailServiceにMailLoggerを統合する

**作業内容**:
- `server/internal/service/email/email_service.go`を開く
- `EmailService`構造体に`logger *logging.MailLogger`フィールドを追加
- `logging`パッケージのインポートを追加

**受け入れ基準**:
- `EmailService`構造体に`logger`フィールドが追加されている
- 必要なパッケージがインポートされている
- 既存のコードスタイルに従っている

---

#### - [ ] タスク 3.2: NewEmailService関数の拡張
**目的**: NewEmailService関数にMailLoggerパラメータを追加する

**作業内容**:
- `NewEmailService(cfg *config.EmailConfig, mailLogger *logging.MailLogger) (*EmailService, error)`にシグネチャを変更
- `EmailService`構造体の初期化時に`logger: mailLogger`を設定
- 既存の送信実装選択ロジックは変更しない

**受け入れ基準**:
- `NewEmailService`関数のシグネチャが変更されている
- `MailLogger`が正しく設定されている
- 既存の機能に影響がない

---

#### - [ ] タスク 3.3: getSenderTypeメソッドの実装
**目的**: 送信実装の種類を取得するメソッドを実装する

**作業内容**:
- `getSenderType() string`メソッドを実装
- `s.sender`の型を判定:
  - `*MockSender`の場合: `"mock"`を返す
  - `*MailpitSender`の場合: `"mailpit"`を返す
  - `*SESSender`の場合: `"ses"`を返す
  - それ以外の場合: `"unknown"`を返す

**受け入れ基準**:
- `getSenderType`メソッドが実装されている
- 各送信実装の種類が正しく判定される
- 既存のコードスタイルに従っている

---

#### - [ ] タスク 3.4: SendEmailメソッドへのログ出力統合
**目的**: SendEmailメソッドにログ出力処理を追加する

**作業内容**:
- `SendEmail`メソッドを開く
- メール送信前のログ出力（オプション）:
  - `s.logger != nil`の場合、`s.logger.LogMail(to, subject, body, s.getSenderType(), false, nil)`を呼び出し
  - 送信結果はまだ不明のため、`success=false`で記録
- メール送信処理（既存のコード）:
  - `err := s.sender.Send(ctx, to, subject, body)`
- メール送信後のログ出力:
  - `s.logger != nil`の場合、`s.logger.LogMail(to, subject, body, s.getSenderType(), err == nil, err)`を呼び出し
  - `success`は`err == nil`で判定
- エラーを返す（既存の動作）

**受け入れ基準**:
- `SendEmail`メソッドにログ出力処理が追加されている
- 送信前後のログ出力が正しく実行される
- 送信成功・失敗の両方でログが出力される
- 既存のメール送信機能に影響がない

---

### Phase 4: アプリケーション起動時の統合

#### - [ ] タスク 4.1: 環境判定処理の実装
**目的**: アプリケーション起動時に環境を判定してMailLoggerを初期化する

**作業内容**:
- `server/cmd/server/main.go`を開く
- 設定読み込み後に環境判定処理を追加:
  - `appEnv := os.Getenv("APP_ENV")`
  - `appEnv == ""`の場合は`appEnv = "develop"`（デフォルト）
  - `mailLogEnabled := appEnv == "develop" || appEnv == "staging"`
- `logging`パッケージのインポートを追加

**受け入れ基準**:
- 環境判定処理が実装されている
- develop/staging環境では有効、production環境では無効になる
- デフォルトはdevelop環境として扱われる

---

#### - [ ] タスク 4.2: MailLoggerの初期化
**目的**: アプリケーション起動時にMailLoggerを初期化する

**作業内容**:
- `server/cmd/server/main.go`でMailLoggerの初期化処理を追加:
  - `mailLogger, err := logging.NewMailLogger(cfg.Logging.MailLogOutputDir, mailLogEnabled)`
  - エラーが発生した場合は警告を出力して`mailLogger = nil`に設定
  - `defer mailLogger.Close()`でアプリケーション終了時にクローズ
- エラーハンドリング:
  - `log.Printf("Warning: Failed to initialize mail logger: %v", err)`
  - `log.Println("Mail logging will be disabled")`

**受け入れ基準**:
- MailLoggerの初期化処理が実装されている
- エラーハンドリングが適切に実装されている
- アプリケーション終了時にリソースが解放される

---

#### - [ ] タスク 4.3: EmailService作成時のMailLogger統合
**目的**: EmailService作成時にMailLoggerを渡す

**作業内容**:
- `server/cmd/server/main.go`でEmailService作成処理を確認
- `NewEmailService(cfg.Email, mailLogger)`に変更
- 既存の`NewEmailService`呼び出しを更新

**受け入れ基準**:
- EmailService作成時にMailLoggerが渡されている
- 既存の機能に影響がない

---

### Phase 5: テスト実装

#### - [ ] タスク 5.1: MailLoggerのユニットテスト
**目的**: MailLoggerの各メソッドをテストする

**作業内容**:
- `server/internal/logging/mail_logger_test.go`を新規作成
- `NewMailLogger`のテスト:
  - 有効化時の正常系テスト
  - 無効化時のnil返却テスト
  - ディレクトリ作成失敗時のエラーテスト
- `LogMail`のテスト:
  - 正常なログ出力のテスト
  - メール本文の200文字制限テスト
  - JSON形式での出力テスト
  - ログ出力失敗時のエラーハンドリングテスト
- `Close`のテスト:
  - 正常なクローズ処理のテスト
  - nilの場合のテスト

**受け入れ基準**:
- すべてのテストが実装されている
- テストが成功する
- テストカバレッジが適切である

---

#### - [ ] タスク 5.2: EmailServiceのユニットテスト拡張
**目的**: EmailServiceのログ出力機能をテストする

**作業内容**:
- 既存の`server/internal/service/email/email_service_test.go`を確認
- `SendEmail`メソッドのテストにログ出力の検証を追加:
  - ログ出力が正しく実行されることのテスト
  - 送信成功時のログ内容のテスト
  - 送信失敗時のログ内容のテスト
  - ログ出力失敗時もメール送信が継続することのテスト
- MockLoggerを作成してテストに使用

**受け入れ基準**:
- ログ出力に関するテストが追加されている
- すべてのテストが成功する
- 既存のテストに影響がない

---

#### - [ ] タスク 5.3: 統合テストの実装
**目的**: メール送信とログ出力の統合テストを実装する

**作業内容**:
- `server/internal/service/email/integration_test.go`を新規作成（必要に応じて）
- 実際のメール送信処理とログ出力の統合テスト:
  - MockSenderでのログ出力テスト
  - MailpitSenderでのログ出力テスト（Mailpitが起動している場合）
  - 環境別制御のテスト（develop/staging/production）
- ログファイルの内容を確認するテスト

**受け入れ基準**:
- 統合テストが実装されている
- すべてのテストが成功する
- ログファイルに正しく記録されることが確認できる

---

#### - [ ] タスク 5.4: E2Eテストの実装
**目的**: メール送信APIのE2Eテストを実装する

**作業内容**:
- 既存のE2Eテストを確認
- メール送信APIエンドポイントのE2Eテストにログ出力の検証を追加:
  - メール送信APIを呼び出す
  - ログファイルに正しく記録されることを確認
  - JSON形式で正しく出力されることを確認
  - メール本文の200文字制限が正しく動作することを確認

**受け入れ基準**:
- E2Eテストが実装されている
- すべてのテストが成功する
- ログファイルの内容が正しいことが確認できる

---

### Phase 6: 動作確認とドキュメント

#### - [ ] タスク 6.1: 動作確認
**目的**: 実装した機能が正しく動作することを確認する

**作業内容**:
- 開発環境でメール送信を実行
- ログファイル（`logs/mail-YYYY-MM-DD.log`）が作成されることを確認
- ログファイルの内容を確認:
  - JSON形式で出力されていること
  - 送信時刻、送信先、件名、本文が記録されていること
  - メール本文が200文字以上の場合、200文字で切り捨てられていること
  - 送信成功・失敗の状態が記録されていること
- 環境別制御の確認:
  - develop環境: ログが出力されること
  - staging環境: ログが出力されること
  - production環境: ログが出力されないこと

**受け入れ基準**:
- すべての動作確認項目が確認できる
- ログファイルが正しく作成される
- ログ内容が正しい

---

#### - [ ] タスク 6.2: 既存テストの確認
**目的**: 既存のテストがすべて成功することを確認する

**作業内容**:
- `go test ./...`を実行
- すべてのテストが成功することを確認
- テストエラーが発生した場合は修正

**受け入れ基準**:
- すべての既存テストが成功する
- 新しいテストも成功する

---

#### - [ ] タスク 6.3: リントとフォーマットの確認
**目的**: コードがリントとフォーマットの規約に従っていることを確認する

**作業内容**:
- `golangci-lint run`を実行（またはプロジェクトのリントツール）
- `go fmt ./...`を実行
- リントエラーがないことを確認

**受け入れ基準**:
- リントエラーがない
- コードが適切にフォーマットされている

---

## 受け入れ基準（全体）

### 機能要件
1. **メール送信ログ出力**: メール送信時にログが出力されること
2. **日付別ファイル分割**: 日付が変わったら自動的に新しいログファイルに切り替わること
3. **JSON形式**: ログがJSON形式で出力されること
4. **メール内容の切り捨て**: メール本文が200文字以上の場合、200文字で切り捨てられること
5. **環境別制御**: 開発環境とステージング環境でのみログが出力され、本番環境では出力されないこと
6. **設定による出力先変更**: 設定ファイルでログ出力先を変更できること

### 非機能要件
1. **パフォーマンス**: ログ出力がメール送信処理のパフォーマンスに影響を与えないこと
2. **エラーハンドリング**: ログ出力に失敗してもメール送信処理が継続すること
3. **既存機能との統合**: 既存のログ機能と統合され、同じ`logs`ディレクトリを使用すること

## 注意事項

- ログ出力に失敗してもメール送信処理は継続すること
- 本番環境ではログを出力しないこと
- メール本文は200文字で切り捨てること
- 既存のメール送信機能に影響を与えないこと
