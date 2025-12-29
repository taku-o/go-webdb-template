package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/taku-o/go-webdb-template/internal/db"
	"github.com/taku-o/go-webdb-template/internal/model"
)

// DmUserRepository はユーザーのデータアクセスを担当
type DmUserRepository struct {
	groupManager  *db.GroupManager
	tableSelector *db.TableSelector
}

// NewDmUserRepository は新しいDmUserRepositoryを作成
func NewDmUserRepository(groupManager *db.GroupManager) *DmUserRepository {
	return &DmUserRepository{
		groupManager:  groupManager,
		tableSelector: db.NewTableSelector(32, 8),
	}
}

// Create はユーザーを作成
func (r *DmUserRepository) Create(ctx context.Context, req *model.CreateDmUserRequest) (*model.DmUser, error) {
	now := time.Now()
	user := &model.DmUser{
		Name:      req.Name,
		Email:     req.Email,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// IDを生成（タイムスタンプベース）
	user.ID = now.UnixNano()

	// テーブル名の生成
	tableName := r.tableSelector.GetTableName("dm_users", user.ID)

	// 接続の取得
	conn, err := r.groupManager.GetShardingConnectionByID(user.ID, "dm_users")
	if err != nil {
		return nil, fmt.Errorf("failed to get sharding connection: %w", err)
	}

	// sql.DBを取得
	sqlDB, err := conn.DB.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	query := fmt.Sprintf(`
		INSERT INTO %s (id, name, email, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`, tableName)

	_, err = sqlDB.ExecContext(ctx, query, user.ID, user.Name, user.Email, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// GetByID はIDでユーザーを取得
func (r *DmUserRepository) GetByID(ctx context.Context, id int64) (*model.DmUser, error) {
	// テーブル名の生成
	tableName := r.tableSelector.GetTableName("dm_users", id)

	// 接続の取得
	conn, err := r.groupManager.GetShardingConnectionByID(id, "dm_users")
	if err != nil {
		return nil, fmt.Errorf("failed to get sharding connection: %w", err)
	}

	// sql.DBを取得
	sqlDB, err := conn.DB.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	query := fmt.Sprintf(`
		SELECT id, name, email, created_at, updated_at
		FROM %s
		WHERE id = ?
	`, tableName)

	var user model.DmUser
	err = sqlDB.QueryRowContext(ctx, query, id).Scan(
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

// List はすべてのユーザーを取得（クロステーブルクエリ）
func (r *DmUserRepository) List(ctx context.Context, limit, offset int) ([]*model.DmUser, error) {
	users := make([]*model.DmUser, 0)

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

		tableName := fmt.Sprintf("dm_users_%03d", tableNum)

		query := fmt.Sprintf(`
			SELECT id, name, email, created_at, updated_at
			FROM %s
			ORDER BY id
			LIMIT ? OFFSET ?
		`, tableName)

		rows, err := sqlDB.QueryContext(ctx, query, limit, offset)
		if err != nil {
			return nil, fmt.Errorf("failed to query table %s: %w", tableName, err)
		}

		for rows.Next() {
			var user model.DmUser
			if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt); err != nil {
				rows.Close()
				return nil, fmt.Errorf("failed to scan user: %w", err)
			}
			users = append(users, &user)
		}

		if err := rows.Err(); err != nil {
			rows.Close()
			return nil, fmt.Errorf("rows error: %w", err)
		}
		rows.Close()
	}

	return users, nil
}

// Update はユーザーを更新
func (r *DmUserRepository) Update(ctx context.Context, id int64, req *model.UpdateDmUserRequest) (*model.DmUser, error) {
	// テーブル名の生成
	tableName := r.tableSelector.GetTableName("dm_users", id)

	// 接続の取得
	conn, err := r.groupManager.GetShardingConnectionByID(id, "dm_users")
	if err != nil {
		return nil, fmt.Errorf("failed to get sharding connection: %w", err)
	}

	// sql.DBを取得
	sqlDB, err := conn.DB.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	// 更新するフィールドを動的に構築
	query := fmt.Sprintf("UPDATE %s SET updated_at = ?", tableName)
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

	result, err := sqlDB.ExecContext(ctx, query, args...)
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
func (r *DmUserRepository) Delete(ctx context.Context, id int64) error {
	// テーブル名の生成
	tableName := r.tableSelector.GetTableName("dm_users", id)

	// 接続の取得
	conn, err := r.groupManager.GetShardingConnectionByID(id, "dm_users")
	if err != nil {
		return fmt.Errorf("failed to get sharding connection: %w", err)
	}

	// sql.DBを取得
	sqlDB, err := conn.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	query := fmt.Sprintf("DELETE FROM %s WHERE id = ?", tableName)
	result, err := sqlDB.ExecContext(ctx, query, id)
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
