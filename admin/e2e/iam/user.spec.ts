import { test, expect } from '../fixtures'

const API_BASE = process.env.E2E_API_URL ?? 'http://localhost:8080/api/v1'
const uid = () => Date.now().toString(36)

test.describe('User Management', () => {
  test.beforeEach(async ({ authedPage: page }) => {
    await page.goto('/system/users')
    await page.waitForURL('**/system/users')
    // Wait for table data to load
    await expect(page.locator('.n-data-table').getByText('admin').first()).toBeVisible()
  })

  test('displays user management page with correct header', async ({ authedPage: page }) => {
    await expect(page.locator('h1', { hasText: '用户管理' })).toBeVisible()
    await expect(page.getByText('维护后台账号、启停状态和角色绑定')).toBeVisible()
  })

  test('shows create user button for admin', async ({ authedPage: page }) => {
    await expect(page.getByRole('button', { name: '+ 新增用户' })).toBeVisible()
  })

  test('displays user table with columns', async ({ authedPage: page }) => {
    const table = page.locator('.n-data-table')
    await expect(table).toBeVisible()
    await expect(table.locator('.n-data-table-th__title', { hasText: '用户' })).toBeVisible()
    await expect(table.locator('.n-data-table-th__title', { hasText: '状态' })).toBeVisible()
    await expect(table.getByText('创建时间')).toBeVisible()
    await expect(table.getByText('操作')).toBeVisible()
  })

  test('shows admin user in the list', async ({ authedPage: page }) => {
    await expect(page.locator('.n-data-table').getByText('admin').first()).toBeVisible()
  })

  test('creates a new user', async ({ authedPage: page }) => {
    const id = uid()
    const username = `e2e_create_${id}`

    await page.getByRole('button', { name: '+ 新增用户' }).click()
    await expect(page.locator('.n-modal')).toBeVisible()
    // Modal heading is inside .n-modal
    await expect(page.locator('.n-modal').getByText('新增用户')).toBeVisible()

    await page.getByPlaceholder('请输入用户名').fill(username)
    await page.getByPlaceholder('至少 8 位').fill('E2eTest@123')
    await page.getByPlaceholder('请输入昵称').fill(`E2E用户_${id}`)

    await page.locator('.n-modal').getByRole('button', { name: '保存' }).click()

    await expect(page.locator('.n-alert').getByText(/成功/)).toBeVisible()
    await expect(page.locator('.n-data-table').getByText(username)).toBeVisible()
  })

  test('opens edit modal with existing data', async ({ authedPage: page }) => {
    const row = page.locator('.n-data-table').locator('tr', { hasText: 'admin' }).first()
    await row.getByRole('button', { name: '编辑' }).click()
    await expect(page.locator('.n-modal')).toBeVisible()
    // Edit mode: username field is hidden, only nickname is shown
    await expect(page.locator('.n-modal').getByText('编辑用户')).toBeVisible()
    const nicknameInput = page.locator('.n-modal').getByPlaceholder('请输入昵称')
    await expect(nicknameInput).toHaveValue(/./)
  })

  test('toggles user status via API', async ({ authedPage: page }) => {
    const id = uid()
    const username = `e2e_status_${id}`

    const token = await getAdminToken(page)
    const resp = await page.request.post(`${API_BASE}/system/users`, {
      headers: { Authorization: `Bearer ${token}` },
      data: {
        username, password: 'E2eStatus@123', nickname: `状态测试_${id}`, status: 1,
      },
    })
    expect(resp.ok()).toBeTruthy()

    // Toggle via API directly to avoid NInput fill() issues
    const createBody = await resp.json()
    const userId = createBody.data?.id
    const toggleResp = await page.request.post(`${API_BASE}/system/users/${userId}/status`, {
      headers: { Authorization: `Bearer ${token}` },
      data: { status: 2 },
    })
    expect(toggleResp.ok()).toBeTruthy()
  })

  test('action buttons visible for admin', async ({ authedPage: page }) => {
    const table = page.locator('.n-data-table')
    await expect(table.getByRole('button', { name: '编辑' }).first()).toBeVisible()
    await expect(table.getByRole('button', { name: /启用|禁用/ }).first()).toBeVisible()
  })
})

async function getAdminToken(page: import('@playwright/test').Page): Promise<string> {
  const response = await page.request.post(`${API_BASE}/auth/login`, {
    data: { username: 'admin', password: 'Admin@123456' },
  })
  const body = await response.json()
  return body.data?.access_token
}
