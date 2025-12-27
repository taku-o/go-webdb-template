import { render, screen } from '@testing-library/react'
import HomePage from '../page'

// Mock @auth0/nextjs-auth0
jest.mock('@auth0/nextjs-auth0', () => ({
  useUser: jest.fn(),
}))

import { useUser } from '@auth0/nextjs-auth0'

const mockUseUser = useUser as jest.MockedFunction<typeof useUser>

describe('HomePage', () => {
  beforeEach(() => {
    // Default mock: not logged in
    mockUseUser.mockReturnValue({
      user: undefined,
      error: undefined,
      isLoading: false,
      checkSession: jest.fn(),
    })
  })

  afterEach(() => {
    jest.clearAllMocks()
  })

  it('renders the main heading', () => {
    render(<HomePage />)

    const heading = screen.getByRole('heading', { level: 1 })
    expect(heading).toBeInTheDocument()
    expect(heading).toHaveTextContent(/Go DB Project Sample/)
  })

  it('renders navigation links', () => {
    render(<HomePage />)

    // Check for link texts (using getAllByText since heading and link may have same text)
    const userLinks = screen.getAllByText(/ユーザー管理/)
    expect(userLinks.length).toBeGreaterThan(0)

    const postLinks = screen.getAllByText(/投稿管理/)
    expect(postLinks.length).toBeGreaterThan(0)

    const userPostLinks = screen.getAllByText(/ユーザーと投稿/)
    expect(userPostLinks.length).toBeGreaterThan(0)
  })

  it('displays tech stack', () => {
    render(<HomePage />)

    // Check for tech stack descriptions
    expect(screen.getByText(/Go \(Sharding対応\)/)).toBeInTheDocument()
    expect(screen.getByText(/Next\.js 14 \(App Router\)/)).toBeInTheDocument()
    expect(screen.getByText(/TypeScript/)).toBeInTheDocument()
  })

  it('has links with correct hrefs', () => {
    render(<HomePage />)

    const links = screen.getAllByRole('link')
    const usersLink = links.find(link => link.getAttribute('href') === '/users')
    expect(usersLink).toBeDefined()

    const postsLink = links.find(link => link.getAttribute('href') === '/posts')
    expect(postsLink).toBeDefined()

    const userPostsLink = links.find(link => link.getAttribute('href') === '/user-posts')
    expect(userPostsLink).toBeDefined()
  })

  // Auth0 login/logout tests
  describe('Authentication UI', () => {
    it('shows loading state when isLoading is true', () => {
      mockUseUser.mockReturnValue({
        user: undefined,
        error: undefined,
        isLoading: true,
        checkSession: jest.fn(),
      })

      render(<HomePage />)

      expect(screen.getByText('Loading...')).toBeInTheDocument()
    })

    it('shows login button when user is not logged in', () => {
      mockUseUser.mockReturnValue({
        user: undefined,
        error: undefined,
        isLoading: false,
        checkSession: jest.fn(),
      })

      render(<HomePage />)

      expect(screen.getByText('ログインしていません')).toBeInTheDocument()
      const loginLink = screen.getByRole('link', { name: 'ログイン' })
      expect(loginLink).toBeInTheDocument()
      expect(loginLink).toHaveAttribute('href', '/auth/login')
    })

    it('shows logout button and user info when user is logged in', () => {
      mockUseUser.mockReturnValue({
        user: {
          name: 'Test User',
          email: 'test@example.com',
          sub: 'auth0|12345',
        },
        error: undefined,
        isLoading: false,
        checkSession: jest.fn(),
      })

      render(<HomePage />)

      expect(screen.getByText('ログイン中: Test User')).toBeInTheDocument()
      expect(screen.getByText('test@example.com')).toBeInTheDocument()
      const logoutLink = screen.getByRole('link', { name: 'ログアウト' })
      expect(logoutLink).toBeInTheDocument()
      expect(logoutLink).toHaveAttribute('href', '/auth/logout')
    })

    it('shows error message when authentication error occurs', () => {
      mockUseUser.mockReturnValue({
        user: undefined,
        error: new Error('Authentication failed'),
        isLoading: false,
        checkSession: jest.fn(),
      })

      render(<HomePage />)

      expect(screen.getByText('認証エラーが発生しました: Authentication failed')).toBeInTheDocument()
      const retryLink = screen.getByRole('link', { name: '再度ログイン' })
      expect(retryLink).toBeInTheDocument()
      expect(retryLink).toHaveAttribute('href', '/auth/login')
    })
  })
})
