package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/taku-o/go-webdb-template/internal/model"
)

// MockDmUserService はDmUserServiceのモック
type MockDmUserService struct {
	CreateDmUserFunc     func(ctx context.Context, req *model.CreateDmUserRequest) (*model.DmUser, error)
	GetDmUserFunc        func(ctx context.Context, id string) (*model.DmUser, error)
	ListDmUsersFunc      func(ctx context.Context, limit, offset int) ([]*model.DmUser, error)
	UpdateDmUserFunc     func(ctx context.Context, id string, req *model.UpdateDmUserRequest) (*model.DmUser, error)
	DeleteDmUserFunc     func(ctx context.Context, id string) error
	CheckEmailExistsFunc func(ctx context.Context, email string) (bool, error)
}

func (m *MockDmUserService) CreateDmUser(ctx context.Context, req *model.CreateDmUserRequest) (*model.DmUser, error) {
	if m.CreateDmUserFunc != nil {
		return m.CreateDmUserFunc(ctx, req)
	}
	return nil, nil
}

func (m *MockDmUserService) GetDmUser(ctx context.Context, id string) (*model.DmUser, error) {
	if m.GetDmUserFunc != nil {
		return m.GetDmUserFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockDmUserService) ListDmUsers(ctx context.Context, limit, offset int) ([]*model.DmUser, error) {
	if m.ListDmUsersFunc != nil {
		return m.ListDmUsersFunc(ctx, limit, offset)
	}
	return nil, nil
}

func (m *MockDmUserService) UpdateDmUser(ctx context.Context, id string, req *model.UpdateDmUserRequest) (*model.DmUser, error) {
	if m.UpdateDmUserFunc != nil {
		return m.UpdateDmUserFunc(ctx, id, req)
	}
	return nil, nil
}

func (m *MockDmUserService) DeleteDmUser(ctx context.Context, id string) error {
	if m.DeleteDmUserFunc != nil {
		return m.DeleteDmUserFunc(ctx, id)
	}
	return nil
}

func (m *MockDmUserService) CheckEmailExists(ctx context.Context, email string) (bool, error) {
	if m.CheckEmailExistsFunc != nil {
		return m.CheckEmailExistsFunc(ctx, email)
	}
	return false, nil
}

func TestDmUserUsecase_CreateDmUser(t *testing.T) {
	tests := []struct {
		name        string
		mockFunc    func(ctx context.Context, req *model.CreateDmUserRequest) (*model.DmUser, error)
		req         *model.CreateDmUserRequest
		want        *model.DmUser
		wantErr     bool
		expectedErr string
	}{
		{
			name: "creates user successfully",
			mockFunc: func(ctx context.Context, req *model.CreateDmUserRequest) (*model.DmUser, error) {
				return &model.DmUser{
					ID:    "user-001",
					Name:  req.Name,
					Email: req.Email,
				}, nil
			},
			req: &model.CreateDmUserRequest{
				Name:  "Test User",
				Email: "test@example.com",
			},
			want: &model.DmUser{
				ID:    "user-001",
				Name:  "Test User",
				Email: "test@example.com",
			},
			wantErr: false,
		},
		{
			name: "returns error when service fails",
			mockFunc: func(ctx context.Context, req *model.CreateDmUserRequest) (*model.DmUser, error) {
				return nil, errors.New("service error")
			},
			req: &model.CreateDmUserRequest{
				Name:  "Test User",
				Email: "test@example.com",
			},
			want:        nil,
			wantErr:     true,
			expectedErr: "service error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockDmUserService{
				CreateDmUserFunc: tt.mockFunc,
			}

			u := NewDmUserUsecase(mockService)
			got, err := u.CreateDmUser(context.Background(), tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want.ID, got.ID)
			assert.Equal(t, tt.want.Name, got.Name)
			assert.Equal(t, tt.want.Email, got.Email)
		})
	}
}

func TestDmUserUsecase_GetDmUser(t *testing.T) {
	tests := []struct {
		name        string
		mockFunc    func(ctx context.Context, id string) (*model.DmUser, error)
		userID      string
		want        *model.DmUser
		wantErr     bool
		expectedErr string
	}{
		{
			name: "gets user successfully",
			mockFunc: func(ctx context.Context, id string) (*model.DmUser, error) {
				return &model.DmUser{
					ID:    id,
					Name:  "Test User",
					Email: "test@example.com",
				}, nil
			},
			userID: "user-001",
			want: &model.DmUser{
				ID:    "user-001",
				Name:  "Test User",
				Email: "test@example.com",
			},
			wantErr: false,
		},
		{
			name: "returns error when user not found",
			mockFunc: func(ctx context.Context, id string) (*model.DmUser, error) {
				return nil, errors.New("user not found")
			},
			userID:      "user-999",
			want:        nil,
			wantErr:     true,
			expectedErr: "user not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockDmUserService{
				GetDmUserFunc: tt.mockFunc,
			}

			u := NewDmUserUsecase(mockService)
			got, err := u.GetDmUser(context.Background(), tt.userID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want.ID, got.ID)
			assert.Equal(t, tt.want.Name, got.Name)
		})
	}
}

func TestDmUserUsecase_ListDmUsers(t *testing.T) {
	tests := []struct {
		name        string
		mockFunc    func(ctx context.Context, limit, offset int) ([]*model.DmUser, error)
		limit       int
		offset      int
		wantCount   int
		wantErr     bool
		expectedErr string
	}{
		{
			name: "lists users successfully",
			mockFunc: func(ctx context.Context, limit, offset int) ([]*model.DmUser, error) {
				return []*model.DmUser{
					{ID: "user-001", Name: "User 1"},
					{ID: "user-002", Name: "User 2"},
				}, nil
			},
			limit:     10,
			offset:    0,
			wantCount: 2,
			wantErr:   false,
		},
		{
			name: "returns error when service fails",
			mockFunc: func(ctx context.Context, limit, offset int) ([]*model.DmUser, error) {
				return nil, errors.New("service error")
			},
			limit:       10,
			offset:      0,
			wantCount:   0,
			wantErr:     true,
			expectedErr: "service error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockDmUserService{
				ListDmUsersFunc: tt.mockFunc,
			}

			u := NewDmUserUsecase(mockService)
			got, err := u.ListDmUsers(context.Background(), tt.limit, tt.offset)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				return
			}

			assert.NoError(t, err)
			assert.Len(t, got, tt.wantCount)
		})
	}
}

func TestDmUserUsecase_UpdateDmUser(t *testing.T) {
	tests := []struct {
		name        string
		mockFunc    func(ctx context.Context, id string, req *model.UpdateDmUserRequest) (*model.DmUser, error)
		userID      string
		req         *model.UpdateDmUserRequest
		want        *model.DmUser
		wantErr     bool
		expectedErr string
	}{
		{
			name: "updates user successfully",
			mockFunc: func(ctx context.Context, id string, req *model.UpdateDmUserRequest) (*model.DmUser, error) {
				return &model.DmUser{
					ID:    id,
					Name:  req.Name,
					Email: req.Email,
				}, nil
			},
			userID: "user-001",
			req: &model.UpdateDmUserRequest{
				Name:  "Updated User",
				Email: "updated@example.com",
			},
			want: &model.DmUser{
				ID:    "user-001",
				Name:  "Updated User",
				Email: "updated@example.com",
			},
			wantErr: false,
		},
		{
			name: "returns error when service fails",
			mockFunc: func(ctx context.Context, id string, req *model.UpdateDmUserRequest) (*model.DmUser, error) {
				return nil, errors.New("service error")
			},
			userID: "user-001",
			req: &model.UpdateDmUserRequest{
				Name: "Updated User",
			},
			want:        nil,
			wantErr:     true,
			expectedErr: "service error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockDmUserService{
				UpdateDmUserFunc: tt.mockFunc,
			}

			u := NewDmUserUsecase(mockService)
			got, err := u.UpdateDmUser(context.Background(), tt.userID, tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want.ID, got.ID)
			assert.Equal(t, tt.want.Name, got.Name)
		})
	}
}

func TestDmUserUsecase_DeleteDmUser(t *testing.T) {
	tests := []struct {
		name        string
		mockFunc    func(ctx context.Context, id string) error
		userID      string
		wantErr     bool
		expectedErr string
	}{
		{
			name: "deletes user successfully",
			mockFunc: func(ctx context.Context, id string) error {
				return nil
			},
			userID:  "user-001",
			wantErr: false,
		},
		{
			name: "returns error when service fails",
			mockFunc: func(ctx context.Context, id string) error {
				return errors.New("service error")
			},
			userID:      "user-001",
			wantErr:     true,
			expectedErr: "service error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockDmUserService{
				DeleteDmUserFunc: tt.mockFunc,
			}

			u := NewDmUserUsecase(mockService)
			err := u.DeleteDmUser(context.Background(), tt.userID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				return
			}

			assert.NoError(t, err)
		})
	}
}
