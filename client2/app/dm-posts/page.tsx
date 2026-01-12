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
    <main className="min-h-screen p-8">
      <div className="max-w-6xl mx-auto">
        <div className="mb-6">
          <Link href="/" className="inline-flex items-center text-blue-600 hover:text-blue-800 hover:underline">
            <ArrowLeft className="mr-2 h-4 w-4" />
            トップページに戻る
          </Link>
        </div>

        <h1 className="text-3xl font-bold mb-8">投稿管理</h1>

        {error && (
          <div className="mb-4">
            <ErrorAlert message={error} />
          </div>
        )}

        {/* 作成フォーム */}
        <Card className="mb-8">
          <CardHeader>
            <CardTitle>新規投稿作成</CardTitle>
            <CardDescription>
              新しい投稿を作成します
            </CardDescription>
          </CardHeader>
          <CardContent>
            <form onSubmit={handleCreate} className="space-y-4">
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
                  <p className="text-sm text-gray-500">
                    先に<Link href="/dm-users" className="text-blue-600 hover:underline">ユーザー</Link>を作成してください
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
            <CardTitle>投稿一覧</CardTitle>
            <CardDescription>
              {dmPosts.length}件の投稿が登録されています
            </CardDescription>
          </CardHeader>
          <CardContent>
            {loading ? (
              <LoadingOverlay message="読み込み中..." />
            ) : dmPosts.length === 0 ? (
              <p className="text-center text-gray-500 py-8">
                投稿がありません。上のフォームから作成してください。
              </p>
            ) : (
              <div className="rounded-md border">
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead>タイトル</TableHead>
                      <TableHead>本文</TableHead>
                      <TableHead>ユーザーID</TableHead>
                      <TableHead>作成日時</TableHead>
                      <TableHead>更新日時</TableHead>
                      <TableHead className="text-right">操作</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {dmPosts.map((dmPost) => (
                      <TableRow key={dmPost.id}>
                        <TableCell className="font-medium">{dmPost.title}</TableCell>
                        <TableCell className="max-w-md">
                          <p className="truncate">{dmPost.content}</p>
                        </TableCell>
                        <TableCell className="font-mono text-xs">{dmPost.user_id}</TableCell>
                        <TableCell className="text-sm text-gray-600">
                          {new Date(dmPost.created_at).toLocaleString('ja-JP')}
                        </TableCell>
                        <TableCell className="text-sm text-gray-600">
                          {new Date(dmPost.updated_at).toLocaleString('ja-JP')}
                        </TableCell>
                        <TableCell className="text-right">
                          <Button
                            onClick={() => handleDelete(dmPost)}
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
