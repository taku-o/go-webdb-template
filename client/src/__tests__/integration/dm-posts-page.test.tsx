import { render, screen, waitFor, act } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import PostsPage from '@/app/dm-posts/page'

// Mock post data
const mockPosts = [
  {
    id: '1',
    user_id: 'user-001',
    title: 'First Post',
    content: 'This is the first post content',
    created_at: '2024-01-15T10:00:00Z',
    updated_at: '2024-01-15T10:00:00Z',
  },
  {
    id: '2',
    user_id: 'user-002',
    title: 'Second Post',
    content: 'This is the second post content',
    created_at: '2024-01-15T11:00:00Z',
    updated_at: '2024-01-15T11:00:00Z',
  },
]

// Mock user data
const mockUsers = [
  {
    id: 'user-001',
    name: 'User 1',
    email: 'user1@example.com',
    created_at: '2024-01-15T10:00:00Z',
    updated_at: '2024-01-15T10:00:00Z',
  },
  {
    id: 'user-002',
    name: 'User 2',
    email: 'user2@example.com',
    created_at: '2024-01-15T11:00:00Z',
    updated_at: '2024-01-15T11:00:00Z',
  },
]

// Mock apiClient
const mockGetDmPosts = jest.fn()
const mockGetDmUsers = jest.fn()
const mockCreateDmPost = jest.fn()
const mockDeleteDmPost = jest.fn()

jest.mock('@/lib/api', () => ({
  apiClient: {
    getDmPosts: (...args: unknown[]) => mockGetDmPosts(...args),
    getDmUsers: (...args: unknown[]) => mockGetDmUsers(...args),
    createDmPost: (...args: unknown[]) => mockCreateDmPost(...args),
    deleteDmPost: (...args: unknown[]) => mockDeleteDmPost(...args),
  },
}))

beforeEach(() => {
  jest.clearAllMocks()
  mockGetDmPosts.mockResolvedValue(mockPosts)
  mockGetDmUsers.mockResolvedValue(mockUsers)
  mockCreateDmPost.mockResolvedValue({
    id: '3',
    user_id: 'user-001',
    title: 'New Post',
    content: 'New post content',
    created_at: '2024-01-15T12:00:00Z',
    updated_at: '2024-01-15T12:00:00Z',
  })
  mockDeleteDmPost.mockResolvedValue({})
})

describe('PostsPage Integration', () => {
  it('displays page title', async () => {
    await act(async () => {
      render(<PostsPage />)
    })

    // Wait for async operations to complete
    await waitFor(() => {
      expect(screen.queryByText('読み込み中...')).not.toBeInTheDocument()
    })

    expect(screen.getByText('投稿管理')).toBeInTheDocument()
  })

  it('displays posts from API', async () => {
    await act(async () => {
      render(<PostsPage />)
    })

    // Wait for loading to complete
    await waitFor(() => {
      expect(screen.queryByText('読み込み中...')).not.toBeInTheDocument()
    }, { timeout: 5000 })

    // Wait for posts to load
    await waitFor(() => {
      expect(screen.getByText('First Post')).toBeInTheDocument()
    }, { timeout: 10000 })

    expect(screen.getByText('Second Post')).toBeInTheDocument()
  })

  it('displays loading state', async () => {
    render(<PostsPage />)

    expect(screen.getByText('読み込み中...')).toBeInTheDocument()

    // Wait for async operations to complete to avoid act() warning
    await waitFor(() => {
      expect(screen.queryByText('読み込み中...')).not.toBeInTheDocument()
    })
  })

  it('displays back link to top page', async () => {
    await act(async () => {
      render(<PostsPage />)
    })

    // Wait for async operations to complete
    await waitFor(() => {
      expect(screen.queryByText('読み込み中...')).not.toBeInTheDocument()
    })

    const backLink = screen.getByRole('link', { name: /トップページに戻る/ })
    expect(backLink).toBeInTheDocument()
    expect(backLink).toHaveAttribute('href', '/')
  })

  it('displays create post form', async () => {
    await act(async () => {
      render(<PostsPage />)
    })

    // Wait for async operations to complete
    await waitFor(() => {
      expect(screen.queryByText('読み込み中...')).not.toBeInTheDocument()
    })

    expect(screen.getByText('新規投稿作成')).toBeInTheDocument()
    expect(screen.getByLabelText('ユーザー')).toBeInTheDocument()
    expect(screen.getByLabelText('タイトル')).toBeInTheDocument()
    expect(screen.getByLabelText('本文')).toBeInTheDocument()
  })

  it('displays post count in card description', async () => {
    await act(async () => {
      render(<PostsPage />)
    })

    // Wait for posts to load
    await waitFor(() => {
      expect(screen.queryByText('読み込み中...')).not.toBeInTheDocument()
    }, { timeout: 5000 })

    await waitFor(() => {
      expect(screen.getByText(/2件の投稿が登録されています/)).toBeInTheDocument()
    }, { timeout: 5000 })
  })

  it('handles API errors gracefully', async () => {
    mockGetDmPosts.mockRejectedValue(new Error('Internal Server Error'))

    await act(async () => {
      render(<PostsPage />)
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

  it('displays empty state when no posts', async () => {
    mockGetDmPosts.mockResolvedValue([])

    await act(async () => {
      render(<PostsPage />)
    })

    // Wait for loading to complete
    await waitFor(() => {
      expect(screen.queryByText('読み込み中...')).not.toBeInTheDocument()
    }, { timeout: 5000 })

    // Empty state message
    await waitFor(() => {
      expect(screen.getByText(/投稿がありません/)).toBeInTheDocument()
    }, { timeout: 5000 })
  })
})
