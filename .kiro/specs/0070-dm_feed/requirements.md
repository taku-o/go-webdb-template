# フィード機能のURL変更の要件定義書

## 1. 概要

### 1.1 プロジェクト情報
- **プロジェクト名**: go-webdb-template
- **Feature名**: 0070-dm_feed
- **作成日**: 2026-01-27
- **関連Issue**: https://github.com/taku-o/go-webdb-template/issues/145

### 1.2 目的
既存のフィード機能のURLパスを変更する。`/feed/`から`/dm_feed/`に変更し、ダミー機能であることを明確にする。

### 1.3 スコープ
- フィード機能のURLパスの変更（`/feed/` → `/dm_feed/`）
- リダイレクト用ページのURL変更（`/feed` → `/dm_feed`）
- フィード一覧ページのURL変更（`/feed/[userId]` → `/dm_feed/[userId]`）
- 返信一覧ページのURL変更（`/feed/[userId]/[postId]` → `/dm_feed/[userId]/[postId]`）
- トップページのリンク先変更（`/feed` → `/dm_feed`）
- コンポーネント内のリンク先変更（`/feed/` → `/dm_feed/`）

**本実装の範囲外**:
- コンポーネント名の変更（`client/components/feed`は変更しない）
- 機能の追加・削除・変更
- UI/UXの変更
- データ構造の変更

## 2. 背景・現状分析

### 2.1 現状
- フィード機能は`/feed/`パスで実装されている
- 以下のページが存在:
  - `/feed` - リダイレクト用ページ（`/feed/[userId]`にリダイレクト）
  - `/feed/[userId]` - フィード一覧ページ
  - `/feed/[userId]/[postId]` - 返信一覧ページ
- トップページ（`client/app/page.tsx`）に「フィード」リンクが存在（リンク先: `/feed`）
- コンポーネント内に`/feed/`へのリンクが存在

### 2.2 問題点
- フィード機能はダミー用の機能であるが、URLが`/feed/`となっており、ダミー機能であることが明確でない

### 2.3 必要性
- ダミー機能であることをURLで明確にするため、`/dm_feed/`に変更する必要がある

### 2.4 実現可否
- 既存のファイルパスとコンポーネント名は変更せず、URLパスのみを変更するため、実装可能
- 既存の機能への影響は最小限

## 3. 機能要件

### 3.1 URLパスの変更

#### 3.1.1 リダイレクト用ページのURL変更
- **変更対象**: `client/app/feed/page.tsx`
- 現在のURL: `/feed`
- 変更後のURL: `/dm_feed`
- リダイレクト先: `/feed/[userId]` → `/dm_feed/[userId]`に変更

#### 3.1.2 フィード一覧ページのURL変更
- **変更対象**: `client/app/feed/[userId]/page.tsx`
- 現在のURL: `/feed/[userId]`
- 変更後のURL: `/dm_feed/[userId]`
- ディレクトリをリネームするため、ファイルパスは`client/app/dm_feed/[userId]/page.tsx`に変更される

#### 3.1.3 返信一覧ページのURL変更
- **変更対象**: `client/app/feed/[userId]/[postId]/page.tsx`
- 現在のURL: `/feed/[userId]/[postId]`
- 変更後のURL: `/dm_feed/[userId]/[postId]`
- ディレクトリをリネームするため、ファイルパスは`client/app/dm_feed/[userId]/[postId]/page.tsx`に変更される

### 3.2 リンク先の変更

#### 3.2.1 トップページのリンク変更
- **変更対象**: `client/app/page.tsx`
- 現在のリンク先: `/feed`
- 変更後のリンク先: `/dm_feed`

#### 3.2.2 コンポーネント内のリンク変更
- **変更対象**: `client/components/feed/`内のコンポーネント
- 現在のリンク先: `/feed/`で始まるパス
- 変更後のリンク先: `/dm_feed/`で始まるパス
- 変更が必要な箇所:
  - `feed-post-card.tsx`: `/feed/${userId}/${post.id}` → `/dm_feed/${userId}/${post.id}`
  - その他、`/feed/`で始まるリンクが存在するコンポーネント

### 3.3 リダイレクト処理の変更
- **変更対象**: `client/app/dm_feed/page.tsx`（ディレクトリリネーム後）
- リダイレクト先を`/feed/[userId]`から`/dm_feed/[userId]`に変更

## 4. 非機能要件

### 4.1 互換性
- 既存の機能への影響を最小限にする
- コンポーネント名やファイル構造は変更しない

### 4.2 保守性
- URLパスの変更のみを行い、既存のコード構造を維持する
- 変更箇所を明確にする

## 5. 制約事項

### 5.1 実装上の制約
- コンポーネント名は変更しない（`client/components/feed`は変更しない）
- `client/app/feed/`ディレクトリを`client/app/dm_feed/`にリネームする
- URLパスを変更する（ディレクトリ名の変更により自動的に変更される）

### 5.2 機能の制約
- 既存の機能は変更しない
- UI/UXは変更しない
- データ構造は変更しない

## 6. 受け入れ基準

### 6.1 URLパスの変更
- [ ] `/dm_feed`にアクセスした際、`/dm_feed/[userId]`にリダイレクトされる
- [ ] `/dm_feed/[userId]`にアクセスした際、フィード一覧ページが表示される
- [ ] `/dm_feed/[userId]/[postId]`にアクセスした際、返信一覧ページが表示される
- [ ] 旧URL（`/feed`、`/feed/[userId]`、`/feed/[userId]/[postId]`）へのアクセスは404エラーとなる（または適切にリダイレクトされる）

### 6.2 リンク先の変更
- [ ] トップページの「フィード」リンクが`/dm_feed`を指している
- [ ] コンポーネント内の全ての`/feed/`で始まるリンクが`/dm_feed/`に変更されている
- [ ] フィード一覧ページから返信一覧ページへの遷移が正常に動作する

### 6.3 動作確認
- [ ] フィード一覧ページが正常に表示される
- [ ] 返信一覧ページが正常に表示される
- [ ] 新規投稿機能が正常に動作する
- [ ] 返信機能が正常に動作する
- [ ] いいね機能が正常に動作する
- [ ] 無限スクロール機能が正常に動作する

## 7. 影響範囲

### 7.1 変更されるファイル
- `client/app/page.tsx`: リンク先を`/feed`から`/dm_feed`に変更
- `client/app/dm_feed/page.tsx`: リダイレクト先を`/feed/[userId]`から`/dm_feed/[userId]`に変更（ディレクトリリネーム後）
- `client/components/feed/feed-post-card.tsx`: リンク先を`/feed/`から`/dm_feed/`に変更
- その他、`/feed/`で始まるリンクが存在するコンポーネント

### 7.2 ディレクトリのリネーム
- `client/app/feed/`ディレクトリを`client/app/dm_feed/`にリネーム

### 7.3 既存ファイルへの影響
- 既存のコンポーネント構造への影響なし（URLパスの変更のみ）
- 既存の機能への影響なし

### 7.4 既存機能への影響
- 既存の機能への影響なし（URLパスの変更のみ）

## 8. 実装上の注意事項

### 8.1 URLパスの変更方法
- Next.jsのApp Routerでは、ディレクトリ名がURLパスになる
- `client/app/feed/`ディレクトリを`client/app/dm_feed/`にリネームする
- これにより、URLパスが`/feed/`から`/dm_feed/`に変更される

### 8.2 リンクの変更
- コンポーネント内の全ての`/feed/`で始まるリンクを`/dm_feed/`に変更する
- ハードコードされたリンクだけでなく、動的に生成されるリンクも確認する

### 8.3 テスト
- 変更後のURLで全ての機能が正常に動作することを確認する
- 旧URLへのアクセスが適切に処理されることを確認する

## 9. 参考情報

### 9.1 関連ドキュメント
- 既存の`client/app/feed/`ディレクトリ
- 既存の`client/components/feed/`ディレクトリ
- Next.js App Routerドキュメント

### 9.2 関連Issue
- https://github.com/taku-o/go-webdb-template/issues/145: 本要件定義書の元となったIssue
- https://github.com/taku-o/go-webdb-template/issues/142: フィード機能の実装Issue（0069-feed-component）

### 9.3 技術スタック
- **フレームワーク**: Next.js 14+ (App Router)
- **言語**: TypeScript 5+
- **UIライブラリ**: shadcn/ui
- **スタイリング**: Tailwind CSS

### 9.4 実装の流れ
1. `client/app/feed/`ディレクトリを`client/app/dm_feed/`にリネーム
2. `client/app/dm_feed/page.tsx`のリダイレクト先を`/dm_feed/[userId]`に変更
3. `client/app/page.tsx`のリンク先を`/dm_feed`に変更
4. `client/components/feed/`内のコンポーネントのリンク先を`/dm_feed/`に変更
5. 動作確認

### 9.5 依存関係
- 既存のNext.jsプロジェクト構造
- 既存のフィード機能（0069-feed-component）
