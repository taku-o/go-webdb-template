package ratelimit

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/taku-o/go-webdb-template/internal/config"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
	redisstore "github.com/ulule/limiter/v3/drivers/store/redis"
)

// NewRateLimitMiddleware はレートリミットミドルウェアを作成
func NewRateLimitMiddleware(cfg *config.Config) (echo.MiddlewareFunc, error) {
	// レートリミットが無効な場合は、常に許可するミドルウェアを返す
	if !cfg.API.RateLimit.Enabled {
		return func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				return next(c)
			}
		}, nil
	}

	// レートリミットの設定（1分あたりのリクエスト数）
	rate := limiter.Rate{
		Period: time.Minute,
		Limit:  int64(cfg.API.RateLimit.RequestsPerMinute),
	}

	// ストレージの初期化
	store, err := initStore(cfg)
	if err != nil {
		// fail-open方式: エラー時はログに記録し、リクエストを許可
		logrus.WithError(err).Error("failed to initialize rate limit store, allowing all requests")
		return func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				return next(c)
			}
		}, nil
	}

	// limiterインスタンスの作成
	instance := limiter.New(store, rate)

	// ミドルウェア関数の返却
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// IPアドレスの取得
			ip := c.RealIP()
			if ip == "" {
				// IPアドレスが取得できない場合は許可
				return next(c)
			}

			// レートリミットチェック
			context, err := instance.Get(c.Request().Context(), ip)
			if err != nil {
				// fail-open方式: エラー時はログに記録し、リクエストを許可
				logrus.WithError(err).WithField("ip", ip).Warn("rate limit check failed, allowing request")
				return next(c)
			}

			// X-RateLimit-*ヘッダーの設定
			c.Response().Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", context.Limit))
			c.Response().Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", context.Remaining))
			c.Response().Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", context.Reset))

			// レートリミット超過時
			if context.Reached {
				return c.JSON(http.StatusTooManyRequests, map[string]interface{}{
					"code":    429,
					"message": "Too Many Requests",
				})
			}

			// レートリミット内の場合は次のハンドラーを実行
			return next(c)
		}
	}, nil
}

// initStore は環境に応じたストレージを初期化
func initStore(cfg *config.Config) (limiter.Store, error) {
	// キャッシュサーバー設定からRedis Clusterのアドレスを取得
	if len(cfg.CacheServer.Redis.Cluster.Addrs) == 0 {
		// In-Memoryストレージを使用
		return memory.NewStore(), nil
	}

	// Redis Clusterを使用
	rdb := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: cfg.CacheServer.Redis.Cluster.Addrs,
	})

	// Redisストアの作成
	return redisstore.NewStoreWithOptions(rdb, limiter.StoreOptions{
		Prefix: "ratelimit",
	})
}
