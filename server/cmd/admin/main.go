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

	_ "github.com/GoAdminGroup/go-admin/modules/db/drivers/sqlite"
	_ "github.com/GoAdminGroup/themes/adminlte"

	gorillaAdapter "github.com/GoAdminGroup/go-admin/adapter/gorilla"
	goadminContext "github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/engine"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/chartjs"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/gorilla/mux"

	"github.com/example/go-webdb-template/internal/admin"
	adminAuth "github.com/example/go-webdb-template/internal/admin/auth"
	"github.com/example/go-webdb-template/internal/admin/pages"
	"github.com/example/go-webdb-template/internal/config"
	appdb "github.com/example/go-webdb-template/internal/db"
	"github.com/example/go-webdb-template/internal/logging"
)

func main() {
	// 設定ファイルの読み込み
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// GORM DB Managerの初期化
	gormManager, err := appdb.NewGORMManager(cfg)
	if err != nil {
		log.Fatalf("Failed to create GORM manager: %v", err)
	}
	defer gormManager.CloseAll()

	// すべてのShardへの接続確認
	if err := gormManager.PingAll(); err != nil {
		log.Fatalf("Failed to ping databases: %v", err)
	}
	log.Println("Successfully connected to all database shards")

	// Gorilla Mux Router
	app := mux.NewRouter()

	// GoAdmin Engineの初期化
	eng := engine.Default()

	// チャートテンプレートを追加
	template.AddComp(chartjs.NewChart())

	// GoAdmin設定構造体の作成
	adminCfg := admin.NewConfig(cfg)
	goadminCfg := adminCfg.GetGoAdminConfig()

	// GoAdmin Engineの設定とテーブルジェネレータの登録
	if err := eng.AddConfig(goadminCfg).
		AddGenerators(admin.Generators).
		Use(app); err != nil {
		log.Fatalf("Failed to initialize GoAdmin: %v", err)
	}

	// データベース接続を取得
	conn := db.GetConnection(eng.Services)

	// 管理者ユーザーの初期化（設定ファイルの認証情報を使用）
	if cfg.Admin.Auth.Username != "" && cfg.Admin.Auth.Password != "" {
		if err := adminAuth.UpdateAdminPassword(conn, cfg.Admin.Auth.Username, cfg.Admin.Auth.Password); err != nil {
			log.Printf("Warning: Failed to update admin password: %v", err)
		}
	}

	// カスタムページの登録（Gorilla Mux用にContent関数を使用）
	app.HandleFunc("/admin", gorillaAdapter.Content(func(ctx gorillaAdapter.Context) (types.Panel, error) {
		return pages.HomePage(goadminContext.NewContext(ctx.Request), conn)
	})).Methods("GET")
	app.HandleFunc("/admin/", gorillaAdapter.Content(func(ctx gorillaAdapter.Context) (types.Panel, error) {
		return pages.HomePage(goadminContext.NewContext(ctx.Request), conn)
	})).Methods("GET")
	app.HandleFunc("/admin/user/register", gorillaAdapter.Content(func(ctx gorillaAdapter.Context) (types.Panel, error) {
		return pages.UserRegisterPage(goadminContext.NewContext(ctx.Request), conn)
	})).Methods("GET", "POST")
	app.HandleFunc("/admin/user/register/new", gorillaAdapter.Content(func(ctx gorillaAdapter.Context) (types.Panel, error) {
		return pages.UserRegisterCompletePage(goadminContext.NewContext(ctx.Request), conn)
	})).Methods("GET")
	app.HandleFunc("/admin/api-key", gorillaAdapter.Content(func(ctx gorillaAdapter.Context) (types.Panel, error) {
		return pages.APIKeyPage(goadminContext.NewContext(ctx.Request), conn)
	})).Methods("GET", "POST")

	// アクセスログの初期化（production環境以外）
	var httpHandler http.Handler = app
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "develop"
	}
	if env != "production" {
		accessLogger, err := logging.NewAccessLogger("admin", cfg.Logging.OutputDir)
		if err != nil {
			log.Printf("Warning: Failed to initialize access logger: %v", err)
			log.Println("Access logging will be disabled")
		} else {
			defer accessLogger.Close()
			// アクセスログミドルウェアを追加
			accessLogMiddleware := logging.NewAccessLogMiddleware(accessLogger)
			httpHandler = accessLogMiddleware.Middleware(app)
			log.Printf("Access logging enabled: %s", cfg.Logging.OutputDir)
		}
	} else {
		log.Println("Access logging disabled in production environment")
	}

	// HTTPサーバーの設定
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Admin.Port),
		Handler:      httpHandler,
		ReadTimeout:  cfg.Admin.ReadTimeout,
		WriteTimeout: cfg.Admin.WriteTimeout,
	}

	// Graceful shutdown
	go func() {
		log.Printf("Starting admin server on port %d", cfg.Admin.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Admin server failed: %v", err)
		}
	}()

	// シグナル待機
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down admin server...")

	// Graceful shutdown (30秒のタイムアウト)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Admin server forced to shutdown: %v", err)
	}

	log.Println("Admin server exited")
}
