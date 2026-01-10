package e2e_test

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

func setupJobqueueE2EServer(t *testing.T, withJobqueueHandler bool) *httptest.Server {
	// Setup test database with GroupManager
	groupManager := testutil.SetupTestGroupManager(t, 4, 8)
	t.Cleanup(func() {
		testutil.CleanupTestGroupManager(groupManager)
	})

	// Initialize layers
	dmUserRepo := repository.NewDmUserRepository(groupManager)
	dmUserService := service.NewDmUserService(dmUserRepo)
	dmUserHandler := testutil.CreateDmUserHandler(dmUserService)

	dmPostRepo := repository.NewDmPostRepository(groupManager)
	dmPostService := service.NewDmPostService(dmPostRepo, dmUserRepo)
	dmPostHandler := testutil.CreateDmPostHandler(dmPostService)

	// TodayHandler
	todayHandler := testutil.CreateTodayHandler()

	// EmailHandler（MockSenderを使用）
	emailCfg := &config.EmailConfig{
		SenderType: "mock",
	}
	emailService, err := email.NewEmailService(emailCfg, nil)
	require.NoError(t, err)
	templateService := email.NewTemplateService()
	emailHandler := testutil.CreateEmailHandler(emailService, templateService)

	// JobqueueHandler（nilクライアントでテスト）
	var dmJobqueueHandler *handler.DmJobqueueHandler
	if withJobqueueHandler {
		dmJobqueueHandler = testutil.CreateDmJobqueueHandler(nil)
	}

	// Setup router with test config
	cfg := testutil.GetTestConfig()
	r := router.NewRouter(dmUserHandler, dmPostHandler, todayHandler, emailHandler, dmJobqueueHandler, cfg)

	return httptest.NewServer(r)
}

// TestJobqueueE2E_RegisterJob_ServiceUnavailable はRedis未接続時のジョブ登録E2Eテスト
func TestJobqueueE2E_RegisterJob_ServiceUnavailable(t *testing.T) {
	server := setupJobqueueE2EServer(t, true)
	defer server.Close()

	// Step 1: ジョブ登録を試行（nilクライアント）
	registerReq := map[string]interface{}{
		"message":       "E2E test job message",
		"delay_seconds": 30,
	}
	registerBody, _ := json.Marshal(registerReq)

	resp, err := doRequestWithAuth("POST", server.URL+"/api/dm-jobqueue/register", registerBody)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Redis未接続の場合は503エラー
	assert.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)

	// エラーレスポンスの確認
	var errorResp map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&errorResp)
	require.NoError(t, err)
	assert.Contains(t, errorResp["detail"], "Redis is not connected")
}

// TestJobqueueE2E_RegisterJob_EndpointNotRegistered はハンドラー未登録時のE2Eテスト
func TestJobqueueE2E_RegisterJob_EndpointNotRegistered(t *testing.T) {
	// ハンドラーなしでサーバーを起動
	server := setupJobqueueE2EServer(t, false)
	defer server.Close()

	// ジョブ登録を試行（エンドポイント未登録）
	registerReq := map[string]interface{}{
		"message": "E2E test job",
	}
	registerBody, _ := json.Marshal(registerReq)

	resp, err := doRequestWithAuth("POST", server.URL+"/api/dm-jobqueue/register", registerBody)
	require.NoError(t, err)
	defer resp.Body.Close()

	// エンドポイントが登録されていない場合は404エラー
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

// TestJobqueueE2E_RegisterJob_AuthenticationRequired は認証必須のE2Eテスト
func TestJobqueueE2E_RegisterJob_AuthenticationRequired(t *testing.T) {
	server := setupJobqueueE2EServer(t, true)
	defer server.Close()

	// 認証なしでリクエスト
	registerReq := map[string]interface{}{
		"message": "E2E test job",
	}
	registerBody, _ := json.Marshal(registerReq)

	req, err := http.NewRequest("POST", server.URL+"/api/dm-jobqueue/register", bytes.NewReader(registerBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// 認証なしの場合は401エラー
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

// TestJobqueueE2E_RegisterJob_WithAllOptions は全オプション指定のE2Eテスト
func TestJobqueueE2E_RegisterJob_WithAllOptions(t *testing.T) {
	server := setupJobqueueE2EServer(t, true)
	defer server.Close()

	// 全オプションを指定してジョブ登録を試行
	registerReq := map[string]interface{}{
		"message":       "Full options test",
		"delay_seconds": 120,
		"max_retry":     5,
	}
	registerBody, _ := json.Marshal(registerReq)

	resp, err := doRequestWithAuth("POST", server.URL+"/api/dm-jobqueue/register", registerBody)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Redis未接続の場合は503エラー（リクエスト自体は正常に処理される）
	assert.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)
}

// TestJobqueueE2E_RegisterJob_EmptyBody は空ボディでのE2Eテスト
func TestJobqueueE2E_RegisterJob_EmptyBody(t *testing.T) {
	server := setupJobqueueE2EServer(t, true)
	defer server.Close()

	// 空のリクエストボディ
	registerReq := map[string]interface{}{}
	registerBody, _ := json.Marshal(registerReq)

	resp, err := doRequestWithAuth("POST", server.URL+"/api/dm-jobqueue/register", registerBody)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Redis未接続の場合は503エラー（空ボディも許可）
	assert.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)
}

// TestJobqueueE2E_RegisterJob_MessageOnlyOption はメッセージのみ指定のE2Eテスト
func TestJobqueueE2E_RegisterJob_MessageOnlyOption(t *testing.T) {
	server := setupJobqueueE2EServer(t, true)
	defer server.Close()

	// メッセージのみ指定
	registerReq := map[string]interface{}{
		"message": "Message only test",
	}
	registerBody, _ := json.Marshal(registerReq)

	resp, err := doRequestWithAuth("POST", server.URL+"/api/dm-jobqueue/register", registerBody)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Redis未接続の場合は503エラー
	assert.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)
}
