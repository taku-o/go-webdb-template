package cli

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// MockGenerateSampleServiceInterface はGenerateSampleServiceInterfaceのモック
type MockGenerateSampleServiceInterface struct {
	GenerateDmUsersFunc func(ctx context.Context, totalCount int) ([]string, error)
	GenerateDmPostsFunc func(ctx context.Context, dmUserIDs []string, totalCount int) error
	GenerateDmNewsFunc  func(ctx context.Context, totalCount int) error
}

func (m *MockGenerateSampleServiceInterface) GenerateDmUsers(ctx context.Context, totalCount int) ([]string, error) {
	if m.GenerateDmUsersFunc != nil {
		return m.GenerateDmUsersFunc(ctx, totalCount)
	}
	return []string{}, nil
}

func (m *MockGenerateSampleServiceInterface) GenerateDmPosts(ctx context.Context, dmUserIDs []string, totalCount int) error {
	if m.GenerateDmPostsFunc != nil {
		return m.GenerateDmPostsFunc(ctx, dmUserIDs, totalCount)
	}
	return nil
}

func (m *MockGenerateSampleServiceInterface) GenerateDmNews(ctx context.Context, totalCount int) error {
	if m.GenerateDmNewsFunc != nil {
		return m.GenerateDmNewsFunc(ctx, totalCount)
	}
	return nil
}

func TestGenerateSampleUsecase_GenerateSampleData(t *testing.T) {
	tests := []struct {
		name                string
		generateDmUsersFunc func(ctx context.Context, totalCount int) ([]string, error)
		generateDmPostsFunc func(ctx context.Context, dmUserIDs []string, totalCount int) error
		generateDmNewsFunc  func(ctx context.Context, totalCount int) error
		totalCount          int
		wantError           bool
		expectedErr         string
	}{
		{
			name: "success",
			generateDmUsersFunc: func(ctx context.Context, totalCount int) ([]string, error) {
				return []string{"user1", "user2"}, nil
			},
			generateDmPostsFunc: func(ctx context.Context, dmUserIDs []string, totalCount int) error {
				return nil
			},
			generateDmNewsFunc: func(ctx context.Context, totalCount int) error {
				return nil
			},
			totalCount: 10,
			wantError:  false,
		},
		{
			name: "zero count",
			generateDmUsersFunc: func(ctx context.Context, totalCount int) ([]string, error) {
				return []string{}, nil
			},
			generateDmPostsFunc: func(ctx context.Context, dmUserIDs []string, totalCount int) error {
				return nil
			},
			generateDmNewsFunc: func(ctx context.Context, totalCount int) error {
				return nil
			},
			totalCount: 0,
			wantError:  false,
		},
		{
			name: "generateDmUsers error",
			generateDmUsersFunc: func(ctx context.Context, totalCount int) ([]string, error) {
				return nil, errors.New("failed to generate users")
			},
			generateDmPostsFunc: nil,
			generateDmNewsFunc:  nil,
			totalCount:          10,
			wantError:           true,
			expectedErr:         "failed to generate users",
		},
		{
			name: "generateDmPosts error",
			generateDmUsersFunc: func(ctx context.Context, totalCount int) ([]string, error) {
				return []string{"user1", "user2"}, nil
			},
			generateDmPostsFunc: func(ctx context.Context, dmUserIDs []string, totalCount int) error {
				return errors.New("failed to generate posts")
			},
			generateDmNewsFunc: nil,
			totalCount:         10,
			wantError:          true,
			expectedErr:        "failed to generate posts",
		},
		{
			name: "generateDmNews error",
			generateDmUsersFunc: func(ctx context.Context, totalCount int) ([]string, error) {
				return []string{"user1", "user2"}, nil
			},
			generateDmPostsFunc: func(ctx context.Context, dmUserIDs []string, totalCount int) error {
				return nil
			},
			generateDmNewsFunc: func(ctx context.Context, totalCount int) error {
				return errors.New("failed to generate news")
			},
			totalCount:  10,
			wantError:   true,
			expectedErr: "failed to generate news",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockGenerateSampleServiceInterface{
				GenerateDmUsersFunc: tt.generateDmUsersFunc,
				GenerateDmPostsFunc: tt.generateDmPostsFunc,
				GenerateDmNewsFunc:  tt.generateDmNewsFunc,
			}

			usecase := NewGenerateSampleUsecase(mockService)

			ctx := context.Background()
			err := usecase.GenerateSampleData(ctx, tt.totalCount)

			if tt.wantError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNewGenerateSampleUsecase(t *testing.T) {
	mockService := &MockGenerateSampleServiceInterface{}

	usecase := NewGenerateSampleUsecase(mockService)

	assert.NotNil(t, usecase)
}
