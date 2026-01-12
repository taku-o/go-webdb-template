# 新clientアプリの作成の設計書

## 1. 概要

### 1.1 設計の目的
要件定義書に基づき、`client2`ディレクトリに新しいクライアントアプリケーションを作成するための詳細設計を定義する。precedentテンプレートをベースに、shadcn/uiとNextAuth (Auth.js)を統合したモダンなNext.jsアプリケーションの構築方針を明確にする。

### 1.2 設計の範囲
- `client2`ディレクトリの構造設計
- precedentテンプレートの統合設計
- shadcn/uiの統合設計
- NextAuth (Auth.js)の統合設計
- 開発環境の設定設計
- ポート設定（3000）の設計

### 1.3 設計方針
- **precedentテンプレートの標準構造を維持**: テンプレートの標準的な構造に従い、必要最小限のカスタマイズのみ行う
- **shadcn/uiの積極的利用**: shadcn/uiコンポーネントを積極的に利用し、UI開発を効率化する
- **NextAuth (Auth.js)の標準的な実装**: NextAuth (Auth.js)の標準的な実装パターンに従う
- **既存のclientディレクトリとの独立性**: 既存の`client`ディレクトリとは完全に独立した実装とする
- **ポート3000の使用**: 開発サーバーはポート3000で起動する

## 2. アーキテクチャ設計

### 2.1 全体構成

#### 2.1.1 アーキテクチャ概要

```
┌─────────────────────────────────────────┐
│         client2/ (新規アプリケーション)   │
│                                         │
│  ┌──────────────────────────────────┐  │
│  │   Next.js 14+ (App Router)        │  │
│  │   - app/                          │  │
│  │   - components/                   │  │
│  │   - lib/                          │  │
│  └──────────────────────────────────┘  │
│                                         │
│  ┌──────────────────────────────────┐  │
│  │   shadcn/ui                       │  │
│  │   - components/ui/                │  │
│  │   - Tailwind CSS                  │  │
│  │   - Radix UI                      │  │
│  └──────────────────────────────────┘  │
│                                         │
│  ┌──────────────────────────────────┐  │
│  │   NextAuth (Auth.js)              │  │
│  │   - app/api/auth/[...nextauth]/   │  │
│  │   - 認証プロバイダー設定          │  │
│  └──────────────────────────────────┘  │
│                                         │
│  ┌──────────────────────────────────┐  │
│  │   TypeScript 5+                  │  │
│  │   - 型安全性                      │  │
│  │   - 型定義                        │  │
│  └──────────────────────────────────┘  │
└─────────────────────────────────────────┘
```

#### 2.1.2 ディレクトリ構造

**client2ディレクトリ構造**:
```
client2/
├── app/                          # Next.js App Router
│   ├── layout.tsx                # ルートレイアウト
│   ├── page.tsx                  # トップページ
│   ├── api/                      # APIルート
│   │   └── auth/
│   │       └── [...nextauth]/
│   │           └── route.ts     # NextAuth (Auth.js)ルート
│   └── (その他のページ)
├── components/                   # Reactコンポーネント
│   ├── ui/                       # shadcn/uiコンポーネント
│   │   ├── button.tsx
│   │   ├── card.tsx
│   │   └── (その他のコンポーネント)
│   └── (カスタムコンポーネント)
├── lib/                          # ユーティリティ関数
│   ├── utils.ts                  # ユーティリティ関数（shadcn/ui用）
│   └── (その他のユーティリティ)
├── public/                       # 静的ファイル
│   └── (画像、フォントなど)
├── types/                        # TypeScript型定義
│   └── (型定義ファイル)
├── .env.example                  # 環境変数テンプレート
├── .env.local                    # ローカル環境変数（gitignore）
├── .gitignore                    # Git除外設定
├── components.json               # shadcn/ui設定
├── next.config.js                # Next.js設定
├── package.json                  # パッケージ設定
├── postcss.config.js             # PostCSS設定（Tailwind CSS用）
├── tailwind.config.js            # Tailwind CSS設定
├── tsconfig.json                 # TypeScript設定
├── prettier.config.js            # Prettier設定（precedentテンプレートに含まれる場合）
├── eslint.config.js              # ESLint設定（precedentテンプレートに含まれる場合）
└── README.md                     # プロジェクトドキュメント
```

**注意**: precedentテンプレートに含まれるPrisma関連のファイル（`prisma/`ディレクトリ、`schema.prisma`など）は、本要件では不要なため、削除または無効化する。

### 2.2 技術スタックの詳細設計

#### 2.2.1 Next.js 14+ (App Router)

**設定ファイル**: `next.config.js`

**主要な設定項目**:
- App Routerの有効化（デフォルトで有効）
- TypeScriptの統合
- 環境変数の設定
- ポート設定（3000）

**実装方針**:
- precedentテンプレートの標準的な`next.config.js`をベースにする
- 必要に応じてポート設定を追加（`-p 3000`または環境変数で指定）

#### 2.2.2 shadcn/ui

**設定ファイル**: `components.json`

**主要な設定項目**:
- スタイル設定（デフォルト、ダークモードなど）
- コンポーネントの配置先（`components/ui/`）
- Tailwind CSSの設定
- TypeScriptの設定

**実装方針**:
- `npx shadcn-ui@latest init`を実行して初期化
- 必要最小限のコンポーネントのみをインストール（Button、Cardなど）
- `components.json`を適切に設定

**インストールするコンポーネント（既存client移行用）**:
既存の`client`ディレクトリからの機能移行を考慮し、以下のコンポーネントをインストールします：

- `alert-dialog`: アラートダイアログコンポーネント
- `alert`: アラートコンポーネント
- `button`: ボタンコンポーネント
- `select`: セレクトボックスコンポーネント
- `input`: 入力フィールドコンポーネント
- `form`: フォームコンポーネント（react-hook-form統合）
- `field`: フィールドコンポーネント（ラベル、入力、エラー表示を統合）
- `card`: カードコンポーネント

**shadcn/uiで利用可能な主要なコンポーネント一覧**:
shadcn/uiには以下のような多くのコンポーネントが用意されています：

**基本コンポーネント**:
- `accordion`: アコーディオン
- `alert`: アラート
- `alert-dialog`: アラートダイアログ
- `avatar`: アバター
- `badge`: バッジ
- `button`: ボタン
- `card`: カード
- `checkbox`: チェックボックス
- `input`: 入力フィールド
- `label`: ラベル
- `radio-group`: ラジオボタングループ
- `select`: セレクトボックス
- `separator`: 区切り線
- `switch`: スイッチ
- `textarea`: テキストエリア
- `toggle`: トグル
- `toggle-group`: トグルグループ

**ナビゲーション・レイアウト**:
- `breadcrumb`: パンくずリスト
- `dropdown-menu`: ドロップダウンメニュー
- `menubar`: メニューバー
- `navigation-menu`: ナビゲーションメニュー
- `pagination`: ページネーション
- `sidebar`: サイドバー
- `tabs`: タブ

**オーバーレイ・ポップアップ**:
- `dialog`: ダイアログ
- `drawer`: ドロワー
- `hover-card`: ホバーカード
- `popover`: ポップオーバー
- `sheet`: シート
- `tooltip`: ツールチップ
- `context-menu`: コンテキストメニュー

**データ表示**:
- `table`: テーブル
- `data-table`: データテーブル（拡張版）
- `calendar`: カレンダー
- `chart`: チャート
- `progress`: プログレスバー
- `skeleton`: スケルトンローダー
- `empty`: 空状態表示

**フォーム・入力**:
- `combobox`: コンボボックス
- `command`: コマンドパレット
- `date-picker`: 日付ピッカー
- `input-otp`: OTP入力
- `slider`: スライダー

**その他**:
- `aspect-ratio`: アスペクト比
- `carousel`: カルーセル
- `resizable`: リサイズ可能パネル
- `scroll-area`: スクロールエリア
- `sonner`: トースト通知（Sonner）

**注意**: 本要件では、既存の`client`ディレクトリからの機能移行を考慮し、上記のコンポーネント（`alert-dialog`、`alert`、`button`、`select`、`input`、`form`、`field`、`card`）を初期セットアップ時にインストールします。他のコンポーネントは、実際の機能実装時に必要に応じて追加インストールします。

#### 2.2.3 NextAuth (Auth.js)

**設定ファイル**: `app/api/auth/[...nextauth]/route.ts`

**主要な設定項目**:
- 認証プロバイダーの設定（必要最小限）
- セッション管理の設定
- コールバックURLの設定
- 環境変数の設定（`AUTH_SECRET`、`AUTH_URL`など）

**実装方針**:
- NextAuth (Auth.js)の標準的な実装パターンに従う
- 必要最小限のプロバイダー（例: Credentials、Googleなど）を設定
- 環境変数を適切に設定

**環境変数**:
- `AUTH_SECRET`: 認証用の秘密鍵
- `AUTH_URL`: 認証URL（開発環境では`http://localhost:3000`）

#### 2.2.4 Tailwind CSS

**設定ファイル**: `tailwind.config.js`, `postcss.config.js`

**主要な設定項目**:
- コンテンツパスの設定（`app/`, `components/`など）
- テーマのカスタマイズ（shadcn/ui用）
- プラグインの設定

**実装方針**:
- shadcn/uiと互換性のあるTailwind CSS設定を使用
- precedentテンプレートの標準的な設定をベースにする

#### 2.2.5 TypeScript

**設定ファイル**: `tsconfig.json`

**主要な設定項目**:
- コンパイラオプション
- パスエイリアス（`@/components`, `@/lib`など）
- 型定義の設定

**実装方針**:
- Next.js 14+の標準的なTypeScript設定を使用
- パスエイリアスを適切に設定

## 3. 実装方針

### 3.1 precedentテンプレートの取得とセットアップ

#### 3.1.1 テンプレートの取得方法

**方法1: GitHubからクローン**
```bash
git clone https://github.com/steven-tey/precedent.git client2
cd client2
rm -rf .git
```

**方法2: degitを使用（推奨）**
```bash
npx degit steven-tey/precedent client2
cd client2
```

**方法3: テンプレートとして使用**
- GitHubの「Use this template」機能を使用

#### 3.1.2 初期セットアップ

1. **依存関係のインストール**
   ```bash
   npm install
   ```

2. **Prisma関連の削除（必須）**
   precedentテンプレートにはPrismaが含まれているが、本要件では不要なため削除する：
   
   **削除するファイル・ディレクトリ**:
   ```bash
   # prismaディレクトリを削除（存在する場合）
   rm -rf prisma
   ```
   
   **削除する依存関係**:
   ```bash
   npm uninstall prisma @prisma/client
   ```
   
   **package.jsonの確認と修正**:
   - `scripts`セクションからPrisma関連のスクリプトを削除
     - `prisma:generate`
     - `prisma:push`
     - `prisma:migrate`
     - `postinstall`（Prisma関連のコマンドが含まれている場合）
   - `devDependencies`または`dependencies`から`prisma`と`@prisma/client`を削除（`npm uninstall`で自動的に削除される）
   
   **環境変数の削除**:
   - `.env.example`から`DATABASE_URL`を削除（存在する場合）
   - `.env.local`から`DATABASE_URL`を削除（存在する場合）
   
   **コードの確認**:
   - Prismaを使用しているコード（`lib/db.ts`など）が存在する場合は削除または無効化
   - Prismaのインポート文が含まれているファイルを確認し、必要に応じて削除

3. **環境変数の設定**
   - `.env.example`を確認
   - `.env.local`を作成して必要な環境変数を設定

### 3.2 shadcn/uiの統合

#### 3.2.1 初期化

```bash
npx shadcn-ui@latest init
```

**設定項目**:
- スタイル: デフォルト（または好みのスタイル）
- ベースカラー: slate（または好みのカラー）
- CSS変数: Yes
- コンポーネントの配置先: `components/ui`

#### 3.2.2 コンポーネントのインストール

```bash
npx shadcn-ui@latest add button
npx shadcn-ui@latest add card
```

**注意**: 必要最小限のコンポーネントのみをインストールする。

#### 3.2.3 Tailwind CSS設定の確認

`tailwind.config.js`がshadcn/uiと互換性があることを確認する。

### 3.3 NextAuth (Auth.js)の統合

#### 3.3.1 インストール

precedentテンプレートには既にNextAuth (Auth.js)が含まれている可能性が高い。含まれていない場合は：

```bash
npm install next-auth@beta
```

#### 3.3.2 認証ルートの作成

`app/api/auth/[...nextauth]/route.ts`を作成または確認：

```typescript
import NextAuth from "next-auth"

const handler = NextAuth({
  // 認証設定
})

export { handler as GET, handler as POST }
```

#### 3.3.3 環境変数の設定

`.env.local`に以下を追加：

```
AUTH_SECRET=your-secret-key-here
AUTH_URL=http://localhost:3000
```

### 3.4 ポート設定

#### 3.4.1 package.jsonでの設定

`package.json`の`scripts`セクションを確認：

```json
{
  "scripts": {
    "dev": "next dev -p 3000",
    "build": "next build",
    "start": "next start -p 3000"
  }
}
```

#### 3.4.2 環境変数での設定（代替案）

`.env.local`に以下を追加：

```
PORT=3000
```

`next.config.js`でポートを読み込む（必要に応じて）。

### 3.5 基本的なページ構造の作成

#### 3.5.1 ルートレイアウト

`app/layout.tsx`を作成または確認：

```typescript
import type { Metadata } from "next"

export const metadata: Metadata = {
  title: "Client2 App",
  description: "New client application",
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="ja">
      <body>{children}</body>
    </html>
  )
}
```

#### 3.5.2 トップページ

`app/page.tsx`を作成または確認：

```typescript
export default function Home() {
  return (
    <main>
      <h1>Client2 App</h1>
    </main>
  )
}
```

## 4. 設定ファイルの詳細

### 4.1 next.config.js

```javascript
/** @type {import('next').NextConfig} */
const nextConfig = {
  // Next.js設定
}

module.exports = nextConfig
```

**主要な設定項目**:
- App Routerの有効化（デフォルト）
- TypeScriptの統合（デフォルト）
- 環境変数の設定（必要に応じて）

### 4.2 tsconfig.json

```json
{
  "compilerOptions": {
    "target": "ES2020",
    "lib": ["dom", "dom.iterable", "esnext"],
    "allowJs": true,
    "skipLibCheck": true,
    "strict": true,
    "noEmit": true,
    "esModuleInterop": true,
    "module": "esnext",
    "moduleResolution": "bundler",
    "resolveJsonModule": true,
    "isolatedModules": true,
    "jsx": "preserve",
    "incremental": true,
    "plugins": [
      {
        "name": "next"
      }
    ],
    "paths": {
      "@/*": ["./*"]
    }
  },
  "include": ["next-env.d.ts", "**/*.ts", "**/*.tsx", ".next/types/**/*.ts"],
  "exclude": ["node_modules"]
}
```

### 4.3 tailwind.config.js

```javascript
/** @type {import('tailwindcss').Config} */
module.exports = {
  darkMode: ["class"],
  content: [
    './pages/**/*.{ts,tsx}',
    './components/**/*.{ts,tsx}',
    './app/**/*.{ts,tsx}',
    './src/**/*.{ts,tsx}',
  ],
  theme: {
    // shadcn/ui用のテーマ設定
  },
  plugins: [require("tailwindcss-animate")],
}
```

### 4.4 components.json

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

### 4.5 package.json

**主要な依存関係**:
- `next`: Next.js 14+
- `react`, `react-dom`: React 18+
- `typescript`: TypeScript 5+
- `tailwindcss`: Tailwind CSS
- `next-auth`: NextAuth (Auth.js)
- `@radix-ui/*`: Radix UIコンポーネント（shadcn/uiの基盤）
- `class-variance-authority`: クラスバリアンス管理（shadcn/ui用）
- `clsx`: クラス名ユーティリティ（shadcn/ui用）
- `tailwind-merge`: Tailwind CSSクラスマージ（shadcn/ui用）
- `react-hook-form`: フォーム管理ライブラリ（`form`コンポーネント用）
- `@hookform/resolvers`: バリデーションリゾルバ（Zodなどと統合する場合）

**主要なスクリプト**:
- `dev`: 開発サーバー起動（ポート3000）
- `build`: プロダクションビルド
- `start`: プロダクションサーバー起動（ポート3000）
- `lint`: ESLint実行
- `type-check`: TypeScript型チェック

### 4.6 .env.example

```
# NextAuth (Auth.js)
AUTH_SECRET=your-secret-key-here
AUTH_URL=http://localhost:3000

# その他の環境変数（必要に応じて）
```

## 5. 統合方法の詳細

### 5.1 precedentテンプレートとshadcn/uiの統合

#### 5.1.1 競合の解決

precedentテンプレートには標準でshadcn/uiが含まれていないため、手動で統合する必要がある。

**手順**:
1. precedentテンプレートをセットアップ
2. shadcn/uiを初期化（`npx shadcn-ui@latest init`）
3. 必要最小限のコンポーネントをインストール
4. Tailwind CSS設定を確認・調整

#### 5.1.2 スタイルの統一

precedentテンプレートのスタイルとshadcn/uiのスタイルを統一する。

### 5.2 NextAuth (Auth.js)の統合

#### 5.2.1 precedentテンプレートとの統合

precedentテンプレートには標準でNextAuth (Auth.js)が含まれている可能性が高い。含まれている場合は、設定を確認・調整する。

#### 5.2.2 認証プロバイダーの設定

必要最小限のプロバイダー（例: Credentials、Googleなど）を設定する。

### 5.3 ポート設定の統合

#### 5.3.1 package.jsonでの設定

`package.json`の`scripts`セクションでポートを指定：

```json
{
  "scripts": {
    "dev": "next dev -p 3000"
  }
}
```

#### 5.3.2 環境変数での設定（代替案）

環境変数`PORT`を設定し、`next.config.js`で読み込む（必要に応じて）。

## 6. 実装手順

### 6.1 初期セットアップ

1. **precedentテンプレートの取得**
   ```bash
   npx degit steven-tey/precedent client2
   cd client2
   ```

2. **依存関係のインストール**
   ```bash
   npm install
   ```

3. **Prisma関連の削除（必須）**
   precedentテンプレートにはPrismaが含まれているが、本要件では不要なため削除する：
   
   **削除するファイル・ディレクトリ**:
   ```bash
   # prismaディレクトリを削除（存在する場合）
   rm -rf prisma
   ```
   
   **削除する依存関係**:
   ```bash
   npm uninstall prisma @prisma/client
   ```
   
   **package.jsonの確認と修正**:
   - `scripts`セクションからPrisma関連のスクリプトを削除
     - `prisma:generate`
     - `prisma:push`
     - `prisma:migrate`
     - `postinstall`（Prisma関連のコマンドが含まれている場合）
   - `devDependencies`または`dependencies`から`prisma`と`@prisma/client`を削除（`npm uninstall`で自動的に削除される）
   
   **環境変数の削除**:
   - `.env.example`から`DATABASE_URL`を削除（存在する場合）
   - `.env.local`から`DATABASE_URL`を削除（存在する場合）
   
   **コードの確認**:
   - Prismaを使用しているコード（`lib/db.ts`など）が存在する場合は削除または無効化
   - Prismaのインポート文が含まれているファイルを確認し、必要に応じて削除

### 6.2 shadcn/uiの統合

1. **shadcn/uiの初期化**
   ```bash
   npx shadcn-ui@latest init
   ```

2. **既存client移行用コンポーネントのインストール**
   ```bash
   npx shadcn-ui@latest add alert-dialog
   npx shadcn-ui@latest add alert
   npx shadcn-ui@latest add button
   npx shadcn-ui@latest add select
   npx shadcn-ui@latest add input
   npx shadcn-ui@latest add form
   npx shadcn-ui@latest add field
   npx shadcn-ui@latest add card
   ```
   
   **注意**: `form`コンポーネントは`react-hook-form`と統合されているため、`react-hook-form`がインストールされていることを確認してください。`field`コンポーネントは新しいコンポーネントで、フォームフィールドを統合的に扱うためのコンポーネントです。

3. **Tailwind CSS設定の確認**
   - `tailwind.config.js`を確認
   - shadcn/uiと互換性があることを確認

### 6.3 NextAuth (Auth.js)の設定

1. **NextAuth (Auth.js)のインストール確認**
   - `package.json`を確認
   - 含まれていない場合はインストール

2. **認証ルートの作成・確認**
   - `app/api/auth/[...nextauth]/route.ts`を作成または確認

3. **環境変数の設定**
   - `.env.example`を確認
   - `.env.local`を作成して必要な環境変数を設定

### 6.4 ポート設定

1. **package.jsonの確認**
   - `scripts`セクションでポート3000を指定

2. **動作確認**
   - `npm run dev`でポート3000で起動することを確認

### 6.5 基本的なページ構造の作成

1. **ルートレイアウトの確認**
   - `app/layout.tsx`を確認または作成

2. **トップページの確認**
   - `app/page.tsx`を確認または作成

3. **動作確認**
   - 開発サーバーを起動
   - ブラウザで`http://localhost:3000`にアクセス
   - 基本的なページが表示されることを確認

### 6.6 動作確認

1. **開発サーバーの起動**
   ```bash
   npm run dev
   ```

2. **ブラウザでの確認**
   - `http://localhost:3000`にアクセス
   - 基本的なページが表示されることを確認

3. **ホットリロードの確認**
   - ファイルを変更
   - ブラウザが自動的に更新されることを確認

4. **型チェックの確認**
   ```bash
   npm run type-check
   ```

5. **ビルドの確認**
   ```bash
   npm run build
   ```

## 7. 注意事項

### 7.1 precedentテンプレートのカスタマイズ

- precedentテンプレートの標準的な構造を維持する
- 必要最小限のカスタマイズのみ行う
- **Prismaの削除（必須）**: Prisma関連のファイル、依存関係、スクリプト、環境変数を削除する

### 7.2 shadcn/uiの統合

- shadcn/uiの標準的な実装パターンに従う
- 必要最小限のコンポーネントのみをインストールする
- `components.json`を適切に設定する

### 7.3 NextAuth (Auth.js)の設定

- NextAuth (Auth.js)の標準的な実装パターンに従う
- 環境変数を適切に設定する
- 必要最小限のプロバイダーのみを設定する

### 7.4 ポート設定

- ポート3000で起動することを確認する
- `package.json`の`scripts`セクションでポートを指定する

### 7.5 既存のclientディレクトリとの関係

- `client2`ディレクトリは既存の`client`ディレクトリとは完全に独立している
- 両方のディレクトリが同時に存在しても問題ないようにする
- 既存の`client`ディレクトリには影響を与えない

## 8. 参考情報

### 8.1 関連ドキュメント
- `.kiro/specs/0063-client2/requirements.md`: 要件定義書
- `.kiro/steering/structure.md`: ファイル組織とコードパターン
- `.kiro/steering/tech.md`: 技術スタックとアーキテクチャ

### 8.2 外部リソース
- **precedentテンプレート**: https://github.com/steven-tey/precedent
- **shadcn/ui**: https://ui.shadcn.com/
- **NextAuth (Auth.js)**: https://authjs.dev/
- **Next.js**: https://nextjs.org/
- **Tailwind CSS**: https://tailwindcss.com/

### 8.3 技術スタック
- **フレームワーク**: Next.js 14+ (App Router)
- **言語**: TypeScript 5+
- **UIコンポーネント**: shadcn/ui
- **認証**: NextAuth (Auth.js)
- **スタイリング**: Tailwind CSS
- **UIプリミティブ**: Radix UI（shadcn/uiの基盤）
