package storage

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewLocalFileStore(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir := t.TempDir()
	uploadPath := filepath.Join(tempDir, "uploads")

	// LocalFileStoreを作成
	store, err := NewLocalFileStore(uploadPath)
	if err != nil {
		t.Fatalf("NewLocalFileStore failed: %v", err)
	}

	// storeがnilでないことを確認
	if store == nil {
		t.Fatal("expected store to be non-nil")
	}

	// ディレクトリが作成されていることを確認
	if _, err := os.Stat(uploadPath); os.IsNotExist(err) {
		t.Errorf("expected directory %s to exist", uploadPath)
	}
}

func TestNewLocalFileStore_ExistingDirectory(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir := t.TempDir()
	uploadPath := filepath.Join(tempDir, "existing_uploads")

	// 事前にディレクトリを作成
	if err := os.MkdirAll(uploadPath, 0755); err != nil {
		t.Fatalf("failed to create directory: %v", err)
	}

	// LocalFileStoreを作成（既存ディレクトリでエラーにならないことを確認）
	store, err := NewLocalFileStore(uploadPath)
	if err != nil {
		t.Fatalf("NewLocalFileStore failed: %v", err)
	}

	if store == nil {
		t.Fatal("expected store to be non-nil")
	}
}

func TestLocalFileStore_GetStore(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir := t.TempDir()
	uploadPath := filepath.Join(tempDir, "uploads")

	// LocalFileStoreを作成
	store, err := NewLocalFileStore(uploadPath)
	if err != nil {
		t.Fatalf("NewLocalFileStore failed: %v", err)
	}

	// GetStoreを呼び出し
	dataStore := store.GetStore()

	// dataStoreがnilでないことを確認
	if dataStore == nil {
		t.Fatal("expected dataStore to be non-nil")
	}
}

func TestLocalFileStore_GetPath(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir := t.TempDir()
	uploadPath := filepath.Join(tempDir, "uploads")

	// LocalFileStoreを作成
	store, err := NewLocalFileStore(uploadPath)
	if err != nil {
		t.Fatalf("NewLocalFileStore failed: %v", err)
	}

	// GetPathを呼び出し
	path := store.GetPath()

	// パスが正しいことを確認
	if path != uploadPath {
		t.Errorf("expected path %s, got %s", uploadPath, path)
	}
}
