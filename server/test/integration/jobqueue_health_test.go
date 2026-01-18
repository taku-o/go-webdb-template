package integration_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupJobQueueHealthTestServer creates a minimal test server that mimics
// the JobQueue server's /health endpoint behavior
func setupJobQueueHealthTestServer(_ *testing.T) *httptest.Server {
	mux := http.NewServeMux()

	// Health check endpoint (認証不要)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	return httptest.NewServer(mux)
}

func TestJobQueueHealth_HealthCheckNoAuth(t *testing.T) {
	server := setupJobQueueHealthTestServer(t)
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
