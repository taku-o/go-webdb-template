'use client'

import { useState } from 'react'
import { Textarea } from '@/components/ui/textarea'
import { Button } from '@/components/ui/button'

const MAX_CONTENT_LENGTH = 280

interface ReplyFormProps {
  onSubmit: (content: string) => Promise<void>
  isSubmitting?: boolean
}

export function ReplyForm({ onSubmit, isSubmitting = false }: ReplyFormProps) {
  const [content, setContent] = useState('')
  const [error, setError] = useState<string | null>(null)

  const contentLength = content.length
  const isOverLimit = contentLength > MAX_CONTENT_LENGTH
  const isEmpty = contentLength === 0
  const isValid = !isEmpty && !isOverLimit

  const handleContentChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    setContent(e.target.value)
    setError(null)
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    if (!isValid) {
      if (isEmpty) {
        setError('返信内容を入力してください')
      } else if (isOverLimit) {
        setError(`返信内容は${MAX_CONTENT_LENGTH}文字以内で入力してください`)
      }
      return
    }

    try {
      await onSubmit(content)
      setContent('')
      setError(null)
    } catch (err) {
      setError(err instanceof Error ? err.message : '返信に失敗しました')
    }
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-3">
      <div className="relative">
        <Textarea
          value={content}
          onChange={handleContentChange}
          placeholder="返信を入力..."
          className={`min-h-[80px] resize-none ${
            isOverLimit ? 'border-destructive focus-visible:ring-destructive' : ''
          }`}
          aria-label="返信内容"
          aria-describedby="reply-counter reply-error"
          disabled={isSubmitting}
        />
      </div>

      <div className="flex items-center justify-between">
        {/* 文字数カウンター */}
        <div
          id="reply-counter"
          className={`text-sm ${
            isOverLimit
              ? 'text-destructive'
              : contentLength > MAX_CONTENT_LENGTH * 0.9
                ? 'text-yellow-600'
                : 'text-muted-foreground'
          }`}
          aria-live="polite"
        >
          {contentLength} / {MAX_CONTENT_LENGTH}
        </div>

        {/* 返信ボタン */}
        <Button
          type="submit"
          disabled={!isValid || isSubmitting}
          aria-busy={isSubmitting}
        >
          {isSubmitting ? '返信中...' : '返信する'}
        </Button>
      </div>

      {/* バリデーションエラー */}
      {error && (
        <p
          id="reply-error"
          className="text-sm text-destructive"
          role="alert"
          aria-live="assertive"
        >
          {error}
        </p>
      )}
    </form>
  )
}
