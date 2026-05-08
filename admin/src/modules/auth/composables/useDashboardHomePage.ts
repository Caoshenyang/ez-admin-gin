import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'

import { getDashboardSummary } from '../api/dashboard'
import type { DashboardData } from '../types/dashboard'
import { authMenus } from '@/router/dynamic-menu'
import { displayText } from '@/utils/format'
import {
  dashboardCurrentDateLabel,
  dashboardCurrentUserLabel,
  dashboardErrorMessage,
  dashboardFocusFacts,
  dashboardHealthItems,
  dashboardHealthTagType,
  dashboardHeroIcon,
  dashboardHeroStatusText,
  dashboardIsHealthy,
  dashboardLoginStatusLabel,
  dashboardLoginStatusTagType,
  dashboardMetricCards,
  dashboardOperationStatusTagType,
  dashboardQuickLinks,
  findMenuPathByTitle,
  flattenPageMenus,
  formatDashboardDateTime,
  formatDashboardRoutePath,
  formatMetricValue,
} from './dashboard-page.utils'

export function useDashboardHomePage() {
  const router = useRouter()
  const loading = ref(false)
  const errorMessage = ref('')
  const dashboard = ref<DashboardData | null>(null)
  const refreshedAt = ref('')

  const currentDateLabel = computed(() => dashboardCurrentDateLabel())
  const currentUserLabel = computed(() => dashboardCurrentUserLabel(dashboard.value))

  const visiblePageMenus = computed(() => flattenPageMenus(authMenus.value))
  const visiblePageTotal = computed(() => visiblePageMenus.value.length + 1)

  const healthPath = computed(() => findMenuPathByTitle(visiblePageMenus.value, '系统状态'))
  const userManagePath = computed(() => findMenuPathByTitle(visiblePageMenus.value, '用户管理'))

  const isHealthy = computed(() => dashboardIsHealthy(dashboard.value))
  const heroStatusText = computed(() => dashboardHeroStatusText(dashboard.value, loading.value))
  const metricCards = computed(() => dashboardMetricCards(dashboard.value, visiblePageTotal.value))
  const healthItems = computed(() => dashboardHealthItems(dashboard.value))
  const quickLinks = computed(() => dashboardQuickLinks(visiblePageMenus.value))
  const recentOperations = computed(() => dashboard.value?.recent_operations ?? [])
  const recentLogins = computed(() => dashboard.value?.recent_logins ?? [])
  const latestNotices = computed(() => dashboard.value?.latest_notices ?? [])

  const refreshedLabel = computed(() => {
    return refreshedAt.value ? formatDashboardDateTime(refreshedAt.value) : '尚未同步'
  })

  const focusFacts = computed(() => {
    return dashboardFocusFacts(dashboard.value, visiblePageTotal.value, refreshedLabel.value)
  })

  async function loadDashboard() {
    loading.value = true
    errorMessage.value = ''

    try {
      dashboard.value = await getDashboardSummary()
      refreshedAt.value = new Date().toISOString()
    } catch (error) {
      errorMessage.value = dashboardErrorMessage(error)
    } finally {
      loading.value = false
    }
  }

  function navigateTo(path: string) {
    if (!path) {
      return
    }

    void router.push(path)
  }

  onMounted(() => {
    void loadDashboard()
  })

  return {
    currentDateLabel,
    currentUserLabel,
    dashboard,
    displayText,
    errorMessage,
    focusFacts,
    formatDashboardDateTime,
    formatDashboardRoutePath,
    formatMetricValue,
    healthItems,
    healthPath,
    heroIcon: dashboardHeroIcon,
    heroStatusText,
    isHealthy,
    latestNotices,
    loading,
    loadDashboard,
    navigateTo,
    quickLinks,
    recentLogins,
    recentOperations,
    refreshedLabel,
    userManagePath,
    visiblePageTotal,
    getHealthTagType: dashboardHealthTagType,
    getLoginStatusLabel: dashboardLoginStatusLabel,
    getLoginStatusTagType: dashboardLoginStatusTagType,
    getStatusTagType: dashboardOperationStatusTagType,
    metricCards,
  }
}
