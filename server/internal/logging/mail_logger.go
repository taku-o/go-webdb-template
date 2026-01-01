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

// MailLogFormatter はメール送信ログのフォーマッター
// logrusのFormatterインターフェースに適合させるため定義
// 実際のJSON出力はLogMailメソッド内で直接行う
type MailLogFormatter struct{}

// Format はlogrus.Formatterインターフェースを実装
func (f *MailLogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	return []byte{}, nil
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
	logger.SetFormatter(&MailLogFormatter{})
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
	jsonBytes, jsonErr := json.Marshal(entry)
	if jsonErr != nil {
		// JSONエンコードに失敗した場合は標準エラー出力に記録
		fmt.Fprintf(os.Stderr, "Failed to encode mail log: %v\n", jsonErr)
		return
	}

	// ログファイルに出力（改行を追加）
	if _, writeErr := m.writer.Write(append(jsonBytes, '\n')); writeErr != nil {
		// ログ出力に失敗した場合は標準エラー出力に記録
		fmt.Fprintf(os.Stderr, "Failed to write mail log: %v\n", writeErr)
	}
}

// Close はロガーをクローズ
func (m *MailLogger) Close() error {
	if m == nil || m.writer == nil {
		return nil
	}
	return m.writer.Close()
}
