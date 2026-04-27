<script setup lang="ts">
import { NAlert, NButton, NCard, NTag } from 'naive-ui'
import { computed, onMounted, ref } from 'vue'

import { getSystemHealth } from '../../api/health'
import type { SystemHealthData } from '../../types/health'

type ServiceKey = 'database' | 'redis'

const loading = ref(false)
const errorMessage = ref('')
const health = ref<SystemHealthData | null>(null)
const lastCheckedAt = ref('')

const dependencyCards = computed(() => {
  return ([
    {
      key: 'database',
      label: '数据库',
      value: health.value?.database,
      description: '验证 PostgreSQL 连接是否可用',
    },
    {
      key: 'redis',
      label: 'Redis',
      value: health.value?.redis,
      description: '验证缓存和会话依赖是否可用',
    },
  ] satisfies Array<{
    key: ServiceKey
    label: string
    value?: string
    description: string
  }>)
})

const endpointCards = [
  {
    title: '后台接口',
    path: '/api/v1/system/health',
    description: '需要登录和权限，适合在管理台里确认认证链路与依赖状态。',
  },
  {
    title: '公开探针',
    path: '/health',
    description: '给 Nginx、容器健康检查和外部监控使用，不依赖登录态。',
  },
]

const envTagType = computed(() => {
  return health.value?.env === 'prod' ? 'success' : 'warning'
})

const lastCheckedLabel = computed(() => {
  return lastCheckedAt.value ? formatTime(lastCheckedAt.value) : '尚未检查'
})

function formatTime(value: string) {
  return value ? new Date(value).toLocaleString() : '-'
}

function formatStatusLabel(value?: string) {
  return value === 'ok' ? '正常' : value || '待检查'
}

function getStatusTagType(value?: string) {
  return value === 'ok' ? 'success' : 'error'
}

function getStatusText(value?: string) {
  if (value === 'ok') {
    return '服务连通性正常'
  }

  if (loading.value) {
    return '正在刷新状态...'
  }

  return '请点击刷新重新检查'
}

function getErrorMessage(error: unknown) {
  if (typeof error === 'object' && error !== null) {
    const response = (error as { response?: { data?: { message?: string } } }).response
    if (typeof response?.data?.message === 'string' && response.data.message) {
      return response.data.message
    }
  }

  return '系统状态获取失败，请稍后重试。'
}

async function loadHealth() {
  loading.value = true
  errorMessage.value = ''

  try {
    health.value = await getSystemHealth()
    lastCheckedAt.value = new Date().toISOString()
  } catch (error) {
    errorMessage.value = getErrorMessage(error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  void loadHealth()
})
</script>

<template>
  <main class="flex h-full flex-col gap-5 overflow-hidden">
    <section class="flex items-center justify-between gap-4">
      <div>
        <h1 class="text-[28px] font-bold text-[#111827]">系统状态</h1>
        <p class="mt-1 text-sm text-[#6B7280]">
          登录后检查后台运行环境、数据库和 Redis 的连通性。
        </p>
      </div>

      <NButton type="primary" color="#2080F0" :loading="loading" @click="void loadHealth()">
        刷新状态
      </NButton>
    </section>

    <NAlert
      v-if="errorMessage"
      type="error"
      title="状态检查失败"
      class="rounded-lg"
      :bordered="false"
    >
      {{ errorMessage }}
    </NAlert>

    <section class="grid gap-4 xl:grid-cols-[minmax(0,1.4fr)_minmax(280px,0.9fr)]">
      <NCard class="rounded-lg" :bordered="false" content-style="padding: 24px;">
        <div class="flex h-full flex-col gap-5">
          <div class="flex items-start justify-between gap-4">
            <div>
              <p class="text-sm font-medium uppercase tracking-[0.24em] text-[#94A3B8]">
                Runtime Snapshot
              </p>
              <h2 class="mt-2 text-2xl font-bold text-[#111827]">
                {{ health ? '核心依赖全部在线' : '等待首次检查结果' }}
              </h2>
              <p class="mt-2 text-sm leading-6 text-[#6B7280]">
                这个页面调用的是受保护的后台接口，适合在登录后确认权限链路和依赖状态都正常。
              </p>
            </div>

            <NTag :type="envTagType" size="large" round :bordered="false">
              {{ health?.env || 'unknown' }}
            </NTag>
          </div>

          <div class="grid gap-4 md:grid-cols-2">
            <article
              v-for="item in dependencyCards"
              :key="item.key"
              class="rounded-2xl border border-[#E5E7EB] bg-[#F8FAFC] px-5 py-4"
            >
              <div class="flex items-center justify-between gap-3">
                <span class="text-sm font-semibold text-[#334155]">{{ item.label }}</span>
                <NTag
                  :type="getStatusTagType(item.value)"
                  size="small"
                  round
                  :bordered="false"
                >
                  {{ formatStatusLabel(item.value) }}
                </NTag>
              </div>
              <p class="mt-3 text-lg font-bold text-[#111827]">{{ item.value || 'pending' }}</p>
              <p class="mt-1 text-sm text-[#64748B]">{{ item.description }}</p>
            </article>
          </div>

          <div class="rounded-2xl bg-[#111827] px-5 py-4 text-white">
            <div class="flex items-center justify-between gap-4">
              <div>
                <p class="text-xs uppercase tracking-[0.2em] text-white/55">Last Check</p>
                <p class="mt-2 text-base font-semibold">{{ lastCheckedLabel }}</p>
              </div>
              <p class="text-sm text-white/72">{{ getStatusText(health?.database) }}</p>
            </div>
          </div>
        </div>
      </NCard>

      <section class="grid gap-4">
        <NCard class="rounded-lg" :bordered="false" content-style="padding: 20px;">
          <div class="flex flex-col gap-3">
            <div>
              <p class="text-sm font-semibold text-[#111827]">运行环境</p>
              <p class="mt-1 text-sm text-[#6B7280]">
                当前后端 `app.env` 返回值，会随部署环境切换为 `dev` 或 `prod`。
              </p>
            </div>

            <div class="rounded-xl bg-[#F8FAFC] px-4 py-3">
              <p class="text-xs uppercase tracking-[0.18em] text-[#94A3B8]">Environment</p>
              <p class="mt-2 text-2xl font-bold text-[#111827]">{{ health?.env || 'unknown' }}</p>
            </div>
          </div>
        </NCard>

        <NCard class="rounded-lg" :bordered="false" content-style="padding: 20px;">
          <div class="flex flex-col gap-3">
            <div>
              <p class="text-sm font-semibold text-[#111827]">接口职责</p>
              <p class="mt-1 text-sm text-[#6B7280]">
                同样是健康检查，公开探针和后台菜单入口分别服务于不同场景。
              </p>
            </div>

            <article
              v-for="endpoint in endpointCards"
              :key="endpoint.path"
              class="rounded-xl border border-[#E5E7EB] px-4 py-3"
            >
              <div class="flex items-center justify-between gap-3">
                <span class="font-semibold text-[#111827]">{{ endpoint.title }}</span>
                <code class="rounded bg-[#F8FAFC] px-2 py-1 text-xs text-[#475569]">
                  {{ endpoint.path }}
                </code>
              </div>
              <p class="mt-2 text-sm leading-6 text-[#64748B]">{{ endpoint.description }}</p>
            </article>
          </div>
        </NCard>
      </section>
    </section>
  </main>
</template>
