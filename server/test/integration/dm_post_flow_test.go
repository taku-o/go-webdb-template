package integration_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/taku-o/go-webdb-template/internal/model"
	"github.com/taku-o/go-webdb-template/internal/repository"
	"github.com/taku-o/go-webdb-template/internal/service"
	"github.com/taku-o/go-webdb-template/test/fixtures"
	"github.com/taku-o/go-webdb-template/test/testutil"
)

func TestDmPostCRUDFlow(t *testing.T) {
	// Setup test database with GroupManager
	groupManager := testutil.SetupTestGroupManager(t, 4, 8)
	defer testutil.CleanupTestGroupManager(groupManager)

	// Initialize services (using GORM repositories)
	dmUserRepo := repository.NewDmUserRepositoryGORM(groupManager)
	dmUserService := service.NewDmUserService(dmUserRepo)
	dmPostRepo := repository.NewDmPostRepositoryGORM(groupManager)
	dmPostService := service.NewDmPostService(dmPostRepo, dmUserRepo)

	// Create a test dm_user first
	dmUser := fixtures.CreateTestDmUser(t, dmUserService, "PostTestUser")

	// Test Create DmPost
	ctx := context.Background()
	t.Run("Create DmPost", func(t *testing.T) {
		createReq := &model.CreateDmPostRequest{
			UserID:  dmUser.ID,
			Title:   "Integration Test Post",
			Content: "This is a test post content",
		}
		dmPost, err := dmPostService.CreateDmPost(ctx, createReq)
		require.NoError(t, err)
		assert.NotZero(t, dmPost.ID)
		assert.Equal(t, dmUser.ID, dmPost.UserID)
		assert.Equal(t, "Integration Test Post", dmPost.Title)
		assert.Equal(t, "This is a test post content", dmPost.Content)

		// Test Read
		t.Run("Get DmPost by ID", func(t *testing.T) {
			retrieved, err := dmPostService.GetDmPost(ctx, dmPost.ID, dmUser.ID)
			require.NoError(t, err)
			assert.Equal(t, dmPost.ID, retrieved.ID)
			assert.Equal(t, dmPost.Title, retrieved.Title)
		})

		// Test Update
		t.Run("Update DmPost", func(t *testing.T) {
			updateReq := &model.UpdateDmPostRequest{
				Title:   "Updated Title",
				Content: "Updated content",
			}
			updated, err := dmPostService.UpdateDmPost(ctx, dmPost.ID, dmUser.ID, updateReq)
			require.NoError(t, err)
			assert.Equal(t, "Updated Title", updated.Title)
			assert.Equal(t, "Updated content", updated.Content)
		})

		// Test Delete
		t.Run("Delete DmPost", func(t *testing.T) {
			err := dmPostService.DeleteDmPost(ctx, dmPost.ID, dmUser.ID)
			assert.NoError(t, err)

			// Verify deletion
			_, err = dmPostService.GetDmPost(ctx, dmPost.ID, dmUser.ID)
			assert.Error(t, err)
		})
	})
}

func TestDmPostCrossShardJoin(t *testing.T) {
	// Setup test database with GroupManager
	groupManager := testutil.SetupTestGroupManager(t, 4, 8)
	defer testutil.CleanupTestGroupManager(groupManager)

	// Initialize services (using GORM repositories)
	dmUserRepo := repository.NewDmUserRepositoryGORM(groupManager)
	dmUserService := service.NewDmUserService(dmUserRepo)
	dmPostRepo := repository.NewDmPostRepositoryGORM(groupManager)
	dmPostService := service.NewDmPostService(dmPostRepo, dmUserRepo)

	// Create multiple dm_users
	dmUser1 := fixtures.CreateTestDmUser(t, dmUserService, "User1")
	dmUser2 := fixtures.CreateTestDmUser(t, dmUserService, "User2")

	// Create dm_posts for each dm_user
	dmPost1 := fixtures.CreateTestDmPost(t, dmPostService, dmUser1.ID, "User1 Post")
	dmPost2 := fixtures.CreateTestDmPost(t, dmPostService, dmUser2.ID, "User2 Post")

	t.Logf("Created DmUser1 (ID=%s)", dmUser1.ID)
	t.Logf("Created DmUser2 (ID=%s)", dmUser2.ID)

	// Test cross-shard JOIN
	ctx := context.Background()
	t.Run("GetDmUserPosts returns data from all shards", func(t *testing.T) {
		dmUserPosts, err := dmPostService.GetDmUserPosts(ctx, 100, 0)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(dmUserPosts), 2)

		// Verify data contains our test dm_posts with dm_user info
		foundDmPost1 := false
		foundDmPost2 := false

		for _, up := range dmUserPosts {
			if up.PostID == dmPost1.ID {
				assert.Equal(t, dmUser1.ID, up.UserID)
				assert.Equal(t, dmUser1.Name, up.UserName)
				assert.Equal(t, dmPost1.Title, up.PostTitle)
				foundDmPost1 = true
			}
			if up.PostID == dmPost2.ID {
				assert.Equal(t, dmUser2.ID, up.UserID)
				assert.Equal(t, dmUser2.Name, up.UserName)
				assert.Equal(t, dmPost2.Title, up.PostTitle)
				foundDmPost2 = true
			}
		}

		assert.True(t, foundDmPost1, "Should find dm_post 1 with dm_user data")
		assert.True(t, foundDmPost2, "Should find dm_post 2 with dm_user data")
	})
}
