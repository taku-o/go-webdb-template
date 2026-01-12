# 既存clientアプリの機能をclient2アプリに移植する設計書

## 1. 概要

### 1.1 設計の目的
要件定義書に基づき、既存の`client`アプリケーションの機能を`client2`アプリケーションに移植するための詳細設計を定義する。Auth0からNextAuth (Auth.js) v5への移行、カスタムUIからshadcn/uiへの移行、デザイン改善を含む包括的な設計方針を明確にする。

### 1.2 設計の範囲
- 認証システムの移行設計（Auth0 → NextAuth (Auth.js) v5）
- APIクライアントの移植・改修設計
- 型定義の移植設計
- ページコンポーネントの移植設計
- UIコンポーネントの移行設計（カスタム → shadcn/ui）
- テストコードの移植設計
- デザインシステムの統一設計

### 1.3 設計方針
- **既存機能の完全移植**: 既存の`client`アプリのすべての機能を`client2`に移植する
- **認証方式の標準化**: NextAuth (Auth.js) v5の標準的な実装パターンに従う
- **UIコンポーネントの統一**: shadcn/uiコンポーネントを積極的に利用し、統一されたUIを実現する
- **デザインの改善**: 既存のデザインを改善し、モダンなUI/UXを提供する
- **段階的な実装**: Phaseごとに段階的に実装を進める
- **既存APIとの互換性維持**: 既存のバックエンドAPIとの互換性を維持する

## 2. アーキテクチャ設計

### 2.1 全体構成

#### 2.1.1 アーキテクチャ概要

```
┌─────────────────────────────────────────────────────────┐
│              client2/ (移植先アプリケーション)            │
│                                                          │
│  ┌──────────────────────────────────────────────────┐  │
│  │   Next.js 14+ (App Router)                        │  │
│  │   - app/ (ページとAPIルート)                       │  │
│  │   - components/ (Reactコンポーネント)              │  │
│  │   - lib/ (ユーティリティ関数)                      │  │
│  │   - types/ (型定義)                                │  │
│  └──────────────────────────────────────────────────┘  │
│                                                          │
│  ┌──────────────────────────────────────────────────┐  │
│  │   NextAuth (Auth.js) v5                           │  │
│  │   - auth.ts (認証設定)                             │  │
│  │   - app/api/auth/[...nextauth]/route.ts            │  │
│  │   - app/api/auth/token/route.ts                    │  │
│  │   - app/api/auth/profile/route.ts                  │  │
│  │   - lib/auth.ts (認証ヘルパー)                      │  │
│  └──────────────────────────────────────────────────┘  │
│                                                          │
│  ┌──────────────────────────────────────────────────┐  │
│  │   shadcn/ui                                        │  │
│  │   - components/ui/ (UIコンポーネント)               │  │
│  │   - Tailwind CSS (スタイリング)                    │  │
│  └──────────────────────────────────────────────────┘  │
│                                                          │
│  ┌──────────────────────────────────────────────────┐  │
│  │   APIクライアント                                   │  │
│  │   - lib/api.ts (バックエンドAPI呼び出し)            │  │
│  │   - NextAuth対応のトークン取得                      │  │
│  └──────────────────────────────────────────────────┘  │
│                                                          │
│  ┌──────────────────────────────────────────────────┐  │
│  │   テスト                                            │  │
│  │   - e2e/ (Playwright E2Eテスト)                    │  │
│  │   - src/__tests__/ (統合・単体テスト)              │  │
│  └──────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────┐
│           バックエンドAPI (既存、変更なし)                 │
│   - 既存のAPIエンドポイント                              │
│   - 認証トークンによる認証                               │
└─────────────────────────────────────────────────────────┘
```

#### 2.1.2 ディレクトリ構造

**client2ディレクトリ構造（移植後）**:
```
client2/
├── app/                          # Next.js App Router
│   ├── layout.tsx                # ルートレイアウト
│   ├── page.tsx                  # トップページ（移植・改修）
│   ├── api/                      # APIルート
│   │   └── auth/
│   │       ├── [...nextauth]/
│   │       │   └── route.ts     # NextAuth (Auth.js)ルート
│   │       ├── token/
│   │       │   └── route.ts     # トークン取得API（新規）
│   │       └── profile/
│   │           └── route.ts     # プロフィール取得API（新規）
│   ├── dm-users/
│   │   └── page.tsx             # ユーザー管理ページ（新規）
│   ├── dm-posts/
│   │   └── page.tsx             # 投稿管理ページ（新規）
│   ├── dm-user-posts/
│   │   └── page.tsx             # ユーザーと投稿のJOINページ（新規）
│   ├── dm_email/
│   │   └── send/
│   │       └── page.tsx         # メール送信ページ（新規）
│   ├── dm-jobqueue/
│   │   └── page.tsx             # ジョブキューページ（新規）
│   └── dm_movie/
│       └── upload/
│           └── page.tsx         # 動画アップロードページ（新規）
├── components/                   # Reactコンポーネント
│   ├── ui/                       # shadcn/uiコンポーネント
│   │   ├── button.tsx
│   │   ├── card.tsx
│   │   ├── form.tsx
│   │   ├── input.tsx
│   │   ├── select.tsx
│   │   ├── table.tsx
│   │   ├── textarea.tsx
│   │   ├── badge.tsx
│   │   ├── dialog.tsx
│   │   └── alert.tsx
│   ├── layout/
│   │   └── navbar.tsx           # ナビゲーションバー（改修）
│   └── TodayApiButton.tsx        # TodayApiButton（移植・改修）
├── lib/                          # ユーティリティ関数
│   ├── utils.ts                  # ユーティリティ関数（shadcn/ui用）
│   ├── api.ts                    # APIクライアント（移植・改修）
│   └── auth.ts                   # 認証ヘルパー（新規）
├── types/                        # TypeScript型定義
│   ├── dm_post.ts                # 投稿関連の型定義（移植）
│   ├── dm_user.ts                # ユーザー関連の型定義（移植）
│   └── jobqueue.ts               # ジョブキュー関連の型定義（移植）
├── e2e/                          # E2Eテスト（Playwright）
│   ├── auth-flow.spec.ts         # 認証フローテスト（移植・改修）
│   ├── user-flow.spec.ts         # ユーザーフローテスト（移植）
│   ├── post-flow.spec.ts         # 投稿フローテスト（移植）
│   ├── cross-shard.spec.ts       # クロスシャードテスト（移植）
│   ├── email-send.spec.ts        # メール送信テスト（移植）
│   └── csv-download.spec.ts      # CSVダウンロードテスト（移植）
├── src/__tests__/                # 統合・単体テスト
│   ├── integration/              # 統合テスト（移植・改修）
│   ├── components/               # コンポーネント単体テスト（移植）
│   └── lib/                      # ライブラリ単体テスト（移植）
├── auth.ts                       # NextAuth設定（拡張）
├── components.json               # shadcn/ui設定
├── next.config.js                # Next.js設定
├── package.json                  # パッケージ設定（依存関係追加）
├── tsconfig.json                 # TypeScript設定（エイリアス追加）
├── playwright.config.ts          # Playwright設定（新規）
├── .env.example                  # 環境変数テンプレート
└── .env.local                    # ローカル環境変数（gitignore）
```

### 2.2 技術スタック

#### 2.2.1 フレームワーク
- **Next.js 14+ (App Router)**: ページルーティングとAPIルート
- **React 18+**: UIコンポーネント
- **TypeScript 5+**: 型安全性

#### 2.2.2 認証
- **NextAuth (Auth.js) v5**: 認証ライブラリ
  - `auth.ts`: 認証設定
  - `app/api/auth/[...nextauth]/route.ts`: 認証ルート
  - `app/api/auth/token/route.ts`: トークン取得API
  - `app/api/auth/profile/route.ts`: プロフィール取得API
  - `lib/auth.ts`: 認証ヘルパー関数

#### 2.2.3 UIコンポーネント
- **shadcn/ui**: UIコンポーネントライブラリ
  - 既存: `alert-dialog`, `alert`, `button`, `select`, `input`, `form`, `field`, `card`
  - 追加: `table`, `textarea`, `badge`, `dialog`（必要に応じて）

#### 2.2.4 スタイリング
- **Tailwind CSS**: ユーティリティファーストのCSSフレームワーク
- **shadcn/ui CSS変数**: テーマ管理

#### 2.2.5 テスト
- **Playwright**: E2Eテスト
- **Jest**: 単体・統合テスト（必要に応じて）

#### 2.2.6 ファイルアップロード
- **Uppy**: ファイルアップロードライブラリ
  - `@uppy/core`
  - `@uppy/react`
  - `@uppy/tus` (TUSプロトコル)
  - `@uppy/dashboard`

## 3. Phase 1: 基盤整備の設計

### 3.1 型定義の移植設計

#### 3.1.1 移植対象ファイル
- `client/src/types/dm_post.ts` → `client2/types/dm_post.ts`
- `client/src/types/dm_user.ts` → `client2/types/dm_user.ts`
- `client/src/types/jobqueue.ts` → `client2/types/jobqueue.ts`

#### 3.1.2 実装方針
- 既存の型定義をそのまま移植する（変更なし）
- `client2/types/`ディレクトリを作成
- `tsconfig.json`に`@/types`エイリアスを追加（必要に応じて）

#### 3.1.3 ディレクトリ構造
```
client2/
└── types/
    ├── dm_post.ts
    ├── dm_user.ts
    └── jobqueue.ts
```

### 3.2 NextAuth設定の拡張設計

#### 3.2.1 認証プロバイダー設定
- `client2/auth.ts`に認証プロバイダーを追加
- 既存のAuth0設定を参考に、NextAuth (Auth.js) v5のプロバイダーを設定

#### 3.2.2 トークン取得API設計

**エンドポイント**: `app/api/auth/token/route.ts`

**実装方針**:
- NextAuth v5の`auth()`関数を使用してセッションを取得
- セッションからトークンを取得
- 既存の`client/src/app/auth/token/route.ts`と同等の機能を実現

**実装例**:
```typescript
import { NextResponse } from "next/server"
import { auth } from "@/auth"

export async function GET() {
  try {
    const session = await auth()
    if (!session?.accessToken) {
      return NextResponse.json(
        { error: "No access token available" },
        { status: 401 }
      )
    }
    return NextResponse.json({ accessToken: session.accessToken })
  } catch (error) {
    console.error("Failed to get access token:", error)
    return NextResponse.json(
      { error: "Failed to get access token" },
      { status: 500 }
    )
  }
}
```

#### 3.2.3 プロフィール取得API設計

**エンドポイント**: `app/api/auth/profile/route.ts`

**実装方針**:
- NextAuth v5の`auth()`関数を使用してセッションを取得
- セッションからユーザー情報を取得
- 既存の`client/src/app/auth/profile/route.ts`と同等の機能を実現

**実装例**:
```typescript
import { NextResponse } from "next/server"
import { auth } from "@/auth"

export async function GET() {
  try {
    const session = await auth()
    if (!session?.user) {
      return NextResponse.json(
        { error: "Not authenticated" },
        { status: 401 }
      )
    }
    return NextResponse.json({ user: session.user })
  } catch (error) {
    console.error("Failed to get profile:", error)
    return NextResponse.json(
      { error: "Failed to get profile" },
      { status: 500 }
    )
  }
}
```

### 3.3 APIクライアントの移植・改修設計

#### 3.3.1 移植対象ファイル
- `client/src/lib/api.ts` → `client2/lib/api.ts`

#### 3.3.2 改修内容

**認証トークン取得の変更**:
- Auth0依存をNextAuth (Auth.js) v5に置き換え
- `getAuthToken`関数をNextAuth対応に変更

**実装方針**:
```typescript
// client2/lib/api.ts
import { getAuthToken } from './auth'

class ApiClient {
  private async request<T>(
    endpoint: string,
    options?: RequestInit
  ): Promise<T> {
    const url = `${this.baseURL}${endpoint}`
    
    // NextAuth対応のトークン取得
    const token = await getAuthToken()
    
    const headers = {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`,
      ...options?.headers,
    }
    
    // ... 既存の実装
  }
}
```

#### 3.3.3 環境変数
- `NEXT_PUBLIC_API_BASE_URL`: バックエンドAPIのベースURL
- `NEXT_PUBLIC_API_KEY`: APIキー（認証なしの場合）

### 3.4 認証ヘルパーの実装設計

#### 3.4.1 ファイル作成
- `client2/lib/auth.ts`を新規作成

#### 3.4.2 実装内容

**NextAuth v5のラッパー関数**:
```typescript
// client2/lib/auth.ts
import { auth, signIn, signOut } from "@/auth"

// サーバー側での認証状態取得
export async function getServerSession() {
  return await auth()
}

// クライアント側での認証状態取得用フック（ダミー実装で処理を差し込む場所を用意）
export function useAuth() {
  // TODO: クライアント側での認証状態取得を実装
  // 現時点ではダミー実装で処理を差し込む場所を用意
  return {
    user: null,
    isLoading: false,
    signIn: async () => {},
    signOut: async () => {},
  }
}
```

**トークン取得関数**:
```typescript
export async function getAuthToken(): Promise<string> {
  const session = await auth()
  if (session?.accessToken) {
    return session.accessToken
  }
  
  // 認証なしの場合はAPIキーを使用
  const apiKey = process.env.NEXT_PUBLIC_API_KEY
  if (!apiKey) {
    throw new Error('NEXT_PUBLIC_API_KEY is not set')
  }
  return apiKey
}
```

## 4. Phase 2: 共通コンポーネントとレイアウトの設計

### 4.1 レイアウトコンポーネントの改修設計

#### 4.1.1 ナビゲーションバーの改修
- `client2/components/layout/navbar.tsx`を改修

**実装内容**:
- 認証状態の表示（NextAuth対応）
- ログイン/ログアウトボタンの実装
- ユーザー情報表示の実装
- shadcn/uiの`button`コンポーネントを使用

**実装例**:
```typescript
// client2/components/layout/navbar.tsx
import { auth, signIn, signOut } from "@/auth"
import { Button } from "@/components/ui/button"

export async function Navbar() {
  const session = await auth()
  
  return (
    <nav>
      {session?.user ? (
        <>
          <span>{session.user.name}</span>
          <form action={async () => {
            "use server"
            await signOut()
          }}>
            <Button type="submit">ログアウト</Button>
          </form>
        </>
      ) : (
        <form action={async () => {
          "use server"
          await signIn()
        }}>
          <Button type="submit">ログイン</Button>
        </form>
      )}
    </nav>
  )
}
```

### 4.2 shadcn/uiコンポーネントの追加インストール設計

#### 4.2.1 追加インストール対象
- `table`: 一覧表示用
- `textarea`: フォーム用（投稿の本文入力など）
- `badge`: ステータス表示用
- `dialog`: モーダル用（既にインストール済みか確認）

#### 4.2.2 インストール方法
```bash
npx shadcn-ui@latest add table
npx shadcn-ui@latest add textarea
npx shadcn-ui@latest add badge
npx shadcn-ui@latest add dialog  # 必要に応じて
```

### 4.3 共通UIコンポーネントの作成設計

#### 4.3.1 エラー表示コンポーネント
- shadcn/uiの`alert`コンポーネントを使用
- エラーメッセージを表示する共通コンポーネント

#### 4.3.2 ローディング表示コンポーネント
- ローディング状態を表示する共通コンポーネント
- shadcn/uiのスタイルに統一

#### 4.3.3 フォームコンポーネント
- shadcn/uiの`form`コンポーネントを使用
- react-hook-formと統合

## 5. Phase 3: ページ移植の設計

### 5.1 トップページの移植・改修設計

#### 5.1.1 移植対象
- `client/src/app/page.tsx` → `client2/app/page.tsx`（既存を置き換え）

#### 5.1.2 実装内容
- 機能一覧の表示（shadcn/uiの`card`コンポーネントを使用）
- 認証状態の表示（NextAuth対応）
- TodayApiButtonコンポーネントの統合
- デザイン改善（モダンなUI）

**移植しない内容**:
- サンプル画像ファイルの参照例（`client/src/app/page.tsx`の125-142行目）は移植しない

### 5.2 TodayApiButtonコンポーネントの移植設計

#### 5.2.1 移植対象
- `client/src/components/TodayApiButton.tsx` → `client2/components/TodayApiButton.tsx`

#### 5.2.2 改修内容
- Auth0依存をNextAuth対応に変更
- shadcn/uiの`button`コンポーネントを使用
- エラー表示にshadcn/uiの`alert`を使用

### 5.3 ユーザー管理ページの移植設計

#### 5.3.1 移植対象
- `client/src/app/dm-users/page.tsx` → `client2/app/dm-users/page.tsx`（新規作成）

#### 5.3.2 実装内容
- CRUD機能の実装（作成、一覧表示、削除）
- CSVダウンロード機能の実装
- shadcn/uiコンポーネントを使用（`form`, `input`, `select`, `button`, `table`, `card`）
- デザイン改善（モダンなUI、レスポンシブ対応）

### 5.4 投稿管理ページの移植設計

#### 5.4.1 移植対象
- `client/src/app/dm-posts/page.tsx` → `client2/app/dm-posts/page.tsx`（新規作成）

#### 5.4.2 実装内容
- CRUD機能の実装（作成、一覧表示、削除）
- ユーザー選択の実装（`select`コンポーネントを使用）
- shadcn/uiコンポーネントを使用（`form`, `input`, `textarea`, `select`, `button`, `card`）
- デザイン改善（モダンなUI、レスポンシブ対応）

### 5.5 ユーザーと投稿のJOINページの移植設計

#### 5.5.1 移植対象
- `client/src/app/dm-user-posts/page.tsx` → `client2/app/dm-user-posts/page.tsx`（新規作成）

#### 5.5.2 実装内容
- クロスシャードクエリの説明を保持
- shadcn/uiの`card`を使用して投稿を表示
- デザイン改善（モダンなUI、レスポンシブ対応）

### 5.6 メール送信ページの移植設計

#### 5.6.1 移植対象
- `client/src/app/dm_email/send/page.tsx` → `client2/app/dm_email/send/page.tsx`（新規作成）

#### 5.6.2 実装内容
- メール送信フォームの実装
- shadcn/uiの`form`, `input`, `button`を使用
- 成功/エラーメッセージの表示（shadcn/uiの`alert`を使用）
- デザイン改善（モダンなUI、レスポンシブ対応）

### 5.7 ジョブキューページの移植設計

#### 5.7.1 移植対象
- `client/src/app/dm-jobqueue/page.tsx` → `client2/app/dm-jobqueue/page.tsx`（新規作成）

#### 5.7.2 実装内容
- ジョブキューの表示と操作
- shadcn/uiコンポーネントを使用（`form`, `input`, `button`, `alert`）
- デザイン改善（モダンなUI、レスポンシブ対応）

### 5.8 動画アップロードページの移植設計

#### 5.8.1 移植対象
- `client/src/app/dm_movie/upload/page.tsx` → `client2/app/dm_movie/upload/page.tsx`（新規作成）

#### 5.8.2 実装内容
- Uppyライブラリのインストール（`@uppy/core`, `@uppy/react`, `@uppy/tus`, `@uppy/dashboard`）
- TUSプロトコル対応の実装
- 認証トークンの取得方法をNextAuth対応に変更
- デザイン改善（モダンなUI、レスポンシブ対応）

**Uppy設定例**:
```typescript
import Uppy from '@uppy/core'
import { Dashboard } from '@uppy/react'
import Tus from '@uppy/tus'
import { getAuthToken } from '@/lib/auth'

const uppy = new Uppy()
  .use(Tus, {
    endpoint: '/api/upload',
    headers: async () => {
      const token = await getAuthToken()
      return {
        Authorization: `Bearer ${token}`,
      }
    },
  })
```

## 6. Phase 4: デザイン改善の設計

### 6.1 デザインシステムの統一設計

#### 6.1.1 カラーパレット
- shadcn/uiのデフォルトテーマをベース
- 全ページで統一されたカラーパレットを使用

#### 6.1.2 タイポグラフィ
- フォントサイズ、行間、フォントファミリーを統一

#### 6.1.3 スペーシング
- 統一されたスペーシングスケールを使用
- Tailwind CSSのスペーシングユーティリティを活用

#### 6.1.4 コンポーネントスタイル
- ボタン、フォーム、カードなどのスタイルを統一
- shadcn/uiのデフォルトスタイルをベースに

### 6.2 レスポンシブデザインの改善設計

#### 6.2.1 ブレークポイント
- モバイル: 640px以下
- タブレット: 641px - 1024px
- デスクトップ: 1025px以上

#### 6.2.2 実装方針
- Tailwind CSSのレスポンシブユーティリティを活用
- モバイルファーストのアプローチ
- 各ページでレスポンシブ対応を確認・改善

### 6.3 アクセシビリティの向上設計

#### 6.3.1 キーボードナビゲーション
- すべてのインタラクティブ要素がキーボードで操作可能
- フォーカス表示を明確に

#### 6.3.2 スクリーンリーダー対応
- 適切なARIA属性の使用
- shadcn/uiのアクセシビリティ機能を活用

## 7. Phase 5: テストの移植設計

### 7.1 テスト環境のセットアップ設計

#### 7.1.1 Playwright設定
- `playwright.config.ts`を作成
- テスト用の環境変数設定

#### 7.1.2 Jest設定（必要に応じて）
- `jest.config.js`を作成（必要に応じて）
- テスト用の環境変数設定

### 7.2 E2Eテストの移植・改修設計

#### 7.2.1 移植対象
- `client/e2e/auth-flow.spec.ts` → NextAuth対応に変更
- `client/e2e/user-flow.spec.ts` → 移植
- `client/e2e/post-flow.spec.ts` → 移植
- `client/e2e/cross-shard.spec.ts` → 移植
- `client/e2e/email-send.spec.ts` → 移植
- `client/e2e/csv-download.spec.ts` → 移植

#### 7.2.2 改修内容
- 認証フローのテストをNextAuth対応に変更
- 既存のテストロジックを維持しつつ、NextAuthの認証方法に合わせて調整

### 7.3 統合テストの移植設計

#### 7.3.1 移植対象
- `client/src/__tests__/integration/`のテストを移植

#### 7.3.2 改修内容
- NextAuth対応に変更
- テストが正常に動作することを確認

### 7.4 単体テストの移植設計

#### 7.4.1 移植対象
- `client/src/components/__tests__/`のテストを移植
- `client/src/lib/__tests__/`のテストを移植

#### 7.4.2 実装方針
- 既存のテストロジックを維持
- テストが正常に動作することを確認

## 8. Phase 6: 最終確認とドキュメントの設計

### 8.1 動作確認の設計

#### 8.1.1 確認項目
- すべてのページの動作確認
- 認証フローの確認（ログイン、ログアウト、トークン取得）
- API呼び出しの確認（すべてのエンドポイント）
- エラーハンドリングの確認

### 8.2 パフォーマンス確認の設計

#### 8.2.1 確認項目
- ページ読み込み速度の確認
- API呼び出しの最適化（必要に応じて）

### 8.3 ドキュメント更新の設計

#### 8.3.1 作成・更新対象
- `docs/Temp-Client2.md`の作成・更新
- 環境変数のドキュメント化
- セットアップ手順のドキュメント化
- 機能説明のドキュメント化

## 9. データフロー設計

### 9.1 認証フロー

```
ユーザー
  │
  ├─→ ログイン
  │     │
  │     └─→ NextAuth (Auth.js) v5
  │           │
  │           └─→ セッション作成
  │                 │
  │                 └─→ トークン取得
  │                       │
  │                       └─→ API呼び出し
  │
  └─→ ログアウト
        │
        └─→ セッション削除
```

### 9.2 API呼び出しフロー

```
ページコンポーネント
  │
  └─→ lib/api.ts (ApiClient)
        │
        ├─→ lib/auth.ts (getAuthToken)
        │     │
        │     └─→ auth.ts (NextAuth)
        │           │
        │           └─→ セッション取得
        │                 │
        │                 └─→ トークン返却
        │
        └─→ バックエンドAPI
              │
              └─→ レスポンス返却
```

## 10. セキュリティ設計

### 10.1 認証セキュリティ
- NextAuth (Auth.js) v5のセキュリティベストプラクティスに従う
- セッション管理を適切に実装
- トークンの安全な管理

### 10.2 APIセキュリティ
- 認証が必要なAPI呼び出しで適切にトークンを使用
- エラーハンドリングで機密情報を漏洩しない

## 11. パフォーマンス設計

### 11.1 ページ読み込み速度
- Next.js 14のApp Routerの最適化機能を活用
- 必要に応じてServer ComponentsとClient Componentsを適切に使い分け

### 11.2 API呼び出しの最適化
- 必要に応じてキャッシュを活用
- 不要なAPI呼び出しを削減

## 12. 実装順序

### 12.1 Phase 1: 基盤整備
1. 型定義の移植
2. NextAuth設定の拡張
3. APIクライアントの移植・改修
4. 認証ヘルパーの実装

### 12.2 Phase 2: 共通コンポーネントとレイアウト
1. レイアウトコンポーネントの改修
2. shadcn/uiコンポーネントの追加インストール
3. 共通UIコンポーネントの作成

### 12.3 Phase 3: ページ移植
1. トップページの移植・改修
2. TodayApiButtonコンポーネントの移植
3. ユーザー管理ページの移植
4. 投稿管理ページの移植
5. ユーザーと投稿のJOINページの移植
6. メール送信ページの移植
7. ジョブキューページの移植
8. 動画アップロードページの移植

### 12.4 Phase 4: デザイン改善
- Phase 3と並行して実施

### 12.5 Phase 5: テストの移植
1. テスト環境のセットアップ
2. E2Eテストの移植・改修
3. 統合テストの移植
4. 単体テストの移植

### 12.6 Phase 6: 最終確認とドキュメント
1. 動作確認
2. パフォーマンス確認
3. ドキュメント更新

## 13. 依存関係

### 13.1 Phase間の依存関係
- Phase 2はPhase 1完了後に開始
- Phase 3はPhase 1とPhase 2完了後に開始
- Phase 4はPhase 3と並行して実施可能
- Phase 5はPhase 3完了後に開始
- Phase 6はすべてのPhase完了後に実施

## 14. リスクと対策

### 14.1 認証移行のリスク
- **リスク**: Auth0からNextAuthへの移行で認証フローが異なる
- **対策**: 既存のAuth0設定を参考に、NextAuthの標準的な実装パターンに従う

### 14.2 UI移行のリスク
- **リスク**: カスタムUIからshadcn/uiへの移行で見た目が変わる
- **対策**: 段階的に移行し、各ページでデザインを改善

### 14.3 テスト移行のリスク
- **リスク**: テストがNextAuth対応で動作しない
- **対策**: テストを段階的に移植し、各テストで動作確認を行う

## 15. 参考情報

### 15.1 関連ドキュメント
- `.kiro/specs/0064-client2-impl/requirements.md`: 要件定義書
- `.kiro/specs/0064-client2-impl/Client2-Setup-Summary.md`: client2アプリのセットアップ作業まとめ
- `.kiro/specs/0063-client2/`: client2アプリ作成の要件定義書・設計書

### 15.2 外部リソース
- **NextAuth (Auth.js)**: https://authjs.dev/
- **shadcn/ui**: https://ui.shadcn.com/
- **Next.js**: https://nextjs.org/
- **Playwright**: https://playwright.dev/
- **Uppy**: https://uppy.io/
