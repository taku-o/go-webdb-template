package router

import (
	"fmt"
	"net/http"
	"os"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humaecho"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"github.com/taku-o/go-webdb-template/internal/api/handler"
	"github.com/taku-o/go-webdb-template/internal/auth"
	"github.com/taku-o/go-webdb-template/internal/config"
	"github.com/taku-o/go-webdb-template/internal/ratelimit"
)

// NewRouter は新しいEchoルーターを作成
func NewRouter(dmUserHandler *handler.DmUserHandler, dmPostHandler *handler.DmPostHandler, todayHandler *handler.TodayHandler, emailHandler *handler.EmailHandler, dmJobqueueHandler *handler.DmJobqueueHandler, cfg *config.Config) *echo.Echo {
	e := echo.New()

	// デバッグモードの設定（開発環境のみ）
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "develop"
	}
	if env == "develop" {
		e.Debug = true
	}

	// AUTH0_ISSUER_BASE_URLが空の場合、サーバー起動時にエラーを発生させる
	if cfg.API.Auth0IssuerBaseURL == "" {
		panic("AUTH0_ISSUER_BASE_URL is required in config")
	}

	// Recoverミドルウェア
	e.Use(middleware.Recover())

	// CORS設定
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     cfg.CORS.AllowedOrigins,
		AllowMethods:     cfg.CORS.AllowedMethods,
		AllowHeaders:     cfg.CORS.AllowedHeaders,
		ExposeHeaders:    cfg.CORS.ExposeHeaders,
		AllowCredentials: true,
	}))

	// レートリミットミドルウェア（認証ミドルウェアの前に適用）
	rateLimitMiddleware, err := ratelimit.NewRateLimitMiddleware(cfg)
	if err != nil {
		// エラー時はログに記録し、サーバー起動を継続（fail-open方式）
		logrus.WithError(err).Error("failed to create rate limit middleware")
	} else {
		e.Use(rateLimitMiddleware)
	}

	// Health check
	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	// Huma API設定
	humaConfig := huma.DefaultConfig("go-webdb-template API", "1.0.0")
	humaConfig.DocsPath = "/docs"
	humaConfig.Servers = []*huma.Server{
		{
			URL:         fmt.Sprintf("http://localhost:%d", cfg.Server.Port),
			Description: "Development server",
		},
	}

	// SecurityScheme設定の追加
	humaConfig.Components = &huma.Components{
		SecuritySchemes: map[string]*huma.SecurityScheme{
			"bearerAuth": {
				Type:         "http",
				Scheme:       "bearer",
				BearerFormat: "JWT",
			},
		},
	}

	// Huma APIインスタンスの作成（ルートレベル、認証なし）
	humaAPI := humaecho.New(e, humaConfig)

	// Humaミドルウェアとして認証を追加（/api/パスのみ）
	humaAPI.UseMiddleware(auth.NewHumaAuthMiddleware(&cfg.API, env, cfg.API.Auth0IssuerBaseURL))

	// Humaエンドポイントの登録
	handler.RegisterDmUserEndpoints(humaAPI, dmUserHandler)
	handler.RegisterDmPostEndpoints(humaAPI, dmPostHandler)
	handler.RegisterTodayEndpoints(humaAPI, todayHandler)

	// EmailHandlerが設定されている場合のみ登録
	if emailHandler != nil {
		handler.RegisterEmailEndpoints(humaAPI, emailHandler)
	}

	// DmJobqueueHandlerが設定されている場合のみ登録
	if dmJobqueueHandler != nil {
		handler.RegisterDmJobqueueEndpoints(humaAPI, dmJobqueueHandler)
	}

	return e
}


// RegisterUploadEndpoints はTUSアップロードエンドポイントを登録する
func RegisterUploadEndpoints(e *echo.Echo, h *handler.UploadHandler, cfg *config.Config) error {
	if h == nil {
		return nil
	}

	uploadCfg := h.GetConfig()
	basePath := uploadCfg.BasePath

	// TUSハンドラーをEchoにマウント
	// StripPrefixを使用してパスを調整
	tusHandler := http.StripPrefix(basePath, h.GetHandler())

	// 環境情報を取得
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "develop"
	}

	// 認証ミドルウェアを作成
	authMiddleware := auth.NewEchoAuthMiddleware(&cfg.API, env, cfg.API.Auth0IssuerBaseURL)

	// ファイル検証ミドルウェアを作成
	validationMiddleware := handler.NewUploadValidationMiddleware(uploadCfg)

	// TUSプロトコルの全メソッドをサポート（認証ミドルウェアとファイル検証ミドルウェアを適用）
	// ミドルウェアは後から追加したものが先に実行される（認証 -> 検証 -> TUSハンドラー）
	e.Any(basePath, echo.WrapHandler(tusHandler), authMiddleware, validationMiddleware)
	e.Any(basePath+"/*", echo.WrapHandler(tusHandler), authMiddleware, validationMiddleware)

	return nil
}
