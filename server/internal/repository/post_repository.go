package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/taku-o/go-webdb-template/internal/db"
	"github.com/taku-o/go-webdb-template/internal/model"
)

// PostRepository は投稿のデータアクセスを担当
type PostRepository struct {
	groupManager  *db.GroupManager
	tableSelector *db.TableSelector
}

// NewPostRepository は新しいPostRepositoryを作成
func NewPostRepository(groupManager *db.GroupManager) *PostRepository {
	return &PostRepository{
		groupManager:  groupManager,
		tableSelector: db.NewTableSelector(32, 8),
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

	// UserIDをキーとしてテーブル/DBを決定（同じユーザーのデータは同じテーブルに配置）
	tableName := r.tableSelector.GetTableName("posts", req.UserID)

	// 接続の取得
	conn, err := r.groupManager.GetShardingConnectionByID(req.UserID, "posts")
	if err != nil {
		return nil, fmt.Errorf("failed to get sharding connection: %w", err)
	}

	// sql.DBを取得
	sqlDB, err := conn.DB.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	query := fmt.Sprintf(`
		INSERT INTO %s (id, user_id, title, content, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`, tableName)

	_, err = sqlDB.ExecContext(ctx, query, post.ID, post.UserID, post.Title, post.Content, post.CreatedAt, post.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create post: %w", err)
	}

	return post, nil
}

// GetByID はIDで投稿を取得
func (r *PostRepository) GetByID(ctx context.Context, id int64, userID int64) (*model.Post, error) {
	// UserIDをキーとしてテーブル/DBを決定
	tableName := r.tableSelector.GetTableName("posts", userID)

	// 接続の取得
	conn, err := r.groupManager.GetShardingConnectionByID(userID, "posts")
	if err != nil {
		return nil, fmt.Errorf("failed to get sharding connection: %w", err)
	}

	// sql.DBを取得
	sqlDB, err := conn.DB.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	query := fmt.Sprintf(`
		SELECT id, user_id, title, content, created_at, updated_at
		FROM %s
		WHERE id = ?
	`, tableName)

	var post model.Post
	err = sqlDB.QueryRowContext(ctx, query, id).Scan(
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
	// UserIDをキーとしてテーブル/DBを決定
	tableName := r.tableSelector.GetTableName("posts", userID)

	// 接続の取得
	conn, err := r.groupManager.GetShardingConnectionByID(userID, "posts")
	if err != nil {
		return nil, fmt.Errorf("failed to get sharding connection: %w", err)
	}

	// sql.DBを取得
	sqlDB, err := conn.DB.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	query := fmt.Sprintf(`
		SELECT id, user_id, title, content, created_at, updated_at
		FROM %s
		WHERE user_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`, tableName)

	rows, err := sqlDB.QueryContext(ctx, query, userID, limit, offset)
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

// List はすべての投稿を取得（クロステーブルクエリ）
func (r *PostRepository) List(ctx context.Context, limit, offset int) ([]*model.Post, error) {
	posts := make([]*model.Post, 0)

	// テーブル数分ループして各テーブルからデータを取得
	tableCount := r.tableSelector.GetTableCount()
	for tableNum := 0; tableNum < tableCount; tableNum++ {
		// テーブル番号から接続を取得
		conn, err := r.groupManager.GetShardingConnection(tableNum)
		if err != nil {
			return nil, fmt.Errorf("failed to get connection for table %d: %w", tableNum, err)
		}

		sqlDB, err := conn.DB.DB()
		if err != nil {
			return nil, fmt.Errorf("failed to get sql.DB for table %d: %w", tableNum, err)
		}

		tableName := fmt.Sprintf("posts_%03d", tableNum)

		query := fmt.Sprintf(`
			SELECT id, user_id, title, content, created_at, updated_at
			FROM %s
			ORDER BY created_at DESC
			LIMIT ? OFFSET ?
		`, tableName)

		rows, err := sqlDB.QueryContext(ctx, query, limit, offset)
		if err != nil {
			return nil, fmt.Errorf("failed to query table %s: %w", tableName, err)
		}

		for rows.Next() {
			var post model.Post
			if err := rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.CreatedAt, &post.UpdatedAt); err != nil {
				rows.Close()
				return nil, fmt.Errorf("failed to scan post: %w", err)
			}
			posts = append(posts, &post)
		}

		if err := rows.Err(); err != nil {
			rows.Close()
			return nil, fmt.Errorf("rows error: %w", err)
		}
		rows.Close()
	}

	return posts, nil
}

// GetUserPosts はユーザーと投稿をJOINして取得（クロステーブルクエリ）
func (r *PostRepository) GetUserPosts(ctx context.Context, limit, offset int) ([]*model.UserPost, error) {
	userPosts := make([]*model.UserPost, 0)

	// テーブル数分ループして各テーブルからデータを取得
	tableCount := r.tableSelector.GetTableCount()
	for tableNum := 0; tableNum < tableCount; tableNum++ {
		// テーブル番号から接続を取得
		conn, err := r.groupManager.GetShardingConnection(tableNum)
		if err != nil {
			return nil, fmt.Errorf("failed to get connection for table %d: %w", tableNum, err)
		}

		sqlDB, err := conn.DB.DB()
		if err != nil {
			return nil, fmt.Errorf("failed to get sql.DB for table %d: %w", tableNum, err)
		}

		postsTable := fmt.Sprintf("posts_%03d", tableNum)
		usersTable := fmt.Sprintf("users_%03d", tableNum)

		query := fmt.Sprintf(`
			SELECT
				p.id as post_id,
				p.title as post_title,
				p.content as post_content,
				u.id as user_id,
				u.name as user_name,
				u.email as user_email,
				p.created_at
			FROM %s p
			INNER JOIN %s u ON p.user_id = u.id
			ORDER BY p.created_at DESC
			LIMIT ? OFFSET ?
		`, postsTable, usersTable)

		rows, err := sqlDB.QueryContext(ctx, query, limit, offset)
		if err != nil {
			return nil, fmt.Errorf("failed to query table %s: %w", postsTable, err)
		}

		for rows.Next() {
			var up model.UserPost
			if err := rows.Scan(&up.PostID, &up.PostTitle, &up.PostContent, &up.UserID, &up.UserName, &up.UserEmail, &up.CreatedAt); err != nil {
				rows.Close()
				return nil, fmt.Errorf("failed to scan user post: %w", err)
			}
			userPosts = append(userPosts, &up)
		}

		if err := rows.Err(); err != nil {
			rows.Close()
			return nil, fmt.Errorf("rows error: %w", err)
		}
		rows.Close()
	}

	return userPosts, nil
}

// Update は投稿を更新
func (r *PostRepository) Update(ctx context.Context, id int64, userID int64, req *model.UpdatePostRequest) (*model.Post, error) {
	// UserIDをキーとしてテーブル/DBを決定
	tableName := r.tableSelector.GetTableName("posts", userID)

	// 接続の取得
	conn, err := r.groupManager.GetShardingConnectionByID(userID, "posts")
	if err != nil {
		return nil, fmt.Errorf("failed to get sharding connection: %w", err)
	}

	// sql.DBを取得
	sqlDB, err := conn.DB.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	query := fmt.Sprintf("UPDATE %s SET updated_at = ?", tableName)
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

	result, err := sqlDB.ExecContext(ctx, query, args...)
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

	// 更新後の投稿を取得
	return r.GetByID(ctx, id, userID)
}

// Delete は投稿を削除
func (r *PostRepository) Delete(ctx context.Context, id int64, userID int64) error {
	// UserIDをキーとしてテーブル/DBを決定
	tableName := r.tableSelector.GetTableName("posts", userID)

	// 接続の取得
	conn, err := r.groupManager.GetShardingConnectionByID(userID, "posts")
	if err != nil {
		return fmt.Errorf("failed to get sharding connection: %w", err)
	}

	// sql.DBを取得
	sqlDB, err := conn.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	query := fmt.Sprintf("DELETE FROM %s WHERE id = ? AND user_id = ?", tableName)
	result, err := sqlDB.ExecContext(ctx, query, id, userID)
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
