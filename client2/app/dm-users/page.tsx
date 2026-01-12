'use client'

import { useState, useEffect } from 'react'
import Link from 'next/link'
import { apiClient } from '@/lib/api'
import { DmUser } from '@/types/dm_user'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { ErrorAlert } from '@/components/shared/error-alert'
import { LoadingSpinner } from '@/components/shared/loading-spinner'
import { LoadingOverlay } from '@/components/shared/loading-overlay'
import { ArrowLeft, Download, Trash2 } from 'lucide-react'

export default function UsersPage() {
  const [dmUsers, setDmUsers] = useState<DmUser[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [name, setName] = useState('')
  const [email, setEmail] = useState('')
  const [creating, setCreating] = useState(false)
  const [downloading, setDownloading] = useState(false)
  const [downloadError, setDownloadError] = useState<string | null>(null)

  const loadUsers = async () => {
    try {
      setLoading(true)
      const data = await apiClient.getDmUsers()
      setDmUsers(data)
      setError(null)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load users')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    loadUsers()
  }, [])

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!name || !email) return

    try {
      setCreating(true)
      setError(null)
      await apiClient.createDmUser({ name, email })
      setName('')
      setEmail('')
      await loadUsers()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to create user')
    } finally {
      setCreating(false)
    }
  }

  const handleDelete = async (id: string) => {
    if (!confirm('本当に削除しますか？')) return

    try {
      setError(null)
      await apiClient.deleteDmUser(id)
      await loadUsers()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to delete user')
    }
  }

  const handleDownloadCSV = async () => {
    try {
      setDownloading(true)
      setDownloadError(null)
      await apiClient.downloadDmUsersCSV()
    } catch (err) {
      setDownloadError(err instanceof Error ? err.message : 'Failed to download CSV')
    } finally {
      setDownloading(false)
    }
  }

  return (
    <main className="min-h-screen p-8">
      <div className="max-w-6xl mx-auto">
        <div className="mb-6">
          <Link href="/" className="inline-flex items-center text-blue-600 hover:text-blue-800 hover:underline">
            <ArrowLeft className="mr-2 h-4 w-4" />
            トップページに戻る
          </Link>
        </div>

        <h1 className="text-3xl font-bold mb-8">ユーザー管理</h1>

        {error && (
          <div className="mb-4">
            <ErrorAlert message={error} />
          </div>
        )}

        {/* 作成フォーム */}
        <Card className="mb-8">
          <CardHeader>
            <CardTitle>新規ユーザー作成</CardTitle>
            <CardDescription>
              新しいユーザーを追加します
            </CardDescription>
          </CardHeader>
          <CardContent>
            <form onSubmit={handleCreate} className="space-y-4">
              <div className="space-y-2">
                <Label htmlFor="name">名前</Label>
                <Input
                  id="name"
                  type="text"
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  placeholder="ユーザー名を入力"
                  required
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="email">メールアドレス</Label>
                <Input
                  id="email"
                  type="email"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  placeholder="example@example.com"
                  required
                />
              </div>
              <Button
                type="submit"
                disabled={creating}
                className="w-full"
              >
                {creating ? (
                  <>
                    <LoadingSpinner size="sm" className="mr-2" />
                    作成中...
                  </>
                ) : (
                  '作成'
                )}
              </Button>
            </form>
          </CardContent>
        </Card>

        {/* ユーザー一覧 */}
        <Card>
          <CardHeader>
            <div className="flex justify-between items-center">
              <div>
                <CardTitle>ユーザー一覧</CardTitle>
                <CardDescription>
                  {dmUsers.length}件のユーザーが登録されています
                </CardDescription>
              </div>
              <Button
                onClick={handleDownloadCSV}
                disabled={downloading}
                variant="outline"
              >
                {downloading ? (
                  <>
                    <LoadingSpinner size="sm" className="mr-2" />
                    ダウンロード中...
                  </>
                ) : (
                  <>
                    <Download className="mr-2 h-4 w-4" />
                    CSVダウンロード
                  </>
                )}
              </Button>
            </div>
          </CardHeader>
          <CardContent>
            {downloadError && (
              <div className="mb-4">
                <ErrorAlert message={downloadError} />
              </div>
            )}
            {loading ? (
              <LoadingOverlay message="読み込み中..." />
            ) : dmUsers.length === 0 ? (
              <p className="text-center text-gray-500 py-8">
                ユーザーがいません。上のフォームから作成してください。
              </p>
            ) : (
              <div className="rounded-md border">
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead>ID</TableHead>
                      <TableHead>名前</TableHead>
                      <TableHead>メールアドレス</TableHead>
                      <TableHead>作成日時</TableHead>
                      <TableHead>更新日時</TableHead>
                      <TableHead className="text-right">操作</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {dmUsers.map((dmUser) => (
                      <TableRow key={dmUser.id}>
                        <TableCell className="font-mono text-xs">{dmUser.id}</TableCell>
                        <TableCell className="font-medium">{dmUser.name}</TableCell>
                        <TableCell>{dmUser.email}</TableCell>
                        <TableCell className="text-sm text-gray-600">
                          {new Date(dmUser.created_at).toLocaleString('ja-JP')}
                        </TableCell>
                        <TableCell className="text-sm text-gray-600">
                          {new Date(dmUser.updated_at).toLocaleString('ja-JP')}
                        </TableCell>
                        <TableCell className="text-right">
                          <Button
                            onClick={() => handleDelete(dmUser.id)}
                            variant="destructive"
                            size="sm"
                          >
                            <Trash2 className="mr-2 h-4 w-4" />
                            削除
                          </Button>
                        </TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </div>
            )}
          </CardContent>
        </Card>
      </div>
    </main>
  )
}
