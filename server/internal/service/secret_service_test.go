package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSecretService_GenerateSecretKey(t *testing.T) {
	tests := []struct {
		name      string
		wantError bool
	}{
		{
			name:      "success",
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewSecretService()
			ctx := context.Background()

			got, err := service.GenerateSecretKey(ctx)

			if tt.wantError {
				assert.Error(t, err)
				assert.Empty(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, got)
			}
		})
	}
}
