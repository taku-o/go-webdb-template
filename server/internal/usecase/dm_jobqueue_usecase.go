package usecase

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/taku-o/go-webdb-template/internal/service/jobqueue"
)

// JobQueueClientInterface はJobQueueClientのインターフェース
type JobQueueClientInterface interface {
	EnqueueJob(ctx context.Context, jobType string, payload []byte, opts *JobOptions) (*JobInfo, error)
}

// JobOptions はジョブオプション
type JobOptions struct {
	DelaySeconds int
	MaxRetry     int
}

// JobInfo はジョブ情報
type JobInfo struct {
	ID string
}

// JobQueueClientAdapter はjobqueue.Clientをラップするアダプター
type JobQueueClientAdapter struct {
	client *jobqueue.Client
}

// NewJobQueueClientAdapter は新しいJobQueueClientAdapterを作成
func NewJobQueueClientAdapter(client *jobqueue.Client) *JobQueueClientAdapter {
	if client == nil {
		return nil
	}
	return &JobQueueClientAdapter{client: client}
}

// EnqueueJob はジョブをキューに登録
func (a *JobQueueClientAdapter) EnqueueJob(ctx context.Context, jobType string, payload []byte, opts *JobOptions) (*JobInfo, error) {
	jobOpts := &jobqueue.JobOptions{
		DelaySeconds: opts.DelaySeconds,
		MaxRetry:     opts.MaxRetry,
	}

	info, err := a.client.EnqueueJob(ctx, jobType, payload, jobOpts)
	if err != nil {
		return nil, err
	}

	return &JobInfo{ID: info.ID}, nil
}

// DmJobqueueUsecase はジョブキュー関連のビジネスロジックを担当するユースケース層
type DmJobqueueUsecase struct {
	jobQueueClient JobQueueClientInterface
}

// NewDmJobqueueUsecase は新しいDmJobqueueUsecaseを作成
func NewDmJobqueueUsecase(jobQueueClient JobQueueClientInterface) *DmJobqueueUsecase {
	return &DmJobqueueUsecase{
		jobQueueClient: jobQueueClient,
	}
}

// RegisterJob はジョブを登録
func (u *DmJobqueueUsecase) RegisterJob(ctx context.Context, message string, delaySeconds int, maxRetry int) (string, error) {
	// Redis接続が利用できない場合のエラーハンドリング
	if u.jobQueueClient == nil {
		return "", errors.New("job queue service is unavailable: Redis is not connected")
	}

	// メッセージの設定（デフォルト値）
	if message == "" {
		message = "Job executed successfully"
	}

	// ペイロードの作成
	payload := jobqueue.DelayPrintPayload{
		Message: message,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	// ジョブオプションの作成
	jobOpts := &JobOptions{
		DelaySeconds: delaySeconds,
		MaxRetry:     maxRetry,
	}

	// ジョブをキューに登録
	info, err := u.jobQueueClient.EnqueueJob(
		ctx,
		jobqueue.JobTypeDelayPrint,
		payloadBytes,
		jobOpts,
	)
	if err != nil {
		return "", err
	}

	return info.ID, nil
}
