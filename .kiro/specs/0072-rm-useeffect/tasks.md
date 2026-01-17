# Next.jsのコードからuseEffect排除の実装タスク一覧

## 概要
Next.jsのコードからuseEffectを排除する実装を段階的に進めるためのタスク一覧。要件定義書と設計書に基づき、実装可能な粒度でタスクを分解。

## 実装フェーズ

### Phase 1: 初期データ読み込みの置き換え

#### - [ ] タスク 1.1: `client/app/dm-posts/page.tsx`のuseEffect排除
**目的**: 投稿とユーザーの初期読み込みのuseEffectを排除し、refコールバック関数に置き換える

**作業内容**:
- `useEffect`のインポートを削除（他のuseEffectがない場合）
- `useRef`を使用して、`containerRef`と`hasLoadedRef`を作成
- `loadInitialData`関数を実装（`hasLoadedRef`で重複読み込みを防ぐ）
- refコールバック関数`setContainerRef`を実装
- メインコンテナ要素に`ref={setContainerRef}`を設定
- `useEffect(() => { loadPosts(); loadUsers(); }, [])`を削除
- 動作確認: ページ読み込み時に投稿とユーザーが正常に読み込まれることを確認

**受け入れ基準**:
- `useEffect`が排除されている
- refコールバック関数が実装されている
- ページ読み込み時に投稿とユーザーが正常に読み込まれる
- 既存の機能（作成、削除）が正常に動作する

_Requirements: 6.1, Design: 1.1_

---

#### - [ ] タスク 1.2: `client/app/dm-users/page.tsx`のuseEffect排除
**目的**: ユーザーの初期読み込みのuseEffectを排除し、refコールバック関数に置き換える

**作業内容**:
- `useEffect`のインポートを削除（他のuseEffectがない場合）
- `useRef`を使用して、`containerRef`と`hasLoadedRef`を作成
- `loadInitialData`関数を実装（`hasLoadedRef`で重複読み込みを防ぐ）
- refコールバック関数`setContainerRef`を実装
- メインコンテナ要素に`ref={setContainerRef}`を設定
- `useEffect(() => { loadUsers(); }, [])`を削除
- 動作確認: ページ読み込み時にユーザーが正常に読み込まれることを確認

**受け入れ基準**:
- `useEffect`が排除されている
- refコールバック関数が実装されている
- ページ読み込み時にユーザーが正常に読み込まれる
- 既存の機能（作成、削除、ダウンロード）が正常に動作する

_Requirements: 6.1, Design: 1.1_

---

#### - [ ] タスク 1.3: `client/app/dm-user-posts/page.tsx`のuseEffect排除
**目的**: ユーザー投稿の初期読み込みのuseEffectを排除し、refコールバック関数に置き換える

**作業内容**:
- `useEffect`のインポートを削除（他のuseEffectがない場合）
- `useRef`を使用して、`containerRef`と`hasLoadedRef`を作成
- `loadInitialData`関数を実装（`hasLoadedRef`で重複読み込みを防ぐ）
- refコールバック関数`setContainerRef`を実装
- メインコンテナ要素に`ref={setContainerRef}`を設定
- `useEffect(() => { loadUserPosts(); }, [])`を削除
- 動作確認: ページ読み込み時にユーザー投稿が正常に読み込まれることを確認

**受け入れ基準**:
- `useEffect`が排除されている
- refコールバック関数が実装されている
- ページ読み込み時にユーザー投稿が正常に読み込まれる

_Requirements: 6.1, Design: 1.1_

---

#### - [ ] タスク 1.4: `client/app/dm_feed/[userId]/page.tsx`の初期データ読み込みのuseEffect排除
**目的**: フィード投稿の初期読み込みのuseEffectを排除し、refコールバック関数とuserId変更検知に置き換える

**作業内容**:
- `useEffect`のインポートを削除（他のuseEffectがある場合は残す）
- `useRef`を使用して、`containerRef`と`currentUserIdRef`を作成
- `loadInitialData`関数を実装
- refコールバック関数`setContainerRef`を実装（userIdの変更も検知）
- メインコンテナ要素に`ref={setContainerRef}`を設定
- userIdが変更されたときの処理を追加（render内で`currentUserIdRef.current !== userId`をチェック）
- `useEffect(() => { loadInitialPosts(); }, [userId])`を削除
- 動作確認: ページ読み込み時とuserId変更時にフィード投稿が正常に読み込まれることを確認

**受け入れ基準**:
- 初期データ読み込みの`useEffect`が排除されている
- refコールバック関数が実装されている
- userIdが変更されたときにデータが正常に読み込まれる
- ページ読み込み時にフィード投稿が正常に読み込まれる

_Requirements: 6.1, Design: 1.1_

---

#### - [ ] タスク 1.5: `client/app/dm_feed/[userId]/[postId]/page.tsx`の初期データ読み込みのuseEffect排除
**目的**: 投稿と返信の初期読み込みのuseEffectを排除し、refコールバック関数とuserId/postId変更検知に置き換える

**作業内容**:
- `useEffect`のインポートを削除（他のuseEffectがある場合は残す）
- `useRef`を使用して、`containerRef`、`currentUserIdRef`、`currentPostIdRef`を作成
- `loadInitialData`関数を実装
- refコールバック関数`setContainerRef`を実装（userIdとpostIdの変更も検知）
- メインコンテナ要素に`ref={setContainerRef}`を設定
- userIdまたはpostIdが変更されたときの処理を追加（render内でチェック）
- `useEffect(() => { loadInitialData(); }, [userId, postId])`を削除
- 動作確認: ページ読み込み時とuserId/postId変更時に投稿と返信が正常に読み込まれることを確認

**受け入れ基準**:
- 初期データ読み込みの`useEffect`が排除されている
- refコールバック関数が実装されている
- userIdまたはpostIdが変更されたときにデータが正常に読み込まれる
- ページ読み込み時に投稿と返信が正常に読み込まれる

_Requirements: 6.1, Design: 1.1_

---

### Phase 2: Intersection Observerの置き換え

#### - [ ] タスク 2.1: `client/app/dm_feed/[userId]/page.tsx`のIntersection Observer関連のuseEffect排除（下方向スクロール）
**目的**: 下方向スクロール検知のIntersection ObserverのuseEffectを排除し、refコールバック関数に置き換える

**作業内容**:
- `useEffect`のインポートを削除（他のuseEffectがない場合）
- `loadMoreRef`を`useRef`から`refコールバック関数`に変更
- `setLoadMoreRef`関数を実装（Intersection Observerを登録し、クリーンアップ関数を返す）
- `useEffect`でIntersection Observerを登録していた部分を削除
- 動作確認: 下方向スクロール時に古い投稿が正常に読み込まれることを確認

**受け入れ基準**:
- Intersection Observer関連の`useEffect`が排除されている（下方向スクロール）
- refコールバック関数が実装されている
- 下方向スクロール時に古い投稿が正常に読み込まれる
- クリーンアップが正しく実行される

_Requirements: 6.2, Design: 2.1_

---

#### - [ ] タスク 2.2: `client/app/dm_feed/[userId]/page.tsx`のIntersection Observer関連のuseEffect排除（上方向スクロール）
**目的**: 上方向スクロール検知のIntersection ObserverのuseEffectを排除し、refコールバック関数に置き換える

**作業内容**:
- `loadNewerRef`を`useRef`から`refコールバック関数`に変更
- `setLoadNewerRef`関数を実装（Intersection Observerを登録し、クリーンアップ関数を返す）
- `useEffect`でIntersection Observerを登録していた部分を削除
- 動作確認: 上方向スクロール時に新しい投稿が正常に読み込まれることを確認

**受け入れ基準**:
- Intersection Observer関連の`useEffect`が排除されている（上方向スクロール）
- refコールバック関数が実装されている
- 上方向スクロール時に新しい投稿が正常に読み込まれる
- クリーンアップが正しく実行される

_Requirements: 6.2, Design: 2.1_

---

#### - [ ] タスク 2.3: `client/app/dm_feed/[userId]/[postId]/page.tsx`のIntersection Observer関連のuseEffect排除（下方向スクロール）
**目的**: 下方向スクロール検知のIntersection ObserverのuseEffectを排除し、refコールバック関数に置き換える

**作業内容**:
- `loadOlderRef`を`useRef`から`refコールバック関数`に変更
- `setLoadOlderRef`関数を実装（Intersection Observerを登録し、クリーンアップ関数を返す）
- `useEffect`でIntersection Observerを登録していた部分を削除
- 動作確認: 下方向スクロール時に古い返信が正常に読み込まれることを確認

**受け入れ基準**:
- Intersection Observer関連の`useEffect`が排除されている（下方向スクロール）
- refコールバック関数が実装されている
- 下方向スクロール時に古い返信が正常に読み込まれる
- クリーンアップが正しく実行される

_Requirements: 6.2, Design: 2.1_

---

#### - [ ] タスク 2.4: `client/app/dm_feed/[userId]/[postId]/page.tsx`のIntersection Observer関連のuseEffect排除（上方向スクロール）
**目的**: 上方向スクロール検知のIntersection ObserverのuseEffectを排除し、refコールバック関数に置き換える

**作業内容**:
- `loadNewerRef`を`useRef`から`refコールバック関数`に変更
- `setLoadNewerRef`関数を実装（Intersection Observerを登録し、クリーンアップ関数を返す）
- `useEffect`でIntersection Observerを登録していた部分を削除
- 動作確認: 上方向スクロール時に新しい返信が正常に読み込まれることを確認

**受け入れ基準**:
- Intersection Observer関連の`useEffect`が排除されている（上方向スクロール）
- refコールバック関数が実装されている
- 上方向スクロール時に新しい返信が正常に読み込まれる
- クリーンアップが正しく実行される

_Requirements: 6.2, Design: 2.1_

---

#### - [ ] タスク 2.5: `client/lib/hooks/use-intersection-observer.ts`のuseEffect排除
**目的**: カスタムフックのuseEffectを排除し、refコールバック関数を返す関数に変更する

**作業内容**:
- `useEffect`のインポートを削除
- `useRef`を使用して、`observerRef`を作成
- `setElementRef`関数を実装（Intersection Observerを登録し、クリーンアップ関数を返す）
- `useEffect`でIntersection Observerを登録していた部分を削除
- フックの戻り値を変更（refコールバック関数を返す、またはentryとsetElementRefの両方を返す）
- このフックを使用しているコンポーネントを確認し、必要に応じて更新
- 動作確認: このフックを使用しているコンポーネントが正常に動作することを確認

**受け入れ基準**:
- `useEffect`が排除されている
- refコールバック関数が実装されている
- このフックを使用しているコンポーネントが正常に動作する
- クリーンアップが正しく実行される

_Requirements: 6.2, Design: 2.1_

---

### Phase 3: イベントリスナーの置き換え

#### - [ ] タスク 3.1: `client/lib/hooks/use-scroll.ts`のuseEffect排除
**目的**: スクロールイベントリスナーのuseEffectを排除し、refコールバック関数に置き換える

**作業内容**:
- `useEffect`のインポートを削除
- `useRef`を使用して、`listenerRef`を作成
- `setScrollRef`関数を実装（スクロールイベントリスナーを登録し、クリーンアップ関数を返す）
- `useEffect`でイベントリスナーを登録していた部分を削除
- フックの戻り値を変更（`scrolled`と`setScrollRef`を返す）
- このフックを使用しているコンポーネントを確認し、必要に応じて更新
- 動作確認: このフックを使用しているコンポーネントが正常に動作することを確認

**受け入れ基準**:
- `useEffect`が排除されている
- refコールバック関数が実装されている
- このフックを使用しているコンポーネントが正常に動作する
- クリーンアップが正しく実行される

_Requirements: 6.3, Design: 3.1_

---

#### - [ ] タスク 3.2: `client/lib/hooks/use-media-query.ts`のuseEffect排除
**目的**: リサイズイベントリスナーのuseEffectを排除し、refコールバック関数に置き換える

**作業内容**:
- `useEffect`のインポートを削除
- `useRef`を使用して、`listenerRef`を作成
- `setMediaQueryRef`関数を実装（リサイズイベントリスナーを登録し、クリーンアップ関数を返す）
- 初期デバイス判定を`setMediaQueryRef`内で実行
- `useEffect`でイベントリスナーを登録していた部分を削除
- フックの戻り値を変更（`device`、`width`、`height`、`isMobile`、`isTablet`、`isDesktop`、`setMediaQueryRef`を返す）
- このフックを使用しているコンポーネントを確認し、必要に応じて更新
- 動作確認: このフックを使用しているコンポーネントが正常に動作することを確認

**受け入れ基準**:
- `useEffect`が排除されている
- refコールバック関数が実装されている
- このフックを使用しているコンポーネントが正常に動作する
- クリーンアップが正しく実行される

_Requirements: 6.3, Design: 3.1_

---

### Phase 4: localStorage読み込みの置き換え

#### - [ ] タスク 4.1: `client/lib/hooks/use-local-storage.ts`のuseEffect排除
**目的**: localStorage読み込みのuseEffectを排除し、初期値の直接読み込みに置き換える

**作業内容**:
- `useEffect`のインポートを削除
- `useState`の初期値関数内でlocalStorageから値を読み込む（SSRを考慮して`typeof window !== 'undefined'`でチェック）
- `useEffect`でlocalStorageから値を読み込んでいた部分を削除
- エラーハンドリングを追加（try-catch）
- このフックを使用しているコンポーネントを確認し、必要に応じて更新
- 動作確認: このフックを使用しているコンポーネントが正常に動作することを確認

**受け入れ基準**:
- `useEffect`が排除されている
- 初期値が直接localStorageから読み込まれる
- SSR時にエラーが発生しない
- このフックを使用しているコンポーネントが正常に動作する

_Requirements: 6.4, Design: 4.1_

---

### Phase 5: その他のuseEffectの置き換え

#### - [ ] タスク 5.1: `client/components/shared/counting-numbers.tsx`のuseEffect排除
**目的**: アニメーション処理のuseEffectを排除し、value変更時の直接処理に置き換える

**作業内容**:
- `useEffect`のインポートを削除
- `useRef`を使用して、`animationRef`と`previousValueRef`を作成
- `startAnimation`関数を実装（requestAnimationFrameを使用）
- `useEffect`でアニメーションを開始していた部分を削除
- render内で`previousValueRef.current !== value`をチェックし、変更時に`startAnimation`を呼び出す
- 動作確認: コンポーネントが正常に動作し、アニメーションが正しく表示されることを確認

**受け入れ基準**:
- `useEffect`が排除されている
- valueが変更されたときにアニメーションが開始される
- アニメーションが正しく表示される
- クリーンアップが正しく実行される（requestAnimationFrameのキャンセル）

_Requirements: 6.5, Design: 5.1_

---

#### - [ ] タスク 5.2: `client/app/dm_movie/upload/page.tsx`のuseEffect排除
**目的**: Uppyインスタンス作成のuseEffectを排除し、refコールバック関数に置き換える

**作業内容**:
- `useEffect`のインポートを削除（他のuseEffectがない場合）
- `useRef`を使用して、`uppyRef`と`containerRef`を作成
- `setContainerRef`関数を実装（Uppyインスタンスを作成し、クリーンアップ関数を返す）
- `useEffect`でUppyインスタンスを作成していた部分を削除
- メインコンテナ要素に`ref={setContainerRef}`を設定
- 動作確認: ページ読み込み時にUppyインスタンスが正常に作成され、アップロード機能が正常に動作することを確認

**受け入れ基準**:
- `useEffect`が排除されている
- refコールバック関数が実装されている
- Uppyインスタンスが正常に作成される
- アップロード機能が正常に動作する
- クリーンアップが正しく実行される

_Requirements: 6.5, Design: 5.1_

---

### Phase 6: テストコードの更新

#### - [ ] タスク 6.1: テストコード内のuseEffect確認と排除
**目的**: テストコード内のuseEffectを確認し、絶対に必要でない限り排除する

**作業内容**:
- `client/`ディレクトリ内のテストファイルを検索（`*.test.ts`、`*.test.tsx`、`*.spec.ts`、`*.spec.tsx`）
- 各テストファイルで`useEffect`の使用を確認
- 絶対に必要でないuseEffectを排除
- 必要に応じて、refコールバック関数やイベントハンドラーに置き換え
- テストが正常に動作することを確認

**受け入れ基準**:
- テストコード内のuseEffectが排除されている（絶対に必要な場合を除く）
- テストが正常に動作する
- テストのカバレッジが維持される

_Requirements: 6.7, Design: テスト戦略_

---

### Phase 7: 動作確認と最終チェック

#### - [ ] タスク 7.1: すべてのページの動作確認
**目的**: すべてのページが正常に表示され、機能が正常に動作することを確認する

**作業内容**:
- 各ページにアクセスして、正常に表示されることを確認
- 各ページの機能（データ読み込み、フォーム送信、削除など）が正常に動作することを確認
- エラーが発生しないことを確認
- ブラウザのコンソールでエラーがないことを確認

**受け入れ基準**:
- すべてのページが正常に表示される
- すべての機能が正常に動作する
- エラーが発生しない

_Requirements: 6.8_

---

#### - [ ] タスク 7.2: 既存テストの動作確認
**目的**: 既存のテストが正常に動作することを確認する

**作業内容**:
- 既存のテストを実行
- テストが失敗する場合は、原因を調査して修正
- テストのカバレッジが維持されることを確認

**受け入れ基準**:
- 既存のテストが正常に動作する
- テストのカバレッジが維持される

_Requirements: 6.8_

---

#### - [ ] タスク 7.3: パフォーマンス確認
**目的**: パフォーマンスの劣化がないことを確認する

**作業内容**:
- ページの読み込み速度を確認
- 不要な再レンダリングが発生していないことを確認
- メモリリークが発生していないことを確認（ブラウザの開発者ツールで確認）

**受け入れ基準**:
- パフォーマンスの劣化がない
- 不要な再レンダリングが発生していない
- メモリリークが発生していない

_Requirements: 6.8_

---

#### - [ ] タスク 7.4: useEffectの使用状況の最終確認
**目的**: すべてのuseEffectが排除されているか、または保持が必要な箇所のみが残っていることを確認する

**作業内容**:
- `client/`ディレクトリ内で`useEffect`を検索
- 各useEffectの使用箇所を確認
- 保持が必要なuseEffect（例: `video-player.tsx`のクリーンアップ）以外が排除されていることを確認
- 保持が必要なuseEffectには、適切な理由があることを確認

**受け入れ基準**:
- 保持が必要なuseEffect以外が排除されている
- 保持が必要なuseEffectには、適切な理由がある

_Requirements: 6.6, 6.8_

---

## 受け入れ基準の確認

### 初期データ読み込みの置き換え
- [ ] `client/app/dm-posts/page.tsx`のuseEffectが排除されている
- [ ] `client/app/dm-users/page.tsx`のuseEffectが排除されている
- [ ] `client/app/dm-user-posts/page.tsx`のuseEffectが排除されている
- [ ] `client/app/dm_feed/[userId]/page.tsx`の初期データ読み込みのuseEffectが排除されている
- [ ] `client/app/dm_feed/[userId]/[postId]/page.tsx`の初期データ読み込みのuseEffectが排除されている
- [ ] 既存の機能が正常に動作する

### Intersection Observerの置き換え
- [ ] `client/app/dm_feed/[userId]/page.tsx`のIntersection Observer関連のuseEffectが排除されている
- [ ] `client/app/dm_feed/[userId]/[postId]/page.tsx`のIntersection Observer関連のuseEffectが排除されている
- [ ] `client/lib/hooks/use-intersection-observer.ts`のuseEffectが排除されている（または、保持が必要な場合は保持）
- [ ] 無限スクロール機能が正常に動作する

### イベントリスナーの置き換え
- [ ] `client/lib/hooks/use-scroll.ts`のuseEffectが排除されている
- [ ] `client/lib/hooks/use-media-query.ts`のuseEffectが排除されている
- [ ] 既存の機能が正常に動作する

### localStorage読み込みの置き換え
- [ ] `client/lib/hooks/use-local-storage.ts`のuseEffectが排除されている
- [ ] localStorageの読み込みが正常に動作する

### その他のuseEffectの置き換え
- [ ] `client/components/shared/counting-numbers.tsx`のuseEffectが排除されている（または、保持が必要な場合は保持）
- [ ] `client/app/dm_movie/upload/page.tsx`のuseEffectが排除されている（または、保持が必要な場合は保持）
- [ ] 既存の機能が正常に動作する

### 保持が必要なuseEffect
- [ ] `client/components/video-player/video-player.tsx`のクリーンアップ用useEffectが保持されている

### テストコードのuseEffect排除
- [ ] テストコード内のuseEffectが排除されている（または、保持が必要な場合は保持）
- [ ] テストが正常に動作する

### 動作確認
- [ ] すべてのページが正常に表示される
- [ ] すべての機能が正常に動作する
- [ ] 既存のテストが正常に動作する
- [ ] パフォーマンスの劣化がない
