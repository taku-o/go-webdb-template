package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/taku-o/go-webdb-template/internal/db"
	"github.com/taku-o/go-webdb-template/internal/model"
	"gorm.io/gorm"
)

// PostRepositoryGORM は投稿のデータアクセスを担当（GORM版）
type PostRepositoryGORM struct {
	groupManager  *db.GroupManager
	tableSelector *db.TableSelector
}

// NewPostRepositoryGORM は新しいPostRepositoryGORMを作成
func NewPostRepositoryGORM(groupManager *db.GroupManager) *PostRepositoryGORM {
	return &PostRepositoryGORM{
		groupManager:  groupManager,
		tableSelector: db.NewTableSelector(32, 8),
	}
}

// Create は投稿を作成
func (r *PostRepositoryGORM) Create(ctx context.Context, req *model.CreatePostRequest) (*model.Post, error) {
	post := &model.Post{
		UserID:  req.UserID,
		Title:   req.Title,
		Content: req.Content,
	}

	// ID生成（タイムスタンプベース、既存ロジック維持）
	post.ID = time.Now().UnixNano()

	// UserIDをキーとしてテーブル/DBを決定（同じユーザーのデータは同じテーブルに配置）
	tableName := r.tableSelector.GetTableName("posts", req.UserID)

	// 接続の取得
	conn, err := r.groupManager.GetShardingConnectionByID(req.UserID, "posts")
	if err != nil {
		return nil, fmt.Errorf("failed to get sharding connection: %w", err)
	}

	// GORM APIで作成（動的テーブル名を使用）
	if err := conn.DB.WithContext(ctx).Table(tableName).Create(post).Error; err != nil {
		return nil, fmt.Errorf("failed to create post: %w", err)
	}

	return post, nil
}

// GetByID はIDで投稿を取得
func (r *PostRepositoryGORM) GetByID(ctx context.Context, id int64, userID int64) (*model.Post, error) {
	// UserIDをキーとしてテーブル/DBを決定
	tableName := r.tableSelector.GetTableName("posts", userID)

	// 接続の取得
	conn, err := r.groupManager.GetShardingConnectionByID(userID, "posts")
	if err != nil {
		return nil, fmt.Errorf("failed to get sharding connection: %w", err)
	}

	var post model.Post
	if err := conn.DB.WithContext(ctx).Table(tableName).First(&post, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("post not found: %d", id)
		}
		return nil, fmt.Errorf("failed to get post: %w", err)
	}

	return &post, nil
}

// ListByUserID はユーザーIDで投稿一覧を取得
func (r *PostRepositoryGORM) ListByUserID(ctx context.Context, userID int64, limit, offset int) ([]*model.Post, error) {
	// UserIDをキーとしてテーブル/DBを決定
	tableName := r.tableSelector.GetTableName("posts", userID)

	// 接続の取得
	conn, err := r.groupManager.GetShardingConnectionByID(userID, "posts")
	if err != nil {
		return nil, fmt.Errorf("failed to get sharding connection: %w", err)
	}

	var posts []*model.Post
	if err := conn.DB.WithContext(ctx).
		Table(tableName).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&posts).Error; err != nil {
		return nil, fmt.Errorf("failed to query posts: %w", err)
	}

	return posts, nil
}

// List はすべての投稿を取得（クロステーブルクエリ）
func (r *PostRepositoryGORM) List(ctx context.Context, limit, offset int) ([]*model.Post, error) {
	posts := make([]*model.Post, 0)

	// テーブル数分ループして各テーブルからデータを取得
	tableCount := r.tableSelector.GetTableCount()
	for tableNum := 0; tableNum < tableCount; tableNum++ {
		// テーブル番号から接続を取得
		conn, err := r.groupManager.GetShardingConnection(tableNum)
		if err != nil {
			return nil, fmt.Errorf("failed to get connection for table %d: %w", tableNum, err)
		}

		tableName := fmt.Sprintf("posts_%03d", tableNum)

		var tablePosts []*model.Post
		if err := conn.DB.WithContext(ctx).
			Table(tableName).
			Order("created_at DESC").
			Limit(limit).
			Offset(offset).
			Find(&tablePosts).Error; err != nil {
			return nil, fmt.Errorf("failed to query table %s: %w", tableName, err)
		}
		posts = append(posts, tablePosts...)
	}

	return posts, nil
}

// GetUserPosts はユーザーと投稿をJOINして取得（クロステーブルクエリ）
func (r *PostRepositoryGORM) GetUserPosts(ctx context.Context, limit, offset int) ([]*model.UserPost, error) {
	userPosts := make([]*model.UserPost, 0)

	// テーブル数分ループして各テーブルからデータを取得
	tableCount := r.tableSelector.GetTableCount()
	for tableNum := 0; tableNum < tableCount; tableNum++ {
		// テーブル番号から接続を取得
		conn, err := r.groupManager.GetShardingConnection(tableNum)
		if err != nil {
			return nil, fmt.Errorf("failed to get connection for table %d: %w", tableNum, err)
		}

		postsTable := fmt.Sprintf("posts_%03d", tableNum)
		usersTable := fmt.Sprintf("users_%03d", tableNum)

		var tableUserPosts []*model.UserPost
		err = conn.DB.WithContext(ctx).
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
			Find(&tableUserPosts).Error

		if err != nil {
			return nil, fmt.Errorf("failed to query table %s: %w", postsTable, err)
		}
		userPosts = append(userPosts, tableUserPosts...)
	}

	return userPosts, nil
}

// Update は投稿を更新
func (r *PostRepositoryGORM) Update(ctx context.Context, id int64, userID int64, req *model.UpdatePostRequest) (*model.Post, error) {
	// UserIDをキーとしてテーブル/DBを決定
	tableName := r.tableSelector.GetTableName("posts", userID)

	// 接続の取得
	conn, err := r.groupManager.GetShardingConnectionByID(userID, "posts")
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

	result := conn.DB.WithContext(ctx).
		Table(tableName).
		Where("id = ? AND user_id = ?", id, userID).
		Updates(updates)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to update post: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("post not found: %d", id)
	}

	return r.GetByID(ctx, id, userID)
}

// Delete は投稿を削除
func (r *PostRepositoryGORM) Delete(ctx context.Context, id int64, userID int64) error {
	// UserIDをキーとしてテーブル/DBを決定
	tableName := r.tableSelector.GetTableName("posts", userID)

	// 接続の取得
	conn, err := r.groupManager.GetShardingConnectionByID(userID, "posts")
	if err != nil {
		return fmt.Errorf("failed to get sharding connection: %w", err)
	}

	result := conn.DB.WithContext(ctx).
		Table(tableName).
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&model.Post{})

	if result.Error != nil {
		return fmt.Errorf("failed to delete post: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("post not found: %d", id)
	}

	return nil
}
