package email

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/taku-o/go-webdb-template/internal/config"
	"github.com/taku-o/go-webdb-template/internal/logging"
)

func TestNewEmailService_Mock(t *testing.T) {
	cfg := &config.EmailConfig{
		SenderType: "mock",
	}

	service, err := NewEmailService(cfg, nil)
	assert.NoError(t, err)
	assert.NotNil(t, service)
}

func TestNewEmailService_Mailpit(t *testing.T) {
	cfg := &config.EmailConfig{
		SenderType: "mailpit",
		Mailpit: config.MailpitConfig{
			SMTPHost: "localhost",
			SMTPPort: 1025,
		},
		SES: config.SESConfig{
			From: "noreply@example.com",
		},
	}

	service, err := NewEmailService(cfg, nil)
	assert.NoError(t, err)
	assert.NotNil(t, service)
}

func TestNewEmailService_SES(t *testing.T) {
	cfg := &config.EmailConfig{
		SenderType: "ses",
		SES: config.SESConfig{
			From:   "sender@example.com",
			Region: "us-east-1",
		},
	}

	service, err := NewEmailService(cfg, nil)
	// AWS認証情報がない環境ではエラーになる可能性がある
	if err != nil {
		t.Skipf("AWS認証情報が設定されていないためスキップ: %v", err)
	}
	assert.NotNil(t, service)
}

func TestNewEmailService_InvalidSenderType(t *testing.T) {
	cfg := &config.EmailConfig{
		SenderType: "invalid",
	}

	service, err := NewEmailService(cfg, nil)
	assert.Error(t, err)
	assert.Nil(t, service)
}

func TestNewEmailService_DefaultSenderType_Develop(t *testing.T) {
	// APP_ENV=developの場合、MockSenderがデフォルト
	os.Setenv("APP_ENV", "develop")
	defer os.Unsetenv("APP_ENV")

	cfg := &config.EmailConfig{
		SenderType: "", // 空の場合はデフォルト
	}

	service, err := NewEmailService(cfg, nil)
	assert.NoError(t, err)
	assert.NotNil(t, service)
}

func TestNewEmailService_DefaultSenderType_Production(t *testing.T) {
	// APP_ENV=productionの場合、SESSenderがデフォルト
	os.Setenv("APP_ENV", "production")
	defer os.Unsetenv("APP_ENV")

	cfg := &config.EmailConfig{
		SenderType: "", // 空の場合はデフォルト
		SES: config.SESConfig{
			From:   "sender@example.com",
			Region: "us-east-1",
		},
	}

	service, err := NewEmailService(cfg, nil)
	// AWS認証情報がない環境ではエラーになる可能性がある
	if err != nil {
		t.Skipf("AWS認証情報が設定されていないためスキップ: %v", err)
	}
	assert.NotNil(t, service)
}

func TestEmailService_SendEmail(t *testing.T) {
	cfg := &config.EmailConfig{
		SenderType: "mock",
	}

	service, err := NewEmailService(cfg, nil)
	assert.NoError(t, err)

	err = service.SendEmail(context.Background(), []string{"test@example.com"}, "Test Subject", "Test Body")
	assert.NoError(t, err)
}

func TestNewEmailService_WithMailLogger(t *testing.T) {
	tmpDir := t.TempDir()

	mailLogger, err := logging.NewMailLogger(tmpDir, true)
	assert.NoError(t, err)
	defer mailLogger.Close()

	cfg := &config.EmailConfig{
		SenderType: "mock",
	}

	service, err := NewEmailService(cfg, mailLogger)
	assert.NoError(t, err)
	assert.NotNil(t, service)
}

func TestEmailService_SendEmail_WithLogging(t *testing.T) {
	tmpDir := t.TempDir()

	mailLogger, err := logging.NewMailLogger(tmpDir, true)
	assert.NoError(t, err)
	defer mailLogger.Close()

	cfg := &config.EmailConfig{
		SenderType: "mock",
	}

	service, err := NewEmailService(cfg, mailLogger)
	assert.NoError(t, err)

	err = service.SendEmail(context.Background(), []string{"test@example.com"}, "Test Subject", "Test Body")
	assert.NoError(t, err)

	// ログファイルが作成されていることを確認
	today := time.Now().Format("2006-01-02")
	logFile := filepath.Join(tmpDir, "mail-"+today+".log")

	content, err := os.ReadFile(logFile)
	assert.NoError(t, err)
	assert.NotEmpty(t, content)

	// ログの内容を確認
	var entry logging.MailLogEntry
	err = json.Unmarshal(content[:len(content)-1], &entry)
	assert.NoError(t, err)
	assert.Equal(t, []string{"test@example.com"}, entry.To)
	assert.Equal(t, "Test Subject", entry.Subject)
	assert.Equal(t, "mock", entry.SenderType)
	assert.True(t, entry.Success)
}

func TestEmailService_GetSenderType(t *testing.T) {
	tests := []struct {
		name       string
		senderType string
		expected   string
	}{
		{"MockSender", "mock", "mock"},
		{"MailpitSender", "mailpit", "mailpit"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.EmailConfig{
				SenderType: tt.senderType,
				Mailpit: config.MailpitConfig{
					SMTPHost: "localhost",
					SMTPPort: 1025,
				},
				SES: config.SESConfig{
					From: "noreply@example.com",
				},
			}

			service, err := NewEmailService(cfg, nil)
			assert.NoError(t, err)

			senderType := service.GetSenderType()
			assert.Equal(t, tt.expected, senderType)
		})
	}
}
