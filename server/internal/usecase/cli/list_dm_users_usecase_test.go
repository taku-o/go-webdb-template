package cli

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/taku-o/go-webdb-template/internal/model"
)

// MockDmUserServiceInterface はDmUserServiceInterfaceのモック
type MockDmUserServiceInterface struct {
	ListDmUsersFunc func(ctx context.Context, limit, offset int) ([]*model.DmUser, error)
}

func (m *MockDmUserServiceInterface) CreateDmUser(ctx context.Context, req *model.CreateDmUserRequest) (*model.DmUser, error) {
	return nil, nil
}

func (m *MockDmUserServiceInterface) GetDmUser(ctx context.Context, id string) (*model.DmUser, error) {
	return nil, nil
}

func (m *MockDmUserServiceInterface) ListDmUsers(ctx context.Context, limit, offset int) ([]*model.DmUser, error) {
	if m.ListDmUsersFunc != nil {
		return m.ListDmUsersFunc(ctx, limit, offset)
	}
	return nil, nil
}

func (m *MockDmUserServiceInterface) UpdateDmUser(ctx context.Context, id string, req *model.UpdateDmUserRequest) (*model.DmUser, error) {
	return nil, nil
}

func (m *MockDmUserServiceInterface) DeleteDmUser(ctx context.Context, id string) error {
	return nil
}

func TestListDmUsersUsecase_ListDmUsers(t *testing.T) {
	tests := []struct {
		name        string
		limit       int
		offset      int
		mockFunc    func(ctx context.Context, limit, offset int) ([]*model.DmUser, error)
		wantUsers   []*model.DmUser
		wantError   bool
		expectedErr string
	}{
		{
			name:   "success with users",
			limit:  20,
			offset: 0,
			mockFunc: func(ctx context.Context, limit, offset int) ([]*model.DmUser, error) {
				return []*model.DmUser{
					{ID: "1", Name: "User 1", Email: "user1@example.com"},
					{ID: "2", Name: "User 2", Email: "user2@example.com"},
				}, nil
			},
			wantUsers: []*model.DmUser{
				{ID: "1", Name: "User 1", Email: "user1@example.com"},
				{ID: "2", Name: "User 2", Email: "user2@example.com"},
			},
			wantError: false,
		},
		{
			name:   "success with empty list",
			limit:  20,
			offset: 0,
			mockFunc: func(ctx context.Context, limit, offset int) ([]*model.DmUser, error) {
				return []*model.DmUser{}, nil
			},
			wantUsers: []*model.DmUser{},
			wantError: false,
		},
		{
			name:   "service error",
			limit:  20,
			offset: 0,
			mockFunc: func(ctx context.Context, limit, offset int) ([]*model.DmUser, error) {
				return nil, errors.New("database error")
			},
			wantUsers:   nil,
			wantError:   true,
			expectedErr: "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockDmUserServiceInterface{
				ListDmUsersFunc: tt.mockFunc,
			}

			usecase := NewListDmUsersUsecase(mockService)

			ctx := context.Background()
			gotUsers, err := usecase.ListDmUsers(ctx, tt.limit, tt.offset)

			if tt.wantError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				assert.Nil(t, gotUsers)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantUsers, gotUsers)
			}
		})
	}
}
