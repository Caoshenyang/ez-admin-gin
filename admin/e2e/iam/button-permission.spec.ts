import { test, expect, getAdminToken, API_BASE } from '../fixtures'
import type { Page } from '@playwright/test'

const uid = () => Date.now().toString(36)

async function setupRestrictedUser(
  page: Page,
  menuIds: number[],
  permissions: Array<{ method: string; path: string }>,
) {
  const token = await getAdminToken()
  const id = uid()

  // Create role
  const roleResp = await page.request.post(`${API_BASE}/system/roles`, {
    headers: { Authorization: `Bearer ${token}` },
    data: {
      code: `e2e_btn_${id}`,
      name: `E2E按钮测试_${id}`,
      status: 1,
      sort: 0,
    },
  })
  expect(roleResp.ok()).toBeTruthy()
  const roleBody = await roleResp.json()
  const roleId = roleBody.data?.id

  // Assign menu permissions (system dir + role menu + specified buttons)
  await page.request.post(`${API_BASE}/system/roles/${roleId}/menus`, {
    headers: { Authorization: `Bearer ${token}` },
    data: { menu_ids: [100, 202, ...menuIds] },
  })

  // Assign API permissions
  await page.request.post(`${API_BASE}/system/roles/${roleId}/permissions`, {
    headers: { Authorization: `Bearer ${token}` },
    data: { permissions },
  })

  // Create user
  const username = `e2e_btn_${id}`
  const userResp = await page.request.post(`${API_BASE}/system/users`, {
    headers: { Authorization: `Bearer ${token}` },
    data: {
      username,
      password: 'E2eBtnTest@123',
      nickname: `按钮测试_${id}`,
      status: 1,
    },
  })
  expect(userResp.ok()).toBeTruthy()
  const userBody = await userResp.json()
  const userId = userBody.data?.id

  // Assign role to user
  const assignResp = await page.request.post(`${API_BASE}/system/users/${userId}/roles`, {
    headers: { Authorization: `Bearer ${token}` },
    data: { role_ids: [roleId] },
  })
  expect(assignResp.ok()).toBeTruthy()

  return { username, password: 'E2eBtnTest@123' }
}

async function loginAsRestricted(page: Page, username: string, password: string) {
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
  await page.goto('/system/roles')
  await page.waitForURL('**/system/roles')
}

test.describe('Button Permission', () => {
  test.describe('Admin user — all buttons visible', () => {
    test.beforeEach(async ({ authedPage: page }) => {
      await page.goto('/system/roles')
      await page.waitForURL('**/system/roles')
      await expect(page.locator('h1', { hasText: '角色权限' })).toBeVisible()
      await expect(page.locator('.role-card').first()).toBeVisible()
    })

    test('admin sees create role button', async ({ authedPage: page }) => {
      await expect(page.getByRole('button', { name: '+ 新增角色' })).toBeVisible()
    })

    test('admin sees edit button on role cards', async ({ authedPage: page }) => {
      await expect(page.locator('.role-card').getByText('编辑').first()).toBeVisible()
    })

    test('admin sees status toggle on non-super-admin roles', async ({ authedPage: page }) => {
      const toggleButtons = page.locator('.role-card').getByRole('button', { name: /启用|禁用/ })
      if ((await toggleButtons.count()) > 0) {
        await expect(toggleButtons.first()).toBeVisible()
      }
    })

    test('admin sees save permission button', async ({ authedPage: page }) => {
      await expect(page.getByRole('button', { name: '保存权限' })).toBeVisible()
    })
  })

  test.describe('Restricted user — limited buttons visible', () => {
    test('user without create permission does not see create button', async ({ page }) => {
      const { username, password } = await setupRestrictedUser(
        page,
        [1020, 1022],
        [{ method: 'GET', path: '/api/v1/system/roles' }],
      )

      await loginAsRestricted(page, username, password)

      await expect(page.getByRole('button', { name: '+ 新增角色' })).not.toBeVisible()
      await expect(page.locator('.role-card').getByText('编辑').first()).toBeVisible()
    })

    test('user without update permission does not see edit button', async ({ page }) => {
      const { username, password } = await setupRestrictedUser(
        page,
        [1020, 1021],
        [{ method: 'GET', path: '/api/v1/system/roles' }],
      )

      await loginAsRestricted(page, username, password)

      await expect(page.locator('.role-card').getByText('编辑')).toHaveCount(0)
      await expect(page.getByRole('button', { name: '+ 新增角色' })).toBeVisible()
    })

    test('user without status permission does not see status toggle', async ({ page }) => {
      const { username, password } = await setupRestrictedUser(
        page,
        [1020, 1022],
        [{ method: 'GET', path: '/api/v1/system/roles' }],
      )

      await loginAsRestricted(page, username, password)

      await expect(page.locator('.role-card').getByRole('button', { name: /启用|禁用/ })).toHaveCount(0)
      await expect(page.locator('.role-card').getByText('编辑').first()).toBeVisible()
    })
  })
})
