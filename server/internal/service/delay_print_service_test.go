package service

import (
	"bytes"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDelayPrintService_PrintMessage(t *testing.T) {
	tests := []struct {
		name    string
		message string
		wantErr bool
	}{
		{
			name:    "正常なメッセージを出力",
			message: "Hello, World!",
			wantErr: false,
		},
		{
			name:    "空のメッセージを出力",
			message: "",
			wantErr: false,
		},
		{
			name:    "日本語メッセージを出力",
			message: "ジョブが正常に実行されました",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 標準出力をキャプチャ
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			service := NewDelayPrintService()
			err := service.PrintMessage(tt.message)

			// 標準出力を元に戻す
			w.Close()
			os.Stdout = oldStdout

			// 出力を読み取る
			var buf bytes.Buffer
			buf.ReadFrom(r)
			output := buf.String()

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				// メッセージが含まれていることを確認
				assert.Contains(t, output, tt.message)
				// タイムスタンプ形式を確認（YYYY-MM-DD HH:MM:SS）
				assert.True(t, strings.Contains(output, time.Now().Format("2006-01-02")))
			}
		})
	}
}

func TestNewDelayPrintService(t *testing.T) {
	service := NewDelayPrintService()
	assert.NotNil(t, service)
}
