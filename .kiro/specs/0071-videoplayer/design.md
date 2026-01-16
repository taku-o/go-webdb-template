# 動画プレイヤーコンポーネントとデモページの設計書

## Overview

### 目的
Twitter風のフィードUIで使用する想定の動画プレイヤーコンポーネントを作成し、そのコンポーネントを表示するデモページを実装する。動画はHLS配信を想定し、HLS.jsを利用してHLSをサポートしないブラウザでも再生可能にする。UIはvideoタグをベースに実装し、必要に応じてカスタムコントロールを追加する。なお、Twitter風のフィードUI自体は別で実装する想定であり、本実装では単体の動画プレイヤーコンポーネントのデモページのみを作成する。

### ユーザー
- **開発者**: 動画プレイヤーコンポーネントの動作を確認する
- **デモンストレーション**: 動画プレイヤーコンポーネントの実装例を確認する

### 影響
現在のシステム状態を以下のように変更する：
- `client/components/video-player/video-player.tsx`: 新規作成（動画プレイヤーコンポーネント）
- `client/app/dm_videoplayer/page.tsx`: 新規作成（デモページ）
- `client/app/page.tsx`: 「動画プレイヤー」リンクを追加
- `.gitignore`: `client/public/demo-videos/`を追加
- `client/public/demo-videos/`: 新規作成（デモ用ファイル配置用、gitにコミットしない）
- `client/package.json`: 依存関係を追加（`hls.js`）

### Goals
- 動画プレイヤーコンポーネントを作成する（HLS.jsとvideoタグを利用）
- デモページを作成する（単体の動画プレイヤーコンポーネントを表示）
- 遅延読み込み機能を実装する（再生ボタンを押すまで動画をダウンロードしない）
- HLS配信に対応する（HLSをサポートするブラウザとサポートしないブラウザの両方に対応）
- videoタグをベースにしたUIを実装する
- レスポンシブデザインに対応する

### Non-Goals
- Twitter風のフィードUIの実装（別で実装する想定）
- 動画のアップロード機能（既存の動画アップロード機能を使用）
- 動画のエンコード機能
- 動画の配信機能（既存の動画ファイルを使用）
- データベース連携（デモ用の静的ファイルのみ）
- 認証機能の新規実装（既存の認証機能を使用）
- コンポーネントの単体テスト
- E2Eテスト

## Architecture

### ディレクトリ構造

#### 新規作成されるディレクトリとファイル
```
client/
├── components/
│   └── video-player/
│       └── video-player.tsx          # 動画プレイヤーコンポーネント
├── app/
│   ├── dm_videoplayer/
│   │   └── page.tsx                   # デモページ
│   └── page.tsx                       # トップページ（リンク追加）
├── public/
│   └── demo-videos/                   # デモ用ファイル（gitにコミットしない）
│       ├── mini-movie-m.mp4
│       └── mini-movie-m.png
└── package.json                        # 依存関係追加
```

### コンポーネント設計

#### 動画プレイヤーコンポーネント (`client/components/video-player/video-player.tsx`)
```
VideoPlayer
├── Props
│   ├── videoUrl: string               # HLSまたはMP4のURL
│   ├── thumbnailUrl: string           # サムネイル画像のURL
│   └── className?: string             # 追加のCSSクラス（任意）
├── State
│   ├── hls: Hls | null                 # HLS.jsインスタンス（表示用）
│   └── error: string | null            # エラーメッセージ
├── Refs
│   ├── videoRef: HTMLVideoElement      # video要素の参照
│   └── hlsRef: Hls | null              # HLS.jsインスタンスの参照（クリーンアップ用）
├── HLS.js統合
│   ├── HLSサポート判定（useMemo）
│   ├── HLS.js初期化（再生ボタンクリック時）
│   └── HLS.jsクリーンアップ（アンマウント時）
├── videoタグ実装
│   ├── video要素
│   ├── poster属性（サムネイル表示）
│   ├── controls属性（標準コントロール）
│   └── preload="none"（遅延読み込み）
└── 遅延読み込み
    ├── preload="none"
    └── 再生ボタンクリック時にHLS.jsを初期化
```

#### デモページ (`client/app/dm_videoplayer/page.tsx`)
```
VideoPlayerDemoPage
├── VideoPlayerコンポーネント
│   ├── videoUrl: "/demo-videos/mini-movie-m.mp4"
│   └── thumbnailUrl: "/demo-videos/mini-movie-m.png"
└── レスポンシブレイアウト
```

## 詳細設計

### 1. 動画プレイヤーコンポーネント

#### 1.1 コンポーネントの基本構造

**ファイル**: `client/components/video-player/video-player.tsx`

```typescript
'use client'

import { useState, useEffect, useRef, useMemo } from 'react'
import Hls from 'hls.js'

interface VideoPlayerProps {
  videoUrl: string
  thumbnailUrl: string
  className?: string
}

export function VideoPlayer({ videoUrl, thumbnailUrl, className }: VideoPlayerProps) {
  const [hls, setHls] = useState<Hls | null>(null)
  const [error, setError] = useState<string | null>(null)
  const videoRef = useRef<HTMLVideoElement>(null)
  const hlsRef = useRef<Hls | null>(null) // HLS.jsインスタンスを保持（クリーンアップ用）

  // HLSサポート判定（useMemoで一度だけ判定）
  const isHlsSupported = useMemo(() => {
    if (typeof document === 'undefined') return false
    const video = document.createElement('video')
    return !!video.canPlayType('application/vnd.apple.mpegurl')
  }, [])

  // HLS.jsの初期化（再生ボタンクリック時）
  const handlePlay = () => {
    const video = videoRef.current
    if (!video) return

    // 既に初期化済みの場合は何もしない
    if (hlsRef.current) return

    // HLSをサポートするブラウザでは、videoタグのsrc属性にHLSのURLを直接設定
    if (isHlsSupported) {
      video.src = videoUrl
      // エラーハンドリングを設定
      video.onerror = () => {
        setError('動画の読み込みに失敗しました')
      }
      return
    }

    // HLSをサポートしないブラウザでは、HLS.jsを使用
    if (Hls.isSupported()) {
      const hlsInstance = new Hls()
      hlsInstance.loadSource(videoUrl)
      hlsInstance.attachMedia(video)

      hlsInstance.on(Hls.Events.ERROR, (event, data) => {
        if (data.fatal) {
          switch (data.type) {
            case Hls.ErrorTypes.NETWORK_ERROR:
              setError('ネットワークエラーが発生しました')
              break
            case Hls.ErrorTypes.MEDIA_ERROR:
              setError('メディアエラーが発生しました')
              break
            default:
              setError('動画の読み込みに失敗しました')
              break
          }
        }
      })

      // video要素のエラーハンドリングも設定
      video.onerror = () => {
        setError('動画の読み込みに失敗しました')
      }

      hlsRef.current = hlsInstance
      setHls(hlsInstance)
    } else {
      // HLS.jsもサポートしない場合は、通常のMP4として扱う
      video.src = videoUrl
      // エラーハンドリングを設定
      video.onerror = () => {
        setError('動画の読み込みに失敗しました')
      }
    }
  }

  // HLS.jsのクリーンアップ（コンポーネントアンマウント時のみ必要）
  useEffect(() => {
    return () => {
      if (hlsRef.current) {
        hlsRef.current.destroy()
        hlsRef.current = null
      }
    }
  }, []) // 依存配列を空にして、アンマウント時のみ実行

  return (
    <div className={`relative w-full ${className || ''}`}>
      {error && (
        <div className="absolute inset-0 flex items-center justify-center bg-black bg-opacity-50 text-white z-10 rounded">
          {error}
        </div>
      )}
      <video
        ref={videoRef}
        className="w-full h-auto rounded"
        poster={thumbnailUrl}
        controls
        preload="none"
        onPlay={handlePlay}
        playsInline
      >
        <source src={videoUrl} type={isHlsSupported ? 'application/vnd.apple.mpegurl' : 'video/mp4'} />
        お使いのブラウザは動画タグをサポートしていません。
      </video>
    </div>
  )
}
```

#### 1.2 HLS対応の実装

**HLSサポート判定**:
- `useMemo`を使用してHLSをサポートするブラウザかどうかを一度だけ判定
- `document.createElement('video')`でvideo要素を作成し、`canPlayType('application/vnd.apple.mpegurl')`で判定
- `useEffect`は使用しない

**HLS.jsの初期化（遅延読み込み対応）**:
- 再生ボタンクリック時（`onPlay`イベント）にHLS.jsを初期化
- 既に初期化済みの場合は何もしない（`hlsRef.current`で重複初期化を防ぐ）
- HLSをサポートするブラウザでは、videoタグのsrc属性にHLSのURLを直接設定
- HLSをサポートしないブラウザでは、HLS.jsを使用:
  - `Hls.isSupported()`でHLS.jsがサポートされているか確認
  - `new Hls()`でHLS.jsインスタンスを作成
  - `hlsInstance.loadSource(videoUrl)`で動画ソースを読み込み
  - `hlsInstance.attachMedia(video)`でvideo要素にアタッチ
- HLS.jsインスタンスは`hlsRef`で保持（クリーンアップ用）

**エラーハンドリング**:
- video要素の`onerror`イベントハンドラーを`handlePlay`内で設定
- HLS.jsのエラーハンドリングも`handlePlay`内で設定
- `useEffect`は使用しない

**HLS.jsのクリーンアップ**:
- コンポーネントアンマウント時のみ`useEffect`を使用（必要最小限）
- `hlsRef.current`を参照して`hlsInstance.destroy()`を呼び出し
- メモリリークを防ぐ

#### 1.3 videoタグの実装

**video要素**:
- HTML5の`<video>`タグを直接使用
- `ref`でvideo要素の参照を取得
- `poster`属性にサムネイル画像のURLを設定
- `controls`属性で標準コントロールを表示
- `preload="none"`で遅延読み込みを実装
- `playsInline`属性を設定（モバイルでのインライン再生対応）

**source要素**:
- `<source>`要素で動画ソースを指定
- `src`属性に動画URLを設定
- `type`属性にMIMEタイプを設定（HLSの場合は`application/vnd.apple.mpegurl`、MP4の場合は`video/mp4`）

**遅延読み込み**:
- `preload="none"`を設定して、動画を再生ボタンを押すまでダウンロードしない
- サムネイル画像を`poster`属性で表示
- 再生ボタンをクリックしたら`handlePlay`が呼び出され、HLS.jsを初期化して動画を読み込む

**スタイリング**:
- Tailwind CSSを使用してスタイリング
- `w-full h-auto`でレスポンシブ対応
- `rounded`で角丸を適用
- 必要に応じて追加のCSSクラスを`className`プロパティで指定可能

### 2. デモページ

#### 2.1 ページの基本構造

**ファイル**: `client/app/dm_videoplayer/page.tsx`

```typescript
import { VideoPlayer } from '@/components/video-player/video-player'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'

export default function VideoPlayerDemoPage() {
  return (
    <main className="min-h-screen p-4 sm:p-6 md:p-8">
      <div className="max-w-4xl mx-auto">
        <Card>
          <CardHeader>
            <CardTitle>動画プレイヤー</CardTitle>
            <CardDescription>
              動画プレイヤーコンポーネントのデモページです。
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="w-full max-w-3xl mx-auto">
              <VideoPlayer
                videoUrl="/demo-videos/mini-movie-m.mp4"
                thumbnailUrl="/demo-videos/mini-movie-m.png"
              />
            </div>
          </CardContent>
        </Card>
      </div>
    </main>
  )
}
```

#### 2.2 レスポンシブデザイン

- `max-w-4xl mx-auto`: 最大幅を4xlに設定し、中央揃え
- `p-4 sm:p-6 md:p-8`: レスポンシブなパディング
- `w-full max-w-3xl mx-auto`: 動画プレイヤーの幅を制限し、中央揃え

### 3. トップページへのリンク追加

#### 3.1 リンクの追加

**ファイル**: `client/app/page.tsx`

`features`配列に以下を追加:

```typescript
{
  title: '動画プレイヤー',
  description: '動画プレイヤーコンポーネントのデモ',
  href: '/dm_videoplayer',
}
```

### 4. デモ用ファイルの配置

#### 4.1 ファイルのコピー

- `~/Desktop/movie/mini-movie-m.mp4`を`client/public/demo-videos/mini-movie-m.mp4`にコピー
- `~/Desktop/movie/mini-movie-m.png`を`client/public/demo-videos/mini-movie-m.png`にコピー

#### 4.2 .gitignoreの更新

**ファイル**: `.gitignore`

以下を追加:

```
# Demo videos (not committed to git)
client/public/demo-videos/
```

### 5. 依存関係の追加

#### 5.1 パッケージのインストール

```bash
npm install hls.js
```

#### 5.2 型定義のインストール（必要に応じて）

```bash
npm install --save-dev @types/hls.js
```

## データフロー

### 動画プレイヤーコンポーネントの動作フロー

```
1. コンポーネントマウント
   ↓
2. HLSサポート判定
   ├─ HLSをサポートするブラウザ
   │  └─ video.src = videoUrl（直接設定）
   └─ HLSをサポートしないブラウザ
      ├─ HLS.jsがサポートされている
      │  └─ HLS.jsインスタンスを作成・初期化
      └─ HLS.jsもサポートしない
         └─ video.src = videoUrl（MP4として扱う）
   ↓
3. サムネイル画像を表示（poster属性）
   ↓
4. ユーザーが再生ボタンをクリック
   ↓
5. 動画を読み込み（preload="none"により、この時点で読み込み開始）
   ↓
6. 動画を再生
   ↓
7. コンポーネントアンマウント時
   └─ HLS.jsインスタンスを破棄（クリーンアップ）
```

## エラーハンドリング

### エラーケース

1. **ネットワークエラー**
   - HLS.jsの`Hls.ErrorTypes.NETWORK_ERROR`を検出
   - エラーメッセージを表示

2. **メディアエラー**
   - HLS.jsの`Hls.ErrorTypes.MEDIA_ERROR`を検出
   - エラーメッセージを表示

3. **動画の読み込みエラー**
   - video要素の`error`イベントを検出
   - エラーメッセージを表示

### エラー表示

- エラーメッセージを動画プレイヤーの上にオーバーレイ表示
- ユーザーフレンドリーなメッセージを表示

## パフォーマンス考慮事項

### 遅延読み込み

- `preload="none"`により、初期ページ読み込み時に動画をダウンロードしない
- 再生ボタンをクリックした時点で動画を読み込む
- ページ読み込み時のパフォーマンスを維持

### メモリ管理

- HLS.jsインスタンスを適切にクリーンアップ
- コンポーネントアンマウント時に`hlsInstance.destroy()`を呼び出し
- メモリリークを防ぐ

### 複数インスタンス対応

- 複数の動画プレイヤーが同時に存在しても問題なく動作する設計
- 各インスタンスが独立して動作する
- 将来のTwitter風フィードUIでの使用を想定

## ブラウザ互換性

### 対応ブラウザ

- **Chrome**: HLS.jsを使用
- **Firefox**: HLS.jsを使用
- **Safari**: ネイティブHLSサポート
- **Edge**: HLS.jsを使用

### フォールバック

- HLS.jsがサポートされない場合は、通常のMP4として扱う
- エラーが発生した場合は、エラーメッセージを表示

## 実装上の注意事項

### 1. HLS.jsの初期化タイミング

- `useEffect`内でHLS.jsを初期化
- 依存配列に`videoUrl`と`isHlsSupported`を含める
- クリーンアップ関数でHLS.jsインスタンスを破棄

### 2. videoタグの実装

- HTML5の`<video>`タグを直接使用
- `ref`でvideo要素の参照を取得
- `poster`属性でサムネイル画像を表示
- `controls`属性で標準コントロールを表示
- `preload="none"`で遅延読み込みを実装
- `onPlay`イベントハンドラーで再生ボタンクリック時にHLS.jsを初期化
- Tailwind CSSでスタイリング

### 3. 型定義

- TypeScriptの型定義を適切に使用
- `Hls`型を使用
- 必要に応じて`@types/hls.js`をインストール

### 4. クリーンアップ

- `useEffect`のクリーンアップ関数でリソースを解放
- HLS.jsインスタンスの破棄
- イベントリスナーの削除

## テスト戦略

### 手動テスト

- 各種ブラウザ（Chrome、Firefox、Safari、Edge）で動作確認
- 動画の再生・停止・一時停止が正常に動作することを確認
- サムネイル画像が表示されることを確認
- 遅延読み込みが正常に動作することを確認
- エラーハンドリングが正常に動作することを確認

### テスト項目

1. **基本動作**
   - 動画が正常に再生される
   - サムネイル画像が表示される
   - 再生ボタンが正常に動作する

2. **HLS対応**
   - HLSをサポートするブラウザで正常に動作する
   - HLSをサポートしないブラウザでHLS.jsが正常に動作する

3. **遅延読み込み**
   - 初期ページ読み込み時に動画がダウンロードされない
   - 再生ボタンをクリックした時点で動画が読み込まれる

4. **エラーハンドリング**
   - ネットワークエラーが適切に表示される
   - メディアエラーが適切に表示される

5. **レスポンシブデザイン**
   - モバイル、タブレット、デスクトップで正常に表示される

## 参考情報

### 関連ドキュメント

- HLS.jsドキュメント: https://github.com/video-dev/hls.js/
- HTML5 videoタグドキュメント: https://developer.mozilla.org/ja/docs/Web/HTML/Element/video
- Next.js App Routerドキュメント: https://nextjs.org/docs/app

### 関連Issue

- https://github.com/taku-o/go-webdb-template/issues/147: 本設計書の元となったIssue

### 技術スタック

- **フレームワーク**: Next.js 14+ (App Router)
- **言語**: TypeScript 5+
- **UIライブラリ**: shadcn/ui
- **スタイリング**: Tailwind CSS
- **動画プレイヤー**: HTML5 videoタグ
- **HLS対応**: HLS.js
