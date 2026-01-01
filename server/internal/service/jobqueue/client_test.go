package jobqueue

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/taku-o/go-webdb-template/internal/config"
)

func TestNewClient_WithDefaultAddr(t *testing.T) {
	// 空のアドレスでクライアントを作成（デフォルト値を使用）
	cfg := &config.Config{
		CacheServer: config.CacheServerConfig{
			Redis: config.RedisConfig{
				JobQueue: config.RedisSingleConfig{
					Addr: "", // 空の場合はlocalhost:6379がデフォルト
				},
			},
		},
	}

	client, err := NewClient(cfg)
	// クライアント作成自体はRedis接続なしでも成功する
	assert.NoError(t, err)
	assert.NotNil(t, client)

	// クリーンアップ
	if client != nil {
		client.Close()
	}
}

func TestNewClient_WithCustomAddr(t *testing.T) {
	// カスタムアドレスでクライアントを作成
	cfg := &config.Config{
		CacheServer: config.CacheServerConfig{
			Redis: config.RedisConfig{
				JobQueue: config.RedisSingleConfig{
					Addr: "localhost:6380", // カスタムポート
				},
			},
		},
	}

	client, err := NewClient(cfg)
	// クライアント作成自体はRedis接続なしでも成功する
	assert.NoError(t, err)
	assert.NotNil(t, client)

	// クリーンアップ
	if client != nil {
		client.Close()
	}
}

func TestJobOptions_DefaultValues(t *testing.T) {
	// JobOptionsの初期値が0であることを確認
	opts := &JobOptions{}
	assert.Equal(t, 0, opts.MaxRetry)
	assert.Equal(t, 0, opts.DelaySeconds)
}

func TestJobOptions_CustomValues(t *testing.T) {
	// JobOptionsにカスタム値を設定できることを確認
	opts := &JobOptions{
		MaxRetry:     5,
		DelaySeconds: 60,
	}
	assert.Equal(t, 5, opts.MaxRetry)
	assert.Equal(t, 60, opts.DelaySeconds)
}

func TestClient_Close(t *testing.T) {
	cfg := &config.Config{
		CacheServer: config.CacheServerConfig{
			Redis: config.RedisConfig{
				JobQueue: config.RedisSingleConfig{
					Addr: "",
				},
			},
		},
	}

	client, err := NewClient(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, client)

	// Closeが正常に実行できること
	err = client.Close()
	assert.NoError(t, err)
}
