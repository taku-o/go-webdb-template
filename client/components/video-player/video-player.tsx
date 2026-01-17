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
