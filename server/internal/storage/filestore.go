package storage

import (
	"os"

	"github.com/tus/tusd/v2/pkg/filestore"
	"github.com/tus/tusd/v2/pkg/handler"
)

// LocalFileStore はローカルファイルシステムストア
type LocalFileStore struct {
	path  string
	store *filestore.FileStore
}

// NewLocalFileStore は新しいLocalFileStoreを作成する
// 指定されたパスのディレクトリが存在しない場合は作成する
func NewLocalFileStore(path string) (*LocalFileStore, error) {
	// ディレクトリが存在しない場合は作成
	if err := os.MkdirAll(path, 0755); err != nil {
		return nil, err
	}

	// tusdのFileStoreを作成
	store := filestore.New(path)

	return &LocalFileStore{
		path:  path,
		store: &store,
	}, nil
}

// GetStore はtusdのDataStoreを返す
func (l *LocalFileStore) GetStore() handler.DataStore {
	return l.store
}

// GetPath はストレージパスを返す
func (l *LocalFileStore) GetPath() string {
	return l.path
}
