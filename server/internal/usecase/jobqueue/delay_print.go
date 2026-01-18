package jobqueue

import (
	"context"
)

// DelayPrintServiceInterface はDelayPrintServiceのインターフェース
type DelayPrintServiceInterface interface {
	PrintMessage(message string) error
}

// DelayPrintUsecase はジョブ処理のビジネスロジックを実装
type DelayPrintUsecase struct {
	service DelayPrintServiceInterface
}

// NewDelayPrintUsecase は新しいDelayPrintUsecaseを作成
func NewDelayPrintUsecase(service DelayPrintServiceInterface) *DelayPrintUsecase {
	return &DelayPrintUsecase{
		service: service,
	}
}

// Execute はジョブ処理を実行
func (u *DelayPrintUsecase) Execute(ctx context.Context, message string) error {
	// デフォルトメッセージの設定（空文字列の場合）
	if message == "" {
		message = "Job executed successfully"
	}

	// サービス層の呼び出し
	return u.service.PrintMessage(message)
}
