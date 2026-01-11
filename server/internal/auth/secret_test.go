package auth

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateSecretKey(t *testing.T) {
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
			got, err := GenerateSecretKey()

			if tt.wantError {
				assert.Error(t, err)
				assert.Empty(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, got)

				// Base64デコードして32バイトであることを確認
				decoded, err := base64.StdEncoding.DecodeString(got)
				assert.NoError(t, err)
				assert.Equal(t, 32, len(decoded))
			}
		})
	}
}

func TestGenerateSecretKey_Uniqueness(t *testing.T) {
	// 複数回生成して、それぞれ異なる秘密鍵が生成されることを確認
	secret1, err1 := GenerateSecretKey()
	assert.NoError(t, err1)

	secret2, err2 := GenerateSecretKey()
	assert.NoError(t, err2)

	assert.NotEqual(t, secret1, secret2, "generated secrets should be unique")
}
