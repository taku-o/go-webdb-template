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

func TestPostRepositoryGORM_Create(t *testing.T) {
	groupManager := testutil.SetupTestGroupManager(t, 4, 8)
	defer testutil.CleanupTestGroupManager(groupManager)

	repo := repository.NewPostRepositoryGORM(groupManager)
	ctx := context.Background()

	req := &model.CreatePostRequest{
		UserID:  1,
		Title:   "Test Post",
		Content: "This is a test post content",
	}

	post, err := repo.Create(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, post)
	assert.NotZero(t, post.ID)
	assert.Equal(t, int64(1), post.UserID)
	assert.Equal(t, "Test Post", post.Title)
	assert.Equal(t, "This is a test post content", post.Content)
	assert.NotZero(t, post.CreatedAt)
	assert.NotZero(t, post.UpdatedAt)
}

func TestPostRepositoryGORM_GetByID(t *testing.T) {
	groupManager := testutil.SetupTestGroupManager(t, 4, 8)
	defer testutil.CleanupTestGroupManager(groupManager)

	repo := repository.NewPostRepositoryGORM(groupManager)
	ctx := context.Background()

	// Create test post first
	req := &model.CreatePostRequest{
		UserID:  1,
		Title:   "Test Post",
		Content: "Test content",
	}
	created, err := repo.Create(ctx, req)
	require.NoError(t, err)

	// Test retrieval
	post, err := repo.GetByID(ctx, created.ID, created.UserID)
	assert.NoError(t, err)
	assert.NotNil(t, post)
	assert.Equal(t, created.ID, post.ID)
	assert.Equal(t, int64(1), post.UserID)
	assert.Equal(t, "Test Post", post.Title)
	assert.Equal(t, "Test content", post.Content)
}

func TestPostRepositoryGORM_GetByID_NotFound(t *testing.T) {
	groupManager := testutil.SetupTestGroupManager(t, 4, 8)
	defer testutil.CleanupTestGroupManager(groupManager)

	repo := repository.NewPostRepositoryGORM(groupManager)
	ctx := context.Background()

	// Test retrieval of non-existent post
	post, err := repo.GetByID(ctx, 999, 1)
	assert.Error(t, err)
	assert.Nil(t, post)
}

func TestPostRepositoryGORM_Update(t *testing.T) {
	groupManager := testutil.SetupTestGroupManager(t, 4, 8)
	defer testutil.CleanupTestGroupManager(groupManager)

	repo := repository.NewPostRepositoryGORM(groupManager)
	ctx := context.Background()

	// Create test post first
	createReq := &model.CreatePostRequest{
		UserID:  1,
		Title:   "Original Title",
		Content: "Original content",
	}
	created, err := repo.Create(ctx, createReq)
	require.NoError(t, err)

	// Update post
	updateReq := &model.UpdatePostRequest{
		Title:   "Updated Title",
		Content: "Updated content",
	}
	updated, err := repo.Update(ctx, created.ID, created.UserID, updateReq)
	assert.NoError(t, err)
	assert.NotNil(t, updated)
	assert.Equal(t, "Updated Title", updated.Title)
	assert.Equal(t, "Updated content", updated.Content)

	// Verify update
	post, err := repo.GetByID(ctx, created.ID, created.UserID)
	require.NoError(t, err)
	assert.Equal(t, "Updated Title", post.Title)
	assert.Equal(t, "Updated content", post.Content)
}

func TestPostRepositoryGORM_Delete(t *testing.T) {
	groupManager := testutil.SetupTestGroupManager(t, 4, 8)
	defer testutil.CleanupTestGroupManager(groupManager)

	repo := repository.NewPostRepositoryGORM(groupManager)
	ctx := context.Background()

	// Create test post first
	req := &model.CreatePostRequest{
		UserID:  1,
		Title:   "Test Post",
		Content: "Test content",
	}
	created, err := repo.Create(ctx, req)
	require.NoError(t, err)

	// Delete post
	err = repo.Delete(ctx, created.ID, created.UserID)
	assert.NoError(t, err)

	// Verify deletion
	post, err := repo.GetByID(ctx, created.ID, created.UserID)
	assert.Error(t, err)
	assert.Nil(t, post)
}

func TestPostRepositoryGORM_ListByUserID(t *testing.T) {
	groupManager := testutil.SetupTestGroupManager(t, 4, 8)
	defer testutil.CleanupTestGroupManager(groupManager)

	repo := repository.NewPostRepositoryGORM(groupManager)
	ctx := context.Background()

	// Create test posts
	req1 := &model.CreatePostRequest{
		UserID:  1,
		Title:   "Post 1",
		Content: "Content 1",
	}
	_, err := repo.Create(ctx, req1)
	require.NoError(t, err)

	req2 := &model.CreatePostRequest{
		UserID:  1,
		Title:   "Post 2",
		Content: "Content 2",
	}
	_, err = repo.Create(ctx, req2)
	require.NoError(t, err)

	// Get posts by user ID
	posts, err := repo.ListByUserID(ctx, 1, 10, 0)
	assert.NoError(t, err)
	assert.Len(t, posts, 2)
}

func TestPostRepositoryGORM_List(t *testing.T) {
	groupManager := testutil.SetupTestGroupManager(t, 4, 8)
	defer testutil.CleanupTestGroupManager(groupManager)

	repo := repository.NewPostRepositoryGORM(groupManager)
	ctx := context.Background()

	// Create test posts with different UserIDs (different tables)
	req1 := &model.CreatePostRequest{
		UserID:  1,
		Title:   "Post 1",
		Content: "Content 1",
	}
	_, err := repo.Create(ctx, req1)
	require.NoError(t, err)

	req2 := &model.CreatePostRequest{
		UserID:  2,
		Title:   "Post 2",
		Content: "Content 2",
	}
	_, err = repo.Create(ctx, req2)
	require.NoError(t, err)

	// List all posts (cross-table query)
	posts, err := repo.List(ctx, 10, 0)
	assert.NoError(t, err)
	assert.Len(t, posts, 2)
}

func TestPostRepositoryGORM_GetUserPosts(t *testing.T) {
	groupManager := testutil.SetupTestGroupManager(t, 4, 8)
	defer testutil.CleanupTestGroupManager(groupManager)

	userRepo := repository.NewUserRepositoryGORM(groupManager)
	postRepo := repository.NewPostRepositoryGORM(groupManager)
	ctx := context.Background()

	// Create test user
	userReq := &model.CreateUserRequest{
		Name:  "Test User",
		Email: "test@example.com",
	}
	user, err := userRepo.Create(ctx, userReq)
	require.NoError(t, err)

	// Create test post
	postReq := &model.CreatePostRequest{
		UserID:  user.ID,
		Title:   "Test Post",
		Content: "Test content",
	}
	_, err = postRepo.Create(ctx, postReq)
	require.NoError(t, err)

	// Get user posts with JOIN (cross-table query)
	userPosts, err := postRepo.GetUserPosts(ctx, 10, 0)
	assert.NoError(t, err)
	assert.NotEmpty(t, userPosts)

	// Verify the first result
	assert.Equal(t, user.ID, userPosts[0].UserID)
	assert.Equal(t, "Test User", userPosts[0].UserName)
	assert.Equal(t, "test@example.com", userPosts[0].UserEmail)
	assert.Equal(t, "Test Post", userPosts[0].PostTitle)
	assert.Equal(t, "Test content", userPosts[0].PostContent)
}
