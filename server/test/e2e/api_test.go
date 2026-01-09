package e2e_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/taku-o/go-webdb-template/internal/api/handler"
	"github.com/taku-o/go-webdb-template/internal/api/router"
	"github.com/taku-o/go-webdb-template/internal/repository"
	"github.com/taku-o/go-webdb-template/internal/service"
	"github.com/taku-o/go-webdb-template/internal/util/idgen"
	"github.com/taku-o/go-webdb-template/test/testutil"
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
	dmUserRepo := repository.NewDmUserRepositoryGORM(groupManager)
	dmUserService := service.NewDmUserService(dmUserRepo)
	dmUserHandler := handler.NewDmUserHandler(dmUserService)

	dmPostRepo := repository.NewDmPostRepositoryGORM(groupManager)
	dmPostService := service.NewDmPostService(dmPostRepo, dmUserRepo)
	dmPostHandler := handler.NewDmPostHandler(dmPostService)

	// TodayHandler
	todayHandler := handler.NewTodayHandler()

	// Setup router with test config
	cfg := testutil.GetTestConfig()
	r := router.NewRouter(dmUserHandler, dmPostHandler, todayHandler, nil, nil, cfg)

	return httptest.NewServer(r)
}

func TestDmUserAPI_CreateAndRetrieve(t *testing.T) {
	server := setupTestServer(t)
	defer server.Close()

	// ユニークなメールアドレスを生成
	uniqueID, err := idgen.GenerateUUIDv7()
	require.NoError(t, err)
	uniqueEmail := fmt.Sprintf("e2e-%s@example.com", uniqueID)

	// Create dm_user
	createReq := map[string]string{
		"name":  "E2E Test User",
		"email": uniqueEmail,
	}
	body, _ := json.Marshal(createReq)

	resp, err := doRequestWithAuth("POST", server.URL+"/api/dm-users", body)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var dmUser map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&dmUser)
	require.NoError(t, err)
	assert.Equal(t, "E2E Test User", dmUser["name"])
	assert.Equal(t, uniqueEmail, dmUser["email"])

	// IDはUUIDv7形式のstring (32文字)
	dmUserID := dmUser["id"].(string)
	assert.Len(t, dmUserID, 32, "ID should be 32 characters (UUIDv7 format)")

	// クリーンアップ: 作成したユーザーを削除
	t.Cleanup(func() {
		_, _ = doRequestWithAuth("DELETE", server.URL+fmt.Sprintf("/api/dm-users/%s", dmUserID), nil)
	})

	// Retrieve dm_user
	resp, err = doRequestWithAuth("GET", server.URL+fmt.Sprintf("/api/dm-users/%s", dmUserID), nil)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Debug: print response if not OK
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Logf("Expected 200 but got %d, body: %s, dmUserID: %s", resp.StatusCode, string(body), dmUserID)
		// Re-create reader for subsequent decode
		resp.Body = io.NopCloser(bytes.NewBuffer(body))
	}

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var retrieved map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&retrieved)
	require.NoError(t, err)
	assert.Equal(t, "E2E Test User", retrieved["name"])
	assert.Equal(t, uniqueEmail, retrieved["email"])
}

func TestDmUserAPI_UpdateAndDelete(t *testing.T) {
	server := setupTestServer(t)
	defer server.Close()

	// ユニークなメールアドレスを生成
	uniqueID, err := idgen.GenerateUUIDv7()
	require.NoError(t, err)
	originalEmail := fmt.Sprintf("original-%s@example.com", uniqueID)
	updatedEmail := fmt.Sprintf("updated-%s@example.com", uniqueID)

	// Create dm_user
	createReq := map[string]string{
		"name":  "Original Name",
		"email": originalEmail,
	}
	body, _ := json.Marshal(createReq)
	resp, err := doRequestWithAuth("POST", server.URL+"/api/dm-users", body)
	require.NoError(t, err)
	defer resp.Body.Close()

	var dmUser map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&dmUser)
	dmUserID := dmUser["id"].(string)

	// Update dm_user
	updateReq := map[string]string{
		"name":  "Updated Name",
		"email": updatedEmail,
	}
	body, _ = json.Marshal(updateReq)
	resp, err = doRequestWithAuth("PUT", server.URL+fmt.Sprintf("/api/dm-users/%s", dmUserID), body)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var updated map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&updated)
	assert.Equal(t, "Updated Name", updated["name"])

	// Delete dm_user
	resp, err = doRequestWithAuth("DELETE", server.URL+fmt.Sprintf("/api/dm-users/%s", dmUserID), nil)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)

	// Verify deletion
	resp, err = doRequestWithAuth("GET", server.URL+fmt.Sprintf("/api/dm-users/%s", dmUserID), nil)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestDmPostAPI_CompleteFlow(t *testing.T) {
	server := setupTestServer(t)
	defer server.Close()

	// ユニークなメールアドレスを生成
	uniqueID, err := idgen.GenerateUUIDv7()
	require.NoError(t, err)
	uniqueEmail := fmt.Sprintf("posttest-%s@example.com", uniqueID)

	// Create dm_user first
	dmUserReq := map[string]string{
		"name":  "Post Test User",
		"email": uniqueEmail,
	}
	body, _ := json.Marshal(dmUserReq)
	resp, err := doRequestWithAuth("POST", server.URL+"/api/dm-users", body)
	require.NoError(t, err)
	defer resp.Body.Close()

	var dmUser map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&dmUser)
	dmUserID := dmUser["id"].(string)

	// クリーンアップ: 作成したユーザーを削除（ポストより後に削除）
	t.Cleanup(func() {
		_, _ = doRequestWithAuth("DELETE", server.URL+fmt.Sprintf("/api/dm-users/%s", dmUserID), nil)
	})

	// Create dm_post
	dmPostReq := map[string]interface{}{
		"user_id": dmUserID,
		"title":   "Test Post",
		"content": "Test content",
	}
	body, _ = json.Marshal(dmPostReq)
	resp, err = doRequestWithAuth("POST", server.URL+"/api/dm-posts", body)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var dmPost map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&dmPost)
	dmPostID := dmPost["id"].(string)
	assert.Equal(t, "Test Post", dmPost["title"])

	// クリーンアップ: 作成したポストを削除（ユーザーより先に削除）
	t.Cleanup(func() {
		_, _ = doRequestWithAuth("DELETE", server.URL+fmt.Sprintf("/api/dm-posts/%s?user_id=%s", dmPostID, dmUserID), nil)
	})

	// Get dm_post
	resp, err = doRequestWithAuth("GET", server.URL+fmt.Sprintf("/api/dm-posts/%s?user_id=%s", dmPostID, dmUserID), nil)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Get dm_user_posts (JOIN)
	resp, err = doRequestWithAuth("GET", server.URL+"/api/dm-user-posts", nil)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var dmUserPosts []map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&dmUserPosts)
	assert.Greater(t, len(dmUserPosts), 0)

	// Find our dm_post in the results
	found := false
	for _, up := range dmUserPosts {
		upPostID := up["post_id"].(string)
		if upPostID == dmPostID {
			upUserID := up["user_id"].(string)
			assert.Equal(t, dmUserID, upUserID)
			assert.Equal(t, "Post Test User", up["user_name"])
			assert.Equal(t, "Test Post", up["post_title"])
			found = true
			break
		}
	}
	assert.True(t, found, "Should find our dm_post in dm_user-posts results")
}
