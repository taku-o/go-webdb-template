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

func setupEmailE2EServer(t *testing.T) *httptest.Server {
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
	r := router.NewRouter(dmUserHandler, dmPostHandler, todayHandler, emailHandler, cfg)

	return httptest.NewServer(r)
}

// TestEmailE2E_SendWelcomeEmail はウェルカムメール送信のE2Eテスト
func TestEmailE2E_SendWelcomeEmail(t *testing.T) {
	server := setupEmailE2EServer(t)
	defer server.Close()

	// Step 1: ユーザーを作成
	createUserReq := map[string]string{
		"name":  "E2E Email Test User",
		"email": "e2e-email@example.com",
	}
	userBody, _ := json.Marshal(createUserReq)

	userResp, err := doRequestWithAuth("POST", server.URL+"/api/dm-users", userBody)
	require.NoError(t, err)
	defer userResp.Body.Close()
	assert.Equal(t, http.StatusCreated, userResp.StatusCode)

	var createdUser map[string]interface{}
	err = json.NewDecoder(userResp.Body).Decode(&createdUser)
	require.NoError(t, err)

	// Step 2: ウェルカムメールを送信
	emailReq := map[string]interface{}{
		"to":       []string{createdUser["email"].(string)},
		"template": "welcome",
		"data": map[string]interface{}{
			"Name":  createdUser["name"],
			"Email": createdUser["email"],
		},
	}
	emailBody, _ := json.Marshal(emailReq)

	emailResp, err := doRequestWithAuth("POST", server.URL+"/api/email/send", emailBody)
	require.NoError(t, err)
	defer emailResp.Body.Close()

	assert.Equal(t, http.StatusOK, emailResp.StatusCode)

	var emailResult map[string]interface{}
	err = json.NewDecoder(emailResp.Body).Decode(&emailResult)
	require.NoError(t, err)
	assert.Equal(t, true, emailResult["success"])
	assert.NotEmpty(t, emailResult["message"])
}

// TestEmailE2E_SendEmailWithoutAuth は認証なしでメール送信できないことを確認
func TestEmailE2E_SendEmailWithoutAuth(t *testing.T) {
	server := setupEmailE2EServer(t)
	defer server.Close()

	// メール送信リクエスト（認証なし）
	emailReq := map[string]interface{}{
		"to":       []string{"test@example.com"},
		"template": "welcome",
		"data": map[string]interface{}{
			"Name":  "Test User",
			"Email": "test@example.com",
		},
	}
	emailBody, _ := json.Marshal(emailReq)

	req, err := http.NewRequest("POST", server.URL+"/api/email/send", bytes.NewReader(emailBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	// Authorization headerなし

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

// TestEmailE2E_SendEmailValidation はメール送信のバリデーションをテスト
func TestEmailE2E_SendEmailValidation(t *testing.T) {
	server := setupEmailE2EServer(t)
	defer server.Close()

	tests := []struct {
		name           string
		request        map[string]interface{}
		expectedStatus int
	}{
		{
			name: "不正なメールアドレス",
			request: map[string]interface{}{
				"to":       []string{"not-an-email"},
				"template": "welcome",
				"data": map[string]interface{}{
					"Name":  "Test",
					"Email": "not-an-email",
				},
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "存在しないテンプレート",
			request: map[string]interface{}{
				"to":       []string{"valid@example.com"},
				"template": "does-not-exist",
				"data": map[string]interface{}{
					"Name":  "Test",
					"Email": "valid@example.com",
				},
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.request)
			resp, err := doRequestWithAuth("POST", server.URL+"/api/email/send", body)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}

// TestEmailE2E_FullWorkflow はユーザー作成からメール送信までの完全なワークフローをテスト
func TestEmailE2E_FullWorkflow(t *testing.T) {
	server := setupEmailE2EServer(t)
	defer server.Close()

	// Step 1: 複数のユーザーを作成
	users := []map[string]string{
		{"name": "User One", "email": "user1@example.com"},
		{"name": "User Two", "email": "user2@example.com"},
		{"name": "User Three", "email": "user3@example.com"},
	}

	var createdUsers []map[string]interface{}
	for _, user := range users {
		body, _ := json.Marshal(user)
		resp, err := doRequestWithAuth("POST", server.URL+"/api/dm-users", body)
		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var created map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&created)
		require.NoError(t, err)
		resp.Body.Close()

		createdUsers = append(createdUsers, created)
	}

	// Step 2: 全ユーザーにウェルカムメールを送信
	var emails []string
	for _, user := range createdUsers {
		emails = append(emails, user["email"].(string))
	}

	emailReq := map[string]interface{}{
		"to":       emails,
		"template": "welcome",
		"data": map[string]interface{}{
			"Name":  "ユーザーの皆様",
			"Email": "users@example.com",
		},
	}
	emailBody, _ := json.Marshal(emailReq)

	emailResp, err := doRequestWithAuth("POST", server.URL+"/api/email/send", emailBody)
	require.NoError(t, err)
	defer emailResp.Body.Close()

	assert.Equal(t, http.StatusOK, emailResp.StatusCode)

	var emailResult map[string]interface{}
	err = json.NewDecoder(emailResp.Body).Decode(&emailResult)
	require.NoError(t, err)
	assert.Equal(t, true, emailResult["success"])
}
