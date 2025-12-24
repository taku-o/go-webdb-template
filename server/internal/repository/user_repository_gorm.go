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
	dbManager *db.GORMManager
}

// NewUserRepositoryGORM は新しいUserRepositoryGORMを作成
func NewUserRepositoryGORM(dbManager *db.GORMManager) *UserRepositoryGORM {
	return &UserRepositoryGORM{
		dbManager: dbManager,
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

	// シャードキーに基づいてGORMインスタンスを取得
	database, err := r.dbManager.GetGORMByKey(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get database: %w", err)
	}

	// GORM APIで作成
	if err := database.WithContext(ctx).Create(user).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// GetByID はIDでユーザーを取得
func (r *UserRepositoryGORM) GetByID(ctx context.Context, id int64) (*model.User, error) {
	database, err := r.dbManager.GetGORMByKey(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get database: %w", err)
	}

	var user model.User
	if err := database.WithContext(ctx).First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found: %d", id)
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// List はすべてのユーザーを取得（クロスシャードクエリ）
func (r *UserRepositoryGORM) List(ctx context.Context, limit, offset int) ([]*model.User, error) {
	connections := r.dbManager.GetAllGORMConnections()
	users := make([]*model.User, 0)

	// 各シャードから並列にデータを取得
	for _, conn := range connections {
		var shardUsers []*model.User
		if err := conn.DB.WithContext(ctx).
			Order("id").
			Limit(limit).
			Offset(offset).
			Find(&shardUsers).Error; err != nil {
			return nil, fmt.Errorf("failed to query shard %d: %w", conn.ShardID, err)
		}
		users = append(users, shardUsers...)
	}

	return users, nil
}

// Update はユーザーを更新
func (r *UserRepositoryGORM) Update(ctx context.Context, id int64, req *model.UpdateUserRequest) (*model.User, error) {
	database, err := r.dbManager.GetGORMByKey(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get database: %w", err)
	}

	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	updates["updated_at"] = time.Now()

	result := database.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).Updates(updates)
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
	database, err := r.dbManager.GetGORMByKey(id)
	if err != nil {
		return fmt.Errorf("failed to get database: %w", err)
	}

	result := database.WithContext(ctx).Delete(&model.User{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete user: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("user not found: %d", id)
	}

	return nil
}
