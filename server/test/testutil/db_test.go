package testutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadTestConfig(t *testing.T) {
	cfg, err := LoadTestConfig()
	require.NoError(t, err)

	// マスターデータベース名を確認
	assert.Equal(t, "webdb_master_test", cfg.Database.Groups.Master[0].Name)

	// シャーディングデータベース名を確認
	assert.Equal(t, "webdb_sharding_1_test", cfg.Database.Groups.Sharding.Databases[0].Name)

	// その他の設定を確認（ホスト、ポート等）
	assert.Equal(t, "localhost", cfg.Database.Groups.Master[0].Host)
	assert.Equal(t, 5432, cfg.Database.Groups.Master[0].Port)
	assert.Equal(t, "webdb", cfg.Database.Groups.Master[0].User)
	assert.Equal(t, "webdb", cfg.Database.Groups.Master[0].Password)
}

func TestClearTestDatabase(t *testing.T) {
	manager := SetupTestGroupManager(t, 4, 8)
	defer CleanupTestGroupManager(manager)

	// テストデータを挿入
	masterConn, err := manager.GetMasterConnection()
	require.NoError(t, err)
	masterConn.DB.Exec("INSERT INTO dm_news (title, content, created_at, updated_at) VALUES ('test', 'test', NOW(), NOW())")

	// データが存在することを確認
	var count int64
	masterConn.DB.Raw("SELECT COUNT(*) FROM dm_news").Scan(&count)
	assert.Greater(t, count, int64(0))

	// データベースをクリア
	ClearTestDatabase(t, manager)

	// データがクリアされたことを確認
	masterConn.DB.Raw("SELECT COUNT(*) FROM dm_news").Scan(&count)
	assert.Equal(t, int64(0), count)

	// テーブルが存在することを確認（構造が維持される）
	var exists bool
	masterConn.DB.Raw(`
		SELECT EXISTS (
			SELECT FROM information_schema.tables
			WHERE table_schema = 'public'
			AND table_name = 'dm_news'
		)
	`).Scan(&exists)
	assert.True(t, exists)
}
