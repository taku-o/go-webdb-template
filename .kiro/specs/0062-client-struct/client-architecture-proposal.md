# クライアントアプリケーション設計改善提案

## 現状の問題点

### 1. 認証処理の分散
- `TodayApiButton.tsx` で直接 `/auth/token` を呼び出している
- 認証ロジックがコンポーネント内に散在している
- 認証トークンの取得方法が統一されていない

### 2. API呼び出しの分散
- `TodayApiButton.tsx` で直接 `fetch` を呼び出している
- `api.ts` に `getToday` メソッドは存在するが、JWTを引数で受け取る形になっている
- コンポーネントがAPI呼び出しの詳細を知る必要がある

### 3. `api.ts` の設計問題
- `ApiClient` は API Key を前提としているが、Auth0 JWT も扱う必要がある
- 認証トークンの取得ロジックが共通化されていない
- 認証方式（API Key vs JWT）の切り替えが不自然

## 提案する設計

### アーキテクチャ概要

```
┌─────────────────────────────────────────┐
│         React Components                 │
│  (TodayApiButton, etc.)                 │
│  - UI表示のみ                           │
│  - 状態管理（useState, etc.）           │
└──────────────┬──────────────────────────┘
               │
               │ apiClient.method()
               ▼
┌─────────────────────────────────────────┐
│         API Client (api.ts)              │
│  - 統一されたAPI呼び出しインターフェース │
│  - 認証トークンの自動取得・付与         │
│  - エラーハンドリング                    │
└──────────────┬──────────────────────────┘
               │
               │ getAuthToken()
               ▼
┌─────────────────────────────────────────┐
│      Auth Service (auth.ts)              │
│  - 認証トークンの取得ロジック            │
│  - Auth0 JWT / API Key の切り替え       │
│  - トークンキャッシュ（オプション）      │
└─────────────────────────────────────────┘
```

### 1. 認証サービスの共通化 (`client/src/lib/auth.ts`)

認証トークンの取得ロジックを一箇所に集約します。Auth0 SDKの標準的な方法を優先的に利用し、コードを簡素化します。

**注意**: `@auth0/nextjs-auth0`は既にインストール済みです。`handleAuth`を使用した標準的な認証エンドポイント（`/api/auth/[auth0]/route.ts`）を積極的に作成します。既存の`/auth/token`エンドポイントは必要に応じて維持または改善します。使えそうなライブラリ、フレームワークがあるなら積極的に利用します。

```typescript
// client/src/lib/auth.ts
import { User } from '@auth0/nextjs-auth0'

// Auth0のUser型をAuth0Userとしてエイリアス定義（他のUser型との衝突を避けるため）
type Auth0User = User

/**
 * 認証トークンを取得する
 * - ログイン中: Auth0 JWTを取得（既存の`/auth/token`エンドポイントを使用）
 * - 未ログイン: Public API Keyを使用
 */
export async function getAuthToken(auth0user: Auth0User | undefined): Promise<string> {
  if (auth0user) {
    // ログイン中: Auth0 JWTを取得（既存の`/auth/token`エンドポイントを使用）
    const response = await fetch('/auth/token')
    if (!response.ok) {
      throw new Error('Failed to get access token')
    }
    const data = await response.json()
    return data.accessToken
  } else {
    // 未ログイン: Public API Keyを使用
    const apiKey = process.env.NEXT_PUBLIC_API_KEY
    if (!apiKey) {
      throw new Error('NEXT_PUBLIC_API_KEY is not set')
    }
    return apiKey
  }
}
```

### 2. API Client の改善 (`client/src/lib/api.ts`)

認証トークンの取得を `ApiClient` 内部で処理するように変更します。

```typescript
// client/src/lib/api.ts
import { getAuthToken } from './auth'
import { User } from '@auth0/nextjs-auth0'

// Auth0のUser型をAuth0Userとしてエイリアス定義（他のUser型との衝突を避けるため）
type Auth0User = User

class ApiClient {
  private baseURL: string
  private apiKey: string | null

  constructor(baseURL: string) {
    this.baseURL = baseURL
    this.apiKey = process.env.NEXT_PUBLIC_API_KEY || null
  }

  /**
   * 認証トークンを取得してリクエストを実行
   * @param endpoint - APIエンドポイント
   * @param options - fetchオプション
   * @param auth0user - Auth0ユーザー情報（オプション）
   */
  private async request<T>(
    endpoint: string,
    options?: RequestInit,
    auth0user?: Auth0User | undefined
  ): Promise<T> {
    const url = `${this.baseURL}${endpoint}`

    // 認証トークンを取得
    const token = await getAuthToken(auth0user)

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

  // 既存のメソッドも必要に応じて修正（後方互換性は不要）
  // 注意: 戻り値のUser[]はアプリケーションのUser型（dm-users）なので、Auth0Userではない
  // メソッド名を本来あるべき名前に変更（getUsers → getDmUsers）
  async getDmUsers(limit = 20, offset = 0, auth0user?: Auth0User | undefined): Promise<User[]> {
    return this.request<User[]>(`/api/dm-users?limit=${limit}&offset=${offset}`, undefined, auth0user)
  }

  // 同様に、getPosts → getDmPosts など、他の既存メソッドも本来あるべき名前に変更

  // Today API (Auth0 JWT使用)
  async getToday(auth0user?: Auth0User | undefined): Promise<{ date: string }> {
    return this.request<{ date: string }>('/api/today', undefined, auth0user)
  }
}

export const apiClient = new ApiClient(API_BASE_URL)
```

### 3. コンポーネントの簡素化 (`client/src/components/TodayApiButton.tsx`)

コンポーネントから認証処理とAPI呼び出しの詳細を削除し、UI表示と状態管理のみに集中します。

```typescript
'use client'

import { useUser } from '@auth0/nextjs-auth0'
import { useState } from 'react'
import { apiClient } from '@/lib/api'

export default function TodayApiButton() {
  const { user: auth0user, isLoading } = useUser()
  const [date, setDate] = useState<string | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)

  const handleClick = async () => {
    setLoading(true)
    setError(null)
    setDate(null)

    try {
      const data = await apiClient.getToday(auth0user || undefined)
      setDate(data.date)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'An error occurred')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="p-4 border rounded-lg bg-gray-50">
      <h3 className="font-semibold mb-4">Today API (Private Endpoint)</h3>
      <p className="text-sm text-gray-600 mb-4">
        Auth0ログイン時のみアクセス可能なプライベートAPIをテストします。
      </p>
      <button
        onClick={handleClick}
        disabled={loading || isLoading}
        className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
      >
        {loading ? 'Loading...' : 'Get Today'}
      </button>
      {date && (
        <p className="mt-4 text-green-600">
          Today: {date}
        </p>
      )}
      {error && (
        <p className="mt-4 text-red-600">
          Error: {error}
        </p>
      )}
    </div>
  )
}
```

## 設計原則

### 1. 関心の分離 (Separation of Concerns)
- **コンポーネント**: UI表示と状態管理のみ
- **API Client**: API呼び出しの抽象化
- **Auth Service**: 認証トークンの取得ロジック

### 2. 単一責任の原則 (Single Responsibility Principle)
- 各モジュールは一つの責任のみを持つ
- 認証処理は `auth.ts` に集約
- API呼び出しは `api.ts` に集約

### 3. DRY原則 (Don't Repeat Yourself)
- 認証トークン取得ロジックを一箇所に集約
- API呼び出しパターンを統一

### 4. 依存性の逆転 (Dependency Inversion)
- コンポーネントは `apiClient` に依存
- `apiClient` は `auth` に依存
- 実装の詳細は隠蔽される

## 実装手順

1. **`client/src/lib/auth.ts` を作成**
   - `getAuthToken` 関数を実装

2. **`client/src/lib/api.ts` を修正**
   - `request` メソッドを修正して `auth0user` パラメータを受け取る
   - `getToday` メソッドを修正して `auth0user` パラメータを受け取る
   - 既存のメソッドも必要に応じて修正（`auth0user` はオプショナル）

3. **`client/src/components/TodayApiButton.tsx` を修正**
   - 認証処理を削除
   - API呼び出しを `apiClient.getToday` に変更

4. **テストの更新**
   - `api.test.ts` を更新
   - `TodayApiButton` のテストを更新（必要に応じて）

5. **他のコンポーネントの確認**
   - 同様の問題がないか確認
   - 必要に応じて修正

## 期待される効果

1. **保守性の向上**
   - 認証ロジックの変更が一箇所で済む
   - API呼び出しパターンが統一される

2. **テスタビリティの向上**
   - 各モジュールを独立してテスト可能
   - モックが容易になる

3. **再利用性の向上**
   - 認証処理を他のコンポーネントでも再利用可能
   - API呼び出しパターンを他のエンドポイントでも適用可能

4. **コードの可読性向上**
   - コンポーネントが簡潔になる
   - 責務が明確になる

## 注意事項

- **基本方針**: 使えそうなライブラリ、フレームワークがあるなら積極的に利用する
- **後方互換性**: 不要（既存コードの変更も可）
- `auth0user` パラメータはオプショナルとし、未指定時は API Key を使用
- エラーハンドリングは統一された形式を維持
- Auth0 SDKの標準的な方法（`handleAuth`など）を積極的に利用する
