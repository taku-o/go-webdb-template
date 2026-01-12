'use client'

import { useState } from 'react'
import { apiClient } from '@/lib/api'
import { Button } from '@/components/ui/button'
import { ErrorAlert } from '@/components/shared/error-alert'
import { LoadingSpinner } from '@/components/shared/loading-spinner'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'

export default function TodayApiButton() {
  const [date, setDate] = useState<string | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)

  const handleClick = async () => {
    setLoading(true)
    setError(null)
    setDate(null)

    try {
      const data = await apiClient.getToday()
      setDate(data.date)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'An error occurred')
    } finally {
      setLoading(false)
    }
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle>Today API (Private Endpoint)</CardTitle>
        <CardDescription>
          NextAuthログイン時のみアクセス可能なプライベートAPIをテストします。
        </CardDescription>
      </CardHeader>
      <CardContent>
        <Button
          onClick={handleClick}
          disabled={loading}
          className="w-full"
        >
          {loading ? (
            <>
              <LoadingSpinner size="sm" className="mr-2" />
              Loading...
            </>
          ) : (
            'Get Today'
          )}
        </Button>
        {date && (
          <p className="mt-4 text-green-600 font-semibold">
            Today: {date}
          </p>
        )}
        {error && (
          <div className="mt-4">
            <ErrorAlert message={error} />
          </div>
        )}
      </CardContent>
    </Card>
  )
}
