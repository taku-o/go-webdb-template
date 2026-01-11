# 0062-client-struct セッションサマリー

## 概要
クライアントアプリケーション（Next.js）の設計改善プロジェクト。認証処理とAPI呼び出しの責務を適切に分離し、コードの保守性、テスタビリティ、再利用性を向上させる。

## 会話の流れと決定事項

### 1. 問題の指摘
- `TodayApiButton.tsx`で認証処理がコンポーネント内に直接実装されている
- API呼び出しがコンポーネント内で直接`fetch`を使用している
- 認証トークン取得ロジックが共通化されていない

### 2. 設計提案の作成
- `client-architecture-proposal.md`を作成
- 認証サービスの共通化（`auth.ts`）
- API Clientの改善（`api.ts`）
- コンポーネントの簡素化（`TodayApiButton.tsx`）

### 3. 要件定義書の作成
- `requirements.md`を作成
- 機能要件、非機能要件、制約事項、受け入れ基準を定義

### 4. 基本方針の決定
- **ライブラリの積極的利用**: 使えそうなライブラリ、フレームワークがあるなら積極的に利用する
- **後方互換性**: 不要（既存コードの変更も可）
- **Auth0 SDKの標準的な方法**: `handleAuth`を使用した標準的な認証エンドポイントを積極的に作成

### 5. 型名・変数名の決定
- **型名**: `User`型を`Auth0User`型に変更（他のUser型との衝突を避けるため）
- **変数名**: `user`変数を`auth0user`に変更（将来的に`user`変数を使う可能性があるため）
- **メソッド名**: `getUsers`、`getPosts`を`getDmUsers`、`getDmPosts`に変更（本来あるべき名前）

### 6. 承認状況
- **要件定義書**: 承認済み（2026-01-27）
- **設計書**: 未作成
- **タスク**: 未作成

## 主要な決定事項

### アーキテクチャ
```
React Components (UI表示・状態管理)
    ↓ apiClient.method()
API Client (api.ts) - 統一されたAPI呼び出しインターフェース
    ↓ getAuthToken()
Auth Service (auth.ts) - 認証トークンの取得ロジック
```

### 実装方針
1. **認証サービスの共通化** (`client/src/lib/auth.ts`)
   - `getAuthToken(auth0user: Auth0User | undefined): Promise<string>`
   - ログイン中: `/auth/token`を呼び出してAuth0 JWTを取得
   - 未ログイン: `process.env.NEXT_PUBLIC_API_KEY`を使用

2. **API Clientの改善** (`client/src/lib/api.ts`)
   - `request`メソッドに`auth0user`パラメータを追加
   - `getToday(auth0user?: Auth0User | undefined)`メソッドを修正
   - 既存メソッド（`getDmUsers`、`getDmPosts`など）も必要に応じて修正

3. **コンポーネントの簡素化** (`client/src/components/TodayApiButton.tsx`)
   - 認証処理とAPI呼び出しの詳細を削除
   - `apiClient.getToday(auth0user || undefined)`を使用
   - `useUser()`から取得した変数名を`auth0user`に変更

4. **Auth0 SDKの標準エンドポイント**
   - `client/src/app/api/auth/[auth0]/route.ts`を作成
   - `handleAuth()`を使用して標準的な認証エンドポイントを自動生成

## 型定義

### Auth0User型
```typescript
import { User } from '@auth0/nextjs-auth0'
type Auth0User = User
```

### 主要な関数・メソッドのシグネチャ
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
async getDmUsers(limit = 20, offset = 0, auth0user?: Auth0User | undefined): Promise<User[]>
async getDmPosts(...): Promise<Post[]>
```

## ファイル構成

### 新規作成が必要なファイル
- `client/src/lib/auth.ts`: 認証トークン取得ロジックを集約
- `client/src/app/api/auth/[auth0]/route.ts`: Auth0 SDKの標準的な認証エンドポイント

### 修正が必要なファイル
- `client/src/lib/api.ts`: `request`メソッドと`getToday`メソッドを修正、既存メソッド名を変更
- `client/src/components/TodayApiButton.tsx`: 認証処理とAPI呼び出しを削除し、`apiClient`を使用
- `client/src/lib/__tests__/api.test.ts`: 新しい認証処理に対応したテストを追加

## 次のステップ
1. 設計書（design.md）の作成
2. タスクリスト（tasks.md）の作成
3. 実装の開始

## 参考情報
- 設計提案: `.kiro/specs/0062-client-struct/client-architecture-proposal.md`
- 要件定義書: `.kiro/specs/0062-client-struct/requirements.md`
- 仕様メタデータ: `.kiro/specs/0062-client-struct/spec.json`
