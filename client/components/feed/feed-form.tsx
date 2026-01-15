'use client'

import { useState } from 'react'
import { Textarea } from '@/components/ui/textarea'
import { Button } from '@/components/ui/button'

const MAX_CONTENT_LENGTH = 280

interface FeedFormProps {
  onSubmit: (content: string) => Promise<void>
  isSubmitting?: boolean
}

export function FeedForm({ onSubmit, isSubmitting = false }: FeedFormProps) {
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
        setError('投稿内容を入力してください')
      } else if (isOverLimit) {
        setError(`投稿内容は${MAX_CONTENT_LENGTH}文字以内で入力してください`)
      }
      return
    }

    try {
      await onSubmit(content)
      setContent('')
      setError(null)
    } catch (err) {
      setError(err instanceof Error ? err.message : '投稿に失敗しました')
    }
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-3">
      <div className="relative">
        <Textarea
          value={content}
          onChange={handleContentChange}
          placeholder="いまどうしてる？"
          className={`min-h-[100px] resize-none ${
            isOverLimit ? 'border-destructive focus-visible:ring-destructive' : ''
          }`}
          aria-label="投稿内容"
          aria-describedby="content-counter content-error"
          disabled={isSubmitting}
        />
      </div>

      <div className="flex items-center justify-between">
        {/* 文字数カウンター */}
        <div
          id="content-counter"
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

        {/* 投稿ボタン */}
        <Button
          type="submit"
          disabled={!isValid || isSubmitting}
          aria-busy={isSubmitting}
        >
          {isSubmitting ? '投稿中...' : '投稿する'}
        </Button>
      </div>

      {/* バリデーションエラー */}
      {error && (
        <p
          id="content-error"
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
