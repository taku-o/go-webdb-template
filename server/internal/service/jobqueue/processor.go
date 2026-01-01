package jobqueue

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/hibiken/asynq"
)

// DelayPrintPayload は遅延出力ジョブのペイロード
type DelayPrintPayload struct {
	Message string `json:"message"`
}

// ProcessDelayPrintJob は遅延出力ジョブを処理
func ProcessDelayPrintJob(ctx context.Context, t *asynq.Task) error {
	// ペイロードの解析
	var payload DelayPrintPayload
	if len(t.Payload()) == 0 {
		// ペイロードがない場合はデフォルトメッセージを使用
		payload.Message = "Job executed successfully"
	} else {
		if err := json.Unmarshal(t.Payload(), &payload); err != nil {
			return fmt.Errorf("failed to unmarshal payload: %w", err)
		}
	}

	// メッセージが空の場合はデフォルトメッセージを使用
	if payload.Message == "" {
		payload.Message = "Job executed successfully"
	}

	// 標準出力に文字列を出力
	fmt.Printf("[%s] %s\n", time.Now().Format("2006-01-02 15:04:05"), payload.Message)
	os.Stdout.Sync() // バッファをフラッシュ

	return nil
}
