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

	"github.com/sirupsen/logrus"
	"gorm.io/gorm/logger"
	"gopkg.in/natefinch/lumberjack.v2"
)

// SQLLogger はGORMのLoggerインターフェースを実装
type SQLLogger struct {
	logrusLogger *logrus.Logger
	writer       io.WriteCloser
	shardID      int
	driver       string
	logLevel     logger.LogLevel
	outputDir    string
}

// NewSQLLogger は新しいSQLLoggerを作成
// SQLログが無効な場合はnilを返す
func NewSQLLogger(shardID int, driver string, outputDir string, enabled bool) (*SQLLogger, error) {
	// SQLログが無効な場合はnilを返す
	if !enabled {
		return nil, nil
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
	logrusLogger := logrus.New()
	logrusLogger.SetOutput(writer)
	logrusLogger.SetFormatter(&SQLTextFormatter{})
	logrusLogger.SetLevel(logrus.InfoLevel)

	return &SQLLogger{
		logrusLogger: logrusLogger,
		writer:       writer,
		shardID:      shardID,
		driver:       driver,
		logLevel:     logger.Info,
		outputDir:    outputDir,
	}, nil
}

// LogMode はログレベルを設定
func (l *SQLLogger) LogMode(level logger.LogLevel) logger.Interface {
	if l == nil {
		return l
	}
	newLogger := *l
	newLogger.logLevel = level
	return &newLogger
}

// Info は情報ログを出力
func (l *SQLLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l == nil || l.logLevel < logger.Info {
		return
	}
	l.logrusLogger.WithFields(logrus.Fields{
		"message": fmt.Sprintf(msg, data...),
	}).Info("sql")
}

// Warn は警告ログを出力
func (l *SQLLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l == nil || l.logLevel < logger.Warn {
		return
	}
	l.logrusLogger.WithFields(logrus.Fields{
		"message": fmt.Sprintf(msg, data...),
	}).Warn("sql")
}

// Error はエラーログを出力
func (l *SQLLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l == nil || l.logLevel < logger.Error {
		return
	}
	l.logrusLogger.WithFields(logrus.Fields{
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
	l.logrusLogger.WithFields(logrus.Fields{
		"shard_id":      l.shardID,
		"driver":        l.driver,
		"table":         tableName,
		"rows_affected": rows,
		"sql":           sql,
		"duration_ms":   float64(duration.Microseconds()) / 1000.0,
	}).Info("sql")
}

// Close はロガーをクローズ
func (l *SQLLogger) Close() error {
	if l == nil || l.writer == nil {
		return nil
	}
	return l.writer.Close()
}

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
		"[%s] [%v] [%v] [%v] %v | %v | %.2fms\n",
		timestamp, shardID, driver, table, rowsAffected, sql, durationMs,
	)

	return []byte(logLine), nil
}

// extractTableName はSQLクエリからテーブル名を抽出
func extractTableName(sql string) string {
	sql = strings.TrimSpace(sql)
	sqlUpper := strings.ToUpper(sql)

	// FROM句から抽出（SELECT, DELETE）
	if strings.HasPrefix(sqlUpper, "SELECT") || strings.HasPrefix(sqlUpper, "DELETE") {
		re := regexp.MustCompile(`(?i)FROM\s+["'\x60]?(\w+)["'\x60]?`)
		matches := re.FindStringSubmatch(sql)
		if len(matches) > 1 {
			return matches[1]
		}
	}

	// INSERT INTO句から抽出
	if strings.HasPrefix(sqlUpper, "INSERT") {
		re := regexp.MustCompile(`(?i)INTO\s+["'\x60]?(\w+)["'\x60]?`)
		matches := re.FindStringSubmatch(sql)
		if len(matches) > 1 {
			return matches[1]
		}
	}

	// UPDATE句から抽出
	if strings.HasPrefix(sqlUpper, "UPDATE") {
		re := regexp.MustCompile(`(?i)UPDATE\s+["'\x60]?(\w+)["'\x60]?`)
		matches := re.FindStringSubmatch(sql)
		if len(matches) > 1 {
			return matches[1]
		}
	}

	return "unknown"
}

// FilterDSN はDSN文字列からpasswordなどの機密情報をフィルタリング
func FilterDSN(dsn string, driver string) string {
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
		// SQLite DSN: 通常password情報は含まれないが、そのまま返す
		return dsn
	default:
		return dsn
	}
}
