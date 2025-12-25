# ログ出力機能設計書

## 1. 概要

### 1.1 設計の目的
要件定義書に基づき、APIサーバーと管理画面サーバーのアクセスログを日付別ファイルに出力する機能の詳細設計を定義する。`logrus`と`lumberjack`ライブラリを活用し、既存のアーキテクチャに統合する。

### 1.2 設計の範囲
- アクセスログ出力機能のアーキテクチャ設計
- HTTPミドルウェアの設計
- ログライブラリ（logrus + lumberjack）の統合設計
- 設定構造体の拡張設計
- 日付別ファイル分割の実装設計
- エラーハンドリング設計
- テスト戦略

## 2. アーキテクチャ設計

### 2.1 ディレクトリ構造

#### 2.1.1 変更前の構造
```
server/
├── cmd/
│   ├── server/
│   │   └── main.go          # サーバー起動コマンド
│   └── admin/
│       └── main.go          # 管理画面起動コマンド
├── internal/
│   ├── api/
│   │   ├── handler/
│   │   └── router/
│   ├── config/
│   └── ...
└── ...
```

#### 2.1.2 変更後の構造
```
server/
├── cmd/
│   ├── server/
│   │   └── main.go          # サーバー起動コマンド（アクセスログミドルウェア追加）
│   └── admin/
│       └── main.go          # 管理画面起動コマンド（アクセスログミドルウェア追加）
├── internal/
│   ├── api/
│   │   ├── handler/
│   │   └── router/           # アクセスログミドルウェア統合
│   ├── config/               # LoggingConfig拡張
│   ├── logging/              # 新規パッケージ
│   │   ├── access_logger.go # アクセスログ出力機能
│   │   └── middleware.go     # HTTPミドルウェア
│   └── ...
├── logs/                     # 新規ディレクトリ
│   └── .gitkeep             # ディレクトリ保持用
└── ...
```

### 2.2 アクセスログ出力の実行フロー

```
┌─────────────────────────────────────────────────────────────┐
│                    HTTPリクエスト受信                        │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              1. アクセスログミドルウェア起動                  │
│              - リクエスト開始時刻を記録                      │
│              - リクエスト情報を取得                          │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              2. 次のハンドラーを実行                          │
│              - 実際のHTTPハンドラー処理                       │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              3. レスポンス情報を取得                           │
│              - ステータスコード                              │
│              - レスポンス時間（ミリ秒）                      │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              4. アクセスログの出力                            │
│              - logrusでログエントリを作成                    │
│              - lumberjackで日付別ファイルに出力              │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              5. レスポンス返却                                │
└─────────────────────────────────────────────────────────────┘
```

### 2.3 既存アーキテクチャとの統合

```
┌─────────────────────────────────────────────────────────────┐
│              HTTPリクエスト                                  │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              Router (internal/api/router)                   │
│              - アクセスログミドルウェア（新規追加）           │
│              - CORSミドルウェア（既存）                      │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              Handler (internal/api/handler)                 │
│              - UserHandler                                  │
│              - PostHandler                                  │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              Service (internal/service)                     │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              Repository (internal/repository)               │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              AccessLogger (internal/logging)                │
│              - logrus + lumberjack                           │
│              - 日付別ファイル出力                            │
└─────────────────────────────────────────────────────────────┘
```

## 3. コンポーネント設計

### 3.1 AccessLoggerの設計

#### 3.1.1 パッケージ構造
```go
package logging

import (
    "os"
    "path/filepath"
    
    "github.com/sirupsen/logrus"
    "gopkg.in/natefinch/lumberjack.v2"
)
```

#### 3.1.2 AccessLogger構造体の設計

```go
// AccessLogger はアクセスログを出力するロガー
type AccessLogger struct {
    logger *logrus.Logger
    writer *lumberjack.Logger
}

// NewAccessLogger は新しいAccessLoggerを作成
// outputDirは絶対パスと相対パスの両方をサポートします
// - 相対パス: サーバーの実行ディレクトリからの相対パス（例: "logs"）
// - 絶対パス: システムの絶対パス（例: "/var/log/go-webdb-template"）
func NewAccessLogger(logType string, outputDir string) (*AccessLogger, error) {
    // 1. 出力ディレクトリの作成
    // os.MkdirAllは絶対パスと相対パスの両方をサポート
    if err := os.MkdirAll(outputDir, 0755); err != nil {
        return nil, fmt.Errorf("failed to create log directory: %w", err)
    }
    
    // 2. lumberjackの設定（日付別ファイル分割）
    // filepath.Joinは絶対パスと相対パスの両方を正しく処理
    writer := &lumberjack.Logger{
        Filename:   filepath.Join(outputDir, logType+"-access-2006-01-02.log"),
        MaxSize:    0,  // サイズ制限なし（日付別分割のみ使用）
        MaxBackups: 0,  // バックアップ保持数なし（将来の拡張用）
        MaxAge:     0,  // 保持日数なし（将来の拡張用）
        LocalTime:  true, // ローカルタイムゾーンを使用
        Compress:   false, // 圧縮なし
    }
    
    // 3. logrusの設定
    logger := logrus.New()
    logger.SetOutput(writer)
    logger.SetFormatter(&logrus.TextFormatter{
        DisableColors:   true,
        FullTimestamp:   true,
        TimestampFormat: "2006-01-02 15:04:05",
    })
    logger.SetLevel(logrus.InfoLevel)
    
    return &AccessLogger{
        logger: logger,
        writer: writer,
    }, nil
}

// LogAccess はアクセスログを出力
func (a *AccessLogger) LogAccess(method, path, protocol string, statusCode int, responseTimeMs float64, remoteIP, userAgent string) {
    a.logger.WithFields(logrus.Fields{
        "method":       method,
        "path":         path,
        "protocol":     protocol,
        "status_code":  statusCode,
        "response_time_ms": responseTimeMs,
        "remote_ip":    remoteIP,
        "user_agent":   userAgent,
    }).Info("access")
}

// Close はロガーをクローズ
func (a *AccessLogger) Close() error {
    if a.writer != nil {
        return a.writer.Close()
    }
    return nil
}
```

#### 3.1.3 ログエントリのフォーマット

logrusのTextFormatterを使用した場合の出力例：
```
time="2025-01-27 14:30:45" level=info msg=access method=GET path=/api/users protocol=HTTP/1.1 status_code=200 response_time_ms=15.2 remote_ip=192.168.1.100 user_agent="Mozilla/5.0..."
```

要件定義書の形式に合わせるため、カスタムフォーマッターを実装：
```go
// CustomTextFormatter はカスタムテキストフォーマッター
type CustomTextFormatter struct {
    logrus.TextFormatter
}

// Format はログエントリをフォーマット
func (f *CustomTextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
    timestamp := entry.Time.Format("2006-01-02 15:04:05")
    method := entry.Data["method"].(string)
    path := entry.Data["path"].(string)
    protocol := entry.Data["protocol"].(string)
    statusCode := entry.Data["status_code"].(int)
    responseTimeMs := entry.Data["response_time_ms"].(float64)
    remoteIP := entry.Data["remote_ip"].(string)
    userAgent := entry.Data["user_agent"].(string)
    
    logLine := fmt.Sprintf("[%s] %s %s %s %d %.1fms %s \"%s\"\n",
        timestamp, method, path, protocol, statusCode, responseTimeMs, remoteIP, userAgent)
    
    return []byte(logLine), nil
}
```

### 3.2 HTTPミドルウェアの設計

#### 3.2.1 AccessLogMiddleware構造体の設計

```go
// AccessLogMiddleware はアクセスログを記録するHTTPミドルウェア
type AccessLogMiddleware struct {
    accessLogger *AccessLogger
}

// NewAccessLogMiddleware は新しいAccessLogMiddlewareを作成
func NewAccessLogMiddleware(accessLogger *AccessLogger) *AccessLogMiddleware {
    return &AccessLogMiddleware{
        accessLogger: accessLogger,
    }
}

// Middleware はHTTPミドルウェア関数を返す
func (m *AccessLogMiddleware) Middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // リクエスト開始時刻を記録
        startTime := time.Now()
        
        // レスポンスライターをラップしてステータスコードを取得
        rw := &responseWriter{
            ResponseWriter: w,
            statusCode:     http.StatusOK,
        }
        
        // 次のハンドラーを実行
        next.ServeHTTP(rw, r)
        
        // レスポンス時間を計算
        responseTime := time.Since(startTime)
        responseTimeMs := float64(responseTime.Nanoseconds()) / 1000000.0
        
        // リモートIPアドレスを取得
        remoteIP := r.RemoteAddr
        if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
            remoteIP = forwarded
        }
        
        // User-Agentを取得
        userAgent := r.Header.Get("User-Agent")
        if userAgent == "" {
            userAgent = "-"
        }
        
        // アクセスログを出力
        m.accessLogger.LogAccess(
            r.Method,
            r.URL.Path,
            r.Proto,
            rw.statusCode,
            responseTimeMs,
            remoteIP,
            userAgent,
        )
    })
}

// responseWriter はステータスコードを記録するResponseWriterラッパー
type responseWriter struct {
    http.ResponseWriter
    statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
    rw.statusCode = code
    rw.ResponseWriter.WriteHeader(code)
}
```

### 3.3 設定構造体の拡張設計

#### 3.3.1 LoggingConfig構造体の拡張

```go
// LoggingConfig はロギング設定
type LoggingConfig struct {
    Level    string `mapstructure:"level"`
    Format   string `mapstructure:"format"`
    Output   string `mapstructure:"output"`
    OutputDir string `mapstructure:"output_dir"` // 新規追加
}
```

#### 3.3.2 設定ファイルの例

```yaml
logging:
  level: debug
  format: json
  output: file
  output_dir: logs  # デフォルト値
```

### 3.4 サーバー起動処理への統合

#### 3.4.1 APIサーバー（server/cmd/server/main.go）への統合

```go
func main() {
    // ... 既存の初期化処理 ...
    
    // Routerの初期化
    r := router.NewRouter(userHandler, postHandler, cfg)
    
    // アクセスログの初期化
    accessLogger, err := logging.NewAccessLogger("api", cfg.Logging.OutputDir)
    if err != nil {
        log.Printf("Warning: Failed to initialize access logger: %v", err)
        log.Println("Access logging will be disabled")
    } else {
        defer accessLogger.Close()
        
        // アクセスログミドルウェアを追加
        accessLogMiddleware := logging.NewAccessLogMiddleware(accessLogger)
        r = accessLogMiddleware.Middleware(r).(http.Handler)
    }
    
    // HTTPサーバーの設定
    srv := &http.Server{
        Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
        Handler:      r,
        ReadTimeout:  cfg.Server.ReadTimeout,
        WriteTimeout: cfg.Server.WriteTimeout,
    }
    
    // ... 既存のサーバー起動処理 ...
}
```

#### 3.4.2 管理画面サーバー（server/cmd/admin/main.go）への統合

```go
func main() {
    // ... 既存の初期化処理 ...
    
    // 環境判定（production環境ではアクセスログを出力しない）
    env := os.Getenv("APP_ENV")
    if env == "" {
        env = "develop"
    }
    
    // Routerの初期化
    app := mux.NewRouter()
    
    // アクセスログの初期化（production環境以外）
    var accessLogger *logging.AccessLogger
    if env != "production" {
        var err error
        accessLogger, err = logging.NewAccessLogger("admin", cfg.Logging.OutputDir)
        if err != nil {
            log.Printf("Warning: Failed to initialize access logger: %v", err)
            log.Println("Access logging will be disabled")
        } else {
            defer accessLogger.Close()
            
            // アクセスログミドルウェアを追加
            accessLogMiddleware := logging.NewAccessLogMiddleware(accessLogger)
            app = accessLogMiddleware.Middleware(app).(http.Handler)
        }
    }
    
    // ... 既存のGoAdmin初期化処理 ...
    
    // HTTPサーバーの設定
    srv := &http.Server{
        Addr:         fmt.Sprintf(":%d", cfg.Admin.Port),
        Handler:      app,
        ReadTimeout:  cfg.Admin.ReadTimeout,
        WriteTimeout: cfg.Admin.WriteTimeout,
    }
    
    // ... 既存のサーバー起動処理 ...
}
```

#### 3.4.3 Routerへの統合（オプション）

`internal/api/router/router.go`でアクセスログミドルウェアを統合する方法：

```go
// NewRouter は新しいルーターを作成
func NewRouter(userHandler *handler.UserHandler, postHandler *handler.PostHandler, cfg *config.Config, accessLogger *logging.AccessLogger) http.Handler {
    r := mux.NewRouter()
    
    // ... 既存のルート定義 ...
    
    // アクセスログミドルウェアを追加
    var handler http.Handler = r
    if accessLogger != nil {
        accessLogMiddleware := logging.NewAccessLogMiddleware(accessLogger)
        handler = accessLogMiddleware.Middleware(r)
    }
    
    // CORS middleware
    c := cors.New(cors.Options{
        AllowedOrigins:   cfg.CORS.AllowedOrigins,
        AllowedMethods:   cfg.CORS.AllowedMethods,
        AllowedHeaders:   cfg.CORS.AllowedHeaders,
        AllowCredentials: true,
    })
    
    return c.Handler(handler)
}
```

## 4. データモデル

### 4.1 ログエントリの構造

#### 4.1.1 ログエントリの形式

```
[YYYY-MM-DD HH:MM:SS] METHOD /path HTTP/1.1 STATUS_CODE RESPONSE_TIME_MS REMOTE_IP USER_AGENT
```

#### 4.1.2 データ型

| フィールド | 型 | 説明 |
|----------|-----|------|
| timestamp | string | リクエスト日時（YYYY-MM-DD HH:MM:SS形式） |
| method | string | HTTPメソッド（GET, POST, PUT, DELETE等） |
| path | string | リクエストパス（URLパス） |
| protocol | string | HTTPプロトコル（HTTP/1.1等） |
| status_code | int | HTTPステータスコード |
| response_time_ms | float64 | レスポンス時間（ミリ秒） |
| remote_ip | string | リモートIPアドレス |
| user_agent | string | User-Agentヘッダー |

#### 4.1.3 ログファイル名の形式

- APIサーバー: `api-access-YYYY-MM-DD.log`
- 管理画面サーバー: `admin-access-YYYY-MM-DD.log`

例:
- `api-access-2025-01-27.log`
- `admin-access-2025-01-27.log`

### 4.2 lumberjackの日付フォーマット

lumberjackはGo言語の日付フォーマット（`2006-01-02`）を使用：

```go
Filename: filepath.Join(outputDir, logType+"-access-2006-01-02.log")
```

- `2006`: 年（4桁）
- `01`: 月（2桁）
- `02`: 日（2桁）

lumberjackが自動的に現在の日付に置き換えて、日付が変わったら新しいファイルに切り替える。

## 5. エラーハンドリング

### 5.1 ログディレクトリ作成のエラー

**エラーケース**:
- ディレクトリの作成に失敗（権限不足、ディスク容量不足等）

**処理**:
```go
if err := os.MkdirAll(outputDir, 0755); err != nil {
    return nil, fmt.Errorf("failed to create log directory: %w", err)
}
```

**フォールバック**:
- エラーが発生した場合は、標準エラー出力にエラーメッセージを出力
- アクセスログ機能を無効化し、HTTPリクエスト処理は継続

### 5.2 ログファイル書き込みのエラー

**エラーケース**:
- ログファイルへの書き込みに失敗（ディスク容量不足、権限不足等）

**処理**:
- `logrus`と`lumberjack`が自動的にエラーを処理
- ログファイルへの書き込みに失敗しても、HTTPリクエスト処理は継続
- エラーは標準エラー出力に出力される（logrusのデフォルト動作）

### 5.3 アクセスログ初期化のエラー

**エラーケース**:
- `NewAccessLogger()`の呼び出しに失敗

**処理**:
```go
accessLogger, err := logging.NewAccessLogger("api", cfg.Logging.OutputDir)
if err != nil {
    log.Printf("Warning: Failed to initialize access logger: %v", err)
    log.Println("Access logging will be disabled")
} else {
    defer accessLogger.Close()
    // ミドルウェアを追加
}
```

**フォールバック**:
- アクセスログ機能を無効化
- HTTPリクエスト処理は正常に継続
- 警告メッセージを標準エラー出力に出力

### 5.4 ミドルウェア実行時のエラー

**エラーケース**:
- ログ出力処理中にエラーが発生

**処理**:
- エラーが発生してもHTTPリクエスト処理は継続
- エラーは`logrus`が自動的に処理（標準エラー出力に出力）

## 6. 設定の拡張

### 6.1 LoggingConfig構造体の拡張

```go
// LoggingConfig はロギング設定
type LoggingConfig struct {
    Level     string `mapstructure:"level"`
    Format    string `mapstructure:"format"`
    Output    string `mapstructure:"output"`
    OutputDir string `mapstructure:"output_dir"` // 新規追加
}
```

### 6.2 設定ファイルの更新

#### 6.2.1 output_dirのパス指定について

`output_dir`は**絶対パスと相対パスの両方をサポート**します：
- **相対パス**: プロジェクトルートからの相対パス（例: `logs`, `./logs`, `../logs`）
- **絶対パス**: システムの絶対パス（例: `/var/log/go-webdb-template`, `C:\logs\go-webdb-template`）

相対パスの場合は、サーバーの実行ディレクトリからの相対パスとして解釈されます。

#### 6.2.2 develop環境（config/develop/config.yaml）

```yaml
logging:
  level: debug
  format: json
  output: file
  output_dir: logs  # 相対パス（デフォルト値）
```

#### 6.2.3 staging環境（config/staging/config.yaml）

```yaml
logging:
  level: info
  format: json
  output: file
  output_dir: logs  # 相対パス
```

#### 6.2.4 production環境（config/production/config.yaml.example）

```yaml
logging:
  level: info
  format: json
  output: file
  output_dir: /var/log/go-webdb-template  # 絶対パス（本番環境では絶対パスを推奨）
```

### 6.3 デフォルト値の設定

`config.Load()`関数でデフォルト値を設定：

```go
// Load は指定された環境の設定ファイルを読み込む
func Load() (*Config, error) {
    // ... 既存の設定読み込み処理 ...
    
    var cfg Config
    if err := viper.Unmarshal(&cfg); err != nil {
        return nil, fmt.Errorf("failed to unmarshal config: %w", err)
    }
    
    // デフォルト値の設定
    if cfg.Logging.OutputDir == "" {
        cfg.Logging.OutputDir = "logs"
    }
    
    return &cfg, nil
}
```

## 7. テスト戦略

### 7.1 ユニットテスト

#### 7.1.1 AccessLoggerのテスト

**テストケース**:
1. **正常系**: アクセスログが正常に出力される
2. **正常系**: 日付が変わったら新しいファイルに切り替わる
3. **異常系**: ディレクトリ作成に失敗した場合、エラーを返す
4. **正常系**: ログエントリのフォーマットが正しい

**テスト実装場所**:
- `server/internal/logging/access_logger_test.go`

**テスト実装例**:
```go
func TestAccessLogger_LogAccess(t *testing.T) {
    // 一時ディレクトリを作成
    tmpDir := t.TempDir()
    
    // AccessLoggerを作成
    logger, err := NewAccessLogger("api", tmpDir)
    if err != nil {
        t.Fatalf("Failed to create access logger: %v", err)
    }
    defer logger.Close()
    
    // アクセスログを出力
    logger.LogAccess("GET", "/api/users", "HTTP/1.1", 200, 15.2, "192.168.1.100", "Mozilla/5.0")
    
    // ログファイルの内容を確認
    // ... (実装)
}
```

#### 7.1.2 AccessLogMiddlewareのテスト

**テストケース**:
1. **正常系**: ミドルウェアが正常に動作する
2. **正常系**: レスポンス時間が正しく記録される
3. **正常系**: ステータスコードが正しく記録される
4. **正常系**: リモートIPアドレスが正しく取得される

**テスト実装場所**:
- `server/internal/logging/middleware_test.go`

### 7.2 統合テスト

#### 7.2.1 アクセスログ機能の統合テスト

**テストケース**:
1. **正常系**: APIサーバーへのリクエストがログファイルに記録される
2. **正常系**: 管理画面サーバーへのリクエストがログファイルに記録される（production環境以外）
3. **正常系**: production環境では管理画面サーバーのログが出力されない
4. **正常系**: 日付が変わったら新しいファイルに切り替わる
5. **正常系**: ログファイル名に日付が含まれる

**テスト実装場所**:
- `server/test/integration/access_log_test.go`

### 7.3 E2Eテスト

#### 7.3.1 アクセスログ機能のE2Eテスト

**テストケース**:
1. サーバー起動時にアクセスログが正常に初期化される
2. HTTPリクエストがログファイルに記録される
3. ログファイルの内容が正しい形式である
4. ログファイルへの書き込みエラーが発生してもHTTPリクエスト処理は継続する

**テスト実装場所**:
- 既存のE2Eテストスイートに追加

## 8. 実装上の注意事項

### 8.1 ログライブラリの依存関係

**重要なポイント**:
- `go.mod`に以下の依存関係を追加：
  - `github.com/sirupsen/logrus`
  - `gopkg.in/natefinch/lumberjack.v2`
- `go mod tidy`を実行して依存関係を更新

### 8.2 日付別ファイル分割の実装

**重要なポイント**:
- `lumberjack`の`Filename`に日付フォーマット（`2006-01-02`）を含める
- `LocalTime: true`を設定してローカルタイムゾーンを使用
- ライブラリが自動的に日付変更を検知して新しいファイルに切り替える

### 8.3 ログフォーマットの実装

**重要なポイント**:
- 要件定義書の形式に合わせてカスタムフォーマッターを実装
- テキスト形式で人間が読みやすい形式にする
- 将来的にJSON形式への拡張も検討可能な設計にする

### 8.4 ミドルウェアの統合

**重要なポイント**:
- 既存のCORSミドルウェアと競合しないようにする
- ミドルウェアの実行順序を考慮する（アクセスログ → CORS）
- エラーが発生してもHTTPリクエスト処理は継続する

### 8.5 環境別制御

**重要なポイント**:
- production環境では管理画面サーバーのログを出力しない
- 環境判定は`APP_ENV`環境変数または設定ファイルから取得
- デフォルト値は`develop`環境とする

### 8.6 ディレクトリ管理

**重要なポイント**:
- `logs`ディレクトリが存在しない場合は自動作成
- `logs/.gitkeep`ファイルを追加してディレクトリを保持
- `.gitignore`で`logs/*`と`logs/**`を除外（`logs`ディレクトリ自体は除外しない）

### 8.7 パフォーマンス考慮

**重要なポイント**:
- ログ出力は同期的に行う（シンプルさを優先）
- ログファイルへの書き込みエラーが発生してもHTTPリクエスト処理には影響を与えない
- 将来的に非同期書き込みやバッファリングを検討する場合は、設計時に拡張可能な構造にする

## 9. 参考情報

### 9.1 関連Issue
- GitHub Issue #10: ログ出力機能の組み込み

### 9.2 既存ドキュメント
- `server/cmd/server/main.go`: サーバー起動コマンドの実装例
- `server/cmd/admin/main.go`: 管理画面サーバー起動コマンドの実装例
- `server/internal/config/config.go`: 設定読み込み実装
- `server/internal/api/router/router.go`: ルーター実装

### 9.3 ログライブラリ
- **logrus**: https://github.com/sirupsen/logrus
  - ドキュメント: https://github.com/sirupsen/logrus#readme
  - テキストフォーマッター: https://pkg.go.dev/github.com/sirupsen/logrus#TextFormatter
- **lumberjack**: https://github.com/natefinch/lumberjack
  - ドキュメント: https://pkg.go.dev/gopkg.in/natefinch/lumberjack.v2
  - 日付フォーマット: `Filename`に`2006-01-02`形式を含める

### 9.4 Go言語標準ライブラリ
- `os`パッケージ: https://pkg.go.dev/os
- `time`パッケージ: https://pkg.go.dev/time
- `net/http`パッケージ: https://pkg.go.dev/net/http
- `path/filepath`パッケージ: https://pkg.go.dev/path/filepath

