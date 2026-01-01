'use client'

import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import SendEmailPage from '../page'
import { apiClient } from '@/lib/api'

// Mock apiClient
jest.mock('@/lib/api', () => ({
  apiClient: {
    sendEmail: jest.fn(),
  },
}))

// Mock next/link
jest.mock('next/link', () => {
  return function MockLink({ children, href }: { children: React.ReactNode; href: string }) {
    return <a href={href}>{children}</a>
  }
})

describe('SendEmailPage', () => {
  beforeEach(() => {
    jest.clearAllMocks()
  })

  describe('レンダリング', () => {
    it('ページタイトルが表示される', () => {
      render(<SendEmailPage />)
      expect(screen.getByText('メール送信')).toBeInTheDocument()
    })

    it('ウェルカムメール送信セクションが表示される', () => {
      render(<SendEmailPage />)
      expect(screen.getByText('ウェルカムメール送信')).toBeInTheDocument()
    })

    it('メールアドレス入力フィールドが表示される', () => {
      render(<SendEmailPage />)
      expect(screen.getByLabelText('送信先メールアドレス')).toBeInTheDocument()
    })

    it('名前入力フィールドが表示される', () => {
      render(<SendEmailPage />)
      expect(screen.getByLabelText('お名前')).toBeInTheDocument()
    })

    it('送信ボタンが表示される', () => {
      render(<SendEmailPage />)
      expect(screen.getByRole('button', { name: 'メールを送信' })).toBeInTheDocument()
    })

    it('トップページへのリンクが表示される', () => {
      render(<SendEmailPage />)
      const link = screen.getByText(/トップページに戻る/)
      expect(link).toBeInTheDocument()
      expect(link).toHaveAttribute('href', '/')
    })
  })

  describe('フォーム入力', () => {
    it('メールアドレスを入力できる', async () => {
      const user = userEvent.setup()
      render(<SendEmailPage />)

      const emailInput = screen.getByLabelText('送信先メールアドレス')
      await user.type(emailInput, 'test@example.com')

      expect(emailInput).toHaveValue('test@example.com')
    })

    it('名前を入力できる', async () => {
      const user = userEvent.setup()
      render(<SendEmailPage />)

      const nameInput = screen.getByLabelText('お名前')
      await user.type(nameInput, 'テスト太郎')

      expect(nameInput).toHaveValue('テスト太郎')
    })
  })

  describe('フォーム送信', () => {
    it('正常に送信が成功した場合、成功メッセージが表示される', async () => {
      const user = userEvent.setup()
      ;(apiClient.sendEmail as jest.Mock).mockResolvedValueOnce({
        success: true,
        message: 'メールを送信しました',
      })

      render(<SendEmailPage />)

      const emailInput = screen.getByLabelText('送信先メールアドレス')
      const nameInput = screen.getByLabelText('お名前')
      const submitButton = screen.getByRole('button', { name: 'メールを送信' })

      await user.type(emailInput, 'test@example.com')
      await user.type(nameInput, 'テスト太郎')
      await user.click(submitButton)

      await waitFor(() => {
        expect(screen.getByText('メールを送信しました')).toBeInTheDocument()
      })
    })

    it('送信成功後、フォームがクリアされる', async () => {
      const user = userEvent.setup()
      ;(apiClient.sendEmail as jest.Mock).mockResolvedValueOnce({
        success: true,
        message: 'メールを送信しました',
      })

      render(<SendEmailPage />)

      const emailInput = screen.getByLabelText('送信先メールアドレス')
      const nameInput = screen.getByLabelText('お名前')
      const submitButton = screen.getByRole('button', { name: 'メールを送信' })

      await user.type(emailInput, 'test@example.com')
      await user.type(nameInput, 'テスト太郎')
      await user.click(submitButton)

      await waitFor(() => {
        expect(emailInput).toHaveValue('')
        expect(nameInput).toHaveValue('')
      })
    })

    it('APIが失敗を返した場合、エラーメッセージが表示される', async () => {
      const user = userEvent.setup()
      ;(apiClient.sendEmail as jest.Mock).mockResolvedValueOnce({
        success: false,
        message: '送信に失敗しました',
      })

      render(<SendEmailPage />)

      const emailInput = screen.getByLabelText('送信先メールアドレス')
      const nameInput = screen.getByLabelText('お名前')
      const submitButton = screen.getByRole('button', { name: 'メールを送信' })

      await user.type(emailInput, 'test@example.com')
      await user.type(nameInput, 'テスト太郎')
      await user.click(submitButton)

      await waitFor(() => {
        expect(screen.getByText('送信に失敗しました')).toBeInTheDocument()
      })
    })

    it('APIがエラーをスローした場合、エラーメッセージが表示される', async () => {
      const user = userEvent.setup()
      ;(apiClient.sendEmail as jest.Mock).mockRejectedValueOnce(new Error('ネットワークエラー'))

      render(<SendEmailPage />)

      const emailInput = screen.getByLabelText('送信先メールアドレス')
      const nameInput = screen.getByLabelText('お名前')
      const submitButton = screen.getByRole('button', { name: 'メールを送信' })

      await user.type(emailInput, 'test@example.com')
      await user.type(nameInput, 'テスト太郎')
      await user.click(submitButton)

      await waitFor(() => {
        expect(screen.getByText('ネットワークエラー')).toBeInTheDocument()
      })
    })

    it('送信中はボタンが無効化され、テキストが変わる', async () => {
      const user = userEvent.setup()
      let resolvePromise: (value: { success: boolean; message: string }) => void
      const promise = new Promise<{ success: boolean; message: string }>((resolve) => {
        resolvePromise = resolve
      })
      ;(apiClient.sendEmail as jest.Mock).mockReturnValueOnce(promise)

      render(<SendEmailPage />)

      const emailInput = screen.getByLabelText('送信先メールアドレス')
      const nameInput = screen.getByLabelText('お名前')
      const submitButton = screen.getByRole('button', { name: 'メールを送信' })

      await user.type(emailInput, 'test@example.com')
      await user.type(nameInput, 'テスト太郎')
      await user.click(submitButton)

      expect(screen.getByRole('button', { name: '送信中...' })).toBeDisabled()

      resolvePromise!({ success: true, message: 'メールを送信しました' })

      await waitFor(() => {
        expect(screen.getByRole('button', { name: 'メールを送信' })).not.toBeDisabled()
      })
    })

    it('正しいパラメータでAPIが呼び出される', async () => {
      const user = userEvent.setup()
      ;(apiClient.sendEmail as jest.Mock).mockResolvedValueOnce({
        success: true,
        message: 'メールを送信しました',
      })

      render(<SendEmailPage />)

      const emailInput = screen.getByLabelText('送信先メールアドレス')
      const nameInput = screen.getByLabelText('お名前')
      const submitButton = screen.getByRole('button', { name: 'メールを送信' })

      await user.type(emailInput, 'test@example.com')
      await user.type(nameInput, 'テスト太郎')
      await user.click(submitButton)

      await waitFor(() => {
        expect(apiClient.sendEmail).toHaveBeenCalledWith(
          ['test@example.com'],
          'welcome',
          { Name: 'テスト太郎', Email: 'test@example.com' }
        )
      })
    })
  })
})
