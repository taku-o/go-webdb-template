'use client'

import { useState, useEffect } from 'react'
import Link from 'next/link'
import Uppy from '@uppy/core'
import Tus from '@uppy/tus'
import Dashboard from '@uppy/react/dashboard'
import '@uppy/core/css/style.min.css'
import '@uppy/dashboard/css/style.min.css'

export default function MovieUploadPage() {
  const [uppy, setUppy] = useState<Uppy | null>(null)
  const [uploadProgress, setUploadProgress] = useState(0)
  const [uploadStatus, setUploadStatus] = useState<'idle' | 'uploading' | 'success' | 'error'>('idle')
  const [errorMessage, setErrorMessage] = useState<string | null>(null)

  useEffect(() => {
    // uppyインスタンスの作成
    const uppyInstance = new Uppy({
      id: 'dm_movie_uploader',
      autoProceed: false,
      restrictions: {
        maxFileSize: 2147483648, // 2GB
        allowedFileTypes: ['.mp4'],
        maxNumberOfFiles: 1,
      },
    })
      .use(Tus, {
        endpoint: (process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8080') + '/api/upload/dm_movie',
        chunkSize: 5 * 1024 * 1024, // 5MB
        retryDelays: [0, 1000, 3000, 5000],
        headers: {
          Authorization: `Bearer ${process.env.NEXT_PUBLIC_API_KEY}`,
        },
      })
      .on('upload-progress', (file, progress) => {
        if (progress.bytesTotal) {
          const percent = Math.round((progress.bytesUploaded / progress.bytesTotal) * 100)
          setUploadProgress(percent)
        }
      })
      .on('upload-success', (file, response) => {
        setUploadStatus('success')
        setUploadProgress(100)
      })
      .on('upload-error', (file, error, response) => {
        setUploadStatus('error')
        setErrorMessage(error?.message || 'アップロードに失敗しました')
      })
      .on('upload', () => {
        setUploadStatus('uploading')
        setUploadProgress(0)
        setErrorMessage(null)
      })

    setUppy(uppyInstance)

    // クリーンアップ
    return () => {
      uppyInstance.destroy()
    }
  }, [])

  return (
    <main className="min-h-screen p-8">
      <div className="max-w-4xl mx-auto">
        <div className="mb-6">
          <Link href="/" className="text-blue-500 hover:underline">
            ← トップページに戻る
          </Link>
        </div>

        <h1 className="text-2xl font-bold mb-4">動画ファイルアップロード</h1>
        <p className="text-gray-600 mb-8">
          MP4形式の動画ファイルをアップロードできます。最大ファイルサイズは2GBです。
        </p>

        {uppy && (
          <Dashboard
            uppy={uppy}
            proudlyDisplayPoweredByUppy={false}
            height={400}
            width="100%"
          />
        )}

        {uploadStatus === 'uploading' && (
          <div className="mt-4">
            <div className="w-full bg-gray-200 rounded-full h-4">
              <div
                className="bg-blue-600 h-4 rounded-full transition-all duration-300"
                style={{ width: `${uploadProgress}%` }}
              ></div>
            </div>
            <p className="text-sm text-gray-600 mt-2">{uploadProgress}% アップロード済み</p>
          </div>
        )}

        {uploadStatus === 'success' && (
          <div className="mt-4 p-4 bg-green-100 text-green-700 rounded">
            アップロードが完了しました。
          </div>
        )}

        {uploadStatus === 'error' && errorMessage && (
          <div className="mt-4 p-4 bg-red-100 text-red-700 rounded">
            {errorMessage}
          </div>
        )}
      </div>
    </main>
  )
}
