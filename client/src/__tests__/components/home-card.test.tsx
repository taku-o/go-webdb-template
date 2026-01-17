import { render, screen } from '@testing-library/react'

// Mock react-markdown
jest.mock('react-markdown', () => {
  return {
    __esModule: true,
    default: ({ children }: { children: string }) => <>{children}</>,
  }
})

import Card from '@/components/home/card'

describe('Card', () => {
  const defaultProps = {
    title: 'Test Card Title',
    description: 'This is a test description for the card component.',
    demo: <div data-testid="demo-content">Demo Content</div>,
  }

  it('renders the card with title and description', () => {
    render(<Card {...defaultProps} />)

    expect(screen.getByText('Test Card Title')).toBeInTheDocument()
    expect(screen.getByText('This is a test description for the card component.')).toBeInTheDocument()
  })

  it('renders the demo content', () => {
    render(<Card {...defaultProps} />)

    expect(screen.getByTestId('demo-content')).toBeInTheDocument()
    expect(screen.getByText('Demo Content')).toBeInTheDocument()
  })

  it('renders description text', () => {
    const propsWithLink = {
      ...defaultProps,
      description: 'Check out this link for more info.',
    }

    render(<Card {...propsWithLink} />)

    expect(screen.getByText('Check out this link for more info.')).toBeInTheDocument()
  })

  it('applies large class when large prop is true', () => {
    const { container } = render(<Card {...defaultProps} large={true} />)

    const cardElement = container.firstChild as HTMLElement
    expect(cardElement).toHaveClass('md:col-span-2')
  })

  it('does not apply large class when large prop is false', () => {
    const { container } = render(<Card {...defaultProps} large={false} />)

    const cardElement = container.firstChild as HTMLElement
    expect(cardElement).not.toHaveClass('md:col-span-2')
  })

  it('does not apply large class when large prop is not provided', () => {
    const { container } = render(<Card {...defaultProps} />)

    const cardElement = container.firstChild as HTMLElement
    expect(cardElement).not.toHaveClass('md:col-span-2')
  })
})
