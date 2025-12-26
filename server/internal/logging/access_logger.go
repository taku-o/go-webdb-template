package logging

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// AccessLogger はアクセスログを出力するロガー
type AccessLogger struct {
	logger *logrus.Logger
	writer io.WriteCloser
}

// AccessLogEntry はアクセスログのJSON構造体
type AccessLogEntry struct {
	Timestamp      string  `json:"timestamp"`
	Method         string  `json:"method"`
	Path           string  `json:"path"`
	Protocol       string  `json:"protocol"`
	StatusCode     int     `json:"status_code"`
	ResponseTimeMs float64 `json:"response_time_ms"`
	RemoteIP       string  `json:"remote_ip"`
	UserAgent      string  `json:"user_agent"`
	Headers        string  `json:"headers,omitempty"`
	RequestBody    string  `json:"request_body,omitempty"`
}

// CustomTextFormatter はカスタムテキストフォーマッター（JSON形式）
type CustomTextFormatter struct{}

// Format はログエントリをJSON形式でフォーマット
func (f *CustomTextFormatter) Format(timestamp time.Time, fields map[string]interface{}) ([]byte, error) {
	entry := AccessLogEntry{
		Timestamp:      timestamp.Format("2006-01-02 15:04:05"),
		Method:         fields["method"].(string),
		Path:           fields["path"].(string),
		Protocol:       fields["protocol"].(string),
		StatusCode:     fields["status_code"].(int),
		ResponseTimeMs: fields["response_time_ms"].(float64),
		RemoteIP:       fields["remote_ip"].(string),
		UserAgent:      fields["user_agent"].(string),
		Headers:        fields["headers"].(string),
		RequestBody:    fields["request_body"].(string),
	}

	jsonBytes, err := json.Marshal(entry)
	if err != nil {
		return nil, err
	}

	// 改行を追加
	return append(jsonBytes, '\n'), nil
}

// logrusFormatterAdapter はCustomTextFormatterをlogrus.Formatterに適合させるアダプター
type logrusFormatterAdapter struct {
	formatter *CustomTextFormatter
}

// Format はlogrus.Formatterインターフェースを実装
func (a *logrusFormatterAdapter) Format(entry *logrus.Entry) ([]byte, error) {
	return a.formatter.Format(entry.Time, map[string]interface{}{
		"method":           entry.Data["method"],
		"path":             entry.Data["path"],
		"protocol":         entry.Data["protocol"],
		"status_code":      entry.Data["status_code"],
		"response_time_ms": entry.Data["response_time_ms"],
		"remote_ip":        entry.Data["remote_ip"],
		"user_agent":       entry.Data["user_agent"],
		"headers":          entry.Data["headers"],
		"request_body":     entry.Data["request_body"],
	})
}

// NewAccessLogger は新しいAccessLoggerを作成
// outputDirは絶対パスと相対パスの両方をサポートします
// - 相対パス: サーバーの実行ディレクトリからの相対パス（例: "logs"）
// - 絶対パス: システムの絶対パス（例: "/var/log/go-webdb-template"）
func NewAccessLogger(logType string, outputDir string) (*AccessLogger, error) {
	// 出力ディレクトリの作成
	// os.MkdirAllは絶対パスと相対パスの両方をサポート
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	// lumberjackの設定（日付別ファイル分割）
	// filepath.Joinは絶対パスと相対パスの両方を正しく処理
	logFileName := fmt.Sprintf("%s-access-%s.log", logType, time.Now().Format("2006-01-02"))
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
	logger.SetFormatter(&logrusFormatterAdapter{formatter: &CustomTextFormatter{}})
	logger.SetLevel(logrus.InfoLevel)

	return &AccessLogger{
		logger: logger,
		writer: writer,
	}, nil
}

// LogAccess はアクセスログを出力
func (a *AccessLogger) LogAccess(method, path, protocol string, statusCode int, responseTimeMs float64, remoteIP, userAgent, headers, requestBody string) {
	a.logger.WithFields(logrus.Fields{
		"method":           method,
		"path":             path,
		"protocol":         protocol,
		"status_code":      statusCode,
		"response_time_ms": responseTimeMs,
		"remote_ip":        remoteIP,
		"user_agent":       userAgent,
		"headers":          headers,
		"request_body":     requestBody,
	}).Info("access")
}

// Close はロガーをクローズ
func (a *AccessLogger) Close() error {
	if a.writer != nil {
		return a.writer.Close()
	}
	return nil
}

// Writer はログの出力先を返す
func (a *AccessLogger) Writer() io.Writer {
	return a.writer
}
