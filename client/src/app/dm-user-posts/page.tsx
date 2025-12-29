'use client'

import { useState, useEffect } from 'react'
import Link from 'next/link'
import { apiClient } from '@/lib/api'
import { UserPost } from '@/types/post'

export default function UserPostsPage() {
  const [userPosts, setUserPosts] = useState<UserPost[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const loadUserPosts = async () => {
    try {
      setLoading(true)
      const data = await apiClient.getUserPosts()
      setUserPosts(data)
      setError(null)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load user posts')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    loadUserPosts()
  }, [])

  return (
    <main className="min-h-screen p-8">
      <div className="max-w-4xl mx-auto">
        <div className="mb-6">
          <Link href="/" className="text-blue-500 hover:underline">
            ← トップページに戻る
          </Link>
        </div>

        <h1 className="text-3xl font-bold mb-4">ユーザーと投稿（JOIN）</h1>

        <div className="mb-6 p-4 bg-blue-50 border border-blue-200 rounded">
          <h2 className="font-semibold mb-2">🔀 クロスシャードクエリ</h2>
          <p className="text-sm text-gray-700">
            このページでは、複数のShardからユーザーと投稿をJOINして取得しています。
            各Shardから並列にデータを取得し、アプリケーション層でマージして表示しています。
          </p>
        </div>

        {error && (
          <div className="mb-4 p-4 bg-red-100 border border-red-400 text-red-700 rounded">
            {error}
          </div>
        )}

        {loading ? (
          <p>読み込み中...</p>
        ) : userPosts.length === 0 ? (
          <div className="text-center py-8">
            <p className="text-gray-500 mb-4">表示する投稿がありません。</p>
            <div className="space-x-4">
              <Link href="/users" className="text-blue-500 hover:underline">
                ユーザーを作成
              </Link>
              <span className="text-gray-400">|</span>
              <Link href="/posts" className="text-blue-500 hover:underline">
                投稿を作成
              </Link>
            </div>
          </div>
        ) : (
          <div className="space-y-4">
            {userPosts.map((item, index) => (
              <div key={index} className="p-6 border rounded-lg hover:shadow-md transition-shadow">
                <div className="flex items-start justify-between mb-3">
                  <div className="flex-1">
                    <h3 className="font-bold text-xl mb-1">{item.post_title}</h3>
                    <div className="flex items-center space-x-3 text-sm text-gray-600">
                      <span className="font-medium">{item.user_name}</span>
                      <span>•</span>
                      <span>{item.user_email}</span>
                    </div>
                  </div>
                </div>

                <p className="text-gray-700 mb-3">{item.post_content}</p>

                <div className="text-xs text-gray-400 space-y-1">
                  <div>投稿ID: {item.post_id} | ユーザーID: {item.user_id}</div>
                  <div>作成日: {new Date(item.created_at).toLocaleString('ja-JP')}</div>
                </div>
              </div>
            ))}
          </div>
        )}

        <div className="mt-8 p-4 border rounded-lg bg-gray-50">
          <h3 className="font-semibold mb-2">Sharding情報</h3>
          <p className="text-sm text-gray-600">
            Hash-based shardingにより、user_idをキーとしてデータが2つのShardに分散されています。
            このページでは両方のShardからデータを取得し、統合して表示しています。
          </p>
        </div>
      </div>
    </main>
  )
}
