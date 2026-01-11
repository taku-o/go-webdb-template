package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAPIKeyService_GenerateAPIKey(t *testing.T) {
	tests := []struct {
		name      string
		secretKey string
		version   string
		env       string
		issuedAt  int64
		wantErr   bool
	}{
		{
			name:      "正常系: APIキーを生成できる",
			secretKey: "test-secret-key-12345678901234567890",
			version:   "v1.0.0",
			env:       "develop",
			issuedAt:  time.Now().Unix(),
			wantErr:   false,
		},
		{
			name:      "正常系: production環境でAPIキーを生成できる",
			secretKey: "test-secret-key-12345678901234567890",
			version:   "v2.0.0",
			env:       "production",
			issuedAt:  time.Now().Unix(),
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewAPIKeyService()
			got, err := s.GenerateAPIKey(tt.secretKey, tt.version, tt.env, tt.issuedAt)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotEmpty(t, got)
			// JWTトークン形式（3つのドット区切り）であることを確認
			assert.Contains(t, got, ".")
		})
	}
}

func TestAPIKeyService_DecodeAPIKeyPayload(t *testing.T) {
	s := NewAPIKeyService()

	// テスト用トークンを生成
	secretKey := "test-secret-key-12345678901234567890"
	version := "v1.0.0"
	env := "develop"
	issuedAt := time.Now().Unix()

	validToken, err := s.GenerateAPIKey(secretKey, version, env, issuedAt)
	assert.NoError(t, err)

	tests := []struct {
		name        string
		token       string
		wantVersion string
		wantEnv     string
		wantErr     bool
	}{
		{
			name:        "正常系: 有効なトークンをデコードできる",
			token:       validToken,
			wantVersion: version,
			wantEnv:     env,
			wantErr:     false,
		},
		{
			name:    "異常系: 無効なトークン",
			token:   "invalid-token",
			wantErr: true,
		},
		{
			name:    "異常系: 空のトークン",
			token:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := s.DecodeAPIKeyPayload(tt.token)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, claims)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, claims)
			assert.Equal(t, tt.wantVersion, claims.Version)
			assert.Equal(t, tt.wantEnv, claims.Env)
			assert.Equal(t, "go-webdb-template", claims.Issuer)
			assert.Equal(t, "public", claims.Type)
		})
	}
}
