import { test, expect, getAdminToken, API_BASE } from '../fixtures'

const uid = () => Date.now().toString(36)

test.describe('Role Authorization', () => {
  test.beforeEach(async ({ authedPage: page }) => {
    await page.goto('/system/roles')
    await page.waitForURL('**/system/roles')
    await expect(page.locator('h1', { hasText: '角色权限' })).toBeVisible()
    await expect(page.locator('.role-card').first()).toBeVisible()
  })

  test('displays role page with correct header', async ({ authedPage: page }) => {
    await expect(page.locator('h1', { hasText: '角色权限' })).toBeVisible()
    await expect(page.getByText('维护角色本身，以及角色拥有的菜单、按钮和接口权限。')).toBeVisible()
  })

  test('shows super_admin role in role list', async ({ authedPage: page }) => {
    await expect(page.locator('.role-card').getByText('超级管理员')).toBeVisible()
    await expect(page.locator('.role-card').getByText('super_admin')).toBeVisible()
  })

  test('super admin role shows protected tag', async ({ authedPage: page }) => {
    await page.locator('.role-card').getByText('超级管理员').click()
    await expect(page.getByText('受保护角色')).toBeVisible()
  })

  test('creates a new role', async ({ authedPage: page }) => {
    const id = uid()
    const roleCode = `e2e_role_${id}`
    const roleName = `E2E角色_${id}`

    await page.getByRole('button', { name: '+ 新增角色' }).click()
    await expect(page.locator('.n-modal')).toBeVisible()
    await expect(page.locator('.n-modal').getByText('新增角色')).toBeVisible()

    await page.getByPlaceholder('demo_operator').fill(roleCode)
    await page.getByPlaceholder('请输入角色名称').fill(roleName)

    await page.locator('.n-modal').getByRole('button', { name: '保存' }).click()

    await expect(page.locator('.n-alert').getByText(/成功/)).toBeVisible()
    await expect(page.locator('.role-card').getByText(roleName)).toBeVisible()
  })

  test('opens edit modal with existing data', async ({ authedPage: page }) => {
    const token = await getAdminToken()
    const id = uid()
    const roleCode = `e2e_edit_${id}`
    const roleName = `编辑测试_${id}`

    await page.request.post(`${API_BASE}/system/roles`, {
      headers: { Authorization: `Bearer ${token}` },
      data: { code: roleCode, name: roleName, status: 1, sort: 0 },
    })

    await page.reload()
    await page.waitForURL('**/system/roles')
    await expect(page.locator('.role-card').getByText(roleName)).toBeVisible()

    const card = page.locator('.role-card', { hasText: roleName })
    await card.getByRole('button', { name: '编辑' }).click()
    await expect(page.locator('.n-modal')).toBeVisible()

    const nameInput = page.locator('.n-modal').getByPlaceholder('请输入角色名称')
    await expect(nameInput).toHaveValue(roleName)
  })

  test('toggles role status via UI', async ({ authedPage: page }) => {
    const token = await getAdminToken()
    const id = uid()
    const roleCode = `e2e_status_${id}`
    const roleName = `状态测试_${id}`

    await page.request.post(`${API_BASE}/system/roles`, {
      headers: { Authorization: `Bearer ${token}` },
      data: { code: roleCode, name: roleName, status: 1, sort: 0 },
    })

    await page.reload()
    await page.waitForURL('**/system/roles')
    await expect(page.locator('.role-card').first()).toBeVisible()
    await expect(page.locator('.role-card').getByText(roleName)).toBeVisible()

    const card = page.locator('.role-card', { hasText: roleName })
    await card.getByRole('button', { name: '禁用' }).click()
    await page.getByRole('button', { name: '确认' }).click()

    await expect(page.locator('.n-alert').getByText(/已禁用|已启用/)).toBeVisible()
    await expect(card.getByRole('button', { name: '启用' })).toBeVisible()
  })

  test('permission panel shows menu tree for selected role', async ({ authedPage: page }) => {
    await page.locator('.role-card').getByText('超级管理员').click()

    await expect(page.getByText('菜单与按钮权限')).toBeVisible()
    await expect(page.getByText('菜单权限', { exact: true })).toBeVisible()
    await expect(page.getByText('按钮权限', { exact: true })).toBeVisible()
    await expect(page.getByText('接口权限', { exact: true })).toBeVisible()
  })

  test('assigns menu permissions to a role', async ({ authedPage: page }) => {
    const token = await getAdminToken()
    const id = uid()
    const roleCode = `e2e_perm_${id}`
    const roleName = `权限测试_${id}`

    const createResp = await page.request.post(`${API_BASE}/system/roles`, {
      headers: { Authorization: `Bearer ${token}` },
      data: { code: roleCode, name: roleName, status: 1, sort: 0 },
    })
    const createBody = await createResp.json()
    void createBody.data?.id

    await page.reload()
    await page.waitForURL('**/system/roles')
    await expect(page.locator('.role-card').first()).toBeVisible()

    const card = page.locator('.role-card', { hasText: roleName })
    await expect(card).toBeVisible()
    await card.click()

    await page.locator('.n-tabs-tab').filter({ hasText: '菜单权限' }).click()
    await page.getByRole('checkbox', { name: '全选' }).click()
    await page.getByRole('button', { name: '保存权限' }).click()
    await expect(page.locator('.n-alert').getByText(/已更新|成功/)).toBeVisible()
  })

  test('adds API permission row', async ({ authedPage: page }) => {
    const token = await getAdminToken()
    const id = uid()
    const roleCode = `e2e_api_${id}`
    const roleName = `API测试_${id}`

    await page.request.post(`${API_BASE}/system/roles`, {
      headers: { Authorization: `Bearer ${token}` },
      data: { code: roleCode, name: roleName, status: 1, sort: 0 },
    })

    await page.reload()
    await page.waitForURL('**/system/roles')
    await expect(page.locator('.role-card').first()).toBeVisible()

    const card = page.locator('.role-card', { hasText: roleName })
    await expect(card).toBeVisible()
    await card.click()

    await page.locator('.n-tabs-tab').filter({ hasText: '接口权限' }).click()
    await page.getByRole('button', { name: '+ 添加接口' }).click()
    await expect(page.getByPlaceholder('/api/v1/system/users')).toBeVisible()
  })
})
