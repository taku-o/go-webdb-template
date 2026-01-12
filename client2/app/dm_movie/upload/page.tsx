'use client'

import { useState, useEffect } from 'react'
import Link from 'next/link'
// @ts-ignore - Uppyの型定義が正しく解決されない場合があるため
import Uppy from '@uppy/core'
// @ts-ignore
import Dashboard from '@uppy/react/dashboard'
import '@uppy/core/css/style.min.css'
import '@uppy/dashboard/css/style.min.css'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { ArrowLeft, Upload, CheckCircle2, AlertCircle } from 'lucide-react'
import { createMovieUploader } from '@/lib/api'

export default function MovieUploadPage() {
  const [uppy, setUppy] = useState<Uppy | null>(null)
  const [uploadProgress, setUploadProgress] = useState(0)
  const [uploadStatus, setUploadStatus] = useState<'idle' | 'uploading' | 'success' | 'error'>('idle')
  const [errorMessage, setErrorMessage] = useState<string | null>(null)

  useEffect(() => {
    // uppyインスタンスの作成
    const uppyInstance = createMovieUploader({
      onUploadProgress: (percent) => {
        setUploadProgress(percent)
      },
      onUploadSuccess: () => {
        setUploadStatus('success')
        setUploadProgress(100)
      },
      onUploadError: (error) => {
        setUploadStatus('error')
        setErrorMessage(error)
      },
      onUploadStart: () => {
        setUploadStatus('uploading')
        setUploadProgress(0)
        setErrorMessage(null)
      },
    })

    setUppy(uppyInstance)

    // クリーンアップ
    return () => {
      uppyInstance.destroy()
    }
  }, [])

  return (
    <main className="min-h-screen p-4 sm:p-6 md:p-8">
      <div className="max-w-4xl mx-auto">
        <nav aria-label="パンくずリスト">
          <div className="mb-4 sm:mb-6">
            <Link href="/" className="inline-flex items-center text-primary hover:underline text-sm sm:text-base" aria-label="トップページに戻る">
              <ArrowLeft className="mr-2 h-4 w-4" aria-hidden="true" />
              トップページに戻る
            </Link>
          </div>
        </nav>

        <h1 className="text-2xl sm:text-3xl font-bold mb-4 sm:mb-6 flex items-center gap-2">
          <Upload className="h-6 w-6 sm:h-8 sm:w-8" />
          動画ファイルアップロード
        </h1>

        <Card className="mb-4 sm:mb-6">
          <CardHeader>
            <CardTitle className="text-lg sm:text-xl">動画アップロード</CardTitle>
            <CardDescription className="text-sm">
              MP4形式の動画ファイルをアップロードできます。最大ファイルサイズは2GBです。
            </CardDescription>
          </CardHeader>
          <CardContent>
            {uppy && (
              <Dashboard
                uppy={uppy}
                proudlyDisplayPoweredByUppy={false}
                height={300}
                width="100%"
              />
            )}

            {uploadStatus === 'uploading' && (
              <div className="mt-4">
                <div className="w-full bg-muted rounded-full h-4">
                  <div
                    className="bg-primary h-4 rounded-full transition-all duration-300"
                    style={{ width: `${uploadProgress}%` }}
                  ></div>
                </div>
                <p className="text-sm text-muted-foreground mt-2">{uploadProgress}% アップロード済み</p>
              </div>
            )}

            {uploadStatus === 'success' && (
              <Alert className="mt-4 border-green-200 bg-green-50 text-green-800" role="status" aria-live="polite">
                <CheckCircle2 className="h-4 w-4" aria-hidden="true" />
                <AlertDescription>アップロードが完了しました。</AlertDescription>
              </Alert>
            )}

            {uploadStatus === 'error' && errorMessage && (
              <Alert variant="destructive" className="mt-4" role="alert" aria-live="assertive">
                <AlertCircle className="h-4 w-4" aria-hidden="true" />
                <AlertDescription>{errorMessage}</AlertDescription>
              </Alert>
            )}
          </CardContent>
        </Card>
      </div>
    </main>
  )
}
