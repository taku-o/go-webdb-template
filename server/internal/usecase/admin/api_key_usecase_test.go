package admin

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/taku-o/go-webdb-template/internal/auth"
	"github.com/taku-o/go-webdb-template/internal/config"
)

// MockAPIKeyService はAPIKeyServiceInterfaceのモック
type MockAPIKeyService struct {
	GenerateAPIKeyFunc      func(ctx context.Context, secretKey, version, env string, issuedAt int64) (string, error)
	DecodeAPIKeyPayloadFunc func(ctx context.Context, token string) (*auth.JWTClaims, error)
}

func (m *MockAPIKeyService) GenerateAPIKey(ctx context.Context, secretKey, version, env string, issuedAt int64) (string, error) {
	if m.GenerateAPIKeyFunc != nil {
		return m.GenerateAPIKeyFunc(ctx, secretKey, version, env, issuedAt)
	}
	return "", nil
}

func (m *MockAPIKeyService) DecodeAPIKeyPayload(ctx context.Context, token string) (*auth.JWTClaims, error) {
	if m.DecodeAPIKeyPayloadFunc != nil {
		return m.DecodeAPIKeyPayloadFunc(ctx, token)
	}
	return nil, nil
}

func TestAPIKeyUsecase_GenerateAPIKey(t *testing.T) {
	tests := []struct {
		name               string
		env                string
		generateAPIKeyFunc func(ctx context.Context, secretKey, version, env string, issuedAt int64) (string, error)
		wantToken          string
		wantErr            bool
		wantErrContains    string
	}{
		{
			name: "正常系: APIキーを生成できる",
			env:  "develop",
			generateAPIKeyFunc: func(ctx context.Context, secretKey, version, env string, issuedAt int64) (string, error) {
				return "generated-token", nil
			},
			wantToken: "generated-token",
			wantErr:   false,
		},
		{
			name: "正常系: 環境が空の場合はdevelopが使用される",
			env:  "",
			generateAPIKeyFunc: func(ctx context.Context, secretKey, version, env string, issuedAt int64) (string, error) {
				assert.Equal(t, "develop", env)
				return "generated-token", nil
			},
			wantToken: "generated-token",
			wantErr:   false,
		},
		{
			name: "異常系: service層からエラーが返された場合",
			env:  "develop",
			generateAPIKeyFunc: func(ctx context.Context, secretKey, version, env string, issuedAt int64) (string, error) {
				return "", errors.New("failed to generate token")
			},
			wantToken:       "",
			wantErr:         true,
			wantErrContains: "failed to generate token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockAPIKeyService{
				GenerateAPIKeyFunc: tt.generateAPIKeyFunc,
			}

			cfg := &config.Config{
				API: config.APIConfig{
					SecretKey:      "test-secret-key",
					CurrentVersion: "v1.0.0",
				},
			}

			u := NewAPIKeyUsecase(mockService, cfg)
			got, err := u.GenerateAPIKey(context.Background(), tt.env)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.wantErrContains != "" {
					assert.Contains(t, err.Error(), tt.wantErrContains)
				}
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.wantToken, got)
		})
	}
}

func TestAPIKeyUsecase_DecodeAPIKeyPayload(t *testing.T) {
	tests := []struct {
		name                    string
		token                   string
		decodeAPIKeyPayloadFunc func(ctx context.Context, token string) (*auth.JWTClaims, error)
		wantClaims              *auth.JWTClaims
		wantErr                 bool
		wantErrContains         string
	}{
		{
			name:  "正常系: トークンをデコードできる",
			token: "valid-token",
			decodeAPIKeyPayloadFunc: func(ctx context.Context, token string) (*auth.JWTClaims, error) {
				return &auth.JWTClaims{
					Issuer:  "go-webdb-template",
					Type:    "public",
					Version: "v1.0.0",
					Env:     "develop",
				}, nil
			},
			wantClaims: &auth.JWTClaims{
				Issuer:  "go-webdb-template",
				Type:    "public",
				Version: "v1.0.0",
				Env:     "develop",
			},
			wantErr: false,
		},
		{
			name:  "異常系: service層からエラーが返された場合",
			token: "invalid-token",
			decodeAPIKeyPayloadFunc: func(ctx context.Context, token string) (*auth.JWTClaims, error) {
				return nil, errors.New("failed to decode token")
			},
			wantClaims:      nil,
			wantErr:         true,
			wantErrContains: "failed to decode token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockAPIKeyService{
				DecodeAPIKeyPayloadFunc: tt.decodeAPIKeyPayloadFunc,
			}

			cfg := &config.Config{
				API: config.APIConfig{
					SecretKey:      "test-secret-key",
					CurrentVersion: "v1.0.0",
				},
			}

			u := NewAPIKeyUsecase(mockService, cfg)
			got, err := u.DecodeAPIKeyPayload(context.Background(), tt.token)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.wantErrContains != "" {
					assert.Contains(t, err.Error(), tt.wantErrContains)
				}
				assert.Nil(t, got)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, got)
			assert.Equal(t, tt.wantClaims.Issuer, got.Issuer)
			assert.Equal(t, tt.wantClaims.Type, got.Type)
			assert.Equal(t, tt.wantClaims.Version, got.Version)
			assert.Equal(t, tt.wantClaims.Env, got.Env)
		})
	}
}
