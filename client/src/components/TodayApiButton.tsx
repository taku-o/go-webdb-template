'use client'

import { useUser } from '@auth0/nextjs-auth0'
import { useState } from 'react'

export default function TodayApiButton() {
  const { user, isLoading } = useUser()
  const [date, setDate] = useState<string | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)

  const handleClick = async () => {
    setLoading(true)
    setError(null)
    setDate(null)

    try {
      // JWTの取得
      let token: string
      if (user) {
        // ログイン中: Auth0 JWTを取得
        const response = await fetch('/auth/token')
        if (!response.ok) {
          throw new Error('Failed to get access token')
        }
        const data = await response.json()
        token = data.accessToken
      } else {
        // 未ログイン: Public API Keyを使用
        token = process.env.NEXT_PUBLIC_API_KEY!
      }

      // API呼び出し
      const apiBaseUrl = process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8080'
      const response = await fetch(`${apiBaseUrl}/api/today`, {
        headers: {
          'Authorization': `Bearer ${token}`,
        },
      })

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}))
        throw new Error(errorData.message || response.statusText)
      }

      const data = await response.json()
      setDate(data.date)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'An error occurred')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="p-4 border rounded-lg bg-gray-50">
      <h3 className="font-semibold mb-4">Today API (Private Endpoint)</h3>
      <p className="text-sm text-gray-600 mb-4">
        Auth0ログイン時のみアクセス可能なプライベートAPIをテストします。
      </p>
      <button
        onClick={handleClick}
        disabled={loading || isLoading}
        className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
      >
        {loading ? 'Loading...' : 'Get Today'}
      </button>
      {date && (
        <p className="mt-4 text-green-600">
          Today: {date}
        </p>
      )}
      {error && (
        <p className="mt-4 text-red-600">
          Error: {error}
        </p>
      )}
    </div>
  )
}
