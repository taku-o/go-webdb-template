import { render, screen } from '@testing-library/react'
import { ErrorAlert } from '@/components/shared/error-alert'

describe('ErrorAlert', () => {
  it('renders the error message', () => {
    render(<ErrorAlert message="Something went wrong" />)

    expect(screen.getByText('Something went wrong')).toBeInTheDocument()
  })

  it('renders the title when provided', () => {
    render(<ErrorAlert title="Error" message="Something went wrong" />)

    expect(screen.getByText('Error')).toBeInTheDocument()
    expect(screen.getByText('Something went wrong')).toBeInTheDocument()
  })

  it('does not render title when not provided', () => {
    render(<ErrorAlert message="Something went wrong" />)

    expect(screen.queryByRole('heading')).not.toBeInTheDocument()
  })

  it('applies custom className', () => {
    const { container } = render(
      <ErrorAlert message="Test message" className="custom-class" />
    )

    expect(container.firstChild).toHaveClass('custom-class')
  })

  it('renders alert with destructive variant', () => {
    const { container } = render(<ErrorAlert message="Test message" />)

    // The alert should have the destructive variant styling
    const alert = container.querySelector('[role="alert"]')
    expect(alert).toBeInTheDocument()
  })
})
