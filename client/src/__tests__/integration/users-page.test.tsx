import { render, screen, waitFor, act } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import UsersPage from '@/app/dm-users/page'

// Mock user data
const mockUsers = [
  {
    id: '1',
    name: 'User 1',
    email: 'user1@example.com',
    created_at: '2024-01-15T10:00:00Z',
    updated_at: '2024-01-15T10:00:00Z',
  },
  {
    id: '2',
    name: 'User 2',
    email: 'user2@example.com',
    created_at: '2024-01-15T11:00:00Z',
    updated_at: '2024-01-15T11:00:00Z',
  },
]

// Mock apiClient
const mockGetDmUsers = jest.fn()
const mockCreateDmUser = jest.fn()
const mockDeleteDmUser = jest.fn()

jest.mock('@/lib/api', () => ({
  apiClient: {
    getDmUsers: (...args: unknown[]) => mockGetDmUsers(...args),
    createDmUser: (...args: unknown[]) => mockCreateDmUser(...args),
    deleteDmUser: (...args: unknown[]) => mockDeleteDmUser(...args),
  },
}))

beforeEach(() => {
  jest.clearAllMocks()
  mockGetDmUsers.mockResolvedValue(mockUsers)
  mockCreateDmUser.mockResolvedValue({
    id: '3',
    name: 'New User',
    email: 'new@example.com',
    created_at: '2024-01-15T12:00:00Z',
    updated_at: '2024-01-15T12:00:00Z',
  })
  mockDeleteDmUser.mockResolvedValue({})
})

describe('UsersPage Integration', () => {
  it('displays users from API', async () => {
    await act(async () => {
      render(<UsersPage />)
    })

    // Wait for loading to complete
    await waitFor(() => {
      expect(screen.queryByText('読み込み中...')).not.toBeInTheDocument()
    }, { timeout: 5000 })

    // Wait for users to load (check for table rows)
    await waitFor(() => {
      expect(screen.getByText('User 1')).toBeInTheDocument()
    }, { timeout: 10000 })

    // Check users are displayed (use getAllByText for multiple matches)
    const user1Emails = screen.getAllByText('user1@example.com')
    expect(user1Emails.length).toBeGreaterThan(0)

    expect(screen.getByText('User 2')).toBeInTheDocument()
    const user2Emails = screen.getAllByText('user2@example.com')
    expect(user2Emails.length).toBeGreaterThan(0)
  })

  it('creates a new user', async () => {
    const newUser = {
      id: '3',
      name: 'New User',
      email: 'new@example.com',
      created_at: '2024-01-15T12:00:00Z',
      updated_at: '2024-01-15T12:00:00Z',
    }

    // Setup mock to return updated users list after creation
    let usersData = [...mockUsers]
    mockGetDmUsers.mockImplementation(() => Promise.resolve(usersData))
    mockCreateDmUser.mockImplementation(() => {
      usersData = [...usersData, newUser]
      return Promise.resolve(newUser)
    })

    const user = userEvent.setup()
    await act(async () => {
      render(<UsersPage />)
    })

    // Wait for loading to complete
    await waitFor(() => {
      expect(screen.queryByText('読み込み中...')).not.toBeInTheDocument()
    }, { timeout: 5000 })

    // Wait for initial load
    await waitFor(() => {
      expect(screen.getByText('User 1')).toBeInTheDocument()
    }, { timeout: 10000 })

    // Find inputs by id (client uses id attributes)
    const nameInput = screen.getByLabelText('名前')
    const emailInput = screen.getByLabelText('メールアドレス')
    const submitButton = screen.getByRole('button', { name: /作成/ })

    await user.type(nameInput, 'New User')
    await user.type(emailInput, 'new@example.com')
    await user.click(submitButton)

    // Wait for new user to appear (after form submission and reload)
    await waitFor(() => {
      expect(screen.getByText('New User')).toBeInTheDocument()
    }, { timeout: 10000 })
    // Check email appears in table (more specific selector)
    await waitFor(() => {
      const emailCells = screen.getAllByText('new@example.com')
      expect(emailCells.length).toBeGreaterThan(0)
    }, { timeout: 5000 })

    // Form should be cleared (wait a bit for state update)
    await waitFor(() => {
      expect(nameInput).toHaveValue('')
      expect(emailInput).toHaveValue('')
    }, { timeout: 5000 })
  })

  it('handles API errors gracefully', async () => {
    mockGetDmUsers.mockRejectedValue(new Error('Internal Server Error'))

    await act(async () => {
      render(<UsersPage />)
    })

    // Wait for loading to complete
    await waitFor(() => {
      expect(screen.queryByText('読み込み中...')).not.toBeInTheDocument()
    }, { timeout: 5000 })

    // Component displays error message (ErrorAlert component)
    await waitFor(() => {
      expect(screen.getByText(/Internal Server Error/i)).toBeInTheDocument()
    }, { timeout: 10000 })
  })

  it('shows loading state', async () => {
    render(<UsersPage />)

    expect(screen.getByText('読み込み中...')).toBeInTheDocument()

    // Wait for async operations to complete to avoid act() warning
    await waitFor(() => {
      expect(screen.queryByText('読み込み中...')).not.toBeInTheDocument()
    })
  })
})
