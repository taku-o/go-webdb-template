# 既存clientアプリの機能をclient2アプリに移植する要件定義書

## 1. 概要

### 1.1 プロジェクト情報
- **プロジェクト名**: go-webdb-template
- **Feature名**: 0064-client2-impl
- **作成日**: 2026-01-27
- **関連Issue**: （該当する場合は記載）

### 1.2 目的
既存の`client`アプリケーションの機能を、新しく作成した`client2`アプリケーションに移植する。NextAuth (Auth.js) v5とshadcn/uiを活用したモダンな実装に移行し、同時にデザインを改善する。

### 1.3 スコープ
- 既存の`client`アプリの全機能を`client2`アプリに移植
- Auth0からNextAuth (Auth.js) v5への認証方式の移行
- カスタムUIからshadcn/uiへの移行
- デザインの改善（モダンなUI/UX）
- テストコードの移植と改修

**本実装の範囲外**:
- 既存の`client`ディレクトリの削除（別タスクで実施）
- 新機能の追加（既存機能の移植のみ）
- バックエンドAPIの変更（既存APIとの互換性を維持）

## 2. 背景・現状分析

### 2.1 既存のclientアプリの状況

#### 2.1.1 技術スタック
- **フレームワーク**: Next.js 14+ (App Router)
- **認証**: Auth0 (`@auth0/nextjs-auth0`)
- **UI**: カスタム実装（Tailwind CSS）
- **言語**: TypeScript 5+

#### 2.1.2 実装されている機能
1. **認証機能**
   - Auth0によるログイン/ログアウト
   - 認証トークンの取得（`/auth/token`）
   - プロフィール取得（`/auth/profile`）

2. **ページ機能**
   - トップページ（機能一覧、認証状態表示）
   - ユーザー管理ページ（CRUD、CSVダウンロード）
   - 投稿管理ページ（CRUD）
   - ユーザーと投稿のJOINページ（クロスシャードクエリ）
   - メール送信ページ
   - 動画アップロードページ（TUSプロトコル、Uppy）
   - ジョブキューページ

3. **コンポーネント**
   - TodayApiButton（プライベートAPIテスト用）

4. **APIクライアント**
   - `lib/api.ts`: バックエンドAPI呼び出し用クライアント
   - Auth0トークン取得対応

5. **型定義**
   - `types/dm_post.ts`: 投稿関連の型定義
   - `types/dm_user.ts`: ユーザー関連の型定義
   - `types/jobqueue.ts`: ジョブキュー関連の型定義

6. **テスト**
   - E2Eテスト（Playwright）: 6つのテストファイル
   - 統合テスト: 2つのテストファイル
   - 単体テスト: 2つのテストファイル

### 2.2 client2アプリの現状

#### 2.2.1 技術スタック
- **フレームワーク**: Next.js 14+ (App Router)
- **認証**: NextAuth (Auth.js) v5（基本設定済み）
- **UI**: shadcn/ui（8コンポーネントインストール済み）
- **言語**: TypeScript 5+
- **ベース**: precedentテンプレート

#### 2.2.2 既に実装済みの内容
- precedentテンプレートベースの基本構造
- NextAuth (Auth.js) v5の基本設定（`auth.ts`, `/api/auth/[...nextauth]/route.ts`）
- shadcn/uiの統合（`components.json`設定済み）
- 基本的なレイアウト構造（`app/layout.tsx`, `components/layout/navbar.tsx`）
- 環境変数の設定（`.env.example`, `.env.local`）

### 2.3 移植の必要性
- 既存の`client`アプリのデザインが雑で改善が必要
- よりモダンなUIコンポーネントライブラリ（shadcn/ui）の採用
- NextAuth (Auth.js)による認証の標準化
- コードの保守性向上

## 3. 機能要件

### 3.1 Phase 1: 基盤整備（認証・API・型定義）

#### 3.1.1 型定義の移植
- **目的**: 既存の型定義をclient2に移植する
- **実装内容**:
  - `client/src/types/dm_post.ts` → `client2/types/dm_post.ts`に移植
  - `client/src/types/dm_user.ts` → `client2/types/dm_user.ts`に移植
  - `client/src/types/jobqueue.ts` → `client2/types/jobqueue.ts`に移植
  - `client2/tsconfig.json`に`@/types`エイリアスを追加（必要に応じて）

#### 3.1.2 NextAuth設定の拡張
- **目的**: NextAuth (Auth.js) v5の設定を拡張し、既存のAuth0機能と同等の機能を実現する
- **実装内容**:
  - `client2/auth.ts`に認証プロバイダー設定を追加（既存のAuth0設定を参考）
  - セッション管理の実装
  - トークン取得用のAPI Route作成（`app/api/auth/token/route.ts`）
    - 既存の`client/src/app/auth/token/route.ts`と同等の機能を実現
  - プロフィール取得用のAPI Route作成（`app/api/auth/profile/route.ts`）
    - 既存の`client/src/app/auth/profile/route.ts`と同等の機能を実現

#### 3.1.3 APIクライアントの移植・改修
- **目的**: 既存のAPIクライアントを移植し、NextAuth対応に変更する
- **実装内容**:
  - `client/src/lib/api.ts`を`client2/lib/api.ts`に移植
  - Auth0依存をNextAuth (Auth.js) v5に置き換え
  - `getAuthToken`関数をNextAuth対応に変更
  - 環境変数の確認（`NEXT_PUBLIC_API_BASE_URL`, `NEXT_PUBLIC_API_KEY`）
  - すべてのAPIメソッドが正常に動作することを確認

#### 3.1.4 認証ヘルパーの実装
- **目的**: NextAuth v5の認証機能をラップしたヘルパー関数を作成する
- **実装内容**:
  - `client2/lib/auth.ts`を作成
  - NextAuth v5の`auth()`, `signIn()`, `signOut()`をラップ
  - クライアント側での認証状態取得用フック（ダミー実装で処理を差し込む場所を用意）

### 3.2 Phase 2: 共通コンポーネントとレイアウト

#### 3.2.1 レイアウトコンポーネントの改修
- **目的**: 認証状態を表示し、ログイン/ログアウト機能を実装する
- **実装内容**:
  - `client2/components/layout/navbar.tsx`に認証状態表示を追加
  - ログイン/ログアウトボタンの実装（NextAuth対応）
  - ユーザー情報表示の実装
  - shadcn/uiの`button`コンポーネントを使用

#### 3.2.2 shadcn/uiコンポーネントの追加インストール
- **目的**: ページ移植に必要なshadcn/uiコンポーネントを追加インストールする
- **実装内容**:
  - `table`: 一覧表示用
  - `textarea`: フォーム用（投稿の本文入力など）
  - `badge`: ステータス表示用
  - `dialog`: モーダル用（既にインストール済みか確認）
  - タスク実行時に必要と判断した場合は追加のコンポーネントをインストール（タスク実行者の判断）

#### 3.2.3 共通UIコンポーネントの作成
- **目的**: ページ間で共通して使用するUIコンポーネントを作成する
- **実装内容**:
  - エラー表示コンポーネント（shadcn/uiの`alert`を使用）
  - ローディング表示コンポーネント
  - フォームコンポーネント（shadcn/uiの`form`を使用）

### 3.3 Phase 3: ページ移植（シンプルなものから順に）

#### 3.3.1 トップページの移植・改修
- **目的**: 既存のトップページ機能を移植し、shadcn/uiでデザインを改善する
- **実装内容**:
  - `client2/app/page.tsx`を既存のトップページ機能に置き換え
  - shadcn/uiの`card`コンポーネントを使用して機能一覧を表示
  - 認証状態の表示（NextAuth対応）
  - TodayApiButtonコンポーネントの統合
  - デザイン改善（モダンなUI）

#### 3.3.2 TodayApiButtonコンポーネントの移植
- **目的**: 既存のTodayApiButtonコンポーネントを移植し、NextAuth対応に変更する
- **実装内容**:
  - `client/src/components/TodayApiButton.tsx` → `client2/components/TodayApiButton.tsx`に移植
  - Auth0依存をNextAuth対応に変更
  - shadcn/uiの`button`コンポーネントを使用
  - エラー表示にshadcn/uiの`alert`を使用

#### 3.3.3 ユーザー管理ページの移植
- **目的**: 既存のユーザー管理ページを移植し、shadcn/uiでデザインを改善する
- **実装内容**:
  - `client2/app/dm-users/page.tsx`を作成
  - shadcn/uiの`form`, `input`, `select`, `button`, `table`, `card`を使用
  - CRUD機能の実装（作成、一覧表示、削除）
  - CSVダウンロード機能の実装
  - デザイン改善（モダンなUI、レスポンシブ対応）

#### 3.3.4 投稿管理ページの移植
- **目的**: 既存の投稿管理ページを移植し、shadcn/uiでデザインを改善する
- **実装内容**:
  - `client2/app/dm-posts/page.tsx`を作成
  - shadcn/uiコンポーネントを使用（`form`, `input`, `textarea`, `select`, `button`, `card`）
  - CRUD機能の実装（作成、一覧表示、削除）
  - ユーザー選択の実装（`select`コンポーネントを使用）
  - デザイン改善（モダンなUI、レスポンシブ対応）

#### 3.3.5 ユーザーと投稿のJOINページの移植
- **目的**: 既存のJOINページを移植し、shadcn/uiでデザインを改善する
- **実装内容**:
  - `client2/app/dm-user-posts/page.tsx`を作成
  - クロスシャードクエリの説明を保持
  - shadcn/uiの`card`を使用して投稿を表示
  - デザイン改善（モダンなUI、レスポンシブ対応）

#### 3.3.6 メール送信ページの移植
- **目的**: 既存のメール送信ページを移植し、shadcn/uiでデザインを改善する
- **実装内容**:
  - `client2/app/dm_email/send/page.tsx`を作成
  - shadcn/uiの`form`, `input`, `button`を使用
  - 成功/エラーメッセージの表示（shadcn/uiの`alert`を使用）
  - デザイン改善（モダンなUI、レスポンシブ対応）

#### 3.3.7 ジョブキューページの移植
- **目的**: 既存のジョブキューページを移植し、shadcn/uiでデザインを改善する
- **実装内容**:
  - `client2/app/dm-jobqueue/page.tsx`を作成
  - shadcn/uiコンポーネントを使用（`form`, `input`, `button`, `alert`）
  - デザイン改善（モダンなUI、レスポンシブ対応）

#### 3.3.8 動画アップロードページの移植
- **目的**: 既存の動画アップロードページを移植し、NextAuth対応に変更する
- **実装内容**:
  - `client2/app/dm_movie/upload/page.tsx`を作成
  - Uppyライブラリのインストール（`@uppy/core`, `@uppy/react`, `@uppy/tus`, `@uppy/dashboard`）
  - TUSプロトコル対応の実装
  - 認証トークンの取得方法をNextAuth対応に変更
  - デザイン改善（モダンなUI、レスポンシブ対応）

### 3.4 Phase 4: デザイン改善

#### 3.4.1 デザインシステムの統一
- **目的**: 全ページで統一されたデザインシステムを適用する
- **実装内容**:
  - カラーパレットの統一（shadcn/uiのデフォルトテーマをベース）
  - タイポグラフィの統一
  - スペーシングの統一
  - ボタン、フォーム、カードなどのスタイル統一

#### 3.4.2 レスポンシブデザインの改善
- **目的**: モバイル・タブレット・デスクトップで適切に表示されるようにする
- **実装内容**:
  - モバイル対応の確認と改善
  - タブレット対応の確認と改善
  - デスクトップ表示の最適化

#### 3.4.3 アクセシビリティの向上
- **目的**: アクセシビリティを向上させる
- **実装内容**:
  - shadcn/uiのアクセシビリティ機能を活用
  - キーボードナビゲーションの確認
  - スクリーンリーダー対応の確認

### 3.5 Phase 5: テストの移植

#### 3.5.1 テスト環境のセットアップ
- **目的**: テスト実行環境を整備する
- **実装内容**:
  - Playwrightのインストールと設定
  - Jestのインストールと設定（タスク実行時に必要と判断した場合は実施。タスク実行者の判断）
  - テスト用の環境変数設定
  - `playwright.config.ts`の作成

#### 3.5.2 E2Eテストの移植・改修
- **目的**: 既存のE2Eテストを移植し、NextAuth対応に変更する
- **実装内容**:
  - `client/e2e/auth-flow.spec.ts` → NextAuth対応に変更
  - `client/e2e/user-flow.spec.ts` → 移植
  - `client/e2e/post-flow.spec.ts` → 移植
  - `client/e2e/cross-shard.spec.ts` → 移植
  - `client/e2e/email-send.spec.ts` → 移植
  - `client/e2e/csv-download.spec.ts` → 移植

#### 3.5.3 統合テストの移植
- **目的**: 既存の統合テストを移植し、NextAuth対応に変更する
- **実装内容**:
  - `client/src/__tests__/integration/`のテストを移植
  - NextAuth対応に変更
  - テストが正常に動作することを確認

#### 3.5.4 単体テストの移植
- **目的**: 既存の単体テストを移植する
- **実装内容**:
  - `client/src/components/__tests__/`のテストを移植
  - `client/src/lib/__tests__/`のテストを移植
  - テストが正常に動作することを確認

### 3.6 Phase 6: 最終確認とドキュメント

#### 3.6.1 動作確認
- **目的**: すべての機能が正常に動作することを確認する
- **実装内容**:
  - すべてのページの動作確認
  - 認証フローの確認（ログイン、ログアウト、トークン取得）
  - API呼び出しの確認（すべてのエンドポイント）
  - エラーハンドリングの確認

#### 3.6.2 パフォーマンス確認
- **目的**: パフォーマンスが適切であることを確認する
- **実装内容**:
  - ページ読み込み速度の確認
  - API呼び出しの最適化（必要に応じて）

#### 3.6.3 ドキュメント更新
- **目的**: ドキュメントを更新し、移植後の状態を反映する
- **実装内容**:
  - `docs/Temp-Client2.md`の作成・更新
  - 環境変数のドキュメント化
  - セットアップ手順のドキュメント化
  - 機能説明のドキュメント化

## 4. 非機能要件

### 4.1 パフォーマンス
- **ページ読み込み速度**: 既存のclientアプリと同等またはそれ以上の速度
- **API呼び出し**: 既存のclientアプリと同等のレスポンス時間
- **ビルド時間**: 既存のclientアプリと同等またはそれ以上のビルド時間

### 4.2 信頼性
- **エラーハンドリング**: すべてのエラーケースを適切に処理する
- **型安全性**: TypeScriptの型チェックが正常に動作する
- **APIエラーハンドリング**: ネットワークエラー、認証エラーなどを適切に処理する

### 4.3 保守性
- **コードの可読性**: shadcn/uiとNextAuthの標準的なパターンに従う
- **一貫性**: 全ページで統一されたコンポーネントとスタイルを使用
- **ドキュメント**: コードの変更点と移植内容をドキュメント化

### 4.4 互換性
- **バックエンドAPI**: 既存のバックエンドAPIと完全に互換性があること
- **ブラウザ**: モダンブラウザ（Chrome、Firefox、Safari、Edge）で動作すること
- **レスポンシブ**: モバイル、タブレット、デスクトップで適切に表示されること

### 4.5 セキュリティ
- **認証**: NextAuth (Auth.js) v5のセキュリティベストプラクティスに従う
- **トークン管理**: 認証トークンを適切に管理する
- **API呼び出し**: 認証が必要なAPI呼び出しで適切にトークンを使用する

## 5. 制約事項

### 5.1 技術的制約
- **既存のclientディレクトリ**: 既存の`client`ディレクトリには影響を与えないこと
- **バックエンドAPI**: 既存のバックエンドAPIの変更は行わないこと
- **Next.js App Router**: Next.js 14+のApp Routerを使用すること
- **TypeScript**: TypeScript 5+を使用すること

### 5.2 実装上の制約
- **ディレクトリ構造**: `client2/`ディレクトリに実装すること
- **認証方式**: NextAuth (Auth.js) v5を使用すること（Auth0は使用しない）
- **UIコンポーネント**: shadcn/uiコンポーネントを使用すること（カスタムコンポーネントは最小限）
- **デザイン**: shadcn/uiのデフォルトテーマをベースにすること

### 5.3 動作環境
- **ローカル環境**: ローカル環境で開発サーバーが正常に起動すること
- **ポート**: 開発サーバーはポート3000で起動すること
- **ブラウザ**: モダンブラウザ（Chrome、Firefox、Safari、Edge）で動作することを前提

### 5.4 既存機能の維持
- **機能の完全性**: 既存のclientアプリのすべての機能を移植すること
- **API互換性**: 既存のバックエンドAPIとの互換性を維持すること
- **データ構造**: 既存のデータ構造（型定義）を維持すること

## 6. 受け入れ基準

### 6.1 Phase 1: 基盤整備
- [ ] 型定義がすべて移植されている（`dm_post.ts`, `dm_user.ts`, `jobqueue.ts`）
- [ ] NextAuth設定が拡張され、トークン取得とプロフィール取得が動作する
- [ ] APIクライアントが移植され、NextAuth対応に変更されている
- [ ] 認証ヘルパーが実装され、正常に動作する
- [ ] すべてのAPIメソッドが正常に動作する

### 6.2 Phase 2: 共通コンポーネントとレイアウト
- [ ] レイアウトコンポーネントに認証状態表示とログイン/ログアウト機能が実装されている
- [ ] 必要なshadcn/uiコンポーネントがすべてインストールされている
- [ ] 共通UIコンポーネント（エラー表示、ローディング、フォーム）が作成されている

### 6.3 Phase 3: ページ移植
- [ ] トップページが移植され、デザインが改善されている
- [ ] TodayApiButtonコンポーネントが移植され、NextAuth対応になっている
- [ ] ユーザー管理ページが移植され、CRUDとCSVダウンロードが動作する
- [ ] 投稿管理ページが移植され、CRUDが動作する
- [ ] ユーザーと投稿のJOINページが移植されている
- [ ] メール送信ページが移植されている
- [ ] ジョブキューページが移植されている
- [ ] 動画アップロードページが移植され、NextAuth対応になっている
- [ ] すべてのページでshadcn/uiコンポーネントが使用されている
- [ ] すべてのページでデザインが改善されている

### 6.4 Phase 4: デザイン改善
- [ ] デザインシステムが統一されている（カラー、タイポグラフィ、スペーシング）
- [ ] レスポンシブデザインが適切に実装されている（モバイル、タブレット、デスクトップ）
- [ ] アクセシビリティが向上している（キーボードナビゲーション、スクリーンリーダー対応）

### 6.5 Phase 5: テストの移植
- [ ] テスト環境がセットアップされている（Playwright、Jest）
- [ ] E2Eテストがすべて移植され、NextAuth対応になっている
- [ ] 統合テストがすべて移植され、NextAuth対応になっている
- [ ] 単体テストがすべて移植されている
- [ ] すべてのテストが正常に動作する

### 6.6 Phase 6: 最終確認とドキュメント
- [ ] すべてのページの動作確認が完了している
- [ ] 認証フローの確認が完了している
- [ ] API呼び出しの確認が完了している
- [ ] エラーハンドリングの確認が完了している
- [ ] パフォーマンス確認が完了している
- [ ] ドキュメントが更新されている（README、環境変数、セットアップ手順）

## 7. 影響範囲

### 7.1 変更が必要なファイル

#### 新規作成が必要なファイル
- `client2/types/`: 型定義ディレクトリ
  - `client2/types/dm_post.ts`
  - `client2/types/dm_user.ts`
  - `client2/types/jobqueue.ts`
- `client2/lib/`: ライブラリディレクトリ
  - `client2/lib/api.ts`（既存の`lib/utils.ts`とは別）
  - `client2/lib/auth.ts`
- `client2/app/api/auth/`: 認証API Route
  - `client2/app/api/auth/token/route.ts`
  - `client2/app/api/auth/profile/route.ts`
- `client2/app/`: ページディレクトリ
  - `client2/app/page.tsx`（既存を置き換え）
  - `client2/app/dm-users/page.tsx`
  - `client2/app/dm-posts/page.tsx`
  - `client2/app/dm-user-posts/page.tsx`
  - `client2/app/dm_email/send/page.tsx`
  - `client2/app/dm-jobqueue/page.tsx`
  - `client2/app/dm_movie/upload/page.tsx`
- `client2/components/`: コンポーネントディレクトリ
  - `client2/components/TodayApiButton.tsx`
  - `client2/components/layout/navbar.tsx`（既存を改修）
- `client2/e2e/`: E2Eテストディレクトリ
  - `client2/e2e/auth-flow.spec.ts`
  - `client2/e2e/user-flow.spec.ts`
  - `client2/e2e/post-flow.spec.ts`
  - `client2/e2e/cross-shard.spec.ts`
  - `client2/e2e/email-send.spec.ts`
  - `client2/e2e/csv-download.spec.ts`
- `client2/src/__tests__/`: テストディレクトリ
  - `client2/src/__tests__/integration/`
  - `client2/src/__tests__/components/`
  - `client2/src/__tests__/lib/`

#### 変更が必要なファイル
- `client2/auth.ts`: NextAuth設定の拡張
- `client2/package.json`: 依存関係の追加（Uppy、Playwright、Jest等）
- `client2/tsconfig.json`: 型定義エイリアスの追加（必要に応じて）
- `client2/playwright.config.ts`: 新規作成（E2Eテスト用）

### 7.2 既存機能への影響
- **既存のclientディレクトリ**: 影響なし（独立したディレクトリに実装）
- **既存のバックエンドAPI**: 影響なし（既存APIとの互換性を維持）
- **既存のデータベース**: 影響なし（データ構造を維持）

## 8. 実装上の注意事項

### 8.1 認証の移行
- **Auth0からNextAuthへの移行**: 認証フローとトークン取得方法が異なるため、慎重に実装する
- **トークン取得**: NextAuth v5のAPIを使用してトークンを取得する
- **セッション管理**: NextAuth v5のセッション管理機能を活用する

### 8.2 UIコンポーネントの移行
- **shadcn/uiの使用**: 既存のカスタムコンポーネントをshadcn/uiコンポーネントに置き換える
- **スタイルの統一**: shadcn/uiのデフォルトテーマをベースに、統一されたスタイルを適用する
- **レスポンシブ対応**: すべてのコンポーネントがレスポンシブに対応していることを確認する

### 8.3 APIクライアントの移行
- **認証トークンの取得**: NextAuth対応の`getAuthToken`関数を実装する
- **エラーハンドリング**: 既存のエラーハンドリングロジックを維持する
- **型安全性**: TypeScriptの型定義を維持し、型安全性を確保する

### 8.4 テストの移行
- **NextAuth対応**: すべてのテストをNextAuth対応に変更する
- **認証フローのテスト**: NextAuthの認証フローをテストする
- **E2Eテスト**: Playwrightを使用してE2Eテストを実装する

### 8.5 デザイン改善
- **段階的な改善**: 各ページ移植時に並行してデザイン改善を実施する
- **ユーザビリティ**: 既存の機能を損なわずに、ユーザビリティを向上させる
- **一貫性**: 全ページで統一されたデザインを適用する

## 9. 参考情報

### 9.1 関連ドキュメント
- `.kiro/steering/structure.md`: ファイル組織とコードパターン
- `.kiro/steering/tech.md`: 技術スタックとアーキテクチャ
- `.kiro/steering/product.md`: プロダクトコンテキストとビジネス目標
- `docs/Client2-Setup-Summary.md`: client2アプリのセットアップ作業まとめ
- `.kiro/specs/0063-client2/`: client2アプリ作成の要件定義書

### 9.2 外部リソース
- **NextAuth (Auth.js)**: https://authjs.dev/
- **shadcn/ui**: https://ui.shadcn.com/
- **Next.js**: https://nextjs.org/
- **Playwright**: https://playwright.dev/
- **Uppy**: https://uppy.io/

### 9.3 技術スタック
- **フレームワーク**: Next.js 14+ (App Router)
- **言語**: TypeScript 5+
- **UIコンポーネント**: shadcn/ui
- **認証**: NextAuth (Auth.js) v5
- **スタイリング**: Tailwind CSS
- **テスト**: Playwright (E2E), Jest (単体・統合)
- **ファイルアップロード**: Uppy (TUSプロトコル)

### 9.4 実装の流れ
1. Phase 1: 基盤整備（認証・API・型定義）
2. Phase 2: 共通コンポーネントとレイアウト
3. Phase 3: ページ移植（シンプルなものから順に）
4. Phase 4: デザイン改善（Phase 3と並行して実施）
5. Phase 5: テストの移植
6. Phase 6: 最終確認とドキュメント

### 9.5 依存関係
- Phase 2はPhase 1完了後に開始
- Phase 3はPhase 1とPhase 2完了後に開始
- Phase 4はPhase 3と並行して実施可能
- Phase 5はPhase 3完了後に開始
- Phase 6はすべてのPhase完了後に実施
