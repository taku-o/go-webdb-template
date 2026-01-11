import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import TodayApiButton from '../TodayApiButton'
import { apiClient } from '@/lib/api'

// Mock @auth0/nextjs-auth0
jest.mock('@auth0/nextjs-auth0', () => ({
  useUser: jest.fn(),
}))

// Mock apiClient
jest.mock('@/lib/api', () => ({
  apiClient: {
    getToday: jest.fn(),
  },
}))

import { useUser } from '@auth0/nextjs-auth0'

const mockUseUser = useUser as jest.MockedFunction<typeof useUser>
const mockGetToday = apiClient.getToday as jest.MockedFunction<typeof apiClient.getToday>

describe('TodayApiButton', () => {
  beforeEach(() => {
    jest.clearAllMocks()
    // Default mock: not logged in
    mockUseUser.mockReturnValue({
      user: undefined,
      error: undefined,
      isLoading: false,
      checkSession: jest.fn(),
    })
  })

  it('renders the component with title and description', () => {
    render(<TodayApiButton />)

    expect(screen.getByText('Today API (Private Endpoint)')).toBeInTheDocument()
    expect(screen.getByText(/Auth0ログイン時のみアクセス可能なプライベートAPIをテストします/)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Get Today' })).toBeInTheDocument()
  })

  it('calls apiClient.getToday with undefined when user is not logged in', async () => {
    const user = userEvent.setup()
    mockGetToday.mockResolvedValueOnce({ date: '2024-01-15' })

    render(<TodayApiButton />)

    const button = screen.getByRole('button', { name: 'Get Today' })
    await user.click(button)

    await waitFor(() => {
      expect(mockGetToday).toHaveBeenCalledWith(undefined)
    }, { timeout: 3000 })
  })

  it('calls apiClient.getToday with auth0user when user is logged in', async () => {
    const user = userEvent.setup()
    const mockAuth0User = {
      sub: 'auth0|123',
      email: 'test@example.com',
      name: 'Test User',
    }

    mockUseUser.mockReturnValue({
      user: mockAuth0User,
      error: undefined,
      isLoading: false,
      checkSession: jest.fn(),
    })

    mockGetToday.mockResolvedValueOnce({ date: '2024-01-15' })

    render(<TodayApiButton />)

    const button = screen.getByRole('button', { name: 'Get Today' })
    await user.click(button)

    await waitFor(() => {
      expect(mockGetToday).toHaveBeenCalledWith(mockAuth0User)
    }, { timeout: 3000 })
  })

  it('displays the date when API call succeeds', async () => {
    const user = userEvent.setup()
    mockGetToday.mockResolvedValueOnce({ date: '2024-01-15' })

    render(<TodayApiButton />)

    const button = screen.getByRole('button', { name: 'Get Today' })
    await user.click(button)

    await waitFor(() => {
      expect(screen.getByText('Today: 2024-01-15')).toBeInTheDocument()
    })
  })

  it('displays error message when API call fails', async () => {
    const user = userEvent.setup()
    const errorMessage = 'Failed to get today'
    mockGetToday.mockRejectedValueOnce(new Error(errorMessage))

    render(<TodayApiButton />)

    const button = screen.getByRole('button', { name: 'Get Today' })
    await user.click(button)

    await waitFor(() => {
      expect(screen.getByText(`Error: ${errorMessage}`)).toBeInTheDocument()
    })
  })

  it('displays generic error message when error is not an Error instance', async () => {
    const user = userEvent.setup()
    mockGetToday.mockRejectedValueOnce('Unknown error')

    render(<TodayApiButton />)

    const button = screen.getByRole('button', { name: 'Get Today' })
    await user.click(button)

    await waitFor(() => {
      expect(screen.getByText('Error: An error occurred')).toBeInTheDocument()
    })
  })

  it('shows loading state while API call is in progress', async () => {
    const user = userEvent.setup()
    let resolvePromise: (value: { date: string }) => void
    const promise = new Promise<{ date: string }>(resolve => {
      resolvePromise = resolve
    })
    mockGetToday.mockReturnValueOnce(promise)

    render(<TodayApiButton />)

    const button = screen.getByRole('button', { name: 'Get Today' })
    await user.click(button)

    // Button should show loading state
    expect(screen.getByText('Loading...')).toBeInTheDocument()
    expect(button).toBeDisabled()

    // Resolve the promise
    resolvePromise!({ date: '2024-01-15' })

    await waitFor(() => {
      expect(screen.getByText('Today: 2024-01-15')).toBeInTheDocument()
    })
    expect(screen.queryByText('Loading...')).not.toBeInTheDocument()
  })

  it('disables button when isLoading is true', () => {
    mockUseUser.mockReturnValue({
      user: undefined,
      error: undefined,
      isLoading: true,
      checkSession: jest.fn(),
    })

    render(<TodayApiButton />)

    const button = screen.getByRole('button', { name: 'Get Today' })
    expect(button).toBeDisabled()
  })

  it('clears previous date and error when button is clicked again', async () => {
    const user = userEvent.setup()
    mockGetToday
      .mockResolvedValueOnce({ date: '2024-01-15' })
      .mockResolvedValueOnce({ date: '2024-01-16' })

    render(<TodayApiButton />)

    const button = screen.getByRole('button', { name: 'Get Today' })

    // First click
    await user.click(button)
    await waitFor(() => {
      expect(screen.getByText('Today: 2024-01-15')).toBeInTheDocument()
    })

    // Second click
    await user.click(button)
    await waitFor(() => {
      expect(screen.queryByText('Today: 2024-01-15')).not.toBeInTheDocument()
      expect(screen.getByText('Today: 2024-01-16')).toBeInTheDocument()
    })
  })
})
