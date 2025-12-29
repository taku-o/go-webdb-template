package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/taku-o/go-webdb-template/internal/api/handler"
	"github.com/taku-o/go-webdb-template/internal/api/router"
	"github.com/taku-o/go-webdb-template/internal/config"
	"github.com/taku-o/go-webdb-template/internal/db"
	"github.com/taku-o/go-webdb-template/internal/logging"
	"github.com/taku-o/go-webdb-template/internal/repository"
	"github.com/taku-o/go-webdb-template/internal/service"
)

func main() {
	// 設定ファイルの読み込み
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// GroupManagerの初期化
	groupManager, err := db.NewGroupManager(cfg)
	if err != nil {
		log.Fatalf("Failed to create group manager: %v", err)
	}
	defer groupManager.CloseAll()

	// すべてのデータベースへの接続確認
	if err := groupManager.PingAll(); err != nil {
		log.Fatalf("Failed to ping databases: %v", err)
	}
	log.Println("Successfully connected to all database groups")

	// Repository層の初期化（GORM版を使用）
	dmUserRepo := repository.NewDmUserRepositoryGORM(groupManager)
	dmPostRepo := repository.NewDmPostRepositoryGORM(groupManager)

	// Service層の初期化
	dmUserService := service.NewDmUserService(dmUserRepo)
	dmPostService := service.NewDmPostService(dmPostRepo, dmUserRepo)

	// Handler層の初期化
	dmUserHandler := handler.NewDmUserHandler(dmUserService)
	dmPostHandler := handler.NewDmPostHandler(dmPostService)
	todayHandler := handler.NewTodayHandler()

	// Echoルーターの初期化
	e := router.NewRouter(dmUserHandler, dmPostHandler, todayHandler, cfg)

	// アクセスログの初期化
	accessLogger, err := logging.NewAccessLogger("api", cfg.Logging.OutputDir)
	if err != nil {
		log.Printf("Warning: Failed to initialize access logger: %v", err)
		log.Println("Access logging will be disabled")
	} else {
		defer accessLogger.Close()
		// Echoのアクセスログミドルウェアを追加
		e.Use(logging.NewEchoAccessLogMiddleware(accessLogger))
		log.Printf("Access logging enabled: %s", cfg.Logging.OutputDir)
	}

	// Echoサーバーのタイムアウト設定
	e.Server.ReadTimeout = cfg.Server.ReadTimeout
	e.Server.WriteTimeout = cfg.Server.WriteTimeout

	// Graceful shutdown
	go func() {
		log.Printf("Starting server on port %d", cfg.Server.Port)
		if err := e.Start(fmt.Sprintf(":%d", cfg.Server.Port)); err != nil {
			log.Printf("Server stopped: %v", err)
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

	if err := e.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
