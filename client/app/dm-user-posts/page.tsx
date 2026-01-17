'use client'

import { useState, useRef } from 'react'
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
  const hasLoadedRef = useRef(false)

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

  const loadInitialData = async () => {
    if (hasLoadedRef.current) return
    hasLoadedRef.current = true
    await loadUserPosts()
  }

  const setContainerRef = (node: HTMLElement | null) => {
    if (node && !hasLoadedRef.current) {
      loadInitialData()
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

        <h1 className="text-2xl sm:text-3xl font-bold mb-4 sm:mb-6">ユーザーと投稿（JOIN）</h1>

        {/* クロスシャードクエリの説明 */}
        <Card className="mb-4 sm:mb-6 border-primary/20 bg-primary/5">
          <CardHeader>
            <CardTitle className="flex items-center gap-2 text-base sm:text-lg">
              <Info className="h-4 w-4 sm:h-5 sm:w-5" />
              クロスシャードクエリ
            </CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-xs sm:text-sm text-foreground">
              このページでは、複数のShardからユーザーと投稿をJOINして取得しています。
              各Shardから並列にデータを取得し、アプリケーション層でマージして表示しています。
            </p>
          </CardContent>
        </Card>

        {error && (
          <div className="mb-4" role="alert" aria-live="assertive">
            <ErrorAlert message={error} />
          </div>
        )}

        {loading ? (
          <Card>
            <CardContent className="py-8">
              <div role="status" aria-live="polite" aria-label="ユーザーと投稿の一覧を読み込み中">
                <LoadingOverlay message="読み込み中..." />
              </div>
            </CardContent>
          </Card>
        ) : dmUserPosts.length === 0 ? (
          <Card>
            <CardContent className="py-8">
              <div className="text-center" role="status">
                <p className="text-muted-foreground mb-4">表示する投稿がありません。</p>
                <nav aria-label="作成ページへのリンク">
                  <div className="space-x-4">
                    <Link href="/dm-users" className="text-primary hover:underline" aria-label="ユーザー作成ページへ">
                      ユーザーを作成
                    </Link>
                    <span className="text-muted-foreground" aria-hidden="true">|</span>
                    <Link href="/dm-posts" className="text-primary hover:underline" aria-label="投稿作成ページへ">
                      投稿を作成
                    </Link>
                  </div>
                </nav>
              </div>
            </CardContent>
          </Card>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4 sm:gap-6 mb-6 sm:mb-8">
            {dmUserPosts.map((item, index) => (
              <Card key={index} className="hover:shadow-lg transition-shadow">
                <CardHeader>
                  <CardTitle className="text-lg sm:text-xl">{item.post_title}</CardTitle>
                  <CardDescription>
                    <div className="flex flex-col sm:flex-row items-start sm:items-center space-y-1 sm:space-y-0 sm:space-x-3 text-xs sm:text-sm">
                      <span className="font-medium">{item.user_name}</span>
                      <span className="hidden sm:inline">•</span>
                      <span>{item.user_email}</span>
                    </div>
                  </CardDescription>
                </CardHeader>
                <CardContent>
                  <p className="text-foreground mb-4 whitespace-pre-wrap text-sm sm:text-base">{item.post_content}</p>
                  <div className="text-xs text-muted-foreground space-y-1 border-t pt-3">
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
        <Card className="bg-muted">
          <CardHeader>
            <CardTitle>Sharding情報</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-sm text-muted-foreground">
              Hash-based shardingにより、user_idをキーとしてデータが2つのShardに分散されています。
              このページでは両方のShardからデータを取得し、統合して表示しています。
            </p>
          </CardContent>
        </Card>
      </div>
    </main>
  )
}
