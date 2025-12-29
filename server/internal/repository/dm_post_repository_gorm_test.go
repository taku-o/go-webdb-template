package repository_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/taku-o/go-webdb-template/internal/model"
	"github.com/taku-o/go-webdb-template/internal/repository"
	"github.com/taku-o/go-webdb-template/test/testutil"
)

func TestDmPostRepositoryGORM_Create(t *testing.T) {
	groupManager := testutil.SetupTestGroupManager(t, 4, 8)
	defer testutil.CleanupTestGroupManager(groupManager)

	dmPostRepo := repository.NewDmPostRepositoryGORM(groupManager)
	ctx := context.Background()

	req := &model.CreateDmPostRequest{
		UserID:  1,
		Title:   "Test Post",
		Content: "This is a test post content",
	}

	dmPost, err := dmPostRepo.Create(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, dmPost)
	assert.NotZero(t, dmPost.ID)
	assert.Equal(t, int64(1), dmPost.UserID)
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

	// Create test post first
	req := &model.CreateDmPostRequest{
		UserID:  1,
		Title:   "Test Post",
		Content: "Test content",
	}
	created, err := dmPostRepo.Create(ctx, req)
	require.NoError(t, err)

	// Test retrieval
	dmPost, err := dmPostRepo.GetByID(ctx, created.ID, created.UserID)
	assert.NoError(t, err)
	assert.NotNil(t, dmPost)
	assert.Equal(t, created.ID, dmPost.ID)
	assert.Equal(t, int64(1), dmPost.UserID)
	assert.Equal(t, "Test Post", dmPost.Title)
	assert.Equal(t, "Test content", dmPost.Content)
}

func TestDmPostRepositoryGORM_GetByID_NotFound(t *testing.T) {
	groupManager := testutil.SetupTestGroupManager(t, 4, 8)
	defer testutil.CleanupTestGroupManager(groupManager)

	dmPostRepo := repository.NewDmPostRepositoryGORM(groupManager)
	ctx := context.Background()

	// Test retrieval of non-existent post
	dmPost, err := dmPostRepo.GetByID(ctx, 999, 1)
	assert.Error(t, err)
	assert.Nil(t, dmPost)
}

func TestDmPostRepositoryGORM_Update(t *testing.T) {
	groupManager := testutil.SetupTestGroupManager(t, 4, 8)
	defer testutil.CleanupTestGroupManager(groupManager)

	dmPostRepo := repository.NewDmPostRepositoryGORM(groupManager)
	ctx := context.Background()

	// Create test post first
	createReq := &model.CreateDmPostRequest{
		UserID:  1,
		Title:   "Original Title",
		Content: "Original content",
	}
	created, err := dmPostRepo.Create(ctx, createReq)
	require.NoError(t, err)

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

	// Create test post first
	req := &model.CreateDmPostRequest{
		UserID:  1,
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

	// Create test posts
	req1 := &model.CreateDmPostRequest{
		UserID:  1,
		Title:   "Post 1",
		Content: "Content 1",
	}
	_, err := dmPostRepo.Create(ctx, req1)
	require.NoError(t, err)

	req2 := &model.CreateDmPostRequest{
		UserID:  1,
		Title:   "Post 2",
		Content: "Content 2",
	}
	_, err = dmPostRepo.Create(ctx, req2)
	require.NoError(t, err)

	// Get posts by user ID
	dmPosts, err := dmPostRepo.ListByUserID(ctx, 1, 10, 0)
	assert.NoError(t, err)
	assert.Len(t, dmPosts, 2)
}

func TestDmPostRepositoryGORM_List(t *testing.T) {
	groupManager := testutil.SetupTestGroupManager(t, 4, 8)
	defer testutil.CleanupTestGroupManager(groupManager)

	dmPostRepo := repository.NewDmPostRepositoryGORM(groupManager)
	ctx := context.Background()

	// Create test posts with different UserIDs (different tables)
	req1 := &model.CreateDmPostRequest{
		UserID:  1,
		Title:   "Post 1",
		Content: "Content 1",
	}
	_, err := dmPostRepo.Create(ctx, req1)
	require.NoError(t, err)

	req2 := &model.CreateDmPostRequest{
		UserID:  2,
		Title:   "Post 2",
		Content: "Content 2",
	}
	_, err = dmPostRepo.Create(ctx, req2)
	require.NoError(t, err)

	// List all posts (cross-table query)
	dmPosts, err := dmPostRepo.List(ctx, 10, 0)
	assert.NoError(t, err)
	assert.Len(t, dmPosts, 2)
}

func TestDmPostRepositoryGORM_GetUserPosts(t *testing.T) {
	groupManager := testutil.SetupTestGroupManager(t, 4, 8)
	defer testutil.CleanupTestGroupManager(groupManager)

	dmUserRepo := repository.NewDmUserRepositoryGORM(groupManager)
	dmPostRepo := repository.NewDmPostRepositoryGORM(groupManager)
	ctx := context.Background()

	// Create test user
	dmUserReq := &model.CreateDmUserRequest{
		Name:  "Test User",
		Email: "test@example.com",
	}
	dmUser, err := dmUserRepo.Create(ctx, dmUserReq)
	require.NoError(t, err)

	// Create test post
	dmPostReq := &model.CreateDmPostRequest{
		UserID:  dmUser.ID,
		Title:   "Test Post",
		Content: "Test content",
	}
	_, err = dmPostRepo.Create(ctx, dmPostReq)
	require.NoError(t, err)

	// Get user posts with JOIN (cross-table query)
	dmUserPosts, err := dmPostRepo.GetUserPosts(ctx, 10, 0)
	assert.NoError(t, err)
	assert.NotEmpty(t, dmUserPosts)

	// Verify the first result
	assert.Equal(t, dmUser.ID, dmUserPosts[0].UserID)
	assert.Equal(t, "Test User", dmUserPosts[0].UserName)
	assert.Equal(t, "test@example.com", dmUserPosts[0].UserEmail)
	assert.Equal(t, "Test Post", dmUserPosts[0].PostTitle)
	assert.Equal(t, "Test content", dmUserPosts[0].PostContent)
}
