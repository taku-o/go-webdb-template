package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/example/go-webdb-template/internal/config"
)

func TestAuthMiddleware_Middleware(t *testing.T) {
	cfg := &config.APIConfig{
		SecretKey:       testSecretKey,
		CurrentVersion:  "v2",
		InvalidVersions: []string{"v1"},
	}

	validClaims := &JWTClaims{
		Issuer:   "go-webdb-template",
		Subject:  "public_client",
		Type:     "public",
		Scope:    []string{"read", "write"},
		IssuedAt: time.Now().Unix(),
		Version:  "v2",
		Env:      "develop",
	}
	validToken := createTestToken(validClaims, testSecretKey)

	readOnlyClaims := &JWTClaims{
		Issuer:   "go-webdb-template",
		Subject:  "public_client",
		Type:     "public",
		Scope:    []string{"read"},
		IssuedAt: time.Now().Unix(),
		Version:  "v2",
		Env:      "develop",
	}
	readOnlyToken := createTestToken(readOnlyClaims, testSecretKey)

	tests := []struct {
		name           string
		method         string
		authHeader     string
		wantStatusCode int
		wantBody       string
	}{
		{
			name:           "valid token GET request",
			method:         "GET",
			authHeader:     "Bearer " + validToken,
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "valid token POST request",
			method:         "POST",
			authHeader:     "Bearer " + validToken,
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "no authorization header",
			method:         "GET",
			authHeader:     "",
			wantStatusCode: http.StatusUnauthorized,
			wantBody:       "Authorization header is required",
		},
		{
			name:           "invalid authorization header format",
			method:         "GET",
			authHeader:     "InvalidFormat",
			wantStatusCode: http.StatusUnauthorized,
			wantBody:       "Invalid authorization header format",
		},
		{
			name:           "invalid token",
			method:         "GET",
			authHeader:     "Bearer invalid-token",
			wantStatusCode: http.StatusUnauthorized,
			wantBody:       "Invalid API key",
		},
		{
			name:           "read only token GET request",
			method:         "GET",
			authHeader:     "Bearer " + readOnlyToken,
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "read only token POST request - insufficient scope",
			method:         "POST",
			authHeader:     "Bearer " + readOnlyToken,
			wantStatusCode: http.StatusForbidden,
			wantBody:       "Insufficient scope",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			middleware := NewAuthMiddleware(cfg, "develop")

			handler := middleware.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("OK"))
			}))

			req := httptest.NewRequest(tt.method, "/api/test", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.wantStatusCode, rr.Code)
			if tt.wantBody != "" {
				assert.Contains(t, rr.Body.String(), tt.wantBody)
			}
		})
	}
}

func TestValidateScope(t *testing.T) {
	cfg := &config.APIConfig{
		SecretKey:       testSecretKey,
		CurrentVersion:  "v2",
		InvalidVersions: []string{},
	}
	middleware := NewAuthMiddleware(cfg, "develop")

	tests := []struct {
		name    string
		scope   []string
		method  string
		wantErr bool
	}{
		{"read scope GET", []string{"read"}, "GET", false},
		{"read scope POST", []string{"read"}, "POST", true},
		{"write scope GET", []string{"write"}, "GET", true},
		{"write scope POST", []string{"write"}, "POST", false},
		{"write scope PUT", []string{"write"}, "PUT", false},
		{"write scope DELETE", []string{"write"}, "DELETE", false},
		{"read and write scope GET", []string{"read", "write"}, "GET", false},
		{"read and write scope POST", []string{"read", "write"}, "POST", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims := &JWTClaims{Scope: tt.scope}
			err := middleware.validateScope(claims, tt.method)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
