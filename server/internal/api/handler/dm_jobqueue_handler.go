package handler

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	usecaseapi "github.com/taku-o/go-webdb-template/internal/usecase/api"
)

// DmJobqueueHandler はジョブキューAPIのハンドラー
type DmJobqueueHandler struct {
	dmJobqueueUsecase *usecaseapi.DmJobqueueUsecase
}

// NewDmJobqueueHandler は新しいDmJobqueueHandlerを作成
func NewDmJobqueueHandler(dmJobqueueUsecase *usecaseapi.DmJobqueueUsecase) *DmJobqueueHandler {
	return &DmJobqueueHandler{
		dmJobqueueUsecase: dmJobqueueUsecase,
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
	// usecase層でジョブ登録を実行
	jobID, err := h.dmJobqueueUsecase.RegisterJob(ctx, req.Message, req.DelaySeconds, req.MaxRetry)
	if err != nil {
		// Redis接続が利用できない場合のエラーハンドリング
		if err.Error() == "job queue service is unavailable: Redis is not connected" {
			return nil, huma.Error503ServiceUnavailable(err.Error())
		}
		return nil, huma.Error500InternalServerError(err.Error())
	}

	return &RegisterJobResponse{
		JobID:  jobID,
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
