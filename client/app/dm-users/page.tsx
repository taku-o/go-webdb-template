'use client'

import { useState, useRef } from 'react'
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
  const hasLoadedRef = useRef(false)

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

  const loadInitialData = async () => {
    if (hasLoadedRef.current) return
    hasLoadedRef.current = true
    await loadUsers()
  }

  const setContainerRef = (node: HTMLElement | null) => {
    if (node && !hasLoadedRef.current) {
      loadInitialData()
    }
  }

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
    <main ref={setContainerRef} className="min-h-screen p-4 sm:p-6 md:p-8">
      <div className="max-w-6xl mx-auto">
        <nav aria-label="パンくずリスト">
          <div className="mb-4 sm:mb-6">
            <Link href="/" className="inline-flex items-center text-primary hover:underline text-sm sm:text-base" aria-label="トップページに戻る">
              <ArrowLeft className="mr-2 h-4 w-4" aria-hidden="true" />
              トップページに戻る
            </Link>
          </div>
        </nav>

        <h1 className="text-2xl sm:text-3xl font-bold mb-6 sm:mb-8">ユーザー管理</h1>

        {error && (
          <div className="mb-4" role="alert" aria-live="assertive">
            <ErrorAlert message={error} />
          </div>
        )}

        {/* 作成フォーム */}
        <Card className="mb-6 sm:mb-8">
          <CardHeader>
            <CardTitle className="text-lg sm:text-xl">新規ユーザー作成</CardTitle>
            <CardDescription className="text-sm">
              新しいユーザーを追加します
            </CardDescription>
          </CardHeader>
          <CardContent>
            <form onSubmit={handleCreate} className="space-y-4" aria-label="新規ユーザー作成フォーム">
              <div className="space-y-2">
                <Label htmlFor="name">名前</Label>
                <Input
                  id="name"
                  type="text"
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  placeholder="ユーザー名を入力"
                  required
                  aria-required="true"
                  aria-invalid={error ? "true" : "false"}
                  aria-describedby={error ? "name-error" : undefined}
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
                  aria-required="true"
                  aria-invalid={error ? "true" : "false"}
                  aria-describedby={error ? "email-error" : undefined}
                />
              </div>
              <Button
                type="submit"
                disabled={creating}
                className="w-full"
                aria-label={creating ? "ユーザー作成中" : "ユーザーを作成"}
                aria-busy={creating}
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
            <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
              <div>
                <CardTitle className="text-lg sm:text-xl">ユーザー一覧</CardTitle>
                <CardDescription className="text-sm">
                  {dmUsers.length}件のユーザーが登録されています
                </CardDescription>
              </div>
              <Button
                onClick={handleDownloadCSV}
                disabled={downloading}
                variant="outline"
                size="sm"
                className="w-full sm:w-auto"
                aria-label={downloading ? "CSVダウンロード中" : "ユーザー一覧をCSV形式でダウンロード"}
                aria-busy={downloading}
              >
                {downloading ? (
                  <>
                    <LoadingSpinner size="sm" className="mr-2" />
                    ダウンロード中...
                  </>
                ) : (
                  <>
                    <Download className="mr-2 h-4 w-4" aria-hidden="true" />
                    CSVダウンロード
                  </>
                )}
              </Button>
            </div>
          </CardHeader>
          <CardContent>
            {downloadError && (
              <div className="mb-4" role="alert" aria-live="assertive">
                <ErrorAlert message={downloadError} />
              </div>
            )}
            {loading ? (
              <div role="status" aria-live="polite" aria-label="ユーザー一覧を読み込み中">
                <LoadingOverlay message="読み込み中..." />
              </div>
            ) : dmUsers.length === 0 ? (
              <p className="text-center text-muted-foreground py-8 text-sm sm:text-base" role="status">
                ユーザーがいません。上のフォームから作成してください。
              </p>
            ) : (
              <div className="rounded-md border overflow-x-auto">
                <Table role="table" aria-label="ユーザー一覧">
                  <TableHeader>
                    <TableRow>
                      <TableHead className="hidden sm:table-cell">ID</TableHead>
                      <TableHead>名前</TableHead>
                      <TableHead className="hidden md:table-cell">メールアドレス</TableHead>
                      <TableHead className="hidden lg:table-cell">作成日時</TableHead>
                      <TableHead className="hidden lg:table-cell">更新日時</TableHead>
                      <TableHead className="text-right">操作</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {dmUsers.map((dmUser) => (
                      <TableRow key={dmUser.id}>
                        <TableCell className="hidden sm:table-cell font-mono text-xs">{dmUser.id}</TableCell>
                        <TableCell className="font-medium">
                          <div>
                            <div>{dmUser.name}</div>
                            <div className="text-xs text-muted-foreground sm:hidden">{dmUser.email}</div>
                          </div>
                        </TableCell>
                        <TableCell className="hidden md:table-cell">{dmUser.email}</TableCell>
                        <TableCell className="hidden lg:table-cell text-sm text-muted-foreground">
                          {new Date(dmUser.created_at).toLocaleString('ja-JP')}
                        </TableCell>
                        <TableCell className="hidden lg:table-cell text-sm text-muted-foreground">
                          {new Date(dmUser.updated_at).toLocaleString('ja-JP')}
                        </TableCell>
                        <TableCell className="text-right">
                          <Button
                            onClick={() => handleDelete(dmUser.id)}
                            variant="destructive"
                            size="sm"
                            aria-label={`${dmUser.name}を削除`}
                          >
                            <Trash2 className="mr-2 h-4 w-4" aria-hidden="true" />
                            <span className="hidden sm:inline">削除</span>
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
