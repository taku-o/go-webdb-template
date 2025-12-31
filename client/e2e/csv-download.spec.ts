import { test, expect } from '@playwright/test'

test.describe('CSV Download Flow', () => {
  test.beforeEach(async ({ page }) => {
    // Navigate to dm-users page
    await page.goto('/dm-users')
  })

  test('should display CSV download button', async ({ page }) => {
    // Check if CSV download button is visible
    const downloadButton = page.locator('button:has-text("CSVダウンロード")')
    await expect(downloadButton).toBeVisible()
  })

  test('should call CSV API when button is clicked', async ({ page }) => {
    // Wait for page to load
    await expect(page.locator('h1')).toContainText('ユーザー管理')

    // Set up request interception
    let csvRequestMade = false
    let csvResponseStatus = 0

    await page.route('**/api/export/dm-users/csv', async route => {
      csvRequestMade = true
      // Continue with actual request
      const response = await route.fetch()
      csvResponseStatus = response.status()
      await route.fulfill({ response })
    })

    // Click CSV download button
    const downloadButton = page.locator('button:has-text("CSVダウンロード")')
    await downloadButton.click()

    // Wait for button to return to normal state (download complete)
    await expect(downloadButton).toContainText('CSVダウンロード', { timeout: 10000 })

    // Verify CSV API was called
    expect(csvRequestMade).toBe(true)
    expect(csvResponseStatus).toBe(200)
  })

  test('should show loading state during download', async ({ page }) => {
    // Wait for page to load
    await expect(page.locator('h1')).toContainText('ユーザー管理')

    // Mock slow response to see loading state
    await page.route('**/api/export/dm-users/csv', async route => {
      // Delay response by 500ms
      await new Promise(resolve => setTimeout(resolve, 500))
      await route.fulfill({
        status: 200,
        headers: {
          'Content-Type': 'text/csv; charset=utf-8',
          'Content-Disposition': 'attachment; filename="dm-users.csv"',
        },
        body: 'ID,Name,Email,Created At,Updated At\n',
      })
    })

    // Click download button
    const downloadButton = page.locator('button:has-text("CSVダウンロード")')
    await downloadButton.click()

    // Check loading state appears
    await expect(page.locator('button:has-text("ダウンロード中...")')).toBeVisible()

    // Wait for download to complete
    await expect(page.locator('button:has-text("CSVダウンロード")')).toBeVisible({ timeout: 5000 })
  })

  test('should verify CSV response content structure', async ({ page }) => {
    // Wait for page to load
    await expect(page.locator('h1')).toContainText('ユーザー管理')

    // Set up request interception to capture response
    let csvContent = ''

    await page.route('**/api/export/dm-users/csv', async route => {
      const response = await route.fetch()
      csvContent = await response.text()
      await route.fulfill({ response })
    })

    // Click CSV download button
    const downloadButton = page.locator('button:has-text("CSVダウンロード")')
    await downloadButton.click()

    // Wait for download to complete
    await expect(downloadButton).toContainText('CSVダウンロード', { timeout: 10000 })

    // Verify CSV header is present
    expect(csvContent).toContain('ID,Name,Email,Created At,Updated At')

    // Verify CSV contains at least header row
    const lines = csvContent.trim().split('\n')
    expect(lines.length).toBeGreaterThanOrEqual(1)
  })

  test('should handle API error gracefully', async ({ page }) => {
    // Wait for page to load
    await expect(page.locator('h1')).toContainText('ユーザー管理')

    // Mock error response
    await page.route('**/api/export/dm-users/csv', async route => {
      await route.fulfill({
        status: 500,
        body: 'Internal server error',
      })
    })

    // Click download button
    const downloadButton = page.locator('button:has-text("CSVダウンロード")')
    await downloadButton.click()

    // Wait for button to return to normal state
    await expect(downloadButton).toContainText('CSVダウンロード', { timeout: 5000 })

    // Check error message is displayed
    await expect(page.locator('text=Internal server error')).toBeVisible()
  })
})
