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
	userRepo := repository.NewUserRepositoryGORM(groupManager)
	postRepo := repository.NewPostRepositoryGORM(groupManager)

	// Create a test user first
	ctx := context.Background()
	user, err := userRepo.Create(ctx, &model.CreateUserRequest{
		Name:  "PostTestUser GORM",
		Email: "posttest.gorm@example.com",
	})
	require.NoError(t, err)

	// Test Create Post
	t.Run("Create Post", func(t *testing.T) {
		createReq := &model.CreatePostRequest{
			UserID:  user.ID,
			Title:   "Integration Test Post GORM",
			Content: "This is a GORM test post content",
		}
		post, err := postRepo.Create(ctx, createReq)
		require.NoError(t, err)
		assert.NotZero(t, post.ID)
		assert.Equal(t, user.ID, post.UserID)
		assert.Equal(t, "Integration Test Post GORM", post.Title)
		assert.Equal(t, "This is a GORM test post content", post.Content)

		// Test Read
		t.Run("Get Post by ID", func(t *testing.T) {
			retrieved, err := postRepo.GetByID(ctx, post.ID, user.ID)
			require.NoError(t, err)
			assert.Equal(t, post.ID, retrieved.ID)
			assert.Equal(t, post.Title, retrieved.Title)
		})

		// Test Update
		t.Run("Update Post", func(t *testing.T) {
			updateReq := &model.UpdatePostRequest{
				Title:   "Updated Title GORM",
				Content: "Updated content GORM",
			}
			updated, err := postRepo.Update(ctx, post.ID, user.ID, updateReq)
			require.NoError(t, err)
			assert.Equal(t, "Updated Title GORM", updated.Title)
			assert.Equal(t, "Updated content GORM", updated.Content)
		})

		// Test Delete
		t.Run("Delete Post", func(t *testing.T) {
			err := postRepo.Delete(ctx, post.ID, user.ID)
			assert.NoError(t, err)

			// Verify deletion
			_, err = postRepo.GetByID(ctx, post.ID, user.ID)
			assert.Error(t, err)
		})
	})
}

func TestCrossShardJoinGORM(t *testing.T) {
	// Setup test database with GroupManager
	groupManager := testutil.SetupTestGroupManager(t, 4, 8)
	defer testutil.CleanupTestGroupManager(groupManager)

	// Initialize GORM repositories
	userRepo := repository.NewUserRepositoryGORM(groupManager)
	postRepo := repository.NewPostRepositoryGORM(groupManager)

	ctx := context.Background()

	// Create multiple users
	user1, err := userRepo.Create(ctx, &model.CreateUserRequest{
		Name:  "User1 GORM",
		Email: "user1.gorm@example.com",
	})
	require.NoError(t, err)

	user2, err := userRepo.Create(ctx, &model.CreateUserRequest{
		Name:  "User2 GORM",
		Email: "user2.gorm@example.com",
	})
	require.NoError(t, err)

	// Create posts for each user
	post1, err := postRepo.Create(ctx, &model.CreatePostRequest{
		UserID:  user1.ID,
		Title:   "User1 Post GORM",
		Content: "Content by User1",
	})
	require.NoError(t, err)

	post2, err := postRepo.Create(ctx, &model.CreatePostRequest{
		UserID:  user2.ID,
		Title:   "User2 Post GORM",
		Content: "Content by User2",
	})
	require.NoError(t, err)

	t.Logf("Created User1 (ID=%d)", user1.ID)
	t.Logf("Created User2 (ID=%d)", user2.ID)

	// Test cross-shard JOIN
	t.Run("GetUserPosts returns data from all shards", func(t *testing.T) {
		userPosts, err := postRepo.GetUserPosts(ctx, 100, 0)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(userPosts), 2)

		// Verify data contains our test posts with user info
		foundPost1 := false
		foundPost2 := false

		for _, up := range userPosts {
			if up.PostID == post1.ID {
				assert.Equal(t, user1.ID, up.UserID)
				assert.Equal(t, user1.Name, up.UserName)
				assert.Equal(t, post1.Title, up.PostTitle)
				foundPost1 = true
			}
			if up.PostID == post2.ID {
				assert.Equal(t, user2.ID, up.UserID)
				assert.Equal(t, user2.Name, up.UserName)
				assert.Equal(t, post2.Title, up.PostTitle)
				foundPost2 = true
			}
		}

		assert.True(t, foundPost1, "Should find post 1 with user data")
		assert.True(t, foundPost2, "Should find post 2 with user data")
	})

	// Test cross-shard List
	t.Run("List returns posts from all shards", func(t *testing.T) {
		allPosts, err := postRepo.List(ctx, 100, 0)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(allPosts), 2)

		postIDs := make(map[int64]bool)
		for _, post := range allPosts {
			postIDs[post.ID] = true
		}

		assert.True(t, postIDs[post1.ID], "Post 1 should be in results")
		assert.True(t, postIDs[post2.ID], "Post 2 should be in results")
	})
}
