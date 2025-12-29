package integration_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/taku-o/go-webdb-template/internal/model"
	"github.com/taku-o/go-webdb-template/internal/repository"
	"github.com/taku-o/go-webdb-template/internal/service"
	"github.com/taku-o/go-webdb-template/test/testutil"
)

func TestUserCRUDFlow(t *testing.T) {
	// Setup test database with GroupManager
	groupManager := testutil.SetupTestGroupManager(t, 4, 8)
	defer testutil.CleanupTestGroupManager(groupManager)

	// Initialize repositories and services (using GORM repositories)
	dmUserRepo := repository.NewDmUserRepositoryGORM(groupManager)
	dmUserService := service.NewDmUserService(dmUserRepo)

	// Test Create
	t.Run("Create User", func(t *testing.T) {
		createReq := &model.CreateDmUserRequest{
			Name:  "Integration Test User",
			Email: "integration@example.com",
		}
		dmUser, err := dmUserService.CreateDmUser(context.Background(), createReq)
		require.NoError(t, err)
		assert.NotZero(t, dmUser.ID)
		assert.Equal(t, "Integration Test User", dmUser.Name)
		assert.Equal(t, "integration@example.com", dmUser.Email)
		assert.NotZero(t, dmUser.CreatedAt)
		assert.NotZero(t, dmUser.UpdatedAt)

		// Test Read
		t.Run("Get User by ID", func(t *testing.T) {
			retrieved, err := dmUserService.GetDmUser(context.Background(), dmUser.ID)
			require.NoError(t, err)
			assert.Equal(t, dmUser.ID, retrieved.ID)
			assert.Equal(t, dmUser.Name, retrieved.Name)
			assert.Equal(t, dmUser.Email, retrieved.Email)
		})

		// Test Update
		t.Run("Update User", func(t *testing.T) {
			updateReq := &model.UpdateDmUserRequest{
				Name:  "Updated Name",
				Email: "updated@example.com",
			}
			updated, err := dmUserService.UpdateDmUser(context.Background(), dmUser.ID, updateReq)
			require.NoError(t, err)
			assert.Equal(t, dmUser.ID, updated.ID)
			assert.Equal(t, "Updated Name", updated.Name)
			assert.Equal(t, "updated@example.com", updated.Email)

			// Verify update persisted
			retrieved, err := dmUserService.GetDmUser(context.Background(), dmUser.ID)
			require.NoError(t, err)
			assert.Equal(t, "Updated Name", retrieved.Name)
		})

		// Test Delete
		t.Run("Delete User", func(t *testing.T) {
			err := dmUserService.DeleteDmUser(context.Background(), dmUser.ID)
			assert.NoError(t, err)

			// Verify deletion
			_, err = dmUserService.GetDmUser(context.Background(), dmUser.ID)
			assert.Error(t, err)
		})
	})
}

func TestUserCrossShardOperations(t *testing.T) {
	// Setup test database with GroupManager
	groupManager := testutil.SetupTestGroupManager(t, 4, 8)
	defer testutil.CleanupTestGroupManager(groupManager)

	dmUserRepo := repository.NewDmUserRepositoryGORM(groupManager)
	dmUserService := service.NewDmUserService(dmUserRepo)

	// Create multiple users
	ctx := context.Background()
	dmUser1, err := dmUserService.CreateDmUser(ctx, &model.CreateDmUserRequest{
		Name:  "User 1",
		Email: "user1@example.com",
	})
	require.NoError(t, err)

	dmUser2, err := dmUserService.CreateDmUser(ctx, &model.CreateDmUserRequest{
		Name:  "User 2",
		Email: "user2@example.com",
	})
	require.NoError(t, err)

	dmUser3, err := dmUserService.CreateDmUser(ctx, &model.CreateDmUserRequest{
		Name:  "User 3",
		Email: "user3@example.com",
	})
	require.NoError(t, err)

	// Log user IDs (shard distribution is internal)
	t.Logf("Created User 1 (ID=%d)", dmUser1.ID)
	t.Logf("Created User 2 (ID=%d)", dmUser2.ID)
	t.Logf("Created User 3 (ID=%d)", dmUser3.ID)

	// Test GetAll returns users from all shards
	t.Run("GetAll returns users from all shards", func(t *testing.T) {
		allDmUsers, err := dmUserService.ListDmUsers(ctx, 100, 0)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(allDmUsers), 3)

		// Verify all our users are in the result
		dmUserIDs := make(map[int64]bool)
		for _, dmUser := range allDmUsers {
			dmUserIDs[dmUser.ID] = true
		}

		assert.True(t, dmUserIDs[dmUser1.ID], "User 1 should be in results")
		assert.True(t, dmUserIDs[dmUser2.ID], "User 2 should be in results")
		assert.True(t, dmUserIDs[dmUser3.ID], "User 3 should be in results")
	})
}
