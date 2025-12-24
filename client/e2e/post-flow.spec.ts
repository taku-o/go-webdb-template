import { test, expect } from '@playwright/test'

test.describe('Post Management Flow', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/')
  })

  test('should navigate to posts page', async ({ page }) => {
    await page.click('text=投稿管理')
    await expect(page).toHaveURL('/posts')
    await expect(page.locator('h1')).toContainText('投稿管理')
  })

  test('should create user and post', async ({ page }) => {
    // First create a user
    await page.goto('/users')
    await page.fill('input[type="text"]', 'Post Author')
    await page.fill('input[type="email"]', 'author@example.com')
    await page.click('button:has-text("作成")')
    await expect(page.locator('text=Post Author')).toBeVisible()

    // Navigate to posts page
    await page.goto('/posts')

    // Select the user from dropdown
    await page.selectOption('select', { label: /Post Author/ })

    // Fill in post form
    await page.fill('input[type="text"]', 'My First Post')
    await page.fill('textarea', 'This is the content of my first post.')

    // Submit form
    await page.click('button:has-text("作成")')

    // Wait for post to appear
    await expect(page.locator('text=My First Post')).toBeVisible()
    await expect(page.locator('text=This is the content of my first post.')).toBeVisible()
  })

  test('should delete a post', async ({ page }) => {
    // First create a user
    await page.goto('/users')
    await page.fill('input[type="text"]', 'Delete Post User')
    await page.fill('input[type="email"]', 'deletepost@example.com')
    await page.click('button:has-text("作成")')

    // Create a post
    await page.goto('/posts')
    await page.selectOption('select', { label: /Delete Post User/ })
    await page.fill('input[type="text"]', 'Post to Delete')
    await page.fill('textarea', 'This post will be deleted')
    await page.click('button:has-text("作成")')

    // Wait for post to appear
    await expect(page.locator('text=Post to Delete')).toBeVisible()

    // Delete the post
    page.on('dialog', dialog => dialog.accept())
    const deleteButton = page.locator('text=Post to Delete')
      .locator('..')
      .locator('..')
      .locator('button:has-text("削除")')
    await deleteButton.click()

    // Post should be removed
    await expect(page.locator('text=Post to Delete')).not.toBeVisible()
  })

  test('should display message when no users exist', async ({ page }) => {
    await page.goto('/posts')

    // Should show message about creating users first
    const message = page.locator('text=先にユーザーを作成してください')
    await expect(message).toBeVisible()
  })
})
