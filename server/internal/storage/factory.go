package storage

import (
	"fmt"

	"github.com/taku-o/go-webdb-template/internal/config"
	"github.com/tus/tusd/v2/pkg/handler"
)

// NewFileStore は設定に基づいて適切なストレージを作成する
func NewFileStore(cfg *config.UploadConfig) (handler.DataStore, error) {
	switch cfg.Storage.Type {
	case "local":
		store, err := NewLocalFileStore(cfg.Storage.Local.Path)
		if err != nil {
			return nil, fmt.Errorf("failed to create local file store: %w", err)
		}
		return store.GetStore(), nil
	case "s3":
		store, err := NewS3FileStore(cfg.Storage.S3.Bucket, cfg.Storage.S3.Region)
		if err != nil {
			return nil, fmt.Errorf("failed to create S3 file store: %w", err)
		}
		return store.GetStore(), nil
	default:
		return nil, fmt.Errorf("unsupported storage type: %s", cfg.Storage.Type)
	}
}
