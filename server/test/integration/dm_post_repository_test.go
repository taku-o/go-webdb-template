package integration_test

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

func TestDmPostRepository_CRUDFlow(t *testing.T) {
	// Setup test database with GroupManager
	groupManager := testutil.SetupTestGroupManager(t, 4, 8)
	defer testutil.CleanupTestGroupManager(groupManager)

	// Initialize repositories
	dmUserRepo := repository.NewDmUserRepository(groupManager)
	dmPostRepo := repository.NewDmPostRepository(groupManager)

	// ユニークなメールアドレスを生成
	uniqueID, err := idgen.GenerateUUIDv7()
	require.NoError(t, err)
	uniqueEmail := fmt.Sprintf("posttest-gorm-%s@example.com", uniqueID)

	// Create a test dm_user first
	ctx := context.Background()
	dmUser, err := dmUserRepo.Create(ctx, &model.CreateDmUserRequest{
		Name:  "PostTestUser GORM",
		Email: uniqueEmail,
	})
	require.NoError(t, err)

	// クリーンアップ
	defer func() {
		_ = dmUserRepo.Delete(ctx, dmUser.ID)
	}()

	// Test Create DmPost
	t.Run("Create DmPost", func(t *testing.T) {
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
		t.Run("Get DmPost by ID", func(t *testing.T) {
			retrieved, err := dmPostRepo.GetByID(ctx, dmPost.ID, dmUser.ID)
			require.NoError(t, err)
			assert.Equal(t, dmPost.ID, retrieved.ID)
			assert.Equal(t, dmPost.Title, retrieved.Title)
		})

		// Test Update
		t.Run("Update DmPost", func(t *testing.T) {
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
		t.Run("Delete DmPost", func(t *testing.T) {
			err := dmPostRepo.Delete(ctx, dmPost.ID, dmUser.ID)
			assert.NoError(t, err)

			// Verify deletion
			_, err = dmPostRepo.GetByID(ctx, dmPost.ID, dmUser.ID)
			assert.Error(t, err)
		})
	})
}

func TestDmPostRepository_CrossShardJoin(t *testing.T) {
	// Setup test database with GroupManager
	groupManager := testutil.SetupTestGroupManager(t, 4, 8)
	defer testutil.CleanupTestGroupManager(groupManager)

	// Initialize repositories
	dmUserRepo := repository.NewDmUserRepository(groupManager)
	dmPostRepo := repository.NewDmPostRepository(groupManager)

	ctx := context.Background()

	// ユニークなメールアドレスを生成
	uniqueID1, err := idgen.GenerateUUIDv7()
	require.NoError(t, err)
	uniqueID2, err := idgen.GenerateUUIDv7()
	require.NoError(t, err)

	// Create multiple dm_users
	dmUser1, err := dmUserRepo.Create(ctx, &model.CreateDmUserRequest{
		Name:  "User1 GORM",
		Email: fmt.Sprintf("user1-gorm-%s@example.com", uniqueID1),
	})
	require.NoError(t, err)

	dmUser2, err := dmUserRepo.Create(ctx, &model.CreateDmUserRequest{
		Name:  "User2 GORM",
		Email: fmt.Sprintf("user2-gorm-%s@example.com", uniqueID2),
	})
	require.NoError(t, err)

	// Create dm_posts for each dm_user
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

	// クリーンアップ
	defer func() {
		_ = dmPostRepo.Delete(ctx, dmPost1.ID, dmPost1.UserID)
		_ = dmPostRepo.Delete(ctx, dmPost2.ID, dmPost2.UserID)
		_ = dmUserRepo.Delete(ctx, dmUser1.ID)
		_ = dmUserRepo.Delete(ctx, dmUser2.ID)
	}()

	t.Logf("Created DmUser1 (ID=%s)", dmUser1.ID)
	t.Logf("Created DmUser2 (ID=%s)", dmUser2.ID)

	// Test cross-shard JOIN
	t.Run("GetDmUserPosts returns data from all shards", func(t *testing.T) {
		dmUserPosts, err := dmPostRepo.GetUserPosts(ctx, 100, 0)
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

	// Test cross-shard List
	t.Run("List returns dm_posts from all shards", func(t *testing.T) {
		allDmPosts, err := dmPostRepo.List(ctx, 100, 0)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(allDmPosts), 2)

		dmPostIDs := make(map[string]bool)
		for _, dmPost := range allDmPosts {
			dmPostIDs[dmPost.ID] = true
		}

		assert.True(t, dmPostIDs[dmPost1.ID], "DmPost 1 should be in results")
		assert.True(t, dmPostIDs[dmPost2.ID], "DmPost 2 should be in results")
	})
}
