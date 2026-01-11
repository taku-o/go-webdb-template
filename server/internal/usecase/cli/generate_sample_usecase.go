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
