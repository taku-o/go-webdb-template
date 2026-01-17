import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { VideoPlayer } from '@/components/video-player/video-player'

// Mock hls.js
jest.mock('hls.js', () => {
  return {
    __esModule: true,
    default: {
      isSupported: jest.fn(() => true),
      Events: {
        MANIFEST_PARSED: 'hlsManifestParsed',
        ERROR: 'hlsError',
      },
      ErrorTypes: {
        NETWORK_ERROR: 'networkError',
        MEDIA_ERROR: 'mediaError',
      },
    },
  }
})

describe('VideoPlayer', () => {
  const defaultProps = {
    videoUrl: '/test-video.mp4',
    thumbnailUrl: '/test-thumbnail.jpg',
  }

  beforeEach(() => {
    jest.clearAllMocks()
  })

  it('renders video element with poster image', () => {
    render(<VideoPlayer {...defaultProps} />)

    const video = document.querySelector('video')
    expect(video).toBeInTheDocument()
    expect(video).toHaveAttribute('poster', '/test-thumbnail.jpg')
  })

  it('renders source element for MP4 videos', () => {
    render(<VideoPlayer {...defaultProps} />)

    const source = document.querySelector('source')
    expect(source).toBeInTheDocument()
    expect(source).toHaveAttribute('src', '/test-video.mp4')
    expect(source).toHaveAttribute('type', 'video/mp4')
  })

  it('shows custom play button for HLS videos', () => {
    render(<VideoPlayer videoUrl="/test-video.m3u8" thumbnailUrl="/thumbnail.jpg" />)

    const playButton = screen.getByRole('button', { name: '動画を再生' })
    expect(playButton).toBeInTheDocument()
  })

  it('does not show custom play button for MP4 videos', () => {
    render(<VideoPlayer {...defaultProps} />)

    const playButton = screen.queryByRole('button', { name: '動画を再生' })
    expect(playButton).not.toBeInTheDocument()
  })

  it('applies custom className', () => {
    const { container } = render(
      <VideoPlayer {...defaultProps} className="custom-video-class" />
    )

    expect(container.firstChild).toHaveClass('custom-video-class')
  })

  it('video has playsInline attribute for mobile devices', () => {
    render(<VideoPlayer {...defaultProps} />)

    const video = document.querySelector('video')
    expect(video).toHaveAttribute('playsinline')
  })

  it('video has nodownload in controlsList', () => {
    render(<VideoPlayer {...defaultProps} />)

    const video = document.querySelector('video')
    expect(video).toHaveAttribute('controlslist', 'nodownload')
  })

  it('video has preload set to none for performance', () => {
    render(<VideoPlayer {...defaultProps} />)

    const video = document.querySelector('video')
    expect(video).toHaveAttribute('preload', 'none')
  })

  it('displays fallback text for unsupported browsers', () => {
    render(<VideoPlayer {...defaultProps} />)

    expect(screen.getByText('お使いのブラウザは動画タグをサポートしていません。')).toBeInTheDocument()
  })
})
