package email

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/taku-o/go-webdb-template/internal/config"
)

func TestNewEmailService_Mock(t *testing.T) {
	cfg := &config.EmailConfig{
		SenderType: "mock",
	}

	service, err := NewEmailService(cfg)
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

	service, err := NewEmailService(cfg)
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

	service, err := NewEmailService(cfg)
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

	service, err := NewEmailService(cfg)
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

	service, err := NewEmailService(cfg)
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

	service, err := NewEmailService(cfg)
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

	service, err := NewEmailService(cfg)
	assert.NoError(t, err)

	err = service.SendEmail(context.Background(), []string{"test@example.com"}, "Test Subject", "Test Body")
	assert.NoError(t, err)
}
