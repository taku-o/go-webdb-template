package cli

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// MockSecretServiceInterface はSecretServiceInterfaceのモック
type MockSecretServiceInterface struct {
	GenerateSecretKeyFunc func(ctx context.Context) (string, error)
}

func (m *MockSecretServiceInterface) GenerateSecretKey(ctx context.Context) (string, error) {
	if m.GenerateSecretKeyFunc != nil {
		return m.GenerateSecretKeyFunc(ctx)
	}
	return "", nil
}

func TestGenerateSecretUsecase_GenerateSecret(t *testing.T) {
	tests := []struct {
		name        string
		mockFunc    func(ctx context.Context) (string, error)
		wantSecret  string
		wantError   bool
		expectedErr string
	}{
		{
			name: "success",
			mockFunc: func(ctx context.Context) (string, error) {
				return "test-secret-key", nil
			},
			wantSecret: "test-secret-key",
			wantError:  false,
		},
		{
			name: "service error",
			mockFunc: func(ctx context.Context) (string, error) {
				return "", errors.New("failed to generate secret key")
			},
			wantSecret:  "",
			wantError:   true,
			expectedErr: "failed to generate secret key",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockSecretServiceInterface{
				GenerateSecretKeyFunc: tt.mockFunc,
			}

			usecase := NewGenerateSecretUsecase(mockService)

			ctx := context.Background()
			gotSecret, err := usecase.GenerateSecret(ctx)

			if tt.wantError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				assert.Empty(t, gotSecret)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantSecret, gotSecret)
			}
		})
	}
}
