# Redis遅延接続・自動再接続機能設計書

## 1. 概要

### 1.1 設計の目的
要件定義書に基づき、Redisへの遅延接続とダウン後の自動再接続を実現するための詳細設計を定義する。`github.com/redis/go-redis/v9`の標準機能を活用し、既存のアーキテクチャに統合する。

### 1.2 設計の範囲
- Redis接続オプションの設定追加（2種類のRedis接続に対応）
  - Cache server用（Cluster接続、`redis.NewClusterClient`）
  - Jobqueue用（単一接続、`asynq.RedisClientOpt`）
- 既存のRedis環境の確認
- 遅延接続と自動再接続の動作確認
- Redis接続エラー時のリトライ機能の確認（標準機能の活用）
- エラーハンドリング設計
- テスト戦略

### 1.3 設計方針
- **`github.com/redis/go-redis/v9`の標準機能の活用**: 遅延接続と自動再接続は`redis.NewClusterClient`と`asynq.RedisClientOpt`の標準機能を活用
- **既存パターンの遵守**: 既存のDocker Compose設定や起動スクリプトのパターンに従う
- **後方互換性の保持**: 既存のRedis接続機能を壊さない
- **設定ファイルからの読み込み**: 接続オプション設定は設定ファイルから読み込む
- **標準リトライ機能の活用**: `github.com/redis/go-redis/v9`の標準リトライ機能を活用

## 2. アーキテクチャ設計

### 2.1 ディレクトリ構造

#### 2.1.1 変更前の構造
```
go-webdb-template/
├── docker-compose.redis.yml     # 既存: Redis用Docker Compose設定（0035-jobqueueで実装済み）
├── scripts/
│   ├── start-redis.sh           # 既存: Redis起動スクリプト（0035-jobqueueで実装済み）
│   └── ...
├── server/
│   ├── internal/
│   │   ├── ratelimit/
│   │   │   └── middleware.go    # 接続オプションがAddrsのみ
│   │   ├── service/
│   │   │   └── jobqueue/
│   │   │       ├── client.go    # 接続オプションがAddrのみ
│   │   │       └── server.go    # 接続オプションがAddrのみ
│   │   └── config/
│   │       └── config.go        # RedisClusterConfig、RedisSingleConfigに接続オプションなし
│   └── ...
└── config/
    └── {env}/
        └── cacheserver.yaml     # Redis接続情報のみ
```

#### 2.1.2 変更後の構造
```
go-webdb-template/
├── docker-compose.redis.yml     # 既存: 確認のみ
├── scripts/
│   ├── start-redis.sh           # 既存: 確認のみ
│   └── ...
├── server/
│   ├── internal/
│   │   ├── ratelimit/
│   │   │   └── middleware.go    # 変更: 接続オプション追加（Cluster接続用）
│   │   ├── service/
│   │   │   └── jobqueue/
│   │   │       ├── client.go    # 変更: 接続オプション追加（単一接続用）
│   │   │       └── server.go    # 変更: 接続オプション追加（単一接続用）
│   │   └── config/
│   │       └── config.go        # 変更: RedisClusterConfig、RedisSingleConfigに接続オプション追加
│   └── ...
└── config/
    └── {env}/
        └── cacheserver.yaml     # 変更: 接続オプション設定追加
```

### 2.2 接続確立の実行フロー

#### 2.2.1 変更前のフロー（接続オプション未設定）
```
┌─────────────────────────────────────────────────────────────┐
│              1. アプリケーション起動                           │
│              server/cmd/server/main.go                        │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              2. RateLimitMiddlewareの初期化                   │
│              ratelimit.NewRateLimitMiddleware()              │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              3. Redis接続の作成（遅延接続）                     │
│              redis.NewClusterClient()                       │
│              - Addrsのみ設定                                 │
│              - 接続オプション未設定                          │
│              - 実際のTCP接続はまだ確立されていない             │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              4. サーバー起動（Redis接続不要）                  │
│              - Redis接続が利用できない場合でも起動成功        │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              5. 最初のRedisコマンド実行時                      │
│              - この時点で接続が確立される（遅延接続）        │
│              - 接続オプション未設定のため、デフォルト動作      │
└─────────────────────────────────────────────────────────────┘
```

#### 2.2.2 変更後のフロー（接続オプション設定済み）
```
┌─────────────────────────────────────────────────────────────┐
│              1. アプリケーション起動                           │
│              server/cmd/server/main.go                        │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              2. RateLimitMiddlewareの初期化                   │
│              ratelimit.NewRateLimitMiddleware()              │
│              - 設定ファイルから接続オプションを読み込み        │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              3. Redis接続の作成（遅延接続）                     │
│              redis.NewClusterClient()                       │
│              - Addrs設定                                     │
│              - 接続オプション設定済み                        │
│                * MaxRetries: 2                              │
│                * MinRetryBackoff: 8ms                        │
│                * MaxRetryBackoff: 512ms                      │
│                * DialTimeout: 5s                            │
│                * ReadTimeout: 3s                             │
│                * PoolSize: CPU数×10                         │
│                * PoolTimeout: 4s                             │
│              - 実際のTCP接続はまだ確立されていない             │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              4. サーバー起動（Redis接続不要）                  │
│              - Redis接続が利用できない場合でも起動成功        │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              5. 最初のRedisコマンド実行時                      │
│              - この時点で接続が確立される（遅延接続）        │
│              - 接続エラー時は自動リトライ（最大2回）         │
│              - リトライ間隔: 8ms～512ms（指数バックオフ）     │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              6. Redisダウン後の自動再接続                      │
│              - Redisが復旧した際に自動的に再接続            │
│              - 接続プール設定により実現                      │
└─────────────────────────────────────────────────────────────┘
```

### 2.3 既存アーキテクチャとの統合

#### 2.3.1 Cache Server用（Cluster接続）

```
┌─────────────────────────────────────────────────────────────┐
│              RateLimitMiddleware (internal/ratelimit)      │
│              - initStore()でRedis接続作成                   │
│              - 接続オプション設定済み                       │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              redis.NewClusterClient (go-redis/v9)          │
│              - 遅延接続対応（デフォルト）                    │
│              - 自動再接続対応（接続オプション設定により）      │
│              - リトライ機能（接続オプション設定により）        │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              Redis Cluster (データベース)                    │
│              - Docker Composeで起動（既存）                │
│              - ポート: 6379                                 │
└─────────────────────────────────────────────────────────────┘
```

#### 2.3.2 Jobqueue用（単一接続）

```
┌─────────────────────────────────────────────────────────────┐
│              JobQueueClient/Server (internal/service/jobqueue)│
│              - NewClient()/NewServer()でRedis接続作成       │
│              - 接続オプション設定済み                       │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              asynq.RedisClientOpt (hibiken/asynq)          │
│              - 内部的にgo-redis/v9を使用                    │
│              - 遅延接続対応（デフォルト）                    │
│              - 自動再接続対応（接続オプション設定により）      │
│              - リトライ機能（接続オプション設定により）        │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              Redis (単一接続、データベース)                    │
│              - Docker Composeで起動（既存）                │
│              - ポート: 6379                                 │
└─────────────────────────────────────────────────────────────┘
```

## 3. コンポーネント設計

### 3.1 Redis接続オプションの設定追加

#### 3.1.1 server/internal/config/config.go の修正

**変更前** (`RedisClusterConfig`構造体):
```go
// RedisClusterConfig はRedis Cluster設定
type RedisClusterConfig struct {
    Addrs []string `mapstructure:"addrs"` // Redis Clusterのアドレスリスト
}
```

**変更後**:
```go
// RedisClusterConfig はRedis Cluster設定
type RedisClusterConfig struct {
    Addrs           []string      `mapstructure:"addrs"`            // Redis Clusterのアドレスリスト
    MaxRetries      int           `mapstructure:"max_retries"`       // コマンド失敗時の最大リトライ数（デフォルト: 2）
    MinRetryBackoff time.Duration `mapstructure:"min_retry_backoff"` // リトライ間隔（最小）（デフォルト: 8ms）
    MaxRetryBackoff time.Duration `mapstructure:"max_retry_backoff"` // リトライ間隔（最大）（デフォルト: 512ms）
    DialTimeout     time.Duration `mapstructure:"dial_timeout"`      // 接続確立のタイムアウト（デフォルト: 5s）
    ReadTimeout     time.Duration `mapstructure:"read_timeout"`      // 読み取りタイムアウト（デフォルト: 3s）
    PoolSize        int           `mapstructure:"pool_size"`          // 接続プールサイズ（デフォルト: CPU数×10）
    PoolTimeout     time.Duration `mapstructure:"pool_timeout"`       // プールから接続を取り出す際の待機時間（デフォルト: 4s）
}
```

**変更前** (`RedisSingleConfig`構造体):
```go
// RedisSingleConfig は単一Redis接続設定（ジョブキュー用）
type RedisSingleConfig struct {
    Addr string `mapstructure:"addr"` // 単一Redis接続アドレス（例: "localhost:6379"）
}
```

**変更後**:
```go
// RedisSingleConfig は単一Redis接続設定（ジョブキュー用）
type RedisSingleConfig struct {
    Addr           string        `mapstructure:"addr"`            // 単一Redis接続アドレス（例: "localhost:6379"）
    MaxRetries     int           `mapstructure:"max_retries"`       // コマンド失敗時の最大リトライ数（デフォルト: 2）
    MinRetryBackoff time.Duration `mapstructure:"min_retry_backoff"` // リトライ間隔（最小）（デフォルト: 8ms）
    MaxRetryBackoff time.Duration `mapstructure:"max_retry_backoff"` // リトライ間隔（最大）（デフォルト: 512ms）
    DialTimeout     time.Duration `mapstructure:"dial_timeout"`      // 接続確立のタイムアウト（デフォルト: 5s）
    ReadTimeout     time.Duration `mapstructure:"read_timeout"`      // 読み取りタイムアウト（デフォルト: 3s）
    WriteTimeout    time.Duration `mapstructure:"write_timeout"`     // 書き込みタイムアウト（デフォルト: 3s）
    PoolSize        int           `mapstructure:"pool_size"`          // 接続プールサイズ（デフォルト: CPU数×10）
    PoolTimeout     time.Duration `mapstructure:"pool_timeout"`       // プールから接続を取り出す際の待機時間（デフォルト: 4s）
}
```

**設計ポイント**:
- 接続オプションを`RedisClusterConfig`と`RedisSingleConfig`構造体に追加
- 設定ファイルから読み込めるように`mapstructure`タグを設定
- デフォルト値は設定ファイルで指定されていない場合に使用
- `RedisSingleConfig`には`WriteTimeout`も追加（`asynq.RedisClientOpt`で使用）

#### 3.1.2 server/internal/ratelimit/middleware.go の修正（Cache Server用）

**変更前** (`initStore`関数):
```go
// Redis Clusterを使用
rdb := redis.NewClusterClient(&redis.ClusterOptions{
    Addrs: cfg.CacheServer.Redis.Default.Cluster.Addrs,
})
```

**変更後**:
```go
// Redis Clusterを使用（接続オプション設定済み）
clusterOpts := &redis.ClusterOptions{
    Addrs: cfg.CacheServer.Redis.Default.Cluster.Addrs,
}

// 接続オプションの設定（設定ファイルから読み込む、未設定の場合はデフォルト値を使用）
if cfg.CacheServer.Redis.Default.Cluster.MaxRetries > 0 {
    clusterOpts.MaxRetries = cfg.CacheServer.Redis.Default.Cluster.MaxRetries
} else {
    clusterOpts.MaxRetries = 2 // デフォルト値
}

if cfg.CacheServer.Redis.Default.Cluster.MinRetryBackoff > 0 {
    clusterOpts.MinRetryBackoff = cfg.CacheServer.Redis.Default.Cluster.MinRetryBackoff
} else {
    clusterOpts.MinRetryBackoff = 8 * time.Millisecond // デフォルト値
}

if cfg.CacheServer.Redis.Default.Cluster.MaxRetryBackoff > 0 {
    clusterOpts.MaxRetryBackoff = cfg.CacheServer.Redis.Default.Cluster.MaxRetryBackoff
} else {
    clusterOpts.MaxRetryBackoff = 512 * time.Millisecond // デフォルト値
}

if cfg.CacheServer.Redis.Default.Cluster.DialTimeout > 0 {
    clusterOpts.DialTimeout = cfg.CacheServer.Redis.Default.Cluster.DialTimeout
} else {
    clusterOpts.DialTimeout = 5 * time.Second // デフォルト値
}

if cfg.CacheServer.Redis.Default.Cluster.ReadTimeout > 0 {
    clusterOpts.ReadTimeout = cfg.CacheServer.Redis.Default.Cluster.ReadTimeout
} else {
    clusterOpts.ReadTimeout = 3 * time.Second // デフォルト値
}

if cfg.CacheServer.Redis.Default.Cluster.PoolSize > 0 {
    clusterOpts.PoolSize = cfg.CacheServer.Redis.Default.Cluster.PoolSize
} else {
    clusterOpts.PoolSize = 10 * runtime.NumCPU() // デフォルト値: CPU数×10
}

if cfg.CacheServer.Redis.Default.Cluster.PoolTimeout > 0 {
    clusterOpts.PoolTimeout = cfg.CacheServer.Redis.Default.Cluster.PoolTimeout
} else {
    clusterOpts.PoolTimeout = 4 * time.Second // デフォルト値
}

rdb := redis.NewClusterClient(clusterOpts)
```

**設計ポイント**:
- 設定ファイルから接続オプションを読み込む
- 設定値が0以下の場合はデフォルト値を使用
- `PoolSize`のデフォルト値は`runtime.NumCPU() * 10`（CPU数×10）
- 既存のコードスタイルに従う

**注意**: `initStore`関数内で2箇所（144行目付近と163行目付近）で`redis.NewClusterClient`が呼び出されているため、両方の箇所で同じ修正を適用する。

#### 3.1.3 server/internal/service/jobqueue/client.go の修正（Jobqueue用）

**変更前** (`NewClient`関数):
```go
redisOpt := asynq.RedisClientOpt{
    Addr: redisAddr,
}

client := asynq.NewClient(redisOpt)
```

**変更後**:
```go
redisOpt := asynq.RedisClientOpt{
    Addr: redisAddr,
}

// 接続オプションの設定（設定ファイルから読み込む、未設定の場合はデフォルト値を使用）
if cfg.CacheServer.Redis.JobQueue.DialTimeout > 0 {
    redisOpt.DialTimeout = cfg.CacheServer.Redis.JobQueue.DialTimeout
} else {
    redisOpt.DialTimeout = 5 * time.Second // デフォルト値
}

if cfg.CacheServer.Redis.JobQueue.ReadTimeout > 0 {
    redisOpt.ReadTimeout = cfg.CacheServer.Redis.JobQueue.ReadTimeout
} else {
    redisOpt.ReadTimeout = 3 * time.Second // デフォルト値
}

if cfg.CacheServer.Redis.JobQueue.WriteTimeout > 0 {
    redisOpt.WriteTimeout = cfg.CacheServer.Redis.JobQueue.WriteTimeout
} else {
    redisOpt.WriteTimeout = 3 * time.Second // デフォルト値
}

client := asynq.NewClient(redisOpt)

// リトライ設定は、asynqが内部的に使用するgo-redisクライアントのオプションを直接設定
// MakeRedisClient()でredis.Clientを取得して設定
if redisClient, ok := redisOpt.MakeRedisClient().(*redis.Client); ok {
    if cfg.CacheServer.Redis.JobQueue.MaxRetries > 0 {
        redisClient.Options().MaxRetries = cfg.CacheServer.Redis.JobQueue.MaxRetries
    } else {
        redisClient.Options().MaxRetries = 2 // デフォルト値
    }

    if cfg.CacheServer.Redis.JobQueue.MinRetryBackoff > 0 {
        redisClient.Options().MinRetryBackoff = cfg.CacheServer.Redis.JobQueue.MinRetryBackoff
    } else {
        redisClient.Options().MinRetryBackoff = 8 * time.Millisecond // デフォルト値
    }

    if cfg.CacheServer.Redis.JobQueue.MaxRetryBackoff > 0 {
        redisClient.Options().MaxRetryBackoff = cfg.CacheServer.Redis.JobQueue.MaxRetryBackoff
    } else {
        redisClient.Options().MaxRetryBackoff = 512 * time.Millisecond // デフォルト値
    }

    if cfg.CacheServer.Redis.JobQueue.PoolSize > 0 {
        redisClient.Options().PoolSize = cfg.CacheServer.Redis.JobQueue.PoolSize
    } else {
        redisClient.Options().PoolSize = 10 * runtime.NumCPU() // デフォルト値: CPU数×10
    }

    if cfg.CacheServer.Redis.JobQueue.PoolTimeout > 0 {
        redisClient.Options().PoolTimeout = cfg.CacheServer.Redis.JobQueue.PoolTimeout
    } else {
        redisClient.Options().PoolTimeout = 4 * time.Second // デフォルト値
    }
}
```

**設計ポイント**:
- `asynq.RedisClientOpt`で直接設定できるオプション（`DialTimeout`, `ReadTimeout`, `WriteTimeout`）を設定
- リトライ設定は`MakeRedisClient()`で取得した`redis.Client`のオプションを直接設定
- 設定値が0以下の場合はデフォルト値を使用
- `PoolSize`のデフォルト値は`runtime.NumCPU() * 10`（CPU数×10）

#### 3.1.4 server/internal/service/jobqueue/server.go の修正（Jobqueue用）

`server.go`の`NewServer`関数でも同様の修正を適用する（`client.go`と同じロジック）。

**設計ポイント**:
- `client.go`と同じ接続オプション設定ロジックを適用
- コードの重複を避けるため、共通のヘルパー関数を作成することも検討可能

### 3.2 設定ファイルの修正

#### 3.2.1 config/{env}/cacheserver.yaml の修正

**変更前**:
```yaml
redis:
  jobqueue:
    addr: "localhost:6379"
  default:
    cluster:
      addrs:
        - "localhost:6379"
```

**変更後**:
```yaml
redis:
  # ジョブキュー用Redis接続設定（単一接続、1台）
  jobqueue:
    addr: "localhost:6379"
    # 接続オプション（オプション、未設定の場合はデフォルト値を使用）
    max_retries: 2                      # コマンド失敗時の最大リトライ数（デフォルト: 2）
    min_retry_backoff: 8ms              # リトライ間隔（最小）（デフォルト: 8ms）
    max_retry_backoff: 512ms            # リトライ間隔（最大）（デフォルト: 512ms）
    dial_timeout: 5s                    # 接続確立のタイムアウト（デフォルト: 5s）
    read_timeout: 3s                    # 読み取りタイムアウト（デフォルト: 3s）
    write_timeout: 3s                   # 書き込みタイムアウト（デフォルト: 3s）
    pool_size: 10                       # 接続プールサイズ（デフォルト: CPU数×10）
    pool_timeout: 4s                    # プールから接続を取り出す際の待機時間（デフォルト: 4s）
  
  # デフォルト用Redis接続設定（複数台対応可能、rate limit等で使用）
  default:
    cluster:
      addrs:
        - "localhost:6379"
      # 接続オプション（オプション、未設定の場合はデフォルト値を使用）
      max_retries: 2                      # コマンド失敗時の最大リトライ数（デフォルト: 2）
      min_retry_backoff: 8ms              # リトライ間隔（最小）（デフォルト: 8ms）
      max_retry_backoff: 512ms            # リトライ間隔（最大）（デフォルト: 512ms）
      dial_timeout: 5s                    # 接続確立のタイムアウト（デフォルト: 5s）
      read_timeout: 3s                    # 読み取りタイムアウト（デフォルト: 3s）
      pool_size: 10                       # 接続プールサイズ（デフォルト: CPU数×10）
      pool_timeout: 4s                    # プールから接続を取り出す際の待機時間（デフォルト: 4s）
```

**設計ポイント**:
- 接続オプションはすべてオプション（未設定の場合はデフォルト値を使用）
- 既存の設定ファイルとの互換性を保つ（`addr`や`addrs`のみでも動作する）
- 設定値の単位は`time.Duration`形式（`8ms`, `5s`など）
- jobqueue用とdefault用で別々に設定可能

### 3.3 Redis環境の確認

#### 3.3.1 docker-compose.redis.yml の確認

既存の`docker-compose.redis.yml`ファイル（0035-jobqueueで実装済み）を確認する。

**確認項目**:
- Redisコンテナの定義が存在するか
- ポート6379が公開されているか
- データ永続化の設定が存在するか
- ネットワーク設定が適切か

#### 3.3.2 scripts/start-redis.sh の確認

既存の`scripts/start-redis.sh`ファイル（0035-jobqueueで実装済み）を確認する。

**確認項目**:
- Docker Composeを使用してRedisを起動するか
- 既存の起動スクリプト（`start-mailpit.sh`など）と同じパターンか
- `start`/`stop`コマンドをサポートしているか

#### 3.3.3 config/{env}/cacheserver.yaml の確認

既存の設定ファイルを確認する。

**確認項目**:
- `redis.jobqueue.addr`が設定されているか
- `redis.default.cluster.addrs`が設定されているか
- 既存の設定が正しく読み込まれているか

**注意**: develop環境では、`redis.default.cluster.addrs`が空配列（`[]`）の場合、`storage_type: "auto"`の設定ではMemoryストレージが使用される。動作確認のためには、以下のいずれかの対応が必要：
- `config/develop/config.yaml`の`api.rate_limit.storage_type`を`"redis"`に変更
- `config/develop/cacheserver.yaml`の`redis.default.cluster.addrs`に`["localhost:6379"]`を設定

#### 3.3.4 config/{env}/config.yaml の確認

Rate Limitのストレージタイプ設定を確認する。

**確認項目**:
- `api.rate_limit.storage_type`が設定されているか
- `storage_type: "auto"`の場合、`cacheserver.yaml`の`redis.default.cluster.addrs`が設定されているか確認
- 動作確認のため、必要に応じて`storage_type: "redis"`に変更

**動作確認前の設定変更手順**:
1. `config/develop/config.yaml`を確認
2. `api.rate_limit.storage_type`が`"auto"`の場合、`cacheserver.yaml`の`redis.default.cluster.addrs`が空でないことを確認
3. 空の場合は、`cacheserver.yaml`の`redis.default.cluster.addrs`に`["localhost:6379"]`を設定するか、`config.yaml`の`storage_type`を`"redis"`に変更

### 3.4 遅延接続と自動再接続の動作確認

#### 3.4.1 遅延接続の確認（Cache Server用）

**確認手順**:
1. **事前準備**: `config/develop/config.yaml`の`api.rate_limit.storage_type`を`"redis"`に設定、または`config/develop/cacheserver.yaml`の`redis.default.cluster.addrs`に`["localhost:6379"]`を設定
2. Redisを停止した状態でAPIサーバーを起動
3. サーバーが正常に起動することを確認（Redis接続なしで起動可能）
4. 最初のAPIリクエストを送信（レートリミットチェックが実行される）
5. ログで接続確立のタイミングを確認
6. Redisが接続されたことを確認

**期待される動作**:
- サーバー起動時にはRedis接続が確立されていない
- 最初のRedisコマンド実行時に接続が確立される
- ログで接続確立のタイミングが確認できる

#### 3.4.2 自動再接続の確認（Cache Server用）

**確認手順**:
1. **事前準備**: `config/develop/config.yaml`の`api.rate_limit.storage_type`を`"redis"`に設定、または`config/develop/cacheserver.yaml`の`redis.default.cluster.addrs`に`["localhost:6379"]`を設定
2. Redisを起動した状態でAPIサーバーを起動
3. 正常に動作することを確認
4. Redisを停止
5. APIリクエストを送信（エラーが発生することを確認）
6. Redisを再起動
7. 次のAPIリクエストを送信（自動的に再接続されることを確認）
8. ログで再接続のタイミングを確認

#### 3.4.3 遅延接続の確認（Jobqueue用）

**確認手順**:
1. Redisを停止した状態でAPIサーバーを起動
2. サーバーが正常に起動することを確認（Redis接続なしで起動可能）
3. Jobqueueクライアントを作成（接続はまだ確立されていない）
4. ジョブを登録（最初のRedisコマンド実行時）
5. ログで接続確立のタイミングを確認
6. Redisが接続されたことを確認

#### 3.4.4 自動再接続の確認（Jobqueue用）

**確認手順**:
1. Redisを起動した状態でAPIサーバーを起動
2. 正常に動作することを確認
3. Redisを停止
4. ジョブ登録を試行（エラーが発生することを確認）
5. Redisを再起動
6. 次のジョブ登録を試行（自動的に再接続されることを確認）
7. ログで再接続のタイミングを確認

**期待される動作**:
- Redis停止時はエラーが発生する
- Redis再起動後、次のコマンド実行時に自動的に再接続される
- ログで再接続のタイミングが確認できる

### 3.5 Redis接続エラー時のリトライ機能の確認

#### 3.5.1 リトライ機能の確認（Cache Server用）

`github.com/redis/go-redis/v9`の標準リトライ機能が動作することを確認する。

**確認手順**:
1. **事前準備**: `config/develop/config.yaml`の`api.rate_limit.storage_type`を`"redis"`に設定、または`config/develop/cacheserver.yaml`の`redis.default.cluster.addrs`に`["localhost:6379"]`を設定
2. Redisを停止した状態でAPIサーバーを起動
3. APIリクエストを送信（レートリミットチェックが実行される）
4. ログでリトライ処理を確認
5. リトライ間隔が`MinRetryBackoff`と`MaxRetryBackoff`の範囲内であることを確認

#### 3.5.2 リトライ機能の確認（Jobqueue用）

`github.com/redis/go-redis/v9`の標準リトライ機能が動作することを確認する。

**確認手順**:
1. Redisを停止した状態でAPIサーバーを起動
2. ジョブ登録を試行（最初のRedisコマンド実行時）
3. ログでリトライ処理を確認
4. リトライ間隔が`MinRetryBackoff`と`MaxRetryBackoff`の範囲内であることを確認

**期待される動作**:
- 接続エラー時に自動的にリトライが実行される
- 最大2回までリトライ（初回 + 2回のリトライ）
- リトライ間隔は8ms～512msの範囲で指数バックオフ
- ログでリトライ処理が確認できる

## 4. データモデル

### 4.1 設定構造体

```go
// RedisClusterConfig はRedis Cluster設定（Cache Server用）
type RedisClusterConfig struct {
    Addrs           []string      `mapstructure:"addrs"`            // Redis Clusterのアドレスリスト
    MaxRetries      int           `mapstructure:"max_retries"`       // コマンド失敗時の最大リトライ数（デフォルト: 2）
    MinRetryBackoff time.Duration `mapstructure:"min_retry_backoff"` // リトライ間隔（最小）（デフォルト: 8ms）
    MaxRetryBackoff time.Duration `mapstructure:"max_retry_backoff"` // リトライ間隔（最大）（デフォルト: 512ms）
    DialTimeout     time.Duration `mapstructure:"dial_timeout"`      // 接続確立のタイムアウト（デフォルト: 5s）
    ReadTimeout     time.Duration `mapstructure:"read_timeout"`      // 読み取りタイムアウト（デフォルト: 3s）
    PoolSize        int           `mapstructure:"pool_size"`         // 接続プールサイズ（デフォルト: CPU数×10）
    PoolTimeout     time.Duration `mapstructure:"pool_timeout"`      // プールから接続を取り出す際の待機時間（デフォルト: 4s）
}

// RedisSingleConfig は単一Redis接続設定（Jobqueue用）
type RedisSingleConfig struct {
    Addr           string        `mapstructure:"addr"`            // 単一Redis接続アドレス（例: "localhost:6379"）
    MaxRetries     int           `mapstructure:"max_retries"`       // コマンド失敗時の最大リトライ数（デフォルト: 2）
    MinRetryBackoff time.Duration `mapstructure:"min_retry_backoff"` // リトライ間隔（最小）（デフォルト: 8ms）
    MaxRetryBackoff time.Duration `mapstructure:"max_retry_backoff"` // リトライ間隔（最大）（デフォルト: 512ms）
    DialTimeout     time.Duration `mapstructure:"dial_timeout"`      // 接続確立のタイムアウト（デフォルト: 5s）
    ReadTimeout     time.Duration `mapstructure:"read_timeout"`      // 読み取りタイムアウト（デフォルト: 3s）
    WriteTimeout    time.Duration `mapstructure:"write_timeout"`     // 書き込みタイムアウト（デフォルト: 3s）
    PoolSize        int           `mapstructure:"pool_size"`          // 接続プールサイズ（デフォルト: CPU数×10）
    PoolTimeout     time.Duration `mapstructure:"pool_timeout"`       // プールから接続を取り出す際の待機時間（デフォルト: 4s）
}
```

### 4.2 設定ファイル構造

```yaml
redis:
  # ジョブキュー用Redis接続設定（単一接続、1台）
  jobqueue:
    addr: "localhost:6379"
    max_retries: 2
    min_retry_backoff: 8ms
    max_retry_backoff: 512ms
    dial_timeout: 5s
    read_timeout: 3s
    write_timeout: 3s
    pool_size: 10
    pool_timeout: 4s
  
  # デフォルト用Redis接続設定（複数台対応可能、rate limit等で使用）
  default:
    cluster:
      addrs:
        - "localhost:6379"
      max_retries: 2
      min_retry_backoff: 8ms
      max_retry_backoff: 512ms
      dial_timeout: 5s
      read_timeout: 3s
      pool_size: 10
      pool_timeout: 4s
```

## 5. エラーハンドリング

### 5.1 接続エラー時の処理

**既存の実装**:
- `initStore`関数でエラーが発生した場合、`fail-open`方式でリクエストを許可
- エラーログを出力して、リクエストを許可

**本実装での変更**:
- 接続オプション設定により、自動リトライが実行される
- リトライ後もエラーが発生した場合、既存の`fail-open`方式で処理
- エラーログにリトライ情報を含める

### 5.2 設定値の検証

**検証項目**:
- 設定値が0以下の場合はデフォルト値を使用
- `PoolSize`が0以下の場合は`runtime.NumCPU() * 10`を使用
- タイムアウト設定が0以下の場合は適切なデフォルト値を使用

## 6. テスト戦略

### 6.1 単体テスト

#### 6.1.1 設定読み込みテスト
- `RedisClusterConfig`構造体が正しく設定ファイルから読み込まれることを確認
- `RedisSingleConfig`構造体が正しく設定ファイルから読み込まれることを確認
- デフォルト値が正しく適用されることを確認

#### 6.1.2 接続オプション設定テスト
- `initStore`関数（Cache Server用）で接続オプションが正しく設定されることを確認
- `NewClient`関数（Jobqueue用）で接続オプションが正しく設定されることを確認
- `NewServer`関数（Jobqueue用）で接続オプションが正しく設定されることを確認
- 設定値が0以下の場合にデフォルト値が使用されることを確認

### 6.2 統合テスト

#### 6.2.1 遅延接続テスト
- Redis停止状態でAPIサーバーが起動できることを確認
- 最初のRedisコマンド実行時に接続が確立されることを確認

#### 6.2.2 自動再接続テスト
- Redis停止→再起動時に自動的に再接続されることを確認
- 再接続後の動作が正常であることを確認

#### 6.2.3 リトライ機能テスト
- 接続エラー時にリトライが実行されることを確認
- リトライ間隔が適切であることを確認

### 6.3 動作確認テスト

#### 6.3.1 手動動作確認
- Redis環境での動作確認
- ログでの接続確立・再接続・リトライの確認

## 7. 実装順序

1. **Phase 1**: 設定構造体の拡張
   - `RedisClusterConfig`構造体に接続オプションフィールドを追加（Cache Server用）
   - `RedisSingleConfig`構造体に接続オプションフィールドを追加（Jobqueue用）
   - 設定ファイルの例を更新

2. **Phase 2**: 接続オプション設定の実装（Cache Server用）
   - `initStore`関数で接続オプションを設定
   - デフォルト値の適用ロジックを実装

3. **Phase 2-2**: 接続オプション設定の実装（Jobqueue用）
   - `NewClient`関数で接続オプションを設定
   - `NewServer`関数で接続オプションを設定
   - デフォルト値の適用ロジックを実装

4. **Phase 3**: 既存環境の確認
   - `docker-compose.redis.yml`の確認
   - `scripts/start-redis.sh`の確認
   - 設定ファイルの確認

5. **Phase 4**: 動作確認
   - 動作確認前の設定確認と変更（`config.yaml`と`cacheserver.yaml`）
   - Cache Server用の遅延接続・自動再接続・リトライ機能の動作確認
   - Jobqueue用の遅延接続・自動再接続・リトライ機能の動作確認
   - 動作確認後の設定復元（必要に応じて）

6. **Phase 5**: テスト実装
   - 単体テストの実装
   - 統合テストの実装

## 8. 注意事項

### 8.1 既存機能への影響
- 既存のRedis接続機能を壊さないこと
- 既存の設定ファイルとの互換性を保つこと（`addr`や`addrs`のみでも動作する）
- Cache Server用とJobqueue用の両方の接続を修正すること

### 8.2 パフォーマンス
- 接続プール設定により、接続の再利用が適切に行われること
- リトライ処理が過度にパフォーマンスに影響を与えないこと

### 8.3 エラーハンドリング
- 既存の`fail-open`方式を維持すること
- エラーログに適切な情報を含めること

### 8.4 動作確認時の設定変更
- develop環境では、デフォルトで`storage_type: "auto"`かつ`redis.default.cluster.addrs: []`のため、Memoryストレージが使用される
- Cache Server用の動作確認のためには、以下のいずれかの対応が必要：
  - `config/develop/config.yaml`の`api.rate_limit.storage_type`を`"redis"`に変更
  - `config/develop/cacheserver.yaml`の`redis.default.cluster.addrs`に`["localhost:6379"]`を設定
- 動作確認後、必要に応じて設定を元に戻すこと
