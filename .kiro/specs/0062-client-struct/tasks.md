# クライアントアプリケーション設計改善の実装タスク一覧

## 概要
クライアントアプリケーション（Next.js）の設計を改善し、認証処理とAPI呼び出しの責務を適切に分離するためのタスク一覧。要件定義書と設計書に基づき、実装可能な粒度でタスクを分解。

## 実装フェーズ

### Phase 0: 型定義ファイルの変更

#### タスク 0.1: user.tsの型名変更とファイル名変更
**目的**: 型定義ファイルの型名を`User`から`DmUser`に変更し、リクエスト型も変更する。また、ファイル名も変更する。

**作業内容**:
- `client/src/types/user.ts`を`client/src/types/dm_user.ts`にリネーム
- `User`型を`DmUser`に変更
- `CreateUserRequest`型を`CreateDmUserRequest`に変更
- `UpdateUserRequest`型を`UpdateDmUserRequest`に変更
- エクスポート名も変更

**実装内容**:
- ファイルリネーム: `client/src/types/user.ts` → `client/src/types/dm_user.ts`
- `export interface User` → `export interface DmUser`に変更
- `export interface CreateUserRequest` → `export interface CreateDmUserRequest`に変更
- `export interface UpdateUserRequest` → `export interface UpdateDmUserRequest`に変更
- 型名を`DmUser`、`CreateDmUserRequest`、`UpdateDmUserRequest`に統一

**実装例**:
```typescript
// 修正前
export interface User {
  id: string
  name: string
  email: string
  created_at: string
  updated_at: string
}

export interface CreateUserRequest {
  name: string
  email: string
}

export interface UpdateUserRequest {
  name?: string
  email?: string
}

// 修正後
export interface DmUser {
  id: string
  name: string
  email: string
  created_at: string
  updated_at: string
}

export interface CreateDmUserRequest {
  name: string
  email: string
}

export interface UpdateDmUserRequest {
  name?: string
  email?: string
}
```

**受け入れ基準**:
- [ ] `client/src/types/user.ts`が`client/src/types/dm_user.ts`にリネームされている
- [ ] `client/src/types/dm_user.ts`の`User`型が`DmUser`に変更されている
- [ ] `CreateUserRequest`型が`CreateDmUserRequest`に変更されている
- [ ] `UpdateUserRequest`型が`UpdateDmUserRequest`に変更されている
- [ ] エクスポート名が`DmUser`、`CreateDmUserRequest`、`UpdateDmUserRequest`になっている

- _Requirements: 3.2.4_
- _Design: 3.2.1_

---

#### タスク 0.2: post.tsの型名変更とファイル名変更
**目的**: 型定義ファイルの型名を`Post`から`DmPost`に変更し、リクエスト型も変更する。また、ファイル名も変更する。

**作業内容**:
- `client/src/types/post.ts`を`client/src/types/dm_post.ts`にリネーム
- `Post`型を`DmPost`に変更
- `CreatePostRequest`型を`CreateDmPostRequest`に変更
- `UpdatePostRequest`型を`UpdateDmPostRequest`に変更
- エクスポート名も変更

**実装内容**:
- ファイルリネーム: `client/src/types/post.ts` → `client/src/types/dm_post.ts`
- `export interface Post` → `export interface DmPost`に変更
- `export interface CreatePostRequest` → `export interface CreateDmPostRequest`に変更
- `export interface UpdatePostRequest` → `export interface UpdateDmPostRequest`に変更
- 型名を`DmPost`、`CreateDmPostRequest`、`UpdateDmPostRequest`に統一

**実装例**:
```typescript
// 修正前
export interface Post {
  id: string
  user_id: string
  title: string
  content: string
  created_at: string
  updated_at: string
}

export interface CreatePostRequest {
  user_id: string
  title: string
  content: string
}

export interface UpdatePostRequest {
  title?: string
  content?: string
}

// 修正後
export interface DmPost {
  id: string
  user_id: string
  title: string
  content: string
  created_at: string
  updated_at: string
}

export interface CreateDmPostRequest {
  user_id: string
  title: string
  content: string
}

export interface UpdateDmPostRequest {
  title?: string
  content?: string
}
```

**受け入れ基準**:
- [ ] `client/src/types/post.ts`が`client/src/types/dm_post.ts`にリネームされている
- [ ] `client/src/types/dm_post.ts`の`Post`型が`DmPost`に変更されている
- [ ] `CreatePostRequest`型が`CreateDmPostRequest`に変更されている
- [ ] `UpdatePostRequest`型が`UpdateDmPostRequest`に変更されている
- [ ] エクスポート名が`DmPost`、`CreateDmPostRequest`、`UpdateDmPostRequest`になっている

- _Requirements: 3.2.4_
- _Design: 3.2.1_

---

#### タスク 0.3: 型定義を使用しているファイルのインポート更新
**目的**: 型定義ファイルの変更（ファイル名変更、型名変更）に伴い、インポートパスと型の使用箇所を更新する。

**作業内容**:
- `client/src/lib/api.ts`のインポートパスと型の使用箇所を更新
- `client/src/lib/__tests__/api.test.ts`のインポートパスと型の使用箇所を更新
- `client/src/app/dm-users/page.tsx`のインポートパスを更新
- `client/src/app/dm-posts/page.tsx`のインポートパスを更新
- `client/src/app/dm-user-posts/page.tsx`のインポートパスを更新

**実装内容**:
- `import { User, CreateUserRequest, UpdateUserRequest } from '@/types/user'` → `import { DmUser, CreateDmUserRequest, UpdateDmUserRequest } from '@/types/dm_user'`に変更
- `import { Post, CreatePostRequest, UpdatePostRequest, UserPost } from '@/types/post'` → `import { DmPost, CreateDmPostRequest, UpdateDmPostRequest, UserPost } from '@/types/dm_post'`に変更
- 型の使用箇所も更新（`User` → `DmUser`、`Post` → `DmPost`、`CreateUserRequest` → `CreateDmUserRequest`、`UpdateUserRequest` → `UpdateDmUserRequest`、`CreatePostRequest` → `CreateDmPostRequest`、`UpdatePostRequest` → `UpdateDmPostRequest`）

**実装例**:
```typescript
// 修正前（api.ts）
import { User, CreateUserRequest, UpdateUserRequest } from '@/types/user'
import { Post, CreatePostRequest, UpdatePostRequest, UserPost } from '@/types/post'

async createUser(data: CreateUserRequest): Promise<User> {
  // ...
}

async updateUser(id: string, data: UpdateUserRequest): Promise<User> {
  // ...
}

async createPost(data: CreatePostRequest): Promise<Post> {
  // ...
}

async updatePost(id: string, userId: string, data: UpdatePostRequest): Promise<Post> {
  // ...
}

// 修正後（api.ts）
import { DmUser, CreateDmUserRequest, UpdateDmUserRequest } from '@/types/dm_user'
import { DmPost, CreateDmPostRequest, UpdateDmPostRequest, UserPost } from '@/types/dm_post'

async createDmUser(data: CreateDmUserRequest): Promise<DmUser> {
  // ...
}

async updateDmUser(id: string, data: UpdateDmUserRequest): Promise<DmUser> {
  // ...
}

async createDmPost(data: CreateDmPostRequest): Promise<DmPost> {
  // ...
}

async updateDmPost(id: string, userId: string, data: UpdateDmPostRequest): Promise<DmPost> {
  // ...
}
```

**受け入れ基準**:
- [ ] `client/src/lib/api.ts`のインポートパスが`@/types/dm_user`、`@/types/dm_post`に更新されている
- [ ] `client/src/lib/api.ts`の型の使用箇所（メソッドのパラメータ、戻り値の型など）が更新されている
- [ ] `client/src/lib/__tests__/api.test.ts`のインポートパスが`@/types/dm_user`、`@/types/dm_post`に更新されている
- [ ] `client/src/lib/__tests__/api.test.ts`の型の使用箇所が更新されている
- [ ] `client/src/app/dm-users/page.tsx`のインポートパスが`@/types/dm_user`に更新されている
- [ ] `client/src/app/dm-posts/page.tsx`のインポートパスが`@/types/dm_post`、`@/types/dm_user`に更新されている
- [ ] `client/src/app/dm-user-posts/page.tsx`のインポートパスが`@/types/dm_post`に更新されている
- [ ] 型の使用箇所（変数宣言、戻り値の型など）が更新されている

- _Requirements: 3.2.4_
- _Design: 3.2.1_

---

### Phase 1: Auth0 SDK標準エンドポイントの作成

#### タスク 1.1: Auth0 SDK標準エンドポイントの作成
**目的**: Auth0 SDKの標準的な方法を利用して、標準的な認証エンドポイントを自動生成する。

**作業内容**:
- `client/src/app/api/auth/[auth0]/route.ts`を新規作成
- `handleAuth`を使用して標準的な認証エンドポイントを自動生成

**実装内容**:
- ファイル作成: `client/src/app/api/auth/[auth0]/route.ts`
- `handleAuth`をインポートしてエクスポート
- これにより、`/api/auth/login`、`/api/auth/logout`などの標準エンドポイントが自動的に作成される

**実装例**:
```typescript
import { handleAuth } from '@auth0/nextjs-auth0'

export const GET = handleAuth()
```

**受け入れ基準**:
- [ ] `client/src/app/api/auth/[auth0]/route.ts`が作成されている
- [ ] `handleAuth`を使用して標準的な認証エンドポイントが自動生成されている
- [ ] `/api/auth/login`、`/api/auth/logout`などのエンドポイントが利用可能である

- _Requirements: 3.1.1_
- _Design: 3.4.1_

---

### Phase 2: 認証サービスの作成

#### タスク 2.1: auth.tsの作成
**目的**: 認証トークン取得ロジックを一箇所に集約する。

**作業内容**:
- `client/src/lib/auth.ts`を新規作成
- `getAuthToken`関数を実装
- `Auth0User`型をエイリアス定義

**実装内容**:
- ファイル作成: `client/src/lib/auth.ts`
- `@auth0/nextjs-auth0`から`User`をインポート
- `type Auth0User = User`としてエイリアス定義
- `getAuthToken(auth0user: Auth0User | undefined): Promise<string>`関数を実装
  - `auth0user`が存在する場合: `/auth/token`を呼び出してAuth0 JWTを取得
  - `auth0user`が存在しない場合: `process.env.NEXT_PUBLIC_API_KEY`を使用
  - エラーハンドリングを実装

**実装例**:
```typescript
import { User } from '@auth0/nextjs-auth0'

type Auth0User = User

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

**受け入れ基準**:
- [ ] `client/src/lib/auth.ts`が作成されている
- [ ] `Auth0User`型がエイリアス定義されている
- [ ] `getAuthToken`関数が実装されている
- [ ] ログイン中に`/auth/token`を呼び出してAuth0 JWTを取得する処理が実装されている
- [ ] 未ログイン時に`process.env.NEXT_PUBLIC_API_KEY`を使用する処理が実装されている
- [ ] エラーハンドリングが適切に実装されている

- _Requirements: 3.1.2, 3.1.3_
- _Design: 3.1.1_

---

### Phase 3: API Clientの修正

#### タスク 3.1: 型定義の追加
**目的**: `Auth0User`型をエイリアス定義し、`getAuthToken`関数をインポートする。また、型定義ファイルの変更に伴い、インポート文を更新する。

**作業内容**:
- `client/src/lib/api.ts`に型定義を追加
- `getAuthToken`関数をインポート
- 型定義ファイルの変更に伴い、インポート文を更新（Phase 0で実施済みの場合は確認のみ）

**実装内容**:
- `@auth0/nextjs-auth0`から`User`をインポート
- `type Auth0User = User`としてエイリアス定義
- `./auth`から`getAuthToken`をインポート
- 型定義のインポート文を確認: `import { DmUser, ... } from '@/types/dm_user'`、`import { DmPost, ... } from '@/types/dm_post'`（Phase 0で実施済み）

**受け入れ基準**:
- [ ] `Auth0User`型がエイリアス定義されている
- [ ] `getAuthToken`関数がインポートされている
- [ ] 型定義のインポートパスが`@/types/dm_user`、`@/types/dm_post`になっている（Phase 0で実施済み）
- [ ] 型定義のインポート文が`DmUser`、`DmPost`になっている（Phase 0で実施済み）

- _Requirements: 3.2.1_
- _Design: 3.2.1_

---

#### タスク 3.2: requestメソッドの修正
**目的**: `request`メソッドに`auth0user`パラメータを追加し、`getAuthToken`を呼び出すように変更する。

**作業内容**:
- `request`メソッドのシグネチャを変更
- `getAuthToken`を呼び出してトークンを取得
- 取得したトークンをAuthorizationヘッダーに付与

**実装内容**:
- `request`メソッドのシグネチャを変更: `request<T>(endpoint: string, options?: RequestInit, auth0user?: Auth0User | undefined): Promise<T>`
- `getAuthToken(auth0user)`を呼び出してトークンを取得
- 取得したトークンをAuthorizationヘッダーに付与
- 既存のエラーハンドリングを維持

**実装例**:
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

**受け入れ基準**:
- [ ] `request`メソッドのシグネチャが`auth0user`パラメータを受け取るように変更されている
- [ ] `getAuthToken(auth0user)`を呼び出してトークンを取得する処理が実装されている
- [ ] 取得したトークンをAuthorizationヘッダーに付与する処理が実装されている
- [ ] 既存のエラーハンドリングが維持されている

- _Requirements: 3.2.2_
- _Design: 3.2.1_

---

#### タスク 3.3: getTodayメソッドの修正
**目的**: `getToday`メソッドを`auth0user`パラメータを受け取るように変更する。

**作業内容**:
- `getToday`メソッドのシグネチャを変更
- `request`メソッドに`auth0user`パラメータを渡す

**実装内容**:
- `getToday`メソッドのシグネチャを変更: `getToday(auth0user?: Auth0User | undefined): Promise<{ date: string }>`
- `request`メソッドに`auth0user`パラメータを渡す
- 既存の`jwt`パラメータを削除

**実装例**:
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

**受け入れ基準**:
- [ ] `getToday`メソッドのシグネチャが`auth0user`パラメータを受け取るように変更されている
- [ ] `request`メソッドに`auth0user`パラメータを渡す処理が実装されている
- [ ] 既存の`jwt`パラメータが削除されている

- _Requirements: 3.2.3_
- _Design: 3.2.1_

---

#### タスク 3.4: User関連メソッドの修正
**目的**: User関連のメソッド名を変更し、`auth0user`パラメータを追加する。

**作業内容**:
- `getUsers` → `getDmUsers`に変更
- `getUser` → `getDmUser`に変更
- `createUser` → `createDmUser`に変更
- `updateUser` → `updateDmUser`に変更
- `deleteUser` → `deleteDmUser`に変更
- 各メソッドに`auth0user`パラメータを追加（オプショナル）
- 戻り値の型を`User` → `DmUser`、`User[]` → `DmUser[]`に変更

**実装内容**:
- メソッド名を変更
- `auth0user?: Auth0User | undefined`パラメータを追加
- `request`メソッドに`auth0user`を渡す
- 戻り値の型を`DmUser`、`DmUser[]`に変更

**実装例**:
```typescript
// 修正前
async getUsers(limit = 20, offset = 0): Promise<User[]> {
  return this.request<User[]>(`/api/dm-users?limit=${limit}&offset=${offset}`)
}

async createUser(data: CreateUserRequest): Promise<User> {
  return this.request<User>('/api/dm-users', {
    method: 'POST',
    body: JSON.stringify(data),
  })
}

async updateUser(id: string, data: UpdateUserRequest): Promise<User> {
  return this.request<User>(`/api/dm-users/${id}`, {
    method: 'PUT',
    body: JSON.stringify(data),
  })
}

// 修正後
async getDmUsers(limit = 20, offset = 0, auth0user?: Auth0User | undefined): Promise<DmUser[]> {
  return this.request<DmUser[]>(`/api/dm-users?limit=${limit}&offset=${offset}`, undefined, auth0user)
}

async createDmUser(data: CreateDmUserRequest, auth0user?: Auth0User | undefined): Promise<DmUser> {
  return this.request<DmUser>('/api/dm-users', {
    method: 'POST',
    body: JSON.stringify(data),
  }, auth0user)
}

async updateDmUser(id: string, data: UpdateDmUserRequest, auth0user?: Auth0User | undefined): Promise<DmUser> {
  return this.request<DmUser>(`/api/dm-users/${id}`, {
    method: 'PUT',
    body: JSON.stringify(data),
  }, auth0user)
}
```

**受け入れ基準**:
- [ ] `getUsers` → `getDmUsers`に変更されている
- [ ] `getUser` → `getDmUser`に変更されている
- [ ] `createUser` → `createDmUser`に変更されている
- [ ] `updateUser` → `updateDmUser`に変更されている
- [ ] `deleteUser` → `deleteDmUser`に変更されている
- [ ] 各メソッドに`auth0user`パラメータが追加されている
- [ ] パラメータの型が`CreateDmUserRequest`、`UpdateDmUserRequest`に変更されている
- [ ] 戻り値の型が`DmUser`、`DmUser[]`に変更されている

- _Requirements: 3.2.4_
- _Design: 3.2.1_

---

#### タスク 3.5: Post関連メソッドの修正
**目的**: Post関連のメソッド名を変更し、`auth0user`パラメータを追加する。

**作業内容**:
- `getPosts` → `getDmPosts`に変更
- `getPost` → `getDmPost`に変更
- `createPost` → `createDmPost`に変更
- `updatePost` → `updateDmPost`に変更
- `deletePost` → `deleteDmPost`に変更
- 各メソッドに`auth0user`パラメータを追加（オプショナル）
- 戻り値の型を`Post` → `DmPost`、`Post[]` → `DmPost[]`に変更

**実装内容**:
- メソッド名を変更
- `auth0user?: Auth0User | undefined`パラメータを追加
- `request`メソッドに`auth0user`を渡す
- 戻り値の型を`DmPost`、`DmPost[]`に変更

**実装例**:
```typescript
// 修正前
async getPosts(limit = 20, offset = 0, userId?: string): Promise<Post[]> {
  // ...
}

async createPost(data: CreatePostRequest): Promise<Post> {
  return this.request<Post>('/api/dm-posts', {
    method: 'POST',
    body: JSON.stringify(data),
  })
}

async updatePost(id: string, userId: string, data: UpdatePostRequest): Promise<Post> {
  return this.request<Post>(`/api/dm-posts/${id}?user_id=${userId}`, {
    method: 'PUT',
    body: JSON.stringify(data),
  })
}

// 修正後
async getDmPosts(limit = 20, offset = 0, userId?: string, auth0user?: Auth0User | undefined): Promise<DmPost[]> {
  // ...
}

async createDmPost(data: CreateDmPostRequest, auth0user?: Auth0User | undefined): Promise<DmPost> {
  return this.request<DmPost>('/api/dm-posts', {
    method: 'POST',
    body: JSON.stringify(data),
  }, auth0user)
}

async updateDmPost(id: string, userId: string, data: UpdateDmPostRequest, auth0user?: Auth0User | undefined): Promise<DmPost> {
  return this.request<DmPost>(`/api/dm-posts/${id}?user_id=${userId}`, {
    method: 'PUT',
    body: JSON.stringify(data),
  }, auth0user)
}
```

**受け入れ基準**:
- [ ] `getPosts` → `getDmPosts`に変更されている
- [ ] `getPost` → `getDmPost`に変更されている
- [ ] `createPost` → `createDmPost`に変更されている
- [ ] `updatePost` → `updateDmPost`に変更されている
- [ ] `deletePost` → `deleteDmPost`に変更されている
- [ ] 各メソッドに`auth0user`パラメータが追加されている
- [ ] パラメータの型が`CreateDmPostRequest`、`UpdateDmPostRequest`に変更されている
- [ ] 戻り値の型が`DmPost`、`DmPost[]`に変更されている

- _Requirements: 3.2.4_
- _Design: 3.2.1_

---

#### タスク 3.6: その他のメソッドの修正
**目的**: その他のメソッド（`getUserPosts`、`sendEmail`、`registerJob`、`downloadUsersCSV`）に`auth0user`パラメータを追加する。

**作業内容**:
- `getUserPosts` → `getDmUserPosts`に変更（必要に応じて）
- `downloadUsersCSV` → `downloadDmUsersCSV`に変更
- 各メソッドに`auth0user`パラメータを追加（オプショナル）

**実装内容**:
- メソッド名を変更（必要に応じて）
- `auth0user?: Auth0User | undefined`パラメータを追加
- `request`メソッドに`auth0user`を渡す（`downloadDmUsersCSV`は直接`fetch`を使用しているため、`getAuthToken`を呼び出す）

**実装例（downloadDmUsersCSV）**:
```typescript
// 修正前
async downloadUsersCSV(): Promise<void> {
  const url = `${this.baseURL}/api/export/dm-users/csv`
  const token = this.apiKey
  // ...
}

// 修正後
async downloadDmUsersCSV(auth0user?: Auth0User | undefined): Promise<void> {
  const url = `${this.baseURL}/api/export/dm-users/csv`
  const token = await getAuthToken(auth0user)
  // ...
}
```

**受け入れ基準**:
- [ ] `getUserPosts` → `getDmUserPosts`に変更されている（必要に応じて）
- [ ] `downloadUsersCSV` → `downloadDmUsersCSV`に変更されている
- [ ] 各メソッドに`auth0user`パラメータが追加されている
- [ ] `downloadDmUsersCSV`で`getAuthToken`を呼び出す処理が実装されている

- _Requirements: 3.2.4_
- _Design: 3.2.1_

---

### Phase 4: コンポーネントの修正

#### タスク 4.1: TodayApiButton.tsxの修正
**目的**: コンポーネントから認証処理とAPI呼び出しの詳細を削除し、`apiClient`を使用するように変更する。

**作業内容**:
- 認証処理（`/auth/token`の呼び出し）を削除
- API呼び出し（`fetch`の直接呼び出し）を削除
- `apiClient.getToday(auth0user || undefined)`を使用
- `useUser()`から取得した変数名を`auth0user`に変更

**実装内容**:
- 認証処理の削除: `/auth/token`の呼び出しを削除
- API呼び出しの削除: `fetch`の直接呼び出しを削除
- `apiClient`のインポートを追加
- `const { user: auth0user, isLoading } = useUser()`に変更
- `apiClient.getToday(auth0user || undefined)`を呼び出す

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

**受け入れ基準**:
- [ ] 認証処理（`/auth/token`の呼び出し）が削除されている
- [ ] API呼び出し（`fetch`の直接呼び出し）が削除されている
- [ ] `apiClient.getToday(auth0user || undefined)`を呼び出すように変更されている
- [ ] `useUser()`から取得した変数名が`auth0user`に変更されている
- [ ] UI表示と状態管理のみに集中している
- [ ] 既存のエラーハンドリングが維持されている

- _Requirements: 3.3.1_
- _Design: 3.3.1_

---

### Phase 5: 既存コンポーネントの修正（メソッド名変更対応）

#### タスク 5.1: dm-users/page.tsxの修正
**目的**: `apiClient`のメソッド名変更と型名変更に対応する。

**作業内容**:
- `apiClient.getUsers()` → `apiClient.getDmUsers()`に変更
- `apiClient.createUser()` → `apiClient.createDmUser()`に変更
- `apiClient.deleteUser()` → `apiClient.deleteDmUser()`に変更
- `apiClient.downloadUsersCSV()` → `apiClient.downloadDmUsersCSV()`に変更
- インポートパスを更新: `import { User } from '@/types/user'` → `import { DmUser } from '@/types/dm_user'`
- 型の使用箇所を更新: `User[]` → `DmUser[]`、`User` → `DmUser`

**実装内容**:
- メソッド呼び出しを変更
- インポートパスを更新
- 型定義を更新（変数宣言、戻り値の型など）

**受け入れ基準**:
- [ ] `apiClient.getUsers()` → `apiClient.getDmUsers()`に変更されている
- [ ] `apiClient.createUser()` → `apiClient.createDmUser()`に変更されている
- [ ] `apiClient.deleteUser()` → `apiClient.deleteDmUser()`に変更されている
- [ ] `apiClient.downloadUsersCSV()` → `apiClient.downloadDmUsersCSV()`に変更されている
- [ ] インポートパスが`@/types/dm_user`に更新されている
- [ ] インポート文が`DmUser`に更新されている
- [ ] 型の使用箇所が`DmUser`、`DmUser[]`に更新されている

- _Requirements: 3.2.4_
- _Design: 3.2.1_

---

#### タスク 5.2: dm-posts/page.tsxの修正
**目的**: `apiClient`のメソッド名変更と型名変更に対応する。

**作業内容**:
- `apiClient.getPosts()` → `apiClient.getDmPosts()`に変更
- `apiClient.getUsers()` → `apiClient.getDmUsers()`に変更
- `apiClient.createPost()` → `apiClient.createDmPost()`に変更
- `apiClient.deletePost()` → `apiClient.deleteDmPost()`に変更
- インポートパスを更新: `import { Post } from '@/types/post'` → `import { DmPost } from '@/types/dm_post'`
- インポートパスを更新: `import { User } from '@/types/user'` → `import { DmUser } from '@/types/dm_user'`
- 型の使用箇所を更新: `Post[]` → `DmPost[]`、`User[]` → `DmUser[]`、`Post` → `DmPost`、`User` → `DmUser`

**実装内容**:
- メソッド呼び出しを変更
- インポートパスを更新
- 型定義を更新（変数宣言、戻り値の型など）

**受け入れ基準**:
- [ ] `apiClient.getPosts()` → `apiClient.getDmPosts()`に変更されている
- [ ] `apiClient.getUsers()` → `apiClient.getDmUsers()`に変更されている
- [ ] `apiClient.createPost()` → `apiClient.createDmPost()`に変更されている
- [ ] `apiClient.deletePost()` → `apiClient.deleteDmPost()`に変更されている
- [ ] インポートパスが`@/types/dm_post`、`@/types/dm_user`に更新されている
- [ ] インポート文が`DmPost`、`DmUser`に更新されている
- [ ] 型の使用箇所が`DmPost`、`DmPost[]`、`DmUser`、`DmUser[]`に更新されている

- _Requirements: 3.2.4_
- _Design: 3.2.1_

---

#### タスク 5.3: dm-user-posts/page.tsxの修正
**目的**: `apiClient`のメソッド名変更に対応する。

**作業内容**:
- `apiClient.getUserPosts()` → `apiClient.getDmUserPosts()`に変更（必要に応じて）

**実装内容**:
- メソッド呼び出しを変更

**受け入れ基準**:
- [ ] `apiClient.getUserPosts()` → `apiClient.getDmUserPosts()`に変更されている（必要に応じて）

- _Requirements: 3.2.4_
- _Design: 3.2.1_

---

### Phase 6: テストの更新

#### タスク 6.1: api.test.tsの更新
**目的**: 新しい認証処理と型名変更に対応したテストを追加・更新する。

**作業内容**:
- `getAuthToken`関数のモックを追加
- `getToday`メソッドのテストを更新
- `auth0user`パラメータが渡された場合のテストを追加
- 既存メソッド名変更に対応したテストを更新
- 型名の変更に対応: `User` → `DmUser`、`Post` → `DmPost`

**実装内容**:
- `getAuthToken`関数のモックを追加
- `getToday`メソッドのテストを更新（`auth0user`パラメータ対応）
- 既存メソッド名変更に対応したテストを更新
- インポート文を更新: `import { User, ... } from '@/types/user'` → `import { DmUser, ... } from '@/types/user'`
- インポート文を更新: `import { Post, ... } from '@/types/post'` → `import { DmPost, ... } from '@/types/post'`
- 型の使用箇所を更新: `User` → `DmUser`、`Post` → `DmPost`、`User[]` → `DmUser[]`、`Post[]` → `DmPost[]`
- ログイン中と未ログイン時の両方のケースをテスト

**受け入れ基準**:
- [ ] `getAuthToken`関数のモックが追加されている
- [ ] `getToday`メソッドのテストが更新されている
- [ ] `auth0user`パラメータが渡された場合のテストが追加されている
- [ ] 既存メソッド名変更に対応したテストが更新されている
- [ ] インポートパスが`@/types/dm_user`、`@/types/dm_post`に更新されている
- [ ] インポート文が`DmUser`、`DmPost`、`CreateDmUserRequest`、`UpdateDmUserRequest`、`CreateDmPostRequest`、`UpdateDmPostRequest`に更新されている
- [ ] 型の使用箇所が`DmUser`、`DmPost`、`DmUser[]`、`DmPost[]`、`CreateDmUserRequest`、`UpdateDmUserRequest`、`CreateDmPostRequest`、`UpdateDmPostRequest`に更新されている
- [ ] ログイン中と未ログイン時の両方のケースがテストされている

- _Requirements: 3.4.1, 6.5_
- _Design: 7.1.2_

---

#### タスク 6.2: TodayApiButtonのテスト更新
**目的**: 新しい実装に対応したテストを更新する。

**作業内容**:
- `apiClient.getToday`のモックを追加
- 認証処理のテストを削除（認証処理は`auth.ts`に移動）
- API呼び出しのテストを更新

**実装内容**:
- `apiClient.getToday`のモックを追加
- 認証処理のテストを削除
- API呼び出しのテストを更新

**受け入れ基準**:
- [ ] `apiClient.getToday`のモックが追加されている
- [ ] 認証処理のテストが削除されている
- [ ] API呼び出しのテストが更新されている

- _Requirements: 3.4.2, 6.5_
- _Design: 7.1.3_

---

### Phase 7: 動作確認

#### タスク 7.1: ローカル環境での動作確認
**目的**: ローカル環境でクライアントアプリが正常に動作することを確認する。

**作業内容**:
- ローカル環境でクライアントアプリを起動
- `TodayApiButton`コンポーネントが正常に動作することを確認
- ログイン中にAuth0 JWTが使用されることを確認
- 未ログイン時にAPI Keyが使用されることを確認
- 全てのAPI呼び出し（`dm-users`、`dm-posts`など）が正常に動作することを確認

**受け入れ基準**:
- [ ] ローカル環境でクライアントアプリが正常に動作する
- [ ] `TodayApiButton`コンポーネントが正常に動作する
- [ ] ログイン中にAuth0 JWTが使用される
- [ ] 未ログイン時にAPI Keyが使用される
- [ ] 全てのAPI呼び出し（`dm-users`、`dm-posts`など）が正常に動作する

- _Requirements: 6.4_
- _Design: 9.2_

---

#### タスク 7.2: テストの実行
**目的**: 全てのテストが通過することを確認する。

**作業内容**:
- 既存のテストを実行
- 新規追加したテストを実行
- 全てのテストが通過することを確認

**受け入れ基準**:
- [ ] 既存のテストが全て通過する
- [ ] 新規追加したテストが全て通過する
- [ ] テストエラーが0件である

- _Requirements: 6.5_
- _Design: 7.1, 7.2_

---

## 実装上の注意事項

### 型定義の注意点
- `Auth0User`型は`@auth0/nextjs-auth0`の`User`型のエイリアス
- アプリケーションの`DmUser`型（`@/types/user`）と混同しない
- 戻り値の`DmUser[]`、`DmPost[]`はアプリケーションの型

### 変数名の注意点
- `useUser()`から取得した変数名は`auth0user`に変更
- 将来的に`user`変数を使う可能性があるため、`auth0user`を使用

### メソッド名の注意点
- `getUsers` → `getDmUsers`、`getPosts` → `getDmPosts`など、本来あるべき名前に変更
- 既存のコンポーネントで使用している場合は、呼び出し側も修正が必要

### Auth0 SDKの利用
- `handleAuth`を使用した標準的な認証エンドポイントを積極的に作成
- 既存の`/auth/token`エンドポイントは維持（アクセストークン取得用）

### 後方互換性
- 後方互換性は不要（既存コードの変更も可）
- 既存のコンポーネントで`apiClient.getUsers`などを使用している場合は、`apiClient.getDmUsers`に変更が必要

## 参考情報

### 関連ドキュメント
- 要件定義書: `.kiro/specs/0062-client-struct/requirements.md`
- 設計書: `.kiro/specs/0062-client-struct/design.md`
- 設計提案: `.kiro/specs/0062-client-struct/client-architecture-proposal.md`
- セッションサマリー: `.kiro/specs/0062-client-struct/session-summary.md`

### 技術スタック
- **フレームワーク**: Next.js 14+ (App Router)
- **言語**: TypeScript 5+
- **認証**: Auth0 Next.js SDK (`@auth0/nextjs-auth0`)
- **テスト**: Jest、React Testing Library
