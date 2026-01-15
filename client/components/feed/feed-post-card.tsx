'use client'

import Link from 'next/link'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Heart, MessageCircle } from 'lucide-react'
import { DmFeedPost } from '@/types/dm_feed'
import { formatRelativeTime } from '@/lib/utils'

interface FeedPostCardProps {
  post: DmFeedPost
  userId: string
  onLikeToggle?: (postId: string) => void
}

export function FeedPostCard({ post, userId, onLikeToggle }: FeedPostCardProps) {
  const handleLikeClick = () => {
    if (onLikeToggle) {
      onLikeToggle(post.id)
    }
  }

  return (
    <Card className="hover:bg-accent/50 transition-colors">
      <CardContent className="p-4">
        {/* 投稿者情報 */}
        <div className="flex items-start gap-3">
          {/* アバター（シンプルな円形） */}
          <div
            className="w-10 h-10 rounded-full bg-primary/10 flex items-center justify-center flex-shrink-0"
            aria-hidden="true"
          >
            <span className="text-sm font-medium text-primary">
              {post.userName.charAt(0)}
            </span>
          </div>

          <div className="flex-1 min-w-0">
            {/* 投稿者名とハンドル */}
            <div className="flex items-center gap-2 flex-wrap">
              <span className="font-bold text-sm truncate">{post.userName}</span>
              <span className="text-sm text-muted-foreground truncate">
                {post.userHandle}
              </span>
              <span className="text-sm text-muted-foreground">·</span>
              <time
                className="text-sm text-muted-foreground"
                dateTime={post.createdAt}
                aria-label={`投稿日時: ${formatRelativeTime(post.createdAt)}`}
              >
                {formatRelativeTime(post.createdAt)}
              </time>
            </div>

            {/* 投稿内容 */}
            <p className="mt-2 text-sm whitespace-pre-wrap break-words">
              {post.content}
            </p>

            {/* アクションボタン */}
            <div className="flex items-center gap-4 mt-3">
              {/* 返信ボタン */}
              <Link
                href={`/feed/${userId}/${post.id}`}
                className="inline-flex items-center gap-1 text-muted-foreground hover:text-primary transition-colors"
                aria-label="返信一覧を表示"
              >
                <MessageCircle className="h-4 w-4" aria-hidden="true" />
                <span className="text-xs">返信</span>
              </Link>

              {/* いいねボタン */}
              <Button
                variant="ghost"
                size="sm"
                className={`h-auto p-0 gap-1 ${
                  post.liked
                    ? 'text-red-500 hover:text-red-600'
                    : 'text-muted-foreground hover:text-red-500'
                }`}
                onClick={handleLikeClick}
                aria-label={post.liked ? 'いいねを取り消す' : 'いいねする'}
                aria-pressed={post.liked}
              >
                <Heart
                  className={`h-4 w-4 ${post.liked ? 'fill-current' : ''}`}
                  aria-hidden="true"
                />
                <span className="text-xs">{post.likeCount}</span>
              </Button>
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  )
}
