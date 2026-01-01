'use client'

import { useState } from 'react'
import Link from 'next/link'
import { apiClient } from '@/lib/api'

export default function SendEmailPage() {
  const [toEmail, setToEmail] = useState('')
  const [name, setName] = useState('')
  const [sending, setSending] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [success, setSuccess] = useState<string | null>(null)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!toEmail || !name) return

    try {
      setSending(true)
      setError(null)
      setSuccess(null)

      const result = await apiClient.sendEmail(
        [toEmail],
        'welcome',
        { Name: name, Email: toEmail }
      )

      if (result.success) {
        setSuccess(result.message)
        setToEmail('')
        setName('')
      } else {
        setError(result.message || 'メール送信に失敗しました')
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : 'メール送信に失敗しました')
    } finally {
      setSending(false)
    }
  }

  return (
    <main className="min-h-screen p-8">
      <div className="max-w-2xl mx-auto">
        <div className="mb-6">
          <Link href="/" className="text-blue-500 hover:underline">
            &larr; トップページに戻る
          </Link>
        </div>

        <h1 className="text-3xl font-bold mb-8">メール送信</h1>

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
          <h2 className="text-xl font-semibold mb-4">ウェルカムメール送信</h2>
          <p className="text-sm text-gray-600 mb-4">
            ユーザーにウェルカムメールを送信します。
          </p>
          <form onSubmit={handleSubmit} className="space-y-4">
            <div>
              <label htmlFor="toEmail" className="block text-sm font-medium mb-1">送信先メールアドレス</label>
              <input
                id="toEmail"
                type="email"
                value={toEmail}
                onChange={(e) => setToEmail(e.target.value)}
                className="w-full px-3 py-2 border rounded"
                placeholder="example@example.com"
                required
              />
            </div>
            <div>
              <label htmlFor="name" className="block text-sm font-medium mb-1">お名前</label>
              <input
                id="name"
                type="text"
                value={name}
                onChange={(e) => setName(e.target.value)}
                className="w-full px-3 py-2 border rounded"
                placeholder="山田 太郎"
                required
              />
            </div>
            <button
              type="submit"
              disabled={sending}
              className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600 disabled:bg-gray-400"
            >
              {sending ? '送信中...' : 'メールを送信'}
            </button>
          </form>
        </div>
      </div>
    </main>
  )
}
