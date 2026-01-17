# Next.jsのコードからuseEffect排除の要件定義書

## 1. 概要

### 1.1 プロジェクト情報
- **プロジェクト名**: go-webdb-template
- **Feature名**: 0072-rm-useeffect
- **作成日**: 2026-01-27
- **関連Issue**: https://github.com/taku-o/go-webdb-template/issues/149

### 1.2 目的
現在のNext.jsのコードからuseEffectを、絶対にuseEffectを利用しなければいけない箇所を除いて、useEffectを排除する。処理はイベントドリブンな処理などに置き換える。

### 1.3 スコープ
- クライアント側（`client/`ディレクトリ）のNext.jsコード内のuseEffectの排除
- テストコード内のuseEffectの排除（絶対に必要でない限り）
- イベントドリブンな処理への置き換え
- 絶対にuseEffectが必要な箇所の特定と保持

**本実装の範囲外**:
- サーバー側（Go）のコード
- 外部ライブラリのコード

## 2. 背景・現状分析

### 2.1 現状
- Next.js 14+ (App Router)を使用したクライアントアプリケーションが存在
- 複数のコンポーネントとカスタムフックでuseEffectが使用されている
- 主なuseEffectの使用パターン:
  1. 初期データ読み込み（ページ読み込み時のAPI呼び出し）
  2. Intersection Observer（無限スクロール実装）
  3. イベントリスナーの登録/解除（スクロール、リサイズ）
  4. localStorageの読み込み
  5. メディアクエリ（ウィンドウサイズの監視）
  6. アニメーション処理
  7. 外部ライブラリ（Uppy、HLS.js）の初期化とクリーンアップ

### 2.2 問題点
- useEffectの多用により、副作用の管理が複雑になっている
- イベントドリブンな処理がuseEffectに依存している
- コンポーネントのライフサイクルに依存した処理が多い

### 2.3 必要性
- コードの可読性と保守性の向上
- イベントドリブンな処理への移行による、より明確な処理フローの実現
- Next.js App Routerの特性を活かした実装への改善

### 2.4 実現可否
- Next.js App RouterのServer Componentsを活用可能
- イベントハンドラーによる処理の置き換えが可能
- カスタムフックの再設計が可能

## 3. 機能要件

### 3.1 初期データ読み込みの置き換え

#### 3.1.1 対象ファイル
- `client/app/dm-posts/page.tsx`: 投稿とユーザーの初期読み込み
- `client/app/dm-users/page.tsx`: ユーザーの初期読み込み
- `client/app/dm-user-posts/page.tsx`: ユーザー投稿の初期読み込み
- `client/app/dm_feed/[userId]/page.tsx`: フィード投稿の初期読み込み
- `client/app/dm_feed/[userId]/[postId]/page.tsx`: 投稿と返信の初期読み込み

#### 3.1.2 置き換え方法
- Server Componentsを活用して、サーバー側でデータを取得
- または、クライアントコンポーネントの場合は、イベントハンドラー（例: ボタンクリック）でデータを読み込む
- ページ読み込み時の自動読み込みが必要な場合は、Server Componentsを使用

### 3.2 Intersection Observerの置き換え

#### 3.2.1 対象ファイル
- `client/app/dm_feed/[userId]/page.tsx`: 無限スクロール（上方向・下方向）
- `client/app/dm_feed/[userId]/[postId]/page.tsx`: 無限スクロール（上方向・下方向）
- `client/lib/hooks/use-intersection-observer.ts`: カスタムフック

#### 3.2.2 置き換え方法
- Intersection Observerのコールバックを直接イベントハンドラーとして実装
- refのコールバック関数を使用して、要素がマウントされた時点でObserverを登録
- または、ボタンクリックなどの明示的なイベントで読み込みをトリガー

### 3.3 イベントリスナーの置き換え

#### 3.3.1 対象ファイル
- `client/lib/hooks/use-scroll.ts`: スクロールイベント
- `client/lib/hooks/use-media-query.ts`: リサイズイベント

#### 3.3.2 置き換え方法
- イベントハンドラーを直接コンポーネントに実装
- または、refのコールバック関数を使用して、要素がマウントされた時点でイベントリスナーを登録
- CSSメディアクエリを活用して、JavaScriptでの監視を不要にする（可能な場合）

### 3.4 localStorage読み込みの置き換え

#### 3.4.1 対象ファイル
- `client/lib/hooks/use-local-storage.ts`: localStorageの読み込み

#### 3.4.2 置き換え方法
- 初期値を直接localStorageから読み込む（SSRを考慮）
- または、イベントハンドラーで読み込みをトリガー
- または、refのコールバック関数を使用して、要素がマウントされた時点で読み込み

### 3.5 その他のuseEffectの置き換え

#### 3.5.1 対象ファイル
- `client/components/shared/counting-numbers.tsx`: アニメーション処理
- `client/app/dm_movie/upload/page.tsx`: Uppyインスタンスの作成

#### 3.5.2 置き換え方法
- アニメーション処理: requestAnimationFrameを直接使用、またはCSSアニメーションに置き換え
- Uppyインスタンス: refのコールバック関数を使用して、要素がマウントされた時点で作成

### 3.6 保持が必要なuseEffect

#### 3.6.1 対象ファイル
- `client/components/video-player/video-player.tsx`: HLS.jsのクリーンアップ

#### 3.6.2 保持理由
- コンポーネントのアンマウント時にリソースをクリーンアップする必要がある
- これはReactのライフサイクルに依存する処理であり、useEffectが適切

## 4. 非機能要件

### 4.1 UI/UX
- 既存の機能が正常に動作することを維持
- パフォーマンスの劣化を避ける
- ユーザー体験に影響を与えない

### 4.2 パフォーマンス
- イベントドリブンな処理への置き換えにより、不要な再レンダリングを避ける
- Server Componentsを活用して、初期読み込みのパフォーマンスを向上

### 4.3 保守性
- コードの可読性を向上
- 処理フローを明確にする
- 既存のプロジェクト構造に沿った実装

### 4.4 互換性
- 既存の機能との互換性を保つ
- 既存のAPIとの互換性を保つ
- 既存のテストが正常に動作することを維持

## 5. 制約事項

### 5.1 技術的制約
- Next.js App Routerの特性を考慮
- Server ComponentsとClient Componentsの使い分け
- SSR（Server-Side Rendering）を考慮した実装

### 5.2 実装上の制約
- 既存のプロジェクト構造に沿って実装する
- 既存のAPIとの互換性を保つ
- 既存のテストが正常に動作することを維持

### 5.3 機能の制約
- 絶対にuseEffectが必要な箇所（例: クリーンアップ処理）は保持する
- テストコードも、絶対に必要でない限りuseEffectを排除する

## 6. 受け入れ基準

### 6.1 初期データ読み込みの置き換え
- [ ] `client/app/dm-posts/page.tsx`のuseEffectが排除されている
- [ ] `client/app/dm-users/page.tsx`のuseEffectが排除されている
- [ ] `client/app/dm-user-posts/page.tsx`のuseEffectが排除されている
- [ ] `client/app/dm_feed/[userId]/page.tsx`の初期データ読み込みのuseEffectが排除されている
- [ ] `client/app/dm_feed/[userId]/[postId]/page.tsx`の初期データ読み込みのuseEffectが排除されている
- [ ] 既存の機能が正常に動作する

### 6.2 Intersection Observerの置き換え
- [ ] `client/app/dm_feed/[userId]/page.tsx`のIntersection Observer関連のuseEffectが排除されている
- [ ] `client/app/dm_feed/[userId]/[postId]/page.tsx`のIntersection Observer関連のuseEffectが排除されている
- [ ] `client/lib/hooks/use-intersection-observer.ts`のuseEffectが排除されている（または、保持が必要な場合は保持）
- [ ] 無限スクロール機能が正常に動作する

### 6.3 イベントリスナーの置き換え
- [ ] `client/lib/hooks/use-scroll.ts`のuseEffectが排除されている
- [ ] `client/lib/hooks/use-media-query.ts`のuseEffectが排除されている
- [ ] 既存の機能が正常に動作する

### 6.4 localStorage読み込みの置き換え
- [ ] `client/lib/hooks/use-local-storage.ts`のuseEffectが排除されている
- [ ] localStorageの読み込みが正常に動作する

### 6.5 その他のuseEffectの置き換え
- [ ] `client/components/shared/counting-numbers.tsx`のuseEffectが排除されている（または、保持が必要な場合は保持）
- [ ] `client/app/dm_movie/upload/page.tsx`のuseEffectが排除されている（または、保持が必要な場合は保持）
- [ ] 既存の機能が正常に動作する

### 6.6 保持が必要なuseEffect
- [ ] `client/components/video-player/video-player.tsx`のクリーンアップ用useEffectが保持されている

### 6.7 テストコードのuseEffect排除
- [ ] テストコード内のuseEffectが排除されている（または、保持が必要な場合は保持）
- [ ] テストが正常に動作する

### 6.8 動作確認
- [ ] すべてのページが正常に表示される
- [ ] すべての機能が正常に動作する
- [ ] 既存のテストが正常に動作する
- [ ] パフォーマンスの劣化がない

## 7. 影響範囲

### 7.1 変更されるファイル
- `client/app/dm-posts/page.tsx`
- `client/app/dm-users/page.tsx`
- `client/app/dm-user-posts/page.tsx`
- `client/app/dm_feed/[userId]/page.tsx`
- `client/app/dm_feed/[userId]/[postId]/page.tsx`
- `client/lib/hooks/use-scroll.ts`
- `client/lib/hooks/use-media-query.ts`
- `client/lib/hooks/use-local-storage.ts`
- `client/lib/hooks/use-intersection-observer.ts`
- `client/components/shared/counting-numbers.tsx`
- `client/app/dm_movie/upload/page.tsx`
- テストコード（useEffectが使用されている場合）

### 7.2 保持されるファイル
- `client/components/video-player/video-player.tsx`: クリーンアップ用useEffectを保持

### 7.3 既存ファイルへの影響
- 既存のAPIとの互換性を保つ
- 既存のテストが正常に動作することを維持

### 7.4 既存機能への影響
- 既存の機能が正常に動作することを維持
- ユーザー体験に影響を与えない

## 8. 実装上の注意事項

### 8.1 Server Componentsの活用
- Next.js App RouterのServer Componentsを活用して、初期データ読み込みをサーバー側で実行
- Client Componentsが必要な場合は、イベントハンドラーでデータを読み込む

### 8.2 イベントドリブンな処理への置き換え
- useEffectで実装されていた処理を、イベントハンドラーやrefのコールバック関数に置き換え
- 処理フローを明確にする

### 8.3 refのコールバック関数の活用
- refのコールバック関数を使用して、要素がマウントされた時点で処理を実行
- これにより、useEffectを使わずに初期化処理を実装可能

### 8.4 クリーンアップ処理の保持
- コンポーネントのアンマウント時にリソースをクリーンアップする必要がある場合は、useEffectを保持

### 8.5 テスト
- 既存のテストが正常に動作することを確認
- テストコードも、絶対に必要でない限りuseEffectを排除する

### 8.6 パフォーマンス
- イベントドリブンな処理への置き換えにより、不要な再レンダリングを避ける
- Server Componentsを活用して、初期読み込みのパフォーマンスを向上

## 9. 参考情報

### 9.1 関連ドキュメント
- Next.js App Routerドキュメント
- React Hooksドキュメント
- 既存のプロジェクトドキュメント

### 9.2 関連Issue
- https://github.com/taku-o/go-webdb-template/issues/149: 本要件定義書の元となったIssue

### 9.3 技術スタック
- **フレームワーク**: Next.js 14+ (App Router)
- **言語**: TypeScript 5+
- **UIライブラリ**: shadcn/ui
- **スタイリング**: Tailwind CSS

### 9.4 実装の流れ
1. 現在のuseEffectの使用状況を確認
2. 各useEffectの置き換え方法を検討
3. 絶対にuseEffectが必要な箇所を特定
4. イベントドリブンな処理への置き換えを実施
5. 動作確認とテスト

### 9.5 依存関係
- 既存のNext.jsプロジェクト構造
- 既存のAPI
- 既存のコンポーネントとカスタムフック
