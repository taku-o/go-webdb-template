import { render, screen, waitFor, act, fireEvent } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { FeedForm } from '@/components/feed/feed-form'

describe('FeedForm', () => {
  const mockOnSubmit = jest.fn()

  beforeEach(() => {
    jest.clearAllMocks()
  })

  it('renders the form with textarea and submit button', () => {
    render(<FeedForm onSubmit={mockOnSubmit} />)

    expect(screen.getByLabelText('投稿内容')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: '投稿する' })).toBeInTheDocument()
    expect(screen.getByText('0 / 280')).toBeInTheDocument()
  })

  it('displays character count as user types', async () => {
    const user = userEvent.setup()
    render(<FeedForm onSubmit={mockOnSubmit} />)

    const textarea = screen.getByLabelText('投稿内容')
    await user.type(textarea, 'Hello World')

    expect(screen.getByText('11 / 280')).toBeInTheDocument()
  })

  it('submit button is disabled when textarea is empty', () => {
    render(<FeedForm onSubmit={mockOnSubmit} />)

    const button = screen.getByRole('button', { name: '投稿する' })
    expect(button).toBeDisabled()
  })

  it('submit button is enabled when textarea has content', async () => {
    const user = userEvent.setup()
    render(<FeedForm onSubmit={mockOnSubmit} />)

    const textarea = screen.getByLabelText('投稿内容')
    await user.type(textarea, 'Test content')

    const button = screen.getByRole('button', { name: '投稿する' })
    expect(button).not.toBeDisabled()
  })

  it('shows error when trying to submit empty content', async () => {
    render(<FeedForm onSubmit={mockOnSubmit} />)

    // Clear the textarea and force submit (button is disabled, so we need to trigger form submit)
    const form = screen.getByLabelText('投稿内容').closest('form')!
    await act(async () => {
      fireEvent.submit(form)
    })

    await waitFor(() => {
      expect(screen.getByRole('alert')).toHaveTextContent('投稿内容を入力してください')
    })
  })

  it('shows error when content exceeds 280 characters', async () => {
    const user = userEvent.setup()
    render(<FeedForm onSubmit={mockOnSubmit} />)

    const textarea = screen.getByLabelText('投稿内容')
    const longContent = 'a'.repeat(281)
    await user.type(textarea, longContent)

    expect(screen.getByText('281 / 280')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: '投稿する' })).toBeDisabled()
  })

  it('calls onSubmit with content and clears form on success', async () => {
    const user = userEvent.setup()
    mockOnSubmit.mockResolvedValueOnce(undefined)

    render(<FeedForm onSubmit={mockOnSubmit} />)

    const textarea = screen.getByLabelText('投稿内容')
    await user.type(textarea, 'Test post content')

    const button = screen.getByRole('button', { name: '投稿する' })
    await user.click(button)

    await waitFor(() => {
      expect(mockOnSubmit).toHaveBeenCalledWith('Test post content')
    })

    // Form should be cleared after successful submit
    await waitFor(() => {
      expect(textarea).toHaveValue('')
    })
  })

  it('shows loading state during submission', async () => {
    const user = userEvent.setup()
    render(<FeedForm onSubmit={mockOnSubmit} isSubmitting={true} />)

    expect(screen.getByRole('button', { name: '投稿中...' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: '投稿中...' })).toBeDisabled()
    expect(screen.getByLabelText('投稿内容')).toBeDisabled()
  })

  it('shows error message when submission fails', async () => {
    const user = userEvent.setup()
    mockOnSubmit.mockRejectedValueOnce(new Error('投稿に失敗しました'))

    render(<FeedForm onSubmit={mockOnSubmit} />)

    const textarea = screen.getByLabelText('投稿内容')
    await user.type(textarea, 'Test content')

    const button = screen.getByRole('button', { name: '投稿する' })
    await user.click(button)

    await waitFor(() => {
      expect(screen.getByRole('alert')).toHaveTextContent('投稿に失敗しました')
    })
  })
})
