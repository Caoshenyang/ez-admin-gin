import { test, expect, getAdminToken, API_BASE } from '../fixtures'

const uid = () => Date.now().toString(36)

async function createUserWithoutMenuPermission(page: import('@playwright/test').Page) {
  const token = await getAdminToken()
  const id = uid()

  // Create role with NO role management menu permission
  const roleResp = await page.request.post(`${API_BASE}/system/roles`, {
    headers: { Authorization: `Bearer ${token}` },
    data: { code: `e2e_nomenu_${id}`, name: `无菜单测试_${id}`, status: 1, sort: 0 },
  })
  expect(roleResp.ok()).toBeTruthy()
  const roleBody = await roleResp.json()
  const roleId = roleBody.data?.id

  // Only assign system directory (100) + dashboard menu, NOT role menu (202)
  await page.request.post(`${API_BASE}/system/roles/${roleId}/menus`, {
    headers: { Authorization: `Bearer ${token}` },
    data: { menu_ids: [100, 200] },
  })

  // Assign minimal API permissions
  await page.request.post(`${API_BASE}/system/roles/${roleId}/permissions`, {
    headers: { Authorization: `Bearer ${token}` },
    data: { permissions: [{ method: 'GET', path: '/api/v1/system/health' }] },
  })

  // Create user
  const username = `e2e_nomenu_${id}`
  const userResp = await page.request.post(`${API_BASE}/system/users`, {
    headers: { Authorization: `Bearer ${token}` },
    data: { username, password: 'E2eNoMenu@123', nickname: `无菜单_${id}`, status: 1 },
  })
  expect(userResp.ok()).toBeTruthy()
  const userBody = await userResp.json()
  const userId = userBody.data?.id

  // Assign role
  const assignResp = await page.request.post(`${API_BASE}/system/users/${userId}/roles`, {
    headers: { Authorization: `Bearer ${token}` },
    data: { role_ids: [roleId] },
  })
  expect(assignResp.ok()).toBeTruthy()

  return { username, password: 'E2eNoMenu@123' }
}

async function loginAsRestricted(page: import('@playwright/test').Page, username: string, password: string) {
  const response = await page.request.post(`${API_BASE}/auth/login`, {
    data: { username, password },
  })
  expect(response.ok()).toBeTruthy()
  const body = await response.json()
  const token = body.data?.access_token

  await page.goto('/')
  await page.evaluate(
    (data: { token: string; username: string }) => {
      localStorage.setItem('ez-admin-access-token', data.token)
      localStorage.setItem('ez-admin-token-type', 'Bearer')
      localStorage.setItem(
        'ez-admin-user-info',
        JSON.stringify({
          userId: 0,
          username: data.username,
          nickname: data.username,
          expiresAt: new Date(Date.now() + 2 * 60 * 60 * 1000).toISOString(),
        }),
      )
    },
    { token, username },
  )
}

test.describe('No Permission Page', () => {
  test('restricted user does not see unauthorized menu in sidebar', async ({ page }) => {
    const { username, password } = await createUserWithoutMenuPermission(page)
    await loginAsRestricted(page, username, password)
    await page.goto('/dashboard')
    await page.waitForURL('**/dashboard')

    // Role management menu should NOT appear in sidebar
    await expect(page.locator('.n-menu-item').getByText('角色管理')).not.toBeVisible()
  })

  test('restricted user navigating to unauthorized route is redirected to dashboard', async ({ page }) => {
    const { username, password } = await createUserWithoutMenuPermission(page)
    await loginAsRestricted(page, username, password)

    // Try to navigate directly to role management page
    await page.goto('/system/roles')

    // Should be redirected to dashboard (fallback route)
    await page.waitForURL('**/dashboard', { timeout: 5000 })
    await expect(page).toHaveURL(/dashboard/)
  })

  test('API request without permission shows error message', async ({ page }) => {
    const { username, password } = await createUserWithoutMenuPermission(page)
    await loginAsRestricted(page, username, password)
    await page.goto('/dashboard')
    await page.waitForURL('**/dashboard')

    // Make an API call without permission via page context
    const response = await page.request.get(`${API_BASE}/system/roles`, {
      headers: { Authorization: `Bearer ${await page.evaluate(() => localStorage.getItem('ez-admin-access-token') || '')}` },
    })

    // Should get 403 Forbidden
    expect(response.status()).toBe(403)
  })
})
