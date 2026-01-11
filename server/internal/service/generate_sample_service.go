package service

import (
	"context"
	"fmt"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/taku-o/go-webdb-template/internal/db"
	"github.com/taku-o/go-webdb-template/internal/model"
	"github.com/taku-o/go-webdb-template/internal/util/idgen"
)

// DmUserRepositoryInterface はDmUserRepositoryのインターフェース
type DmUserRepositoryInterface interface {
	InsertDmUsersBatch(ctx context.Context, tableName string, dmUsers []*model.DmUser) error
}

// DmPostRepositoryInterface はDmPostRepositoryのインターフェース
type DmPostRepositoryInterface interface {
	InsertDmPostsBatch(ctx context.Context, tableName string, dmPosts []*model.DmPost) error
}

// DmNewsRepositoryInterface はDmNewsRepositoryのインターフェース
type DmNewsRepositoryInterface interface {
	InsertDmNewsBatch(ctx context.Context, dmNews []*model.DmNews) error
}

// GenerateSampleServiceInterface はサンプルデータ生成サービスのインターフェース
type GenerateSampleServiceInterface interface {
	GenerateDmUsers(ctx context.Context, totalCount int) ([]string, error)
	GenerateDmPosts(ctx context.Context, dmUserIDs []string, totalCount int) error
	GenerateDmNews(ctx context.Context, totalCount int) error
}

// GenerateSampleService はサンプルデータ生成のビジネスロジックを担当
type GenerateSampleService struct {
	dmUserRepository DmUserRepositoryInterface
	dmPostRepository DmPostRepositoryInterface
	dmNewsRepository DmNewsRepositoryInterface
	tableSelector    *db.TableSelector
}

// NewGenerateSampleService は新しいGenerateSampleServiceを作成
func NewGenerateSampleService(
	dmUserRepository DmUserRepositoryInterface,
	dmPostRepository DmPostRepositoryInterface,
	dmNewsRepository DmNewsRepositoryInterface,
	tableSelector *db.TableSelector,
) *GenerateSampleService {
	return &GenerateSampleService{
		dmUserRepository: dmUserRepository,
		dmPostRepository: dmPostRepository,
		dmNewsRepository: dmNewsRepository,
		tableSelector:    tableSelector,
	}
}

// GenerateDmUsers はdm_usersテーブルにデータを生成
func (s *GenerateSampleService) GenerateDmUsers(ctx context.Context, totalCount int) ([]string, error) {
	var allDmUserIDs []string

	// テーブル番号ごとにユーザーをグループ化するマップ
	usersByTable := make(map[int][]*model.DmUser)

	// 全ユーザーを生成し、IDに基づいて正しいテーブルに振り分け
	for i := 0; i < totalCount; i++ {
		id, err := idgen.GenerateUUIDv7()
		if err != nil {
			return nil, fmt.Errorf("failed to generate UUIDv7: %w", err)
		}

		// UUIDからテーブル番号を計算
		tableNumber, err := s.tableSelector.GetTableNumberFromUUID(id)
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
		// テーブル名を生成
		tableName := fmt.Sprintf("dm_users_%03d", tableNumber)

		// バッチ挿入
		if len(dmUsers) > 0 {
			if err := s.dmUserRepository.InsertDmUsersBatch(ctx, tableName, dmUsers); err != nil {
				return nil, fmt.Errorf("failed to insert batch to %s: %w", tableName, err)
			}
		}
	}

	return allDmUserIDs, nil
}

// GenerateDmPosts はdm_postsテーブルにデータを生成
func (s *GenerateSampleService) GenerateDmPosts(ctx context.Context, dmUserIDs []string, totalCount int) error {
	if len(dmUserIDs) == 0 {
		return fmt.Errorf("no dm_user IDs available for dm_posts generation")
	}

	// テーブル番号ごとに投稿をグループ化するマップ
	postsByTable := make(map[int][]*model.DmPost)

	// 全投稿を生成し、user_idに基づいて正しいテーブルに振り分け
	for i := 0; i < totalCount; i++ {
		id, err := idgen.GenerateUUIDv7()
		if err != nil {
			return fmt.Errorf("failed to generate UUIDv7: %w", err)
		}

		// dm_user_idをランダムに選択
		dmUserID := dmUserIDs[gofakeit.IntRange(0, len(dmUserIDs)-1)]

		// user_idからテーブル番号を計算（dm_postsのシャーディングキーはuser_id）
		tableNumber, err := s.tableSelector.GetTableNumberFromUUID(dmUserID)
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
		// テーブル名を生成
		tableName := fmt.Sprintf("dm_posts_%03d", tableNumber)

		// バッチ挿入
		if len(dmPosts) > 0 {
			if err := s.dmPostRepository.InsertDmPostsBatch(ctx, tableName, dmPosts); err != nil {
				return fmt.Errorf("failed to insert batch to %s: %w", tableName, err)
			}
		}
	}

	return nil
}

// GenerateDmNews はdm_newsテーブルにデータを生成
func (s *GenerateSampleService) GenerateDmNews(ctx context.Context, totalCount int) error {
	// バッチでデータ生成
	var dmNews []*model.DmNews

	// MySQLのTIMESTAMP型は1970-01-01 00:00:01以降の日付のみサポート
	minDate := time.Date(1970, 1, 2, 0, 0, 0, 0, time.UTC)
	maxDate := time.Now()

	for i := 0; i < totalCount; i++ {
		authorID := int64(gofakeit.Int32()) & 0x7FFFFFFF
		publishedAt := gofakeit.DateRange(minDate, maxDate)

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
		if err := s.dmNewsRepository.InsertDmNewsBatch(ctx, dmNews); err != nil {
			return fmt.Errorf("failed to insert batch to dm_news: %w", err)
		}
	}

	return nil
}
