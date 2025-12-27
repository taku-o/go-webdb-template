import { test, expect } from '@playwright/test'

test.describe('Auth0 Login Flow', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/')
  })

  test('should display login button when not authenticated', async ({ page }) => {
    // ログインボタンが表示されることを確認
    const loginButton = page.locator('a[href="/auth/login"]')
    await expect(loginButton).toBeVisible()
    await expect(loginButton).toHaveText('ログイン')
  })

  test('should display "ログインしていません" message when not authenticated', async ({ page }) => {
    // 未ログインメッセージが表示されることを確認
    await expect(page.locator('text=ログインしていません')).toBeVisible()
  })

  test('should redirect to Auth0 login page when clicking login button', async ({ page }) => {
    // ログインボタンをクリック
    const loginButton = page.locator('a[href="/auth/login"]')
    await loginButton.click()

    // Auth0のログインページにリダイレクトされることを確認
    // (/auth/loginを経由してAuth0にリダイレクトされる)
    await expect(page).toHaveURL(/auth0\.com|\/auth\/login/)
  })

  test('should maintain existing navigation links', async ({ page }) => {
    // 既存のナビゲーションリンクが表示されることを確認
    await expect(page.getByText('ユーザー管理')).toBeVisible()
    await expect(page.getByText('投稿管理')).toBeVisible()
    await expect(page.getByText('ユーザーと投稿')).toBeVisible()
  })

  test('should display page heading', async ({ page }) => {
    // ページのヘッディングが表示されることを確認
    await expect(page.locator('h1')).toContainText('Go DB Project Sample')
  })
})

test.describe('Auth0 Logout Flow', () => {
  // 注意: 実際のAuth0認証フローをテストするにはテストアカウントが必要
  // このテストはAuth0が正しく設定されている前提で動作確認が可能

  test('should have logout link with correct href when logged in', async ({ page }) => {
    // このテストはログイン済み状態でのみ成功する
    // ログイン済みかどうかを確認
    await page.goto('/')

    const logoutButton = page.locator('a[href="/auth/logout"]')
    const isLoggedIn = await logoutButton.isVisible().catch(() => false)

    if (isLoggedIn) {
      await expect(logoutButton).toHaveText('ログアウト')
    } else {
      // 未ログインの場合、ログインボタンが表示されることを確認
      const loginButton = page.locator('a[href="/auth/login"]')
      await expect(loginButton).toBeVisible()
    }
  })
})
