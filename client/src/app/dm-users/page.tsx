'use client'

import { useState, useEffect } from 'react'
import Link from 'next/link'
import { apiClient } from '@/lib/api'
import { DmUser } from '@/types/dm_user'

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
      <div className="max-w-4xl mx-auto">
        <div className="mb-6">
          <Link href="/" className="text-blue-500 hover:underline">
            ← トップページに戻る
          </Link>
        </div>

        <h1 className="text-3xl font-bold mb-8">ユーザー管理</h1>

        {error && (
          <div className="mb-4 p-4 bg-red-100 border border-red-400 text-red-700 rounded">
            {error}
          </div>
        )}

        {/* 作成フォーム */}
        <div className="mb-8 p-6 border rounded-lg bg-gray-50">
          <h2 className="text-xl font-semibold mb-4">新規ユーザー作成</h2>
          <form onSubmit={handleCreate} className="space-y-4">
            <div>
              <label className="block text-sm font-medium mb-1">名前</label>
              <input
                type="text"
                value={name}
                onChange={(e) => setName(e.target.value)}
                className="w-full px-3 py-2 border rounded"
                required
              />
            </div>
            <div>
              <label className="block text-sm font-medium mb-1">メールアドレス</label>
              <input
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                className="w-full px-3 py-2 border rounded"
                required
              />
            </div>
            <button
              type="submit"
              disabled={creating}
              className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600 disabled:bg-gray-400"
            >
              {creating ? '作成中...' : '作成'}
            </button>
          </form>
        </div>

        {/* ユーザー一覧 */}
        <div>
          <div className="flex justify-between items-center mb-4">
            <h2 className="text-xl font-semibold">ユーザー一覧</h2>
            <button
              onClick={handleDownloadCSV}
              disabled={downloading}
              className="px-4 py-2 bg-green-500 text-white rounded hover:bg-green-600 disabled:opacity-50"
            >
              {downloading ? 'ダウンロード中...' : 'CSVダウンロード'}
            </button>
          </div>
          {downloadError && (
            <div className="mb-4 p-4 bg-red-100 border border-red-400 text-red-700 rounded">
              {downloadError}
            </div>
          )}
          {loading ? (
            <p>読み込み中...</p>
          ) : dmUsers.length === 0 ? (
            <p className="text-gray-500">ユーザーがいません。上のフォームから作成してください。</p>
          ) : (
            <div className="space-y-2">
              {dmUsers.map((dmUser, index) => (
                <div key={index} className="p-4 border rounded-lg flex justify-between items-center">
                  <div>
                    <div className="font-medium">{dmUser.name}</div>
                    <div className="text-sm text-gray-600">{dmUser.email}</div>
                    <div className="text-xs text-gray-400">ID: {dmUser.id}</div>
                  </div>
                  <button
                    onClick={() => handleDelete(dmUser.id)}
                    className="px-3 py-1 bg-red-500 text-white rounded hover:bg-red-600"
                  >
                    削除
                  </button>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>
    </main>
  )
}
