package admin_test

import (
	"testing"

	"gorm.io/gorm"

	"github.com/example/go-webdb-template/internal/admin"
	"github.com/example/go-webdb-template/internal/config"
	"github.com/example/go-webdb-template/internal/model"
	"github.com/example/go-webdb-template/test/testutil"
)

// TestAdminConfigIntegration はGoAdmin設定の統合テスト
func TestAdminConfigIntegration(t *testing.T) {
	cfg := &config.Config{
		Admin: config.AdminConfig{
			Port: 8081,
			Auth: config.AuthConfig{
				Username: "admin",
				Password: "password",
			},
			Session: config.SessionConfig{
				Lifetime: 7200,
			},
		},
		Database: config.DatabaseConfig{
			Shards: []config.ShardConfig{
				{
					ID:     1,
					Driver: "sqlite3",
					DSN:    "test.db",
				},
			},
		},
		Logging: config.LoggingConfig{
			Level: "debug",
		},
	}

	// Config構造体の作成
	adminCfg := admin.NewConfig(cfg)
	if adminCfg == nil {
		t.Fatal("NewConfig returned nil")
	}

	// GetAdminPort
	if adminCfg.GetAdminPort() != 8081 {
		t.Errorf("expected port 8081, got %d", adminCfg.GetAdminPort())
	}

	// GetGoAdminConfig
	goadminCfg := adminCfg.GetGoAdminConfig()
	if goadminCfg == nil {
		t.Fatal("GetGoAdminConfig returned nil")
	}

	if goadminCfg.SessionLifeTime != 7200 {
		t.Errorf("expected SessionLifeTime 7200, got %d", goadminCfg.SessionLifeTime)
	}

	if goadminCfg.Debug != true {
		t.Error("expected Debug to be true for debug logging level")
	}
}

// TestAdminTableGenerators はテーブルジェネレータの統合テスト
func TestAdminTableGenerators(t *testing.T) {
	// Generatorsマップの確認
	if admin.Generators == nil {
		t.Fatal("Generators map is nil")
	}

	// usersジェネレータの確認
	if _, ok := admin.Generators["users"]; !ok {
		t.Error("users generator not found in Generators map")
	}

	// postsジェネレータの確認
	if _, ok := admin.Generators["posts"]; !ok {
		t.Error("posts generator not found in Generators map")
	}
}

// TestShardingWithCRUDOperations はシャーディングを使用したCRUD操作の統合テスト
func TestShardingWithCRUDOperations(t *testing.T) {
	manager := testutil.SetupTestGORMShards(t, 2)
	defer testutil.CleanupTestGORMDB(manager)

	// テストデータを挿入
	connections := manager.GetAllGORMConnections()

	// Shard 1にユーザーを作成
	user1 := &model.User{
		ID:    1,
		Name:  "User 1",
		Email: "user1@example.com",
	}
	if err := connections[0].DB.Create(user1).Error; err != nil {
		t.Fatalf("Failed to create user1: %v", err)
	}

	// Shard 2にユーザーを作成
	user2 := &model.User{
		ID:    2,
		Name:  "User 2",
		Email: "user2@example.com",
	}
	if err := connections[1].DB.Create(user2).Error; err != nil {
		t.Fatalf("Failed to create user2: %v", err)
	}

	// 全シャードからユーザーを取得
	users, err := admin.FindUserAcrossShards(manager, func(db *gorm.DB) *gorm.DB {
		return db.Model(&model.User{})
	})
	if err != nil {
		t.Fatalf("FindUserAcrossShards failed: %v", err)
	}

	if len(users) != 2 {
		t.Errorf("expected 2 users, got %d", len(users))
	}

	// ユーザー数をカウント
	count, err := admin.CountUsersAcrossShards(manager)
	if err != nil {
		t.Fatalf("CountUsersAcrossShards failed: %v", err)
	}

	if count != 2 {
		t.Errorf("expected 2 users, got %d", count)
	}

	// 投稿を作成
	post1 := &model.Post{
		ID:      1,
		UserID:  1,
		Title:   "Post 1",
		Content: "Content 1",
	}
	if err := connections[0].DB.Create(post1).Error; err != nil {
		t.Fatalf("Failed to create post1: %v", err)
	}

	post2 := &model.Post{
		ID:      2,
		UserID:  2,
		Title:   "Post 2",
		Content: "Content 2",
	}
	if err := connections[1].DB.Create(post2).Error; err != nil {
		t.Fatalf("Failed to create post2: %v", err)
	}

	// 投稿数をカウント
	postCount, err := admin.CountPostsAcrossShards(manager)
	if err != nil {
		t.Fatalf("CountPostsAcrossShards failed: %v", err)
	}

	if postCount != 2 {
		t.Errorf("expected 2 posts, got %d", postCount)
	}

	// シャード統計情報を取得
	stats, err := admin.GetShardStats(manager)
	if err != nil {
		t.Fatalf("GetShardStats failed: %v", err)
	}

	if len(stats) != 2 {
		t.Errorf("expected 2 shard stats, got %d", len(stats))
	}

	// 各シャードのデータを確認
	var totalUsers, totalPosts int64
	for _, s := range stats {
		totalUsers += s.UserCount
		totalPosts += s.PostCount
	}

	if totalUsers != 2 {
		t.Errorf("expected 2 total users, got %d", totalUsers)
	}

	if totalPosts != 2 {
		t.Errorf("expected 2 total posts, got %d", totalPosts)
	}
}

// TestInsertToShardWithRouting はシャードルーティングを使用した挿入テスト
func TestInsertToShardWithRouting(t *testing.T) {
	manager := testutil.SetupTestGORMShards(t, 2)
	defer testutil.CleanupTestGORMDB(manager)

	// InsertToShardを使用してデータを挿入
	user := &model.User{
		ID:    100,
		Name:  "Routed User",
		Email: "routed@example.com",
	}

	err := admin.InsertToShard(manager, 100, user)
	if err != nil {
		t.Fatalf("InsertToShard failed: %v", err)
	}

	// 挿入されたデータを確認
	users, err := admin.FindUserAcrossShards(manager, func(db *gorm.DB) *gorm.DB {
		return db.Model(&model.User{})
	})
	if err != nil {
		t.Fatalf("FindUserAcrossShards failed: %v", err)
	}

	found := false
	for _, u := range users {
		if u.Name == "Routed User" {
			found = true
			break
		}
	}

	if !found {
		t.Error("Inserted user not found across shards")
	}
}
