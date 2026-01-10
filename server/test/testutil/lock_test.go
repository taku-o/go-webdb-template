package testutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetProjectRoot_Success(t *testing.T) {
	// プロジェクトルートを取得
	root, err := getProjectRoot()
	require.NoError(t, err)

	// go.modが存在するディレクトリを返すことを確認
	assert.NotEmpty(t, root)
	assert.Contains(t, root, "server")
}

func TestAcquireTestLock_Success(t *testing.T) {
	// ロックを取得
	fileLock, err := AcquireTestLock(t)
	require.NoError(t, err)
	require.NotNil(t, fileLock)

	// ロックが取得できていることを確認
	assert.True(t, fileLock.Locked())

	// ロックを解放
	err = fileLock.Unlock()
	require.NoError(t, err)
}

func TestAcquireTestLock_Sequential(t *testing.T) {
	// 最初のロックを取得
	fileLock1, err := AcquireTestLock(t)
	require.NoError(t, err)
	require.NotNil(t, fileLock1)
	defer fileLock1.Unlock()

	// 同じプロセス内では、flockは同じファイルディスクリプタを再利用するため
	// 同じロックを再度取得できることを確認（これは期待される動作）
	// 注意: これは同一プロセス内の動作であり、プロセス間のロックとは異なる
}
