package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/example/go-webdb-template/internal/config"
)

const testSecretKey = "test-secret-key-for-jwt-signing"

func createTestToken(claims *JWTClaims, secretKey string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(secretKey))
	return tokenString
}

func TestJWTValidator_ValidateJWT(t *testing.T) {
	cfg := &config.APIConfig{
		SecretKey:       testSecretKey,
		CurrentVersion:  "v2",
		InvalidVersions: []string{"v1"},
	}

	tests := []struct {
		name      string
		claims    *JWTClaims
		secretKey string
		env       string
		wantErr   bool
		errMsg    string
	}{
		{
			name: "valid public token",
			claims: &JWTClaims{
				Issuer:   "go-webdb-template",
				Subject:  "public_client",
				Type:     "public",
				Scope:    []string{"read", "write"},
				IssuedAt: time.Now().Unix(),
				Version:  "v2",
				Env:      "develop",
			},
			secretKey: testSecretKey,
			env:       "develop",
			wantErr:   false,
		},
		{
			name: "valid private token",
			claims: &JWTClaims{
				Issuer:   "go-webdb-template",
				Subject:  "user123",
				Type:     "private",
				Scope:    []string{"read", "write"},
				IssuedAt: time.Now().Unix(),
				Version:  "v2",
				Env:      "develop",
			},
			secretKey: testSecretKey,
			env:       "develop",
			wantErr:   false,
		},
		{
			name: "invalid signature",
			claims: &JWTClaims{
				Issuer:   "go-webdb-template",
				Subject:  "public_client",
				Type:     "public",
				Scope:    []string{"read", "write"},
				IssuedAt: time.Now().Unix(),
				Version:  "v2",
				Env:      "develop",
			},
			secretKey: "wrong-secret-key",
			env:       "develop",
			wantErr:   true,
			errMsg:    "failed to parse JWT",
		},
		{
			name: "invalid issuer",
			claims: &JWTClaims{
				Issuer:   "wrong-issuer",
				Subject:  "public_client",
				Type:     "public",
				Scope:    []string{"read", "write"},
				IssuedAt: time.Now().Unix(),
				Version:  "v2",
				Env:      "develop",
			},
			secretKey: testSecretKey,
			env:       "develop",
			wantErr:   true,
			errMsg:    "invalid issuer",
		},
		{
			name: "invalid type",
			claims: &JWTClaims{
				Issuer:   "go-webdb-template",
				Subject:  "public_client",
				Type:     "invalid",
				Scope:    []string{"read", "write"},
				IssuedAt: time.Now().Unix(),
				Version:  "v2",
				Env:      "develop",
			},
			secretKey: testSecretKey,
			env:       "develop",
			wantErr:   true,
			errMsg:    "invalid token type",
		},
		{
			name: "invalid version (v1 is invalidated)",
			claims: &JWTClaims{
				Issuer:   "go-webdb-template",
				Subject:  "public_client",
				Type:     "public",
				Scope:    []string{"read", "write"},
				IssuedAt: time.Now().Unix(),
				Version:  "v1",
				Env:      "develop",
			},
			secretKey: testSecretKey,
			env:       "develop",
			wantErr:   true,
			errMsg:    "invalid token version",
		},
		{
			name: "environment mismatch",
			claims: &JWTClaims{
				Issuer:   "go-webdb-template",
				Subject:  "public_client",
				Type:     "public",
				Scope:    []string{"read", "write"},
				IssuedAt: time.Now().Unix(),
				Version:  "v2",
				Env:      "production",
			},
			secretKey: testSecretKey,
			env:       "develop",
			wantErr:   true,
			errMsg:    "token environment mismatch",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := NewJWTValidator(cfg, tt.env)
			tokenString := createTestToken(tt.claims, tt.secretKey)

			claims, err := validator.ValidateJWT(tokenString)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				assert.Nil(t, claims)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, claims)
				assert.Equal(t, tt.claims.Issuer, claims.Issuer)
				assert.Equal(t, tt.claims.Subject, claims.Subject)
				assert.Equal(t, tt.claims.Type, claims.Type)
				assert.Equal(t, tt.claims.Scope, claims.Scope)
				assert.Equal(t, tt.claims.Version, claims.Version)
				assert.Equal(t, tt.claims.Env, claims.Env)
			}
		})
	}
}

func TestJWTValidator_IsVersionInvalid(t *testing.T) {
	cfg := &config.APIConfig{
		SecretKey:       testSecretKey,
		CurrentVersion:  "v2",
		InvalidVersions: []string{"v1", "v0"},
	}

	validator := NewJWTValidator(cfg, "develop")

	tests := []struct {
		version string
		want    bool
	}{
		{"v1", true},
		{"v0", true},
		{"v2", false},
		{"v3", false},
	}

	for _, tt := range tests {
		t.Run(tt.version, func(t *testing.T) {
			got := validator.IsVersionInvalid(tt.version)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestParseJWTClaims(t *testing.T) {
	claims := &JWTClaims{
		Issuer:   "go-webdb-template",
		Subject:  "public_client",
		Type:     "public",
		Scope:    []string{"read", "write"},
		IssuedAt: time.Now().Unix(),
		Version:  "v2",
		Env:      "develop",
	}

	tokenString := createTestToken(claims, testSecretKey)

	parsedClaims, err := ParseJWTClaims(tokenString)
	require.NoError(t, err)
	assert.Equal(t, claims.Issuer, parsedClaims.Issuer)
	assert.Equal(t, claims.Subject, parsedClaims.Subject)
	assert.Equal(t, claims.Type, parsedClaims.Type)
	assert.Equal(t, claims.Scope, parsedClaims.Scope)
	assert.Equal(t, claims.Version, parsedClaims.Version)
	assert.Equal(t, claims.Env, parsedClaims.Env)
}

func TestParseJWTClaims_InvalidToken(t *testing.T) {
	_, err := ParseJWTClaims("invalid-token")
	require.Error(t, err)
}
