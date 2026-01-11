package admin

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/taku-o/go-webdb-template/internal/model"
)

// MockDmUserService はDmUserServiceInterfaceのモック
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

func TestDmUserRegisterUsecase_RegisterDmUser(t *testing.T) {
	tests := []struct {
		name                 string
		inputName            string
		inputEmail           string
		checkEmailExistsFunc func(ctx context.Context, email string) (bool, error)
		createDmUserFunc     func(ctx context.Context, req *model.CreateDmUserRequest) (*model.DmUser, error)
		wantID               string
		wantErr              bool
		wantErrContains      string
	}{
		{
			name:       "正常系: ユーザー登録が成功する",
			inputName:  "Test User",
			inputEmail: "test@example.com",
			checkEmailExistsFunc: func(ctx context.Context, email string) (bool, error) {
				return false, nil
			},
			createDmUserFunc: func(ctx context.Context, req *model.CreateDmUserRequest) (*model.DmUser, error) {
				return &model.DmUser{
					ID:    "user-001",
					Name:  req.Name,
					Email: req.Email,
				}, nil
			},
			wantID:  "user-001",
			wantErr: false,
		},
		{
			name:       "異常系: メールアドレスが既に登録されている",
			inputName:  "Test User",
			inputEmail: "existing@example.com",
			checkEmailExistsFunc: func(ctx context.Context, email string) (bool, error) {
				return true, nil
			},
			createDmUserFunc: nil,
			wantID:           "",
			wantErr:          true,
			wantErrContains:  "このメールアドレスは既に登録されています",
		},
		{
			name:       "異常系: メールアドレスチェックでエラーが発生",
			inputName:  "Test User",
			inputEmail: "test@example.com",
			checkEmailExistsFunc: func(ctx context.Context, email string) (bool, error) {
				return false, errors.New("database connection error")
			},
			createDmUserFunc: nil,
			wantID:           "",
			wantErr:          true,
			wantErrContains:  "failed to check email",
		},
		{
			name:       "異常系: ユーザー作成でエラーが発生",
			inputName:  "Test User",
			inputEmail: "test@example.com",
			checkEmailExistsFunc: func(ctx context.Context, email string) (bool, error) {
				return false, nil
			},
			createDmUserFunc: func(ctx context.Context, req *model.CreateDmUserRequest) (*model.DmUser, error) {
				return nil, errors.New("failed to create user")
			},
			wantID:          "",
			wantErr:         true,
			wantErrContains: "failed to create user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockDmUserService{
				CheckEmailExistsFunc: tt.checkEmailExistsFunc,
				CreateDmUserFunc:     tt.createDmUserFunc,
			}

			u := NewDmUserRegisterUsecase(mockService)
			gotID, err := u.RegisterDmUser(context.Background(), tt.inputName, tt.inputEmail)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.wantErrContains != "" {
					assert.Contains(t, err.Error(), tt.wantErrContains)
				}
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.wantID, gotID)
		})
	}
}
