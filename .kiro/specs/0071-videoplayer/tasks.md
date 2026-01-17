# 動画プレイヤーコンポーネントとデモページの実装タスク一覧

## 概要
動画プレイヤーコンポーネントとデモページの実装を段階的に進めるためのタスク一覧。要件定義書と設計書に基づき、実装可能な粒度でタスクを分解。

## 実装フェーズ

### Phase 1: 依存関係の追加

#### - [ ] タスク 1.1: 必要なパッケージのインストール
**目的**: 動画プレイヤーコンポーネントに必要な依存関係を追加する

**作業内容**:
- `client/`ディレクトリに移動
- 以下のパッケージをインストール:
  - `hls.js` - HLS.jsライブラリ
- 型定義が必要な場合は`@types/hls.js`もインストール
- `package.json`に依存関係が追加されていることを確認

**受け入れ基準**:
- `hls.js`がインストールされている
- `package.json`に依存関係が追加されている
- `node_modules`にパッケージがインストールされている

_Requirements: 3.3, 6.4, Design: 5.1_

---

### Phase 2: 動画プレイヤーコンポーネントの作成

#### - [ ] タスク 2.1: コンポーネントディレクトリとファイルの作成
**目的**: 動画プレイヤーコンポーネントのディレクトリとファイルを作成する

**作業内容**:
- `client/components/video-player/`ディレクトリを作成
- `client/components/video-player/video-player.tsx`ファイルを作成
- 基本的なコンポーネント構造を作成（Props、State、基本的なreturn文）

**受け入れ基準**:
- `client/components/video-player/`ディレクトリが存在する
- `client/components/video-player/video-player.tsx`ファイルが存在する
- TypeScriptのコンパイルエラーが発生しない

_Requirements: 3.1.1, 6.1, Design: 1.1_

---

#### - [ ] タスク 2.2: PropsとStateの定義
**目的**: コンポーネントのPropsとStateを定義する

**作業内容**:
- `VideoPlayerProps`インターフェースを定義:
  - `videoUrl: string` - HLSまたはMP4のURL
  - `thumbnailUrl: string` - サムネイル画像のURL
  - `className?: string` - 追加のCSSクラス（任意）
- Stateを定義:
  - `hls: Hls | null` - HLS.jsインスタンス（表示用）
  - `error: string | null` - エラーメッセージ
- `useRef`で参照を取得:
  - `videoRef: HTMLVideoElement` - video要素の参照
  - `hlsRef: Hls | null` - HLS.jsインスタンスの参照（クリーンアップ用）

**受け入れ基準**:
- `VideoPlayerProps`インターフェースが定義されている
- Stateが適切に定義されている
- `useRef`でvideo要素の参照を取得している
- TypeScriptのコンパイルエラーが発生しない

_Requirements: 3.1.1, Design: 1.1_

---

#### - [ ] タスク 2.3: HLSサポート判定の実装
**目的**: HLSをサポートするブラウザかどうかを判定する

**作業内容**:
- `useMemo`を使用してHLSサポート判定を実装（`useEffect`は使用しない）
- `document.createElement('video')`でvideo要素を作成
- `video.canPlayType('application/vnd.apple.mpegurl')`を使用してHLSをサポートするブラウザかどうかを判定
- 判定結果を`isHlsSupported`として`useMemo`で一度だけ計算
- `typeof document === 'undefined'`のチェックを追加（SSR対応）

**受け入れ基準**:
- HLSサポート判定が実装されている
- `useMemo`を使用して一度だけ判定される
- `useEffect`を使用していない
- TypeScriptのコンパイルエラーが発生しない

_Requirements: 3.1.4, Design: 1.2_

---

#### - [ ] タスク 2.4: HLS.jsの統合（遅延読み込み対応）
**目的**: HLS.jsを統合してHLSをサポートしないブラウザでも動画を再生できるようにする（再生ボタンクリック時に初期化）

**作業内容**:
- `handlePlay`関数を実装（再生ボタンクリック時にHLS.jsを初期化）
- 既に初期化済みの場合は何もしない（`hlsRef.current`で重複初期化を防ぐ）
- HLSをサポートするブラウザでは、videoタグのsrc属性にHLSのURLを直接設定
- HLSをサポートしないブラウザでは、HLS.jsを使用して動画を再生
  - `Hls.isSupported()`でHLS.jsがサポートされているか確認
  - `new Hls()`でHLS.jsインスタンスを作成
  - `hlsInstance.loadSource(videoUrl)`で動画ソースを読み込み
  - `hlsInstance.attachMedia(video)`でvideo要素にアタッチ
- HLS.jsのエラーハンドリングを実装（`handlePlay`内で）:
  - `Hls.Events.ERROR`イベントをリッスン
  - ネットワークエラー、メディアエラーを検出してエラーメッセージを設定
- video要素のエラーハンドリングを実装（`handlePlay`内で）:
  - `video.onerror`イベントハンドラーを設定
  - エラーが発生した場合、エラーメッセージを`error`ステートに設定
- HLS.jsインスタンスを`hlsRef.current`に保存（クリーンアップ用）
- HLS.jsインスタンスを`hls`ステートにも保存（表示用、任意）
- videoタグの`onPlay`イベントハンドラーに`handlePlay`を設定

**受け入れ基準**:
- `handlePlay`関数が実装されている
- 再生ボタンクリック時にHLS.jsが初期化される
- 既に初期化済みの場合は重複初期化されない（`hlsRef.current`で判定）
- HLSをサポートするブラウザでは、videoタグのsrc属性にHLSのURLが直接設定される
- HLSをサポートしないブラウザでは、HLS.jsを使用して動画を再生できる
- HLS.jsのエラーハンドリングが実装されている（`handlePlay`内）
- video要素のエラーハンドリングが実装されている（`handlePlay`内）
- `onPlay`イベントハンドラーが設定されている
- TypeScriptのコンパイルエラーが発生しない

_Requirements: 3.1.4, 6.1, Design: 1.2_

---

#### - [ ] タスク 2.5: HLS.jsのクリーンアップ実装
**目的**: コンポーネントアンマウント時にHLS.jsインスタンスを破棄する

**作業内容**:
- `useEffect`を使用してコンポーネントアンマウント時のクリーンアップを実装（必要最小限のuseEffect）
- 依存配列を空にして、アンマウント時のみ実行されるようにする
- クリーンアップ関数で`hlsRef.current`を参照して`hlsInstance.destroy()`を呼び出し
- クリーンアップ後、`hlsRef.current = null`を設定

**受け入れ基準**:
- `useEffect`がアンマウント時のみ実行される（依存配列が空）
- クリーンアップ関数でHLS.jsインスタンスが破棄される
- `hlsRef.current`が`null`に設定される
- TypeScriptのコンパイルエラーが発生しない

**注意**: このタスクは必要最小限の`useEffect`使用（コンポーネントアンマウント時のクリーンアップのため）

_Requirements: 3.1.4, 6.1, Design: 1.2_

---

#### - [ ] タスク 2.6: videoタグの実装
**目的**: HTML5のvideoタグを使用して動画プレイヤーを実装する

**作業内容**:
- HTML5の`<video>`タグを使用して動画プレイヤーを実装
- `ref`でvideo要素の参照を取得（`videoRef`を使用）
- `poster`属性にサムネイル画像のURLを設定
- `controls`属性で標準コントロールを表示
- `preload="none"`で遅延読み込みを実装
- `playsInline`属性を設定（モバイルでのインライン再生対応）
- `<source>`要素で動画ソースを指定:
  - `src`属性に動画URLを設定
  - `type`属性にMIMEタイプを設定（HLSの場合は`application/vnd.apple.mpegurl`、MP4の場合は`video/mp4`）
- `onPlay`イベントハンドラーに`handlePlay`関数を設定（再生ボタンクリック時にHLS.jsを初期化）
- Tailwind CSSでスタイリング（`w-full h-auto rounded`など）
- エラーメッセージをオーバーレイ表示

**受け入れ基準**:
- HTML5の`<video>`タグが使用されている
- `ref`でvideo要素の参照を取得している
- `poster`属性でサムネイル画像が表示される
- `controls`属性で標準コントロールが表示される
- `preload="none"`で遅延読み込みが実装されている
- `<source>`要素で動画ソースが指定されている
- `onPlay`イベントハンドラーが設定されている
- エラーメッセージがオーバーレイ表示される
- TypeScriptのコンパイルエラーが発生しない

_Requirements: 3.1.3, 6.1, Design: 1.3_

---

#### - [ ] タスク 2.7: コンポーネントのエクスポート
**目的**: コンポーネントをエクスポートして他のファイルから使用できるようにする

**作業内容**:
- `VideoPlayer`コンポーネントを`export`する
- 必要に応じて、型定義もエクスポートする

**受け入れ基準**:
- `VideoPlayer`コンポーネントがエクスポートされている
- 他のファイルからインポートできる
- TypeScriptのコンパイルエラーが発生しない

_Requirements: 3.1.1, Design: 1.1_

---

### Phase 3: デモ用ファイルの配置

#### - [ ] タスク 3.1: デモ用ディレクトリの作成
**目的**: デモ用ファイルを配置するディレクトリを作成する

**作業内容**:
- `client/public/demo-videos/`ディレクトリを作成
- ディレクトリが存在することを確認

**受け入れ基準**:
- `client/public/demo-videos/`ディレクトリが存在する

_Requirements: 3.2.2, 6.3, Design: 4.1_

---

#### - [ ] タスク 3.2: デモ用動画ファイルのコピー
**目的**: デモ用の動画ファイルをコピーする

**作業内容**:
- `~/Desktop/movie/mini-movie-m.mp4`を`client/public/demo-videos/mini-movie-m.mp4`にコピー
- ファイルが正しくコピーされていることを確認

**受け入れ基準**:
- `client/public/demo-videos/mini-movie-m.mp4`が存在する
- ファイルサイズが正しいことを確認

_Requirements: 3.2.2, 6.3, Design: 4.1_

---

#### - [ ] タスク 3.3: デモ用サムネイル画像ファイルのコピー
**目的**: デモ用のサムネイル画像ファイルをコピーする

**作業内容**:
- `~/Desktop/movie/mini-movie-m.png`を`client/public/demo-videos/mini-movie-m.png`にコピー
- ファイルが正しくコピーされていることを確認

**受け入れ基準**:
- `client/public/demo-videos/mini-movie-m.png`が存在する
- ファイルサイズが正しいことを確認

_Requirements: 3.2.2, 6.3, Design: 4.1_

---

#### - [ ] タスク 3.4: HLSファイルのコピー（動作確認用）
**目的**: HLS動作確認用のファイルをコピーする

**作業内容**:
- `~/Desktop/movie/mini-movie-hls.m3u8`を`client/public/demo-videos/mini-movie-hls.m3u8`にコピー
- すべての`.segments`ディレクトリを`client/public/demo-videos/`にコピー:
  - `mini-movie-hls - HEVCモバイル通信（小、3G以下）.segments/`
  - `mini-movie-hls - HEVCモバイル通信（中、3G以下）.segments/`
  - `mini-movie-hls - HEVCモバイル通信（大、3G以下）.segments/`
  - `mini-movie-hls - HEVCブロードバンドHD（および4G LTE以上）.segments/`
  - `mini-movie-hls - HEVCブロードバンドUHD（および4G LTE以上）.segments/`
  - `mini-movie-hls - オーディオ（標準）.segments/`
- ファイルが正しくコピーされていることを確認

**受け入れ基準**:
- `client/public/demo-videos/mini-movie-hls.m3u8`が存在する
- すべての`.segments`ディレクトリが存在する
- 各`.segments`ディレクトリ内に必要なファイル（`prog_index.m3u8`、セグメントファイル）が存在する

**注意**: これらのファイルは動作確認用であり、デモページの最終コードでは使用しない。

_Requirements: 3.2.2, 6.3, Design: 4.1_

---

#### - [ ] タスク 3.5: .gitignoreの更新
**目的**: デモ用ファイルをgitにコミットしないようにする

**作業内容**:
- `.gitignore`ファイルを開く
- `client/public/demo-videos/`を追加
- 既存のエントリと重複していないことを確認

**受け入れ基準**:
- `.gitignore`に`client/public/demo-videos/`が追加されている
- `git status`でデモ用ファイルが表示されないことを確認

_Requirements: 3.2.2, 6.3, Design: 4.2_

---

### Phase 4: デモページの作成

#### - [ ] タスク 4.1: デモページディレクトリとファイルの作成
**目的**: デモページのディレクトリとファイルを作成する

**作業内容**:
- `client/app/dm_videoplayer/`ディレクトリを作成
- `client/app/dm_videoplayer/page.tsx`ファイルを作成
- 基本的なページ構造を作成

**受け入れ基準**:
- `client/app/dm_videoplayer/`ディレクトリが存在する
- `client/app/dm_videoplayer/page.tsx`ファイルが存在する
- TypeScriptのコンパイルエラーが発生しない

_Requirements: 3.2.1, 6.2, Design: 2.1_

---

#### - [ ] タスク 4.2: デモページの実装
**目的**: デモページを実装して動画プレイヤーコンポーネントを表示する

**作業内容**:
- `VideoPlayer`コンポーネントをインポート
- `Card`、`CardContent`、`CardDescription`、`CardHeader`、`CardTitle`をインポート（shadcn/ui）
- デモページのコンポーネントを実装:
  - ページタイトル: 「動画プレイヤー」
  - 説明文を追加
  - `VideoPlayer`コンポーネントを表示:
    - `videoUrl: "/demo-videos/mini-movie-m.mp4"`（MP4ファイルを使用）
    - `thumbnailUrl: "/demo-videos/mini-movie-m.png"`
  - レスポンシブデザインに対応（`max-w-4xl mx-auto`、`p-4 sm:p-6 md:p-8`など）

**受け入れ基準**:
- `VideoPlayer`コンポーネントがインポートされている
- デモページが実装されている
- 動画プレイヤーコンポーネントが表示される
- MP4ファイルを使用している（HLSファイルは使用しない）
- レスポンシブデザインに対応している
- TypeScriptのコンパイルエラーが発生しない

**注意**: デモページの最終コードはMP4ファイルを使用する。HLSファイルは動作確認時のみ一時的に使用する。

_Requirements: 3.2.1, 3.2.3, 6.2, Design: 2.2_

---

### Phase 5: トップページへのリンク追加

#### - [ ] タスク 5.1: トップページにリンクを追加
**目的**: トップページの機能一覧に「動画プレイヤー」リンクを追加する

**作業内容**:
- `client/app/page.tsx`を開く
- `features`配列に以下を追加:
  ```typescript
  {
    title: '動画プレイヤー',
    description: '動画プレイヤーコンポーネントのデモ',
    href: '/dm_videoplayer',
  }
  ```
- 既存のリンクと同様の形式で追加
- リンクが正しく表示されることを確認

**受け入れ基準**:
- `features`配列に「動画プレイヤー」リンクが追加されている
- トップページに「動画プレイヤー」リンクが表示される
- リンクから`/dm_videoplayer`に遷移できる
- TypeScriptのコンパイルエラーが発生しない

_Requirements: 3.2.1, 6.2, Design: 3_

---

### Phase 6: 動作確認

#### - [ ] タスク 6.1: 基本的な動作確認
**目的**: 動画プレイヤーコンポーネントとデモページが正常に動作することを確認する

**作業内容**:
- 開発サーバーを起動（`npm run dev`）
- トップページにアクセス
- 「動画プレイヤー」リンクをクリック
- デモページが表示されることを確認
- 動画プレイヤーコンポーネントが表示されることを確認
- サムネイル画像が表示されることを確認
- 再生ボタンをクリックして動画が再生されることを確認
- コントロールが正常に動作することを確認（再生、停止、音量、フルスクリーンなど）

**受け入れ基準**:
- トップページからデモページに遷移できる
- デモページが正常に表示される
- 動画プレイヤーコンポーネントが表示される
- サムネイル画像が表示される
- 動画が正常に再生される
- コントロールが正常に動作する

_Requirements: 6.1, 6.2, Design: テスト戦略_

---

#### - [ ] タスク 6.2: HLS対応の動作確認
**目的**: HLS対応が正常に動作することを確認する

**作業内容**:
- **注意**: このタスクは一時的にコードを書き換えて実施する。動作確認後は元に戻す。
- `client/app/dm_videoplayer/page.tsx`を開く
- `VideoPlayer`コンポーネントの`videoUrl`を一時的に変更:
  - 変更前: `videoUrl="/demo-videos/mini-movie-m.mp4"`
  - 変更後: `videoUrl="/demo-videos/mini-movie-hls.m3u8"`
- HLSをサポートするブラウザ（Safari等）で動作確認:
  - 動画が正常に再生されることを確認
  - HLS.jsを使用せずにネイティブで再生されることを確認
  - ネットワークタブで`.m3u8`ファイルと`.m4s`セグメントファイルがダウンロードされることを確認
- HLSをサポートしないブラウザ（Chrome、Firefox、Edge等）で動作確認:
  - 動画が正常に再生されることを確認
  - HLS.jsを使用して再生されることを確認
  - ネットワークタブで`.m3u8`ファイルと`.m4s`セグメントファイルがダウンロードされることを確認
- 動作確認後、`videoUrl`を元に戻す（MP4に戻す）

**受け入れ基準**:
- HLSをサポートするブラウザで動画が正常に再生される
- HLSをサポートしないブラウザで動画が正常に再生される（HLS.jsを使用）
- ネットワークタブでHLSファイル（`.m3u8`、`.m4s`）がダウンロードされることを確認
- エラーが発生しない
- 動作確認後、コードが元に戻されている（MP4を使用）

_Requirements: 6.1, Design: ブラウザ互換性_

---

#### - [ ] タスク 6.3: 遅延読み込みの動作確認
**目的**: 遅延読み込みが正常に動作することを確認する

**作業内容**:
- デモページを開く
- ネットワークタブを開いて動画ファイルのダウンロードを監視
- ページ読み込み時に動画ファイルがダウンロードされないことを確認
- 再生ボタンをクリックした時点で動画ファイルがダウンロードされることを確認

**受け入れ基準**:
- ページ読み込み時に動画ファイルがダウンロードされない
- 再生ボタンをクリックした時点で動画ファイルがダウンロードされる
- `preload="none"`が設定されていることを確認

_Requirements: 3.1.2, 6.1, Design: 1.3_

---

#### - [ ] タスク 6.4: レスポンシブデザインの動作確認
**目的**: レスポンシブデザインが正常に動作することを確認する

**作業内容**:
- モバイルサイズ（375px）で表示確認
- タブレットサイズ（768px）で表示確認
- デスクトップサイズ（1024px以上）で表示確認
- 各サイズで動画プレイヤーが適切に表示されることを確認
- コントロールが適切に表示されることを確認

**受け入れ基準**:
- モバイルサイズで適切に表示される
- タブレットサイズで適切に表示される
- デスクトップサイズで適切に表示される
- 各サイズでコントロールが適切に表示される

_Requirements: 4.1, 6.2, Design: 2.2_

---

#### - [ ] タスク 6.5: エラーハンドリングの動作確認
**目的**: エラーハンドリングが正常に動作することを確認する

**作業内容**:
- 存在しない動画URLを指定してエラーが表示されることを確認
- ネットワークエラーが発生した場合にエラーメッセージが表示されることを確認
- エラーメッセージがユーザーフレンドリーであることを確認

**受け入れ基準**:
- エラーが発生した場合、エラーメッセージが表示される
- エラーメッセージがユーザーフレンドリーである
- エラーが適切にハンドリングされる

_Requirements: 4.1, 6.5, Design: エラーハンドリング_

---

#### - [ ] タスク 6.6: 各種ブラウザでの動作確認
**目的**: 各種ブラウザで正常に動作することを確認する

**作業内容**:
- Chromeで動作確認
- Firefoxで動作確認
- Safariで動作確認
- Edgeで動作確認
- 各ブラウザで動画が正常に再生されることを確認
- 各ブラウザでコントロールが正常に動作することを確認

**受け入れ基準**:
- Chromeで正常に動作する
- Firefoxで正常に動作する
- Safariで正常に動作する
- Edgeで正常に動作する
- 各ブラウザでエラーが発生しない

_Requirements: 4.4, 6.5, Design: ブラウザ互換性_

---

## 受け入れ基準の確認

### 動画プレイヤーコンポーネント
- [ ] `client/components/video-player/video-player.tsx`が作成されている
- [ ] videoタグをベースに実装されている
- [ ] HLS.jsが適切に統合されている
- [ ] HLSをサポートするブラウザではそのまま動画を表示できる
- [ ] HLSをサポートしないブラウザではHLS.jsを利用して動画を表示できる
- [ ] 動画は再生ボタンを押すまで、ダウンロード開始しない（`preload="none"`）
- [ ] サムネイル画像が表示される
- [ ] videoタグが適切に実装されている
- [ ] 動画プレイヤーが表示される
- [ ] 複数の動画プレイヤーが同時に存在しても問題なく動作する

### デモページ
- [ ] `client/app/dm_videoplayer/page.tsx`が作成されている
- [ ] `/dm_videoplayer`にアクセスした際、デモページが表示される
- [ ] トップページに「動画プレイヤー」リンクが追加されている
- [ ] リンクからデモページに遷移できる
- [ ] 単体の動画プレイヤーコンポーネントが表示される
- [ ] 動画プレイヤーコンポーネントが正常に動作する
- [ ] レスポンシブデザインに対応している

### デモ用ファイルの配置
- [ ] `client/public/demo-videos/mini-movie-m.mp4`が存在する（gitにコミットされていない）
- [ ] `client/public/demo-videos/mini-movie-m.png`が存在する（gitにコミットされていない）
- [ ] `.gitignore`に`client/public/demo-videos/`が追加されている
- [ ] デモページでこれらのファイルが正しく参照されている

### 依存関係
- [ ] `hls.js`がインストールされている
- [ ] 必要な型定義がインストールされている（該当する場合）

### UI/UX
- [ ] videoタグが適切に使用されている
- [ ] レスポンシブデザインに対応している
- [ ] アクセシビリティに配慮されている
- [ ] ローディング状態が適切に表示される
- [ ] エラー状態が適切に表示される
