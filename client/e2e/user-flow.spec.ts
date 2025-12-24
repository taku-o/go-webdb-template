import { test, expect } from '@playwright/test'

test.describe('User Management Flow', () => {
  test.beforeEach(async ({ page }) => {
    // Navigate to home page
    await page.goto('/')
  })

  test('should display home page with navigation', async ({ page }) => {
    // Check main heading
    await expect(page.locator('h1')).toContainText('データベースシャーディング')

    // Check navigation links
    await expect(page.getByText('ユーザー管理')).toBeVisible()
    await expect(page.getByText('投稿管理')).toBeVisible()
    await expect(page.getByText('ユーザーと投稿')).toBeVisible()
  })

  test('should navigate to users page', async ({ page }) => {
    // Click on users link
    await page.click('text=ユーザー管理')

    // Wait for navigation
    await expect(page).toHaveURL('/users')

    // Check page title
    await expect(page.locator('h1')).toContainText('ユーザー管理')
  })

  test('should create a new user', async ({ page }) => {
    // Navigate to users page
    await page.goto('/users')

    // Fill in user creation form
    await page.fill('input[type="text"]', 'E2E Test User')
    await page.fill('input[type="email"]', 'e2e-test@example.com')

    // Submit form
    await page.click('button:has-text("作成")')

    // Wait for user to appear in list
    await expect(page.locator('text=E2E Test User')).toBeVisible()
    await expect(page.locator('text=e2e-test@example.com')).toBeVisible()

    // Form should be cleared
    await expect(page.locator('input[type="text"]')).toHaveValue('')
    await expect(page.locator('input[type="email"]')).toHaveValue('')
  })

  test('should delete a user', async ({ page }) => {
    // Navigate to users page
    await page.goto('/users')

    // Create a user first
    await page.fill('input[type="text"]', 'User to Delete')
    await page.fill('input[type="email"]', 'delete@example.com')
    await page.click('button:has-text("作成")')

    // Wait for user to appear
    await expect(page.locator('text=User to Delete')).toBeVisible()

    // Click delete button
    page.on('dialog', dialog => dialog.accept())
    const deleteButton = page.locator('text=User to Delete')
      .locator('..')
      .locator('..')
      .locator('button:has-text("削除")')
    await deleteButton.click()

    // User should be removed
    await expect(page.locator('text=User to Delete')).not.toBeVisible()
  })

  test('should show validation for empty form', async ({ page }) => {
    await page.goto('/users')

    // Try to submit empty form
    const submitButton = page.locator('button:has-text("作成")')

    // Check if button is disabled or if HTML5 validation prevents submission
    const isDisabled = await submitButton.isDisabled()
    const nameInput = page.locator('input[type="text"]')
    const isRequired = await nameInput.getAttribute('required')

    expect(isDisabled || isRequired).toBeTruthy()
  })
})
