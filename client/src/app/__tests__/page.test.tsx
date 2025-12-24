import { render, screen } from '@testing-library/react'
import HomePage from '../page'

describe('HomePage', () => {
  it('renders the main heading', () => {
    render(<HomePage />)

    const heading = screen.getByRole('heading', { level: 1 })
    expect(heading).toBeInTheDocument()
    expect(heading).toHaveTextContent(/Go DB Project Sample/)
  })

  it('renders navigation links', () => {
    render(<HomePage />)

    // Check for link texts (using getAllByText since heading and link may have same text)
    const userLinks = screen.getAllByText(/ユーザー管理/)
    expect(userLinks.length).toBeGreaterThan(0)

    const postLinks = screen.getAllByText(/投稿管理/)
    expect(postLinks.length).toBeGreaterThan(0)

    const userPostLinks = screen.getAllByText(/ユーザーと投稿/)
    expect(userPostLinks.length).toBeGreaterThan(0)
  })

  it('displays tech stack', () => {
    render(<HomePage />)

    // Check for tech stack descriptions
    expect(screen.getByText(/Go \(Sharding対応\)/)).toBeInTheDocument()
    expect(screen.getByText(/Next\.js 14 \(App Router\)/)).toBeInTheDocument()
    expect(screen.getByText(/TypeScript/)).toBeInTheDocument()
  })

  it('has links with correct hrefs', () => {
    render(<HomePage />)

    const links = screen.getAllByRole('link')
    const usersLink = links.find(link => link.getAttribute('href') === '/users')
    expect(usersLink).toBeDefined()

    const postsLink = links.find(link => link.getAttribute('href') === '/posts')
    expect(postsLink).toBeDefined()

    const userPostsLink = links.find(link => link.getAttribute('href') === '/user-posts')
    expect(userPostsLink).toBeDefined()
  })
})
