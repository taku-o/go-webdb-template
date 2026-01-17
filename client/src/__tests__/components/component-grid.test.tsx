import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'

// Mock window.matchMedia
beforeAll(() => {
  Object.defineProperty(window, 'matchMedia', {
    writable: true,
    value: jest.fn().mockImplementation(query => ({
      matches: false,
      media: query,
      onchange: null,
      addListener: jest.fn(),
      removeListener: jest.fn(),
      addEventListener: jest.fn(),
      removeEventListener: jest.fn(),
      dispatchEvent: jest.fn(),
    })),
  })
})

import ComponentGrid from '@/components/home/component-grid'

describe('ComponentGrid', () => {
  it('renders Modal button', () => {
    render(<ComponentGrid />)

    expect(screen.getByRole('button', { name: 'Modal' })).toBeInTheDocument()
  })

  it('renders Popover button', () => {
    render(<ComponentGrid />)

    expect(screen.getByRole('button', { name: /Popover/ })).toBeInTheDocument()
  })

  it('renders Tooltip element', () => {
    render(<ComponentGrid />)

    expect(screen.getByText('Tooltip')).toBeInTheDocument()
  })

  it('opens modal when Modal button is clicked', async () => {
    const user = userEvent.setup()
    render(<ComponentGrid />)

    const modalButton = screen.getByRole('button', { name: 'Modal' })
    await user.click(modalButton)

    // Modal content should be visible (heading "Precedent")
    expect(screen.getByRole('heading', { name: 'Precedent' })).toBeInTheDocument()
  })

  it('opens popover when Popover button is clicked', async () => {
    const user = userEvent.setup()
    render(<ComponentGrid />)

    const popoverButton = screen.getByRole('button', { name: /Popover/ })
    await user.click(popoverButton)

    // Popover items should be visible
    expect(screen.getByRole('button', { name: 'Item 1' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Item 2' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Item 3' })).toBeInTheDocument()
  })
})
