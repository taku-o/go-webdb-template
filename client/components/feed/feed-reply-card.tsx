'use client'

import { Card, CardContent } from '@/components/ui/card'
import { DmFeedReply } from '@/types/dm_feed'
import { formatRelativeTime } from '@/lib/utils'

interface FeedReplyCardProps {
  reply: DmFeedReply
}

export function FeedReplyCard({ reply }: FeedReplyCardProps) {
  return (
    <Card className="hover:bg-accent/50 transition-colors">
      <CardContent className="p-4">
        {/* 返信者情報 */}
        <div className="flex items-start gap-3">
          {/* アバター（シンプルな円形） */}
          <div
            className="w-10 h-10 rounded-full bg-primary/10 flex items-center justify-center flex-shrink-0"
            aria-hidden="true"
          >
            <span className="text-sm font-medium text-primary">
              {reply.userName.charAt(0)}
            </span>
          </div>

          <div className="flex-1 min-w-0">
            {/* 返信者名とハンドル */}
            <div className="flex items-center gap-2 flex-wrap">
              <span className="font-bold text-sm truncate">{reply.userName}</span>
              <span className="text-sm text-muted-foreground truncate">
                {reply.userHandle}
              </span>
              <span className="text-sm text-muted-foreground">·</span>
              <time
                className="text-sm text-muted-foreground"
                dateTime={reply.createdAt}
                aria-label={`返信日時: ${formatRelativeTime(reply.createdAt)}`}
              >
                {formatRelativeTime(reply.createdAt)}
              </time>
            </div>

            {/* 返信内容 */}
            <p className="mt-2 text-sm whitespace-pre-wrap break-words">
              {reply.content}
            </p>
          </div>
        </div>
      </CardContent>
    </Card>
  )
}
