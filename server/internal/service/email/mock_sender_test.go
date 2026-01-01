package email

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMockSender_Send(t *testing.T) {
	tests := []struct {
		name    string
		to      []string
		subject string
		body    string
		wantErr bool
	}{
		{
			name:    "正常なメール送信",
			to:      []string{"test@example.com"},
			subject: "Test Subject",
			body:    "Test Body",
			wantErr: false,
		},
		{
			name:    "複数の送信先",
			to:      []string{"test1@example.com", "test2@example.com"},
			subject: "Multiple Recipients",
			body:    "Test Body for multiple recipients",
			wantErr: false,
		},
		{
			name:    "空の送信先",
			to:      []string{},
			subject: "Empty Recipients",
			body:    "Test Body",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 標準出力をキャプチャ
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			sender := NewMockSender()
			err := sender.Send(tt.to, tt.subject, tt.body)

			// 標準出力を復元してキャプチャした内容を取得
			w.Close()
			os.Stdout = oldStdout
			var buf bytes.Buffer
			io.Copy(&buf, r)
			output := buf.String()

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				// 出力フォーマットの確認
				assert.True(t, strings.Contains(output, "[Mock Email]"), "出力に [Mock Email] が含まれるべき")
				assert.True(t, strings.Contains(output, tt.subject), "出力に件名が含まれるべき")
				assert.True(t, strings.Contains(output, tt.body), "出力に本文が含まれるべき")
			}
		})
	}
}

func TestMockSender_ImplementsEmailSender(t *testing.T) {
	// MockSenderがEmailSenderインターフェースを実装していることを確認
	var _ EmailSender = (*MockSender)(nil)
}

func TestNewMockSender(t *testing.T) {
	sender := NewMockSender()
	assert.NotNil(t, sender)
}
