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

// MailLogger はメール送信ログを出力するロガー
type MailLogger struct {
	logger  *logrus.Logger
	writer  io.WriteCloser
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

// MailLogFormatter はメール送信ログのJSON形式フォーマッター
type MailLogFormatter struct{}

// Format はログエントリをJSON形式でフォーマット
func (f *MailLogFormatter) Format(timestamp time.Time, fields map[string]interface{}) ([]byte, error) {
	// エラーメッセージの取得
	errorMsg := ""
	if fields["error"] != nil {
		errorMsg = fields["error"].(string)
	}

	entry := MailLogEntry{
		Timestamp:     timestamp.Format("2006-01-02 15:04:05"),
		To:            fields["to"].([]string),
		Subject:       fields["subject"].(string),
		Body:          fields["body"].(string),
		BodyTruncated: fields["body_truncated"].(bool),
		SenderType:    fields["sender_type"].(string),
		Success:       fields["success"].(bool),
		Error:         errorMsg,
	}

	jsonBytes, err := json.Marshal(entry)
	if err != nil {
		return nil, err
	}

	// 改行を追加
	return append(jsonBytes, '\n'), nil
}

// mailLogFormatterAdapter はMailLogFormatterをlogrus.Formatterに適合させるアダプター
type mailLogFormatterAdapter struct {
	formatter *MailLogFormatter
}

// Format はlogrus.Formatterインターフェースを実装
func (a *mailLogFormatterAdapter) Format(entry *logrus.Entry) ([]byte, error) {
	return a.formatter.Format(entry.Time, map[string]interface{}{
		"to":             entry.Data["to"],
		"subject":        entry.Data["subject"],
		"body":           entry.Data["body"],
		"body_truncated": entry.Data["body_truncated"],
		"sender_type":    entry.Data["sender_type"],
		"success":        entry.Data["success"],
		"error":          entry.Data["error"],
	})
}

// NewMailLogger は新しいMailLoggerを作成
// outputDirは絶対パスと相対パスの両方をサポート
// enabledがfalseの場合はnilを返す
func NewMailLogger(outputDir string, enabled bool) (*MailLogger, error) {
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
		MaxBackups: 0,     // バックアップ保持数なし
		MaxAge:     0,     // 保持日数なし
		LocalTime:  true,  // ローカルタイムゾーンを使用
		Compress:   false, // 圧縮なし
	}

	// logrusの設定
	logger := logrus.New()
	logger.SetOutput(writer)
	logger.SetFormatter(&mailLogFormatterAdapter{formatter: &MailLogFormatter{}})
	logger.SetLevel(logrus.InfoLevel)

	return &MailLogger{
		logger:  logger,
		writer:  writer,
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

	// logrus経由でログ出力
	m.logger.WithFields(logrus.Fields{
		"to":             to,
		"subject":        subject,
		"body":           body,
		"body_truncated": bodyTruncated,
		"sender_type":    senderType,
		"success":        success,
		"error":          errorMsg,
	}).Info("mail")
}

// Close はロガーをクローズ
func (m *MailLogger) Close() error {
	if m == nil || m.writer == nil {
		return nil
	}
	return m.writer.Close()
}
