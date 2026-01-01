package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/taku-o/go-webdb-template/internal/config"
)

// テスト用の設定値（ミドルウェアテスト用）
const (
	mwTestSecretKey          = "test-secret-key-for-jwt-signing"
	mwTestEnv                = "develop"
	mwTestAuth0IssuerBaseURL = "https://dev-oaa5vtzmld4dsxtd.jp.auth0.com"
)

// getTestAPIConfig はテスト用のAPI設定を返す
func getTestAPIConfig() *config.APIConfig {
	return &config.APIConfig{
		CurrentVersion:     "v2",
		SecretKey:          mwTestSecretKey,
		InvalidVersions:    []string{"v1"},
		Auth0IssuerBaseURL: mwTestAuth0IssuerBaseURL,
	}
}

// getTestAPIToken はテスト用のAPIトークンを生成
func getTestAPIToken() (string, error) {
	return GeneratePublicAPIKey(mwTestSecretKey, "v2", mwTestEnv, time.Now().Unix())
}

func TestValidateScope(t *testing.T) {
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
			err := validateScope(claims, tt.method)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// TestNewEchoAuthMiddleware はEcho用認証ミドルウェアの作成を確認
func TestNewEchoAuthMiddleware(t *testing.T) {
	cfg := getTestAPIConfig()

	middleware := NewEchoAuthMiddleware(cfg, mwTestEnv, mwTestAuth0IssuerBaseURL)
	require.NotNil(t, middleware)
}

// TestEchoAuthMiddleware_NoAuthHeader は認証ヘッダーがない場合に401を返すことを確認
func TestEchoAuthMiddleware_NoAuthHeader(t *testing.T) {
	cfg := getTestAPIConfig()

	middleware := NewEchoAuthMiddleware(cfg, mwTestEnv, mwTestAuth0IssuerBaseURL)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/upload/dm_movie", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := middleware(func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	err := handler(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

// TestEchoAuthMiddleware_InvalidAuthHeader は無効な認証ヘッダーの場合に401を返すことを確認
func TestEchoAuthMiddleware_InvalidAuthHeader(t *testing.T) {
	cfg := getTestAPIConfig()

	middleware := NewEchoAuthMiddleware(cfg, mwTestEnv, mwTestAuth0IssuerBaseURL)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/upload/dm_movie", nil)
	req.Header.Set("Authorization", "InvalidToken")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := middleware(func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	err := handler(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

// TestEchoAuthMiddleware_ValidToken は有効なトークンで認証が成功することを確認
func TestEchoAuthMiddleware_ValidToken(t *testing.T) {
	cfg := getTestAPIConfig()

	middleware := NewEchoAuthMiddleware(cfg, mwTestEnv, mwTestAuth0IssuerBaseURL)

	token, err := getTestAPIToken()
	require.NoError(t, err)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/upload/dm_movie", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := middleware(func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	err = handler(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

// TestEchoAuthMiddleware_AllTUSMethods はTUSプロトコルの全メソッドで認証が動作することを確認
func TestEchoAuthMiddleware_AllTUSMethods(t *testing.T) {
	cfg := getTestAPIConfig()

	middleware := NewEchoAuthMiddleware(cfg, mwTestEnv, mwTestAuth0IssuerBaseURL)

	token, err := getTestAPIToken()
	require.NoError(t, err)

	methods := []string{
		http.MethodOptions,
		http.MethodPost,
		http.MethodPatch,
		http.MethodHead,
		http.MethodDelete,
	}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(method, "/api/upload/dm_movie", nil)
			req.Header.Set("Authorization", "Bearer "+token)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			handler := middleware(func(c echo.Context) error {
				return c.String(http.StatusOK, "OK")
			})

			err := handler(c)
			require.NoError(t, err)
			assert.Equal(t, http.StatusOK, rec.Code, "method %s should succeed with valid token", method)
		})
	}
}
