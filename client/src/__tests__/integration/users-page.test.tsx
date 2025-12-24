import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { rest } from 'msw'
import { setupServer } from 'msw/node'
import UsersPage from '@/app/users/page'

const server = setupServer(
  rest.get('http://localhost:8080/api/users', (req, res, ctx) => {
    return res(
      ctx.json([
        {
          id: 1,
          name: 'User 1',
          email: 'user1@example.com',
          created_at: '2024-01-15T10:00:00Z',
          updated_at: '2024-01-15T10:00:00Z',
        },
        {
          id: 2,
          name: 'User 2',
          email: 'user2@example.com',
          created_at: '2024-01-15T11:00:00Z',
          updated_at: '2024-01-15T11:00:00Z',
        },
      ])
    )
  }),

  rest.post('http://localhost:8080/api/users', async (req, res, ctx) => {
    const body = await req.json()
    return res(
      ctx.status(201),
      ctx.json({
        id: 3,
        name: body.name,
        email: body.email,
        created_at: '2024-01-15T12:00:00Z',
        updated_at: '2024-01-15T12:00:00Z',
      })
    )
  }),

  rest.delete('http://localhost:8080/api/users/:id', (req, res, ctx) => {
    return res(ctx.status(204))
  })
)

beforeAll(() => server.listen())
afterEach(() => server.resetHandlers())
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

    // Fill in form
    const nameInput = screen.getByLabelText('名前')
    const emailInput = screen.getByLabelText('メールアドレス')
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
      rest.get('http://localhost:8080/api/users', (req, res, ctx) => {
        return res(ctx.status(500))
      })
    )

    render(<UsersPage />)

    await waitFor(() => {
      expect(screen.getByText(/Failed to load users/i)).toBeInTheDocument()
    })
  })

  it('shows loading state', () => {
    render(<UsersPage />)

    expect(screen.getByText('読み込み中...')).toBeInTheDocument()
  })
})
