package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/taku-o/go-webdb-template/internal/db"
	"github.com/taku-o/go-webdb-template/internal/model"
	"github.com/taku-o/go-webdb-template/internal/service"
)

// MockDmUserRepository はDmUserRepositoryのモック
type MockDmUserRepository struct {
	InsertDmUsersBatchFunc func(ctx context.Context, tableName string, dmUsers []*model.DmUser) error
}

func (m *MockDmUserRepository) InsertDmUsersBatch(ctx context.Context, tableName string, dmUsers []*model.DmUser) error {
	if m.InsertDmUsersBatchFunc != nil {
		return m.InsertDmUsersBatchFunc(ctx, tableName, dmUsers)
	}
	return nil
}

// MockDmPostRepository はDmPostRepositoryのモック
type MockDmPostRepository struct {
	InsertDmPostsBatchFunc func(ctx context.Context, tableName string, dmPosts []*model.DmPost) error
}

func (m *MockDmPostRepository) InsertDmPostsBatch(ctx context.Context, tableName string, dmPosts []*model.DmPost) error {
	if m.InsertDmPostsBatchFunc != nil {
		return m.InsertDmPostsBatchFunc(ctx, tableName, dmPosts)
	}
	return nil
}

// MockDmNewsRepository はDmNewsRepositoryのモック
type MockDmNewsRepository struct {
	InsertDmNewsBatchFunc func(ctx context.Context, dmNews []*model.DmNews) error
}

func (m *MockDmNewsRepository) InsertDmNewsBatch(ctx context.Context, dmNews []*model.DmNews) error {
	if m.InsertDmNewsBatchFunc != nil {
		return m.InsertDmNewsBatchFunc(ctx, dmNews)
	}
	return nil
}

func TestGenerateSampleService_GenerateDmUsers(t *testing.T) {
	tests := []struct {
		name        string
		mockFunc    func(ctx context.Context, tableName string, dmUsers []*model.DmUser) error
		totalCount  int
		wantError   bool
		expectedLen int
	}{
		{
			name: "success",
			mockFunc: func(ctx context.Context, tableName string, dmUsers []*model.DmUser) error {
				return nil
			},
			totalCount:  10,
			wantError:   false,
			expectedLen: 10,
		},
		{
			name: "zero count",
			mockFunc: func(ctx context.Context, tableName string, dmUsers []*model.DmUser) error {
				return nil
			},
			totalCount:  0,
			wantError:   false,
			expectedLen: 0,
		},
		{
			name: "repository error",
			mockFunc: func(ctx context.Context, tableName string, dmUsers []*model.DmUser) error {
				return errors.New("failed to insert batch")
			},
			totalCount:  10,
			wantError:   true,
			expectedLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepo := &MockDmUserRepository{
				InsertDmUsersBatchFunc: tt.mockFunc,
			}
			mockPostRepo := &MockDmPostRepository{}
			mockNewsRepo := &MockDmNewsRepository{}

			tableSelector := db.NewTableSelector(32, 8)

			svc := service.NewGenerateSampleService(
				mockUserRepo,
				mockPostRepo,
				mockNewsRepo,
				tableSelector,
			)

			ctx := context.Background()
			got, err := svc.GenerateDmUsers(ctx, tt.totalCount)

			if tt.wantError {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Len(t, got, tt.expectedLen)
			}
		})
	}
}

func TestGenerateSampleService_GenerateDmPosts(t *testing.T) {
	tests := []struct {
		name       string
		mockFunc   func(ctx context.Context, tableName string, dmPosts []*model.DmPost) error
		dmUserIDs  []string
		totalCount int
		wantError  bool
	}{
		{
			name: "success",
			mockFunc: func(ctx context.Context, tableName string, dmPosts []*model.DmPost) error {
				return nil
			},
			dmUserIDs:  []string{"0194e79d-4fb6-7af2-a20b-f9a8b29a2d58", "0194e79d-4fb6-7af2-a20b-f9a8b29a2d59"},
			totalCount: 5,
			wantError:  false,
		},
		{
			name: "empty user IDs",
			mockFunc: func(ctx context.Context, tableName string, dmPosts []*model.DmPost) error {
				return nil
			},
			dmUserIDs:  []string{},
			totalCount: 5,
			wantError:  true,
		},
		{
			name: "zero count",
			mockFunc: func(ctx context.Context, tableName string, dmPosts []*model.DmPost) error {
				return nil
			},
			dmUserIDs:  []string{"0194e79d-4fb6-7af2-a20b-f9a8b29a2d58"},
			totalCount: 0,
			wantError:  false,
		},
		{
			name: "repository error",
			mockFunc: func(ctx context.Context, tableName string, dmPosts []*model.DmPost) error {
				return errors.New("failed to insert batch")
			},
			dmUserIDs:  []string{"0194e79d-4fb6-7af2-a20b-f9a8b29a2d58"},
			totalCount: 5,
			wantError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepo := &MockDmUserRepository{}
			mockPostRepo := &MockDmPostRepository{
				InsertDmPostsBatchFunc: tt.mockFunc,
			}
			mockNewsRepo := &MockDmNewsRepository{}

			tableSelector := db.NewTableSelector(32, 8)

			svc := service.NewGenerateSampleService(
				mockUserRepo,
				mockPostRepo,
				mockNewsRepo,
				tableSelector,
			)

			ctx := context.Background()
			err := svc.GenerateDmPosts(ctx, tt.dmUserIDs, tt.totalCount)

			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGenerateSampleService_GenerateDmNews(t *testing.T) {
	tests := []struct {
		name       string
		mockFunc   func(ctx context.Context, dmNews []*model.DmNews) error
		totalCount int
		wantError  bool
	}{
		{
			name: "success",
			mockFunc: func(ctx context.Context, dmNews []*model.DmNews) error {
				return nil
			},
			totalCount: 10,
			wantError:  false,
		},
		{
			name: "zero count",
			mockFunc: func(ctx context.Context, dmNews []*model.DmNews) error {
				return nil
			},
			totalCount: 0,
			wantError:  false,
		},
		{
			name: "repository error",
			mockFunc: func(ctx context.Context, dmNews []*model.DmNews) error {
				return errors.New("failed to insert batch")
			},
			totalCount: 10,
			wantError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepo := &MockDmUserRepository{}
			mockPostRepo := &MockDmPostRepository{}
			mockNewsRepo := &MockDmNewsRepository{
				InsertDmNewsBatchFunc: tt.mockFunc,
			}

			tableSelector := db.NewTableSelector(32, 8)

			svc := service.NewGenerateSampleService(
				mockUserRepo,
				mockPostRepo,
				mockNewsRepo,
				tableSelector,
			)

			ctx := context.Background()
			err := svc.GenerateDmNews(ctx, tt.totalCount)

			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNewGenerateSampleService(t *testing.T) {
	mockUserRepo := &MockDmUserRepository{}
	mockPostRepo := &MockDmPostRepository{}
	mockNewsRepo := &MockDmNewsRepository{}
	tableSelector := db.NewTableSelector(32, 8)

	svc := service.NewGenerateSampleService(
		mockUserRepo,
		mockPostRepo,
		mockNewsRepo,
		tableSelector,
	)

	assert.NotNil(t, svc)
}
