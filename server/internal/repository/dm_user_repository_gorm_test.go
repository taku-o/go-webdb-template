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

func TestDmUserRepositoryGORM_Create(t *testing.T) {
	groupManager := testutil.SetupTestGroupManager(t, 4, 8)
	defer testutil.CleanupTestGroupManager(groupManager)

	dmUserRepo := repository.NewDmUserRepositoryGORM(groupManager)
	ctx := context.Background()

	req := &model.CreateDmUserRequest{
		Name:  "Test User",
		Email: "test@example.com",
	}

	dmUser, err := dmUserRepo.Create(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, dmUser)
	assert.NotZero(t, dmUser.ID)
	assert.Equal(t, "Test User", dmUser.Name)
	assert.Equal(t, "test@example.com", dmUser.Email)
	assert.NotZero(t, dmUser.CreatedAt)
	assert.NotZero(t, dmUser.UpdatedAt)
}

func TestDmUserRepositoryGORM_GetByID(t *testing.T) {
	groupManager := testutil.SetupTestGroupManager(t, 4, 8)
	defer testutil.CleanupTestGroupManager(groupManager)

	dmUserRepo := repository.NewDmUserRepositoryGORM(groupManager)
	ctx := context.Background()

	// Create test user first
	req := &model.CreateDmUserRequest{
		Name:  "Test User",
		Email: "test@example.com",
	}
	created, err := dmUserRepo.Create(ctx, req)
	require.NoError(t, err)

	// Test retrieval
	dmUser, err := dmUserRepo.GetByID(ctx, created.ID)
	assert.NoError(t, err)
	assert.NotNil(t, dmUser)
	assert.Equal(t, created.ID, dmUser.ID)
	assert.Equal(t, "Test User", dmUser.Name)
	assert.Equal(t, "test@example.com", dmUser.Email)
}

func TestDmUserRepositoryGORM_GetByID_NotFound(t *testing.T) {
	groupManager := testutil.SetupTestGroupManager(t, 4, 8)
	defer testutil.CleanupTestGroupManager(groupManager)

	dmUserRepo := repository.NewDmUserRepositoryGORM(groupManager)
	ctx := context.Background()

	// Test retrieval of non-existent user
	dmUser, err := dmUserRepo.GetByID(ctx, "00000000000000000000000000000000")
	assert.Error(t, err)
	assert.Nil(t, dmUser)
}

func TestDmUserRepositoryGORM_Update(t *testing.T) {
	groupManager := testutil.SetupTestGroupManager(t, 4, 8)
	defer testutil.CleanupTestGroupManager(groupManager)

	dmUserRepo := repository.NewDmUserRepositoryGORM(groupManager)
	ctx := context.Background()

	// Create test user first
	createReq := &model.CreateDmUserRequest{
		Name:  "Original Name",
		Email: "original@example.com",
	}
	created, err := dmUserRepo.Create(ctx, createReq)
	require.NoError(t, err)

	// Update user
	updateReq := &model.UpdateDmUserRequest{
		Name:  "Updated Name",
		Email: "updated@example.com",
	}
	updated, err := dmUserRepo.Update(ctx, created.ID, updateReq)
	assert.NoError(t, err)
	assert.NotNil(t, updated)
	assert.Equal(t, "Updated Name", updated.Name)
	assert.Equal(t, "updated@example.com", updated.Email)

	// Verify update
	dmUser, err := dmUserRepo.GetByID(ctx, created.ID)
	require.NoError(t, err)
	assert.Equal(t, "Updated Name", dmUser.Name)
	assert.Equal(t, "updated@example.com", dmUser.Email)
}

func TestDmUserRepositoryGORM_Delete(t *testing.T) {
	groupManager := testutil.SetupTestGroupManager(t, 4, 8)
	defer testutil.CleanupTestGroupManager(groupManager)

	dmUserRepo := repository.NewDmUserRepositoryGORM(groupManager)
	ctx := context.Background()

	// Create test user first
	req := &model.CreateDmUserRequest{
		Name:  "Test User",
		Email: "test@example.com",
	}
	created, err := dmUserRepo.Create(ctx, req)
	require.NoError(t, err)

	// Delete user
	err = dmUserRepo.Delete(ctx, created.ID)
	assert.NoError(t, err)

	// Verify deletion
	dmUser, err := dmUserRepo.GetByID(ctx, created.ID)
	assert.Error(t, err)
	assert.Nil(t, dmUser)
}

func TestDmUserRepositoryGORM_List(t *testing.T) {
	groupManager := testutil.SetupTestGroupManager(t, 4, 8)
	defer testutil.CleanupTestGroupManager(groupManager)

	dmUserRepo := repository.NewDmUserRepositoryGORM(groupManager)
	ctx := context.Background()

	// Create test users
	req1 := &model.CreateDmUserRequest{
		Name:  "User 1",
		Email: "user1@example.com",
	}
	_, err := dmUserRepo.Create(ctx, req1)
	require.NoError(t, err)

	req2 := &model.CreateDmUserRequest{
		Name:  "User 2",
		Email: "user2@example.com",
	}
	_, err = dmUserRepo.Create(ctx, req2)
	require.NoError(t, err)

	// List users (cross-table query)
	dmUsers, err := dmUserRepo.List(ctx, 10, 0)
	assert.NoError(t, err)
	assert.Len(t, dmUsers, 2)
}
