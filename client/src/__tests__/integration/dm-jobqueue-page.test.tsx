import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { http, HttpResponse } from 'msw'
import { setupServer } from 'msw/node'
import DmJobqueuePage from '@/app/dm-jobqueue/page'

const server = setupServer(
  http.post('http://localhost:8080/api/dm-jobqueue/register', async () => {
    return HttpResponse.json({
      job_id: 'test-job-id-123',
      status: 'registered',
    }, { status: 201 })
  })
)

beforeAll(() => server.listen())
afterEach(() => server.resetHandlers())
afterAll(() => server.close())

describe('DmJobqueuePage Integration', () => {
  it('displays job registration form', async () => {
    render(<DmJobqueuePage />)

    expect(screen.getByRole('heading', { name: /ジョブキュー/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /ジョブを登録/i })).toBeInTheDocument()
  })

  it('registers a job successfully', async () => {
    const user = userEvent.setup()
    render(<DmJobqueuePage />)

    const submitButton = screen.getByRole('button', { name: /ジョブを登録/i })
    await user.click(submitButton)

    // Wait for success message
    await waitFor(() => {
      expect(screen.getByText(/登録されました/i)).toBeInTheDocument()
    })
  })

  it('handles API errors gracefully', async () => {
    server.use(
      http.post('http://localhost:8080/api/dm-jobqueue/register', () => {
        return new HttpResponse('Service Unavailable', { status: 503 })
      })
    )

    const user = userEvent.setup()
    render(<DmJobqueuePage />)

    const submitButton = screen.getByRole('button', { name: /ジョブを登録/i })
    await user.click(submitButton)

    await waitFor(() => {
      expect(screen.getByText(/Service Unavailable/i)).toBeInTheDocument()
    })
  })

  it('allows custom message input', async () => {
    const user = userEvent.setup()
    render(<DmJobqueuePage />)

    const messageInput = screen.getByLabelText(/メッセージ/i)
    await user.type(messageInput, 'Custom test message')

    expect(messageInput).toHaveValue('Custom test message')
  })
})
