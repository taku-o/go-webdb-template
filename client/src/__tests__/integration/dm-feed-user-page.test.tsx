import { render, screen, waitFor, act } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import FeedPage from '@/app/dm_feed/[userId]/page'

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
  useParams: () => ({ userId: 'test-user-001' }),
}))

// Mock feed posts data
const mockFeedPosts = [
  {
    id: 'post-001',
    userId: 'test-user-001',
    userName: 'テストユーザー1',
    userHandle: '@test_user_1',
    content: 'First feed post content',
    createdAt: '2024-01-15T10:00:00Z',
    liked: false,
    likeCount: 5,
  },
  {
    id: 'post-002',
    userId: 'test-user-001',
    userName: 'テストユーザー1',
    userHandle: '@test_user_1',
    content: 'Second feed post content',
    createdAt: '2024-01-15T11:00:00Z',
    liked: true,
    likeCount: 10,
  },
]

// Mock apiClient
const mockGetDmFeedPosts = jest.fn()
const mockCreateDmFeedPost = jest.fn()
const mockToggleLikeDmPost = jest.fn()

jest.mock('@/lib/api', () => ({
  apiClient: {
    getDmFeedPosts: (...args: unknown[]) => mockGetDmFeedPosts(...args),
    createDmFeedPost: (...args: unknown[]) => mockCreateDmFeedPost(...args),
    toggleLikeDmPost: (...args: unknown[]) => mockToggleLikeDmPost(...args),
  },
}))

beforeEach(() => {
  jest.clearAllMocks()
  mockGetDmFeedPosts.mockResolvedValue(mockFeedPosts)
  mockCreateDmFeedPost.mockResolvedValue({
    id: 'post-003',
    userId: 'test-user-001',
    userName: 'あなた',
    userHandle: '@you',
    content: 'New post content',
    createdAt: new Date().toISOString(),
    liked: false,
    likeCount: 0,
  })
  mockToggleLikeDmPost.mockResolvedValue({ liked: true, likeCount: 6 })
})

describe('FeedPage Integration', () => {
  it('displays page title', async () => {
    await act(async () => {
      render(<FeedPage />)
    })

    // Wait for async operations to complete
    await waitFor(() => {
      expect(screen.getByText('First feed post content')).toBeInTheDocument()
    })

    expect(screen.getByText('フィード')).toBeInTheDocument()
  })

  it('displays feed posts from API', async () => {
    await act(async () => {
      render(<FeedPage />)
    })

    // Wait for posts to load
    await waitFor(() => {
      expect(screen.getByText('First feed post content')).toBeInTheDocument()
    })

    expect(screen.getByText('Second feed post content')).toBeInTheDocument()
  })

  it('displays back link to top page', async () => {
    await act(async () => {
      render(<FeedPage />)
    })

    // Wait for async operations to complete
    await waitFor(() => {
      expect(screen.getByText('First feed post content')).toBeInTheDocument()
    })

    const backLink = screen.getByRole('link', { name: /トップページに戻る/ })
    expect(backLink).toBeInTheDocument()
    expect(backLink).toHaveAttribute('href', '/')
  })

  it('displays new post form section', async () => {
    await act(async () => {
      render(<FeedPage />)
    })

    // Wait for async operations to complete
    await waitFor(() => {
      expect(screen.getByText('First feed post content')).toBeInTheDocument()
    })

    expect(screen.getByRole('region', { name: /新規投稿フォーム/ })).toBeInTheDocument()
  })

  it('handles API errors gracefully', async () => {
    mockGetDmFeedPosts.mockRejectedValue(new Error('Internal Server Error'))

    await act(async () => {
      render(<FeedPage />)
    })

    // Wait for error message
    await waitFor(() => {
      expect(screen.getByText(/Internal Server Error/i)).toBeInTheDocument()
    })
  })

  it('displays empty state when no posts', async () => {
    mockGetDmFeedPosts.mockResolvedValue([])

    await act(async () => {
      render(<FeedPage />)
    })

    // Wait for empty state
    await waitFor(() => {
      expect(screen.getByText(/投稿がありません/)).toBeInTheDocument()
    })
  })
})
