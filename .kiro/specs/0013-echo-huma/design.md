# Echo・Huma導入設計書

## 1. 概要

### 1.1 設計の目的
要件定義書に基づき、既存のGorilla MuxベースのAPIサーバーをEchoとHumaフレームワークに置き換えるための詳細設計を定義する。既存のService層とRepository層は変更せず、Handler層とRouter層のみを変更することで、既存機能の互換性を維持しつつ、自動バリデーションとOpenAPI仕様の自動生成を実現する。

### 1.2 設計の範囲
- Echoフレームワークの導入と設定
- Humaフレームワークの導入と設定（humaechoアダプター）
- 既存エンドポイントのHuma化
- リクエスト/レスポンス構造体のHumaタグ対応
- 認証ミドルウェアのEcho形式への変換
- CORS設定のEcho形式への移行
- アクセスログのEcho形式への統合
- OpenAPI仕様の自動生成

### 1.3 設計方針
- **既存アーキテクチャの維持**: Service層とRepository層は変更しない
- **段階的な移行**: 既存の実装を維持しつつ、新しい実装に置き換える
- **後方互換性の維持**: 既存のAPIエンドポイントの動作とリクエスト/レスポンス形式を維持
- **ライブラリの知見活用**: EchoとHumaのベストプラクティスを活用
- **型安全性の向上**: Humaの型システムを活用してコンパイル時にエラーを検出

## 2. アーキテクチャ設計

### 2.1 ディレクトリ構造

#### 2.1.1 変更前の構造
```
server/
└── internal/
    ├── api/
    │   ├── router/
    │   │   └── router.go          # Gorilla Muxルーター
    │   └── handler/
    │       ├── user_handler.go     # 標準http.Handler
    │       └── post_handler.go     # 標準http.Handler
    ├── auth/
    │   └── middleware.go          # 標準http.Handlerミドルウェア
    └── logging/
        └── middleware.go          # 標準http.Handlerミドルウェア
```

#### 2.1.2 変更後の構造
```
server/
└── internal/
    ├── api/
    │   ├── router/
    │   │   └── router.go          # Echo/Humaルーター
    │   ├── handler/
    │   │   ├── user_handler.go     # Humaハンドラー
    │   │   └── post_handler.go   # Humaハンドラー
    │   └── huma/                  # 新規（オプション）
    │       ├── inputs.go          # リクエスト構造体定義
    │       └── outputs.go         # レスポンス構造体定義
    ├── auth/
    │   └── middleware.go          # Echoミドルウェア形式
    └── logging/
        └── middleware.go          # Echoミドルウェア形式（またはEcho標準ミドルウェアを使用）
```

### 2.2 リクエスト処理フロー

#### 2.2.1 変更前のフロー
```
HTTPリクエスト
    ↓
Gorilla Muxルーター
    ↓
認証ミドルウェア（標準http.Handler）
    ↓
CORSミドルウェア（github.com/rs/cors）
    ↓
アクセスログミドルウェア（標準http.Handler）
    ↓
ハンドラー（標準http.Handler）
    ├─ 手動でJSONデコード
    ├─ 手動でバリデーション
    ├─ Service層呼び出し
    ├─ 手動でエラーハンドリング
    └─ 手動でJSONエンコード
    ↓
HTTPレスポンス
```

#### 2.2.2 変更後のフロー
```
HTTPリクエスト
    ↓
Echoルーター
    ↓
認証ミドルウェア（Echo形式）
    ↓
CORSミドルウェア（Echo標準）
    ↓
アクセスログミドルウェア（Echo形式）
    ↓
Humaエンドポイント（huma.Register）
    ├─ 自動でJSONデコード
    ├─ 自動でバリデーション（Humaタグ）
    ├─ ハンドラー関数実行
    │   ├─ Service層呼び出し
    │   └─ エラー返却（Humaが自動処理）
    └─ 自動でJSONエンコード
    ↓
HTTPレスポンス
```

### 2.3 既存アーキテクチャとの統合

#### 2.3.1 Service層との統合
- Service層は変更しない
- Humaハンドラーから既存のService層メソッドを呼び出す
- Service層のエラーはHumaのエラーレスポンス形式に変換

#### 2.3.2 Repository層との統合
- Repository層は変更しない
- Service層経由でRepository層にアクセス

#### 2.3.3 Model層との統合
- 既存のModel構造体は維持
- Huma用のInput/Output構造体を新規作成（既存Modelを埋め込む）

## 3. コンポーネント設計

### 3.1 Echoフレームワークの設定

#### 3.1.1 Echoインスタンスの作成
```go
// server/cmd/server/main.go

package main

import (
    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
)

func main() {
    // Echoインスタンスの作成
    e := echo.New()
    
    // デバッグモードの設定（開発環境のみ）
    if os.Getenv("APP_ENV") == "develop" {
        e.Debug = true
    }
    
    // リカバリーミドルウェア
    e.Use(middleware.Recover())
    
    // ... その他の設定
}
```

#### 3.1.2 ミドルウェアの設定順序
1. Recover（最外側）
2. Logger（アクセスログ）
3. CORS
4. 認証ミドルウェア（APIエンドポイントのみ）
5. Humaエンドポイント

### 3.2 Humaフレームワークの設定

#### 3.2.1 Huma APIインスタンスの作成
```go
// server/internal/api/router/router.go

package router

import (
    "net/http"
    
    "github.com/danielgtaylor/huma/v2"
    "github.com/danielgtaylor/huma/v2/adapters/humaecho"
    "github.com/labstack/echo/v4"
)

func NewRouter(e *echo.Echo, userService *service.UserService, postService *service.PostService, cfg *config.Config) {
    // Huma設定
    config := huma.DefaultConfig("go-webdb-template API", "1.0.0")
    config.DocsPath = "/openapi.json"
    config.Servers = []*huma.Server{
        {
            URL:         fmt.Sprintf("http://localhost:%d", cfg.Server.Port),
            Description: "Development server",
        },
    }
    
    // Huma APIインスタンスの作成
    api := humaecho.New(e, config)
    
    // エンドポイントの登録
    registerUserEndpoints(api, userService)
    registerPostEndpoints(api, postService)
}
```

### 3.3 リクエスト/レスポンス構造体の設計

#### 3.3.1 ユーザーエンドポイントの構造体

**CreateUserInput**:
```go
// server/internal/api/huma/inputs.go

package huma

type CreateUserInput struct {
    Body struct {
        Name  string `json:"name" maxLength:"100" doc:"ユーザー名"`
        Email string `json:"email" format:"email" maxLength:"255" doc:"メールアドレス"`
    } `json:"body"`
}
```

**GetUserInput**:
```go
type GetUserInput struct {
    Path struct {
        ID int64 `path:"id" doc:"ユーザーID"`
    } `json:"path"`
}
```

**UpdateUserInput**:
```go
type UpdateUserInput struct {
    Path struct {
        ID int64 `path:"id" doc:"ユーザーID"`
    } `json:"path"`
    Body struct {
        Name  string `json:"name,omitempty" maxLength:"100" doc:"ユーザー名"`
        Email string `json:"email,omitempty" format:"email" maxLength:"255" doc:"メールアドレス"`
    } `json:"body"`
}
```

**DeleteUserInput**:
```go
type DeleteUserInput struct {
    Path struct {
        ID int64 `path:"id" doc:"ユーザーID"`
    } `json:"path"`
}
```

**ListUsersInput**:
```go
type ListUsersInput struct {
    Query struct {
        Limit  int `query:"limit" default:"20" minimum:"1" maximum:"100" doc:"取得件数"`
        Offset int `query:"offset" default:"0" minimum:"0" doc:"オフセット"`
    } `json:"query"`
}
```

**UserOutput**:
```go
// server/internal/api/huma/outputs.go

package huma

import "github.com/example/go-webdb-template/internal/model"

type UserOutput struct {
    Body struct {
        model.User
    } `json:"body"`
}
```

**UsersOutput**:
```go
type UsersOutput struct {
    Body struct {
        Users []*model.User `json:"users" doc:"ユーザー一覧"`
    } `json:"body"`
}
```

#### 3.3.2 投稿エンドポイントの構造体

**CreatePostInput**:
```go
type CreatePostInput struct {
    Body struct {
        UserID  int64  `json:"user_id" minimum:"1" doc:"ユーザーID"`
        Title   string `json:"title" maxLength:"200" doc:"タイトル"`
        Content string `json:"content" doc:"内容"`
    } `json:"body"`
}
```

**GetPostInput**:
```go
type GetPostInput struct {
    Path struct {
        ID int64 `path:"id" doc:"投稿ID"`
    } `json:"path"`
    Query struct {
        UserID int64 `query:"user_id" required:"true" minimum:"1" doc:"ユーザーID"`
    } `json:"query"`
}
```

**UpdatePostInput**:
```go
type UpdatePostInput struct {
    Path struct {
        ID int64 `path:"id" doc:"投稿ID"`
    } `json:"path"`
    Query struct {
        UserID int64 `query:"user_id" required:"true" minimum:"1" doc:"ユーザーID"`
    } `json:"query"`
    Body struct {
        Title   string `json:"title,omitempty" maxLength:"200" doc:"タイトル"`
        Content string `json:"content,omitempty" doc:"内容"`
    } `json:"body"`
}
```

**DeletePostInput**:
```go
type DeletePostInput struct {
    Path struct {
        ID int64 `path:"id" doc:"投稿ID"`
    } `json:"path"`
    Query struct {
        UserID int64 `query:"user_id" required:"true" minimum:"1" doc:"ユーザーID"`
    } `json:"query"`
}
```

**ListPostsInput**:
```go
type ListPostsInput struct {
    Query struct {
        Limit  int    `query:"limit" default:"20" minimum:"1" maximum:"100" doc:"取得件数"`
        Offset int    `query:"offset" default:"0" minimum:"0" doc:"オフセット"`
        UserID *int64 `query:"user_id,omitempty" minimum:"1" doc:"ユーザーID（オプション）"`
    } `json:"query"`
}
```

**GetUserPostsInput**:
```go
type GetUserPostsInput struct {
    Query struct {
        Limit  int `query:"limit" default:"20" minimum:"1" maximum:"100" doc:"取得件数"`
        Offset int `query:"offset" default:"0" minimum:"0" doc:"オフセット"`
    } `json:"query"`
}
```

**PostOutput**:
```go
type PostOutput struct {
    Body struct {
        model.Post
    } `json:"body"`
}
```

**PostsOutput**:
```go
type PostsOutput struct {
    Body struct {
        Posts []*model.Post `json:"posts" doc:"投稿一覧"`
    } `json:"body"`
}
```

**UserPostsOutput**:
```go
type UserPostsOutput struct {
    Body struct {
        UserPosts []*model.UserPost `json:"user_posts" doc:"ユーザーと投稿のJOIN結果"`
    } `json:"body"`
}
```

### 3.4 ハンドラー関数の設計

#### 3.4.1 ユーザーエンドポイントのハンドラー

**CreateUser**:
```go
// server/internal/api/handler/user_handler.go

package handler

import (
    "context"
    "net/http"
    
    "github.com/danielgtaylor/huma/v2"
    "github.com/example/go-webdb-template/internal/api/huma"
    "github.com/example/go-webdb-template/internal/service"
)

func registerCreateUser(api huma.API, userService *service.UserService) {
    huma.Register(api, huma.Operation{
        Method: http.MethodPost,
        Path:   "/api/users",
        Summary: "ユーザーを作成",
    }, func(ctx context.Context, input *humaapi.CreateUserInput) (*humaapi.UserOutput, error) {
        // Service層の呼び出し
        req := &model.CreateUserRequest{
            Name:  input.Body.Name,
            Email: input.Body.Email,
        }
        
        user, err := userService.CreateUser(ctx, req)
        if err != nil {
            return nil, huma.Error400BadRequest(err.Error())
        }
        
        return &humaapi.UserOutput{
            Body: struct {
                model.User
            }{
                User: *user,
            },
        }, nil
    })
}
```

**GetUser**:
```go
func registerGetUser(api huma.API, userService *service.UserService) {
    huma.Register(api, huma.Operation{
        Method: http.MethodGet,
        Path:   "/api/users/{id}",
        Summary: "ユーザーを取得",
    }, func(ctx context.Context, input *humaapi.GetUserInput) (*humaapi.UserOutput, error) {
        user, err := userService.GetUser(ctx, input.Path.ID)
        if err != nil {
            return nil, huma.Error404NotFound(err.Error())
        }
        
        return &humaapi.UserOutput{
            Body: struct {
                model.User
            }{
                User: *user,
            },
        }, nil
    })
}
```

**ListUsers**:
```go
func registerListUsers(api huma.API, userService *service.UserService) {
    huma.Register(api, huma.Operation{
        Method: http.MethodGet,
        Path:   "/api/users",
        Summary: "ユーザー一覧を取得",
    }, func(ctx context.Context, input *humaapi.ListUsersInput) (*humaapi.UsersOutput, error) {
        users, err := userService.ListUsers(ctx, input.Query.Limit, input.Query.Offset)
        if err != nil {
            return nil, huma.Error500InternalServerError(err.Error())
        }
        
        return &humaapi.UsersOutput{
            Body: struct {
                Users []*model.User `json:"users"`
            }{
                Users: users,
            },
        }, nil
    })
}
```

**UpdateUser**:
```go
func registerUpdateUser(api huma.API, userService *service.UserService) {
    huma.Register(api, huma.Operation{
        Method: http.MethodPut,
        Path:   "/api/users/{id}",
        Summary: "ユーザーを更新",
    }, func(ctx context.Context, input *humaapi.UpdateUserInput) (*humaapi.UserOutput, error) {
        req := &model.UpdateUserRequest{
            Name:  input.Body.Name,
            Email: input.Body.Email,
        }
        
        user, err := userService.UpdateUser(ctx, input.Path.ID, req)
        if err != nil {
            return nil, huma.Error500InternalServerError(err.Error())
        }
        
        return &humaapi.UserOutput{
            Body: struct {
                model.User
            }{
                User: *user,
            },
        }, nil
    })
}
```

**DeleteUser**:
```go
func registerDeleteUser(api huma.API, userService *service.UserService) {
    huma.Register(api, huma.Operation{
        Method: http.MethodDelete,
        Path:   "/api/users/{id}",
        Summary: "ユーザーを削除",
    }, func(ctx context.Context, input *humaapi.DeleteUserInput) (*humaapi.DeleteUserOutput, error) {
        err := userService.DeleteUser(ctx, input.Path.ID)
        if err != nil {
            return nil, huma.Error500InternalServerError(err.Error())
        }
        
        return &humaapi.DeleteUserOutput{
            Status: http.StatusNoContent,
        }, nil
    })
}
```

#### 3.4.2 投稿エンドポイントのハンドラー

同様のパターンで投稿エンドポイントも実装する。

### 3.5 認証ミドルウェアの設計

#### 3.5.1 Echoミドルウェア形式への変換
```go
// server/internal/auth/middleware.go

package auth

import (
    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
)

// NewAuthMiddleware は新しい認証ミドルウェアを作成
func NewAuthMiddleware(cfg *config.APIConfig, env string) echo.MiddlewareFunc {
    validator := NewJWTValidator(cfg, env)
    
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            // AuthorizationヘッダーからJWTトークンを取得
            authHeader := c.Request().Header.Get("Authorization")
            if authHeader == "" {
                return c.JSON(http.StatusUnauthorized, map[string]interface{}{
                    "code":    http.StatusUnauthorized,
                    "message": "Authorization header is required",
                })
            }
            
            // Bearerトークンの抽出
            parts := strings.Split(authHeader, " ")
            if len(parts) != 2 || parts[0] != "Bearer" {
                return c.JSON(http.StatusUnauthorized, map[string]interface{}{
                    "code":    http.StatusUnauthorized,
                    "message": "Invalid authorization header format",
                })
            }
            
            tokenString := parts[1]
            
            // JWT検証
            claims, err := validator.ValidateJWT(tokenString)
            if err != nil {
                return c.JSON(http.StatusUnauthorized, map[string]interface{}{
                    "code":    http.StatusUnauthorized,
                    "message": "Invalid API key",
                })
            }
            
            // スコープ検証
            if err := validateScope(claims, c.Request().Method); err != nil {
                return c.JSON(http.StatusForbidden, map[string]interface{}{
                    "code":    http.StatusForbidden,
                    "message": "Insufficient scope",
                })
            }
            
            // 次のハンドラーを実行
            return next(c)
        }
    }
}
```

### 3.6 CORS設定の設計

#### 3.6.1 Echo CORSミドルウェアの設定
```go
// server/internal/api/router/router.go

func setupMiddleware(e *echo.Echo, cfg *config.Config) {
    // CORS設定
    e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
        AllowOrigins:     cfg.CORS.AllowedOrigins,
        AllowMethods:     cfg.CORS.AllowedMethods,
        AllowHeaders:     cfg.CORS.AllowedHeaders,
        AllowCredentials: true,
    }))
}
```

### 3.7 アクセスログの設計

#### 3.7.1 Echoログミドルウェアの設定
既存のアクセスログ機能をEchoのログミドルウェアに統合するか、既存のログ機能をEchoミドルウェア形式に変換する。

**オプション1: Echo標準ログミドルウェアを使用**
```go
e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
    Format: `${time_rfc3339} ${status} ${method} ${uri} ${latency_human} ${bytes_in} ${bytes_out}`,
    Output: accessLogger.Writer(),
}))
```

**オプション2: 既存のログ機能をEcho形式に変換**
```go
// server/internal/logging/middleware.go

func NewAccessLogMiddleware(accessLogger *AccessLogger) echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            startTime := time.Now()
            
            // リクエスト情報の記録
            // ...
            
            // 次のハンドラーを実行
            err := next(c)
            
            // レスポンス情報の記録
            // ...
            
            return err
        }
    }
}
```

### 3.8 ヘルスチェックエンドポイントの設計

#### 3.8.1 Echoハンドラーとして実装
```go
// server/internal/api/router/router.go

func registerHealthCheck(e *echo.Echo) {
    e.GET("/health", func(c echo.Context) error {
        return c.String(http.StatusOK, "OK")
    })
}
```

## 4. エラーハンドリング設計

### 4.1 Humaのエラーレスポンス形式

Humaは標準的なエラーレスポンス形式を提供：

```go
// 400 Bad Request
return nil, huma.Error400BadRequest("Invalid request")

// 404 Not Found
return nil, huma.Error404NotFound("Resource not found")

// 500 Internal Server Error
return nil, huma.Error500InternalServerError("Internal server error")
```

### 4.2 Service層エラーの変換

Service層から返されるエラーをHumaのエラーレスポンス形式に変換：

```go
func convertServiceError(err error) error {
    if err == nil {
        return nil
    }
    
    // エラーメッセージから適切なHTTPステータスコードを判定
    errMsg := err.Error()
    
    if strings.Contains(errMsg, "not found") {
        return huma.Error404NotFound(errMsg)
    }
    
    if strings.Contains(errMsg, "invalid") || strings.Contains(errMsg, "required") {
        return huma.Error400BadRequest(errMsg)
    }
    
    // デフォルトは500エラー
    return huma.Error500InternalServerError(errMsg)
}
```

### 4.3 バリデーションエラー

Humaは自動的にバリデーションエラーを処理し、適切なエラーレスポンスを返す。

## 5. OpenAPI仕様の自動生成

### 5.1 OpenAPIエンドポイント

Humaは自動的にOpenAPI仕様を生成し、以下のエンドポイントでアクセス可能：

- `GET /openapi.json`: OpenAPI 3.0仕様（JSON形式）
- `GET /openapi.yaml`: OpenAPI 3.0仕様（YAML形式）

### 5.2 仕様の内容

- 全エンドポイントの定義
- リクエスト/レスポンススキーマ
- バリデーションルール
- エラーレスポンス定義
- 認証方式の定義（必要に応じて）

## 6. テスト戦略

### 6.1 単体テスト

#### 6.1.1 ハンドラーのテスト
- Humaハンドラー関数の単体テスト
- リクエスト/レスポンス構造体のテスト
- エラーハンドリングのテスト

#### 6.1.2 ミドルウェアのテスト
- 認証ミドルウェアのテスト
- CORSミドルウェアのテスト
- アクセスログミドルウェアのテスト

### 6.2 統合テスト

#### 6.2.1 E2Eテスト
- 既存のE2Eテストを更新
- Echo/Humaのテストユーティリティを使用
- 全エンドポイントの動作確認

### 6.3 OpenAPI仕様の検証

#### 6.3.1 仕様の検証
- OpenAPI仕様の妥当性検証
- 既存のAPI仕様との整合性確認

## 7. 移行計画

### 7.1 段階的な移行

1. **Phase 1: Echoの導入**
   - Echoフレームワークの導入
   - 基本的なルーティングの移行
   - ミドルウェアの移行

2. **Phase 2: Humaの導入**
   - Humaフレームワークの導入
   - リクエスト/レスポンス構造体の定義
   - エンドポイントのHuma化

3. **Phase 3: テストと検証**
   - 既存テストの更新
   - E2Eテストの実行
   - OpenAPI仕様の検証

### 7.2 後方互換性の維持

- 既存のAPIエンドポイントの動作を維持
- 既存のリクエスト/レスポンス形式を維持
- 既存のテストコードを可能な限り維持

## 8. パフォーマンス考慮事項

### 8.1 オーバーヘッド

- EchoとHumaのオーバーヘッドは最小限
- 既存のパフォーマンスを維持

### 8.2 最適化

- 必要に応じて、Echoの設定を最適化
- Humaの設定を最適化

## 9. セキュリティ考慮事項

### 9.1 認証の維持

- 既存のJWT認証方式を維持
- 認証ミドルウェアのセキュリティを維持

### 9.2 バリデーション

- Humaの自動バリデーションにより、セキュリティが向上
- 入力値の検証が自動化される

## 10. 参考情報

### 10.1 関連ドキュメント
- 要件定義書: `.kiro/specs/0013-echo-huma/requirements.md`
- Issue #24: APIサーバーのリクエストの処理にEcho、および、Humaを導入

### 10.2 技術スタック
- **Go**: 1.23.4
- **Echo**: v4
- **Huma**: v2
- **humaecho**: v2

### 10.3 参考リンク
- Echo公式ドキュメント: https://echo.labstack.com/
- Huma公式ドキュメント: https://huma.rocks/
- Huma Echoアダプター: https://pkg.go.dev/github.com/danielgtaylor/huma/v2/adapters/humaecho

