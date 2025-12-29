package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/taku-o/go-webdb-template/internal/config"
	"github.com/taku-o/go-webdb-template/internal/db"
	"github.com/taku-o/go-webdb-template/internal/model"
	"gorm.io/gorm"
)

const (
	batchSize  = 500
	tableCount = 32
	totalCount = 100
)

func main() {
	// 1. 設定ファイルの読み込み
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2. GroupManagerの初期化
	groupManager, err := db.NewGroupManager(cfg)
	if err != nil {
		log.Fatalf("Failed to create group manager: %v", err)
	}
	defer groupManager.CloseAll()

	// 3. データベース接続確認
	if err := groupManager.PingAll(); err != nil {
		log.Fatalf("Failed to ping databases: %v", err)
	}

	log.Println("Starting sample data generation...")

	// 4. dm_usersテーブルへのデータ生成
	userIDs, err := generateUsers(groupManager, totalCount)
	if err != nil {
		log.Fatalf("Failed to generate users: %v", err)
	}

	// 5. dm_postsテーブルへのデータ生成
	if err := generatePosts(groupManager, userIDs, totalCount); err != nil {
		log.Fatalf("Failed to generate posts: %v", err)
	}

	// 6. dm_newsテーブルへのデータ生成
	if err := generateNews(groupManager, totalCount); err != nil {
		log.Fatalf("Failed to generate news: %v", err)
	}

	// 7. 生成完了メッセージ
	log.Println("Sample data generation completed successfully")
	os.Exit(0)
}

// generateUsers はdm_usersテーブルにデータを生成
// 戻り値: 生成されたuser_idのリスト（posts生成時に使用）
func generateUsers(groupManager *db.GroupManager, totalCount int) ([]int64, error) {
	countPerTable := totalCount / tableCount // 各テーブルに約3～4件

	var allUserIDs []int64

	// 各テーブル（0～31）に対してデータ生成
	for tableNumber := 0; tableNumber < tableCount; tableNumber++ {
		// 接続を取得
		conn, err := groupManager.GetShardingConnection(tableNumber)
		if err != nil {
			return nil, fmt.Errorf("failed to get connection for table %d: %w", tableNumber, err)
		}

		// テーブル名を生成
		tableName := fmt.Sprintf("dm_users_%03d", tableNumber)

		// バッチでデータ生成
		var users []*model.DmUser
		for i := 0; i < countPerTable; i++ {
			user := &model.DmUser{
				Name:      gofakeit.Name(),
				Email:     gofakeit.Email(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			users = append(users, user)
		}

		// バッチ挿入
		if len(users) > 0 {
			if err := insertUsersBatch(conn.DB, tableName, users); err != nil {
				return nil, fmt.Errorf("failed to insert batch to %s: %w", tableName, err)
			}

			// 挿入後にuser_idを取得
			ids, err := fetchUserIDs(conn.DB, tableName, len(users))
			if err != nil {
				return nil, fmt.Errorf("failed to fetch user IDs from %s: %w", tableName, err)
			}
			allUserIDs = append(allUserIDs, ids...)
		}

		log.Printf("Generated %d users in %s", countPerTable, tableName)
	}

	return allUserIDs, nil
}

// generatePosts はdm_postsテーブルにデータを生成
// userIDs: 既存のdm_usersテーブルから取得したuser_idのリスト
func generatePosts(groupManager *db.GroupManager, userIDs []int64, totalCount int) error {
	countPerTable := totalCount / tableCount // 各テーブルに約3～4件

	if len(userIDs) == 0 {
		return fmt.Errorf("no user IDs available for posts generation")
	}

	// 各テーブル（0～31）に対してデータ生成
	for tableNumber := 0; tableNumber < tableCount; tableNumber++ {
		// 接続を取得
		conn, err := groupManager.GetShardingConnection(tableNumber)
		if err != nil {
			return fmt.Errorf("failed to get connection for table %d: %w", tableNumber, err)
		}

		// テーブル名を生成
		tableName := fmt.Sprintf("dm_posts_%03d", tableNumber)

		// バッチでデータ生成
		var posts []*model.DmPost
		for i := 0; i < countPerTable; i++ {
			// user_idをランダムに選択
			userID := userIDs[gofakeit.IntRange(0, len(userIDs)-1)]

			post := &model.DmPost{
				UserID:    userID,
				Title:     gofakeit.Sentence(5),
				Content:   gofakeit.Paragraph(3, 5, 10, "\n"),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			posts = append(posts, post)
		}

		// バッチ挿入
		if len(posts) > 0 {
			if err := insertPostsBatch(conn.DB, tableName, posts); err != nil {
				return fmt.Errorf("failed to insert batch to %s: %w", tableName, err)
			}
		}

		log.Printf("Generated %d posts in %s", countPerTable, tableName)
	}

	return nil
}

// generateNews はdm_newsテーブルにデータを生成
func generateNews(groupManager *db.GroupManager, totalCount int) error {
	// master接続を取得
	conn, err := groupManager.GetMasterConnection()
	if err != nil {
		return fmt.Errorf("failed to get master connection: %w", err)
	}

	// バッチでデータ生成
	var news []*model.DmNews
	for i := 0; i < totalCount; i++ {
		authorID := gofakeit.Int64()
		publishedAt := gofakeit.Date()

		n := &model.DmNews{
			Title:       gofakeit.Sentence(5),
			Content:     gofakeit.Paragraph(3, 5, 10, "\n"),
			AuthorID:    &authorID,
			PublishedAt: &publishedAt,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		news = append(news, n)
	}

	// バッチ挿入
	if len(news) > 0 {
		if err := insertNewsBatch(conn.DB, news); err != nil {
			return fmt.Errorf("failed to insert batch to news: %w", err)
		}
	}

	log.Printf("Generated %d news articles", totalCount)
	return nil
}

// insertUsersBatch はdm_usersテーブルにバッチでデータを挿入
func insertUsersBatch(db *gorm.DB, tableName string, users []*model.DmUser) error {
	if len(users) == 0 {
		return nil
	}

	// バッチサイズを考慮して分割
	for i := 0; i < len(users); i += batchSize {
		end := i + batchSize
		if end > len(users) {
			end = len(users)
		}
		batch := users[i:end]

		// 生SQLでバッチ挿入（動的テーブル名対応）
		query := fmt.Sprintf("INSERT INTO %s (name, email, created_at, updated_at) VALUES ", tableName)
		var values []interface{}
		var placeholders []string

		for _, user := range batch {
			placeholders = append(placeholders, "(?, ?, ?, ?)")
			values = append(values, user.Name, user.Email, user.CreatedAt, user.UpdatedAt)
		}

		query += strings.Join(placeholders, ", ")

		if err := db.Exec(query, values...).Error; err != nil {
			return fmt.Errorf("failed to insert batch: %w", err)
		}
	}

	return nil
}

// insertPostsBatch はdm_postsテーブルにバッチでデータを挿入
func insertPostsBatch(db *gorm.DB, tableName string, posts []*model.DmPost) error {
	if len(posts) == 0 {
		return nil
	}

	// バッチサイズを考慮して分割
	for i := 0; i < len(posts); i += batchSize {
		end := i + batchSize
		if end > len(posts) {
			end = len(posts)
		}
		batch := posts[i:end]

		// 生SQLでバッチ挿入（動的テーブル名対応）
		query := fmt.Sprintf("INSERT INTO %s (user_id, title, content, created_at, updated_at) VALUES ", tableName)
		var values []interface{}
		var placeholders []string

		for _, post := range batch {
			placeholders = append(placeholders, "(?, ?, ?, ?, ?)")
			values = append(values, post.UserID, post.Title, post.Content, post.CreatedAt, post.UpdatedAt)
		}

		query += strings.Join(placeholders, ", ")

		if err := db.Exec(query, values...).Error; err != nil {
			return fmt.Errorf("failed to insert batch: %w", err)
		}
	}

	return nil
}

// insertNewsBatch はdm_newsテーブルにバッチでデータを挿入
func insertNewsBatch(db *gorm.DB, news []*model.DmNews) error {
	if len(news) == 0 {
		return nil
	}

	// バッチサイズを考慮して分割
	for i := 0; i < len(news); i += batchSize {
		end := i + batchSize
		if end > len(news) {
			end = len(news)
		}
		batch := news[i:end]

		// GORMのCreateInBatchesを使用（固定テーブル名）
		if err := db.Table("dm_news").CreateInBatches(batch, len(batch)).Error; err != nil {
			return fmt.Errorf("failed to insert batch: %w", err)
		}
	}

	return nil
}

// fetchUserIDs は挿入後のuser_idを取得
func fetchUserIDs(db *gorm.DB, tableName string, limit int) ([]int64, error) {
	var ids []int64
	query := fmt.Sprintf("SELECT id FROM %s ORDER BY id DESC LIMIT ?", tableName)
	if err := db.Raw(query, limit).Scan(&ids).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch user IDs: %w", err)
	}
	return ids, nil
}
