# Auth0セッション完全ログアウト要件定義書

## 1. 概要

### 1.1 プロジェクト情報
- **プロジェクト名**: go-webdb-template
- **Feature名**: 0082-logout-session
- **作成日**: 2026-01-17
- **関連Issue**: https://github.com/taku-o/go-webdb-template/issues/168

### 1.2 目的
ログアウト時にAuth0のセッションも完全にクリアし、次回ログイン時に前回のログイン情報が残らないようにする。

### 1.3 スコープ
- `client/lib/actions/auth-actions.ts`の`signOutAction`を修正
- next-authの`signOut`にパラメータを渡して、Auth0のログアウトURLにリダイレクト
- Auth0ログアウト後、アプリケーションに戻るためのURLパラメータを設定
- 環境変数`AUTH0_ISSUER`からAuth0ログアウトURLを計算（`/v2/logout`を追加）
- 環境変数`client/.env.local`からリダイレクト先URLを取得
- 環境変数が設定されていない場合はエラーを発生させる

**本実装の範囲外**:
- ログイン機能の変更
- 認証ロジックの変更（Auth0プロバイダー設定の変更）
- その他の認証関連コンポーネントの変更
- UIデザインの変更

## 2. 背景・現状分析

### 2.1 現在の状況

#### 2.1.1 ログアウト処理（auth-actions.ts）
- **ファイル**: `client/lib/actions/auth-actions.ts`
- **実装**: `signOutAction`が定義されている（9-11行目）
- **コード**:
  ```ts
  export async function signOutAction() {
    await signOut()
  }
  ```
- **問題点**: `signOut()`をパラメータなしで呼び出しているため、Auth0のセッションがクリアされない

#### 2.1.2 認証設定（auth.ts）
- **ファイル**: `client/auth.ts`
- **実装**: NextAuthの設定でAuth0プロバイダーを使用
- **環境変数**: `AUTH0_CLIENT_ID`, `AUTH0_CLIENT_SECRET`, `AUTH0_ISSUER`, `AUTH0_AUDIENCE`を使用

### 2.2 課題点
1. **Auth0セッションの残存**: ログアウトしてもAuth0のログイン情報が残ってしまう
2. **次回ログイン時の問題**: 前回Auth0ログインした時の情報が残ってしまっている
3. **不完全なログアウト**: next-authのセッションのみクリアされ、Auth0のセッションがクリアされない

### 2.3 本実装による改善点
1. **完全なログアウト**: Auth0のセッションも完全にクリアされる
2. **次回ログイン時の正常化**: 次回ログイン時に前回のログイン情報が残らない
3. **セキュリティ向上**: ログアウト時にすべてのセッション情報がクリアされる

## 3. 機能要件

### 3.1 ログアウト処理の修正

#### 3.1.1 signOutActionの修正

**概要**: `signOutAction`を修正して、Auth0のログアウトURLにリダイレクトするようにする。

**変更内容**:
- `signOut()`にパラメータを渡す
- `redirect: true`と`redirectTo`を指定
- `redirectTo`にはAuth0のログアウトURLを設定
- Auth0ログアウトURLには`returnTo`パラメータを追加して、ログアウト後にアプリケーションに戻るURLを指定

**変更前**:
```ts
export async function signOutAction() {
  await signOut()
}
```

**変更後**:
```ts
export async function signOutAction() {
  const auth0Issuer = process.env.AUTH0_ISSUER
  const appBaseUrl = process.env.NEXT_PUBLIC_APP_BASE_URL
  
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
  
  await signOut({
    redirect: true,
    redirectTo: logoutUrl
  })
}
```

#### 3.1.2 環境変数の設定

**概要**: `client/.env.local`に必要な環境変数を定義する。

**必要な環境変数**:
- `AUTH0_ISSUER`: Auth0のIssuer URL（既存の環境変数、例: `https://{domain}.auth0.com`）
- `NEXT_PUBLIC_APP_BASE_URL`: アプリケーションのベースURL（例: `http://localhost:3000`）

**環境変数の形式**:
```env
AUTH0_ISSUER=https://your-domain.auth0.com
NEXT_PUBLIC_APP_BASE_URL=http://localhost:3000
```

**注意事項**:
- `AUTH0_ISSUER`は既存の環境変数（`client/auth.ts`で使用されている）
- `AUTH0_ISSUER`から`/v2/logout`を追加してAuth0ログアウトURLを構築する（コード内で計算）
- `NEXT_PUBLIC_APP_BASE_URL`は`NEXT_PUBLIC_`プレフィックスが必要（クライアント側で使用するため）

### 3.2 エラーハンドリング

#### 3.2.1 環境変数の検証

**要件**: 環境変数が設定されていない場合はエラーを発生させる。

**実装**:
- `AUTH0_ISSUER`が設定されていない場合: `Error('AUTH0_ISSUER is not set')`をthrow
- `NEXT_PUBLIC_APP_BASE_URL`が設定されていない場合: `Error('NEXT_PUBLIC_APP_BASE_URL is not set')`をthrow
- フォールバック機能は実装しない（エラーを発生させる）

#### 3.2.2 エラーメッセージ

**要件**: エラーメッセージは明確で、設定が必要な環境変数を示す。

**実装**:
- エラーメッセージは日本語でも英語でも良い（既存のコードスタイルに合わせる）
- エラーメッセージには環境変数名を含める

### 3.3 Auth0ログアウトURLの構築

#### 3.3.1 Auth0ログアウトURLの構築

**要件**: `AUTH0_ISSUER`からAuth0ログアウトURLを構築し、`returnTo`パラメータを追加する。

**実装**:
- `AUTH0_ISSUER`から`/v2/logout`を追加してAuth0ログアウトURLを構築
- `returnTo`パラメータに`NEXT_PUBLIC_APP_BASE_URL`を使用
- URLエンコーディングを適用（`encodeURIComponent`を使用）

**URL形式**:
```
{AUTH0_ISSUER}/v2/logout?returnTo={NEXT_PUBLIC_APP_BASE_URL}/
```

**例**:
```
https://your-domain.auth0.com/v2/logout?returnTo=http%3A%2F%2Flocalhost%3A3000%2F
```

#### 3.3.2 リダイレクト先の決定

**要件**: ログアウト後、アプリケーションのトップページ（`/`）に戻る。

**実装**:
- `returnTo`パラメータには`${appBaseUrl}/`を設定
- スラッシュ（`/`）を含めることで、トップページにリダイレクトされる

## 4. 非機能要件

### 4.1 パフォーマンス
- **変更による影響**: パフォーマンスへの影響はない（ログアウト処理はユーザー操作に応じた1回限りの処理）
- **実行時間**: 既存の実装と同等の実行時間を維持（リダイレクト処理が追加されるが、ユーザー体験への影響は最小限）

### 4.2 互換性
- **既存機能の維持**: ログアウト機能の動作は変更しない（Auth0ログアウトURLへのリダイレクトが追加されるのみ）
- **ブラウザ互換性**: 既存のブラウザ互換性を維持
- **Next.js互換性**: 既存のNext.jsバージョンとの互換性を維持
- **next-auth互換性**: next-authの`signOut`関数の仕様に準拠

### 4.3 セキュリティ
- **認証処理**: 既存の認証処理を維持（Auth0ログアウトURLへのリダイレクトが追加されるのみ）
- **セキュリティリスク**: 変更によるセキュリティリスクはない（むしろ、Auth0セッションの完全なクリアによりセキュリティが向上）
- **環境変数の管理**: 環境変数は`client/.env.local`に保存（gitから除外される）

### 4.4 保守性
- **コードの一貫性**: 既存のコードスタイルを維持
- **変更の容易性**: 環境変数の変更により、Auth0ログアウトURLやリダイレクト先を変更可能
- **可読性**: コードの可読性を維持（環境変数の検証とURL構築を明確に実装）

### 4.5 動作環境
- **開発環境**: 既存の開発環境で動作する（環境変数の設定が必要）
- **本番環境**: 既存の本番環境で動作する（環境変数の設定が必要）
- **依存関係**: 既存の依存関係を維持

## 5. 制約事項

### 5.1 既存システムとの関係
- **認証システム**: 既存の認証システム（Auth0）との互換性を維持
- **Next.js Server Actions**: 既存のNext.js Server Actionsの仕組みを維持
- **next-auth**: 既存のnext-authの設定を維持（`signOut`関数のパラメータを追加するのみ）

### 5.2 技術スタック
- **Next.js**: 既存のNext.jsバージョンを維持
- **TypeScript**: 既存のTypeScriptバージョンを維持
- **next-auth**: 既存のnext-authバージョンを維持

### 5.3 依存関係
- **auth-actions.ts**: `signOutAction`が存在する必要がある
- **@/auth**: `signOut`関数が存在する必要がある
- **環境変数**: `AUTH0_ISSUER`（既存）と`NEXT_PUBLIC_APP_BASE_URL`が設定されている必要がある

### 5.4 運用上の制約
- **環境変数の設定**: 開発環境と本番環境の両方で`NEXT_PUBLIC_APP_BASE_URL`を設定する必要がある（`AUTH0_ISSUER`は既存）
- **テスト**: 既存のテストが正常に動作することを確認する必要がある（環境変数のモックが必要な場合がある）
- **デプロイ**: 既存のデプロイプロセスに影響を与えない（`NEXT_PUBLIC_APP_BASE_URL`の設定が必要）

## 6. 受け入れ基準

### 6.1 signOutActionの修正
- [ ] `signOutAction`で`signOut()`にパラメータを渡している
- [ ] `redirect: true`と`redirectTo`を指定している
- [ ] `AUTH0_ISSUER`から`/v2/logout`を追加してAuth0ログアウトURLを構築している
- [ ] `redirectTo`には構築したAuth0のログアウトURLを設定している
- [ ] Auth0ログアウトURLに`returnTo`パラメータを追加している
- [ ] `returnTo`パラメータにはアプリケーションのベースURLを設定している
- [ ] URLエンコーディングを適用している

### 6.2 環境変数の設定
- [ ] `AUTH0_ISSUER`が設定されている（既存の環境変数）
- [ ] `client/.env.local`に`NEXT_PUBLIC_APP_BASE_URL`が定義されている
- [ ] 環境変数の形式が正しい（`AUTH0_ISSUER`は`/v2/logout`を含めない）

### 6.3 エラーハンドリング
- [ ] `AUTH0_ISSUER`が設定されていない場合、エラーを発生させる
- [ ] `NEXT_PUBLIC_APP_BASE_URL`が設定されていない場合、エラーを発生させる
- [ ] エラーメッセージに環境変数名が含まれている
- [ ] フォールバック機能が実装されていない

### 6.4 動作確認
- [ ] ログアウトボタンをクリックすると、Auth0のログアウトURLにリダイレクトされる
- [ ] Auth0ログアウト後、アプリケーションのトップページに戻る
- [ ] ログアウト後、Auth0のセッションがクリアされている
- [ ] 次回ログイン時に前回のログイン情報が残らない
- [ ] 既存のログアウト機能に影響がない

### 6.5 コード品質
- [ ] コードの一貫性が保たれている
- [ ] 環境変数の検証が適切に実装されている
- [ ] URL構築が適切に実装されている
- [ ] エラーハンドリングが適切に実装されている

## 7. 影響範囲

### 7.1 変更が必要なファイル

#### 更新するファイル
- `client/lib/actions/auth-actions.ts`: `signOutAction`を修正

#### 新規作成するファイル
- なし

#### 環境変数ファイル
- `client/.env.local`: `NEXT_PUBLIC_APP_BASE_URL`を追加（既存ファイルに追加、または新規作成）
- 注意: `AUTH0_ISSUER`は既存の環境変数（`client/auth.ts`で使用されている）

### 7.2 既存ファイルの扱い
- `client/auth.ts`: 変更なし（既存のAuth0設定を維持）
- `client/components/auth/auth-buttons.tsx`: 変更なし（既存の`signOutAction`の使用を維持）
- `client/components/layout/navbar.tsx`: 変更なし（既存の`signOutAction`の使用を維持）

### 7.3 既存機能への影響
- **既存のログアウト機能**: 動作は変更される（Auth0ログアウトURLへのリダイレクトが追加される）
- **既存のコンポーネント**: 影響なし（`signOutAction`の呼び出し方法は変更されない）
- **既存のテスト**: 影響の可能性あり（環境変数のモックが必要な場合がある）
- **UIデザイン**: 影響なし（ログアウト処理の内部実装のみ変更）

## 8. 実装上の注意事項

### 8.1 signOutActionの修正
- **パラメータの形式**: `signOut({ redirect: true, redirectTo: logoutUrl })`の形式を使用
- **URL構築**: `AUTH0_ISSUER`から`/v2/logout`を追加してAuth0ログアウトURLを構築し、`returnTo`パラメータを追加する際、URLエンコーディングを適用
- **環境変数の取得**: `process.env.AUTH0_ISSUER`と`process.env.NEXT_PUBLIC_APP_BASE_URL`を使用
- **エラーハンドリング**: 環境変数が設定されていない場合はエラーをthrow（フォールバックは実装しない）

### 8.2 環境変数の設定
- **ファイル**: `client/.env.local`に環境変数を定義
- **形式**: `NEXT_PUBLIC_APP_BASE_URL=http://localhost:3000`（`NEXT_PUBLIC_`プレフィックスが必要）
- **注意**: `AUTH0_ISSUER`は既存の環境変数（`client/auth.ts`で使用されている）で、`/v2/logout`を含めない
- **注意**: `.env.local`はgitから除外されるため、開発環境と本番環境で個別に設定する必要がある

### 8.3 テストの確認
- **既存テストの動作確認**: 既存のテストが正常に動作することを確認
- **環境変数のモック**: テストで環境変数をモックする必要がある場合は、適切にモックを設定
- **統合テスト**: ログアウト処理の統合テストを実行し、Auth0ログアウトURLへのリダイレクトを確認

### 8.4 動作確認
- **ログアウト機能**: ログアウトボタンをクリックして、Auth0ログアウトURLにリダイレクトされることを確認
- **リダイレクト**: Auth0ログアウト後、アプリケーションのトップページに戻ることを確認
- **セッションクリア**: ログアウト後、Auth0のセッションがクリアされていることを確認（次回ログイン時に前回のログイン情報が残らないことを確認）
- **既存機能**: 既存のログアウト機能に影響がないことを確認

### 8.5 コードレビュー
- **実装の確認**: `signOutAction`で`signOut()`にパラメータを渡していることを確認
- **環境変数の検証**: 環境変数の検証が適切に実装されていることを確認
- **URL構築**: URL構築が適切に実装されていることを確認
- **エラーハンドリング**: エラーハンドリングが適切に実装されていることを確認

## 9. 既知の問題と対応計画

### 9.1 型チェックエラーの修正計画

#### 9.1.1 問題の概要
`cd client && npm run type-check 2>&1`を実行すると、以下の型チェックエラーが発生しています：

```
src/__tests__/integration/page-page.test.tsx(35,36): error TS2345: Argument of type 'null' is not assignable to parameter of type 'never'.
src/__tests__/integration/page-page.test.tsx(44,36): error TS2345: Argument of type 'null' is not assignable to parameter of type 'never'.
...
```

#### 9.1.2 原因分析
1. **jest.setup.jsのモック定義の問題**: `client/jest.setup.js`の53-64行目で`auth`関数をモックしているが、`jest.fn(() => null)`となっており、`Promise`を返していない
2. **型推論の問題**: `auth`関数の実際の戻り値は`Promise<Session | null>`であるが、モックの型定義が適切でないため、TypeScriptが型を正しく推論できていない
3. **テストファイルでの型アサーション**: `src/__tests__/integration/page-page.test.tsx`で`mockAuth.mockResolvedValueOnce(null)`を使用しているが、型が`never`として推論されている

#### 9.1.3 修正計画

**修正対象ファイル**:
1. `client/jest.setup.js`: `auth`関数のモック定義を修正
2. `client/src/__tests__/integration/page-page.test.tsx`: 型アサーションを修正

**修正内容**:

1. **jest.setup.jsの修正**:
   - `jest.fn(() => null)`を`jest.fn(() => Promise.resolve(null))`に変更
   - `auth`関数が`Promise<Session | null>`を返すようにモックを定義

2. **テストファイルの修正**:
   - `auth`関数の型を適切にアサーション
   - `Session`型をインポートして使用

**修正後のコード例**:

```javascript
// jest.setup.js
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

```typescript
// src/__tests__/integration/page-page.test.tsx
import { render, screen } from '@testing-library/react'
import Home from '@/app/page'
import { auth } from '@/auth'
import type { Session } from 'next-auth'

// Mock auth function
const mockAuth = auth as unknown as jest.MockedFunction<() => Promise<Session | null>>
```

#### 9.1.4 実装タイミング
- **実装フェーズ**: 設計フェーズまたは実装フェーズで対応
- **優先度**: 中（型チェックエラーは実装前に修正が必要）
- **影響範囲**: テストファイルのみ（本番コードへの影響なし）

#### 9.1.5 検証方法
修正後、以下のコマンドで型チェックが通ることを確認：
```bash
cd client && npm run type-check
```

## 10. 参考情報

### 9.1 関連Issue
- GitHub Issue #168: ログアウトしてもAuth0のログイン情報が残っている

### 10.2 既存ファイル
- `client/lib/actions/auth-actions.ts`: ログアウトアクション（修正対象）
- `client/auth.ts`: NextAuthの設定（`AUTH0_ISSUER`環境変数を使用、参考）
- `client/components/auth/auth-buttons.tsx`: ログアウトボタン（`signOutAction`を使用）
- `client/components/layout/navbar.tsx`: ヘッダーのログアウトボタン（`signOutAction`を使用）

### 10.3 技術スタック
- **Next.js**: Server Actionsを使用
- **TypeScript**: 型安全性を維持
- **next-auth**: 認証ライブラリ
- **Auth0**: 認証プロバイダー

### 10.4 参考リンク
- Next.js Server Actions: https://nextjs.org/docs/app/building-your-application/data-fetching/server-actions-and-mutations
- Next.js Authentication: https://nextjs.org/docs/app/building-your-application/authentication
- next-auth signOut: https://next-auth.js.org/getting-started/client#signout
- Auth0 Logout: https://auth0.com/docs/authenticate/login/logout

### 10.5 Auth0ログアウトURLの形式
Auth0のログアウトURLは以下の形式です：
```
{AUTH0_ISSUER}/v2/logout?returnTo={returnUrl}&client_id={clientId}
```

ただし、本実装では`client_id`パラメータは含めません（next-authが管理するため）。
`returnTo`パラメータのみを追加します。

`AUTH0_ISSUER`は既存の環境変数で、`client/auth.ts`で使用されています。
`AUTH0_ISSUER`から`/v2/logout`を追加してAuth0ログアウトURLを構築します。
