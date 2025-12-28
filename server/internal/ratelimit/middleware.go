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

	// 分あたりのレートリミット設定
	minuteRate := limiter.Rate{
		Period: time.Minute,
		Limit:  int64(cfg.API.RateLimit.RequestsPerMinute),
	}

	// ストレージの初期化（分制限用）
	minuteStore, err := initStore(cfg, "ratelimit")
	if err != nil {
		// fail-open方式: エラー時はログに記録し、リクエストを許可
		logrus.WithError(err).Error("failed to initialize rate limit store, allowing all requests")
		return func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				return next(c)
			}
		}, nil
	}

	// 分制限limiterインスタンスの作成
	minuteLimiter := limiter.New(minuteStore, minuteRate)

	// 時間制限limiterインスタンスの作成（設定されている場合のみ）
	var hourLimiter *limiter.Limiter
	if cfg.API.RateLimit.RequestsPerHour > 0 {
		hourRate := limiter.Rate{
			Period: time.Hour,
			Limit:  int64(cfg.API.RateLimit.RequestsPerHour),
		}

		hourStore, err := initStore(cfg, "ratelimit_hour")
		if err != nil {
			logrus.WithError(err).Error("failed to initialize hourly rate limit store, allowing all requests")
			return func(next echo.HandlerFunc) echo.HandlerFunc {
				return func(c echo.Context) error {
					return next(c)
				}
			}, nil
		}

		hourLimiter = limiter.New(hourStore, hourRate)
	}

	// ミドルウェア関数の返却
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// IPアドレスの取得
			ip := c.RealIP()
			if ip == "" {
				// IPアドレスが取得できない場合は許可
				return next(c)
			}

			// 分制限のレートリミットチェック
			minuteContext, err := minuteLimiter.Get(c.Request().Context(), ip)
			if err != nil {
				// fail-open方式: エラー時はログに記録し、リクエストを許可
				logrus.WithError(err).WithField("ip", ip).Warn("minute rate limit check failed, allowing request")
				return next(c)
			}

			// X-RateLimit-*ヘッダーの設定（分制限の情報）
			c.Response().Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", minuteContext.Limit))
			c.Response().Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", minuteContext.Remaining))
			c.Response().Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", minuteContext.Reset))

			// 分制限超過時
			if minuteContext.Reached {
				return c.JSON(http.StatusTooManyRequests, map[string]interface{}{
					"code":    429,
					"message": "Too Many Requests",
				})
			}

			// 時間制限のレートリミットチェック（設定されている場合のみ）
			if hourLimiter != nil {
				hourContext, err := hourLimiter.Get(c.Request().Context(), ip)
				if err != nil {
					// fail-open方式: エラー時はログに記録し、リクエストを許可
					logrus.WithError(err).WithField("ip", ip).Warn("hourly rate limit check failed, allowing request")
					return next(c)
				}

				// X-RateLimit-Hour-*ヘッダーの設定（時間制限の情報）
				c.Response().Header().Set("X-RateLimit-Hour-Limit", fmt.Sprintf("%d", hourContext.Limit))
				c.Response().Header().Set("X-RateLimit-Hour-Remaining", fmt.Sprintf("%d", hourContext.Remaining))
				c.Response().Header().Set("X-RateLimit-Hour-Reset", fmt.Sprintf("%d", hourContext.Reset))

				// 時間制限超過時
				if hourContext.Reached {
					return c.JSON(http.StatusTooManyRequests, map[string]interface{}{
						"code":    429,
						"message": "Too Many Requests",
					})
				}
			}

			// レートリミット内の場合は次のハンドラーを実行
			return next(c)
		}
	}, nil
}

// initStore は環境に応じたストレージを初期化
func initStore(cfg *config.Config, prefix string) (limiter.Store, error) {
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
		Prefix: prefix,
	})
}
