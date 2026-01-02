package jobqueue

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	"github.com/taku-o/go-webdb-template/internal/config"
)

// JobOptions はジョブ登録時のオプション
type JobOptions struct {
	MaxRetry     int // 最大リトライ回数（0の場合はDefaultMaxRetryを使用）
	DelaySeconds int // 遅延時間（秒、0の場合はDefaultDelaySecondsを使用）
}

// Client はAsynqクライアントをラップする構造体
type Client struct {
	client      *asynq.Client
	redisClient *redis.Client
}

// NewClient は新しいJobQueueClientを作成
func NewClient(cfg *config.Config) (*Client, error) {
	// cacheserver.yamlからジョブキュー用Redis接続設定を取得
	// ジョブキュー用は単一Redis接続（1台）を使用
	redisAddr := cfg.CacheServer.Redis.JobQueue.Addr
	if redisAddr == "" {
		redisAddr = "localhost:6379" // デフォルト値
	}

	// go-redisクライアントを直接作成し、全ての接続オプションを設定
	// asynq.NewClientFromRedisClient()を使用することで、設定が確実に反映される
	redisOpts := &redis.Options{
		Addr: redisAddr,
	}

	// 接続オプションの設定（設定ファイルから読み込む、未設定の場合はデフォルト値を使用）
	if cfg.CacheServer.Redis.JobQueue.MaxRetries > 0 {
		redisOpts.MaxRetries = cfg.CacheServer.Redis.JobQueue.MaxRetries
	} else {
		redisOpts.MaxRetries = 2 // デフォルト値
	}

	if cfg.CacheServer.Redis.JobQueue.MinRetryBackoff > 0 {
		redisOpts.MinRetryBackoff = cfg.CacheServer.Redis.JobQueue.MinRetryBackoff
	} else {
		redisOpts.MinRetryBackoff = 8 * time.Millisecond // デフォルト値
	}

	if cfg.CacheServer.Redis.JobQueue.MaxRetryBackoff > 0 {
		redisOpts.MaxRetryBackoff = cfg.CacheServer.Redis.JobQueue.MaxRetryBackoff
	} else {
		redisOpts.MaxRetryBackoff = 512 * time.Millisecond // デフォルト値
	}

	if cfg.CacheServer.Redis.JobQueue.DialTimeout > 0 {
		redisOpts.DialTimeout = cfg.CacheServer.Redis.JobQueue.DialTimeout
	} else {
		redisOpts.DialTimeout = 5 * time.Second // デフォルト値
	}

	if cfg.CacheServer.Redis.JobQueue.ReadTimeout > 0 {
		redisOpts.ReadTimeout = cfg.CacheServer.Redis.JobQueue.ReadTimeout
	} else {
		redisOpts.ReadTimeout = 3 * time.Second // デフォルト値
	}

	if cfg.CacheServer.Redis.JobQueue.WriteTimeout > 0 {
		redisOpts.WriteTimeout = cfg.CacheServer.Redis.JobQueue.WriteTimeout
	} else {
		redisOpts.WriteTimeout = 3 * time.Second // デフォルト値
	}

	if cfg.CacheServer.Redis.JobQueue.PoolSize > 0 {
		redisOpts.PoolSize = cfg.CacheServer.Redis.JobQueue.PoolSize
	} else {
		redisOpts.PoolSize = 10 * runtime.NumCPU() // デフォルト値: CPU数×10
	}

	if cfg.CacheServer.Redis.JobQueue.PoolTimeout > 0 {
		redisOpts.PoolTimeout = cfg.CacheServer.Redis.JobQueue.PoolTimeout
	} else {
		redisOpts.PoolTimeout = 4 * time.Second // デフォルト値
	}

	// go-redisクライアントを作成
	redisClient := redis.NewClient(redisOpts)

	// asynq.NewClientFromRedisClient()を使用して、設定済みのRedisクライアントを渡す
	client := asynq.NewClientFromRedisClient(redisClient)

	return &Client{
		client:      client,
		redisClient: redisClient,
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
// NewClientFromRedisClientを使用しているため、Redisクライアントを直接クローズする
func (c *Client) Close() error {
	return c.redisClient.Close()
}
