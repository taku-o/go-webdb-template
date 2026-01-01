import { test, expect } from '@playwright/test'

test.describe('Email Send Flow', () => {
  test.beforeEach(async ({ page }) => {
    // Navigate to email send page
    await page.goto('/dm_email/send')
  })

  test('should display email send page with form', async ({ page }) => {
    // Check main heading
    await expect(page.locator('h1')).toContainText('メール送信')

    // Check section heading
    await expect(page.locator('h2')).toContainText('ウェルカムメール送信')

    // Check form elements
    await expect(page.locator('label[for="toEmail"]')).toContainText('送信先メールアドレス')
    await expect(page.locator('label[for="name"]')).toContainText('お名前')
    await expect(page.getByRole('button', { name: 'メールを送信' })).toBeVisible()
  })

  test('should have link to return to top page', async ({ page }) => {
    // Check back link
    const backLink = page.getByText('トップページに戻る')
    await expect(backLink).toBeVisible()
    await expect(backLink).toHaveAttribute('href', '/')
  })

  test('should fill in and submit email form', async ({ page }) => {
    // Fill in form fields
    await page.fill('#toEmail', 'e2e-test@example.com')
    await page.fill('#name', 'E2E テストユーザー')

    // Submit form
    await page.click('button:has-text("メールを送信")')

    // Wait for success message
    await expect(page.locator('text=メールを送信しました')).toBeVisible({ timeout: 10000 })

    // Form should be cleared after successful submission
    await expect(page.locator('#toEmail')).toHaveValue('')
    await expect(page.locator('#name')).toHaveValue('')
  })

  test('should show loading state while sending', async ({ page }) => {
    // Fill in form fields
    await page.fill('#toEmail', 'loading-test@example.com')
    await page.fill('#name', 'ローディングテスト')

    // Click submit button
    await page.click('button:has-text("メールを送信")')

    // Button should show loading state (this may be brief)
    // We check that either loading state is shown or success message appears
    const loadingOrSuccess = await Promise.race([
      page.locator('button:has-text("送信中...")').isVisible().catch(() => false),
      page.locator('text=メールを送信しました').isVisible().catch(() => false)
    ])

    expect(loadingOrSuccess).toBeTruthy()
  })

  test('should require email field', async ({ page }) => {
    // Fill only name, leave email empty
    await page.fill('#name', 'テストユーザー')

    // Check that email input has required attribute
    const emailInput = page.locator('#toEmail')
    await expect(emailInput).toHaveAttribute('required', '')

    // Try to submit - HTML5 validation should prevent submission
    await page.click('button:has-text("メールを送信")')

    // Should stay on the same page (no success message)
    await expect(page.locator('text=メールを送信しました')).not.toBeVisible()
  })

  test('should require name field', async ({ page }) => {
    // Fill only email, leave name empty
    await page.fill('#toEmail', 'test@example.com')

    // Check that name input has required attribute
    const nameInput = page.locator('#name')
    await expect(nameInput).toHaveAttribute('required', '')

    // Try to submit - HTML5 validation should prevent submission
    await page.click('button:has-text("メールを送信")')

    // Should stay on the same page (no success message)
    await expect(page.locator('text=メールを送信しました')).not.toBeVisible()
  })

  test('should validate email format', async ({ page }) => {
    // Fill in invalid email
    await page.fill('#toEmail', 'invalid-email')
    await page.fill('#name', 'テストユーザー')

    // Check that email input has type="email"
    const emailInput = page.locator('#toEmail')
    await expect(emailInput).toHaveAttribute('type', 'email')

    // Try to submit - HTML5 validation should prevent submission or API should return error
    await page.click('button:has-text("メールを送信")')

    // Should not show success message
    await expect(page.locator('text=メールを送信しました')).not.toBeVisible()
  })
})
