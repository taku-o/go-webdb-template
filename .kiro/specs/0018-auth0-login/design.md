# Auth0ログイン機能設計書

## 1. 概要

### 1.1 設計の目的
要件定義書に基づき、Auth0を使用してPARTNERアプリのアカウントでログインできる機能を実装する詳細設計を定義する。ログイン成功時にJWTを取得し、HTTP-only Cookieに保存する。ログイン状態の表示とログアウト機能も実装する。

### 1.2 設計の範囲
- Auth0 SDK（@auth0/nextjs-auth0）の導入と設定
- ログイン機能の実装（Auth0へのリダイレクト、コールバック処理）
- JWTの取得・保存（HTTP-only Cookieへの自動保存）
- ログイン状態の表示（ログイン済み/未ログインのUI切り替え）
- ログアウト機能の実装（HTTP-only CookieからJWTを削除）
- 環境別設定の管理（develop/staging/production）
- JWT取得機能の実装（`getAccessToken()`で取得可能にする）

**本設計の範囲外**:
- JWTを使ってAPIを叩く機能（次のissueで対応）
- アカウント情報のデータベース保存（別issueで対応）

### 1.3 設計方針
- **Auth0 SDKの活用**: Next.js App Routerに対応した`@auth0/nextjs-auth0`を使用し、標準的な実装パターンに従う
- **セキュリティ優先**: JWTはHTTP-only Cookieに保存し、Client Secretは環境変数のみで管理
- **環境別設定**: 開発、ステージング、本番環境で適切に設定を切り替え可能にする
- **将来の拡張性**: JWT取得機能を実装し、次のissueでAPI呼び出しに使用できる準備を整える
- **既存システムとの互換性**: 既存のAPIキー方式（`NEXT_PUBLIC_API_KEY`）は維持し、変更しない

## 2. アーキテクチャ設計

### 2.1 ディレクトリ構造

#### 2.1.1 変更前の構造
```
client/
├── src/
│   ├── app/
│   │   ├── layout.tsx
│   │   ├── page.tsx
│   │   ├── users/
│   │   ├── posts/
│   │   └── user-posts/
│   ├── lib/
│   │   └── api.ts
│   └── types/
├── package.json
└── .env.local (存在しない)
```

#### 2.1.2 変更後の構造
```
client/
├── src/
│   ├── app/
│   │   ├── layout.tsx                    # 修正: UserProviderの追加
│   │   ├── page.tsx                      # 修正: ログイン/ログアウトUIの追加
│   │   ├── api/
│   │   │   └── auth/
│   │   │       └── [...auth0]/
│   │   │           └── route.ts          # 新規: Auth0 SDKのハンドラー
│   │   ├── users/
│   │   ├── posts/
│   │   └── user-posts/
│   ├── components/                       # 新規（オプション）
│   │   ├── LoginButton.tsx               # 新規: ログインボタンコンポーネント（オプション）
│   │   └── LogoutButton.tsx              # 新規: ログアウトボタンコンポーネント（オプション）
│   ├── lib/
│   │   └── api.ts                        # 変更なし（次のissueで修正）
│   └── types/
├── package.json                          # 修正: @auth0/nextjs-auth0の追加
├── .env.local                            # 新規: 環境変数設定ファイル（開発用）
├── .env.development                      # 新規: 環境変数設定ファイル（開発環境用）
├── .env.staging                          # 新規: 環境変数設定ファイル（ステージング環境用）
└── .env.production                       # 新規: 環境変数設定ファイル（本番環境用）
```

### 2.2 ファイル構成

#### 2.2.1 Auth0 SDKのハンドラー

**`client/src/app/api/auth/[...auth0]/route.ts`**: Auth0 SDKのAPI Routesハンドラー
- Auth0 SDKの`handleAuth()`関数をエクスポート
- ログイン、ログアウト、コールバック処理を自動的に処理

#### 2.2.2 環境変数設定ファイル（クライアント側：Next.js用）

**`client/.env.local`**: クライアント側のローカル開発用環境変数（Gitにコミットしない）
- 全環境（開発/ステージング/本番）で読み込まれる
- ローカル開発時の設定を記載
- `AUTH0_SECRET`: セッション暗号化用の秘密鍵
- `AUTH0_BASE_URL`: アプリケーションのベースURL
- `AUTH0_ISSUER_BASE_URL`: Auth0のドメイン
- `AUTH0_CLIENT_ID`: Client ID
- `AUTH0_CLIENT_SECRET`: Client Secret

**`client/.env.development`**: クライアント側の開発環境用環境変数
- 開発環境（`npm run dev`）でのみ読み込まれる
- 開発環境用のAuth0設定

**`client/.env.staging`**: クライアント側のステージング環境用環境変数（オプション）
- ステージング環境用のAuth0設定

**`client/.env.production`**: クライアント側の本番環境用環境変数（オプション）
- 本番環境（`npm run build && npm start`）でのみ読み込まれる
- 本番環境用のAuth0設定

**注意**: 
- これらの環境変数ファイルは**クライアント側（Next.js）用**である
- **今回の実装では不要**: サーバー側（Go）でのAuth0設定は今回の実装範囲外です（`config/{env}/config.yaml`への追加は行いません）
- Client Secretは全環境で環境変数のみで管理し、Gitにコミットしない
- Next.jsの環境変数の優先順位: `.env.local` > `.env.development` / `.env.production`

#### 2.2.3 UIコンポーネント（オプション）

**`client/src/components/LoginButton.tsx`**: ログインボタンコンポーネント
- 未ログイン時のみ表示
- Auth0のログインページにリダイレクト

**`client/src/components/LogoutButton.tsx`**: ログアウトボタンコンポーネント
- ログイン済み時のみ表示
- ログアウト処理を実行

### 2.3 システム構成図

```
┌─────────────────────────────────────────────────────────┐
│                    ユーザー                              │
└──────────────────┬────────────────────────────────────┘
                    │
                    │ ログインボタンクリック
                    │
                    ▼
┌─────────────────────────────────────────────────────────┐
│         client/src/app/page.tsx                          │
│  - useUser()で認証状態を確認                             │
│  - 未ログイン時: ログインボタンを表示                    │
│  - ログイン済み時: ログアウトボタンを表示                │
└──────────────────┬────────────────────────────────────┘
                    │
                    │ /api/auth/login にリダイレクト
                    │
                    ▼
┌─────────────────────────────────────────────────────────┐
│    client/src/app/api/auth/[...auth0]/route.ts          │
│    - handleAuth()でAuth0 SDKのハンドラーを提供           │
│    - ログイン、ログアウト、コールバック処理を自動処理     │
└──────────────────┬────────────────────────────────────┘
                    │
                    │ Auth0へのリダイレクト
                    │
                    ▼
┌─────────────────────────────────────────────────────────┐
│                    Auth0                                 │
│  - PARTNERアプリのアカウントでログイン                   │
│  - 認証成功後、コールバックURLにリダイレクト             │
└──────────────────┬────────────────────────────────────┘
                    │
                    │ /api/auth/callback にリダイレクト
                    │
                    ▼
┌─────────────────────────────────────────────────────────┐
│    client/src/app/api/auth/[...auth0]/route.ts          │
│    - コールバック処理                                     │
│    - JWTを取得し、HTTP-only Cookieに保存                 │
└──────────────────┬────────────────────────────────────┘
                    │
                    │ ログイン成功
                    │
                    ▼
┌─────────────────────────────────────────────────────────┐
│         client/src/app/page.tsx                          │
│  - useUser()で認証状態を確認                             │
│  - ログイン済み状態を表示                                 │
│  - ログアウトボタンを表示                                 │
└─────────────────────────────────────────────────────────┘
```

### 2.4 データフロー

#### 2.4.1 ログインフロー
```
ユーザーがログインボタンをクリック
    ↓
/api/auth/login にリダイレクト
    ↓
Auth0 SDKがAuth0のログインページにリダイレクト
    ↓
ユーザーがPARTNERアプリのアカウントでログイン
    ↓
Auth0が認証を処理
    ↓
/api/auth/callback にリダイレクト
    ↓
Auth0 SDKがコールバック処理を実行
    ↓
JWTを取得し、HTTP-only Cookieに保存
    ↓
ログイン成功、元のページにリダイレクト
    ↓
useUser()で認証状態を確認し、ログイン済み状態を表示
```

#### 2.4.2 ログアウトフロー
```
ユーザーがログアウトボタンをクリック
    ↓
/api/auth/logout にリダイレクト
    ↓
Auth0 SDKがログアウト処理を実行
    ↓
HTTP-only CookieからJWTを削除
    ↓
Auth0側のセッションも無効化（オプション）
    ↓
ログアウト成功、指定されたページにリダイレクト
    ↓
useUser()で認証状態を確認し、未ログイン状態を表示
```

#### 2.4.3 JWT取得フロー（将来のAPI呼び出し用）
```
Server ComponentsまたはClient Componentsで
getAccessToken()を呼び出し
    ↓
Auth0 SDKがHTTP-only Cookieからセッション情報を取得
    ↓
JWT（accessToken）を返却
    ↓
API呼び出しのHeaderに設定（次のissueで実装）
```

## 3. コンポーネント設計

### 3.1 Auth0 SDKの設定

#### 3.1.1 API Routesハンドラーの実装

**`client/src/app/api/auth/[...auth0]/route.ts`**:
```typescript
import { handleAuth } from '@auth0/nextjs-auth0'

export const GET = handleAuth()
```

- Auth0 SDKの`handleAuth()`関数を使用
- ログイン、ログアウト、コールバック処理を自動的に処理
- Next.js App RouterのRoute Handler形式で実装

#### 3.1.2 UserProviderの設定

**`client/src/app/layout.tsx`**:
```typescript
import { UserProvider } from '@auth0/nextjs-auth0/client'

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="ja">
      <body>
        <UserProvider>
          {children}
        </UserProvider>
      </body>
    </html>
  )
}
```

- `UserProvider`でアプリケーション全体をラップ
- 認証状態をアプリケーション全体で共有

### 3.2 ログイン機能の実装

#### 3.2.1 ログインボタンの実装

**方法1: リンクを使用（シンプル）**
```typescript
<a href="/api/auth/login">Login</a>
```

**方法2: useUserフックを使用（推奨）**
```typescript
'use client'

import { useUser } from '@auth0/nextjs-auth0/client'

export default function LoginButton() {
  const { loginWithRedirect } = useUser()

  return (
    <button onClick={() => loginWithRedirect()}>
      Login
    </button>
  )
}
```

#### 3.2.2 トップページへの統合

**`client/src/app/page.tsx`**:
```typescript
'use client'

import { useUser } from '@auth0/nextjs-auth0/client'
import Link from 'next/link'

export default function Home() {
  const { user, error, isLoading } = useUser()

  if (isLoading) return <div>Loading...</div>
  if (error) return <div>{error.message}</div>

  return (
    <main className="min-h-screen p-8">
      <div className="max-w-4xl mx-auto">
        <h1 className="text-4xl font-bold mb-8">Go DB Project Sample</h1>
        
        {/* ログイン/ログアウトUI */}
        <div className="mb-8">
          {user ? (
            <div>
              <p>Logged in as: {user.name}</p>
              <a href="/api/auth/logout">Logout</a>
            </div>
          ) : (
            <a href="/api/auth/login">Login</a>
          )}
        </div>

        {/* 既存のコンテンツ */}
        {/* ... */}
      </div>
    </main>
  )
}
```

### 3.3 ログイン状態の表示

#### 3.3.1 useUserフックの使用

```typescript
import { useUser } from '@auth0/nextjs-auth0/client'

const { user, error, isLoading } = useUser()
```

- `user`: ログイン済みの場合、ユーザー情報を含むオブジェクト
- `error`: 認証エラーが発生した場合、エラーオブジェクト
- `isLoading`: 認証状態の確認中は`true`

#### 3.3.2 UIの条件分岐

```typescript
if (isLoading) {
  // ローディング状態の表示
  return <div>Loading...</div>
}

if (error) {
  // エラー状態の表示
  return <div>Error: {error.message}</div>
}

if (user) {
  // ログイン済み状態の表示
  return (
    <div>
      <p>Logged in as: {user.name}</p>
      <a href="/api/auth/logout">Logout</a>
    </div>
  )
}

// 未ログイン状態の表示
return <a href="/api/auth/login">Login</a>
```

### 3.4 ログアウト機能の実装

#### 3.4.1 ログアウトボタンの実装

**方法1: リンクを使用（シンプル）**
```typescript
<a href="/api/auth/logout">Logout</a>
```

**方法2: useUserフックを使用（推奨）**
```typescript
'use client'

import { useUser } from '@auth0/nextjs-auth0/client'

export default function LogoutButton() {
  const { logout } = useUser()

  return (
    <button onClick={() => logout({ returnTo: '/' })}>
      Logout
    </button>
  )
}
```

#### 3.4.2 ログアウト後のリダイレクト

- `logout({ returnTo: '/' })`: ログアウト後、指定されたURLにリダイレクト
- デフォルトでは、ログアウト後にAuth0のログアウトページにリダイレクトされる

### 3.5 JWT取得機能の実装

#### 3.5.1 Server Componentsでの取得

```typescript
import { getAccessToken } from '@auth0/nextjs-auth0'

export default async function ServerComponent() {
  const { accessToken } = await getAccessToken()
  
  // accessTokenを使用（次のissueでAPI呼び出しに使用）
  return <div>Server Component</div>
}
```

#### 3.5.2 Client Componentsでの取得

```typescript
'use client'

import { useUser } from '@auth0/nextjs-auth0/client'

export default function ClientComponent() {
  const { user, getAccessToken } = useUser()

  const handleApiCall = async () => {
    const token = await getAccessToken()
    // tokenを使用（次のissueでAPI呼び出しに使用）
  }

  return <div>Client Component</div>
}
```

**注意**: 本実装では、JWT取得機能を実装するが、API呼び出しでの使用は次のissueで対応する。

### 3.6 環境変数の管理

#### 3.6.1 環境変数の設定（クライアント側：Next.js用）

**`client/.env.local`** (ローカル開発用):
```env
AUTH0_SECRET=your-secret-key-here
AUTH0_BASE_URL=http://localhost:3000
AUTH0_ISSUER_BASE_URL=https://your-domain.auth0.com
AUTH0_CLIENT_ID=your-client-id
AUTH0_CLIENT_SECRET=your-client-secret
```

**`client/.env.development`** (開発環境用):
```env
AUTH0_SECRET=${AUTH0_SECRET}
AUTH0_BASE_URL=http://localhost:3000
AUTH0_ISSUER_BASE_URL=https://dev-domain.auth0.com
AUTH0_CLIENT_ID=dev-client-id
AUTH0_CLIENT_SECRET=${AUTH0_CLIENT_SECRET}
```

**`client/.env.staging`** (ステージング環境用):
```env
AUTH0_SECRET=${AUTH0_SECRET}
AUTH0_BASE_URL=https://staging.example.com
AUTH0_ISSUER_BASE_URL=https://staging-domain.auth0.com
AUTH0_CLIENT_ID=staging-client-id
AUTH0_CLIENT_SECRET=${AUTH0_CLIENT_SECRET}
```

**`client/.env.production`** (本番環境用):
```env
AUTH0_SECRET=${AUTH0_SECRET}
AUTH0_BASE_URL=https://example.com
AUTH0_ISSUER_BASE_URL=https://prod-domain.auth0.com
AUTH0_CLIENT_ID=prod-client-id
AUTH0_CLIENT_SECRET=${AUTH0_CLIENT_SECRET}
```

#### 3.6.2 環境変数の読み込み

- Next.jsが自動的に環境変数を読み込む（`client/`ディレクトリ内のファイル）
- `.env.local`は全環境で読み込まれ、Gitにコミットしない（ローカル開発用）
- `.env.development`は開発環境（`npm run dev`）でのみ読み込まれる
- `.env.production`は本番環境（`npm run build && npm start`）でのみ読み込まれる
- 優先順位: `.env.local` > `.env.development` / `.env.production`

#### 3.6.3 Client Secretの管理

- Client Secretは全環境で環境変数のみで管理
- `client/.env.local`、`client/.env.development`、`client/.env.staging`、`client/.env.production`に直接記載しない
- 環境変数として設定（例: `export AUTH0_CLIENT_SECRET=your-secret`）
- CI/CD環境では、シークレット管理システム（GitHub Secrets等）を使用
- **今回の実装では不要**: サーバー側（Go）でのAuth0設定は今回の実装範囲外です（`config/{env}/config.yaml`への追加は行いません）

## 4. 実装詳細

### 4.1 パッケージのインストール

```bash
cd client
npm install @auth0/nextjs-auth0
```

### 4.2 API Routesハンドラーの作成

**`client/src/app/api/auth/[...auth0]/route.ts`**を作成:
```typescript
import { handleAuth } from '@auth0/nextjs-auth0'

export const GET = handleAuth()
```

### 4.3 UserProviderの設定

**`client/src/app/layout.tsx`**を修正:
```typescript
import { UserProvider } from '@auth0/nextjs-auth0/client'

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="ja">
      <body>
        <UserProvider>
          {children}
        </UserProvider>
      </body>
    </html>
  )
}
```

### 4.4 トップページへのログイン/ログアウトUIの追加

**`client/src/app/page.tsx`**を修正:
- `'use client'`ディレクティブを追加
- `useUser()`フックをインポート
- ログイン状態に応じてUIを切り替え

### 4.5 環境変数の設定（クライアント側：Next.js用）

- `client/.env.local`を作成し、ローカル開発用の環境変数を設定（全環境で読み込まれる）
- `client/.env.development`を作成し、開発環境用の環境変数を設定（`npm run dev`でのみ読み込まれる）
- `client/.env.staging`、`client/.env.production`を作成（必要に応じて）
- Client Secretは環境変数として設定（ファイルに直接記載しない）
- **今回の実装では不要**: サーバー側（Go）でのAuth0設定は今回の実装範囲外です（`config/{env}/config.yaml`への追加は行いません）

## 5. エラーハンドリング

### 5.1 認証エラーの処理

```typescript
const { user, error, isLoading } = useUser()

if (error) {
  return (
    <div className="error">
      <p>認証エラーが発生しました: {error.message}</p>
      <a href="/api/auth/login">再度ログイン</a>
    </div>
  )
}
```

### 5.2 ネットワークエラーの処理

- Auth0への接続に失敗した場合、適切なエラーメッセージを表示
- リトライ機能を実装（オプション）

### 5.3 ログイン失敗時の処理

- Auth0からエラーが返された場合、エラーメッセージを表示
- ユーザーに再度ログインを促す

### 5.4 JWT取得エラーの処理

```typescript
try {
  const { accessToken } = await getAccessToken()
  // JWTを使用
} catch (error) {
  // エラーハンドリング
  console.error('Failed to get access token:', error)
}
```

## 6. セキュリティ考慮事項

### 6.1 JWTの安全な保存

- HTTP-only Cookieに保存することで、JavaScriptからアクセスできないようにする
- XSS攻撃からJWTを保護

### 6.2 Client Secretの管理

- 環境変数のみで管理し、Gitにコミットしない
- `.gitignore`に`.env.local`が含まれていることを確認
- CI/CD環境では、シークレット管理システムを使用

### 6.3 HTTPSの使用

- 本番環境ではHTTPSを使用
- 開発環境ではHTTPも可（localhost）

### 6.4 セッション管理

- Auth0 SDKが自動的にセッションを管理
- 適切な有効期限を設定

## 7. テスト戦略

### 7.1 ユニットテスト

#### 7.1.1 ログインボタンのテスト

```typescript
import { render, screen } from '@testing-library/react'
import { useUser } from '@auth0/nextjs-auth0/client'
import LoginButton from '@/components/LoginButton'

jest.mock('@auth0/nextjs-auth0/client')

describe('LoginButton', () => {
  it('renders login button when user is not logged in', () => {
    (useUser as jest.Mock).mockReturnValue({
      user: null,
      error: null,
      isLoading: false,
    })

    render(<LoginButton />)
    expect(screen.getByText('Login')).toBeInTheDocument()
  })
})
```

#### 7.1.2 ログアウトボタンのテスト

```typescript
import { render, screen } from '@testing-library/react'
import { useUser } from '@auth0/nextjs-auth0/client'
import LogoutButton from '@/components/LogoutButton'

jest.mock('@auth0/nextjs-auth0/client')

describe('LogoutButton', () => {
  it('renders logout button when user is logged in', () => {
    (useUser as jest.Mock).mockReturnValue({
      user: { name: 'Test User' },
      error: null,
      isLoading: false,
    })

    render(<LogoutButton />)
    expect(screen.getByText('Logout')).toBeInTheDocument()
  })
})
```

### 7.2 統合テスト

#### 7.2.1 ログインフローのテスト

- ログインボタンをクリック
- Auth0のログインページにリダイレクトされることを確認
- ログイン成功後、コールバック処理が正常に動作することを確認

#### 7.2.2 ログアウトフローのテスト

- ログアウトボタンをクリック
- HTTP-only CookieからJWTが削除されることを確認
- 未ログイン状態に戻ることを確認

### 7.3 E2Eテスト

#### 7.3.1 ログインE2Eテスト

```typescript
import { test, expect } from '@playwright/test'

test('user can login', async ({ page }) => {
  await page.goto('http://localhost:3000')
  
  // ログインボタンをクリック
  await page.click('text=Login')
  
  // Auth0のログインページにリダイレクトされることを確認
  await expect(page).toHaveURL(/auth0\.com/)
  
  // ログイン処理（テスト用の認証情報を使用）
  // ...
  
  // ログイン成功後、元のページに戻ることを確認
  await expect(page).toHaveURL('http://localhost:3000')
  
  // ログアウトボタンが表示されることを確認
  await expect(page.locator('text=Logout')).toBeVisible()
})
```

#### 7.3.2 ログアウトE2Eテスト

```typescript
import { test, expect } from '@playwright/test'

test('user can logout', async ({ page }) => {
  // ログイン状態で開始
  // ...
  
  // ログアウトボタンをクリック
  await page.click('text=Logout')
  
  // ログアウト後、未ログイン状態になることを確認
  await expect(page.locator('text=Login')).toBeVisible()
  await expect(page.locator('text=Logout')).not.toBeVisible()
})
```

## 8. パフォーマンス考慮事項

### 8.1 認証状態の管理

- `useUser()`フックは効率的に認証状態を管理
- 不要な再認証を避ける

### 8.2 ローディング状態の表示

- 認証状態の確認中は適切なローディング表示を行う
- ユーザー体験を向上させる

## 9. 参考情報

### 9.1 関連ドキュメント

- 要件定義書: `.kiro/specs/0018-auth0-login/requirements.md`
- Auth0設定手順: `docs/Partner-Idp-Auth0-Login.md`

### 9.2 技術資料

- [Auth0 Next.js SDK Documentation](https://auth0.com/docs/quickstart/webapp/nextjs)
- [Auth0 Next.js SDK - Getting an Access Token](https://auth0.com/docs/quickstart/webapp/nextjs/01-login#get-an-access-token)
- [Next.js App Router Documentation](https://nextjs.org/docs/app)

### 9.3 既存実装の参考

- 既存の設計書: `.kiro/specs/0016-fix-tablesplit/design.md`（フォーマット参考）
