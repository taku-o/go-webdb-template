package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/taku-o/go-webdb-template/internal/db"
	"github.com/taku-o/go-webdb-template/internal/model"
	"github.com/taku-o/go-webdb-template/internal/util/idgen"
	"gorm.io/gorm"
)

// DmPostRepositoryGORM は投稿のデータアクセスを担当（GORM版）
type DmPostRepositoryGORM struct {
	groupManager  *db.GroupManager
	tableSelector *db.TableSelector
}

// NewDmPostRepositoryGORM は新しいDmPostRepositoryGORMを作成
func NewDmPostRepositoryGORM(groupManager *db.GroupManager) *DmPostRepositoryGORM {
	return &DmPostRepositoryGORM{
		groupManager:  groupManager,
		tableSelector: db.NewTableSelector(db.DBShardingTableCount, db.DBShardingTablesPerDB),
	}
}

// Create は投稿を作成
func (r *DmPostRepositoryGORM) Create(ctx context.Context, req *model.CreateDmPostRequest) (*model.DmPost, error) {
	// ID生成（UUIDv7）
	id, err := idgen.GenerateUUIDv7()
	if err != nil {
		return nil, fmt.Errorf("failed to generate ID: %w", err)
	}

	post := &model.DmPost{
		ID:      id,
		UserID:  req.UserID,
		Title:   req.Title,
		Content: req.Content,
	}

	// UserIDをキーとしてテーブル/DBを決定（同じユーザーのデータは同じテーブルに配置）
	tableName, err := r.tableSelector.GetTableNameFromUUID("dm_posts", req.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get table name: %w", err)
	}

	// 接続の取得
	conn, err := r.groupManager.GetShardingConnectionByUUID(req.UserID, "dm_posts")
	if err != nil {
		return nil, fmt.Errorf("failed to get sharding connection: %w", err)
	}

	// リトライ機能付きでGORM APIで作成（動的テーブル名を使用）
	err = db.ExecuteWithRetry(func() error {
		return conn.DB.WithContext(ctx).Table(tableName).Create(post).Error
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create post: %w", err)
	}

	return post, nil
}

// GetByID はIDで投稿を取得
func (r *DmPostRepositoryGORM) GetByID(ctx context.Context, id string, userID string) (*model.DmPost, error) {
	// UserIDをキーとしてテーブル/DBを決定
	tableName, err := r.tableSelector.GetTableNameFromUUID("dm_posts", userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get table name: %w", err)
	}

	// 接続の取得
	conn, err := r.groupManager.GetShardingConnectionByUUID(userID, "dm_posts")
	if err != nil {
		return nil, fmt.Errorf("failed to get sharding connection: %w", err)
	}

	var post model.DmPost
	// リトライ機能付きでクエリ実行
	err = db.ExecuteWithRetry(func() error {
		return conn.DB.WithContext(ctx).Table(tableName).Where("id = ?", id).First(&post).Error
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("post not found: %s", id)
		}
		return nil, fmt.Errorf("failed to get post: %w", err)
	}

	return &post, nil
}

// ListByUserID はユーザーIDで投稿一覧を取得
func (r *DmPostRepositoryGORM) ListByUserID(ctx context.Context, userID string, limit, offset int) ([]*model.DmPost, error) {
	// UserIDをキーとしてテーブル/DBを決定
	tableName, err := r.tableSelector.GetTableNameFromUUID("dm_posts", userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get table name: %w", err)
	}

	// 接続の取得
	conn, err := r.groupManager.GetShardingConnectionByUUID(userID, "dm_posts")
	if err != nil {
		return nil, fmt.Errorf("failed to get sharding connection: %w", err)
	}

	var posts []*model.DmPost
	// リトライ機能付きでクエリ実行
	err = db.ExecuteWithRetry(func() error {
		return conn.DB.WithContext(ctx).
			Table(tableName).
			Where("user_id = ?", userID).
			Order("created_at DESC").
			Limit(limit).
			Offset(offset).
			Find(&posts).Error
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query posts: %w", err)
	}

	return posts, nil
}

// List はすべての投稿を取得（クロステーブルクエリ）
func (r *DmPostRepositoryGORM) List(ctx context.Context, limit, offset int) ([]*model.DmPost, error) {
	posts := make([]*model.DmPost, 0)

	// テーブル数分ループして各テーブルからデータを取得
	tableCount := r.tableSelector.GetTableCount()
	for tableNum := 0; tableNum < tableCount; tableNum++ {
		// テーブル番号から接続を取得
		conn, err := r.groupManager.GetShardingConnection(tableNum)
		if err != nil {
			return nil, fmt.Errorf("failed to get connection for table %d: %w", tableNum, err)
		}

		tableName := fmt.Sprintf("dm_posts_%03d", tableNum)

		var tablePosts []*model.DmPost
		// リトライ機能付きでクエリ実行
		err = db.ExecuteWithRetry(func() error {
			return conn.DB.WithContext(ctx).
				Table(tableName).
				Order("created_at DESC").
				Limit(limit).
				Offset(offset).
				Find(&tablePosts).Error
		})
		if err != nil {
			return nil, fmt.Errorf("failed to query table %s: %w", tableName, err)
		}
		posts = append(posts, tablePosts...)
	}

	return posts, nil
}

// GetUserPosts はユーザーと投稿をJOINして取得（クロステーブルクエリ）
func (r *DmPostRepositoryGORM) GetUserPosts(ctx context.Context, limit, offset int) ([]*model.DmUserPost, error) {
	userPosts := make([]*model.DmUserPost, 0)

	// テーブル数分ループして各テーブルからデータを取得
	tableCount := r.tableSelector.GetTableCount()
	for tableNum := 0; tableNum < tableCount; tableNum++ {
		// テーブル番号から接続を取得
		conn, err := r.groupManager.GetShardingConnection(tableNum)
		if err != nil {
			return nil, fmt.Errorf("failed to get connection for table %d: %w", tableNum, err)
		}

		postsTable := fmt.Sprintf("dm_posts_%03d", tableNum)
		usersTable := fmt.Sprintf("dm_users_%03d", tableNum)

		var tableDmUserPosts []*model.DmUserPost
		// リトライ機能付きでクエリ実行
		err = db.ExecuteWithRetry(func() error {
			return conn.DB.WithContext(ctx).
				Table(postsTable+" p").
				Select(`
					p.id as post_id,
					p.title as post_title,
					p.content as post_content,
					u.id as user_id,
					u.name as user_name,
					u.email as user_email,
					p.created_at
				`).
				Joins(fmt.Sprintf("INNER JOIN %s u ON p.user_id = u.id", usersTable)).
				Order("p.created_at DESC").
				Limit(limit).
				Offset(offset).
				Find(&tableDmUserPosts).Error
		})
		if err != nil {
			return nil, fmt.Errorf("failed to query table %s: %w", postsTable, err)
		}
		userPosts = append(userPosts, tableDmUserPosts...)
	}

	return userPosts, nil
}

// Update は投稿を更新
func (r *DmPostRepositoryGORM) Update(ctx context.Context, id string, userID string, req *model.UpdateDmPostRequest) (*model.DmPost, error) {
	// UserIDをキーとしてテーブル/DBを決定
	tableName, err := r.tableSelector.GetTableNameFromUUID("dm_posts", userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get table name: %w", err)
	}

	// 接続の取得
	conn, err := r.groupManager.GetShardingConnectionByUUID(userID, "dm_posts")
	if err != nil {
		return nil, fmt.Errorf("failed to get sharding connection: %w", err)
	}

	updates := make(map[string]interface{})
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Content != "" {
		updates["content"] = req.Content
	}
	updates["updated_at"] = time.Now()

	var result *gorm.DB
	// リトライ機能付きでクエリ実行
	err = db.ExecuteWithRetry(func() error {
		result = conn.DB.WithContext(ctx).
			Table(tableName).
			Where("id = ? AND user_id = ?", id, userID).
			Updates(updates)
		return result.Error
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update post: %w", err)
	}
	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("post not found: %s", id)
	}

	return r.GetByID(ctx, id, userID)
}

// Delete は投稿を削除
func (r *DmPostRepositoryGORM) Delete(ctx context.Context, id string, userID string) error {
	// UserIDをキーとしてテーブル/DBを決定
	tableName, err := r.tableSelector.GetTableNameFromUUID("dm_posts", userID)
	if err != nil {
		return fmt.Errorf("failed to get table name: %w", err)
	}

	// 接続の取得
	conn, err := r.groupManager.GetShardingConnectionByUUID(userID, "dm_posts")
	if err != nil {
		return fmt.Errorf("failed to get sharding connection: %w", err)
	}

	var result *gorm.DB
	// リトライ機能付きでクエリ実行
	err = db.ExecuteWithRetry(func() error {
		result = conn.DB.WithContext(ctx).
			Table(tableName).
			Where("id = ? AND user_id = ?", id, userID).
			Delete(&model.DmPost{})
		return result.Error
	})
	if err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("post not found: %s", id)
	}

	return nil
}
