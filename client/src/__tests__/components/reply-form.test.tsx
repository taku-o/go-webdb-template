import { render, screen, waitFor, act, fireEvent } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { ReplyForm } from '@/components/feed/reply-form'

describe('ReplyForm', () => {
  const mockOnSubmit = jest.fn()

  beforeEach(() => {
    jest.clearAllMocks()
  })

  it('renders the form with textarea and submit button', () => {
    render(<ReplyForm onSubmit={mockOnSubmit} />)

    expect(screen.getByLabelText('返信内容')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: '返信する' })).toBeInTheDocument()
    expect(screen.getByText('0 / 280')).toBeInTheDocument()
  })

  it('displays character count as user types', async () => {
    const user = userEvent.setup()
    render(<ReplyForm onSubmit={mockOnSubmit} />)

    const textarea = screen.getByLabelText('返信内容')
    await user.type(textarea, 'Hello World')

    expect(screen.getByText('11 / 280')).toBeInTheDocument()
  })

  it('submit button is disabled when textarea is empty', () => {
    render(<ReplyForm onSubmit={mockOnSubmit} />)

    const button = screen.getByRole('button', { name: '返信する' })
    expect(button).toBeDisabled()
  })

  it('submit button is enabled when textarea has content', async () => {
    const user = userEvent.setup()
    render(<ReplyForm onSubmit={mockOnSubmit} />)

    const textarea = screen.getByLabelText('返信内容')
    await user.type(textarea, 'Test content')

    const button = screen.getByRole('button', { name: '返信する' })
    expect(button).not.toBeDisabled()
  })

  it('shows error when trying to submit empty content', async () => {
    render(<ReplyForm onSubmit={mockOnSubmit} />)

    // Force submit via form event
    const form = screen.getByLabelText('返信内容').closest('form')!
    await act(async () => {
      fireEvent.submit(form)
    })

    await waitFor(() => {
      expect(screen.getByRole('alert')).toHaveTextContent('返信内容を入力してください')
    })
  })

  it('shows error when content exceeds 280 characters', async () => {
    const user = userEvent.setup()
    render(<ReplyForm onSubmit={mockOnSubmit} />)

    const textarea = screen.getByLabelText('返信内容')
    const longContent = 'a'.repeat(281)
    await user.type(textarea, longContent)

    expect(screen.getByText('281 / 280')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: '返信する' })).toBeDisabled()
  })

  it('calls onSubmit with content and clears form on success', async () => {
    const user = userEvent.setup()
    mockOnSubmit.mockResolvedValueOnce(undefined)

    render(<ReplyForm onSubmit={mockOnSubmit} />)

    const textarea = screen.getByLabelText('返信内容')
    await user.type(textarea, 'Test reply content')

    const button = screen.getByRole('button', { name: '返信する' })
    await user.click(button)

    await waitFor(() => {
      expect(mockOnSubmit).toHaveBeenCalledWith('Test reply content')
    })

    // Form should be cleared after successful submit
    await waitFor(() => {
      expect(textarea).toHaveValue('')
    })
  })

  it('shows loading state during submission', async () => {
    render(<ReplyForm onSubmit={mockOnSubmit} isSubmitting={true} />)

    expect(screen.getByRole('button', { name: '返信中...' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: '返信中...' })).toBeDisabled()
    expect(screen.getByLabelText('返信内容')).toBeDisabled()
  })

  it('shows error message when submission fails', async () => {
    const user = userEvent.setup()
    mockOnSubmit.mockRejectedValueOnce(new Error('返信に失敗しました'))

    render(<ReplyForm onSubmit={mockOnSubmit} />)

    const textarea = screen.getByLabelText('返信内容')
    await user.type(textarea, 'Test content')

    const button = screen.getByRole('button', { name: '返信する' })
    await user.click(button)

    await waitFor(() => {
      expect(screen.getByRole('alert')).toHaveTextContent('返信に失敗しました')
    })
  })

  it('displays correct placeholder text', () => {
    render(<ReplyForm onSubmit={mockOnSubmit} />)

    expect(screen.getByPlaceholderText('返信を入力...')).toBeInTheDocument()
  })
})
