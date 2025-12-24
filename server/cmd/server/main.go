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

	"github.com/example/go-webdb-template/internal/api/handler"
	"github.com/example/go-webdb-template/internal/api/router"
	"github.com/example/go-webdb-template/internal/config"
	"github.com/example/go-webdb-template/internal/db"
	"github.com/example/go-webdb-template/internal/repository"
	"github.com/example/go-webdb-template/internal/service"
)

func main() {
	// 設定ファイルの読み込み
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// DB Managerの初期化
	dbManager, err := db.NewManager(cfg)
	if err != nil {
		log.Fatalf("Failed to create DB manager: %v", err)
	}
	defer dbManager.CloseAll()

	// すべてのShardへの接続確認
	if err := dbManager.PingAll(); err != nil {
		log.Fatalf("Failed to ping databases: %v", err)
	}
	log.Println("Successfully connected to all database shards")

	// Repository層の初期化
	userRepo := repository.NewUserRepository(dbManager)
	postRepo := repository.NewPostRepository(dbManager)

	// Service層の初期化
	userService := service.NewUserService(userRepo)
	postService := service.NewPostService(postRepo, userRepo)

	// Handler層の初期化
	userHandler := handler.NewUserHandler(userService)
	postHandler := handler.NewPostHandler(postService)

	// Routerの初期化
	r := router.NewRouter(userHandler, postHandler, cfg)

	// HTTPサーバーの設定
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// Graceful shutdown
	go func() {
		log.Printf("Starting server on port %d", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// シグナル待機
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown (30秒のタイムアウト)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
