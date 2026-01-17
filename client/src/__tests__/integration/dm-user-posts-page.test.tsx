import { render, screen, waitFor, act } from '@testing-library/react'
import UserPostsPage from '@/app/dm-user-posts/page'

// Mock user-posts data
const mockUserPosts = [
  {
    user_id: 'user-001',
    user_name: 'User 1',
    user_email: 'user1@example.com',
    post_id: 'post-001',
    post_title: 'First Post Title',
    post_content: 'This is the first post content',
    created_at: '2024-01-15T10:00:00Z',
  },
  {
    user_id: 'user-002',
    user_name: 'User 2',
    user_email: 'user2@example.com',
    post_id: 'post-002',
    post_title: 'Second Post Title',
    post_content: 'This is the second post content',
    created_at: '2024-01-15T11:00:00Z',
  },
]

// Mock apiClient
const mockGetDmUserPosts = jest.fn()

jest.mock('@/lib/api', () => ({
  apiClient: {
    getDmUserPosts: (...args: unknown[]) => mockGetDmUserPosts(...args),
  },
}))

beforeEach(() => {
  jest.clearAllMocks()
  mockGetDmUserPosts.mockResolvedValue(mockUserPosts)
})

describe('UserPostsPage Integration', () => {
  it('displays page title', async () => {
    await act(async () => {
      render(<UserPostsPage />)
    })

    // Wait for async operations to complete
    await waitFor(() => {
      expect(screen.queryByText('読み込み中...')).not.toBeInTheDocument()
    })

    expect(screen.getByText('ユーザーと投稿（JOIN）')).toBeInTheDocument()
  })

  it('displays cross-shard query explanation', async () => {
    await act(async () => {
      render(<UserPostsPage />)
    })

    // Wait for async operations to complete
    await waitFor(() => {
      expect(screen.queryByText('読み込み中...')).not.toBeInTheDocument()
    })

    expect(screen.getByText('クロスシャードクエリ')).toBeInTheDocument()
    expect(screen.getByText(/複数のShardからユーザーと投稿をJOINして取得/)).toBeInTheDocument()
  })

  it('displays user posts from API', async () => {
    await act(async () => {
      render(<UserPostsPage />)
    })

    // Wait for loading to complete
    await waitFor(() => {
      expect(screen.queryByText('読み込み中...')).not.toBeInTheDocument()
    }, { timeout: 5000 })

    // Wait for user posts to load
    await waitFor(() => {
      expect(screen.getByText('First Post Title')).toBeInTheDocument()
    }, { timeout: 10000 })

    expect(screen.getByText('Second Post Title')).toBeInTheDocument()
    expect(screen.getByText('User 1')).toBeInTheDocument()
    expect(screen.getByText('User 2')).toBeInTheDocument()
  })

  it('displays loading state', async () => {
    render(<UserPostsPage />)

    expect(screen.getByText('読み込み中...')).toBeInTheDocument()

    // Wait for async operations to complete to avoid act() warning
    await waitFor(() => {
      expect(screen.queryByText('読み込み中...')).not.toBeInTheDocument()
    })
  })

  it('displays back link to top page', async () => {
    await act(async () => {
      render(<UserPostsPage />)
    })

    // Wait for async operations to complete
    await waitFor(() => {
      expect(screen.queryByText('読み込み中...')).not.toBeInTheDocument()
    })

    const backLink = screen.getByRole('link', { name: /トップページに戻る/ })
    expect(backLink).toBeInTheDocument()
    expect(backLink).toHaveAttribute('href', '/')
  })

  it('displays sharding info section', async () => {
    await act(async () => {
      render(<UserPostsPage />)
    })

    // Wait for async operations to complete
    await waitFor(() => {
      expect(screen.queryByText('読み込み中...')).not.toBeInTheDocument()
    })

    expect(screen.getByText('Sharding情報')).toBeInTheDocument()
    expect(screen.getByText(/Hash-based sharding/)).toBeInTheDocument()
  })

  it('handles API errors gracefully', async () => {
    mockGetDmUserPosts.mockRejectedValue(new Error('Internal Server Error'))

    await act(async () => {
      render(<UserPostsPage />)
    })

    // Wait for loading to complete
    await waitFor(() => {
      expect(screen.queryByText('読み込み中...')).not.toBeInTheDocument()
    }, { timeout: 5000 })

    // Error message should be displayed
    await waitFor(() => {
      expect(screen.getByText(/Internal Server Error/i)).toBeInTheDocument()
    }, { timeout: 10000 })
  })

  it('displays empty state when no user posts', async () => {
    mockGetDmUserPosts.mockResolvedValue([])

    await act(async () => {
      render(<UserPostsPage />)
    })

    // Wait for loading to complete
    await waitFor(() => {
      expect(screen.queryByText('読み込み中...')).not.toBeInTheDocument()
    }, { timeout: 5000 })

    // Empty state message
    await waitFor(() => {
      expect(screen.getByText(/表示する投稿がありません/)).toBeInTheDocument()
    }, { timeout: 5000 })

    // Links to create pages
    expect(screen.getByRole('link', { name: /ユーザー作成ページへ/ })).toHaveAttribute('href', '/dm-users')
    expect(screen.getByRole('link', { name: /投稿作成ページへ/ })).toHaveAttribute('href', '/dm-posts')
  })
})
