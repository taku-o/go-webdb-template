// Package email はメール送信機能を提供するパッケージです。
// 標準出力、Mailpit、AWS SESの3つの送信方式をサポートします。
package email

// EmailSender はメール送信のインターフェース
type EmailSender interface {
	// Send はメールを送信します
	// to: 送信先メールアドレスのリスト
	// subject: メールの件名
	// body: メールの本文
	// 戻り値: 送信に失敗した場合はエラー
	Send(to []string, subject, body string) error
}
