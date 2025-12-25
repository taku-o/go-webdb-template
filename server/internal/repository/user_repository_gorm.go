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

// UserRepositoryGORM はユーザーのデータアクセスを担当（GORM版）
type UserRepositoryGORM struct {
	groupManager  *db.GroupManager
	tableSelector *db.TableSelector
}

// NewUserRepositoryGORM は新しいUserRepositoryGORMを作成
func NewUserRepositoryGORM(groupManager *db.GroupManager) *UserRepositoryGORM {
	return &UserRepositoryGORM{
		groupManager:  groupManager,
		tableSelector: db.NewTableSelector(32, 8),
	}
}

// Create はユーザーを作成
func (r *UserRepositoryGORM) Create(ctx context.Context, req *model.CreateUserRequest) (*model.User, error) {
	user := &model.User{
		Name:  req.Name,
		Email: req.Email,
	}

	// ID生成（タイムスタンプベース、既存ロジック維持）
	user.ID = time.Now().UnixNano()

	// テーブル名の生成
	tableName := r.tableSelector.GetTableName("users", user.ID)

	// 接続の取得
	conn, err := r.groupManager.GetShardingConnectionByID(user.ID, "users")
	if err != nil {
		return nil, fmt.Errorf("failed to get sharding connection: %w", err)
	}

	// GORM APIで作成（動的テーブル名を使用）
	if err := conn.DB.WithContext(ctx).Table(tableName).Create(user).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// GetByID はIDでユーザーを取得
func (r *UserRepositoryGORM) GetByID(ctx context.Context, id int64) (*model.User, error) {
	// テーブル名の生成
	tableName := r.tableSelector.GetTableName("users", id)

	// 接続の取得
	conn, err := r.groupManager.GetShardingConnectionByID(id, "users")
	if err != nil {
		return nil, fmt.Errorf("failed to get sharding connection: %w", err)
	}

	var user model.User
	if err := conn.DB.WithContext(ctx).Table(tableName).First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found: %d", id)
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// List はすべてのユーザーを取得（クロステーブルクエリ）
func (r *UserRepositoryGORM) List(ctx context.Context, limit, offset int) ([]*model.User, error) {
	connections := r.groupManager.GetAllShardingConnections()
	users := make([]*model.User, 0)

	// 各データベースの各テーブルからデータを取得
	for _, conn := range connections {
		// このデータベースに含まれるテーブル（8つずつ）
		startTable := (conn.ShardID - 1) * 8
		endTable := startTable + 7

		for tableNum := startTable; tableNum <= endTable; tableNum++ {
			tableName := fmt.Sprintf("users_%03d", tableNum)

			var tableUsers []*model.User
			if err := conn.DB.WithContext(ctx).
				Table(tableName).
				Order("id").
				Limit(limit).
				Offset(offset).
				Find(&tableUsers).Error; err != nil {
				return nil, fmt.Errorf("failed to query table %s: %w", tableName, err)
			}
			users = append(users, tableUsers...)
		}
	}

	return users, nil
}

// Update はユーザーを更新
func (r *UserRepositoryGORM) Update(ctx context.Context, id int64, req *model.UpdateUserRequest) (*model.User, error) {
	// テーブル名の生成
	tableName := r.tableSelector.GetTableName("users", id)

	// 接続の取得
	conn, err := r.groupManager.GetShardingConnectionByID(id, "users")
	if err != nil {
		return nil, fmt.Errorf("failed to get sharding connection: %w", err)
	}

	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	updates["updated_at"] = time.Now()

	result := conn.DB.WithContext(ctx).Table(tableName).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to update user: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("user not found: %d", id)
	}

	return r.GetByID(ctx, id)
}

// Delete はユーザーを削除
func (r *UserRepositoryGORM) Delete(ctx context.Context, id int64) error {
	// テーブル名の生成
	tableName := r.tableSelector.GetTableName("users", id)

	// 接続の取得
	conn, err := r.groupManager.GetShardingConnectionByID(id, "users")
	if err != nil {
		return fmt.Errorf("failed to get sharding connection: %w", err)
	}

	result := conn.DB.WithContext(ctx).Table(tableName).Delete(&model.User{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete user: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("user not found: %d", id)
	}

	return nil
}
