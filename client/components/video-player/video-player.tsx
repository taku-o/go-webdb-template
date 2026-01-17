'use client'

import { useState, useRef, useEffect } from 'react'
import Hls from 'hls.js'

interface VideoPlayerProps {
  videoUrl: string
  thumbnailUrl: string
  className?: string
}

export function VideoPlayer({ videoUrl, thumbnailUrl, className }: VideoPlayerProps) {
  const [error, setError] = useState<string | null>(null)
  const [isPlaying, setIsPlaying] = useState(false) // 再生開始済みフラグ（UI制御用）
  const videoRef = useRef<HTMLVideoElement>(null)
  const hlsRef = useRef<Hls | null>(null) // HLS.jsインスタンスを保持（クリーンアップ用）
  const initializedRef = useRef<boolean>(false) // 初期化済みフラグ

  // HLSファイルかどうかを判定
  const isHlsUrl = videoUrl.endsWith('.m3u8')

  // ブラウザがHLSをネイティブサポートしているかどうかを判定（再生時に呼び出す）
  const checkCanPlayHlsNatively = () => {
    const video = document.createElement('video')
    return !!video.canPlayType('application/vnd.apple.mpegurl')
  }

  // HLSファイルの場合は常にカスタム再生ボタンを表示（ハイドレーションエラー回避）
  // 実際のブラウザ判定は再生時に行う
  const needsCustomPlayButton = isHlsUrl

  // videoのエラーハンドラー
  const handleError = () => {
    setError('動画の読み込みに失敗しました')
  }

  // 動画の初期化と再生
  const initializeAndPlay = () => {
    const video = videoRef.current
    if (!video) return

    // 既に初期化済みの場合は何もしない
    if (initializedRef.current) return
    initializedRef.current = true
    setIsPlaying(true)

    // HLSファイルでない場合（MP4等）は、そのまま再生
    if (!isHlsUrl) {
      // MP4の場合は既にsourceタグで設定されているので、追加の処理は不要
      return
    }

    // HLSファイルの場合
    // ブラウザがHLSをネイティブサポートしている場合（Safari等）
    if (checkCanPlayHlsNatively()) {
      video.src = videoUrl
      video.play()
      return
    }

    // HLS.jsを使用してHLSを再生
    if (Hls.isSupported()) {
      const hlsInstance = new Hls()
      hlsInstance.loadSource(videoUrl)
      hlsInstance.attachMedia(video)

      hlsInstance.on(Hls.Events.MANIFEST_PARSED, () => {
        video.play()
      })

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

      hlsRef.current = hlsInstance
    } else {
      // HLS.jsもサポートしない場合
      setError('お使いのブラウザはHLS再生をサポートしていません')
    }
  }

  // カスタム再生ボタンのクリックハンドラー（HLS.jsが必要な場合）
  const handleCustomPlayClick = () => {
    initializeAndPlay()
  }

  // 標準の再生イベントハンドラー（MP4またはネイティブHLS対応ブラウザの場合）
  const handlePlay = () => {
    initializeAndPlay()
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

  // MIMEタイプを取得
  const mimeType = 'video/mp4'

  return (
    <div className={`relative w-full ${className || ''}`}>
      {error && (
        <div className="absolute inset-0 flex items-center justify-center bg-black bg-opacity-50 text-white z-10 rounded">
          {error}
        </div>
      )}

      {/* HLSファイルで、まだ再生していない場合はカスタム再生ボタンを表示 */}
      {needsCustomPlayButton && !isPlaying && (
        <button
          onClick={handleCustomPlayClick}
          className="absolute inset-0 flex items-center justify-center z-10 bg-transparent cursor-pointer"
          aria-label="動画を再生"
        >
          <div className="w-16 h-16 bg-black bg-opacity-60 rounded-full flex items-center justify-center hover:bg-opacity-80 transition-all">
            <svg
              className="w-8 h-8 text-white ml-1"
              fill="currentColor"
              viewBox="0 0 24 24"
            >
              <path d="M8 5v14l11-7z" />
            </svg>
          </div>
        </button>
      )}

      <video
        ref={videoRef}
        className="w-full h-auto rounded"
        poster={thumbnailUrl}
        controls={isPlaying || !needsCustomPlayButton}
        controlsList="nodownload"
        preload="none"
        onPlay={handlePlay}
        onError={handleError}
        playsInline
      >
        {/* HLSファイルの場合はsourceタグを設定しない（Chromeで再生ボタンが無効になるため） */}
        {!needsCustomPlayButton && <source src={videoUrl} type={mimeType} onError={handleError} />}
        お使いのブラウザは動画タグをサポートしていません。
      </video>
    </div>
  )
}
