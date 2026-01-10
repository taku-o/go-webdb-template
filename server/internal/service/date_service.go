package service

import (
	"context"
	"time"
)

// DateService は日付関連のドメインロジックを担当
type DateService struct{}

// NewDateService は新しいDateServiceを作成
func NewDateService() *DateService {
	return &DateService{}
}

// GetToday は今日の日付をYYYY-MM-DD形式で取得
func (s *DateService) GetToday(ctx context.Context) (string, error) {
	return time.Now().Format("2006-01-02"), nil
}
