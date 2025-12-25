import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { http, HttpResponse } from 'msw'
import { setupServer } from 'msw/node'
import UsersPage from '@/app/users/page'

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
  http.get('http://localhost:8080/api/users', () => {
    return HttpResponse.json(mockUsers)
  }),

  http.post('http://localhost:8080/api/users', async ({ request }) => {
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

  http.delete('http://localhost:8080/api/users/:id', () => {
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

    // Wait for users to load
    await waitFor(() => {
      expect(screen.getByText('User 1')).toBeInTheDocument()
    })

    expect(screen.getByText('user1@example.com')).toBeInTheDocument()
    expect(screen.getByText('User 2')).toBeInTheDocument()
    expect(screen.getByText('user2@example.com')).toBeInTheDocument()
  })

  it('creates a new user', async () => {
    const user = userEvent.setup()
    render(<UsersPage />)

    // Wait for initial load
    await waitFor(() => {
      expect(screen.getByText('User 1')).toBeInTheDocument()
    })

    // Find inputs by role
    const inputs = screen.getAllByRole('textbox')
    const nameInput = inputs[0] // First text input is name
    const emailInput = inputs[1] // Second text input is email
    const submitButton = screen.getByRole('button', { name: '作成' })

    await user.type(nameInput, 'New User')
    await user.type(emailInput, 'new@example.com')
    await user.click(submitButton)

    // Wait for new user to appear
    await waitFor(() => {
      expect(screen.getByText('New User')).toBeInTheDocument()
    })
    expect(screen.getByText('new@example.com')).toBeInTheDocument()

    // Form should be cleared
    expect(nameInput).toHaveValue('')
    expect(emailInput).toHaveValue('')
  })

  it('handles API errors gracefully', async () => {
    server.use(
      http.get('http://localhost:8080/api/users', () => {
        return new HttpResponse(null, { status: 500, statusText: 'Internal Server Error' })
      })
    )

    render(<UsersPage />)

    // Component displays err.message which comes from the fetch error
    await waitFor(() => {
      expect(screen.getByText(/Internal Server Error/i)).toBeInTheDocument()
    })
  })

  it('shows loading state', () => {
    render(<UsersPage />)

    expect(screen.getByText('読み込み中...')).toBeInTheDocument()
  })
})
