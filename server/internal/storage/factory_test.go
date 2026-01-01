package storage

import (
	"path/filepath"
	"testing"

	"github.com/taku-o/go-webdb-template/internal/config"
)

func TestNewFileStore_Local(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir := t.TempDir()
	uploadPath := filepath.Join(tempDir, "uploads")

	cfg := &config.UploadConfig{
		Storage: config.StorageConfig{
			Type: "local",
			Local: config.LocalStorageConfig{
				Path: uploadPath,
			},
		},
	}

	store, err := NewFileStore(cfg)
	if err != nil {
		t.Fatalf("NewFileStore failed: %v", err)
	}

	if store == nil {
		t.Fatal("expected store to be non-nil")
	}
}

func TestNewFileStore_InvalidType(t *testing.T) {
	cfg := &config.UploadConfig{
		Storage: config.StorageConfig{
			Type: "invalid",
		},
	}

	_, err := NewFileStore(cfg)
	if err == nil {
		t.Fatal("expected error for invalid storage type")
	}
}

func TestNewFileStore_S3(t *testing.T) {
	// S3テストはAWS認証情報がない環境ではスキップ
	cfg := &config.UploadConfig{
		Storage: config.StorageConfig{
			Type: "s3",
			S3: config.S3StorageConfig{
				Bucket: "test-bucket",
				Region: "ap-northeast-1",
			},
		},
	}

	store, err := NewFileStore(cfg)
	// AWS認証情報がない場合はエラーになる可能性がある
	if err != nil {
		t.Logf("NewFileStore with S3 returned error (expected in dev environment without AWS credentials): %v", err)
		return
	}

	if store == nil {
		t.Fatal("expected store to be non-nil")
	}
}
