# ヘッダーのログアウトボタン統一要件定義書

## 1. 概要

### 1.1 プロジェクト情報
- **プロジェクト名**: go-webdb-template
- **Feature名**: 0081-header-logout
- **作成日**: 2026-01-17
- **関連Issue**: https://github.com/taku-o/go-webdb-template/issues/166

### 1.2 目的
ヘッダーのログアウトボタンと画面内のログアウトボタンの実装を統一し、一貫性のある認証処理を実現する。

### 1.3 スコープ
- `client/components/layout/navbar.tsx`のログアウトボタン実装を修正
- 画面内のログアウトボタン（`client/components/auth/auth-buttons.tsx`）と同じ実装方法に統一
- **実装方法**: `client/components/auth/auth-buttons.tsx`の`AuthButtons`コンポーネントを直接使用
  - 完全に同じ実装になる（画面内のログアウトボタンと全く同じ）
  - UIデザインが変わる可能性がある（`variant="destructive" size="sm"`、レスポンシブレイアウト）

**本実装の範囲外**:
- ログイン機能の変更
- 認証ロジックの変更
- その他の認証関連コンポーネントの変更
- UIデザインの調整（デザインが崩れる可能性は許容する）

## 2. 背景・現状分析

### 2.1 現在の状況

#### 2.1.1 ヘッダーのログアウトボタン（navbar.tsx）
- **ファイル**: `client/components/layout/navbar.tsx`
- **実装方法**: 直接`signOut()`を呼び出している（34-39行目）
- **コード**:
  ```tsx
  <form action={async () => {
    "use server"
    await signOut()
  }}>
    <Button type="submit" variant="outline">ログアウト</Button>
  </form>
  ```
- **問題点**: `auth-actions.ts`の`signOutAction`を使用していない

#### 2.1.2 画面内のログアウトボタン（auth-buttons.tsx）
- **ファイル**: `client/components/auth/auth-buttons.tsx`
- **実装方法**: `signOutAction`を使用している（26行目）
- **コード**:
  ```tsx
  <form action={signOutAction}>
    <Button type="submit" variant="destructive" size="sm" className="w-full sm:w-auto" aria-label="ログアウト">
      ログアウト
    </Button>
  </form>
  ```
- **使用場所**: `client/app/page.tsx`の92行目で使用

#### 2.1.3 認証アクション（auth-actions.ts）
- **ファイル**: `client/lib/actions/auth-actions.ts`
- **実装**: `signOutAction`が定義されている（9-11行目）
- **コード**:
  ```ts
  export async function signOutAction() {
    await signOut()
  }
  ```

### 2.2 課題点
1. **実装の不統一**: ヘッダーのログアウトボタンと画面内のログアウトボタンで異なる実装方法を使用している
2. **保守性の低下**: 認証処理の変更時に複数の箇所を修正する必要がある
3. **一貫性の欠如**: 同じ機能なのに異なる実装方法を使用しているため、コードの一貫性が損なわれている

### 2.3 本実装による改善点
1. **実装の統一**: すべてのログアウトボタンで`signOutAction`を使用するようになる
2. **保守性の向上**: 認証処理の変更時に`auth-actions.ts`のみを修正すればよくなる
3. **一貫性の確保**: 同じ機能で同じ実装方法を使用するようになる

## 3. 機能要件

### 3.1 ヘッダーのログアウトボタンの修正

#### 3.1.1 navbar.tsxの修正

**概要**: `AuthButtons`コンポーネントを直接使用して、完全に同じ実装にする。

**変更内容**:
- navbar.tsxのユーザー情報表示部分とログアウト/ログインボタンを`AuthButtons`コンポーネントに置き換え
- **変更前**:
  ```tsx
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
  ```
- **変更後**:
  ```tsx
  <AuthButtons user={session?.user ?? null} />
  ```
- **インポート追加**: `AuthButtons`を`@/components/auth/auth-buttons`からインポート
- **インポート削除**: `signOut`、`signIn`を`@/auth`からのインポートから削除（他で使用されていない場合）

#### 3.1.2 UIデザインの変更

**変更点**:
- ボタンのスタイル: `variant="outline"` → `variant="destructive" size="sm"`
- レイアウト: 横並び（`flex items-center gap-4`） → レスポンシブ（`flex-col sm:flex-row`）
- ユーザー情報表示: インライン表示 → ブロック表示（「ログイン中: {user.name}」形式）

**注意事項**:
- UIデザインが変わる可能性があるが、これは許容する
- ヘッダーのレイアウト構造が変わる可能性があるが、これは許容する
- 完全に同じ実装になることで、Issue #166の要件を満たす

### 3.2 実装方法の統一

#### 3.2.1 コンポーネントの直接使用
- **統一方法**: すべてのログアウトボタンで`AuthButtons`コンポーネントを使用
- **メリット**:
  - 完全に同じ実装になる（画面内のログアウトボタンと全く同じ）
  - コンポーネントの再利用性が高い
  - 認証関連のUI変更時に`AuthButtons`コンポーネントのみを修正すればよい
  - Issue #166の「作りを同じにしたい」という要件を完全に満たす

#### 3.2.2 既存実装との整合性
- **auth-buttons.tsx**: 既存の実装を維持（変更不要）
- **auth-actions.ts**: 既存の実装を維持（変更不要）
- **page.tsx**: 既存の実装を維持（変更不要）

## 4. 非機能要件

### 4.1 パフォーマンス
- **変更による影響**: パフォーマンスへの影響はない（`AuthButtons`コンポーネントは既存の実装と同じ処理を行うため）
- **実行時間**: 既存の実装と同等の実行時間を維持

### 4.2 互換性
- **既存機能の維持**: ログアウト機能の動作は変更しない
- **ブラウザ互換性**: 既存のブラウザ互換性を維持
- **Next.js互換性**: 既存のNext.jsバージョンとの互換性を維持

### 4.3 セキュリティ
- **認証処理**: 既存の認証処理を維持（`signOut()`の呼び出し方法のみ変更）
- **セキュリティリスク**: 変更によるセキュリティリスクはない

### 4.4 保守性
- **コードの一貫性**: すべてのログアウトボタンで同じコンポーネント（`AuthButtons`）を使用
- **変更の容易性**: 認証関連のUI変更時に`AuthButtons`コンポーネントのみを修正すればよい
- **可読性**: コードの可読性が向上（統一されたコンポーネント使用）
- **コンポーネントの再利用性**: `AuthButtons`コンポーネントの再利用性が向上

### 4.5 動作環境
- **開発環境**: 既存の開発環境で動作する
- **本番環境**: 既存の本番環境で動作する
- **依存関係**: 既存の依存関係を維持

## 5. 制約事項

### 5.1 既存システムとの関係
- **認証システム**: 既存の認証システム（Auth0）との互換性を維持
- **Next.js Server Actions**: 既存のNext.js Server Actionsの仕組みを維持
- **既存コンポーネント**: 既存のコンポーネント（`auth-buttons.tsx`）への影響はない

### 5.2 技術スタック
- **Next.js**: 既存のNext.jsバージョンを維持
- **TypeScript**: 既存のTypeScriptバージョンを維持
- **React**: 既存のReactバージョンを維持

### 5.3 依存関係
- **auth-buttons.tsx**: `AuthButtons`コンポーネントが存在する必要がある
- **auth-actions.ts**: `signOutAction`と`signInAction`が存在する必要がある（`AuthButtons`コンポーネント内で使用）
- **@/auth**: `signOut`と`signIn`関数が存在する必要がある（`signOutAction`と`signInAction`内で使用）

### 5.4 運用上の制約
- **テスト**: 既存のテストが正常に動作することを確認する必要がある
- **デプロイ**: 既存のデプロイプロセスに影響を与えない

## 6. 受け入れ基準

### 6.1 navbar.tsxの修正
- [ ] `navbar.tsx`で`AuthButtons`コンポーネントを使用している
- [ ] `AuthButtons`を`@/components/auth/auth-buttons`からインポートしている
- [ ] 直接`signOut()`を呼び出す実装が削除されている
- [ ] 直接`signIn()`を呼び出す実装が削除されている
- [ ] ユーザー情報表示部分が`AuthButtons`コンポーネントに置き換えられている
- [ ] 不要な`signOut`や`signIn`のインポートが削除されている（他で使用されていない場合）
- [ ] ログアウト機能が正常に動作する
- [ ] ログイン機能が正常に動作する（未ログイン時）

### 6.2 実装の統一
- [ ] ヘッダーのログアウトボタンと画面内のログアウトボタンで同じコンポーネント（`AuthButtons`）を使用している
- [ ] 完全に同じ実装になっている（UIデザインも含む）
- [ ] 直接`signOut()`を呼び出す実装が存在しない
- [ ] 直接`signIn()`を呼び出す実装が存在しない（`AuthButtons`コンポーネント内で`signInAction`を使用）

### 6.3 動作確認
- [ ] ヘッダーのログアウトボタンが正常に動作する
- [ ] ログアウト後に適切なページにリダイレクトされる
- [ ] 既存のログアウト機能に影響がない
- [ ] 既存のテストが正常に動作する

### 6.4 コード品質
- [ ] コードの一貫性が保たれている
- [ ] 不要なインポートが削除されている
- [ ] コードの可読性が向上している

## 7. 影響範囲

### 7.1 変更が必要なファイル

#### 更新するファイル
- `client/components/layout/navbar.tsx`: ユーザー情報表示部分とログアウト/ログインボタンを`AuthButtons`コンポーネントに置き換え

### 7.2 既存ファイルの扱い
- `client/components/auth/auth-buttons.tsx`: 変更なし（既に`signOutAction`を使用している）
- `client/lib/actions/auth-actions.ts`: 変更なし（既に`signOutAction`が定義されている）
- `client/app/page.tsx`: 変更なし（既に`AuthButtons`を使用している）

### 7.3 既存機能への影響
- **既存のログアウト機能**: 影響なし（`AuthButtons`コンポーネントは既存の実装と同じ処理を行うため）
- **既存のコンポーネント**: 影響なし（`auth-buttons.tsx`は変更しない）
- **既存のテスト**: 影響の可能性あり（テストが`signOut`や`signIn`を直接モックしている場合、またはnavbar.tsxのUI構造に依存している場合）
- **UIデザイン**: ヘッダーのUIデザインが変わる可能性がある（許容する）

## 8. 実装上の注意事項

### 8.1 navbar.tsxの修正
- **インポートの追加**: `AuthButtons`を`@/components/auth/auth-buttons`からインポート
- **インポートの削除**: `signOut`、`signIn`を`@/auth`からのインポートから削除（他で使用されていない場合）
- **コンポーネントの置き換え**: ユーザー情報表示部分とログアウト/ログインボタンを`<AuthButtons user={session?.user ?? null} />`に置き換え
- **UIの変更**: `AuthButtons`コンポーネントのUIデザインに従う（`variant="destructive" size="sm"`、レスポンシブレイアウト）
- **UIデザインの許容**: デザインが崩れる可能性があるが、これは許容する（完全に同じ実装になることを優先）

### 8.2 テストの確認
- **既存テストの動作確認**: 既存のテストが正常に動作することを確認
- **モックの確認**: テストで`signOut`を直接モックしている場合は、`signOutAction`をモックするように変更が必要な可能性がある

### 8.3 動作確認
- **ログアウト機能**: ヘッダーのログアウトボタンが正常に動作することを確認
- **リダイレクト**: ログアウト後に適切なページにリダイレクトされることを確認
- **既存機能**: 既存のログアウト機能に影響がないことを確認

### 8.4 コードレビュー
- **実装の統一**: すべてのログアウトボタンで`signOutAction`を使用していることを確認
- **不要なコード**: 不要なインポートやコードが削除されていることを確認
- **コードの可読性**: コードの可読性が向上していることを確認

## 9. 参考情報

### 9.1 関連Issue
- GitHub Issue #166: headerのログアウトボタンの作りが、画面内のログアウトボタンの作りと異なる

### 9.2 既存ファイル
- `client/components/layout/navbar.tsx`: ヘッダーのログアウトボタン（修正対象）
- `client/components/auth/auth-buttons.tsx`: 画面内のログアウトボタン（参考実装）
- `client/lib/actions/auth-actions.ts`: 認証アクション（使用する関数）
- `client/app/page.tsx`: ホームページ（`AuthButtons`の使用例）

### 9.3 技術スタック
- **Next.js**: Server Actionsを使用
- **TypeScript**: 型安全性を維持
- **React**: コンポーネントベースの実装

### 9.4 参考リンク
- Next.js Server Actions: https://nextjs.org/docs/app/building-your-application/data-fetching/server-actions-and-mutations
- Next.js Authentication: https://nextjs.org/docs/app/building-your-application/authentication
