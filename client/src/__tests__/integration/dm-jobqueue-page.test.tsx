import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import DmJobqueuePage from '@/app/dm-jobqueue/page'

// Mock apiClient
const mockRegisterJob = jest.fn()

jest.mock('@/lib/api', () => ({
  apiClient: {
    registerJob: (...args: unknown[]) => mockRegisterJob(...args),
  },
}))

beforeEach(() => {
  jest.clearAllMocks()
  mockRegisterJob.mockResolvedValue({
    job_id: 'test-job-id-123',
    status: 'registered',
  })
})

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

    // Wait for success message (displayed in Alert with role="status")
    await waitFor(() => {
      expect(screen.getByRole('status')).toBeInTheDocument()
    }, { timeout: 10000 })
  })

  it('handles API errors gracefully', async () => {
    mockRegisterJob.mockRejectedValue(new Error('Service Unavailable'))

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
