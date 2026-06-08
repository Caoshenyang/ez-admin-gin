import { defineStore } from 'pinia'
import { computed, ref } from 'vue'

// useAdminShellStore 管理后台布局的侧栏菜单高亮、展开状态和工作标签页。
export interface WorkTab {
  affix?: boolean
  fullPath: string
  icon?: string
  keepAlive?: boolean
  key: string
  name?: string
  params?: Record<string, unknown>
  path: string
  query?: Record<string, unknown>
  title: string
  closable: boolean
}

const STORAGE_KEY = 'ez-admin-work-tabs'

const dashboardTab: WorkTab = {
  affix: true,
  fullPath: '/dashboard',
  key: 'dashboard',
  path: '/dashboard',
  title: '工作台',
  closable: false,
}

function normalizeTab(tab: WorkTab): WorkTab {
  const fullPath = tab.fullPath || tab.path
  return {
    ...tab,
    fullPath,
    key: tab.key || fullPath,
    path: tab.path || fullPath,
    closable: tab.affix ? false : tab.closable,
  }
}

function safeParseTabs() {
  try {
    const raw = localStorage.getItem(STORAGE_KEY)
    if (!raw) return []
    const parsed = JSON.parse(raw) as WorkTab[]
    if (!Array.isArray(parsed)) return []
    return parsed
      .filter((tab) => tab && typeof tab.fullPath === 'string' && typeof tab.title === 'string')
      .map(normalizeTab)
  } catch {
    return []
  }
}

export const useAdminShellStore = defineStore('admin-shell', () => {
  const activeMenuKey = ref('')
  const expandedMenuKeys = ref<string[]>([])
  const openTabs = ref<WorkTab[]>([dashboardTab])
  const routeRefreshNonce = ref<Record<string, number>>({})
  const sidebarCollapsed = ref(false)

  const tabMap = computed(() => {
    return new Map(openTabs.value.map((tab) => [tab.fullPath, tab]))
  })

  function persistTabs() {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(openTabs.value.map(normalizeTab)))
  }

  function hydrateTabs() {
    const storedTabs = safeParseTabs()
    const mergedTabs = [
      dashboardTab,
      ...storedTabs.filter((tab) => tab.fullPath !== dashboardTab.fullPath),
    ]
    openTabs.value = mergedTabs.map(normalizeTab)
    persistTabs()
  }

  function setActiveMenuKey(key: string) {
    activeMenuKey.value = key
  }

  function setExpandedMenuKeys(keys: string[]) {
    expandedMenuKeys.value = keys
  }

  function ensureExpandedMenuKeys(keys: string[]) {
    if (keys.length === 0) {
      return
    }

    // 合并而非替换，防止路由切换时折叠用户手动展开的菜单。
    const merged = new Set(expandedMenuKeys.value)
    for (const key of keys) {
      merged.add(key)
    }
    expandedMenuKeys.value = Array.from(merged)
  }

  function ensureTab(tab: WorkTab) {
    const normalized = normalizeTab(tab)
    const existing = tabMap.value.get(normalized.fullPath)
    if (existing) {
      existing.title = normalized.title
      existing.path = normalized.path
      existing.name = normalized.name
      existing.icon = normalized.icon
      persistTabs()
      return
    }
    openTabs.value.push(normalized)
    persistTabs()
  }

  function findFallbackPath(closedFullPath: string, activeFullPath: string) {
    const index = openTabs.value.findIndex((tab) => tab.fullPath === closedFullPath)
    if (closedFullPath !== activeFullPath) {
      return activeFullPath
    }
    return (
      openTabs.value[index + 1]?.fullPath ??
      openTabs.value[index - 1]?.fullPath ??
      dashboardTab.fullPath
    )
  }

  function closeTab(fullPath: string, activeFullPath = fullPath) {
    const target = tabMap.value.get(fullPath)
    if (!target?.closable) {
      return activeFullPath
    }

    const fallback = findFallbackPath(fullPath, activeFullPath)
    openTabs.value = openTabs.value.filter((tab) => tab.fullPath !== fullPath)
    // 至少保留工作台标签页，防止标签栏为空。
    if (openTabs.value.length === 0) {
      openTabs.value = [dashboardTab]
    }
    persistTabs()
    return fallback
  }

  function closeOtherTabs(currentFullPath: string) {
    const current = tabMap.value.get(currentFullPath)
    openTabs.value = openTabs.value.filter((tab) => tab.affix)
    if (current && !current.affix) {
      openTabs.value.push(current)
    }
    persistTabs()
  }

  function closeLeftTabs(currentFullPath: string) {
    const currentIndex = openTabs.value.findIndex((tab) => tab.fullPath === currentFullPath)
    if (currentIndex < 0) return
    openTabs.value = openTabs.value.filter((tab, index) => tab.affix || index >= currentIndex)
    persistTabs()
  }

  function closeRightTabs(currentFullPath: string) {
    const currentIndex = openTabs.value.findIndex((tab) => tab.fullPath === currentFullPath)
    if (currentIndex < 0) return
    openTabs.value = openTabs.value.filter((tab, index) => tab.affix || index <= currentIndex)
    persistTabs()
  }

  function closeAllTabs() {
    openTabs.value = openTabs.value.filter((tab) => tab.affix)
    if (!openTabs.value.some((tab) => tab.fullPath === dashboardTab.fullPath)) {
      openTabs.value.unshift(dashboardTab)
    }
    persistTabs()
  }

  function pinCurrentTab(currentFullPath: string) {
    const tab = tabMap.value.get(currentFullPath)
    if (!tab || tab.fullPath === dashboardTab.fullPath) return
    tab.affix = !tab.affix
    tab.closable = !tab.affix
    persistTabs()
  }

  // refreshRoute 通过递增 nonce 强制路由组件重新挂载，实现页面刷新。
  function refreshRoute(path: string) {
    routeRefreshNonce.value = {
      ...routeRefreshNonce.value,
      [path]: (routeRefreshNonce.value[path] ?? 0) + 1,
    }
  }

  function getRouteViewKey(path: string) {
    return `${path}::${routeRefreshNonce.value[path] ?? 0}`
  }

  function toggleSidebar() {
    sidebarCollapsed.value = !sidebarCollapsed.value
  }

  function reset() {
    activeMenuKey.value = ''
    expandedMenuKeys.value = []
    openTabs.value = [dashboardTab]
    routeRefreshNonce.value = {}
    sidebarCollapsed.value = false
    persistTabs()
  }

  return {
    activeMenuKey,
    expandedMenuKeys,
    openTabs,
    sidebarCollapsed,
    ensureExpandedMenuKeys,
    ensureTab,
    closeLeftTabs,
    closeAllTabs,
    closeOtherTabs,
    closeRightTabs,
    closeTab,
    getRouteViewKey,
    hydrateTabs,
    refreshRoute,
    persistTabs,
    pinCurrentTab,
    reset,
    setActiveMenuKey,
    setExpandedMenuKeys,
    toggleSidebar,
  }
})
