package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/taku-o/go-webdb-template/internal/config"
)

// TestConfigLoad は設定ファイルの読み込みをテスト
func TestConfigLoad(t *testing.T) {
	// 環境変数を設定
	os.Setenv("APP_ENV", "test")
	defer os.Unsetenv("APP_ENV")

	cfg, err := config.Load()
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
}

// TestConfigLoadWithDevelopEnv はdevelop環境での設定読み込みをテスト
func TestConfigLoadWithDevelopEnv(t *testing.T) {
	// 環境変数を設定
	os.Setenv("APP_ENV", "develop")
	defer os.Unsetenv("APP_ENV")

	cfg, err := config.Load()
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
}

// TestRedisConfig はRedis設定の読み込みをテスト
func TestRedisConfig(t *testing.T) {
	// 環境変数を設定
	os.Setenv("APP_ENV", "test")
	defer os.Unsetenv("APP_ENV")

	cfg, err := config.Load()
	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	// Redis設定が存在することを確認（空でもOK）
	// ジョブキュー用Redisアドレスが設定されているか確認
	// 設定が空の場合はデフォルト値が使用される
	t.Logf("JobQueue Redis Addr: %s", cfg.CacheServer.Redis.JobQueue.Addr)
}
