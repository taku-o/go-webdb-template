package integration_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/taku-o/go-webdb-template/internal/api/handler"
	"github.com/taku-o/go-webdb-template/internal/api/router"
	"github.com/taku-o/go-webdb-template/internal/auth"
	"github.com/taku-o/go-webdb-template/internal/repository"
	"github.com/taku-o/go-webdb-template/internal/service"
	"github.com/taku-o/go-webdb-template/test/testutil"
)

func setupAuthTestServer(t *testing.T) *httptest.Server {
	// Setup test database with GroupManager
	groupManager := testutil.SetupTestGroupManager(t, 4, 8)
	t.Cleanup(func() {
		testutil.CleanupTestGroupManager(groupManager)
	})

	// Initialize layers (using GORM repositories)
	dmUserRepo := repository.NewDmUserRepository(groupManager)
	dmUserService := service.NewDmUserService(dmUserRepo)
	dmUserHandler := handler.NewDmUserHandler(dmUserService)

	dmPostRepo := repository.NewDmPostRepository(groupManager)
	dmPostService := service.NewDmPostService(dmPostRepo, dmUserRepo)
	dmPostHandler := handler.NewDmPostHandler(dmPostService)

	// TodayHandler
	todayHandler := handler.NewTodayHandler()

	// Setup router with test config
	cfg := testutil.GetTestConfig()
	r := router.NewRouter(dmUserHandler, dmPostHandler, todayHandler, nil, nil, cfg)

	return httptest.NewServer(r)
}

func TestAPIAuth_ValidToken(t *testing.T) {
	server := setupAuthTestServer(t)
	defer server.Close()

	// Get valid token
	token, err := testutil.GetTestAPIToken()
	require.NoError(t, err)

	// Access API with valid token
	req, err := http.NewRequest("GET", server.URL+"/api/dm-users", nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestAPIAuth_NoToken(t *testing.T) {
	server := setupAuthTestServer(t)
	defer server.Close()

	// Access API without token
	resp, err := http.Get(server.URL + "/api/dm-users")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestAPIAuth_InvalidToken(t *testing.T) {
	server := setupAuthTestServer(t)
	defer server.Close()

	// Access API with invalid token
	req, err := http.NewRequest("GET", server.URL+"/api/dm-users", nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer invalid-token")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestAPIAuth_InvalidVersion(t *testing.T) {
	server := setupAuthTestServer(t)
	defer server.Close()

	// Generate token with invalid version (v1)
	token, err := auth.GeneratePublicAPIKey(testutil.TestSecretKey, "v1", testutil.TestEnv, time.Now().Unix())
	require.NoError(t, err)

	// Access API with invalid version token
	req, err := http.NewRequest("GET", server.URL+"/api/dm-users", nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestAPIAuth_HealthCheckNoAuth(t *testing.T) {
	server := setupAuthTestServer(t)
	defer server.Close()

	// Health check should not require auth
	resp, err := http.Get(server.URL + "/health")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
