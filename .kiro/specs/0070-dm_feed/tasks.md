# フィード機能のURL変更の実装タスク一覧

## 概要
フィード機能のURLパスを`/feed/`から`/dm_feed/`に変更するためのタスク一覧。要件定義書と設計書に基づき、実装可能な粒度でタスクを分解。

## 実装フェーズ

### Phase 1: ディレクトリのリネーム

#### - [ ] タスク 1.1: ディレクトリのリネーム
**目的**: `client/app/feed/`ディレクトリを`client/app/dm_feed/`にリネームする

**作業内容**:
- `client/app/feed/`ディレクトリを`client/app/dm_feed/`にリネーム
- これにより、以下のファイルパスが自動的に変更される:
  - `client/app/feed/page.tsx` → `client/app/dm_feed/page.tsx`
  - `client/app/feed/[userId]/page.tsx` → `client/app/dm_feed/[userId]/page.tsx`
  - `client/app/feed/[userId]/[postId]/page.tsx` → `client/app/dm_feed/[userId]/[postId]/page.tsx`

**受け入れ基準**:
- `client/app/feed/`ディレクトリが存在しない
- `client/app/dm_feed/`ディレクトリが存在する
- `client/app/dm_feed/page.tsx`が存在する
- `client/app/dm_feed/[userId]/page.tsx`が存在する
- `client/app/dm_feed/[userId]/[postId]/page.tsx`が存在する

_Requirements: 3.1, 7.2, Design: 1.1_

---

### Phase 2: リダイレクト用ページの変更

#### - [ ] タスク 2.1: リダイレクト先の変更
**目的**: リダイレクト用ページのリダイレクト先を`/feed/[userId]`から`/dm_feed/[userId]`に変更する

**作業内容**:
- `client/app/dm_feed/page.tsx`を開く
- `redirect()`関数の引数を`/feed/${DUMMY_USER_ID}`から`/dm_feed/${DUMMY_USER_ID}`に変更
- その他の処理は変更しない

**受け入れ基準**:
- `redirect()`関数の引数が`/dm_feed/${DUMMY_USER_ID}`になっている
- `/dm_feed`にアクセスした際、`/dm_feed/[userId]`にリダイレクトされる
- TypeScriptのコンパイルエラーが発生しない

_Requirements: 3.1.1, 3.3, 6.1, Design: 2_

---

### Phase 3: トップページのリンク変更

#### - [ ] タスク 3.1: トップページのリンク先変更
**目的**: トップページの「フィード」リンクのリンク先を`/feed`から`/dm_feed`に変更する

**作業内容**:
- `client/app/page.tsx`を開く
- `features`配列内の`href: '/feed'`を`href: '/dm_feed'`に変更
- その他の処理は変更しない

**受け入れ基準**:
- `features`配列内の`href`プロパティが`/dm_feed`になっている
- トップページの「フィード」リンクが`/dm_feed`を指している
- リンクから`/dm_feed`に遷移できる
- TypeScriptのコンパイルエラーが発生しない

_Requirements: 3.2.1, 6.2, Design: 3_

---

### Phase 4: コンポーネント内のリンク変更

#### - [ ] タスク 4.1: FeedPostCardコンポーネントのリンク変更
**目的**: FeedPostCardコンポーネントのリンク先を`/feed/`から`/dm_feed/`に変更する

**作業内容**:
- `client/components/feed/feed-post-card.tsx`を開く
- `Link`コンポーネントの`href`プロパティを`/feed/${userId}/${post.id}`から`/dm_feed/${userId}/${post.id}`に変更
- その他の処理は変更しない

**受け入れ基準**:
- `Link`コンポーネントの`href`プロパティが`/dm_feed/${userId}/${post.id}`になっている
- フィード一覧ページから返信一覧ページへの遷移が正常に動作する
- TypeScriptのコンパイルエラーが発生しない

_Requirements: 3.2.2, 6.2, Design: 4.1_

---

#### - [ ] タスク 4.2: 返信一覧ページのリンク変更
**目的**: 返信一覧ページの「フィードに戻る」リンクのリンク先を`/feed/`から`/dm_feed/`に変更する

**作業内容**:
- `client/app/dm_feed/[userId]/[postId]/page.tsx`を開く
- `Link`コンポーネントの`href`プロパティを`/feed/${userId}`から`/dm_feed/${userId}`に変更
- その他の処理は変更しない

**受け入れ基準**:
- `Link`コンポーネントの`href`プロパティが`/dm_feed/${userId}`になっている
- 返信一覧ページからフィード一覧ページへの遷移が正常に動作する
- TypeScriptのコンパイルエラーが発生しない

_Requirements: 3.2.2, 6.2, Design: 4.2_

---

### Phase 5: リンクの確認

#### - [ ] タスク 5.1: リンクの確認
**目的**: `/feed/`で始まるリンクが全て`/dm_feed/`に変更されていることを確認する

**作業内容**:
- `grep`コマンドなどで`/feed`を含む全てのファイルを検索
- 変更漏れがないことを確認
- 必要に応じて追加の変更を実施

**受け入れ基準**:
- `/feed`を含むファイルが存在しない（または、コメントやドキュメント内のみ）
- 全てのリンクが`/dm_feed/`に変更されている
- 変更漏れがない

_Requirements: 3.2.2, 6.2, 8.2, Design: 5_

---

### Phase 6: 動作確認

#### - [ ] タスク 6.1: URLパスの動作確認
**目的**: 変更後のURLパスが正常に動作することを確認する

**作業内容**:
- 開発サーバーを起動
- `/dm_feed`にアクセスし、`/dm_feed/[userId]`にリダイレクトされることを確認
- `/dm_feed/[userId]`にアクセスし、フィード一覧ページが表示されることを確認
- `/dm_feed/[userId]/[postId]`にアクセスし、返信一覧ページが表示されることを確認
- 旧URL（`/feed`、`/feed/[userId]`、`/feed/[userId]/[postId]`）へのアクセスが404エラーになることを確認

**受け入れ基準**:
- `/dm_feed`にアクセスした際、`/dm_feed/[userId]`にリダイレクトされる
- `/dm_feed/[userId]`にアクセスした際、フィード一覧ページが表示される
- `/dm_feed/[userId]/[postId]`にアクセスした際、返信一覧ページが表示される
- 旧URLへのアクセスが404エラーになる

_Requirements: 6.1, Design: 実装手順 ステップ6_

---

#### - [ ] タスク 6.2: リンクの動作確認
**目的**: 全てのリンクが正常に動作することを確認する

**作業内容**:
- トップページの「フィード」リンクから`/dm_feed`に遷移できることを確認
- フィード一覧ページから返信一覧ページへの遷移が正常に動作することを確認
- 返信一覧ページからフィード一覧ページへの遷移が正常に動作することを確認

**受け入れ基準**:
- トップページの「フィード」リンクが`/dm_feed`を指している
- フィード一覧ページから返信一覧ページへの遷移が正常に動作する
- 返信一覧ページからフィード一覧ページへの遷移が正常に動作する

_Requirements: 6.2, Design: 実装手順 ステップ6_

---

#### - [ ] タスク 6.3: 機能の動作確認
**目的**: 全ての機能が正常に動作することを確認する

**作業内容**:
- フィード一覧ページが正常に表示されることを確認
- 返信一覧ページが正常に表示されることを確認
- 新規投稿機能が正常に動作することを確認
- 返信機能が正常に動作することを確認
- いいね機能が正常に動作することを確認
- 無限スクロール機能が正常に動作することを確認

**受け入れ基準**:
- フィード一覧ページが正常に表示される
- 返信一覧ページが正常に表示される
- 新規投稿機能が正常に動作する
- 返信機能が正常に動作する
- いいね機能が正常に動作する
- 無限スクロール機能が正常に動作する

_Requirements: 6.3, Design: 実装手順 ステップ6_

---

## 実装順序

1. Phase 1: ディレクトリのリネーム（タスク 1.1）
2. Phase 2: リダイレクト用ページの変更（タスク 2.1）
3. Phase 3: トップページのリンク変更（タスク 3.1）
4. Phase 4: コンポーネント内のリンク変更（タスク 4.1, 4.2）
5. Phase 5: リンクの確認（タスク 5.1）
6. Phase 6: 動作確認（タスク 6.1, 6.2, 6.3）

## 注意事項

### 実装時の注意点
- ディレクトリのリネームは最初に行う（他のタスクはリネーム後に実施）
- リンクの変更は全てのファイルで実施する必要がある
- 変更漏れがないよう、`grep`コマンドなどで確認する
- 動作確認は全ての機能で実施する

### テスト方針
- 単体テストは不要（URLパスの変更のみで、ロジックの変更はない）
- 統合テストは不要（URLパスの変更のみで、機能の変更はない）
- 動作確認は手動で実施する

## 関連ドキュメント

- 要件定義書: `.kiro/specs/0070-dm_feed/requirements.md`
- 設計書: `.kiro/specs/0070-dm_feed/design.md`
- 関連Issue: https://github.com/taku-o/go-webdb-template/issues/145
