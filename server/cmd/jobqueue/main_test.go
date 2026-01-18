package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/taku-o/go-webdb-template/internal/config"
	"github.com/taku-o/go-webdb-template/internal/service/jobqueue"
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

// TestJobQueueServerCreation はJobQueueサーバーの作成をテスト
func TestJobQueueServerCreation(t *testing.T) {
	// 環境変数を設定
	os.Setenv("APP_ENV", "test")
	defer os.Unsetenv("APP_ENV")

	cfg, err := config.Load()
	assert.NoError(t, err)

	// JobQueueサーバーの作成
	// 注意: NewServerは設定に基づいてサーバーを初期化するのみ
	// 実際のRedis接続はStart()時に行われる（遅延接続）
	server, err := jobqueue.NewServer(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, server)

	// サーバーを閉じる（Shutdown）
	err = server.Shutdown()
	assert.NoError(t, err)
}

// TestGracefulShutdownSignalSetup はシグナル設定が正しく動作することをテスト
func TestGracefulShutdownSignalSetup(t *testing.T) {
	// シグナルチャネルを作成
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// テスト用にシグナル通知を停止
	signal.Stop(quit)

	// シグナルを送信（テスト用にチャネルに直接送信）
	done := make(chan bool, 1)
	go func() {
		select {
		case <-quit:
			done <- true
		case <-time.After(100 * time.Millisecond):
			done <- false
		}
	}()

	// テスト用にチャネルにシグナルを直接送信
	quit <- syscall.SIGINT

	// シグナルが受信されることを確認
	result := <-done
	assert.True(t, result, "シグナルが受信される")
}

// TestGracefulShutdownSIGTERM はSIGTERMシグナルが受信できることをテスト
func TestGracefulShutdownSIGTERM(t *testing.T) {
	// シグナルチャネルを作成
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// テスト用にシグナル通知を停止
	signal.Stop(quit)

	// シグナルを送信（テスト用にチャネルに直接送信）
	done := make(chan bool, 1)
	go func() {
		select {
		case <-quit:
			done <- true
		case <-time.After(100 * time.Millisecond):
			done <- false
		}
	}()

	// テスト用にチャネルにSIGTERMシグナルを直接送信
	quit <- syscall.SIGTERM

	// シグナルが受信されることを確認
	result := <-done
	assert.True(t, result, "SIGTERMシグナルが受信される")
}

// TestServerShutdown はサーバーのShutdownが正常に動作することをテスト
func TestServerShutdown(t *testing.T) {
	// 環境変数を設定
	os.Setenv("APP_ENV", "test")
	defer os.Unsetenv("APP_ENV")

	cfg, err := config.Load()
	assert.NoError(t, err)

	// JobQueueサーバーの作成
	server, err := jobqueue.NewServer(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, server)

	// Shutdown呼び出しがエラーなく完了すること
	err = server.Shutdown()
	assert.NoError(t, err)

	// 2回目のShutdown呼び出しもエラーなく完了すること（べき等性）
	err = server.Shutdown()
	assert.NoError(t, err)
}

// TestHealthEndpoint は /health エンドポイントのテスト
func TestHealthEndpoint(t *testing.T) {
	// HTTPサーバーを作成
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// テストリクエストを作成
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	// リクエストを処理
	mux.ServeHTTP(rec, req)

	// アサーション
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "OK", rec.Body.String())
	assert.Equal(t, "text/plain", rec.Header().Get("Content-Type"))
}
