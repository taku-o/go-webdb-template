package admin

import (
	"sync"

	"gorm.io/gorm"

	"github.com/taku-o/go-webdb-template/internal/db"
	"github.com/taku-o/go-webdb-template/internal/model"
)

// ShardResult はシャードクエリの結果を表す
type ShardResult[T any] struct {
	ShardID int
	Data    []T
	Error   error
}

// QueryAllShards は全シャードからデータを取得してマージする
func QueryAllShards[T any](manager *db.GORMManager, queryFn func(*gorm.DB) *gorm.DB) ([]T, error) {
	connections := manager.GetAllGORMConnections()
	results := make(chan ShardResult[T], len(connections))
	var wg sync.WaitGroup

	for _, conn := range connections {
		wg.Add(1)
		go func(conn *db.GORMConnection) {
			defer wg.Done()
			var data []T
			query := queryFn(conn.DB)
			err := query.Find(&data).Error
			results <- ShardResult[T]{
				ShardID: conn.ShardID,
				Data:    data,
				Error:   err,
			}
		}(conn)
	}

	wg.Wait()
	close(results)

	var allData []T
	for result := range results {
		if result.Error != nil {
			return nil, result.Error
		}
		allData = append(allData, result.Data...)
	}

	return allData, nil
}

// FindUserAcrossShards は全シャードからユーザーを検索する
func FindUserAcrossShards(manager *db.GORMManager, queryFn func(*gorm.DB) *gorm.DB) ([]model.DmUser, error) {
	return QueryAllShards[model.DmUser](manager, queryFn)
}

// FindPostAcrossShards は全シャードから投稿を検索する
func FindPostAcrossShards(manager *db.GORMManager, queryFn func(*gorm.DB) *gorm.DB) ([]model.DmPost, error) {
	return QueryAllShards[model.DmPost](manager, queryFn)
}

// CountAcrossShards は全シャードの件数を取得する
func CountAcrossShards[T any](manager *db.GORMManager) (int64, error) {
	connections := manager.GetAllGORMConnections()
	results := make(chan struct {
		Count int64
		Error error
	}, len(connections))
	var wg sync.WaitGroup

	for _, conn := range connections {
		wg.Add(1)
		go func(conn *db.GORMConnection) {
			defer wg.Done()
			var count int64
			var model T
			err := conn.DB.Model(&model).Count(&count).Error
			results <- struct {
				Count int64
				Error error
			}{Count: count, Error: err}
		}(conn)
	}

	wg.Wait()
	close(results)

	var totalCount int64
	for result := range results {
		if result.Error != nil {
			return 0, result.Error
		}
		totalCount += result.Count
	}

	return totalCount, nil
}

// CountUsersAcrossShards は全シャードのユーザー数を取得する
func CountUsersAcrossShards(manager *db.GORMManager) (int64, error) {
	return CountAcrossShards[model.DmUser](manager)
}

// CountPostsAcrossShards は全シャードの投稿数を取得する
func CountPostsAcrossShards(manager *db.GORMManager) (int64, error) {
	return CountAcrossShards[model.DmPost](manager)
}

// ShardStats はシャードの統計情報を表す
type ShardStats struct {
	ShardID    int
	UserCount  int64
	PostCount  int64
	TotalCount int64
}

// GetShardStats は各シャードの統計情報を取得する
func GetShardStats(manager *db.GORMManager) ([]ShardStats, error) {
	connections := manager.GetAllGORMConnections()
	results := make(chan ShardStats, len(connections))
	errChan := make(chan error, len(connections))
	var wg sync.WaitGroup

	for _, conn := range connections {
		wg.Add(1)
		go func(conn *db.GORMConnection) {
			defer wg.Done()

			var userCount, postCount int64
			if err := conn.DB.Model(&model.DmUser{}).Count(&userCount).Error; err != nil {
				errChan <- err
				return
			}
			if err := conn.DB.Model(&model.DmPost{}).Count(&postCount).Error; err != nil {
				errChan <- err
				return
			}

			results <- ShardStats{
				ShardID:    conn.ShardID,
				UserCount:  userCount,
				PostCount:  postCount,
				TotalCount: userCount + postCount,
			}
		}(conn)
	}

	wg.Wait()
	close(results)
	close(errChan)

	// エラーチェック
	for err := range errChan {
		if err != nil {
			return nil, err
		}
	}

	var stats []ShardStats
	for stat := range results {
		stats = append(stats, stat)
	}

	return stats, nil
}

// GetShardForUserID はユーザーIDに基づいてシャードDBを取得する
func GetShardForUserID(manager *db.GORMManager, userID int64) (*gorm.DB, error) {
	return manager.GetGORMByKey(userID)
}

// InsertToShard は指定されたシャードにデータを挿入する
func InsertToShard[T any](manager *db.GORMManager, shardKey int64, data *T) error {
	gormDB, err := manager.GetGORMByKey(shardKey)
	if err != nil {
		// デフォルトシャードを使用
		connections := manager.GetAllGORMConnections()
		if len(connections) == 0 {
			return err
		}
		gormDB = connections[0].DB
	}
	return gormDB.Create(data).Error
}
