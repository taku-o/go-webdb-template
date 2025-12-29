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

func TestDmUserRepository_Create(t *testing.T) {
	groupManager := testutil.SetupTestGroupManager(t, 4, 8)
	defer testutil.CleanupTestGroupManager(groupManager)

	repo := repository.NewDmUserRepository(groupManager)
	ctx := context.Background()

	req := &model.CreateDmUserRequest{
		Name:  "Test User",
		Email: "test@example.com",
	}

	dmUser, err := repo.Create(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, dmUser)
	assert.NotZero(t, dmUser.ID)
	assert.Equal(t, "Test User", dmUser.Name)
	assert.Equal(t, "test@example.com", dmUser.Email)
	assert.NotZero(t, dmUser.CreatedAt)
	assert.NotZero(t, dmUser.UpdatedAt)
}

func TestDmUserRepository_GetByID(t *testing.T) {
	groupManager := testutil.SetupTestGroupManager(t, 4, 8)
	defer testutil.CleanupTestGroupManager(groupManager)

	repo := repository.NewDmUserRepository(groupManager)
	ctx := context.Background()

	// Create test dm_user first
	req := &model.CreateDmUserRequest{
		Name:  "Test User",
		Email: "test@example.com",
	}
	created, err := repo.Create(ctx, req)
	require.NoError(t, err)

	// Test retrieval
	dmUser, err := repo.GetByID(ctx, created.ID)
	assert.NoError(t, err)
	assert.NotNil(t, dmUser)
	assert.Equal(t, created.ID, dmUser.ID)
	assert.Equal(t, "Test User", dmUser.Name)
	assert.Equal(t, "test@example.com", dmUser.Email)
}

func TestDmUserRepository_GetByID_NotFound(t *testing.T) {
	groupManager := testutil.SetupTestGroupManager(t, 4, 8)
	defer testutil.CleanupTestGroupManager(groupManager)

	repo := repository.NewDmUserRepository(groupManager)
	ctx := context.Background()

	// Test retrieval of non-existent dm_user
	dmUser, err := repo.GetByID(ctx, 999)
	assert.Error(t, err)
	assert.Nil(t, dmUser)
}

func TestDmUserRepository_Update(t *testing.T) {
	groupManager := testutil.SetupTestGroupManager(t, 4, 8)
	defer testutil.CleanupTestGroupManager(groupManager)

	repo := repository.NewDmUserRepository(groupManager)
	ctx := context.Background()

	// Create test dm_user first
	createReq := &model.CreateDmUserRequest{
		Name:  "Original Name",
		Email: "original@example.com",
	}
	created, err := repo.Create(ctx, createReq)
	require.NoError(t, err)

	// Update dm_user
	updateReq := &model.UpdateDmUserRequest{
		Name:  "Updated Name",
		Email: "updated@example.com",
	}
	updated, err := repo.Update(ctx, created.ID, updateReq)
	assert.NoError(t, err)
	assert.NotNil(t, updated)
	assert.Equal(t, "Updated Name", updated.Name)
	assert.Equal(t, "updated@example.com", updated.Email)

	// Verify update
	dmUser, err := repo.GetByID(ctx, created.ID)
	require.NoError(t, err)
	assert.Equal(t, "Updated Name", dmUser.Name)
	assert.Equal(t, "updated@example.com", dmUser.Email)
}

func TestDmUserRepository_Delete(t *testing.T) {
	groupManager := testutil.SetupTestGroupManager(t, 4, 8)
	defer testutil.CleanupTestGroupManager(groupManager)

	repo := repository.NewDmUserRepository(groupManager)
	ctx := context.Background()

	// Create test dm_user first
	req := &model.CreateDmUserRequest{
		Name:  "Test User",
		Email: "test@example.com",
	}
	created, err := repo.Create(ctx, req)
	require.NoError(t, err)

	// Delete dm_user
	err = repo.Delete(ctx, created.ID)
	assert.NoError(t, err)

	// Verify deletion
	dmUser, err := repo.GetByID(ctx, created.ID)
	assert.Error(t, err)
	assert.Nil(t, dmUser)
}

func TestDmUserRepository_List(t *testing.T) {
	groupManager := testutil.SetupTestGroupManager(t, 4, 8)
	defer testutil.CleanupTestGroupManager(groupManager)

	repo := repository.NewDmUserRepository(groupManager)
	ctx := context.Background()

	// Create test dm_users
	req1 := &model.CreateDmUserRequest{
		Name:  "User 1",
		Email: "user1@example.com",
	}
	_, err := repo.Create(ctx, req1)
	require.NoError(t, err)

	req2 := &model.CreateDmUserRequest{
		Name:  "User 2",
		Email: "user2@example.com",
	}
	_, err = repo.Create(ctx, req2)
	require.NoError(t, err)

	// List dm_users
	dmUsers, err := repo.List(ctx, 10, 0)
	assert.NoError(t, err)
	assert.Len(t, dmUsers, 2)
}
