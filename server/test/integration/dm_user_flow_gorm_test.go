package integration_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/taku-o/go-webdb-template/internal/model"
	"github.com/taku-o/go-webdb-template/internal/repository"
	"github.com/taku-o/go-webdb-template/test/testutil"
)

func TestDmUserCRUDFlowGORM(t *testing.T) {
	// Setup test database with GroupManager
	groupManager := testutil.SetupTestGroupManager(t, 4, 8)
	defer testutil.CleanupTestGroupManager(groupManager)

	// Initialize GORM repository
	dmUserRepo := repository.NewDmUserRepositoryGORM(groupManager)
	ctx := context.Background()

	// Test Create
	t.Run("Create DmUser", func(t *testing.T) {
		createReq := &model.CreateDmUserRequest{
			Name:  "Integration Test User GORM",
			Email: "integration.gorm@example.com",
		}
		dmUser, err := dmUserRepo.Create(ctx, createReq)
		require.NoError(t, err)
		assert.NotZero(t, dmUser.ID)
		assert.Equal(t, "Integration Test User GORM", dmUser.Name)
		assert.Equal(t, "integration.gorm@example.com", dmUser.Email)
		assert.NotZero(t, dmUser.CreatedAt)
		assert.NotZero(t, dmUser.UpdatedAt)

		// Test Read
		t.Run("Get DmUser by ID", func(t *testing.T) {
			retrieved, err := dmUserRepo.GetByID(ctx, dmUser.ID)
			require.NoError(t, err)
			assert.Equal(t, dmUser.ID, retrieved.ID)
			assert.Equal(t, dmUser.Name, retrieved.Name)
			assert.Equal(t, dmUser.Email, retrieved.Email)
		})

		// Test Update
		t.Run("Update DmUser", func(t *testing.T) {
			updateReq := &model.UpdateDmUserRequest{
				Name:  "Updated Name GORM",
				Email: "updated.gorm@example.com",
			}
			updated, err := dmUserRepo.Update(ctx, dmUser.ID, updateReq)
			require.NoError(t, err)
			assert.Equal(t, dmUser.ID, updated.ID)
			assert.Equal(t, "Updated Name GORM", updated.Name)
			assert.Equal(t, "updated.gorm@example.com", updated.Email)

			// Verify update persisted
			retrieved, err := dmUserRepo.GetByID(ctx, dmUser.ID)
			require.NoError(t, err)
			assert.Equal(t, "Updated Name GORM", retrieved.Name)
		})

		// Test Delete
		t.Run("Delete DmUser", func(t *testing.T) {
			err := dmUserRepo.Delete(ctx, dmUser.ID)
			assert.NoError(t, err)

			// Verify deletion
			_, err = dmUserRepo.GetByID(ctx, dmUser.ID)
			assert.Error(t, err)
		})
	})
}

func TestDmUserCrossShardOperationsGORM(t *testing.T) {
	// Setup test database with GroupManager
	groupManager := testutil.SetupTestGroupManager(t, 4, 8)
	defer testutil.CleanupTestGroupManager(groupManager)

	dmUserRepo := repository.NewDmUserRepositoryGORM(groupManager)

	// Create multiple dm_users
	ctx := context.Background()
	dmUser1, err := dmUserRepo.Create(ctx, &model.CreateDmUserRequest{
		Name:  "User 1 GORM",
		Email: "user1.gorm@example.com",
	})
	require.NoError(t, err)

	dmUser2, err := dmUserRepo.Create(ctx, &model.CreateDmUserRequest{
		Name:  "User 2 GORM",
		Email: "user2.gorm@example.com",
	})
	require.NoError(t, err)

	dmUser3, err := dmUserRepo.Create(ctx, &model.CreateDmUserRequest{
		Name:  "User 3 GORM",
		Email: "user3.gorm@example.com",
	})
	require.NoError(t, err)

	// Log dm_user IDs (shard distribution is internal)
	t.Logf("Created DmUser 1 (ID=%d)", dmUser1.ID)
	t.Logf("Created DmUser 2 (ID=%d)", dmUser2.ID)
	t.Logf("Created DmUser 3 (ID=%d)", dmUser3.ID)

	// Test GetAll returns dm_users from all shards
	t.Run("GetAll returns dm_users from all shards", func(t *testing.T) {
		allDmUsers, err := dmUserRepo.List(ctx, 100, 0)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(allDmUsers), 3)

		// Verify all our dm_users are in the result
		dmUserIDs := make(map[int64]bool)
		for _, dmUser := range allDmUsers {
			dmUserIDs[dmUser.ID] = true
		}

		assert.True(t, dmUserIDs[dmUser1.ID], "DmUser 1 should be in results")
		assert.True(t, dmUserIDs[dmUser2.ID], "DmUser 2 should be in results")
		assert.True(t, dmUserIDs[dmUser3.ID], "DmUser 3 should be in results")
	})
}
