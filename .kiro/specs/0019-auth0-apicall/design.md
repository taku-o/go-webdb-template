# Auth0 API呼び出し機能設計書

## 1. 概要

### 1.1 設計の目的
要件定義書に基づき、Auth0から受け取ったJWTをAPIサーバーとの通信で利用する機能を実装する詳細設計を定義する。クライアント側でAuth0 JWTとPublic API Keyの切り替えを実装し、サーバー側でJWT種類の判別、Auth0 JWT検証、API公開レベルの実装を行う。

### 1.2 設計の範囲
- クライアント側: Auth0 JWTとPublic API Keyの切り替え、API呼び出し
- サーバー側: JWT種類の判別、Auth0 JWT検証、Public API Key JWT検証、API公開レベルの実装
- 新規private APIエンドポイントの追加（`/api/today`）

**本設計の範囲外**:
- Auth0ログイン機能の実装（Issue #30で対応済み）
- アカウント情報のデータベース保存（別issueで対応）

### 1.3 設計方針
- **ライブラリの積極的利用**: `github.com/MicahParks/keyfunc`を使用してJWKS取得とキャッシュを実装
- **既存実装の維持**: Public API Key JWT検証機能は変更せず、Auth0 JWT検証を追加
- **API公開レベルの柔軟な定義**: `huma.Operation`に直接追加できる場合はコード内定義、困難な場合はマップで管理
- **後方互換性の維持**: 既存のAPIエンドポイントの動作は変更しない（公開レベルのみ追加）
- **シンプルな実装**: private/public判定の実装例として`/api/today`エンドポイントを追加

## 2. アーキテクチャ設計

### 2.1 ディレクトリ構造

#### 2.1.1 変更前の構造
```
client/
├── src/
│   ├── app/
│   │   ├── page.tsx
│   │   └── ...
│   ├── lib/
│   │   └── api.ts
│   └── types/
server/
├── internal/
│   ├── auth/
│   │   ├── jwt.go
│   │   └── middleware.go
│   ├── api/
│   │   ├── handler/
│   │   │   ├── user_handler.go
│   │   │   └── post_handler.go
│   │   └── router/
│   │       └── router.go
│   └── config/
│       └── config.go
config/
└── {env}/
    └── config.yaml
```

#### 2.1.2 変更後の構造
```
client/
├── src/
│   ├── app/
│   │   ├── page.tsx                      # 修正: TodayApiButtonの追加
│   │   └── ...
│   ├── components/
│   │   └── TodayApiButton.tsx            # 新規: private API呼び出しボタン
│   ├── lib/
│   │   └── api.ts                        # 修正: JWT取得ロジックの追加
│   └── types/
server/
├── internal/
│   ├── auth/
│   │   ├── jwt.go                        # 修正: JWT種類の判別機能の追加
│   │   ├── middleware.go                 # 修正: アクセス制御機能の追加
│   │   └── auth0_validator.go            # 新規: Auth0 JWT検証機能
│   ├── api/
│   │   ├── handler/
│   │   │   ├── user_handler.go           # 修正: 公開レベル定義（public）
│   │   │   ├── post_handler.go            # 修正: 公開レベル定義（public）
│   │   │   └── today_handler.go          # 新規: /api/todayエンドポイント
│   │   └── router/
│   │       └── router.go
│   └── config/
│       └── config.go                     # 修正: AUTH0_ISSUER_BASE_URL設定の追加
config/
└── {env}/
    └── config.yaml                       # 修正: AUTH0_ISSUER_BASE_URLの追加
```

### 2.2 ファイル構成

#### 2.2.1 クライアント側（Next.js）

**`client/src/lib/api.ts`**: APIクライアントクラス
- JWT取得ロジックの追加（Auth0 JWTとPublic API Keyの切り替え）
- `useUser()`フックを使用してログイン状態を確認
- ログイン中は`getAccessToken()`でJWTを取得
- 未ログイン時は`NEXT_PUBLIC_API_KEY`を使用

**`client/src/components/TodayApiButton.tsx`**: private API呼び出しボタンコンポーネント
- `/api/today`エンドポイントを呼び出す
- エラーハンドリングとエラーメッセージの表示
- private/public判定の実装例として機能

**`client/src/app/page.tsx`**: トップページ
- `TodayApiButton`コンポーネントの追加

#### 2.2.2 サーバー側（Go）

**`server/internal/auth/auth0_validator.go`**: Auth0 JWT検証機能
- `github.com/MicahParks/keyfunc`ライブラリを使用
- JWKS取得とキャッシュ（12時間）
- Auth0 JWTの検証（RS256署名）

**`server/internal/auth/jwt.go`**: JWT検証機能
- JWT種類の判別機能の追加（`iss`による判別）
- Public API Key JWT検証機能（既存機能の維持）

**`server/internal/auth/middleware.go`**: 認証ミドルウェア
- アクセス制御機能の追加（API公開レベルのチェック）
- Public API Key JWTでprivateなAPIにアクセスした場合、403 Forbiddenを返す

**`server/internal/api/handler/today_handler.go`**: 新規private APIのハンドラー
- `GET /api/today`エンドポイントの実装
- 今日の日付をYYYY-MM-DD形式で返す

**`server/internal/api/handler/user_handler.go`**: ユーザーAPIハンドラー
- 既存APIの公開レベル定義（public）

**`server/internal/api/handler/post_handler.go`**: 投稿APIハンドラー
- 既存APIの公開レベル定義（public）

**`server/internal/config/config.go`**: 設定管理
- `AUTH0_ISSUER_BASE_URL`設定の追加

#### 2.2.3 設定ファイル

**`config/{env}/config.yaml`**: 環境別設定ファイル
- `AUTH0_ISSUER_BASE_URL`の追加（公開情報のため設定ファイルで管理）
  - develop: `https://dev-oaa5vtzmld4dsxtd.jp.auth0.com`
  - staging: （設定値は環境に応じて決定）
  - production: （設定値は環境に応じて決定）

### 2.3 システム構成図

```
┌─────────────────────────────────────────────────────────┐
│                    クライアント（Next.js）                │
│                                                          │
│  ┌──────────────────────────────────────────────────┐  │
│  │  client/src/lib/api.ts                           │  │
│  │  - useUser()でログイン状態を確認                 │  │
│  │  - ログイン中: getAccessToken()でJWT取得         │  │
│  │  - 未ログイン: NEXT_PUBLIC_API_KEYを使用         │  │
│  └──────────────────┬───────────────────────────────┘  │
│                     │                                    │
│                     │ Authorization: Bearer <JWT>       │
│                     │                                    │
└─────────────────────┼────────────────────────────────────┘
                      │
                      ▼
┌─────────────────────────────────────────────────────────┐
│              サーバー（Go）                              │
│                                                          │
│  ┌──────────────────────────────────────────────────┐  │
│  │  server/internal/auth/middleware.go              │  │
│  │  - AuthorizationヘッダーからJWTを取得            │  │
│  │  - JWT種類の判別（issによる判別）                │  │
│  └──────────────────┬───────────────────────────────┘  │
│                     │                                    │
│         ┌───────────┴───────────┐                       │
│         │                       │                       │
│         ▼                       ▼                       │
│  ┌──────────────┐      ┌──────────────┐                │
│  │ Auth0 JWT    │      │ Public API   │                │
│  │ 検証         │      │ Key JWT検証  │                │
│  │              │      │              │                │
│  │ auth0_      │      │ jwt.go       │                │
│  │ validator.go │      │ (既存)       │                │
│  └──────┬───────┘      └──────┬───────┘                │
│         │                     │                        │
│         └───────────┬─────────┘                        │
│                     │                                    │
│                     ▼                                    │
│  ┌──────────────────────────────────────────────────┐  │
│  │  API公開レベルのチェック                         │  │
│  │  - Public API Key → publicなAPIのみ許可          │  │
│  │  - Auth0 JWT → publicとprivateの両方許可        │  │
│  └──────────────────┬───────────────────────────────┘  │
│                     │                                    │
│                     ▼                                    │
│  ┌──────────────────────────────────────────────────┐  │
│  │  APIハンドラー                                   │  │
│  │  - /api/users (public)                          │  │
│  │  - /api/posts (public)                          │  │
│  │  - /api/user-posts (public)                     │  │
│  │  - /api/today (private)                         │  │
│  └──────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────┘
```

### 2.4 データフロー

#### 2.4.1 API呼び出しフロー（Auth0 JWT使用時）
```
クライアント: useUser()でログイン状態を確認
    ↓
クライアント: getAccessToken()でJWTを取得
    ↓
クライアント: Authorization: Bearer <Auth0 JWT>ヘッダーでAPIリクエスト送信
    ↓
サーバー: ミドルウェアでJWTを取得
    ↓
サーバー: JWT種類の判別（issがAuth0ドメイン）
    ↓
サーバー: Auth0 JWT検証（keyfuncライブラリでJWKS取得・検証）
    ↓
サーバー: JWTの許容する公開レベルを判定（"private" = public/private両方許可）
    ↓
サーバー: 許容する公開レベルをコンテキストに設定
    ↓
サーバー: APIハンドラーを実行（エンドポイントの公開レベルと比較）
    ↓
サーバー: レスポンスを返す
```

#### 2.4.2 API呼び出しフロー（Public API Key使用時）
```
クライアント: useUser()でログイン状態を確認（未ログイン）
    ↓
クライアント: NEXT_PUBLIC_API_KEYを使用
    ↓
クライアント: Authorization: Bearer <Public API Key JWT>ヘッダーでAPIリクエスト送信
    ↓
サーバー: ミドルウェアでJWTを取得
    ↓
サーバー: JWT種類の判別（issが"go-webdb-template"）
    ↓
サーバー: Public API Key JWT検証（既存のjwt.goで検証）
    ↓
サーバー: JWTの許容する公開レベルを判定（"public" = publicのみ許可）
    ↓
サーバー: 許容する公開レベルをコンテキストに設定
    ↓
サーバー: APIハンドラーを実行（エンドポイントの公開レベルと比較）
    ↓
    ├─ publicなAPI → 許可
    └─ privateなAPI → 403 Forbiddenを返す
    ↓
サーバー: レスポンスを返す
```

#### 2.4.3 JWKS取得・キャッシュフロー
```
サーバー起動時
    ↓
Auth0Validatorの初期化
    ↓
keyfunc.New()でJWKS取得開始
    ↓
JWKSをメモリキャッシュに保存
    ↓
12時間ごとに定期更新（RefreshInterval）
    ↓
未知のKIDが来たら再取得（RefreshUnknownKID）
    ↓
JWT検証時にキャッシュからJWKSを取得
```

## 3. コンポーネント設計

### 3.1 クライアント側コンポーネント

#### 3.1.1 ApiClientクラスの修正

**`client/src/lib/api.ts`**:
```typescript
import { User, CreateUserRequest, UpdateUserRequest } from '@/types/user'
import { Post, CreatePostRequest, UpdatePostRequest, UserPost } from '@/types/post'
import { useUser } from '@auth0/nextjs-auth0/client'

const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8080'

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

  // JWTを取得（Auth0 JWTまたはPublic API Key）
  private async getJWT(): Promise<string> {
    // Client ComponentsでのみuseUser()を使用可能
    // Server ComponentsではgetAccessToken()を使用
    // ここではClient Components用の実装を想定
    
    // 注意: useUser()はReactフックのため、コンポーネント内でのみ使用可能
    // ApiClientクラス内では直接使用できない
    // そのため、getJWT()メソッドは呼び出し側からJWTを渡す方式に変更する必要がある
    
    // 実装方法の選択肢:
    // 1. ApiClientのrequest()メソッドにJWTを引数として渡す
    // 2. ApiClientのインスタンス作成時にJWT取得関数を渡す
    // 3. ApiClientを関数として実装し、useUser()を使用できるようにする
    
    // 要件定義書の意図を考慮すると、方法2が適切
    // ただし、既存のapiClientインスタンスとの互換性を考慮する必要がある
    
    // 暫定的な実装: 常にNEXT_PUBLIC_API_KEYを使用（後で修正）
    return this.apiKey!
  }

  private async request<T>(
    endpoint: string,
    options?: RequestInit,
    jwt?: string
  ): Promise<T> {
    const url = `${this.baseURL}${endpoint}`

    // JWTの取得（引数で渡された場合はそれを使用、なければgetJWT()で取得）
    const token = jwt || await this.getJWT()

    // Authorizationヘッダーを追加
    const headers = {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`,
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

  // 既存のメソッドは変更なし（互換性の維持）
  // ただし、内部的にrequest()メソッドを使用するため、JWT取得ロジックが適用される
  async getUsers(limit = 20, offset = 0): Promise<User[]> {
    return this.request<User[]>(`/api/users?limit=${limit}&offset=${offset}`)
  }

  // ... 他の既存メソッド ...
}
```

**実装上の注意**: 
- `useUser()`はReactフックのため、クラスメソッド内では直接使用できない
- 実装方法の選択肢を検討し、最適な方法を採用する必要がある
- 既存の`apiClient`インスタンスとの互換性を維持する

#### 3.1.2 TodayApiButtonコンポーネント

**`client/src/components/TodayApiButton.tsx`**:
```typescript
'use client'

import { useUser } from '@auth0/nextjs-auth0/client'
import { useState } from 'react'

export default function TodayApiButton() {
  const { user, getAccessToken } = useUser()
  const [date, setDate] = useState<string | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)

  const handleClick = async () => {
    setLoading(true)
    setError(null)
    setDate(null)

    try {
      // JWTの取得
      let token: string
      if (user) {
        // ログイン中: Auth0 JWTを取得
        token = await getAccessToken()
      } else {
        // 未ログイン: Public API Keyを使用
        token = process.env.NEXT_PUBLIC_API_KEY!
      }

      // API呼び出し
      const response = await fetch(`${process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8080'}/api/today`, {
        headers: {
          'Authorization': `Bearer ${token}`,
        },
      })

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}))
        throw new Error(errorData.message || response.statusText)
      }

      const data = await response.json()
      setDate(data.date)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'エラーが発生しました')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div>
      <button onClick={handleClick} disabled={loading}>
        {loading ? '取得中...' : '今日の日付を取得'}
      </button>
      {date && <p>今日の日付: {date}</p>}
      {error && <p style={{ color: 'red' }}>エラー: {error}</p>}
    </div>
  )
}
```

### 3.2 サーバー側コンポーネント

#### 3.2.1 JWT種類の判別機能

**`server/internal/auth/jwt.go`**に追加:
```go
// JWTType はJWTの種類
type JWTType string

const (
    JWTTypeAuth0        JWTType = "auth0"
    JWTTypePublicAPIKey JWTType = "public_api_key"
    JWTTypeUnknown      JWTType = "unknown"
)

// DetectJWTType はJWTの種類を判別（署名検証前）
func DetectJWTType(tokenString string) (JWTType, error) {
    // 署名検証なしでパース
    parser := jwt.NewParser()
    token, _, err := parser.ParseUnverified(tokenString, &JWTClaims{})
    if err != nil {
        return JWTTypeUnknown, fmt.Errorf("failed to parse JWT: %w", err)
    }

    claims, ok := token.Claims.(*JWTClaims)
    if !ok {
        return JWTTypeUnknown, errors.New("invalid token claims")
    }

    // issによる判別
    if claims.Issuer == "go-webdb-template" {
        return JWTTypePublicAPIKey, nil
    }

    // Auth0のドメインパターンをチェック
    if strings.HasPrefix(claims.Issuer, "https://") && 
       (strings.Contains(claims.Issuer, ".auth0.com") || 
        strings.Contains(claims.Issuer, ".auth0.jp")) {
        return JWTTypeAuth0, nil
    }

    return JWTTypeUnknown, fmt.Errorf("unknown issuer: %s", claims.Issuer)
}
```

#### 3.2.2 Auth0 JWT検証機能

**`server/internal/auth/auth0_validator.go`**（新規）:
```go
package auth

import (
    "context"
    "fmt"
    "time"

    "github.com/MicahParks/keyfunc"
    "github.com/golang-jwt/jwt/v5"
    "github.com/example/go-webdb-template/internal/config"
)

// Auth0Validator はAuth0 JWT検証機能を提供
type Auth0Validator struct {
    jwks *keyfunc.JWKS
}

// NewAuth0Validator は新しいAuth0Validatorを作成
func NewAuth0Validator(issuerBaseURL string) (*Auth0Validator, error) {
    // JWKS URLの構築
    jwksURL := fmt.Sprintf("%s/.well-known/jwks.json", issuerBaseURL)

    // keyfuncのオプション設定
    options := keyfunc.Options{
        RefreshInterval:   time.Hour * 12,   // 12時間ごとに定期更新
        RefreshRateLimit:  time.Minute * 5,  // 再取得は最低5分あける（DoS対策）
        RefreshTimeout:    time.Second * 10, // 取得時のタイムアウト
        RefreshUnknownKID: true,             // 未知のKIDが来たら再取得する（重要！）
    }

    // JWKSの取得とキャッシュ
    jwks, err := keyfunc.Get(jwksURL, options)
    if err != nil {
        return nil, fmt.Errorf("failed to get JWKS: %w", err)
    }

    return &Auth0Validator{
        jwks: jwks,
    }, nil
}

// ValidateAuth0JWT はAuth0 JWTを検証
func (v *Auth0Validator) ValidateAuth0JWT(tokenString string) (*jwt.Token, error) {
    // JWTの検証
    token, err := jwt.Parse(tokenString, v.jwks.Keyfunc)
    if err != nil {
        return nil, fmt.Errorf("failed to validate Auth0 JWT: %w", err)
    }

    if !token.Valid {
        return nil, errors.New("invalid Auth0 JWT")
    }

    return token, nil
}

// Close はリソースを解放
func (v *Auth0Validator) Close() {
    if v.jwks != nil {
        v.jwks.EndBackground()
    }
}
```

#### 3.2.3 認証ミドルウェアの拡張

**`server/internal/auth/middleware.go`**の修正:
```go
// NewHumaAuthMiddleware は新しいHuma形式の認証ミドルウェアを作成
func NewHumaAuthMiddleware(cfg *config.APIConfig, env string, auth0IssuerBaseURL string) func(ctx huma.Context, next func(huma.Context)) {
    validator := NewJWTValidator(cfg, env)
    
    // Auth0Validatorの初期化
    var auth0Validator *Auth0Validator
    if auth0IssuerBaseURL != "" {
        var err error
        auth0Validator, err = NewAuth0Validator(auth0IssuerBaseURL)
        if err != nil {
            // エラーハンドリング（ログ出力など）
            // 起動時エラーとして処理
        }
    }

    return func(ctx huma.Context, next func(huma.Context)) {
        path := ctx.URL().Path

        // OpenAPIドキュメントのパスは認証をスキップ
        if isOpenAPIPath(path) {
            next(ctx)
            return
        }

        // /api/で始まるパスのみ認証を適用
        if !strings.HasPrefix(path, "/api/") {
            next(ctx)
            return
        }

        // AuthorizationヘッダーからJWTトークンを取得
        authHeader := ctx.Header("Authorization")
        if authHeader == "" {
            writeHumaError(ctx, http.StatusUnauthorized, "Authorization header is required")
            return
        }

        // Bearerトークンの抽出
        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            writeHumaError(ctx, http.StatusUnauthorized, "Invalid authorization header format")
            return
        }

        tokenString := parts[1]

        // JWT種類の判別
        jwtType, err := DetectJWTType(tokenString)
        if err != nil {
            writeHumaError(ctx, http.StatusUnauthorized, "Invalid token format")
            return
        }

        var claims *JWTClaims
        var token *jwt.Token

        // JWT種類に応じた検証
        switch jwtType {
        case JWTTypeAuth0:
            if auth0Validator == nil {
                writeHumaError(ctx, http.StatusUnauthorized, "Auth0 JWT validation is not configured")
                return
            }
            token, err = auth0Validator.ValidateAuth0JWT(tokenString)
            if err != nil {
                writeHumaError(ctx, http.StatusUnauthorized, "Invalid Auth0 JWT")
                return
            }
            // Auth0 JWTのクレームを取得（必要に応じて）
            // 注意: Auth0 JWTのクレーム構造は異なる可能性がある
            // ここでは検証のみ行い、クレームの詳細な処理は必要に応じて実装

        case JWTTypePublicAPIKey:
            claims, err = validator.ValidateJWT(tokenString)
            if err != nil {
                writeHumaError(ctx, http.StatusUnauthorized, "Invalid API key")
                return
            }

        default:
            writeHumaError(ctx, http.StatusUnauthorized, "Unknown JWT type")
            return
        }

        // JWTの許容する公開レベルを判定
        var allowedAccessLevel string
        switch jwtType {
        case JWTTypeAuth0:
            // Auth0 JWTはpublicとprivateの両方にアクセス可能
            allowedAccessLevel = "private"
        case JWTTypePublicAPIKey:
            // Public API Key JWTはpublicなAPIのみアクセス可能
            allowedAccessLevel = "public"
        default:
            writeHumaError(ctx, http.StatusUnauthorized, "Unknown JWT type")
            return
        }

        // JWTの許容する公開レベルをコンテキストに設定
        // huma.Contextは標準のcontext.Contextを拡張しているため、context.WithValue()を使用
        ctx = huma.ContextWithValue(ctx, "allowed_access_level", allowedAccessLevel)

        // スコープ検証（Public API Key JWTの場合のみ）
        if jwtType == JWTTypePublicAPIKey {
            if err := validateScope(claims, ctx.Method()); err != nil {
                writeHumaError(ctx, http.StatusForbidden, "Insufficient scope")
                return
            }
        }

        // 次のハンドラーを実行
        next(ctx)
    }
}
```

#### 3.2.4 API公開レベルの定義とチェック

**設計方針**:
- APIの公開レベルは各ハンドラーファイル（`user_handler.go`、`post_handler.go`など）で定義
- エンドポイント登録時に公開レベルを指定し、ミドルウェアでチェック
- JWT検証後、JWTの許容する公開レベルをコンテキストに設定
- ハンドラー実行前に、エンドポイントの公開レベルとJWTの許容する公開レベルを比較

**実装方法**:
1. 各ハンドラーファイルで、エンドポイント登録時に公開レベルを指定（カスタムフィールドまたはメタデータ）
2. ミドルウェアで、JWT検証後に許容する公開レベルをコンテキストに設定
3. ハンドラー実行前に、エンドポイントの公開レベルとJWTの許容する公開レベルを比較
4. Public API Key JWTでprivateなAPIにアクセスした場合、403 Forbiddenを返す

**注意**: Humaの`huma.Operation`構造体に直接公開レベルを追加できない場合は、カスタムメタデータやコンテキストを使用する方法を検討する必要があります。

#### 3.2.5 ハンドラーでの公開レベル定義

**`server/internal/api/handler/user_handler.go`**の修正例:
```go
// RegisterUserEndpoints はHuma APIにユーザーエンドポイントを登録
func RegisterUserEndpoints(api huma.API, h *UserHandler) {
    // POST /api/users - ユーザー作成（public）
    huma.Register(api, huma.Operation{
        OperationID:   "create-user",
        Method:        http.MethodPost,
        Path:          "/api/users",
        Summary:       "ユーザーを作成",
        Tags:          []string{"users"},
        DefaultStatus: http.StatusCreated,
        // 公開レベル: public（カスタムメタデータまたはコンテキストで管理）
    }, func(ctx context.Context, input *humaapi.CreateUserInput) (*humaapi.UserOutput, error) {
        // ミドルウェアで設定されたJWTの許容する公開レベルを取得
        allowedAccessLevel := ctx.Value("allowed_access_level").(string)
        endpointAccessLevel := "public" // このエンドポイントの公開レベル
        
        // 公開レベルのチェック（ミドルウェアで既にチェック済みだが、念のため）
        if endpointAccessLevel == "private" && allowedAccessLevel == "public" {
            return nil, huma.Error403Forbidden("Private API requires Auth0 authentication")
        }
        
        // 既存の処理...
        req := &model.CreateUserRequest{
            Name:  input.Body.Name,
            Email: input.Body.Email,
        }

        user, err := h.userService.CreateUser(ctx, req)
        if err != nil {
            return nil, huma.Error500InternalServerError(err.Error())
        }

        resp := &humaapi.UserOutput{}
        resp.Body = *user
        return resp, nil
    })
    
    // 他のエンドポイントも同様に公開レベルを定義（全てpublic）
}
```

**`server/internal/api/handler/today_handler.go`**（新規）:
```go
package handler

import (
    "context"
    "net/http"
    "time"

    "github.com/danielgtaylor/huma/v2"
    humaapi "github.com/example/go-webdb-template/internal/api/huma"
)

// TodayHandler は今日の日付APIのハンドラー
type TodayHandler struct {
}

// NewTodayHandler は新しいTodayHandlerを作成
func NewTodayHandler() *TodayHandler {
    return &TodayHandler{}
}

// RegisterTodayEndpoints はHuma APIに今日の日付エンドポイントを登録
func RegisterTodayEndpoints(api huma.API, h *TodayHandler) {
    // GET /api/today - 今日の日付を取得（private）
    huma.Register(api, huma.Operation{
        OperationID: "get-today",
        Method:      http.MethodGet,
        Path:        "/api/today",
        Summary:     "今日の日付を取得",
        Tags:        []string{"today"},
        // 公開レベル: private（カスタムメタデータまたはコンテキストで管理）
    }, func(ctx context.Context, input *humaapi.GetTodayInput) (*humaapi.TodayOutput, error) {
        // ミドルウェアで設定されたJWTの許容する公開レベルを取得
        allowedAccessLevel := ctx.Value("allowed_access_level").(string)
        endpointAccessLevel := "private" // このエンドポイントの公開レベル
        
        // 公開レベルのチェック（ミドルウェアで既にチェック済みだが、念のため）
        if endpointAccessLevel == "private" && allowedAccessLevel == "public" {
            return nil, huma.Error403Forbidden("Private API requires Auth0 authentication")
        }
        
        // 今日の日付をYYYY-MM-DD形式で取得
        today := time.Now().Format("2006-01-02")

        resp := &humaapi.TodayOutput{}
        resp.Body.Date = today
        return resp, nil
    })
}
```

**注意**: 上記の実装では、各ハンドラー内で公開レベルのチェックを行います。ミドルウェアではJWTの許容する公開レベルをコンテキストに設定するのみで、実際の公開レベルチェックは各ハンドラーで行います。これにより、APIのパス定義が各ハンドラーファイルに集約され、散在を防ぐことができます。

**ミドルウェアでの処理（637-653行目）**:
- JWTの許容する公開レベルを判定（Auth0 JWT → "private"、Public API Key JWT → "public"）
- 許容する公開レベルをコンテキストに設定（`ctx = huma.ContextWithValue(ctx, "allowed_access_level", allowedAccessLevel)`）

**ハンドラーでの処理（700-710行目、764-774行目）**:
- コンテキストからJWTの許容する公開レベルを取得
- エンドポイントの公開レベルと比較
- Public API Key JWTでprivateなAPIにアクセスした場合、403 Forbiddenを返す

#### 3.2.6 設定構造体の拡張

**`server/internal/config/config.go`**の修正:
```go
// APIConfig はAPIキー設定
type APIConfig struct {
    CurrentVersion     string   `mapstructure:"current_version"`
    PublicKey          string   `mapstructure:"public_key"`
    SecretKey          string   `mapstructure:"secret_key"`
    InvalidVersions    []string `mapstructure:"invalid_versions"`
    Auth0IssuerBaseURL string   `mapstructure:"auth0_issuer_base_url"` // 新規追加
}
```

**`config/{env}/config.yaml`**の修正:
```yaml
api:
  current_version: "v2"
  secret_key: "RNrxs7Rt1ZViughEGb8J08Uc1uQobSOZRRb+BmnGaag="
  invalid_versions:
    - "v1"
  auth0_issuer_base_url: "https://dev-oaa5vtzmld4dsxtd.jp.auth0.com"  # 新規追加
```

## 4. データモデル設計

### 4.1 JWT種類の判別

#### 4.1.1 JWTType列挙型
```go
type JWTType string

const (
    JWTTypeAuth0        JWTType = "auth0"
    JWTTypePublicAPIKey JWTType = "public_api_key"
    JWTTypeUnknown      JWTType = "unknown"
)
```

#### 4.1.2 判別ロジック
- `iss`（Issuer）による判別:
  - `iss == "go-webdb-template"` → `JWTTypePublicAPIKey`
  - `iss`がAuth0ドメイン（`https://*.auth0.com`または`https://*.auth0.jp`） → `JWTTypeAuth0`
  - その他 → `JWTTypeUnknown`

### 4.2 API公開レベルの管理

#### 4.2.1 AccessLevel型
```go
type AccessLevel string

const (
    AccessLevelPublic  AccessLevel = "public"
    AccessLevelPrivate AccessLevel = "private"
)
```

#### 4.2.2 公開レベルの定義方法

**設計方針**:
- APIの公開レベルは各ハンドラーファイル（`user_handler.go`、`post_handler.go`など）で定義
- エンドポイント登録時に公開レベルを指定（カスタムメタデータまたはコンテキストで管理）
- JWT検証後、JWTの許容する公開レベルをコンテキストに設定
- ハンドラー実行時に、エンドポイントの公開レベルとJWTの許容する公開レベルを比較

**JWTの許容する公開レベル**:
- Auth0 JWT: `"private"`（publicとprivateの両方にアクセス可能）
- Public API Key JWT: `"public"`（publicなAPIのみアクセス可能）

### 4.3 Auth0 JWT検証用のデータ構造

#### 4.3.1 Auth0Validator構造体
```go
type Auth0Validator struct {
    jwks *keyfunc.JWKS
}
```

#### 4.3.2 JWKSキャッシュ設定
```go
options := keyfunc.Options{
    RefreshInterval:   time.Hour * 12,   // 12時間ごとに定期更新
    RefreshRateLimit:  time.Minute * 5,  // 再取得は最低5分あける（DoS対策）
    RefreshTimeout:    time.Second * 10, // 取得時のタイムアウト
    RefreshUnknownKID: true,             // 未知のKIDが来たら再取得する（重要！）
}
```

## 5. エラーハンドリング設計

### 5.1 JWT検証エラー

#### 5.1.1 エラーケース
- Authorizationヘッダーが存在しない: 401 Unauthorized
- Bearerトークンの形式が不正: 401 Unauthorized
- JWT種類の判別失敗: 401 Unauthorized
- Auth0 JWT検証失敗: 401 Unauthorized
- Public API Key JWT検証失敗: 401 Unauthorized

#### 5.1.2 エラーレスポンス
```json
{
  "code": 401,
  "message": "Invalid API key"
}
```

### 5.2 アクセス制御エラー

#### 5.2.1 エラーケース
- Public API Key JWTでprivateなAPIにアクセス: 403 Forbidden

#### 5.2.2 エラーレスポンス
```json
{
  "code": 403,
  "message": "Private API requires Auth0 authentication"
}
```

### 5.3 JWKS取得エラー

#### 5.3.1 エラーケース
- JWKS取得のネットワークエラー: ログ出力、リトライ（keyfuncライブラリが自動処理）
- JWKS取得のタイムアウト: ログ出力、リトライ（keyfuncライブラリが自動処理）

#### 5.3.2 エラーハンドリング
- keyfuncライブラリが自動的にリトライを処理
- エラー時はログに記録し、JWT検証を拒否

### 5.4 クライアント側エラーハンドリング

#### 5.4.1 API呼び出しエラー
```typescript
try {
  const response = await fetch(url, options)
  if (!response.ok) {
    if (response.status === 401 || response.status === 403) {
      const errorData = await response.json().catch(() => ({}))
      throw new Error(errorData.message || response.statusText)
    }
    throw new Error(response.statusText)
  }
  return await response.json()
} catch (error) {
  // エラーハンドリング
  console.error('API call failed:', error)
  throw error
}
```

## 6. 設定設計

### 6.1 サーバー側設定

#### 6.1.1 config.yamlへの追加
```yaml
api:
  current_version: "v2"
  secret_key: "RNrxs7Rt1ZViughEGb8J08Uc1uQobSOZRRb+BmnGaag="
  invalid_versions:
    - "v1"
  auth0_issuer_base_url: "https://dev-oaa5vtzmld4dsxtd.jp.auth0.com"
```

#### 6.1.2 環境別設定
- develop: `https://dev-oaa5vtzmld4dsxtd.jp.auth0.com`
- staging: （設定値は環境に応じて決定）
- production: （設定値は環境に応じて決定）

### 6.2 クライアント側設定

#### 6.2.1 環境変数（変更なし）
- Issue #30で設定済みの環境変数をそのまま使用
- `AUTH0_SECRET`, `AUTH0_BASE_URL`, `AUTH0_ISSUER_BASE_URL`, `AUTH0_CLIENT_ID`, `AUTH0_CLIENT_SECRET`

## 7. セキュリティ考慮事項

### 7.1 JWTの安全な検証
- Auth0 JWTはRS256署名を適切に検証
- Public API Key JWTはHS256署名を適切に検証
- JWKSはHTTPSを使用して取得

### 7.2 アクセス制御の厳格な実装
- Public API KeyはprivateなAPIにアクセスできない
- エラーメッセージは適切に管理し、セキュリティ上の情報漏洩を防ぐ

### 7.3 JWKSキャッシュのセキュリティ
- 未知のKIDが来たら再取得（`RefreshUnknownKID: true`）
- DoS対策として再取得間隔を制限（`RefreshRateLimit: time.Minute * 5`）

## 8. パフォーマンス考慮事項

### 8.1 JWKSキャッシュ
- メモリキャッシュを使用してJWKS取得のオーバーヘッドを削減
- 12時間ごとの定期更新により、最新のJWKSを維持

### 8.2 効率的なJWT検証
- JWT種類の判別は署名検証前に実行（`ParseUnverified`を使用）
- 不要な検証処理を回避

## 9. テスト戦略

### 9.1 ユニットテスト

#### 9.1.1 JWT種類の判別テスト
```go
func TestDetectJWTType(t *testing.T) {
    tests := []struct {
        name      string
        token     string
        wantType  JWTType
        wantError bool
    }{
        {
            name:     "Public API Key JWT",
            token:    generatePublicAPIKeyJWT(),
            wantType: JWTTypePublicAPIKey,
        },
        {
            name:     "Auth0 JWT",
            token:    generateAuth0JWT(),
            wantType: JWTTypeAuth0,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            gotType, err := DetectJWTType(tt.token)
            if (err != nil) != tt.wantError {
                t.Errorf("DetectJWTType() error = %v, wantError %v", err, tt.wantError)
                return
            }
            if gotType != tt.wantType {
                t.Errorf("DetectJWTType() = %v, want %v", gotType, tt.wantType)
            }
        })
    }
}
```

#### 9.1.2 API公開レベルの取得テスト
```go
func TestGetAPIAccessLevel(t *testing.T) {
    tests := []struct {
        name     string
        path     string
        method   string
        wantLevel string
    }{
        {
            name:      "Public API",
            path:      "/api/users",
            method:    "GET",
            wantLevel: "public",
        },
        {
            name:      "Private API",
            path:      "/api/today",
            method:    "GET",
            wantLevel: "private",
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            gotLevel := getAPIAccessLevel(tt.path, tt.method)
            if gotLevel != tt.wantLevel {
                t.Errorf("getAPIAccessLevel() = %v, want %v", gotLevel, tt.wantLevel)
            }
        })
    }
}
```

### 9.2 統合テスト

#### 9.2.1 Auth0 JWT検証のテスト
- Auth0 JWTを使用してAPIを呼び出す
- publicなAPIへのアクセスが正常に動作することを確認
- privateなAPIへのアクセスが正常に動作することを確認

#### 9.2.2 Public API Key JWT検証のテスト
- Public API Key JWTを使用してAPIを呼び出す
- publicなAPIへのアクセスが正常に動作することを確認
- privateなAPIへのアクセスが403 Forbiddenを返すことを確認

### 9.3 E2Eテスト

#### 9.3.1 クライアント側のE2Eテスト

**Auth0のJWTテストをPlaywrightで行う場合の選択肢**:

1. **Auth0のテスト用アカウントを使用する（可能）**:
   - テスト用のAuth0アカウントを作成し、実際のAuth0ログインフローをテストする
   - メリット: 実際の動作をテストできる
   - デメリット: 外部サービスに依存するためテストが不安定になる可能性、テストの実行時間が長くなる可能性
   - 実装方法: テスト用の認証情報を環境変数で管理し、Playwrightで実際のログインフローを実行

2. **統合テストとしてサーバー側でテストする（推奨）**:
   - Auth0 JWTの検証はサーバー側の機能のため、サーバー側の統合テストで検証する
   - メリット: テストが高速で安定、外部サービスへの依存が少ない
   - 実装方法: テスト用のAuth0 JWTを生成するか、モックを使用

3. **PlaywrightでUIの動作のみテストする**:
   - Auth0ログイン後のUI動作（ボタンの表示、エラーメッセージの表示など）をテストする
   - メリット: UIの動作を確認できる
   - デメリット: 実際のAuth0ログインフローは含めない

4. **モックを使用する（オプション）**:
   - テスト環境でAuth0のログインフローをモックし、テスト用のJWTを設定する
   - メリット: テストが高速で安定
   - デメリット: 実際のAuth0との統合をテストできない

**Playwrightでの実装例（UI動作のテスト）**:
```typescript
import { test, expect } from '@playwright/test'

test('Today API button shows error when not logged in', async ({ page }) => {
  // 未ログイン状態で開始
  await page.goto('http://localhost:3000')
  
  // Today APIボタンをクリック
  await page.click('text=今日の日付を取得')
  
  // エラーメッセージが表示されることを確認
  await expect(page.locator('text=エラー:')).toBeVisible()
  await expect(page.locator('text=Private API requires Auth0 authentication')).toBeVisible()
})

test('Today API button works after Auth0 login', async ({ page }) => {
  // 前提条件: テスト用のAuth0アカウントを作成し、環境変数に設定
  // TEST_AUTH0_USERNAME: テスト用のAuth0ユーザー名（メールアドレス）
  // TEST_AUTH0_PASSWORD: テスト用のAuth0パスワード
  // 
  // 注意: テスト用認証情報の取り扱い
  // - テスト環境と本番環境が完全に分離されている場合のみ、公開情報として扱うことが可能
  // - テスト用アカウントには最小限の権限のみを付与すること
  // - テスト環境が本番環境と分離されていない場合は、認証情報を機密情報として扱うこと
  // - 一般的には、環境変数で管理し、Gitにコミットしないことを推奨
  
  await page.goto('http://localhost:3000')
  
  // Auth0ログイン（テスト用の認証情報を使用）
  await page.click('text=Login')
  
  // Auth0のログインページにリダイレクトされることを確認
  await page.waitForURL(/auth0\.com/)
  
  // テスト用の認証情報でログイン
  const username = process.env.TEST_AUTH0_USERNAME
  const password = process.env.TEST_AUTH0_PASSWORD
  
  if (!username || !password) {
    test.skip('TEST_AUTH0_USERNAME and TEST_AUTH0_PASSWORD are not set')
  }
  
  await page.fill('input[name="username"], input[type="email"]', username)
  await page.fill('input[name="password"], input[type="password"]', password)
  await page.click('button[type="submit"], button:has-text("Log In")')
  
  // ログイン成功を待つ（コールバックURLにリダイレクトされる）
  await page.waitForURL('http://localhost:3000', { timeout: 10000 })
  
  // ログイン状態を確認（ログアウトボタンが表示されることを確認）
  await expect(page.locator('text=Logout')).toBeVisible()
  
  // Today APIボタンをクリック
  await page.click('text=今日の日付を取得')
  
  // 日付が表示されることを確認
  await expect(page.locator('text=今日の日付:')).toBeVisible()
  
  // エラーメッセージが表示されないことを確認
  await expect(page.locator('text=エラー:')).not.toBeVisible()
})
```

**統合テストとしてサーバー側でテストする（推奨）**:
```go
// server/internal/auth/middleware_test.go
func TestAuth0JWTValidation(t *testing.T) {
    // Auth0 JWTを使用してAPIを呼び出すテスト
    // 実際のAuth0 JWTを生成するか、テスト用のJWTを使用
    // publicなAPIへのアクセスが正常に動作することを確認
    // privateなAPIへのアクセスが正常に動作することを確認
}

func TestPublicAPIKeyJWTValidation(t *testing.T) {
    // Public API Key JWTを使用してAPIを呼び出すテスト
    // publicなAPIへのアクセスが正常に動作することを確認
    // privateなAPIへのアクセスが403 Forbiddenを返すことを確認
}
```

## 10. 参考情報

### 10.1 関連ドキュメント
- 要件定義書: `.kiro/specs/0019-auth0-apicall/requirements.md`
- Auth0ログイン機能設計書: `.kiro/specs/0018-auth0-login/design.md`
- Public API Key機能設計書: `.kiro/specs/0011-public-api-key/design.md`

### 10.2 技術資料
- [github.com/MicahParks/keyfunc](https://github.com/MicahParks/keyfunc)
- [JWKS (JSON Web Key Set) - Auth0](https://auth0.com/docs/secure/tokens/json-web-tokens/json-web-key-sets)
- [Auth0 Next.js SDK - Getting an Access Token](https://auth0.com/docs/quickstart/webapp/nextjs/01-login#get-an-access-token)

### 10.3 既存実装の参考
- 既存の設計書: `.kiro/specs/0018-auth0-login/design.md`（フォーマット参考）
- 既存のJWT検証実装: `server/internal/auth/jwt.go`
- 既存の認証ミドルウェア: `server/internal/auth/middleware.go`
