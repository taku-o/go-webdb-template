package email

import (
	"context"

	"gopkg.in/mail.v2"
)

// MailpitSender はMailpitにメールを送信する実装
type MailpitSender struct {
	smtpHost string
	smtpPort int
	from     string
}

// NewMailpitSender は新しいMailpitSenderを作成
func NewMailpitSender(smtpHost string, smtpPort int, from string) *MailpitSender {
	return &MailpitSender{
		smtpHost: smtpHost,
		smtpPort: smtpPort,
		from:     from,
	}
}

// Send はgomailを使用してMailpitにメールを送信
func (s *MailpitSender) Send(ctx context.Context, to []string, subject, body string) error {
	m := mail.NewMessage()
	m.SetHeader("From", s.from)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	// Mailpit用のダイアラー（認証なし）
	d := mail.NewDialer(s.smtpHost, s.smtpPort, "", "")

	return d.DialAndSend(m)
}
