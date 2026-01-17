import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { FeedPostCard } from '@/components/feed/feed-post-card'
import { DmFeedPost } from '@/types/dm_feed'

// Mock next/link
jest.mock('next/link', () => {
  return ({ children, href }: { children: React.ReactNode; href: string }) => (
    <a href={href}>{children}</a>
  )
})

describe('FeedPostCard', () => {
  const mockPost: DmFeedPost = {
    id: 'post-001',
    userId: 'user-001',
    userName: 'Test User',
    userHandle: '@testuser',
    content: 'This is a test post content',
    createdAt: '2024-01-15T10:00:00Z',
    likeCount: 5,
    liked: false,
  }

  const mockUserId = 'user-001'

  beforeEach(() => {
    jest.clearAllMocks()
  })

  it('renders post content correctly', () => {
    render(<FeedPostCard post={mockPost} userId={mockUserId} />)

    expect(screen.getByText('Test User')).toBeInTheDocument()
    expect(screen.getByText('@testuser')).toBeInTheDocument()
    expect(screen.getByText('This is a test post content')).toBeInTheDocument()
  })

  it('displays user avatar with first character of userName', () => {
    render(<FeedPostCard post={mockPost} userId={mockUserId} />)

    expect(screen.getByText('T')).toBeInTheDocument()
  })

  it('displays like count', () => {
    render(<FeedPostCard post={mockPost} userId={mockUserId} />)

    expect(screen.getByText('5')).toBeInTheDocument()
  })

  it('displays reply link with correct href', () => {
    render(<FeedPostCard post={mockPost} userId={mockUserId} />)

    const replyLink = screen.getByRole('link', { name: '返信' })
    expect(replyLink).toHaveAttribute('href', '/dm_feed/user-001/post-001')
  })

  it('calls onLikeToggle when like button is clicked', async () => {
    const user = userEvent.setup()
    const mockOnLikeToggle = jest.fn()

    render(
      <FeedPostCard
        post={mockPost}
        userId={mockUserId}
        onLikeToggle={mockOnLikeToggle}
      />
    )

    const likeButton = screen.getByRole('button', { name: 'いいねする' })
    await user.click(likeButton)

    expect(mockOnLikeToggle).toHaveBeenCalledWith('post-001')
  })

  it('displays correct aria-label for liked post', () => {
    const likedPost: DmFeedPost = {
      ...mockPost,
      liked: true,
      likeCount: 6,
    }

    render(<FeedPostCard post={likedPost} userId={mockUserId} />)

    const likeButton = screen.getByRole('button', { name: 'いいねを取り消す' })
    expect(likeButton).toHaveAttribute('aria-pressed', 'true')
  })

  it('displays correct aria-label for not liked post', () => {
    render(<FeedPostCard post={mockPost} userId={mockUserId} />)

    const likeButton = screen.getByRole('button', { name: 'いいねする' })
    expect(likeButton).toHaveAttribute('aria-pressed', 'false')
  })

  it('does not call onLikeToggle when not provided', async () => {
    const user = userEvent.setup()

    render(<FeedPostCard post={mockPost} userId={mockUserId} />)

    const likeButton = screen.getByRole('button', { name: 'いいねする' })
    await user.click(likeButton)

    // Should not throw error when onLikeToggle is undefined
  })
})
