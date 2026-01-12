'use client'

import { useState, useEffect } from 'react'
import Link from 'next/link'
import { apiClient } from '@/lib/api'
import { DmUserPost } from '@/types/dm_post'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { ErrorAlert } from '@/components/shared/error-alert'
import { LoadingOverlay } from '@/components/shared/loading-overlay'
import { ArrowLeft, Info } from 'lucide-react'

export default function UserPostsPage() {
  const [dmUserPosts, setDmUserPosts] = useState<DmUserPost[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const loadUserPosts = async () => {
    try {
      setLoading(true)
      const data = await apiClient.getDmUserPosts()
      setDmUserPosts(data)
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
      <div className="max-w-6xl mx-auto">
        <div className="mb-6">
          <Link href="/" className="inline-flex items-center text-blue-600 hover:text-blue-800 hover:underline">
            <ArrowLeft className="mr-2 h-4 w-4" />
            トップページに戻る
          </Link>
        </div>

        <h1 className="text-3xl font-bold mb-4">ユーザーと投稿（JOIN）</h1>

        {/* クロスシャードクエリの説明 */}
        <Card className="mb-6 border-blue-200 bg-blue-50">
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Info className="h-5 w-5" />
              クロスシャードクエリ
            </CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-sm text-gray-700">
              このページでは、複数のShardからユーザーと投稿をJOINして取得しています。
              各Shardから並列にデータを取得し、アプリケーション層でマージして表示しています。
            </p>
          </CardContent>
        </Card>

        {error && (
          <div className="mb-4">
            <ErrorAlert message={error} />
          </div>
        )}

        {loading ? (
          <Card>
            <CardContent className="py-8">
              <LoadingOverlay message="読み込み中..." />
            </CardContent>
          </Card>
        ) : dmUserPosts.length === 0 ? (
          <Card>
            <CardContent className="py-8">
              <div className="text-center">
                <p className="text-gray-500 mb-4">表示する投稿がありません。</p>
                <div className="space-x-4">
                  <Link href="/dm-users" className="text-blue-600 hover:underline">
                    ユーザーを作成
                  </Link>
                  <span className="text-gray-400">|</span>
                  <Link href="/dm-posts" className="text-blue-600 hover:underline">
                    投稿を作成
                  </Link>
                </div>
              </div>
            </CardContent>
          </Card>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-8">
            {dmUserPosts.map((item, index) => (
              <Card key={index} className="hover:shadow-lg transition-shadow">
                <CardHeader>
                  <CardTitle className="text-xl">{item.post_title}</CardTitle>
                  <CardDescription>
                    <div className="flex items-center space-x-3 text-sm">
                      <span className="font-medium">{item.user_name}</span>
                      <span>•</span>
                      <span>{item.user_email}</span>
                    </div>
                  </CardDescription>
                </CardHeader>
                <CardContent>
                  <p className="text-gray-700 mb-4 whitespace-pre-wrap">{item.post_content}</p>
                  <div className="text-xs text-gray-400 space-y-1 border-t pt-3">
                    <div className="font-mono">
                      投稿ID: {item.post_id} | ユーザーID: {item.user_id}
                    </div>
                    <div>
                      作成日: {new Date(item.created_at).toLocaleString('ja-JP')}
                    </div>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        )}

        {/* Sharding情報 */}
        <Card className="bg-gray-50">
          <CardHeader>
            <CardTitle>Sharding情報</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-sm text-gray-600">
              Hash-based shardingにより、user_idをキーとしてデータが2つのShardに分散されています。
              このページでは両方のShardからデータを取得し、統合して表示しています。
            </p>
          </CardContent>
        </Card>
      </div>
    </main>
  )
}
