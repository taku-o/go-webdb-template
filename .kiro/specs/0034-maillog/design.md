# メール送信ログ機能設計書

## 1. 概要

### 1.1 設計の目的
要件定義書に基づき、メール送信時にデバッグ用のログを出力する機能の詳細設計を定義する。既存のアクセスログ機能やSQLログ機能と同じパターンで実装し、`logrus`と`lumberjack`ライブラリを活用して、既存のアーキテクチャに統合する。

### 1.2 設計の範囲
- メール送信ログ出力機能のアーキテクチャ設計
- メール送信実装（MockSender、MailpitSender、SESSender）への統合設計
- ログライブラリ（logrus + lumberjack）の統合設計
- 設定構造体の拡張設計
- 日付別ファイル分割の実装設計
- 環境別制御（develop/staging/production）の設計
- メール本文の200文字制限処理の設計
- JSON形式でのログ出力設計
- エラーハンドリング設計
- テスト戦略

### 1.3 設計方針
- **既存ライブラリの活用**: `logrus`と`lumberjack`ライブラリを既存のアクセスログ機能と同じパターンで使用
- **既存機能との統合**: 既存のアクセスログ機能やSQLログ機能と同じ`logs`ディレクトリを使用
- **環境別制御**: 開発環境とステージング環境でのみメール送信ログを出力し、本番環境では出力しない
- **セキュリティ**: メール本文を200文字で切り捨てることで機密情報の漏洩リスクを軽減
- **エラーハンドリング**: ログ出力に失敗してもメール送信処理は継続する

## 2. アーキテクチャ設計

### 2.1 ディレクトリ構造

#### 2.1.1 変更前の構造
```
server/
├── internal/
│   ├── service/
│   │   └── email/
│   │       ├── email_service.go    # EmailService
│   │       ├── email_sender.go    # EmailSenderインターフェース
│   │       ├── mock_sender.go     # MockSender実装
│   │       ├── mailpit_sender.go  # MailpitSender実装
│   │       └── ses_sender.go     # SESSender実装
│   ├── config/
│   │   └── config.go              # LoggingConfig（OutputDir, SQLLogOutputDirのみ）
│   └── logging/
│       └── access_logger.go       # アクセスログ機能
└── ...
```

#### 2.1.2 変更後の構造
```
server/
├── internal/
│   ├── service/
│   │   └── email/
│   │       ├── email_service.go    # EmailService（ロガー統合）
│   │       ├── email_sender.go     # EmailSenderインターフェース（変更不要）
│   │       ├── mock_sender.go      # MockSender実装（ログ出力追加）
│   │       ├── mailpit_sender.go   # MailpitSender実装（ログ出力追加）
│   │       └── ses_sender.go       # SESSender実装（ログ出力追加）
│   ├── config/
│   │   └── config.go               # LoggingConfig拡張（MailLogOutputDir追加）
│   └── logging/
│       ├── access_logger.go        # アクセスログ機能（変更不要）
│       └── mail_logger.go          # 新規: メール送信ログ機能
└── ...
```

### 2.2 メール送信ログ出力の実行フロー

```
┌─────────────────────────────────────────────────────────────┐
│              1. アプリケーション起動                           │
│              server/cmd/server/main.go                        │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              2. 設定ファイル読み込み                           │
│              config.Load()                                  │
│              - APP_ENV環境変数を確認                         │
│              - LoggingConfigを読み込み                       │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              3. メール送信ログLogger初期化（環境判定）          │
│              - develop/staging: MailLogger有効化              │
│              - production: MailLogger無効化                 │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              4. EmailService初期化                            │
│              NewEmailService(cfg)                            │
│              - 送信実装（MockSender/MailpitSender/SESSender）│
│              - MailLoggerを注入                              │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              5. メール送信リクエスト                           │
│              EmailService.SendEmail()                       │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              6. 送信実装のSend()呼び出し                       │
│              MockSender.Send() / MailpitSender.Send() /      │
│              SESSender.Send()                                │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              7. メール送信ログ出力（送信前）                    │
│              - 送信時刻、送信先、件名、本文を記録              │
│              - 本文が200文字以上の場合、200文字で切り捨て      │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              8. メール送信処理実行                             │
│              - MockSender: 標準出力に出力                    │
│              - MailpitSender: Mailpitに送信                  │
│              - SESSender: AWS SESに送信                      │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              9. メール送信ログ出力（送信後）                    │
│              - 送信結果（成功/失敗）を記録                     │
│              - エラーが発生した場合はエラーメッセージを記録    │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              10. レスポンス返却                                │
└─────────────────────────────────────────────────────────────┘
```

### 2.3 既存アーキテクチャとの統合

```
┌─────────────────────────────────────────────────────────────┐
│              EmailService (internal/service/email)          │
│              - MailLoggerを保持                              │
│              - SendEmail()でログ出力                          │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              EmailSender実装                                │
│              - MockSender                                    │
│              - MailpitSender                                 │
│              - SESSender                                     │
│              - 各Send()メソッド内でログ出力                  │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              MailLogger (internal/logging)                   │
│              - logrus + lumberjack                          │
│              - 日付別ファイル出力                            │
│              - JSON形式で出力                                │
└─────────────────────────────────────────────────────────────┘
```

## 3. コンポーネント設計

### 3.1 MailLoggerの設計

#### 3.1.1 パッケージ構造
```go
package logging

import (
    "encoding/json"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "strings"
    "time"

    "github.com/sirupsen/logrus"
    "gopkg.in/natefinch/lumberjack.v2"
)
```

#### 3.1.2 MailLogger構造体の設計

```go
// MailLogger はメール送信ログを出力するロガー
type MailLogger struct {
    logger *logrus.Logger
    writer io.WriteCloser
    enabled bool
}

// MailLogEntry はメール送信ログのJSON構造体
type MailLogEntry struct {
    Timestamp     string   `json:"timestamp"`
    To            []string `json:"to"`
    Subject       string   `json:"subject"`
    Body          string   `json:"body"`
    BodyTruncated bool     `json:"body_truncated"`
    SenderType    string   `json:"sender_type"`
    Success       bool     `json:"success"`
    Error         string   `json:"error,omitempty"`
}

// NewMailLogger は新しいMailLoggerを作成
// outputDirは絶対パスと相対パスの両方をサポートします
// - 相対パス: サーバーの実行ディレクトリからの相対パス（例: "logs"）
// - 絶対パス: システムの絶対パス（例: "/var/log/go-webdb-template"）
// enabledがfalseの場合はnilを返します
func NewMailLogger(outputDir string, enabled bool) (*MailLogger, error) {
    // ログが無効な場合はnilを返す
    if !enabled {
        return nil, nil
    }

    // 出力ディレクトリの作成
    if err := os.MkdirAll(outputDir, 0755); err != nil {
        return nil, fmt.Errorf("failed to create log directory: %w", err)
    }

    // lumberjackの設定（日付別ファイル分割）
    logFileName := fmt.Sprintf("mail-%s.log", time.Now().Format("2006-01-02"))
    writer := &lumberjack.Logger{
        Filename:   filepath.Join(outputDir, logFileName),
        MaxSize:    0,     // サイズ制限なし（日付別分割のみ使用）
        MaxBackups: 0,     // バックアップ保持数なし（将来の拡張用）
        MaxAge:     0,     // 保持日数なし（将来の拡張用）
        LocalTime:  true,  // ローカルタイムゾーンを使用
        Compress:   false, // 圧縮なし
    }

    // logrusの設定
    logger := logrus.New()
    logger.SetOutput(writer)
    logger.SetFormatter(&MailLogFormatter{})
    logger.SetLevel(logrus.InfoLevel)

    return &MailLogger{
        logger: logger,
        writer: writer,
        enabled: true,
    }, nil
}

// LogMail はメール送信ログを出力
func (m *MailLogger) LogMail(to []string, subject, body, senderType string, success bool, err error) {
    if m == nil || !m.enabled {
        return
    }

    // メール本文の切り捨て処理
    bodyTruncated := false
    if len(body) > 200 {
        body = body[:200] + "..."
        bodyTruncated = true
    }

    // エラーメッセージの取得
    errorMsg := ""
    if err != nil {
        errorMsg = err.Error()
    }

    // ログエントリの作成
    entry := MailLogEntry{
        Timestamp:     time.Now().Format("2006-01-02 15:04:05"),
        To:            to,
        Subject:       subject,
        Body:          body,
        BodyTruncated: bodyTruncated,
        SenderType:    senderType,
        Success:       success,
        Error:         errorMsg,
    }

    // JSON形式で出力
    jsonBytes, err := json.Marshal(entry)
    if err != nil {
        // JSONエンコードに失敗した場合は標準エラー出力に記録
        fmt.Fprintf(os.Stderr, "Failed to encode mail log: %v\n", err)
        return
    }

    // ログファイルに出力（改行を追加）
    if _, err := m.writer.Write(append(jsonBytes, '\n')); err != nil {
        // ログ出力に失敗した場合は標準エラー出力に記録
        fmt.Fprintf(os.Stderr, "Failed to write mail log: %v\n", err)
    }
}

// Close はロガーをクローズ
func (m *MailLogger) Close() error {
    if m == nil || m.writer == nil {
        return nil
    }
    return m.writer.Close()
}
```

#### 3.1.3 MailLogFormatterの設計

```go
// MailLogFormatter はメール送信ログのJSONフォーマッター
type MailLogFormatter struct{}

// Format はlogrus.Formatterインターフェースを実装
// このフォーマッターは使用しないが、logrusのインターフェースに適合させるため定義
func (f *MailLogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
    // 実際のJSON出力はLogMailメソッド内で直接行うため、ここでは使用しない
    return []byte{}, nil
}
```

### 3.2 EmailServiceへの統合設計

#### 3.2.1 EmailService構造体の拡張

```go
// EmailService はメール送信サービス
type EmailService struct {
    sender EmailSender
    logger *logging.MailLogger  // 新規追加
}

// NewEmailService は新しいEmailServiceを作成
func NewEmailService(cfg *config.EmailConfig, mailLogger *logging.MailLogger) (*EmailService, error) {
    // 既存の送信実装選択ロジックは変更なし
    senderType := cfg.SenderType
    if senderType == "" {
        appEnv := os.Getenv("APP_ENV")
        switch appEnv {
        case "staging", "production":
            senderType = "ses"
        default:
            senderType = "mock"
        }
    }

    var sender EmailSender
    var err error

    switch senderType {
    case "mock":
        sender = NewMockSender()
    case "mailpit":
        from := cfg.SES.From
        if from == "" {
            from = "noreply@example.com"
        }
        sender = NewMailpitSender(cfg.Mailpit.SMTPHost, cfg.Mailpit.SMTPPort, from)
    case "ses":
        sender, err = NewSESSender(cfg.SES.Region, cfg.SES.From)
        if err != nil {
            return nil, fmt.Errorf("failed to create SES sender: %w", err)
        }
    default:
        return nil, fmt.Errorf("invalid sender type: %s", senderType)
    }

    return &EmailService{
        sender: sender,
        logger: mailLogger,  // 新規追加
    }, nil
}

// SendEmail はメールを送信します
func (s *EmailService) SendEmail(ctx context.Context, to []string, subject, body string) error {
    // 送信前のログ出力
    if s.logger != nil {
        // 送信前のログ（送信結果はまだ不明のため、success=falseで記録）
        s.logger.LogMail(to, subject, body, s.getSenderType(), false, nil)
    }

    // メール送信
    err := s.sender.Send(ctx, to, subject, body)
    
    // 送信後のログ出力
    if s.logger != nil {
        success := err == nil
        s.logger.LogMail(to, subject, body, s.getSenderType(), success, err)
    }

    return err
}

// getSenderType は送信実装の種類を取得
func (s *EmailService) getSenderType() string {
    switch s.sender.(type) {
    case *MockSender:
        return "mock"
    case *MailpitSender:
        return "mailpit"
    case *SESSender:
        return "ses"
    default:
        return "unknown"
    }
}
```

### 3.3 設定構造体の拡張設計

#### 3.3.1 LoggingConfig構造体の拡張

```go
// LoggingConfig はロギング設定
type LoggingConfig struct {
    Level           string `mapstructure:"level"`
    Format          string `mapstructure:"format"`
    Output          string `mapstructure:"output"`
    OutputDir       string `mapstructure:"output_dir"`
    SQLLogEnabled   bool   `mapstructure:"sql_log_enabled"`      // SQLログの有効/無効（オプション）
    SQLLogOutputDir string `mapstructure:"sql_log_output_dir"`   // SQLログ出力先ディレクトリ（オプション）
    MailLogOutputDir string `mapstructure:"mail_log_output_dir"` // 新規追加: メール送信ログ出力先ディレクトリ（オプション）
}
```

#### 3.3.2 設定ファイルの読み込み処理

```go
// config.Load()内で処理
if cfg.Logging.MailLogOutputDir == "" {
    cfg.Logging.MailLogOutputDir = cfg.Logging.OutputDir
}
```

### 3.4 環境別制御の設計

#### 3.4.1 環境判定処理

```go
// server/cmd/server/main.go または EmailService初期化時

// 環境判定
appEnv := os.Getenv("APP_ENV")
if appEnv == "" {
    appEnv = "develop" // デフォルトはdevelop
}

// メール送信ログの有効/無効判定
mailLogEnabled := appEnv == "develop" || appEnv == "staging"

// MailLoggerの作成
mailLogger, err := logging.NewMailLogger(cfg.Logging.MailLogOutputDir, mailLogEnabled)
if err != nil {
    log.Printf("Warning: Failed to initialize mail logger: %v", err)
    log.Println("Mail logging will be disabled")
    mailLogger = nil
}
```

## 4. データモデル

### 4.1 MailLogEntry構造体

```go
type MailLogEntry struct {
    Timestamp     string   `json:"timestamp"`      // 送信時刻（YYYY-MM-DD HH:MM:SS形式）
    To            []string `json:"to"`            // 送信先メールアドレスリスト
    Subject       string   `json:"subject"`       // メール件名
    Body          string   `json:"body"`           // メール本文（200文字制限）
    BodyTruncated bool     `json:"body_truncated"` // 本文が切り捨てられたかどうか
    SenderType    string   `json:"sender_type"`   // 送信実装の種類（mock/mailpit/ses）
    Success       bool     `json:"success"`        // 送信成功/失敗
    Error         string   `json:"error,omitempty"` // エラーメッセージ（エラー時のみ）
}
```

### 4.2 JSON出力例

#### 4.2.1 送信成功の場合
```json
{
  "timestamp": "2025-01-27 14:30:45",
  "to": ["recipient@example.com"],
  "subject": "Welcome to our service",
  "body": "This is a welcome email...",
  "body_truncated": false,
  "sender_type": "ses",
  "success": true
}
```

#### 4.2.2 送信失敗の場合
```json
{
  "timestamp": "2025-01-27 14:30:45",
  "to": ["recipient@example.com"],
  "subject": "Welcome to our service",
  "body": "This is a welcome email...",
  "body_truncated": false,
  "sender_type": "ses",
  "success": false,
  "error": "SES error: Invalid email address"
}
```

#### 4.2.3 本文が200文字以上の場合
```json
{
  "timestamp": "2025-01-27 14:30:45",
  "to": ["recipient@example.com"],
  "subject": "Welcome to our service",
  "body": "This is a very long email body that contains more than 200 characters. When the body is longer than 200 characters, it will be truncated and the truncated flag will be set to true. This helps to keep log files manageable and prevents them from growing too large...",
  "body_truncated": true,
  "sender_type": "ses",
  "success": true
}
```

## 5. エラーハンドリング

### 5.1 ログ出力エラーの処理

#### 5.1.1 ログ出力失敗時の処理
- ログ出力に失敗してもメール送信処理は継続する
- ログ出力エラーは標準エラー出力に記録する
- メール送信のエラーとして扱わない

```go
// LogMailメソッド内でのエラーハンドリング
func (m *MailLogger) LogMail(...) {
    // JSONエンコードに失敗した場合
    jsonBytes, err := json.Marshal(entry)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to encode mail log: %v\n", err)
        return // メール送信処理は継続
    }

    // ログファイルへの書き込みに失敗した場合
    if _, err := m.writer.Write(append(jsonBytes, '\n')); err != nil {
        fmt.Fprintf(os.Stderr, "Failed to write mail log: %v\n", err)
        // メール送信処理は継続
    }
}
```

### 5.2 ログディレクトリ作成エラーの処理

#### 5.2.1 ディレクトリ作成失敗時の処理
- ディレクトリ作成に失敗した場合はエラーを返す
- MailLoggerの作成に失敗した場合は、nilを返してログ機能を無効化
- メール送信機能自体は利用可能

```go
// NewMailLogger内でのエラーハンドリング
func NewMailLogger(outputDir string, enabled bool) (*MailLogger, error) {
    if !enabled {
        return nil, nil
    }

    // ディレクトリ作成に失敗した場合はエラーを返す
    if err := os.MkdirAll(outputDir, 0755); err != nil {
        return nil, fmt.Errorf("failed to create log directory: %w", err)
    }

    // 以降の処理...
}
```

### 5.3 メール送信エラーの処理

#### 5.3.1 送信失敗時のログ出力
- メール送信に失敗した場合もログに記録する
- エラーメッセージを含めてログに出力する
- `success`フィールドを`false`に設定する

```go
// EmailService.SendEmail内でのエラーハンドリング
func (s *EmailService) SendEmail(ctx context.Context, to []string, subject, body string) error {
    // メール送信
    err := s.sender.Send(ctx, to, subject, body)
    
    // 送信後のログ出力（成功・失敗に関わらず）
    if s.logger != nil {
        success := err == nil
        s.logger.LogMail(to, subject, body, s.getSenderType(), success, err)
    }

    return err
}
```

## 6. テスト戦略

### 6.1 ユニットテスト

#### 6.1.1 MailLoggerのテスト
- `NewMailLogger`のテスト
  - 有効化時の正常系テスト
  - 無効化時のnil返却テスト
  - ディレクトリ作成失敗時のエラーテスト
- `LogMail`のテスト
  - 正常なログ出力のテスト
  - メール本文の200文字制限テスト
  - JSONエンコードエラーのテスト
  - ログファイル書き込みエラーのテスト

#### 6.1.2 EmailServiceのテスト
- `SendEmail`のテスト
  - ログ出力が正しく実行されることのテスト
  - 送信成功時のログ内容のテスト
  - 送信失敗時のログ内容のテスト
  - ログ出力失敗時もメール送信が継続することのテスト

### 6.2 統合テスト

#### 6.2.1 メール送信ログ出力の統合テスト
- 実際のメール送信処理とログ出力の統合テスト
- 各送信実装（MockSender、MailpitSender、SESSender）でのログ出力テスト
- 環境別制御のテスト（develop/staging/production）

### 6.3 E2Eテスト

#### 6.3.1 メール送信APIのE2Eテスト
- メール送信APIエンドポイントのテスト
- ログファイルに正しく記録されることのテスト
- JSON形式で正しく出力されることのテスト

## 7. 実装の詳細

### 7.1 メール本文の切り捨て処理

```go
// LogMailメソッド内での処理
bodyTruncated := false
if len(body) > 200 {
    body = body[:200] + "..."
    bodyTruncated = true
}
```

### 7.2 日付別ファイル分割

```go
// lumberjackの設定
logFileName := fmt.Sprintf("mail-%s.log", time.Now().Format("2006-01-02"))
writer := &lumberjack.Logger{
    Filename:   filepath.Join(outputDir, logFileName),
    MaxSize:    0,     // サイズ制限なし（日付別分割のみ使用）
    MaxBackups: 0,     // バックアップ保持数なし
    MaxAge:     0,     // 保持日数なし
    LocalTime:  true,  // ローカルタイムゾーンを使用
    Compress:   false, // 圧縮なし
}
```

### 7.3 環境別制御の実装

```go
// 環境判定
appEnv := os.Getenv("APP_ENV")
if appEnv == "" {
    appEnv = "develop" // デフォルトはdevelop
}

// メール送信ログの有効/無効判定
mailLogEnabled := appEnv == "develop" || appEnv == "staging"

// MailLoggerの作成
mailLogger, err := logging.NewMailLogger(cfg.Logging.MailLogOutputDir, mailLogEnabled)
```

## 8. 既存コードへの影響

### 8.1 変更が必要なファイル

1. **`server/internal/logging/mail_logger.go`** (新規作成)
   - MailLogger構造体の実装
   - MailLogEntry構造体の定義
   - LogMailメソッドの実装

2. **`server/internal/config/config.go`** (変更)
   - LoggingConfig構造体に`MailLogOutputDir`フィールドを追加
   - 設定ファイル読み込み処理で`MailLogOutputDir`のデフォルト値を設定

3. **`server/internal/service/email/email_service.go`** (変更)
   - EmailService構造体に`logger`フィールドを追加
   - NewEmailService関数に`mailLogger`パラメータを追加
   - SendEmailメソッドにログ出力処理を追加
   - getSenderTypeメソッドを追加

4. **`server/cmd/server/main.go`** (変更)
   - MailLoggerの初期化処理を追加
   - EmailService作成時にMailLoggerを渡す

### 8.2 変更が不要なファイル

- `server/internal/service/email/email_sender.go` (変更不要)
- `server/internal/service/email/mock_sender.go` (変更不要)
- `server/internal/service/email/mailpit_sender.go` (変更不要)
- `server/internal/service/email/ses_sender.go` (変更不要)

## 9. パフォーマンス考慮事項

### 9.1 ログ出力のパフォーマンス
- ログ出力は同期的に処理される
- メール送信処理のパフォーマンスに影響を与えないよう、ログ出力は軽量に実装
- JSONエンコードとファイル書き込みは高速に処理される

### 9.2 メモリ使用量
- メール本文の200文字制限により、メモリ使用量を抑制
- ログエントリは1行ずつ出力するため、メモリに大量のログを保持しない

## 10. セキュリティ考慮事項

### 10.1 機密情報の保護
- メール本文を200文字で切り捨てることで、機密情報の漏洩リスクを軽減
- 本番環境ではログを出力しないことで、機密情報の漏洩リスクを回避

### 10.2 ログファイルのアクセス制御
- ログファイルは適切なファイル権限で作成される（0755）
- ログファイルの内容は機密情報を含む可能性があるため、適切に管理する
