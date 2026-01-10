package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestHealthEndpoint(t *testing.T) {
	// Gorilla Mux Routerを作成
	app := mux.NewRouter()

	// Health check endpointを登録
	app.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	// テストリクエストを作成
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	// リクエストを処理
	app.ServeHTTP(rec, req)

	// アサーション
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "OK", rec.Body.String())
	assert.Equal(t, "text/plain", rec.Header().Get("Content-Type"))
}
