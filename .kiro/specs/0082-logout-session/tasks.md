# Auth0セッション完全ログアウト実装タスク一覧

## 概要
ログアウト時にAuth0のセッションも完全にクリアするための実装タスク一覧。要件定義書と設計書に基づき、実装可能な粒度でタスクを分解。

## 実装フェーズ

### Phase 1: 型チェックエラーの修正

#### - [ ] タスク 1.1: jest.setup.jsの修正
**目的**: `auth`関数のモック定義を修正して、`Promise<Session | null>`を返すようにする。

**作業内容**:
- `client/jest.setup.js`の53-64行目を修正
- `jest.fn(() => null)`を`jest.fn(() => Promise.resolve(null))`に変更
- `auth`関数が`Promise<Session | null>`を返すようにモックを定義

**実装コード**:
```javascript
// 変更前
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

// 変更後
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

**確認事項**:
- `jest.fn(() => Promise.resolve(null))`に変更されている
- `auth`関数が`Promise`を返すようにモックが定義されている
- TypeScriptの型チェックでエラーが発生しない

#### - [ ] タスク 1.2: テストファイルの修正
**目的**: `auth`関数の型を適切にアサーションして、型チェックエラーを解消する。

**作業内容**:
- `client/src/__tests__/integration/page-page.test.tsx`を修正
- `Session`型をインポート
- `auth`関数の型を適切にアサーション

**実装コード**:
```typescript
// 変更前
import { render, screen } from '@testing-library/react'
import Home from '@/app/page'
import { auth } from '@/auth'

// Mock auth function
const mockAuth = auth as jest.MockedFunction<typeof auth>

// 変更後
import { render, screen } from '@testing-library/react'
import Home from '@/app/page'
import { auth } from '@/auth'
import type { Session } from 'next-auth'

// Mock auth function
const mockAuth = auth as unknown as jest.MockedFunction<() => Promise<Session | null>>
```

**確認事項**:
- `Session`型がインポートされている
- `auth`関数の型が適切にアサーションされている
- TypeScriptの型チェックでエラーが発生しない

#### - [ ] タスク 1.3: 型チェックの確認
**目的**: 型チェックエラーが解消されていることを確認する。

**作業内容**:
- `cd client && npm run type-check`を実行
- 型チェックエラーが0件であることを確認

**確認事項**:
- 型チェックコマンドが正常に実行される
- 型チェックエラーが0件である
- エラーメッセージが表示されない

### Phase 2: signOutActionの修正

#### - [ ] タスク 2.1: 環境変数の取得と検証の実装
**目的**: 環境変数を取得し、設定されていない場合はエラーをthrowする。

**作業内容**:
- `client/lib/actions/auth-actions.ts`の`signOutAction`を修正
- `process.env.AUTH0_ISSUER`を取得
- `process.env.NEXT_PUBLIC_APP_BASE_URL`を取得
- 環境変数が設定されていない場合、エラーをthrow

**実装コード**:
```typescript
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
  
  // 以下、タスク2.2で実装
}
```

**確認事項**:
- `AUTH0_ISSUER`が正しく取得されている
- `NEXT_PUBLIC_APP_BASE_URL`が正しく取得されている
- 環境変数が設定されていない場合、適切なエラーメッセージでエラーがthrowされる
- エラーメッセージに環境変数名が含まれている

#### - [ ] タスク 2.2: Auth0ログアウトURLの構築
**目的**: `AUTH0_ISSUER`からAuth0ログアウトURLを構築し、`returnTo`パラメータを追加する。

**作業内容**:
- `AUTH0_ISSUER`から`/v2/logout`を追加してAuth0ログアウトURLを構築
- `returnTo`パラメータに`NEXT_PUBLIC_APP_BASE_URL`を使用
- URLエンコーディングを適用（`encodeURIComponent`を使用）

**実装コード**:
```typescript
// AUTH0_ISSUERから/v2/logoutを追加してAuth0ログアウトURLを構築
const auth0LogoutUrl = `${auth0Issuer}/v2/logout`

// Auth0ログアウトURLにreturnToパラメータを追加
const returnToUrl = `${appBaseUrl}/`
const logoutUrl = `${auth0LogoutUrl}?returnTo=${encodeURIComponent(returnToUrl)}`
```

**確認事項**:
- `auth0LogoutUrl`が正しく構築されている（`{AUTH0_ISSUER}/v2/logout`）
- `returnToUrl`が正しく構築されている（`{NEXT_PUBLIC_APP_BASE_URL}/`）
- `logoutUrl`が正しく構築されている（`{auth0LogoutUrl}?returnTo={encoded_returnToUrl}`）
- URLエンコーディングが正しく適用されている

#### - [ ] タスク 2.3: signOut関数の呼び出し
**目的**: next-authの`signOut`関数にパラメータを渡して、Auth0ログアウトURLにリダイレクトする。

**作業内容**:
- `signOut`関数に`redirect: true`と`redirectTo: logoutUrl`を渡す
- `signOut`関数を呼び出す

**実装コード**:
```typescript
// next-authのsignOutにパラメータを渡してリダイレクト
await signOut({
  redirect: true,
  redirectTo: logoutUrl
})
```

**確認事項**:
- `signOut`関数が正しいパラメータで呼び出されている
- `redirect: true`が設定されている
- `redirectTo: logoutUrl`が設定されている
- 完全な実装コードが正しく動作する

#### - [ ] タスク 2.4: 完全な実装コードの確認
**目的**: `signOutAction`の完全な実装コードが正しく動作することを確認する。

**作業内容**:
- 完全な実装コードを確認
- コードの一貫性を確認
- エラーハンドリングが適切に実装されていることを確認

**完全な実装コード**:
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

**確認事項**:
- 完全な実装コードが正しく動作する
- コードの一貫性が保たれている
- 環境変数の検証が適切に実装されている
- URL構築が適切に実装されている
- エラーハンドリングが適切に実装されている
- TypeScriptの型チェックでエラーが発生しない

### Phase 3: 環境変数の設定

#### - [ ] タスク 3.1: 環境変数ファイルの確認
**目的**: `client/.env.local`ファイルが存在することを確認する。

**作業内容**:
- `client/.env.local`ファイルが存在するか確認
- 存在しない場合は新規作成する準備をする

**確認事項**:
- `client/.env.local`ファイルが存在する（または新規作成可能である）

#### - [ ] タスク 3.2: NEXT_PUBLIC_APP_BASE_URLの追加
**目的**: `client/.env.local`に`NEXT_PUBLIC_APP_BASE_URL`環境変数を追加する。

**作業内容**:
- `client/.env.local`に`NEXT_PUBLIC_APP_BASE_URL`を追加
- 開発環境の場合は`http://localhost:3000`を設定
- 本番環境の場合は適切なURLを設定

**実装例**:
```env
# client/.env.local

# 既存の環境変数
AUTH0_ISSUER=https://your-domain.auth0.com

# 新規追加の環境変数
NEXT_PUBLIC_APP_BASE_URL=http://localhost:3000
```

**確認事項**:
- `NEXT_PUBLIC_APP_BASE_URL`が正しく追加されている
- 環境変数の形式が正しい（`NEXT_PUBLIC_`プレフィックスが付いている）
- 値が適切に設定されている（開発環境: `http://localhost:3000`）

### Phase 4: 動作確認

#### - [ ] タスク 4.1: 開発サーバーの起動
**目的**: 開発サーバーを起動して、実装が正しく動作することを確認する準備をする。

**作業内容**:
- クライアントの開発サーバーを起動
- サーバーが正常に起動することを確認

**確認事項**:
- 開発サーバーが正常に起動する
- エラーが発生しない

#### - [ ] タスク 4.2: ログアウト機能の動作確認
**目的**: ログアウトボタンをクリックして、Auth0ログアウトURLにリダイレクトされることを確認する。

**作業内容**:
- ログイン状態でアプリケーションにアクセス
- ログアウトボタンをクリック
- Auth0ログアウトURLにリダイレクトされることを確認
- ブラウザの開発者ツールでリダイレクト先URLを確認

**確認事項**:
- ログアウトボタンが表示されている（ログイン状態時）
- ログアウトボタンをクリックすると、Auth0ログアウトURLにリダイレクトされる
- リダイレクト先URLが正しい形式である（`{AUTH0_ISSUER}/v2/logout?returnTo={encoded_url}`）

#### - [ ] タスク 4.3: Auth0ログアウト後のリダイレクト確認
**目的**: Auth0ログアウト後、アプリケーションのトップページに戻ることを確認する。

**作業内容**:
- Auth0ログアウトURLにリダイレクトされた後、Auth0のログアウト処理が完了するまで待つ
- `returnTo`パラメータで指定されたURL（アプリケーションのトップページ）にリダイレクトされることを確認

**確認事項**:
- Auth0ログアウト後、アプリケーションのトップページ（`/`）にリダイレクトされる
- リダイレクト先URLが正しい（`NEXT_PUBLIC_APP_BASE_URL/`）

#### - [ ] タスク 4.4: Auth0セッションのクリア確認
**目的**: ログアウト後、Auth0のセッションがクリアされていることを確認する。

**作業内容**:
- ログアウト後、再度ログインする
- 前回のログイン情報が残っていないことを確認（Auth0のログイン画面で前回のアカウント情報が表示されない）

**確認事項**:
- ログアウト後、Auth0のセッションがクリアされている
- 次回ログイン時に前回のログイン情報が残らない
- Auth0のログイン画面で前回のアカウント情報が表示されない

#### - [ ] タスク 4.5: 既存機能への影響確認
**目的**: 既存のログアウト機能に影響がないことを確認する。

**作業内容**:
- 既存のログアウトボタン（`AuthButtons`コンポーネント）が正常に動作することを確認
- ヘッダーのログアウトボタンが正常に動作することを確認（0081-header-logoutで実装済み）

**確認事項**:
- 既存のログアウト機能が正常に動作する
- ヘッダーのログアウトボタンが正常に動作する
- 画面内のログアウトボタンが正常に動作する

### Phase 5: テスト

#### - [ ] タスク 5.1: 型チェックの実行
**目的**: 型チェックエラーが発生していないことを確認する。

**作業内容**:
- `cd client && npm run type-check`を実行
- 型チェックエラーが0件であることを確認

**確認事項**:
- 型チェックコマンドが正常に実行される
- 型チェックエラーが0件である

#### - [ ] タスク 5.2: 既存テストの実行
**目的**: 既存のテストが正常に動作することを確認する。

**作業内容**:
- 既存のテストを実行
- すべてのテストが正常に動作することを確認

**確認事項**:
- 既存のテストが正常に実行される
- すべてのテストがパスする
- テストエラーが発生しない

#### - [ ] タスク 5.3: ユニットテストの実装
**目的**: `signOutAction`関数のユニットテストを実装する

**作業内容**:
- `signOutAction`関数のユニットテストを実装
- 環境変数が設定されている場合のテスト
- 環境変数が設定されていない場合のテスト
- URL構築のテスト

**確認事項**:
- ユニットテストが実装されている
- すべてのテストケースがカバーされている
- テストが正常に実行される

## 受け入れ基準の確認

### 要件定義書の受け入れ基準

#### 6.1 signOutActionの修正
- [ ] `signOutAction`で`signOut()`にパラメータを渡している
- [ ] `redirect: true`と`redirectTo`を指定している
- [ ] `AUTH0_ISSUER`から`/v2/logout`を追加してAuth0ログアウトURLを構築している
- [ ] `redirectTo`には構築したAuth0のログアウトURLを設定している
- [ ] Auth0ログアウトURLに`returnTo`パラメータを追加している
- [ ] `returnTo`パラメータにはアプリケーションのベースURLを設定している
- [ ] URLエンコーディングを適用している

#### 6.2 環境変数の設定
- [ ] `AUTH0_ISSUER`が設定されている（既存の環境変数）
- [ ] `client/.env.local`に`NEXT_PUBLIC_APP_BASE_URL`が定義されている
- [ ] 環境変数の形式が正しい（`AUTH0_ISSUER`は`/v2/logout`を含めない）

#### 6.3 エラーハンドリング
- [ ] `AUTH0_ISSUER`が設定されていない場合、エラーを発生させる
- [ ] `NEXT_PUBLIC_APP_BASE_URL`が設定されていない場合、エラーを発生させる
- [ ] エラーメッセージに環境変数名が含まれている
- [ ] フォールバック機能が実装されていない

#### 6.4 動作確認
- [ ] ログアウトボタンをクリックすると、Auth0のログアウトURLにリダイレクトされる
- [ ] Auth0ログアウト後、アプリケーションのトップページに戻る
- [ ] ログアウト後、Auth0のセッションがクリアされている
- [ ] 次回ログイン時に前回のログイン情報が残らない
- [ ] 既存のログアウト機能に影響がない

#### 6.5 コード品質
- [ ] コードの一貫性が保たれている
- [ ] 環境変数の検証が適切に実装されている
- [ ] URL構築が適切に実装されている
- [ ] エラーハンドリングが適切に実装されている

#### 9.1 型チェックエラーの修正計画
- [ ] `jest.setup.js`の修正が完了している
- [ ] テストファイルの修正が完了している
- [ ] 型チェックが通ることを確認している

## 実装時の注意事項

### 環境変数の設定
- `NEXT_PUBLIC_APP_BASE_URL`は`NEXT_PUBLIC_`プレフィックスが必要（クライアント側で使用するため）
- `.env.local`はgitから除外されるため、開発環境と本番環境で個別に設定する必要がある

### エラーハンドリング
- 環境変数が設定されていない場合は、フォールバック機能を実装せず、エラーをthrowする
- エラーメッセージは明確で、設定が必要な環境変数を示す

### URL構築
- `AUTH0_ISSUER`から`/v2/logout`を追加する際、スラッシュの重複に注意する
- `returnTo`パラメータには`/`を含める（トップページにリダイレクトするため）
- URLエンコーディングを必ず適用する（`encodeURIComponent`を使用）

### テスト
- 環境変数のモックが必要な場合は、適切にモックを設定する
- 既存のテストが正常に動作することを確認する

## 参考情報

### 関連ドキュメント
- 要件定義書: `.kiro/specs/0082-logout-session/requirements.md`
- 設計書: `.kiro/specs/0082-logout-session/design.md`
- GitHub Issue: https://github.com/taku-o/go-webdb-template/issues/168

### 技術リファレンス
- next-auth signOut: https://next-auth.js.org/getting-started/client#signout
- Auth0 Logout: https://auth0.com/docs/authenticate/login/logout
- Next.js Server Actions: https://nextjs.org/docs/app/building-your-application/data-fetching/server-actions-and-mutations
