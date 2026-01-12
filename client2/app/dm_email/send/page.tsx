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
import { ArrowLeft, CheckCircle2, AlertCircle, Mail } from 'lucide-react'

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
          <Link href="/" className="inline-flex items-center text-blue-600 hover:text-blue-800 hover:underline">
            <ArrowLeft className="mr-2 h-4 w-4" />
            トップページに戻る
          </Link>
        </div>

        <h1 className="text-3xl font-bold mb-8 flex items-center gap-2">
          <Mail className="h-8 w-8" />
          メール送信
        </h1>

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
            <CardTitle>ウェルカムメール送信</CardTitle>
            <CardDescription>
              ユーザーにウェルカムメールを送信します
            </CardDescription>
          </CardHeader>
          <CardContent>
            <form onSubmit={handleSubmit} className="space-y-4">
              <div className="space-y-2">
                <Label htmlFor="toEmail">送信先メールアドレス</Label>
                <Input
                  id="toEmail"
                  type="email"
                  value={toEmail}
                  onChange={(e) => setToEmail(e.target.value)}
                  placeholder="example@example.com"
                  required
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="name">お名前</Label>
                <Input
                  id="name"
                  type="text"
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  placeholder="山田 太郎"
                  required
                />
              </div>
              <Button
                type="submit"
                disabled={sending}
                className="w-full"
              >
                {sending ? (
                  <>
                    <LoadingSpinner size="sm" className="mr-2" />
                    送信中...
                  </>
                ) : (
                  <>
                    <Mail className="mr-2 h-4 w-4" />
                    メールを送信
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
