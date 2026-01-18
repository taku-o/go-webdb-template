package jobqueue

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

// MockDelayPrintService はテスト用のモックサービス
type MockDelayPrintService struct {
	PrintedMessage string
	PrintError     error
}

func (m *MockDelayPrintService) PrintMessage(message string) error {
	m.PrintedMessage = message
	return m.PrintError
}

func TestDelayPrintUsecase_Execute(t *testing.T) {
	tests := []struct {
		name           string
		message        string
		printError     error
		wantMessage    string
		wantErr        bool
	}{
		{
			name:        "正常なメッセージを処理",
			message:     "Hello, World!",
			printError:  nil,
			wantMessage: "Hello, World!",
			wantErr:     false,
		},
		{
			name:        "空のメッセージはデフォルト値を使用",
			message:     "",
			printError:  nil,
			wantMessage: "Job executed successfully",
			wantErr:     false,
		},
		{
			name:        "日本語メッセージを処理",
			message:     "ジョブが実行されました",
			printError:  nil,
			wantMessage: "ジョブが実行されました",
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockDelayPrintService{
				PrintError: tt.printError,
			}

			usecase := NewDelayPrintUsecase(mockService)
			err := usecase.Execute(context.Background(), tt.message)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantMessage, mockService.PrintedMessage)
			}
		})
	}
}

func TestNewDelayPrintUsecase(t *testing.T) {
	mockService := &MockDelayPrintService{}
	usecase := NewDelayPrintUsecase(mockService)
	assert.NotNil(t, usecase)
}
