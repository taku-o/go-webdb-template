# フィード機能のURL変更の設計書

## Overview

### 目的
既存のフィード機能のURLパスを`/feed/`から`/dm_feed/`に変更し、ダミー機能であることを明確にする。

### ユーザー
- **開発者**: ダミー機能であることをURLで識別できる
- **ユーザー**: ダミー機能であることをURLで認識できる

### 影響
現在のシステム状態を以下のように変更する：
- `client/app/feed/`ディレクトリを`client/app/dm_feed/`にリネーム
- `client/app/page.tsx`: リンク先を`/feed`から`/dm_feed`に変更
- `client/app/dm_feed/page.tsx`: リダイレクト先を`/feed/[userId]`から`/dm_feed/[userId]`に変更
- `client/components/feed/feed-post-card.tsx`: リンク先を`/feed/`から`/dm_feed/`に変更
- `client/app/dm_feed/[userId]/[postId]/page.tsx`: リンク先を`/feed/`から`/dm_feed/`に変更

### Goals
- フィード機能のURLパスを`/feed/`から`/dm_feed/`に変更する
- 全てのリンク先を`/dm_feed/`に統一する
- 既存の機能を維持する（URLパスの変更のみ）

### Non-Goals
- コンポーネント名の変更（`client/components/feed`は変更しない）
- 機能の追加・削除・変更
- UI/UXの変更
- データ構造の変更
- 旧URLへのリダイレクト処理（404エラーで問題なし）

## Architecture

### ディレクトリ構造の変更

#### 変更前
```
client/app/
├── feed/
│   ├── page.tsx                    # リダイレクト用ページ
│   └── [userId]/
│       ├── page.tsx                # フィード一覧ページ
│       └── [postId]/
│           └── page.tsx            # 返信一覧ページ
└── page.tsx                        # トップページ（リンク先: /feed）
```

#### 変更後
```
client/app/
├── dm_feed/                        # feed/からリネーム
│   ├── page.tsx                    # リダイレクト用ページ（リダイレクト先を変更）
│   └── [userId]/
│       ├── page.tsx                # フィード一覧ページ（変更なし）
│       └── [postId]/
│           └── page.tsx            # 返信一覧ページ（リンク先を変更）
└── page.tsx                        # トップページ（リンク先: /dm_feed）
```

### URLパスの変更

#### 変更前のURL
- `/feed` - リダイレクト用ページ
- `/feed/[userId]` - フィード一覧ページ
- `/feed/[userId]/[postId]` - 返信一覧ページ

#### 変更後のURL
- `/dm_feed` - リダイレクト用ページ
- `/dm_feed/[userId]` - フィード一覧ページ
- `/dm_feed/[userId]/[postId]` - 返信一覧ページ

## 詳細設計

### 1. ディレクトリのリネーム

#### 1.1 ディレクトリリネーム処理
- **対象**: `client/app/feed/`ディレクトリ
- **操作**: `client/app/feed/`を`client/app/dm_feed/`にリネーム
- **方法**: ファイルシステムのリネーム操作を使用
- **影響**: Next.jsのApp Routerでは、ディレクトリ名がURLパスになるため、自動的にURLパスが変更される

#### 1.2 リネーム後のファイル構造
```
client/app/dm_feed/
├── page.tsx                        # リダイレクト用ページ
└── [userId]/
    ├── page.tsx                    # フィード一覧ページ
    └── [postId]/
        └── page.tsx                # 返信一覧ページ
```

### 2. リダイレクト用ページの変更

#### 2.1 ファイルパス
- **変更前**: `client/app/feed/page.tsx`
- **変更後**: `client/app/dm_feed/page.tsx`（ディレクトリリネームにより自動的に変更）

#### 2.2 変更内容
```typescript
// 変更前
redirect(`/feed/${DUMMY_USER_ID}`)

// 変更後
redirect(`/dm_feed/${DUMMY_USER_ID}`)
```

#### 2.3 実装詳細
- `redirect()`関数の引数を`/feed/${DUMMY_USER_ID}`から`/dm_feed/${DUMMY_USER_ID}`に変更
- その他の処理は変更なし

### 3. トップページのリンク変更

#### 3.1 ファイルパス
- **対象**: `client/app/page.tsx`

#### 3.2 変更内容
```typescript
// 変更前
{
  title: 'フィード',
  description: 'Twitter風のフィードUIコンポーネント',
  href: '/feed',
}

// 変更後
{
  title: 'フィード',
  description: 'Twitter風のフィードUIコンポーネント',
  href: '/dm_feed',
}
```

#### 3.3 実装詳細
- `features`配列内の`href`プロパティを`/feed`から`/dm_feed`に変更
- その他の処理は変更なし

### 4. コンポーネント内のリンク変更

#### 4.1 FeedPostCardコンポーネント

##### 4.1.1 ファイルパス
- **対象**: `client/components/feed/feed-post-card.tsx`

##### 4.1.2 変更内容
```typescript
// 変更前
<Link
  href={`/feed/${userId}/${post.id}`}
  className="inline-flex items-center gap-1 text-muted-foreground hover:text-primary transition-colors"
  aria-label="返信一覧を表示"
>

// 変更後
<Link
  href={`/dm_feed/${userId}/${post.id}`}
  className="inline-flex items-center gap-1 text-muted-foreground hover:text-primary transition-colors"
  aria-label="返信一覧を表示"
>
```

##### 4.1.3 実装詳細
- `Link`コンポーネントの`href`プロパティを`/feed/${userId}/${post.id}`から`/dm_feed/${userId}/${post.id}`に変更
- その他の処理は変更なし

#### 4.2 返信一覧ページのリンク

##### 4.2.1 ファイルパス
- **対象**: `client/app/dm_feed/[userId]/[postId]/page.tsx`（ディレクトリリネーム後）

##### 4.2.2 変更内容
```typescript
// 変更前
<Link
  href={`/feed/${userId}`}
  className="inline-flex items-center text-primary hover:underline text-sm sm:text-base"
  aria-label="フィードに戻る"
>

// 変更後
<Link
  href={`/dm_feed/${userId}`}
  className="inline-flex items-center text-primary hover:underline text-sm sm:text-base"
  aria-label="フィードに戻る"
>
```

##### 4.2.3 実装詳細
- `Link`コンポーネントの`href`プロパティを`/feed/${userId}`から`/dm_feed/${userId}`に変更
- その他の処理は変更なし

### 5. 変更箇所の一覧

#### 5.1 ディレクトリのリネーム
- `client/app/feed/` → `client/app/dm_feed/`

#### 5.2 ファイル内の変更
1. `client/app/dm_feed/page.tsx`（リネーム後）
   - リダイレクト先: `/feed/${DUMMY_USER_ID}` → `/dm_feed/${DUMMY_USER_ID}`

2. `client/app/page.tsx`
   - リンク先: `/feed` → `/dm_feed`

3. `client/components/feed/feed-post-card.tsx`
   - リンク先: `/feed/${userId}/${post.id}` → `/dm_feed/${userId}/${post.id}`

4. `client/app/dm_feed/[userId]/[postId]/page.tsx`（リネーム後）
   - リンク先: `/feed/${userId}` → `/dm_feed/${userId}`

### 6. 変更不要な箇所

#### 6.1 コンポーネントディレクトリ
- `client/components/feed/`は変更しない（コンポーネント名は変更しない）

#### 6.2 コンポーネントの実装
- コンポーネントの機能やUIは変更しない
- 型定義やAPI関数は変更しない

## 実装手順

### ステップ1: ディレクトリのリネーム
1. `client/app/feed/`ディレクトリを`client/app/dm_feed/`にリネーム
2. これにより、以下のファイルパスが自動的に変更される:
   - `client/app/feed/page.tsx` → `client/app/dm_feed/page.tsx`
   - `client/app/feed/[userId]/page.tsx` → `client/app/dm_feed/[userId]/page.tsx`
   - `client/app/feed/[userId]/[postId]/page.tsx` → `client/app/dm_feed/[userId]/[postId]/page.tsx`

### ステップ2: リダイレクト用ページの変更
1. `client/app/dm_feed/page.tsx`を開く
2. `redirect()`関数の引数を`/feed/${DUMMY_USER_ID}`から`/dm_feed/${DUMMY_USER_ID}`に変更

### ステップ3: トップページのリンク変更
1. `client/app/page.tsx`を開く
2. `features`配列内の`href: '/feed'`を`href: '/dm_feed'`に変更

### ステップ4: FeedPostCardコンポーネントのリンク変更
1. `client/components/feed/feed-post-card.tsx`を開く
2. `Link`コンポーネントの`href`プロパティを`/feed/${userId}/${post.id}`から`/dm_feed/${userId}/${post.id}`に変更

### ステップ5: 返信一覧ページのリンク変更
1. `client/app/dm_feed/[userId]/[postId]/page.tsx`を開く
2. `Link`コンポーネントの`href`プロパティを`/feed/${userId}`から`/dm_feed/${userId}`に変更

### ステップ6: 動作確認
1. 開発サーバーを起動
2. `/dm_feed`にアクセスし、`/dm_feed/[userId]`にリダイレクトされることを確認
3. `/dm_feed/[userId]`にアクセスし、フィード一覧ページが表示されることを確認
4. `/dm_feed/[userId]/[postId]`にアクセスし、返信一覧ページが表示されることを確認
5. トップページの「フィード」リンクから`/dm_feed`に遷移できることを確認
6. フィード一覧ページから返信一覧ページへの遷移が正常に動作することを確認
7. 全ての機能（新規投稿、返信、いいね、無限スクロール）が正常に動作することを確認

## 技術的詳細

### Next.js App Routerの仕様
- Next.jsのApp Routerでは、`app`ディレクトリ内のディレクトリ名がURLパスになる
- ディレクトリをリネームすると、自動的にURLパスが変更される
- 例: `app/feed/` → `app/dm_feed/`にリネームすると、URLパスが`/feed/` → `/dm_feed/`に変更される

### リンクの変更方法
- Next.jsの`Link`コンポーネントの`href`プロパティを変更する
- 動的に生成されるリンク（テンプレートリテラルを使用）も同様に変更する

### リダイレクトの変更方法
- Next.jsの`redirect()`関数の引数を変更する
- サーバーコンポーネント内で使用される`redirect()`関数は、リダイレクト先のURLを指定する

## テスト方針

### 単体テスト
- 不要（URLパスの変更のみで、ロジックの変更はない）

### 統合テスト
- 不要（URLパスの変更のみで、機能の変更はない）

### 動作確認
- 手動で以下の項目を確認:
  1. `/dm_feed`にアクセスし、リダイレクトが正常に動作することを確認
  2. `/dm_feed/[userId]`にアクセスし、フィード一覧ページが表示されることを確認
  3. `/dm_feed/[userId]/[postId]`にアクセスし、返信一覧ページが表示されることを確認
  4. トップページの「フィード」リンクから`/dm_feed`に遷移できることを確認
  5. フィード一覧ページから返信一覧ページへの遷移が正常に動作することを確認
  6. 全ての機能（新規投稿、返信、いいね、無限スクロール）が正常に動作することを確認

## リスクと対策

### リスク1: 旧URLへのアクセス
- **リスク**: 旧URL（`/feed`、`/feed/[userId]`、`/feed/[userId]/[postId]`）へのアクセスが404エラーになる
- **対策**: 要件定義書で「旧URLへのアクセスは404エラーとなる（または適切にリダイレクトされる）」と明記されているため、問題なし

### リスク2: リンクの見落とし
- **リスク**: 一部のリンクが変更されず、旧URLを指し続ける可能性
- **対策**: `grep`コマンドなどで`/feed`を含む全てのファイルを検索し、変更漏れがないことを確認する

### リスク3: ディレクトリリネーム時のエラー
- **リスク**: ディレクトリリネーム時にファイルが失われる可能性
- **対策**: Gitで管理されているため、問題が発生した場合は元に戻すことができる

## 依存関係

### 既存の機能への依存
- 既存のフィード機能（0069-feed-component）に依存
- 既存のコンポーネント（`client/components/feed/`）に依存

### 外部ライブラリへの依存
- Next.js 14+ (App Router)
- React
- TypeScript

## 参考情報

### 関連ドキュメント
- 要件定義書: `.kiro/specs/0070-dm_feed/requirements.md`
- 既存のフィード機能の設計書: `.kiro/specs/0069-feed-component/design.md`

### 関連Issue
- https://github.com/taku-o/go-webdb-template/issues/145: 本設計書の元となったIssue

### 技術スタック
- **フレームワーク**: Next.js 14+ (App Router)
- **言語**: TypeScript 5+
- **UIライブラリ**: shadcn/ui
- **スタイリング**: Tailwind CSS
