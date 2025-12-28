package e2e_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/example/go-webdb-template/internal/api/handler"
	"github.com/example/go-webdb-template/internal/api/router"
	"github.com/example/go-webdb-template/internal/repository"
	"github.com/example/go-webdb-template/internal/service"
	"github.com/example/go-webdb-template/test/testutil"
)

// testAPIToken はテスト用のAPIトークン
var testAPIToken string

func init() {
	var err error
	testAPIToken, err = testutil.GetTestAPIToken()
	if err != nil {
		panic("Failed to generate test API token: " + err.Error())
	}
}

// doRequestWithAuth は認証ヘッダー付きのリクエストを実行
func doRequestWithAuth(method, url string, body []byte) (*http.Response, error) {
	var bodyReader io.Reader
	if body != nil {
		bodyReader = bytes.NewBuffer(body)
	}
	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+testAPIToken)
	return http.DefaultClient.Do(req)
}

func setupTestServer(t *testing.T) *httptest.Server {
	// Setup test database with GroupManager
	groupManager := testutil.SetupTestGroupManager(t, 4, 8)
	t.Cleanup(func() {
		testutil.CleanupTestGroupManager(groupManager)
	})

	// Initialize layers (using GORM repositories)
	userRepo := repository.NewUserRepositoryGORM(groupManager)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	postRepo := repository.NewPostRepositoryGORM(groupManager)
	postService := service.NewPostService(postRepo, userRepo)
	postHandler := handler.NewPostHandler(postService)

	// TodayHandler
	todayHandler := handler.NewTodayHandler()

	// Setup router with test config
	cfg := testutil.GetTestConfig()
	r := router.NewRouter(userHandler, postHandler, todayHandler, cfg)

	return httptest.NewServer(r)
}

func TestUserAPI_CreateAndRetrieve(t *testing.T) {
	server := setupTestServer(t)
	defer server.Close()

	// Create user
	createReq := map[string]string{
		"name":  "E2E Test User",
		"email": "e2e@example.com",
	}
	body, _ := json.Marshal(createReq)

	resp, err := doRequestWithAuth("POST", server.URL+"/api/users", body)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var user map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&user)
	require.NoError(t, err)
	assert.Equal(t, "E2E Test User", user["name"])
	assert.Equal(t, "e2e@example.com", user["email"])

	userIDStr := user["id"].(string)
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	require.NoError(t, err)

	// Retrieve user
	resp, err = doRequestWithAuth("GET", server.URL+fmt.Sprintf("/api/users/%d", userID), nil)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Debug: print response if not OK
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Logf("Expected 200 but got %d, body: %s, userID: %d", resp.StatusCode, string(body), userID)
		// Re-create reader for subsequent decode
		resp.Body = io.NopCloser(bytes.NewBuffer(body))
	}

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var retrieved map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&retrieved)
	require.NoError(t, err)
	assert.Equal(t, "E2E Test User", retrieved["name"])
	assert.Equal(t, "e2e@example.com", retrieved["email"])
}

func TestUserAPI_UpdateAndDelete(t *testing.T) {
	server := setupTestServer(t)
	defer server.Close()

	// Create user
	createReq := map[string]string{
		"name":  "Original Name",
		"email": "original@example.com",
	}
	body, _ := json.Marshal(createReq)
	resp, err := doRequestWithAuth("POST", server.URL+"/api/users", body)
	require.NoError(t, err)
	defer resp.Body.Close()

	var user map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&user)
	userIDStr := user["id"].(string)
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	require.NoError(t, err)

	// Update user
	updateReq := map[string]string{
		"name":  "Updated Name",
		"email": "updated@example.com",
	}
	body, _ = json.Marshal(updateReq)
	resp, err = doRequestWithAuth("PUT", server.URL+fmt.Sprintf("/api/users/%d", userID), body)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var updated map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&updated)
	assert.Equal(t, "Updated Name", updated["name"])

	// Delete user
	resp, err = doRequestWithAuth("DELETE", server.URL+fmt.Sprintf("/api/users/%d", userID), nil)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)

	// Verify deletion
	resp, err = doRequestWithAuth("GET", server.URL+fmt.Sprintf("/api/users/%d", userID), nil)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestPostAPI_CompleteFlow(t *testing.T) {
	server := setupTestServer(t)
	defer server.Close()

	// Create user first
	userReq := map[string]string{
		"name":  "Post Test User",
		"email": "posttest@example.com",
	}
	body, _ := json.Marshal(userReq)
	resp, err := doRequestWithAuth("POST", server.URL+"/api/users", body)
	require.NoError(t, err)
	defer resp.Body.Close()

	var user map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&user)
	userIDStr := user["id"].(string)
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	require.NoError(t, err)

	// Create post
	postReq := map[string]interface{}{
		"user_id": userID,
		"title":   "Test Post",
		"content": "Test content",
	}
	body, _ = json.Marshal(postReq)
	resp, err = doRequestWithAuth("POST", server.URL+"/api/posts", body)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var post map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&post)
	postIDStr := post["id"].(string)
	postID, err := strconv.ParseInt(postIDStr, 10, 64)
	require.NoError(t, err)
	assert.Equal(t, "Test Post", post["title"])

	// Get post
	resp, err = doRequestWithAuth("GET", server.URL+fmt.Sprintf("/api/posts/%d?user_id=%d", postID, userID), nil)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Get user posts (JOIN)
	resp, err = doRequestWithAuth("GET", server.URL+"/api/user-posts", nil)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var userPosts []map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&userPosts)
	assert.Greater(t, len(userPosts), 0)

	// Find our post in the results
	found := false
	for _, up := range userPosts {
		upPostIDStr := up["post_id"].(string)
		upPostID, _ := strconv.ParseInt(upPostIDStr, 10, 64)
		if upPostID == postID {
			upUserIDStr := up["user_id"].(string)
			upUserID, _ := strconv.ParseInt(upUserIDStr, 10, 64)
			assert.Equal(t, userID, upUserID)
			assert.Equal(t, "Post Test User", up["user_name"])
			assert.Equal(t, "Test Post", up["post_title"])
			found = true
			break
		}
	}
	assert.True(t, found, "Should find our post in user-posts results")
}
