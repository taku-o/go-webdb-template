# APIドキュメント改善要件定義書

## 1. 概要

### 1.1 プロジェクト情報
- **プロジェクト名**: go-webdb-template
- **Issue番号**: #37
- **Issueタイトル**: 8080/docsのAPIドキュメントにAuthorization Header、APIのアクセスレベルの情報を追加する
- **Feature名**: 0020-apidoc
- **作成日**: 2025-01-27

### 1.2 目的
Humaフレームワークで生成されるAPIドキュメント（http://localhost:8080/docs）に、Authorization HeaderとAPIアクセスレベルの情報を追加する。これにより、APIドキュメント上で認証が必要であることが明確になり、Swagger UIから直接APIをテストできるようになる。

### 1.3 スコープ
- HumaのSecurityScheme設定によるAuthorization Headerの表示
- 各エンドポイントへのSecurityプロパティの適用
- TagsによるAPIアクセスレベルの表示（Public API / Private API）
- 既存の認証ロジックやAPI動作への影響なし（ドキュメント表示のみの改善）

**本実装の範囲外**:
- 認証ロジックの変更（既存の認証ミドルウェアは変更しない）
- APIエンドポイントの追加・削除
- アクセス制御ロジックの変更

## 2. 背景・現状分析

### 2.1 現在の実装
- **Huma APIドキュメント**: `http://localhost:8080/docs`でアクセス可能
- **認証機能**: Issue #19（0019-auth0-apicall）で実装済み
  - Auth0 JWTとPublic API Key JWTの両方をサポート
  - 認証ミドルウェア（`server/internal/auth/middleware.go`）でJWT検証を実施
  - アクセスレベル（public/private）のチェック機能あり
- **APIエンドポイント**: 
  - Userエンドポイント（5つ）：全てpublicレベル
  - Postエンドポイント（6つ）：全てpublicレベル
  - Todayエンドポイント（1つ）：privateレベル
- **Huma設定**: `server/internal/api/router/router.go`で`huma.DefaultConfig()`を使用

### 2.2 課題点
1. **Authorization Headerの未表示**: Request Sample（curl）に`Authorization: Bearer <token>`が表示されない
2. **Swagger UIの「Authorize」ボタンがない**: ブラウザ上でトークンを入力してAPIをテストできない
3. **APIアクセスレベルの未表示**: 各APIがpublicレベルかprivateレベルかをドキュメントで確認できない
4. **API利用者への情報不足**: 認証が必要なAPIかどうか、どのレベルの認証が必要かをドキュメントから判断できない

### 2.3 本実装による改善点
1. **Authorization Headerの表示**: Request Sampleに自動的に`-H 'Authorization: Bearer <token>'`が追加される
2. **Swagger UIの「Authorize」ボタン**: ブラウザ上でトークンを入力してAPIをテストできるようになる
3. **APIアクセスレベルの表示**: Tagsで「Public API」または「Private API」が表示される
4. **API利用者への情報提供**: ドキュメントから認証要件とアクセスレベルを確認できる

## 3. 機能要件

### 3.1 SecuritySchemeの設定

#### 3.1.1 Bearer認証スキームの定義
- `huma.Config`の`Components.SecuritySchemes`にBearer認証スキームを追加
- 設定内容：
  - Type: "http"
  - Scheme: "bearer"
  - BearerFormat: "JWT"
- 実装場所：`server/internal/api/router/router.go`の`humaConfig`設定部分

#### 3.1.2 実装方法
```go
config := huma.DefaultConfig("go-webdb-template API", "1.0.0")

// セキュリティスキームの定義
config.Components = &huma.Components{
    SecuritySchemes: map[string]*huma.SecurityScheme{
        "bearerAuth": {
            Type:         "http",
            Scheme:       "bearer",
            BearerFormat: "JWT",
        },
    },
}
```

### 3.2 各エンドポイントへのSecurity適用

#### 3.2.1 Securityプロパティの追加
- すべてのエンドポイントに`Security`プロパティを追加
- 対象エンドポイント：
  - Userエンドポイント（5つ）：全てpublic
    - POST /api/users
    - GET /api/users/{id}
    - GET /api/users
    - PUT /api/users/{id}
    - DELETE /api/users/{id}
  - Postエンドポイント（6つ）：全てpublic
    - POST /api/posts
    - GET /api/posts/{id}
    - GET /api/posts
    - PUT /api/posts/{id}
    - DELETE /api/posts/{id}
    - GET /api/user-posts
  - Todayエンドポイント（1つ）：private
    - GET /api/today

#### 3.2.2 実装方法
各`huma.Register()`の`huma.Operation`構造体に`Security`プロパティを追加：

```go
huma.Register(api, huma.Operation{
    OperationID: "get-user",
    Method:      http.MethodGet,
    Path:        "/api/users/{id}",
    Summary:     "ユーザーを取得",
    Tags:        []string{"users", "Public API"},
    Security: []map[string][]string{
        {"bearerAuth": []},
    },
}, ...)
```

#### 3.2.3 実装場所
- `server/internal/api/handler/user_handler.go` - Userエンドポイント（5箇所）
- `server/internal/api/handler/post_handler.go` - Postエンドポイント（6箇所）
- `server/internal/api/handler/today_handler.go` - Todayエンドポイント（1箇所）

### 3.3 APIアクセスレベルの表示

#### 3.3.1 Tagsによるアクセスレベル表示
- 既存のTags（"users", "posts", "today"）に加えて、アクセスレベルのTagを追加
- Tag名：
  - "Public API" - publicレベルのAPI（Public API Key JWTとAuth0 JWTの両方でアクセス可能）
  - "Private API" - privateレベルのAPI（Auth0 JWTのみでアクセス可能）

#### 3.3.2 実装方法
既存のTags配列にアクセスレベルのTagを追加：

```go
// Public APIの場合
Tags: []string{"users", "Public API"}

// Private APIの場合
Tags: []string{"today", "Private API"}
```

#### 3.3.3 各エンドポイントのTag設定
- Userエンドポイント（5つ）：`[]string{"users", "Public API"}`
- Postエンドポイント（6つ）：`[]string{"posts", "Public API"}`
- Todayエンドポイント（1つ）：`[]string{"today", "Private API"}`

## 4. 非機能要件

### 4.1 互換性
- **既存機能への影響なし**: 認証ロジック、API動作、アクセス制御は一切変更しない
- **後方互換性の維持**: 既存のAPIクライアントへの影響なし
- **既存のTags維持**: 機能的なグループ化タグ（"users", "posts", "today"）は維持

### 4.2 ドキュメント品質
- **明確な情報表示**: APIドキュメント上で認証要件とアクセスレベルが明確に分かる
- **一貫性**: すべてのエンドポイントに一貫して適用
- **可読性**: Swagger UIで視覚的に分かりやすい表示

### 4.3 メンテナンス性
- **設定の一元管理**: SecuritySchemeは`router.go`で一元管理
- **明確な実装**: 各エンドポイントでSecurityとTagsが明示的に設定されている

## 5. 制約事項

### 5.1 実装範囲の制約
- **認証ロジックの変更禁止**: 既存の認証ミドルウェア（`server/internal/auth/middleware.go`）は変更しない
- **API動作の変更禁止**: 既存のAPIエンドポイントの動作は変更しない
- **アクセス制御の変更禁止**: 既存の`auth.CheckAccessLevel()`によるアクセス制御は維持

### 5.2 技術的制約
- **Humaフレームワーク**: 既存のHuma v2の機能を使用
- **Go言語**: 既存のGo言語実装を維持
- **OpenAPI仕様**: Humaが生成するOpenAPI仕様に準拠

### 5.3 設定の制約
- **SecurityScheme名**: "bearerAuth"を使用（Issue #37の記載に従う）
- **Tag名**: "Public API"と"Private API"を使用（大文字小文字を区別）

## 6. 受け入れ基準

### 6.1 SecurityScheme設定
- [ ] `huma.Config`の`Components.SecuritySchemes`にBearer認証スキームが定義されている
- [ ] SecuritySchemeの設定内容が正しい（Type: "http", Scheme: "bearer", BearerFormat: "JWT"）

### 6.2 各エンドポイントへのSecurity適用
- [ ] すべてのエンドポイント（12個）に`Security`プロパティが追加されている
- [ ] Securityプロパティの形式が正しい（`[]map[string][]string{{"bearerAuth": []}}`）

### 6.3 APIアクセスレベルの表示
- [ ] すべてのエンドポイントに適切なアクセスレベルTagが追加されている
  - Userエンドポイント（5つ）："Public API"
  - Postエンドポイント（6つ）："Public API"
  - Todayエンドポイント（1つ）："Private API"
- [ ] 既存の機能タグ（"users", "posts", "today"）が維持されている

### 6.4 APIドキュメントの表示確認
- [ ] `http://localhost:8080/docs`にアクセスして、Request Sample（curl）に`-H 'Authorization: Bearer <token>'`が表示される
- [ ] Swagger UIに「Authorize」ボタンが表示される
- [ ] 「Authorize」ボタンをクリックしてトークンを入力できる
- [ ] 各APIにアクセスレベルのTag（"Public API" / "Private API"）が表示される
- [ ] 既存の機能タグ（"users", "posts", "today"）も表示される

### 6.5 既存機能の動作確認
- [ ] 既存のAPI動作に影響がないことを確認（認証、アクセス制御が正常に動作）
- [ ] 既存のAPIクライアントへの影響がないことを確認

## 7. 影響範囲

### 7.1 変更が必要なファイル

#### サーバー側（Go）
- `server/internal/api/router/router.go` - SecurityScheme設定の追加
- `server/internal/api/handler/user_handler.go` - SecurityとTags追加（5箇所）
- `server/internal/api/handler/post_handler.go` - SecurityとTags追加（6箇所）
- `server/internal/api/handler/today_handler.go` - SecurityとTags追加（1箇所）

### 7.2 変更なしのファイル
- `server/internal/auth/middleware.go` - 認証ミドルウェア（変更なし）
- `server/internal/auth/jwt.go` - JWT検証ロジック（変更なし）
- その他の認証関連ファイル（変更なし）

### 7.3 新規ファイル
- `.kiro/specs/0020-apidoc/requirements.md` - 本要件定義書
- `.kiro/specs/0020-apidoc/spec.json` - 仕様書メタデータ

## 8. 実装上の注意事項

### 8.1 SecurityScheme設定
- `huma.DefaultConfig()`で作成したConfigに`Components`を追加する際は、既存の設定を上書きしないよう注意
- SecurityScheme名は"bearerAuth"を使用（Issue #37の記載に従う）

### 8.2 Securityプロパティの適用
- すべてのエンドポイントに一貫して適用する
- Securityプロパティの形式は`[]map[string][]string{{"bearerAuth": []}}`とする
- 空の配列`[]`を指定することで、すべてのスコープを許可する

### 8.3 Tagsによるアクセスレベル表示
- 既存の機能タグ（"users", "posts", "today"）は維持し、アクセスレベルTagを追加する
- Tag名は大文字小文字を区別するため、"Public API"と"Private API"を正確に記述する
- 各エンドポイントのアクセスレベルは、既存の`auth.CheckAccessLevel()`の呼び出しから判断する

### 8.4 既存機能との互換性
- 認証ロジックやAPI動作は一切変更しない
- 既存の`auth.CheckAccessLevel()`によるアクセス制御は維持
- ドキュメント表示のみを改善する

### 8.5 動作確認
- 実装後は`http://localhost:8080/docs`にアクセスして、以下を確認：
  - Request SampleにAuthorization Headerが表示される
  - Swagger UIに「Authorize」ボタンが表示される
  - 各APIにアクセスレベルのTagが表示される
  - 既存のAPI動作に影響がない

## 9. 参考情報

### 9.1 関連Issue
- GitHub Issue #37: 8080/docsのAPIドキュメントにAuthorization Header、APIのアクセスレベルの情報を追加する（0020-apidoc）
- GitHub Issue #40: Auth0から受け取ったJWTをAPIサーバーとの通信で利用する（0019-auth0-apicall）

### 9.2 既存ドキュメント
- `.kiro/specs/0019-auth0-apicall/requirements.md`: Auth0 API呼び出し機能の要件定義書
- `.kiro/specs/0013-echo-huma/requirements.md`: Echo + Huma導入の要件定義書

### 9.3 技術スタック
- **Go言語**: 1.24+
- **Humaフレームワーク**: v2
- **OpenAPI仕様**: Humaが自動生成

### 9.4 参考資料
- [Huma Documentation](https://huma.rocks/)
- [OpenAPI Security Schemes](https://swagger.io/specification/#security-scheme-object)
- Issue #37のコメントに記載されている実装方法
