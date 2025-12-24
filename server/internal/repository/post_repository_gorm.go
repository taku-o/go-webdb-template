package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/example/go-webdb-template/internal/db"
	"github.com/example/go-webdb-template/internal/model"
	"gorm.io/gorm"
)

// PostRepositoryGORM は投稿のデータアクセスを担当（GORM版）
type PostRepositoryGORM struct {
	dbManager *db.GORMManager
}

// NewPostRepositoryGORM は新しいPostRepositoryGORMを作成
func NewPostRepositoryGORM(dbManager *db.GORMManager) *PostRepositoryGORM {
	return &PostRepositoryGORM{
		dbManager: dbManager,
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

	// UserIDをキーとしてShardを決定（同じユーザーのデータは同じShardに配置）
	database, err := r.dbManager.GetGORMByKey(req.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get database: %w", err)
	}

	// GORM APIで作成
	if err := database.WithContext(ctx).Create(post).Error; err != nil {
		return nil, fmt.Errorf("failed to create post: %w", err)
	}

	return post, nil
}

// GetByID はIDで投稿を取得
func (r *PostRepositoryGORM) GetByID(ctx context.Context, id int64, userID int64) (*model.Post, error) {
	database, err := r.dbManager.GetGORMByKey(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get database: %w", err)
	}

	var post model.Post
	if err := database.WithContext(ctx).First(&post, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("post not found: %d", id)
		}
		return nil, fmt.Errorf("failed to get post: %w", err)
	}

	return &post, nil
}

// ListByUserID はユーザーIDで投稿一覧を取得
func (r *PostRepositoryGORM) ListByUserID(ctx context.Context, userID int64, limit, offset int) ([]*model.Post, error) {
	database, err := r.dbManager.GetGORMByKey(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get database: %w", err)
	}

	var posts []*model.Post
	if err := database.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&posts).Error; err != nil {
		return nil, fmt.Errorf("failed to query posts: %w", err)
	}

	return posts, nil
}

// List はすべての投稿を取得（クロスシャードクエリ）
func (r *PostRepositoryGORM) List(ctx context.Context, limit, offset int) ([]*model.Post, error) {
	connections := r.dbManager.GetAllGORMConnections()
	posts := make([]*model.Post, 0)

	// 各シャードから並列にデータを取得
	for _, conn := range connections {
		var shardPosts []*model.Post
		if err := conn.DB.WithContext(ctx).
			Order("created_at DESC").
			Limit(limit).
			Offset(offset).
			Find(&shardPosts).Error; err != nil {
			return nil, fmt.Errorf("failed to query shard %d: %w", conn.ShardID, err)
		}
		posts = append(posts, shardPosts...)
	}

	return posts, nil
}

// GetUserPosts はユーザーと投稿をJOINして取得（クロスシャードクエリ）
func (r *PostRepositoryGORM) GetUserPosts(ctx context.Context, limit, offset int) ([]*model.UserPost, error) {
	connections := r.dbManager.GetAllGORMConnections()
	userPosts := make([]*model.UserPost, 0)

	// GORMでJOINクエリを実行
	for _, conn := range connections {
		var shardUserPosts []*model.UserPost
		err := conn.DB.WithContext(ctx).
			Table("posts p").
			Select(`
				p.id as post_id,
				p.title as post_title,
				p.content as post_content,
				u.id as user_id,
				u.name as user_name,
				u.email as user_email,
				p.created_at
			`).
			Joins("INNER JOIN users u ON p.user_id = u.id").
			Order("p.created_at DESC").
			Limit(limit).
			Offset(offset).
			Find(&shardUserPosts).Error

		if err != nil {
			return nil, fmt.Errorf("failed to query shard %d: %w", conn.ShardID, err)
		}
		userPosts = append(userPosts, shardUserPosts...)
	}

	return userPosts, nil
}

// Update は投稿を更新
func (r *PostRepositoryGORM) Update(ctx context.Context, id int64, userID int64, req *model.UpdatePostRequest) (*model.Post, error) {
	database, err := r.dbManager.GetGORMByKey(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get database: %w", err)
	}

	updates := make(map[string]interface{})
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Content != "" {
		updates["content"] = req.Content
	}
	updates["updated_at"] = time.Now()

	result := database.WithContext(ctx).
		Model(&model.Post{}).
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
	database, err := r.dbManager.GetGORMByKey(userID)
	if err != nil {
		return fmt.Errorf("failed to get database: %w", err)
	}

	result := database.WithContext(ctx).
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
