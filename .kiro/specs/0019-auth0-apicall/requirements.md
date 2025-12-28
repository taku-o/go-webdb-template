# Auth0 API呼び出し機能要件定義書

## 1. 概要

### 1.1 プロジェクト情報
- **プロジェクト名**: go-webdb-template
- **Issue番号**: #40
- **Issueタイトル**: Auth0から受け取ったJWTをAPIサーバーとの通信で利用する
- **Feature名**: 0019-auth0-apicall
- **作成日**: 2025-01-27

### 1.2 目的
Auth0でログインしたユーザーが、Auth0から受け取ったJWTを使用してAPIサーバーと通信できるようにする。未ログイン時は既存の`NEXT_PUBLIC_API_KEY`を使用してAPIサーバーと通信する。サーバー側でJWTの種類（Auth0 JWT / Public API Key JWT）を判別し、適切に検証する。APIの公開レベル（public/private）を定義し、適切なアクセス制御を実装する。

### 1.3 スコープ
- クライアント側: Auth0 JWTとPublic API Keyの切り替え、API呼び出し
- サーバー側: JWT種類の判別、Auth0 JWT検証、Public API Key JWT検証、API公開レベルの実装
- 新規private APIエンドポイントの追加（動作確認用）

**本実装の範囲外**:
- Auth0ログイン機能の実装（Issue #30で対応済み）
- アカウント情報のデータベース保存（別issueで対応）

## 2. 背景・現状分析

### 2.1 現在の実装
- **Auth0ログイン機能**: Issue #30で実装済み（`getAccessToken()`でJWT取得可能）
- **API呼び出し**: `NEXT_PUBLIC_API_KEY`を使用してAPIを呼び出している（`client/src/lib/api.ts`）
- **サーバー側JWT検証**: Public API Key JWTのみ検証可能（`server/internal/auth/jwt.go`）
  - HS256（HMAC署名）を使用
  - `iss: "go-webdb-template"`を検証
- **APIエンドポイント**: `/api/users`, `/api/posts`, `/api/user-posts`が実装済み

### 2.2 課題点
1. **Auth0 JWTの未対応**: Auth0でログインしても、API呼び出しでAuth0 JWTを使用できない
2. **JWT種類の判別機能がない**: サーバー側でAuth0 JWTとPublic API Key JWTを判別できない
3. **Auth0 JWT検証機能がない**: サーバー側でAuth0 JWTを検証できない（RS256署名の検証が必要）
4. **API公開レベルの定義がない**: すべてのAPIが同じアクセスレベルで公開されている
5. **アクセス制御の不足**: Public API KeyとAuth0 JWTで異なるアクセス権限を設定できない

### 2.3 本実装による改善点
1. **Auth0 JWTの利用**: Auth0でログインしたユーザーが、Auth0 JWTを使用してAPIを呼び出せるようになる
2. **JWT種類の判別**: サーバー側でJWTの種類を自動判別し、適切な検証方法を適用
3. **Auth0 JWT検証**: Auth0の公開鍵（JWKS）を使用してRS256署名を検証
4. **API公開レベルの定義**: public/privateの2段階のアクセスレベルを定義
5. **適切なアクセス制御**: Public API KeyはpublicなAPIのみ、Auth0 JWTはpublicとprivateの両方にアクセス可能

## 3. 機能要件

### 3.1 クライアント側要件

#### 3.1.1 JWT取得と切り替え機能
- Auth0でログイン中は、`getAccessToken()`を使用してAuth0から受け取ったJWTを取得
- 未ログイン時は、`NEXT_PUBLIC_API_KEY`を使用
- `ApiClient`クラスの修正: JWT取得ロジックの追加
  - ログイン状態を確認（`useUser()`フックを使用）
  - ログイン中は`getAccessToken()`でJWTを取得
  - 未ログイン時は`NEXT_PUBLIC_API_KEY`を使用
  - すべてのAPIリクエストに適切なJWTを`Authorization: Bearer <JWT>`ヘッダーで送信

#### 3.1.2 動作確認用UIの追加
- privateなAPI（`/api/today`）を呼び出すリンク/ボタンを追加（必須）
- `NEXT_PUBLIC_API_KEY`使用時にそのリンクをクリックすると、APIサーバーからエラー（403 Forbidden）が返ることを確認できる
- private/public判定の実装例として機能する

### 3.2 サーバー側要件

#### 3.2.1 JWT種類の判別機能
- JWTのペイロード（`iss`など）を見て、Auth0 JWTかPublic API Key JWTかを判別
- 判別方法:
  - Auth0 JWT: `iss`がAuth0のドメイン（例: `https://{domain}.auth0.com/`）
  - Public API Key JWT: `iss`が`"go-webdb-template"`

#### 3.2.2 Auth0 JWT検証機能
- Auth0 JWTの検証: Auth0の公開鍵（JWKS）を使用した検証（RS256）
- JWKS取得方法:
  - `github.com/MicahParks/keyfunc`ライブラリを使用
  - `iss`（Issuer）から自動的に`.well-known/jwks.json`のURLを導出
  - 環境ごとに異なる場合は、`config/{env}/config.yaml`に`AUTH0_ISSUER_BASE_URL`を追加
  - JWKS URLは`{AUTH0_ISSUER_BASE_URL}/.well-known/jwks.json`として導出
- JWKSキャッシュ: メモリキャッシュを使用（`keyfunc`ライブラリのキャッシュ機能を利用）
  - キャッシュ期間: 12時間（`RefreshInterval: time.Hour * 12`）
  - キャッシュ設定オプション:
    ```go
    options := keyfunc.Options{
        RefreshInterval:   time.Hour * 12,   // 12時間ごとに定期更新
        RefreshRateLimit:  time.Minute * 5,  // 再取得は最低5分あける（DoS対策）
        RefreshTimeout:    time.Second * 10, // 取得時のタイムアウト
        RefreshUnknownKID: true,             // 未知のKIDが来たら再取得する（重要！）
    }
    ```

#### 3.2.3 Public API Key JWT検証機能
- 既存の実装を維持
- HS256（HMAC署名）を使用
- `iss: "go-webdb-template"`を確認

#### 3.2.4 API公開レベルの定義
- APIの公開レベルを定義: public/private
- 定義方法（優先順位）:
  1. **コード内での定義（優先）**: 
     - `huma.Operation`構造体に直接追加できる場合は、各エンドポイント登録時に公開レベルを指定
     - 例: `AccessLevel: "public"` または `AccessLevel: "private"` のようなフィールドを追加
     - 実現困難な場合は、別のマップや構造体で管理（例: `map[string]string`でパスと公開レベルの対応を管理）
  2. **代替案（実現困難な場合）**:
     - 設定ファイル（`config/{env}/config.yaml`）で公開レベルを定義
     - 別のGoファイル（例: `server/internal/api/access_level.go`）でマップとして定義
- 既存APIの公開レベル:
  - 既存のAPI（`/api/users`, `/api/posts`, `/api/user-posts`）は全て**public**として定義
- 新規追加APIの公開レベル:
  - 新規追加する動作確認用API（`/api/today`）は**private**として定義

#### 3.2.5 アクセス制御機能
- 認証ミドルウェア（`server/internal/auth/middleware.go`）でJWT検証後に、APIの公開レベルをチェック
- Public API Key JWTの場合:
  - publicなAPIへのアクセス: 許可
  - privateなAPIへのアクセス: 403 Forbiddenを返す
- Auth0 JWTの場合:
  - publicなAPIへのアクセス: 許可
  - privateなAPIへのアクセス: 許可

#### 3.2.6 新規追加するprivate APIエンドポイント

**エンドポイント**: `GET /api/today`

**機能**: 今日の日付をYYYY-MM-DDのフォーマットで返す

**レスポンス形式**:
```json
{
  "date": "2025-01-27"
}
```

**公開レベル**: private（Auth0 JWTでのみアクセス可能）

**実装場所**: 新規ハンドラーファイル `server/internal/api/handler/today_handler.go`として実装
- テスト用のAPIではなく、private/public判定の実装例として機能する

### 3.3 環境別設定

#### 3.3.1 サーバー側（Go）の設定
- `config/{env}/config.yaml`に`AUTH0_ISSUER_BASE_URL`を追加（公開情報のため設定ファイルで管理）
  - develop: `https://dev-oaa5vtzmld4dsxtd.jp.auth0.com`
  - staging: （設定値は環境に応じて決定）
  - production: （設定値は環境に応じて決定）

#### 3.3.2 クライアント側（Next.js）の環境変数
- **変更なし**: クライアント側の環境変数はIssue #30で設定済み

## 4. 非機能要件

### 4.1 セキュリティ
- **JWTの安全な検証**: Auth0 JWTはRS256署名を適切に検証
- **JWKSの安全な取得**: HTTPSを使用してJWKSを取得
- **アクセス制御の厳格な実装**: Public API KeyはprivateなAPIにアクセスできない
- **エラーメッセージの適切な管理**: セキュリティ上の情報漏洩を防ぐため、エラーメッセージは適切に管理

### 4.2 パフォーマンス
- **JWKSキャッシュ**: メモリキャッシュを使用してJWKS取得のオーバーヘッドを削減
- **効率的なJWT検証**: JWT種類の判別と検証を効率的に実行
- **不要なネットワークリクエストの回避**: JWKSキャッシュにより、不要なJWKS取得リクエストを回避

### 4.3 エラーハンドリング
- **JWT検証エラー**: JWT検証に失敗した場合、適切なエラーレスポンス（401 Unauthorized）を返す
- **アクセス制御エラー**: privateなAPIにPublic API Keyでアクセスした場合、適切なエラーレスポンス（403 Forbidden）を返す
- **JWKS取得エラー**: JWKS取得に失敗した場合、適切なエラーハンドリングとリトライ機能を実装
- **ネットワークエラー**: ネットワークエラーが発生した場合、適切なエラーハンドリングを行う

### 4.4 ユーザビリティ
- **シームレスな切り替え**: ログイン/ログアウト時に、API呼び出しで使用するJWTが自動的に切り替わる
- **明確なエラーメッセージ**: アクセス権限がない場合、ユーザーに分かりやすいエラーメッセージを表示

## 5. 制約事項

### 5.1 実装範囲の制約
- **Auth0ログイン機能**: Issue #30で実装済みの機能を利用（変更なし）
- **既存のAPIキー方式**: `NEXT_PUBLIC_API_KEY`を使用したAPI呼び出しは維持（変更なし）
- **既存のAPIエンドポイント**: 既存のAPIエンドポイントの動作は変更しない（公開レベルのみ追加）

### 5.2 設定の制約
- **環境別設定**: develop、staging、productionで異なるAuth0テナントを使用する可能性がある
- **JWKSエンドポイント**: 環境ごとに異なるJWKSエンドポイントを使用

### 5.3 技術的制約
- **Go言語**: 既存のGo言語実装を維持
- **Next.js**: 既存のNext.js実装を維持
- **JWT検証ライブラリ**: 既存の`github.com/golang-jwt/jwt/v5`を維持し、Auth0 JWT検証用に追加のライブラリを導入

## 6. 受け入れ基準

### 6.1 クライアント側機能
- [ ] Auth0でログイン中は、Auth0から受け取ったJWTをAPI呼び出しで使用できる
- [ ] 未ログイン時は、`NEXT_PUBLIC_API_KEY`をAPI呼び出しで使用できる
- [ ] ログイン/ログアウト時に、API呼び出しで使用するJWTが自動的に切り替わる
- [ ] privateなAPI（`/api/today`）を呼び出すリンク/ボタンが表示される
- [ ] `NEXT_PUBLIC_API_KEY`使用時に`/api/today`を呼び出すと、エラー（403 Forbidden）が返る

### 6.2 サーバー側機能
- [ ] JWTの種類（Auth0 JWT / Public API Key JWT）を正しく判別できる
- [ ] Auth0 JWTを正しく検証できる（RS256署名の検証）
- [ ] Public API Key JWTを正しく検証できる（既存機能の維持）
- [ ] JWKSをメモリキャッシュに保存し、効率的に利用できる
- [ ] APIの公開レベル（public/private）を正しく定義できる
- [ ] Public API Key JWTでprivateなAPIにアクセスした場合、403 Forbiddenを返す
- [ ] Auth0 JWTでpublicなAPIにアクセスした場合、正常に処理される
- [ ] Auth0 JWTでprivateなAPIにアクセスした場合、正常に処理される
- [ ] 新規追加した`/api/today`エンドポイントが正しく動作する

### 6.3 環境別設定
- [ ] 開発環境（develop）で正しく動作する
- [ ] 環境ごとに異なる`AUTH0_ISSUER_BASE_URL`が正しく読み込まれる
- [ ] staging環境とproduction環境の設定ファイルに`AUTH0_ISSUER_BASE_URL`を追加できる（動作確認は不要）

## 7. 影響範囲

### 7.1 新規追加が必要なファイル

#### クライアント側（Next.js）
- `client/src/components/TodayApiButton.tsx`: private API（`/api/today`）を呼び出すボタンコンポーネント

#### サーバー側（Go）
- `server/internal/api/handler/today_handler.go`: 新規追加するprivate APIのハンドラー（`/api/today`）
- `server/internal/auth/auth0_validator.go`: Auth0 JWT検証機能（JWKS取得、キャッシュ、検証）
- `server/internal/api/access_level.go`: API公開レベルの定義（`huma.Operation`に直接追加できない場合）

#### 設定ファイル
- `config/{env}/config.yaml`: `AUTH0_ISSUER_BASE_URL`の追加（環境ごとに異なる場合）

#### ドキュメント
- `.kiro/specs/0019-auth0-apicall/requirements.md`: 本要件定義書
- `.kiro/specs/0019-auth0-apicall/spec.json`: 仕様書メタデータ

### 7.2 変更が必要なファイル

#### クライアント側（Next.js）
- `client/src/lib/api.ts`: JWT取得ロジックの追加（Auth0 JWTとPublic API Keyの切り替え）
- `client/src/app/page.tsx`: private API（`/api/today`）を呼び出すリンク/ボタンの追加

#### サーバー側（Go）
- `server/internal/auth/jwt.go`: JWT種類の判別機能の追加
- `server/internal/auth/middleware.go`: アクセス制御機能の追加（API公開レベルのチェック）
- `server/internal/api/handler/user_handler.go`: 既存APIの公開レベル定義（public）
- `server/internal/api/handler/post_handler.go`: 既存APIの公開レベル定義（public）
- `server/internal/config/config.go`: `AUTH0_ISSUER_BASE_URL`設定の追加
- `server/go.mod`: Auth0 JWT検証用ライブラリの追加（`github.com/MicahParks/keyfunc`）

### 7.3 変更なしのファイル
- `client/src/app/api/auth/[...auth0]/route.ts`: Auth0 SDKのハンドラー（Issue #30で実装済み）
- その他のAuth0ログイン機能関連ファイル（Issue #30で実装済み）

## 8. 実装上の注意事項

### 8.1 JWT種類の判別
- JWTのペイロードをパースして`iss`を確認する際は、署名検証前にパースする（`ParseUnverified`を使用）
- `iss`の形式を正しく判定する（Auth0のドメインパターンを考慮）

### 8.2 Auth0 JWT検証
- `github.com/MicahParks/keyfunc`ライブラリを使用し、実装の複雑さを軽減
- JWKSキャッシュの実装は、ライブラリのキャッシュ機能を利用
- キャッシュ設定: 12時間ごとの定期更新、5分間隔の再取得制限、10秒のタイムアウト、未知のKIDが来たら再取得
- 環境ごとに異なるJWKSエンドポイントに対応（設定ファイルで`AUTH0_ISSUER_BASE_URL`を管理）

### 8.3 API公開レベルの定義
- `huma.Operation`構造体に直接追加できる場合は、各エンドポイント登録時に公開レベルを指定
- 実現困難な場合は、別のマップや構造体で管理し、ミドルウェアで参照
- 既存APIの公開レベルは全てpublicとして定義（後方互換性の維持）

### 8.4 アクセス制御の実装
- 認証ミドルウェアでJWT検証後に、APIの公開レベルをチェック
- Public API Key JWTでprivateなAPIにアクセスした場合、適切なエラーレスポンス（403 Forbidden）を返す
- エラーメッセージは適切に管理し、セキュリティ上の情報漏洩を防ぐ

### 8.5 新規追加するprivate API
- `/api/today`エンドポイントはprivate/public判定の実装例として機能するため、シンプルな実装とする
- 日付の取得はサーバー側の現在時刻を使用
- レスポンス形式はJSON形式で統一

### 8.6 既存機能との互換性
- 既存のPublic API Key JWT検証機能は変更しない
- 既存のAPIエンドポイントの動作は変更しない（公開レベルのみ追加）
- 既存のクライアント側のAPI呼び出しは、JWT取得ロジックの追加のみで動作する

## 9. 参考情報

### 9.1 関連Issue
- GitHub Issue #30: Auth0によるログイン機能を用意する（0018-auth0-login）
- GitHub Issue #40: Auth0から受け取ったJWTをAPIサーバーとの通信で利用する（0019-auth0-apicall）

### 9.2 既存ドキュメント
- `.kiro/specs/0018-auth0-login/requirements.md`: Auth0ログイン機能の要件定義書
- `.kiro/specs/0011-public-api-key/requirements.md`: Public API Key機能の要件定義書

### 9.3 技術スタック
- **Go言語**: 1.24+
- **JWT検証ライブラリ**: `github.com/golang-jwt/jwt/v5`（既存）
- **Auth0 JWT検証ライブラリ**: `github.com/MicahParks/keyfunc`
- **Next.js**: 14 (App Router)
- **Auth0 SDK**: `@auth0/nextjs-auth0`（Issue #30で導入済み）

### 9.4 参考資料
- [Auth0 Next.js SDK Documentation](https://auth0.com/docs/quickstart/webapp/nextjs)
- [Auth0 Next.js SDK - Getting an Access Token](https://auth0.com/docs/quickstart/webapp/nextjs/01-login#get-an-access-token)
- [JWKS (JSON Web Key Set) - Auth0](https://auth0.com/docs/secure/tokens/json-web-tokens/json-web-key-sets)
- [github.com/MicahParks/keyfunc](https://github.com/MicahParks/keyfunc)
