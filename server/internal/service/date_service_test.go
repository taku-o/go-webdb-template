package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDateService_GetToday(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "returns today's date in YYYY-MM-DD format",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewDateService()
			got, err := s.GetToday(context.Background())

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			// Check date format is YYYY-MM-DD
			_, parseErr := time.Parse("2006-01-02", got)
			assert.NoError(t, parseErr, "date should be in YYYY-MM-DD format")

			// Check the date matches today
			expected := time.Now().Format("2006-01-02")
			assert.Equal(t, expected, got)
		})
	}
}
