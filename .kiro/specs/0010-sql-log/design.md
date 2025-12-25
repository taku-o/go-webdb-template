# SQLログ出力機能設計書

## 1. 概要

### 1.1 設計の目的
要件定義書に基づき、データベースに発行するSQLクエリをログに出力する機能の詳細設計を定義する。GORMのLoggerインターフェースを実装し、既存の`logrus`と`lumberjack`ライブラリを活用して、既存のアーキテクチャに統合する。

### 1.2 設計の範囲
- GORM Loggerインターフェースの実装設計
- SQLログ出力機能のアーキテクチャ設計
- DSN文字列からの機密情報フィルタリング設計
- 環境別制御（develop/staging/production）の設計
- 日付別ログファイル分割の実装設計
- 設定構造体の拡張設計
- シャーディング対応（複数シャード）の設計
- エラーハンドリング設計
- テスト戦略

### 1.3 設計方針
- **既存ライブラリの活用**: `logrus`と`lumberjack`ライブラリを既存のアクセスログ機能と同じパターンで使用
- **GORM Loggerインターフェース**: GORMの標準インターフェース（`gorm.io/gorm/logger.Interface`）を実装
- **環境別制御**: 開発環境とステージング環境でのみSQLログを出力し、本番環境では出力しない
- **セキュリティ**: DSN文字列からpasswordなどの機密情報をフィルタリング
- **既存機能との統合**: 既存のアクセスログ機能と同じ`logs`ディレクトリを使用
- **シャーディング対応**: すべてのシャードに対してLoggerを設定

## 2. アーキテクチャ設計

### 2.1 ディレクトリ構造

#### 2.1.1 変更前の構造
```
server/
├── internal/
│   ├── db/
│   │   ├── connection.go      # GORM接続作成（Logger未設定）
│   │   └── manager.go         # GORMManager
│   ├── config/
│   │   └── config.go          # LoggingConfig（OutputDirのみ）
│   └── logging/
│       └── access_logger.go   # アクセスログ機能
└── ...
```

#### 2.1.2 変更後の構造
```
server/
├── internal/
│   ├── db/
│   │   ├── connection.go      # GORM接続作成（Logger設定追加）
│   │   ├── manager.go         # GORMManager（変更不要）
│   │   └── logger.go          # 新規: GORM Logger実装
│   ├── config/
│   │   └── config.go          # LoggingConfig拡張（SQLLogEnabled, SQLLogOutputDir追加）
│   └── logging/
│       └── access_logger.go   # アクセスログ機能（変更不要）
└── ...
```

### 2.2 SQLログ出力の実行フロー

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
│              3. SQL Logger初期化（環境判定）                    │
│              - develop/staging: SQL Logger有効化             │
│              - production: SQL Logger無効化                 │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              4. GORMManager初期化                            │
│              NewGORMManager(cfg)                            │
│              - 各シャードに対してGORM接続を作成               │
│              - 各接続にSQL Loggerを設定                      │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              5. SQLクエリ実行                                 │
│              gormDB.Find(), Create(), Update()等            │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              6. GORM Logger.Trace()呼び出し                   │
│              - SQLクエリ情報を取得                           │
│              - 実行時間を計算                                │
│              - 結果件数を取得                                │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              7. SQLログ出力                                   │
│              - logrusでログエントリを作成                    │
│              - lumberjackで日付別ファイルに出力                │
│              - フォーマット: [YYYY-MM-DD HH:MM:SS] [SHARD_ID] │
│                [DRIVER] [TABLE] ROWS_AFFECTED | SQL_QUERY | │
│                DURATION_MS                                   │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              8. SQLクエリ結果返却                             │
└─────────────────────────────────────────────────────────────┘
```

### 2.3 既存アーキテクチャとの統合

#### 2.3.1 GORM接続作成処理への統合
- `server/internal/db/connection.go`の`createGORMConnection`関数を拡張
- GORM接続作成時に、環境判定とSQL Logger設定を追加
- 既存の接続プール設定は変更しない

#### 2.3.2 設定構造体への統合
- `server/internal/config/config.go`の`LoggingConfig`構造体を拡張
- `SQLLogEnabled`と`SQLLogOutputDir`フィールドを追加（オプション）
- 既存の`OutputDir`フィールドは維持

#### 2.3.3 ログ出力機能との統合
- 既存のアクセスログ機能（`server/internal/logging/access_logger.go`）と同じパターンで実装
- 同じ`logs`ディレクトリを使用
- 同じ`lumberjack`と`logrus`ライブラリを使用

## 3. コンポーネント設計

### 3.1 SQLLogger構造体

#### 3.1.1 構造体定義
```go
// server/internal/db/logger.go

package db

import (
    "context"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "regexp"
    "strings"
    "time"

    "github.com/example/go-webdb-template/internal/config"
    "github.com/sirupsen/logrus"
    "gorm.io/gorm/logger"
    "gopkg.in/natefinch/lumberjack.v2"
)

// SQLLogger はGORMのLoggerインターフェースを実装
type SQLLogger struct {
    logger      *logrus.Logger
    writer      io.WriteCloser
    shardID     int
    driver      string
    logLevel    logger.LogLevel
    outputDir   string
}

// NewSQLLogger は新しいSQLLoggerを作成
func NewSQLLogger(shardID int, driver string, cfg *config.LoggingConfig, env string) (*SQLLogger, error) {
    // 環境判定: production環境ではLoggerを作成しない
    if env == "production" {
        return nil, nil // nilを返してLoggerを無効化
    }

    // 設定からログ出力先を取得
    outputDir := cfg.SQLLogOutputDir
    if outputDir == "" {
        outputDir = cfg.OutputDir
    }
    if outputDir == "" {
        outputDir = "logs" // デフォルト値
    }

    // 出力ディレクトリの作成
    if err := os.MkdirAll(outputDir, 0755); err != nil {
        return nil, fmt.Errorf("failed to create log directory: %w", err)
    }

    // lumberjackの設定（日付別ファイル分割）
    logFileName := fmt.Sprintf("sql-%s.log", time.Now().Format("2006-01-02"))
    writer := &lumberjack.Logger{
        Filename:   filepath.Join(outputDir, logFileName),
        MaxSize:    0,     // サイズ制限なし（日付別分割のみ使用）
        MaxBackups: 0,     // バックアップ保持数なし
        MaxAge:     0,     // 保持日数なし
        LocalTime:  true,  // ローカルタイムゾーンを使用
        Compress:   false, // 圧縮なし
    }

    // logrusの設定
    logger := logrus.New()
    logger.SetOutput(writer)
    logger.SetFormatter(&SQLTextFormatter{})
    logger.SetLevel(logrus.InfoLevel)

    return &SQLLogger{
        logger:    logger,
        writer:    writer,
        shardID:   shardID,
        driver:    driver,
        logLevel:  logger.Info,
        outputDir: outputDir,
    }, nil
}
```

#### 3.1.2 GORM Loggerインターフェースの実装
```go
// LogMode はログレベルを設定
func (l *SQLLogger) LogMode(level logger.LogLevel) logger.Interface {
    newLogger := *l
    newLogger.logLevel = level
    return &newLogger
}

// Info は情報ログを出力
func (l *SQLLogger) Info(ctx context.Context, msg string, data ...interface{}) {
    if l == nil || l.logLevel < logger.Info {
        return
    }
    l.logger.WithFields(logrus.Fields{
        "message": fmt.Sprintf(msg, data...),
    }).Info("sql")
}

// Warn は警告ログを出力
func (l *SQLLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
    if l == nil || l.logLevel < logger.Warn {
        return
    }
    l.logger.WithFields(logrus.Fields{
        "message": fmt.Sprintf(msg, data...),
    }).Warn("sql")
}

// Error はエラーログを出力
func (l *SQLLogger) Error(ctx context.Context, msg string, data ...interface{}) {
    if l == nil || l.logLevel < logger.Error {
        return
    }
    l.logger.WithFields(logrus.Fields{
        "message": fmt.Sprintf(msg, data...),
    }).Error("sql")
}

// Trace はSQLクエリトレースを出力（最重要メソッド）
func (l *SQLLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
    if l == nil || l.logLevel <= logger.Silent {
        return
    }

    // SQLクエリと結果件数を取得
    sql, rows := fc()
    duration := time.Since(begin)

    // テーブル名を抽出（SQLクエリから）
    tableName := extractTableName(sql)

    // ログエントリを作成
    l.logger.WithFields(logrus.Fields{
        "shard_id":      l.shardID,
        "driver":        l.driver,
        "table":         tableName,
        "rows_affected": rows,
        "sql":           sql,
        "duration_ms":   duration.Milliseconds(),
    }).Info("sql")
}

// Close はロガーをクローズ
func (l *SQLLogger) Close() error {
    if l == nil || l.writer == nil {
        return nil
    }
    return l.writer.Close()
}
```

### 3.2 DSNフィルタリング機能

#### 3.2.1 フィルタリング関数
```go
// filterDSN はDSN文字列からpasswordなどの機密情報をフィルタリング
func filterDSN(dsn string, driver string) string {
    switch driver {
    case "postgres":
        // PostgreSQL DSN: "host=xxx port=xxx user=xxx password=xxx dbname=xxx"
        re := regexp.MustCompile(`password=[^ ]+`)
        return re.ReplaceAllString(dsn, "password=***")
    case "mysql":
        // MySQL DSN: "user:password@tcp(host:port)/dbname"
        re := regexp.MustCompile(`:[^@]+@`)
        return re.ReplaceAllString(dsn, ":***@")
    case "sqlite3":
        // SQLite DSN: 通常password情報は含まれないが、念のためチェック
        return dsn
    default:
        return dsn
    }
}
```

### 3.3 テーブル名抽出機能

#### 3.3.1 抽出関数
```go
// extractTableName はSQLクエリからテーブル名を抽出
func extractTableName(sql string) string {
    sql = strings.TrimSpace(sql)
    sqlUpper := strings.ToUpper(sql)

    // FROM句から抽出
    if strings.HasPrefix(sqlUpper, "SELECT") {
        re := regexp.MustCompile(`(?i)FROM\s+(\w+)`)
        matches := re.FindStringSubmatch(sql)
        if len(matches) > 1 {
            return matches[1]
        }
    }

    // INSERT INTO句から抽出
    if strings.HasPrefix(sqlUpper, "INSERT") {
        re := regexp.MustCompile(`(?i)INTO\s+(\w+)`)
        matches := re.FindStringSubmatch(sql)
        if len(matches) > 1 {
            return matches[1]
        }
    }

    // UPDATE句から抽出
    if strings.HasPrefix(sqlUpper, "UPDATE") {
        re := regexp.MustCompile(`(?i)UPDATE\s+(\w+)`)
        matches := re.FindStringSubmatch(sql)
        if len(matches) > 1 {
            return matches[1]
        }
    }

    // DELETE FROM句から抽出
    if strings.HasPrefix(sqlUpper, "DELETE") {
        re := regexp.MustCompile(`(?i)FROM\s+(\w+)`)
        matches := re.FindStringSubmatch(sql)
        if len(matches) > 1 {
            return matches[1]
        }
    }

    return "unknown"
}
```

### 3.4 SQLTextFormatter構造体

#### 3.4.1 フォーマッター実装
```go
// SQLTextFormatter はSQLログのテキストフォーマッター
type SQLTextFormatter struct{}

// Format はログエントリをテキスト形式でフォーマット
func (f *SQLTextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
    timestamp := entry.Time.Format("2006-01-02 15:04:05")
    shardID := entry.Data["shard_id"]
    driver := entry.Data["driver"]
    table := entry.Data["table"]
    rowsAffected := entry.Data["rows_affected"]
    sql := entry.Data["sql"]
    durationMs := entry.Data["duration_ms"]

    logLine := fmt.Sprintf(
        "[%s] [%d] [%s] [%s] %d | %s | %dms\n",
        timestamp, shardID, driver, table, rowsAffected, sql, durationMs,
    )

    return []byte(logLine), nil
}
```

### 3.5 GORM接続作成処理への統合

#### 3.5.1 createGORMConnection関数の拡張
```go
// createGORMConnection はGORM接続を作成するヘルパー関数（拡張版）
func createGORMConnection(cfg *config.ShardConfig, isWriter bool, loggingCfg *config.LoggingConfig, env string) (*gorm.DB, error) {
    var dialector gorm.Dialector

    dsn := cfg.GetWriterDSN()
    if !isWriter && len(cfg.ReaderDSNs) > 0 {
        dsn = cfg.ReaderDSNs[0]
    }

    switch cfg.Driver {
    case "sqlite3":
        dialector = sqlite.Open(dsn)
    case "postgres":
        dialector = postgres.Open(dsn)
    case "mysql":
        dialector = mysql.Open(dsn)
    default:
        return nil, fmt.Errorf("unsupported driver: %s", cfg.Driver)
    }

    // SQL Loggerの作成
    sqlLogger, err := NewSQLLogger(cfg.ID, cfg.Driver, loggingCfg, env)
    if err != nil {
        // Logger作成エラーは警告のみ（SQLクエリ実行には影響しない）
        log.Printf("Warning: Failed to create SQL logger for shard %d: %v", cfg.ID, err)
    }

    // GORM ConfigにLoggerを設定
    gormConfig := &gorm.Config{}
    if sqlLogger != nil {
        gormConfig.Logger = sqlLogger
    }

    db, err := gorm.Open(dialector, gormConfig)
    if err != nil {
        return nil, err
    }

    // 接続プール設定
    sqlDB, err := db.DB()
    if err != nil {
        return nil, err
    }
    sqlDB.SetMaxOpenConns(cfg.MaxConnections)
    sqlDB.SetMaxIdleConns(cfg.MaxIdleConnections)
    sqlDB.SetConnMaxLifetime(cfg.ConnectionMaxLifetime)

    return db, nil
}
```

#### 3.5.2 NewGORMConnection関数の拡張
```go
// NewGORMConnection は新しいGORM接続を作成（拡張版）
func NewGORMConnection(cfg *config.ShardConfig, loggingCfg *config.LoggingConfig, env string) (*GORMConnection, error) {
    // 1. Writer接続を作成（Logger設定付き）
    writerDB, err := createGORMConnection(cfg, true, loggingCfg, env)
    if err != nil {
        return nil, fmt.Errorf("failed to create writer connection: %w", err)
    }

    // 2. Reader接続を作成（複数可）
    // ... （既存のコードと同じ）

    // 3. dbresolverプラグインを設定
    // ... （既存のコードと同じ）

    return &GORMConnection{
        DB:      writerDB,
        ShardID: cfg.ID,
        Driver:  cfg.Driver,
        config:  cfg,
    }, nil
}
```

#### 3.5.3 NewGORMManager関数の拡張
```go
// NewGORMManager は新しいGORM Managerを作成（拡張版）
func NewGORMManager(cfg *config.Config) (*GORMManager, error) {
    manager := &GORMManager{
        connections: make(map[int]*GORMConnection),
        strategy:    NewHashBasedSharding(len(cfg.Database.Shards)),
    }

    // 環境判定
    env := os.Getenv("APP_ENV")
    if env == "" {
        env = "develop"
    }

    // 各シャードへの接続を確立（Logger設定付き）
    for _, shardCfg := range cfg.Database.Shards {
        conn, err := NewGORMConnection(&shardCfg, &cfg.Logging, env)
        if err != nil {
            manager.CloseAll()
            return nil, fmt.Errorf("failed to create connection for shard %d: %w", shardCfg.ID, err)
        }
        manager.connections[shardCfg.ID] = conn
    }

    return manager, nil
}
```

## 4. データモデル

### 4.1 ログエントリ構造

#### 4.1.1 ログフォーマット
```
[YYYY-MM-DD HH:MM:SS] [SHARD_ID] [DRIVER] [TABLE] ROWS_AFFECTED | SQL_QUERY | DURATION_MS
```

#### 4.1.2 ログエントリの例
```
[2025-01-27 14:30:45] [1] [sqlite3] [users] 1 | SELECT * FROM users WHERE id = ? | 2.5ms
[2025-01-27 14:30:46] [2] [sqlite3] [posts] 1 | INSERT INTO posts (user_id, title, content) VALUES (?, ?, ?) | 3.2ms
[2025-01-27 14:30:47] [1] [sqlite3] [users] 0 | UPDATE users SET email = ? WHERE id = ? | 1.8ms
```

### 4.2 設定構造体の拡張

#### 4.2.1 LoggingConfig構造体
```go
// LoggingConfig はロギング設定
type LoggingConfig struct {
    Level          string `mapstructure:"level"`
    Format         string `mapstructure:"format"`
    Output         string `mapstructure:"output"`
    OutputDir      string `mapstructure:"output_dir"`        // 既存: アクセスログ用
    SQLLogEnabled  bool   `mapstructure:"sql_log_enabled"`   // 新規: SQLログ有効/無効（オプション）
    SQLLogOutputDir string `mapstructure:"sql_log_output_dir"` // 新規: SQLログ出力先（オプション）
}
```

#### 4.2.2 設定ファイル例
```yaml
# config/develop/config.yaml
logging:
  level: debug
  format: json
  output: file
  output_dir: logs  # アクセスログとSQLログの共通出力先
  sql_log_enabled: true  # オプション（環境別に自動判定する場合は省略可能）
  sql_log_output_dir: logs  # オプション（output_dirと同じ場合は省略可能）
```

## 5. エラーハンドリング設計

### 5.1 Logger初期化エラー
- Loggerの初期化に失敗した場合、警告を出力してLoggerなしで動作を継続
- SQLクエリ実行には影響を与えない

### 5.2 ログファイル書き込みエラー
- ログファイルへの書き込みエラーが発生した場合、標準エラー出力にエラーメッセージを出力
- SQLクエリ実行は継続する

### 5.3 DSNフィルタリングエラー
- DSNフィルタリングに失敗した場合、元のDSNをそのまま使用（セキュリティリスクを最小化）
- 警告を出力して処理を継続

### 5.4 テーブル名抽出エラー
- テーブル名の抽出に失敗した場合、"unknown"を返す
- SQLクエリ実行には影響を与えない

## 6. テスト戦略

### 6.1 単体テスト

#### 6.1.1 SQLLoggerのテスト
- `NewSQLLogger`関数のテスト
  - 正常系: Loggerの作成成功
  - 異常系: ディレクトリ作成失敗
  - 環境判定: production環境ではnilを返す
- `Trace`メソッドのテスト
  - SQLクエリのログ出力
  - 実行時間の記録
  - 結果件数の記録
- `LogMode`メソッドのテスト
  - ログレベルの変更

#### 6.1.2 DSNフィルタリングのテスト
- PostgreSQL DSNのフィルタリング
- MySQL DSNのフィルタリング
- SQLite DSNのフィルタリング（変更なし）

#### 6.1.3 テーブル名抽出のテスト
- SELECT文からのテーブル名抽出
- INSERT文からのテーブル名抽出
- UPDATE文からのテーブル名抽出
- DELETE文からのテーブル名抽出
- 抽出失敗時の"unknown"返却

### 6.2 統合テスト

#### 6.2.1 GORM接続作成のテスト
- Logger設定付きGORM接続の作成
- 複数シャードへのLogger設定
- 環境別のLogger有効/無効化

#### 6.2.2 SQLログ出力のテスト
- 実際のSQLクエリ実行時のログ出力
- ログファイルの作成確認
- ログフォーマットの確認

### 6.3 E2Eテスト

#### 6.3.1 環境別動作確認
- 開発環境でのSQLログ出力確認
- ステージング環境でのSQLログ出力確認
- 本番環境でのSQLログ非出力確認

#### 6.3.2 シャーディング対応確認
- すべてのシャード（shard1, shard2, shard3, shard4）でのSQLログ出力確認
- シャードIDが正しく記録されることを確認

## 7. パフォーマンス考慮

### 7.1 ログ出力のオーバーヘッド
- SQLログ出力は同期的に行う（シンプルさを優先）
- ログ出力がSQLクエリ実行のボトルネックにならないよう、軽量な実装を心がける
- 本番環境ではSQLログを出力しないため、パフォーマンスへの影響はない

### 7.2 ログファイルのI/O
- `lumberjack`ライブラリがバッファリングを提供
- 日付別ファイル分割により、1ファイルあたりのサイズを制限

## 8. セキュリティ考慮

### 8.1 機密情報のフィルタリング
- DSN文字列からpasswordを確実にフィルタリング
- 正規表現による確実なマスク処理
- フィルタリング失敗時のフォールバック処理

### 8.2 ログファイルのアクセス制御
- ログファイルのパーミッション設定（0755）
- ログファイルへのアクセス制御（必要に応じて）

## 9. 将来の拡張性

### 9.1 JSON形式への拡張
- 現在はテキスト形式だが、将来的にJSON形式への拡張も可能
- `SQLTextFormatter`を`SQLJSONFormatter`に置き換えることで実現可能

### 9.2 シャード別ログファイル分割
- 現在はすべてのシャードのログを同じファイルに出力
- 将来的にシャード別にログファイルを分割することも可能

### 9.3 ログローテーション機能
- 現在は日付別ファイル分割のみ
- 将来的にサイズベースのローテーションや保持期間設定も追加可能

## 10. 参考情報

### 10.1 GORM Loggerドキュメント
- **公式ドキュメント**: https://gorm.io/docs/logger.html
- **インターフェース定義**: `gorm.io/gorm/logger.Interface`
- **実装例**: GORMのドキュメントに実装例が記載されている

### 10.2 既存実装の参考
- `server/internal/logging/access_logger.go`: アクセスログ実装（ログ出力パターンの参考）
- `server/internal/db/connection.go`: GORM接続作成実装
- `server/internal/config/config.go`: 設定読み込み実装

### 10.3 ログライブラリ
- **`logrus`**: https://github.com/sirupsen/logrus
- **`lumberjack`**: https://github.com/natefinch/lumberjack

