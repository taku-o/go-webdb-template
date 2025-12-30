package fixtures

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/taku-o/go-webdb-template/internal/model"
	"github.com/taku-o/go-webdb-template/internal/service"
)

// CreateTestDmPost creates a test dm_post using the service layer
func CreateTestDmPost(t *testing.T, svc *service.DmPostService, userID string, title string) *model.DmPost {
	req := &model.CreateDmPostRequest{
		UserID:  userID,
		Title:   title,
		Content: "Content for " + title,
	}
	dmPost, err := svc.CreateDmPost(context.Background(), req)
	require.NoError(t, err)
	return dmPost
}

// CreateTestDmPostWithContent creates a test dm_post with specific content
func CreateTestDmPostWithContent(t *testing.T, svc *service.DmPostService, userID string, title, content string) *model.DmPost {
	req := &model.CreateDmPostRequest{
		UserID:  userID,
		Title:   title,
		Content: content,
	}
	dmPost, err := svc.CreateDmPost(context.Background(), req)
	require.NoError(t, err)
	return dmPost
}

// CreateMultipleTestDmPosts creates multiple test dm_posts for a dm_user
func CreateMultipleTestDmPosts(t *testing.T, svc *service.DmPostService, userID string, count int) []*model.DmPost {
	dmPosts := make([]*model.DmPost, count)
	for i := 0; i < count; i++ {
		title := fmt.Sprintf("Post %d", i+1)
		dmPosts[i] = CreateTestDmPost(t, svc, userID, title)
	}
	return dmPosts
}
