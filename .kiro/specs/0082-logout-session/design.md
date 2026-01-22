# Auth0セッション完全ログアウト設計書

## 1. 概要

### 1.1 設計の目的
要件定義書に基づき、ログアウト時にAuth0のセッションも完全にクリアするための詳細設計を定義する。next-authの`signOut`にパラメータを渡して、Auth0のログアウトURLにリダイレクトし、ログアウト後にアプリケーションに戻る仕組みを実装する。

### 1.2 設計の範囲
- `client/lib/actions/auth-actions.ts`の`signOutAction`の修正
- 環境変数の設定（`NEXT_PUBLIC_APP_BASE_URL`）
- 型チェックエラーの修正（`jest.setup.js`、テストファイル）

### 1.3 設計方針
- **完全なログアウト**: Auth0のセッションも完全にクリアする
- **環境変数の検証**: 環境変数が設定されていない場合はエラーを発生させる（フォールバックなし）
- **URL構築**: `AUTH0_ISSUER`から`/v2/logout`を追加してAuth0ログアウトURLを構築
- **リダイレクト**: Auth0ログアウト後、アプリケーションのトップページに戻る
- **既存機能の維持**: 既存のログアウト機能の呼び出し方法は変更しない
- **型安全性**: TypeScriptの型チェックエラーを修正する

## 2. 関数設計

### 2.1 signOutActionの修正設計

#### 2.1.1 変更前の実装

```typescript
"use server"

import { signIn, signOut } from '@/auth'

export async function signInAction() {
  await signIn('auth0')
}

export async function signOutAction() {
  await signOut()
}
```

#### 2.1.2 変更後の実装

```typescript
"use server"

import { signIn, signOut } from '@/auth'

export async function signInAction() {
  await signIn('auth0')
}

export async function signOutAction() {
  // 環境変数の取得
  const auth0Issuer = process.env.AUTH0_ISSUER
  const appBaseUrl = process.env.NEXT_PUBLIC_APP_BASE_URL
  
  // 環境変数の検証
  if (!auth0Issuer) {
    throw new Error('AUTH0_ISSUER is not set')
  }
  if (!appBaseUrl) {
    throw new Error('NEXT_PUBLIC_APP_BASE_URL is not set')
  }
  
  // AUTH0_ISSUERから/v2/logoutを追加してAuth0ログアウトURLを構築
  const auth0LogoutUrl = `${auth0Issuer}/v2/logout`
  
  // Auth0ログアウトURLにreturnToパラメータを追加
  const returnToUrl = `${appBaseUrl}/`
  const logoutUrl = `${auth0LogoutUrl}?returnTo=${encodeURIComponent(returnToUrl)}`
  
  // next-authのsignOutにパラメータを渡してリダイレクト
  await signOut({
    redirect: true,
    redirectTo: logoutUrl
  })
}
```

#### 2.1.3 実装の詳細

**環境変数の取得**:
- `process.env.AUTH0_ISSUER`: Auth0のIssuer URL（既存の環境変数）
- `process.env.NEXT_PUBLIC_APP_BASE_URL`: アプリケーションのベースURL（新規追加）

**環境変数の検証**:
- `AUTH0_ISSUER`が設定されていない場合: `Error('AUTH0_ISSUER is not set')`をthrow
- `NEXT_PUBLIC_APP_BASE_URL`が設定されていない場合: `Error('NEXT_PUBLIC_APP_BASE_URL is not set')`をthrow
- フォールバック機能は実装しない（エラーを発生させる）

**URL構築**:
1. `AUTH0_ISSUER`から`/v2/logout`を追加してAuth0ログアウトURLを構築
   - 例: `https://your-domain.auth0.com` → `https://your-domain.auth0.com/v2/logout`
2. `returnTo`パラメータに`NEXT_PUBLIC_APP_BASE_URL`を使用
   - 例: `http://localhost:3000` → `http://localhost:3000/`
3. URLエンコーディングを適用（`encodeURIComponent`を使用）
   - 例: `http://localhost:3000/` → `http%3A%2F%2Flocalhost%3A3000%2F`
4. 最終的なURL: `{AUTH0_ISSUER}/v2/logout?returnTo={encoded_returnToUrl}`

**signOut関数の呼び出し**:
- `redirect: true`: リダイレクトを有効にする
- `redirectTo: logoutUrl`: Auth0ログアウトURLにリダイレクト

### 2.2 エラーハンドリング設計

#### 2.2.1 エラーの種類

1. **環境変数未設定エラー**:
   - `AUTH0_ISSUER`が設定されていない場合
   - `NEXT_PUBLIC_APP_BASE_URL`が設定されていない場合

2. **エラーメッセージ**:
   - 明確で、設定が必要な環境変数を示す
   - 環境変数名を含める
   - 例: `'AUTH0_ISSUER is not set'`

#### 2.2.2 エラーの伝播

- Server Action内でエラーが発生した場合、Next.jsが自動的にエラーハンドリングを行う
- クライアント側では、エラーが発生した場合、適切なエラーメッセージが表示される（Next.jsのデフォルト動作）

## 3. 環境変数設計

### 3.1 必要な環境変数

#### 3.1.1 既存の環境変数

- `AUTH0_ISSUER`: Auth0のIssuer URL
  - **ファイル**: `client/.env.local`（既存）
  - **形式**: `https://{domain}.auth0.com`
  - **例**: `https://your-domain.auth0.com`
  - **注意**: `/v2/logout`を含めない（コード内で追加）

#### 3.1.2 新規追加の環境変数

- `NEXT_PUBLIC_APP_BASE_URL`: アプリケーションのベースURL
  - **ファイル**: `client/.env.local`（新規追加）
  - **形式**: `http://localhost:3000`（開発環境）または`https://your-domain.com`（本番環境）
  - **例**: `http://localhost:3000`
  - **注意**: `NEXT_PUBLIC_`プレフィックスが必要（クライアント側で使用するため）

### 3.2 環境変数の設定例

```env
# client/.env.local

# 既存の環境変数
AUTH0_ISSUER=https://your-domain.auth0.com

# 新規追加の環境変数
NEXT_PUBLIC_APP_BASE_URL=http://localhost:3000
```

### 3.3 環境変数の検証

- **検証タイミング**: `signOutAction`実行時
- **検証方法**: 環境変数が`undefined`または空文字列の場合、エラーをthrow
- **フォールバック**: なし（エラーを発生させる）

## 4. URL構築設計

### 4.1 Auth0ログアウトURLの構築

#### 4.1.1 URL構築の流れ

1. **ベースURLの取得**: `AUTH0_ISSUER`から取得
2. **ログアウトパスの追加**: `/v2/logout`を追加
3. **returnToパラメータの構築**: `NEXT_PUBLIC_APP_BASE_URL`に`/`を追加
4. **URLエンコーディング**: `encodeURIComponent`を使用してエンコード
5. **最終URLの構築**: クエリパラメータとして`returnTo`を追加

#### 4.1.2 URL構築の例

**入力**:
- `AUTH0_ISSUER`: `https://your-domain.auth0.com`
- `NEXT_PUBLIC_APP_BASE_URL`: `http://localhost:3000`

**処理**:
1. `auth0LogoutUrl = "https://your-domain.auth0.com/v2/logout"`
2. `returnToUrl = "http://localhost:3000/"`
3. `encodedReturnTo = encodeURIComponent("http://localhost:3000/")` → `"http%3A%2F%2Flocalhost%3A3000%2F"`
4. `logoutUrl = "https://your-domain.auth0.com/v2/logout?returnTo=http%3A%2F%2Flocalhost%3A3000%2F"`

**出力**:
```
https://your-domain.auth0.com/v2/logout?returnTo=http%3A%2F%2Flocalhost%3A3000%2F
```

### 4.2 リダイレクトフロー

1. **ユーザーがログアウトボタンをクリック**
2. **`signOutAction`が実行される**
3. **環境変数の検証**
4. **Auth0ログアウトURLの構築**
5. **`signOut({ redirect: true, redirectTo: logoutUrl })`が実行される**
6. **Auth0のログアウトURLにリダイレクト**
7. **Auth0がセッションをクリア**
8. **`returnTo`パラメータで指定されたURL（アプリケーションのトップページ）にリダイレクト**

## 5. 型チェックエラー修正設計

### 5.1 問題の概要

`cd client && npm run type-check 2>&1`を実行すると、以下の型チェックエラーが発生しています：

```
src/__tests__/integration/page-page.test.tsx(35,36): error TS2345: Argument of type 'null' is not assignable to parameter of type 'never'.
...
```

### 5.2 原因分析

1. **jest.setup.jsのモック定義の問題**: `auth`関数をモックしているが、`jest.fn(() => null)`となっており、`Promise`を返していない
2. **型推論の問題**: `auth`関数の実際の戻り値は`Promise<Session | null>`であるが、モックの型定義が適切でないため、TypeScriptが型を正しく推論できていない

### 5.3 修正設計

#### 5.3.1 jest.setup.jsの修正

**変更前**:
```javascript
jest.mock('@/auth', () => {
  const mockAuth = jest.fn(() => null)
  return {
    handlers: {
      GET: {},
      POST: {},
    },
    auth: mockAuth,
    signIn: jest.fn(),
    signOut: jest.fn(),
  }
})
```

**変更後**:
```javascript
jest.mock('@/auth', () => {
  const mockAuth = jest.fn(() => Promise.resolve(null))
  return {
    handlers: {
      GET: {},
      POST: {},
    },
    auth: mockAuth,
    signIn: jest.fn(),
    signOut: jest.fn(),
  }
})
```

**変更内容**:
- `jest.fn(() => null)`を`jest.fn(() => Promise.resolve(null))`に変更
- `auth`関数が`Promise<Session | null>`を返すようにモックを定義

#### 5.3.2 テストファイルの修正

**変更前**:
```typescript
import { render, screen } from '@testing-library/react'
import Home from '@/app/page'
import { auth } from '@/auth'

// Mock auth function
const mockAuth = auth as jest.MockedFunction<typeof auth>
```

**変更後**:
```typescript
import { render, screen } from '@testing-library/react'
import Home from '@/app/page'
import { auth } from '@/auth'
import type { Session } from 'next-auth'

// Mock auth function
const mockAuth = auth as unknown as jest.MockedFunction<() => Promise<Session | null>>
```

**変更内容**:
- `Session`型をインポート
- `auth`関数の型を適切にアサーション（`as unknown as jest.MockedFunction<() => Promise<Session | null>>`）

### 5.4 検証方法

修正後、以下のコマンドで型チェックが通ることを確認：
```bash
cd client && npm run type-check
```

## 6. テスト設計

### 6.1 テストの種類

#### 6.1.1 ユニットテスト

- **対象**: `signOutAction`関数
- **テスト内容**:
  - 環境変数が設定されている場合、正しくAuth0ログアウトURLを構築する
  - 環境変数が設定されていない場合、エラーをthrowする
  - URLエンコーディングが正しく適用される
  - `signOut`関数が正しいパラメータで呼び出される

#### 6.1.2 統合テスト

- **対象**: ログアウト処理全体
- **テスト内容**:
  - ログアウトボタンをクリックすると、Auth0ログアウトURLにリダイレクトされる
  - Auth0ログアウト後、アプリケーションのトップページに戻る

#### 6.1.3 E2Eテスト

- **対象**: ログアウトフロー全体
- **テスト内容**:
  - ログイン → ログアウト → 次回ログイン時に前回のログイン情報が残らない

### 6.2 モック設計

#### 6.2.1 環境変数のモック

- **テスト環境**: `process.env.AUTH0_ISSUER`と`process.env.NEXT_PUBLIC_APP_BASE_URL`をモック
- **モック方法**: Jestの`process.env`を直接設定

#### 6.2.2 signOut関数のモック

- **モック方法**: `@/auth`の`signOut`関数をモック
- **検証内容**: `signOut`関数が正しいパラメータで呼び出されることを確認

## 7. 実装順序

### 7.1 実装の優先順位

1. **型チェックエラーの修正**（優先度: 中）
   - `jest.setup.js`の修正
   - テストファイルの修正
   - 型チェックの確認

2. **signOutActionの修正**（優先度: 高）
   - 環境変数の取得と検証
   - URL構築ロジックの実装
   - `signOut`関数の呼び出し

3. **環境変数の設定**（優先度: 高）
   - `client/.env.local`に`NEXT_PUBLIC_APP_BASE_URL`を追加

4. **テストの実装**（優先度: 中）
   - ユニットテストの実装
   - 統合テストの実装

5. **動作確認**（優先度: 高）
   - ログアウト機能の動作確認
   - Auth0セッションのクリア確認

### 7.2 実装の依存関係

- 型チェックエラーの修正は、他の実装に影響しないため、独立して実装可能
- `signOutAction`の修正は、環境変数の設定に依存する
- テストの実装は、`signOutAction`の修正後に実装

## 8. リスクと対策

### 8.1 リスク

1. **環境変数が設定されていない場合のエラー**
   - **リスク**: ログアウト機能が動作しない
   - **対策**: 環境変数の検証を実装し、明確なエラーメッセージを表示

2. **Auth0ログアウトURLの構築ミス**
   - **リスク**: リダイレクトが失敗する
   - **対策**: URL構築ロジックを明確に実装し、テストで検証

3. **型チェックエラーの修正漏れ**
   - **リスク**: 型安全性が損なわれる
   - **対策**: 型チェックエラーの修正計画を実装し、型チェックを実行して確認

### 8.2 対策

- **環境変数の検証**: 実装時に必ず検証ロジックを追加
- **URL構築のテスト**: ユニットテストでURL構築ロジックを検証
- **型チェックの実行**: 実装後に必ず型チェックを実行

## 9. 参考情報

### 9.1 関連ドキュメント

- 要件定義書: `.kiro/specs/0082-logout-session/requirements.md`
- GitHub Issue: https://github.com/taku-o/go-webdb-template/issues/168

### 9.2 技術リファレンス

- next-auth signOut: https://next-auth.js.org/getting-started/client#signout
- Auth0 Logout: https://auth0.com/docs/authenticate/login/logout
- Next.js Server Actions: https://nextjs.org/docs/app/building-your-application/data-fetching/server-actions-and-mutations

### 9.3 既存ファイル

- `client/lib/actions/auth-actions.ts`: 修正対象
- `client/auth.ts`: NextAuthの設定（参考）
- `client/jest.setup.js`: 型チェックエラー修正対象
- `client/src/__tests__/integration/page-page.test.tsx`: 型チェックエラー修正対象
