# APIドキュメント改善設計書

## 1. 概要

### 1.1 設計の目的
要件定義書に基づき、Humaフレームワークで生成されるAPIドキュメントにAuthorization HeaderとAPIアクセスレベルの情報を追加する詳細設計を定義する。SecurityScheme設定と各エンドポイントへのSecurity適用、Tagsによるアクセスレベル表示を実装する。

### 1.2 設計の範囲
- Huma ConfigへのSecurityScheme設定の追加
- すべてのエンドポイントへのSecurityプロパティの適用
- TagsによるAPIアクセスレベルの表示（Public API / Private API）
- 既存の認証ロジックやAPI動作への影響なし（ドキュメント表示のみの改善）

**本設計の範囲外**:
- 認証ロジックの変更（既存の認証ミドルウェアは変更しない）
- APIエンドポイントの追加・削除
- アクセス制御ロジックの変更

### 1.3 設計方針
- **ドキュメント表示のみの改善**: 認証ロジックやAPI動作は一切変更しない
- **一貫性の維持**: すべてのエンドポイントに一貫して適用
- **既存機能の維持**: 既存のTags（"users", "posts", "today"）は維持し、アクセスレベルTagを追加
- **シンプルな実装**: Humaの標準機能を使用し、複雑な実装を避ける

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
```

#### 2.1.2 変更後の構造
```
server/
├── internal/
│   ├── api/
│   │   ├── handler/
│   │   │   ├── user_handler.go          # 修正: SecurityとTags追加
│   │   │   ├── post_handler.go           # 修正: SecurityとTags追加
│   │   │   └── today_handler.go          # 修正: SecurityとTags追加
│   │   └── router/
│   │       └── router.go                 # 修正: SecurityScheme設定の追加
```

### 2.2 ファイル構成

#### 2.2.1 サーバー側（Go）

**`server/internal/api/router/router.go`**: Huma API設定
- SecurityScheme設定の追加（Bearer認証スキーム）
- `huma.Config`の`Components.SecuritySchemes`に設定を追加

**`server/internal/api/handler/user_handler.go`**: ユーザーAPIハンドラー
- すべてのエンドポイント（5つ）に`Security`プロパティを追加
- すべてのエンドポイントに`Tags`に"Public API"を追加

**`server/internal/api/handler/post_handler.go`**: 投稿APIハンドラー
- すべてのエンドポイント（6つ）に`Security`プロパティを追加
- すべてのエンドポイントに`Tags`に"Public API"を追加

**`server/internal/api/handler/today_handler.go`**: 今日の日付APIハンドラー
- エンドポイント（1つ）に`Security`プロパティを追加
- エンドポイントに`Tags`に"Private API"を追加

### 2.3 システム構成図

```
┌─────────────────────────────────────────────────────────┐
│              サーバー（Go）                                │
│                                                           │
│  ┌──────────────────────────────────────────────────┐   │
│  │  server/internal/api/router/router.go            │   │
│  │  - Huma Config設定                                │   │
│  │  - SecurityScheme設定（Bearer認証）               │   │
│  └──────────────────┬───────────────────────────────┘   │
│                     │                                     │
│                     ▼                                     │
│  ┌──────────────────────────────────────────────────┐   │
│  │  Huma APIインスタンス                             │   │
│  │  - SecurityScheme: bearerAuth                     │   │
│  └──────────────────┬───────────────────────────────┘   │
│                     │                                     │
│         ┌───────────┴───────────┐                      │
│         │                         │                      │
│         ▼                         ▼                      │
│  ┌──────────────┐        ┌──────────────┐              │
│  │ User         │        │ Post          │              │
│  │ Endpoints    │        │ Endpoints     │              │
│  │ (5個)        │        │ (6個)        │              │
│  │              │        │              │              │
│  │ Security:    │        │ Security:    │              │
│  │ bearerAuth   │        │ bearerAuth   │              │
│  │              │        │              │              │
│  │ Tags:        │        │ Tags:        │              │
│  │ ["users",    │        │ ["posts",    │              │
│  │  "Public     │        │  "Public     │              │
│  │   API"]      │        │   API"]      │              │
│  └──────────────┘        └──────────────┘              │
│         │                         │                      │
│         └───────────┬─────────────┘                      │
│                     │                                     │
│                     ▼                                     │
│  ┌──────────────┐                                        │
│  │ Today        │                                        │
│  │ Endpoint     │                                        │
│  │ (1個)        │                                        │
│  │              │                                        │
│  │ Security:    │                                        │
│  │ bearerAuth   │                                        │
│  │              │                                        │
│  │ Tags:        │                                        │
│  │ ["today",    │                                        │
│  │  "Private    │                                        │
│  │   API"]      │                                        │
│  └──────────────┘                                        │
│                     │                                     │
│                     ▼                                     │
│  ┌──────────────────────────────────────────────────┐   │
│  │  OpenAPI仕様生成                                 │   │
│  │  - SecurityScheme定義                            │   │
│  │  - 各エンドポイントのSecurity定義                │   │
│  │  - Tags定義                                      │   │
│  └──────────────────┬───────────────────────────────┘   │
│                     │                                     │
│                     ▼                                     │
│  ┌──────────────────────────────────────────────────┐   │
│  │  Swagger UI (/docs)                               │   │
│  │  - "Authorize"ボタン表示                          │   │
│  │  - Request SampleにAuthorization Header表示      │   │
│  │  - Tagsでアクセスレベル表示                       │   │
│  └──────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────┘
```

### 2.4 データフロー

#### 2.4.1 APIドキュメント生成フロー
```
サーバー起動時
    ↓
router.go: Huma Config設定
    ↓
router.go: SecurityScheme設定（bearerAuth）
    ↓
各ハンドラー: エンドポイント登録
    ↓
各ハンドラー: Securityプロパティ設定（bearerAuth）
    ↓
各ハンドラー: Tags設定（既存Tag + アクセスレベルTag）
    ↓
Huma: OpenAPI仕様生成
    ↓
Swagger UI: ドキュメント表示
    ↓
- "Authorize"ボタン表示
- Request SampleにAuthorization Header表示
- Tagsでアクセスレベル表示
```

## 3. コンポーネント設計

### 3.1 SecurityScheme設定

#### 3.1.1 router.goの修正

**`server/internal/api/router/router.go`**:
```go
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

// Huma APIインスタンスの作成
humaAPI := humaecho.New(e, humaConfig)
```

**実装上の注意**:
- `huma.DefaultConfig()`で作成したConfigに`Components`を追加する際は、既存の設定を上書きしないよう注意
- SecurityScheme名は"bearerAuth"を使用（Issue #37の記載に従う）

### 3.2 各エンドポイントへのSecurity適用

#### 3.2.1 user_handler.goの修正

**`server/internal/api/handler/user_handler.go`**の修正例:
```go
// RegisterUserEndpoints はHuma APIにユーザーエンドポイントを登録
func RegisterUserEndpoints(api huma.API, h *UserHandler) {
    // POST /api/users - ユーザー作成
    huma.Register(api, huma.Operation{
        OperationID:   "create-user",
        Method:        http.MethodPost,
        Path:          "/api/users",
        Summary:       "ユーザーを作成",
        Tags:          []string{"users", "Public API"},
        DefaultStatus: http.StatusCreated,
        Security: []map[string][]string{
            {"bearerAuth": []},
        },
    }, func(ctx context.Context, input *humaapi.CreateUserInput) (*humaapi.UserOutput, error) {
        // 既存の処理（変更なし）
        // ...
    })

    // GET /api/users/{id} - ユーザー取得
    huma.Register(api, huma.Operation{
        OperationID: "get-user",
        Method:      http.MethodGet,
        Path:        "/api/users/{id}",
        Summary:     "ユーザーを取得",
        Tags:        []string{"users", "Public API"},
        Security: []map[string][]string{
            {"bearerAuth": []},
        },
    }, func(ctx context.Context, input *humaapi.GetUserInput) (*humaapi.UserOutput, error) {
        // 既存の処理（変更なし）
        // ...
    })

    // GET /api/users - ユーザー一覧取得
    huma.Register(api, huma.Operation{
        OperationID: "list-users",
        Method:      http.MethodGet,
        Path:        "/api/users",
        Summary:     "ユーザー一覧を取得",
        Tags:        []string{"users", "Public API"},
        Security: []map[string][]string{
            {"bearerAuth": []},
        },
    }, func(ctx context.Context, input *humaapi.ListUsersInput) (*humaapi.UsersOutput, error) {
        // 既存の処理（変更なし）
        // ...
    })

    // PUT /api/users/{id} - ユーザー更新
    huma.Register(api, huma.Operation{
        OperationID: "update-user",
        Method:      http.MethodPut,
        Path:        "/api/users/{id}",
        Summary:     "ユーザーを更新",
        Tags:        []string{"users", "Public API"},
        Security: []map[string][]string{
            {"bearerAuth": []},
        },
    }, func(ctx context.Context, input *humaapi.UpdateUserInput) (*humaapi.UserOutput, error) {
        // 既存の処理（変更なし）
        // ...
    })

    // DELETE /api/users/{id} - ユーザー削除
    huma.Register(api, huma.Operation{
        OperationID:   "delete-user",
        Method:        http.MethodDelete,
        Path:          "/api/users/{id}",
        Summary:       "ユーザーを削除",
        Tags:          []string{"users", "Public API"},
        DefaultStatus: http.StatusNoContent,
        Security: []map[string][]string{
            {"bearerAuth": []},
        },
    }, func(ctx context.Context, input *humaapi.DeleteUserInput) (*struct{}, error) {
        // 既存の処理（変更なし）
        // ...
    })
}
```

#### 3.2.2 post_handler.goの修正

**`server/internal/api/handler/post_handler.go`**の修正例:
```go
// RegisterPostEndpoints はHuma APIに投稿エンドポイントを登録
func RegisterPostEndpoints(api huma.API, h *PostHandler) {
    // POST /api/posts - 投稿作成
    huma.Register(api, huma.Operation{
        OperationID:   "create-post",
        Method:        http.MethodPost,
        Path:          "/api/posts",
        Summary:       "投稿を作成",
        Tags:          []string{"posts", "Public API"},
        DefaultStatus: http.StatusCreated,
        Security: []map[string][]string{
            {"bearerAuth": []},
        },
    }, func(ctx context.Context, input *humaapi.CreatePostInput) (*humaapi.PostOutput, error) {
        // 既存の処理（変更なし）
        // ...
    })

    // GET /api/posts/{id} - 投稿取得
    huma.Register(api, huma.Operation{
        OperationID: "get-post",
        Method:      http.MethodGet,
        Path:        "/api/posts/{id}",
        Summary:     "投稿を取得",
        Tags:        []string{"posts", "Public API"},
        Security: []map[string][]string{
            {"bearerAuth": []},
        },
    }, func(ctx context.Context, input *humaapi.GetPostInput) (*humaapi.PostOutput, error) {
        // 既存の処理（変更なし）
        // ...
    })

    // GET /api/posts - 投稿一覧取得
    huma.Register(api, huma.Operation{
        OperationID: "list-posts",
        Method:      http.MethodGet,
        Path:        "/api/posts",
        Summary:     "投稿一覧を取得",
        Tags:        []string{"posts", "Public API"},
        Security: []map[string][]string{
            {"bearerAuth": []},
        },
    }, func(ctx context.Context, input *humaapi.ListPostsInput) (*humaapi.PostsOutput, error) {
        // 既存の処理（変更なし）
        // ...
    })

    // PUT /api/posts/{id} - 投稿更新
    huma.Register(api, huma.Operation{
        OperationID: "update-post",
        Method:      http.MethodPut,
        Path:        "/api/posts/{id}",
        Summary:     "投稿を更新",
        Tags:        []string{"posts", "Public API"},
        Security: []map[string][]string{
            {"bearerAuth": []},
        },
    }, func(ctx context.Context, input *humaapi.UpdatePostInput) (*humaapi.PostOutput, error) {
        // 既存の処理（変更なし）
        // ...
    })

    // DELETE /api/posts/{id} - 投稿削除
    huma.Register(api, huma.Operation{
        OperationID:   "delete-post",
        Method:        http.MethodDelete,
        Path:          "/api/posts/{id}",
        Summary:       "投稿を削除",
        Tags:          []string{"posts", "Public API"},
        DefaultStatus: http.StatusNoContent,
        Security: []map[string][]string{
            {"bearerAuth": []},
        },
    }, func(ctx context.Context, input *humaapi.DeletePostInput) (*struct{}, error) {
        // 既存の処理（変更なし）
        // ...
    })

    // GET /api/user-posts - ユーザーと投稿のJOIN結果取得
    huma.Register(api, huma.Operation{
        OperationID: "get-user-posts",
        Method:      http.MethodGet,
        Path:        "/api/user-posts",
        Summary:     "ユーザーと投稿のJOIN結果を取得",
        Tags:        []string{"posts", "Public API"},
        Security: []map[string][]string{
            {"bearerAuth": []},
        },
    }, func(ctx context.Context, input *humaapi.GetUserPostsInput) (*humaapi.UserPostsOutput, error) {
        // 既存の処理（変更なし）
        // ...
    })
}
```

#### 3.2.3 today_handler.goの修正

**`server/internal/api/handler/today_handler.go`**の修正例:
```go
// RegisterTodayEndpoints はHuma APIにToday APIエンドポイントを登録
func RegisterTodayEndpoints(api huma.API, h *TodayHandler) {
    // GET /api/today - 今日の日付取得（private API）
    huma.Register(api, huma.Operation{
        OperationID: "get-today",
        Method:      http.MethodGet,
        Path:        "/api/today",
        Summary:     "今日の日付を取得（Auth0認証必須）",
        Tags:        []string{"today", "Private API"},
        Security: []map[string][]string{
            {"bearerAuth": []},
        },
    }, func(ctx context.Context, input *humaapi.GetTodayInput) (*humaapi.TodayOutput, error) {
        // 既存の処理（変更なし）
        // ...
    })
}
```

### 3.3 Tagsによるアクセスレベル表示

#### 3.3.1 Tag設定の実装

**実装方針**:
- 既存の機能タグ（"users", "posts", "today"）は維持
- アクセスレベルのTagを追加（"Public API" / "Private API"）
- すべてのエンドポイントに一貫して適用

**Tag名**:
- "Public API" - publicレベルのAPI（Public API Key JWTとAuth0 JWTの両方でアクセス可能）
- "Private API" - privateレベルのAPI（Auth0 JWTのみでアクセス可能）

**各エンドポイントのTag設定**:
- Userエンドポイント（5つ）：`[]string{"users", "Public API"}`
- Postエンドポイント（6つ）：`[]string{"posts", "Public API"}`
- Todayエンドポイント（1つ）：`[]string{"today", "Private API"}`

## 4. データモデル設計

### 4.1 SecurityScheme定義

#### 4.1.1 SecurityScheme構造
```go
huma.SecurityScheme{
    Type:         "http",      // HTTP認証スキーム
    Scheme:       "bearer",    // Bearer認証
    BearerFormat: "JWT",       // JWT形式
}
```

#### 4.1.2 SecurityScheme名
- "bearerAuth" - Issue #37の記載に従う

### 4.2 Securityプロパティ

#### 4.2.1 Securityプロパティの形式
```go
Security: []map[string][]string{
    {"bearerAuth": []},
}
```

**説明**:
- `[]map[string][]string` - Securityプロパティの形式
- `{"bearerAuth": []}` - bearerAuthスキームを使用、スコープは空（すべてのスコープを許可）

### 4.3 Tags定義

#### 4.3.1 Tag構造
```go
Tags: []string{"users", "Public API"}
```

**説明**:
- 最初の要素: 機能的なグループ化タグ（既存のTag）
- 2番目の要素: アクセスレベルのTag（新規追加）

#### 4.3.2 アクセスレベルTagの定義
- "Public API" - publicレベルのAPI
- "Private API" - privateレベルのAPI

## 5. エラーハンドリング設計

### 5.1 実装上の注意事項

#### 5.1.1 SecurityScheme設定
- `huma.DefaultConfig()`で作成したConfigに`Components`を追加する際は、既存の設定を上書きしないよう注意
- SecurityScheme名は"bearerAuth"を使用（Issue #37の記載に従う）

#### 5.1.2 Securityプロパティの適用
- すべてのエンドポイントに一貫して適用する
- Securityプロパティの形式は`[]map[string][]string{{"bearerAuth": []}}`とする
- 空の配列`[]`を指定することで、すべてのスコープを許可する

#### 5.1.3 Tagsによるアクセスレベル表示
- 既存の機能タグ（"users", "posts", "today"）は維持し、アクセスレベルTagを追加する
- Tag名は大文字小文字を区別するため、"Public API"と"Private API"を正確に記述する
- 各エンドポイントのアクセスレベルは、既存の`auth.CheckAccessLevel()`の呼び出しから判断する

### 5.2 既存機能との互換性

#### 5.2.1 認証ロジック
- 認証ロジックやAPI動作は一切変更しない
- 既存の`auth.CheckAccessLevel()`によるアクセス制御は維持
- ドキュメント表示のみを改善する

#### 5.2.2 API動作
- 既存のAPIエンドポイントの動作は変更しない
- 既存のAPIクライアントへの影響なし

## 6. 設定設計

### 6.1 Huma Config設定

#### 6.1.1 SecurityScheme設定
```go
humaConfig.Components = &huma.Components{
    SecuritySchemes: map[string]*huma.SecurityScheme{
        "bearerAuth": {
            Type:         "http",
            Scheme:       "bearer",
            BearerFormat: "JWT",
        },
    },
}
```

#### 6.1.2 設定場所
- `server/internal/api/router/router.go`の`humaConfig`設定部分

### 6.2 エンドポイント設定

#### 6.2.1 Securityプロパティ
- すべてのエンドポイントに`Security`プロパティを追加
- 形式: `[]map[string][]string{{"bearerAuth": []}}`

#### 6.2.2 Tags設定
- 既存の機能タグに加えて、アクセスレベルTagを追加
- Userエンドポイント: `[]string{"users", "Public API"}`
- Postエンドポイント: `[]string{"posts", "Public API"}`
- Todayエンドポイント: `[]string{"today", "Private API"}`

## 7. セキュリティ考慮事項

### 7.1 ドキュメント表示のみの改善
- 認証ロジックやAPI動作は一切変更しない
- 既存の認証ミドルウェアは変更しない
- ドキュメント表示のみを改善する

### 7.2 情報の明確化
- APIドキュメント上で認証が必要であることが明確になる
- 各APIのアクセスレベルが明確になる
- Swagger UIから直接APIをテストできるようになる

## 8. パフォーマンス考慮事項

### 8.1 ドキュメント生成への影響
- SecurityScheme設定とSecurityプロパティの追加は、ドキュメント生成時の処理のみに影響
- 実行時のパフォーマンスへの影響はない

### 8.2 メモリ使用量
- SecurityScheme設定とSecurityプロパティの追加によるメモリ使用量の増加は無視できる程度

## 9. テスト戦略

### 9.1 ユニットテスト

#### 9.1.1 SecurityScheme設定のテスト
- `huma.Config`の`Components.SecuritySchemes`にBearer認証スキームが定義されていることを確認
- SecuritySchemeの設定内容が正しいことを確認（Type: "http", Scheme: "bearer", BearerFormat: "JWT"）

#### 9.1.2 Securityプロパティのテスト
- すべてのエンドポイント（12個）に`Security`プロパティが追加されていることを確認
- Securityプロパティの形式が正しいことを確認（`[]map[string][]string{{"bearerAuth": []}}`）

#### 9.1.3 Tags設定のテスト
- すべてのエンドポイントに適切なアクセスレベルTagが追加されていることを確認
- 既存の機能タグ（"users", "posts", "today"）が維持されていることを確認

### 9.2 統合テスト

#### 9.2.1 APIドキュメントの表示確認
- `http://localhost:8080/docs`にアクセスして、Request Sample（curl）に`-H 'Authorization: Bearer <token>'`が表示されることを確認
- Swagger UIに「Authorize」ボタンが表示されることを確認
- 「Authorize」ボタンをクリックしてトークンを入力できることを確認
- 各APIにアクセスレベルのTag（"Public API" / "Private API"）が表示されることを確認
- 既存の機能タグ（"users", "posts", "today"）も表示されることを確認

#### 9.2.2 既存機能の動作確認
- 既存のAPI動作に影響がないことを確認（認証、アクセス制御が正常に動作）
- 既存のAPIクライアントへの影響がないことを確認

### 9.3 E2Eテスト

#### 9.3.1 Swagger UIでのAPIテスト
- Swagger UIの「Authorize」ボタンからトークンを入力
- 各APIをSwagger UIから直接テスト
- Request Sampleのcurlコマンドが正しく生成されることを確認

## 10. 参考情報

### 10.1 関連ドキュメント
- 要件定義書: `.kiro/specs/0020-apidoc/requirements.md`
- Auth0 API呼び出し機能設計書: `.kiro/specs/0019-auth0-apicall/design.md`
- Echo + Huma導入設計書: `.kiro/specs/0013-echo-huma/design.md`

### 10.2 技術資料
- [Huma Documentation](https://huma.rocks/)
- [OpenAPI Security Schemes](https://swagger.io/specification/#security-scheme-object)
- Issue #37のコメントに記載されている実装方法

### 10.3 既存実装の参考
- 既存のHuma設定: `server/internal/api/router/router.go`
- 既存のエンドポイント登録: `server/internal/api/handler/user_handler.go`, `post_handler.go`, `today_handler.go`
