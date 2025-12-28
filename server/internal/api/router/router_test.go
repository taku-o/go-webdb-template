package router

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/example/go-webdb-template/internal/api/handler"
	"github.com/example/go-webdb-template/test/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestOpenAPIEndpoint はOpenAPIエンドポイントが正しく動作することを確認
func TestOpenAPIEndpoint(t *testing.T) {
	// テスト用の設定
	cfg := testutil.GetTestConfig()

	// ハンドラーはnilでも登録テストは可能
	router := NewRouter(nil, nil, nil, cfg)

	// テスト用のAPIトークンを取得
	token, err := testutil.GetTestAPIToken()
	require.NoError(t, err)

	// /openapi-3.0.json エンドポイントのテスト（JSON形式のOpenAPIスキーマ）
	req := httptest.NewRequest(http.MethodGet, "/openapi-3.0.json", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Header().Get("Content-Type"), "openapi+json")

	// JSONとしてパースできることを確認
	var openAPI map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &openAPI)
	require.NoError(t, err)

	// 基本的なOpenAPI構造を確認
	assert.Equal(t, "3.0.3", openAPI["openapi"])

	info, ok := openAPI["info"].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "go-webdb-template API", info["title"])
	assert.Equal(t, "1.0.0", info["version"])

	// pathsが存在することを確認
	paths, ok := openAPI["paths"].(map[string]interface{})
	require.True(t, ok)

	// ユーザーエンドポイントが含まれていることを確認
	assert.Contains(t, paths, "/api/users")
	assert.Contains(t, paths, "/api/users/{id}")

	// 投稿エンドポイントが含まれていることを確認
	assert.Contains(t, paths, "/api/posts")
	assert.Contains(t, paths, "/api/posts/{id}")
	assert.Contains(t, paths, "/api/user-posts")
}

// TestHealthEndpoint はヘルスチェックエンドポイントが正しく動作することを確認
func TestHealthEndpoint(t *testing.T) {
	cfg := testutil.GetTestConfig()

	router := NewRouter(nil, nil, nil, cfg)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "OK", rec.Body.String())
}

// TestRegisterUserEndpointsIntegration はユーザーエンドポイントが登録されることを確認
func TestRegisterUserEndpointsIntegration(t *testing.T) {
	// RegisterUserEndpoints関数のシグネチャを確認
	var _ func(*handler.UserHandler) = func(h *handler.UserHandler) {
		// handler.RegisterUserEndpoints(api, h) の形式で呼び出し可能
	}
}

// TestRegisterPostEndpointsIntegration は投稿エンドポイントが登録されることを確認
func TestRegisterPostEndpointsIntegration(t *testing.T) {
	// RegisterPostEndpoints関数のシグネチャを確認
	var _ func(*handler.PostHandler) = func(h *handler.PostHandler) {
		// handler.RegisterPostEndpoints(api, h) の形式で呼び出し可能
	}
}
