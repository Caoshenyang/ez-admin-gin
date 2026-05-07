import { defineStore } from 'pinia'
import { computed, ref } from 'vue'

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
