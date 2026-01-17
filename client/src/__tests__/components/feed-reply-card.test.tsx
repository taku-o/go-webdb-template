import { render, screen } from '@testing-library/react'
import { FeedReplyCard } from '@/components/feed/feed-reply-card'
import { DmFeedReply } from '@/types/dm_feed'

describe('FeedReplyCard', () => {
  const mockReply: DmFeedReply = {
    id: 'reply-001',
    postId: 'post-001',
    userId: 'user-002',
    userName: 'Reply User',
    userHandle: '@replyuser',
    content: 'This is a test reply content',
    createdAt: '2024-01-15T11:00:00Z',
  }

  it('renders reply content correctly', () => {
    render(<FeedReplyCard reply={mockReply} />)

    expect(screen.getByText('Reply User')).toBeInTheDocument()
    expect(screen.getByText('@replyuser')).toBeInTheDocument()
    expect(screen.getByText('This is a test reply content')).toBeInTheDocument()
  })

  it('displays user avatar with first character of userName', () => {
    render(<FeedReplyCard reply={mockReply} />)

    expect(screen.getByText('R')).toBeInTheDocument()
  })

  it('displays time element with correct dateTime attribute', () => {
    render(<FeedReplyCard reply={mockReply} />)

    const timeElement = screen.getByRole('time')
    expect(timeElement).toHaveAttribute('dateTime', '2024-01-15T11:00:00Z')
  })

  it('handles long content correctly', () => {
    const longContent = 'This is a very long reply content that should be displayed correctly without any issues.'
    const longReply: DmFeedReply = {
      ...mockReply,
      content: longContent,
    }

    render(<FeedReplyCard reply={longReply} />)

    expect(screen.getByText(longContent)).toBeInTheDocument()
  })

  it('handles special characters in userName', () => {
    const specialReply: DmFeedReply = {
      ...mockReply,
      userName: '特殊文字ユーザー',
    }

    render(<FeedReplyCard reply={specialReply} />)

    expect(screen.getByText('特殊文字ユーザー')).toBeInTheDocument()
    expect(screen.getByText('特')).toBeInTheDocument()
  })
})
