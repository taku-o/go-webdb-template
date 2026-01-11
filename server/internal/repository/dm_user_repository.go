package repository

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/taku-o/go-webdb-template/internal/db"
	"github.com/taku-o/go-webdb-template/internal/model"
	"github.com/taku-o/go-webdb-template/internal/util/idgen"
	"gorm.io/gorm"
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
		tableSelector: db.NewTableSelector(db.DBShardingTableCount, db.DBShardingTablesPerDB),
	}
}

// Create はユーザーを作成
func (r *DmUserRepository) Create(ctx context.Context, req *model.CreateDmUserRequest) (*model.DmUser, error) {
	// ID生成（UUIDv7）
	id, err := idgen.GenerateUUIDv7()
	if err != nil {
		return nil, fmt.Errorf("failed to generate ID: %w", err)
	}

	user := &model.DmUser{
		ID:    id,
		Name:  req.Name,
		Email: req.Email,
	}

	// テーブル名の生成
	tableName, err := r.tableSelector.GetTableNameFromUUID("dm_users", user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get table name: %w", err)
	}

	// 接続の取得
	conn, err := r.groupManager.GetShardingConnectionByUUID(user.ID, "dm_users")
	if err != nil {
		return nil, fmt.Errorf("failed to get sharding connection: %w", err)
	}

	// リトライ機能付きでGORM APIで作成（動的テーブル名を使用）
	err = db.ExecuteWithRetry(func() error {
		return conn.DB.WithContext(ctx).Table(tableName).Create(user).Error
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// GetByID はIDでユーザーを取得
func (r *DmUserRepository) GetByID(ctx context.Context, id string) (*model.DmUser, error) {
	// テーブル名の生成
	tableName, err := r.tableSelector.GetTableNameFromUUID("dm_users", id)
	if err != nil {
		return nil, fmt.Errorf("failed to get table name: %w", err)
	}

	// 接続の取得
	conn, err := r.groupManager.GetShardingConnectionByUUID(id, "dm_users")
	if err != nil {
		return nil, fmt.Errorf("failed to get sharding connection: %w", err)
	}

	var user model.DmUser
	// リトライ機能付きでクエリ実行
	err = db.ExecuteWithRetry(func() error {
		return conn.DB.WithContext(ctx).Table(tableName).Where("id = ?", id).First(&user).Error
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found: %s", id)
		}
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

		tableName := fmt.Sprintf("dm_users_%03d", tableNum)

		var tableUsers []*model.DmUser
		// リトライ機能付きでクエリ実行
		err = db.ExecuteWithRetry(func() error {
			return conn.DB.WithContext(ctx).
				Table(tableName).
				Order("id").
				Limit(limit).
				Offset(offset).
				Find(&tableUsers).Error
		})
		if err != nil {
			return nil, fmt.Errorf("failed to query table %s: %w", tableName, err)
		}
		users = append(users, tableUsers...)
	}

	return users, nil
}

// Update はユーザーを更新
func (r *DmUserRepository) Update(ctx context.Context, id string, req *model.UpdateDmUserRequest) (*model.DmUser, error) {
	// テーブル名の生成
	tableName, err := r.tableSelector.GetTableNameFromUUID("dm_users", id)
	if err != nil {
		return nil, fmt.Errorf("failed to get table name: %w", err)
	}

	// 接続の取得
	conn, err := r.groupManager.GetShardingConnectionByUUID(id, "dm_users")
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

	var result *gorm.DB
	// リトライ機能付きでクエリ実行
	err = db.ExecuteWithRetry(func() error {
		result = conn.DB.WithContext(ctx).Table(tableName).Where("id = ?", id).Updates(updates)
		return result.Error
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}
	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("user not found: %s", id)
	}

	return r.GetByID(ctx, id)
}

// Delete はユーザーを削除
func (r *DmUserRepository) Delete(ctx context.Context, id string) error {
	// テーブル名の生成
	tableName, err := r.tableSelector.GetTableNameFromUUID("dm_users", id)
	if err != nil {
		return fmt.Errorf("failed to get table name: %w", err)
	}

	// 接続の取得
	conn, err := r.groupManager.GetShardingConnectionByUUID(id, "dm_users")
	if err != nil {
		return fmt.Errorf("failed to get sharding connection: %w", err)
	}

	var result *gorm.DB
	// リトライ機能付きでクエリ実行
	err = db.ExecuteWithRetry(func() error {
		result = conn.DB.WithContext(ctx).Table(tableName).Where("id = ?", id).Delete(&model.DmUser{})
		return result.Error
	})
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("user not found: %s", id)
	}

	return nil
}

// InsertDmUsersBatch はdm_usersテーブルにバッチでデータを挿入
func (r *DmUserRepository) InsertDmUsersBatch(ctx context.Context, tableName string, dmUsers []*model.DmUser) error {
	if len(dmUsers) == 0 {
		return nil
	}

	// テーブル番号から接続を取得
	// tableNameからテーブル番号を抽出（例: "dm_users_001" -> 1）
	tableNumber, err := extractTableNumber(tableName, "dm_users_")
	if err != nil {
		return fmt.Errorf("failed to extract table number from %s: %w", tableName, err)
	}

	conn, err := r.groupManager.GetShardingConnection(tableNumber)
	if err != nil {
		return fmt.Errorf("failed to get connection for table %d: %w", tableNumber, err)
	}

	// バッチサイズを考慮して分割
	for i := 0; i < len(dmUsers); i += db.BatchSize {
		end := i + db.BatchSize
		if end > len(dmUsers) {
			end = len(dmUsers)
		}
		batch := dmUsers[i:end]

		// GORMのCreateInBatchesを使用（動的テーブル名対応）
		err = db.ExecuteWithRetry(func() error {
			return conn.DB.WithContext(ctx).Table(tableName).CreateInBatches(batch, len(batch)).Error
		})
		if err != nil {
			return fmt.Errorf("failed to insert batch: %w", err)
		}
	}

	return nil
}

// extractTableNumber はテーブル名からテーブル番号を抽出
func extractTableNumber(tableName, prefix string) (int, error) {
	if !strings.HasPrefix(tableName, prefix) {
		return 0, fmt.Errorf("table name %s does not start with %s", tableName, prefix)
	}
	suffix := tableName[len(prefix):]
	tableNumber, err := strconv.Atoi(suffix)
	if err != nil {
		return 0, fmt.Errorf("failed to parse table number from %s: %w", suffix, err)
	}
	return tableNumber, nil
}
