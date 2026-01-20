# ヘッダーのログアウトボタン統一設計書

## 1. 概要

### 1.1 設計の目的
要件定義書に基づき、ヘッダーのログアウトボタンと画面内のログアウトボタンの実装を統一するための詳細設計を定義する。`AuthButtons`コンポーネントを直接使用することで、完全に同じ実装を実現する。

### 1.2 設計の範囲
- `client/components/layout/navbar.tsx`の修正
  - ユーザー情報表示部分とログアウト/ログインボタンを`AuthButtons`コンポーネントに置き換え
  - 不要なインポートの削除

### 1.3 設計方針
- **コンポーネントの再利用**: 既存の`AuthButtons`コンポーネントを直接使用
- **完全な統一**: 画面内のログアウトボタンと全く同じ実装にする
- **UIデザインの許容**: デザインが崩れる可能性があるが、完全に同じ実装になることを優先
- **既存機能の維持**: 既存の認証機能を完全に維持
- **最小限の変更**: 必要な変更のみを行い、既存コードへの影響を最小化

## 2. コンポーネント設計

### 2.1 navbar.tsxの修正設計

#### 2.1.1 変更前の構造

```tsx
import { auth, signIn, signOut } from "@/auth";
import { Button } from "@/components/ui/button";

export default async function NavBar() {
  const session = await auth();

  return (
    <div className="flex items-center gap-4">
      {session?.user ? (
        <>
          <div className="flex items-center gap-2">
            <span className="text-sm font-semibold">{session.user.name}</span>
            {session.user.email && (
              <span className="text-xs text-gray-600">{session.user.email}</span>
            )}
          </div>
          <form action={async () => {
            "use server"
            await signOut()
          }}>
            <Button type="submit" variant="outline">ログアウト</Button>
          </form>
        </>
      ) : (
        <form action={async () => {
          "use server"
          await signIn('auth0')
        }}>
          <Button type="submit" aria-label="ナビゲーションバーからログイン">ログイン</Button>
        </form>
      )}
    </div>
  );
}
```

#### 2.1.2 変更後の構造

```tsx
import { auth } from "@/auth";
import { AuthButtons } from "@/components/auth/auth-buttons";

export default async function NavBar() {
  const session = await auth();

  return (
    <div className="flex items-center gap-4">
      <AuthButtons user={session?.user ?? null} />
    </div>
  );
}
```

#### 2.1.3 変更内容の詳細

**インポートの変更**:
- **追加**: `AuthButtons`を`@/components/auth/auth-buttons`からインポート
- **削除**: `signIn`、`signOut`を`@/auth`からのインポートから削除
- **維持**: `auth`は`@/auth`からインポート（セッション取得に使用）

**コンポーネントの置き換え**:
- **削除**: ユーザー情報表示部分（`<div className="flex items-center gap-2">...</div>`）
- **削除**: ログアウトボタンのform（`<form action={async () => { "use server"; await signOut() }}>...</form>`）
- **削除**: ログインボタンのform（`<form action={async () => { "use server"; await signIn('auth0') }}>...</form>`）
- **追加**: `<AuthButtons user={session?.user ?? null} />`

**プロップの渡し方**:
- `session?.user ?? null`を`AuthButtons`コンポーネントに渡す
- `session?.user`が`undefined`の場合は`null`を渡す（`AuthButtons`コンポーネントの型定義に合わせる）

### 2.2 AuthButtonsコンポーネントの仕様

#### 2.2.1 コンポーネントのインターフェース

```tsx
interface AuthButtonsProps {
  user: {
    name?: string | null
    email?: string | null
  } | null
}
```

#### 2.2.2 コンポーネントの動作

**ログイン済みの場合**:
- ユーザー情報を表示（「ログイン中: {user.name}」形式）
- メールアドレスがあれば表示
- ログアウトボタンを表示（`variant="destructive" size="sm"`）
- `signOutAction`を使用してログアウト処理を実行

**未ログインの場合**:
- 「ログインしていません」を表示
- ログインボタンを表示（`size="sm"`）
- `signInAction`を使用してログイン処理を実行

#### 2.2.3 UIデザイン

**レイアウト**:
- レスポンシブレイアウト（`flex-col sm:flex-row`）
- 小画面では縦並び、大画面では横並び

**ボタンスタイル**:
- ログアウトボタン: `variant="destructive" size="sm"`
- ログインボタン: `size="sm"`（デフォルトvariant）
- ボタン幅: `w-full sm:w-auto`（小画面では全幅、大画面では自動幅）

**ユーザー情報表示**:
- フォントサイズ: `text-sm sm:text-base`
- メールアドレス: `text-xs sm:text-sm text-muted-foreground`

### 2.3 既存コンポーネントとの関係

#### 2.3.1 AuthButtonsコンポーネント
- **変更**: なし（既存の実装をそのまま使用）
- **場所**: `client/components/auth/auth-buttons.tsx`

#### 2.3.2 auth-actions.ts
- **変更**: なし（既存の実装をそのまま使用）
- **場所**: `client/lib/actions/auth-actions.ts`
- **使用関数**: `signOutAction`、`signInAction`

## 3. インポート設計

### 3.1 変更前のインポート

```tsx
import { auth, signIn, signOut } from "@/auth";
import { Button } from "@/components/ui/button";
```

### 3.2 変更後のインポート

```tsx
import { auth } from "@/auth";
import { AuthButtons } from "@/components/auth/auth-buttons";
```

### 3.3 インポートの変更理由

- **`auth`**: セッション取得に使用するため維持
- **`signIn`**: `AuthButtons`コンポーネント内で`signInAction`を使用するため不要
- **`signOut`**: `AuthButtons`コンポーネント内で`signOutAction`を使用するため不要
- **`Button`**: `AuthButtons`コンポーネント内で使用されるため不要
- **`AuthButtons`**: 新規追加（使用するコンポーネント）

## 4. UIデザインの変更

### 4.1 変更前のUIデザイン

**ログイン済み時**:
- ユーザー情報: インライン表示（横並び）
  - 名前: `text-sm font-semibold`
  - メール: `text-xs text-gray-600`
- ログアウトボタン: `variant="outline"`（デフォルトサイズ）
- レイアウト: `flex items-center gap-4`（常に横並び）

**未ログイン時**:
- ログインボタン: デフォルトvariant（デフォルトサイズ）

### 4.2 変更後のUIデザイン

**ログイン済み時**:
- ユーザー情報: ブロック表示（縦並び）
  - 「ログイン中: {user.name}」形式
  - フォントサイズ: `text-sm sm:text-base`
  - メール: `text-xs sm:text-sm text-muted-foreground`
- ログアウトボタン: `variant="destructive" size="sm"`
- レイアウト: `flex-col sm:flex-row`（レスポンシブ）

**未ログイン時**:
- 「ログインしていません」を表示
- ログインボタン: `size="sm"`（デフォルトvariant）
- レイアウト: `flex-col sm:flex-row`（レスポンシブ）

### 4.3 UIデザイン変更の影響

**ヘッダーのレイアウト**:
- 小画面では縦並びになる可能性がある
- 大画面では横並びを維持

**ボタンの見た目**:
- ログアウトボタンが赤色（`destructive`）になる
- ボタンサイズが小さくなる（`sm`）

**ユーザー情報の表示**:
- 「ログイン中: {user.name}」形式に変更
- メールアドレスの表示位置が変わる

**注意事項**:
- UIデザインが変わる可能性があるが、これは許容する
- 完全に同じ実装になることで、Issue #166の要件を満たす

## 5. データフロー設計

### 5.1 セッション取得から表示までの流れ

```
1. NavBarコンポーネントがレンダリングされる
   ↓
2. `auth()`でセッション情報を取得
   ↓
3. `session?.user ?? null`を`AuthButtons`コンポーネントに渡す
   ↓
4. `AuthButtons`コンポーネントがユーザー情報とボタンを表示
   ↓
5. ユーザーがログアウト/ログインボタンをクリック
   ↓
6. `signOutAction`または`signInAction`が実行される
   ↓
7. `@/auth`の`signOut()`または`signIn()`が実行される
   ↓
8. 認証状態が更新され、ページがリロードされる
```

### 5.2 プロップの型定義

```tsx
// navbar.tsx
const session = await auth();
// session?.user の型: { name?: string | null; email?: string | null } | undefined

// AuthButtonsに渡す値
<AuthButtons user={session?.user ?? null} />
// user の型: { name?: string | null; email?: string | null } | null
```

## 6. エラーハンドリング設計

### 6.1 セッション取得エラー

**想定されるエラー**:
- `auth()`の実行時にエラーが発生する可能性

**対応**:
- 既存の実装と同様に、エラーが発生した場合は`session`が`undefined`になる
- `session?.user ?? null`により、`null`が`AuthButtons`に渡される
- `AuthButtons`コンポーネントは`user`が`null`の場合、未ログイン状態として表示する

### 6.2 認証アクションエラー

**想定されるエラー**:
- `signOutAction`または`signInAction`の実行時にエラーが発生する可能性

**対応**:
- 既存の`AuthButtons`コンポーネントの実装に従う
- `auth-actions.ts`のエラーハンドリングに依存
- Next.jsのServer Actionsのエラーハンドリングに従う

## 7. テスト設計

### 7.1 既存テストへの影響

**想定される影響**:
- navbar.tsxのUI構造に依存するテストが失敗する可能性
- `signOut`や`signIn`を直接モックしているテストが失敗する可能性

**対応方針**:
- 既存のテストを確認し、必要に応じて修正する
- `AuthButtons`コンポーネントをモックする必要がある場合は、モックを追加する

### 7.2 新規テストの検討

**単体テスト**:
- navbar.tsxが`AuthButtons`コンポーネントを正しく使用していることを確認
- セッション情報が正しく`AuthButtons`に渡されていることを確認

**統合テスト**:
- ログアウト機能が正常に動作することを確認
- ログイン機能が正常に動作することを確認（未ログイン時）
- UIデザインが正しく表示されることを確認

**E2Eテスト（Playwright）**:
- ヘッダーのログアウトボタンが正常に動作することを確認
- ログアウト後に適切なページにリダイレクトされることを確認

### 7.3 テストの実装方針

**既存テストの修正**:
- navbar.tsxのUI構造に依存するテストを修正
- `signOut`や`signIn`を直接モックしているテストを修正

**新規テストの追加**:
- 要件定義書の受け入れ基準に基づいてテストを追加
- 既存の`AuthButtons`コンポーネントのテストを参考にする

## 8. 実装上の注意事項

### 8.1 インポートの整理

- `signIn`と`signOut`のインポートを削除する際、他で使用されていないことを確認する
- `Button`のインポートを削除する際、他で使用されていないことを確認する

### 8.2 型の整合性

- `session?.user ?? null`の型が`AuthButtons`コンポーネントの`user`プロップの型と一致することを確認する
- TypeScriptの型チェックでエラーが発生しないことを確認する

### 8.3 UIデザインの確認

- 変更後のUIデザインが意図通りに表示されることを確認する
- レスポンシブレイアウトが正常に動作することを確認する
- デザインが崩れている場合は、要件定義書の「UIデザインの調整は範囲外」に従い、許容する

### 8.4 既存機能の確認

- ログアウト機能が正常に動作することを確認する
- ログイン機能が正常に動作することを確認する（未ログイン時）
- 既存のテストが正常に動作することを確認する

## 9. 参考情報

### 9.1 関連ファイル

- `client/components/layout/navbar.tsx`: 修正対象ファイル
- `client/components/auth/auth-buttons.tsx`: 使用するコンポーネント
- `client/lib/actions/auth-actions.ts`: 認証アクション（`AuthButtons`内で使用）
- `client/app/page.tsx`: `AuthButtons`の使用例

### 9.2 技術スタック

- **Next.js**: Server Actionsを使用
- **TypeScript**: 型安全性を維持
- **React**: コンポーネントベースの実装
- **Tailwind CSS**: UIスタイリング

### 9.3 参考リンク

- Next.js Server Actions: https://nextjs.org/docs/app/building-your-application/data-fetching/server-actions-and-mutations
- Next.js Authentication: https://nextjs.org/docs/app/building-your-application/authentication
