package jobqueue

import (
	"runtime"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/taku-o/go-webdb-template/internal/config"
)

// buildRedisOptions はRedis接続オプションを構築する
// 設定ファイルから読み込んだ値を使用し、未設定の場合はデフォルト値を使用する
func buildRedisOptions(cfg *config.RedisSingleConfig, addr string) *redis.Options {
	opts := &redis.Options{
		Addr: addr,
	}

	// MaxRetries: コマンド失敗時の最大リトライ数
	if cfg.MaxRetries > 0 {
		opts.MaxRetries = cfg.MaxRetries
	} else {
		opts.MaxRetries = 2 // デフォルト値
	}

	// MinRetryBackoff: リトライ間隔（最小）
	if cfg.MinRetryBackoff > 0 {
		opts.MinRetryBackoff = cfg.MinRetryBackoff
	} else {
		opts.MinRetryBackoff = 8 * time.Millisecond // デフォルト値
	}

	// MaxRetryBackoff: リトライ間隔（最大）
	if cfg.MaxRetryBackoff > 0 {
		opts.MaxRetryBackoff = cfg.MaxRetryBackoff
	} else {
		opts.MaxRetryBackoff = 512 * time.Millisecond // デフォルト値
	}

	// DialTimeout: 接続確立のタイムアウト
	if cfg.DialTimeout > 0 {
		opts.DialTimeout = cfg.DialTimeout
	} else {
		opts.DialTimeout = 5 * time.Second // デフォルト値
	}

	// ReadTimeout: 読み取りタイムアウト
	if cfg.ReadTimeout > 0 {
		opts.ReadTimeout = cfg.ReadTimeout
	} else {
		opts.ReadTimeout = 3 * time.Second // デフォルト値
	}

	// WriteTimeout: 書き込みタイムアウト
	if cfg.WriteTimeout > 0 {
		opts.WriteTimeout = cfg.WriteTimeout
	} else {
		opts.WriteTimeout = 3 * time.Second // デフォルト値
	}

	// PoolSize: 接続プールサイズ
	if cfg.PoolSize > 0 {
		opts.PoolSize = cfg.PoolSize
	} else {
		opts.PoolSize = 10 * runtime.NumCPU() // デフォルト値: CPU数×10
	}

	// PoolTimeout: プールから接続を取り出す際の待機時間
	if cfg.PoolTimeout > 0 {
		opts.PoolTimeout = cfg.PoolTimeout
	} else {
		opts.PoolTimeout = 4 * time.Second // デフォルト値
	}

	return opts
}
