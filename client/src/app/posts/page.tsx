'use client'

import { useState, useEffect } from 'react'
import Link from 'next/link'
import { apiClient } from '@/lib/api'
import { Post, User } from '@/types/user'

export default function PostsPage() {
  const [posts, setPosts] = useState<Post[]>([])
  const [users, setUsers] = useState<User[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [userId, setUserId] = useState('')
  const [title, setTitle] = useState('')
  const [content, setContent] = useState('')
  const [creating, setCreating] = useState(false)

  const loadPosts = async () => {
    try {
      setLoading(true)
      const data = await apiClient.getPosts()
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
      const data = await apiClient.getUsers()
      setUsers(data)
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
      await apiClient.createPost({
        user_id: parseInt(userId),
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

  const handleDelete = async (post: Post) => {
    if (!confirm('本当に削除しますか？')) return

    try {
      await apiClient.deletePost(post.id, post.user_id)
      await loadPosts()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to delete post')
    }
  }

  return (
    <main className="min-h-screen p-8">
      <div className="max-w-4xl mx-auto">
        <div className="mb-6">
          <Link href="/" className="text-blue-500 hover:underline">
            ← トップページに戻る
          </Link>
        </div>

        <h1 className="text-3xl font-bold mb-8">投稿管理</h1>

        {error && (
          <div className="mb-4 p-4 bg-red-100 border border-red-400 text-red-700 rounded">
            {error}
          </div>
        )}

        {/* 作成フォーム */}
        <div className="mb-8 p-6 border rounded-lg bg-gray-50">
          <h2 className="text-xl font-semibold mb-4">新規投稿作成</h2>
          <form onSubmit={handleCreate} className="space-y-4">
            <div>
              <label className="block text-sm font-medium mb-1">ユーザー</label>
              <select
                value={userId}
                onChange={(e) => setUserId(e.target.value)}
                className="w-full px-3 py-2 border rounded"
                required
              >
                <option value="">選択してください</option>
                {users.map((user) => (
                  <option key={user.id} value={user.id}>
                    {user.name} ({user.email})
                  </option>
                ))}
              </select>
              {users.length === 0 && (
                <p className="text-sm text-gray-500 mt-1">
                  先に<Link href="/users" className="text-blue-500 hover:underline">ユーザー</Link>を作成してください
                </p>
              )}
            </div>
            <div>
              <label className="block text-sm font-medium mb-1">タイトル</label>
              <input
                type="text"
                value={title}
                onChange={(e) => setTitle(e.target.value)}
                className="w-full px-3 py-2 border rounded"
                required
              />
            </div>
            <div>
              <label className="block text-sm font-medium mb-1">本文</label>
              <textarea
                value={content}
                onChange={(e) => setContent(e.target.value)}
                className="w-full px-3 py-2 border rounded"
                rows={4}
                required
              />
            </div>
            <button
              type="submit"
              disabled={creating || users.length === 0}
              className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600 disabled:bg-gray-400"
            >
              {creating ? '作成中...' : '作成'}
            </button>
          </form>
        </div>

        {/* 投稿一覧 */}
        <div>
          <h2 className="text-xl font-semibold mb-4">投稿一覧</h2>
          {loading ? (
            <p>読み込み中...</p>
          ) : posts.length === 0 ? (
            <p className="text-gray-500">投稿がありません。上のフォームから作成してください。</p>
          ) : (
            <div className="space-y-4">
              {posts.map((post) => (
                <div key={post.id} className="p-4 border rounded-lg">
                  <div className="flex justify-between items-start mb-2">
                    <h3 className="font-bold text-lg">{post.title}</h3>
                    <button
                      onClick={() => handleDelete(post)}
                      className="px-3 py-1 bg-red-500 text-white rounded hover:bg-red-600 text-sm"
                    >
                      削除
                    </button>
                  </div>
                  <p className="text-gray-700 mb-2">{post.content}</p>
                  <div className="text-xs text-gray-400">
                    <div>投稿ID: {post.id} | ユーザーID: {post.user_id}</div>
                    <div>作成日: {new Date(post.created_at).toLocaleString('ja-JP')}</div>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>
    </main>
  )
}
