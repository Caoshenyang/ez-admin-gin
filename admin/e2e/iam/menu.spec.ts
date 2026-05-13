import { test, expect } from '../fixtures'

const uid = () => Date.now().toString(36)
test.describe('Menu Permission', () => {
  test.beforeEach(async ({ authedPage: page }) => {
    await page.goto('/system/menus')
    await page.waitForURL('**/system/menus')
    await expect(page.locator('.menu-table').getByText('系统管理')).toBeVisible()
  })

  test('displays menu management page with correct header', async ({ authedPage: page }) => {
    await expect(page.locator('h1', { hasText: '菜单管理' })).toBeVisible()
    await expect(page.getByText('维护侧边栏目录、页面菜单和页面内按钮权限')).toBeVisible()
  })

  test('shows create root directory button for admin', async ({ authedPage: page }) => {
    await expect(page.getByRole('button', { name: '+ 新增根目录' })).toBeVisible()
  })

  test('displays menu table with seed data columns', async ({ authedPage: page }) => {
    const table = page.locator('.menu-table')
    await expect(table).toBeVisible()
    await expect(table.getByText('菜单名称')).toBeVisible()
    await expect(table.getByText('路由')).toBeVisible()
    await expect(table.getByText('权限标识')).toBeVisible()
    await expect(table.locator('.n-data-table-th__title', { hasText: '状态' })).toBeVisible()
    await expect(table.locator('.n-data-table-th__title', { hasText: '操作' })).toBeVisible()
  })

  test('shows seed menu items in table', async ({ authedPage: page }) => {
    await expect(page.locator('.menu-table').getByText('系统管理')).toBeVisible()
    await expect(page.locator('.menu-table').getByText('用户管理')).toBeVisible()
    await expect(page.locator('.menu-table').getByText('菜单管理')).toBeVisible()
  })

  test('creates a new root directory menu', async ({ authedPage: page }) => {
    const id = uid()
    const menuName = `E2E目录_${id}`
    const menuCode = `e2e:create:${id}`

    await page.getByRole('button', { name: '+ 新增根目录' }).click()
    await expect(page.locator('.n-modal')).toBeVisible()

    await page.getByPlaceholder('请输入菜单名称').fill(menuName)
    await page.getByPlaceholder('system:example:list').fill(menuCode)
    await page.getByPlaceholder('/system/example').fill(`/e2e/${id}`)

    await page.locator('.n-modal').getByRole('button', { name: '保存' }).click()

    await expect(page.locator('.n-alert').getByText(/成功/)).toBeVisible()
    await expect(page.locator('.menu-table').getByText(menuName)).toBeVisible()
  })

  test('opens edit modal with correct data', async ({ authedPage: page }) => {
    await page.getByText('展开全部').click()
    // Click edit on first menu row (系统管理)
    await page.locator('.menu-table').getByRole('button', { name: '编辑' }).first().click()
    await expect(page.locator('.n-modal')).toBeVisible()
    // Verify edit modal shows existing data
    const nameInput = page.locator('.n-modal').getByPlaceholder('请输入菜单名称')
    await expect(nameInput).toHaveValue(/./) // Has some value
    // Verify code field is disabled in edit mode
    await expect(page.locator('.n-modal').getByPlaceholder('system:example:list')).toBeDisabled()
  })

  test('deletes a menu item', async ({ authedPage: page }) => {
    const id = uid()
    const menuName = `待删除_${id}`
    const menuCode = `e2e:del:${id}`

    await page.getByRole('button', { name: '+ 新增根目录' }).click()
    await expect(page.locator('.n-modal')).toBeVisible()
    await page.getByPlaceholder('请输入菜单名称').fill(menuName)
    await page.getByPlaceholder('system:example:list').fill(menuCode)
    await page.locator('.n-modal').getByRole('button', { name: '保存' }).click()
    await expect(page.locator('.n-alert').getByText(/成功/)).toBeVisible()

    const row = page.locator('.menu-table').locator('tr', { hasText: menuName })
    await row.getByRole('button', { name: '删除' }).click()
    await page.getByRole('button', { name: '确认' }).click()

    await expect(page.locator('.n-alert').getByText(/成功/)).toBeVisible()
    await expect(page.locator('.menu-table').getByText(menuName)).not.toBeVisible()
  })

  test('action buttons are visible for admin user', async ({ authedPage: page }) => {
    await page.getByText('展开全部').click()
    await expect(page.locator('.menu-table').getByRole('button', { name: '编辑' }).first()).toBeVisible()
    await expect(page.locator('.menu-table').getByRole('button', { name: '删除' }).first()).toBeVisible()
  })
})
