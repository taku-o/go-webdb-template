package integration_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/example/go-db-prj-sample/internal/model"
	"github.com/example/go-db-prj-sample/internal/repository"
	"github.com/example/go-db-prj-sample/internal/service"
	"github.com/example/go-db-prj-sample/test/fixtures"
	"github.com/example/go-db-prj-sample/test/testutil"
)

func TestPostCRUDFlow(t *testing.T) {
	// Setup test database with sharding
	dbManager := testutil.SetupTestShards(t, 2)
	defer testutil.CleanupTestDB(dbManager)

	// Initialize services
	userRepo := repository.NewUserRepository(dbManager)
	userService := service.NewUserService(userRepo)
	postRepo := repository.NewPostRepository(dbManager)
	postService := service.NewPostService(postRepo, userRepo)

	// Create a test user first
	user := fixtures.CreateTestUser(t, userService, "PostTestUser")

	// Test Create Post
	ctx := context.Background()
	t.Run("Create Post", func(t *testing.T) {
		createReq := &model.CreatePostRequest{
			UserID:  user.ID,
			Title:   "Integration Test Post",
			Content: "This is a test post content",
		}
		post, err := postService.CreatePost(ctx, createReq)
		require.NoError(t, err)
		assert.NotZero(t, post.ID)
		assert.Equal(t, user.ID, post.UserID)
		assert.Equal(t, "Integration Test Post", post.Title)
		assert.Equal(t, "This is a test post content", post.Content)

		// Test Read
		t.Run("Get Post by ID", func(t *testing.T) {
			retrieved, err := postService.GetPost(ctx, post.ID, user.ID)
			require.NoError(t, err)
			assert.Equal(t, post.ID, retrieved.ID)
			assert.Equal(t, post.Title, retrieved.Title)
		})

		// Test Update
		t.Run("Update Post", func(t *testing.T) {
			updateReq := &model.UpdatePostRequest{
				Title:   "Updated Title",
				Content: "Updated content",
			}
			updated, err := postService.UpdatePost(ctx, post.ID, user.ID, updateReq)
			require.NoError(t, err)
			assert.Equal(t, "Updated Title", updated.Title)
			assert.Equal(t, "Updated content", updated.Content)
		})

		// Test Delete
		t.Run("Delete Post", func(t *testing.T) {
			err := postService.DeletePost(ctx, post.ID, user.ID)
			assert.NoError(t, err)

			// Verify deletion
			_, err = postService.GetPost(ctx, post.ID, user.ID)
			assert.Error(t, err)
		})
	})
}

func TestCrossShardJoin(t *testing.T) {
	// Setup test database with sharding
	dbManager := testutil.SetupTestShards(t, 2)
	defer testutil.CleanupTestDB(dbManager)

	// Initialize services
	userRepo := repository.NewUserRepository(dbManager)
	userService := service.NewUserService(userRepo)
	postRepo := repository.NewPostRepository(dbManager)
	postService := service.NewPostService(postRepo, userRepo)

	// Create multiple users
	user1 := fixtures.CreateTestUser(t, userService, "User1")
	user2 := fixtures.CreateTestUser(t, userService, "User2")

	// Create posts for each user
	post1 := fixtures.CreateTestPost(t, postService, user1.ID, "User1 Post")
	post2 := fixtures.CreateTestPost(t, postService, user2.ID, "User2 Post")

	t.Logf("Created User1 (ID=%d)", user1.ID)
	t.Logf("Created User2 (ID=%d)", user2.ID)

	// Test cross-shard JOIN
	ctx := context.Background()
	t.Run("GetUserPosts returns data from all shards", func(t *testing.T) {
		userPosts, err := postService.GetUserPosts(ctx, 100, 0)
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
}
