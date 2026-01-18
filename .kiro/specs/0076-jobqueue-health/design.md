# JobQueueサーバー死活監視エンドポイントの設計書

## 1. 概要

### 1.1 設計の目的
要件定義書に基づき、JobQueueサーバーにHTTPサーバーを追加し、`/health`エンドポイントを実装するための詳細設計を定義する。Adminサーバーと同様のシンプルな実装を行い、他のサーバー（API、Admin）と一貫性を保つ。

### 1.2 設計の範囲
- JobQueueサーバーにHTTPサーバーを追加する
- `/health`エンドポイントを実装する
- 設定ファイルに`jobqueue`セクションを追加する
- HTTPサーバーとAsynqサーバーを並行して動作させる
- Graceful shutdownを実装する
- 単体テストと統合テストを実装する

### 1.3 設計方針
- **シンプルな実装**: Adminサーバーと同様のシンプルな実装を維持する
- **認証不要**: ヘルスチェック用のため、認証ミドルウェアを通過しない
- **一貫性**: 他のサーバー（API、Admin）の実装パターンに合わせる
- **並行実行**: HTTPサーバーとAsynqサーバーをgoroutineで並行して動作させる
- **テスト**: 単体テストと統合テストを実装する

## 2. アーキテクチャ設計

### 2.1 サーバー構成

```
┌─────────────────────────────────────────────────────────────┐
│              JobQueueサーバー (Port 8082)                    │
└─────────────────────────────────────────────────────────────┘
                          │
        ┌─────────────────┴─────────────────┐
        │                                   │
        ▼                                   ▼
┌──────────────────┐            ┌──────────────────┐
│  HTTPサーバー     │            │  Asynqサーバー     │
│  (goroutine)      │            │  (goroutine)      │
│                  │            │                  │
│  - /health       │            │  - ジョブ処理      │
│  - 200 OK        │            │  - Redis接続      │
│  - "OK"          │            │  - ジョブ消化      │
└──────────────────┘            └──────────────────┘
```

### 2.2 リクエストフロー

```
┌─────────────────────────────────────────────────────────────┐
│              /health エンドポイントのリクエストフロー           │
└─────────────────────────────────────────────────────────────┘
                          │
                          ▼
        ┌─────────────────────────────────┐
        │  HTTPリクエスト: GET /health     │
        └─────────────────────────────────┘
                          │
                          ▼
        ┌─────────────────────────────────┐
        │  標準ライブラリ net/http          │
        │  - ルーティング処理               │
        │  - 認証ミドルウェアなし            │
        └─────────────────────────────────┘
                          │
                          ▼
        ┌─────────────────────────────────┐
        │  ハンドラー関数                   │
        │  - 200 OKを返す                  │
        │  - "OK"を返す                   │
        └─────────────────────────────────┘
                          │
                          ▼
        ┌─────────────────────────────────┐
        │  HTTPレスポンス                  │
        │  - Status: 200 OK               │
        │  - Body: "OK"                   │
        │  - Content-Type: text/plain      │
        └─────────────────────────────────┘
```

### 2.3 サーバー起動フロー

```
┌─────────────────────────────────────────────────────────────┐
│              JobQueueサーバーの起動フロー                      │
└─────────────────────────────────────────────────────────────┘
                          │
                          ▼
        ┌─────────────────────────────────┐
        │  1. 設定ファイルの読み込み          │
        │     - config.Load()              │
        │     - JobQueueConfigの取得        │
        └─────────────────────────────────┘
                          │
                          ▼
        ┌─────────────────────────────────┐
        │  2. Asynqサーバーの初期化          │
        │     - jobqueue.NewServer(cfg)    │
        └─────────────────────────────────┘
                          │
                          ▼
        ┌─────────────────────────────────┐
        │  3. HTTPサーバーの初期化           │
        │     - http.Serverの作成          │
        │     - /healthエンドポイント登録  │
        └─────────────────────────────────┘
                          │
        ┌─────────────────┴─────────────────┐
        │                                   │
        ▼                                   ▼
┌──────────────────┐            ┌──────────────────┐
│  4a. Asynqサーバー │            │  4b. HTTPサーバー  │
│  起動 (goroutine) │            │  起動 (goroutine) │
│  - Start()       │            │  - ListenAndServe│
└──────────────────┘            └──────────────────┘
                          │
                          ▼
        ┌─────────────────────────────────┐
        │  5. シグナル待機                  │
        │     - SIGINT, SIGTERM            │
        └─────────────────────────────────┘
                          │
                          ▼
        ┌─────────────────────────────────┐
        │  6. Graceful shutdown            │
        │     - Asynqサーバー停止            │
        │     - HTTPサーバー停止             │
        └─────────────────────────────────┘
```

## 3. 実装設計

### 3.1 設定ファイルの拡張

#### 3.1.1 設定構造体の追加

**ファイル**: `server/internal/config/config.go`

```go
// JobQueueConfig はJobQueueサーバー設定
type JobQueueConfig struct {
	Port         int           `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

// Config構造体に追加
type Config struct {
	// ... 既存のフィールド ...
	JobQueue JobQueueConfig `mapstructure:"jobqueue"`
}
```

#### 3.1.2 設定ファイルの追加

**ファイル**: `config/develop/config.yaml`

```yaml
jobqueue:
  port: 8082
  read_timeout: 30s
  write_timeout: 30s
```

**ファイル**: `config/staging/config.yaml`

```yaml
jobqueue:
  port: 8082
  read_timeout: 30s
  write_timeout: 30s
```

**ファイル**: `config/production/config.yaml.example`

```yaml
jobqueue:
  port: 8082
  read_timeout: 30s
  write_timeout: 30s
```

#### 3.1.3 デフォルト値の設定

設定ファイルに`jobqueue`セクションが存在しない場合のデフォルト値：
- `port`: 8082
- `read_timeout`: 30s
- `write_timeout`: 30s

### 3.2 HTTPサーバーの実装

#### 3.2.1 実装場所
- **ファイル**: `server/cmd/jobqueue/main.go`
- **実装位置**: Asynqサーバーの初期化後、起動前

#### 3.2.2 実装コード

```go
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/taku-o/go-webdb-template/internal/config"
	"github.com/taku-o/go-webdb-template/internal/service/jobqueue"
)

func main() {
	log.Println("Starting JobQueue server...")

	// 1. 設定ファイルの読み込み
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	log.Println("Configuration loaded successfully")

	// 2. Asynqサーバーの初期化
	jobQueueServer, err := jobqueue.NewServer(cfg)
	if err != nil {
		log.Fatalf("Failed to create job queue server: %v", err)
	}

	// 3. HTTPサーバーの初期化
	mux := http.NewServeMux()
	
	// Health check endpoint (認証不要)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.JobQueue.Port),
		Handler:      mux,
		ReadTimeout:  cfg.JobQueue.ReadTimeout,
		WriteTimeout: cfg.JobQueue.WriteTimeout,
	}

	// 4. Asynqサーバーの起動（バックグラウンド）
	go func() {
		log.Println("Starting job queue processing...")
		if err := jobQueueServer.Start(); err != nil {
			log.Printf("ERROR: Failed to start job queue server: %v", err)
		}
	}()

	// 5. HTTPサーバーの起動（バックグラウンド）
	go func() {
		log.Printf("Starting HTTP server on port %d", cfg.JobQueue.Port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server failed: %v", err)
		}
	}()

	log.Println("JobQueue server started successfully")

	// 6. Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down JobQueue server...")

	// 7. Asynqサーバーの停止
	if err := jobQueueServer.Shutdown(); err != nil {
		log.Printf("JobQueue server shutdown error: %v", err)
	}

	// 8. HTTPサーバーの停止（30秒のタイムアウト）
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Printf("HTTP server forced to shutdown: %v", err)
	}

	log.Println("JobQueue server exited")
}
```

#### 3.2.3 実装の詳細
- **HTTPサーバー**: 標準ライブラリの`net/http`を使用
- **ルーター**: `http.NewServeMux()`を使用（シンプルな実装）
- **エンドポイント**: `/health`のみ
- **認証**: 不要（認証ミドルウェアなし）
- **レスポンス**: 
  - ステータスコード: `200 OK`
  - レスポンスボディ: `"OK"`（文字列）
  - Content-Type: `text/plain`

#### 3.2.4 実装位置の決定理由
- Asynqサーバーの初期化後: 設定ファイルが読み込まれた後にHTTPサーバーを初期化
- 起動前: 両方のサーバーをgoroutineで並行して起動
- Graceful shutdown: 両方のサーバーを適切に停止

### 3.3 コード変更箇所

#### 3.3.1 変更前のコード構造

```go
func main() {
	log.Println("Starting JobQueue server...")

	// 1. 設定ファイルの読み込み
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2. Asynqサーバーの初期化
	jobQueueServer, err := jobqueue.NewServer(cfg)
	if err != nil {
		log.Fatalf("Failed to create job queue server: %v", err)
	}

	// 3. サーバー起動（バックグラウンド）
	go func() {
		log.Println("Starting job queue processing...")
		if err := jobQueueServer.Start(); err != nil {
			log.Printf("ERROR: Failed to start job queue server: %v", err)
		}
	}()

	// 4. Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down JobQueue server...")
	if err := jobQueueServer.Shutdown(); err != nil {
		log.Printf("JobQueue server shutdown error: %v", err)
	}

	log.Println("JobQueue server exited")
}
```

#### 3.3.2 変更後のコード構造

```go
func main() {
	log.Println("Starting JobQueue server...")

	// 1. 設定ファイルの読み込み
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2. Asynqサーバーの初期化
	jobQueueServer, err := jobqueue.NewServer(cfg)
	if err != nil {
		log.Fatalf("Failed to create job queue server: %v", err)
	}

	// 3. HTTPサーバーの初期化
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.JobQueue.Port),
		Handler:      mux,
		ReadTimeout:  cfg.JobQueue.ReadTimeout,
		WriteTimeout: cfg.JobQueue.WriteTimeout,
	}

	// 4. Asynqサーバーの起動（バックグラウンド）
	go func() {
		log.Println("Starting job queue processing...")
		if err := jobQueueServer.Start(); err != nil {
			log.Printf("ERROR: Failed to start job queue server: %v", err)
		}
	}()

	// 5. HTTPサーバーの起動（バックグラウンド）
	go func() {
		log.Printf("Starting HTTP server on port %d", cfg.JobQueue.Port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server failed: %v", err)
		}
	}()

	// 6. Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down JobQueue server...")

	// 7. Asynqサーバーの停止
	if err := jobQueueServer.Shutdown(); err != nil {
		log.Printf("JobQueue server shutdown error: %v", err)
	}

	// 8. HTTPサーバーの停止
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Printf("HTTP server forced to shutdown: %v", err)
	}

	log.Println("JobQueue server exited")
}
```

### 3.4 実装の注意事項

#### 3.4.1 並行実行
- HTTPサーバーとAsynqサーバーをgoroutineで並行して起動
- 両方のサーバーが独立して動作する
- 一方のサーバーがエラーで停止しても、もう一方は動作を継続

#### 3.4.2 Graceful shutdown
- SIGINT、SIGTERMシグナルを受信した場合、両方のサーバーを停止
- Asynqサーバーを先に停止し、その後HTTPサーバーを停止
- HTTPサーバーの停止には30秒のタイムアウトを設定

#### 3.4.3 エラーハンドリング
- HTTPサーバーの起動エラーは`log.Fatalf`で処理（起動失敗時は全体を停止）
- Asynqサーバーの起動エラーは`log.Printf`で処理（エラーでも起動を継続）
- Graceful shutdown時のエラーはログに記録するのみ

#### 3.4.4 パフォーマンス
- シンプルな実装のため、レスポンス時間は1ms以下を想定
- 追加のリソース消費は不要
- HTTPサーバーのオーバーヘッドは最小限

## 4. テスト設計

### 4.1 単体テスト

#### 4.1.1 テストファイル
- **ファイル**: `server/cmd/jobqueue/main_test.go`（新規作成）
- **テスト関数**: `TestHealthEndpoint`

#### 4.1.2 テスト内容
- エンドポイントが`200 OK`を返すことを確認
- レスポンスボディが`"OK"`であることを確認
- Content-Typeが`text/plain`であることを確認
- 認証なしでアクセス可能であることを確認

#### 4.1.3 テストコード例

```go
package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
```

### 4.2 統合テスト

#### 4.2.1 テストファイル
- **ファイル**: `server/test/integration/jobqueue_health_test.go`（新規作成）
- **テスト関数**: `TestJobQueueHealth_HealthCheckNoAuth`

#### 4.2.2 テスト内容
- 実際のJobQueueサーバーを起動してテスト（可能な場合）
- 認証なしで`/health`エンドポイントにアクセス可能であることを確認
- レスポンスが正しいことを確認
- HTTPサーバーとAsynqサーバーが並行して動作することを確認

#### 4.2.3 テストコード例

```go
package integration

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJobQueueHealth_HealthCheckNoAuth(t *testing.T) {
	// 注意: 実際の実装では、テスト用のサーバー起動関数を使用
	// ここでは例として記載

	// Health check should not require auth
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get("http://localhost:8082/health")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "text/plain", resp.Header.Get("Content-Type"))

	// レスポンスボディを確認
	body := make([]byte, 1024)
	n, err := resp.Body.Read(body)
	require.NoError(t, err)
	assert.Equal(t, "OK", string(body[:n]))
}
```

### 4.3 テスト実行方法

#### 4.3.1 単体テストの実行
```bash
cd server/cmd/jobqueue
go test -v -run TestHealthEndpoint
```

#### 4.3.2 統合テストの実行
```bash
cd server/test/integration
go test -v -run TestJobQueueHealth_HealthCheckNoAuth
```

## 5. 既存機能への影響

### 5.1 既存のAsynqサーバーへの影響
- **影響なし**: HTTPサーバーは独立して動作
- **既存の機能**: Asynqサーバーの機能は影響を受けない
- **ジョブ処理**: 既存のジョブ処理機能は影響を受けない

### 5.2 設定ファイルへの影響
- **新規追加**: `jobqueue`セクションを追加（既存の設定項目に影響なし）
- **既存の設定**: `server`、`admin`等の既存の設定項目は影響を受けない

### 5.3 テストへの影響
- **既存のテスト**: 影響なし（新規エンドポイントの追加のみ）
- **新規テスト**: `/health`エンドポイントのテストを追加

## 6. 実装上の注意事項

### 6.1 HTTPサーバー実装の注意事項
- **並行実行**: HTTPサーバーとAsynqサーバーをgoroutineで並行して起動
- **Graceful shutdown**: 両方のサーバーを適切に停止する
- **エラーハンドリング**: HTTPサーバーの起動エラーを適切に処理する
- **実装の簡潔性**: Adminサーバーと同様のシンプルな実装を維持する

### 6.2 設定ファイル実装の注意事項
- **設定構造**: 既存の`ServerConfig`、`AdminConfig`と同様の構造を使用
- **デフォルト値**: ポート番号のデフォルト値を適切に設定
- **環境別設定**: 各環境の設定ファイルに適切な値を設定

### 6.3 テストの注意事項
- **単体テスト**: エンドポイントが正常に動作することを確認するテストを追加する
- **統合テスト**: HTTPサーバーとAsynqサーバーが並行して動作することを確認する
- **既存テスト**: 既存のテストが全て失敗しないことを確認する

### 6.4 動作確認の注意事項
- **ローカル環境**: ローカル環境で`curl http://localhost:8082/health`が正常に動作することを確認
- **既存機能**: 既存のAsynqサーバーの機能が正常に動作することを確認
- **Graceful shutdown**: 両方のサーバーが適切に停止することを確認

## 7. 参考実装

### 7.1 Adminサーバーの実装

#### 7.1.1 エンドポイント実装
```go
// Health check endpoint (認証不要)
app.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/plain")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("OK"))
}).Methods("GET")
```

#### 7.1.2 HTTPサーバーの起動
```go
httpServer := &http.Server{
    Addr:         fmt.Sprintf(":%d", cfg.Admin.Port),
    Handler:      httpHandler,
    ReadTimeout:  cfg.Admin.ReadTimeout,
    WriteTimeout: cfg.Admin.WriteTimeout,
}

go func() {
    log.Printf("Starting admin server on port %d", cfg.Admin.Port)
    if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        log.Fatalf("Admin server failed: %v", err)
    }
}()
```

### 7.2 JobQueueサーバーでの実装方針
- Adminサーバーと同様のシンプルな実装を維持
- 標準ライブラリの`net/http`を使用（Gorilla Mux Routerは不要）
- Asynqサーバーと並行して動作させる
- テストも同様のパターンで実装

## 8. 実装チェックリスト

### 8.1 実装項目
- [ ] `server/internal/config/config.go`に`JobQueueConfig`構造体を追加
- [ ] `config/develop/config.yaml`に`jobqueue`セクションを追加
- [ ] `config/staging/config.yaml`に`jobqueue`セクションを追加
- [ ] `config/production/config.yaml.example`に`jobqueue`セクションを追加
- [ ] `server/cmd/jobqueue/main.go`にHTTPサーバーの起動処理を追加
- [ ] `server/cmd/jobqueue/main.go`に`/health`エンドポイントを追加
- [ ] HTTPサーバーとAsynqサーバーを並行して起動する実装を追加
- [ ] Graceful shutdownの実装を追加（HTTPサーバーとAsynqサーバーの両方）

### 8.2 テスト項目
- [ ] 単体テストを実装（`server/cmd/jobqueue/main_test.go`）
- [ ] 統合テストを実装（`server/test/integration/jobqueue_health_test.go`）
- [ ] 既存のテストが全て失敗しないことを確認

### 8.3 動作確認項目
- [ ] ローカル環境で`curl http://localhost:8082/health`が正常に動作する
- [ ] JobQueueサーバーが起動した時、HTTPサーバーとAsynqサーバーの両方が動作する
- [ ] 既存のAsynqサーバーの機能が正常に動作することを確認
- [ ] Graceful shutdownが正常に動作する（両方のサーバーが適切に停止する）

## 9. 参考情報

### 9.1 関連Issue
- Feature名: 0076-jobqueue-health

### 9.2 既存実装の参考
- **Adminサーバー**: `server/cmd/admin/main.go`の`/health`エンドポイント実装
- **Adminサーバーのテスト**: `server/cmd/admin/main_test.go`の`TestHealthEndpoint`
- **APIサーバー**: `server/internal/api/router/router.go`の`/health`エンドポイント実装

### 9.3 技術スタック
- **言語**: Go
- **HTTPサーバー**: 標準ライブラリの`net/http`
- **設定管理**: `github.com/spf13/viper`（既存システムと同様）
- **ジョブキュー**: `github.com/hibiken/asynq`（既存システムと同様）

### 9.4 関連ドキュメント
- `server/cmd/jobqueue/main.go`: JobQueueサーバーのメインエントリーポイント
- `server/cmd/admin/main.go`: Adminサーバーのメインエントリーポイント（参考）
- `server/internal/config/config.go`: 設定構造体の定義
- `config/develop/config.yaml`: 開発環境の設定ファイル
