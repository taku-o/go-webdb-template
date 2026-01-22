import { render, screen } from '@testing-library/react'
import Home from '@/app/page'
import { auth } from '@/auth'
import type { Session } from 'next-auth'

// Mock auth function
const mockAuth = auth as unknown as jest.MockedFunction<() => Promise<Session | null>>

// Mock AuthButtons component to avoid Server Actions warning
jest.mock('@/components/auth/auth-buttons', () => ({
  AuthButtons: ({ user }: { user: { name?: string | null; email?: string | null } | null }) => {
    if (user) {
      return (
        <div>
          <p>ログイン中: {user.name}</p>
          {user.email && <p>{user.email}</p>}
          <button aria-label="ログアウト">ログアウト</button>
        </div>
      )
    }
    return (
      <div>
        <p>ログインしていません</p>
        <button aria-label="ログイン">ログイン</button>
      </div>
    )
  },
}))

describe('Home Page Integration', () => {
  beforeEach(() => {
    jest.clearAllMocks()
  })

  it('displays main title', async () => {
    mockAuth.mockResolvedValueOnce(null)

    const Component = await Home()
    render(Component)

    expect(screen.getByText('Go DB Project Sample')).toBeInTheDocument()
  })

  it('displays project description', async () => {
    mockAuth.mockResolvedValueOnce(null)

    const Component = await Home()
    render(Component)

    expect(screen.getByText(/Go \+ Next\.js \+ Sharding対応のサンプルプロジェクト/)).toBeInTheDocument()
  })

  it('displays feature cards', async () => {
    mockAuth.mockResolvedValueOnce(null)

    const Component = await Home()
    render(Component)

    // Check feature titles
    expect(screen.getByText('ユーザー管理')).toBeInTheDocument()
    expect(screen.getByText('投稿管理')).toBeInTheDocument()
    expect(screen.getByText('ユーザーと投稿')).toBeInTheDocument()
    expect(screen.getByText('動画アップロード')).toBeInTheDocument()
    expect(screen.getByText('メール送信')).toBeInTheDocument()
    expect(screen.getByText('ジョブキュー')).toBeInTheDocument()
    expect(screen.getByText('フィード')).toBeInTheDocument()
    expect(screen.getByText('動画プレイヤー')).toBeInTheDocument()
  })

  it('displays login button when user is not authenticated', async () => {
    mockAuth.mockResolvedValueOnce(null)

    const Component = await Home()
    render(Component)

    expect(screen.getByText('ログインしていません')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'ログイン' })).toBeInTheDocument()
  })

  it('displays user info and logout button when user is authenticated', async () => {
    mockAuth.mockResolvedValueOnce({
      user: {
        name: 'Test User',
        email: 'test@example.com',
      },
      expires: '2024-12-31T23:59:59.999Z',
    })

    const Component = await Home()
    render(Component)

    expect(screen.getByText(/ログイン中: Test User/)).toBeInTheDocument()
    expect(screen.getByText('test@example.com')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'ログアウト' })).toBeInTheDocument()
  })

  it('displays TodayApiButton component', async () => {
    mockAuth.mockResolvedValueOnce(null)

    const Component = await Home()
    render(Component)

    expect(screen.getByText('Today API (Private Endpoint)')).toBeInTheDocument()
  })

  it('has correct navigation links', async () => {
    mockAuth.mockResolvedValueOnce(null)

    const Component = await Home()
    render(Component)

    // Check that feature cards have correct links
    const userManagementLink = screen.getByRole('link', { name: /ユーザー管理/ })
    expect(userManagementLink).toHaveAttribute('href', '/dm-users')

    const postManagementLink = screen.getByRole('link', { name: /投稿管理/ })
    expect(postManagementLink).toHaveAttribute('href', '/dm-posts')

    const userPostsLink = screen.getByRole('link', { name: /ユーザーと投稿/ })
    expect(userPostsLink).toHaveAttribute('href', '/dm-user-posts')

    const videoUploadLink = screen.getByRole('link', { name: /動画アップロード/ })
    expect(videoUploadLink).toHaveAttribute('href', '/dm_movie/upload')

    const emailSendLink = screen.getByRole('link', { name: /メール送信/ })
    expect(emailSendLink).toHaveAttribute('href', '/dm_email/send')

    const jobQueueLink = screen.getByRole('link', { name: /ジョブキュー/ })
    expect(jobQueueLink).toHaveAttribute('href', '/dm-jobqueue')

    const feedLink = screen.getByRole('link', { name: /フィード/ })
    expect(feedLink).toHaveAttribute('href', '/dm_feed')

    const videoPlayerLink = screen.getByRole('link', { name: /動画プレイヤー/ })
    expect(videoPlayerLink).toHaveAttribute('href', '/dm_videoplayer')
  })
})
