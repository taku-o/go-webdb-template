package e2e_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/example/go-db-prj-sample/internal/api/handler"
	"github.com/example/go-db-prj-sample/internal/api/router"
	"github.com/example/go-db-prj-sample/internal/repository"
	"github.com/example/go-db-prj-sample/internal/service"
	"github.com/example/go-db-prj-sample/test/testutil"
)

func setupTestServer(t *testing.T) *httptest.Server {
	// Setup test database
	dbManager := testutil.SetupTestShards(t, 2)
	t.Cleanup(func() {
		testutil.CleanupTestDB(dbManager)
	})

	// Initialize layers
	userRepo := repository.NewUserRepository(dbManager)
	userService := service.NewUserService(userRepo, dbManager)
	userHandler := handler.NewUserHandler(userService)

	postRepo := repository.NewPostRepository(dbManager, userRepo)
	postService := service.NewPostService(postRepo, dbManager)
	postHandler := handler.NewPostHandler(postService)

	// Setup router
	r := router.SetupRouter(userHandler, postHandler)

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

	resp, err := http.Post(
		server.URL+"/api/users",
		"application/json",
		bytes.NewBuffer(body),
	)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var user map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&user)
	require.NoError(t, err)
	assert.Equal(t, "E2E Test User", user["name"])
	assert.Equal(t, "e2e@example.com", user["email"])

	userID := int(user["id"].(float64))

	// Retrieve user
	resp, err = http.Get(server.URL + fmt.Sprintf("/api/users/%d", userID))
	require.NoError(t, err)
	defer resp.Body.Close()
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
	resp, err := http.Post(server.URL+"/api/users", "application/json", bytes.NewBuffer(body))
	require.NoError(t, err)
	defer resp.Body.Close()

	var user map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&user)
	userID := int(user["id"].(float64))

	// Update user
	updateReq := map[string]string{
		"name":  "Updated Name",
		"email": "updated@example.com",
	}
	body, _ = json.Marshal(updateReq)
	req, _ := http.NewRequest(
		"PUT",
		server.URL+fmt.Sprintf("/api/users/%d", userID),
		bytes.NewBuffer(body),
	)
	req.Header.Set("Content-Type", "application/json")
	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var updated map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&updated)
	assert.Equal(t, "Updated Name", updated["name"])

	// Delete user
	req, _ = http.NewRequest("DELETE", server.URL+fmt.Sprintf("/api/users/%d", userID), nil)
	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)

	// Verify deletion
	resp, err = http.Get(server.URL + fmt.Sprintf("/api/users/%d", userID))
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
	resp, err := http.Post(server.URL+"/api/users", "application/json", bytes.NewBuffer(body))
	require.NoError(t, err)
	defer resp.Body.Close()

	var user map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&user)
	userID := int(user["id"].(float64))

	// Create post
	postReq := map[string]interface{}{
		"user_id": userID,
		"title":   "Test Post",
		"content": "Test content",
	}
	body, _ = json.Marshal(postReq)
	resp, err = http.Post(server.URL+"/api/posts", "application/json", bytes.NewBuffer(body))
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var post map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&post)
	postID := int(post["id"].(float64))
	assert.Equal(t, "Test Post", post["title"])

	// Get post
	resp, err = http.Get(server.URL + fmt.Sprintf("/api/posts/%d?user_id=%d", postID, userID))
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Get user posts (JOIN)
	resp, err = http.Get(server.URL + "/api/user-posts")
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var userPosts []map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&userPosts)
	assert.Greater(t, len(userPosts), 0)

	// Find our post in the results
	found := false
	for _, up := range userPosts {
		if int(up["post_id"].(float64)) == postID {
			assert.Equal(t, float64(userID), up["user_id"].(float64))
			assert.Equal(t, "Post Test User", up["user_name"])
			assert.Equal(t, "Test Post", up["post_title"])
			found = true
			break
		}
	}
	assert.True(t, found, "Should find our post in user-posts results")
}
