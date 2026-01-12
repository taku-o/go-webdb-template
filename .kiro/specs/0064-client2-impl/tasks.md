# 既存clientアプリの機能をclient2アプリに移植する実装タスク一覧

## 概要
既存の`client`アプリケーションの機能を`client2`アプリケーションに移植するためのタスク一覧。Auth0からNextAuth (Auth.js) v5への移行、カスタムUIからshadcn/uiへの移行、デザイン改善を含む包括的な実装タスクを定義する。要件定義書と設計書に基づき、実装可能な粒度でタスクを分解。

## 実装フェーズ

### Phase 1: 基盤整備（認証・API・型定義）

#### タスク 1.1: 型定義の移植
**目的**: 既存の型定義をclient2に移植する。

**作業内容**:
- `client/src/types/`ディレクトリの型定義ファイルを`client2/types/`に移植
- `tsconfig.json`に`@/types`エイリアスを追加（必要に応じて）

**実装内容**:
1. `client2/types/`ディレクトリを作成
2. 以下のファイルを移植:
   - `client/src/types/dm_post.ts` → `client2/types/dm_post.ts`
   - `client/src/types/dm_user.ts` → `client2/types/dm_user.ts`
   - `client/src/types/jobqueue.ts` → `client2/types/jobqueue.ts`
3. `tsconfig.json`に`@/types`エイリアスを追加（必要に応じて）

**受け入れ基準**:
- [ ] `client2/types/`ディレクトリが作成されている
- [ ] `dm_post.ts`, `dm_user.ts`, `jobqueue.ts`が移植されている
- [ ] 型定義ファイルにエラーがない（TypeScript型チェックが通る）
- [ ] `tsconfig.json`に`@/types`エイリアスが追加されている（必要に応じて）

- _Requirements: 3.1.1_
- _Design: 3.1_

---

#### タスク 1.2: NextAuth設定の拡張
**目的**: NextAuth (Auth.js) v5の設定を拡張し、既存のAuth0機能と同等の機能を実現する。

**作業内容**:
- `client2/auth.ts`に認証プロバイダー設定を追加
- トークン取得用のAPI Route作成
- プロフィール取得用のAPI Route作成

**実装内容**:
1. `client2/auth.ts`を確認・拡張（既存のAuth0設定を参考）
2. `app/api/auth/token/route.ts`を作成:
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
3. `app/api/auth/profile/route.ts`を作成:
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

**受け入れ基準**:
- [ ] `client2/auth.ts`に認証プロバイダー設定が追加されている
- [ ] `app/api/auth/token/route.ts`が作成され、トークン取得が動作する
- [ ] `app/api/auth/profile/route.ts`が作成され、プロフィール取得が動作する
- [ ] 既存の`client/src/app/auth/token/route.ts`と同等の機能を実現している
- [ ] 既存の`client/src/app/auth/profile/route.ts`と同等の機能を実現している

- _Requirements: 3.1.2_
- _Design: 3.2_

---

#### タスク 1.3: APIクライアントの移植・改修
**目的**: 既存のAPIクライアントを移植し、NextAuth対応に変更する。

**作業内容**:
- `client/src/lib/api.ts`を`client2/lib/api.ts`に移植
- Auth0依存をNextAuth (Auth.js) v5に置き換え
- `getAuthToken`関数をNextAuth対応に変更
- 環境変数の確認

**実装内容**:
1. `client/src/lib/api.ts`を`client2/lib/api.ts`にコピー
2. Auth0依存を削除:
   - `import { User } from '@auth0/nextjs-auth0/types'`を削除
   - `Auth0User`型を削除
3. `getAuthToken`関数の呼び出しをNextAuth対応に変更:
   - `lib/auth.ts`の`getAuthToken`関数を使用
4. 環境変数の確認:
   - `NEXT_PUBLIC_API_BASE_URL`
   - `NEXT_PUBLIC_API_KEY`
5. すべてのAPIメソッドが正常に動作することを確認

**受け入れ基準**:
- [ ] `client2/lib/api.ts`が作成されている
- [ ] Auth0依存が削除されている
- [ ] `getAuthToken`関数がNextAuth対応に変更されている
- [ ] 環境変数が適切に設定されている
- [ ] すべてのAPIメソッドが正常に動作する

- _Requirements: 3.1.3_
- _Design: 3.3_

---

#### タスク 1.4: 認証ヘルパーの実装
**目的**: NextAuth v5の認証機能をラップしたヘルパー関数を作成する。

**作業内容**:
- `client2/lib/auth.ts`を作成
- NextAuth v5の`auth()`, `signIn()`, `signOut()`をラップ
- クライアント側での認証状態取得用フック（ダミー実装で処理を差し込む場所を用意）

**実装内容**:
1. `client2/lib/auth.ts`を作成:
   ```typescript
   import { auth, signIn, signOut } from "@/auth"

   // サーバー側での認証状態取得
   export async function getServerSession() {
     return await auth()
   }

   // トークン取得関数
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

**受け入れ基準**:
- [ ] `client2/lib/auth.ts`が作成されている
- [ ] `getServerSession`関数が実装されている
- [ ] `getAuthToken`関数が実装されている
- [ ] `useAuth`フックがダミー実装で用意されている（処理を差し込む場所を用意）

- _Requirements: 3.1.4_
- _Design: 3.4_

---

### Phase 2: 共通コンポーネントとレイアウト

#### タスク 2.1: レイアウトコンポーネントの改修
**目的**: 認証状態を表示し、ログイン/ログアウト機能を実装する。

**作業内容**:
- `client2/components/layout/navbar.tsx`に認証状態表示を追加
- ログイン/ログアウトボタンの実装
- ユーザー情報表示の実装

**実装内容**:
1. `client2/components/layout/navbar.tsx`を確認・改修
2. NextAuthの`auth()`関数を使用して認証状態を取得
3. ログイン/ログアウトボタンを実装（shadcn/uiの`button`コンポーネントを使用）
4. ユーザー情報を表示

**受け入れ基準**:
- [ ] `navbar.tsx`に認証状態表示が追加されている
- [ ] ログイン/ログアウトボタンが実装されている
- [ ] ユーザー情報が表示される
- [ ] shadcn/uiの`button`コンポーネントが使用されている

- _Requirements: 3.2.1_
- _Design: 4.1_

---

#### タスク 2.2: shadcn/uiコンポーネントの追加インストール
**目的**: ページ移植に必要なshadcn/uiコンポーネントを追加インストールする。

**作業内容**:
- 必要なshadcn/uiコンポーネントをインストール
- 既存のコンポーネントを確認

**実装内容**:
```bash
npx shadcn-ui@latest add table
npx shadcn-ui@latest add textarea
npx shadcn-ui@latest add badge
npx shadcn-ui@latest add dialog  # 既にインストール済みか確認
```

**注意**: タスク実行時に必要と判断した場合は追加のコンポーネントをインストール（タスク実行者の判断）

**受け入れ基準**:
- [ ] `table`コンポーネントがインストールされている
- [ ] `textarea`コンポーネントがインストールされている
- [ ] `badge`コンポーネントがインストールされている
- [ ] `dialog`コンポーネントがインストールされている（既にインストール済みか確認）

- _Requirements: 3.2.2_
- _Design: 4.2_

---

#### タスク 2.3: 共通UIコンポーネントの作成
**目的**: ページ間で共通して使用するUIコンポーネントを作成する。

**作業内容**:
- エラー表示コンポーネントの作成
- ローディング表示コンポーネントの作成
- フォームコンポーネントの作成

**実装内容**:
1. エラー表示コンポーネント（shadcn/uiの`alert`を使用）
2. ローディング表示コンポーネント
3. フォームコンポーネント（shadcn/uiの`form`を使用、react-hook-formと統合）

**受け入れ基準**:
- [ ] エラー表示コンポーネントが作成されている
- [ ] ローディング表示コンポーネントが作成されている
- [ ] フォームコンポーネントが作成されている
- [ ] shadcn/uiコンポーネントが使用されている

- _Requirements: 3.2.3_
- _Design: 4.3_

---

### Phase 3: ページ移植（シンプルなものから順に）

#### タスク 3.1: トップページの移植・改修
**目的**: 既存のトップページ機能を移植し、shadcn/uiでデザインを改善する。

**作業内容**:
- `client2/app/page.tsx`を既存のトップページ機能に置き換え
- shadcn/uiの`card`コンポーネントを使用
- 認証状態の表示（NextAuth対応）
- TodayApiButtonコンポーネントの統合

**実装内容**:
1. `client/src/app/page.tsx`を確認
2. `client2/app/page.tsx`を既存のトップページ機能に置き換え
3. 機能一覧をshadcn/uiの`card`コンポーネントで表示
4. 認証状態の表示をNextAuth対応に変更
5. TodayApiButtonコンポーネントを統合
6. デザイン改善（モダンなUI）

**注意**: サンプル画像ファイルの参照例（`client/src/app/page.tsx`の125-142行目）は移植しない

**受け入れ基準**:
- [ ] `client2/app/page.tsx`が既存のトップページ機能に置き換えられている
- [ ] 機能一覧がshadcn/uiの`card`コンポーネントで表示されている
- [ ] 認証状態がNextAuth対応で表示されている
- [ ] TodayApiButtonコンポーネントが統合されている
- [ ] デザインが改善されている（モダンなUI）
- [ ] サンプル画像ファイルの参照例が移植されていない

- _Requirements: 3.3.1_
- _Design: 5.1_

---

#### タスク 3.2: TodayApiButtonコンポーネントの移植
**目的**: 既存のTodayApiButtonコンポーネントを移植し、NextAuth対応に変更する。

**作業内容**:
- `client/src/components/TodayApiButton.tsx`を`client2/components/TodayApiButton.tsx`に移植
- Auth0依存をNextAuth対応に変更
- shadcn/uiコンポーネントを使用

**実装内容**:
1. `client/src/components/TodayApiButton.tsx`を確認
2. `client2/components/TodayApiButton.tsx`に移植
3. Auth0依存を削除し、NextAuth対応に変更
4. shadcn/uiの`button`コンポーネントを使用
5. エラー表示にshadcn/uiの`alert`を使用

**受け入れ基準**:
- [ ] `client2/components/TodayApiButton.tsx`が作成されている
- [ ] Auth0依存が削除されている
- [ ] NextAuth対応に変更されている
- [ ] shadcn/uiの`button`コンポーネントが使用されている
- [ ] エラー表示にshadcn/uiの`alert`が使用されている

- _Requirements: 3.3.2_
- _Design: 5.2_

---

#### タスク 3.3: ユーザー管理ページの移植
**目的**: 既存のユーザー管理ページを移植し、shadcn/uiでデザインを改善する。

**作業内容**:
- `client2/app/dm-users/page.tsx`を作成
- CRUD機能の実装
- CSVダウンロード機能の実装
- shadcn/uiコンポーネントを使用

**実装内容**:
1. `client/src/app/dm-users/page.tsx`を確認
2. `client2/app/dm-users/page.tsx`を作成
3. CRUD機能を実装（作成、一覧表示、削除）
4. CSVダウンロード機能を実装
5. shadcn/uiコンポーネントを使用（`form`, `input`, `select`, `button`, `table`, `card`）
6. デザイン改善（モダンなUI、レスポンシブ対応）

**受け入れ基準**:
- [ ] `client2/app/dm-users/page.tsx`が作成されている
- [ ] CRUD機能が動作する（作成、一覧表示、削除）
- [ ] CSVダウンロード機能が動作する
- [ ] shadcn/uiコンポーネントが使用されている
- [ ] デザインが改善されている（モダンなUI、レスポンシブ対応）

- _Requirements: 3.3.3_
- _Design: 5.3_

---

#### タスク 3.4: 投稿管理ページの移植
**目的**: 既存の投稿管理ページを移植し、shadcn/uiでデザインを改善する。

**作業内容**:
- `client2/app/dm-posts/page.tsx`を作成
- CRUD機能の実装
- ユーザー選択の実装
- shadcn/uiコンポーネントを使用

**実装内容**:
1. `client/src/app/dm-posts/page.tsx`を確認
2. `client2/app/dm-posts/page.tsx`を作成
3. CRUD機能を実装（作成、一覧表示、削除）
4. ユーザー選択を実装（`select`コンポーネントを使用）
5. shadcn/uiコンポーネントを使用（`form`, `input`, `textarea`, `select`, `button`, `card`）
6. デザイン改善（モダンなUI、レスポンシブ対応）

**受け入れ基準**:
- [ ] `client2/app/dm-posts/page.tsx`が作成されている
- [ ] CRUD機能が動作する（作成、一覧表示、削除）
- [ ] ユーザー選択が動作する
- [ ] shadcn/uiコンポーネントが使用されている
- [ ] デザインが改善されている（モダンなUI、レスポンシブ対応）

- _Requirements: 3.3.4_
- _Design: 5.4_

---

#### タスク 3.5: ユーザーと投稿のJOINページの移植
**目的**: 既存のJOINページを移植し、shadcn/uiでデザインを改善する。

**作業内容**:
- `client2/app/dm-user-posts/page.tsx`を作成
- クロスシャードクエリの説明を保持
- shadcn/uiの`card`を使用して投稿を表示

**実装内容**:
1. `client/src/app/dm-user-posts/page.tsx`を確認
2. `client2/app/dm-user-posts/page.tsx`を作成
3. クロスシャードクエリの説明を保持
4. shadcn/uiの`card`を使用して投稿を表示
5. デザイン改善（モダンなUI、レスポンシブ対応）

**受け入れ基準**:
- [ ] `client2/app/dm-user-posts/page.tsx`が作成されている
- [ ] クロスシャードクエリの説明が保持されている
- [ ] shadcn/uiの`card`が使用されている
- [ ] デザインが改善されている（モダンなUI、レスポンシブ対応）

- _Requirements: 3.3.5_
- _Design: 5.5_

---

#### タスク 3.6: メール送信ページの移植
**目的**: 既存のメール送信ページを移植し、shadcn/uiでデザインを改善する。

**作業内容**:
- `client2/app/dm_email/send/page.tsx`を作成
- メール送信フォームの実装
- 成功/エラーメッセージの表示

**実装内容**:
1. `client/src/app/dm_email/send/page.tsx`を確認
2. `client2/app/dm_email/send/page.tsx`を作成
3. メール送信フォームを実装
4. shadcn/uiの`form`, `input`, `button`を使用
5. 成功/エラーメッセージの表示（shadcn/uiの`alert`を使用）
6. デザイン改善（モダンなUI、レスポンシブ対応）

**受け入れ基準**:
- [ ] `client2/app/dm_email/send/page.tsx`が作成されている
- [ ] メール送信フォームが動作する
- [ ] shadcn/uiコンポーネントが使用されている
- [ ] 成功/エラーメッセージが表示される
- [ ] デザインが改善されている（モダンなUI、レスポンシブ対応）

- _Requirements: 3.3.6_
- _Design: 5.6_

---

#### タスク 3.7: ジョブキューページの移植
**目的**: 既存のジョブキューページを移植し、shadcn/uiでデザインを改善する。

**作業内容**:
- `client2/app/dm-jobqueue/page.tsx`を作成
- ジョブキューの表示と操作
- shadcn/uiコンポーネントを使用

**実装内容**:
1. `client/src/app/dm-jobqueue/page.tsx`を確認
2. `client2/app/dm-jobqueue/page.tsx`を作成
3. ジョブキューの表示と操作を実装
4. shadcn/uiコンポーネントを使用（`form`, `input`, `button`, `alert`）
5. デザイン改善（モダンなUI、レスポンシブ対応）

**受け入れ基準**:
- [ ] `client2/app/dm-jobqueue/page.tsx`が作成されている
- [ ] ジョブキューの表示と操作が動作する
- [ ] shadcn/uiコンポーネントが使用されている
- [ ] デザインが改善されている（モダンなUI、レスポンシブ対応）

- _Requirements: 3.3.7_
- _Design: 5.7_

---

#### タスク 3.8: 動画アップロードページの移植
**目的**: 既存の動画アップロードページを移植し、NextAuth対応に変更する。

**作業内容**:
- `client2/app/dm_movie/upload/page.tsx`を作成
- Uppyライブラリのインストール
- TUSプロトコル対応の実装
- 認証トークンの取得方法をNextAuth対応に変更

**実装内容**:
1. `client/src/app/dm_movie/upload/page.tsx`を確認
2. Uppyライブラリをインストール:
   ```bash
   npm install @uppy/core @uppy/react @uppy/tus @uppy/dashboard --legacy-peer-deps
   ```
3. `client2/app/dm_movie/upload/page.tsx`を作成
4. TUSプロトコル対応を実装
5. 認証トークンの取得方法をNextAuth対応に変更（`lib/auth.ts`の`getAuthToken`を使用）
6. デザイン改善（モダンなUI、レスポンシブ対応）

**受け入れ基準**:
- [ ] `client2/app/dm_movie/upload/page.tsx`が作成されている
- [ ] Uppyライブラリがインストールされている
- [ ] TUSプロトコル対応が実装されている
- [ ] 認証トークンの取得がNextAuth対応になっている
- [ ] デザインが改善されている（モダンなUI、レスポンシブ対応）

- _Requirements: 3.3.8_
- _Design: 5.8_

---

### Phase 4: デザイン改善

#### タスク 4.1: デザインシステムの統一
**目的**: 全ページで統一されたデザインシステムを適用する。

**作業内容**:
- カラーパレットの統一
- タイポグラフィの統一
- スペーシングの統一
- コンポーネントスタイルの統一

**実装内容**:
1. shadcn/uiのデフォルトテーマをベースにカラーパレットを統一
2. フォントサイズ、行間、フォントファミリーを統一
3. Tailwind CSSのスペーシングスケールを統一
4. ボタン、フォーム、カードなどのスタイルを統一

**受け入れ基準**:
- [ ] カラーパレットが統一されている
- [ ] タイポグラフィが統一されている
- [ ] スペーシングが統一されている
- [ ] コンポーネントスタイルが統一されている

- _Requirements: 3.4.1_
- _Design: 6.1_

---

#### タスク 4.2: レスポンシブデザインの改善
**目的**: モバイル・タブレット・デスクトップで適切に表示されるようにする。

**作業内容**:
- モバイル対応の確認と改善
- タブレット対応の確認と改善
- デスクトップ表示の最適化

**実装内容**:
1. 各ページでモバイル（640px以下）の表示を確認・改善
2. 各ページでタブレット（641px - 1024px）の表示を確認・改善
3. 各ページでデスクトップ（1025px以上）の表示を最適化
4. Tailwind CSSのレスポンシブユーティリティを活用

**受け入れ基準**:
- [ ] モバイル対応が適切に実装されている
- [ ] タブレット対応が適切に実装されている
- [ ] デスクトップ表示が最適化されている

- _Requirements: 3.4.2_
- _Design: 6.2_

---

#### タスク 4.3: アクセシビリティの向上
**目的**: アクセシビリティを向上させる。

**作業内容**:
- キーボードナビゲーションの確認
- スクリーンリーダー対応の確認
- shadcn/uiのアクセシビリティ機能を活用

**実装内容**:
1. すべてのインタラクティブ要素がキーボードで操作可能であることを確認
2. フォーカス表示を明確にする
3. 適切なARIA属性を使用
4. shadcn/uiのアクセシビリティ機能を活用

**受け入れ基準**:
- [ ] キーボードナビゲーションが適切に実装されている
- [ ] スクリーンリーダー対応が適切に実装されている
- [ ] shadcn/uiのアクセシビリティ機能が活用されている

- _Requirements: 3.4.3_
- _Design: 6.3_

---

### Phase 5: テストの移植

#### タスク 5.1: テスト環境のセットアップ
**目的**: テスト実行環境を整備する。

**作業内容**:
- Playwrightのインストールと設定
- Jestのインストールと設定（タスク実行時に必要と判断した場合は実施。タスク実行者の判断）
- テスト用の環境変数設定
- `playwright.config.ts`の作成

**実装内容**:
1. Playwrightをインストール:
   ```bash
   npm install -D @playwright/test --legacy-peer-deps
   npx playwright install
   ```
2. `playwright.config.ts`を作成
3. Jestをインストール（必要に応じて）:
   ```bash
   npm install -D jest @types/jest ts-jest --legacy-peer-deps
   ```
4. テスト用の環境変数を設定

**受け入れ基準**:
- [ ] Playwrightがインストールされている
- [ ] `playwright.config.ts`が作成されている
- [ ] テスト用の環境変数が設定されている
- [ ] Jestがインストールされている（必要に応じて）

- _Requirements: 3.5.1_
- _Design: 7.1_

---

#### タスク 5.2: E2Eテストの移植・改修
**目的**: 既存のE2Eテストを移植し、NextAuth対応に変更する。

**作業内容**:
- `client/e2e/`のテストファイルを`client2/e2e/`に移植
- NextAuth対応に変更

**実装内容**:
1. `client2/e2e/`ディレクトリを作成
2. 以下のテストファイルを移植・改修:
   - `client/e2e/auth-flow.spec.ts` → `client2/e2e/auth-flow.spec.ts`（NextAuth対応に変更）
   - `client/e2e/user-flow.spec.ts` → `client2/e2e/user-flow.spec.ts`
   - `client/e2e/post-flow.spec.ts` → `client2/e2e/post-flow.spec.ts`
   - `client/e2e/cross-shard.spec.ts` → `client2/e2e/cross-shard.spec.ts`
   - `client/e2e/email-send.spec.ts` → `client2/e2e/email-send.spec.ts`
   - `client/e2e/csv-download.spec.ts` → `client2/e2e/csv-download.spec.ts`
3. 認証フローのテストをNextAuth対応に変更

**受け入れ基準**:
- [ ] `client2/e2e/`ディレクトリが作成されている
- [ ] すべてのE2Eテストファイルが移植されている
- [ ] 認証フローのテストがNextAuth対応になっている
- [ ] すべてのテストが正常に動作する

- _Requirements: 3.5.2_
- _Design: 7.2_

---

#### タスク 5.3: 統合テストの移植
**目的**: 既存の統合テストを移植し、NextAuth対応に変更する。

**作業内容**:
- `client/src/__tests__/integration/`のテストを`client2/src/__tests__/integration/`に移植
- NextAuth対応に変更

**実装内容**:
1. `client2/src/__tests__/integration/`ディレクトリを作成
2. `client/src/__tests__/integration/`のテストファイルを移植
3. NextAuth対応に変更
4. テストが正常に動作することを確認

**受け入れ基準**:
- [ ] `client2/src/__tests__/integration/`ディレクトリが作成されている
- [ ] すべての統合テストファイルが移植されている
- [ ] NextAuth対応に変更されている
- [ ] すべてのテストが正常に動作する

- _Requirements: 3.5.3_
- _Design: 7.3_

---

#### タスク 5.4: 単体テストの移植
**目的**: 既存の単体テストを移植する。

**作業内容**:
- `client/src/__tests__/components/`のテストを移植
- `client/src/__tests__/lib/`のテストを移植

**実装内容**:
1. `client2/src/__tests__/components/`ディレクトリを作成
2. `client2/src/__tests__/lib/`ディレクトリを作成
3. `client/src/__tests__/components/`のテストファイルを移植
4. `client/src/__tests__/lib/`のテストファイルを移植
5. テストが正常に動作することを確認

**受け入れ基準**:
- [ ] `client2/src/__tests__/components/`ディレクトリが作成されている
- [ ] `client2/src/__tests__/lib/`ディレクトリが作成されている
- [ ] すべての単体テストファイルが移植されている
- [ ] すべてのテストが正常に動作する

- _Requirements: 3.5.4_
- _Design: 7.4_

---

### Phase 6: 最終確認とドキュメント

#### タスク 6.1: 動作確認
**目的**: すべての機能が正常に動作することを確認する。

**作業内容**:
- すべてのページの動作確認
- 認証フローの確認
- API呼び出しの確認
- エラーハンドリングの確認

**実装内容**:
1. すべてのページを確認:
   - トップページ
   - ユーザー管理ページ
   - 投稿管理ページ
   - ユーザーと投稿のJOINページ
   - メール送信ページ
   - ジョブキューページ
   - 動画アップロードページ
2. 認証フローの確認:
   - ログイン
   - ログアウト
   - トークン取得
3. API呼び出しの確認（すべてのエンドポイント）
4. エラーハンドリングの確認

**受け入れ基準**:
- [ ] すべてのページが正常に動作する
- [ ] 認証フローが正常に動作する
- [ ] API呼び出しが正常に動作する
- [ ] エラーハンドリングが適切に実装されている

- _Requirements: 3.6.1_
- _Design: 8.1_

---

#### タスク 6.2: パフォーマンス確認
**目的**: パフォーマンスが適切であることを確認する。

**作業内容**:
- ページ読み込み速度の確認
- API呼び出しの最適化（必要に応じて）

**実装内容**:
1. 各ページの読み込み速度を確認
2. API呼び出しのレスポンス時間を確認
3. 必要に応じて最適化を実施

**受け入れ基準**:
- [ ] ページ読み込み速度が適切である
- [ ] API呼び出しのレスポンス時間が適切である
- [ ] 必要に応じて最適化が実施されている

- _Requirements: 3.6.2_
- _Design: 8.2_

---

#### タスク 6.3: ドキュメント更新
**目的**: ドキュメントを更新し、移植後の状態を反映する。

**作業内容**:
- `docs/Temp-Client2.md`の作成・更新
- 環境変数のドキュメント化
- セットアップ手順のドキュメント化
- 機能説明のドキュメント化

**実装内容**:
1. `docs/Temp-Client2.md`を作成・更新:
   - プロジェクトの概要
   - セットアップ手順
   - 環境変数の説明
   - 機能説明
   - 技術スタック
2. 環境変数のドキュメント化
3. セットアップ手順のドキュメント化
4. 機能説明のドキュメント化

**受け入れ基準**:
- [ ] `docs/Temp-Client2.md`が作成・更新されている
- [ ] 環境変数がドキュメント化されている
- [ ] セットアップ手順がドキュメント化されている
- [ ] 機能説明がドキュメント化されている

- _Requirements: 3.6.3_
- _Design: 8.3_

---

## 実装順序

### Phase 1: 基盤整備（認証・API・型定義）
1. タスク 1.1: 型定義の移植
2. タスク 1.2: NextAuth設定の拡張
3. タスク 1.3: APIクライアントの移植・改修
4. タスク 1.4: 認証ヘルパーの実装

### Phase 2: 共通コンポーネントとレイアウト
1. タスク 2.1: レイアウトコンポーネントの改修
2. タスク 2.2: shadcn/uiコンポーネントの追加インストール
3. タスク 2.3: 共通UIコンポーネントの作成

### Phase 3: ページ移植（シンプルなものから順に）
1. タスク 3.1: トップページの移植・改修
2. タスク 3.2: TodayApiButtonコンポーネントの移植
3. タスク 3.3: ユーザー管理ページの移植
4. タスク 3.4: 投稿管理ページの移植
5. タスク 3.5: ユーザーと投稿のJOINページの移植
6. タスク 3.6: メール送信ページの移植
7. タスク 3.7: ジョブキューページの移植
8. タスク 3.8: 動画アップロードページの移植

### Phase 4: デザイン改善
- タスク 4.1: デザインシステムの統一
- タスク 4.2: レスポンシブデザインの改善
- タスク 4.3: アクセシビリティの向上

**注意**: Phase 4はPhase 3と並行して実施可能

### Phase 5: テストの移植
1. タスク 5.1: テスト環境のセットアップ
2. タスク 5.2: E2Eテストの移植・改修
3. タスク 5.3: 統合テストの移植
4. タスク 5.4: 単体テストの移植

### Phase 6: 最終確認とドキュメント
1. タスク 6.1: 動作確認
2. タスク 6.2: パフォーマンス確認
3. タスク 6.3: ドキュメント更新

## 依存関係

- Phase 2はPhase 1完了後に開始
- Phase 3はPhase 1とPhase 2完了後に開始
- Phase 4はPhase 3と並行して実施可能
- Phase 5はPhase 3完了後に開始
- Phase 6はすべてのPhase完了後に実施

## 注意事項

### 認証の移行
- Auth0からNextAuthへの移行で認証フローが異なるため、慎重に実装する
- トークン取得方法が異なるため、`lib/api.ts`と`lib/auth.ts`の実装を変更する

### UIコンポーネントの移行
- 既存のカスタムコンポーネントをshadcn/uiコンポーネントに置き換える
- スタイルの統一が必要

### テストの移行
- すべてのテストをNextAuth対応に変更する
- 認証フローのテストをNextAuthの認証フローに合わせて調整する

### デザイン改善
- 各ページ移植時に並行してデザイン改善を実施する
- 既存の機能を損なわずに、ユーザビリティを向上させる
