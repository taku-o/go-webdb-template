package jobqueue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/taku-o/go-webdb-template/internal/service"
	usecasejobqueue "github.com/taku-o/go-webdb-template/internal/usecase/jobqueue"
)

// DelayPrintPayload は遅延出力ジョブのペイロード
type DelayPrintPayload struct {
	Message string `json:"message"`
}

// ProcessDelayPrintJob は遅延出力ジョブを処理
// 入出力制御とusecase層の呼び出しを担当
func ProcessDelayPrintJob(ctx context.Context, t *asynq.Task) error {
	// ペイロードの解析
	var payload DelayPrintPayload
	if len(t.Payload()) == 0 {
		// ペイロードがない場合は空のメッセージを設定（usecase層でデフォルト値を設定）
		payload.Message = ""
	} else {
		if err := json.Unmarshal(t.Payload(), &payload); err != nil {
			return fmt.Errorf("failed to unmarshal payload: %w", err)
		}
	}

	// usecase層の呼び出し
	delayPrintService := service.NewDelayPrintService()
	usecase := usecasejobqueue.NewDelayPrintUsecase(delayPrintService)
	return usecase.Execute(ctx, payload.Message)
}
