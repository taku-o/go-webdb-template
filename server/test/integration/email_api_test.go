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

func setupEmailTestServer(t *testing.T) *httptest.Server {
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

	// Setup router with test config
	cfg := testutil.GetTestConfig()
	r := router.NewRouter(dmUserHandler, dmPostHandler, todayHandler, emailHandler, nil, cfg)

	return httptest.NewServer(r)
}

func TestEmailAPI_SendEmail_Success(t *testing.T) {
	server := setupEmailTestServer(t)
	defer server.Close()

	// Get valid token
	token, err := testutil.GetTestAPIToken()
	require.NoError(t, err)

	// Prepare request body
	requestBody := map[string]interface{}{
		"to":       []string{"test@example.com"},
		"template": "welcome",
		"data": map[string]interface{}{
			"Name":  "テスト太郎",
			"Email": "test@example.com",
		},
	}
	body, err := json.Marshal(requestBody)
	require.NoError(t, err)

	// Send request
	req, err := http.NewRequest("POST", server.URL+"/api/email/send", bytes.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Parse response
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)
	assert.Equal(t, true, result["success"])
}

func TestEmailAPI_SendEmail_NoAuth(t *testing.T) {
	server := setupEmailTestServer(t)
	defer server.Close()

	// Prepare request body
	requestBody := map[string]interface{}{
		"to":       []string{"test@example.com"},
		"template": "welcome",
		"data": map[string]interface{}{
			"Name":  "テスト太郎",
			"Email": "test@example.com",
		},
	}
	body, err := json.Marshal(requestBody)
	require.NoError(t, err)

	// Send request without auth
	req, err := http.NewRequest("POST", server.URL+"/api/email/send", bytes.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestEmailAPI_SendEmail_InvalidEmail(t *testing.T) {
	server := setupEmailTestServer(t)
	defer server.Close()

	// Get valid token
	token, err := testutil.GetTestAPIToken()
	require.NoError(t, err)

	// Prepare request body with invalid email
	requestBody := map[string]interface{}{
		"to":       []string{"invalid-email"},
		"template": "welcome",
		"data": map[string]interface{}{
			"Name":  "テスト太郎",
			"Email": "invalid-email",
		},
	}
	body, err := json.Marshal(requestBody)
	require.NoError(t, err)

	// Send request
	req, err := http.NewRequest("POST", server.URL+"/api/email/send", bytes.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestEmailAPI_SendEmail_InvalidTemplate(t *testing.T) {
	server := setupEmailTestServer(t)
	defer server.Close()

	// Get valid token
	token, err := testutil.GetTestAPIToken()
	require.NoError(t, err)

	// Prepare request body with non-existent template
	requestBody := map[string]interface{}{
		"to":       []string{"test@example.com"},
		"template": "nonexistent",
		"data": map[string]interface{}{
			"Name":  "テスト太郎",
			"Email": "test@example.com",
		},
	}
	body, err := json.Marshal(requestBody)
	require.NoError(t, err)

	// Send request
	req, err := http.NewRequest("POST", server.URL+"/api/email/send", bytes.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestEmailAPI_SendEmail_MultipleRecipients(t *testing.T) {
	server := setupEmailTestServer(t)
	defer server.Close()

	// Get valid token
	token, err := testutil.GetTestAPIToken()
	require.NoError(t, err)

	// Prepare request body with multiple recipients
	requestBody := map[string]interface{}{
		"to":       []string{"user1@example.com", "user2@example.com", "user3@example.com"},
		"template": "welcome",
		"data": map[string]interface{}{
			"Name":  "テスト太郎",
			"Email": "user1@example.com",
		},
	}
	body, err := json.Marshal(requestBody)
	require.NoError(t, err)

	// Send request
	req, err := http.NewRequest("POST", server.URL+"/api/email/send", bytes.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Parse response
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)
	assert.Equal(t, true, result["success"])
}
