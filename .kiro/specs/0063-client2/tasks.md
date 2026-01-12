# 新clientアプリの作成の実装タスク一覧

## 概要
既存の`client`ディレクトリとは別に、新しいクライアントアプリケーションを`client2`ディレクトリに作成するためのタスク一覧。precedentテンプレートをベースに、shadcn/uiとNextAuth (Auth.js)を統合したモダンなNext.jsアプリケーションを構築する。要件定義書と設計書に基づき、実装可能な粒度でタスクを分解。

## 実装フェーズ

### Phase 1: precedentテンプレートの取得と初期セットアップ

#### タスク 1.1: precedentテンプレートの取得
**目的**: precedentテンプレートを`client2`ディレクトリに取得する。

**作業内容**:
- `npx degit`を使用してprecedentテンプレートを取得
- `client2`ディレクトリに配置

**実装内容**:
```bash
npx degit steven-tey/precedent client2
cd client2
```

**受け入れ基準**:
- [ ] `client2/`ディレクトリが作成されている
- [ ] precedentテンプレートのファイルが`client2/`ディレクトリに配置されている
- [ ] `.git`ディレクトリが削除されている（degitを使用した場合、自動的に削除される）

- _Requirements: 3.1.1_
- _Design: 3.1.1, 6.1_

---

#### タスク 1.2: 依存関係のインストール
**目的**: precedentテンプレートの依存関係をインストールする。

**作業内容**:
- `package.json`の依存関係をインストール

**実装内容**:
```bash
cd client2
npm install
```

**受け入れ基準**:
- [ ] `npm install`が正常に完了している
- [ ] `node_modules/`ディレクトリが作成されている
- [ ] `package-lock.json`が生成されている（npmを使用している場合）

- _Requirements: 3.3.1_
- _Design: 6.1_

---

#### タスク 1.3: Prisma関連の削除
**目的**: precedentテンプレートに含まれるPrisma関連のファイル、依存関係、スクリプト、環境変数を削除する。

**作業内容**:
- `prisma/`ディレクトリの削除
- Prisma関連の依存関係の削除
- `package.json`からPrisma関連のスクリプトを削除
- 環境変数から`DATABASE_URL`を削除
- Prismaを使用しているコードの確認と削除

**実装内容**:
```bash
# prismaディレクトリを削除（存在する場合）
rm -rf prisma

# Prisma関連の依存関係を削除
npm uninstall prisma @prisma/client
```

**package.jsonの修正**:
- `scripts`セクションから以下を削除:
  - `prisma:generate`
  - `prisma:push`
  - `prisma:migrate`
  - `postinstall`（Prisma関連のコマンドが含まれている場合）

**環境変数の削除**:
- `.env.example`から`DATABASE_URL`を削除（存在する場合）
- `.env.local`から`DATABASE_URL`を削除（存在する場合）

**コードの確認**:
- Prismaを使用しているコード（`lib/db.ts`など）が存在する場合は削除または無効化
- Prismaのインポート文が含まれているファイルを確認し、必要に応じて削除

**受け入れ基準**:
- [ ] `prisma/`ディレクトリが削除されている
- [ ] `package.json`から`prisma`と`@prisma/client`が削除されている
- [ ] `package.json`の`scripts`セクションからPrisma関連のスクリプトが削除されている
- [ ] `.env.example`から`DATABASE_URL`が削除されている（存在していた場合）
- [ ] Prismaを使用しているコードが削除または無効化されている

- _Requirements: 8.1_
- _Design: 3.1.2, 6.1, 7.1_

---

### Phase 2: shadcn/uiの統合

#### タスク 2.1: shadcn/uiの初期化
**目的**: shadcn/uiを初期化し、プロジェクトに統合する。

**作業内容**:
- `npx shadcn-ui@latest init`を実行
- `components.json`の設定
- Tailwind CSS設定の確認

**実装内容**:
```bash
npx shadcn-ui@latest init
```

**設定項目**:
- スタイル: デフォルト（または好みのスタイル）
- ベースカラー: slate（または好みのカラー）
- CSS変数: Yes
- コンポーネントの配置先: `components/ui`

**受け入れ基準**:
- [ ] `npx shadcn-ui@latest init`が正常に完了している
- [ ] `components.json`が作成されている
- [ ] `components.json`が適切に設定されている
- [ ] `tailwind.config.js`がshadcn/uiと互換性がある設定になっている

- _Requirements: 3.2.2_
- _Design: 3.2.1, 6.2_

---

#### タスク 2.2: 既存client移行用コンポーネントのインストール
**目的**: 既存の`client`ディレクトリからの機能移行を考慮し、必要なshadcn/uiコンポーネントをインストールする。

**作業内容**:
- 以下のコンポーネントをインストール:
  - `alert-dialog`
  - `alert`
  - `button`
  - `select`
  - `input`
  - `form`
  - `field`
  - `card`

**実装内容**:
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

**注意**: `form`コンポーネントは`react-hook-form`と統合されているため、`react-hook-form`がインストールされていることを確認する。含まれていない場合はインストールする。

**受け入れ基準**:
- [ ] 8つのコンポーネント（`alert-dialog`, `alert`, `button`, `select`, `input`, `form`, `field`, `card`）がインストールされている
- [ ] `components/ui/`ディレクトリに各コンポーネントのファイルが作成されている
- [ ] `react-hook-form`がインストールされている（`form`コンポーネント用）
- [ ] 各コンポーネントが正常にインポートできることを確認

- _Requirements: 3.2.2_
- _Design: 3.2.2, 6.2_

---

### Phase 3: NextAuth (Auth.js)の統合

#### タスク 3.1: NextAuth (Auth.js)のインストール確認
**目的**: NextAuth (Auth.js)がインストールされていることを確認し、含まれていない場合はインストールする。

**作業内容**:
- `package.json`を確認してNextAuth (Auth.js)が含まれているか確認
- 含まれていない場合はインストール

**実装内容**:
```bash
# package.jsonを確認
# next-authが含まれていない場合
npm install next-auth@beta
```

**受け入れ基準**:
- [ ] `package.json`に`next-auth`が含まれている
- [ ] `node_modules/`に`next-auth`がインストールされている

- _Requirements: 3.2.3_
- _Design: 3.3.1, 6.3_

---

#### タスク 3.2: NextAuth (Auth.js)認証ルートの作成
**目的**: NextAuth (Auth.js)の認証ルートを作成または確認する。

**作業内容**:
- `app/api/auth/[...nextauth]/route.ts`を作成または確認
- 基本的な認証設定を実装

**実装内容**:
`app/api/auth/[...nextauth]/route.ts`を作成または確認：

```typescript
import NextAuth from "next-auth"

const handler = NextAuth({
  // 認証設定
  providers: [
    // 必要最小限のプロバイダーを設定
  ],
})

export { handler as GET, handler as POST }
```

**受け入れ基準**:
- [ ] `app/api/auth/[...nextauth]/route.ts`が存在している
- [ ] NextAuth (Auth.js)の基本的な設定が実装されている
- [ ] 認証ルートが正常に動作する（`/api/auth/signin`などにアクセス可能）

- _Requirements: 3.2.3_
- _Design: 3.3.2, 6.3_

---

#### タスク 3.3: 環境変数の設定
**目的**: NextAuth (Auth.js)に必要な環境変数を設定する。

**作業内容**:
- `.env.example`に必要な環境変数を追加
- `.env.local`を作成して必要な環境変数を設定

**実装内容**:
`.env.example`に以下を追加：
```
# NextAuth (Auth.js)
AUTH_SECRET=your-secret-key-here
AUTH_URL=http://localhost:3000
```

`.env.local`を作成（gitignoreに含まれる）：
```
AUTH_SECRET=your-secret-key-here
AUTH_URL=http://localhost:3000
```

**注意**: `AUTH_SECRET`は適切な秘密鍵を生成する必要があります。開発環境では`openssl rand -base64 32`などで生成できます。

**受け入れ基準**:
- [ ] `.env.example`に`AUTH_SECRET`と`AUTH_URL`が記載されている
- [ ] `.env.local`が作成されている
- [ ] `.env.local`に`AUTH_SECRET`と`AUTH_URL`が設定されている
- [ ] `.env.local`が`.gitignore`に含まれている

- _Requirements: 3.3.2, 8.4_
- _Design: 3.3.3, 6.3_

---

### Phase 4: ポート設定と基本的なページ構造

#### タスク 4.1: ポート設定（3000）
**目的**: 開発サーバーがポート3000で起動するように設定する。

**作業内容**:
- `package.json`の`scripts`セクションでポートを指定

**実装内容**:
`package.json`の`scripts`セクションを確認・修正：

```json
{
  "scripts": {
    "dev": "next dev -p 3000",
    "build": "next build",
    "start": "next start -p 3000"
  }
}
```

**受け入れ基準**:
- [ ] `package.json`の`scripts`セクションに`-p 3000`が指定されている
- [ ] `npm run dev`でポート3000で起動することを確認

- _Requirements: 5.3, 8.5_
- _Design: 3.4, 5.3, 6.1_

---

#### タスク 4.2: ルートレイアウトの作成
**目的**: Next.js App Routerのルートレイアウトを作成または確認する。

**作業内容**:
- `app/layout.tsx`を作成または確認
- 基本的なメタデータを設定

**実装内容**:
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

**受け入れ基準**:
- [ ] `app/layout.tsx`が存在している
- [ ] メタデータが適切に設定されている
- [ ] 基本的なHTML構造が実装されている

- _Requirements: 3.1.2_
- _Design: 3.5.1, 6.5_

---

#### タスク 4.3: トップページの作成
**目的**: 基本的なトップページを作成または確認する。

**作業内容**:
- `app/page.tsx`を作成または確認
- 基本的なページコンテンツを実装

**実装内容**:
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

**受け入れ基準**:
- [ ] `app/page.tsx`が存在している
- [ ] 基本的なページコンテンツが実装されている
- [ ] ブラウザで表示されることを確認

- _Requirements: 3.1.2_
- _Design: 3.5.2, 6.5_

---

### Phase 5: 動作確認とドキュメント

#### タスク 5.1: 開発サーバーの起動確認
**目的**: 開発サーバーが正常に起動することを確認する。

**作業内容**:
- `npm run dev`を実行
- ポート3000で起動することを確認
- ブラウザで`http://localhost:3000`にアクセスしてページが表示されることを確認

**実装内容**:
```bash
cd client2
npm run dev
```

**受け入れ基準**:
- [ ] `npm run dev`が正常に実行される
- [ ] 開発サーバーがポート3000で起動する
- [ ] ブラウザで`http://localhost:3000`にアクセスしてページが表示される
- [ ] エラーが発生しない

- _Requirements: 6.4_
- _Design: 6.6_

---

#### タスク 5.2: ホットリロードの確認
**目的**: ファイル変更時のホットリロードが正常に動作することを確認する。

**作業内容**:
- ファイルを変更（例: `app/page.tsx`のテキストを変更）
- ブラウザが自動的に更新されることを確認

**実装内容**:
1. `app/page.tsx`のテキストを変更
2. ブラウザが自動的に更新されることを確認

**受け入れ基準**:
- [ ] ファイルを変更するとブラウザが自動的に更新される
- [ ] エラーが発生しない

- _Requirements: 6.4_
- _Design: 6.6_

---

#### タスク 5.3: TypeScript型チェックの確認
**目的**: TypeScriptの型チェックが正常に動作することを確認する。

**作業内容**:
- TypeScriptの型チェックを実行
- エラーがないことを確認

**実装内容**:
```bash
npm run type-check
# または
npx tsc --noEmit
```

**受け入れ基準**:
- [ ] TypeScriptの型チェックが正常に実行される
- [ ] 型エラーが発生しない

- _Requirements: 6.4_
- _Design: 6.6_

---

#### タスク 5.4: ビルドの確認
**目的**: プロダクションビルドが正常に完了することを確認する。

**作業内容**:
- `npm run build`を実行
- ビルドが正常に完了することを確認

**実装内容**:
```bash
npm run build
```

**受け入れ基準**:
- [ ] `npm run build`が正常に実行される
- [ ] ビルドが正常に完了する
- [ ] ビルドエラーが発生しない

- _Requirements: 6.4_
- _Design: 6.6_

---

#### タスク 5.5: 一時的なドキュメントの作成
**目的**: 移行完了後にREADMEに移植するための一時的なドキュメントを作成する。

**作業内容**:
- `docs/Temp-Client2.md`を作成
- プロジェクトの概要、セットアップ方法、環境変数の設定方法を記載
- 移行完了後、この内容をREADMEに移植する想定

**実装内容**:
`docs/Temp-Client2.md`を作成：

```markdown
# Client2 App

新clientアプリケーション（precedentテンプレートベース）

## セットアップ

1. 依存関係のインストール
   ```bash
   npm install
   ```

2. 環境変数の設定
   `.env.local`を作成して以下の環境変数を設定：
   ```
   AUTH_SECRET=your-secret-key-here
   AUTH_URL=http://localhost:3000
   ```

3. 開発サーバーの起動
   ```bash
   npm run dev
   ```

## 技術スタック

- Next.js 14+ (App Router)
- TypeScript 5+
- shadcn/ui
- NextAuth (Auth.js)
- Tailwind CSS
```

**注意**: このドキュメントは一時的なもので、`client`から`client2`への移行が完了したら、この内容をREADMEに移植する想定です。

**受け入れ基準**:
- [ ] `docs/Temp-Client2.md`が作成されている
- [ ] プロジェクトの概要が記載されている
- [ ] セットアップ方法が記載されている
- [ ] 環境変数の設定方法が記載されている

- _Requirements: 6.5_
- _Design: 6.6_

---

## 受け入れ基準（全体）

### プロジェクト構造
- [ ] `client2/`ディレクトリが作成されている
- [ ] precedentテンプレートをベースとした初期セットアップが完了している
- [ ] 基本的なディレクトリ構造が整備されている（`app/`, `components/`, `lib/`, `public/`など）

### 技術スタックの統合
- [ ] Next.js 14+ (App Router)がセットアップされている
- [ ] shadcn/uiが統合されている（初期化が完了している）
- [ ] 8つのコンポーネント（`alert-dialog`, `alert`, `button`, `select`, `input`, `form`, `field`, `card`）がインストールされている
- [ ] NextAuth (Auth.js)が統合されている（基本的な設定が完了している）
- [ ] Tailwind CSSが設定されている
- [ ] TypeScriptが設定されている

### Prismaの削除
- [ ] Prisma関連のファイル、依存関係、スクリプト、環境変数が削除されている

### 基本的なプロジェクト設定
- [ ] `package.json`が作成され、必要な依存関係がインストールされている
- [ ] `.env.example`が作成されている
- [ ] 開発サーバーが正常に起動する（`npm run dev`が動作する）
- [ ] 開発サーバーがポート3000で起動する
- [ ] 基本的なページ（トップページなど）が表示される

### 動作確認
- [ ] ローカル環境で開発サーバーが正常に起動する（ポート3000で起動する）
- [ ] ブラウザで基本的なページが表示される（`http://localhost:3000`でアクセス可能）
- [ ] ホットリロードが正常に動作する
- [ ] TypeScriptの型チェックが正常に動作する
- [ ] ビルドが正常に完了する（`npm run build`が動作する）

### ドキュメント
- [ ] `docs/Temp-Client2.md`が作成されている
- [ ] プロジェクトの概要、セットアップ方法、環境変数の設定方法が記載されている
- [ ] 移行完了後にREADMEに移植する想定であることが明確になっている

- _Requirements: 6.1, 6.2, 6.3, 6.4, 6.5_
- _Design: 6.1, 6.2, 6.3, 6.4, 6.5, 6.6_
