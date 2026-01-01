package logging

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestMailLogEntry_Structure はMailLogEntry構造体の構造をテスト
func TestMailLogEntry_Structure(t *testing.T) {
	entry := MailLogEntry{
		Timestamp:     "2025-01-27 14:30:45",
		To:            []string{"recipient@example.com"},
		Subject:       "Test Subject",
		Body:          "Test Body",
		BodyTruncated: false,
		SenderType:    "mock",
		Success:       true,
		Error:         "",
	}

	if entry.Timestamp != "2025-01-27 14:30:45" {
		t.Errorf("expected Timestamp '2025-01-27 14:30:45', got %s", entry.Timestamp)
	}
	if len(entry.To) != 1 || entry.To[0] != "recipient@example.com" {
		t.Errorf("expected To ['recipient@example.com'], got %v", entry.To)
	}
	if entry.Subject != "Test Subject" {
		t.Errorf("expected Subject 'Test Subject', got %s", entry.Subject)
	}
	if entry.Body != "Test Body" {
		t.Errorf("expected Body 'Test Body', got %s", entry.Body)
	}
	if entry.BodyTruncated != false {
		t.Errorf("expected BodyTruncated false, got %v", entry.BodyTruncated)
	}
	if entry.SenderType != "mock" {
		t.Errorf("expected SenderType 'mock', got %s", entry.SenderType)
	}
	if entry.Success != true {
		t.Errorf("expected Success true, got %v", entry.Success)
	}
	if entry.Error != "" {
		t.Errorf("expected Error '', got %s", entry.Error)
	}
}

// TestMailLogEntry_JSONSerialization はMailLogEntryのJSON変換をテスト
func TestMailLogEntry_JSONSerialization(t *testing.T) {
	entry := MailLogEntry{
		Timestamp:     "2025-01-27 14:30:45",
		To:            []string{"recipient@example.com"},
		Subject:       "Test Subject",
		Body:          "Test Body",
		BodyTruncated: false,
		SenderType:    "ses",
		Success:       true,
		Error:         "",
	}

	jsonBytes, err := json.Marshal(entry)
	if err != nil {
		t.Fatalf("failed to marshal MailLogEntry: %v", err)
	}

	jsonStr := string(jsonBytes)

	// JSONフィールドの確認
	if !strings.Contains(jsonStr, `"timestamp":"2025-01-27 14:30:45"`) {
		t.Errorf("expected timestamp field in JSON, got %s", jsonStr)
	}
	if !strings.Contains(jsonStr, `"to":["recipient@example.com"]`) {
		t.Errorf("expected to field in JSON, got %s", jsonStr)
	}
	if !strings.Contains(jsonStr, `"subject":"Test Subject"`) {
		t.Errorf("expected subject field in JSON, got %s", jsonStr)
	}
	if !strings.Contains(jsonStr, `"sender_type":"ses"`) {
		t.Errorf("expected sender_type field in JSON, got %s", jsonStr)
	}
	if !strings.Contains(jsonStr, `"success":true`) {
		t.Errorf("expected success field in JSON, got %s", jsonStr)
	}
	// エラーが空の場合はomitemptyでフィールドが省略される
	if strings.Contains(jsonStr, `"error"`) {
		t.Errorf("expected error field to be omitted when empty, got %s", jsonStr)
	}
}

// TestMailLogEntry_JSONWithError はエラーがある場合のJSON変換をテスト
func TestMailLogEntry_JSONWithError(t *testing.T) {
	entry := MailLogEntry{
		Timestamp:     "2025-01-27 14:30:45",
		To:            []string{"recipient@example.com"},
		Subject:       "Test Subject",
		Body:          "Test Body",
		BodyTruncated: false,
		SenderType:    "ses",
		Success:       false,
		Error:         "SES error: Invalid email address",
	}

	jsonBytes, err := json.Marshal(entry)
	if err != nil {
		t.Fatalf("failed to marshal MailLogEntry: %v", err)
	}

	jsonStr := string(jsonBytes)

	if !strings.Contains(jsonStr, `"error":"SES error: Invalid email address"`) {
		t.Errorf("expected error field in JSON, got %s", jsonStr)
	}
}

// TestNewMailLogger_Enabled は有効化時のMailLogger作成をテスト
func TestNewMailLogger_Enabled(t *testing.T) {
	tmpDir := t.TempDir()

	logger, err := NewMailLogger(tmpDir, true)
	if err != nil {
		t.Fatalf("failed to create MailLogger: %v", err)
	}
	defer logger.Close()

	if logger == nil {
		t.Error("expected non-nil logger when enabled")
	}
}

// TestNewMailLogger_Disabled は無効化時のMailLogger作成をテスト（nilを返す）
func TestNewMailLogger_Disabled(t *testing.T) {
	tmpDir := t.TempDir()

	logger, err := NewMailLogger(tmpDir, false)
	if err != nil {
		t.Fatalf("unexpected error when disabled: %v", err)
	}

	if logger != nil {
		t.Error("expected nil logger when disabled")
	}
}

// TestNewMailLogger_CreatesDirectory はディレクトリが作成されることをテスト
func TestNewMailLogger_CreatesDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	logDir := filepath.Join(tmpDir, "newdir", "logs")

	logger, err := NewMailLogger(logDir, true)
	if err != nil {
		t.Fatalf("failed to create MailLogger: %v", err)
	}
	defer logger.Close()

	// ディレクトリが作成されたことを確認
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		t.Errorf("expected log directory to be created: %s", logDir)
	}
}

// TestLogMail_Success は正常なログ出力をテスト
func TestLogMail_Success(t *testing.T) {
	tmpDir := t.TempDir()

	logger, err := NewMailLogger(tmpDir, true)
	if err != nil {
		t.Fatalf("failed to create MailLogger: %v", err)
	}
	defer logger.Close()

	// ログを出力
	logger.LogMail(
		[]string{"recipient@example.com"},
		"Test Subject",
		"Test Body",
		"mock",
		true,
		nil,
	)

	// ログファイルの内容を確認
	today := time.Now().Format("2006-01-02")
	logFile := filepath.Join(tmpDir, "mail-"+today+".log")

	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("failed to read log file: %v", err)
	}

	// JSON形式で正しく出力されていることを確認
	var entry MailLogEntry
	if err := json.Unmarshal(content[:len(content)-1], &entry); err != nil { // 改行を除去
		t.Fatalf("failed to unmarshal log entry: %v", err)
	}

	if entry.To[0] != "recipient@example.com" {
		t.Errorf("expected To 'recipient@example.com', got %v", entry.To)
	}
	if entry.Subject != "Test Subject" {
		t.Errorf("expected Subject 'Test Subject', got %s", entry.Subject)
	}
	if entry.Success != true {
		t.Errorf("expected Success true, got %v", entry.Success)
	}
}

// TestLogMail_BodyTruncation は200文字以上のメール本文が切り捨てられることをテスト
func TestLogMail_BodyTruncation(t *testing.T) {
	tmpDir := t.TempDir()

	logger, err := NewMailLogger(tmpDir, true)
	if err != nil {
		t.Fatalf("failed to create MailLogger: %v", err)
	}
	defer logger.Close()

	// 201文字のメール本文
	longBody := strings.Repeat("a", 201)

	logger.LogMail(
		[]string{"recipient@example.com"},
		"Test Subject",
		longBody,
		"mock",
		true,
		nil,
	)

	// ログファイルの内容を確認
	today := time.Now().Format("2006-01-02")
	logFile := filepath.Join(tmpDir, "mail-"+today+".log")

	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("failed to read log file: %v", err)
	}

	var entry MailLogEntry
	if err := json.Unmarshal(content[:len(content)-1], &entry); err != nil {
		t.Fatalf("failed to unmarshal log entry: %v", err)
	}

	// 200文字で切り捨てられていることを確認
	expectedBody := strings.Repeat("a", 200) + "..."
	if entry.Body != expectedBody {
		t.Errorf("expected Body to be truncated to 200 chars + '...', got length %d", len(entry.Body))
	}
	if entry.BodyTruncated != true {
		t.Errorf("expected BodyTruncated true, got %v", entry.BodyTruncated)
	}
}

// TestLogMail_BodyExactly200Chars は200文字ちょうどのメール本文が切り捨てられないことをテスト
func TestLogMail_BodyExactly200Chars(t *testing.T) {
	tmpDir := t.TempDir()

	logger, err := NewMailLogger(tmpDir, true)
	if err != nil {
		t.Fatalf("failed to create MailLogger: %v", err)
	}
	defer logger.Close()

	// 200文字のメール本文
	body := strings.Repeat("a", 200)

	logger.LogMail(
		[]string{"recipient@example.com"},
		"Test Subject",
		body,
		"mock",
		true,
		nil,
	)

	// ログファイルの内容を確認
	today := time.Now().Format("2006-01-02")
	logFile := filepath.Join(tmpDir, "mail-"+today+".log")

	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("failed to read log file: %v", err)
	}

	var entry MailLogEntry
	if err := json.Unmarshal(content[:len(content)-1], &entry); err != nil {
		t.Fatalf("failed to unmarshal log entry: %v", err)
	}

	// 200文字ちょうどは切り捨てられない
	if entry.Body != body {
		t.Errorf("expected Body to remain unchanged, got length %d", len(entry.Body))
	}
	if entry.BodyTruncated != false {
		t.Errorf("expected BodyTruncated false, got %v", entry.BodyTruncated)
	}
}

// TestLogMail_WithError はエラーがある場合のログ出力をテスト
func TestLogMail_WithError(t *testing.T) {
	tmpDir := t.TempDir()

	logger, err := NewMailLogger(tmpDir, true)
	if err != nil {
		t.Fatalf("failed to create MailLogger: %v", err)
	}
	defer logger.Close()

	sendErr := errors.New("SES error: Invalid email address")

	logger.LogMail(
		[]string{"recipient@example.com"},
		"Test Subject",
		"Test Body",
		"ses",
		false,
		sendErr,
	)

	// ログファイルの内容を確認
	today := time.Now().Format("2006-01-02")
	logFile := filepath.Join(tmpDir, "mail-"+today+".log")

	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("failed to read log file: %v", err)
	}

	var entry MailLogEntry
	if err := json.Unmarshal(content[:len(content)-1], &entry); err != nil {
		t.Fatalf("failed to unmarshal log entry: %v", err)
	}

	if entry.Success != false {
		t.Errorf("expected Success false, got %v", entry.Success)
	}
	if entry.Error != "SES error: Invalid email address" {
		t.Errorf("expected Error 'SES error: Invalid email address', got %s", entry.Error)
	}
}

// TestLogMail_NilLogger はnilロガーに対するログ出力がパニックしないことをテスト
func TestLogMail_NilLogger(t *testing.T) {
	var logger *MailLogger = nil

	// パニックしないことを確認
	logger.LogMail(
		[]string{"recipient@example.com"},
		"Test Subject",
		"Test Body",
		"mock",
		true,
		nil,
	)
}

// TestClose_Success は正常なクローズをテスト
func TestClose_Success(t *testing.T) {
	tmpDir := t.TempDir()

	logger, err := NewMailLogger(tmpDir, true)
	if err != nil {
		t.Fatalf("failed to create MailLogger: %v", err)
	}

	err = logger.Close()
	if err != nil {
		t.Errorf("expected no error on close, got %v", err)
	}
}

// TestClose_NilLogger はnilロガーのクローズをテスト
func TestClose_NilLogger(t *testing.T) {
	var logger *MailLogger = nil

	err := logger.Close()
	if err != nil {
		t.Errorf("expected no error on nil close, got %v", err)
	}
}

// TestLogMail_MultipleTo は複数の送信先アドレスのログ出力をテスト
func TestLogMail_MultipleTo(t *testing.T) {
	tmpDir := t.TempDir()

	logger, err := NewMailLogger(tmpDir, true)
	if err != nil {
		t.Fatalf("failed to create MailLogger: %v", err)
	}
	defer logger.Close()

	to := []string{"user1@example.com", "user2@example.com", "user3@example.com"}

	logger.LogMail(
		to,
		"Test Subject",
		"Test Body",
		"mock",
		true,
		nil,
	)

	// ログファイルの内容を確認
	today := time.Now().Format("2006-01-02")
	logFile := filepath.Join(tmpDir, "mail-"+today+".log")

	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("failed to read log file: %v", err)
	}

	var entry MailLogEntry
	if err := json.Unmarshal(content[:len(content)-1], &entry); err != nil {
		t.Fatalf("failed to unmarshal log entry: %v", err)
	}

	if len(entry.To) != 3 {
		t.Errorf("expected 3 recipients, got %d", len(entry.To))
	}
	if entry.To[0] != "user1@example.com" {
		t.Errorf("expected first recipient 'user1@example.com', got %s", entry.To[0])
	}
}
