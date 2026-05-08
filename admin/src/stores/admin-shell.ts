import { defineStore } from 'pinia'
import { computed, ref } from 'vue'

// useAdminShellStore 管理后台布局的侧栏菜单高亮、展开状态和工作标签页。
export interface WorkTab {
  title: string
  to: string
  closable: boolean
}

const dashboardTab: WorkTab = {
  title: '工作台',
  to: '/dashboard',
  closable: false,
}

export const useAdminShellStore = defineStore('admin-shell', () => {
  const activeMenuKey = ref('')
  const expandedMenuKeys = ref<string[]>([])
  const openTabs = ref<WorkTab[]>([dashboardTab])
  const routeRefreshNonce = ref<Record<string, number>>({})

  const tabMap = computed(() => {
    return new Map(openTabs.value.map((tab) => [tab.to, tab]))
  })

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
    if (tabMap.value.has(tab.to)) {
      return
    }
    openTabs.value.push(tab)
  }

  function closeTab(path: string) {
    openTabs.value = openTabs.value.filter((tab) => tab.to !== path)
    // 至少保留工作台标签页，防止标签栏为空。
    if (openTabs.value.length === 0) {
      openTabs.value = [dashboardTab]
    }
  }

  function closeOtherTabs(currentPath: string) {
    const current = tabMap.value.get(currentPath)
    openTabs.value = [dashboardTab]
    if (current && current.to !== dashboardTab.to) {
      openTabs.value.push(current)
    }
  }

  function closeAllTabs() {
    openTabs.value = [dashboardTab]
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

  function reset() {
    activeMenuKey.value = ''
    expandedMenuKeys.value = []
    openTabs.value = [dashboardTab]
    routeRefreshNonce.value = {}
  }

  return {
    activeMenuKey,
    expandedMenuKeys,
    openTabs,
    ensureExpandedMenuKeys,
    ensureTab,
    closeAllTabs,
    closeOtherTabs,
    closeTab,
    getRouteViewKey,
    refreshRoute,
    reset,
    setActiveMenuKey,
    setExpandedMenuKeys,
  }
})
