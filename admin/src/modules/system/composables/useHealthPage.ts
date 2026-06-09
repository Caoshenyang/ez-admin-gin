import { computed, onMounted, ref } from 'vue'

import { formatTime } from '@/utils/format'
import { getSystemHealth } from '../api/health'
import type { SystemHealthData } from '../types/health'
import {
  healthCheckItems,
  healthDependencyCards,
  healthEndpointCards,
  healthErrorMessage,
  healthIsReady,
  healthOverview,
  healthReadinessScore,
  healthSignalIcons,
  healthSummaryCards,
} from './health-page.utils'

export function useHealthPage() {
  const loading = ref(false)
  const errorMessage = ref('')
  const health = ref<SystemHealthData | null>(null)
  const lastCheckedAt = ref('')

  const checkItems = computed(() => healthCheckItems(health.value, loading.value))
  const dependencyCards = computed(() => healthDependencyCards(health.value, loading.value))
  const endpointCards = healthEndpointCards
  const isHealthy = computed(() => healthIsReady(health.value))
  const overview = computed(() => healthOverview(health.value, loading.value, errorMessage.value))
  const readinessScore = computed(() => healthReadinessScore(health.value, loading.value))
  const signalIcons = computed(() => healthSignalIcons(health.value, loading.value))

  const lastCheckedLabel = computed(() => {
    return lastCheckedAt.value ? formatTime(lastCheckedAt.value) : '尚未检查'
  })

  const summaryCards = computed(() => {
    return healthSummaryCards(health.value, lastCheckedLabel.value, loading.value)
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

  onMounted(() => {
    void loadHealth()
  })

  return {
    checkItems,
    dependencyCards,
    endpointCards,
    errorMessage,
    health,
    isHealthy,
    lastCheckedLabel,
    loadHealth,
    loading,
    overview,
    readinessScore,
    signalIcons,
    summaryCards,
  }
}
