import { test as base, expect, type Page } from '@playwright/test'

const API_BASE = process.env.E2E_API_URL ?? 'http://localhost:8080/api/v1'

type TestFixtures = {
  authedPage: Page
}

export const test = base.extend<TestFixtures>({
  authedPage: async ({ page }, use) => {
    await loginViaApi(page, 'admin', 'Admin@123456')
    await page.goto('/')
    await page.waitForURL('**/dashboard')
    await use(page)
  },
})

export { expect }

export async function loginViaApi(page: Page, username: string, password: string) {
  const response = await page.request.post(`${API_BASE}/auth/login`, {
    data: { username, password },
  })
  expect(response.ok()).toBeTruthy()
  const body = await response.json()
  const accessToken = body.data?.access_token
  expect(accessToken).toBeTruthy()

  await page.evaluate((token: string) => {
    localStorage.setItem('ez-admin-access-token', token)
    localStorage.setItem('ez-admin-token-type', 'Bearer')
    const expiresIn = new Date(Date.now() + 2 * 60 * 60 * 1000).toISOString()
    localStorage.setItem(
      'ez-admin-user-info',
      JSON.stringify({
        userId: 1,
        username: 'admin',
        nickname: 'Admin',
        expiresAt: expiresIn,
      }),
    )
  }, accessToken)
}

export async function clearAuth(page: Page) {
  await page.evaluate(() => {
    localStorage.clear()
    sessionStorage.clear()
  })
}
