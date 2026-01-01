package email

import (
	"fmt"
	"os"

	"github.com/taku-o/go-webdb-template/internal/config"
)

// EmailService はメール送信サービス
type EmailService struct {
	sender EmailSender
}

// NewEmailService は新しいEmailServiceを作成
// 設定ファイルのsender_typeに基づいて適切な送信実装を選択します
// sender_typeが空の場合は環境変数APP_ENVに基づいてデフォルトを選択します
func NewEmailService(cfg *config.EmailConfig) (*EmailService, error) {
	senderType := cfg.SenderType

	// sender_typeが空の場合は環境に基づいてデフォルトを選択
	if senderType == "" {
		appEnv := os.Getenv("APP_ENV")
		switch appEnv {
		case "staging", "production":
			senderType = "ses"
		default:
			// develop環境やその他の場合はmockをデフォルト
			senderType = "mock"
		}
	}

	var sender EmailSender
	var err error

	switch senderType {
	case "mock":
		sender = NewMockSender()
	case "mailpit":
		from := cfg.SES.From
		if from == "" {
			from = "noreply@example.com"
		}
		sender = NewMailpitSender(cfg.Mailpit.SMTPHost, cfg.Mailpit.SMTPPort, from)
	case "ses":
		sender, err = NewSESSender(cfg.SES.Region, cfg.SES.From)
		if err != nil {
			return nil, fmt.Errorf("failed to create SES sender: %w", err)
		}
	default:
		return nil, fmt.Errorf("invalid sender type: %s", senderType)
	}

	return &EmailService{
		sender: sender,
	}, nil
}

// SendEmail はメールを送信します
func (s *EmailService) SendEmail(to []string, subject, body string) error {
	return s.sender.Send(to, subject, body)
}
