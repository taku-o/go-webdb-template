package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/example/go-webdb-template/internal/db"
	"github.com/example/go-webdb-template/internal/model"
)

// PostRepository は投稿のデータアクセスを担当
type PostRepository struct {
	dbManager *db.Manager
}

// NewPostRepository は新しいPostRepositoryを作成
func NewPostRepository(dbManager *db.Manager) *PostRepository {
	return &PostRepository{
		dbManager: dbManager,
	}
}

// Create は投稿を作成
func (r *PostRepository) Create(ctx context.Context, req *model.CreatePostRequest) (*model.Post, error) {
	now := time.Now()
	post := &model.Post{
		ID:        now.UnixNano(),
		UserID:    req.UserID,
		Title:     req.Title,
		Content:   req.Content,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// UserIDをキーとしてShardを決定（同じユーザーのデータは同じShardに配置）
	database, err := r.dbManager.GetDBByKey(req.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get database: %w", err)
	}

	query := `
		INSERT INTO posts (id, user_id, title, content, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err = database.ExecContext(ctx, query, post.ID, post.UserID, post.Title, post.Content, post.CreatedAt, post.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create post: %w", err)
	}

	return post, nil
}

// GetByID はIDで投稿を取得
func (r *PostRepository) GetByID(ctx context.Context, id int64, userID int64) (*model.Post, error) {
	database, err := r.dbManager.GetDBByKey(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get database: %w", err)
	}

	query := `
		SELECT id, user_id, title, content, created_at, updated_at
		FROM posts
		WHERE id = ?
	`

	var post model.Post
	err = database.QueryRowContext(ctx, query, id).Scan(
		&post.ID,
		&post.UserID,
		&post.Title,
		&post.Content,
		&post.CreatedAt,
		&post.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("post not found: %d", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get post: %w", err)
	}

	return &post, nil
}

// ListByUserID はユーザーIDで投稿一覧を取得
func (r *PostRepository) ListByUserID(ctx context.Context, userID int64, limit, offset int) ([]*model.Post, error) {
	database, err := r.dbManager.GetDBByKey(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get database: %w", err)
	}

	query := `
		SELECT id, user_id, title, content, created_at, updated_at
		FROM posts
		WHERE user_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := database.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query posts: %w", err)
	}
	defer rows.Close()

	posts := make([]*model.Post, 0)
	for rows.Next() {
		var post model.Post
		if err := rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.CreatedAt, &post.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan post: %w", err)
		}
		posts = append(posts, &post)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return posts, nil
}

// List はすべての投稿を取得（クロスシャードクエリ）
func (r *PostRepository) List(ctx context.Context, limit, offset int) ([]*model.Post, error) {
	connections := r.dbManager.GetAllConnections()
	posts := make([]*model.Post, 0)

	query := `
		SELECT id, user_id, title, content, created_at, updated_at
		FROM posts
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	for _, conn := range connections {
		rows, err := conn.DB.QueryContext(ctx, query, limit, offset)
		if err != nil {
			return nil, fmt.Errorf("failed to query shard %d: %w", conn.ShardID, err)
		}
		defer rows.Close()

		for rows.Next() {
			var post model.Post
			if err := rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.CreatedAt, &post.UpdatedAt); err != nil {
				return nil, fmt.Errorf("failed to scan post: %w", err)
			}
			posts = append(posts, &post)
		}

		if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("rows error: %w", err)
		}
	}

	return posts, nil
}

// GetUserPosts はユーザーと投稿をJOINして取得（クロスシャードクエリ）
func (r *PostRepository) GetUserPosts(ctx context.Context, limit, offset int) ([]*model.UserPost, error) {
	connections := r.dbManager.GetAllConnections()
	userPosts := make([]*model.UserPost, 0)

	query := `
		SELECT
			p.id as post_id,
			p.title as post_title,
			p.content as post_content,
			u.id as user_id,
			u.name as user_name,
			u.email as user_email,
			p.created_at
		FROM posts p
		INNER JOIN users u ON p.user_id = u.id
		ORDER BY p.created_at DESC
		LIMIT ? OFFSET ?
	`

	for _, conn := range connections {
		rows, err := conn.DB.QueryContext(ctx, query, limit, offset)
		if err != nil {
			return nil, fmt.Errorf("failed to query shard %d: %w", conn.ShardID, err)
		}
		defer rows.Close()

		for rows.Next() {
			var up model.UserPost
			if err := rows.Scan(&up.PostID, &up.PostTitle, &up.PostContent, &up.UserID, &up.UserName, &up.UserEmail, &up.CreatedAt); err != nil {
				return nil, fmt.Errorf("failed to scan user post: %w", err)
			}
			userPosts = append(userPosts, &up)
		}

		if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("rows error: %w", err)
		}
	}

	return userPosts, nil
}

// Update は投稿を更新
func (r *PostRepository) Update(ctx context.Context, id int64, userID int64, req *model.UpdatePostRequest) (*model.Post, error) {
	database, err := r.dbManager.GetDBByKey(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get database: %w", err)
	}

	query := "UPDATE posts SET updated_at = ?"
	args := []interface{}{time.Now()}

	if req.Title != "" {
		query += ", title = ?"
		args = append(args, req.Title)
	}
	if req.Content != "" {
		query += ", content = ?"
		args = append(args, req.Content)
	}

	query += " WHERE id = ? AND user_id = ?"
	args = append(args, id, userID)

	result, err := database.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to update post: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return nil, fmt.Errorf("post not found: %d", id)
	}

	return r.GetByID(ctx, id, userID)
}

// Delete は投稿を削除
func (r *PostRepository) Delete(ctx context.Context, id int64, userID int64) error {
	database, err := r.dbManager.GetDBByKey(userID)
	if err != nil {
		return fmt.Errorf("failed to get database: %w", err)
	}

	query := "DELETE FROM posts WHERE id = ? AND user_id = ?"
	result, err := database.ExecContext(ctx, query, id, userID)
	if err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("post not found: %d", id)
	}

	return nil
}
