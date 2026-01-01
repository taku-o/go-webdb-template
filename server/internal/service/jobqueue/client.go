package jobqueue

import (
	"context"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	"github.com/taku-o/go-webdb-template/internal/config"
)

// JobOptions はジョブ登録時のオプション
type JobOptions struct {
	MaxRetry     int // 最大リトライ回数（0の場合はDefaultMaxRetryを使用）
	DelaySeconds int // 遅延時間（秒、0の場合はDefaultDelaySecondsを使用）
}

// Client はAsynqクライアントをラップする構造体
type Client struct {
	client *asynq.Client
}

// NewClient は新しいJobQueueClientを作成
func NewClient(cfg *config.Config) (*Client, error) {
	// cacheserver.yamlからジョブキュー用Redis接続設定を取得
	// ジョブキュー用は単一Redis接続（1台）を使用
	redisAddr := cfg.CacheServer.Redis.JobQueue.Addr
	if redisAddr == "" {
		redisAddr = "localhost:6379" // デフォルト値
	}

	redisOpt := asynq.RedisClientOpt{
		Addr: redisAddr,
	}

	client := asynq.NewClient(redisOpt)

	return &Client{
		client: client,
	}, nil
}

// EnqueueJob はジョブをキューに登録
// optsがnilの場合はデフォルト値を使用
func (c *Client) EnqueueJob(ctx context.Context, jobType string, payload []byte, opts *JobOptions) (*asynq.TaskInfo, error) {
	task := asynq.NewTask(jobType, payload)

	// オプションの設定
	asynqOpts := []asynq.Option{}

	// 遅延時間の設定
	delaySeconds := DefaultDelaySeconds
	if opts != nil && opts.DelaySeconds > 0 {
		delaySeconds = opts.DelaySeconds
	}
	asynqOpts = append(asynqOpts, asynq.ProcessIn(time.Duration(delaySeconds)*time.Second))

	// 最大リトライ回数の設定
	maxRetry := DefaultMaxRetry
	if opts != nil && opts.MaxRetry > 0 {
		maxRetry = opts.MaxRetry
	}
	asynqOpts = append(asynqOpts, asynq.MaxRetry(maxRetry))

	info, err := c.client.Enqueue(task, asynqOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed to enqueue job: %w", err)
	}

	return info, nil
}

// Close はクライアントをクローズ
func (c *Client) Close() error {
	return c.client.Close()
}
