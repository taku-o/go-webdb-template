package handler

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/taku-o/go-webdb-template/internal/config"
	"github.com/taku-o/go-webdb-template/internal/storage"
	"github.com/tus/tusd/v2/pkg/handler"
)

// UploadHandler はファイルアップロードAPIのハンドラー
type UploadHandler struct {
	tusdHandler *handler.Handler
	config      *config.UploadConfig
}

// NewUploadHandler は新しいUploadHandlerを作成する
func NewUploadHandler(cfg *config.UploadConfig) (*UploadHandler, error) {
	// ストレージファクトリーを使用してストレージを取得
	dataStore, err := storage.NewFileStore(cfg)
	if err != nil {
		return nil, err
	}

	// ストアコンポーザーを作成
	composer := handler.NewStoreComposer()
	composer.UseCore(dataStore)

	// TUSハンドラーを作成
	tusdHandler, err := handler.NewHandler(handler.Config{
		BasePath:      cfg.BasePath,
		StoreComposer: composer,
	})
	if err != nil {
		return nil, err
	}

	return &UploadHandler{
		tusdHandler: tusdHandler,
		config:      cfg,
	}, nil
}


// UploadCompleteCallback はアップロード完了時に呼び出されるコールバック関数の型
type UploadCompleteCallback func(fileID string, filePath string, fileSize int64)

// NewUploadHandlerWithHook はフック付きの新しいUploadHandlerを作成する
func NewUploadHandlerWithHook(cfg *config.UploadConfig, onComplete UploadCompleteCallback) (*UploadHandler, error) {
	// ストレージファクトリーを使用してストレージを取得
	dataStore, err := storage.NewFileStore(cfg)
	if err != nil {
		return nil, err
	}

	// ストアコンポーザーを作成
	composer := handler.NewStoreComposer()
	composer.UseCore(dataStore)

	// TUSハンドラーを作成
	tusdHandler, err := handler.NewHandler(handler.Config{
		BasePath:              cfg.BasePath,
		StoreComposer:         composer,
		NotifyCompleteUploads: true, // アップロード完了通知を有効化
	})
	if err != nil {
		return nil, err
	}

	// アップロード完了時のフック処理を開始（ゴルーチン）
	if onComplete != nil {
		go func() {
			for event := range tusdHandler.CompleteUploads {
				info := event.Upload
				onComplete(info.ID, info.Storage["Path"], info.Size)
			}
		}()
	}

	return &UploadHandler{
		tusdHandler: tusdHandler,
		config:      cfg,
	}, nil
}

// GetHandler はtusdのHTTPハンドラーを返す
func (h *UploadHandler) GetHandler() http.Handler {
	return h.tusdHandler
}

// GetConfig はアップロード設定を返す
func (h *UploadHandler) GetConfig() *config.UploadConfig {
	return h.config
}


// NewUploadValidationMiddleware はファイルアップロードの検証ミドルウェアを作成する
// TUS POSTリクエスト時にファイルサイズと拡張子を検証する
func NewUploadValidationMiddleware(cfg *config.UploadConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// POST以外のリクエストはスキップ
			if c.Request().Method != http.MethodPost {
				return next(c)
			}

			// ファイルサイズの検証
			uploadLength := c.Request().Header.Get("Upload-Length")
			if uploadLength != "" {
				size, err := strconv.ParseInt(uploadLength, 10, 64)
				if err != nil {
					return c.JSON(http.StatusBadRequest, map[string]string{
						"error": "Invalid Upload-Length header",
					})
				}

				if size > cfg.MaxFileSize {
					return c.JSON(http.StatusRequestEntityTooLarge, map[string]string{
						"error": fmt.Sprintf("File size exceeds maximum allowed size of %d bytes", cfg.MaxFileSize),
					})
				}
			}

			// ファイル拡張子の検証
			uploadMetadata := c.Request().Header.Get("Upload-Metadata")
			if uploadMetadata != "" {
				filename := parseFilenameFromMetadata(uploadMetadata)
				if filename != "" {
					ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(filename), "."))
					if !isAllowedExtension(ext, cfg.AllowedExtensions) {
						return c.JSON(http.StatusBadRequest, map[string]string{
							"error": fmt.Sprintf("File extension '%s' is not allowed. Allowed extensions: %v", ext, cfg.AllowedExtensions),
						})
					}
				}
			}

			return next(c)
		}
	}
}

// parseFilenameFromMetadata はUpload-MetadataヘッダーからBase64エンコードされたファイル名を取得する
func parseFilenameFromMetadata(metadata string) string {
	// Upload-Metadataは "key value,key value" 形式
	// valueはBase64エンコードされている
	pairs := strings.Split(metadata, ",")
	for _, pair := range pairs {
		pair = strings.TrimSpace(pair)
		parts := strings.SplitN(pair, " ", 2)
		if len(parts) == 2 && parts[0] == "filename" {
			decoded, err := base64.StdEncoding.DecodeString(parts[1])
			if err != nil {
				return ""
			}
			return string(decoded)
		}
	}
	return ""
}

// isAllowedExtension は拡張子が許可リストに含まれているかチェックする
func isAllowedExtension(ext string, allowedExtensions []string) bool {
	for _, allowed := range allowedExtensions {
		if strings.EqualFold(ext, allowed) {
			return true
		}
	}
	return false
}
