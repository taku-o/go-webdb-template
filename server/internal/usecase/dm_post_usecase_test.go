package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/taku-o/go-webdb-template/internal/model"
)

// MockDmPostService はテスト用のDmPostServiceモック
type MockDmPostService struct {
	CreateDmPostFunc      func(ctx context.Context, req *model.CreateDmPostRequest) (*model.DmPost, error)
	GetDmPostFunc         func(ctx context.Context, id string, userID string) (*model.DmPost, error)
	ListDmPostsFunc       func(ctx context.Context, limit, offset int) ([]*model.DmPost, error)
	ListDmPostsByUserFunc func(ctx context.Context, userID string, limit, offset int) ([]*model.DmPost, error)
	GetDmUserPostsFunc    func(ctx context.Context, limit, offset int) ([]*model.DmUserPost, error)
	UpdateDmPostFunc      func(ctx context.Context, id string, userID string, req *model.UpdateDmPostRequest) (*model.DmPost, error)
	DeleteDmPostFunc      func(ctx context.Context, id string, userID string) error
}

func (m *MockDmPostService) CreateDmPost(ctx context.Context, req *model.CreateDmPostRequest) (*model.DmPost, error) {
	if m.CreateDmPostFunc != nil {
		return m.CreateDmPostFunc(ctx, req)
	}
	return nil, nil
}

func (m *MockDmPostService) GetDmPost(ctx context.Context, id string, userID string) (*model.DmPost, error) {
	if m.GetDmPostFunc != nil {
		return m.GetDmPostFunc(ctx, id, userID)
	}
	return nil, nil
}

func (m *MockDmPostService) ListDmPosts(ctx context.Context, limit, offset int) ([]*model.DmPost, error) {
	if m.ListDmPostsFunc != nil {
		return m.ListDmPostsFunc(ctx, limit, offset)
	}
	return nil, nil
}

func (m *MockDmPostService) ListDmPostsByUser(ctx context.Context, userID string, limit, offset int) ([]*model.DmPost, error) {
	if m.ListDmPostsByUserFunc != nil {
		return m.ListDmPostsByUserFunc(ctx, userID, limit, offset)
	}
	return nil, nil
}

func (m *MockDmPostService) GetDmUserPosts(ctx context.Context, limit, offset int) ([]*model.DmUserPost, error) {
	if m.GetDmUserPostsFunc != nil {
		return m.GetDmUserPostsFunc(ctx, limit, offset)
	}
	return nil, nil
}

func (m *MockDmPostService) UpdateDmPost(ctx context.Context, id string, userID string, req *model.UpdateDmPostRequest) (*model.DmPost, error) {
	if m.UpdateDmPostFunc != nil {
		return m.UpdateDmPostFunc(ctx, id, userID, req)
	}
	return nil, nil
}

func (m *MockDmPostService) DeleteDmPost(ctx context.Context, id string, userID string) error {
	if m.DeleteDmPostFunc != nil {
		return m.DeleteDmPostFunc(ctx, id, userID)
	}
	return nil
}

func TestDmPostUsecase_CreateDmPost(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	tests := []struct {
		name        string
		req         *model.CreateDmPostRequest
		mockFunc    func(ctx context.Context, req *model.CreateDmPostRequest) (*model.DmPost, error)
		wantErr     bool
		wantPost    *model.DmPost
	}{
		{
			name: "creates post successfully",
			req: &model.CreateDmPostRequest{
				UserID:  "user123",
				Title:   "Test Post",
				Content: "Test Content",
			},
			mockFunc: func(ctx context.Context, req *model.CreateDmPostRequest) (*model.DmPost, error) {
				return &model.DmPost{
					ID:        "post123",
					UserID:    req.UserID,
					Title:     req.Title,
					Content:   req.Content,
					CreatedAt: now,
					UpdatedAt: now,
				}, nil
			},
			wantErr: false,
			wantPost: &model.DmPost{
				ID:        "post123",
				UserID:    "user123",
				Title:     "Test Post",
				Content:   "Test Content",
				CreatedAt: now,
				UpdatedAt: now,
			},
		},
		{
			name: "returns error when service fails",
			req: &model.CreateDmPostRequest{
				UserID:  "user123",
				Title:   "Test Post",
				Content: "Test Content",
			},
			mockFunc: func(ctx context.Context, req *model.CreateDmPostRequest) (*model.DmPost, error) {
				return nil, errors.New("service error")
			},
			wantErr:  true,
			wantPost: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockDmPostService{
				CreateDmPostFunc: tt.mockFunc,
			}
			usecase := NewDmPostUsecase(mockService)

			got, err := usecase.CreateDmPost(ctx, tt.req)

			if tt.wantErr {
				require.Error(t, err)
				assert.Nil(t, got)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantPost.ID, got.ID)
				assert.Equal(t, tt.wantPost.UserID, got.UserID)
				assert.Equal(t, tt.wantPost.Title, got.Title)
				assert.Equal(t, tt.wantPost.Content, got.Content)
			}
		})
	}
}

func TestDmPostUsecase_GetDmPost(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	tests := []struct {
		name     string
		id       string
		userID   string
		mockFunc func(ctx context.Context, id string, userID string) (*model.DmPost, error)
		wantErr  bool
		wantPost *model.DmPost
	}{
		{
			name:   "gets post successfully",
			id:     "post123",
			userID: "user123",
			mockFunc: func(ctx context.Context, id string, userID string) (*model.DmPost, error) {
				return &model.DmPost{
					ID:        id,
					UserID:    userID,
					Title:     "Test Post",
					Content:   "Test Content",
					CreatedAt: now,
					UpdatedAt: now,
				}, nil
			},
			wantErr: false,
			wantPost: &model.DmPost{
				ID:        "post123",
				UserID:    "user123",
				Title:     "Test Post",
				Content:   "Test Content",
				CreatedAt: now,
				UpdatedAt: now,
			},
		},
		{
			name:   "returns error when post not found",
			id:     "nonexistent",
			userID: "user123",
			mockFunc: func(ctx context.Context, id string, userID string) (*model.DmPost, error) {
				return nil, errors.New("post not found")
			},
			wantErr:  true,
			wantPost: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockDmPostService{
				GetDmPostFunc: tt.mockFunc,
			}
			usecase := NewDmPostUsecase(mockService)

			got, err := usecase.GetDmPost(ctx, tt.id, tt.userID)

			if tt.wantErr {
				require.Error(t, err)
				assert.Nil(t, got)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantPost.ID, got.ID)
				assert.Equal(t, tt.wantPost.UserID, got.UserID)
				assert.Equal(t, tt.wantPost.Title, got.Title)
			}
		})
	}
}

func TestDmPostUsecase_ListDmPosts(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	tests := []struct {
		name      string
		limit     int
		offset    int
		mockFunc  func(ctx context.Context, limit, offset int) ([]*model.DmPost, error)
		wantErr   bool
		wantCount int
	}{
		{
			name:   "lists posts successfully",
			limit:  10,
			offset: 0,
			mockFunc: func(ctx context.Context, limit, offset int) ([]*model.DmPost, error) {
				return []*model.DmPost{
					{ID: "post1", UserID: "user1", Title: "Post 1", Content: "Content 1", CreatedAt: now, UpdatedAt: now},
					{ID: "post2", UserID: "user2", Title: "Post 2", Content: "Content 2", CreatedAt: now, UpdatedAt: now},
				}, nil
			},
			wantErr:   false,
			wantCount: 2,
		},
		{
			name:   "returns error when service fails",
			limit:  10,
			offset: 0,
			mockFunc: func(ctx context.Context, limit, offset int) ([]*model.DmPost, error) {
				return nil, errors.New("service error")
			},
			wantErr:   true,
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockDmPostService{
				ListDmPostsFunc: tt.mockFunc,
			}
			usecase := NewDmPostUsecase(mockService)

			got, err := usecase.ListDmPosts(ctx, tt.limit, tt.offset)

			if tt.wantErr {
				require.Error(t, err)
				assert.Nil(t, got)
			} else {
				require.NoError(t, err)
				assert.Len(t, got, tt.wantCount)
			}
		})
	}
}

func TestDmPostUsecase_ListDmPostsByUser(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	tests := []struct {
		name      string
		userID    string
		limit     int
		offset    int
		mockFunc  func(ctx context.Context, userID string, limit, offset int) ([]*model.DmPost, error)
		wantErr   bool
		wantCount int
	}{
		{
			name:   "lists user posts successfully",
			userID: "user123",
			limit:  10,
			offset: 0,
			mockFunc: func(ctx context.Context, userID string, limit, offset int) ([]*model.DmPost, error) {
				return []*model.DmPost{
					{ID: "post1", UserID: userID, Title: "Post 1", Content: "Content 1", CreatedAt: now, UpdatedAt: now},
					{ID: "post2", UserID: userID, Title: "Post 2", Content: "Content 2", CreatedAt: now, UpdatedAt: now},
				}, nil
			},
			wantErr:   false,
			wantCount: 2,
		},
		{
			name:   "returns error when service fails",
			userID: "user123",
			limit:  10,
			offset: 0,
			mockFunc: func(ctx context.Context, userID string, limit, offset int) ([]*model.DmPost, error) {
				return nil, errors.New("service error")
			},
			wantErr:   true,
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockDmPostService{
				ListDmPostsByUserFunc: tt.mockFunc,
			}
			usecase := NewDmPostUsecase(mockService)

			got, err := usecase.ListDmPostsByUser(ctx, tt.userID, tt.limit, tt.offset)

			if tt.wantErr {
				require.Error(t, err)
				assert.Nil(t, got)
			} else {
				require.NoError(t, err)
				assert.Len(t, got, tt.wantCount)
			}
		})
	}
}

func TestDmPostUsecase_GetDmUserPosts(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	tests := []struct {
		name      string
		limit     int
		offset    int
		mockFunc  func(ctx context.Context, limit, offset int) ([]*model.DmUserPost, error)
		wantErr   bool
		wantCount int
	}{
		{
			name:   "gets user posts successfully",
			limit:  10,
			offset: 0,
			mockFunc: func(ctx context.Context, limit, offset int) ([]*model.DmUserPost, error) {
				return []*model.DmUserPost{
					{UserID: "user1", UserName: "User 1", PostID: "post1", PostTitle: "Post 1", CreatedAt: now},
					{UserID: "user2", UserName: "User 2", PostID: "post2", PostTitle: "Post 2", CreatedAt: now},
				}, nil
			},
			wantErr:   false,
			wantCount: 2,
		},
		{
			name:   "returns error when service fails",
			limit:  10,
			offset: 0,
			mockFunc: func(ctx context.Context, limit, offset int) ([]*model.DmUserPost, error) {
				return nil, errors.New("service error")
			},
			wantErr:   true,
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockDmPostService{
				GetDmUserPostsFunc: tt.mockFunc,
			}
			usecase := NewDmPostUsecase(mockService)

			got, err := usecase.GetDmUserPosts(ctx, tt.limit, tt.offset)

			if tt.wantErr {
				require.Error(t, err)
				assert.Nil(t, got)
			} else {
				require.NoError(t, err)
				assert.Len(t, got, tt.wantCount)
			}
		})
	}
}

func TestDmPostUsecase_UpdateDmPost(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	tests := []struct {
		name     string
		id       string
		userID   string
		req      *model.UpdateDmPostRequest
		mockFunc func(ctx context.Context, id string, userID string, req *model.UpdateDmPostRequest) (*model.DmPost, error)
		wantErr  bool
		wantPost *model.DmPost
	}{
		{
			name:   "updates post successfully",
			id:     "post123",
			userID: "user123",
			req: &model.UpdateDmPostRequest{
				Title:   "Updated Title",
				Content: "Updated Content",
			},
			mockFunc: func(ctx context.Context, id string, userID string, req *model.UpdateDmPostRequest) (*model.DmPost, error) {
				return &model.DmPost{
					ID:        id,
					UserID:    userID,
					Title:     req.Title,
					Content:   req.Content,
					CreatedAt: now,
					UpdatedAt: now,
				}, nil
			},
			wantErr: false,
			wantPost: &model.DmPost{
				ID:        "post123",
				UserID:    "user123",
				Title:     "Updated Title",
				Content:   "Updated Content",
				CreatedAt: now,
				UpdatedAt: now,
			},
		},
		{
			name:   "returns error when service fails",
			id:     "post123",
			userID: "user123",
			req: &model.UpdateDmPostRequest{
				Title: "Updated Title",
			},
			mockFunc: func(ctx context.Context, id string, userID string, req *model.UpdateDmPostRequest) (*model.DmPost, error) {
				return nil, errors.New("service error")
			},
			wantErr:  true,
			wantPost: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockDmPostService{
				UpdateDmPostFunc: tt.mockFunc,
			}
			usecase := NewDmPostUsecase(mockService)

			got, err := usecase.UpdateDmPost(ctx, tt.id, tt.userID, tt.req)

			if tt.wantErr {
				require.Error(t, err)
				assert.Nil(t, got)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantPost.ID, got.ID)
				assert.Equal(t, tt.wantPost.Title, got.Title)
				assert.Equal(t, tt.wantPost.Content, got.Content)
			}
		})
	}
}

func TestDmPostUsecase_DeleteDmPost(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name     string
		id       string
		userID   string
		mockFunc func(ctx context.Context, id string, userID string) error
		wantErr  bool
	}{
		{
			name:   "deletes post successfully",
			id:     "post123",
			userID: "user123",
			mockFunc: func(ctx context.Context, id string, userID string) error {
				return nil
			},
			wantErr: false,
		},
		{
			name:   "returns error when service fails",
			id:     "post123",
			userID: "user123",
			mockFunc: func(ctx context.Context, id string, userID string) error {
				return errors.New("service error")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockDmPostService{
				DeleteDmPostFunc: tt.mockFunc,
			}
			usecase := NewDmPostUsecase(mockService)

			err := usecase.DeleteDmPost(ctx, tt.id, tt.userID)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
