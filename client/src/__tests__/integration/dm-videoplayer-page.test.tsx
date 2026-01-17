import { render, screen } from '@testing-library/react'
import VideoPlayerDemoPage from '@/app/dm_videoplayer/page'

describe('VideoPlayerDemoPage Integration', () => {
  it('displays page title', async () => {
    render(<VideoPlayerDemoPage />)

    expect(screen.getByText('動画プレイヤー')).toBeInTheDocument()
  })

  it('displays back link to top page', async () => {
    render(<VideoPlayerDemoPage />)

    const backLink = screen.getByRole('link', { name: /トップページに戻る/ })
    expect(backLink).toBeInTheDocument()
    expect(backLink).toHaveAttribute('href', '/')
  })

  it('displays card description', async () => {
    render(<VideoPlayerDemoPage />)

    expect(screen.getByText(/動画プレイヤーコンポーネントのデモページです/)).toBeInTheDocument()
  })

  it('renders video player component', async () => {
    render(<VideoPlayerDemoPage />)

    // VideoPlayer component should be rendered (video element)
    const videoElement = document.querySelector('video')
    expect(videoElement).toBeInTheDocument()
  })
})
