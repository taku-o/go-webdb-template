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

func TestPostRepository_Create(t *testing.T) {
	groupManager := testutil.SetupTestGroupManager(t, 4, 8)
	defer testutil.CleanupTestGroupManager(groupManager)

	repo := repository.NewDmPostRepository(groupManager)
	ctx := context.Background()

	req := &model.CreateDmPostRequest{
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

func TestPostRepository_GetByID(t *testing.T) {
	groupManager := testutil.SetupTestGroupManager(t, 4, 8)
	defer testutil.CleanupTestGroupManager(groupManager)

	repo := repository.NewDmPostRepository(groupManager)
	ctx := context.Background()

	// Create test post first
	req := &model.CreateDmPostRequest{
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

func TestPostRepository_GetByID_NotFound(t *testing.T) {
	groupManager := testutil.SetupTestGroupManager(t, 4, 8)
	defer testutil.CleanupTestGroupManager(groupManager)

	repo := repository.NewDmPostRepository(groupManager)
	ctx := context.Background()

	// Test retrieval of non-existent post
	post, err := repo.GetByID(ctx, 999, 1)
	assert.Error(t, err)
	assert.Nil(t, post)
}

func TestPostRepository_Update(t *testing.T) {
	groupManager := testutil.SetupTestGroupManager(t, 4, 8)
	defer testutil.CleanupTestGroupManager(groupManager)

	repo := repository.NewDmPostRepository(groupManager)
	ctx := context.Background()

	// Create test post first
	createReq := &model.CreateDmPostRequest{
		UserID:  1,
		Title:   "Original Title",
		Content: "Original content",
	}
	created, err := repo.Create(ctx, createReq)
	require.NoError(t, err)

	// Update post
	updateReq := &model.UpdateDmPostRequest{
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

func TestPostRepository_Delete(t *testing.T) {
	groupManager := testutil.SetupTestGroupManager(t, 4, 8)
	defer testutil.CleanupTestGroupManager(groupManager)

	repo := repository.NewDmPostRepository(groupManager)
	ctx := context.Background()

	// Create test post first
	req := &model.CreateDmPostRequest{
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

func TestPostRepository_ListByUserID(t *testing.T) {
	groupManager := testutil.SetupTestGroupManager(t, 4, 8)
	defer testutil.CleanupTestGroupManager(groupManager)

	repo := repository.NewDmPostRepository(groupManager)
	ctx := context.Background()

	// Create test posts
	req1 := &model.CreateDmPostRequest{
		UserID:  1,
		Title:   "Post 1",
		Content: "Content 1",
	}
	_, err := repo.Create(ctx, req1)
	require.NoError(t, err)

	req2 := &model.CreateDmPostRequest{
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
