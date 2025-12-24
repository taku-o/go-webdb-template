import { test, expect } from '@playwright/test'

test.describe('Cross-Shard JOIN Flow', () => {
  test('should display joined user and post data', async ({ page }) => {
    // Create multiple users
    await page.goto('/users')

    // User 1
    await page.fill('input[type="text"]', 'Alice')
    await page.fill('input[type="email"]', 'alice@example.com')
    await page.click('button:has-text("作成")')
    await expect(page.locator('text=Alice')).toBeVisible()

    // User 2
    await page.fill('input[type="text"]', 'Bob')
    await page.fill('input[type="email"]', 'bob@example.com')
    await page.click('button:has-text("作成")')
    await expect(page.locator('text=Bob')).toBeVisible()

    // Create posts for both users
    await page.goto('/posts')

    // Post for Alice
    await page.selectOption('select', { label: /Alice/ })
    await page.fill('input[type="text"]', 'Alice\'s Post')
    await page.fill('textarea', 'Content from Alice')
    await page.click('button:has-text("作成")')
    await expect(page.locator('text=Alice\'s Post')).toBeVisible()

    // Post for Bob
    await page.selectOption('select', { label: /Bob/ })
    await page.fill('input[type="text"]', 'Bob\'s Post')
    await page.fill('textarea', 'Content from Bob')
    await page.click('button:has-text("作成")')
    await expect(page.locator('text=Bob\'s Post')).toBeVisible()

    // Navigate to user-posts page (cross-shard JOIN)
    await page.goto('/user-posts')

    // Check page title
    await expect(page.locator('h1')).toContainText('ユーザーと投稿')

    // Should display both posts with user information
    await expect(page.locator('text=Alice\'s Post')).toBeVisible()
    await expect(page.locator('text=Alice')).toBeVisible()
    await expect(page.locator('text=alice@example.com')).toBeVisible()

    await expect(page.locator('text=Bob\'s Post')).toBeVisible()
    await expect(page.locator('text=Bob')).toBeVisible()
    await expect(page.locator('text=bob@example.com')).toBeVisible()

    // Verify the cross-shard message is displayed
    await expect(page.locator('text=クロスシャードクエリ')).toBeVisible()
    await expect(page.locator('text=複数のShardからユーザーと投稿をJOINして取得')).toBeVisible()
  })

  test('should show empty state when no posts exist', async ({ page }) => {
    await page.goto('/user-posts')

    // Should show empty state message
    await expect(page.locator('text=表示する投稿がありません')).toBeVisible()

    // Should show links to create users and posts
    await expect(page.locator('a[href="/users"]')).toBeVisible()
    await expect(page.locator('a[href="/posts"]')).toBeVisible()
  })

  test('should navigate back to home page', async ({ page }) => {
    await page.goto('/user-posts')

    // Click back to home link
    await page.click('text=トップページに戻る')

    // Should navigate to home
    await expect(page).toHaveURL('/')
    await expect(page.locator('h1')).toContainText('データベースシャーディング')
  })
})
