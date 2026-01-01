package email

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSESSender_ImplementsEmailSender(t *testing.T) {
	// SESSenderがEmailSenderインターフェースを実装していることを確認
	var _ EmailSender = (*SESSender)(nil)
}

func TestNewSESSender(t *testing.T) {
	tests := []struct {
		name   string
		region string
		from   string
	}{
		{
			name:   "デフォルト設定",
			region: "us-east-1",
			from:   "sender@example.com",
		},
		{
			name:   "東京リージョン",
			region: "ap-northeast-1",
			from:   "noreply@example.co.jp",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sender, err := NewSESSender(tt.region, tt.from)
			// AWS認証情報がない環境ではエラーになる可能性がある
			if err != nil {
				t.Skipf("AWS認証情報が設定されていないためスキップ: %v", err)
			}
			assert.NotNil(t, sender)
			assert.Equal(t, tt.from, sender.from)
		})
	}
}
