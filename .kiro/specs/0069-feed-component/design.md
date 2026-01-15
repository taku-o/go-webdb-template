# フィードUIコンポーネントの設計書

## Overview

### 目的
Twitter風のメッセージのエントリーが並ぶUIをNext.jsで作成し、投稿一覧の表示、新規投稿、返信、いいね機能を提供する。**本実装の主な目的は、コンポーネントの挙動と見た目を重視することである。** データ管理はメモリで行い、データベース連携は行わない。クライアント側でダミーデータを返すAPI実装を使用し、初期の投稿データ・返信データはランダムで生成する。

### ユーザー
- **開発者**: コンポーネントの挙動と見た目を確認する
- **デモンストレーション**: ソーシャルメディア風のUIコンポーネントの実装例を確認する

### 影響
現在のシステム状態を以下のように変更する：
- `client/app/feed/page.tsx`: 新規作成（リダイレクト用ページ、`/feed/[userId]`にリダイレクト）
- `client/app/feed/[userId]/page.tsx`: 新規作成（フィード一覧ページ）
- `client/app/feed/[userId]/[postId]/page.tsx`: 新規作成（返信一覧ページ）
- `client/app/page.tsx`: 「フィード」リンクを追加
- `client/lib/api.ts`: フィード関連のAPI関数を追加
- `client/types/dm_feed.ts`: 新規作成（型定義）
- `client/components/feed/`: 新規作成（必要に応じて、フィード関連コンポーネント）

### Goals
- Twitter風のフィードUIコンポーネントを作成する（挙動と見た目を重視）
- 投稿一覧の表示機能を実装する（初期10件、無限スクロール）
- 新規投稿機能を実装する
- 返信機能を実装する（返信一覧ページ）
- いいね機能を実装する
- shadcn/uiコンポーネントを活用した統一感のあるデザインを実装する
- レスポンシブデザインに対応する

### Non-Goals
- データベース連携（メモリ管理のみ）
- サーバー側API実装
- 認証機能の新規実装（既存の認証機能を使用）
- 返信のネスト機能（返信に対する返信）
- 返信へのいいね機能
- その他のSNS機能（リツイート、ブックマーク等）
- データの永続化（メモリ管理のみ）
- コンポーネントの単体テスト
- E2Eテスト

## Architecture

### ページ構成

#### リダイレクト用ページ (`/feed`)
```
client/app/feed/page.tsx
└── リダイレクト処理
    └── /feed/[userId] にリダイレクト（固定のダミーuserIdを使用）
```

#### フィード一覧ページ (`/feed/[userId]`)
```
client/app/feed/[userId]/page.tsx
├── 新規投稿フォーム
│   ├── テキストエリア（最大280文字）
│   ├── 文字数カウンター
│   └── 投稿ボタン
└── 投稿一覧
    ├── 投稿カード（10件ずつ表示）
    │   ├── 投稿者情報（表示名、ユーザー名）
    │   ├── 投稿日時（相対時刻）
    │   ├── 投稿内容
    │   ├── いいね数・いいねボタン
    │   └── 返信ボタン（返信一覧ページへのリンク）
    ├── 上方向スクロール検知（最新投稿読み込み）
    └── 下方向スクロール検知（古い投稿読み込み）
```

#### 返信一覧ページ (`/feed/[userId]/[postId]`)
```
client/app/feed/[userId]/[postId]/page.tsx
├── 返信元のエントリー（一番上）
│   ├── 投稿者情報（表示名、ユーザー名）
│   ├── 投稿日時（相対時刻）
│   ├── 投稿内容
│   ├── いいね数・いいねボタン
│   └── 返信ボタン（返信フォーム表示用）
├── 返信フォーム
│   ├── テキストエリア（最大280文字）
│   ├── 文字数カウンター
│   └── 返信ボタン
└── 返信一覧
    ├── 返信カード（古いものが上、新しいものが下）
    │   ├── 返信者情報（表示名、ユーザー名）
    │   ├── 返信日時（相対時刻）
    │   └── 返信内容
    ├── 上方向スクロール検知（古い返信読み込み）
    └── 下方向スクロール検知（新しい返信読み込み）
```

### コンポーネント設計

#### リダイレクト用ページのコンポーネント構成
```
RedirectPage (client/app/feed/page.tsx)
└── redirect() 関数を使用して /feed/[userId] にリダイレクト
```

#### フィード一覧ページのコンポーネント構成
```
FeedPage (client/app/feed/[userId]/page.tsx)
├── FeedForm (投稿フォーム)
│   ├── Textarea (shadcn/ui)
│   ├── 文字数カウンター
│   └── Button (shadcn/ui)
└── FeedPostList (投稿一覧)
    ├── FeedPostCard (投稿カード) × N
    │   ├── Card (shadcn/ui)
    │   ├── 投稿者情報表示
    │   ├── 相対時刻表示
    │   ├── 投稿内容表示
    │   ├── LikeButton (いいねボタン)
    │   └── ReplyButton (返信ボタン、Link)
    ├── IntersectionObserver (上方向スクロール検知)
    └── IntersectionObserver (下方向スクロール検知)
```

#### 返信一覧ページのコンポーネント構成
```
ReplyPage (client/app/feed/[userId]/[postId]/page.tsx)
├── FeedPostCard (返信元のエントリー)
│   └── (フィード一覧ページと同じ構造)
├── ReplyForm (返信フォーム)
│   ├── Textarea (shadcn/ui)
│   ├── 文字数カウンター
│   └── Button (shadcn/ui)
└── ReplyList (返信一覧)
    ├── ReplyCard (返信カード) × N
    │   ├── Card (shadcn/ui)
    │   ├── 返信者情報表示
    │   ├── 相対時刻表示
    │   └── 返信内容表示
    ├── IntersectionObserver (上方向スクロール検知)
    └── IntersectionObserver (下方向スクロール検知)
```

### データフロー

#### フィード一覧の表示フロー
```
1. ページ読み込み
   ↓
2. getDmFeedPosts(userId, 10, "") を呼び出し
   ↓
3. メモリからランダム生成された投稿データを取得（初期10件）
   ↓
4. 投稿一覧を表示（時系列順、新しいものが上）
   ↓
5. 上方向スクロール検知
   ↓
6. getDmFeedPosts(userId, 10, "") を呼び出し（最新投稿取得）
   ↓
7. 既存の投稿一覧の上に追加表示（重複チェック）
   ↓
8. 下方向スクロール検知
   ↓
9. getDmFeedPosts(userId, 10, lastPostId) を呼び出し（古い投稿取得）
   ↓
10. 既存の投稿一覧の下に追加表示
```

#### 新規投稿のフロー
```
1. ユーザーが投稿フォームに入力
   ↓
2. バリデーション（空チェック、280文字制限）
   ↓
3. createDmFeedPost(userId, content) を呼び出し
   ↓
4. メモリに新規投稿を追加
   ↓
5. 投稿一覧の一番上に新規投稿を追加表示
   ↓
6. フォームをクリア
```

#### 返信一覧の表示フロー
```
1. 返信一覧ページ読み込み
   ↓
2. getDmFeedPosts(userId, 1, postId) を呼び出し（返信元のエントリー取得）
   ↓
3. getDmFeedReplies(userId, postId, 10, "") を呼び出し（初期10件の返信取得）
   ↓
4. 返信元のエントリーを一番上に表示
   ↓
5. 返信一覧を表示（時系列順、古いものが上、新しいものが下）
   ↓
6. 上方向スクロール検知
   ↓
7. getDmFeedReplies(userId, postId, 10, firstReplyId) を呼び出し（古い返信取得）
   ↓
8. 既存の返信一覧の上に追加表示
   ↓
9. 下方向スクロール検知
   ↓
10. getDmFeedReplies(userId, postId, 10, "") を呼び出し（新しい返信取得）
    ↓
11. 既存の返信一覧の下に追加表示（重複チェック）
```

#### 返信の投稿フロー
```
1. ユーザーが返信フォームに入力
   ↓
2. バリデーション（空チェック、280文字制限）
   ↓
3. replyToDmPost(userId, postId, content) を呼び出し
   ↓
4. メモリに新規返信を追加
   ↓
5. 返信一覧の一番下に新規返信を追加表示
   ↓
6. フォームをクリア
```

#### いいねのフロー
```
1. ユーザーがいいねボタンをクリック
   ↓
2. toggleLikeDmPost(userId, postId) を呼び出し
   ↓
3. メモリのいいね状態を更新（userIdをキーとして管理）
   ↓
4. いいね数を更新
   ↓
5. UIを更新（いいね済みの場合は視覚的に区別）
```

### データ管理

#### メモリ内データ構造
```typescript
// メモリ内のデータストア（簡易実装）
interface FeedDataStore {
  dmFeedPosts: Map<string, DmFeedPost>;  // 投稿IDをキーとする投稿マップ
  replies: Map<string, DmFeedReply[]>;  // 投稿IDをキーとする返信配列マップ
  likes: Map<string, Set<string>>;  // 投稿IDをキーとする、いいね済みユーザーIDのセット
}

// 初期データの生成
- 投稿データ: ランダム生成（投稿者名、投稿内容、日時など）
- 返信データ: ランダム生成（返信者名、返信内容、日時など）
- いいねデータ: ランダム生成（いいね数、いいね済み状態など）
```

#### データのライフサイクル
- **初期化**: ページ読み込み時にランダムデータを生成
- **追加**: 新規投稿・返信時にメモリに追加
- **更新**: いいね状態の更新
- **削除**: ページリロード時に全て初期化

## 詳細設計

### 1. 型定義 (`client/types/dm_feed.ts`)

```typescript
// 投稿の型
export interface DmFeedPost {
  id: string;
  userId: string;
  userName: string;
  userHandle: string;
  content: string;
  createdAt: string;  // ISO 8601形式
  likeCount: number;
  liked: boolean;
}

// 返信の型
export interface DmFeedReply {
  id: string;
  postId: string;
  userId: string;
  userName: string;
  userHandle: string;
  content: string;
  createdAt: string;  // ISO 8601形式
}

// リクエストの型
export interface CreateDmFeedPostRequest {
  content: string;
}

export interface CreateDmFeedReplyRequest {
  content: string;
}
```

### 2. API関数 (`client/lib/api.ts`)

#### 2.1 フィード一覧取得
```typescript
async getDmFeedPosts(
  userId: string,
  limit: number,
  fromPostId: string
): Promise<DmFeedPost[]>
```
- **パラメータ**:
  - `userId`: 操作ユーザーID
  - `limit`: 取得件数（通常10件）
  - `fromPostId`: カーソルベースのページネーション用（空文字列の場合は最新を取得）
- **戻り値**: 投稿配列（時系列順、新しいものが先頭）
- **実装**: メモリ内のデータストアから取得、ランダム生成も含む

#### 2.2 新規投稿作成
```typescript
async createDmFeedPost(
  userId: string,
  content: string
): Promise<DmFeedPost>
```
- **パラメータ**:
  - `userId`: 操作ユーザーID
  - `content`: 投稿内容（最大280文字）
- **戻り値**: 作成された投稿
- **実装**: メモリ内のデータストアに追加

#### 2.3 返信一覧取得
```typescript
async getDmFeedReplies(
  userId: string,
  postId: string,
  limit: number,
  fromReplyId: string
): Promise<DmFeedReply[]>
```
- **パラメータ**:
  - `userId`: 操作ユーザーID
  - `postId`: 投稿ID
  - `limit`: 取得件数（通常10件）
  - `fromReplyId`: カーソルベースのページネーション用（空文字列の場合は最新を取得）
- **戻り値**: 返信配列（時系列順、古いものが先頭、新しいものが末尾）
- **実装**: メモリ内のデータストアから取得、ランダム生成も含む

#### 2.4 返信作成
```typescript
async replyToDmPost(
  userId: string,
  postId: string,
  content: string
): Promise<DmFeedReply>
```
- **パラメータ**:
  - `userId`: 操作ユーザーID
  - `postId`: 投稿ID
  - `content`: 返信内容（最大280文字）
- **戻り値**: 作成された返信
- **実装**: メモリ内のデータストアに追加

#### 2.5 いいねのON/OFF
```typescript
async toggleLikeDmPost(
  userId: string,
  postId: string
): Promise<{ liked: boolean; likeCount: number }>
```
- **パラメータ**:
  - `userId`: 操作ユーザーID
  - `postId`: 投稿ID
- **戻り値**: いいね状態といいね数
- **実装**: メモリ内のデータストアを更新（userIdをキーとして管理）

### 3. リダイレクト用ページ (`client/app/feed/page.tsx`)

#### 3.1 実装方法
```typescript
import { redirect } from 'next/navigation'

export default function FeedRedirectPage() {
  // 固定のダミーuserIdを使用（ダミー実装のため）
  // 将来的には、認証ユーザーのIDを使用する想定
  const dummyUserId = 'dummy-user-001'
  redirect(`/feed/${dummyUserId}`)
}
```

### 4. フィード一覧ページ (`client/app/feed/[userId]/page.tsx`)

#### 4.1 状態管理
```typescript
const [dmFeedPosts, setDmFeedPosts] = useState<DmFeedPost[]>([]);
const [isLoading, setIsLoading] = useState(false);
const [hasMore, setHasMore] = useState(true);
const [newestPostId, setNewestPostId] = useState<string>("");
const [oldestPostId, setOldestPostId] = useState<string>("");
```

#### 4.2 無限スクロール実装
- **上方向スクロール検知**: ページ上部にIntersection Observerを設定
- **下方向スクロール検知**: ページ下部にIntersection Observerを設定
- **スクロール位置の維持**: 上方向スクロール時に最新投稿を追加する際、スクロール位置を維持

#### 4.3 重複チェック
- 投稿IDをキーとして管理し、既存の投稿と重複しないようにする

### 5. 返信一覧ページ (`client/app/feed/[userId]/[postId]/page.tsx`)

#### 5.1 状態管理
```typescript
const [dmFeedPost, setDmFeedPost] = useState<DmFeedPost | null>(null);
const [dmFeedReplies, setDmFeedReplies] = useState<DmFeedReply[]>([]);
const [isLoading, setIsLoading] = useState(false);
const [hasMore, setHasMore] = useState(true);
const [newestReplyId, setNewestReplyId] = useState<string>("");
const [oldestReplyId, setOldestReplyId] = useState<string>("");
```

#### 4.2 無限スクロール実装
- **上方向スクロール検知**: ページ上部にIntersection Observerを設定（古い返信を読み込む）
- **下方向スクロール検知**: ページ下部にIntersection Observerを設定（新しい返信を読み込む）
- 注: 返信一覧は投稿一覧と逆の順序（下が新しい）のため、スクロール挙動も逆になる

### 5. UIコンポーネント

#### 5.1 投稿カード (`FeedPostCard`)
- shadcn/uiの`Card`コンポーネントを使用
- 投稿者情報、投稿内容、いいねボタン、返信ボタンを表示
- レスポンシブデザインに対応

#### 5.2 返信カード (`ReplyCard`)
- shadcn/uiの`Card`コンポーネントを使用
- 返信者情報、返信内容を表示
- レスポンシブデザインに対応

#### 5.3 投稿フォーム (`FeedForm`)
- shadcn/uiの`Textarea`と`Button`を使用
- 文字数カウンターを表示
- バリデーション機能を実装

#### 5.4 返信フォーム (`ReplyForm`)
- shadcn/uiの`Textarea`と`Button`を使用
- 文字数カウンターを表示
- バリデーション機能を実装

### 6. 相対時刻表示

#### 6.1 実装方法
- `date-fns`ライブラリを使用（必要に応じて）
- または、簡易的な相対時刻計算関数を実装

#### 6.2 表示形式
- "5分前"、"1時間前"、"3日前"など
- 1週間以上前の場合は日付を表示

### 7. ランダムデータ生成

#### 7.1 投稿データの生成
- 投稿者名: ランダムな名前のリストから選択
- ユーザーハンドル: ランダムな文字列から生成（@username形式）
- 投稿内容: ランダムなテキストを生成（280文字以内）
- 投稿日時: ランダムな日時を生成（現在時刻から過去に遡る）
- いいね数: ランダムな数値を生成
- いいね済み状態: ランダムに決定

#### 7.2 返信データの生成
- 返信者名: ランダムな名前のリストから選択
- ユーザーハンドル: ランダムな文字列から生成（@username形式）
- 返信内容: ランダムなテキストを生成（280文字以内）
- 返信日時: ランダムな日時を生成（投稿日時から現在時刻まで）

## 実装上の考慮事項

### 1. パフォーマンス
- 無限スクロールにより、大量の投稿でもパフォーマンスを維持
- 仮想スクロールは不要（本実装では）
- メモリ管理は簡易実装で十分

### 2. アクセシビリティ
- 適切な`aria-label`を設定
- キーボード操作に対応
- フォーカス管理を適切に実装

### 3. エラーハンドリング
- API呼び出し時のエラーを適切にハンドリング
- ユーザーフレンドリーなエラーメッセージを表示
- ネットワークエラーなどの適切な処理（ダミー実装では不要だが、将来の拡張を考慮）

### 4. レスポンシブデザイン
- モバイル、タブレット、デスクトップに対応
- Tailwind CSSのレスポンシブクラスを活用
- shadcn/uiコンポーネントのレスポンシブ機能を活用

## 実装順序

1. 型定義の作成（`client/types/dm_feed.ts`）
2. API関数の追加（`client/lib/api.ts`）
   - メモリ内データストアの実装
   - ランダムデータ生成関数の実装
   - 各API関数の実装
3. リダイレクト用ページの作成（`client/app/feed/page.tsx`）
   - リダイレクト処理の実装
4. フィード一覧ページの作成（`client/app/feed/[userId]/page.tsx`）
   - 投稿フォームの実装
   - 投稿一覧の表示
   - 無限スクロールの実装
5. 返信一覧ページの作成（`client/app/feed/[userId]/[postId]/page.tsx`）
   - 返信元エントリーの表示
   - 返信フォームの実装
   - 返信一覧の表示
   - 無限スクロールの実装
6. コンポーネントの作成（必要に応じて）
   - 投稿カードコンポーネント
   - 返信カードコンポーネント
   - 投稿フォームコンポーネント
   - 返信フォームコンポーネント
7. トップページへのリンク追加（`client/app/page.tsx`）
8. 動作確認
