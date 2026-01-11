package repository

import (
	"context"
	"fmt"

	"github.com/taku-o/go-webdb-template/internal/db"
	"github.com/taku-o/go-webdb-template/internal/model"
)

// DmNewsRepository はニュースのデータアクセスを担当
type DmNewsRepository struct {
	groupManager *db.GroupManager
}

// NewDmNewsRepository は新しいDmNewsRepositoryを作成
func NewDmNewsRepository(groupManager *db.GroupManager) *DmNewsRepository {
	return &DmNewsRepository{
		groupManager: groupManager,
	}
}

// InsertDmNewsBatch はdm_newsテーブルにバッチでデータを挿入
func (r *DmNewsRepository) InsertDmNewsBatch(ctx context.Context, dmNews []*model.DmNews) error {
	if len(dmNews) == 0 {
		return nil
	}

	const batchSize = 500

	// master接続を取得
	conn, err := r.groupManager.GetMasterConnection()
	if err != nil {
		return fmt.Errorf("failed to get master connection: %w", err)
	}

	// バッチサイズを考慮して分割
	for i := 0; i < len(dmNews); i += batchSize {
		end := i + batchSize
		if end > len(dmNews) {
			end = len(dmNews)
		}
		batch := dmNews[i:end]

		// GORMのCreateInBatchesを使用（固定テーブル名）
		err = db.ExecuteWithRetry(func() error {
			return conn.DB.WithContext(ctx).Table("dm_news").CreateInBatches(batch, len(batch)).Error
		})
		if err != nil {
			return fmt.Errorf("failed to insert batch: %w", err)
		}
	}

	return nil
}
