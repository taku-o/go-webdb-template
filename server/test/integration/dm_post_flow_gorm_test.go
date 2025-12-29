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

func TestPostCRUDFlowGORM(t *testing.T) {
	// Setup test database with GroupManager
	groupManager := testutil.SetupTestGroupManager(t, 4, 8)
	defer testutil.CleanupTestGroupManager(groupManager)

	// Initialize GORM repositories
	dmUserRepo := repository.NewDmUserRepositoryGORM(groupManager)
	dmPostRepo := repository.NewDmPostRepositoryGORM(groupManager)

	// Create a test user first
	ctx := context.Background()
	dmUser, err := dmUserRepo.Create(ctx, &model.CreateDmUserRequest{
		Name:  "PostTestUser GORM",
		Email: "posttest.gorm@example.com",
	})
	require.NoError(t, err)

	// Test Create Post
	t.Run("Create Post", func(t *testing.T) {
		createReq := &model.CreateDmPostRequest{
			UserID:  dmUser.ID,
			Title:   "Integration Test Post GORM",
			Content: "This is a GORM test post content",
		}
		dmPost, err := dmPostRepo.Create(ctx, createReq)
		require.NoError(t, err)
		assert.NotZero(t, dmPost.ID)
		assert.Equal(t, dmUser.ID, dmPost.UserID)
		assert.Equal(t, "Integration Test Post GORM", dmPost.Title)
		assert.Equal(t, "This is a GORM test post content", dmPost.Content)

		// Test Read
		t.Run("Get Post by ID", func(t *testing.T) {
			retrieved, err := dmPostRepo.GetByID(ctx, dmPost.ID, dmUser.ID)
			require.NoError(t, err)
			assert.Equal(t, dmPost.ID, retrieved.ID)
			assert.Equal(t, dmPost.Title, retrieved.Title)
		})

		// Test Update
		t.Run("Update Post", func(t *testing.T) {
			updateReq := &model.UpdateDmPostRequest{
				Title:   "Updated Title GORM",
				Content: "Updated content GORM",
			}
			updated, err := dmPostRepo.Update(ctx, dmPost.ID, dmUser.ID, updateReq)
			require.NoError(t, err)
			assert.Equal(t, "Updated Title GORM", updated.Title)
			assert.Equal(t, "Updated content GORM", updated.Content)
		})

		// Test Delete
		t.Run("Delete Post", func(t *testing.T) {
			err := dmPostRepo.Delete(ctx, dmPost.ID, dmUser.ID)
			assert.NoError(t, err)

			// Verify deletion
			_, err = dmPostRepo.GetByID(ctx, dmPost.ID, dmUser.ID)
			assert.Error(t, err)
		})
	})
}

func TestCrossShardJoinGORM(t *testing.T) {
	// Setup test database with GroupManager
	groupManager := testutil.SetupTestGroupManager(t, 4, 8)
	defer testutil.CleanupTestGroupManager(groupManager)

	// Initialize GORM repositories
	dmUserRepo := repository.NewDmUserRepositoryGORM(groupManager)
	dmPostRepo := repository.NewDmPostRepositoryGORM(groupManager)

	ctx := context.Background()

	// Create multiple users
	dmUser1, err := dmUserRepo.Create(ctx, &model.CreateDmUserRequest{
		Name:  "User1 GORM",
		Email: "user1.gorm@example.com",
	})
	require.NoError(t, err)

	dmUser2, err := dmUserRepo.Create(ctx, &model.CreateDmUserRequest{
		Name:  "User2 GORM",
		Email: "user2.gorm@example.com",
	})
	require.NoError(t, err)

	// Create posts for each user
	dmPost1, err := dmPostRepo.Create(ctx, &model.CreateDmPostRequest{
		UserID:  dmUser1.ID,
		Title:   "User1 Post GORM",
		Content: "Content by User1",
	})
	require.NoError(t, err)

	dmPost2, err := dmPostRepo.Create(ctx, &model.CreateDmPostRequest{
		UserID:  dmUser2.ID,
		Title:   "User2 Post GORM",
		Content: "Content by User2",
	})
	require.NoError(t, err)

	t.Logf("Created User1 (ID=%d)", dmUser1.ID)
	t.Logf("Created User2 (ID=%d)", dmUser2.ID)

	// Test cross-shard JOIN
	t.Run("GetUserPosts returns data from all shards", func(t *testing.T) {
		dmUserPosts, err := dmPostRepo.GetUserPosts(ctx, 100, 0)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(dmUserPosts), 2)

		// Verify data contains our test posts with user info
		foundPost1 := false
		foundPost2 := false

		for _, up := range dmUserPosts {
			if up.PostID == dmPost1.ID {
				assert.Equal(t, dmUser1.ID, up.UserID)
				assert.Equal(t, dmUser1.Name, up.UserName)
				assert.Equal(t, dmPost1.Title, up.PostTitle)
				foundPost1 = true
			}
			if up.PostID == dmPost2.ID {
				assert.Equal(t, dmUser2.ID, up.UserID)
				assert.Equal(t, dmUser2.Name, up.UserName)
				assert.Equal(t, dmPost2.Title, up.PostTitle)
				foundPost2 = true
			}
		}

		assert.True(t, foundPost1, "Should find post 1 with user data")
		assert.True(t, foundPost2, "Should find post 2 with user data")
	})

	// Test cross-shard List
	t.Run("List returns posts from all shards", func(t *testing.T) {
		allDmPosts, err := dmPostRepo.List(ctx, 100, 0)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(allDmPosts), 2)

		dmPostIDs := make(map[int64]bool)
		for _, dmPost := range allDmPosts {
			dmPostIDs[dmPost.ID] = true
		}

		assert.True(t, dmPostIDs[dmPost1.ID], "Post 1 should be in results")
		assert.True(t, dmPostIDs[dmPost2.ID], "Post 2 should be in results")
	})
}
