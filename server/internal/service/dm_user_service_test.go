package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/taku-o/go-webdb-template/internal/model"
)

// MockDmUserRepository はテスト用のモックリポジトリ
type MockDmUserRepository struct {
	CreateFunc           func(ctx context.Context, req *model.CreateDmUserRequest) (*model.DmUser, error)
	GetByIDFunc          func(ctx context.Context, id string) (*model.DmUser, error)
	ListFunc             func(ctx context.Context, limit, offset int) ([]*model.DmUser, error)
	UpdateFunc           func(ctx context.Context, id string, req *model.UpdateDmUserRequest) (*model.DmUser, error)
	DeleteFunc           func(ctx context.Context, id string) error
	CheckEmailExistsFunc func(ctx context.Context, email string) (bool, error)
}

func (m *MockDmUserRepository) Create(ctx context.Context, req *model.CreateDmUserRequest) (*model.DmUser, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, req)
	}
	return nil, nil
}

func (m *MockDmUserRepository) GetByID(ctx context.Context, id string) (*model.DmUser, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockDmUserRepository) List(ctx context.Context, limit, offset int) ([]*model.DmUser, error) {
	if m.ListFunc != nil {
		return m.ListFunc(ctx, limit, offset)
	}
	return nil, nil
}

func (m *MockDmUserRepository) Update(ctx context.Context, id string, req *model.UpdateDmUserRequest) (*model.DmUser, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, id, req)
	}
	return nil, nil
}

func (m *MockDmUserRepository) Delete(ctx context.Context, id string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

func (m *MockDmUserRepository) CheckEmailExists(ctx context.Context, email string) (bool, error) {
	if m.CheckEmailExistsFunc != nil {
		return m.CheckEmailExistsFunc(ctx, email)
	}
	return false, nil
}

func TestDmUserService_CreateDmUser(t *testing.T) {
	tests := []struct {
		name       string
		req        *model.CreateDmUserRequest
		setupMock  func() *MockDmUserRepository
		wantErr    bool
		errContain string
	}{
		{
			name: "正常系: ユーザーを作成できる",
			req: &model.CreateDmUserRequest{
				Name:  "Test User",
				Email: "test@example.com",
			},
			setupMock: func() *MockDmUserRepository {
				return &MockDmUserRepository{
					CreateFunc: func(ctx context.Context, req *model.CreateDmUserRequest) (*model.DmUser, error) {
						return &model.DmUser{
							ID:        "user-001",
							Name:      req.Name,
							Email:     req.Email,
							CreatedAt: time.Now(),
							UpdatedAt: time.Now(),
						}, nil
					},
				}
			},
			wantErr: false,
		},
		{
			name: "異常系: 名前が空の場合エラー",
			req: &model.CreateDmUserRequest{
				Name:  "",
				Email: "test@example.com",
			},
			setupMock: func() *MockDmUserRepository {
				return &MockDmUserRepository{}
			},
			wantErr:    true,
			errContain: "name is required",
		},
		{
			name: "異常系: メールが空の場合エラー",
			req: &model.CreateDmUserRequest{
				Name:  "Test User",
				Email: "",
			},
			setupMock: func() *MockDmUserRepository {
				return &MockDmUserRepository{}
			},
			wantErr:    true,
			errContain: "email is required",
		},
		{
			name: "異常系: リポジトリエラー",
			req: &model.CreateDmUserRequest{
				Name:  "Test User",
				Email: "test@example.com",
			},
			setupMock: func() *MockDmUserRepository {
				return &MockDmUserRepository{
					CreateFunc: func(ctx context.Context, req *model.CreateDmUserRequest) (*model.DmUser, error) {
						return nil, errors.New("database error")
					},
				}
			},
			wantErr:    true,
			errContain: "failed to create user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := tt.setupMock()
			s := NewDmUserService(mockRepo)

			got, err := s.CreateDmUser(context.Background(), tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContain != "" {
					assert.Contains(t, err.Error(), tt.errContain)
				}
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, got)
			assert.Equal(t, tt.req.Name, got.Name)
			assert.Equal(t, tt.req.Email, got.Email)
		})
	}
}

func TestDmUserService_GetDmUser(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		setupMock  func() *MockDmUserRepository
		wantErr    bool
		errContain string
	}{
		{
			name: "正常系: ユーザーを取得できる",
			id:   "user-001",
			setupMock: func() *MockDmUserRepository {
				return &MockDmUserRepository{
					GetByIDFunc: func(ctx context.Context, id string) (*model.DmUser, error) {
						return &model.DmUser{
							ID:        id,
							Name:      "Test User",
							Email:     "test@example.com",
							CreatedAt: time.Now(),
							UpdatedAt: time.Now(),
						}, nil
					},
				}
			},
			wantErr: false,
		},
		{
			name: "異常系: IDが空の場合エラー",
			id:   "",
			setupMock: func() *MockDmUserRepository {
				return &MockDmUserRepository{}
			},
			wantErr:    true,
			errContain: "user id is required",
		},
		{
			name: "異常系: リポジトリエラー",
			id:   "user-001",
			setupMock: func() *MockDmUserRepository {
				return &MockDmUserRepository{
					GetByIDFunc: func(ctx context.Context, id string) (*model.DmUser, error) {
						return nil, errors.New("not found")
					},
				}
			},
			wantErr:    true,
			errContain: "failed to get user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := tt.setupMock()
			s := NewDmUserService(mockRepo)

			got, err := s.GetDmUser(context.Background(), tt.id)

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

func TestDmUserService_ListDmUsers(t *testing.T) {
	tests := []struct {
		name           string
		limit          int
		offset         int
		setupMock      func() *MockDmUserRepository
		expectedLimit  int
		expectedOffset int
		wantErr        bool
		errContain     string
	}{
		{
			name:   "正常系: ユーザー一覧を取得できる",
			limit:  10,
			offset: 0,
			setupMock: func() *MockDmUserRepository {
				return &MockDmUserRepository{
					ListFunc: func(ctx context.Context, limit, offset int) ([]*model.DmUser, error) {
						return []*model.DmUser{
							{ID: "user-001", Name: "User 1", Email: "user1@example.com"},
							{ID: "user-002", Name: "User 2", Email: "user2@example.com"},
						}, nil
					},
				}
			},
			expectedLimit:  10,
			expectedOffset: 0,
			wantErr:        false,
		},
		{
			name:   "正常系: limitが0以下の場合デフォルト20が適用される",
			limit:  0,
			offset: 0,
			setupMock: func() *MockDmUserRepository {
				return &MockDmUserRepository{
					ListFunc: func(ctx context.Context, limit, offset int) ([]*model.DmUser, error) {
						assert.Equal(t, 20, limit)
						return []*model.DmUser{}, nil
					},
				}
			},
			expectedLimit:  20,
			expectedOffset: 0,
			wantErr:        false,
		},
		{
			name:   "正常系: limitが100を超える場合100が適用される",
			limit:  200,
			offset: 0,
			setupMock: func() *MockDmUserRepository {
				return &MockDmUserRepository{
					ListFunc: func(ctx context.Context, limit, offset int) ([]*model.DmUser, error) {
						assert.Equal(t, 100, limit)
						return []*model.DmUser{}, nil
					},
				}
			},
			expectedLimit:  100,
			expectedOffset: 0,
			wantErr:        false,
		},
		{
			name:   "正常系: offsetが負の場合0が適用される",
			limit:  10,
			offset: -5,
			setupMock: func() *MockDmUserRepository {
				return &MockDmUserRepository{
					ListFunc: func(ctx context.Context, limit, offset int) ([]*model.DmUser, error) {
						assert.Equal(t, 0, offset)
						return []*model.DmUser{}, nil
					},
				}
			},
			expectedLimit:  10,
			expectedOffset: 0,
			wantErr:        false,
		},
		{
			name:   "異常系: リポジトリエラー",
			limit:  10,
			offset: 0,
			setupMock: func() *MockDmUserRepository {
				return &MockDmUserRepository{
					ListFunc: func(ctx context.Context, limit, offset int) ([]*model.DmUser, error) {
						return nil, errors.New("database error")
					},
				}
			},
			wantErr:    true,
			errContain: "failed to list users",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := tt.setupMock()
			s := NewDmUserService(mockRepo)

			got, err := s.ListDmUsers(context.Background(), tt.limit, tt.offset)

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

func TestDmUserService_UpdateDmUser(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		req        *model.UpdateDmUserRequest
		setupMock  func() *MockDmUserRepository
		wantErr    bool
		errContain string
	}{
		{
			name: "正常系: ユーザーを更新できる",
			id:   "user-001",
			req: &model.UpdateDmUserRequest{
				Name:  "Updated Name",
				Email: "updated@example.com",
			},
			setupMock: func() *MockDmUserRepository {
				return &MockDmUserRepository{
					UpdateFunc: func(ctx context.Context, id string, req *model.UpdateDmUserRequest) (*model.DmUser, error) {
						return &model.DmUser{
							ID:        id,
							Name:      req.Name,
							Email:     req.Email,
							CreatedAt: time.Now(),
							UpdatedAt: time.Now(),
						}, nil
					},
				}
			},
			wantErr: false,
		},
		{
			name: "正常系: 名前のみ更新できる",
			id:   "user-001",
			req: &model.UpdateDmUserRequest{
				Name: "Updated Name",
			},
			setupMock: func() *MockDmUserRepository {
				return &MockDmUserRepository{
					UpdateFunc: func(ctx context.Context, id string, req *model.UpdateDmUserRequest) (*model.DmUser, error) {
						return &model.DmUser{
							ID:        id,
							Name:      req.Name,
							Email:     "original@example.com",
							CreatedAt: time.Now(),
							UpdatedAt: time.Now(),
						}, nil
					},
				}
			},
			wantErr: false,
		},
		{
			name: "異常系: IDが空の場合エラー",
			id:   "",
			req: &model.UpdateDmUserRequest{
				Name: "Updated Name",
			},
			setupMock: func() *MockDmUserRepository {
				return &MockDmUserRepository{}
			},
			wantErr:    true,
			errContain: "user id is required",
		},
		{
			name: "異常系: 更新フィールドがない場合エラー",
			id:   "user-001",
			req: &model.UpdateDmUserRequest{
				Name:  "",
				Email: "",
			},
			setupMock: func() *MockDmUserRepository {
				return &MockDmUserRepository{}
			},
			wantErr:    true,
			errContain: "no fields to update",
		},
		{
			name: "異常系: リポジトリエラー",
			id:   "user-001",
			req: &model.UpdateDmUserRequest{
				Name: "Updated Name",
			},
			setupMock: func() *MockDmUserRepository {
				return &MockDmUserRepository{
					UpdateFunc: func(ctx context.Context, id string, req *model.UpdateDmUserRequest) (*model.DmUser, error) {
						return nil, errors.New("database error")
					},
				}
			},
			wantErr:    true,
			errContain: "failed to update user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := tt.setupMock()
			s := NewDmUserService(mockRepo)

			got, err := s.UpdateDmUser(context.Background(), tt.id, tt.req)

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

func TestDmUserService_DeleteDmUser(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		setupMock  func() *MockDmUserRepository
		wantErr    bool
		errContain string
	}{
		{
			name: "正常系: ユーザーを削除できる",
			id:   "user-001",
			setupMock: func() *MockDmUserRepository {
				return &MockDmUserRepository{
					DeleteFunc: func(ctx context.Context, id string) error {
						return nil
					},
				}
			},
			wantErr: false,
		},
		{
			name: "異常系: IDが空の場合エラー",
			id:   "",
			setupMock: func() *MockDmUserRepository {
				return &MockDmUserRepository{}
			},
			wantErr:    true,
			errContain: "user id is required",
		},
		{
			name: "異常系: リポジトリエラー",
			id:   "user-001",
			setupMock: func() *MockDmUserRepository {
				return &MockDmUserRepository{
					DeleteFunc: func(ctx context.Context, id string) error {
						return errors.New("database error")
					},
				}
			},
			wantErr:    true,
			errContain: "failed to delete user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := tt.setupMock()
			s := NewDmUserService(mockRepo)

			err := s.DeleteDmUser(context.Background(), tt.id)

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

func TestDmUserService_CheckEmailExists(t *testing.T) {
	tests := []struct {
		name       string
		email      string
		setupMock  func() *MockDmUserRepository
		wantResult bool
		wantErr    bool
	}{
		{
			name:  "正常系: メールが存在する場合trueを返す",
			email: "existing@example.com",
			setupMock: func() *MockDmUserRepository {
				return &MockDmUserRepository{
					CheckEmailExistsFunc: func(ctx context.Context, email string) (bool, error) {
						return true, nil
					},
				}
			},
			wantResult: true,
			wantErr:    false,
		},
		{
			name:  "正常系: メールが存在しない場合falseを返す",
			email: "notexisting@example.com",
			setupMock: func() *MockDmUserRepository {
				return &MockDmUserRepository{
					CheckEmailExistsFunc: func(ctx context.Context, email string) (bool, error) {
						return false, nil
					},
				}
			},
			wantResult: false,
			wantErr:    false,
		},
		{
			name:  "異常系: リポジトリエラー",
			email: "test@example.com",
			setupMock: func() *MockDmUserRepository {
				return &MockDmUserRepository{
					CheckEmailExistsFunc: func(ctx context.Context, email string) (bool, error) {
						return false, errors.New("database error")
					},
				}
			},
			wantResult: false,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := tt.setupMock()
			s := NewDmUserService(mockRepo)

			got, err := s.CheckEmailExists(context.Background(), tt.email)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.wantResult, got)
		})
	}
}
