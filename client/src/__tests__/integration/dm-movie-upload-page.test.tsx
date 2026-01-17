import { render, screen } from '@testing-library/react'
import MovieUploadPage from '@/app/dm_movie/upload/page'

// Note: Uppy components are mocked in jest.setup.js

describe('MovieUploadPage Integration', () => {
  it('displays page title', async () => {
    render(<MovieUploadPage />)

    expect(screen.getByText('動画ファイルアップロード')).toBeInTheDocument()
  })

  it('displays back link to top page', async () => {
    render(<MovieUploadPage />)

    const backLink = screen.getByRole('link', { name: /トップページに戻る/ })
    expect(backLink).toBeInTheDocument()
    expect(backLink).toHaveAttribute('href', '/')
  })

  it('displays upload card section', async () => {
    render(<MovieUploadPage />)

    expect(screen.getByText('動画アップロード')).toBeInTheDocument()
    expect(screen.getByText(/MP4形式の動画ファイルをアップロードできます/)).toBeInTheDocument()
    expect(screen.getByText(/最大ファイルサイズは2GB/)).toBeInTheDocument()
  })
})
