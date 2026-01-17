'use client'

import { useState, useEffect, useRef, useCallback } from 'react'
import { useParams } from 'next/navigation'
import Link from 'next/link'
import { ArrowLeft } from 'lucide-react'
import { DmFeedPost } from '@/types/dm_feed'
import { apiClient } from '@/lib/api'
import { LoadingSpinner } from '@/components/shared/loading-spinner'
import { ErrorAlert } from '@/components/shared/error-alert'
import { FeedPostCard } from '@/components/feed/feed-post-card'
import { FeedForm } from '@/components/feed/feed-form'

const POSTS_PER_PAGE = 10

export default function FeedPage() {
  const params = useParams()
  const userId = params.userId as string

  // 投稿一覧の状態
  const [dmFeedPosts, setDmFeedPosts] = useState<DmFeedPost[]>([])
  // ローディング状態
  const [isLoading, setIsLoading] = useState(false)
  // さらに古い投稿があるかどうか
  const [hasMore, setHasMore] = useState(true)
  // 最新の投稿ID（上方向スクロール用）
  const [newestPostId, setNewestPostId] = useState<string>('')
  // 最古の投稿ID（下方向スクロール用）
  const [oldestPostId, setOldestPostId] = useState<string>('')
  // エラー状態
  const [error, setError] = useState<string | null>(null)
  // 古い投稿を読み込み中かどうか
  const [isLoadingMore, setIsLoadingMore] = useState(false)
  // 新しい投稿を読み込み中かどうか
  const [isLoadingNewer, setIsLoadingNewer] = useState(false)
  // 新規投稿中かどうか
  const [isSubmitting, setIsSubmitting] = useState(false)

  // 下方向スクロール検知用のref
  const loadMoreRef = useRef<HTMLDivElement>(null)
  // 上方向スクロール検知用のref
  const loadNewerRef = useRef<HTMLDivElement>(null)
  // 投稿一覧のコンテナref（スクロール位置維持用）
  const postsContainerRef = useRef<HTMLDivElement>(null)
  // 初期化済みフラグ
  const hasLoadedRef = useRef(false)
  // 前回のuserId
  const currentUserIdRef = useRef<string | null>(null)

  // 初期データの読み込み
  const loadInitialPosts = async () => {
    try {
      setIsLoading(true)
      setError(null)
      const posts = await apiClient.getDmFeedPosts(userId, POSTS_PER_PAGE, '')
      setDmFeedPosts(posts)

      if (posts.length > 0) {
        setNewestPostId(posts[0].id)
        setOldestPostId(posts[posts.length - 1].id)
      }

      setHasMore(posts.length >= POSTS_PER_PAGE)
    } catch (err) {
      setError(err instanceof Error ? err.message : '投稿の読み込みに失敗しました')
    } finally {
      setIsLoading(false)
    }
  }

  // 古い投稿の読み込み
  const loadOlderPosts = useCallback(async () => {
    if (isLoadingMore || !hasMore || !oldestPostId) return

    try {
      setIsLoadingMore(true)
      const posts = await apiClient.getDmFeedPosts(userId, POSTS_PER_PAGE, oldestPostId)

      if (posts.length > 0) {
        setDmFeedPosts((prev) => [...prev, ...posts])
        setOldestPostId(posts[posts.length - 1].id)
      }

      setHasMore(posts.length >= POSTS_PER_PAGE)
    } catch (err) {
      setError(err instanceof Error ? err.message : '投稿の読み込みに失敗しました')
    } finally {
      setIsLoadingMore(false)
    }
  }, [userId, oldestPostId, isLoadingMore, hasMore])

  // 最新投稿の読み込み
  const loadNewerPosts = useCallback(async () => {
    if (isLoadingNewer || !newestPostId) return

    try {
      setIsLoadingNewer(true)

      // スクロール位置維持のため、現在のスクロール位置とコンテナの高さを記録
      const container = postsContainerRef.current
      const previousScrollHeight = container?.scrollHeight || 0

      // 最新投稿を取得
      const posts = await apiClient.getDmFeedPosts(userId, POSTS_PER_PAGE, '')

      if (posts.length > 0) {
        // 既存の投稿IDのセットを作成して重複チェック
        setDmFeedPosts((prev) => {
          const existingIds = new Set(prev.map((p) => p.id))
          const newPosts = posts.filter((p) => !existingIds.has(p.id))

          if (newPosts.length > 0) {
            // 新しい投稿を上部に追加
            return [...newPosts, ...prev]
          }
          return prev
        })

        // newestPostIdを更新
        setNewestPostId(posts[0].id)

        // スクロール位置を維持
        requestAnimationFrame(() => {
          if (container) {
            const newScrollHeight = container.scrollHeight
            const heightDiff = newScrollHeight - previousScrollHeight
            window.scrollBy(0, heightDiff)
          }
        })
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : '投稿の読み込みに失敗しました')
    } finally {
      setIsLoadingNewer(false)
    }
  }, [userId, newestPostId, isLoadingNewer])

  // 新規投稿の処理
  const handleSubmitPost = async (content: string) => {
    try {
      setIsSubmitting(true)
      setError(null)
      const newPost = await apiClient.createDmFeedPost(userId, content)

      // 投稿一覧の一番上に新規投稿を追加
      setDmFeedPosts((prev) => [newPost, ...prev])

      // newestPostIdを更新
      setNewestPostId(newPost.id)

      // 投稿が空だった場合はoldestPostIdも設定
      if (!oldestPostId) {
        setOldestPostId(newPost.id)
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : '投稿に失敗しました')
      throw err
    } finally {
      setIsSubmitting(false)
    }
  }

  // 初期化処理
  const loadInitialData = async () => {
    if (hasLoadedRef.current && currentUserIdRef.current === userId) return
    hasLoadedRef.current = true
    currentUserIdRef.current = userId
    await loadInitialPosts()
  }

  // refコールバック関数
  const setContainerRef = (node: HTMLElement | null) => {
    if (node && (!hasLoadedRef.current || currentUserIdRef.current !== userId)) {
      loadInitialData()
    }
  }

  // userIdが変更された場合の処理
  if (currentUserIdRef.current !== null && currentUserIdRef.current !== userId) {
    currentUserIdRef.current = userId
    hasLoadedRef.current = false
    loadInitialPosts()
  }

  // いいねのトグル処理
  const handleLikeToggle = async (postId: string) => {
    try {
      const updatedPost = await apiClient.toggleLikeDmPost(userId, postId)

      // 投稿一覧の該当投稿を更新
      setDmFeedPosts((prev) =>
        prev.map((post) =>
          post.id === postId
            ? { ...post, liked: updatedPost.liked, likeCount: updatedPost.likeCount }
            : post
        )
      )
    } catch (err) {
      setError(err instanceof Error ? err.message : 'いいねに失敗しました')
    }
  }

  // Intersection Observer で下方向スクロールを検知
  useEffect(() => {
    const observer = new IntersectionObserver(
      (entries) => {
        if (entries[0].isIntersecting && hasMore && !isLoadingMore) {
          loadOlderPosts()
        }
      },
      { threshold: 0.1 }
    )

    const currentRef = loadMoreRef.current
    if (currentRef) {
      observer.observe(currentRef)
    }

    return () => {
      if (currentRef) {
        observer.unobserve(currentRef)
      }
    }
  }, [hasMore, isLoadingMore, loadOlderPosts])

  // Intersection Observer で上方向スクロールを検知
  useEffect(() => {
    const observer = new IntersectionObserver(
      (entries) => {
        if (entries[0].isIntersecting && !isLoadingNewer && dmFeedPosts.length > 0) {
          loadNewerPosts()
        }
      },
      { threshold: 0.1 }
    )

    const currentRef = loadNewerRef.current
    if (currentRef) {
      observer.observe(currentRef)
    }

    return () => {
      if (currentRef) {
        observer.unobserve(currentRef)
      }
    }
  }, [isLoadingNewer, loadNewerPosts, dmFeedPosts.length])

  return (
    <main ref={setContainerRef} className="min-h-screen p-4 sm:p-6 md:p-8">
      <div className="max-w-2xl mx-auto">
        <nav aria-label="パンくずリスト">
          <div className="mb-4 sm:mb-6">
            <Link
              href="/"
              className="inline-flex items-center text-primary hover:underline text-sm sm:text-base"
              aria-label="トップページに戻る"
            >
              <ArrowLeft className="mr-2 h-4 w-4" aria-hidden="true" />
              トップページに戻る
            </Link>
          </div>
        </nav>

        <h1 className="text-2xl sm:text-3xl font-bold mb-6 sm:mb-8">フィード</h1>

        {/* エラー表示 */}
        {error && (
          <div className="mb-4" role="alert" aria-live="assertive">
            <ErrorAlert message={error} />
          </div>
        )}

        {/* 新規投稿フォーム */}
        <section aria-label="新規投稿フォーム" className="mb-6">
          <FeedForm onSubmit={handleSubmitPost} isSubmitting={isSubmitting} />
        </section>

        {/* 投稿一覧 */}
        <section aria-label="投稿一覧">
          {isLoading && dmFeedPosts.length === 0 ? (
            <div className="flex justify-center py-8" role="status" aria-live="polite">
              <LoadingSpinner size="lg" />
              <span className="sr-only">投稿を読み込み中...</span>
            </div>
          ) : dmFeedPosts.length === 0 ? (
            <p className="text-center text-muted-foreground py-8">
              投稿がありません
            </p>
          ) : (
            <div ref={postsContainerRef} className="space-y-4">
              {/* 上方向スクロール検知用の要素 */}
              <div ref={loadNewerRef} className="py-2">
                {isLoadingNewer && (
                  <div className="flex justify-center" role="status" aria-live="polite">
                    <LoadingSpinner size="md" />
                    <span className="sr-only">最新の投稿を読み込み中...</span>
                  </div>
                )}
              </div>

              {dmFeedPosts.map((post) => (
                <FeedPostCard
                  key={post.id}
                  post={post}
                  userId={userId}
                  onLikeToggle={handleLikeToggle}
                />
              ))}

              {/* 下方向スクロール検知用の要素 */}
              <div ref={loadMoreRef} className="py-4">
                {isLoadingMore ? (
                  <div className="flex justify-center" role="status" aria-live="polite">
                    <LoadingSpinner size="md" />
                    <span className="sr-only">古い投稿を読み込み中...</span>
                  </div>
                ) : !hasMore ? (
                  <p className="text-center text-muted-foreground text-sm">
                    これ以上投稿はありません
                  </p>
                ) : null}
              </div>
            </div>
          )}
        </section>
      </div>
    </main>
  )
}
