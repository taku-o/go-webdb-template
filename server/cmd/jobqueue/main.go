package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/taku-o/go-webdb-template/internal/config"
	"github.com/taku-o/go-webdb-template/internal/service/jobqueue"
)

func main() {
	log.Println("Starting JobQueue server...")

	// 1. 設定ファイルの読み込み
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	log.Println("Configuration loaded successfully")

	// 2. Asynqサーバーの初期化
	// NewServerは設定に基づいてサーバーを初期化する
	// Redis接続は遅延接続であり、Start()時にAsynqが内部で接続を確立する
	// Redisが起動していない場合でも、Asynqは自動的に再接続を試みる
	jobQueueServer, err := jobqueue.NewServer(cfg)
	if err != nil {
		// サーバー初期化エラー（設定エラー等）の場合は起動を中止
		log.Fatalf("Failed to create job queue server: %v", err)
	}

	// 3. HTTPサーバーの初期化
	mux := http.NewServeMux()

	// Health check endpoint (認証不要)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.JobQueue.Port),
		Handler:      mux,
		ReadTimeout:  cfg.JobQueue.ReadTimeout,
		WriteTimeout: cfg.JobQueue.WriteTimeout,
	}

	// 4. Asynqサーバーの起動（バックグラウンド）
	go func() {
		log.Println("Starting job queue processing...")
		if err := jobQueueServer.Start(); err != nil {
			// ジョブ処理サーバーの起動エラーを標準エラー出力に記録
			log.Printf("ERROR: Failed to start job queue server: %v", err)
		}
	}()

	// 5. HTTPサーバーの起動（バックグラウンド）
	go func() {
		log.Printf("Starting HTTP server on port %d", cfg.JobQueue.Port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server failed: %v", err)
		}
	}()

	log.Println("JobQueue server started successfully")

	// 6. Graceful shutdown
	// シグナル待機（SIGINT、SIGTERMを受信した場合、Graceful shutdownを実行）
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down JobQueue server...")

	// 7. Asynqサーバーの停止
	if err := jobQueueServer.Shutdown(); err != nil {
		log.Printf("JobQueue server shutdown error: %v", err)
	}

	// 8. HTTPサーバーの停止（30秒のタイムアウト）
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Printf("HTTP server forced to shutdown: %v", err)
	}

	log.Println("JobQueue server exited")
}
