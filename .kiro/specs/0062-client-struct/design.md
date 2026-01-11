# クライアントアプリケーション設計改善の設計書

## 1. 概要

### 1.1 設計の目的
要件定義書に基づき、クライアントアプリケーション（Next.js）の設計を改善し、認証処理とAPI呼び出しの責務を適切に分離するための詳細設計を定義する。これにより、コードの保守性、テスタビリティ、再利用性を向上させ、コンポーネントとビジネスロジックの関心を分離する。

### 1.2 設計の範囲
- 認証サービスの共通化設計（`client/src/lib/auth.ts`の作成）
- API Clientの改善設計（`client/src/lib/api.ts`の修正）
- コンポーネントの簡素化設計（`client/src/components/TodayApiButton.tsx`の修正）
- Auth0 SDKの標準エンドポイント設計（`client/src/app/api/auth/[auth0]/route.ts`の作成）
- テスト設計
- 既存メソッド名の変更設計（`getUsers` → `getDmUsers`、`getPosts` → `getDmPosts`など）

### 1.3 設計方針
- **関心の分離**: コンポーネント、API Client、Auth Serviceの責務を明確化
- **ライブラリの積極的利用**: 使えそうなライブラリ、フレームワークがあるなら積極的に利用する
- **後方互換性**: 不要（既存コードの変更も可）
- **Auth0 SDKの標準的な方法**: `handleAuth`を使用した標準的な認証エンドポイントを積極的に作成
- **型名・変数名の明確化**: `User`型を`Auth0User`型に、`user`変数を`auth0user`に変更

## 2. アーキテクチャ設計

### 2.1 全体構成

#### 2.1.1 アーキテクチャ概要

```
┌─────────────────────────────────────────┐
│         React Components                 │
│  (TodayApiButton, etc.)                 │
│  - UI表示のみ                           │
│  - 状態管理（useState, etc.）           │
│  - イベントハンドリング                 │
└──────────────┬──────────────────────────┘
               │
               │ apiClient.method(auth0user)
               ▼
┌─────────────────────────────────────────┐
│         API Client (api.ts)              │
│  - 統一されたAPI呼び出しインターフェース │
│  - 認証トークンの自動取得・付与         │
│  - エラーハンドリング                    │
└──────────────┬──────────────────────────┘
               │
               │ getAuthToken(auth0user)
               ▼
┌─────────────────────────────────────────┐
│      Auth Service (auth.ts)              │
│  - 認証トークンの取得ロジック            │
│  - Auth0 JWT / API Key の切り替え       │
└─────────────────────────────────────────┘
               │
               │ /auth/token (既存)
               ▼
┌─────────────────────────────────────────┐
│   Auth0 SDK Standard Endpoints           │
│  /api/auth/[auth0]/route.ts              │
│  - handleAuth()                          │
│  - /api/auth/login                       │
│  - /api/auth/logout                      │
└─────────────────────────────────────────┘
```

#### 2.1.2 ディレクトリ構造

**変更前**:
```
client/src/
├── lib/
│   ├── api.ts                    # API Client（認証処理が分散）
│   └── auth0.ts                  # サーバー側のAuth0設定
├── components/
│   └── TodayApiButton.tsx        # 認証処理とAPI呼び出しが直接実装
└── app/
    └── auth/
        └── token/
            └── route.ts           # アクセストークン取得エンドポイント
```

**変更後**:
```
client/src/
├── lib/
│   ├── api.ts                    # API Client（認証処理を内部で処理）
│   ├── auth.ts                   # 認証トークン取得ロジック（新規作成）
│   └── auth0.ts                  # サーバー側のAuth0設定（既存）
├── components/
│   └── TodayApiButton.tsx        # UI表示と状態管理のみ
└── app/
    ├── api/
    │   └── auth/
    │       └── [auth0]/
    │           └── route.ts      # Auth0 SDK標準エンドポイント（新規作成）
    └── auth/
        └── token/
            └── route.ts           # アクセストークン取得エンドポイント（既存、維持）
```

### 2.2 モジュール設計

#### 2.2.1 Auth Service (`client/src/lib/auth.ts`)

**責務**: 認証トークンの取得ロジックを一箇所に集約

**型定義**:
```typescript
import { User } from '@auth0/nextjs-auth0'

// Auth0のUser型をAuth0Userとしてエイリアス定義（他のUser型との衝突を避けるため）
type Auth0User = User
```

**関数シグネチャ**:
```typescript
export async function getAuthToken(auth0user: Auth0User | undefined): Promise<string>
```

**処理フロー**:
1. `auth0user`が存在する場合:
   - `/auth/token`エンドポイントを呼び出し
   - レスポンスから`accessToken`を取得
   - エラー時は適切なエラーメッセージを投げる
2. `auth0user`が存在しない場合:
   - `process.env.NEXT_PUBLIC_API_KEY`を取得
   - 環境変数が設定されていない場合はエラーを投げる

**エラーハンドリング**:
- `/auth/token`の呼び出し失敗: `Error('Failed to get access token')`
- `NEXT_PUBLIC_API_KEY`が未設定: `Error('NEXT_PUBLIC_API_KEY is not set')`

#### 2.2.2 API Client (`client/src/lib/api.ts`)

**責務**: 統一されたAPI呼び出しインターフェースを提供

**型定義**:
```typescript
import { User } from '@auth0/nextjs-auth0'

// Auth0のUser型をAuth0Userとしてエイリアス定義
type Auth0User = User
```

**クラス構造**:
```typescript
class ApiClient {
  private baseURL: string
  private apiKey: string | null

  constructor(baseURL: string)
  private async request<T>(endpoint: string, options?: RequestInit, auth0user?: Auth0User | undefined): Promise<T>
  
  // 既存メソッド（メソッド名を変更）
  async getDmUsers(limit?: number, offset?: number, auth0user?: Auth0User | undefined): Promise<DmUser[]>
  async getDmUser(id: string, auth0user?: Auth0User | undefined): Promise<DmUser>
  async createDmUser(data: CreateUserRequest, auth0user?: Auth0User | undefined): Promise<DmUser>
  async updateDmUser(id: string, data: UpdateUserRequest, auth0user?: Auth0User | undefined): Promise<DmUser>
  async deleteDmUser(id: string, auth0user?: Auth0User | undefined): Promise<void>
  
  async getDmPosts(limit?: number, offset?: number, userId?: string, auth0user?: Auth0User | undefined): Promise<DmPost[]>
  async getDmPost(id: string, userId: string, auth0user?: Auth0User | undefined): Promise<DmPost>
  async createDmPost(data: CreatePostRequest, auth0user?: Auth0User | undefined): Promise<DmPost>
  async updateDmPost(id: string, userId: string, data: UpdatePostRequest, auth0user?: Auth0User | undefined): Promise<DmPost>
  async deleteDmPost(id: string, userId: string, auth0user?: Auth0User | undefined): Promise<void>
  
  async getDmUserPosts(limit?: number, offset?: number, auth0user?: Auth0User | undefined): Promise<UserPost[]>
  
  async getToday(auth0user?: Auth0User | undefined): Promise<{ date: string }>
  
  async sendEmail(to: string[], template: string, data: Record<string, unknown>, auth0user?: Auth0User | undefined): Promise<{ success: boolean; message: string }>
  
  async registerJob(data?: RegisterJobRequest, auth0user?: Auth0User | undefined): Promise<RegisterJobResponse>
  
  async downloadDmUsersCSV(auth0user?: Auth0User | undefined): Promise<void>
}
```

**requestメソッドの処理フロー**:
1. `getAuthToken(auth0user)`を呼び出してトークンを取得
2. Authorizationヘッダーに`Bearer ${token}`を付与
3. `fetch`でAPI呼び出しを実行
4. レスポンスのステータスコードを確認
5. エラー時は適切なエラーメッセージを投げる
6. 成功時はJSONをパースして返却

**エラーハンドリング**:
- 401/403エラー: `Error(errorData.message || response.statusText)`
- その他のエラー: `Error(error || response.statusText)`

#### 2.2.3 コンポーネント (`client/src/components/TodayApiButton.tsx`)

**責務**: UI表示と状態管理のみ

**変更内容**:
- 認証処理（`/auth/token`の呼び出し）を削除
- API呼び出し（`fetch`の直接呼び出し）を削除
- `apiClient.getToday(auth0user || undefined)`を使用
- `useUser()`から取得した変数名を`auth0user`に変更

**実装例**:
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

  // UI表示のみ
  return (
    // ... JSX ...
  )
}
```

#### 2.2.4 Auth0 SDK標準エンドポイント (`client/src/app/api/auth/[auth0]/route.ts`)

**責務**: Auth0 SDKの標準的な認証エンドポイントを提供

**実装**:
```typescript
import { handleAuth } from '@auth0/nextjs-auth0'

export const GET = handleAuth()
```

**自動生成されるエンドポイント**:
- `/api/auth/login`: ログイン
- `/api/auth/logout`: ログアウト
- `/api/auth/callback`: コールバック
- `/api/auth/me`: ユーザー情報取得

## 3. 詳細設計

### 3.1 認証サービスの実装設計

#### 3.1.1 auth.tsの実装

**ファイル**: `client/src/lib/auth.ts`

**実装内容**:
```typescript
import { User } from '@auth0/nextjs-auth0'

// Auth0のUser型をAuth0Userとしてエイリアス定義（他のUser型との衝突を避けるため）
type Auth0User = User

/**
 * 認証トークンを取得する
 * - ログイン中: Auth0 JWTを取得（既存の`/auth/token`エンドポイントを使用）
 * - 未ログイン: Public API Keyを使用
 * 
 * @param auth0user - Auth0のユーザー情報（オプション）
 * @returns 認証トークン（JWTまたはAPI Key）
 * @throws Error - トークン取得失敗時
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

**エクスポート**:
- `getAuthToken`: 認証トークン取得関数
- `Auth0User`: 型エイリアス（必要に応じて）

### 3.2 API Clientの実装設計

#### 3.2.1 api.tsの修正

**ファイル**: `client/src/lib/api.ts`

**変更内容**:

1. **型定義の追加**:
```typescript
import { User } from '@auth0/nextjs-auth0'
import { getAuthToken } from './auth'

// Auth0のUser型をAuth0Userとしてエイリアス定義
type Auth0User = User
```

2. **requestメソッドの修正**:
```typescript
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
```

3. **既存メソッドの修正**:
- メソッド名を変更: `getUsers` → `getDmUsers`、`getPosts` → `getDmPosts`など
- `auth0user`パラメータを追加（オプショナル）
- `request`メソッドに`auth0user`を渡す

**修正例（getDmUsers）**:
```typescript
// 修正前
async getUsers(limit = 20, offset = 0): Promise<User[]> {
  return this.request<User[]>(`/api/dm-users?limit=${limit}&offset=${offset}`)
}

// 修正後
async getDmUsers(limit = 20, offset = 0, auth0user?: Auth0User | undefined): Promise<DmUser[]> {
  return this.request<DmUser[]>(`/api/dm-users?limit=${limit}&offset=${offset}`, undefined, auth0user)
}
```

**修正例（getToday）**:
```typescript
// 修正前
async getToday(jwt: string): Promise<{ date: string }> {
  return this.request<{ date: string }>('/api/today', undefined, jwt)
}

// 修正後
async getToday(auth0user?: Auth0User | undefined): Promise<{ date: string }> {
  return this.request<{ date: string }>('/api/today', undefined, auth0user)
}
```

### 3.3 コンポーネントの実装設計

#### 3.3.1 TodayApiButton.tsxの修正

**ファイル**: `client/src/components/TodayApiButton.tsx`

**変更前の実装**:
- 認証処理（`/auth/token`の呼び出し）がコンポーネント内に直接実装
- API呼び出し（`fetch`の直接呼び出し）がコンポーネント内に直接実装

**変更後の実装**:
- 認証処理を削除
- API呼び出しを`apiClient.getToday(auth0user || undefined)`に変更
- `useUser()`から取得した変数名を`auth0user`に変更

**実装例**:
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

### 3.4 Auth0 SDK標準エンドポイントの実装設計

#### 3.4.1 /api/auth/[auth0]/route.tsの作成

**ファイル**: `client/src/app/api/auth/[auth0]/route.ts`

**実装**:
```typescript
import { handleAuth } from '@auth0/nextjs-auth0'

export const GET = handleAuth()
```

**自動生成されるエンドポイント**:
- `GET /api/auth/login`: ログインページにリダイレクト
- `GET /api/auth/logout`: ログアウト処理
- `GET /api/auth/callback`: OAuthコールバック処理
- `GET /api/auth/me`: 現在のユーザー情報を取得

**設定**:
- Auth0 SDKの設定は既存の`Auth0Provider`（`layout.tsx`）で管理
- 環境変数は既存の設定を使用

## 4. データフロー設計

### 4.1 認証トークン取得のフロー

```
1. コンポーネント
   ↓ apiClient.getToday(auth0user)
   
2. API Client
   ↓ request('/api/today', undefined, auth0user)
   ↓ getAuthToken(auth0user)
   
3. Auth Service
   ├─ auth0userが存在する場合
   │  ↓ fetch('/auth/token')
   │  ↓ サーバー側エンドポイント
   │  ↓ auth0.getAccessToken()
   │  ↓ Auth0 JWTを取得
   │
   └─ auth0userが存在しない場合
      ↓ process.env.NEXT_PUBLIC_API_KEY
      ↓ API Keyを取得
   
4. API Client
   ↓ Authorization: Bearer ${token}
   ↓ fetch(`${baseURL}/api/today`)
   
5. サーバー
   ↓ 認証・認可処理
   ↓ レスポンス返却
```

### 4.2 エラーハンドリングのフロー

```
1. 認証トークン取得エラー
   getAuthToken()
   ├─ /auth/token呼び出し失敗
   │  ↓ Error('Failed to get access token')
   │
   └─ NEXT_PUBLIC_API_KEY未設定
      ↓ Error('NEXT_PUBLIC_API_KEY is not set')

2. API呼び出しエラー
   request()
   ├─ 401/403エラー
   │  ↓ Error(errorData.message || response.statusText)
   │
   └─ その他のエラー
      ↓ Error(error || response.statusText)

3. コンポーネントでのエラー処理
   try-catch
   ↓ setError(err.message)
   ↓ UIにエラーメッセージを表示
```

## 5. 型設計

### 5.1 型定義

#### 5.1.1 Auth0User型
```typescript
import { User } from '@auth0/nextjs-auth0'
type Auth0User = User
```

**理由**: アプリケーションの`User`型（dm-users）とAuth0の`User`型を区別するため

#### 5.1.2 主要な関数・メソッドの型シグネチャ

```typescript
// auth.ts
export async function getAuthToken(auth0user: Auth0User | undefined): Promise<string>

// api.ts
private async request<T>(
  endpoint: string,
  options?: RequestInit,
  auth0user?: Auth0User | undefined
): Promise<T>

async getToday(auth0user?: Auth0User | undefined): Promise<{ date: string }>
async getDmUsers(limit?: number, offset?: number, auth0user?: Auth0User | undefined): Promise<DmUser[]>
async getDmPosts(limit?: number, offset?: number, userId?: string, auth0user?: Auth0User | undefined): Promise<DmPost[]>
// ... 他のメソッドも同様
```

**注意**: 戻り値の`DmUser[]`、`DmPost[]`はアプリケーションの型（`@/types/user`、`@/types/post`）であり、Auth0Userではない

## 6. エラーハンドリング設計

### 6.1 エラー種別

#### 6.1.1 認証トークン取得エラー
- **エラー**: `/auth/token`の呼び出し失敗
- **メッセージ**: `'Failed to get access token'`
- **発生箇所**: `auth.ts`の`getAuthToken`関数

#### 6.1.2 環境変数エラー
- **エラー**: `NEXT_PUBLIC_API_KEY`が未設定
- **メッセージ**: `'NEXT_PUBLIC_API_KEY is not set'`
- **発生箇所**: `auth.ts`の`getAuthToken`関数

#### 6.1.3 API呼び出しエラー
- **401/403エラー**: 認証・認可エラー
- **メッセージ**: `errorData.message || response.statusText`
- **発生箇所**: `api.ts`の`request`メソッド

#### 6.1.4 その他のAPIエラー
- **エラー**: その他のHTTPエラー
- **メッセージ**: `error || response.statusText`
- **発生箇所**: `api.ts`の`request`メソッド

### 6.2 エラーハンドリングの実装

#### 6.2.1 auth.tsでのエラーハンドリング
```typescript
export async function getAuthToken(auth0user: Auth0User | undefined): Promise<string> {
  if (auth0user) {
    const response = await fetch('/auth/token')
    if (!response.ok) {
      throw new Error('Failed to get access token')
    }
    const data = await response.json()
    return data.accessToken
  } else {
    const apiKey = process.env.NEXT_PUBLIC_API_KEY
    if (!apiKey) {
      throw new Error('NEXT_PUBLIC_API_KEY is not set')
    }
    return apiKey
  }
}
```

#### 6.2.2 api.tsでのエラーハンドリング
```typescript
if (!response.ok) {
  if (response.status === 401 || response.status === 403) {
    const errorData = await response.json().catch(() => ({}))
    throw new Error(errorData.message || response.statusText)
  }
  const error = await response.text()
  throw new Error(error || response.statusText)
}
```

#### 6.2.3 コンポーネントでのエラーハンドリング
```typescript
try {
  const data = await apiClient.getToday(auth0user || undefined)
  setDate(data.date)
} catch (err) {
  setError(err instanceof Error ? err.message : 'An error occurred')
} finally {
  setLoading(false)
}
```

## 7. テスト設計

### 7.1 単体テスト

#### 7.1.1 auth.tsのテスト
- **テスト対象**: `getAuthToken`関数
- **テストケース**:
  - `auth0user`が存在する場合、`/auth/token`を呼び出してJWTを取得
  - `auth0user`が存在しない場合、`NEXT_PUBLIC_API_KEY`を返す
  - `/auth/token`の呼び出し失敗時にエラーを投げる
  - `NEXT_PUBLIC_API_KEY`が未設定の場合にエラーを投げる

#### 7.1.2 api.tsのテスト
- **テスト対象**: `request`メソッド、各APIメソッド
- **テストケース**:
  - `auth0user`が渡された場合、`getAuthToken`が呼び出される
  - `auth0user`が未指定の場合、`getAuthToken(undefined)`が呼び出される
  - 401/403エラー時に適切なエラーメッセージを投げる
  - その他のエラー時に適切なエラーメッセージを投げる
  - 成功時にJSONをパースして返却する

#### 7.1.3 TodayApiButton.tsxのテスト
- **テスト対象**: コンポーネントの動作
- **テストケース**:
  - `apiClient.getToday`が呼び出される
  - 成功時に日付が表示される
  - エラー時にエラーメッセージが表示される
  - ローディング状態が適切に管理される

### 7.2 統合テスト

#### 7.2.1 認証フローのテスト
- **テスト対象**: 認証トークン取得からAPI呼び出しまでの一連の流れ
- **テストケース**:
  - ログイン中にAuth0 JWTが使用される
  - 未ログイン時にAPI Keyが使用される
  - エラーが適切に伝播される

## 8. 実装上の注意事項

### 8.1 型定義の注意点
- `Auth0User`型は`@auth0/nextjs-auth0`の`User`型のエイリアス
- アプリケーションの`DmUser`型（`@/types/user`）と混同しない
- 戻り値の`DmUser[]`、`DmPost[]`はアプリケーションの型

### 8.2 変数名の注意点
- `useUser()`から取得した変数名は`auth0user`に変更
- 将来的に`user`変数を使う可能性があるため、`auth0user`を使用

### 8.3 メソッド名の注意点
- `getUsers` → `getDmUsers`、`getPosts` → `getDmPosts`など、本来あるべき名前に変更
- 既存のコンポーネントで使用している場合は、呼び出し側も修正が必要

### 8.4 Auth0 SDKの利用
- `handleAuth`を使用した標準的な認証エンドポイントを積極的に作成
- 既存の`/auth/token`エンドポイントは維持（アクセストークン取得用）

### 8.5 後方互換性
- 後方互換性は不要（既存コードの変更も可）
- 既存のコンポーネントで`apiClient.getUsers`などを使用している場合は、`apiClient.getDmUsers`に変更が必要

## 9. 実装順序

### 9.1 実装の優先順位

1. **Auth0 SDK標準エンドポイントの作成** (`/api/auth/[auth0]/route.ts`)
   - 最も簡単で、標準的な方法を利用

2. **認証サービスの作成** (`auth.ts`)
   - 他のモジュールの基盤となる

3. **API Clientの修正** (`api.ts`)
   - 認証サービスに依存

4. **コンポーネントの修正** (`TodayApiButton.tsx`)
   - API Clientに依存

5. **テストの更新**
   - 各モジュールの実装後にテストを追加

### 9.2 実装手順

1. `client/src/app/api/auth/[auth0]/route.ts`を作成
2. `client/src/lib/auth.ts`を作成
3. `client/src/lib/api.ts`を修正
   - 型定義を追加
   - `request`メソッドを修正
   - 既存メソッド名を変更し、`auth0user`パラメータを追加
4. `client/src/components/TodayApiButton.tsx`を修正
5. 既存のコンポーネントで`apiClient.getUsers`などを使用している場合は修正
6. テストを更新・追加

## 10. 参考情報

### 10.1 関連ドキュメント
- 要件定義書: `.kiro/specs/0062-client-struct/requirements.md`
- 設計提案: `.kiro/specs/0062-client-struct/client-architecture-proposal.md`
- セッションサマリー: `.kiro/specs/0062-client-struct/session-summary.md`

### 10.2 技術スタック
- **フレームワーク**: Next.js 14+ (App Router)
- **言語**: TypeScript 5+
- **認証**: Auth0 Next.js SDK (`@auth0/nextjs-auth0`)
- **テスト**: Jest、React Testing Library

### 10.3 Auth0 SDKの標準的な方法
- **`handleAuth`**: 標準的な認証エンドポイントを自動生成
- **`useUser`**: クライアントコンポーネントでユーザー情報を取得
- **`getAccessToken`**: サーバーコンポーネントでアクセストークンを取得
