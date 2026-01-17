package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/taku-o/go-webdb-template/internal/model"
)

// MockDmPostRepository はテスト用のモックリポジトリ
type MockDmPostRepository struct {
	CreateFunc       func(ctx context.Context, req *model.CreateDmPostRequest) (*model.DmPost, error)
	GetByIDFunc      func(ctx context.Context, id string, userID string) (*model.DmPost, error)
	ListByUserIDFunc func(ctx context.Context, userID string, limit, offset int) ([]*model.DmPost, error)
	ListFunc         func(ctx context.Context, limit, offset int) ([]*model.DmPost, error)
	GetUserPostsFunc func(ctx context.Context, limit, offset int) ([]*model.DmUserPost, error)
	UpdateFunc       func(ctx context.Context, id string, userID string, req *model.UpdateDmPostRequest) (*model.DmPost, error)
	DeleteFunc       func(ctx context.Context, id string, userID string) error
}

func (m *MockDmPostRepository) Create(ctx context.Context, req *model.CreateDmPostRequest) (*model.DmPost, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, req)
	}
	return nil, nil
}

func (m *MockDmPostRepository) GetByID(ctx context.Context, id string, userID string) (*model.DmPost, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id, userID)
	}
	return nil, nil
}

func (m *MockDmPostRepository) ListByUserID(ctx context.Context, userID string, limit, offset int) ([]*model.DmPost, error) {
	if m.ListByUserIDFunc != nil {
		return m.ListByUserIDFunc(ctx, userID, limit, offset)
	}
	return nil, nil
}

func (m *MockDmPostRepository) List(ctx context.Context, limit, offset int) ([]*model.DmPost, error) {
	if m.ListFunc != nil {
		return m.ListFunc(ctx, limit, offset)
	}
	return nil, nil
}

func (m *MockDmPostRepository) GetUserPosts(ctx context.Context, limit, offset int) ([]*model.DmUserPost, error) {
	if m.GetUserPostsFunc != nil {
		return m.GetUserPostsFunc(ctx, limit, offset)
	}
	return nil, nil
}

func (m *MockDmPostRepository) Update(ctx context.Context, id string, userID string, req *model.UpdateDmPostRequest) (*model.DmPost, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, id, userID, req)
	}
	return nil, nil
}

func (m *MockDmPostRepository) Delete(ctx context.Context, id string, userID string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id, userID)
	}
	return nil
}

func TestDmPostService_CreateDmPost(t *testing.T) {
	tests := []struct {
		name          string
		req           *model.CreateDmPostRequest
		setupMockPost func() *MockDmPostRepository
		setupMockUser func() *MockDmUserRepository
		wantErr       bool
		errContain    string
	}{
		{
			name: "正常系: 投稿を作成できる",
			req: &model.CreateDmPostRequest{
				UserID:  "user-001",
				Title:   "Test Title",
				Content: "Test Content",
			},
			setupMockPost: func() *MockDmPostRepository {
				return &MockDmPostRepository{
					CreateFunc: func(ctx context.Context, req *model.CreateDmPostRequest) (*model.DmPost, error) {
						return &model.DmPost{
							ID:        "post-001",
							UserID:    req.UserID,
							Title:     req.Title,
							Content:   req.Content,
							CreatedAt: time.Now(),
							UpdatedAt: time.Now(),
						}, nil
					},
				}
			},
			setupMockUser: func() *MockDmUserRepository {
				return &MockDmUserRepository{
					GetByIDFunc: func(ctx context.Context, id string) (*model.DmUser, error) {
						return &model.DmUser{ID: id, Name: "Test User"}, nil
					},
				}
			},
			wantErr: false,
		},
		{
			name: "異常系: ユーザーIDが空の場合エラー",
			req: &model.CreateDmPostRequest{
				UserID:  "",
				Title:   "Test Title",
				Content: "Test Content",
			},
			setupMockPost: func() *MockDmPostRepository {
				return &MockDmPostRepository{}
			},
			setupMockUser: func() *MockDmUserRepository {
				return &MockDmUserRepository{}
			},
			wantErr:    true,
			errContain: "user id is required",
		},
		{
			name: "異常系: タイトルが空の場合エラー",
			req: &model.CreateDmPostRequest{
				UserID:  "user-001",
				Title:   "",
				Content: "Test Content",
			},
			setupMockPost: func() *MockDmPostRepository {
				return &MockDmPostRepository{}
			},
			setupMockUser: func() *MockDmUserRepository {
				return &MockDmUserRepository{}
			},
			wantErr:    true,
			errContain: "title is required",
		},
		{
			name: "異常系: コンテンツが空の場合エラー",
			req: &model.CreateDmPostRequest{
				UserID:  "user-001",
				Title:   "Test Title",
				Content: "",
			},
			setupMockPost: func() *MockDmPostRepository {
				return &MockDmPostRepository{}
			},
			setupMockUser: func() *MockDmUserRepository {
				return &MockDmUserRepository{}
			},
			wantErr:    true,
			errContain: "content is required",
		},
		{
			name: "異常系: ユーザーが存在しない場合エラー",
			req: &model.CreateDmPostRequest{
				UserID:  "user-001",
				Title:   "Test Title",
				Content: "Test Content",
			},
			setupMockPost: func() *MockDmPostRepository {
				return &MockDmPostRepository{}
			},
			setupMockUser: func() *MockDmUserRepository {
				return &MockDmUserRepository{
					GetByIDFunc: func(ctx context.Context, id string) (*model.DmUser, error) {
						return nil, errors.New("user not found")
					},
				}
			},
			wantErr:    true,
			errContain: "user not found",
		},
		{
			name: "異常系: リポジトリエラー",
			req: &model.CreateDmPostRequest{
				UserID:  "user-001",
				Title:   "Test Title",
				Content: "Test Content",
			},
			setupMockPost: func() *MockDmPostRepository {
				return &MockDmPostRepository{
					CreateFunc: func(ctx context.Context, req *model.CreateDmPostRequest) (*model.DmPost, error) {
						return nil, errors.New("database error")
					},
				}
			},
			setupMockUser: func() *MockDmUserRepository {
				return &MockDmUserRepository{
					GetByIDFunc: func(ctx context.Context, id string) (*model.DmUser, error) {
						return &model.DmUser{ID: id}, nil
					},
				}
			},
			wantErr:    true,
			errContain: "failed to create post",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPostRepo := tt.setupMockPost()
			mockUserRepo := tt.setupMockUser()
			s := NewDmPostService(mockPostRepo, mockUserRepo)

			got, err := s.CreateDmPost(context.Background(), tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContain != "" {
					assert.Contains(t, err.Error(), tt.errContain)
				}
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, got)
			assert.Equal(t, tt.req.Title, got.Title)
			assert.Equal(t, tt.req.Content, got.Content)
		})
	}
}

func TestDmPostService_GetDmPost(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		userID        string
		setupMockPost func() *MockDmPostRepository
		wantErr       bool
		errContain    string
	}{
		{
			name:   "正常系: 投稿を取得できる",
			id:     "post-001",
			userID: "user-001",
			setupMockPost: func() *MockDmPostRepository {
				return &MockDmPostRepository{
					GetByIDFunc: func(ctx context.Context, id string, userID string) (*model.DmPost, error) {
						return &model.DmPost{
							ID:        id,
							UserID:    userID,
							Title:     "Test Title",
							Content:   "Test Content",
							CreatedAt: time.Now(),
							UpdatedAt: time.Now(),
						}, nil
					},
				}
			},
			wantErr: false,
		},
		{
			name:   "異常系: 投稿IDが空の場合エラー",
			id:     "",
			userID: "user-001",
			setupMockPost: func() *MockDmPostRepository {
				return &MockDmPostRepository{}
			},
			wantErr:    true,
			errContain: "post id is required",
		},
		{
			name:   "異常系: ユーザーIDが空の場合エラー",
			id:     "post-001",
			userID: "",
			setupMockPost: func() *MockDmPostRepository {
				return &MockDmPostRepository{}
			},
			wantErr:    true,
			errContain: "user id is required",
		},
		{
			name:   "異常系: リポジトリエラー",
			id:     "post-001",
			userID: "user-001",
			setupMockPost: func() *MockDmPostRepository {
				return &MockDmPostRepository{
					GetByIDFunc: func(ctx context.Context, id string, userID string) (*model.DmPost, error) {
						return nil, errors.New("not found")
					},
				}
			},
			wantErr:    true,
			errContain: "failed to get post",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPostRepo := tt.setupMockPost()
			s := NewDmPostService(mockPostRepo, &MockDmUserRepository{})

			got, err := s.GetDmPost(context.Background(), tt.id, tt.userID)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContain != "" {
					assert.Contains(t, err.Error(), tt.errContain)
				}
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, got)
			assert.Equal(t, tt.id, got.ID)
		})
	}
}

func TestDmPostService_ListDmPosts(t *testing.T) {
	tests := []struct {
		name          string
		limit         int
		offset        int
		setupMockPost func() *MockDmPostRepository
		wantErr       bool
		errContain    string
	}{
		{
			name:   "正常系: 投稿一覧を取得できる",
			limit:  10,
			offset: 0,
			setupMockPost: func() *MockDmPostRepository {
				return &MockDmPostRepository{
					ListFunc: func(ctx context.Context, limit, offset int) ([]*model.DmPost, error) {
						return []*model.DmPost{
							{ID: "post-001", Title: "Post 1"},
							{ID: "post-002", Title: "Post 2"},
						}, nil
					},
				}
			},
			wantErr: false,
		},
		{
			name:   "正常系: limitが0以下の場合デフォルト20が適用される",
			limit:  0,
			offset: 0,
			setupMockPost: func() *MockDmPostRepository {
				return &MockDmPostRepository{
					ListFunc: func(ctx context.Context, limit, offset int) ([]*model.DmPost, error) {
						assert.Equal(t, 20, limit)
						return []*model.DmPost{}, nil
					},
				}
			},
			wantErr: false,
		},
		{
			name:   "正常系: limitが100を超える場合100が適用される",
			limit:  200,
			offset: 0,
			setupMockPost: func() *MockDmPostRepository {
				return &MockDmPostRepository{
					ListFunc: func(ctx context.Context, limit, offset int) ([]*model.DmPost, error) {
						assert.Equal(t, 100, limit)
						return []*model.DmPost{}, nil
					},
				}
			},
			wantErr: false,
		},
		{
			name:   "異常系: リポジトリエラー",
			limit:  10,
			offset: 0,
			setupMockPost: func() *MockDmPostRepository {
				return &MockDmPostRepository{
					ListFunc: func(ctx context.Context, limit, offset int) ([]*model.DmPost, error) {
						return nil, errors.New("database error")
					},
				}
			},
			wantErr:    true,
			errContain: "failed to list posts",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPostRepo := tt.setupMockPost()
			s := NewDmPostService(mockPostRepo, &MockDmUserRepository{})

			got, err := s.ListDmPosts(context.Background(), tt.limit, tt.offset)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContain != "" {
					assert.Contains(t, err.Error(), tt.errContain)
				}
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, got)
		})
	}
}

func TestDmPostService_ListDmPostsByUser(t *testing.T) {
	tests := []struct {
		name          string
		userID        string
		limit         int
		offset        int
		setupMockPost func() *MockDmPostRepository
		wantErr       bool
		errContain    string
	}{
		{
			name:   "正常系: ユーザーの投稿一覧を取得できる",
			userID: "user-001",
			limit:  10,
			offset: 0,
			setupMockPost: func() *MockDmPostRepository {
				return &MockDmPostRepository{
					ListByUserIDFunc: func(ctx context.Context, userID string, limit, offset int) ([]*model.DmPost, error) {
						return []*model.DmPost{
							{ID: "post-001", UserID: userID, Title: "Post 1"},
							{ID: "post-002", UserID: userID, Title: "Post 2"},
						}, nil
					},
				}
			},
			wantErr: false,
		},
		{
			name:   "異常系: ユーザーIDが空の場合エラー",
			userID: "",
			limit:  10,
			offset: 0,
			setupMockPost: func() *MockDmPostRepository {
				return &MockDmPostRepository{}
			},
			wantErr:    true,
			errContain: "user id is required",
		},
		{
			name:   "正常系: limitが0以下の場合デフォルト20が適用される",
			userID: "user-001",
			limit:  0,
			offset: 0,
			setupMockPost: func() *MockDmPostRepository {
				return &MockDmPostRepository{
					ListByUserIDFunc: func(ctx context.Context, userID string, limit, offset int) ([]*model.DmPost, error) {
						assert.Equal(t, 20, limit)
						return []*model.DmPost{}, nil
					},
				}
			},
			wantErr: false,
		},
		{
			name:   "異常系: リポジトリエラー",
			userID: "user-001",
			limit:  10,
			offset: 0,
			setupMockPost: func() *MockDmPostRepository {
				return &MockDmPostRepository{
					ListByUserIDFunc: func(ctx context.Context, userID string, limit, offset int) ([]*model.DmPost, error) {
						return nil, errors.New("database error")
					},
				}
			},
			wantErr:    true,
			errContain: "failed to list posts by user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPostRepo := tt.setupMockPost()
			s := NewDmPostService(mockPostRepo, &MockDmUserRepository{})

			got, err := s.ListDmPostsByUser(context.Background(), tt.userID, tt.limit, tt.offset)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContain != "" {
					assert.Contains(t, err.Error(), tt.errContain)
				}
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, got)
		})
	}
}

func TestDmPostService_GetDmUserPosts(t *testing.T) {
	tests := []struct {
		name          string
		limit         int
		offset        int
		setupMockPost func() *MockDmPostRepository
		wantErr       bool
		errContain    string
	}{
		{
			name:   "正常系: ユーザー投稿一覧を取得できる",
			limit:  10,
			offset: 0,
			setupMockPost: func() *MockDmPostRepository {
				return &MockDmPostRepository{
					GetUserPostsFunc: func(ctx context.Context, limit, offset int) ([]*model.DmUserPost, error) {
						return []*model.DmUserPost{
							{PostID: "post-001", PostTitle: "Post 1", UserID: "user-001", UserName: "User 1"},
							{PostID: "post-002", PostTitle: "Post 2", UserID: "user-002", UserName: "User 2"},
						}, nil
					},
				}
			},
			wantErr: false,
		},
		{
			name:   "正常系: limitが0以下の場合デフォルト20が適用される",
			limit:  0,
			offset: 0,
			setupMockPost: func() *MockDmPostRepository {
				return &MockDmPostRepository{
					GetUserPostsFunc: func(ctx context.Context, limit, offset int) ([]*model.DmUserPost, error) {
						assert.Equal(t, 20, limit)
						return []*model.DmUserPost{}, nil
					},
				}
			},
			wantErr: false,
		},
		{
			name:   "異常系: リポジトリエラー",
			limit:  10,
			offset: 0,
			setupMockPost: func() *MockDmPostRepository {
				return &MockDmPostRepository{
					GetUserPostsFunc: func(ctx context.Context, limit, offset int) ([]*model.DmUserPost, error) {
						return nil, errors.New("database error")
					},
				}
			},
			wantErr:    true,
			errContain: "failed to get user posts",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPostRepo := tt.setupMockPost()
			s := NewDmPostService(mockPostRepo, &MockDmUserRepository{})

			got, err := s.GetDmUserPosts(context.Background(), tt.limit, tt.offset)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContain != "" {
					assert.Contains(t, err.Error(), tt.errContain)
				}
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, got)
		})
	}
}

func TestDmPostService_UpdateDmPost(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		userID        string
		req           *model.UpdateDmPostRequest
		setupMockPost func() *MockDmPostRepository
		wantErr       bool
		errContain    string
	}{
		{
			name:   "正常系: 投稿を更新できる",
			id:     "post-001",
			userID: "user-001",
			req: &model.UpdateDmPostRequest{
				Title:   "Updated Title",
				Content: "Updated Content",
			},
			setupMockPost: func() *MockDmPostRepository {
				return &MockDmPostRepository{
					UpdateFunc: func(ctx context.Context, id string, userID string, req *model.UpdateDmPostRequest) (*model.DmPost, error) {
						return &model.DmPost{
							ID:        id,
							UserID:    userID,
							Title:     req.Title,
							Content:   req.Content,
							CreatedAt: time.Now(),
							UpdatedAt: time.Now(),
						}, nil
					},
				}
			},
			wantErr: false,
		},
		{
			name:   "正常系: タイトルのみ更新できる",
			id:     "post-001",
			userID: "user-001",
			req: &model.UpdateDmPostRequest{
				Title: "Updated Title",
			},
			setupMockPost: func() *MockDmPostRepository {
				return &MockDmPostRepository{
					UpdateFunc: func(ctx context.Context, id string, userID string, req *model.UpdateDmPostRequest) (*model.DmPost, error) {
						return &model.DmPost{
							ID:      id,
							UserID:  userID,
							Title:   req.Title,
							Content: "Original Content",
						}, nil
					},
				}
			},
			wantErr: false,
		},
		{
			name:   "異常系: 投稿IDが空の場合エラー",
			id:     "",
			userID: "user-001",
			req: &model.UpdateDmPostRequest{
				Title: "Updated Title",
			},
			setupMockPost: func() *MockDmPostRepository {
				return &MockDmPostRepository{}
			},
			wantErr:    true,
			errContain: "post id is required",
		},
		{
			name:   "異常系: ユーザーIDが空の場合エラー",
			id:     "post-001",
			userID: "",
			req: &model.UpdateDmPostRequest{
				Title: "Updated Title",
			},
			setupMockPost: func() *MockDmPostRepository {
				return &MockDmPostRepository{}
			},
			wantErr:    true,
			errContain: "user id is required",
		},
		{
			name:   "異常系: 更新フィールドがない場合エラー",
			id:     "post-001",
			userID: "user-001",
			req: &model.UpdateDmPostRequest{
				Title:   "",
				Content: "",
			},
			setupMockPost: func() *MockDmPostRepository {
				return &MockDmPostRepository{}
			},
			wantErr:    true,
			errContain: "no fields to update",
		},
		{
			name:   "異常系: リポジトリエラー",
			id:     "post-001",
			userID: "user-001",
			req: &model.UpdateDmPostRequest{
				Title: "Updated Title",
			},
			setupMockPost: func() *MockDmPostRepository {
				return &MockDmPostRepository{
					UpdateFunc: func(ctx context.Context, id string, userID string, req *model.UpdateDmPostRequest) (*model.DmPost, error) {
						return nil, errors.New("database error")
					},
				}
			},
			wantErr:    true,
			errContain: "failed to update post",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPostRepo := tt.setupMockPost()
			s := NewDmPostService(mockPostRepo, &MockDmUserRepository{})

			got, err := s.UpdateDmPost(context.Background(), tt.id, tt.userID, tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContain != "" {
					assert.Contains(t, err.Error(), tt.errContain)
				}
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, got)
		})
	}
}

func TestDmPostService_DeleteDmPost(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		userID        string
		setupMockPost func() *MockDmPostRepository
		wantErr       bool
		errContain    string
	}{
		{
			name:   "正常系: 投稿を削除できる",
			id:     "post-001",
			userID: "user-001",
			setupMockPost: func() *MockDmPostRepository {
				return &MockDmPostRepository{
					DeleteFunc: func(ctx context.Context, id string, userID string) error {
						return nil
					},
				}
			},
			wantErr: false,
		},
		{
			name:   "異常系: 投稿IDが空の場合エラー",
			id:     "",
			userID: "user-001",
			setupMockPost: func() *MockDmPostRepository {
				return &MockDmPostRepository{}
			},
			wantErr:    true,
			errContain: "post id is required",
		},
		{
			name:   "異常系: ユーザーIDが空の場合エラー",
			id:     "post-001",
			userID: "",
			setupMockPost: func() *MockDmPostRepository {
				return &MockDmPostRepository{}
			},
			wantErr:    true,
			errContain: "user id is required",
		},
		{
			name:   "異常系: リポジトリエラー",
			id:     "post-001",
			userID: "user-001",
			setupMockPost: func() *MockDmPostRepository {
				return &MockDmPostRepository{
					DeleteFunc: func(ctx context.Context, id string, userID string) error {
						return errors.New("database error")
					},
				}
			},
			wantErr:    true,
			errContain: "failed to delete post",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPostRepo := tt.setupMockPost()
			s := NewDmPostService(mockPostRepo, &MockDmUserRepository{})

			err := s.DeleteDmPost(context.Background(), tt.id, tt.userID)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContain != "" {
					assert.Contains(t, err.Error(), tt.errContain)
				}
				return
			}

			assert.NoError(t, err)
		})
	}
}
