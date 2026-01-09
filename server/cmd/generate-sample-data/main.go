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
	"github.com/taku-o/go-webdb-template/internal/util/idgen"
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
	dmUserIDs, err := generateDmUsers(groupManager, totalCount)
	if err != nil {
		log.Fatalf("Failed to generate users: %v", err)
	}

	// 5. dm_postsテーブルへのデータ生成
	if err := generateDmPosts(groupManager, dmUserIDs, totalCount); err != nil {
		log.Fatalf("Failed to generate posts: %v", err)
	}

	// 6. dm_newsテーブルへのデータ生成
	if err := generateDmNews(groupManager, totalCount); err != nil {
		log.Fatalf("Failed to generate news: %v", err)
	}

	// 7. 生成完了メッセージ
	log.Println("Sample data generation completed successfully")
	os.Exit(0)
}

// generateDmUsers はdm_usersテーブルにデータを生成
// 戻り値: 生成されたdm_user_idのリスト（dm_posts生成時に使用）
func generateDmUsers(groupManager *db.GroupManager, totalCount int) ([]string, error) {
	var allDmUserIDs []string

	// テーブル番号ごとにユーザーをグループ化するマップ
	usersByTable := make(map[int][]*model.DmUser)

	// TableSelectorを作成
	tableSelector := db.NewTableSelector(tableCount, db.DBShardingTablesPerDB)

	// 全ユーザーを生成し、IDに基づいて正しいテーブルに振り分け
	for i := 0; i < totalCount; i++ {
		id, err := idgen.GenerateUUIDv7()
		if err != nil {
			return nil, fmt.Errorf("failed to generate UUIDv7: %w", err)
		}

		// UUIDからテーブル番号を計算
		tableNumber, err := tableSelector.GetTableNumberFromUUID(id)
		if err != nil {
			return nil, fmt.Errorf("failed to get table number from UUID: %w", err)
		}

		dmUser := &model.DmUser{
			ID:        id,
			Name:      gofakeit.Name(),
			Email:     gofakeit.Email(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		usersByTable[tableNumber] = append(usersByTable[tableNumber], dmUser)
		allDmUserIDs = append(allDmUserIDs, id)
	}

	// 各テーブルにデータを挿入
	for tableNumber, dmUsers := range usersByTable {
		// 接続を取得
		conn, err := groupManager.GetShardingConnection(tableNumber)
		if err != nil {
			return nil, fmt.Errorf("failed to get connection for table %d: %w", tableNumber, err)
		}

		// テーブル名を生成
		tableName := fmt.Sprintf("dm_users_%03d", tableNumber)

		// バッチ挿入
		if len(dmUsers) > 0 {
			if err := insertDmUsersBatch(conn.DB, tableName, dmUsers); err != nil {
				return nil, fmt.Errorf("failed to insert batch to %s: %w", tableName, err)
			}
		}

		log.Printf("Generated %d dm_users in %s", len(dmUsers), tableName)
	}

	return allDmUserIDs, nil
}

// generateDmPosts はdm_postsテーブルにデータを生成
// dmUserIDs: 既存のdm_usersテーブルから取得したdm_user_idのリスト
func generateDmPosts(groupManager *db.GroupManager, dmUserIDs []string, totalCount int) error {
	if len(dmUserIDs) == 0 {
		return fmt.Errorf("no dm_user IDs available for dm_posts generation")
	}

	// テーブル番号ごとに投稿をグループ化するマップ
	postsByTable := make(map[int][]*model.DmPost)

	// TableSelectorを作成
	tableSelector := db.NewTableSelector(tableCount, db.DBShardingTablesPerDB)

	// 全投稿を生成し、user_idに基づいて正しいテーブルに振り分け
	for i := 0; i < totalCount; i++ {
		id, err := idgen.GenerateUUIDv7()
		if err != nil {
			return fmt.Errorf("failed to generate UUIDv7: %w", err)
		}

		// dm_user_idをランダムに選択
		dmUserID := dmUserIDs[gofakeit.IntRange(0, len(dmUserIDs)-1)]

		// user_idからテーブル番号を計算（dm_postsのシャーディングキーはuser_id）
		tableNumber, err := tableSelector.GetTableNumberFromUUID(dmUserID)
		if err != nil {
			return fmt.Errorf("failed to get table number from UUID: %w", err)
		}

		dmPost := &model.DmPost{
			ID:        id,
			UserID:    dmUserID,
			Title:     gofakeit.Sentence(5),
			Content:   gofakeit.Paragraph(3, 5, 10, "\n"),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		postsByTable[tableNumber] = append(postsByTable[tableNumber], dmPost)
	}

	// 各テーブルにデータを挿入
	for tableNumber, dmPosts := range postsByTable {
		// 接続を取得
		conn, err := groupManager.GetShardingConnection(tableNumber)
		if err != nil {
			return fmt.Errorf("failed to get connection for table %d: %w", tableNumber, err)
		}

		// テーブル名を生成
		tableName := fmt.Sprintf("dm_posts_%03d", tableNumber)

		// バッチ挿入
		if len(dmPosts) > 0 {
			if err := insertDmPostsBatch(conn.DB, tableName, dmPosts); err != nil {
				return fmt.Errorf("failed to insert batch to %s: %w", tableName, err)
			}
		}

		log.Printf("Generated %d dm_posts in %s", len(dmPosts), tableName)
	}

	return nil
}

// generateDmNews はdm_newsテーブルにデータを生成
func generateDmNews(groupManager *db.GroupManager, totalCount int) error {
	// master接続を取得
	conn, err := groupManager.GetMasterConnection()
	if err != nil {
		return fmt.Errorf("failed to get master connection: %w", err)
	}

	// バッチでデータ生成
	var dmNews []*model.DmNews
	for i := 0; i < totalCount; i++ {
		authorID := int64(gofakeit.Int32()) & 0x7FFFFFFF
		publishedAt := gofakeit.Date()

		n := &model.DmNews{
			Title:       gofakeit.Sentence(5),
			Content:     gofakeit.Paragraph(3, 5, 10, "\n"),
			AuthorID:    &authorID,
			PublishedAt: &publishedAt,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		dmNews = append(dmNews, n)
	}

	// バッチ挿入
	if len(dmNews) > 0 {
		if err := insertDmNewsBatch(conn.DB, dmNews); err != nil {
			return fmt.Errorf("failed to insert batch to dm_news: %w", err)
		}
	}

	log.Printf("Generated %d dm_news articles", totalCount)
	return nil
}

// insertDmUsersBatch はdm_usersテーブルにバッチでデータを挿入
func insertDmUsersBatch(db *gorm.DB, tableName string, dmUsers []*model.DmUser) error {
	if len(dmUsers) == 0 {
		return nil
	}

	// バッチサイズを考慮して分割
	for i := 0; i < len(dmUsers); i += batchSize {
		end := i + batchSize
		if end > len(dmUsers) {
			end = len(dmUsers)
		}
		batch := dmUsers[i:end]

		// 生SQLでバッチ挿入（動的テーブル名対応、IDを含む）
		query := fmt.Sprintf("INSERT INTO %s (id, name, email, created_at, updated_at) VALUES ", tableName)
		var values []interface{}
		var placeholders []string

		for _, dmUser := range batch {
			placeholders = append(placeholders, "(?, ?, ?, ?, ?)")
			values = append(values, dmUser.ID, dmUser.Name, dmUser.Email, dmUser.CreatedAt, dmUser.UpdatedAt)
		}

		query += strings.Join(placeholders, ", ")

		if err := db.Exec(query, values...).Error; err != nil {
			return fmt.Errorf("failed to insert batch: %w", err)
		}
	}

	return nil
}

// insertDmPostsBatch はdm_postsテーブルにバッチでデータを挿入
func insertDmPostsBatch(db *gorm.DB, tableName string, dmPosts []*model.DmPost) error {
	if len(dmPosts) == 0 {
		return nil
	}

	// バッチサイズを考慮して分割
	for i := 0; i < len(dmPosts); i += batchSize {
		end := i + batchSize
		if end > len(dmPosts) {
			end = len(dmPosts)
		}
		batch := dmPosts[i:end]

		// 生SQLでバッチ挿入（動的テーブル名対応、IDを含む）
		query := fmt.Sprintf("INSERT INTO %s (id, user_id, title, content, created_at, updated_at) VALUES ", tableName)
		var values []interface{}
		var placeholders []string

		for _, dmPost := range batch {
			placeholders = append(placeholders, "(?, ?, ?, ?, ?, ?)")
			values = append(values, dmPost.ID, dmPost.UserID, dmPost.Title, dmPost.Content, dmPost.CreatedAt, dmPost.UpdatedAt)
		}

		query += strings.Join(placeholders, ", ")

		if err := db.Exec(query, values...).Error; err != nil {
			return fmt.Errorf("failed to insert batch: %w", err)
		}
	}

	return nil
}

// insertDmNewsBatch はdm_newsテーブルにバッチでデータを挿入
func insertDmNewsBatch(db *gorm.DB, dmNews []*model.DmNews) error {
	if len(dmNews) == 0 {
		return nil
	}

	// バッチサイズを考慮して分割
	for i := 0; i < len(dmNews); i += batchSize {
		end := i + batchSize
		if end > len(dmNews) {
			end = len(dmNews)
		}
		batch := dmNews[i:end]

		// GORMのCreateInBatchesを使用（固定テーブル名）
		if err := db.Table("dm_news").CreateInBatches(batch, len(batch)).Error; err != nil {
			return fmt.Errorf("failed to insert batch: %w", err)
		}
	}

	return nil
}
