'use client'

import { useState, useRef, useCallback } from 'react'
import { useParams } from 'next/navigation'
import Link from 'next/link'
import { ArrowLeft } from 'lucide-react'
import { DmFeedPost, DmFeedReply } from '@/types/dm_feed'
import { apiClient } from '@/lib/api'
import { LoadingSpinner } from '@/components/shared/loading-spinner'
import { ErrorAlert } from '@/components/shared/error-alert'
import { FeedPostCard } from '@/components/feed/feed-post-card'
import { FeedReplyCard } from '@/components/feed/feed-reply-card'
import { ReplyForm } from '@/components/feed/reply-form'

const REPLIES_PER_PAGE = 10

export default function ReplyPage() {
  const params = useParams()
  const userId = params.userId as string
  const postId = params.postId as string

  // 返信元の投稿の状態
  const [dmFeedPost, setDmFeedPost] = useState<DmFeedPost | null>(null)
  // 返信一覧の状態
  const [dmFeedReplies, setDmFeedReplies] = useState<DmFeedReply[]>([])
  // ローディング状態
  const [isLoading, setIsLoading] = useState(false)
  // さらに古い返信があるかどうか
  const [hasMore, setHasMore] = useState(true)
  // 最新の返信ID（下方向スクロール用）
  const [newestReplyId, setNewestReplyId] = useState<string>('')
  // 最古の返信ID（上方向スクロール用）
  const [oldestReplyId, setOldestReplyId] = useState<string>('')
  // エラー状態
  const [error, setError] = useState<string | null>(null)
  // 新しい返信を読み込み中かどうか
  const [isLoadingNewer, setIsLoadingNewer] = useState(false)
  // 古い返信を読み込み中かどうか
  const [isLoadingOlder, setIsLoadingOlder] = useState(false)
  // 返信投稿中かどうか
  const [isSubmittingReply, setIsSubmittingReply] = useState(false)

  // 下方向スクロール検知用のObserverのref（新しい返信読み込み）
  const loadNewerObserverRef = useRef<IntersectionObserver | null>(null)
  // 上方向スクロール検知用のObserverのref（古い返信読み込み）
  const loadOlderObserverRef = useRef<IntersectionObserver | null>(null)
  // 返信一覧のコンテナref（スクロール位置維持用）
  const repliesContainerRef = useRef<HTMLDivElement>(null)
  // 初期化済みフラグ
  const hasLoadedRef = useRef(false)
  // 前回のuserIdとpostId
  const currentUserIdRef = useRef<string | null>(null)
  const currentPostIdRef = useRef<string | null>(null)

  // 初期データの読み込み（返信元の投稿と返信一覧）
  const loadInitialDataInternal = async () => {
    try {
      setIsLoading(true)
      setError(null)

      // 返信元の投稿を取得
      const post = await apiClient.getDmFeedPostById(userId, postId)
      if (post) {
        setDmFeedPost(post)
      } else {
        setError('投稿が見つかりませんでした')
        return
      }

      // 返信一覧を取得（古いものが上、新しいものが下）
      const replies = await apiClient.getDmFeedReplies(userId, postId, REPLIES_PER_PAGE, '')
      setDmFeedReplies(replies)

      if (replies.length > 0) {
        setOldestReplyId(replies[0].id)
        setNewestReplyId(replies[replies.length - 1].id)
      }

      setHasMore(replies.length >= REPLIES_PER_PAGE)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'データの読み込みに失敗しました')
    } finally {
      setIsLoading(false)
    }
  }

  // 初期化処理
  const loadInitialData = async () => {
    if (hasLoadedRef.current && currentUserIdRef.current === userId && currentPostIdRef.current === postId) return
    hasLoadedRef.current = true
    currentUserIdRef.current = userId
    currentPostIdRef.current = postId
    await loadInitialDataInternal()
  }

  // refコールバック関数
  const setContainerRef = (node: HTMLElement | null) => {
    if (node && (!hasLoadedRef.current || currentUserIdRef.current !== userId || currentPostIdRef.current !== postId)) {
      loadInitialData()
    }
  }

  // userIdまたはpostIdが変更された場合の処理
  if (currentUserIdRef.current !== null && currentPostIdRef.current !== null &&
      (currentUserIdRef.current !== userId || currentPostIdRef.current !== postId)) {
    currentUserIdRef.current = userId
    currentPostIdRef.current = postId
    hasLoadedRef.current = false
    loadInitialDataInternal()
  }

  // いいねのトグル処理
  const handleLikeToggle = async (targetPostId: string) => {
    try {
      const updatedPost = await apiClient.toggleLikeDmPost(userId, targetPostId)

      // 返信元の投稿を更新
      if (dmFeedPost && dmFeedPost.id === targetPostId) {
        setDmFeedPost({
          ...dmFeedPost,
          liked: updatedPost.liked,
          likeCount: updatedPost.likeCount,
        })
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : 'いいねに失敗しました')
    }
  }

  // 返信の投稿処理
  const handleSubmitReply = async (content: string) => {
    try {
      setIsSubmittingReply(true)
      setError(null)
      const newReply = await apiClient.replyToDmPost(userId, postId, content)

      // 返信一覧の一番下に新規返信を追加
      setDmFeedReplies((prev) => [...prev, newReply])

      // newestReplyIdを更新
      setNewestReplyId(newReply.id)

      // 返信が空だった場合はoldestReplyIdも設定
      if (!oldestReplyId) {
        setOldestReplyId(newReply.id)
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : '返信に失敗しました')
      throw err
    } finally {
      setIsSubmittingReply(false)
    }
  }

  // 新しい返信の読み込み（下方向スクロール時）
  const loadNewerReplies = useCallback(async () => {
    if (isLoadingNewer || !newestReplyId) return

    try {
      setIsLoadingNewer(true)

      // 最新の返信を取得（newestReplyIdより新しいものを取得するため、空文字で最新から取得）
      const replies = await apiClient.getDmFeedReplies(userId, postId, REPLIES_PER_PAGE, '')

      if (replies.length > 0) {
        // 既存の返信IDのセットを作成して重複チェック
        setDmFeedReplies((prev) => {
          const existingIds = new Set(prev.map((r) => r.id))
          const newReplies = replies.filter((r) => !existingIds.has(r.id))

          if (newReplies.length > 0) {
            // 新しい返信を下部に追加（返信は下が新しい）
            return [...prev, ...newReplies]
          }
          return prev
        })

        // newestReplyIdを更新
        setNewestReplyId(replies[replies.length - 1].id)
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : '返信の読み込みに失敗しました')
    } finally {
      setIsLoadingNewer(false)
    }
  }, [userId, postId, newestReplyId, isLoadingNewer])

  // 古い返信の読み込み（上方向スクロール時）
  const loadOlderReplies = useCallback(async () => {
    if (isLoadingOlder || !hasMore || !oldestReplyId) return

    try {
      setIsLoadingOlder(true)

      // スクロール位置維持のため、現在のスクロール位置とコンテナの高さを記録
      const container = repliesContainerRef.current
      const previousScrollHeight = container?.scrollHeight || 0

      // 古い返信を取得
      const replies = await apiClient.getDmFeedReplies(userId, postId, REPLIES_PER_PAGE, oldestReplyId)

      if (replies.length > 0) {
        // 古い返信を上部に追加
        setDmFeedReplies((prev) => [...replies, ...prev])

        // oldestReplyIdを更新
        setOldestReplyId(replies[0].id)

        // スクロール位置を維持
        requestAnimationFrame(() => {
          if (container) {
            const newScrollHeight = container.scrollHeight
            const heightDiff = newScrollHeight - previousScrollHeight
            window.scrollBy(0, heightDiff)
          }
        })
      }

      setHasMore(replies.length >= REPLIES_PER_PAGE)
    } catch (err) {
      setError(err instanceof Error ? err.message : '返信の読み込みに失敗しました')
    } finally {
      setIsLoadingOlder(false)
    }
  }, [userId, postId, oldestReplyId, isLoadingOlder, hasMore])

  // 上方向スクロール検知用のrefコールバック関数（古い返信読み込み）
  const setLoadOlderRef = useCallback((node: HTMLDivElement | null) => {
    // 既存のObserverをクリーンアップ
    if (loadOlderObserverRef.current) {
      loadOlderObserverRef.current.disconnect()
      loadOlderObserverRef.current = null
    }

    if (node) {
      const observer = new IntersectionObserver(
        (entries) => {
          if (entries[0].isIntersecting && hasMore && !isLoadingOlder) {
            loadOlderReplies()
          }
        },
        { threshold: 0.1 }
      )
      observer.observe(node)
      loadOlderObserverRef.current = observer
    }
  }, [hasMore, isLoadingOlder, loadOlderReplies])

  // 下方向スクロール検知用のrefコールバック関数（新しい返信読み込み）
  const setLoadNewerRef = useCallback((node: HTMLDivElement | null) => {
    // 既存のObserverをクリーンアップ
    if (loadNewerObserverRef.current) {
      loadNewerObserverRef.current.disconnect()
      loadNewerObserverRef.current = null
    }

    if (node) {
      const observer = new IntersectionObserver(
        (entries) => {
          if (entries[0].isIntersecting && !isLoadingNewer && dmFeedReplies.length > 0) {
            loadNewerReplies()
          }
        },
        { threshold: 0.1 }
      )
      observer.observe(node)
      loadNewerObserverRef.current = observer
    }
  }, [isLoadingNewer, loadNewerReplies, dmFeedReplies.length])

  return (
    <main ref={setContainerRef} className="min-h-screen p-4 sm:p-6 md:p-8">
      <div className="max-w-2xl mx-auto">
        <nav aria-label="パンくずリスト">
          <div className="mb-4 sm:mb-6">
            <Link
              href={`/dm_feed/${userId}`}
              className="inline-flex items-center text-primary hover:underline text-sm sm:text-base"
              aria-label="フィードに戻る"
            >
              <ArrowLeft className="mr-2 h-4 w-4" aria-hidden="true" />
              フィードに戻る
            </Link>
          </div>
        </nav>

        <h1 className="text-2xl sm:text-3xl font-bold mb-6 sm:mb-8">返信</h1>

        {/* エラー表示 */}
        {error && (
          <div className="mb-4" role="alert" aria-live="assertive">
            <ErrorAlert message={error} />
          </div>
        )}

        {/* ローディング表示 */}
        {isLoading ? (
          <div className="flex justify-center py-8" role="status" aria-live="polite">
            <LoadingSpinner size="lg" />
            <span className="sr-only">読み込み中...</span>
          </div>
        ) : (
          <>
            {/* 返信元のエントリー */}
            <section aria-label="返信元の投稿" className="mb-6">
              {dmFeedPost && (
                <FeedPostCard
                  post={dmFeedPost}
                  userId={userId}
                  onLikeToggle={handleLikeToggle}
                />
              )}
            </section>

            {/* 返信フォーム */}
            <section aria-label="返信フォーム" className="mb-6">
              <ReplyForm onSubmit={handleSubmitReply} isSubmitting={isSubmittingReply} />
            </section>

            {/* 返信一覧 */}
            <section aria-label="返信一覧">
              {dmFeedReplies.length === 0 ? (
                <p className="text-center text-muted-foreground py-8">
                  返信がありません
                </p>
              ) : (
                <div ref={repliesContainerRef} className="space-y-4">
                  {/* 上方向スクロール検知用の要素（古い返信読み込み） */}
                  <div ref={setLoadOlderRef} className="py-2">
                    {isLoadingOlder ? (
                      <div className="flex justify-center" role="status" aria-live="polite">
                        <LoadingSpinner size="md" />
                        <span className="sr-only">古い返信を読み込み中...</span>
                      </div>
                    ) : !hasMore ? (
                      <p className="text-center text-muted-foreground text-sm">
                        これ以上返信はありません
                      </p>
                    ) : null}
                  </div>

                  {dmFeedReplies.map((reply) => (
                    <FeedReplyCard key={reply.id} reply={reply} />
                  ))}

                  {/* 下方向スクロール検知用の要素（新しい返信読み込み） */}
                  <div ref={setLoadNewerRef} className="py-4">
                    {isLoadingNewer && (
                      <div className="flex justify-center" role="status" aria-live="polite">
                        <LoadingSpinner size="md" />
                        <span className="sr-only">新しい返信を読み込み中...</span>
                      </div>
                    )}
                  </div>
                </div>
              )}
            </section>
          </>
        )}
      </div>
    </main>
  )
}
