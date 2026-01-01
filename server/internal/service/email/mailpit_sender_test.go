package email

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMailpitSender_ImplementsEmailSender(t *testing.T) {
	// MailpitSenderがEmailSenderインターフェースを実装していることを確認
	var _ EmailSender = (*MailpitSender)(nil)
}

func TestNewMailpitSender(t *testing.T) {
	tests := []struct {
		name     string
		smtpHost string
		smtpPort int
		from     string
	}{
		{
			name:     "デフォルト設定",
			smtpHost: "localhost",
			smtpPort: 1025,
			from:     "noreply@example.com",
		},
		{
			name:     "カスタム設定",
			smtpHost: "mailpit.local",
			smtpPort: 2025,
			from:     "custom@example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sender := NewMailpitSender(tt.smtpHost, tt.smtpPort, tt.from)
			assert.NotNil(t, sender)
			assert.Equal(t, tt.smtpHost, sender.smtpHost)
			assert.Equal(t, tt.smtpPort, sender.smtpPort)
			assert.Equal(t, tt.from, sender.from)
		})
	}
}

func TestMailpitSender_Send_ConnectionError(t *testing.T) {
	// Mailpitが起動していない場合のテスト（接続エラー）
	sender := NewMailpitSender("localhost", 19999, "noreply@example.com") // 存在しないポート
	err := sender.Send(context.Background(), []string{"test@example.com"}, "Test Subject", "Test Body")
	assert.Error(t, err, "Mailpitが起動していない場合はエラーになるべき")
}
