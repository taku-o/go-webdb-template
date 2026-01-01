package jobqueue

import (
	"fmt"

	"github.com/hibiken/asynq"
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
