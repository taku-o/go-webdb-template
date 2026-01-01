package ratelimit

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/taku-o/go-webdb-template/internal/config"
)

// TestNewRateLimitMiddleware_Disabled はレートリミットが無効な場合のテスト
func TestNewRateLimitMiddleware_Disabled(t *testing.T) {
	cfg := &config.Config{
		API: config.APIConfig{
			RateLimit: config.RateLimitConfig{
				Enabled:           false,
				RequestsPerMinute: 60,
			},
		},
	}

	middleware, err := NewRateLimitMiddleware(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, middleware)

	// リクエストが常に許可されることを確認
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := middleware(func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	err = handler(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

// TestNewRateLimitMiddleware_Enabled はレートリミットが有効な場合のテスト
func TestNewRateLimitMiddleware_Enabled(t *testing.T) {
	cfg := &config.Config{
		API: config.APIConfig{
			RateLimit: config.RateLimitConfig{
				Enabled:           true,
				RequestsPerMinute: 60,
			},
		},
		CacheServer: config.CacheServerConfig{
			Redis: config.RedisConfig{
				Default: config.RedisDefaultConfig{
					Cluster: config.RedisClusterConfig{
						Addrs: []string{}, // In-Memoryストレージを使用
					},
				},
			},
		},
	}

	middleware, err := NewRateLimitMiddleware(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, middleware)

	// リクエストが許可され、X-RateLimit-*ヘッダーが付与されることを確認
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Real-IP", "192.168.1.1")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := middleware(func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	err = handler(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	// X-RateLimit-*ヘッダーが付与されていることを確認
	assert.NotEmpty(t, rec.Header().Get("X-RateLimit-Limit"))
	assert.NotEmpty(t, rec.Header().Get("X-RateLimit-Remaining"))
	assert.NotEmpty(t, rec.Header().Get("X-RateLimit-Reset"))
}

// TestNewRateLimitMiddleware_RateLimitExceeded はレートリミット超過時のテスト
func TestNewRateLimitMiddleware_RateLimitExceeded(t *testing.T) {
	cfg := &config.Config{
		API: config.APIConfig{
			RateLimit: config.RateLimitConfig{
				Enabled:           true,
				RequestsPerMinute: 2, // 低い閾値を設定
			},
		},
		CacheServer: config.CacheServerConfig{
			Redis: config.RedisConfig{
				Default: config.RedisDefaultConfig{
					Cluster: config.RedisClusterConfig{
						Addrs: []string{}, // In-Memoryストレージを使用
					},
				},
			},
		},
	}

	middleware, err := NewRateLimitMiddleware(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, middleware)

	e := echo.New()

	// 同じIPから複数のリクエストを送信
	for i := 0; i < 3; i++ {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("X-Real-IP", "192.168.1.2")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handler := middleware(func(c echo.Context) error {
			return c.String(http.StatusOK, "OK")
		})

		err = handler(c)

		if i < 2 {
			// 最初の2回は許可される
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, rec.Code)
		} else {
			// 3回目はレートリミット超過
			assert.NoError(t, err) // ハンドラー自体はエラーを返さない
			assert.Equal(t, http.StatusTooManyRequests, rec.Code)
		}
	}
}

// TestNewRateLimitMiddleware_DifferentIPs は異なるIPアドレスが独立してカウントされることを確認
func TestNewRateLimitMiddleware_DifferentIPs(t *testing.T) {
	cfg := &config.Config{
		API: config.APIConfig{
			RateLimit: config.RateLimitConfig{
				Enabled:           true,
				RequestsPerMinute: 1, // 低い閾値を設定
			},
		},
		CacheServer: config.CacheServerConfig{
			Redis: config.RedisConfig{
				Default: config.RedisDefaultConfig{
					Cluster: config.RedisClusterConfig{
						Addrs: []string{}, // In-Memoryストレージを使用
					},
				},
			},
		},
	}

	middleware, err := NewRateLimitMiddleware(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, middleware)

	e := echo.New()

	// 異なるIPからのリクエストは独立してカウントされる
	ips := []string{"192.168.1.10", "192.168.1.11"}
	for _, ip := range ips {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("X-Real-IP", ip)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handler := middleware(func(c echo.Context) error {
			return c.String(http.StatusOK, "OK")
		})

		err = handler(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code, "IP %s should be allowed", ip)
	}
}

// TestNewRateLimitMiddleware_HourlyRateLimitEnabled は時間制限が有効な場合のテスト
func TestNewRateLimitMiddleware_HourlyRateLimitEnabled(t *testing.T) {
	cfg := &config.Config{
		API: config.APIConfig{
			RateLimit: config.RateLimitConfig{
				Enabled:           true,
				RequestsPerMinute: 60,
				RequestsPerHour:   1000,
			},
		},
		CacheServer: config.CacheServerConfig{
			Redis: config.RedisConfig{
				Default: config.RedisDefaultConfig{
					Cluster: config.RedisClusterConfig{
						Addrs: []string{}, // In-Memoryストレージを使用
					},
				},
			},
		},
	}

	middleware, err := NewRateLimitMiddleware(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, middleware)

	// リクエストが許可され、X-RateLimit-Hour-*ヘッダーが付与されることを確認
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Real-IP", "192.168.1.20")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := middleware(func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	err = handler(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	// X-RateLimit-*ヘッダー（分制限）が付与されていることを確認
	assert.NotEmpty(t, rec.Header().Get("X-RateLimit-Limit"))
	assert.NotEmpty(t, rec.Header().Get("X-RateLimit-Remaining"))
	assert.NotEmpty(t, rec.Header().Get("X-RateLimit-Reset"))

	// X-RateLimit-Hour-*ヘッダー（時間制限）が付与されていることを確認
	assert.NotEmpty(t, rec.Header().Get("X-RateLimit-Hour-Limit"))
	assert.NotEmpty(t, rec.Header().Get("X-RateLimit-Hour-Remaining"))
	assert.NotEmpty(t, rec.Header().Get("X-RateLimit-Hour-Reset"))
}

// TestNewRateLimitMiddleware_HourlyRateLimitExceeded は時間制限超過時のテスト
func TestNewRateLimitMiddleware_HourlyRateLimitExceeded(t *testing.T) {
	cfg := &config.Config{
		API: config.APIConfig{
			RateLimit: config.RateLimitConfig{
				Enabled:           true,
				RequestsPerMinute: 100, // 分制限は高く設定
				RequestsPerHour:   2,   // 時間制限は低く設定
			},
		},
		CacheServer: config.CacheServerConfig{
			Redis: config.RedisConfig{
				Default: config.RedisDefaultConfig{
					Cluster: config.RedisClusterConfig{
						Addrs: []string{}, // In-Memoryストレージを使用
					},
				},
			},
		},
	}

	middleware, err := NewRateLimitMiddleware(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, middleware)

	e := echo.New()

	// 同じIPから複数のリクエストを送信
	for i := 0; i < 3; i++ {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("X-Real-IP", "192.168.1.30")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handler := middleware(func(c echo.Context) error {
			return c.String(http.StatusOK, "OK")
		})

		err = handler(c)

		if i < 2 {
			// 最初の2回は許可される
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, rec.Code)
		} else {
			// 3回目は時間制限超過
			assert.NoError(t, err) // ハンドラー自体はエラーを返さない
			assert.Equal(t, http.StatusTooManyRequests, rec.Code)
		}
	}
}

// TestNewRateLimitMiddleware_HourlyRateLimitNotConfigured は時間制限が未設定の場合のテスト
func TestNewRateLimitMiddleware_HourlyRateLimitNotConfigured(t *testing.T) {
	cfg := &config.Config{
		API: config.APIConfig{
			RateLimit: config.RateLimitConfig{
				Enabled:           true,
				RequestsPerMinute: 60,
				RequestsPerHour:   0, // 時間制限は未設定
			},
		},
		CacheServer: config.CacheServerConfig{
			Redis: config.RedisConfig{
				Default: config.RedisDefaultConfig{
					Cluster: config.RedisClusterConfig{
						Addrs: []string{}, // In-Memoryストレージを使用
					},
				},
			},
		},
	}

	middleware, err := NewRateLimitMiddleware(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, middleware)

	// リクエストが許可され、X-RateLimit-Hour-*ヘッダーが付与されないことを確認
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Real-IP", "192.168.1.40")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := middleware(func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	err = handler(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	// X-RateLimit-*ヘッダー（分制限）が付与されていることを確認
	assert.NotEmpty(t, rec.Header().Get("X-RateLimit-Limit"))

	// X-RateLimit-Hour-*ヘッダーは付与されないことを確認
	assert.Empty(t, rec.Header().Get("X-RateLimit-Hour-Limit"))
}
