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
	"github.com/taku-o/go-webdb-template/internal/service/email"
	"github.com/taku-o/go-webdb-template/internal/service/jobqueue"
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

	// メール送信ログの初期化
	var mailLogger *logging.MailLogger
	if cfg.Logging.MailLogEnabled {
		var err error
		mailLogger, err = logging.NewMailLogger(cfg.Logging.MailLogOutputDir, true)
		if err != nil {
			log.Printf("Warning: Failed to initialize mail logger: %v", err)
			log.Println("Mail logging will be disabled")
			mailLogger = nil
		} else {
			defer mailLogger.Close()
			log.Printf("Mail logging enabled: %s", cfg.Logging.MailLogOutputDir)
		}
	}

	// EmailServiceとTemplateServiceの初期化
	emailService, err := email.NewEmailService(&cfg.Email, mailLogger)
	if err != nil {
		log.Fatalf("Failed to create email service: %v", err)
	}
	templateService := email.NewTemplateService()

	// EmailHandlerの初期化
	emailHandler := handler.NewEmailHandler(emailService, templateService)

	// Asynqクライアントの初期化
	// Redisが起動していない場合でも、APIサーバーの起動は継続する
	var jobQueueClient *jobqueue.Client
	jobQueueClient, err = jobqueue.NewClient(cfg)
	if err != nil {
		// Redis接続エラーを標準エラー出力に記録（起動処理は継続）
		log.Printf("WARNING: Failed to create job queue client: %v", err)
		log.Printf("WARNING: Job queue functionality will be unavailable until Redis is started")
		jobQueueClient = nil
	} else {
		defer jobQueueClient.Close()
	}

	// Asynqサーバーの初期化と起動
	// Redisが起動していない場合でも、APIサーバーの起動は継続する
	var jobQueueServer *jobqueue.Server
	if jobQueueClient != nil {
		jobQueueServer, err = jobqueue.NewServer(cfg)
		if err != nil {
			// Redis接続エラーを標準エラー出力に記録（起動処理は継続）
			log.Printf("WARNING: Failed to create job queue server: %v", err)
			log.Printf("WARNING: Job queue processing will be unavailable until Redis is started")
		} else {
			// バックグラウンドでジョブ処理サーバーを起動
			go func() {
				if err := jobQueueServer.Start(); err != nil {
					// ジョブ処理サーバーの起動エラーを標準エラー出力に記録
					log.Printf("ERROR: Failed to start job queue server: %v", err)
				}
			}()
		}
	}

	// ジョブキューハンドラーの初期化（jobQueueClientがnilの場合も許可）
	dmJobqueueHandler := handler.NewDmJobqueueHandler(jobQueueClient)

	// Echoルーターの初期化
	e := router.NewRouter(dmUserHandler, dmPostHandler, todayHandler, emailHandler, dmJobqueueHandler, cfg)

	// UploadHandlerの初期化（設定がある場合のみ）
	if cfg.Upload.BasePath != "" {
		uploadHandler, err := handler.NewUploadHandler(&cfg.Upload)
		if err != nil {
			log.Fatalf("Failed to create upload handler: %v", err)
		}
		// TUSアップロードエンドポイントの登録
		if err := router.RegisterUploadEndpoints(e, uploadHandler, cfg); err != nil {
			log.Fatalf("Failed to register upload endpoints: %v", err)
		}
		log.Printf("Upload endpoint enabled: %s", cfg.Upload.BasePath)
	}

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
	e.Server.IdleTimeout = cfg.Server.IdleTimeout

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

	// ジョブキューサーバーの停止
	if jobQueueServer != nil {
		log.Println("Shutting down job queue server...")
		if err := jobQueueServer.Shutdown(); err != nil {
			log.Printf("Job queue server shutdown error: %v", err)
		}
	}

	// Graceful shutdown (30秒のタイムアウト)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
