# Public API Key認証機能設計書

## 1. 概要

### 1.1 設計の目的
要件定義書に基づき、JWT形式のPublic APIキーによる認証機能の詳細設計を定義する。既存のアーキテクチャに統合し、セキュアなAPIアクセスを実現する。

### 1.2 設計の範囲
- JWT検証機能の実装設計
- 認証ミドルウェアの実装設計
- GoAdmin管理画面でのキー発行機能の設計
- クライアント側（TypeScript/Next.js）の実装設計
- 設定構造体の拡張設計
- 秘密鍵生成ツールの設計
- エラーハンドリング設計
- テスト戦略

### 1.3 設計方針
- **既存ライブラリの活用**: `github.com/golang-jwt/jwt/v5`を使用
- **既存アーキテクチャの維持**: レイヤードアーキテクチャを維持
- **設定ファイルベース**: DB管理は行わない（Issue #21の要件）
- **環境別分離**: 秘密鍵とAPIキーは環境別に分離
- **既存機能との統合**: 既存のAPIルーター、GoAdmin管理画面、クライアント実装との統合

## 2. アーキテクチャ設計

### 2.1 ディレクトリ構造

#### 2.1.1 変更前の構造
```
server/
├── internal/
│   ├── api/
│   │   ├── router/
│   │   │   └── router.go      # ルーティング（認証なし）
│   │   └── handler/            # ハンドラー
│   ├── config/
│   │   └── config.go          # 設定構造体（API設定なし）
│   └── admin/
│       └── pages/              # GoAdminカスタムページ
├── cmd/
│   ├── server/
│   │   └── main.go            # APIサーバー
│   └── admin/
│       └── main.go             # 管理画面サーバー
└── ...
client/
└── src/
    └── lib/
        └── api.ts              # APIクライアント（認証なし）
```

#### 2.1.2 変更後の構造
```
server/
├── internal/
│   ├── auth/                   # 新規: 認証機能
│   │   ├── jwt.go             # JWT検証機能
│   │   └── middleware.go      # 認証ミドルウェア
│   ├── api/
│   │   ├── router/
│   │   │   └── router.go      # ルーティング（認証ミドルウェア追加）
│   │   └── handler/            # ハンドラー（変更不要）
│   ├── config/
│   │   ├── config.go          # 設定構造体（APIConfig追加）
│   │   └── testdata/
│   │       └── develop/
│   │           └── api_key.yaml  # テスト用ダミーAPIキー設定
│   └── admin/
│       └── pages/
│           └── api_key.go      # 新規: APIキー発行ページ
├── cmd/
│   ├── server/
│   │   └── main.go            # APIサーバー（変更不要）
│   ├── admin/
│   │   └── main.go             # 管理画面サーバー（ページ登録追加）
│   └── generate-secret/        # 新規: 秘密鍵生成ツール
│       └── main.go
└── ...
client/
└── src/
    └── lib/
        └── api.ts              # APIクライアント（認証ヘッダー追加）
```

### 2.2 認証処理の実行フロー

```
┌─────────────────────────────────────────────────────────────┐
│              1. クライアントからのリクエスト                    │
│              GET /api/users                                 │
│              Authorization: Bearer <JWT_API_KEY>            │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              2. 認証ミドルウェア（AuthMiddleware）              │
│              - AuthorizationヘッダーからJWTトークンを抽出      │
│              - JWTトークンが存在しない場合: 401エラー          │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              3. JWT検証（ValidateJWT）                         │
│              - JWT署名の検証（秘密鍵による検証）               │
│              - クレームの検証（iss, type, version, env）      │
│              - 無効バージョンリストとの照合                   │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              4. スコープ検証                                  │
│              - HTTPメソッドに応じたスコープチェック            │
│              - GET: "read"スコープが必要                      │
│              - POST/PUT/DELETE: "write"スコープが必要         │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              5. リクエスト処理                                │
│              - Handler層での処理                              │
│              - Service層での処理                              │
│              - Repository層での処理                           │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              6. レスポンス返却                                 │
│              - 正常レスポンス（200, 201等）                   │
│              - またはエラーレスポンス（400, 404, 500等）       │
└─────────────────────────────────────────────────────────────┘
```

### 2.3 既存アーキテクチャとの統合

#### 2.3.1 APIルーターへの統合
- `server/internal/api/router/router.go`の`NewRouter`関数を拡張
- `/api/*`パスに認証ミドルウェアを適用
- 既存のCORSミドルウェアとの共存
- ミドルウェアの適用順序: CORS → 認証 → ハンドラー

#### 2.3.2 設定構造体への統合
- `server/internal/config/config.go`の`Config`構造体に`APIConfig`フィールドを追加
- 既存の設定読み込み処理（`config.Load()`）を維持
- 環境別設定ファイル（`config/{env}/config.yaml`）に`api`セクションを追加

#### 2.3.3 GoAdmin管理画面への統合
- 既存のカスタムページ実装パターン（`server/internal/admin/pages/`）に従う
- `RegisterCustomPages`関数にキー発行ページを追加
- データベースマイグレーションでメニュー項目を追加

#### 2.3.4 クライアント側への統合
- 既存の`ApiClient`クラス（`client/src/lib/api.ts`）を拡張
- `request`メソッドに`Authorization`ヘッダーを追加
- 環境変数（`NEXT_PUBLIC_API_KEY`）からAPIキーを取得

## 3. コンポーネント設計

### 3.1 JWT検証機能（jwt.go）

#### 3.1.1 構造体定義
```go
// server/internal/auth/jwt.go

package auth

import (
    "errors"
    "fmt"
    "time"

    "github.com/golang-jwt/jwt/v5"
    "github.com/example/go-webdb-template/internal/config"
)

// JWTClaims はJWTのクレーム構造
type JWTClaims struct {
    Issuer   string   `json:"iss"`
    Subject  string   `json:"sub"`
    Type     string   `json:"type"`     // "public" | "private"
    Scope    []string `json:"scope"`    // ["read", "write"]
    IssuedAt int64    `json:"iat"`
    Version  string   `json:"version"`
    Env      string   `json:"env"`
    jwt.RegisteredClaims
}

// JWTValidator はJWT検証機能を提供
type JWTValidator struct {
    secretKey      string
    invalidVersions []string
    currentEnv     string
}

// NewJWTValidator は新しいJWTValidatorを作成
func NewJWTValidator(cfg *config.APIConfig, env string) *JWTValidator {
    return &JWTValidator{
        secretKey:      cfg.SecretKey,
        invalidVersions: cfg.InvalidVersions,
        currentEnv:     env,
    }
}
```

#### 3.1.2 JWT検証関数
```go
// ValidateJWT はJWTトークンを検証
func (v *JWTValidator) ValidateJWT(tokenString string) (*JWTClaims, error) {
    // JWTトークンをパース
    token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
        // 署名アルゴリズムの検証
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(v.secretKey), nil
    })

    if err != nil {
        return nil, fmt.Errorf("failed to parse JWT: %w", err)
    }

    // クレームの取得
    claims, ok := token.Claims.(*JWTClaims)
    if !ok || !token.Valid {
        return nil, errors.New("invalid token")
    }

    // クレームの検証
    if err := v.validateClaims(claims); err != nil {
        return nil, err
    }

    return claims, nil
}

// validateClaims はクレームを検証
func (v *JWTValidator) validateClaims(claims *JWTClaims) error {
    // issの検証
    if claims.Issuer != "go-webdb-template" {
        return errors.New("invalid issuer")
    }

    // typeの検証
    if claims.Type != "public" && claims.Type != "private" {
        return errors.New("invalid token type")
    }

    // versionの検証（無効バージョンリストとの照合）
    if v.IsVersionInvalid(claims.Version) {
        return errors.New("invalid token version")
    }

    // envの検証
    if claims.Env != v.currentEnv {
        return errors.New("token environment mismatch")
    }

    return nil
}

// IsVersionInvalid はバージョンが無効かどうかを判定
func (v *JWTValidator) IsVersionInvalid(version string) bool {
    for _, invalidVersion := range v.invalidVersions {
        if version == invalidVersion {
            return true
        }
    }
    return false
}

// ParseJWTClaims はJWTトークンからクレームをパース（表示用）
func ParseJWTClaims(tokenString string) (*JWTClaims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
        // 署名検証なしでパース（表示用）
        return nil, nil
    })

    if err != nil {
        return nil, fmt.Errorf("failed to parse JWT: %w", err)
    }

    claims, ok := token.Claims.(*JWTClaims)
    if !ok {
        return nil, errors.New("invalid token claims")
    }

    return claims, nil
}
```

### 3.2 認証ミドルウェア（middleware.go）

#### 3.2.1 ミドルウェア構造体
```go
// server/internal/auth/middleware.go

package auth

import (
    "encoding/json"
    "net/http"
    "strings"

    "github.com/example/go-webdb-template/internal/config"
)

// AuthMiddleware は認証ミドルウェア
type AuthMiddleware struct {
    validator *JWTValidator
}

// NewAuthMiddleware は新しい認証ミドルウェアを作成
func NewAuthMiddleware(cfg *config.APIConfig, env string) *AuthMiddleware {
    return &AuthMiddleware{
        validator: NewJWTValidator(cfg, env),
    }
}

// Middleware はHTTPミドルウェア関数を返す
func (m *AuthMiddleware) Middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // AuthorizationヘッダーからJWTトークンを取得
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            m.writeErrorResponse(w, http.StatusUnauthorized, "Authorization header is required")
            return
        }

        // Bearerトークンの抽出
        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            m.writeErrorResponse(w, http.StatusUnauthorized, "Invalid authorization header format")
            return
        }

        tokenString := parts[1]

        // JWT検証
        claims, err := m.validator.ValidateJWT(tokenString)
        if err != nil {
            m.writeErrorResponse(w, http.StatusUnauthorized, "Invalid API key")
            return
        }

        // スコープ検証
        if err := m.validateScope(claims, r.Method); err != nil {
            m.writeErrorResponse(w, http.StatusForbidden, "Insufficient scope")
            return
        }

        // 次のハンドラーを実行
        next.ServeHTTP(w, r)
    })
}

// validateScope はスコープを検証
func (m *AuthMiddleware) validateScope(claims *JWTClaims, method string) error {
    hasRead := false
    hasWrite := false

    for _, scope := range claims.Scope {
        if scope == "read" {
            hasRead = true
        }
        if scope == "write" {
            hasWrite = true
        }
    }

    // GETリクエストにはreadスコープが必要
    if method == "GET" && !hasRead {
        return errors.New("read scope required")
    }

    // POST/PUT/DELETEリクエストにはwriteスコープが必要
    if (method == "POST" || method == "PUT" || method == "DELETE") && !hasWrite {
        return errors.New("write scope required")
    }

    return nil
}

// writeErrorResponse はエラーレスポンスを書き込む
func (m *AuthMiddleware) writeErrorResponse(w http.ResponseWriter, statusCode int, message string) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    
    response := map[string]interface{}{
        "code":    statusCode,
        "message": message,
    }
    
    json.NewEncoder(w).Encode(response)
}
```

### 3.3 GoAdminキー発行ページ（api_key.go）

#### 3.3.1 ページハンドラー
```go
// server/internal/admin/pages/api_key.go

package pages

import (
    "encoding/json"
    "fmt"
    "html/template"
    "net/http"
    "os"
    "time"

    "github.com/GoAdminGroup/go-admin/context"
    "github.com/GoAdminGroup/go-admin/modules/db"
    "github.com/GoAdminGroup/go-admin/template/types"
    "github.com/golang-jwt/jwt/v5"
    "github.com/example/go-webdb-template/internal/config"
    "github.com/example/go-webdb-template/internal/auth"
)

// APIKeyPage はAPIキー発行ページを返す
// 注意: RegisterCustomPagesで"/api-key"と登録すると、実際のURLは"/admin/api-key"になる
// HTML内のリンクも"/admin/api-key"とする必要がある
func APIKeyPage(ctx *context.Context, conn db.Connection) (types.Panel, error) {
    // 設定を取得
    cfg, err := config.Load()
    if err != nil {
        return types.Panel{}, err
    }

    // ダウンロードリクエストの処理
    if ctx.Query("download") == "true" {
        return handleDownload(ctx, cfg)
    }

    // POSTリクエスト: キー生成
    if ctx.Method() == http.MethodPost {
        return handleGenerateKey(ctx, cfg)
    }

    // GETリクエスト: フォーム表示
    return renderAPIKeyPage(ctx, cfg)
}

// handleGenerateKey はAPIキーを生成
func handleGenerateKey(ctx *context.Context, cfg *config.Config) (types.Panel, error) {
    // 現在の環境を取得
    env := os.Getenv("APP_ENV")
    if env == "" {
        env = "develop"
    }

    // JWTトークンを生成
    token, err := generatePublicAPIKey(cfg, env)
    if err != nil {
        return types.Panel{}, err
    }

    // ペイロードをデコード
    claims, err := auth.ParseJWTClaims(token)
    if err != nil {
        return types.Panel{}, err
    }

    // 生成結果を表示
    return renderAPIKeyResult(ctx, token, claims)
}

// generatePublicAPIKey はPublic JWTキーを生成
func generatePublicAPIKey(cfg *config.Config, env string) (string, error) {
    now := time.Now()
    
    claims := &auth.JWTClaims{
        Issuer:   "go-webdb-template",
        Subject:  "public_client",
        Type:     "public",
        Scope:    []string{"read", "write"},
        IssuedAt: now.Unix(),
        Version:  cfg.API.CurrentVersion,
        Env:      env,
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString([]byte(cfg.API.SecretKey))
    if err != nil {
        return "", fmt.Errorf("failed to sign token: %w", err)
    }

    return tokenString, nil
}

// handleDownload はJWTトークンをダウンロード
func handleDownload(ctx *context.Context, cfg *config.Config) (types.Panel, error) {
    // セッションまたはクエリパラメータからトークンを取得
    // （実装詳細は省略）
    
    token := ctx.Query("token")
    if token == "" {
        return types.Panel{}, fmt.Errorf("token not found")
    }

    // ファイル名を生成
    timestamp := time.Now().Format("20060102-150405")
    filename := fmt.Sprintf("api-key-%s.txt", timestamp)

    // ダウンロードレスポンスを設定
    ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
    ctx.Header("Content-Type", "text/plain")
    ctx.Write([]byte(token))

    return types.Panel{}, nil
}

// renderAPIKeyPage はAPIキー発行ページをレンダリング
func renderAPIKeyPage(ctx *context.Context, cfg *config.Config) (types.Panel, error) {
    content := `
<div class="box box-primary">
    <div class="box-header with-border">
        <h3 class="box-title">Public APIキー発行</h3>
    </div>
    <div class="box-body">
        <p>新しいPublic APIキーを発行します。</p>
        <form action="/admin/api-key" method="POST">
            <button type="submit" class="btn btn-primary">
                <i class="fa fa-key"></i> APIキーを発行
            </button>
        </form>
    </div>
</div>
`

    return types.Panel{
        Title:       "APIキー発行",
        Description: "Public APIキーを発行します",
        Content:     template.HTML(content),
    }, nil
}

// renderAPIKeyResult は生成結果をレンダリング
func renderAPIKeyResult(ctx *context.Context, token string, claims *auth.JWTClaims) (types.Panel, error) {
    // ペイロードをJSON形式で整形
    payloadJSON, _ := json.MarshalIndent(claims, "", "  ")
    
    // iatを人間が読める形式に変換
    issuedAt := time.Unix(claims.IssuedAt, 0).Format("2006-01-02 15:04:05")

    content := fmt.Sprintf(`
<div class="box box-success">
    <div class="box-header with-border">
        <h3 class="box-title">APIキー発行完了</h3>
    </div>
    <div class="box-body">
        <div class="form-group">
            <label>JWTトークン</label>
            <textarea class="form-control" rows="3" readonly>%s</textarea>
        </div>
        <div class="form-group">
            <label>JWTペイロード</label>
            <pre class="form-control" style="height: 300px; overflow-y: auto;">%s</pre>
        </div>
        <div class="form-group">
            <label>発行日時</label>
            <p>%s</p>
        </div>
        <div class="form-group">
            <label>バージョン</label>
            <p>%s</p>
        </div>
        <div class="form-group">
            <label>環境</label>
            <p>%s</p>
        </div>
        <div class="form-group">
            <a href="/admin/api-key?download=true&token=%s" class="btn btn-success">
                <i class="fa fa-download"></i> ダウンロード
            </a>
        </div>
    </div>
</div>
`, template.HTMLEscapeString(token), template.HTMLEscapeString(string(payloadJSON)), issuedAt, claims.Version, claims.Env, template.URLQueryEscaper(token))

    return types.Panel{
        Title:       "APIキー発行完了",
        Description: "Public APIキーが発行されました",
        Content:     template.HTML(content),
    }, nil
}
```

### 3.4 設定構造体の拡張

#### 3.4.1 APIConfig構造体
```go
// server/internal/config/config.go

// APIConfig はAPIキー設定
type APIConfig struct {
    CurrentVersion  string   `mapstructure:"current_version"`
    PublicKey       string   `mapstructure:"public_key"`       // オプション
    SecretKey       string   `mapstructure:"secret_key"`      // 必須
    InvalidVersions []string `mapstructure:"invalid_versions"`
}

// Config はアプリケーション全体の設定を保持する構造体
type Config struct {
    Server   ServerConfig   `mapstructure:"server"`
    Admin    AdminConfig    `mapstructure:"admin"`
    Database DatabaseConfig `mapstructure:"database"`
    Logging  LoggingConfig  `mapstructure:"logging"`
    CORS     CORSConfig     `mapstructure:"cors"`
    API      APIConfig      `mapstructure:"api"`  // 新規追加
}
```

### 3.5 認証ミドルウェアの適用

#### 3.5.1 Routerへの統合
```go
// server/internal/api/router/router.go

package router

import (
    "net/http"
    "os"

    "github.com/example/go-webdb-template/internal/api/handler"
    "github.com/example/go-webdb-template/internal/auth"
    "github.com/example/go-webdb-template/internal/config"
    "github.com/gorilla/mux"
    "github.com/rs/cors"
)

// NewRouter は新しいルーターを作成
func NewRouter(userHandler *handler.UserHandler, postHandler *handler.PostHandler, cfg *config.Config) http.Handler {
    r := mux.NewRouter()

    // API routes
    api := r.PathPrefix("/api").Subrouter()

    // 認証ミドルウェアの適用
    env := os.Getenv("APP_ENV")
    if env == "" {
        env = "develop"
    }
    authMiddleware := auth.NewAuthMiddleware(&cfg.API, env)
    api.Use(authMiddleware.Middleware)

    // User routes
    api.HandleFunc("/users", userHandler.CreateUser).Methods("POST")
    api.HandleFunc("/users", userHandler.ListUsers).Methods("GET")
    // ... 既存のルート定義 ...

    // Health check（認証不要）
    r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("OK"))
    }).Methods("GET")

    // CORS middleware
    c := cors.New(cors.Options{
        AllowedOrigins:   cfg.CORS.AllowedOrigins,
        AllowedMethods:   cfg.CORS.AllowedMethods,
        AllowedHeaders:   cfg.CORS.AllowedHeaders,
        AllowCredentials: true,
    })

    return c.Handler(r)
}
```

### 3.6 クライアント側の実装

#### 3.6.1 ApiClientの拡張
```typescript
// client/src/lib/api.ts

class ApiClient {
  private baseURL: string
  private apiKey: string | null

  constructor(baseURL: string) {
    this.baseURL = baseURL
    this.apiKey = process.env.NEXT_PUBLIC_API_KEY || null
    
    // APIキーが設定されていない場合、エラーを投げる
    if (!this.apiKey) {
      throw new Error('NEXT_PUBLIC_API_KEY is not set')
    }
  }

  private async request<T>(
    endpoint: string,
    options?: RequestInit
  ): Promise<T> {
    const url = `${this.baseURL}${endpoint}`
    
    // Authorizationヘッダーを追加
    const headers = {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${this.apiKey}`,
      ...options?.headers,
    }

    const response = await fetch(url, {
      ...options,
      headers,
    })

    if (!response.ok) {
      // エラーレスポンスの処理
      if (response.status === 401 || response.status === 403) {
        const errorData = await response.json().catch(() => ({}))
        throw new Error(errorData.message || response.statusText)
      }
      const error = await response.text()
      throw new Error(error || response.statusText)
    }

    if (response.status === 204) {
      return {} as T
    }

    return response.json()
  }

  // ... 既存のメソッド ...
}
```

### 3.7 秘密鍵生成ツール

#### 3.7.1 generate-secret/main.go
```go
// server/cmd/generate-secret/main.go

package main

import (
    "crypto/rand"
    "encoding/base64"
    "fmt"
    "os"
)

func main() {
    // ランダムな秘密鍵を生成（32バイト = 256ビット）
    secretKey := make([]byte, 32)
    if _, err := rand.Read(secretKey); err != nil {
        fmt.Fprintf(os.Stderr, "Failed to generate secret key: %v\n", err)
        os.Exit(1)
    }

    // Base64エンコード
    encodedKey := base64.URLEncoding.EncodeToString(secretKey)

    // 標準出力に表示
    fmt.Println(encodedKey)
}
```

## 4. データモデル設計

### 4.1 JWTクレーム構造

```go
type JWTClaims struct {
    Issuer   string   `json:"iss"`    // "go-webdb-template"（固定）
    Subject  string   `json:"sub"`    // "public_client"（Publicキー）またはユーザーID（Privateキー）
    Type     string   `json:"type"`    // "public" | "private"
    Scope    []string `json:"scope"`   // ["read", "write"]（固定）
    IssuedAt int64    `json:"iat"`    // Unix timestamp
    Version  string   `json:"version"` // "v1", "v2"等（設定ファイルから取得）
    Env      string   `json:"env"`    // "develop" | "staging" | "production"
    // expは未定義（Public APIキーは無期限）
}
```

### 4.2 設定ファイル構造

```yaml
# config/{env}/config.yaml

api:
  current_version: "v2"              # 現在のバージョン
  public_key: "<JWT_TOKEN>"          # 発行済みのPublic APIキー（オプション）
  secret_key: "<SECRET_KEY>"         # JWT署名用の秘密鍵（必須）
  invalid_versions:                  # 無効化されたバージョンリスト
    - "v1"
```

### 4.3 エラーレスポンス構造

```json
{
  "code": 401,
  "message": "Invalid API key"
}
```

または

```json
{
  "code": 403,
  "message": "Insufficient scope"
}
```

## 5. エラーハンドリング設計

### 5.1 認証エラー

#### 5.1.1 エラーケース
- Authorizationヘッダーが存在しない: 401 Unauthorized
- Bearerトークンの形式が不正: 401 Unauthorized
- JWT署名の検証失敗: 401 Unauthorized
- クレームの検証失敗（iss, type, version, env）: 401 Unauthorized
- 無効バージョンのキー: 401 Unauthorized

#### 5.1.2 エラーレスポンス
```json
{
  "code": 401,
  "message": "Invalid API key"
}
```

### 5.2 スコープエラー

#### 5.2.1 エラーケース
- GETリクエストに`read`スコープがない: 403 Forbidden
- POST/PUT/DELETEリクエストに`write`スコープがない: 403 Forbidden

#### 5.2.2 エラーレスポンス
```json
{
  "code": 403,
  "message": "Insufficient scope"
}
```

### 5.3 クライアント側のエラーハンドリング

#### 5.3.1 APIキー未設定時
- `NEXT_PUBLIC_API_KEY`が設定されていない場合、`ApiClient`のコンストラクタでエラーを投げる
- エラーメッセージ: `"NEXT_PUBLIC_API_KEY is not set"`

#### 5.3.2 認証失敗時
- 401/403エラーを受信した場合、エラーレスポンスの`message`を表示
- 適切なエラーメッセージをユーザーに提示

## 6. テスト戦略

### 6.1 ユニットテスト

#### 6.1.1 JWT検証機能のテスト
- `server/internal/auth/jwt_test.go`を作成
- テストケース:
  - 正常系: 有効なJWTトークンの検証
  - 異常系: 無効な署名のJWTトークン
  - 異常系: 不正なiss
  - 異常系: 不正なtype
  - 異常系: 無効バージョンのキー
  - 異常系: 環境不一致

#### 6.1.2 認証ミドルウェアのテスト
- `server/internal/auth/middleware_test.go`を作成
- テストケース:
  - 正常系: 有効なAPIキーでのリクエスト
  - 異常系: Authorizationヘッダーなし
  - 異常系: 無効なAPIキー
  - 異常系: スコープ不足（GETリクエストにreadスコープなし）
  - 異常系: スコープ不足（POSTリクエストにwriteスコープなし）

#### 6.1.3 GoAdminキー発行ページのテスト
- `server/internal/admin/pages/api_key_test.go`を作成
- テストケース:
  - 正常系: APIキーの生成
  - 正常系: JWTペイロードのデコード
  - 正常系: ダウンロード機能

### 6.2 統合テスト

#### 6.2.1 API認証の統合テスト
- `server/test/integration/api_auth_test.go`を作成
- テストケース:
  - 正常系: 有効なAPIキーでAPIアクセス
  - 異常系: 無効なAPIキーでAPIアクセス
  - 異常系: APIキーなしでAPIアクセス
  - 異常系: スコープ不足でのAPIアクセス

### 6.3 E2Eテスト

#### 6.3.1 クライアント側のテスト
- `client/src/lib/__tests__/api.test.ts`を更新
- テストケース:
  - 正常系: APIキーを設定してAPIアクセス
  - 異常系: APIキー未設定時のエラー
  - 異常系: 401/403エラー時の処理

## 7. 実装上の注意事項

### 7.1 JWTライブラリの使用
- `github.com/golang-jwt/jwt/v5`を使用
- 署名アルゴリズム: HS256（HMAC-SHA256）
- 秘密鍵は環境別に分離

### 7.2 秘密鍵の管理
- 秘密鍵はコマンドラインツール（`generate-secret`）で生成
- 生成された秘密鍵を設定ファイルに手動で記述
- 複数のWebサーバーで動作する場合、同じ秘密鍵を使用

### 7.3 認証ミドルウェアの適用順序
- CORSミドルウェア → 認証ミドルウェア → ハンドラー
- `/api/*`パスにのみ適用
- `/health`エンドポイントは認証不要

### 7.4 GoAdminカスタムページの実装
- 既存のカスタムページ実装パターン（`server/internal/admin/pages/`）に従う
- `RegisterCustomPages`関数に`/api-key`パスを追加
  - 注意: `RegisterCustomPages`で`/api-key`と登録すると、GoAdminが自動的に`/admin`プレフィックスを追加し、実際のURLは`/admin/api-key`になる
  - HTML内のリンク（form action、a href等）も`/admin/api-key`とする必要がある
- メニュー項目はデータベースマイグレーションで追加

### 7.5 クライアント側の実装
- 環境変数（`NEXT_PUBLIC_API_KEY`）からAPIキーを取得
- APIキーが設定されていない場合、コンストラクタでエラーを投げる
- すべてのAPIリクエストに`Authorization: Bearer <API_KEY>`ヘッダーを付与

### 7.6 テスト用ダミーAPIキー
- `server/internal/config/testdata/develop/api_key.yaml`に配置
- テスト用の秘密鍵を使用してJWTを生成
- テスト実行時は`testdata/`の設定ファイルを使用

## 8. セキュリティ考慮事項

### 8.1 秘密鍵の管理
- 秘密鍵は環境別に分離（`config/{env}/`に保存）
- staging/productionの秘密鍵は`.gitignore`でcommit不可
- 秘密鍵の漏洩を防ぐため、適切な権限管理を実施

### 8.2 JWT署名の検証
- JWT署名の検証を必須とする
- 秘密鍵による検証を毎回実施（キャッシュしない）

### 8.3 無効化機能
- versionベースの無効化により、脆弱性発見時に迅速に対応可能
- 設定ファイルの`invalid_versions`リストを更新するだけで無効化可能

### 8.4 環境分離
- 環境（develop/staging/production）ごとに異なる秘密鍵を使用
- JWTの`env`フィールドで環境を検証

## 9. パフォーマンス考慮事項

### 9.1 認証ミドルウェアのオーバーヘッド
- JWT検証処理は高速に実行（毎回検証、キャッシュなし）
- 認証処理がAPIレスポンス時間に大きな影響を与えない

### 9.2 スコープ検証の効率化
- スコープ配列の線形探索（固定サイズなので問題なし）
- 必要に応じてマップを使用して高速化可能

## 10. 参考情報

### 10.1 既存実装
- `server/internal/api/router/router.go`: ルーター実装
- `server/internal/logging/middleware.go`: ミドルウェア実装例
- `server/internal/admin/pages/pages.go`: カスタムページ実装例
- `client/src/lib/api.ts`: TypeScriptクライアント実装

### 10.2 JWTライブラリ
- **`github.com/golang-jwt/jwt/v5`**: Go言語用JWTライブラリ
- **公式ドキュメント**: https://github.com/golang-jwt/jwt
- **署名アルゴリズム**: HS256（HMAC-SHA256）

### 10.3 Go言語標準ライブラリ
- `net/http`パッケージ: HTTPミドルウェア実装
- `crypto/rand`パッケージ: ランダムな秘密鍵生成
- `encoding/base64`パッケージ: Base64エンコード

