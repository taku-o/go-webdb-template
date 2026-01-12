'use client'

import { useState, useEffect } from 'react'
import Link from 'next/link'
import { apiClient } from '@/lib/api'
import { DmPost } from '@/types/dm_post'
import { DmUser } from '@/types/dm_user'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { ErrorAlert } from '@/components/shared/error-alert'
import { LoadingSpinner } from '@/components/shared/loading-spinner'
import { LoadingOverlay } from '@/components/shared/loading-overlay'
import { ArrowLeft, Trash2 } from 'lucide-react'

export default function PostsPage() {
  const [dmPosts, setPosts] = useState<DmPost[]>([])
  const [dmUsers, setDmUsers] = useState<DmUser[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [userId, setUserId] = useState('')
  const [title, setTitle] = useState('')
  const [content, setContent] = useState('')
  const [creating, setCreating] = useState(false)

  const loadPosts = async () => {
    try {
      setLoading(true)
      const data = await apiClient.getDmPosts()
      setPosts(data)
      setError(null)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load posts')
    } finally {
      setLoading(false)
    }
  }

  const loadUsers = async () => {
    try {
      const data = await apiClient.getDmUsers()
      setDmUsers(data)
    } catch (err) {
      console.error('Failed to load users:', err)
    }
  }

  useEffect(() => {
    loadPosts()
    loadUsers()
  }, [])

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!userId || !title || !content) return

    try {
      setCreating(true)
      setError(null)
      await apiClient.createDmPost({
        user_id: userId,
        title,
        content,
      })
      setUserId('')
      setTitle('')
      setContent('')
      await loadPosts()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to create post')
    } finally {
      setCreating(false)
    }
  }

  const handleDelete = async (dmPost: DmPost) => {
    if (!confirm('本当に削除しますか？')) return

    try {
      setError(null)
      await apiClient.deleteDmPost(dmPost.id, dmPost.user_id)
      await loadPosts()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to delete post')
    }
  }

  return (
    <main className="min-h-screen p-4 sm:p-6 md:p-8">
      <div className="max-w-6xl mx-auto">
        <nav aria-label="パンくずリスト">
          <div className="mb-4 sm:mb-6">
            <Link href="/" className="inline-flex items-center text-primary hover:underline text-sm sm:text-base" aria-label="トップページに戻る">
              <ArrowLeft className="mr-2 h-4 w-4" aria-hidden="true" />
              トップページに戻る
            </Link>
          </div>
        </nav>

        <h1 className="text-2xl sm:text-3xl font-bold mb-6 sm:mb-8">投稿管理</h1>

        {error && (
          <div className="mb-4" role="alert" aria-live="assertive">
            <ErrorAlert message={error} />
          </div>
        )}

        {/* 作成フォーム */}
        <Card className="mb-6 sm:mb-8">
          <CardHeader>
            <CardTitle className="text-lg sm:text-xl">新規投稿作成</CardTitle>
            <CardDescription className="text-sm">
              新しい投稿を作成します
            </CardDescription>
          </CardHeader>
          <CardContent>
            <form onSubmit={handleCreate} className="space-y-4" aria-label="新規投稿作成フォーム">
              <div className="space-y-2">
                <Label htmlFor="userId">ユーザー</Label>
                <Select value={userId} onValueChange={setUserId} required>
                  <SelectTrigger id="userId">
                    <SelectValue placeholder="選択してください" />
                  </SelectTrigger>
                  <SelectContent>
                    {dmUsers.map((dmUser) => (
                      <SelectItem key={dmUser.id} value={dmUser.id}>
                        {dmUser.name} ({dmUser.email})
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
                {dmUsers.length === 0 && (
                  <p className="text-sm text-muted-foreground">
                    先に<Link href="/dm-users" className="text-primary hover:underline">ユーザー</Link>を作成してください
                  </p>
                )}
              </div>
              <div className="space-y-2">
                <Label htmlFor="title">タイトル</Label>
                <Input
                  id="title"
                  type="text"
                  value={title}
                  onChange={(e) => setTitle(e.target.value)}
                  placeholder="投稿のタイトルを入力"
                  required
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="content">本文</Label>
                <Textarea
                  id="content"
                  value={content}
                  onChange={(e) => setContent(e.target.value)}
                  placeholder="投稿の本文を入力"
                  rows={6}
                  required
                />
              </div>
              <Button
                type="submit"
                disabled={creating || dmUsers.length === 0}
                className="w-full"
                aria-label={creating ? "投稿作成中" : "投稿を作成"}
                aria-busy={creating}
                aria-disabled={dmUsers.length === 0}
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

        {/* 投稿一覧 */}
        <Card>
          <CardHeader>
            <CardTitle className="text-lg sm:text-xl">投稿一覧</CardTitle>
            <CardDescription className="text-sm">
              {dmPosts.length}件の投稿が登録されています
            </CardDescription>
          </CardHeader>
          <CardContent>
            {loading ? (
              <div role="status" aria-live="polite" aria-label="投稿一覧を読み込み中">
                <LoadingOverlay message="読み込み中..." />
              </div>
            ) : dmPosts.length === 0 ? (
              <p className="text-center text-muted-foreground py-8 text-sm sm:text-base" role="status">
                投稿がありません。上のフォームから作成してください。
              </p>
            ) : (
              <div className="rounded-md border overflow-x-auto">
                <Table role="table" aria-label="投稿一覧">
                  <TableHeader>
                    <TableRow>
                      <TableHead>タイトル</TableHead>
                      <TableHead className="hidden md:table-cell">本文</TableHead>
                      <TableHead className="hidden sm:table-cell">ユーザーID</TableHead>
                      <TableHead className="hidden lg:table-cell">作成日時</TableHead>
                      <TableHead className="hidden lg:table-cell">更新日時</TableHead>
                      <TableHead className="text-right">操作</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {dmPosts.map((dmPost) => (
                      <TableRow key={dmPost.id}>
                        <TableCell className="font-medium">
                          <div>
                            <div>{dmPost.title}</div>
                            <div className="text-xs text-muted-foreground md:hidden mt-1 line-clamp-2">{dmPost.content}</div>
                          </div>
                        </TableCell>
                        <TableCell className="hidden md:table-cell max-w-md">
                          <p className="truncate">{dmPost.content}</p>
                        </TableCell>
                        <TableCell className="hidden sm:table-cell font-mono text-xs">{dmPost.user_id}</TableCell>
                        <TableCell className="hidden lg:table-cell text-sm text-muted-foreground">
                          {new Date(dmPost.created_at).toLocaleString('ja-JP')}
                        </TableCell>
                        <TableCell className="hidden lg:table-cell text-sm text-muted-foreground">
                          {new Date(dmPost.updated_at).toLocaleString('ja-JP')}
                        </TableCell>
                        <TableCell className="text-right">
                          <Button
                            onClick={() => handleDelete(dmPost)}
                            variant="destructive"
                            size="sm"
                            aria-label={`${dmPost.title}を削除`}
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
