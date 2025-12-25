package admin

import (
	"testing"

	"gorm.io/gorm"

	"github.com/example/go-webdb-template/internal/model"
	"github.com/example/go-webdb-template/test/testutil"
)

func TestQueryAllShards(t *testing.T) {
	manager := testutil.SetupTestGORMShards(t, 2)
	defer testutil.CleanupTestGORMDB(manager)

	// テストデータを挿入
	connections := manager.GetAllGORMConnections()
	for i, conn := range connections {
		user := &model.User{
			ID:    int64(i + 1),
			Name:  "User " + string(rune('A'+i)),
			Email: "user" + string(rune('a'+i)) + "@example.com",
		}
		if err := conn.DB.Create(user).Error; err != nil {
			t.Fatalf("Failed to create test user: %v", err)
		}
	}

	// 全シャードからデータを取得
	users, err := QueryAllShards[model.User](manager, func(gormDB *gorm.DB) *gorm.DB {
		return gormDB.Model(&model.User{})
	})

	if err != nil {
		t.Fatalf("QueryAllShards failed: %v", err)
	}

	if len(users) != 2 {
		t.Errorf("expected 2 users, got %d", len(users))
	}
}

func TestFindUserAcrossShards(t *testing.T) {
	manager := testutil.SetupTestGORMShards(t, 2)
	defer testutil.CleanupTestGORMDB(manager)

	// テストデータを挿入（Shard 1に挿入）
	connections := manager.GetAllGORMConnections()
	testUser := &model.User{
		ID:    100,
		Name:  "Test User",
		Email: "test@example.com",
	}
	if err := connections[0].DB.Create(testUser).Error; err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// 全シャードから検索
	users, err := FindUserAcrossShards(manager, func(gormDB *gorm.DB) *gorm.DB {
		return gormDB.Model(&model.User{}).Where("id = ?", 100)
	})

	if err != nil {
		t.Fatalf("FindUserAcrossShards failed: %v", err)
	}

	if len(users) != 1 {
		t.Errorf("expected 1 user, got %d", len(users))
	}

	if len(users) > 0 && users[0].Name != "Test User" {
		t.Errorf("expected user name 'Test User', got '%s'", users[0].Name)
	}
}

func TestFindPostAcrossShards(t *testing.T) {
	manager := testutil.SetupTestGORMShards(t, 2)
	defer testutil.CleanupTestGORMDB(manager)

	// テストデータを挿入
	connections := manager.GetAllGORMConnections()
	testPost := &model.Post{
		ID:      1,
		UserID:  100,
		Title:   "Test Post",
		Content: "Test Content",
	}
	if err := connections[0].DB.Create(testPost).Error; err != nil {
		t.Fatalf("Failed to create test post: %v", err)
	}

	// 全シャードから検索
	posts, err := FindPostAcrossShards(manager, func(gormDB *gorm.DB) *gorm.DB {
		return gormDB.Model(&model.Post{}).Where("user_id = ?", 100)
	})

	if err != nil {
		t.Fatalf("FindPostAcrossShards failed: %v", err)
	}

	if len(posts) != 1 {
		t.Errorf("expected 1 post, got %d", len(posts))
	}

	if len(posts) > 0 && posts[0].Title != "Test Post" {
		t.Errorf("expected post title 'Test Post', got '%s'", posts[0].Title)
	}
}

func TestCountUsersAcrossShards(t *testing.T) {
	manager := testutil.SetupTestGORMShards(t, 2)
	defer testutil.CleanupTestGORMDB(manager)

	// テストデータを挿入
	connections := manager.GetAllGORMConnections()
	for i, conn := range connections {
		// 各シャードに2人ずつ
		for j := 0; j < 2; j++ {
			user := &model.User{
				ID:    int64(i*10 + j + 1),
				Name:  "User",
				Email: "user" + string(rune('a'+i*10+j)) + "@example.com",
			}
			if err := conn.DB.Create(user).Error; err != nil {
				t.Fatalf("Failed to create test user: %v", err)
			}
		}
	}

	// 全シャードのユーザー数をカウント
	count, err := CountUsersAcrossShards(manager)

	if err != nil {
		t.Fatalf("CountUsersAcrossShards failed: %v", err)
	}

	if count != 4 {
		t.Errorf("expected 4 users, got %d", count)
	}
}

func TestCountPostsAcrossShards(t *testing.T) {
	manager := testutil.SetupTestGORMShards(t, 2)
	defer testutil.CleanupTestGORMDB(manager)

	// テストデータを挿入
	connections := manager.GetAllGORMConnections()
	for i, conn := range connections {
		// 各シャードに3投稿ずつ
		for j := 0; j < 3; j++ {
			post := &model.Post{
				ID:      int64(i*10 + j + 1),
				UserID:  int64(i + 1),
				Title:   "Post",
				Content: "Content",
			}
			if err := conn.DB.Create(post).Error; err != nil {
				t.Fatalf("Failed to create test post: %v", err)
			}
		}
	}

	// 全シャードの投稿数をカウント
	count, err := CountPostsAcrossShards(manager)

	if err != nil {
		t.Fatalf("CountPostsAcrossShards failed: %v", err)
	}

	if count != 6 {
		t.Errorf("expected 6 posts, got %d", count)
	}
}

func TestGetShardStats(t *testing.T) {
	manager := testutil.SetupTestGORMShards(t, 2)
	defer testutil.CleanupTestGORMDB(manager)

	// テストデータを挿入
	connections := manager.GetAllGORMConnections()
	// Shard 1: 2 users, 3 posts
	for i := 0; i < 2; i++ {
		user := &model.User{
			ID:    int64(i + 1),
			Name:  "User",
			Email: "user" + string(rune('a'+i)) + "@example.com",
		}
		if err := connections[0].DB.Create(user).Error; err != nil {
			t.Fatalf("Failed to create test user: %v", err)
		}
	}
	for i := 0; i < 3; i++ {
		post := &model.Post{
			ID:      int64(i + 1),
			UserID:  1,
			Title:   "Post",
			Content: "Content",
		}
		if err := connections[0].DB.Create(post).Error; err != nil {
			t.Fatalf("Failed to create test post: %v", err)
		}
	}

	// Shard 2: 1 user, 2 posts
	user := &model.User{
		ID:    100,
		Name:  "User",
		Email: "user100@example.com",
	}
	if err := connections[1].DB.Create(user).Error; err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
	for i := 0; i < 2; i++ {
		post := &model.Post{
			ID:      int64(100 + i),
			UserID:  100,
			Title:   "Post",
			Content: "Content",
		}
		if err := connections[1].DB.Create(post).Error; err != nil {
			t.Fatalf("Failed to create test post: %v", err)
		}
	}

	// 統計情報を取得
	stats, err := GetShardStats(manager)

	if err != nil {
		t.Fatalf("GetShardStats failed: %v", err)
	}

	if len(stats) != 2 {
		t.Errorf("expected 2 shard stats, got %d", len(stats))
	}

	// 合計を確認
	var totalUsers, totalPosts int64
	for _, s := range stats {
		totalUsers += s.UserCount
		totalPosts += s.PostCount
	}

	if totalUsers != 3 {
		t.Errorf("expected 3 total users, got %d", totalUsers)
	}

	if totalPosts != 5 {
		t.Errorf("expected 5 total posts, got %d", totalPosts)
	}
}

func TestGetShardForUserID(t *testing.T) {
	manager := testutil.SetupTestGORMShards(t, 2)
	defer testutil.CleanupTestGORMDB(manager)

	// ユーザーIDに基づいてシャードを取得
	gormDB, err := GetShardForUserID(manager, 1)

	if err != nil {
		t.Fatalf("GetShardForUserID failed: %v", err)
	}

	if gormDB == nil {
		t.Error("GetShardForUserID returned nil")
	}
}

func TestInsertToShard(t *testing.T) {
	manager := testutil.SetupTestGORMShards(t, 2)
	defer testutil.CleanupTestGORMDB(manager)

	// シャードにデータを挿入
	user := &model.User{
		ID:    1,
		Name:  "Test User",
		Email: "test@example.com",
	}

	err := InsertToShard(manager, 1, user)

	if err != nil {
		t.Fatalf("InsertToShard failed: %v", err)
	}

	// 挿入されたデータを確認
	gormDB, err := GetShardForUserID(manager, 1)
	if err != nil {
		t.Fatalf("GetShardForUserID failed: %v", err)
	}

	var foundUser model.User
	if err := gormDB.Where("id = ?", 1).First(&foundUser).Error; err != nil {
		t.Fatalf("Failed to find inserted user: %v", err)
	}

	if foundUser.Name != "Test User" {
		t.Errorf("expected user name 'Test User', got '%s'", foundUser.Name)
	}
}
