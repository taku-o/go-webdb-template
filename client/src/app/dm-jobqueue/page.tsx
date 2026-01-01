'use client'

import { useState } from 'react'
import Link from 'next/link'
import { apiClient } from '@/lib/api'

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
          <Link href="/" className="text-blue-500 hover:underline">
            &larr; トップページに戻る
          </Link>
        </div>

        <h1 className="text-3xl font-bold mb-8">ジョブキュー (参考コード)</h1>

        <p className="mb-6 text-gray-600">
          このページは参考コードです。ボタンをクリックすると、3分後に標準出力にメッセージが出力されるジョブが登録されます。
        </p>

        {error && (
          <div className="mb-4 p-4 bg-red-100 border border-red-400 text-red-700 rounded">
            {error}
          </div>
        )}

        {success && (
          <div className="mb-4 p-4 bg-green-100 border border-green-400 text-green-700 rounded">
            {success}
          </div>
        )}

        <div className="p-6 border rounded-lg bg-gray-50">
          <h2 className="text-xl font-semibold mb-4">ジョブ登録</h2>
          <form onSubmit={handleSubmit} className="space-y-4">
            <div>
              <label htmlFor="message" className="block text-sm font-medium mb-1">
                メッセージ (オプション)
              </label>
              <input
                id="message"
                type="text"
                value={message}
                onChange={(e) => setMessage(e.target.value)}
                className="w-full px-3 py-2 border rounded"
                placeholder="出力するメッセージを入力"
              />
            </div>
            <button
              type="submit"
              disabled={loading}
              className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600 disabled:bg-gray-400"
            >
              {loading ? '登録中...' : 'ジョブを登録'}
            </button>
          </form>
        </div>
      </div>
    </main>
  )
}
