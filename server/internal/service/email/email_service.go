package email

import (
	"context"
	"fmt"
	"os"

	"github.com/taku-o/go-webdb-template/internal/config"
	"github.com/taku-o/go-webdb-template/internal/logging"
)

// EmailService はメール送信サービス
type EmailService struct {
	sender EmailSender
	logger *logging.MailLogger
}

// NewEmailService は新しいEmailServiceを作成
// 設定ファイルのsender_typeに基づいて適切な送信実装を選択します
// sender_typeが空の場合は環境変数APP_ENVに基づいてデフォルトを選択します
// mailLoggerがnilの場合はログ出力を行いません
func NewEmailService(cfg *config.EmailConfig, mailLogger *logging.MailLogger) (*EmailService, error) {
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
		logger: mailLogger,
	}, nil
}

// SendEmail はメールを送信します
func (s *EmailService) SendEmail(ctx context.Context, to []string, subject, body string) error {
	// メール送信
	err := s.sender.Send(ctx, to, subject, body)

	// 送信後のログ出力
	if s.logger != nil {
		success := err == nil
		s.logger.LogMail(to, subject, body, s.GetSenderType(), success, err)
	}

	return err
}

// GetSenderType は送信実装の種類を取得
func (s *EmailService) GetSenderType() string {
	switch s.sender.(type) {
	case *MockSender:
		return "mock"
	case *MailpitSender:
		return "mailpit"
	case *SESSender:
		return "ses"
	default:
		return "unknown"
	}
}
