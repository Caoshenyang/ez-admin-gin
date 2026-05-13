import { test, expect } from '../fixtures'

test.describe('Token Expiration', () => {
  test('expired access token triggers redirect to login page', async ({ page }) => {
    // Set a fake expired token
    await page.goto('/')
    await page.evaluate(() => {
      localStorage.setItem('ez-admin-access-token', 'expired.jwt.token')
      localStorage.setItem('ez-admin-token-type', 'Bearer')
      localStorage.setItem(
        'ez-admin-user-info',
        JSON.stringify({
          userId: 1,
          username: 'admin',
          nickname: 'Admin',
          expiresAt: new Date(Date.now() - 60 * 1000).toISOString(), // expired
        }),
      )
    })

    // Navigate to a protected page — should trigger 401 → redirect to /login
    await page.goto('/dashboard')

    // Should end up on login page because both access and refresh fail
    await page.waitForURL(/\/login/, { timeout: 10000 })
    await expect(page).toHaveURL(/\/login/)
  })

  test('valid token allows normal page access', async ({ authedPage: page }) => {
    // Admin with valid token should stay on dashboard
    await page.goto('/dashboard')
    await page.waitForURL('**/dashboard')
    await expect(page).toHaveURL(/dashboard/)
  })

  test('removing token redirects to login on next navigation', async ({ authedPage: page }) => {
    await page.goto('/dashboard')
    await page.waitForURL('**/dashboard')

    // Clear auth tokens
    await page.evaluate(() => {
      localStorage.removeItem('ez-admin-access-token')
      localStorage.removeItem('ez-admin-user-info')
    })

    // Navigate to another protected page
    await page.goto('/system/roles')

    // Should redirect to login
    await page.waitForURL(/\/login/, { timeout: 10000 })
    await expect(page).toHaveURL(/\/login/)
  })
})
