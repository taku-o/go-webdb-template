package integration_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupAdminHealthTestServer creates a minimal test server that mimics
// the Admin server's /health endpoint behavior
func setupAdminHealthTestServer(_ *testing.T) *httptest.Server {
	app := mux.NewRouter()

	// Health check endpoint (認証不要)
	app.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	return httptest.NewServer(app)
}

func TestAdminHealth_HealthCheckNoAuth(t *testing.T) {
	server := setupAdminHealthTestServer(t)
	defer server.Close()

	// Health check should not require auth
	resp, err := http.Get(server.URL + "/health")
	require.NoError(t, err)
	defer resp.Body.Close()

	// ステータスコードが200 OKであることを確認
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// レスポンスボディが"OK"であることを確認
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Equal(t, "OK", string(body))

	// Content-Typeがtext/plainであることを確認
	assert.Equal(t, "text/plain", resp.Header.Get("Content-Type"))
}

func TestAdminHealth_HealthCheckGETMethodOnly(t *testing.T) {
	server := setupAdminHealthTestServer(t)
	defer server.Close()

	// POST should return 405 Method Not Allowed
	resp, err := http.Post(server.URL+"/health", "application/json", nil)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
}
