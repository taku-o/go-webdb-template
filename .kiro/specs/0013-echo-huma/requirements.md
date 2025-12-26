# Echo・Huma導入要件定義書

## 1. 概要

### 1.1 プロジェクト情報
- **プロジェクト名**: go-webdb-template
- **Issue番号**: #24
- **Issueタイトル**: APIサーバーのリクエストの処理にEcho、および、Humaを導入
- **Feature名**: 0013-echo-huma
- **作成日**: 2025-01-27

### 1.2 目的
APIサーバーのリクエスト処理にEchoとHumaを導入し、既存のGorilla Muxベースの実装をHumaベースに置き換える。
これにより、ライブラリの知見を取り込み、自動バリデーション、OpenAPI仕様の自動生成などの機能を活用し、コードの保守性と開発効率を向上させる。

### 1.3 スコープ
- Echoフレームワークの導入と設定
- Humaフレームワークの導入と設定（humaechoアダプター使用）
- 既存のAPIエンドポイントをHuma (huma.Register) で置き換え
- リクエスト/レスポンス構造体のHumaタグ対応
- 認証ミドルウェアのEcho/Huma統合
- CORS設定の維持
- アクセスログの統合
- OpenAPI仕様の自動生成

## 2. 背景・現状分析

### 2.1 現在の実装
- **ルーティング**: Gorilla Mux (`github.com/gorilla/mux v1.8.1`)
- **ハンドラー**: 標準の`http.Handler`インターフェースを使用
- **リクエスト処理**: 手動でJSONデコード、バリデーション、エラーハンドリング
- **認証**: JWT認証ミドルウェア（`server/internal/auth/middleware.go`）
- **CORS**: `github.com/rs/cors`を使用
- **アクセスログ**: カスタムミドルウェア（`server/internal/logging`）
- **エンドポイント**:
  - `/api/users` (POST, GET)
  - `/api/users/{id}` (GET, PUT, DELETE)
  - `/api/posts` (POST, GET)
  - `/api/posts/{id}` (GET, PUT, DELETE)
  - `/api/user-posts` (GET)
  - `/health` (GET)

### 2.2 課題点
1. **手動実装の多さ**: リクエストのデコード、バリデーション、エラーハンドリングを手動で実装している
2. **バリデーションの不統一**: 各ハンドラーで個別にバリデーション処理を実装している
3. **OpenAPI仕様の欠如**: API仕様が手動で管理されており、コードと仕様の不整合が発生しやすい
4. **エラーハンドリングの重複**: 各ハンドラーで同様のエラーハンドリングコードが重複している
5. **ライブラリの知見活用不足**: 標準ライブラリのみを使用しており、成熟したフレームワークの知見を活用できていない

### 2.3 本実装による改善点
1. **自動バリデーション**: Humaのバリデーションタグにより、リクエストの自動バリデーションが可能
2. **OpenAPI仕様の自動生成**: HumaによりOpenAPI 3.0仕様が自動生成され、コードと仕様の整合性が保証される
3. **コードの簡潔化**: ハンドラー関数が簡潔になり、ビジネスロジックに集中できる
4. **ライブラリの知見活用**: EchoとHumaのベストプラクティスを活用できる
5. **型安全性の向上**: リクエスト/レスポンスの型定義が明確になり、コンパイル時にエラーを検出できる

## 3. 機能要件

### 3.1 Echoフレームワークの導入

#### 3.1.1 Echoの設定
- **パッケージ**: `github.com/labstack/echo/v4`
- **初期化**: `echo.New()`でEchoインスタンスを作成
- **ミドルウェア**: 既存の認証、CORS、アクセスログミドルウェアをEcho形式に変換
- **サーバー起動**: `e.Start(":8080")`でサーバーを起動

#### 3.1.2 Echoへの移行
- 既存のGorilla MuxルーターをEchoルーターに置き換え
- 既存の`http.Handler`インターフェースをEchoのハンドラー形式に変換
- パスパラメータの取得方法を変更（`mux.Vars(r)` → `c.Param()`）
- クエリパラメータの取得方法を変更（`r.URL.Query().Get()` → `c.QueryParam()`）

### 3.2 Humaフレームワークの導入

#### 3.2.1 Humaの設定
- **パッケージ**: `github.com/danielgtaylor/huma/v2`
- **アダプター**: `github.com/danielgtaylor/huma/v2/adapters/humaecho`
- **初期化**: `humaecho.New(e, config)`でHuma APIインスタンスを作成
- **設定**: `huma.DefaultConfig("API Name", "1.0.0")`でデフォルト設定を使用

#### 3.2.2 Humaによるエンドポイント登録
既存のエンドポイントを`huma.Register()`で登録：

```go
huma.Register(api, huma.Operation{
    Method: http.MethodGet,
    Path:   "/users/{id}",
    Summary: "ユーザーを取得",
}, func(ctx context.Context, input *GetUserInput) (*GetUserOutput, error) {
    // ハンドラー実装
})
```

### 3.3 リクエスト/レスポンス構造体の定義

#### 3.3.1 Humaタグの追加
既存のリクエスト/レスポンス構造体にHumaタグを追加：

**ユーザーエンドポイント**:
- `CreateUserInput`: `CreateUserRequest`をベースに、Humaタグを追加
- `GetUserInput`: パスパラメータ`id`を定義
- `UpdateUserInput`: パスパラメータ`id`とリクエストボディを定義
- `DeleteUserInput`: パスパラメータ`id`を定義
- `ListUsersInput`: クエリパラメータ`limit`, `offset`を定義
- `UserOutput`: `User`モデルをベースにレスポンス構造体を定義
- `UsersOutput`: `[]*User`をベースにレスポンス構造体を定義

**投稿エンドポイント**:
- `CreatePostInput`: `CreatePostRequest`をベースに、Humaタグを追加
- `GetPostInput`: パスパラメータ`id`とクエリパラメータ`user_id`を定義
- `UpdatePostInput`: パスパラメータ`id`、クエリパラメータ`user_id`、リクエストボディを定義
- `DeletePostInput`: パスパラメータ`id`とクエリパラメータ`user_id`を定義
- `ListPostsInput`: クエリパラメータ`limit`, `offset`, `user_id`を定義
- `GetUserPostsInput`: クエリパラメータ`limit`, `offset`を定義
- `PostOutput`: `Post`モデルをベースにレスポンス構造体を定義
- `PostsOutput`: `[]*Post`をベースにレスポンス構造体を定義
- `UserPostsOutput`: `[]*UserPost`をベースにレスポンス構造体を定義

#### 3.3.2 バリデーションタグ
Humaのバリデーションタグを使用：

```go
type CreateUserInput struct {
    Body struct {
        Name  string `json:"name" maxLength:"100" doc:"ユーザー名"`
        Email string `json:"email" format:"email" maxLength:"255" doc:"メールアドレス"`
    } `json:"body"`
}
```

### 3.4 既存エンドポイントのHuma化

#### 3.4.1 ユーザーエンドポイント
- **POST /api/users**: ユーザー作成
- **GET /api/users**: ユーザー一覧取得
- **GET /api/users/{id}**: ユーザー取得
- **PUT /api/users/{id}**: ユーザー更新
- **DELETE /api/users/{id}**: ユーザー削除

#### 3.4.2 投稿エンドポイント
- **POST /api/posts**: 投稿作成
- **GET /api/posts**: 投稿一覧取得
- **GET /api/posts/{id}**: 投稿取得（`user_id`クエリパラメータ必須）
- **PUT /api/posts/{id}**: 投稿更新（`user_id`クエリパラメータ必須）
- **DELETE /api/posts/{id}**: 投稿削除（`user_id`クエリパラメータ必須）
- **GET /api/user-posts**: ユーザーと投稿のJOIN取得

#### 3.4.3 ヘルスチェックエンドポイント
- **GET /health**: ヘルスチェック（Huma化は不要、Echoの通常ハンドラーで実装）

### 3.5 認証ミドルウェアの統合

#### 3.5.1 Echoミドルウェア形式への変換
既存のJWT認証ミドルウェアをEchoのミドルウェア形式に変換：

```go
func AuthMiddleware() echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            // 認証処理
            return next(c)
        }
    }
}
```

#### 3.5.2 Humaとの統合
- Humaのエンドポイントに認証ミドルウェアを適用
- 認証エラーはHumaのエラーレスポンス形式で返す

### 3.6 CORS設定の維持

#### 3.6.1 Echo CORSミドルウェア
既存のCORS設定をEchoのCORSミドルウェアに移行：

```go
e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
    AllowOrigins: cfg.CORS.AllowedOrigins,
    AllowMethods: cfg.CORS.AllowedMethods,
    AllowHeaders: cfg.CORS.AllowedHeaders,
    AllowCredentials: true,
}))
```

### 3.7 アクセスログの統合

#### 3.7.1 Echoログミドルウェア
既存のアクセスログ機能をEchoのログミドルウェアに統合：

```go
e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
    // 既存のアクセスログ設定を反映
}))
```

### 3.8 OpenAPI仕様の自動生成

#### 3.8.1 OpenAPIエンドポイント
Humaにより自動生成されるOpenAPI仕様にアクセスできるエンドポイントを追加：

- **GET /openapi.json**: OpenAPI 3.0仕様（JSON形式）
- **GET /openapi.yaml**: OpenAPI 3.0仕様（YAML形式）

#### 3.8.2 仕様の内容
- 全エンドポイントの定義
- リクエスト/レスポンススキーマ
- バリデーションルール
- エラーレスポンス定義

## 4. 非機能要件

### 4.1 既存機能の互換性
- 既存のAPIエンドポイントの動作は維持されること
- 既存のリクエスト/レスポンス形式は維持されること
- 既存の認証方式（JWT）は維持されること
- 既存のCORS設定は維持されること

### 4.2 パフォーマンス
- 既存のパフォーマンスを維持すること
- レスポンスタイムの劣化は許容範囲内であること

### 4.3 エラーハンドリング
- Humaの標準エラーレスポンス形式を使用すること
- 既存のエラーメッセージと互換性を保つこと

### 4.4 テスト
- 既存のテストコードが動作すること
- 新しいテストコードでHumaの動作を検証すること

## 5. 制約事項

### 5.1 既存コードの変更範囲
- **Service層**: 変更しない（既存のService層をそのまま使用）
- **Repository層**: 変更しない（既存のRepository層をそのまま使用）
- **Model層**: Humaタグを追加するが、既存の構造体定義は維持
- **Handler層**: Huma形式に置き換えるが、既存のService層呼び出しは維持

### 5.2 既存API仕様の維持
- 既存のリクエスト/レスポンス形式は維持すること
- 既存のエンドポイントパスは維持すること
- 既存のHTTPメソッドは維持すること

### 5.3 認証方式の維持
- 既存のJWT認証方式は維持すること
- 既存の認証ミドルウェアのロジックは維持すること

### 5.4 技術スタック
- **Go**: 1.23.4（既存バージョンを維持）
- **Echo**: v4（最新の安定版）
- **Huma**: v2（最新の安定版）
- 既存の依存関係は可能な限り維持すること

## 6. 受け入れ基準

### 6.1 Echoの導入
- [ ] Echoフレームワークが導入されている
- [ ] 既存のルーティングがEchoで実装されている
- [ ] 既存のミドルウェアがEcho形式に変換されている

### 6.2 Humaの導入
- [ ] Humaフレームワークが導入されている
- [ ] humaechoアダプターが設定されている
- [ ] Huma APIインスタンスが作成されている

### 6.3 エンドポイントのHuma化
- [ ] 全ユーザーエンドポイントがHumaで実装されている
- [ ] 全投稿エンドポイントがHumaで実装されている
- [ ] ヘルスチェックエンドポイントがEchoで実装されている

### 6.4 リクエスト/レスポンス構造体
- [ ] 全エンドポイントのリクエスト構造体にHumaタグが追加されている
- [ ] 全エンドポイントのレスポンス構造体が定義されている
- [ ] バリデーションタグが適切に設定されている

### 6.5 認証・CORS・アクセスログ
- [ ] 認証ミドルウェアがEcho形式に変換されている
- [ ] CORS設定がEcho形式に移行されている
- [ ] アクセスログがEcho形式に統合されている

### 6.6 OpenAPI仕様
- [ ] OpenAPI仕様が自動生成されている
- [ ] `/openapi.json`エンドポイントが動作している
- [ ] `/openapi.yaml`エンドポイントが動作している
- [ ] 全エンドポイントがOpenAPI仕様に含まれている

### 6.7 テスト
- [ ] 既存のテストが全て通過している
- [ ] 新しいテストでHumaの動作を検証している
- [ ] E2Eテストが正常に動作している

### 6.8 ドキュメント
- [ ] `docs/API.md`が更新されている（OpenAPI仕様への参照を追加）
- [ ] 変更内容が文書化されている

## 7. 影響範囲

### 7.1 新規追加が必要なディレクトリ・ファイル

#### ディレクトリ
- なし（既存のディレクトリ構造を維持）

#### ファイル
- `server/internal/api/huma/`: Huma関連の構造体定義（新規、オプション）
  - `inputs.go`: リクエスト構造体定義
  - `outputs.go`: レスポンス構造体定義

### 7.2 変更が必要なファイル

#### ルーティング
- `server/internal/api/router/router.go`: Gorilla MuxからEcho/Humaに置き換え

#### ハンドラー
- `server/internal/api/handler/user_handler.go`: Huma形式に置き換え
- `server/internal/api/handler/post_handler.go`: Huma形式に置き換え

#### 認証
- `server/internal/auth/middleware.go`: Echoミドルウェア形式に変換

#### ログ
- `server/internal/logging/`: Echoログミドルウェア形式に変換（該当ファイル）

#### メイン
- `server/cmd/server/main.go`: Echoサーバー起動に変更

#### 依存関係
- `server/go.mod`: EchoとHumaの依存関係を追加

#### テストファイル
- `server/test/e2e/api_test.go`: テストコードの更新（必要に応じて）
- `server/internal/api/handler/*_test.go`: テストコードの更新（必要に応じて）

### 7.3 削除されるファイル
- なし（既存ファイルは変更のみ）

### 7.4 ドキュメント
- `docs/API.md`: OpenAPI仕様への参照を追加
- `README.md`: 依存関係の更新（必要に応じて）

## 8. 実装上の注意事項

### 8.1 Echoへの移行
- パスパラメータの取得方法を変更（`mux.Vars(r)` → `c.Param()`）
- クエリパラメータの取得方法を変更（`r.URL.Query().Get()` → `c.QueryParam()`）
- リクエストボディの取得方法を変更（`json.NewDecoder(r.Body)` → `c.Bind()`）
- レスポンスの書き込み方法を変更（`json.NewEncoder(w)` → `c.JSON()`）

### 8.2 Humaの使用
- `huma.Register()`でエンドポイントを登録
- リクエスト構造体にはHumaタグを適切に設定
- レスポンス構造体は`Body`フィールドを持つ構造体として定義
- エラーハンドリングはHumaの標準形式を使用

### 8.3 バリデーション
- Humaのバリデーションタグを使用（`maxLength`, `minLength`, `format`, `required`など）
- 既存の`validate`タグはHumaタグに変換（必要に応じて）

### 8.4 認証ミドルウェア
- Echoのミドルウェア形式に変換
- Humaのエンドポイントに適用する場合は、Echoのミドルウェアとして登録
- 認証エラーはHumaのエラーレスポンス形式で返す

### 8.5 テスト
- 既存のテストコードを可能な限り維持
- Echo/Humaのテストユーティリティを使用（必要に応じて）
- OpenAPI仕様の検証テストを追加（オプション）

## 9. 参考情報

### 9.1 関連Issue
- GitHub Issue #24: APIサーバーのリクエストの処理にEcho、および、Humaを導入

### 9.2 既存ドキュメント
- `docs/API.md`: 既存のAPI仕様
- `server/internal/api/router/router.go`: 既存のルーター実装
- `server/internal/api/handler/`: 既存のハンドラー実装

### 9.3 技術スタック
- **Go**: 1.23.4
- **Echo**: v4（導入予定）
- **Huma**: v2（導入予定）
- **Gorilla Mux**: v1.8.1（置き換え予定）

### 9.4 参考リンク
- Echo公式ドキュメント: https://echo.labstack.com/
- Huma公式ドキュメント: https://huma.rocks/
- Huma Echoアダプター: https://pkg.go.dev/github.com/danielgtaylor/huma/v2/adapters/humaecho

