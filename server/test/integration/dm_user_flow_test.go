package integration_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/taku-o/go-webdb-template/internal/model"
	"github.com/taku-o/go-webdb-template/internal/repository"
	"github.com/taku-o/go-webdb-template/internal/service"
	"github.com/taku-o/go-webdb-template/internal/util/idgen"
	"github.com/taku-o/go-webdb-template/test/testutil"
)

func TestDmUserCRUDFlow(t *testing.T) {
	// Setup test database with GroupManager
	groupManager := testutil.SetupTestGroupManager(t, 4, 8)
	defer testutil.CleanupTestGroupManager(groupManager)

	// Initialize repositories and services (using GORM repositories)
	dmUserRepo := repository.NewDmUserRepositoryGORM(groupManager)
	dmUserService := service.NewDmUserService(dmUserRepo)

	// ユニークなメールアドレスを生成
	uniqueID, err := idgen.GenerateUUIDv7()
	require.NoError(t, err)
	uniqueEmail := fmt.Sprintf("integration-%s@example.com", uniqueID)
	updatedEmail := fmt.Sprintf("updated-%s@example.com", uniqueID)

	// Test Create
	t.Run("Create DmUser", func(t *testing.T) {
		createReq := &model.CreateDmUserRequest{
			Name:  "Integration Test User",
			Email: uniqueEmail,
		}
		dmUser, err := dmUserService.CreateDmUser(context.Background(), createReq)
		require.NoError(t, err)
		assert.NotZero(t, dmUser.ID)
		assert.Equal(t, "Integration Test User", dmUser.Name)
		assert.Equal(t, uniqueEmail, dmUser.Email)
		assert.NotZero(t, dmUser.CreatedAt)
		assert.NotZero(t, dmUser.UpdatedAt)

		// Test Read
		t.Run("Get DmUser by ID", func(t *testing.T) {
			retrieved, err := dmUserService.GetDmUser(context.Background(), dmUser.ID)
			require.NoError(t, err)
			assert.Equal(t, dmUser.ID, retrieved.ID)
			assert.Equal(t, dmUser.Name, retrieved.Name)
			assert.Equal(t, dmUser.Email, retrieved.Email)
		})

		// Test Update
		t.Run("Update DmUser", func(t *testing.T) {
			updateReq := &model.UpdateDmUserRequest{
				Name:  "Updated Name",
				Email: updatedEmail,
			}
			updated, err := dmUserService.UpdateDmUser(context.Background(), dmUser.ID, updateReq)
			require.NoError(t, err)
			assert.Equal(t, dmUser.ID, updated.ID)
			assert.Equal(t, "Updated Name", updated.Name)
			assert.Equal(t, updatedEmail, updated.Email)

			// Verify update persisted
			retrieved, err := dmUserService.GetDmUser(context.Background(), dmUser.ID)
			require.NoError(t, err)
			assert.Equal(t, "Updated Name", retrieved.Name)
		})

		// Test Delete
		t.Run("Delete DmUser", func(t *testing.T) {
			err := dmUserService.DeleteDmUser(context.Background(), dmUser.ID)
			assert.NoError(t, err)

			// Verify deletion
			_, err = dmUserService.GetDmUser(context.Background(), dmUser.ID)
			assert.Error(t, err)
		})
	})
}

func TestDmUserCrossShardOperations(t *testing.T) {
	// Setup test database with GroupManager
	groupManager := testutil.SetupTestGroupManager(t, 4, 8)
	defer testutil.CleanupTestGroupManager(groupManager)

	dmUserRepo := repository.NewDmUserRepositoryGORM(groupManager)
	dmUserService := service.NewDmUserService(dmUserRepo)

	// ユニークなメールアドレスを生成
	uniqueID1, err := idgen.GenerateUUIDv7()
	require.NoError(t, err)
	uniqueID2, err := idgen.GenerateUUIDv7()
	require.NoError(t, err)
	uniqueID3, err := idgen.GenerateUUIDv7()
	require.NoError(t, err)

	// Create multiple dm_users
	ctx := context.Background()
	dmUser1, err := dmUserService.CreateDmUser(ctx, &model.CreateDmUserRequest{
		Name:  "User 1",
		Email: fmt.Sprintf("user1-%s@example.com", uniqueID1),
	})
	require.NoError(t, err)

	dmUser2, err := dmUserService.CreateDmUser(ctx, &model.CreateDmUserRequest{
		Name:  "User 2",
		Email: fmt.Sprintf("user2-%s@example.com", uniqueID2),
	})
	require.NoError(t, err)

	dmUser3, err := dmUserService.CreateDmUser(ctx, &model.CreateDmUserRequest{
		Name:  "User 3",
		Email: fmt.Sprintf("user3-%s@example.com", uniqueID3),
	})
	require.NoError(t, err)

	// クリーンアップ
	defer func() {
		_ = dmUserService.DeleteDmUser(ctx, dmUser1.ID)
		_ = dmUserService.DeleteDmUser(ctx, dmUser2.ID)
		_ = dmUserService.DeleteDmUser(ctx, dmUser3.ID)
	}()

	// Log dm_user IDs (shard distribution is internal)
	t.Logf("Created DmUser 1 (ID=%s)", dmUser1.ID)
	t.Logf("Created DmUser 2 (ID=%s)", dmUser2.ID)
	t.Logf("Created DmUser 3 (ID=%s)", dmUser3.ID)

	// Test GetAll returns dm_users from all shards
	t.Run("GetAll returns dm_users from all shards", func(t *testing.T) {
		allDmUsers, err := dmUserService.ListDmUsers(ctx, 100, 0)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(allDmUsers), 3)

		// Verify all our dm_users are in the result
		dmUserIDs := make(map[string]bool)
		for _, dmUser := range allDmUsers {
			dmUserIDs[dmUser.ID] = true
		}

		assert.True(t, dmUserIDs[dmUser1.ID], "DmUser 1 should be in results")
		assert.True(t, dmUserIDs[dmUser2.ID], "DmUser 2 should be in results")
		assert.True(t, dmUserIDs[dmUser3.ID], "DmUser 3 should be in results")
	})
}
