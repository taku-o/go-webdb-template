'use client'

import { useState } from 'react'
import Link from 'next/link'
import { apiClient } from '@/lib/api'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { LoadingSpinner } from '@/components/shared/loading-spinner'
import { ArrowLeft, CheckCircle2, AlertCircle, Clock } from 'lucide-react'

export default function DmJobqueuePage() {
  const [message, setMessage] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [success, setSuccess] = useState<string | null>(null)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError(null)
    setSuccess(null)

    try {
      setLoading(true)
      const response = await apiClient.registerJob({
        message: message || undefined,
      })
      setSuccess(`ジョブが登録されました (ID: ${response.job_id})`)
      setMessage('')
    } catch (err) {
      setError(err instanceof Error ? err.message : 'ジョブの登録に失敗しました')
    } finally {
      setLoading(false)
    }
  }

  return (
    <main className="min-h-screen p-8">
      <div className="max-w-4xl mx-auto">
        <div className="mb-6">
          <Link href="/" className="inline-flex items-center text-blue-600 hover:text-blue-800 hover:underline">
            <ArrowLeft className="mr-2 h-4 w-4" />
            トップページに戻る
          </Link>
        </div>

        <h1 className="text-3xl font-bold mb-4 flex items-center gap-2">
          <Clock className="h-8 w-8" />
          ジョブキュー (参考コード)
        </h1>

        <p className="mb-6 text-gray-600">
          このページは参考コードです。ボタンをクリックすると、3分後に標準出力にメッセージが出力されるジョブが登録されます。
        </p>

        {error && (
          <Alert variant="destructive" className="mb-4">
            <AlertCircle className="h-4 w-4" />
            <AlertDescription>{error}</AlertDescription>
          </Alert>
        )}

        {success && (
          <Alert className="mb-4 border-green-200 bg-green-50 text-green-800">
            <CheckCircle2 className="h-4 w-4" />
            <AlertDescription>{success}</AlertDescription>
          </Alert>
        )}

        <Card>
          <CardHeader>
            <CardTitle>ジョブ登録</CardTitle>
            <CardDescription>
              3分後に実行されるジョブを登録します
            </CardDescription>
          </CardHeader>
          <CardContent>
            <form onSubmit={handleSubmit} className="space-y-4">
              <div className="space-y-2">
                <Label htmlFor="message">メッセージ (オプション)</Label>
                <Input
                  id="message"
                  type="text"
                  value={message}
                  onChange={(e) => setMessage(e.target.value)}
                  placeholder="出力するメッセージを入力"
                />
              </div>
              <Button
                type="submit"
                disabled={loading}
                className="w-full"
              >
                {loading ? (
                  <>
                    <LoadingSpinner size="sm" className="mr-2" />
                    登録中...
                  </>
                ) : (
                  <>
                    <Clock className="mr-2 h-4 w-4" />
                    ジョブを登録
                  </>
                )}
              </Button>
            </form>
          </CardContent>
        </Card>
      </div>
    </main>
  )
}
