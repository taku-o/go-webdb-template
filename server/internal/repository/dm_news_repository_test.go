package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/taku-o/go-webdb-template/internal/model"
	"github.com/taku-o/go-webdb-template/internal/repository"
	"github.com/taku-o/go-webdb-template/test/testutil"
)

func TestDmNewsRepository_InsertDmNewsBatch(t *testing.T) {
	groupManager := testutil.SetupTestGroupManager(t, 4, 8)
	defer testutil.CleanupTestGroupManager(groupManager)

	dmNewsRepo := repository.NewDmNewsRepository(groupManager)
	ctx := context.Background()

	tests := []struct {
		name    string
		dmNews  []*model.DmNews
		wantErr bool
	}{
		{
			name:    "empty slice",
			dmNews:  []*model.DmNews{},
			wantErr: false,
		},
		{
			name: "single news",
			dmNews: func() []*model.DmNews {
				authorID := int64(12345)
				publishedAt := time.Now()
				return []*model.DmNews{
					{
						Title:       "Batch News 1",
						Content:     "Batch content 1",
						AuthorID:    &authorID,
						PublishedAt: &publishedAt,
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
					},
				}
			}(),
			wantErr: false,
		},
		{
			name: "multiple news",
			dmNews: func() []*model.DmNews {
				var news []*model.DmNews
				for i := 0; i < 3; i++ {
					authorID := int64(10000 + i)
					publishedAt := time.Now()
					news = append(news, &model.DmNews{
						Title:       "Batch News",
						Content:     "Batch content",
						AuthorID:    &authorID,
						PublishedAt: &publishedAt,
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
					})
				}
				return news
			}(),
			wantErr: false,
		},
		{
			name: "news without optional fields",
			dmNews: func() []*model.DmNews {
				return []*model.DmNews{
					{
						Title:     "News without optional",
						Content:   "Content without optional",
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
				}
			}(),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := dmNewsRepo.InsertDmNewsBatch(ctx, tt.dmNews)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				// 注: dm_newsテーブルのクリーンアップはmaster接続で行う必要がある
				// GORMのCreateInBatchesを使用しているため、IDはauto incrementで生成される
				// 本テストでは挿入のみを確認し、クリーンアップは行わない
			}
		})
	}
}

func TestNewDmNewsRepository(t *testing.T) {
	groupManager := testutil.SetupTestGroupManager(t, 4, 8)
	defer testutil.CleanupTestGroupManager(groupManager)

	dmNewsRepo := repository.NewDmNewsRepository(groupManager)
	assert.NotNil(t, dmNewsRepo)
}
