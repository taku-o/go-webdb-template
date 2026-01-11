# server/cmd/generate-sample-dataの構造修正の設計書

## 1. 概要

### 1.1 設計の目的
要件定義書に基づき、`server/cmd/generate-sample-data`の実装を、APIサーバーと同じレイヤー構造（usecase -> service -> repository -> db）に変更するための詳細設計を定義する。これにより、CLIコマンドとAPIサーバーで一貫したアーキテクチャを実現し、サンプルデータ生成処理を共通化してコードの保守性と再利用性を向上させる。

### 1.2 設計の範囲
- Repository層の拡張（既存のrepositoryファイルにバッチ挿入メソッドを追加、`dm_news_repository.go`を新規作成）の設計
- Service層（`server/internal/service/generate_sample_service.go`）の設計
- CLI用usecase層（`server/internal/usecase/cli/generate_sample_usecase.go`）の設計
- `server/cmd/generate-sample-data/main.go`の簡素化設計
- 依存関係の注入設計
- テスト設計
- ドキュメント更新の設計

### 1.3 設計方針
- **一貫性**: APIサーバーと同じレイヤー構造を採用
- **既存コードの活用**: 既存のrepositoryファイルを拡張して使用
- **責務の明確化**: 各レイヤーの責務を明確に分離
- **テスト容易性**: 各層を独立してテストできる設計
- **後方互換性**: 既存のCLIコマンドの動作（出力形式、エラーメッセージ）を維持
- **データ生成ロジックの維持**: 既存のデータ生成ロジック（gofakeitの使用、バッチサイズ、テーブル分割など）を維持

## 2. アーキテクチャ設計

### 2.1 全体構成

```
┌─────────────────────────────────────────────────────────────┐
│          CLI Layer (cmd/generate-sample-data/main.go)         │
│  • エントリーポイント                                         │
│  • 設定ファイルの読み込み                                     │
│  • GroupManagerの初期化                                      │
│  • レイヤーの初期化（Repository → Service → Usecase）        │
│  • usecase層の呼び出し                                       │
│  • 結果の出力（標準出力にログ出力）                           │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│      Usecase Layer (internal/usecase/cli)                     │
│  • GenerateSampleUsecase                                     │
│  • ビジネスロジックの調整（CLI用）                           │
│  • GenerateDmUsers() → GenerateDmPosts() → GenerateDmNews()  │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│         Service Layer (internal/service)                        │
│  • GenerateSampleService                                     │
│  • ドメインロジック                                           │
│  • gofakeitを使用したデータ生成                              │
│  • UUID生成、テーブル番号計算                                 │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│      Repository Layer (internal/repository)                   │
│  • DmUserRepository.InsertDmUsersBatch()                     │
│  • DmPostRepository.InsertDmPostsBatch()                     │
│  • DmNewsRepository.InsertDmNewsBatch()                      │
│  • バッチ挿入処理（500件ずつ）                               │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│         DB Layer (internal/db)                                 │
│  • GroupManager                                              │
│  • Sharding接続管理                                          │
│  • Master接続管理                                            │
└─────────────────────────────────────────────────────────────┘
```

### 2.2 データフロー

```
main.go
  ↓
設定ファイルの読み込み（config.Load()）
  ↓
GroupManagerの初期化（db.NewGroupManager(cfg)）
  ↓
Repository層の初期化
  - repository.NewDmUserRepository(groupManager)
  - repository.NewDmPostRepository(groupManager)
  - repository.NewDmNewsRepository(groupManager)
  ↓
Service層の初期化
  - service.NewGenerateSampleService(dmUserRepository, dmPostRepository, dmNewsRepository, tableSelector)
  ↓
Usecase層の初期化
  - cli.NewGenerateSampleUsecase(generateSampleService)
  ↓
usecase.GenerateSampleData(ctx, totalCount)
  ↓
service.GenerateDmUsers(ctx, totalCount)
  ↓
  - UUID生成（idgen.GenerateUUIDv7()）
  - テーブル番号計算（tableSelector.GetTableNumberFromUUID()）
  - データ生成（gofakeit）
  - repository.InsertDmUsersBatch()
  ↓
service.GenerateDmPosts(ctx, dmUserIDs, totalCount)
  ↓
  - UUID生成（idgen.GenerateUUIDv7()）
  - user_idからテーブル番号計算（tableSelector.GetTableNumberFromUUID()）
  - データ生成（gofakeit）
  - repository.InsertDmPostsBatch()
  ↓
service.GenerateDmNews(ctx, totalCount)
  ↓
  - データ生成（gofakeit）
  - repository.InsertDmNewsBatch()
  ↓
結果の出力（標準出力にログ出力）
```

### 2.3 レイヤー構造の比較

#### 修正前
```
main.go
  ↓ (直接実装)
generateDmUsers()
  ↓
insertDmUsersBatch()
  ↓
generateDmPosts()
  ↓
insertDmPostsBatch()
  ↓
generateDmNews()
  ↓
insertDmNewsBatch()
  ↓
標準出力にログ出力
```

#### 修正後
```
main.go
  ↓
usecase/cli.GenerateSampleUsecase.GenerateSampleData()
  ↓
service.GenerateSampleService.GenerateDmUsers()
  ↓
repository.DmUserRepository.InsertDmUsersBatch()
  ↓
service.GenerateSampleService.GenerateDmPosts()
  ↓
repository.DmPostRepository.InsertDmPostsBatch()
  ↓
service.GenerateSampleService.GenerateDmNews()
  ↓
repository.DmNewsRepository.InsertDmNewsBatch()
  ↓
標準出力にログ出力
```

## 3. 詳細設計

### 3.1 Repository層の拡張設計

#### 3.1.1 `dm_user_repository.go`へのメソッド追加

**ファイルパス**: `server/internal/repository/dm_user_repository.go`

**追加するメソッド**:

```go
// InsertDmUsersBatch はdm_usersテーブルにバッチでデータを挿入
func (r *DmUserRepository) InsertDmUsersBatch(ctx context.Context, tableName string, dmUsers []*model.DmUser) error {
	if len(dmUsers) == 0 {
		return nil
	}

	const batchSize = 500

	// バッチサイズを考慮して分割
	for i := 0; i < len(dmUsers); i += batchSize {
		end := i + batchSize
		if end > len(dmUsers) {
			end = len(dmUsers)
		}
		batch := dmUsers[i:end]

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

		// 生SQLでバッチ挿入（動的テーブル名対応、IDを含む）
		query := fmt.Sprintf("INSERT INTO %s (id, name, email, created_at, updated_at) VALUES ", tableName)
		var values []interface{}
		var placeholders []string

		for _, dmUser := range batch {
			placeholders = append(placeholders, "(?, ?, ?, ?, ?)")
			values = append(values, dmUser.ID, dmUser.Name, dmUser.Email, dmUser.CreatedAt, dmUser.UpdatedAt)
		}

		query += strings.Join(placeholders, ", ")

		// リトライ機能付きでクエリ実行
		err = db.ExecuteWithRetry(func() error {
			return conn.DB.WithContext(ctx).Exec(query, values...).Error
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
```

**設計のポイント**:
- 既存の`DmUserRepository`構造体にメソッドを追加
- バッチサイズ500件ずつ挿入する処理を実装
- 既存の`insertDmUsersBatch`関数のロジックを移行
- リトライ機能（`db.ExecuteWithRetry`）を使用
- 動的テーブル名に対応
- `context.Context`を受け取る

#### 3.1.2 `dm_post_repository.go`へのメソッド追加

**ファイルパス**: `server/internal/repository/dm_post_repository.go`

**追加するメソッド**:

```go
// InsertDmPostsBatch はdm_postsテーブルにバッチでデータを挿入
func (r *DmPostRepository) InsertDmPostsBatch(ctx context.Context, tableName string, dmPosts []*model.DmPost) error {
	if len(dmPosts) == 0 {
		return nil
	}

	const batchSize = 500

	// バッチサイズを考慮して分割
	for i := 0; i < len(dmPosts); i += batchSize {
		end := i + batchSize
		if end > len(dmPosts) {
			end = len(dmPosts)
		}
		batch := dmPosts[i:end]

		// テーブル番号から接続を取得
		tableNumber, err := extractTableNumber(tableName, "dm_posts_")
		if err != nil {
			return fmt.Errorf("failed to extract table number from %s: %w", tableName, err)
		}

		conn, err := r.groupManager.GetShardingConnection(tableNumber)
		if err != nil {
			return fmt.Errorf("failed to get connection for table %d: %w", tableNumber, err)
		}

		// 生SQLでバッチ挿入（動的テーブル名対応、IDを含む）
		query := fmt.Sprintf("INSERT INTO %s (id, user_id, title, content, created_at, updated_at) VALUES ", tableName)
		var values []interface{}
		var placeholders []string

		for _, dmPost := range batch {
			placeholders = append(placeholders, "(?, ?, ?, ?, ?, ?)")
			values = append(values, dmPost.ID, dmPost.UserID, dmPost.Title, dmPost.Content, dmPost.CreatedAt, dmPost.UpdatedAt)
		}

		query += strings.Join(placeholders, ", ")

		// リトライ機能付きでクエリ実行
		err = db.ExecuteWithRetry(func() error {
			return conn.DB.WithContext(ctx).Exec(query, values...).Error
		})
		if err != nil {
			return fmt.Errorf("failed to insert batch: %w", err)
		}
	}

	return nil
}
```

**設計のポイント**:
- 既存の`DmPostRepository`構造体にメソッドを追加
- バッチサイズ500件ずつ挿入する処理を実装
- 既存の`insertDmPostsBatch`関数のロジックを移行
- リトライ機能（`db.ExecuteWithRetry`）を使用
- 動的テーブル名に対応
- `context.Context`を受け取る

#### 3.1.3 `dm_news_repository.go`の新規作成

**ファイルパス**: `server/internal/repository/dm_news_repository.go`

**実装内容**:

```go
package repository

import (
	"context"
	"fmt"

	"github.com/taku-o/go-webdb-template/internal/db"
	"github.com/taku-o/go-webdb-template/internal/model"
	"gorm.io/gorm"
)

// DmNewsRepository はニュースのデータアクセスを担当
type DmNewsRepository struct {
	groupManager *db.GroupManager
}

// NewDmNewsRepository は新しいDmNewsRepositoryを作成
func NewDmNewsRepository(groupManager *db.GroupManager) *DmNewsRepository {
	return &DmNewsRepository{
		groupManager: groupManager,
	}
}

// InsertDmNewsBatch はdm_newsテーブルにバッチでデータを挿入
func (r *DmNewsRepository) InsertDmNewsBatch(ctx context.Context, dmNews []*model.DmNews) error {
	if len(dmNews) == 0 {
		return nil
	}

	const batchSize = 500

	// master接続を取得
	conn, err := r.groupManager.GetMasterConnection()
	if err != nil {
		return fmt.Errorf("failed to get master connection: %w", err)
	}

	// バッチサイズを考慮して分割
	for i := 0; i < len(dmNews); i += batchSize {
		end := i + batchSize
		if end > len(dmNews) {
			end = len(dmNews)
		}
		batch := dmNews[i:end]

		// GORMのCreateInBatchesを使用（固定テーブル名）
		err = db.ExecuteWithRetry(func() error {
			return conn.DB.WithContext(ctx).Table("dm_news").CreateInBatches(batch, len(batch)).Error
		})
		if err != nil {
			return fmt.Errorf("failed to insert batch: %w", err)
		}
	}

	return nil
}
```

**設計のポイント**:
- `DmNewsRepository`構造体を新規作成
- master接続を使用（シャーディング不要）
- バッチサイズ500件ずつ挿入する処理を実装
- 既存の`insertDmNewsBatch`関数のロジックを移行
- リトライ機能（`db.ExecuteWithRetry`）を使用
- 固定テーブル名（`dm_news`）を使用
- `context.Context`を受け取る

#### 3.1.4 Repository層のテスト設計

**テストファイル**:
- `server/internal/repository/dm_user_repository_test.go`（既存ファイルに追加、存在する場合）
- `server/internal/repository/dm_post_repository_test.go`（既存ファイルに追加、存在する場合）
- `server/internal/repository/dm_news_repository_test.go`（新規作成）

**テスト内容**:
- バッチ挿入の正常系テスト
- 空のスライスを渡した場合のテスト
- エラーハンドリングのテスト
- バッチサイズを超えるデータのテスト

### 3.2 Service層の設計

#### 3.2.1 `generate_sample_service.go`の設計

**ファイルパス**: `server/internal/service/generate_sample_service.go`

**実装内容**:

```go
package service

import (
	"context"
	"fmt"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/taku-o/go-webdb-template/internal/db"
	"github.com/taku-o/go-webdb-template/internal/model"
	"github.com/taku-o/go-webdb-template/internal/repository"
	"github.com/taku-o/go-webdb-template/internal/util/idgen"
)

// GenerateSampleServiceInterface はサンプルデータ生成サービスのインターフェース
type GenerateSampleServiceInterface interface {
	GenerateDmUsers(ctx context.Context, totalCount int) ([]string, error)
	GenerateDmPosts(ctx context.Context, dmUserIDs []string, totalCount int) error
	GenerateDmNews(ctx context.Context, totalCount int) error
}

// GenerateSampleService はサンプルデータ生成のビジネスロジックを担当
type GenerateSampleService struct {
	dmUserRepository *repository.DmUserRepository
	dmPostRepository *repository.DmPostRepository
	dmNewsRepository *repository.DmNewsRepository
	tableSelector    *db.TableSelector
}

// NewGenerateSampleService は新しいGenerateSampleServiceを作成
func NewGenerateSampleService(
	dmUserRepository *repository.DmUserRepository,
	dmPostRepository *repository.DmPostRepository,
	dmNewsRepository *repository.DmNewsRepository,
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
		if err := s.dmNewsRepository.InsertDmNewsBatch(ctx, dmNews); err != nil {
			return fmt.Errorf("failed to insert batch to dm_news: %w", err)
		}
	}

	return nil
}
```

**設計のポイント**:
- `GenerateSampleServiceInterface`を新規作成（usecase層で使用するため）
- `DmUserRepository`、`DmPostRepository`、`DmNewsRepository`を依存として注入
- `TableSelector`を依存として注入（テーブル番号計算用）
- 既存の`generateDmUsers`、`generateDmPosts`、`generateDmNews`関数のロジックを移行
- `gofakeit`を使用したデータ生成ロジックを維持
- UUID生成、テーブル番号計算、データ生成、repository層への呼び出しを実装
- エラーハンドリング: repository層から返されたエラーをそのまま返す（エラーのラップは不要）
- `context.Context`を受け取る

#### 3.2.2 `generate_sample_service_test.go`の設計

**ファイルパス**: `server/internal/service/generate_sample_service_test.go`

**実装内容**:

```go
package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/taku-o/go-webdb-template/internal/model"
	"github.com/taku-o/go-webdb-template/internal/repository"
)

// MockDmUserRepository はDmUserRepositoryのモック
type MockDmUserRepository struct {
	InsertDmUsersBatchFunc func(ctx context.Context, tableName string, dmUsers []*model.DmUser) error
}

func (m *MockDmUserRepository) InsertDmUsersBatch(ctx context.Context, tableName string, dmUsers []*model.DmUser) error {
	if m.InsertDmUsersBatchFunc != nil {
		return m.InsertDmUsersBatchFunc(ctx, tableName, dmUsers)
	}
	return nil
}

// MockDmPostRepository はDmPostRepositoryのモック
type MockDmPostRepository struct {
	InsertDmPostsBatchFunc func(ctx context.Context, tableName string, dmPosts []*model.DmPost) error
}

func (m *MockDmPostRepository) InsertDmPostsBatch(ctx context.Context, tableName string, dmPosts []*model.DmPost) error {
	if m.InsertDmPostsBatchFunc != nil {
		return m.InsertDmPostsBatchFunc(ctx, tableName, dmPosts)
	}
	return nil
}

// MockDmNewsRepository はDmNewsRepositoryのモック
type MockDmNewsRepository struct {
	InsertDmNewsBatchFunc func(ctx context.Context, dmNews []*model.DmNews) error
}

func (m *MockDmNewsRepository) InsertDmNewsBatch(ctx context.Context, dmNews []*model.DmNews) error {
	if m.InsertDmNewsBatchFunc != nil {
		return m.InsertDmNewsBatchFunc(ctx, dmNews)
	}
	return nil
}

func TestGenerateSampleService_GenerateDmUsers(t *testing.T) {
	tests := []struct {
		name        string
		mockFunc    func(ctx context.Context, tableName string, dmUsers []*model.DmUser) error
		totalCount  int
		wantError   bool
		expectedErr string
	}{
		{
			name: "success",
			mockFunc: func(ctx context.Context, tableName string, dmUsers []*model.DmUser) error {
				return nil
			},
			totalCount: 10,
			wantError:  false,
		},
		{
			name: "repository error",
			mockFunc: func(ctx context.Context, tableName string, dmUsers []*model.DmUser) error {
				return errors.New("failed to insert batch")
			},
			totalCount:  10,
			wantError:  true,
			expectedErr: "failed to insert batch",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepo := &MockDmUserRepository{
				InsertDmUsersBatchFunc: tt.mockFunc,
			}
			mockPostRepo := &MockDmPostRepository{}
			mockNewsRepo := &MockDmNewsRepository{}

			// TableSelectorのモックは実際の実装を使用（またはモックを作成）
			// ここでは簡略化のため、実際のTableSelectorを使用する想定

			service := NewGenerateSampleService(
				&repository.DmUserRepository{}, // 実際の実装ではモックを使用
				&repository.DmPostRepository{},
				&repository.DmNewsRepository{},
				nil, // TableSelectorは実際の実装を使用
			)

			ctx := context.Background()
			got, err := service.GenerateDmUsers(ctx, tt.totalCount)

			if tt.wantError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				assert.Empty(t, got)
			} else {
				assert.NoError(t, err)
				assert.Len(t, got, tt.totalCount)
			}
		})
	}
}

// GenerateDmPosts、GenerateDmNewsのテストも同様に実装
```

**設計のポイント**:
- 関数ポインタを使用してモックを実装（既存のテストパターンに合わせる）
- テーブル駆動テストを使用
- 正常系と異常系の両方をテスト
- `github.com/stretchr/testify/assert`を使用してアサーション
- `DmUserRepository`、`DmPostRepository`、`DmNewsRepository`のモックを使用

### 3.3 CLI用usecase層の設計

#### 3.3.1 ディレクトリ構造

```
server/internal/usecase/cli/
├── list_dm_users_usecase.go
├── list_dm_users_usecase_test.go
├── generate_secret_usecase.go
├── generate_secret_usecase_test.go
├── generate_sample_usecase.go
└── generate_sample_usecase_test.go
```

#### 3.3.2 `generate_sample_usecase.go`の設計

**ファイルパス**: `server/internal/usecase/cli/generate_sample_usecase.go`

**実装内容**:

```go
package cli

import (
	"context"

	"github.com/taku-o/go-webdb-template/internal/service"
)

// GenerateSampleUsecase はCLI用のサンプルデータ生成usecase
type GenerateSampleUsecase struct {
	generateSampleService service.GenerateSampleServiceInterface
}

// NewGenerateSampleUsecase は新しいGenerateSampleUsecaseを作成
func NewGenerateSampleUsecase(generateSampleService service.GenerateSampleServiceInterface) *GenerateSampleUsecase {
	return &GenerateSampleUsecase{
		generateSampleService: generateSampleService,
	}
}

// GenerateSampleData はサンプルデータを生成
func (u *GenerateSampleUsecase) GenerateSampleData(ctx context.Context, totalCount int) error {
	// 1. dm_usersテーブルへのデータ生成
	dmUserIDs, err := u.generateSampleService.GenerateDmUsers(ctx, totalCount)
	if err != nil {
		return err
	}

	// 2. dm_postsテーブルへのデータ生成
	if err := u.generateSampleService.GenerateDmPosts(ctx, dmUserIDs, totalCount); err != nil {
		return err
	}

	// 3. dm_newsテーブルへのデータ生成
	if err := u.generateSampleService.GenerateDmNews(ctx, totalCount); err != nil {
		return err
	}

	return nil
}
```

**設計のポイント**:
- `GenerateSampleServiceInterface`を使用（依存関係の注入）
- コンストラクタで依存関係を注入
- service層のメソッドを順次呼び出す（`GenerateDmUsers()` → `GenerateDmPosts()` → `GenerateDmNews()`）
- エラーハンドリング: service層から返されたエラーをそのまま返す（エラーのラップは不要）
- `context.Context`を受け取る

#### 3.3.3 `generate_sample_usecase_test.go`の設計

**ファイルパス**: `server/internal/usecase/cli/generate_sample_usecase_test.go`

**実装内容**:

```go
package cli

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/taku-o/go-webdb-template/internal/service"
)

// MockGenerateSampleServiceInterface はGenerateSampleServiceInterfaceのモック
type MockGenerateSampleServiceInterface struct {
	GenerateDmUsersFunc func(ctx context.Context, totalCount int) ([]string, error)
	GenerateDmPostsFunc func(ctx context.Context, dmUserIDs []string, totalCount int) error
	GenerateDmNewsFunc  func(ctx context.Context, totalCount int) error
}

func (m *MockGenerateSampleServiceInterface) GenerateDmUsers(ctx context.Context, totalCount int) ([]string, error) {
	if m.GenerateDmUsersFunc != nil {
		return m.GenerateDmUsersFunc(ctx, totalCount)
	}
	return []string{}, nil
}

func (m *MockGenerateSampleServiceInterface) GenerateDmPosts(ctx context.Context, dmUserIDs []string, totalCount int) error {
	if m.GenerateDmPostsFunc != nil {
		return m.GenerateDmPostsFunc(ctx, dmUserIDs, totalCount)
	}
	return nil
}

func (m *MockGenerateSampleServiceInterface) GenerateDmNews(ctx context.Context, totalCount int) error {
	if m.GenerateDmNewsFunc != nil {
		return m.GenerateDmNewsFunc(ctx, totalCount)
	}
	return nil
}

func TestGenerateSampleUsecase_GenerateSampleData(t *testing.T) {
	tests := []struct {
		name        string
		mockFuncs   struct {
			generateDmUsers func(ctx context.Context, totalCount int) ([]string, error)
			generateDmPosts func(ctx context.Context, dmUserIDs []string, totalCount int) error
			generateDmNews  func(ctx context.Context, totalCount int) error
		}
		totalCount  int
		wantError   bool
		expectedErr string
	}{
		{
			name: "success",
			mockFuncs: struct {
				generateDmUsers func(ctx context.Context, totalCount int) ([]string, error)
				generateDmPosts func(ctx context.Context, dmUserIDs []string, totalCount int) error
				generateDmNews  func(ctx context.Context, totalCount int) error
			}{
				generateDmUsers: func(ctx context.Context, totalCount int) ([]string, error) {
					return []string{"user1", "user2"}, nil
				},
				generateDmPosts: func(ctx context.Context, dmUserIDs []string, totalCount int) error {
					return nil
				},
				generateDmNews: func(ctx context.Context, totalCount int) error {
					return nil
				},
			},
			totalCount: 10,
			wantError:  false,
		},
		{
			name: "generateDmUsers error",
			mockFuncs: struct {
				generateDmUsers func(ctx context.Context, totalCount int) ([]string, error)
				generateDmPosts func(ctx context.Context, dmUserIDs []string, totalCount int) error
				generateDmNews  func(ctx context.Context, totalCount int) error
			}{
				generateDmUsers: func(ctx context.Context, totalCount int) ([]string, error) {
					return nil, errors.New("failed to generate users")
				},
				generateDmPosts: nil,
				generateDmNews:  nil,
			},
			totalCount:  10,
			wantError:  true,
			expectedErr: "failed to generate users",
		},
		{
			name: "generateDmPosts error",
			mockFuncs: struct {
				generateDmUsers func(ctx context.Context, totalCount int) ([]string, error)
				generateDmPosts func(ctx context.Context, dmUserIDs []string, totalCount int) error
				generateDmNews  func(ctx context.Context, totalCount int) error
			}{
				generateDmUsers: func(ctx context.Context, totalCount int) ([]string, error) {
					return []string{"user1", "user2"}, nil
				},
				generateDmPosts: func(ctx context.Context, dmUserIDs []string, totalCount int) error {
					return errors.New("failed to generate posts")
				},
				generateDmNews: nil,
			},
			totalCount:  10,
			wantError:  true,
			expectedErr: "failed to generate posts",
		},
		{
			name: "generateDmNews error",
			mockFuncs: struct {
				generateDmUsers func(ctx context.Context, totalCount int) ([]string, error)
				generateDmPosts func(ctx context.Context, dmUserIDs []string, totalCount int) error
				generateDmNews  func(ctx context.Context, totalCount int) error
			}{
				generateDmUsers: func(ctx context.Context, totalCount int) ([]string, error) {
					return []string{"user1", "user2"}, nil
				},
				generateDmPosts: func(ctx context.Context, dmUserIDs []string, totalCount int) error {
					return nil
				},
				generateDmNews: func(ctx context.Context, totalCount int) error {
					return errors.New("failed to generate news")
				},
			},
			totalCount:  10,
			wantError:  true,
			expectedErr: "failed to generate news",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockGenerateSampleServiceInterface{
				GenerateDmUsersFunc: tt.mockFuncs.generateDmUsers,
				GenerateDmPostsFunc: tt.mockFuncs.generateDmPosts,
				GenerateDmNewsFunc:  tt.mockFuncs.generateDmNews,
			}

			usecase := NewGenerateSampleUsecase(mockService)

			ctx := context.Background()
			err := usecase.GenerateSampleData(ctx, tt.totalCount)

			if tt.wantError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
```

**設計のポイント**:
- 関数ポインタを使用してモックを実装（既存のテストパターンに合わせる）
- テーブル駆動テストを使用
- 正常系と異常系の両方をテスト（各処理でエラーが発生する場合）
- `github.com/stretchr/testify/assert`を使用してアサーション

### 3.4 main.goの簡素化設計

#### 3.4.1 修正後のmain.goの構造

**ファイルパス**: `server/cmd/generate-sample-data/main.go`

**実装内容**:

```go
package main

import (
	"context"
	"log"
	"os"

	"github.com/taku-o/go-webdb-template/internal/config"
	"github.com/taku-o/go-webdb-template/internal/db"
	"github.com/taku-o/go-webdb-template/internal/repository"
	"github.com/taku-o/go-webdb-template/internal/service"
	"github.com/taku-o/go-webdb-template/internal/usecase/cli"
)

const (
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

	// 4. Repository層の初期化
	dmUserRepository := repository.NewDmUserRepository(groupManager)
	dmPostRepository := repository.NewDmPostRepository(groupManager)
	dmNewsRepository := repository.NewDmNewsRepository(groupManager)

	// 5. Service層の初期化
	tableSelector := db.NewTableSelector(db.DBShardingTableCount, db.DBShardingTablesPerDB)
	generateSampleService := service.NewGenerateSampleService(
		dmUserRepository,
		dmPostRepository,
		dmNewsRepository,
		tableSelector,
	)

	// 6. Usecase層の初期化
	generateSampleUsecase := cli.NewGenerateSampleUsecase(generateSampleService)

	// 7. サンプルデータの生成
	ctx := context.Background()
	if err := generateSampleUsecase.GenerateSampleData(ctx, totalCount); err != nil {
		log.Fatalf("Failed to generate sample data: %v", err)
	}

	// 8. 生成完了メッセージ
	log.Println("Sample data generation completed successfully")
	os.Exit(0)
}
```

**変更点**:
1. `internal/repository`パッケージをインポート
2. `internal/service`パッケージをインポート
3. `internal/usecase/cli`パッケージをインポート
4. Repository層の初期化を追加（`repository.NewDmUserRepository()`、`repository.NewDmPostRepository()`、`repository.NewDmNewsRepository()`）
5. Service層の初期化を追加（`service.NewGenerateSampleService()`）
6. Usecase層の初期化を追加（`cli.NewGenerateSampleUsecase()`）
7. サンプルデータ生成処理をusecase層の呼び出しに変更
8. 既存の出力形式（標準出力にログ出力）を維持
9. 既存のエラーハンドリング（`log.Fatalf`）を維持

**削除した内容**:
- `gofakeit`の直接使用（service層に移動）
- `idgen`の直接使用（service層に移動）
- データ生成処理の直接実装（service層に移動）
- バッチ挿入処理の直接実装（repository層に移動）
- `generateDmUsers`、`generateDmPosts`、`generateDmNews`関数（service層に移動）
- `insertDmUsersBatch`、`insertDmPostsBatch`、`insertDmNewsBatch`関数（repository層に移動）

#### 3.4.2 ビルド方法

**ビルドコマンド**:
```bash
cd server
go build -o bin/generate-sample-data ./cmd/generate-sample-data
```

**ビルド出力先**: `server/bin/generate-sample-data`

**実行方法**:
```bash
cd server
APP_ENV=develop ./bin/generate-sample-data
```

**注意事項**:
- ビルド出力先は`server/bin/`ディレクトリに統一
- `server/bin/`ディレクトリは`.gitignore`に含まれているため、ビルド成果物はGitにコミットされない
- ビルド前に`server/bin/`ディレクトリが存在しない場合は作成する必要がある

### 3.5 依存関係の注入設計

#### 3.5.1 初期化の順序

```
1. config.Load()
   ↓
2. db.NewGroupManager(cfg)
   ↓
3. repository.NewDmUserRepository(groupManager)
   ↓
4. repository.NewDmPostRepository(groupManager)
   ↓
5. repository.NewDmNewsRepository(groupManager)
   ↓
6. db.NewTableSelector(...)
   ↓
7. service.NewGenerateSampleService(dmUserRepository, dmPostRepository, dmNewsRepository, tableSelector)
   ↓
8. cli.NewGenerateSampleUsecase(generateSampleService)
   ↓
9. usecase.GenerateSampleData(ctx, totalCount)
```

#### 3.5.2 依存関係の図

```
GenerateSampleUsecase
  └── GenerateSampleServiceInterface
        └── GenerateSampleService
              ├── DmUserRepository
              │     └── GroupManager
              ├── DmPostRepository
              │     └── GroupManager
              ├── DmNewsRepository
              │     └── GroupManager
              └── TableSelector
```

## 4. テスト設計

### 4.1 Repository層のテスト

#### 4.1.1 `dm_user_repository_test.go`への追加

**テスト内容**:
- `InsertDmUsersBatch`の正常系テスト
- 空のスライスを渡した場合のテスト
- バッチサイズを超えるデータのテスト
- エラーハンドリングのテスト

#### 4.1.2 `dm_post_repository_test.go`への追加

**テスト内容**:
- `InsertDmPostsBatch`の正常系テスト
- 空のスライスを渡した場合のテスト
- バッチサイズを超えるデータのテスト
- エラーハンドリングのテスト

#### 4.1.3 `dm_news_repository_test.go`の新規作成

**テスト内容**:
- `InsertDmNewsBatch`の正常系テスト
- 空のスライスを渡した場合のテスト
- バッチサイズを超えるデータのテスト
- エラーハンドリングのテスト

### 4.2 Service層のテスト

#### 4.2.1 `generate_sample_service_test.go`の設計

**テスト内容**:
- `GenerateDmUsers`の正常系テスト
- `GenerateDmPosts`の正常系テスト
- `GenerateDmNews`の正常系テスト
- repository層のエラーハンドリングのテスト
- UUID生成エラーのテスト
- テーブル番号計算エラーのテスト

### 4.3 Usecase層のテスト

#### 4.3.1 `generate_sample_usecase_test.go`の設計

**テスト内容**:
- `GenerateSampleData`の正常系テスト
- `GenerateDmUsers`でエラーが発生する場合のテスト
- `GenerateDmPosts`でエラーが発生する場合のテスト
- `GenerateDmNews`でエラーが発生する場合のテスト

## 5. エラーハンドリング設計

### 5.1 エラーの伝播

1. **Repository層**: Go errorsを返却
2. **Service層**: エラーをそのまま返す（エラーのラップは不要）
3. **Usecase層**: エラーをそのまま返す（エラーのラップは不要）
4. **CLI層（main.go）**: エラーを`log.Fatalf`で出力して終了

### 5.2 エラーメッセージ

- 既存のエラーメッセージを維持
- エラーメッセージは既存の形式に合わせる（例: "Failed to generate users: %v"）

## 6. 定数の定義

### 6.1 定数の配置

- `batchSize = 500`: repository層の各メソッド内で定数として定義（または共通定数として定義）
- `tableCount = 32`: service層で`db.DBShardingTableCount`を使用（既存の定数を使用）
- `totalCount = 100`: main.goで定数として定義（またはコマンドライン引数として受け取る）

## 7. ドキュメント更新の設計

### 7.1 `docs/Architecture.md`の更新

**更新内容**:
- CLIコマンドのレイヤー構造を追加（usecase層を含む）
- CLIコマンドのアーキテクチャ図を更新
- CLI用usecase層の説明を追加

### 7.2 `docs/Project-Structure.md`の更新

**更新内容**:
- `server/internal/usecase/cli/generate_sample_usecase.go`を追加
- `server/internal/service/generate_sample_service.go`を追加
- `server/internal/repository/dm_news_repository.go`を追加（新規作成）

### 7.3 `docs/Generate-Sample-Data.md`の更新

**更新内容**:
- CLIコマンドのアーキテクチャ図を更新（usecase層を追加）
- レイヤー構造の説明を更新（main.go → usecase → service → repository → db → 出力）

### 7.4 `.kiro/steering/structure.md`の更新

**更新内容**:
- `server/internal/usecase/cli/generate_sample_usecase.go`を追加
- `server/internal/service/generate_sample_service.go`を追加
- `server/internal/repository/dm_news_repository.go`を追加（新規作成）

## 8. 実装上の注意事項

### 8.1 Repository層の実装

- **既存のrepositoryファイルへの追加**: `dm_user_repository.go`、`dm_post_repository.go`にバッチ挿入メソッドを追加する際は、既存のコードスタイルに合わせる
- **新規repositoryファイルの作成**: `dm_news_repository.go`を新規作成する際は、既存のrepositoryファイルの構造に合わせる
- **バッチサイズ**: 既存のバッチサイズ（500件）を維持
- **エラーハンドリング**: 既存のエラーハンドリングを維持
- **GORMの使用**: 既存のGORMの使用方法を維持
- **動的テーブル名**: 既存の動的テーブル名の生成方法を維持（dm_users、dm_postsは動的テーブル名、dm_newsは固定テーブル名）

### 8.2 Service層の実装

- **依存関係の注入**: コンストラクタで依存関係を注入（`DmUserRepository`、`DmPostRepository`、`DmNewsRepository`、`TableSelector`）
- **エラーハンドリング**: repository層から返されたエラーをそのまま返す（エラーのラップは不要）
- **データ生成ロジック**: 既存の`gofakeit`を使用したデータ生成ロジックを維持
- **UUID生成**: 既存のUUID生成方法（`idgen.GenerateUUIDv7()`）を維持
- **テーブル番号計算**: 既存のテーブル番号計算方法（`tableSelector.GetTableNumberFromUUID()`）を維持

### 8.3 usecase層の実装

- **インターフェースの使用**: 新規作成する`GenerateSampleServiceInterface`を使用
- **依存関係の注入**: コンストラクタで`GenerateSampleServiceInterface`を注入
- **エラーハンドリング**: service層から返されたエラーをそのまま返す（エラーのラップは不要）
- **処理の順序**: `GenerateDmUsers()` → `GenerateDmPosts()` → `GenerateDmNews()`の順で実行

### 8.4 main.goの修正

- **usecase層の初期化**: Repository → Service → Usecaseの順で初期化
- **既存の出力**: 標準出力にログ出力（既存の動作を維持）
- **エラーハンドリング**: 既存のエラーハンドリングを維持
- **ビルド出力先**: ビルド時のバイナリは`server/bin/generate-sample-data`に出力する
- **設定ファイルの読み込み**: 既存の設定ファイル読み込み方法を維持
- **GroupManagerの初期化**: 既存のGroupManager初期化方法を維持

## 9. 参考情報

### 9.1 既存実装の参考

- `server/internal/usecase/cli/list_dm_users_usecase.go`: 既存のCLI用usecase層の実装パターン
- `server/internal/usecase/cli/generate_secret_usecase.go`: 既存のCLI用usecase層の実装パターン
- `server/cmd/generate-sample-data/main.go`: 既存のCLI実装
- `server/internal/service/secret_service.go`: 既存のservice層の実装パターン
- `server/internal/repository/dm_user_repository.go`: 既存のrepository層の実装パターン
- `server/internal/repository/dm_post_repository.go`: 既存のrepository層の実装パターン

### 9.2 技術スタック

- **言語**: Go
- **アーキテクチャ**: レイヤードアーキテクチャ（usecase -> service -> repository -> db -> 出力）
- **テスト**: `testing`（標準ライブラリ）、`github.com/stretchr/testify`（アサーション、モック）
- **データ生成**: `github.com/brianvoe/gofakeit/v6`（ランダムデータ生成）
- **UUID生成**: `github.com/taku-o/go-webdb-template/internal/util/idgen`（UUIDv7生成）
