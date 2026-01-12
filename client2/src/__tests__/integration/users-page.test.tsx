import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { http, HttpResponse } from 'msw'
import { setupServer } from 'msw/node'
import UsersPage from '@/app/dm-users/page'

// Mock user data that can be modified during tests
let mockUsers = [
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

const server = setupServer(
  // Mock NextAuth token endpoint (relative path)
  http.get('*/api/auth/token', () => {
    return new HttpResponse(null, { status: 401 })
  }),

  http.get('http://localhost:8080/api/dm-users', () => {
    return HttpResponse.json(mockUsers)
  }),

  http.post('http://localhost:8080/api/dm-users', async ({ request }) => {
    const body = (await request.json()) as { name: string; email: string }
    const newUser = {
      id: '3',
      name: body.name,
      email: body.email,
      created_at: '2024-01-15T12:00:00Z',
      updated_at: '2024-01-15T12:00:00Z',
    }
    // Add to mock users so next GET returns it
    mockUsers = [...mockUsers, newUser]
    return HttpResponse.json(newUser, { status: 201 })
  }),

  http.delete('http://localhost:8080/api/dm-users/:id', () => {
    return new HttpResponse(null, { status: 204 })
  })
)

beforeAll(() => server.listen())
afterEach(() => {
  server.resetHandlers()
  // Reset mock users
  mockUsers = [
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
})
afterAll(() => server.close())

describe('UsersPage Integration', () => {
  it('displays users from API', async () => {
    render(<UsersPage />)

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
    const user = userEvent.setup()
    render(<UsersPage />)

    // Wait for loading to complete
    await waitFor(() => {
      expect(screen.queryByText('読み込み中...')).not.toBeInTheDocument()
    }, { timeout: 5000 })

    // Wait for initial load
    await waitFor(() => {
      expect(screen.getByText('User 1')).toBeInTheDocument()
    }, { timeout: 10000 })

    // Find inputs by id (client2 uses id attributes)
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
    server.use(
      http.get('http://localhost:8080/api/dm-users', () => {
        return new HttpResponse(null, { status: 500, statusText: 'Internal Server Error' })
      })
    )

    render(<UsersPage />)

    // Wait for loading to complete
    await waitFor(() => {
      expect(screen.queryByText('読み込み中...')).not.toBeInTheDocument()
    }, { timeout: 5000 })

    // Component displays error message (ErrorAlert component)
    await waitFor(() => {
      expect(screen.getByText(/Internal Server Error/i)).toBeInTheDocument()
    }, { timeout: 10000 })
  })

  it('shows loading state', () => {
    render(<UsersPage />)

    expect(screen.getByText('読み込み中...')).toBeInTheDocument()
  })
})
