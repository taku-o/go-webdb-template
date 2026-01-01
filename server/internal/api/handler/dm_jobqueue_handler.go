package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/taku-o/go-webdb-template/internal/service/jobqueue"
)

// DmJobqueueHandler はジョブキューAPIのハンドラー
type DmJobqueueHandler struct {
	jobQueueClient *jobqueue.Client
}

// NewDmJobqueueHandler は新しいDmJobqueueHandlerを作成
func NewDmJobqueueHandler(jobQueueClient *jobqueue.Client) *DmJobqueueHandler {
	return &DmJobqueueHandler{
		jobQueueClient: jobQueueClient,
	}
}

// RegisterJobInput はジョブ登録リクエストのHuma入力
type RegisterJobInput struct {
	Body RegisterJobRequest
}

// RegisterJobRequest はジョブ登録リクエスト
type RegisterJobRequest struct {
	Message      string `json:"message,omitempty"`                      // 出力するメッセージ（オプション）
	DelaySeconds int    `json:"delay_seconds,omitempty" required:"false"` // 遅延時間（秒、オプション、0の場合はデフォルト値を使用）
	MaxRetry     int    `json:"max_retry,omitempty" required:"false"`     // 最大リトライ回数（オプション、0の場合はデフォルト値を使用）
}

// RegisterJobOutput はジョブ登録レスポンスのHuma出力
type RegisterJobOutput struct {
	Body RegisterJobResponse
}

// RegisterJobResponse はジョブ登録レスポンス
type RegisterJobResponse struct {
	JobID  string `json:"job_id"`
	Status string `json:"status"`
}

// RegisterJob はジョブを登録
func (h *DmJobqueueHandler) RegisterJob(ctx context.Context, req *RegisterJobRequest) (*RegisterJobResponse, error) {
	// Redis接続が利用できない場合のエラーハンドリング
	if h.jobQueueClient == nil {
		return nil, huma.Error503ServiceUnavailable("Job queue service is unavailable: Redis is not connected")
	}

	// メッセージの設定（デフォルト値）
	message := req.Message
	if message == "" {
		message = "Job executed successfully"
	}

	// ペイロードの作成
	payload := jobqueue.DelayPrintPayload{
		Message: message,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, huma.Error500InternalServerError("failed to marshal payload")
	}

	// ジョブオプションの作成
	jobOpts := &jobqueue.JobOptions{
		DelaySeconds: req.DelaySeconds,
		MaxRetry:     req.MaxRetry,
	}

	// ジョブをキューに登録
	info, err := h.jobQueueClient.EnqueueJob(
		ctx,
		jobqueue.JobTypeDelayPrint,
		payloadBytes,
		jobOpts,
	)
	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	return &RegisterJobResponse{
		JobID:  info.ID,
		Status: "registered",
	}, nil
}

// RegisterDmJobqueueEndpoints はHuma APIにジョブキューエンドポイントを登録
func RegisterDmJobqueueEndpoints(api huma.API, h *DmJobqueueHandler) {
	// POST /api/dm-jobqueue/register - ジョブ登録
	// 参考コードとして利用するため、将来の実装に影響しない名前を使用
	huma.Register(api, huma.Operation{
		OperationID:   "register-demo-job",
		Method:        http.MethodPost,
		Path:          "/api/dm-jobqueue/register",
		Summary:       "ジョブを登録（参考コード）",
		Description:   "**参考コード**: 将来の本実装に影響しない名前を使用",
		Tags:          []string{"jobqueue-demo"},
		DefaultStatus: http.StatusCreated,
	}, func(ctx context.Context, input *RegisterJobInput) (*RegisterJobOutput, error) {
		resp, err := h.RegisterJob(ctx, &input.Body)
		if err != nil {
			return nil, err
		}
		return &RegisterJobOutput{Body: *resp}, nil
	})
}
