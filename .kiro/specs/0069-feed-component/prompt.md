/kiro:spec-requirements "https://github.com/taku-o/go-webdb-template/issues/142
に対応するための要件定義書を作成してください。
GitHub CLIは入っています。
cc-sddのfeature名は0069-feed-componentとしてください。"
think.


各オブジェクトの名前にはDmをつける。
DmFeedPost
DmFeedReply

API関数のパラメータはこうなる。
> #### 3.6.1 API関数の追加
- `getDmFeedPosts(limit: number, fromPostId: string): Promise<DmFeedPost[]>`
- `createDmFeedPost(content: string): Promise<DmFeedPost>`
- `replyToDmPost(postId: string, content: string): Promise<DmFeedReply>`
- `toggleLikeDmPost(postId: string): Promise<{ liked: boolean; likeCount: number }>`


下にスクロールしたら古いフィードを読み込むが、
上にスクロールしたら最新のフィードを読み込む挙動にしたい。
(もちろんデータはダミーで良い)


画面は2種類用意する。
* フィードの一覧
  * 各エントリーは
      - 投稿者名（表示名）
      - 投稿者ID（ユーザー名）
      - 投稿日時（相対時刻表示、例: "5分前"）
      - 投稿内容（テキスト）
      - いいね数
      - いいねボタン
      - 返信ボタン

* 返信の一覧
  * 一番上は返信元のエントリー
  * その下に返信の一覧が並ぶ

投稿一覧は上が新しいが、
返信一覧は下の方が新しい返信とする。


関数に操作ユーザーのパラメータが必要。
誰がやったという情報と、権限まわりの管理が必要なため。
> #### 3.6.1 API関数の追加
  - `getDmFeedPosts(userId: string, limit: number, fromPostId: string): Promise<DmFeedPost[]>` - フィード一覧取得
  - `createDmFeedPost(userId: string, content: string): Promise<DmFeedPost>` - 新規投稿作成
  - `getDmFeedReplies(userId: string, postId: string, limit: number, fromReplyId: string): Promise<DmFeedReply[]>` - 返信一覧取得
  - `replyToDmPost(userId: string, postId: string, content: string): Promise<DmFeedReply>` - 返信作成
  - `toggleLikeDmPost(userId: string, postId: string): Promise<{ liked: boolean; likeCount: number }>` - いいねのON/OFF

このissueの目的は
コンポーネントの挙動と、見た目を重視したい。
データ管理はメモリで良い。
初期の投稿データ、返信データもランダムで生成しちゃって良い。

要件定義書を承認します。

/kiro:spec-design 0069-feed-component

フィード一覧ページは、他のユーザーの一覧ページを参照することを想定する。
よって、基本的な構成は、
- `client/app/feed/page.tsx`: 新規作成（フィード一覧ページ）
->
- `client/app/feed/[userId]/page.tsx`: 新規作成（フィード一覧ページ）
の方が都合が良い。

しかし、userIdが分からないとURLがわからなくて、フィードの画面に遷移しづらいので、
`client/app/feed/page.tsx` にアクセスしたら、
`client/app/feed/[userId]/page.tsx` にリダイレクトするような挙動にして欲しい。


postという名前が使われているので、
const [dmFeedPosts, setDmFeedPosts] = useState<DmFeedPost[]>([]);
とする

>#### 4.1 状態管理
>```typescript
>const [posts, setPosts] = useState<DmFeedPost[]>([]);

設計書を承認します。

/kiro:spec-tasks 0069-feed-component

タスクリストを承認します。

/sdd-fix-plan

_serena_indexing
/serena-initialize

/kiro:spec-impl 0069-feed-component 1

/kiro:spec-impl 0069-feed-component 2.1
/kiro:spec-impl 0069-feed-component 2.2
/kiro:spec-impl 0069-feed-component 2.3
/kiro:spec-impl 0069-feed-component 2.4
/kiro:spec-impl 0069-feed-component 2.5 2.6 2.7

いったんgit commitしてください。

/kiro:spec-impl 0069-feed-component 3

いったんgit commitしてください。

/kiro:spec-impl 0069-feed-component 4.1
/kiro:spec-impl 0069-feed-component 4.2
/kiro:spec-impl 0069-feed-component 4.3
/kiro:spec-impl 0069-feed-component 4.4
/kiro:spec-impl 0069-feed-component 4.5

/kiro:spec-impl 0069-feed-component 4.6
/kiro:spec-impl 0069-feed-component 4.7
/kiro:spec-impl 0069-feed-component 4.8
/kiro:spec-impl 0069-feed-component 4.9
/kiro:spec-impl 0069-feed-component 4.10

いったんgit commitしてください。

/kiro:spec-impl 0069-feed-component 5.1
/kiro:spec-impl 0069-feed-component 5.2
/kiro:spec-impl 0069-feed-component 5.3
/kiro:spec-impl 0069-feed-component 5.4


ここなんだけど、これは投稿を書いた直後に
データを取得しようとしている？
client/app/feed/[userId]/[postId]/page.tsx
  // ページ読み込み時に返信元の投稿を取得
  useEffect(() => {
    loadPost()
  }, [userId, postId])

消えたのか。ならOK。

/kiro:spec-impl 0069-feed-component 5.5
/kiro:spec-impl 0069-feed-component 5.6
/kiro:spec-impl 0069-feed-component 5.7
/kiro:spec-impl 0069-feed-component 5.8
/kiro:spec-impl 0069-feed-component 5.9

いったんgit commitしてください。






