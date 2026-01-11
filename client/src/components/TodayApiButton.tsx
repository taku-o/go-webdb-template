'use client'

import { useUser } from '@auth0/nextjs-auth0'
import { useState } from 'react'
import { apiClient } from '@/lib/api'

export default function TodayApiButton() {
  const { user: auth0user, isLoading } = useUser()
  const [date, setDate] = useState<string | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)

  const handleClick = async () => {
    setLoading(true)
    setError(null)
    setDate(null)

    try {
      const data = await apiClient.getToday(auth0user || undefined)
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
