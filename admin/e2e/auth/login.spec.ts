import { test, expect, clearAuth } from '../fixtures'

test.describe('Login Flow', () => {
  test('redirects to login page when not authenticated', async ({ page }) => {
    await page.goto('/login')
    await clearAuth(page)
    await page.goto('/dashboard')
    await page.waitForURL(/\/login/)
    await expect(page.getByText('登录控制台')).toBeVisible()
  })

  test('shows error on wrong password', async ({ page }) => {
    await page.goto('/login')
    await clearAuth(page)
    await page.goto('/login')
    await page.getByPlaceholder('请输入用户名').fill('admin')
    await page.getByPlaceholder('请输入密码').fill('WrongPassword123')
    await page.getByPlaceholder('验证码').fill('abcd')
    await page.getByRole('button', { name: '登录' }).click()
    await expect(page.getByText(/密码错误|用户名或密码/)).toBeVisible()
  })

  test('logs in successfully and redirects to dashboard', async ({ page }) => {
    await page.goto('/login')
    await clearAuth(page)
    await page.goto('/login')
    await page.getByPlaceholder('请输入用户名').fill('admin')
    await page.getByPlaceholder('请输入密码').fill('Admin@123456')
    await page.getByPlaceholder('验证码').fill('abcd')
    await page.getByRole('button', { name: '登录' }).click()
    await page.waitForURL('**/dashboard', { timeout: 15_000 })
    await expect(page).toHaveURL(/\/dashboard/)
  })

  test('already logged in user is redirected to dashboard from /login', async ({ page }) => {
    const API_BASE = process.env.E2E_API_URL ?? 'http://localhost:8080/api/v1'
    const response = await page.request.post(`${API_BASE}/auth/login`, {
      data: { username: 'admin', password: 'Admin@123456' },
    })
    const body = await response.json()
    const token = body.data?.access_token

    await page.goto('/login')
    await page.evaluate((t: string) => {
      localStorage.setItem('ez-admin-access-token', t)
      localStorage.setItem('ez-admin-token-type', 'Bearer')
      localStorage.setItem(
        'ez-admin-user-info',
        JSON.stringify({
          userId: 1,
          username: 'admin',
          nickname: 'Admin',
          expiresAt: new Date(Date.now() + 2 * 60 * 60 * 1000).toISOString(),
        }),
      )
    }, token)

    await page.goto('/login')
    await page.waitForURL('**/dashboard')
    await expect(page).toHaveURL(/\/dashboard/)
  })

  test('form validation shows errors for empty fields', async ({ page }) => {
    await page.goto('/login')
    await clearAuth(page)
    await page.goto('/login')
    await page.getByPlaceholder('请输入用户名').clear()
    await page.getByPlaceholder('请输入用户名').blur()
    await expect(page.locator('.n-form-item-feedback__line')).toContainText('请输入用户名')
  })
})
