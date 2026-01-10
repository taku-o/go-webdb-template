package testutil

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/gofrs/flock"
)

// getProjectRoot returns the project root directory (where go.mod exists)
func getProjectRoot() (string, error) {
	// 現在のファイルのディレクトリから開始
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("failed to get current file path")
	}

	dir := filepath.Dir(filename)

	// go.modを探す（最大5階層まで上に遡る）
	for i := 0; i < 5; i++ {
		goModPath := filepath.Join(dir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return "", fmt.Errorf("go.mod not found")
}

// AcquireTestLock acquires a file lock for database tests
// Returns a flock.Flock object that should be unlocked with defer fileLock.Unlock()
func AcquireTestLock(t *testing.T) (*flock.Flock, error) {
	// プロジェクトルートを取得
	projectRoot, err := getProjectRoot()
	if err != nil {
		return nil, fmt.Errorf("failed to get project root: %w", err)
	}

	// ロックファイルのディレクトリを構築
	lockDir := filepath.Join(projectRoot, ".test-lock")

	// ディレクトリが存在しない場合は作成
	if err := os.MkdirAll(lockDir, 0755); err != nil {
		return nil, fmt.Errorf("ロックディレクトリの作成に失敗しました (%s): %w", lockDir, err)
	}

	// ロックファイルのパスを構築
	lockPath := filepath.Join(lockDir, "test-db.lock")

	// ロックオブジェクトを作成
	fileLock := flock.New(lockPath)

	// タイムアウト付きコンテキストを作成
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// ロック取得を試行
	locked, err := fileLock.TryLockContext(ctx, 100*time.Millisecond)
	if err != nil {
		// エラーハンドリング
		if err == context.DeadlineExceeded {
			return nil, fmt.Errorf("%sのロックが取れなかったのでタイムアウトしました", lockPath)
		}
		return nil, fmt.Errorf("ロックファイルの取得に失敗しました (%s): %w", lockPath, err)
	}

	if !locked {
		return nil, fmt.Errorf("%sのロックが取れなかったのでタイムアウトしました", lockPath)
	}

	return fileLock, nil
}
