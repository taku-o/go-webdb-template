package repository_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/example/go-webdb-template/internal/model"
	"github.com/example/go-webdb-template/internal/repository"
	"github.com/example/go-webdb-template/test/testutil"
)

func TestUserRepositoryGORM_Create(t *testing.T) {
	groupManager := testutil.SetupTestGroupManager(t, 4, 8)
	defer testutil.CleanupTestGroupManager(groupManager)

	repo := repository.NewUserRepositoryGORM(groupManager)
	ctx := context.Background()

	req := &model.CreateUserRequest{
		Name:  "Test User",
		Email: "test@example.com",
	}

	user, err := repo.Create(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.NotZero(t, user.ID)
	assert.Equal(t, "Test User", user.Name)
	assert.Equal(t, "test@example.com", user.Email)
	assert.NotZero(t, user.CreatedAt)
	assert.NotZero(t, user.UpdatedAt)
}

func TestUserRepositoryGORM_GetByID(t *testing.T) {
	groupManager := testutil.SetupTestGroupManager(t, 4, 8)
	defer testutil.CleanupTestGroupManager(groupManager)

	repo := repository.NewUserRepositoryGORM(groupManager)
	ctx := context.Background()

	// Create test user first
	req := &model.CreateUserRequest{
		Name:  "Test User",
		Email: "test@example.com",
	}
	created, err := repo.Create(ctx, req)
	require.NoError(t, err)

	// Test retrieval
	user, err := repo.GetByID(ctx, created.ID)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, created.ID, user.ID)
	assert.Equal(t, "Test User", user.Name)
	assert.Equal(t, "test@example.com", user.Email)
}

func TestUserRepositoryGORM_GetByID_NotFound(t *testing.T) {
	groupManager := testutil.SetupTestGroupManager(t, 4, 8)
	defer testutil.CleanupTestGroupManager(groupManager)

	repo := repository.NewUserRepositoryGORM(groupManager)
	ctx := context.Background()

	// Test retrieval of non-existent user
	user, err := repo.GetByID(ctx, 999)
	assert.Error(t, err)
	assert.Nil(t, user)
}

func TestUserRepositoryGORM_Update(t *testing.T) {
	groupManager := testutil.SetupTestGroupManager(t, 4, 8)
	defer testutil.CleanupTestGroupManager(groupManager)

	repo := repository.NewUserRepositoryGORM(groupManager)
	ctx := context.Background()

	// Create test user first
	createReq := &model.CreateUserRequest{
		Name:  "Original Name",
		Email: "original@example.com",
	}
	created, err := repo.Create(ctx, createReq)
	require.NoError(t, err)

	// Update user
	updateReq := &model.UpdateUserRequest{
		Name:  "Updated Name",
		Email: "updated@example.com",
	}
	updated, err := repo.Update(ctx, created.ID, updateReq)
	assert.NoError(t, err)
	assert.NotNil(t, updated)
	assert.Equal(t, "Updated Name", updated.Name)
	assert.Equal(t, "updated@example.com", updated.Email)

	// Verify update
	user, err := repo.GetByID(ctx, created.ID)
	require.NoError(t, err)
	assert.Equal(t, "Updated Name", user.Name)
	assert.Equal(t, "updated@example.com", user.Email)
}

func TestUserRepositoryGORM_Delete(t *testing.T) {
	groupManager := testutil.SetupTestGroupManager(t, 4, 8)
	defer testutil.CleanupTestGroupManager(groupManager)

	repo := repository.NewUserRepositoryGORM(groupManager)
	ctx := context.Background()

	// Create test user first
	req := &model.CreateUserRequest{
		Name:  "Test User",
		Email: "test@example.com",
	}
	created, err := repo.Create(ctx, req)
	require.NoError(t, err)

	// Delete user
	err = repo.Delete(ctx, created.ID)
	assert.NoError(t, err)

	// Verify deletion
	user, err := repo.GetByID(ctx, created.ID)
	assert.Error(t, err)
	assert.Nil(t, user)
}

func TestUserRepositoryGORM_List(t *testing.T) {
	groupManager := testutil.SetupTestGroupManager(t, 4, 8)
	defer testutil.CleanupTestGroupManager(groupManager)

	repo := repository.NewUserRepositoryGORM(groupManager)
	ctx := context.Background()

	// Create test users
	req1 := &model.CreateUserRequest{
		Name:  "User 1",
		Email: "user1@example.com",
	}
	_, err := repo.Create(ctx, req1)
	require.NoError(t, err)

	req2 := &model.CreateUserRequest{
		Name:  "User 2",
		Email: "user2@example.com",
	}
	_, err = repo.Create(ctx, req2)
	require.NoError(t, err)

	// List users (cross-table query)
	users, err := repo.List(ctx, 10, 0)
	assert.NoError(t, err)
	assert.Len(t, users, 2)
}
