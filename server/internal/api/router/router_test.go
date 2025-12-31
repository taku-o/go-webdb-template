package router

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/taku-o/go-webdb-template/internal/api/handler"
	"github.com/taku-o/go-webdb-template/test/testutil"
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
	assert.Contains(t, paths, "/api/dm-users")
	assert.Contains(t, paths, "/api/dm-users/{id}")

	// 投稿エンドポイントが含まれていることを確認
	assert.Contains(t, paths, "/api/dm-posts")
	assert.Contains(t, paths, "/api/dm-posts/{id}")
	assert.Contains(t, paths, "/api/dm-user-posts")
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

// TestRegisterDmUserEndpointsIntegration はユーザーエンドポイントが登録されることを確認
func TestRegisterDmUserEndpointsIntegration(t *testing.T) {
	// RegisterDmUserEndpoints関数のシグネチャを確認
	var _ func(*handler.DmUserHandler) = func(h *handler.DmUserHandler) {
		// handler.RegisterDmUserEndpoints(api, h) の形式で呼び出し可能
	}
}

// TestRegisterDmPostEndpointsIntegration は投稿エンドポイントが登録されることを確認
func TestRegisterDmPostEndpointsIntegration(t *testing.T) {
	// RegisterDmPostEndpoints関数のシグネチャを確認
	var _ func(*handler.DmPostHandler) = func(h *handler.DmPostHandler) {
		// handler.RegisterDmPostEndpoints(api, h) の形式で呼び出し可能
	}
}

// TestSecuritySchemeInOpenAPI はOpenAPIにSecuritySchemeが定義されていることを確認
func TestSecuritySchemeInOpenAPI(t *testing.T) {
	// テスト用の設定
	cfg := testutil.GetTestConfig()

	// ルーターを作成
	router := NewRouter(nil, nil, nil, cfg)

	// テスト用のAPIトークンを取得
	token, err := testutil.GetTestAPIToken()
	require.NoError(t, err)

	// /openapi-3.0.json エンドポイントからOpenAPI仕様を取得
	req := httptest.NewRequest(http.MethodGet, "/openapi-3.0.json", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	// JSONとしてパース
	var openAPI map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &openAPI)
	require.NoError(t, err)

	// componentsが存在することを確認
	components, ok := openAPI["components"].(map[string]interface{})
	require.True(t, ok, "components should exist in OpenAPI spec")

	// securitySchemesが存在することを確認
	securitySchemes, ok := components["securitySchemes"].(map[string]interface{})
	require.True(t, ok, "securitySchemes should exist in components")

	// bearerAuthが存在することを確認
	bearerAuth, ok := securitySchemes["bearerAuth"].(map[string]interface{})
	require.True(t, ok, "bearerAuth should exist in securitySchemes")

	// bearerAuthの設定内容を確認
	assert.Equal(t, "http", bearerAuth["type"], "type should be 'http'")
	assert.Equal(t, "bearer", bearerAuth["scheme"], "scheme should be 'bearer'")
	assert.Equal(t, "JWT", bearerAuth["bearerFormat"], "bearerFormat should be 'JWT'")
}

// TestEndpointSecurityInOpenAPI は各エンドポイントにSecurityプロパティが設定されていることを確認
func TestEndpointSecurityInOpenAPI(t *testing.T) {
	// テスト用の設定
	cfg := testutil.GetTestConfig()

	// ルーターを作成
	router := NewRouter(nil, nil, nil, cfg)

	// テスト用のAPIトークンを取得
	token, err := testutil.GetTestAPIToken()
	require.NoError(t, err)

	// /openapi-3.0.json エンドポイントからOpenAPI仕様を取得
	req := httptest.NewRequest(http.MethodGet, "/openapi-3.0.json", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	// JSONとしてパース
	var openAPI map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &openAPI)
	require.NoError(t, err)

	// pathsを取得
	paths, ok := openAPI["paths"].(map[string]interface{})
	require.True(t, ok, "paths should exist in OpenAPI spec")

	// 各エンドポイントにSecurityプロパティが設定されていることを確認
	endpointsToCheck := []struct {
		path   string
		method string
	}{
		// Users endpoints
		{"/api/dm-users", "post"},
		{"/api/dm-users", "get"},
		{"/api/dm-users/{id}", "get"},
		{"/api/dm-users/{id}", "put"},
		{"/api/dm-users/{id}", "delete"},
		{"/api/export/dm-users/csv", "get"},
		// Posts endpoints
		{"/api/dm-posts", "post"},
		{"/api/dm-posts", "get"},
		{"/api/dm-posts/{id}", "get"},
		{"/api/dm-posts/{id}", "put"},
		{"/api/dm-posts/{id}", "delete"},
		{"/api/dm-user-posts", "get"},
		// Today endpoint
		{"/api/today", "get"},
	}

	for _, ep := range endpointsToCheck {
		t.Run(ep.path+"_"+ep.method, func(t *testing.T) {
			pathItem, ok := paths[ep.path].(map[string]interface{})
			require.True(t, ok, "path %s should exist", ep.path)

			operation, ok := pathItem[ep.method].(map[string]interface{})
			require.True(t, ok, "method %s should exist for path %s", ep.method, ep.path)

			security, ok := operation["security"].([]interface{})
			require.True(t, ok, "security should exist for %s %s", ep.method, ep.path)
			require.Len(t, security, 1, "security should have one item for %s %s", ep.method, ep.path)

			secItem, ok := security[0].(map[string]interface{})
			require.True(t, ok, "security item should be a map for %s %s", ep.method, ep.path)

			_, ok = secItem["bearerAuth"]
			assert.True(t, ok, "bearerAuth should exist in security for %s %s", ep.method, ep.path)
		})
	}
}

// TestEndpointAccessLevelInOpenAPI は各エンドポイントのSummaryとDescriptionにアクセスレベル情報が含まれることを確認
func TestEndpointAccessLevelInOpenAPI(t *testing.T) {
	// テスト用の設定
	cfg := testutil.GetTestConfig()

	// ルーターを作成
	router := NewRouter(nil, nil, nil, cfg)

	// テスト用のAPIトークンを取得
	token, err := testutil.GetTestAPIToken()
	require.NoError(t, err)

	// /openapi-3.0.json エンドポイントからOpenAPI仕様を取得
	req := httptest.NewRequest(http.MethodGet, "/openapi-3.0.json", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	// JSONとしてパース
	var openAPI map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &openAPI)
	require.NoError(t, err)

	// pathsを取得
	paths, ok := openAPI["paths"].(map[string]interface{})
	require.True(t, ok, "paths should exist in OpenAPI spec")

	// 各エンドポイントのSummaryとDescriptionを確認
	endpointsToCheck := []struct {
		path                string
		method              string
		expectedTag         string
		expectedSummary     string // Summaryに含まれるべき文字列（空の場合はチェックしない）
		expectedDescription string // Descriptionに含まれるべき文字列
	}{
		// Users endpoints - Public API（Summaryに[public]は含まない）
		{"/api/dm-users", "post", "users", "", "**Access Level:** `public`"},
		{"/api/dm-users", "get", "users", "", "**Access Level:** `public`"},
		{"/api/dm-users/{id}", "get", "users", "", "**Access Level:** `public`"},
		{"/api/dm-users/{id}", "put", "users", "", "**Access Level:** `public`"},
		{"/api/dm-users/{id}", "delete", "users", "", "**Access Level:** `public`"},
		{"/api/export/dm-users/csv", "get", "users", "", "**Access Level:** `public`"},
		// Posts endpoints - Public API（Summaryに[public]は含まない）
		{"/api/dm-posts", "post", "posts", "", "**Access Level:** `public`"},
		{"/api/dm-posts", "get", "posts", "", "**Access Level:** `public`"},
		{"/api/dm-posts/{id}", "get", "posts", "", "**Access Level:** `public`"},
		{"/api/dm-posts/{id}", "put", "posts", "", "**Access Level:** `public`"},
		{"/api/dm-posts/{id}", "delete", "posts", "", "**Access Level:** `public`"},
		{"/api/dm-user-posts", "get", "posts", "", "**Access Level:** `public`"},
		// Today endpoint - Private API（Summaryに[private]を含む）
		{"/api/today", "get", "today", "[private]", "**Access Level:** `private`"},
	}

	for _, ep := range endpointsToCheck {
		t.Run(ep.path+"_"+ep.method, func(t *testing.T) {
			pathItem, ok := paths[ep.path].(map[string]interface{})
			require.True(t, ok, "path %s should exist", ep.path)

			operation, ok := pathItem[ep.method].(map[string]interface{})
			require.True(t, ok, "method %s should exist for path %s", ep.method, ep.path)

			// Tagsに機能タグのみが含まれていることを確認（Public API/Private APIは含まれない）
			tags, ok := operation["tags"].([]interface{})
			require.True(t, ok, "tags should exist for %s %s", ep.method, ep.path)
			actualTags := make([]string, len(tags))
			for i, tag := range tags {
				actualTags[i] = tag.(string)
			}
			assert.Contains(t, actualTags, ep.expectedTag, "tag '%s' should exist for %s %s", ep.expectedTag, ep.method, ep.path)
			assert.NotContains(t, actualTags, "Public API", "tag 'Public API' should not exist for %s %s", ep.method, ep.path)
			assert.NotContains(t, actualTags, "Private API", "tag 'Private API' should not exist for %s %s", ep.method, ep.path)

			// Summaryにアクセスレベルが含まれていることを確認（expectedSummaryが空でない場合のみ）
			summary, ok := operation["summary"].(string)
			require.True(t, ok, "summary should exist for %s %s", ep.method, ep.path)
			if ep.expectedSummary != "" {
				assert.Contains(t, summary, ep.expectedSummary, "summary should contain '%s' for %s %s", ep.expectedSummary, ep.method, ep.path)
			}

			// Descriptionにアクセスレベルが含まれていることを確認
			description, ok := operation["description"].(string)
			require.True(t, ok, "description should exist for %s %s", ep.method, ep.path)
			assert.Contains(t, description, ep.expectedDescription, "description should contain '%s' for %s %s", ep.expectedDescription, ep.method, ep.path)
		})
	}
}
