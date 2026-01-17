import { render, screen, waitFor, act } from '@testing-library/react'
import ReplyPage from '@/app/dm_feed/[userId]/[postId]/page'

// Mock IntersectionObserver
const mockIntersectionObserver = jest.fn()
mockIntersectionObserver.mockImplementation((callback) => ({
  observe: jest.fn(),
  unobserve: jest.fn(),
  disconnect: jest.fn(),
}))
window.IntersectionObserver = mockIntersectionObserver

// Mock next/navigation
jest.mock('next/navigation', () => ({
  ...jest.requireActual('next/navigation'),
  useParams: () => ({ userId: 'test-user-001', postId: 'test-post-001' }),
}))

// Mock post data
const mockPost = {
  id: 'test-post-001',
  userId: 'test-user-001',
  userName: 'テストユーザー',
  userHandle: '@test_user',
  content: 'Original post content for reply test',
  createdAt: '2024-01-15T10:00:00Z',
  liked: false,
  likeCount: 5,
}

// Mock replies data
const mockReplies = [
  {
    id: 'reply-001',
    postId: 'test-post-001',
    userId: 'test-user-001',
    userName: 'テストユーザー1',
    userHandle: '@test_user_1',
    content: 'First reply content',
    createdAt: '2024-01-15T11:00:00Z',
  },
  {
    id: 'reply-002',
    postId: 'test-post-001',
    userId: 'test-user-002',
    userName: 'テストユーザー2',
    userHandle: '@test_user_2',
    content: 'Second reply content',
    createdAt: '2024-01-15T12:00:00Z',
  },
]

// Mock apiClient
const mockGetDmFeedPostById = jest.fn()
const mockGetDmFeedReplies = jest.fn()
const mockReplyToDmPost = jest.fn()
const mockToggleLikeDmPost = jest.fn()

jest.mock('@/lib/api', () => ({
  apiClient: {
    getDmFeedPostById: (...args: unknown[]) => mockGetDmFeedPostById(...args),
    getDmFeedReplies: (...args: unknown[]) => mockGetDmFeedReplies(...args),
    replyToDmPost: (...args: unknown[]) => mockReplyToDmPost(...args),
    toggleLikeDmPost: (...args: unknown[]) => mockToggleLikeDmPost(...args),
  },
}))

beforeEach(() => {
  jest.clearAllMocks()
  mockGetDmFeedPostById.mockResolvedValue(mockPost)
  mockGetDmFeedReplies.mockResolvedValue(mockReplies)
  mockReplyToDmPost.mockResolvedValue({
    id: 'reply-003',
    postId: 'test-post-001',
    userId: 'test-user-001',
    userName: 'あなた',
    userHandle: '@you',
    content: 'New reply content',
    createdAt: new Date().toISOString(),
  })
  mockToggleLikeDmPost.mockResolvedValue({ liked: true, likeCount: 6 })
})

describe('ReplyPage Integration', () => {
  it('displays page title', async () => {
    await act(async () => {
      render(<ReplyPage />)
    })

    // Wait for async operations to complete
    await waitFor(() => {
      expect(screen.getByText('Original post content for reply test')).toBeInTheDocument()
    })

    expect(screen.getByRole('heading', { name: '返信', level: 1 })).toBeInTheDocument()
  })

  it('displays original post and replies', async () => {
    await act(async () => {
      render(<ReplyPage />)
    })

    // Wait for data to load
    await waitFor(() => {
      expect(screen.getByText('Original post content for reply test')).toBeInTheDocument()
    })

    // Wait for replies to load
    await waitFor(() => {
      expect(screen.getByText('First reply content')).toBeInTheDocument()
    })

    expect(screen.getByText('Second reply content')).toBeInTheDocument()
  })

  it('displays back link to feed', async () => {
    await act(async () => {
      render(<ReplyPage />)
    })

    // Wait for async operations to complete
    await waitFor(() => {
      expect(screen.getByText('Original post content for reply test')).toBeInTheDocument()
    })

    const backLink = screen.getByRole('link', { name: /フィードに戻る/ })
    expect(backLink).toBeInTheDocument()
    expect(backLink).toHaveAttribute('href', '/dm_feed/test-user-001')
  })

  it('displays reply form section', async () => {
    await act(async () => {
      render(<ReplyPage />)
    })

    // Wait for loading to complete
    await waitFor(() => {
      expect(screen.getByText('Original post content for reply test')).toBeInTheDocument()
    })

    expect(screen.getByRole('region', { name: /返信フォーム/ })).toBeInTheDocument()
  })

  it('handles API errors gracefully', async () => {
    mockGetDmFeedPostById.mockRejectedValue(new Error('Internal Server Error'))

    await act(async () => {
      render(<ReplyPage />)
    })

    // Wait for error message
    await waitFor(() => {
      expect(screen.getByText(/Internal Server Error/i)).toBeInTheDocument()
    })
  })

  it('displays empty state when no replies', async () => {
    mockGetDmFeedReplies.mockResolvedValue([])

    await act(async () => {
      render(<ReplyPage />)
    })

    // Wait for loading to complete
    await waitFor(() => {
      expect(screen.getByText('Original post content for reply test')).toBeInTheDocument()
    })

    // Wait for empty state
    await waitFor(() => {
      expect(screen.getByText(/返信がありません/)).toBeInTheDocument()
    })
  })
})
