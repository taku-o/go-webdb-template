package fixtures

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/taku-o/go-webdb-template/internal/model"
	"github.com/taku-o/go-webdb-template/internal/service"
)

// CreateTestPost creates a test post using the service layer
func CreateTestPost(t *testing.T, svc *service.DmPostService, userID int64, title string) *model.DmPost {
	req := &model.CreateDmPostRequest{
		UserID:  userID,
		Title:   title,
		Content: "Content for " + title,
	}
	post, err := svc.CreateDmPost(context.Background(), req)
	require.NoError(t, err)
	return post
}

// CreateTestPostWithContent creates a test post with specific content
func CreateTestPostWithContent(t *testing.T, svc *service.DmPostService, userID int64, title, content string) *model.DmPost {
	req := &model.CreateDmPostRequest{
		UserID:  userID,
		Title:   title,
		Content: content,
	}
	post, err := svc.CreateDmPost(context.Background(), req)
	require.NoError(t, err)
	return post
}

// CreateMultipleTestPosts creates multiple test posts for a user
func CreateMultipleTestPosts(t *testing.T, svc *service.DmPostService, userID int64, count int) []*model.DmPost {
	posts := make([]*model.DmPost, count)
	for i := 0; i < count; i++ {
		title := fmt.Sprintf("Post %d", i+1)
		posts[i] = CreateTestPost(t, svc, userID, title)
	}
	return posts
}
