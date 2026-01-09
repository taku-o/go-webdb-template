package integration_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/taku-o/go-webdb-template/internal/api/handler"
	"github.com/taku-o/go-webdb-template/internal/api/router"
	"github.com/taku-o/go-webdb-template/internal/config"
	"github.com/taku-o/go-webdb-template/internal/repository"
	"github.com/taku-o/go-webdb-template/internal/service"
	"github.com/taku-o/go-webdb-template/internal/service/email"
	"github.com/taku-o/go-webdb-template/test/testutil"
)

func setupJobqueueTestServer(t *testing.T, withJobqueueHandler bool) *httptest.Server {
	// Setup test database with GroupManager
	groupManager := testutil.SetupTestGroupManager(t, 4, 8)
	t.Cleanup(func() {
		testutil.CleanupTestGroupManager(groupManager)
	})

	// Initialize layers (using GORM repositories)
	dmUserRepo := repository.NewDmUserRepositoryGORM(groupManager)
	dmUserService := service.NewDmUserService(dmUserRepo)
	dmUserHandler := handler.NewDmUserHandler(dmUserService)

	dmPostRepo := repository.NewDmPostRepositoryGORM(groupManager)
	dmPostService := service.NewDmPostService(dmPostRepo, dmUserRepo)
	dmPostHandler := handler.NewDmPostHandler(dmPostService)

	// TodayHandler
	todayHandler := handler.NewTodayHandler()

	// EmailHandler（MockSenderを使用）
	emailCfg := &config.EmailConfig{
		SenderType: "mock",
	}
	emailService, err := email.NewEmailService(emailCfg, nil)
	require.NoError(t, err)
	templateService := email.NewTemplateService()
	emailHandler := handler.NewEmailHandler(emailService, templateService)

	// JobqueueHandler（Redisなしでテストする場合はnil）
	var dmJobqueueHandler *handler.DmJobqueueHandler
	if withJobqueueHandler {
		// nilクライアントでハンドラーを作成（Redis接続なしのテスト用）
		dmJobqueueHandler = handler.NewDmJobqueueHandler(nil)
	}

	// Setup router with test config
	cfg := testutil.GetTestConfig()
	r := router.NewRouter(dmUserHandler, dmPostHandler, todayHandler, emailHandler, dmJobqueueHandler, cfg)

	return httptest.NewServer(r)
}

func TestJobqueueAPI_RegisterJob_ServiceUnavailable(t *testing.T) {
	// nilクライアントでサーバーを起動（Redisなし）
	server := setupJobqueueTestServer(t, true)
	defer server.Close()

	// Get valid token
	token, err := testutil.GetTestAPIToken()
	require.NoError(t, err)

	// Prepare request body
	requestBody := map[string]interface{}{
		"message":       "Test job message",
		"delay_seconds": 10,
	}
	body, err := json.Marshal(requestBody)
	require.NoError(t, err)

	// Send request
	req, err := http.NewRequest("POST", server.URL+"/api/dm-jobqueue/register", bytes.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Redis接続なしの場合は503エラー
	assert.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)
}

func TestJobqueueAPI_RegisterJob_WithoutHandler(t *testing.T) {
	// ハンドラーなしでサーバーを起動
	server := setupJobqueueTestServer(t, false)
	defer server.Close()

	// Get valid token
	token, err := testutil.GetTestAPIToken()
	require.NoError(t, err)

	// Prepare request body
	requestBody := map[string]interface{}{
		"message": "Test job message",
	}
	body, err := json.Marshal(requestBody)
	require.NoError(t, err)

	// Send request
	req, err := http.NewRequest("POST", server.URL+"/api/dm-jobqueue/register", bytes.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// ハンドラーが登録されていない場合は404エラー
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestJobqueueAPI_RegisterJob_EmptyBody(t *testing.T) {
	// nilクライアントでサーバーを起動
	server := setupJobqueueTestServer(t, true)
	defer server.Close()

	// Get valid token
	token, err := testutil.GetTestAPIToken()
	require.NoError(t, err)

	// 空のリクエストボディ
	requestBody := map[string]interface{}{}
	body, err := json.Marshal(requestBody)
	require.NoError(t, err)

	// Send request
	req, err := http.NewRequest("POST", server.URL+"/api/dm-jobqueue/register", bytes.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Redis接続なしの場合は503エラー（ボディが空でもnilクライアントチェックが先に行われる）
	assert.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)
}

func TestJobqueueAPI_RegisterJob_WithCustomDelaySeconds(t *testing.T) {
	// nilクライアントでサーバーを起動
	server := setupJobqueueTestServer(t, true)
	defer server.Close()

	// Get valid token
	token, err := testutil.GetTestAPIToken()
	require.NoError(t, err)

	// カスタム遅延時間を含むリクエストボディ
	requestBody := map[string]interface{}{
		"message":       "Test with custom delay",
		"delay_seconds": 60,
		"max_retry":     3,
	}
	body, err := json.Marshal(requestBody)
	require.NoError(t, err)

	// Send request
	req, err := http.NewRequest("POST", server.URL+"/api/dm-jobqueue/register", bytes.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Redis接続なしの場合は503エラー
	assert.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)
}

func TestJobqueueAPI_RegisterJob_UnauthorizedWithoutToken(t *testing.T) {
	// サーバーを起動
	server := setupJobqueueTestServer(t, true)
	defer server.Close()

	// リクエストボディ
	requestBody := map[string]interface{}{
		"message": "Test job message",
	}
	body, err := json.Marshal(requestBody)
	require.NoError(t, err)

	// トークンなしでリクエスト
	req, err := http.NewRequest("POST", server.URL+"/api/dm-jobqueue/register", bytes.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// 認証エラー
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}
