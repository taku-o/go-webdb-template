package repository_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/taku-o/go-webdb-template/internal/model"
	"github.com/taku-o/go-webdb-template/internal/repository"
	"github.com/taku-o/go-webdb-template/internal/util/idgen"
	"github.com/taku-o/go-webdb-template/test/testutil"
)

func TestDmPostRepositoryGORM_Create(t *testing.T) {
	groupManager := testutil.SetupTestGroupManager(t, 4, 8)
	defer testutil.CleanupTestGroupManager(groupManager)

	dmPostRepo := repository.NewDmPostRepositoryGORM(groupManager)
	ctx := context.Background()

	// テスト用のユーザーIDを生成
	userID, err := idgen.GenerateUUIDv7()
	require.NoError(t, err)

	req := &model.CreateDmPostRequest{
		UserID:  userID,
		Title:   "Test Post",
		Content: "This is a test post content",
	}

	dmPost, err := dmPostRepo.Create(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, dmPost)

	// クリーンアップ
	defer func() {
		if dmPost != nil {
			_ = dmPostRepo.Delete(ctx, dmPost.ID, dmPost.UserID)
		}
	}()

	assert.NotEmpty(t, dmPost.ID)
	assert.Equal(t, userID, dmPost.UserID)
	assert.Equal(t, "Test Post", dmPost.Title)
	assert.Equal(t, "This is a test post content", dmPost.Content)
	assert.NotZero(t, dmPost.CreatedAt)
	assert.NotZero(t, dmPost.UpdatedAt)
}

func TestDmPostRepositoryGORM_GetByID(t *testing.T) {
	groupManager := testutil.SetupTestGroupManager(t, 4, 8)
	defer testutil.CleanupTestGroupManager(groupManager)

	dmPostRepo := repository.NewDmPostRepositoryGORM(groupManager)
	ctx := context.Background()

	// テスト用のユーザーIDを生成
	userID, err := idgen.GenerateUUIDv7()
	require.NoError(t, err)

	// Create test post first
	req := &model.CreateDmPostRequest{
		UserID:  userID,
		Title:   "Test Post",
		Content: "Test content",
	}
	created, err := dmPostRepo.Create(ctx, req)
	require.NoError(t, err)

	// クリーンアップ
	defer func() {
		_ = dmPostRepo.Delete(ctx, created.ID, created.UserID)
	}()

	// Test retrieval
	dmPost, err := dmPostRepo.GetByID(ctx, created.ID, created.UserID)
	assert.NoError(t, err)
	assert.NotNil(t, dmPost)
	assert.Equal(t, created.ID, dmPost.ID)
	assert.Equal(t, userID, dmPost.UserID)
	assert.Equal(t, "Test Post", dmPost.Title)
	assert.Equal(t, "Test content", dmPost.Content)
}

func TestDmPostRepositoryGORM_GetByID_NotFound(t *testing.T) {
	groupManager := testutil.SetupTestGroupManager(t, 4, 8)
	defer testutil.CleanupTestGroupManager(groupManager)

	dmPostRepo := repository.NewDmPostRepositoryGORM(groupManager)
	ctx := context.Background()

	// テスト用のユーザーIDを生成
	userID, err := idgen.GenerateUUIDv7()
	require.NoError(t, err)

	// Test retrieval of non-existent post
	dmPost, err := dmPostRepo.GetByID(ctx, "00000000000000000000000000000000", userID)
	assert.Error(t, err)
	assert.Nil(t, dmPost)
}

func TestDmPostRepositoryGORM_Update(t *testing.T) {
	groupManager := testutil.SetupTestGroupManager(t, 4, 8)
	defer testutil.CleanupTestGroupManager(groupManager)

	dmPostRepo := repository.NewDmPostRepositoryGORM(groupManager)
	ctx := context.Background()

	// テスト用のユーザーIDを生成
	userID, err := idgen.GenerateUUIDv7()
	require.NoError(t, err)

	// Create test post first
	createReq := &model.CreateDmPostRequest{
		UserID:  userID,
		Title:   "Original Title",
		Content: "Original content",
	}
	created, err := dmPostRepo.Create(ctx, createReq)
	require.NoError(t, err)

	// クリーンアップ
	defer func() {
		_ = dmPostRepo.Delete(ctx, created.ID, created.UserID)
	}()

	// Update post
	updateReq := &model.UpdateDmPostRequest{
		Title:   "Updated Title",
		Content: "Updated content",
	}
	updated, err := dmPostRepo.Update(ctx, created.ID, created.UserID, updateReq)
	assert.NoError(t, err)
	assert.NotNil(t, updated)
	assert.Equal(t, "Updated Title", updated.Title)
	assert.Equal(t, "Updated content", updated.Content)

	// Verify update
	dmPost, err := dmPostRepo.GetByID(ctx, created.ID, created.UserID)
	require.NoError(t, err)
	assert.Equal(t, "Updated Title", dmPost.Title)
	assert.Equal(t, "Updated content", dmPost.Content)
}

func TestDmPostRepositoryGORM_Delete(t *testing.T) {
	groupManager := testutil.SetupTestGroupManager(t, 4, 8)
	defer testutil.CleanupTestGroupManager(groupManager)

	dmPostRepo := repository.NewDmPostRepositoryGORM(groupManager)
	ctx := context.Background()

	// テスト用のユーザーIDを生成
	userID, err := idgen.GenerateUUIDv7()
	require.NoError(t, err)

	// Create test post first
	req := &model.CreateDmPostRequest{
		UserID:  userID,
		Title:   "Test Post",
		Content: "Test content",
	}
	created, err := dmPostRepo.Create(ctx, req)
	require.NoError(t, err)

	// Delete post
	err = dmPostRepo.Delete(ctx, created.ID, created.UserID)
	assert.NoError(t, err)

	// Verify deletion
	dmPost, err := dmPostRepo.GetByID(ctx, created.ID, created.UserID)
	assert.Error(t, err)
	assert.Nil(t, dmPost)
}

func TestDmPostRepositoryGORM_ListByUserID(t *testing.T) {
	groupManager := testutil.SetupTestGroupManager(t, 4, 8)
	defer testutil.CleanupTestGroupManager(groupManager)

	dmPostRepo := repository.NewDmPostRepositoryGORM(groupManager)
	ctx := context.Background()

	// テスト用のユーザーIDを生成
	userID, err := idgen.GenerateUUIDv7()
	require.NoError(t, err)

	// Create test posts
	req1 := &model.CreateDmPostRequest{
		UserID:  userID,
		Title:   "Post 1",
		Content: "Content 1",
	}
	post1, err := dmPostRepo.Create(ctx, req1)
	require.NoError(t, err)

	req2 := &model.CreateDmPostRequest{
		UserID:  userID,
		Title:   "Post 2",
		Content: "Content 2",
	}
	post2, err := dmPostRepo.Create(ctx, req2)
	require.NoError(t, err)

	// クリーンアップ
	defer func() {
		_ = dmPostRepo.Delete(ctx, post1.ID, post1.UserID)
		_ = dmPostRepo.Delete(ctx, post2.ID, post2.UserID)
	}()

	// Get posts by user ID
	dmPosts, err := dmPostRepo.ListByUserID(ctx, userID, 10, 0)
	assert.NoError(t, err)
	assert.Len(t, dmPosts, 2)
}

func TestDmPostRepositoryGORM_List(t *testing.T) {
	groupManager := testutil.SetupTestGroupManager(t, 4, 8)
	defer testutil.CleanupTestGroupManager(groupManager)

	dmPostRepo := repository.NewDmPostRepositoryGORM(groupManager)
	ctx := context.Background()

	// テスト用のユーザーIDを生成（同一テーブルに対してテスト）
	userID, err := idgen.GenerateUUIDv7()
	require.NoError(t, err)

	// テスト前の件数を取得（特定テーブルのみ）
	initialPosts, err := dmPostRepo.ListByUserID(ctx, userID, 1000, 0)
	require.NoError(t, err)
	initialCount := len(initialPosts)

	// Create test posts with same UserID (same table)
	req1 := &model.CreateDmPostRequest{
		UserID:  userID,
		Title:   "Post 1",
		Content: "Content 1",
	}
	post1, err := dmPostRepo.Create(ctx, req1)
	require.NoError(t, err)

	req2 := &model.CreateDmPostRequest{
		UserID:  userID,
		Title:   "Post 2",
		Content: "Content 2",
	}
	post2, err := dmPostRepo.Create(ctx, req2)
	require.NoError(t, err)

	// クリーンアップ
	defer func() {
		_ = dmPostRepo.Delete(ctx, post1.ID, post1.UserID)
		_ = dmPostRepo.Delete(ctx, post2.ID, post2.UserID)
	}()

	// List posts by user ID (single table query)
	dmPosts, err := dmPostRepo.ListByUserID(ctx, userID, 1000, 0)
	assert.NoError(t, err)
	// 2件増えたことを確認
	assert.Equal(t, initialCount+2, len(dmPosts))
}

func TestDmPostRepositoryGORM_GetUserPosts(t *testing.T) {
	groupManager := testutil.SetupTestGroupManager(t, 4, 8)
	defer testutil.CleanupTestGroupManager(groupManager)

	dmUserRepo := repository.NewDmUserRepositoryGORM(groupManager)
	dmPostRepo := repository.NewDmPostRepositoryGORM(groupManager)
	ctx := context.Background()

	// ユニークなメールアドレスを生成
	uniqueID, err := idgen.GenerateUUIDv7()
	require.NoError(t, err)
	uniqueEmail := fmt.Sprintf("test-%s@example.com", uniqueID)

	// Create test user
	dmUserReq := &model.CreateDmUserRequest{
		Name:  "Test User",
		Email: uniqueEmail,
	}
	dmUser, err := dmUserRepo.Create(ctx, dmUserReq)
	require.NoError(t, err)

	// Create test post
	dmPostReq := &model.CreateDmPostRequest{
		UserID:  dmUser.ID,
		Title:   "Test Post",
		Content: "Test content",
	}
	dmPost, err := dmPostRepo.Create(ctx, dmPostReq)
	require.NoError(t, err)

	// クリーンアップ
	defer func() {
		_ = dmPostRepo.Delete(ctx, dmPost.ID, dmPost.UserID)
		_ = dmUserRepo.Delete(ctx, dmUser.ID)
	}()

	// Get user posts with JOIN (cross-table query)
	dmUserPosts, err := dmPostRepo.GetUserPosts(ctx, 1000, 0)
	assert.NoError(t, err)
	assert.NotEmpty(t, dmUserPosts)

	// 作成したユーザーの投稿を検索
	var found bool
	for _, up := range dmUserPosts {
		if up.UserID == dmUser.ID {
			assert.Equal(t, "Test User", up.UserName)
			assert.Equal(t, uniqueEmail, up.UserEmail)
			assert.Equal(t, "Test Post", up.PostTitle)
			assert.Equal(t, "Test content", up.PostContent)
			found = true
			break
		}
	}
	assert.True(t, found, "Created user post should be found in the list")
}
