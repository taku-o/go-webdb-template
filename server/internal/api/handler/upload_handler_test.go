package handler

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/taku-o/go-webdb-template/internal/config"
)

func TestNewUploadHandler(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir := t.TempDir()
	uploadPath := filepath.Join(tempDir, "uploads")

	cfg := &config.UploadConfig{
		BasePath:          "/api/upload/dm_movie",
		MaxFileSize:       2147483648,
		AllowedExtensions: []string{"mp4"},
		Storage: config.StorageConfig{
			Type: "local",
			Local: config.LocalStorageConfig{
				Path: uploadPath,
			},
		},
	}

	handler, err := NewUploadHandler(cfg)
	if err != nil {
		t.Fatalf("NewUploadHandler failed: %v", err)
	}

	if handler == nil {
		t.Fatal("expected handler to be non-nil")
	}
}

func TestUploadHandler_GetConfig(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir := t.TempDir()
	uploadPath := filepath.Join(tempDir, "uploads")

	cfg := &config.UploadConfig{
		BasePath:          "/api/upload/dm_movie",
		MaxFileSize:       2147483648,
		AllowedExtensions: []string{"mp4"},
		Storage: config.StorageConfig{
			Type: "local",
			Local: config.LocalStorageConfig{
				Path: uploadPath,
			},
		},
	}

	handler, err := NewUploadHandler(cfg)
	if err != nil {
		t.Fatalf("NewUploadHandler failed: %v", err)
	}

	// GetConfigを呼び出し
	handlerCfg := handler.GetConfig()

	if handlerCfg.BasePath != "/api/upload/dm_movie" {
		t.Errorf("expected BasePath '/api/upload/dm_movie', got %s", handlerCfg.BasePath)
	}
	if handlerCfg.MaxFileSize != 2147483648 {
		t.Errorf("expected MaxFileSize 2147483648, got %d", handlerCfg.MaxFileSize)
	}
}

func TestUploadHandler_GetHandler(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir := t.TempDir()
	uploadPath := filepath.Join(tempDir, "uploads")

	cfg := &config.UploadConfig{
		BasePath:          "/api/upload/dm_movie",
		MaxFileSize:       2147483648,
		AllowedExtensions: []string{"mp4"},
		Storage: config.StorageConfig{
			Type: "local",
			Local: config.LocalStorageConfig{
				Path: uploadPath,
			},
		},
	}

	handler, err := NewUploadHandler(cfg)
	if err != nil {
		t.Fatalf("NewUploadHandler failed: %v", err)
	}

	// GetHandlerを呼び出し
	httpHandler := handler.GetHandler()

	if httpHandler == nil {
		t.Fatal("expected httpHandler to be non-nil")
	}
}

// TestNewUploadValidationMiddleware はUploadValidationMiddlewareの作成を確認
func TestNewUploadValidationMiddleware(t *testing.T) {
	cfg := &config.UploadConfig{
		MaxFileSize:       1024 * 1024, // 1MB
		AllowedExtensions: []string{"mp4"},
	}

	middleware := NewUploadValidationMiddleware(cfg)
	require.NotNil(t, middleware)
}

// TestUploadValidationMiddleware_FileSizeExceeded はファイルサイズ超過時に413を返すことを確認
func TestUploadValidationMiddleware_FileSizeExceeded(t *testing.T) {
	cfg := &config.UploadConfig{
		MaxFileSize:       1024 * 1024, // 1MB
		AllowedExtensions: []string{"mp4"},
	}

	middleware := NewUploadValidationMiddleware(cfg)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/upload/dm_movie", nil)
	// Upload-Lengthヘッダーで2MBを指定（制限超過）
	req.Header.Set("Upload-Length", "2097152")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := middleware(func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	err := handler(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusRequestEntityTooLarge, rec.Code)
}

// TestUploadValidationMiddleware_FileSizeValid は有効なファイルサイズで成功することを確認
func TestUploadValidationMiddleware_FileSizeValid(t *testing.T) {
	cfg := &config.UploadConfig{
		MaxFileSize:       1024 * 1024, // 1MB
		AllowedExtensions: []string{"mp4"},
	}

	middleware := NewUploadValidationMiddleware(cfg)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/upload/dm_movie", nil)
	// Upload-Lengthヘッダーで512KBを指定（制限内）
	req.Header.Set("Upload-Length", "524288")
	// Upload-Metadataヘッダーにファイル名を設定（Base64エンコード）
	// filename test.mp4 -> dGVzdC5tcDQ=
	req.Header.Set("Upload-Metadata", "filename dGVzdC5tcDQ=")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := middleware(func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	err := handler(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

// TestUploadValidationMiddleware_InvalidExtension は無効な拡張子で400を返すことを確認
func TestUploadValidationMiddleware_InvalidExtension(t *testing.T) {
	cfg := &config.UploadConfig{
		MaxFileSize:       1024 * 1024, // 1MB
		AllowedExtensions: []string{"mp4"},
	}

	middleware := NewUploadValidationMiddleware(cfg)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/upload/dm_movie", nil)
	req.Header.Set("Upload-Length", "524288")
	// Upload-Metadataヘッダーに無効な拡張子のファイル名を設定
	// filename test.txt -> dGVzdC50eHQ=
	req.Header.Set("Upload-Metadata", "filename dGVzdC50eHQ=")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := middleware(func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	err := handler(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// TestUploadValidationMiddleware_ValidExtension は有効な拡張子で成功することを確認
func TestUploadValidationMiddleware_ValidExtension(t *testing.T) {
	cfg := &config.UploadConfig{
		MaxFileSize:       1024 * 1024, // 1MB
		AllowedExtensions: []string{"mp4", "mov"},
	}

	middleware := NewUploadValidationMiddleware(cfg)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/upload/dm_movie", nil)
	req.Header.Set("Upload-Length", "524288")
	// filename video.mov -> dmlkZW8ubW92
	req.Header.Set("Upload-Metadata", "filename dmlkZW8ubW92")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := middleware(func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	err := handler(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

// TestUploadValidationMiddleware_NonPostRequest はPOST以外のリクエストでスキップすることを確認
func TestUploadValidationMiddleware_NonPostRequest(t *testing.T) {
	cfg := &config.UploadConfig{
		MaxFileSize:       1024, // 1KB
		AllowedExtensions: []string{"mp4"},
	}

	middleware := NewUploadValidationMiddleware(cfg)

	methods := []string{
		http.MethodOptions,
		http.MethodPatch,
		http.MethodHead,
		http.MethodDelete,
	}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(method, "/api/upload/dm_movie", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			handler := middleware(func(c echo.Context) error {
				return c.String(http.StatusOK, "OK")
			})

			err := handler(c)
			require.NoError(t, err)
			assert.Equal(t, http.StatusOK, rec.Code, "method %s should skip validation", method)
		})
	}
}

// TestNewUploadHandlerWithHook はフック付きUploadHandlerの作成を確認
func TestNewUploadHandlerWithHook(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir := t.TempDir()
	uploadPath := filepath.Join(tempDir, "uploads")

	cfg := &config.UploadConfig{
		BasePath:          "/api/upload/dm_movie",
		MaxFileSize:       2147483648,
		AllowedExtensions: []string{"mp4"},
		Storage: config.StorageConfig{
			Type: "local",
			Local: config.LocalStorageConfig{
				Path: uploadPath,
			},
		},
	}

	// カスタムフックを設定してハンドラーを作成
	hookCalled := false
	onComplete := func(fileID string, filePath string, fileSize int64) {
		hookCalled = true
	}

	handler, err := NewUploadHandlerWithHook(cfg, onComplete)
	require.NoError(t, err)
	require.NotNil(t, handler)

	// フック関数が設定されたことは、内部状態として確認（実際のフック実行はE2Eテストで確認）
	// ここでは作成されたハンドラーがnilでないことだけ確認
	assert.NotNil(t, handler.GetHandler())
	// hookCalledはここでは変更されない（アップロード完了時に呼ばれる）
	assert.False(t, hookCalled)
}
