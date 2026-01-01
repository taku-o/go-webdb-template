# ジョブキュー機能設計書

## 1. 概要

### 1.1 設計の目的
要件定義書に基づき、キューにジョブを登録し、バックグラウンドで処理を行う機能の詳細設計を定義する。Asynqライブラリを使用してRedisベースのジョブキューシステムを実装し、既存のアーキテクチャに統合する。

### 1.2 設計の範囲
- Redis環境の構築（Docker Compose設定）
- Redis Insight環境の構築（Docker Compose設定）
- 起動スクリプトの実装
- Asynqライブラリの導入と統合
- ジョブキューシステムの実装（クライアント/サーバー）
- ジョブ処理の実装（遅延時間対応）
- ジョブ登録APIの実装
- クライアント側UIの実装
- 設定管理（Redis接続設定）
- エラーハンドリング設計
- テスト戦略

### 1.3 設計方針
- **Asynqライブラリの活用**: `github.com/hibiken/asynq`を使用して堅牢なジョブキューシステムを構築
- **既存パターンの遵守**: 既存のDocker Compose設定や起動スクリプトのパターンに従う
- **参考コードとしての実装**: 将来の本実装に影響しない名前を使用
- **柔軟な遅延時間設定**: ジョブごとに遅延時間を設定可能、デフォルトは3分
- **エラーハンドリング**: Redis接続失敗時も適切にエラーを返す

## 2. アーキテクチャ設計

### 2.1 ディレクトリ構造

#### 2.1.1 変更前の構造
```
go-webdb-template/
├── docker-compose.mailpit.yml
├── docker-compose.metabase.yml
├── docker-compose.cloudbeaver.yml
├── scripts/
│   ├── start-mailpit.sh
│   ├── metabase-start.sh
│   └── cloudbeaver-start.sh
├── server/
│   ├── internal/
│   │   ├── api/
│   │   │   ├── handler/
│   │   │   │   ├── dm_user_handler.go
│   │   │   │   └── dm_post_handler.go
│   │   │   └── router/
│   │   │       └── router.go
│   │   ├── config/
│   │   │   └── config.go
│   │   └── service/
│   └── cmd/
│       └── server/
│           └── main.go
└── ...
```

#### 2.1.2 変更後の構造
```
go-webdb-template/
├── docker-compose.redis.yml              # 新規: Redis用Docker Compose設定
├── docker-compose.redis-insight.yml      # 新規: Redis Insight用Docker Compose設定
├── docker-compose.mailpit.yml
├── docker-compose.metabase.yml
├── docker-compose.cloudbeaver.yml
├── scripts/
│   ├── start-redis.sh                    # 新規: Redis起動スクリプト
│   ├── start-redis-insight.sh            # 新規: Redis Insight起動スクリプト
│   ├── start-mailpit.sh
│   ├── metabase-start.sh
│   └── cloudbeaver-start.sh
├── server/
│   ├── internal/
│   │   ├── api/
│   │   │   ├── handler/
│   │   │   │   ├── dm_user_handler.go
│   │   │   │   ├── dm_post_handler.go
│   │   │   │   └── dm_jobqueue_handler.go  # 新規: ジョブキューAPIハンドラー
│   │   │   └── router/
│   │   │       └── router.go              # 変更: ジョブキューエンドポイント登録
│   │   ├── config/
│   │   │   └── config.go                 # 変更: Redis設定追加
│   │   └── service/
│   │       └── jobqueue/                 # 新規: ジョブキューサービス
│   │           ├── client.go             # Asynqクライアント
│   │           ├── server.go             # Asynqサーバー
│   │           ├── processor.go          # ジョブ処理実装
│   │           └── constants.go          # 定数定義（デフォルト遅延時間など）
│   └── cmd/
│       └── server/
│           └── main.go                   # 変更: ジョブキューサーバー起動
└── ...
```

### 2.2 ジョブ登録・処理の実行フロー

```
┌─────────────────────────────────────────────────────────────┐
│              1. アプリケーション起動                           │
│              server/cmd/server/main.go                        │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              2. 設定ファイル読み込み                           │
│              config.Load()                                  │
│              - Redis接続設定を読み込み                        │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              3. Asynqクライアント初期化                        │
│              jobqueue.NewClient()                            │
│              - Redis接続                                     │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              4. Asynqサーバー初期化                           │
│              jobqueue.NewServer()                            │
│              - Redis接続                                     │
│              - ジョブハンドラー登録                           │
│              - バックグラウンドでジョブ処理開始                │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              5. クライアントからジョブ登録リクエスト             │
│              POST /api/dm-jobqueue/register                  │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              6. APIハンドラー処理                              │
│              DmJobqueueHandler.RegisterJob()                 │
│              - リクエストから遅延時間を取得（オプション）        │
│              - デフォルト値（3分）を使用する場合は定数参照        │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              7. ジョブをRedisキューに登録                      │
│              AsynqClient.Enqueue()                           │
│              - ジョブタイプ: "demo:delay_print"               │
│              - 遅延時間: 指定値またはデフォルト（3分）          │
│              - ペイロード: 出力する文字列                      │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              8. APIレスポンス返却                              │
│              - ジョブIDを返す                                 │
│              - 登録成功/失敗のステータス                       │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              9. Asynqサーバーがジョブを処理                    │
│              - 指定された遅延時間待機                          │
│              - ジョブハンドラー実行                            │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              10. 標準出力に文字列を出力                        │
│              fmt.Println("Job executed: ...")                │
└─────────────────────────────────────────────────────────────┘
```

### 2.3 既存アーキテクチャとの統合

```
┌─────────────────────────────────────────────────────────────┐
│              DmJobqueueHandler (internal/api/handler)       │
│              - AsynqClientを保持                            │
│              - RegisterJob()でジョブ登録                      │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              JobQueueClient (internal/service/jobqueue)     │
│              - Asynqクライアントをラップ                      │
│              - EnqueueJob()でジョブ登録                       │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              Redis (Docker Compose)                         │
│              - ジョブキューを保存                             │
│              - データ永続化（RDB/AOF）                        │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              JobQueueServer (internal/service/jobqueue)     │
│              - Asynqサーバーをラップ                          │
│              - バックグラウンドでジョブ処理                     │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              JobProcessor (internal/service/jobqueue)       │
│              - ProcessDelayPrintJob()でジョブ処理             │
│              - 遅延時間待機後、標準出力に文字列を出力           │
└─────────────────────────────────────────────────────────────┘
```

## 3. コンポーネント設計

### 3.1 Docker Compose設定

#### 3.1.1 docker-compose.redis.yml

```yaml
version: '3.8'

services:
  redis:
    image: redis:7-alpine
    container_name: redis
    ports:
      - "6379:6379"
    volumes:
      - ./redis/data/jobqueue:/data
    command: >
      redis-server
      --appendonly yes
      --appendfsync everysec
    restart: unless-stopped
    networks:
      - redis-network

networks:
  redis-network:
    name: redis-network  # 外部ネットワークとして使用するため、名前を明示的に指定
    driver: bridge
```

**設計ポイント**:
- Redis 7 Alpineイメージを使用（軽量）
- ポート6379を公開
- データ永続化: AOF（Append Only File）を有効化
- データ保存先: プロジェクトルートの`redis/data/jobqueue`ディレクトリに保存
  - ホストマシンから直接アクセス可能
  - データは`{プロジェクトルート}/redis/data/jobqueue`に保存される
- ネットワーク`redis-network`を作成（外部ネットワークとして使用可能）
- Docker Composeと起動スクリプトは開発用途（本番・staging環境では別の方法でRedisを起動）
- Redis自体は本番環境・staging環境でも使用される（1台のRedisサーバー）

#### 3.1.2 docker-compose.redis-insight.yml

```yaml
version: '3.8'

services:
  redis-insight:
    image: redis/redis-insight:latest
    container_name: redis-insight
    ports:
      - "8001:8001"
    environment:
      # docker-compose.redis.ymlで起動した1台のRedisサーバーと接続
      # 起動スクリプトでconfig/{env}/cacheserver.yamlから読み取って環境変数に設定
      - REDIS_HOSTS=${REDIS_HOSTS:-local:redis:6379}  # 起動スクリプトで設定（デフォルト値: local:redis:6379）
    volumes:
      - redis-insight-data:/data
    restart: unless-stopped
    networks:
      - redis-network  # docker-compose.redis.ymlのネットワークに接続

volumes:
  redis-insight-data:
    driver: local

networks:
  redis-network:
    external: true  # docker-compose.redis.ymlで作成されたネットワークを使用
    name: redis-network  # docker-compose.redis.ymlのネットワーク名と一致させる
```

**注意**: `docker-compose.redis.yml`で作成されたネットワーク`redis-network`を外部ネットワークとして使用します。`docker-compose.redis.yml`を先に起動してネットワークを作成する必要があります。

**設計ポイント**:
- Redis Insight最新イメージを使用
- ポート8001を公開（Web UI）
- `docker-compose.redis.yml`で起動した1台のRedisサーバーと接続
- 起動スクリプトで`config/{env}/cacheserver.yaml`から読み取って環境変数`REDIS_HOSTS`に設定
- ボリューム`redis-insight-data`で設定を永続化
- Docker Composeと起動スクリプトは開発用途（本番・staging環境では別の方法でRedis Insightを起動）
- Redis Insight自体は本番環境・staging環境でも使用される（データビューワとして使用）

### 3.2 起動スクリプト（開発用途）

#### 3.2.1 scripts/start-redis.sh

```bash
#!/bin/bash

# Redis起動スクリプト
# 使用方法: ./scripts/start-redis.sh {start|stop}

SCRIPT_DIR=$(cd "$(dirname "$0")/.." && pwd)
COMPOSE_FILE="$SCRIPT_DIR/docker-compose.redis.yml"

case "$1" in
  start)
    echo "Starting Redis..."
    docker-compose -f "$COMPOSE_FILE" up -d
    echo "Redis started. Port: 6379"
    ;;
  stop)
    echo "Stopping Redis..."
    docker-compose -f "$COMPOSE_FILE" down
    echo "Redis stopped."
    ;;
  *)
    echo "Usage: $0 {start|stop}"
    exit 1
    ;;
esac
```

**設計ポイント**:
- 既存の`start-mailpit.sh`と同じパターン
- `start`/`stop`コマンドをサポート
- 適切なフィードバックを提供
- 開発用途（本番・staging環境では別の方法でRedisを起動）

#### 3.2.2 scripts/start-redis-insight.sh

```bash
#!/bin/bash

# Redis Insight起動スクリプト
# 使用方法: ./scripts/start-redis-insight.sh {start|stop}
# docker-compose.redis.ymlで起動した1台のRedisサーバーと接続

SCRIPT_DIR=$(cd "$(dirname "$0")/.." && pwd)
COMPOSE_FILE="$SCRIPT_DIR/docker-compose.redis-insight.yml"

# 環境変数の取得（デフォルト: develop）
ENV=${APP_ENV:-develop}
CACHESERVER_CONFIG="$SCRIPT_DIR/config/$ENV/cacheserver.yaml"

# Redis接続アドレスの取得（cacheserver.yamlから）
# docker-compose.redis.ymlで起動したRedisコンテナ名（redis）を使用
REDIS_HOST="redis"  # docker-compose.redis.ymlのコンテナ名
REDIS_PORT="6379"

if [ -f "$CACHESERVER_CONFIG" ]; then
    # cacheserver.yamlからジョブキュー用Redis接続アドレスを取得（簡易的な実装）
    # より堅牢な実装にはyqなどのツールを使用することを推奨
    # redis.jobqueue.addrを取得
    ADDR=$(grep -A 1 "jobqueue:" "$CACHESERVER_CONFIG" | grep "addr:" | sed 's/.*addr: *//' | tr -d ' "')
    if [ -z "$ADDR" ]; then
        # フォールバック: 旧形式のcluster.addrsから取得
        ADDR=$(grep -A 2 "addrs:" "$CACHESERVER_CONFIG" | grep -v "addrs:" | head -1 | sed 's/.*- //' | tr -d ' "')
    fi
    if [ -n "$ADDR" ]; then
        # アドレスがlocalhost:6379の場合は、コンテナ名（redis）を使用
        if [[ "$ADDR" == "localhost:6379" ]]; then
            REDIS_HOST="redis"
        else
            # 別のアドレスの場合はそのまま使用（ホスト名:ポート形式を想定）
            REDIS_HOST=$(echo "$ADDR" | cut -d: -f1)
            REDIS_PORT=$(echo "$ADDR" | cut -d: -f2)
        fi
    fi
fi

# REDIS_HOSTS環境変数の設定（Redis Insight用）
# 形式: local:ホスト名:ポート
export REDIS_HOSTS="local:$REDIS_HOST:$REDIS_PORT"

case "$1" in
  start)
    echo "Starting Redis Insight..."
    echo "Connecting to Redis at: $REDIS_HOST:$REDIS_PORT"
    docker-compose -f "$COMPOSE_FILE" up -d
    echo "Redis Insight started. Web UI: http://localhost:8001"
    ;;
  stop)
    echo "Stopping Redis Insight..."
    docker-compose -f "$COMPOSE_FILE" down
    echo "Redis Insight stopped."
    ;;
  *)
    echo "Usage: $0 {start|stop}"
    exit 1
    ;;
esac
```

**設計ポイント**:
- 既存の起動スクリプトと同じパターン
- `docker-compose.redis.yml`で起動した1台のRedisサーバーと接続
- `config/{env}/cacheserver.yaml`からRedis接続アドレスを読み取る（オプション）
- デフォルトでは`docker-compose.redis.yml`のコンテナ名（`redis`）を使用
- Web UIのURLを表示
- 簡易的な実装（より堅牢な実装には`yq`などのツールを使用することを推奨）
- 開発用途（本番・staging環境では別の方法でRedis Insightを起動）

### 3.3 Asynqクライアント/サーバー実装

#### 3.3.1 server/internal/service/jobqueue/constants.go

```go
package jobqueue

// ジョブタイプ定数
const (
    // JobTypeDelayPrint は遅延出力ジョブのタイプ
    // 参考コードとして利用するため、将来の実装に影響しない名前を使用
    JobTypeDelayPrint = "demo:delay_print"
)

// デフォルトの遅延時間（3分 = 180秒）
const DefaultDelaySeconds = 180

// デフォルトの最大リトライ回数
const DefaultMaxRetry = 10
```

#### 3.3.2 server/internal/service/jobqueue/client.go

```go
package jobqueue

import (
    "context"
    "fmt"
    "time"

    "github.com/hibiken/asynq"
    "github.com/taku-o/go-webdb-template/internal/config"
)

// JobOptions はジョブ登録時のオプション
type JobOptions struct {
    MaxRetry     int // 最大リトライ回数（0の場合はDefaultMaxRetryを使用）
    DelaySeconds int // 遅延時間（秒、0の場合はDefaultDelaySecondsを使用）
}

// Client はAsynqクライアントをラップする構造体
type Client struct {
    client *asynq.Client
}

// NewClient は新しいJobQueueClientを作成
func NewClient(cfg *config.Config) (*Client, error) {
    // cacheserver.yamlからジョブキュー用Redis接続設定を取得
    // ジョブキュー用は単一Redis接続（1台）を使用
    redisAddr := cfg.CacheServer.Redis.JobQueue.Addr
    if redisAddr == "" {
        redisAddr = "localhost:6379" // デフォルト値
    }

    redisOpt := asynq.RedisClientOpt{
        Addr: redisAddr,
    }

    client := asynq.NewClient(redisOpt)

    return &Client{
        client: client,
    }, nil
}

// EnqueueJob はジョブをキューに登録
// optsがnilの場合はデフォルト値を使用
func (c *Client) EnqueueJob(ctx context.Context, jobType string, payload []byte, opts *JobOptions) (*asynq.TaskInfo, error) {
    task := asynq.NewTask(jobType, payload)
    
    // オプションの設定
    asynqOpts := []asynq.Option{}
    
    // 遅延時間の設定
    delaySeconds := DefaultDelaySeconds
    if opts != nil && opts.DelaySeconds > 0 {
        delaySeconds = opts.DelaySeconds
    }
    asynqOpts = append(asynqOpts, asynq.ProcessIn(time.Duration(delaySeconds)*time.Second))
    
    // 最大リトライ回数の設定
    maxRetry := DefaultMaxRetry
    if opts != nil && opts.MaxRetry > 0 {
        maxRetry = opts.MaxRetry
    }
    asynqOpts = append(asynqOpts, asynq.MaxRetry(maxRetry))

    info, err := c.client.Enqueue(task, asynqOpts...)
    if err != nil {
        return nil, fmt.Errorf("failed to enqueue job: %w", err)
    }

    return info, nil
}

// Close はクライアントをクローズ
func (c *Client) Close() error {
    return c.client.Close()
}
```

#### 3.3.3 server/internal/service/jobqueue/processor.go

```go
package jobqueue

import (
    "context"
    "fmt"
    "time"

    "github.com/hibiken/asynq"
)

// DelayPrintPayload は遅延出力ジョブのペイロード
type DelayPrintPayload struct {
    Message string `json:"message"`
}

// ProcessDelayPrintJob は遅延出力ジョブを処理
func ProcessDelayPrintJob(ctx context.Context, t *asynq.Task) error {
    // ペイロードの解析
    var payload DelayPrintPayload
    if err := t.Payload(); err != nil {
        // ペイロードがない場合はデフォルトメッセージを使用
        payload.Message = "Job executed successfully"
    } else {
        if err := json.Unmarshal(t.Payload(), &payload); err != nil {
            return fmt.Errorf("failed to unmarshal payload: %w", err)
        }
    }

    // 標準出力に文字列を出力
    fmt.Printf("[%s] %s\n", time.Now().Format("2006-01-02 15:04:05"), payload.Message)

    return nil
}
```

#### 3.3.4 server/internal/service/jobqueue/server.go

```go
package jobqueue

import (
    "context"
    "fmt"

    "github.com/hibiken/asynq"
    "github.com/taku-o/go-webdb-template/internal/config"
)

// Server はAsynqサーバーをラップする構造体
type Server struct {
    server *asynq.Server
    mux    *asynq.ServeMux
}

// NewServer は新しいJobQueueServerを作成
func NewServer(cfg *config.Config) (*Server, error) {
    // cacheserver.yamlからジョブキュー用Redis接続設定を取得
    // ジョブキュー用は単一Redis接続（1台）を使用
    redisAddr := cfg.CacheServer.Redis.JobQueue.Addr
    if redisAddr == "" {
        redisAddr = "localhost:6379" // デフォルト値
    }

    redisOpt := asynq.RedisClientOpt{
        Addr: redisAddr,
    }

    // Asynqサーバーの設定
    srv := asynq.NewServer(
        redisOpt,
        asynq.Config{
            Concurrency: 10, // 同時実行数
            Queues: map[string]int{
                "default": 10, // デフォルトキュー
            },
        },
    )

    // ジョブハンドラーの登録
    mux := asynq.NewServeMux()
    mux.HandleFunc(JobTypeDelayPrint, ProcessDelayPrintJob)

    return &Server{
        server: srv,
        mux:    mux,
    }, nil
}

// Start はサーバーを起動（バックグラウンドで実行）
func (s *Server) Start() error {
    if err := s.server.Run(s.mux); err != nil {
        return fmt.Errorf("failed to start job queue server: %w", err)
    }
    return nil
}

// Shutdown はサーバーを停止
func (s *Server) Shutdown() error {
    s.server.Shutdown()
    return nil
}
```

### 3.4 APIハンドラー実装

#### 3.4.1 server/internal/api/handler/dm_jobqueue_handler.go

```go
package handler

import (
    "context"
    "encoding/json"
    "net/http"

    "github.com/danielgtaylor/huma/v2"
    "github.com/taku-o/go-webdb-template/internal/service/jobqueue"
)

// DmJobqueueHandler はジョブキューAPIのハンドラー
type DmJobqueueHandler struct {
    jobQueueClient *jobqueue.Client
}

// NewDmJobqueueHandler は新しいDmJobqueueHandlerを作成
func NewDmJobqueueHandler(jobQueueClient *jobqueue.Client) *DmJobqueueHandler {
    return &DmJobqueueHandler{
        jobQueueClient: jobQueueClient,
    }
}

// RegisterJobRequest はジョブ登録リクエスト
type RegisterJobRequest struct {
    Message     string `json:"message"`      // 出力するメッセージ（オプション）
    DelaySeconds int   `json:"delay_seconds"` // 遅延時間（秒、オプション、0の場合はデフォルト値を使用）
    MaxRetry    int   `json:"max_retry"`     // 最大リトライ回数（オプション、0の場合はデフォルト値を使用）
}

// RegisterJobResponse はジョブ登録レスポンス
type RegisterJobResponse struct {
    JobID  string `json:"job_id"`
    Status string `json:"status"`
}

// RegisterJob はジョブを登録
func (h *DmJobqueueHandler) RegisterJob(ctx context.Context, req *RegisterJobRequest) (*RegisterJobResponse, error) {
    // Redis接続が利用できない場合のエラーハンドリング
    if h.jobQueueClient == nil {
        return nil, huma.Error503ServiceUnavailable("Job queue service is unavailable: Redis is not connected")
    }

    // メッセージの設定（デフォルト値）
    message := req.Message
    if message == "" {
        message = "Job executed successfully"
    }

    // ペイロードの作成
    payload := jobqueue.DelayPrintPayload{
        Message: message,
    }
    payloadBytes, err := json.Marshal(payload)
    if err != nil {
        return nil, huma.Error500InternalServerError("failed to marshal payload")
    }

    // ジョブオプションの作成
    jobOpts := &jobqueue.JobOptions{
        DelaySeconds: req.DelaySeconds,
        MaxRetry:     req.MaxRetry,
    }

    // ジョブをキューに登録
    info, err := h.jobQueueClient.EnqueueJob(
        ctx,
        jobqueue.JobTypeDelayPrint,
        payloadBytes,
        jobOpts,
    )
    if err != nil {
        return nil, huma.Error500InternalServerError(err.Error())
    }

    return &RegisterJobResponse{
        JobID:  info.ID,
        Status: "registered",
    }, nil
}

// RegisterDmJobqueueEndpoints はHuma APIにジョブキューエンドポイントを登録
func RegisterDmJobqueueEndpoints(api huma.API, h *DmJobqueueHandler) {
    // POST /api/dm-jobqueue/register - ジョブ登録
    // 参考コードとして利用するため、将来の実装に影響しない名前を使用
    huma.Register(api, huma.Operation{
        OperationID:   "register-demo-job",
        Method:        http.MethodPost,
        Path:          "/api/dm-jobqueue/register",
        Summary:       "ジョブを登録（参考コード）",
        Description:   "**参考コード**: 将来の本実装に影響しない名前を使用",
        Tags:          []string{"jobqueue-demo"},
        DefaultStatus: http.StatusCreated,
    }, func(ctx context.Context, input *struct {
        Body RegisterJobRequest `json:"body"`
    }) (*RegisterJobResponse, error) {
        return h.RegisterJob(ctx, &input.Body)
    })
}
```

### 3.5 設定管理

#### 3.5.1 Redis接続設定（cacheserver.yaml）

ジョブキュー用とrate limit用のRedis接続設定を分離します。`config/{env}/cacheserver.yaml`に2つのRedis設定を定義します。

**server/internal/config/config.go への変更**:

```go
// CacheServerConfig はキャッシュサーバー設定
type CacheServerConfig struct {
    Redis RedisConfig `mapstructure:"redis"`
}

// RedisConfig はRedis設定
type RedisConfig struct {
    JobQueue RedisSingleConfig  `mapstructure:"jobqueue"` // ジョブキュー用（単一接続）
    Default  RedisClusterConfig `mapstructure:"default"` // デフォルト用（複数台対応、rate limit等で使用）
}

// RedisSingleConfig は単一Redis接続設定（ジョブキュー用）
type RedisSingleConfig struct {
    Addr string `mapstructure:"addr"` // 単一Redis接続アドレス（例: "localhost:6379"）
}

// RedisClusterConfig はRedis Cluster設定（デフォルト用、rate limit等で使用）
type RedisClusterConfig struct {
    Cluster RedisClusterOptions `mapstructure:"cluster"`
}

// RedisClusterOptions はRedis Clusterオプション
type RedisClusterOptions struct {
    Addrs []string `mapstructure:"addrs"` // Redis Clusterのアドレスリスト
}
```

**設定ファイル例（config/develop/cacheserver.yaml）**:

```yaml
# Redis接続設定
redis:
  # ジョブキュー用Redis接続設定（単一接続、1台）
  jobqueue:
    addr: "localhost:6379"
  
  # デフォルト用Redis接続設定（複数台対応可能、rate limit等で使用）
  default:
    cluster:
      addrs:
        - "localhost:6379"
        # - "localhost:6380"  # 複数台の場合は追加
```

**設計ポイント**:
- ジョブキュー用は単一Redis接続（1台）を想定
- デフォルト用は複数台対応（Cluster設定、rate limit等の他の処理でも使用）
- 2つのRedis設定を分離することで、異なるRedis環境を使用可能
- 既存の`CacheServerConfig`構造体を拡張

#### 3.5.2 Rate Limitストレージタイプ設定（config.yaml）

既存のrate limit機能に、ストレージタイプ（Redis/InMemory）を明示的に指定する設定を追加します。

**server/internal/config/config.go への追加**:

```go
// RateLimitConfig はレートリミット設定
type RateLimitConfig struct {
    Enabled           bool   `mapstructure:"enabled"`
    RequestsPerMinute int    `mapstructure:"requests_per_minute"`
    RequestsPerHour   int    `mapstructure:"requests_per_hour"` // オプション
    StorageType       string `mapstructure:"storage_type"`       // 新規追加: "redis" or "memory"（デフォルト: "auto"）
}
```

**設定ファイル例（config/develop/config.yaml）**:

```yaml
api:
  # ... 既存の設定 ...
  rate_limit:
    enabled: true
    requests_per_minute: 60
    requests_per_hour: 1000
    storage_type: "auto"  # 新規追加: "auto"（自動判定）、"redis"（強制Redis）、"memory"（強制InMemory）
```

**設計ポイント**:
- `storage_type`のデフォルト値は`"auto"`（既存の動作を維持）
- `"auto"`: `cacheserver.yaml`の`redis.default.cluster.addrs`が空の場合はInMemory、そうでない場合はRedis
- `"redis"`: 強制的にRedisを使用（`cacheserver.yaml`の`redis.default`設定が必要）
- `"memory"`: 強制的にInMemoryを使用（Redis接続不要）
- デフォルト用のRedis設定は複数台対応（Cluster設定、rate limit等の他の処理でも使用）

#### 3.5.3 rate limit実装の修正（server/internal/ratelimit/middleware.go）

```go
// initStore は環境に応じたストレージを初期化
func initStore(cfg *config.Config, prefix string) (limiter.Store, error) {
    // storage_type設定を確認
    storageType := cfg.API.RateLimit.StorageType
    if storageType == "" {
        storageType = "auto" // デフォルト値
    }

    // "memory"が指定された場合はIn-Memoryストレージを使用
    if storageType == "memory" {
        return memory.NewStore(), nil
    }

    // "redis"が指定された場合はRedisを使用（設定が必要）
    if storageType == "redis" {
        // デフォルト用のRedis設定を使用
        if len(cfg.CacheServer.Redis.Default.Cluster.Addrs) == 0 {
            return nil, fmt.Errorf("redis storage type specified but no redis addresses configured for default")
        }
        rdb := redis.NewClusterClient(&redis.ClusterOptions{
            Addrs: cfg.CacheServer.Redis.Default.Cluster.Addrs,
        })
        return redisstore.NewStoreWithOptions(rdb, limiter.StoreOptions{
            Prefix: prefix,
        })
    }

    // "auto"の場合は既存の動作（cacheserver.yamlの設定に基づいて自動判定）
    // デフォルト用のRedis設定を使用
    if len(cfg.CacheServer.Redis.Default.Cluster.Addrs) == 0 {
        // In-Memoryストレージを使用
        return memory.NewStore(), nil
    }

    // Redis Clusterを使用（デフォルト用設定）
    rdb := redis.NewClusterClient(&redis.ClusterOptions{
        Addrs: cfg.CacheServer.Redis.Default.Cluster.Addrs,
    })

    // Redisストアの作成
    return redisstore.NewStoreWithOptions(rdb, limiter.StoreOptions{
        Prefix: prefix,
    })
}
```

### 3.6 main.goへの統合

#### 3.6.1 server/cmd/server/main.go への追加

```go
// Asynqクライアントの初期化
// Redisが起動していない場合でも、APIサーバーの起動は継続する
jobQueueClient, err := jobqueue.NewClient(cfg)
if err != nil {
    // Redis接続エラーを標準エラー出力に記録（起動処理は継続）
    log.Printf("WARNING: Failed to create job queue client: %v", err)
    log.Printf("WARNING: Job queue functionality will be unavailable until Redis is started")
    jobQueueClient = nil // nilを設定して、後続処理でエラーを回避
} else {
    defer jobQueueClient.Close()
}

// Asynqサーバーの初期化と起動
// Redisが起動していない場合でも、APIサーバーの起動は継続する
var jobQueueServer *jobqueue.Server
if jobQueueClient != nil {
    server, err := jobqueue.NewServer(cfg)
    if err != nil {
        // Redis接続エラーを標準エラー出力に記録（起動処理は継続）
        log.Printf("WARNING: Failed to create job queue server: %v", err)
        log.Printf("WARNING: Job queue processing will be unavailable until Redis is started")
    } else {
        jobQueueServer = server
        // バックグラウンドでジョブ処理サーバーを起動
        go func() {
            if err := jobQueueServer.Start(); err != nil {
                // ジョブ処理サーバーの起動エラーを標準エラー出力に記録
                log.Printf("ERROR: Failed to start job queue server: %v", err)
            }
        }()
    }
}

// ジョブキューハンドラーの初期化（jobQueueClientがnilの場合も許可）
jobQueueHandler := handler.NewDmJobqueueHandler(jobQueueClient)

// ルーターにジョブキューエンドポイントを登録
handler.RegisterDmJobqueueEndpoints(humaAPI, jobQueueHandler)
```

**設計ポイント**:
- Redisが起動していない場合、標準エラー出力に警告ログを出力
- エラーが発生してもAPIサーバーの起動処理は継続する（`log.Fatalf`を使用しない）
- `jobQueueClient`がnilの場合でも、ハンドラーは初期化可能（ハンドラー側でnilチェックが必要）
- ジョブ処理サーバーの起動エラーも標準エラー出力に記録するが、APIサーバーの起動は継続

## 4. データモデル

### 4.1 ジョブペイロード

```go
// DelayPrintPayload は遅延出力ジョブのペイロード
type DelayPrintPayload struct {
    Message string `json:"message"` // 出力するメッセージ
}
```

### 4.2 APIリクエスト/レスポンス

```go
// RegisterJobRequest はジョブ登録リクエスト
type RegisterJobRequest struct {
    Message      string `json:"message"`       // 出力するメッセージ（オプション）
    DelaySeconds int    `json:"delay_seconds"` // 遅延時間（秒、オプション、0の場合はデフォルト値を使用）
    MaxRetry     int    `json:"max_retry"`     // 最大リトライ回数（オプション、0の場合はデフォルト値を使用）
}

// RegisterJobResponse はジョブ登録レスポンス
type RegisterJobResponse struct {
    JobID  string `json:"job_id"`  // ジョブID
    Status string `json:"status"`  // ステータス（"registered"）
}
```

## 5. エラーハンドリング

### 5.1 Redis接続エラー

- APIサーバー起動時にRedisが起動していない場合:
  - 標準エラー出力に警告ログを出力
  - 起動処理は継続する（`log.Fatalf`を使用しない）
  - ジョブキュー機能は利用不可となる
- ジョブ登録API呼び出し時にRedis接続が利用できない場合:
  - 503 Service Unavailableエラーを返す
  - エラーメッセージ: "Job queue service is unavailable: Redis is not connected"
- Redis接続失敗時は、APIハンドラーで適切なエラーレスポンスを返す

### 5.2 ジョブ登録エラー

- ジョブ登録失敗時は、500エラーを返す
- エラーメッセージ: "Failed to enqueue job: {error}"

### 5.3 ジョブ処理エラー

- ジョブ処理失敗時は、Asynqの再試行機能を使用（デフォルト設定）
- エラーログを標準エラー出力に記録

## 6. テスト戦略

### 6.1 単体テスト

- `jobqueue.Client.EnqueueJob()`のテスト
- `jobqueue.ProcessDelayPrintJob()`のテスト
- デフォルト遅延時間のテスト

### 6.2 統合テスト

- Redis接続のテスト
- ジョブ登録から処理までのフローのテスト
- 遅延時間のテスト

### 6.3 E2Eテスト

- APIエンドポイントのテスト
- クライアント側UIのテスト

## 7. 参考コードとしての命名規則

### 7.1 エンドポイント名

- `/api/dm-jobqueue/register`: 参考コード用のエンドポイント
- 将来の本実装では別のエンドポイント名を使用

### 7.2 ジョブタイプ

- `demo:delay_print`: 参考コード用のジョブタイプ
- 将来の本実装では別のジョブタイプ名を使用

### 7.3 タグ名

- `jobqueue-demo`: 参考コード用のタグ
- 将来の本実装では別のタグ名を使用
