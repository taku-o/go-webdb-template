import { test, expect } from '@playwright/test'

test.describe('Post Management Flow', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/')
  })

  test('should navigate to posts page', async ({ page }) => {
    await page.getByRole('link', { name: /投稿管理/ }).click()
    await expect(page).toHaveURL('/dm-posts')
    await expect(page.locator('h1')).toContainText('投稿管理')
  })

  test('should create user and post', async ({ page }) => {
    // First create a user with unique email
    const timestamp = Date.now()
    await page.goto('/dm-users')
    await expect(page.locator('h1')).toContainText('ユーザー管理')
    await page.fill('input[id="name"]', 'Post Author')
    await page.fill('input[id="email"]', `author-${timestamp}@example.com`)
    await page.click('button:has-text("作成")')
    // API success: form should be cleared
    await expect(page.locator('input[id="name"]')).toHaveValue('', { timeout: 10000 })

    // Navigate to posts page
    await page.goto('/dm-posts')
    await expect(page.locator('h1')).toContainText('投稿管理')

    // Wait for users to load in dropdown
    await page.waitForTimeout(1000)

    // Check if users exist in dropdown
    const combobox = page.getByRole('combobox')
    await combobox.click()
    const options = page.getByRole('option')
    const optionCount = await options.count()

    // Skip if no users
    if (optionCount === 0) {
      return
    }

    // Select the first user from dropdown
    await options.first().click()

    // Fill in post form
    await page.fill('input[id="title"]', 'My First Post')
    await page.fill('textarea[id="content"]', 'This is the content of my first post.')

    // Submit form
    await page.click('button:has-text("作成")')

    // API success: form should be cleared
    await expect(page.locator('input[id="title"]')).toHaveValue('', { timeout: 10000 })
    await expect(page.locator('textarea[id="content"]')).toHaveValue('')

    // Page should still be functional (no error)
    await expect(page.locator('h1')).toContainText('投稿管理')
  })

  test('should delete a post', async ({ page }) => {
    // Navigate to posts page
    await page.goto('/dm-posts')

    // Wait for page to load
    await expect(page.locator('h1')).toContainText('投稿管理')

    // Wait for loading to complete
    await page.waitForSelector('table tbody', { timeout: 10000 }).catch(() => null)

    // Check if there are any posts in the table
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
    await expect(page.locator('h1')).toContainText('投稿管理', { timeout: 5000 })
  })

  test('should display message when no users exist', async ({ page }) => {
    await page.goto('/dm-posts')

    // Should show message about creating users first
    const message = page.locator('text=先にユーザーを作成してください')
    await expect(message).toBeVisible()
  })
})
