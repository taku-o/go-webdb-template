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

func TestDmUserRepositoryGORM_Create(t *testing.T) {
	groupManager := testutil.SetupTestGroupManager(t, 4, 8)
	defer testutil.CleanupTestGroupManager(groupManager)

	dmUserRepo := repository.NewDmUserRepository(groupManager)
	ctx := context.Background()

	// ユニークなメールアドレスを生成
	uniqueID, err := idgen.GenerateUUIDv7()
	require.NoError(t, err)
	uniqueEmail := fmt.Sprintf("test-%s@example.com", uniqueID)

	req := &model.CreateDmUserRequest{
		Name:  "Test User",
		Email: uniqueEmail,
	}

	dmUser, err := dmUserRepo.Create(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, dmUser)

	// クリーンアップ
	defer func() {
		if dmUser != nil {
			_ = dmUserRepo.Delete(ctx, dmUser.ID)
		}
	}()

	assert.NotZero(t, dmUser.ID)
	assert.Equal(t, "Test User", dmUser.Name)
	assert.Equal(t, uniqueEmail, dmUser.Email)
	assert.NotZero(t, dmUser.CreatedAt)
	assert.NotZero(t, dmUser.UpdatedAt)
}

func TestDmUserRepositoryGORM_GetByID(t *testing.T) {
	groupManager := testutil.SetupTestGroupManager(t, 4, 8)
	defer testutil.CleanupTestGroupManager(groupManager)

	dmUserRepo := repository.NewDmUserRepository(groupManager)
	ctx := context.Background()

	// ユニークなメールアドレスを生成
	uniqueID, err := idgen.GenerateUUIDv7()
	require.NoError(t, err)
	uniqueEmail := fmt.Sprintf("test-%s@example.com", uniqueID)

	// Create test user first
	req := &model.CreateDmUserRequest{
		Name:  "Test User",
		Email: uniqueEmail,
	}
	created, err := dmUserRepo.Create(ctx, req)
	require.NoError(t, err)

	// クリーンアップ
	defer func() {
		_ = dmUserRepo.Delete(ctx, created.ID)
	}()

	// Test retrieval
	dmUser, err := dmUserRepo.GetByID(ctx, created.ID)
	assert.NoError(t, err)
	assert.NotNil(t, dmUser)
	assert.Equal(t, created.ID, dmUser.ID)
	assert.Equal(t, "Test User", dmUser.Name)
	assert.Equal(t, uniqueEmail, dmUser.Email)
}

func TestDmUserRepositoryGORM_GetByID_NotFound(t *testing.T) {
	groupManager := testutil.SetupTestGroupManager(t, 4, 8)
	defer testutil.CleanupTestGroupManager(groupManager)

	dmUserRepo := repository.NewDmUserRepository(groupManager)
	ctx := context.Background()

	// Test retrieval of non-existent user
	dmUser, err := dmUserRepo.GetByID(ctx, "00000000000000000000000000000000")
	assert.Error(t, err)
	assert.Nil(t, dmUser)
}

func TestDmUserRepositoryGORM_Update(t *testing.T) {
	groupManager := testutil.SetupTestGroupManager(t, 4, 8)
	defer testutil.CleanupTestGroupManager(groupManager)

	dmUserRepo := repository.NewDmUserRepository(groupManager)
	ctx := context.Background()

	// ユニークなメールアドレスを生成
	uniqueID, err := idgen.GenerateUUIDv7()
	require.NoError(t, err)
	originalEmail := fmt.Sprintf("original-%s@example.com", uniqueID)
	updatedEmail := fmt.Sprintf("updated-%s@example.com", uniqueID)

	// Create test user first
	createReq := &model.CreateDmUserRequest{
		Name:  "Original Name",
		Email: originalEmail,
	}
	created, err := dmUserRepo.Create(ctx, createReq)
	require.NoError(t, err)

	// クリーンアップ
	defer func() {
		_ = dmUserRepo.Delete(ctx, created.ID)
	}()

	// Update user
	updateReq := &model.UpdateDmUserRequest{
		Name:  "Updated Name",
		Email: updatedEmail,
	}
	updated, err := dmUserRepo.Update(ctx, created.ID, updateReq)
	assert.NoError(t, err)
	assert.NotNil(t, updated)
	assert.Equal(t, "Updated Name", updated.Name)
	assert.Equal(t, updatedEmail, updated.Email)

	// Verify update
	dmUser, err := dmUserRepo.GetByID(ctx, created.ID)
	require.NoError(t, err)
	assert.Equal(t, "Updated Name", dmUser.Name)
	assert.Equal(t, updatedEmail, dmUser.Email)
}

func TestDmUserRepositoryGORM_Delete(t *testing.T) {
	groupManager := testutil.SetupTestGroupManager(t, 4, 8)
	defer testutil.CleanupTestGroupManager(groupManager)

	dmUserRepo := repository.NewDmUserRepository(groupManager)
	ctx := context.Background()

	// ユニークなメールアドレスを生成
	uniqueID, err := idgen.GenerateUUIDv7()
	require.NoError(t, err)
	uniqueEmail := fmt.Sprintf("test-%s@example.com", uniqueID)

	// Create test user first
	req := &model.CreateDmUserRequest{
		Name:  "Test User",
		Email: uniqueEmail,
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

func TestDmUserRepositoryGORM_CreateAndRetrieve(t *testing.T) {
	groupManager := testutil.SetupTestGroupManager(t, 4, 8)
	defer testutil.CleanupTestGroupManager(groupManager)

	dmUserRepo := repository.NewDmUserRepository(groupManager)
	ctx := context.Background()

	// ユニークなメールアドレスを生成
	uniqueID1, err := idgen.GenerateUUIDv7()
	require.NoError(t, err)
	uniqueID2, err := idgen.GenerateUUIDv7()
	require.NoError(t, err)

	// Create test users
	req1 := &model.CreateDmUserRequest{
		Name:  "User 1",
		Email: fmt.Sprintf("user1-%s@example.com", uniqueID1),
	}
	user1, err := dmUserRepo.Create(ctx, req1)
	require.NoError(t, err)

	req2 := &model.CreateDmUserRequest{
		Name:  "User 2",
		Email: fmt.Sprintf("user2-%s@example.com", uniqueID2),
	}
	user2, err := dmUserRepo.Create(ctx, req2)
	require.NoError(t, err)

	// クリーンアップ
	defer func() {
		_ = dmUserRepo.Delete(ctx, user1.ID)
		_ = dmUserRepo.Delete(ctx, user2.ID)
	}()

	// Verify users can be retrieved by ID (single shard queries)
	retrieved1, err := dmUserRepo.GetByID(ctx, user1.ID)
	assert.NoError(t, err)
	assert.Equal(t, user1.ID, retrieved1.ID)
	assert.Equal(t, "User 1", retrieved1.Name)

	retrieved2, err := dmUserRepo.GetByID(ctx, user2.ID)
	assert.NoError(t, err)
	assert.Equal(t, user2.ID, retrieved2.ID)
	assert.Equal(t, "User 2", retrieved2.Name)
}
