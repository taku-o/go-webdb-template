package usecase

import (
	"context"
)

// DateServiceInterface はDateServiceのインターフェース
type DateServiceInterface interface {
	GetToday(ctx context.Context) (string, error)
}

// TodayUsecase はtoday関連のビジネスロジックを担当
type TodayUsecase struct {
	dateService DateServiceInterface
}

// NewTodayUsecase は新しいTodayUsecaseを作成
func NewTodayUsecase(dateService DateServiceInterface) *TodayUsecase {
	return &TodayUsecase{
		dateService: dateService,
	}
}

// GetToday は今日の日付を取得
func (u *TodayUsecase) GetToday(ctx context.Context) (string, error) {
	return u.dateService.GetToday(ctx)
}
