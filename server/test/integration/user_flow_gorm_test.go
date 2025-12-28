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

func TestUserCRUDFlowGORM(t *testing.T) {
	// Setup test database with GroupManager
	groupManager := testutil.SetupTestGroupManager(t, 4, 8)
	defer testutil.CleanupTestGroupManager(groupManager)

	// Initialize GORM repository
	userRepo := repository.NewUserRepositoryGORM(groupManager)
	ctx := context.Background()

	// Test Create
	t.Run("Create User", func(t *testing.T) {
		createReq := &model.CreateUserRequest{
			Name:  "Integration Test User GORM",
			Email: "integration.gorm@example.com",
		}
		user, err := userRepo.Create(ctx, createReq)
		require.NoError(t, err)
		assert.NotZero(t, user.ID)
		assert.Equal(t, "Integration Test User GORM", user.Name)
		assert.Equal(t, "integration.gorm@example.com", user.Email)
		assert.NotZero(t, user.CreatedAt)
		assert.NotZero(t, user.UpdatedAt)

		// Test Read
		t.Run("Get User by ID", func(t *testing.T) {
			retrieved, err := userRepo.GetByID(ctx, user.ID)
			require.NoError(t, err)
			assert.Equal(t, user.ID, retrieved.ID)
			assert.Equal(t, user.Name, retrieved.Name)
			assert.Equal(t, user.Email, retrieved.Email)
		})

		// Test Update
		t.Run("Update User", func(t *testing.T) {
			updateReq := &model.UpdateUserRequest{
				Name:  "Updated Name GORM",
				Email: "updated.gorm@example.com",
			}
			updated, err := userRepo.Update(ctx, user.ID, updateReq)
			require.NoError(t, err)
			assert.Equal(t, user.ID, updated.ID)
			assert.Equal(t, "Updated Name GORM", updated.Name)
			assert.Equal(t, "updated.gorm@example.com", updated.Email)

			// Verify update persisted
			retrieved, err := userRepo.GetByID(ctx, user.ID)
			require.NoError(t, err)
			assert.Equal(t, "Updated Name GORM", retrieved.Name)
		})

		// Test Delete
		t.Run("Delete User", func(t *testing.T) {
			err := userRepo.Delete(ctx, user.ID)
			assert.NoError(t, err)

			// Verify deletion
			_, err = userRepo.GetByID(ctx, user.ID)
			assert.Error(t, err)
		})
	})
}

func TestUserCrossShardOperationsGORM(t *testing.T) {
	// Setup test database with GroupManager
	groupManager := testutil.SetupTestGroupManager(t, 4, 8)
	defer testutil.CleanupTestGroupManager(groupManager)

	userRepo := repository.NewUserRepositoryGORM(groupManager)

	// Create multiple users
	ctx := context.Background()
	user1, err := userRepo.Create(ctx, &model.CreateUserRequest{
		Name:  "User 1 GORM",
		Email: "user1.gorm@example.com",
	})
	require.NoError(t, err)

	user2, err := userRepo.Create(ctx, &model.CreateUserRequest{
		Name:  "User 2 GORM",
		Email: "user2.gorm@example.com",
	})
	require.NoError(t, err)

	user3, err := userRepo.Create(ctx, &model.CreateUserRequest{
		Name:  "User 3 GORM",
		Email: "user3.gorm@example.com",
	})
	require.NoError(t, err)

	// Log user IDs (shard distribution is internal)
	t.Logf("Created User 1 (ID=%d)", user1.ID)
	t.Logf("Created User 2 (ID=%d)", user2.ID)
	t.Logf("Created User 3 (ID=%d)", user3.ID)

	// Test GetAll returns users from all shards
	t.Run("GetAll returns users from all shards", func(t *testing.T) {
		allUsers, err := userRepo.List(ctx, 100, 0)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(allUsers), 3)

		// Verify all our users are in the result
		userIDs := make(map[int64]bool)
		for _, user := range allUsers {
			userIDs[user.ID] = true
		}

		assert.True(t, userIDs[user1.ID], "User 1 should be in results")
		assert.True(t, userIDs[user2.ID], "User 2 should be in results")
		assert.True(t, userIDs[user3.ID], "User 3 should be in results")
	})
}
