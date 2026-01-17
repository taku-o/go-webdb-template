import { render, screen } from '@testing-library/react'
import { LoadingOverlay } from '@/components/shared/loading-overlay'

describe('LoadingOverlay', () => {
  it('renders with default loading message', () => {
    render(<LoadingOverlay />)

    // Multiple elements with "Loading..." text exist (sr-only and visible message)
    const loadingTexts = screen.getAllByText('Loading...')
    expect(loadingTexts.length).toBeGreaterThanOrEqual(1)
  })

  it('renders with custom message', () => {
    render(<LoadingOverlay message="Please wait..." />)

    expect(screen.getByText('Please wait...')).toBeInTheDocument()
  })

  it('renders loading spinner', () => {
    render(<LoadingOverlay />)

    // The LoadingSpinner component renders with role="status"
    const spinners = screen.getAllByRole('status')
    expect(spinners.length).toBeGreaterThan(0)
  })

  it('has accessible live region', () => {
    const { container } = render(<LoadingOverlay />)

    // Find the container with aria-live
    const liveRegion = container.querySelector('[aria-live="polite"]')
    expect(liveRegion).toBeInTheDocument()
  })

  it('applies custom aria-label', () => {
    const { container } = render(<LoadingOverlay aria-label="Loading content" />)

    const statusElement = container.querySelector('[aria-label="Loading content"]')
    expect(statusElement).toBeInTheDocument()
  })
})
