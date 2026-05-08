import { computed, onMounted, ref } from 'vue'

import { formatTime } from '@/utils/format'
import { getSystemHealth } from '../api/health'
import type { SystemHealthData } from '../types/health'
import {
  formatHealthStatusLabel,
  healthDependencyCards,
  healthEndpointCards,
  healthErrorMessage,
  healthStatusTagType,
  healthStatusText,
} from './health-page.utils'

export function useHealthPage() {
  const loading = ref(false)
  const errorMessage = ref('')
  const health = ref<SystemHealthData | null>(null)
  const lastCheckedAt = ref('')

  const dependencyCards = computed(() => healthDependencyCards(health.value))
  const endpointCards = healthEndpointCards

  const envTagType = computed(() => {
    return health.value?.env === 'prod' ? 'success' : 'warning'
  })

  const lastCheckedLabel = computed(() => {
    return lastCheckedAt.value ? formatTime(lastCheckedAt.value) : '尚未检查'
  })

  async function loadHealth() {
    loading.value = true
    errorMessage.value = ''

    try {
      health.value = await getSystemHealth()
      lastCheckedAt.value = new Date().toISOString()
    } catch (error) {
      errorMessage.value = healthErrorMessage(error)
    } finally {
      loading.value = false
    }
  }

  function formatStatusLabel(value?: string) {
    return formatHealthStatusLabel(value)
  }

  function getStatusTagType(value?: string) {
    return healthStatusTagType(value)
  }

  function getStatusText(value?: string) {
    return healthStatusText(value, loading.value)
  }

  onMounted(() => {
    void loadHealth()
  })

  return {
    dependencyCards,
    endpointCards,
    envTagType,
    errorMessage,
    formatStatusLabel,
    getStatusTagType,
    getStatusText,
    health,
    lastCheckedLabel,
    loadHealth,
    loading,
  }
}
