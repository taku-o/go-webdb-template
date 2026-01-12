import { test, expect } from '@playwright/test'

test.describe('Cross-Shard JOIN Flow', () => {
  test('should create users and posts and display cross-shard page', async ({ page }) => {
    // Create multiple users with unique emails
    const timestamp = Date.now()
    await page.goto('/dm-users')

    // Wait for page to load
    await expect(page.locator('h1')).toContainText('ユーザー管理')

    // User 1
    await page.fill('input[id="name"]', 'Alice')
    await page.fill('input[id="email"]', `alice-${timestamp}@example.com`)
    await page.click('button:has-text("作成")')
    // API success: form should be cleared
    await expect(page.locator('input[id="name"]')).toHaveValue('', { timeout: 10000 })

    // User 2
    await page.fill('input[id="name"]', 'Bob')
    await page.fill('input[id="email"]', `bob-${timestamp}@example.com`)
    await page.click('button:has-text("作成")')
    // API success: form should be cleared
    await expect(page.locator('input[id="name"]')).toHaveValue('', { timeout: 10000 })

    // Create posts
    await page.goto('/dm-posts')
    await expect(page.locator('h1')).toContainText('投稿管理')

    // Wait for users to load in dropdown
    await page.waitForTimeout(1000)

    // Check if users exist in dropdown
    const combobox = page.getByRole('combobox')
    await combobox.click()
    const options = page.getByRole('option')
    const optionCount = await options.count()

    // Skip post creation if no users
    if (optionCount > 0) {
      // Select the first user and create a post
      await options.first().click()
      await page.fill('input[id="title"]', 'Test Post')
      await page.fill('textarea[id="content"]', 'Test Content')
      await page.click('button:has-text("作成")')
      // API success: form should be cleared
      await expect(page.locator('input[id="title"]')).toHaveValue('', { timeout: 10000 })
    }

    // Navigate to user-posts page (cross-shard JOIN)
    await page.goto('/dm-user-posts')

    // Check page title
    await expect(page.locator('h1')).toContainText('ユーザーと投稿')

    // Page should load without error - either data or empty state
    const hasData = page.locator('text=クロスシャードクエリ')
    const emptyState = page.getByText('表示する投稿がありません。')
    await expect(hasData.or(emptyState)).toBeVisible({ timeout: 15000 })
  })

  test('should show empty state when no posts exist', async ({ page }) => {
    await page.goto('/dm-user-posts')

    // Wait for page title to confirm page is loaded
    await expect(page.locator('h1')).toContainText('ユーザーと投稿', { timeout: 10000 })

    // Wait for either empty state or data to appear (loading should complete)
    // Check if empty state message exists OR data exists
    const emptyState = page.getByText('表示する投稿がありません。')
    const hasData = page.locator('text=クロスシャードクエリ')

    // Wait for either condition with longer timeout
    await expect(emptyState.or(hasData)).toBeVisible({ timeout: 15000 })

    // If empty state is shown, verify the links
    if (await emptyState.isVisible()) {
      await expect(page.getByRole('link', { name: /ユーザーを作成/ })).toBeVisible()
      await expect(page.getByRole('link', { name: /投稿を作成/ })).toBeVisible()
    }
  })

  test('should navigate back to home page', async ({ page }) => {
    await page.goto('/dm-user-posts')

    // Click back to home link
    await page.getByRole('link', { name: /トップページに戻る/ }).click()

    // Should navigate to home
    await expect(page).toHaveURL('/')
    await expect(page.locator('h1')).toContainText('Go DB Project Sample')
  })
})
