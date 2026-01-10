package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// MockDateService はDateServiceのモック
type MockDateService struct {
	GetTodayFunc func(ctx context.Context) (string, error)
}

func (m *MockDateService) GetToday(ctx context.Context) (string, error) {
	if m.GetTodayFunc != nil {
		return m.GetTodayFunc(ctx)
	}
	return "", nil
}

func TestTodayUsecase_GetToday(t *testing.T) {
	tests := []struct {
		name        string
		mockFunc    func(ctx context.Context) (string, error)
		want        string
		wantErr     bool
		expectedErr string
	}{
		{
			name: "returns today's date from DateService",
			mockFunc: func(ctx context.Context) (string, error) {
				return "2026-01-10", nil
			},
			want:    "2026-01-10",
			wantErr: false,
		},
		{
			name: "returns error when DateService fails",
			mockFunc: func(ctx context.Context) (string, error) {
				return "", errors.New("service error")
			},
			want:        "",
			wantErr:     true,
			expectedErr: "service error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockDateService{
				GetTodayFunc: tt.mockFunc,
			}

			u := NewTodayUsecase(mockService)
			got, err := u.GetToday(context.Background())

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
