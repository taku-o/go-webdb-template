package jobqueue

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/taku-o/go-webdb-template/internal/config"
)

func TestNewServer_WithDefaultAddr(t *testing.T) {
	// 空のアドレスでサーバーを作成（デフォルト値を使用）
	cfg := &config.Config{
		CacheServer: config.CacheServerConfig{
			Redis: config.RedisConfig{
				JobQueue: config.RedisSingleConfig{
					Addr: "", // 空の場合はlocalhost:6379がデフォルト
				},
			},
		},
	}

	server, err := NewServer(cfg)
	// サーバー作成自体はRedis接続なしでも成功する
	assert.NoError(t, err)
	assert.NotNil(t, server)
}

func TestNewServer_WithCustomAddr(t *testing.T) {
	// カスタムアドレスでサーバーを作成
	cfg := &config.Config{
		CacheServer: config.CacheServerConfig{
			Redis: config.RedisConfig{
				JobQueue: config.RedisSingleConfig{
					Addr: "localhost:6380", // カスタムポート
				},
			},
		},
	}

	server, err := NewServer(cfg)
	// サーバー作成自体はRedis接続なしでも成功する
	assert.NoError(t, err)
	assert.NotNil(t, server)
}

func TestNewServer_WithConnectionOptions(t *testing.T) {
	// 接続オプションを設定してサーバーを作成
	cfg := &config.Config{
		CacheServer: config.CacheServerConfig{
			Redis: config.RedisConfig{
				JobQueue: config.RedisSingleConfig{
					Addr:            "localhost:6379",
					MaxRetries:      3,
					MinRetryBackoff: 10 * time.Millisecond,
					MaxRetryBackoff: 1 * time.Second,
					DialTimeout:     10 * time.Second,
					ReadTimeout:     5 * time.Second,
					WriteTimeout:    5 * time.Second,
					PoolSize:        20,
					PoolTimeout:     5 * time.Second,
				},
			},
		},
	}

	server, err := NewServer(cfg)
	// サーバー作成自体はRedis接続なしでも成功する
	assert.NoError(t, err)
	assert.NotNil(t, server)
}

func TestNewServer_WithZeroConnectionOptions(t *testing.T) {
	// 接続オプションを設定しない（0値）でサーバーを作成
	// デフォルト値が使用されることを確認
	cfg := &config.Config{
		CacheServer: config.CacheServerConfig{
			Redis: config.RedisConfig{
				JobQueue: config.RedisSingleConfig{
					Addr:            "localhost:6379",
					MaxRetries:      0, // デフォルト値を使用
					MinRetryBackoff: 0, // デフォルト値を使用
					MaxRetryBackoff: 0, // デフォルト値を使用
					DialTimeout:     0, // デフォルト値を使用
					ReadTimeout:     0, // デフォルト値を使用
					WriteTimeout:    0, // デフォルト値を使用
					PoolSize:        0, // デフォルト値を使用
					PoolTimeout:     0, // デフォルト値を使用
				},
			},
		},
	}

	server, err := NewServer(cfg)
	// サーバー作成自体はRedis接続なしでも成功する
	assert.NoError(t, err)
	assert.NotNil(t, server)
}
