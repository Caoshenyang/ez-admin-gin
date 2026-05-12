import { test as base, expect, type Page } from '@playwright/test'

const API_BASE = process.env.E2E_API_URL ?? 'http://localhost:8080/api/v1'

type TestFixtures = {
  authedPage: Page
}

let cachedToken: string | null = null

async function fetchAdminToken(): Promise<string> {
  if (cachedToken) return cachedToken

  const response = await fetch(`${API_BASE}/auth/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ username: 'admin', password: 'Admin@123456' }),
  })
  const body = await response.json()
  if (!response.ok || !body.data?.access_token) {
    throw new Error(`Login failed: ${response.status} ${JSON.stringify(body)}`)
  }
  cachedToken = body.data.access_token
  return cachedToken
}

export const test = base.extend<TestFixtures>({
  authedPage: async ({ page }, use) => {
    const accessToken = await fetchAdminToken()
    await page.goto('/')
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

  // Navigate to app origin first so localStorage is accessible
  await page.goto('/')
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
