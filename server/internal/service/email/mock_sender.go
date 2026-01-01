package email

import (
	"context"
	"fmt"
)

// MockSender は標準出力にメールを出力する送信実装
type MockSender struct{}

// NewMockSender は新しいMockSenderを作成
func NewMockSender() *MockSender {
	return &MockSender{}
}

// Send はメール内容を標準出力に出力
func (s *MockSender) Send(ctx context.Context, to []string, subject, body string) error {
	fmt.Printf("[Mock Email] To: %v | Subject: %s\nBody: %s\n", to, subject, body)
	return nil
}
