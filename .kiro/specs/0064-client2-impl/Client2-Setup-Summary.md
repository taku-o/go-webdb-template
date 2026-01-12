# Client2 セットアップ作業まとめ

## 概要

このドキュメントは、`client2`ディレクトリに新しいクライアントアプリケーションを作成した際の作業内容をまとめたものです。`client`から`client2`への機能移植作業で参照することを想定しています。

## 作業完了日

2026年1月12日

## 実装内容

### Phase 1: precedentテンプレートの取得と初期セットアップ

#### タスク 1.1: precedentテンプレートの取得
- **実施内容**: `npx degit steven-tey/precedent client2`でテンプレートを取得
- **結果**: `client2/`ディレクトリが作成され、precedentテンプレートのファイルが配置された

#### タスク 1.2: 依存関係のインストール
- **実施内容**: `npm install --legacy-peer-deps`で依存関係をインストール
- **注意**: peer dependencyの競合があったため`--legacy-peer-deps`フラグを使用
- **結果**: `node_modules/`と`package-lock.json`が生成された

#### タスク 1.3: Prisma関連の削除
- **実施内容**: 
  - `prisma/`ディレクトリを削除
  - `package.json`から`prisma`と`@prisma/client`を削除
  - Prisma関連のスクリプトを削除
  - `.env.example`から`DATABASE_URL`を削除
- **結果**: Prisma関連の依存関係とファイルが完全に削除された

### Phase 2: shadcn/uiの統合

#### タスク 2.1: shadcn/uiの初期化
- **実施内容**: 
  - `components.json`を手動作成（CLIのインタラクティブプロンプトを回避）
  - `tailwind.config.js`にshadcn/ui用の設定を追加
  - `app/globals.css`にshadcn/uiのCSS変数を追加
- **設定内容**:
  - `components.json`: style=default, baseColor=slate, cssVariables=true
  - `tailwind.config.js`: `darkMode: ["class"]`, `tailwindcss-animate`プラグイン追加
  - `tailwind.config.js`: shadcn/ui用のカラー設定（border, input, ring, background, foreground, primary, secondary, destructive, muted, accent, popover, card）
  - `app/globals.css`: CSS変数の定義（:rootと.dark）
- **注意**: shadcn CLIが`pnpm`を使用しようとしたため、手動で`components.json`を作成した

#### タスク 2.2: 既存client移行用コンポーネントのインストール
- **実施内容**: 以下の8つのコンポーネントをインストール
  - `alert-dialog`
  - `alert`
  - `button`
  - `select`
  - `input`
  - `form`
  - `field`
  - `card`
- **依存関係**: 
  - `react-hook-form`と`@hookform/resolvers`を追加（formコンポーネント用）
  - `class-variance-authority`を追加（コンポーネントのvariant管理用）
- **注意**: shadcn CLIが`@/components/ui/`というリテラルパスにファイルを作成したため、`components/ui/`に手動で移動した
- **結果**: 10個のコンポーネントファイルが`components/ui/`に配置された（依存関係含む）

### Phase 3: NextAuth (Auth.js)の統合

#### タスク 3.1: NextAuth (Auth.js)のインストール確認
- **実施内容**: `npm install next-auth@^5.0.0-beta.30 --legacy-peer-deps`
- **結果**: NextAuth v5 (Auth.js)がインストールされた

#### タスク 3.2: NextAuth (Auth.js)認証ルートの作成
- **実施内容**: 
  - `auth.ts`を作成（NextAuth v5の推奨構造）
  - `app/api/auth/[...nextauth]/route.ts`を作成
  - `tsconfig.json`に`@/auth`エイリアスを追加
- **実装内容**:
  ```typescript
  // auth.ts
  import NextAuth from "next-auth"
  
  export const { handlers, auth, signIn, signOut } = NextAuth({
    providers: [],
  })
  
  // app/api/auth/[...nextauth]/route.ts
  import { handlers } from "@/auth"
  export const { GET, POST } = handlers
  ```
- **結果**: NextAuthの認証ルートが作成され、ビルドで正常に動作することを確認

#### タスク 3.3: 環境変数の設定
- **実施内容**: 
  - `.env.example`に`AUTH_SECRET`と`AUTH_URL`を追加
  - `.env.local`を作成し、`AUTH_SECRET`（`openssl rand -base64 32`で生成）と`AUTH_URL=http://localhost:3000`を設定
- **結果**: 環境変数が適切に設定され、`.env.local`は`.gitignore`に含まれていることを確認

### Phase 4: ポート設定と基本的なページ構造

#### タスク 4.1: ポート設定（3000）
- **実施内容**: `package.json`の`scripts`セクションを修正
  - `"dev": "next dev -p 3000"`
  - `"start": "next start -p 3000"`
- **結果**: 開発サーバーとプロダクションサーバーがポート3000で起動することを確認

#### タスク 4.2: ルートレイアウトの作成
- **実施内容**: `app/layout.tsx`を確認・修正
  - `metadata.title`を`"Client2 App"`に変更
  - `metadata.description`を`"New client application"`に変更
  - `<html lang="ja">`に変更
  - `ClerkProvider`を削除（NextAuthを使用するため）
- **結果**: ルートレイアウトが適切に設定された

#### タスク 4.3: トップページの作成
- **実施内容**: `app/page.tsx`を確認（precedentテンプレートから既存のコンテンツが存在）
- **結果**: 基本的なページコンテンツが実装されていることを確認

### Phase 5: 動作確認とドキュメント

#### タスク 5.1: 開発サーバーの起動確認
- **実施内容**: `npm run dev`を実行
- **発生した問題と対応**:
  1. `border-border`クラスが存在しないエラー
     - **対応**: `tailwind.config.js`にshadcn/ui用のカラー設定を追加
  2. ClerkProviderエラー
     - **対応**: `app/layout.tsx`から`ClerkProvider`を削除
  3. Clerkコンポーネント（SignedOut, SignedIn等）のエラー
     - **対応**: `components/layout/navbar.tsx`からClerkコンポーネントを削除
  4. Clerk middlewareエラー
     - **対応**: `app/middleware.ts`を削除
- **結果**: 開発サーバーが正常に起動し、ポート3000で動作することを確認

#### タスク 5.2: ホットリロードの確認
- **実施内容**: ファイル変更時にブラウザが自動更新されることを確認
- **結果**: Next.jsのデフォルト機能によりホットリロードが正常に動作することを確認

#### タスク 5.3: TypeScript型チェックの確認
- **実施内容**: `npm run type-check`を実行（`package.json`に`type-check`スクリプトを追加）
- **発生した問題と対応**:
  - `class-variance-authority`が見つからないエラー
    - **対応**: `npm install class-variance-authority --legacy-peer-deps`でインストール
- **結果**: TypeScript型チェックが正常に動作し、型エラーがないことを確認

#### タスク 5.4: ビルドの確認
- **実施内容**: `npm run build`を実行
- **結果**: プロダクションビルドが正常に完了することを確認

#### タスク 5.5: 一時的なドキュメントの作成
- **実施内容**: `docs/Temp-Client2.md`を作成
- **内容**: プロジェクトの概要、セットアップ方法、環境変数の設定方法、技術スタック、利用可能なスクリプトを記載
- **結果**: 移行完了後にREADMEに移植する想定のドキュメントが作成された

### レビュー対応

#### Clerk依存の削除
- **指摘内容**: `package.json`に`@clerk/nextjs`が残存している
- **実施内容**: `npm uninstall @clerk/nextjs --legacy-peer-deps`を実行
- **結果**: Clerk関連の20パッケージを削除し、ビルド・開発サーバー・型チェックが正常に動作することを確認

## 技術スタック

- **フレームワーク**: Next.js 14+ (App Router)
- **言語**: TypeScript 5+
- **UIコンポーネント**: shadcn/ui
- **認証**: NextAuth (Auth.js) v5
- **スタイリング**: Tailwind CSS
- **フォーム管理**: react-hook-form
- **バリデーション**: zod

## ディレクトリ構造

```
client2/
├── app/
│   ├── api/
│   │   └── auth/
│   │       └── [...nextauth]/
│   │           └── route.ts
│   ├── layout.tsx
│   ├── page.tsx
│   └── globals.css
├── components/
│   ├── ui/          # shadcn/uiコンポーネント
│   └── layout/       # レイアウトコンポーネント
├── lib/
│   ├── hooks/        # カスタムフック
│   └── utils.ts
├── auth.ts           # NextAuth設定
├── components.json   # shadcn/ui設定
├── tailwind.config.js
├── tsconfig.json
├── next.config.js
├── package.json
├── .env.example
└── .env.local        # gitignoreに含まれる
```

## 重要な設定ファイル

### `components.json` (shadcn/ui設定)
```json
{
  "$schema": "https://ui.shadcn.com/schema.json",
  "style": "default",
  "rsc": true,
  "tsx": true,
  "tailwind": {
    "config": "tailwind.config.js",
    "css": "app/globals.css",
    "baseColor": "slate",
    "cssVariables": true
  },
  "aliases": {
    "components": "@/components",
    "utils": "@/lib/utils"
  }
}
```

### `tailwind.config.js` (重要な設定)
- `darkMode: ["class"]` - shadcn/ui用
- `tailwindcss-animate`プラグイン
- shadcn/ui用のカラー設定（colors.extend）
- `borderRadius`の設定（`var(--radius)`）

### `next.config.js` (重要な設定)
- `images.remotePatterns`を使用（`images.domains`は非推奨）

### `package.json` (重要なスクリプト)
- `"dev": "next dev -p 3000"`
- `"start": "next start -p 3000"`
- `"type-check": "tsc --noEmit"`

## 環境変数

### `.env.example`
```
# NextAuth (Auth.js)
AUTH_SECRET=your-secret-key-here
AUTH_URL=http://localhost:3000
```

### `.env.local` (gitignoreに含まれる)
```
# NextAuth (Auth.js)
AUTH_SECRET=<生成された秘密鍵>
AUTH_URL=http://localhost:3000
```

**注意**: `AUTH_SECRET`は`openssl rand -base64 32`で生成

## インストール済みshadcn/uiコンポーネント

1. `alert-dialog` - アラートダイアログ
2. `alert` - アラート
3. `button` - ボタン
4. `select` - セレクト
5. `input` - 入力フィールド
6. `form` - フォーム（react-hook-form統合）
7. `field` - フィールド（form依存）
8. `card` - カード
9. `label` - ラベル（form依存）
10. `separator` - セパレータ（field依存）

## 削除されたもの

### Prisma関連
- `prisma/`ディレクトリ
- `package.json`の`prisma`と`@prisma/client`依存関係
- Prisma関連のスクリプト
- `.env.example`の`DATABASE_URL`

### Clerk関連
- `@clerk/nextjs`パッケージ
- `app/layout.tsx`の`ClerkProvider`
- `components/layout/navbar.tsx`のClerkコンポーネント（SignedIn, SignedOut, SignInButton, UserButton）
- `app/middleware.ts`（Clerk middleware）

## 発生した問題と対応

### 1. npm install時のpeer dependency競合
- **問題**: `npm install`でERESOLVEエラー
- **対応**: `npm install --legacy-peer-deps`を使用

### 2. shadcn CLIのpnpm使用
- **問題**: shadcn CLIが`pnpm`を使用しようとした
- **対応**: `components.json`を手動作成

### 3. shadcn CLIのパス解決問題
- **問題**: shadcn CLIが`@/components/ui/`というリテラルパスにファイルを作成
- **対応**: 作成されたファイルを`components/ui/`に手動で移動

### 4. border-borderクラスエラー
- **問題**: `The 'border-border' class does not exist`
- **対応**: `tailwind.config.js`にshadcn/ui用のカラー設定を追加

### 5. Clerk関連エラー
- **問題**: ClerkProvider、Clerkコンポーネント、Clerk middlewareのエラー
- **対応**: Clerk関連のコードをすべて削除

### 6. class-variance-authorityが見つからない
- **問題**: TypeScript型チェックで`class-variance-authority`が見つからない
- **対応**: `npm install class-variance-authority --legacy-peer-deps`でインストール

### 7. images.domainsの非推奨警告
- **問題**: `images.domains`が非推奨
- **対応**: `images.remotePatterns`に置き換え

## 次のステップ（client1からの機能移植）

### 移植時に参照すべき情報

1. **認証**: NextAuth (Auth.js) v5を使用（client1はAuth0を使用）
2. **UIコンポーネント**: shadcn/uiコンポーネントを使用（既に8つのコンポーネントがインストール済み）
3. **フォーム**: react-hook-formを使用（shadcn/uiのformコンポーネントと統合）
4. **バリデーション**: zodを使用
5. **API呼び出し**: client1の`lib/api.ts`を参考に、client2用のAPIクライアントを作成

### 移植時の注意点

1. **認証の違い**: client1はAuth0、client2はNextAuth (Auth.js)を使用
   - 認証フローとAPI呼び出し時のトークン取得方法が異なる
   - `lib/auth.ts`の実装をNextAuth v5のAPIに合わせて調整が必要

2. **コンポーネントの違い**: client1は独自コンポーネント、client2はshadcn/uiを使用
   - 既存のコンポーネントをshadcn/uiコンポーネントに置き換える必要がある
   - スタイリングはTailwind CSSで統一されているため、比較的容易

3. **データフェッチング**: client1とclient2でAPI呼び出し方法が異なる可能性
   - Next.js 14のApp RouterではServer ComponentsとClient Componentsの使い分けが重要
   - `lib/api.ts`の実装をApp Routerに合わせて調整が必要

4. **型定義**: client1の`src/types/`をclient2に移植
   - `client/src/types/dm_post.ts`
   - `client/src/types/dm_user.ts`
   - その他の型定義ファイル

5. **ページ構造**: client1はPages Router、client2はApp Routerを使用
   - ページのルーティング構造が異なる
   - 動的ルートの実装方法が異なる

## 参考リンク

- [Next.js 14 App Router](https://nextjs.org/docs/app)
- [NextAuth (Auth.js) v5](https://authjs.dev/)
- [shadcn/ui](https://ui.shadcn.com/)
- [react-hook-form](https://react-hook-form.com/)
- [zod](https://zod.dev/)

## コミット履歴

主要なコミット:
1. `2264501` - fix(client2): 開発サーバー起動時のエラーと警告を修正
2. `c92741f` - feat(client2): TypeScript型チェックとビルド確認の実装
3. `bcebcc6` - docs(client2): 一時的なドキュメントを作成
4. `ea0038e` - fix(client2): Clerk依存を削除

## 備考

- すべてのタスクの受け入れ基準を満たしている
- ブラウザでの動作確認（表示、ホットリロード）を完了
- TypeScript型チェック、ビルド、開発サーバー起動が正常に動作することを確認
