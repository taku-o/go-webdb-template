package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/example/go-db-prj-sample/internal/db"
	"github.com/example/go-db-prj-sample/internal/model"
)

// UserRepository はユーザーのデータアクセスを担当
type UserRepository struct {
	dbManager *db.Manager
}

// NewUserRepository は新しいUserRepositoryを作成
func NewUserRepository(dbManager *db.Manager) *UserRepository {
	return &UserRepository{
		dbManager: dbManager,
	}
}

// Create はユーザーを作成
func (r *UserRepository) Create(ctx context.Context, req *model.CreateUserRequest) (*model.User, error) {
	now := time.Now()
	user := &model.User{
		Name:      req.Name,
		Email:     req.Email,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// まだIDが決まっていないので、仮のIDを使ってShardを決定
	// 実際のアプリケーションではID生成戦略を工夫する必要がある
	// ここでは簡易的にタイムスタンプベースのIDを生成
	user.ID = now.UnixNano()

	database, err := r.dbManager.GetDBByKey(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get database: %w", err)
	}

	query := `
		INSERT INTO users (id, name, email, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`

	_, err = database.ExecContext(ctx, query, user.ID, user.Name, user.Email, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// GetByID はIDでユーザーを取得
func (r *UserRepository) GetByID(ctx context.Context, id int64) (*model.User, error) {
	database, err := r.dbManager.GetDBByKey(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get database: %w", err)
	}

	query := `
		SELECT id, name, email, created_at, updated_at
		FROM users
		WHERE id = ?
	`

	var user model.User
	err = database.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found: %d", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// List はすべてのユーザーを取得（クロスシャードクエリ）
func (r *UserRepository) List(ctx context.Context, limit, offset int) ([]*model.User, error) {
	connections := r.dbManager.GetAllConnections()
	users := make([]*model.User, 0)

	query := `
		SELECT id, name, email, created_at, updated_at
		FROM users
		ORDER BY id
		LIMIT ? OFFSET ?
	`

	// 各Shardから並列にデータを取得してマージ
	for _, conn := range connections {
		rows, err := conn.DB.QueryContext(ctx, query, limit, offset)
		if err != nil {
			return nil, fmt.Errorf("failed to query shard %d: %w", conn.ShardID, err)
		}
		defer rows.Close()

		for rows.Next() {
			var user model.User
			if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt); err != nil {
				return nil, fmt.Errorf("failed to scan user: %w", err)
			}
			users = append(users, &user)
		}

		if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("rows error: %w", err)
		}
	}

	return users, nil
}

// Update はユーザーを更新
func (r *UserRepository) Update(ctx context.Context, id int64, req *model.UpdateUserRequest) (*model.User, error) {
	database, err := r.dbManager.GetDBByKey(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get database: %w", err)
	}

	// 更新するフィールドを動的に構築
	query := "UPDATE users SET updated_at = ?"
	args := []interface{}{time.Now()}

	if req.Name != "" {
		query += ", name = ?"
		args = append(args, req.Name)
	}
	if req.Email != "" {
		query += ", email = ?"
		args = append(args, req.Email)
	}

	query += " WHERE id = ?"
	args = append(args, id)

	result, err := database.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return nil, fmt.Errorf("user not found: %d", id)
	}

	// 更新後のユーザーを取得
	return r.GetByID(ctx, id)
}

// Delete はユーザーを削除
func (r *UserRepository) Delete(ctx context.Context, id int64) error {
	database, err := r.dbManager.GetDBByKey(id)
	if err != nil {
		return fmt.Errorf("failed to get database: %w", err)
	}

	query := "DELETE FROM users WHERE id = ?"
	result, err := database.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("user not found: %d", id)
	}

	return nil
}
