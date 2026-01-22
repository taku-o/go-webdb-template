import { signOutAction } from '@/lib/actions/auth-actions'
import { signOut } from '@/auth'
import { redirect } from 'next/navigation'

// Mock signOut from @/auth
const mockSignOut = signOut as jest.MockedFunction<typeof signOut>

// Mock redirect from next/navigation
jest.mock('next/navigation', () => ({
  ...jest.requireActual('next/navigation'),
  redirect: jest.fn(),
}))
const mockRedirect = redirect as jest.MockedFunction<typeof redirect>

describe('signOutAction', () => {
  const originalEnv = process.env

  beforeEach(() => {
    jest.clearAllMocks()
    // 環境変数をリセット
    process.env = { ...originalEnv }
  })

  afterAll(() => {
    process.env = originalEnv
  })

  describe('環境変数の検証', () => {
    it('AUTH0_ISSUERが設定されていない場合、エラーをthrowする', async () => {
      delete process.env.AUTH0_ISSUER
      process.env.AUTH0_CLIENT_ID = 'test-client-id'
      process.env.NEXT_PUBLIC_APP_BASE_URL = 'http://localhost:3000'

      await expect(signOutAction()).rejects.toThrow('AUTH0_ISSUER is not set')
    })

    it('AUTH0_CLIENT_IDが設定されていない場合、エラーをthrowする', async () => {
      process.env.AUTH0_ISSUER = 'https://example.auth0.com'
      delete process.env.AUTH0_CLIENT_ID
      process.env.NEXT_PUBLIC_APP_BASE_URL = 'http://localhost:3000'

      await expect(signOutAction()).rejects.toThrow('AUTH0_CLIENT_ID is not set')
    })

    it('NEXT_PUBLIC_APP_BASE_URLが設定されていない場合、エラーをthrowする', async () => {
      process.env.AUTH0_ISSUER = 'https://example.auth0.com'
      process.env.AUTH0_CLIENT_ID = 'test-client-id'
      delete process.env.NEXT_PUBLIC_APP_BASE_URL

      await expect(signOutAction()).rejects.toThrow('NEXT_PUBLIC_APP_BASE_URL is not set')
    })
  })

  describe('Auth0ログアウトURLの構築', () => {
    it('正しいAuth0ログアウトURLを構築してsignOutを呼び出す', async () => {
      process.env.AUTH0_ISSUER = 'https://example.auth0.com'
      process.env.AUTH0_CLIENT_ID = 'test-client-id'
      process.env.NEXT_PUBLIC_APP_BASE_URL = 'http://localhost:3000'

      await signOutAction()

      const expectedLogoutUrl = 'https://example.auth0.com/v2/logout?client_id=test-client-id&returnTo=http%3A%2F%2Flocalhost%3A3000'
      expect(mockSignOut).toHaveBeenCalledWith({
        redirect: false,
        redirectTo: expectedLogoutUrl,
      })
      expect(mockRedirect).toHaveBeenCalledWith(expectedLogoutUrl)
    })

    it('URLエンコーディングが正しく適用される', async () => {
      process.env.AUTH0_ISSUER = 'https://my-domain.auth0.com'
      process.env.AUTH0_CLIENT_ID = 'my-client-id'
      process.env.NEXT_PUBLIC_APP_BASE_URL = 'https://my-app.example.com'

      await signOutAction()

      const expectedLogoutUrl = 'https://my-domain.auth0.com/v2/logout?client_id=my-client-id&returnTo=https%3A%2F%2Fmy-app.example.com'
      expect(mockSignOut).toHaveBeenCalledWith({
        redirect: false,
        redirectTo: expectedLogoutUrl,
      })
      expect(mockRedirect).toHaveBeenCalledWith(expectedLogoutUrl)
    })
  })
})
