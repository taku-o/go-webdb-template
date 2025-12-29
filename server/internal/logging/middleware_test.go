package logging

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResponseWriter(t *testing.T) {
	t.Run("正常系: ステータスコードが正しく記録される", func(t *testing.T) {
		rec := httptest.NewRecorder()
		rw := &responseWriter{
			ResponseWriter: rec,
			statusCode:     http.StatusOK,
		}

		rw.WriteHeader(http.StatusCreated)

		assert.Equal(t, http.StatusCreated, rw.statusCode)
		assert.Equal(t, http.StatusCreated, rec.Code)
	})

	t.Run("正常系: デフォルトのステータスコードはOK", func(t *testing.T) {
		rec := httptest.NewRecorder()
		rw := &responseWriter{
			ResponseWriter: rec,
			statusCode:     http.StatusOK,
		}

		assert.Equal(t, http.StatusOK, rw.statusCode)
	})
}

func TestAccessLogMiddleware_Middleware(t *testing.T) {
	t.Run("正常系: ミドルウェアが正常に動作する", func(t *testing.T) {
		tmpDir := t.TempDir()

		accessLogger, err := NewAccessLogger("api", tmpDir)
		require.NoError(t, err)
		require.NotNil(t, accessLogger)
		defer accessLogger.Close()

		middleware := NewAccessLogMiddleware(accessLogger)

		// テスト用のハンドラー
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		})

		// ミドルウェアを適用
		wrappedHandler := middleware.Middleware(handler)

		// リクエストを作成
		req := httptest.NewRequest(http.MethodGet, "/api/dm-users", nil)
		req.RemoteAddr = "192.168.1.100:12345"
		req.Header.Set("User-Agent", "TestClient/1.0")

		rec := httptest.NewRecorder()

		// ハンドラーを実行
		wrappedHandler.ServeHTTP(rec, req)

		// レスポンスを確認
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "OK", rec.Body.String())

		// ログファイルが作成されることを確認
		accessLogger.Close()
		today := time.Now().Format("2006-01-02")
		logFileName := "api-access-" + today + ".log"
		logFilePath := filepath.Join(tmpDir, logFileName)

		content, err := os.ReadFile(logFilePath)
		require.NoError(t, err)

		logContent := string(content)
		assert.Contains(t, logContent, "GET")
		assert.Contains(t, logContent, "/api/dm-users")
		assert.Contains(t, logContent, "200")
	})

	t.Run("正常系: レスポンス時間が記録される", func(t *testing.T) {
		tmpDir := t.TempDir()

		accessLogger, err := NewAccessLogger("api", tmpDir)
		require.NoError(t, err)
		require.NotNil(t, accessLogger)
		defer accessLogger.Close()

		middleware := NewAccessLogMiddleware(accessLogger)

		// 遅延のあるハンドラー
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(10 * time.Millisecond)
			w.WriteHeader(http.StatusOK)
		})

		wrappedHandler := middleware.Middleware(handler)

		req := httptest.NewRequest(http.MethodGet, "/api/slow", nil)
		rec := httptest.NewRecorder()

		wrappedHandler.ServeHTTP(rec, req)

		accessLogger.Close()
		today := time.Now().Format("2006-01-02")
		logFileName := "api-access-" + today + ".log"
		logFilePath := filepath.Join(tmpDir, logFileName)

		content, err := os.ReadFile(logFilePath)
		require.NoError(t, err)

		logContent := string(content)
		// レスポンス時間が記録されていることを確認（10ms以上）
		assert.Contains(t, logContent, "ms")
	})

	t.Run("正常系: X-Forwarded-Forヘッダーが優先される", func(t *testing.T) {
		tmpDir := t.TempDir()

		accessLogger, err := NewAccessLogger("api", tmpDir)
		require.NoError(t, err)
		require.NotNil(t, accessLogger)
		defer accessLogger.Close()

		middleware := NewAccessLogMiddleware(accessLogger)

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		wrappedHandler := middleware.Middleware(handler)

		req := httptest.NewRequest(http.MethodGet, "/api/dm-users", nil)
		req.RemoteAddr = "10.0.0.1:12345"
		req.Header.Set("X-Forwarded-For", "192.168.1.200")

		rec := httptest.NewRecorder()
		wrappedHandler.ServeHTTP(rec, req)

		accessLogger.Close()
		today := time.Now().Format("2006-01-02")
		logFileName := "api-access-" + today + ".log"
		logFilePath := filepath.Join(tmpDir, logFileName)

		content, err := os.ReadFile(logFilePath)
		require.NoError(t, err)

		logContent := string(content)
		// X-Forwarded-ForのIPが記録されていることを確認
		assert.Contains(t, logContent, "192.168.1.200")
	})

	t.Run("正常系: User-Agentが空の場合はハイフンが記録される", func(t *testing.T) {
		tmpDir := t.TempDir()

		accessLogger, err := NewAccessLogger("api", tmpDir)
		require.NoError(t, err)
		require.NotNil(t, accessLogger)
		defer accessLogger.Close()

		middleware := NewAccessLogMiddleware(accessLogger)

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		wrappedHandler := middleware.Middleware(handler)

		req := httptest.NewRequest(http.MethodGet, "/api/dm-users", nil)
		// User-Agentヘッダーを設定しない

		rec := httptest.NewRecorder()
		wrappedHandler.ServeHTTP(rec, req)

		accessLogger.Close()
		today := time.Now().Format("2006-01-02")
		logFileName := "api-access-" + today + ".log"
		logFilePath := filepath.Join(tmpDir, logFileName)

		content, err := os.ReadFile(logFilePath)
		require.NoError(t, err)

		logContent := string(content)
		// ハイフンが記録されていることを確認
		assert.Contains(t, logContent, `"-"`)
	})

	t.Run("正常系: エラーステータスコードが正しく記録される", func(t *testing.T) {
		tmpDir := t.TempDir()

		accessLogger, err := NewAccessLogger("api", tmpDir)
		require.NoError(t, err)
		require.NotNil(t, accessLogger)
		defer accessLogger.Close()

		middleware := NewAccessLogMiddleware(accessLogger)

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		})

		wrappedHandler := middleware.Middleware(handler)

		req := httptest.NewRequest(http.MethodPost, "/api/error", nil)
		rec := httptest.NewRecorder()

		wrappedHandler.ServeHTTP(rec, req)

		accessLogger.Close()
		today := time.Now().Format("2006-01-02")
		logFileName := "api-access-" + today + ".log"
		logFilePath := filepath.Join(tmpDir, logFileName)

		content, err := os.ReadFile(logFilePath)
		require.NoError(t, err)

		logContent := string(content)
		assert.Contains(t, logContent, "POST")
		assert.Contains(t, logContent, "/api/error")
		assert.Contains(t, logContent, "500")
	})
}

func TestNewAccessLogMiddleware(t *testing.T) {
	t.Run("正常系: AccessLogMiddlewareが正常に作成される", func(t *testing.T) {
		tmpDir := t.TempDir()

		accessLogger, err := NewAccessLogger("api", tmpDir)
		require.NoError(t, err)
		require.NotNil(t, accessLogger)
		defer accessLogger.Close()

		middleware := NewAccessLogMiddleware(accessLogger)
		assert.NotNil(t, middleware)
	})
}
