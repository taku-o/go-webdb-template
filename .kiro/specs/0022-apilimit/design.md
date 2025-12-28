# APIレートリミット機能設計書

## 1. 概要

### 1.1 設計の目的
要件定義書に基づき、APIの過剰な呼び出しを防止するためのレートリミット機能の詳細設計を定義する。IPアドレス単位でリクエスト数を制限し、環境別ストレージ（In-Memory/Redis Cluster）を使用してレートリミットを管理する。

### 1.2 設計の範囲
- ulule/limiterライブラリを使用したレートリミットミドルウェアの実装
- IPアドレス単位でのリクエスト数制限
- 環境別ストレージ（開発: In-Memory、staging/production: Redis Cluster）
- 設定ファイルによる閾値管理
- レートリミット超過時の適切なレスポンス（HTTP 429、X-RateLimit-*ヘッダー）
- /apiエンドポイント全体へのレートリミット適用

**本設計の範囲外**:
- 認証ロジックの変更（既存の認証ミドルウェアは変更しない）
- APIエンドポイントの追加・削除
- アクセス制御ロジックの変更
- ユーザー単位のレートリミット（IPアドレスのみ）

### 1.3 設計方針
- **ulule/limiterライブラリの利用**: Issue #15の記載に従い、利用可能なライブラリを使用
- **fail-open方式**: レートリミットチェック時のエラー時はAPIへのリクエストを許可（補助的な機能のため）
- **設定による無効化**: レートリミット設定が未指定または`enabled: false`の場合は機能を無効化
- **既存機能への影響なし**: 認証ロジックやAPI動作は一切変更しない
- **シンプルな実装**: /apiエンドポイント全体に適用（Public APIのみの判定は複雑なため除外）

## 2. アーキテクチャ設計

### 2.1 ディレクトリ構造

#### 2.1.1 変更前の構造
```
server/
├── internal/
│   ├── api/
│   │   ├── handler/
│   │   │   ├── user_handler.go
│   │   │   ├── post_handler.go
│   │   │   └── today_handler.go
│   │   └── router/
│   │       └── router.go
│   ├── auth/
│   │   └── middleware.go
│   └── config/
│       └── config.go
```

#### 2.1.2 変更後の構造
```
server/
├── internal/
│   ├── api/
│   │   ├── handler/
│   │   │   ├── user_handler.go
│   │   │   ├── post_handler.go
│   │   │   └── today_handler.go
│   │   └── router/
│   │       └── router.go          # 修正: レートリミットミドルウェアの適用
│   ├── auth/
│   │   └── middleware.go           # 変更なし
│   ├── config/
│   │   └── config.go               # 修正: レートリミット設定の追加
│   └── ratelimit/                  # 新規: レートリミットパッケージ
│       └── middleware.go           # 新規: レートリミットミドルウェア
```

### 2.2 ファイル構成

#### 2.2.1 サーバー側（Go）

**`server/internal/ratelimit/middleware.go`**: レートリミットミドルウェア
- ulule/limiterライブラリを使用したレートリミット実装
- Echoミドルウェアとして実装
- IPアドレス単位でのリクエスト数制限
- 環境別ストレージ（In-Memory/Redis Cluster）の管理
- X-RateLimit-*ヘッダーの付与
- エラーハンドリング（fail-open方式）

**`server/internal/api/router/router.go`**: ルーター設定
- レートリミットミドルウェアの適用
- 認証ミドルウェアの前に適用（認証前に制限）

**`server/internal/config/config.go`**: 設定管理
- `APIConfig`構造体に`RateLimit`設定を追加
- `Config`構造体に`CacheServer`設定を追加
- `cacheserver.yaml`の読み込み処理を追加
- レートリミット設定とキャッシュサーバー設定の読み込み

#### 2.2.2 設定ファイル

**`config/develop/config.yaml`**: 開発環境設定
- レートリミット設定の追加（`enabled: true`、閾値設定）

**`config/staging/config.yaml`**: Staging環境設定
- レートリミット設定の追加（`enabled: true`、閾値設定）

**`config/production/config.yaml.example`**: 本番環境設定例
- レートリミット設定の追加（`enabled: true`、閾値設定）

### 2.3 システム構成図

```
┌─────────────────────────────────────────────────────────┐
│              サーバー（Go）                                │
│                                                           │
│  ┌──────────────────────────────────────────────────┐   │
│  │  server/internal/api/router/router.go            │   │
│  │  - Echoルーター設定                                │   │
│  │  - ミドルウェア適用順序:                           │   │
│  │    1. Recover                                     │   │
│  │    2. CORS                                        │   │
│  │    3. RateLimit (新規)                            │   │
│  │    4. Auth (既存)                                 │   │
│  └──────────────────┬───────────────────────────────┘   │
│                     │                                     │
│                     ▼                                     │
│  ┌──────────────────────────────────────────────────┐   │
│  │  server/internal/ratelimit/middleware.go         │   │
│  │  - IPアドレス取得 (echo.Context.RealIP())         │   │
│  │  - レートリミットチェック                          │   │
│  │  - X-RateLimit-*ヘッダー付与                       │   │
│  └──────────────────┬───────────────────────────────┘   │
│                     │                                     │
│         ┌───────────┴───────────┐                      │
│         │                         │                      │
│         ▼                         ▼                      │
│  ┌──────────────┐        ┌──────────────┐              │
│  │ In-Memory    │        │ Redis        │              │
│  │ Store        │        │ Cluster      │              │
│  │              │        │              │              │
│  │ (開発環境)    │        │ (Staging/    │              │
│  │              │        │  Production) │              │
│  └──────────────┘        └──────────────┘              │
│                                                           │
│  ┌──────────────────────────────────────────────────┐   │
│  │  server/internal/auth/middleware.go              │   │
│  │  - JWT認証 (既存、変更なし)                       │   │
│  └──────────────────┬───────────────────────────────┘   │
│                     │                                     │
│                     ▼                                     │
│  ┌──────────────────────────────────────────────────┐   │
│  │  APIエンドポイント                              │   │
│  │  - /api/users/*                                │   │
│  │  - /api/posts/*                                │   │
│  │  - /api/today                                  │   │
│  └──────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────┘
```

## 3. 詳細設計

### 3.1 レートリミットミドルウェアの実装

#### 3.1.1 ミドルウェア構造

**`server/internal/ratelimit/middleware.go`**の実装:

```go
package ratelimit

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
	redisstore "github.com/ulule/limiter/v3/drivers/store/redis"
	"github.com/taku-o/go-webdb-template/internal/config"
	"github.com/redis/go-redis/v9"
)

// RateLimitConfig はレートリミット設定
type RateLimitConfig struct {
	Enabled           bool
	RequestsPerMinute int
	RequestsPerHour   int // オプション
}

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

	// レートリミットの設定（1分あたりのリクエスト数）
	rate := limiter.Rate{
		Period: time.Minute,
		Limit:  int64(cfg.API.RateLimit.RequestsPerMinute),
	}

	// ストレージの初期化
	store, err := initStore(cfg)
	if err != nil {
		// fail-open方式: エラー時はログに記録し、リクエストを許可
		// ログ出力は実装時に追加
		return func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				return next(c)
			}
		}, nil
	}

	// limiterインスタンスの作成
	instance := limiter.New(store, rate)

	// ミドルウェア関数の返却
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// IPアドレスの取得
			ip := c.RealIP()
			if ip == "" {
				// IPアドレスが取得できない場合は許可
				return next(c)
			}

			// レートリミットチェック
			context, err := instance.Get(c.Request().Context(), ip)
			if err != nil {
				// fail-open方式: エラー時はログに記録し、リクエストを許可
				// ログ出力は実装時に追加
				return next(c)
			}

			// X-RateLimit-*ヘッダーの設定
			c.Response().Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", context.Limit))
			c.Response().Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", context.Remaining))
			c.Response().Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", context.Reset))

			// レートリミット超過時
			if context.Reached {
				return c.JSON(http.StatusTooManyRequests, map[string]interface{}{
					"code":    429,
					"message": "Too Many Requests",
				})
			}

			// レートリミット内の場合は次のハンドラーを実行
			return next(c)
		}
	}, nil
}

// initStore は環境に応じたストレージを初期化
func initStore(cfg *config.Config) (limiter.Store, error) {
	// キャッシュサーバー設定からRedis Clusterのアドレスを取得
	if len(cfg.CacheServer.Redis.Cluster.Addrs) == 0 {
		// In-Memoryストレージを使用
		return memory.NewStore()
	}

	// Redis Clusterを使用
	rdb := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: cfg.CacheServer.Redis.Cluster.Addrs,
	})

	// Redisストアの作成
	return redisstore.NewStoreWithOptions(rdb, limiter.StoreOptions{
		Prefix: "ratelimit",
	})
}
```

#### 3.1.2 実装上の注意事項

- **IPアドレスの取得**: `echo.Context.RealIP()`を使用（X-Forwarded-Forヘッダーを自動処理）
- **エラーハンドリング**: fail-open方式（エラー時はリクエストを許可）
- **ログ出力**: エラー時はログに記録（実装時に追加）
- **ストレージの切り替え**: `cacheserver.yaml`のRedis Cluster設定の有無で自動切り替え

### 3.2 設定管理の実装

#### 3.2.1 APIConfig構造体の拡張

**`server/internal/config/config.go`**の修正:

```go
// Config はアプリケーション全体の設定を保持する構造体
type Config struct {
	Server      ServerConfig      `mapstructure:"server"`
	Admin       AdminConfig       `mapstructure:"admin"`
	Database    DatabaseConfig    `mapstructure:"database"`
	Logging     LoggingConfig     `mapstructure:"logging"`
	CORS        CORSConfig        `mapstructure:"cors"`
	API         APIConfig         `mapstructure:"api"`
	CacheServer CacheServerConfig `mapstructure:"cache_server"` // 新規追加
}

// APIConfig はAPIキー設定
type APIConfig struct {
	CurrentVersion     string   `mapstructure:"current_version"`
	PublicKey          string   `mapstructure:"public_key"`
	SecretKey          string   `mapstructure:"secret_key"`
	InvalidVersions    []string `mapstructure:"invalid_versions"`
	Auth0IssuerBaseURL string   `mapstructure:"auth0_issuer_base_url"` // Auth0のIssuer Base URL
	RateLimit          RateLimitConfig `mapstructure:"rate_limit"`     // 新規追加
}

// RateLimitConfig はレートリミット設定
type RateLimitConfig struct {
	Enabled           bool `mapstructure:"enabled"`
	RequestsPerMinute int  `mapstructure:"requests_per_minute"`
	RequestsPerHour   int  `mapstructure:"requests_per_hour"` // オプション
}

// CacheServerConfig はキャッシュサーバー設定
type CacheServerConfig struct {
	Redis RedisConfig `mapstructure:"redis"`
}

// RedisConfig はRedis設定
type RedisConfig struct {
	Cluster RedisClusterConfig `mapstructure:"cluster"`
}

// RedisClusterConfig はRedis Cluster設定
type RedisClusterConfig struct {
	Addrs []string `mapstructure:"addrs"` // Redis Clusterのアドレスリスト
}
```

#### 3.2.2 cacheserver.yamlの読み込み

**`server/internal/config/config.go`**の`Load()`関数を修正:

```go
// Load は指定された環境の設定ファイルを読み込む
func Load() (*Config, error) {
	// ... 既存のコード ...

	// データベース設定ファイルのマージ
	viper.SetConfigName("database")
	if err := viper.MergeInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read database config file: %w", err)
	}

	// キャッシュサーバー設定ファイルのマージ（新規追加）
	viper.SetConfigName("cacheserver")
	if err := viper.MergeInConfig(); err != nil {
		// cacheserver.yamlが存在しない場合はエラーにしない（オプショナル）
		// ログに記録するか、エラーを無視する
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// ... 既存のコード ...
}
```

#### 3.2.2 設定ファイルの構造

**`config/develop/config.yaml`**の追加:

```yaml
api:
  current_version: "v2"
  secret_key: "RNrxs7Rt1ZViughEGb8J08Uc1uQobSOZRRb+BmnGaag="
  invalid_versions:
    - "v1"
  auth0_issuer_base_url: "https://dev-oaa5vtzmld4dsxtd.jp.auth0.com"
  rate_limit:                                    # 新規追加
    enabled: true
    requests_per_minute: 60
    requests_per_hour: 1000
```

**`config/staging/config.yaml`**と**`config/production/config.yaml.example`**にも同様の設定を追加。

**`config/develop/cacheserver.yaml`**の作成（新規）:

```yaml
# 開発環境ではRedis Clusterを使用しない（In-Memoryストレージを使用）
redis:
  cluster:
    addrs: []
```

**`config/staging/cacheserver.yaml`**の作成（新規）:

```yaml
redis:
  cluster:
    addrs:
      - host1:6379
      - host2:6379
      - host3:6379
```

**`config/production/cacheserver.yaml.example`**の作成（新規）:

```yaml
redis:
  cluster:
    addrs:
      - host1:6379
      - host2:6379
      - host3:6379
```

### 3.3 ルーターへの統合

#### 3.3.1 router.goの修正

**`server/internal/api/router/router.go`**の修正:

```go
package router

import (
	"fmt"
	"net/http"
	"os"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humaecho"
	"github.com/taku-o/go-webdb-template/internal/api/handler"
	"github.com/taku-o/go-webdb-template/internal/auth"
	"github.com/taku-o/go-webdb-template/internal/config"
	"github.com/taku-o/go-webdb-template/internal/ratelimit"  // 新規追加
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// NewRouter は新しいEchoルーターを作成
func NewRouter(userHandler *handler.UserHandler, postHandler *handler.PostHandler, todayHandler *handler.TodayHandler, cfg *config.Config) *echo.Echo {
	e := echo.New()

	// デバッグモードの設定（開発環境のみ）
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "develop"
	}
	if env == "develop" {
		e.Debug = true
	}

	// AUTH0_ISSUER_BASE_URLが空の場合、サーバー起動時にエラーを発生させる
	if cfg.API.Auth0IssuerBaseURL == "" {
		panic("AUTH0_ISSUER_BASE_URL is required in config")
	}

	// Recoverミドルウェア
	e.Use(middleware.Recover())

	// CORS設定
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     cfg.CORS.AllowedOrigins,
		AllowMethods:     cfg.CORS.AllowedMethods,
		AllowHeaders:     cfg.CORS.AllowedHeaders,
		AllowCredentials: true,
	}))

	// レートリミットミドルウェア（認証ミドルウェアの前に適用）
	rateLimitMiddleware, err := ratelimit.NewRateLimitMiddleware(cfg)
	if err != nil {
		// エラー時はログに記録し、サーバー起動を継続（fail-open方式）
		// ログ出力は実装時に追加
	} else {
		e.Use(rateLimitMiddleware)
	}

	// Health check
	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	// Huma API設定
	humaConfig := huma.DefaultConfig("go-webdb-template API", "1.0.0")
	humaConfig.DocsPath = "/docs"
	humaConfig.Servers = []*huma.Server{
		{
			URL:         fmt.Sprintf("http://localhost:%d", cfg.Server.Port),
			Description: "Development server",
		},
	}

	// SecurityScheme設定の追加
	humaConfig.Components = &huma.Components{
		SecuritySchemes: map[string]*huma.SecurityScheme{
			"bearerAuth": {
				Type:         "http",
				Scheme:       "bearer",
				BearerFormat: "JWT",
			},
		},
	}

	// Huma APIインスタンスの作成（ルートレベル、認証なし）
	humaAPI := humaecho.New(e, humaConfig)

	// Humaミドルウェアとして認証を追加（/api/パスのみ）
	humaAPI.UseMiddleware(auth.NewHumaAuthMiddleware(&cfg.API, env, cfg.API.Auth0IssuerBaseURL))

	// Humaエンドポイントの登録
	handler.RegisterUserEndpoints(humaAPI, userHandler)
	handler.RegisterPostEndpoints(humaAPI, postHandler)
	handler.RegisterTodayEndpoints(humaAPI, todayHandler)

	return e
}
```

#### 3.3.2 ミドルウェアの適用順序

1. **Recover**: パニック回復
2. **CORS**: クロスオリジンリクエスト処理
3. **RateLimit**: レートリミットチェック（新規追加、認証前に適用）
4. **Auth**: JWT認証（既存）

### 3.4 エラーハンドリング設計

#### 3.4.1 fail-open方式の実装

レートリミットチェック時のエラー（Redisサーバーと通信できない場合など）は、以下のように処理：

1. **エラーログの記録**: エラー内容をログに記録
2. **リクエストの許可**: エラー時はAPIへのリクエストを許可
3. **機能の継続**: レートリミット機能が失敗しても、APIサーバーは正常に動作

#### 3.4.2 エラーケース

- **Redis接続エラー**: Redis Clusterへの接続に失敗した場合
- **Redis操作エラー**: レートリミットチェック時のRedis操作エラー
- **IPアドレス取得エラー**: IPアドレスが取得できない場合（許可）
- **ストレージ初期化エラー**: ストレージの初期化に失敗した場合（In-Memoryにフォールバック）

### 3.5 レスポンスヘッダー設計

#### 3.5.1 X-RateLimit-*ヘッダー

すべてのレスポンスに以下のヘッダーを付与：

- **X-RateLimit-Limit**: 制限値（例: `60`）
- **X-RateLimit-Remaining**: 残りリクエスト数（例: `45`）
- **X-RateLimit-Reset**: リセット時刻（Unix timestamp、例: `1706342400`）

#### 3.5.2 レスポンスボディ

レートリミット超過時（HTTP 429）:

```json
{
  "code": 429,
  "message": "Too Many Requests"
}
```

## 4. データフロー

### 4.1 リクエスト処理フロー

```
1. クライアント → サーバー
   ↓
2. Echoルーター
   ↓
3. Recoverミドルウェア
   ↓
4. CORSミドルウェア
   ↓
5. RateLimitミドルウェア（新規）
   ├─ IPアドレス取得
   ├─ レートリミットチェック
   │  ├─ In-Memory/Redis Clusterからカウント取得
   │  ├─ カウント増加
   │  └─ 制限超過判定
   ├─ X-RateLimit-*ヘッダー設定
   └─ 制限超過時: HTTP 429返却
   ↓
6. Authミドルウェア（既存）
   ├─ JWT検証
   └─ アクセスレベルチェック
   ↓
7. APIハンドラー
   └─ レスポンス返却
```

### 4.2 ストレージ選択フロー

```
1. サーバー起動時
   ↓
2. config.goでcacheserver.yamlを読み込み
   ↓
3. CacheServerConfigからRedis Clusterのアドレスを取得
   ├─ addrsが設定されている場合（空でない）
   │  └─ Redis Clusterストレージを初期化
   │     └─ redis.NewClusterClientを使用
   └─ addrsが設定されていない場合（空）
      └─ In-Memoryストレージを初期化
         └─ memory.NewStoreを使用
   ↓
4. レートリミットミドルウェア作成
```

## 5. 依存関係

### 5.1 追加する依存関係

**`go.mod`**に追加:

```
github.com/ulule/limiter/v3 v3.11.2
github.com/redis/go-redis/v9 v9.5.0
```

### 5.2 既存の依存関係

- `github.com/labstack/echo/v4`: Echoフレームワーク（既存）
- `github.com/danielgtaylor/huma/v2`: Humaフレームワーク（既存）
- `github.com/spf13/viper`: 設定管理（既存）

## 6. テスト設計

### 6.1 単体テスト

#### 6.1.1 レートリミットミドルウェアのテスト

- IPアドレス単位でのリクエスト数制限の確認
- レートリミット超過時のHTTP 429返却の確認
- X-RateLimit-*ヘッダーの付与確認
- エラー時のfail-open動作の確認

#### 6.1.2 ストレージのテスト

- In-Memoryストレージの動作確認
- Redis Clusterストレージの動作確認（統合テスト）
- ストレージ切り替えの確認

### 6.2 統合テスト

- レートリミットミドルウェアと認証ミドルウェアの連携確認
- 複数のIPアドレスからの同時リクエストの確認
- Redis Cluster接続エラー時の動作確認

## 7. 実装上の注意事項

### 7.1 パフォーマンス

- レートリミットチェックによるレスポンス時間への影響を最小化
- Redis Cluster使用時は接続プールを適切に管理
- In-Memory使用時はメモリ使用量を監視

### 7.2 セキュリティ

- IPアドレスの偽造対策: `echo.Context.RealIP()`を使用（X-Forwarded-Forヘッダーを適切に処理）
- Redis接続のセキュリティ: 必要に応じてTLS接続を検討

### 7.3 メンテナンス性

- 設定の一元管理: `config.yaml`で管理
- ログ出力: エラー時は適切にログに記録
- コードの可読性: 明確な実装とコメント

## 8. 参考情報

### 8.1 関連Issue
- GitHub Issue #15: APIサーバーへのリクエストにrate limit機能をつける（0022-apilimit）

### 8.2 既存ドキュメント
- `.kiro/specs/0022-apilimit/requirements.md`: 要件定義書

### 8.3 技術スタック
- **Go言語**: 1.24+
- **Echoフレームワーク**: v4
- **ulule/limiter**: v3
- **Redis**: Redis Cluster対応（`redis.NewClusterClient`を使用）

### 8.4 参考資料
- [ulule/limiter Documentation](https://github.com/ulule/limiter)
- [Echo Middleware](https://echo.labstack.com/docs/middleware)
- [Redis Go Client](https://github.com/redis/go-redis)
