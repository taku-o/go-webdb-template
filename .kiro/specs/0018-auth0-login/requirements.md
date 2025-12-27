# Auth0ログイン機能要件定義書

## 1. 概要

### 1.1 プロジェクト情報
- **プロジェクト名**: go-webdb-template
- **Issue番号**: #30
- **Issueタイトル**: Auth0によるログイン機能を用意する
- **Feature名**: 0018-auth0-login
- **作成日**: 2025-01-27

### 1.2 目的
Auth0を使用してPARTNERアプリのアカウントでログインできる機能を実装する。ログイン成功時にJWTを取得し、HTTP-only Cookieに保存する。ログイン状態の表示とログアウト機能も実装する。

### 1.3 スコープ
- Auth0 SDK（@auth0/nextjs-auth0）の導入
- ログイン機能の実装（Auth0へのリダイレクト、コールバック処理）
- JWTの取得・保存（HTTP-only Cookieへの自動保存）
- ログイン状態の表示（ログイン済み/未ログインのUI切り替え）
- ログアウト機能の実装（HTTP-only CookieからJWTを削除）
- 環境別設定の管理（develop/staging/production）
- JWT取得機能の実装（`getAccessToken()`で取得可能にする）

**本実装の範囲外**:
- JWTを使ってAPIを叩く機能（次のissueで対応）
- アカウント情報のデータベース保存（別issueで対応）

## 2. 背景・現状分析

### 2.1 現在の実装
- **認証機能**: 認証機能は実装されていない
- **API呼び出し**: `NEXT_PUBLIC_API_KEY`を使用してAPIを呼び出している（`client/src/lib/api.ts`）
- **Auth0ダッシュボード側の設定**: 既に設定済み（`docs/Partner-Idp-Auth0-Login.md`参照）
  - 疑似PARTNER（IdP役）の作成済み
  - MINE側の接続設定（Enterprise Connection）済み
  - 相互認証の許可設定済み
  - アプリケーションへの紐付け済み

### 2.2 課題点
1. **ログイン機能がない**: ユーザーがPARTNERアプリのアカウントでログインできない
2. **認証状態の管理がない**: ログイン状態を確認・表示する機能がない
3. **JWTの取得・保存機能がない**: ログイン成功時にJWTを取得・保存する機能がない
4. **ログアウト機能がない**: ログアウトする機能がない

### 2.3 本実装による改善点
1. **ログイン機能の実現**: Auth0を使用してPARTNERアプリのアカウントでログインできるようになる
2. **JWTの取得・保存**: ログイン成功時にJWTを取得し、HTTP-only Cookieに安全に保存する
3. **認証状態の可視化**: ログイン状態を表示し、ユーザーが現在の状態を確認できる
4. **ログアウト機能**: ログアウトボタンから簡単にログアウトできる
5. **将来の拡張性**: JWT取得機能（`getAccessToken()`）を実装し、次のissueでAPI呼び出しに使用できる準備を整える

## 3. 機能要件

### 3.1 Auth0 SDKの導入

#### 3.1.1 パッケージのインストール
- `@auth0/nextjs-auth0`パッケージをインストール
- Next.js 14 (App Router)に対応したバージョンを使用

#### 3.1.2 API Routesの設定
- `client/src/app/api/auth/[...auth0]/route.ts`を作成
- Auth0 SDKのハンドラーをエクスポート

### 3.2 環境変数の設定

#### 3.2.1 クライアント側（Next.js）の環境変数
以下の環境変数を設定（`.env.local`、`.env.development`、`.env.production`）:

- `AUTH0_SECRET`: セッション暗号化用の秘密鍵（環境変数のみ、必須）
- `AUTH0_BASE_URL`: アプリケーションのベースURL（例: `http://localhost:3000`）
- `AUTH0_ISSUER_BASE_URL`: Auth0のドメイン（例: `https://{domain}.auth0.com`）
- `AUTH0_CLIENT_ID`: Client ID（公開情報、環境変数で管理）
- `AUTH0_CLIENT_SECRET`: Client Secret（機密情報、環境変数のみ、必須）

**注意**: Client Secretは外部サービスの情報のため、全環境（develop/staging/production）で環境変数のみで管理する。config.yamlには記載しない。

#### 3.2.2 サーバー側（Go）の環境変数
**今回の実装では不要**: サーバー側（Go）でAuth0を使用する機能は今回の実装範囲外です。
- クライアント側（Next.js）のみでAuth0を使用してログイン機能を実装します
- サーバー側でAuth0のJWTを検証する機能は次のissueで対応予定です
- そのため、`config/{env}/config.yaml`へのAuth0設定追加は今回の実装では行いません

### 3.3 ログイン機能の実装

#### 3.3.1 ログインボタンの実装
- トップページ（`client/src/app/page.tsx`）にログインボタンを追加
- または専用のログインページ（`client/src/app/login/page.tsx`）を作成
- 未ログイン時のみ表示

#### 3.3.2 ログイン処理の実装
- ログインボタンクリック時にAuth0のログインページにリダイレクト
- 実装方法:
  ```typescript
  <a href="/api/auth/login">Login</a>
  // または
  import { useUser } from '@auth0/nextjs-auth0/client'
  const { loginWithRedirect } = useUser()
  loginWithRedirect()
  ```

#### 3.3.3 コールバック処理の実装
- Auth0からリダイレクトされた後のコールバック処理
- Auth0 SDKが自動的に処理（`/api/auth/callback`）
- ログイン成功時にJWTを取得し、HTTP-only Cookieに保存

### 3.4 ログイン状態の表示

#### 3.4.1 認証状態の確認
- `useUser()`フックを使用して認証状態を確認
- 実装方法:
  ```typescript
  import { useUser } from '@auth0/nextjs-auth0/client'
  const { user, error, isLoading } = useUser()
  ```

#### 3.4.2 UIの切り替え
- **未ログイン時**: ログインボタンを表示
- **ログイン済み時**: ログアウトボタンとユーザー情報を表示（オプション）
- ローディング状態の表示（`isLoading`が`true`の間）

### 3.5 ログアウト機能の実装

#### 3.5.1 ログアウトボタンの実装
- ログイン済み時のみ表示
- トップページまたは専用のUIコンポーネントに配置

#### 3.5.2 ログアウト処理の実装
- Auth0 SDKの`/api/auth/logout`エンドポイントを使用
- HTTP-only CookieからJWT（セッション）を削除
- Auth0側のセッションも無効化（オプション）
- 実装方法:
  ```typescript
  // 方法1: ログアウトリンク
  <a href="/api/auth/logout">Logout</a>
  // 方法2: useUserフックを使用
  import { useUser } from '@auth0/nextjs-auth0/client'
  const { logout } = useUser()
  logout({ returnTo: '/' })
  ```

### 3.6 JWT取得機能の実装

#### 3.6.1 Server Componentsでの取得
- `getAccessToken()`を使用してJWTを取得
- 実装方法:
  ```typescript
  import { getAccessToken } from '@auth0/nextjs-auth0'
  const { accessToken } = await getAccessToken()
  ```

#### 3.6.2 Client Componentsでの取得
- `useUser()`フックと`getAccessToken()`を使用
- 実装方法:
  ```typescript
  import { useUser } from '@auth0/nextjs-auth0/client'
  const { user, getAccessToken } = useUser()
  const token = await getAccessToken()
  ```

**注意**: 本実装では、JWT取得機能を実装するが、API呼び出しでの使用は次のissueで対応する。

### 3.7 環境別設定の管理

#### 3.7.1 開発環境（develop）
- `.env.development`に環境変数を設定
- Auth0の開発用テナントの設定を使用

#### 3.7.2 ステージング環境（staging）
- `.env.staging`に環境変数を設定
- Auth0のステージング用テナントの設定を使用

#### 3.7.3 本番環境（production）
- `.env.production`に環境変数を設定
- Auth0の本番用テナントの設定を使用

## 4. 非機能要件

### 4.1 セキュリティ
- **JWTの安全な保存**: HTTP-only Cookieに保存し、JavaScriptからアクセスできないようにする
- **HTTPS必須**: 本番環境ではHTTPSを使用する（開発環境ではHTTPも可）
- **Client Secretの管理**: 環境変数のみで管理し、Gitにコミットしない
- **セッション管理**: Auth0 SDKが自動的にセッションを管理し、適切な有効期限を設定

### 4.2 パフォーマンス
- **認証状態の管理**: 認証状態を効率的に管理し、不要な再認証を避ける
- **ローディング状態の表示**: 認証状態の確認中は適切なローディング表示を行う

### 4.3 エラーハンドリング
- **ログイン失敗時の処理**: ログインに失敗した場合、適切なエラーメッセージを表示
- **ネットワークエラーの処理**: ネットワークエラーが発生した場合、適切なエラーハンドリングを行う
- **認証エラーの処理**: 認証エラーが発生した場合、ユーザーに分かりやすいメッセージを表示

### 4.4 ユーザビリティ
- **ログイン状態の明確な表示**: ユーザーが現在のログイン状態を一目で確認できる
- **簡単なログアウト操作**: ログアウトボタンから簡単にログアウトできる
- **適切なリダイレクト**: ログイン後、適切なページにリダイレクトする

## 5. 制約事項

### 5.1 実装範囲の制約
- **JWTを使ってAPIを叩く機能**: 次のissueで対応（本実装ではJWT取得・保存まで）
- **アカウント情報のデータベース保存**: 別issueで対応
- **現在のAPIキー方式**: `NEXT_PUBLIC_API_KEY`を使用したAPI呼び出しは維持（変更なし）

### 5.2 設定の制約
- **Client Secret**: 環境変数のみで管理（config.yamlには記載しない）
- **Auth0ダッシュボード側の設定**: 既に完了しているため、変更不要

### 5.3 技術的制約
- **Next.js App Router**: Next.js 14のApp Routerを使用
- **Auth0 SDK**: `@auth0/nextjs-auth0`を使用（Next.js App Router対応版）

## 6. 受け入れ基準

### 6.1 ログイン機能
- [ ] ログインボタンが表示される（未ログイン時）
- [ ] ログインボタンをクリックするとAuth0のログインページにリダイレクトされる
- [ ] ログイン処理が正常に動作する
- [ ] ログイン成功時にJWTが取得できる（`getAccessToken()`で取得可能）
- [ ] JWTがHTTP-only Cookieに保存される

### 6.2 ログイン状態の表示
- [ ] ログイン状態が正しく表示される（ログイン済み/未ログイン）
- [ ] 未ログイン時はログインボタンが表示される
- [ ] ログイン済み時はログアウトボタンが表示される
- [ ] ローディング状態が適切に表示される

### 6.3 ログアウト機能
- [ ] ログアウトボタンが表示される（ログイン済み時）
- [ ] ログアウトボタンをクリックするとログアウト処理が実行される
- [ ] ログアウト処理が正常に動作する（HTTP-only CookieからJWTが削除される）
- [ ] ログアウト後、未ログイン状態に戻る

### 6.4 環境別設定
- [ ] 開発環境（develop）で正しく動作する
- [ ] ステージング環境（staging）で正しく動作する（設定可能な場合）
- [ ] 本番環境（production）で正しく動作する（設定可能な場合）
- [ ] Client Secretが環境変数で正しく読み込まれる

### 6.5 JWT取得機能
- [ ] Server Componentsで`getAccessToken()`を使用してJWTを取得できる
- [ ] Client Componentsで`useUser()`と`getAccessToken()`を使用してJWTを取得できる

## 7. 影響範囲

### 7.1 新規追加が必要なファイル

#### クライアント側（Next.js）
- `client/src/app/api/auth/[...auth0]/route.ts`: Auth0 SDKのハンドラー
- `client/src/app/login/page.tsx`: 専用ログインページ（オプション）
- `client/src/components/LoginButton.tsx`: ログインボタンコンポーネント（オプション）
- `client/src/components/LogoutButton.tsx`: ログアウトボタンコンポーネント（オプション）
- `client/.env.local`: 環境変数設定ファイル（開発用）
- `client/.env.development`: 環境変数設定ファイル（開発環境用）
- `client/.env.staging`: 環境変数設定ファイル（ステージング環境用）
- `client/.env.production`: 環境変数設定ファイル（本番環境用）

#### ドキュメント
- `.kiro/specs/0018-auth0-login/requirements.md`: 本要件定義書
- `.kiro/specs/0018-auth0-login/spec.json`: 仕様書メタデータ

### 7.2 変更が必要なファイル

#### クライアント側（Next.js）
- `client/src/app/page.tsx`: ログイン/ログアウトUIの追加
- `client/src/app/layout.tsx`: Auth0 SDKのProvider設定（必要に応じて）
- `client/package.json`: `@auth0/nextjs-auth0`パッケージの追加

#### サーバー側（Go）
- **変更なし**: サーバー側（Go）でのAuth0設定は今回の実装範囲外です

### 7.3 変更なしのファイル
- `client/src/lib/api.ts`: apiClientの修正は次のissueで対応

## 8. 実装上の注意事項

### 8.1 Next.js App RouterでのAuth0実装
- Next.js 14のApp Routerに対応した`@auth0/nextjs-auth0`を使用
- API Routesの設定（`[...auth0]/route.ts`）を正しく実装する
- Server ComponentsとClient Componentsで適切にAuth0 SDKを使用する

### 8.2 環境変数の管理
- Client Secretは環境変数のみで管理し、Gitにコミットしない
- `.env.local`は`.gitignore`に追加されていることを確認
- 環境別の環境変数ファイル（`.env.development`、`.env.staging`、`.env.production`）を適切に管理

### 8.3 セキュリティ考慮事項
- Client Secretの取り扱いに注意（外部サービスの情報のため、全環境で環境変数のみで管理）
- HTTP-only Cookieの使用により、XSS攻撃からJWTを保護
- 本番環境ではHTTPSを使用する

### 8.4 JWTの取得方法
- Server Components: `getAccessToken()`を使用
- Client Components: `useUser()`フックと`getAccessToken()`を使用
- 本実装では、JWT取得機能を実装するが、API呼び出しでの使用は次のissueで対応

### 8.5 ログイン状態の管理
- `useUser()`フックを使用して認証状態を確認
- ローディング状態（`isLoading`）を適切に処理
- エラー状態（`error`）を適切に処理

### 8.6 ログアウト処理の実装
- `/api/auth/logout`エンドポイントを使用
- HTTP-only CookieからJWTを削除
- ログアウト後のリダイレクト先を適切に設定

### 8.7 本実装の範囲
- **本実装の範囲**: ログインしてJWTを取得・保存するまで、ログイン状態の表示、ログアウト機能
- **次のissueで対応**: JWTを使ってAPIを叩く機能（apiClientの修正）

## 9. 参考情報

### 9.1 関連Issue
- GitHub Issue #30: Auth0によるログイン機能を用意する

### 9.2 既存ドキュメント
- `docs/Partner-Idp-Auth0-Login.md`: Auth0ダッシュボード側の設定手順
- `.kiro/specs/0016-fix-tablesplit/requirements.md`: 既存の要件定義書（フォーマット参考）

### 9.3 技術スタック
- **Auth0 SDK**: `@auth0/nextjs-auth0`（Next.js App Router対応）
- **Next.js**: 14 (App Router)
- **TypeScript**: 5+
- **参考資料**:
  - [Auth0 Next.js SDK Documentation](https://auth0.com/docs/quickstart/webapp/nextjs)
  - [Auth0 Next.js SDK - Getting an Access Token](https://auth0.com/docs/quickstart/webapp/nextjs/01-login#get-an-access-token)
