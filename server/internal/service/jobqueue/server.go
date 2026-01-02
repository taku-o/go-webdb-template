package jobqueue

import (
	"fmt"
	"runtime"
	"time"

	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	"github.com/taku-o/go-webdb-template/internal/config"
)

// Server はAsynqサーバーをラップする構造体
type Server struct {
	server *asynq.Server
	mux    *asynq.ServeMux
}

// NewServer は新しいJobQueueServerを作成
func NewServer(cfg *config.Config) (*Server, error) {
	// cacheserver.yamlからジョブキュー用Redis接続設定を取得
	// ジョブキュー用は単一Redis接続（1台）を使用
	redisAddr := cfg.CacheServer.Redis.JobQueue.Addr
	if redisAddr == "" {
		redisAddr = "localhost:6379" // デフォルト値
	}

	redisOpt := asynq.RedisClientOpt{
		Addr: redisAddr,
	}

	// 接続オプションの設定（設定ファイルから読み込む、未設定の場合はデフォルト値を使用）
	if cfg.CacheServer.Redis.JobQueue.DialTimeout > 0 {
		redisOpt.DialTimeout = cfg.CacheServer.Redis.JobQueue.DialTimeout
	} else {
		redisOpt.DialTimeout = 5 * time.Second // デフォルト値
	}

	if cfg.CacheServer.Redis.JobQueue.ReadTimeout > 0 {
		redisOpt.ReadTimeout = cfg.CacheServer.Redis.JobQueue.ReadTimeout
	} else {
		redisOpt.ReadTimeout = 3 * time.Second // デフォルト値
	}

	if cfg.CacheServer.Redis.JobQueue.WriteTimeout > 0 {
		redisOpt.WriteTimeout = cfg.CacheServer.Redis.JobQueue.WriteTimeout
	} else {
		redisOpt.WriteTimeout = 3 * time.Second // デフォルト値
	}

	// Asynqサーバーの設定
	srv := asynq.NewServer(
		redisOpt,
		asynq.Config{
			Concurrency: 10, // 同時実行数
			Queues: map[string]int{
				"default": 10, // デフォルトキュー
			},
		},
	)

	// リトライ設定は、asynqが内部的に使用するgo-redisクライアントのオプションを直接設定
	// MakeRedisClient()でredis.Clientを取得して設定
	if redisClient, ok := redisOpt.MakeRedisClient().(*redis.Client); ok {
		if cfg.CacheServer.Redis.JobQueue.MaxRetries > 0 {
			redisClient.Options().MaxRetries = cfg.CacheServer.Redis.JobQueue.MaxRetries
		} else {
			redisClient.Options().MaxRetries = 2 // デフォルト値
		}

		if cfg.CacheServer.Redis.JobQueue.MinRetryBackoff > 0 {
			redisClient.Options().MinRetryBackoff = cfg.CacheServer.Redis.JobQueue.MinRetryBackoff
		} else {
			redisClient.Options().MinRetryBackoff = 8 * time.Millisecond // デフォルト値
		}

		if cfg.CacheServer.Redis.JobQueue.MaxRetryBackoff > 0 {
			redisClient.Options().MaxRetryBackoff = cfg.CacheServer.Redis.JobQueue.MaxRetryBackoff
		} else {
			redisClient.Options().MaxRetryBackoff = 512 * time.Millisecond // デフォルト値
		}

		if cfg.CacheServer.Redis.JobQueue.PoolSize > 0 {
			redisClient.Options().PoolSize = cfg.CacheServer.Redis.JobQueue.PoolSize
		} else {
			redisClient.Options().PoolSize = 10 * runtime.NumCPU() // デフォルト値: CPU数×10
		}

		if cfg.CacheServer.Redis.JobQueue.PoolTimeout > 0 {
			redisClient.Options().PoolTimeout = cfg.CacheServer.Redis.JobQueue.PoolTimeout
		} else {
			redisClient.Options().PoolTimeout = 4 * time.Second // デフォルト値
		}
	}

	// ジョブハンドラーの登録
	mux := asynq.NewServeMux()
	mux.HandleFunc(JobTypeDelayPrint, ProcessDelayPrintJob)

	return &Server{
		server: srv,
		mux:    mux,
	}, nil
}

// Start はサーバーを起動（バックグラウンドで実行）
func (s *Server) Start() error {
	if err := s.server.Run(s.mux); err != nil {
		return fmt.Errorf("failed to start job queue server: %w", err)
	}
	return nil
}

// Shutdown はサーバーを停止
func (s *Server) Shutdown() error {
	s.server.Shutdown()
	return nil
}
