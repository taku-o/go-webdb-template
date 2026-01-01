package handler

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/taku-o/go-webdb-template/internal/config"
	"github.com/taku-o/go-webdb-template/internal/service/email"
)

func TestNewEmailHandler(t *testing.T) {
	// MockSenderを使用したEmailServiceを作成
	cfg := &config.EmailConfig{
		SenderType: "mock",
	}
	emailService, err := email.NewEmailService(cfg)
	assert.NoError(t, err)

	templateService := email.NewTemplateService()

	handler := NewEmailHandler(emailService, templateService)
	assert.NotNil(t, handler)
}

func TestEmailHandler_ValidateEmail(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{
			name:    "正常なメールアドレス",
			email:   "test@example.com",
			wantErr: false,
		},
		{
			name:    "日本語名前付きメールアドレス",
			email:   "山田太郎 <yamada@example.com>",
			wantErr: false,
		},
		{
			name:    "不正なメールアドレス",
			email:   "invalid-email",
			wantErr: true,
		},
		{
			name:    "空のメールアドレス",
			email:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateEmail(tt.email)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
