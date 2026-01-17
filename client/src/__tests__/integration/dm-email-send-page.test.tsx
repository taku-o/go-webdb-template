import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import SendEmailPage from '@/app/dm_email/send/page'

// Mock apiClient
const mockSendEmail = jest.fn()

jest.mock('@/lib/api', () => ({
  apiClient: {
    sendEmail: (...args: unknown[]) => mockSendEmail(...args),
  },
}))

beforeEach(() => {
  jest.clearAllMocks()
  mockSendEmail.mockResolvedValue({ success: true, message: 'メールを送信しました' })
})

describe('SendEmailPage Integration', () => {
  it('displays page title', async () => {
    render(<SendEmailPage />)

    expect(screen.getByText('メール送信')).toBeInTheDocument()
  })

  it('displays email send form', async () => {
    render(<SendEmailPage />)

    expect(screen.getByText('ウェルカムメール送信')).toBeInTheDocument()
    expect(screen.getByLabelText('送信先メールアドレス')).toBeInTheDocument()
    expect(screen.getByLabelText('お名前')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /メールを送信/ })).toBeInTheDocument()
  })

  it('displays back link to top page', async () => {
    render(<SendEmailPage />)

    const backLink = screen.getByRole('link', { name: /トップページに戻る/ })
    expect(backLink).toBeInTheDocument()
    expect(backLink).toHaveAttribute('href', '/')
  })

  it('sends email successfully', async () => {
    const user = userEvent.setup()
    render(<SendEmailPage />)

    // Fill the form
    const emailInput = screen.getByLabelText('送信先メールアドレス')
    const nameInput = screen.getByLabelText('お名前')
    const submitButton = screen.getByRole('button', { name: /メールを送信/ })

    await user.type(emailInput, 'test@example.com')
    await user.type(nameInput, 'Test User')
    await user.click(submitButton)

    // Wait for success message
    await waitFor(() => {
      expect(screen.getByText('メールを送信しました')).toBeInTheDocument()
    }, { timeout: 5000 })

    // Form should be cleared
    await waitFor(() => {
      expect(emailInput).toHaveValue('')
      expect(nameInput).toHaveValue('')
    }, { timeout: 5000 })
  })

  it('handles API errors gracefully', async () => {
    mockSendEmail.mockRejectedValue(new Error('Internal Server Error'))

    const user = userEvent.setup()
    render(<SendEmailPage />)

    // Fill the form
    const emailInput = screen.getByLabelText('送信先メールアドレス')
    const nameInput = screen.getByLabelText('お名前')
    const submitButton = screen.getByRole('button', { name: /メールを送信/ })

    await user.type(emailInput, 'test@example.com')
    await user.type(nameInput, 'Test User')
    await user.click(submitButton)

    // Wait for error message
    await waitFor(() => {
      expect(screen.getByText(/Internal Server Error/i)).toBeInTheDocument()
    }, { timeout: 5000 })
  })

  it('shows loading state during email submission', async () => {
    let resolvePromise: (value: unknown) => void
    const promise = new Promise((resolve) => {
      resolvePromise = resolve
    })

    mockSendEmail.mockImplementation(() => promise)

    const user = userEvent.setup()
    render(<SendEmailPage />)

    // Fill the form
    const emailInput = screen.getByLabelText('送信先メールアドレス')
    const nameInput = screen.getByLabelText('お名前')
    const submitButton = screen.getByRole('button', { name: /メールを送信/ })

    await user.type(emailInput, 'test@example.com')
    await user.type(nameInput, 'Test User')
    await user.click(submitButton)

    // Check loading state
    expect(screen.getByText('送信中...')).toBeInTheDocument()
    expect(submitButton).toBeDisabled()

    // Resolve the promise
    resolvePromise!({ success: true, message: 'メールを送信しました' })

    // Wait for success message
    await waitFor(() => {
      expect(screen.getByText('メールを送信しました')).toBeInTheDocument()
    }, { timeout: 5000 })
  })
})
