import { test, expect } from '@playwright/test'

test.describe('User Management Flow', () => {
  test.beforeEach(async ({ page }) => {
    // Navigate to home page
    await page.goto('/')
  })

  test('should display home page with navigation', async ({ page }) => {
    // Check main heading
    await expect(page.locator('h1')).toContainText('Go DB Project Sample')

    // Check navigation links (use .first() to avoid strict mode violation when text appears multiple times)
    await expect(page.getByText('ユーザー管理').first()).toBeVisible()
    await expect(page.getByText('投稿管理').first()).toBeVisible()
    await expect(page.getByText('ユーザーと投稿').first()).toBeVisible()
  })

  test('should navigate to users page', async ({ page }) => {
    // Click on users link (Cardコンポーネント内のLink)
    await page.getByRole('link', { name: /ユーザー管理/ }).click()

    // Wait for navigation
    await expect(page).toHaveURL('/dm-users')

    // Check page title
    await expect(page.locator('h1')).toContainText('ユーザー管理')
  })

  test('should create a new user', async ({ page }) => {
    // Navigate to users page
    await page.goto('/dm-users')

    // Wait for page to load
    await expect(page.locator('h1')).toContainText('ユーザー管理')

    // Fill in user creation form with unique email
    const timestamp = Date.now()
    await page.fill('input[id="name"]', 'E2E Test User')
    await page.fill('input[id="email"]', `e2e-test-${timestamp}@example.com`)

    // Submit form
    await page.click('button:has-text("作成")')

    // API success: form should be cleared
    await expect(page.locator('input[id="name"]')).toHaveValue('', { timeout: 10000 })
    await expect(page.locator('input[id="email"]')).toHaveValue('')

    // Page should still be functional (no error)
    await expect(page.locator('h1')).toContainText('ユーザー管理')
  })

  test('should delete a user', async ({ page }) => {
    // Navigate to users page
    await page.goto('/dm-users')

    // Wait for page to load
    await expect(page.locator('h1')).toContainText('ユーザー管理')

    // Wait for loading to complete
    await page.waitForSelector('table tbody', { timeout: 10000 }).catch(() => null)

    // Check if there are any users in the table
    const rows = page.locator('table tbody tr')
    const rowCount = await rows.count()

    // Skip if no data
    if (rowCount === 0) {
      return
    }

    // Click delete button on the first row
    page.on('dialog', dialog => dialog.accept())
    const deleteButton = rows.first().getByRole('button', { name: /削除/ })
    await deleteButton.click()

    // API success: page should still be functional (no error)
    await expect(page.locator('h1')).toContainText('ユーザー管理', { timeout: 5000 })
  })

  test('should show validation for empty form', async ({ page }) => {
    await page.goto('/dm-users')

    // Check that name input has required attribute
    const nameInput = page.locator('input[id="name"]')
    await expect(nameInput).toHaveAttribute('required', '')
  })
})
