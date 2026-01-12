# 新clientアプリの作成の要件定義書

## 1. 概要

### 1.1 プロジェクト情報
- **プロジェクト名**: go-webdb-template
- **Feature名**: 0063-client2
- **作成日**: 2026-01-27
- **関連Issue**: https://github.com/taku-o/go-webdb-template/issues/130

### 1.2 目的
既存の`client`ディレクトリとは別に、新しいクライアントアプリケーションを`client2`ディレクトリに作成する。shadcn/uiとNextAuth (Auth.js)を組み合わせた構成で、precedentテンプレートをベースに構築する。

### 1.3 スコープ
- `client2`ディレクトリの作成
- precedentテンプレート（https://github.com/steven-tey/precedent）をベースとした初期セットアップ
- shadcn/uiの統合
- NextAuth (Auth.js)の統合
- 基本的なプロジェクト構造の確立

**本実装の範囲外**:
- 既存の`client`ディレクトリの機能移行（別タスクで実施）
- 既存のAPIエンドポイントとの統合（別タスクで実施）
- 詳細な機能実装（別タスクで実施）

## 2. 背景・現状分析

### 2.1 現在の状況

#### 2.1.1 既存のclientディレクトリ
- **実装場所**: `client/`ディレクトリ
- **技術スタック**: Next.js 14+ (App Router), TypeScript 5+, Auth0 Next.js SDK
- **認証方式**: Auth0 (`@auth0/nextjs-auth0`)
- **UI**: カスタム実装

#### 2.1.2 新clientアプリの必要性
- よりモダンなUIコンポーネントライブラリ（shadcn/ui）の採用
- NextAuth (Auth.js)による認証の標準化
- precedentテンプレートによる開発効率の向上
- 既存の`client`ディレクトリとは独立した新規アプリケーションとして構築

### 2.2 precedentテンプレートについて

precedentテンプレート（https://github.com/steven-tey/precedent）は、以下の技術スタックを含むオピニオネイテッドなNext.jsプロジェクトテンプレートです：

- **Next.js**: Reactフレームワーク
- **Auth.js (NextAuth.js)**: 認証ライブラリ
- **Prisma**: TypeScript-first ORM
- **Tailwind CSS**: ユーティリティファーストのCSSフレームワーク
- **Radix UI**: アクセシブルなUIコンポーネント
- **Framer Motion**: アニメーションライブラリ
- **TypeScript**: 型安全性
- **Prettier & ESLint**: コード品質ツール
- **Vercel Analytics**: パフォーマンス分析

**注意**: precedentテンプレートには標準でshadcn/uiは含まれていませんが、本要件ではshadcn/uiを追加で統合する。

## 3. 機能要件

### 3.1 プロジェクト構造の作成

#### 3.1.1 client2ディレクトリの作成
- **目的**: 既存の`client`ディレクトリとは独立した新規アプリケーションを作成
- **実装内容**:
  - `client2/`ディレクトリをプロジェクトルートに作成
  - precedentテンプレートをベースとした初期セットアップ
  - 既存の`client`ディレクトリには影響を与えない

#### 3.1.2 ディレクトリ構造
- **基本構造**: precedentテンプレートの標準構造に従う
- **主要ディレクトリ**:
  - `client2/app/`: Next.js App Routerのページとルート
  - `client2/components/`: Reactコンポーネント
  - `client2/lib/`: ユーティリティ関数
  - `client2/public/`: 静的ファイル
  - `client2/types/`: TypeScript型定義

### 3.2 技術スタックの統合

#### 3.2.1 Next.jsのセットアップ
- **目的**: Next.js 14+ (App Router)のセットアップ
- **実装内容**:
  - precedentテンプレートからNext.js設定を取得
  - `next.config.js`の設定
  - `tsconfig.json`の設定
  - 基本的なページ構造の作成

#### 3.2.2 shadcn/uiの統合
- **目的**: shadcn/uiコンポーネントライブラリの統合
- **実装内容**:
  - shadcn/uiの初期化（`npx shadcn-ui@latest init`）
  - `components.json`の設定
  - Tailwind CSSの設定
  - 基本的なコンポーネントのインストール（必要最小限）

#### 3.2.3 NextAuth (Auth.js)の統合
- **目的**: NextAuth (Auth.js)による認証機能の統合
- **実装内容**:
  - NextAuth (Auth.js)のインストール
  - 認証設定ファイルの作成
  - 認証プロバイダーの設定（必要最小限）
  - 認証ルートの作成（`/api/auth/[...nextauth]/route.ts`など）

#### 3.2.4 その他の依存関係
- **目的**: precedentテンプレートに含まれるその他の依存関係の統合
- **実装内容**:
  - Prismaのセットアップ（必要に応じて）
  - Tailwind CSSの設定
  - TypeScriptの設定
  - Prettier & ESLintの設定
  - その他、precedentテンプレートに含まれる標準的な依存関係

### 3.3 基本的なプロジェクト設定

#### 3.3.1 パッケージ管理
- **目的**: 依存関係の管理
- **実装内容**:
  - `package.json`の作成
  - 必要な依存関係のインストール
  - スクリプトの設定（`dev`, `build`, `start`など）

#### 3.3.2 環境変数の設定
- **目的**: 環境変数の管理
- **実装内容**:
  - `.env.example`の作成
  - `.env.local`の設定（開発環境用）
  - 必要な環境変数のドキュメント化

#### 3.3.3 開発環境の設定
- **目的**: 開発環境の整備
- **実装内容**:
  - 開発サーバーの起動確認
  - ホットリロードの動作確認
  - 基本的なページの表示確認

## 4. 非機能要件

### 4.1 パフォーマンス
- **開発サーバーの起動**: 開発サーバーが正常に起動すること
- **ビルド時間**: 初回ビルドが正常に完了すること
- **ホットリロード**: ファイル変更時のホットリロードが正常に動作すること

### 4.2 信頼性
- **エラーハンドリング**: 基本的なエラーハンドリングが実装されていること
- **型安全性**: TypeScriptの型チェックが正常に動作すること

### 4.3 保守性
- **コードの可読性**: precedentテンプレートの標準的な構造に従うこと
- **一貫性**: 既存のプロジェクト構造と整合性があること（可能な範囲で）
- **ドキュメント**: 基本的なREADMEまたはドキュメントの作成

### 4.4 互換性
- **Node.jsバージョン**: プロジェクトで使用しているNode.jsバージョンと互換性があること
- **TypeScript**: TypeScript 5+を使用すること
- **Next.js**: Next.js 14+ (App Router)を使用すること

## 5. 制約事項

### 5.1 技術的制約
- **既存のclientディレクトリ**: 既存の`client`ディレクトリには影響を与えないこと
- **precedentテンプレート**: precedentテンプレートの標準的な構造に従うこと
- **Next.js App Router**: Next.js 14+のApp Routerを使用すること
- **TypeScript**: TypeScript 5+を使用すること

### 5.2 実装上の制約
- **ディレクトリ構造**: `client2/`ディレクトリに作成すること
- **命名規則**: precedentテンプレートの標準的な命名規則に従うこと
- **ライブラリの利用**: precedentテンプレートに含まれる標準的なライブラリを優先的に利用すること

### 5.3 動作環境
- **ローカル環境**: ローカル環境で開発サーバーが正常に起動すること
- **ポート**: 開発サーバーはポート3000で起動すること
- **ブラウザ**: モダンブラウザ（Chrome、Firefox、Safari、Edge）で動作することを前提

## 6. 受け入れ基準

### 6.1 プロジェクト構造の作成
- [ ] `client2/`ディレクトリが作成されている
- [ ] precedentテンプレートをベースとした初期セットアップが完了している
- [ ] 基本的なディレクトリ構造が整備されている（`app/`, `components/`, `lib/`, `public/`, `types/`など）

### 6.2 技術スタックの統合
- [ ] Next.js 14+ (App Router)がセットアップされている
- [ ] shadcn/uiが統合されている（初期化が完了している）
- [ ] NextAuth (Auth.js)が統合されている（基本的な設定が完了している）
- [ ] Tailwind CSSが設定されている
- [ ] TypeScriptが設定されている

### 6.3 基本的なプロジェクト設定
- [ ] `package.json`が作成され、必要な依存関係がインストールされている
- [ ] `.env.example`が作成されている
- [ ] 開発サーバーが正常に起動する（`npm run dev`が動作する）
- [ ] 基本的なページ（トップページなど）が表示される

### 6.4 動作確認
- [ ] ローカル環境で開発サーバーが正常に起動する（ポート3000で起動する）
- [ ] ブラウザで基本的なページが表示される（`http://localhost:3000`でアクセス可能）
- [ ] ホットリロードが正常に動作する
- [ ] TypeScriptの型チェックが正常に動作する
- [ ] ビルドが正常に完了する（`npm run build`が動作する）

### 6.5 ドキュメント
- [ ] 基本的なREADMEまたはドキュメントが作成されている（必要最小限）
- [ ] 環境変数の設定方法がドキュメント化されている

## 7. 影響範囲

### 7.1 変更が必要なファイル

#### 新規作成が必要なファイル
- `client2/`: 新規ディレクトリ全体
  - `client2/package.json`: パッケージ設定
  - `client2/next.config.js`: Next.js設定
  - `client2/tsconfig.json`: TypeScript設定
  - `client2/tailwind.config.js`: Tailwind CSS設定
  - `client2/components.json`: shadcn/ui設定
  - `client2/.env.example`: 環境変数テンプレート
  - `client2/app/`: Next.js App Routerのページ
  - `client2/components/`: Reactコンポーネント
  - `client2/lib/`: ユーティリティ関数
  - その他、precedentテンプレートに含まれる標準的なファイル

### 7.2 既存機能への影響
- **既存のclientディレクトリ**: 影響なし（独立したディレクトリに作成）
- **既存のAPIエンドポイント**: 影響なし（本実装ではAPI統合は行わない）
- **既存のサーバー**: 影響なし（クライアント側のみの実装）

## 8. 実装上の注意事項

### 8.1 precedentテンプレートの利用
- **テンプレートの取得**: precedentテンプレートを適切な方法で取得する（GitHubからクローン、またはテンプレートとして使用）
- **カスタマイズ**: precedentテンプレートの標準的な構造を維持しつつ、必要最小限のカスタマイズを行う
- **不要な機能の削除**: 本要件で不要な機能（Prismaなど）は削除または無効化する（必要に応じて）

### 8.2 shadcn/uiの統合
- **初期化**: `npx shadcn-ui@latest init`を実行して初期化する
- **コンポーネントのインストール**: 必要最小限のコンポーネントのみをインストールする（Button、Cardなど）
- **設定**: `components.json`を適切に設定する

### 8.3 NextAuth (Auth.js)の統合
- **インストール**: NextAuth (Auth.js)を適切にインストールする
- **設定**: 認証設定ファイルを適切に作成する
- **プロバイダー**: 必要最小限のプロバイダー（例: Credentials、Googleなど）を設定する
- **ルート**: 認証ルート（`/api/auth/[...nextauth]/route.ts`など）を適切に作成する

### 8.4 環境変数の管理
- **`.env.example`**: 必要な環境変数を`.env.example`に記載する
- **`.env.local`**: 開発環境用の`.env.local`を作成する（gitignoreに含まれる）
- **ドキュメント**: 環境変数の設定方法をドキュメント化する

### 8.5 既存のclientディレクトリとの関係
- **独立性**: `client2`ディレクトリは既存の`client`ディレクトリとは完全に独立している
- **共存**: 両方のディレクトリが同時に存在しても問題ないようにする
- **ポート**: `client2`の開発サーバーはポート3000を使用する（`next.config.js`または`package.json`で設定）

## 9. 参考情報

### 9.1 関連ドキュメント
- `.kiro/steering/structure.md`: ファイル組織とコードパターン
- `.kiro/steering/tech.md`: 技術スタックとアーキテクチャ
- `.kiro/steering/product.md`: プロダクトコンテキストとビジネス目標

### 9.2 外部リソース
- **precedentテンプレート**: https://github.com/steven-tey/precedent
- **shadcn/ui**: https://ui.shadcn.com/
- **NextAuth (Auth.js)**: https://authjs.dev/
- **Next.js**: https://nextjs.org/

### 9.3 技術スタック
- **フレームワーク**: Next.js 14+ (App Router)
- **言語**: TypeScript 5+
- **UIコンポーネント**: shadcn/ui
- **認証**: NextAuth (Auth.js)
- **スタイリング**: Tailwind CSS
- **UIプリミティブ**: Radix UI（shadcn/uiの基盤）

### 9.4 precedentテンプレートの特徴
- **オピニオネイテッド**: 開発者による推奨設定とパターンが含まれている
- **モダンなスタック**: 最新のNext.js、TypeScript、Tailwind CSSを使用
- **認証統合**: NextAuth (Auth.js)が標準で統合されている
- **開発体験**: Prettier、ESLint、Vercel Analyticsなどが含まれている

### 9.5 実装の流れ
1. precedentテンプレートをベースに`client2`ディレクトリを作成
2. shadcn/uiを初期化して統合
3. NextAuth (Auth.js)の設定を確認・調整
4. 基本的なページ構造を作成
5. 開発サーバーの起動確認
6. 基本的な動作確認
