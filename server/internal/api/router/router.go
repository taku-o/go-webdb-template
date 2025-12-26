package router

import (
	"fmt"
	"net/http"
	"os"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humaecho"
	"github.com/example/go-webdb-template/internal/api/handler"
	"github.com/example/go-webdb-template/internal/auth"
	"github.com/example/go-webdb-template/internal/config"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// NewRouter は新しいEchoルーターを作成
func NewRouter(userHandler *handler.UserHandler, postHandler *handler.PostHandler, cfg *config.Config) *echo.Echo {
	e := echo.New()

	// デバッグモードの設定（開発環境のみ）
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "develop"
	}
	if env == "develop" {
		e.Debug = true
	}

	// Recoverミドルウェア
	e.Use(middleware.Recover())

	// CORS設定
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     cfg.CORS.AllowedOrigins,
		AllowMethods:     cfg.CORS.AllowedMethods,
		AllowHeaders:     cfg.CORS.AllowedHeaders,
		AllowCredentials: true,
	}))

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

	// Huma APIインスタンスの作成（ルートレベル、認証なし）
	humaAPI := humaecho.New(e, humaConfig)

	// Humaミドルウェアとして認証を追加（/api/パスのみ）
	humaAPI.UseMiddleware(auth.NewHumaAuthMiddleware(&cfg.API, env))

	// Humaエンドポイントの登録
	handler.RegisterUserEndpoints(humaAPI, userHandler)
	handler.RegisterPostEndpoints(humaAPI, postHandler)

	return e
}
