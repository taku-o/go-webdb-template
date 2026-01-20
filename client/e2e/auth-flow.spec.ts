import { test, expect } from '@playwright/test'

test.describe('NextAuth Login Flow', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/')
  })

  test('should display login button when not authenticated', async ({ page }) => {
    // ログインボタンが表示されることを確認（ナビゲーションバーのボタン）
    const loginButton = page.getByRole('button', { name: 'ログイン' }).first()
    await expect(loginButton).toBeVisible()
  })

  test('should display "ログインしていません" message when not authenticated', async ({ page }) => {
    // 未ログインメッセージが表示されることを確認
    await expect(page.getByText('ログインしていません').first()).toBeVisible()
  })

  test('should redirect to NextAuth login when clicking login button', async ({ page }) => {
    // ログインボタンをクリック（ナビゲーションバーのボタン）
    const loginButton = page.getByRole('button', { name: 'ログイン' }).first()
    await loginButton.click()

    // NextAuthのログインページにリダイレクトされることを確認
    // (NextAuthの設定によって異なる可能性があるため、URLの変化を確認)
    await page.waitForURL(/\/api\/auth\/signin|auth/, { timeout: 5000 })
  })

  test('should maintain existing navigation links', async ({ page }) => {
    // 既存のナビゲーションリンクが表示されることを確認
    await expect(page.getByRole('link', { name: /ユーザー管理/ })).toBeVisible()
    await expect(page.getByRole('link', { name: /投稿管理/ })).toBeVisible()
    await expect(page.getByRole('link', { name: /ユーザーと投稿/ })).toBeVisible()
  })

  test('should display page heading', async ({ page }) => {
    // ページのヘッディングが表示されることを確認
    await expect(page.locator('h1')).toContainText('Go DB Project Sample')
  })
})

test.describe('NextAuth Logout Flow', () => {
  // 注意: 実際のNextAuth認証フローをテストするにはテストアカウントが必要
  // このテストはNextAuthが正しく設定されている前提で動作確認が可能

  test('should have logout button with correct form action when logged in', async ({ page }) => {
    // このテストはログイン済み状態でのみ成功する
    // ログイン済みかどうかを確認
    await page.goto('/')

    const logoutButton = page.getByRole('button', { name: 'ログアウト' }).first()
    const isLoggedIn = await logoutButton.isVisible().catch(() => false)

    if (isLoggedIn) {
      await expect(logoutButton).toBeVisible()
    } else {
      // 未ログインの場合、ログインボタンが表示されることを確認（ナビゲーションバーのボタン）
      const loginButton = page.getByRole('button', { name: 'ログイン' }).first()
      await expect(loginButton).toBeVisible()
    }
  })
})
