<script setup lang="ts">
import {
  LayersOutline,
  PulseOutline,
  ShieldCheckmarkOutline,
  TimeOutline,
} from '@vicons/ionicons5'
import { NAlert, NButton, NCard, NEmpty, NIcon, NTag } from 'naive-ui'
import type { Component } from 'vue'
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'

import { getDashboardSummary } from '../../api/dashboard'
import { authMenus } from '../../router/dynamic-menu'
import { DashboardLoginStatus, type DashboardData } from '../../types/dashboard'
import { MenuType, type AuthMenu } from '../../types/menu'

interface MetricCard {
  label: string
  value: string
  hint: string
  accent: string
  panelClass: string
  icon: Component
}

interface QuickLink {
  title: string
  path: string
  description: string
}

const router = useRouter()
const loading = ref(false)
const errorMessage = ref('')
const dashboard = ref<DashboardData | null>(null)
const refreshedAt = ref('')

const currentDateLabel = computed(() => {
  return new Intl.DateTimeFormat('zh-CN', {
    month: 'long',
    day: 'numeric',
    weekday: 'long',
  }).format(new Date())
})

const currentUserLabel = computed(() => {
  const user = dashboard.value?.current_user
  if (!user) {
    return '管理员'
  }

  return user.nickname || user.username
})

const visiblePageMenus = computed(() => {
  return flattenPageMenus(authMenus.value)
})

const visiblePageTotal = computed(() => {
  return visiblePageMenus.value.length + 1
})

const healthPath = computed(() => findMenuPathByTitle(visiblePageMenus.value, '系统状态'))
const userManagePath = computed(() => findMenuPathByTitle(visiblePageMenus.value, '用户管理'))

const isHealthy = computed(() => {
  const health = dashboard.value?.health
  if (!health) {
    return false
  }

  return health.database === 'ok' && health.redis === 'ok'
})

const heroStatusText = computed(() => {
  const health = dashboard.value?.health
  if (!health) {
    return loading.value ? '正在拉取当前项目的实时概览...' : '等待首次同步项目运行数据'
  }

  if (health.database === 'ok' && health.redis === 'ok') {
    return '环境、数据库和缓存依赖都在线，可以直接作为管理员首页使用。'
  }

  if (health.database === 'ok') {
    return '数据库正常，Redis 有异常信号，建议优先检查缓存和登录态链路。'
  }

  return '核心依赖出现异常，请先处理环境问题，再继续做业务操作。'
})

const metricCards = computed<MetricCard[]>(() => {
  const metrics = dashboard.value?.metrics

  return [
    {
      label: '启用账号',
      value: formatMetricValue(metrics?.enabled_user_total),
      hint: `总账号 ${formatMetricValue(metrics?.user_total)}`,
      accent: '#1677ff',
      panelClass: 'bg-[#eff6ff]',
      icon: ShieldCheckmarkOutline,
    },
    {
      label: '启用角色',
      value: formatMetricValue(metrics?.enabled_role_total),
      hint: `可访问页面 ${formatMetricValue(visiblePageTotal.value)}`,
      accent: '#0f766e',
      panelClass: 'bg-[#ecfeff]',
      icon: LayersOutline,
    },
    {
      label: '今日操作',
      value: formatMetricValue(metrics?.today_operation_total),
      hint: `失败操作 ${formatMetricValue(metrics?.today_risk_operation_total)}`,
      accent: '#b45309',
      panelClass: 'bg-[#fff7ed]',
      icon: TimeOutline,
    },
    {
      label: '文件沉淀',
      value: formatMetricValue(metrics?.file_total),
      hint: `公告 ${formatMetricValue(metrics?.notice_total)} / 配置 ${formatMetricValue(metrics?.config_total)}`,
      accent: '#047857',
      panelClass: 'bg-[#ecfdf5]',
      icon: PulseOutline,
    },
  ]
})

const healthItems = computed(() => {
  const health = dashboard.value?.health

  return [
    {
      label: '运行环境',
      value: health?.env || 'unknown',
      status: health ? 'ok' : 'pending',
    },
    {
      label: '数据库',
      value: health?.database || 'pending',
      status: health?.database || 'pending',
    },
    {
      label: 'Redis',
      value: health?.redis || 'pending',
      status: health?.redis || 'pending',
    },
  ]
})

const quickLinks = computed<QuickLink[]>(() => {
  return visiblePageMenus.value
    .filter((menu) => menu.path && menu.path !== '/dashboard')
    .slice(0, 6)
    .map((menu) => ({
      title: menu.title,
      path: menu.path,
      description: getQuickLinkDescription(menu.title),
    }))
})

const recentOperations = computed(() => {
  return dashboard.value?.recent_operations ?? []
})

const recentLogins = computed(() => {
  return dashboard.value?.recent_logins ?? []
})

const latestNotices = computed(() => {
  return dashboard.value?.latest_notices ?? []
})

const refreshedLabel = computed(() => {
  return refreshedAt.value ? formatDateTime(refreshedAt.value) : '尚未同步'
})

const focusFacts = computed(() => {
  return [
    {
      label: '可访问页面',
      value: formatMetricValue(visiblePageTotal.value),
      hint: '按当前角色实时计算',
    },
    {
      label: '失败登录',
      value: formatMetricValue(dashboard.value?.metrics.today_login_failed_total),
      hint: '今日累计失败次数',
    },
    {
      label: '最近刷新',
      value: refreshedLabel.value,
      hint: '已同步当前项目快照',
    },
  ]
})

function flattenPageMenus(menus: AuthMenu[]) {
  const result: AuthMenu[] = []

  for (const menu of menus) {
    if (menu.type === MenuType.Menu && menu.path) {
      result.push(menu)
    }

    result.push(...flattenPageMenus(menu.children ?? []))
  }

  return result
}

function findMenuPathByTitle(menus: AuthMenu[], title: string) {
  return menus.find((menu) => menu.title === title)?.path || ''
}

function getQuickLinkDescription(title: string) {
  const descriptionMap: Record<string, string> = {
    系统状态: '检查当前环境、数据库和 Redis 状态',
    用户管理: '维护后台账号、状态和角色绑定',
    角色管理: '调整角色权限和菜单分配',
    菜单管理: '维护动态菜单、按钮权限和路由出口',
    系统配置: '查看系统参数和运行期配置项',
    文件管理: '检查上传文件和资源沉淀',
    操作日志: '回看后台操作链路和失败记录',
    登录日志: '检查最近登录结果和来源 IP',
    公告管理: '维护首页公告和系统通知内容',
  }

  return descriptionMap[title] || '进入对应系统页面继续处理业务'
}

function formatMetricValue(value?: number) {
  return typeof value === 'number' ? new Intl.NumberFormat('zh-CN').format(value) : '--'
}

function formatDateTime(value: string) {
  return value ? new Date(value).toLocaleString() : '-'
}

function formatRoutePath(path: string) {
  return path.replace(/^\/api\/v1/, '') || path
}

function getHealthTagType(value: string) {
  return value === 'ok' ? 'success' : value === 'pending' ? 'default' : 'error'
}

function getStatusTagType(success: boolean) {
  return success ? 'success' : 'error'
}

function getLoginStatusTagType(status: number) {
  return status === DashboardLoginStatus.Success ? 'success' : 'error'
}

function getLoginStatusLabel(status: number) {
  return status === DashboardLoginStatus.Success ? '成功' : '失败'
}

function getErrorMessage(error: unknown) {
  if (typeof error === 'object' && error !== null) {
    const response = (error as { response?: { data?: { message?: string } } }).response
    if (typeof response?.data?.message === 'string' && response.data.message) {
      return response.data.message
    }
  }

  return '工作台数据获取失败，请稍后重试。'
}

async function loadDashboard() {
  loading.value = true
  errorMessage.value = ''

  try {
    dashboard.value = await getDashboardSummary()
    refreshedAt.value = new Date().toISOString()
  } catch (error) {
    errorMessage.value = getErrorMessage(error)
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
</script>

<template>
  <main class="flex h-full flex-col gap-4 overflow-auto pr-1">
    <NAlert
      v-if="errorMessage"
      type="error"
      title="工作台同步失败"
      class="rounded-2xl"
      :bordered="false"
    >
      {{ errorMessage }}
    </NAlert>

    <section class="flex flex-wrap items-start justify-between gap-4">
      <div class="max-w-[760px]">
        <p class="text-xs font-semibold uppercase tracking-[0.28em] text-[#94a3b8]">
          Project Workbench
        </p>
        <h1 class="mt-2 text-[30px] font-semibold tracking-[-0.03em] text-[#0f172a]">
          {{ currentUserLabel }}，今天是 {{ currentDateLabel }}
        </h1>
        <p class="mt-3 text-sm leading-7 text-[#475569]">
          先看关键指标，再看异常和待办，这是后台首页最顺手的阅读顺序。
          当前首页直接读取真实项目数据，适合登录后第一眼判断系统是否正常可用。
        </p>
      </div>

      <div class="flex flex-wrap items-center gap-2">
        <NTag
          :type="isHealthy ? 'success' : 'warning'"
          size="large"
          round
          :bordered="false"
        >
          {{ isHealthy ? '运行稳定' : '存在待检查项' }}
        </NTag>
        <NButton
          v-if="userManagePath"
          type="primary"
          color="#1677ff"
          @click="navigateTo(userManagePath)"
        >
          进入用户管理
        </NButton>
        <NButton v-else-if="healthPath" secondary @click="navigateTo(healthPath)">
          查看系统状态
        </NButton>
        <NButton :loading="loading" @click="void loadDashboard()">刷新工作台</NButton>
      </div>
    </section>

    <section class="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
      <NCard
        v-for="item in metricCards"
        :key="item.label"
        class="rounded-[24px]"
        :bordered="false"
        content-style="padding: 0;"
      >
        <div class="metric-card px-5 py-5" :class="item.panelClass">
          <div class="flex items-start justify-between gap-4">
            <div>
              <p class="text-sm font-medium text-[#475569]">{{ item.label }}</p>
              <p class="mt-3 text-[32px] font-semibold tracking-[-0.03em] text-[#0f172a]">
                {{ item.value }}
              </p>
              <p class="mt-2 text-sm" :style="{ color: item.accent }">{{ item.hint }}</p>
            </div>
            <div
              class="flex h-11 w-11 items-center justify-center rounded-2xl bg-white/88 shadow-[0_10px_24px_rgba(15,23,42,0.06)]"
              :style="{ color: item.accent }"
            >
              <NIcon :component="item.icon" :size="20" />
            </div>
          </div>
        </div>
      </NCard>
    </section>

    <section class="grid gap-4 xl:grid-cols-[minmax(0,1.5fr)_360px]">
      <NCard class="rounded-[28px]" :bordered="false" content-style="padding: 0;">
        <div class="overview-panel px-6 py-6">
          <div class="flex flex-wrap items-start justify-between gap-5">
            <div class="max-w-[680px]">
              <div class="flex items-center gap-3">
                <div
                  class="flex h-12 w-12 items-center justify-center rounded-2xl bg-[#0f172a] text-white"
                >
                  <NIcon :component="PulseOutline" :size="22" />
                </div>
                <div>
                  <p class="text-sm font-semibold uppercase tracking-[0.24em] text-[#94a3b8]">
                    Overview
                  </p>
                  <p class="mt-1 text-[24px] font-semibold tracking-[-0.03em] text-[#0f172a]">
                    {{ heroStatusText }}
                  </p>
                </div>
              </div>

              <div class="mt-5 flex flex-wrap gap-2">
                <div
                  v-for="item in healthItems"
                  :key="item.label"
                  class="rounded-2xl border border-white/70 bg-white/82 px-4 py-3"
                >
                  <div class="flex items-center gap-2">
                    <span class="text-xs font-semibold uppercase tracking-[0.18em] text-[#94a3b8]">
                      {{ item.label }}
                    </span>
                    <NTag
                      :type="getHealthTagType(item.status)"
                      size="small"
                      round
                      :bordered="false"
                    >
                      {{ item.value }}
                    </NTag>
                  </div>
                </div>
              </div>

              <div class="mt-5 flex flex-wrap items-center gap-2 text-sm text-[#64748b]">
                <span class="rounded-full bg-white/82 px-3 py-1.5">
                  当前身份 {{ dashboard?.current_user.username || 'waiting' }}
                </span>
                <span class="rounded-full bg-white/82 px-3 py-1.5">
                  健康页 {{ healthPath ? '已接入' : '未开放' }}
                </span>
              </div>
            </div>

            <div class="grid min-w-[220px] gap-3">
              <div
                v-for="item in focusFacts"
                :key="item.label"
                class="rounded-2xl bg-white/82 px-4 py-4 shadow-[0_10px_24px_rgba(15,23,42,0.05)]"
              >
                <p class="text-xs font-semibold uppercase tracking-[0.18em] text-[#94a3b8]">
                  {{ item.label }}
                </p>
                <p class="mt-2 text-lg font-semibold text-[#0f172a]">{{ item.value }}</p>
                <p class="mt-1 text-sm leading-6 text-[#64748b]">{{ item.hint }}</p>
              </div>
            </div>
          </div>
        </div>
      </NCard>

      <section class="grid gap-4">
        <NCard class="rounded-[24px]" :bordered="false" content-style="padding: 22px;">
          <div class="flex items-center justify-between gap-4">
            <div>
              <p class="text-sm font-semibold text-[#111827]">快捷入口</p>
              <p class="mt-1 text-sm text-[#6B7280]">只保留当前角色最常用的几个落点。</p>
            </div>
            <NTag round :bordered="false" type="info">{{ quickLinks.length }} 项</NTag>
          </div>

          <div v-if="quickLinks.length > 0" class="mt-4 grid gap-3">
            <button
              v-for="item in quickLinks"
              :key="item.path"
              type="button"
              class="quick-link-button"
              @click="navigateTo(item.path)"
            >
              <div class="min-w-0">
                <p class="text-sm font-semibold text-[#111827]">{{ item.title }}</p>
                <p class="mt-1 truncate text-sm text-[#64748B]">{{ item.description }}</p>
              </div>
              <span class="shrink-0 text-xs text-[#94A3B8]">进入</span>
            </button>
          </div>

          <NEmpty v-else class="mt-4" description="当前角色没有额外页面可跳转" />
        </NCard>
      </section>
    </section>

    <section class="grid gap-4 xl:grid-cols-[minmax(0,1.3fr)_minmax(320px,0.95fr)]">
      <NCard class="rounded-[24px]" :bordered="false" content-style="padding: 22px;">
        <div class="flex items-center justify-between gap-4">
          <div>
            <p class="text-sm font-semibold text-[#111827]">最近操作</p>
            <p class="mt-1 text-sm text-[#6B7280]">直接读取 `sys_operation_log`，用于判断后台最近发生了什么。</p>
          </div>
          <NTag round :bordered="false" type="info">{{ recentOperations.length }} 条</NTag>
        </div>

        <div v-if="recentOperations.length > 0" class="mt-4 space-y-3">
          <article
            v-for="item in recentOperations"
            :key="item.id"
            class="rounded-2xl border border-[#e5e7eb] px-4 py-3"
          >
            <div class="flex flex-wrap items-center justify-between gap-3">
              <div class="flex min-w-0 items-center gap-2">
                <NTag
                  :type="getStatusTagType(item.success)"
                  size="small"
                  round
                  :bordered="false"
                >
                  {{ item.success ? '成功' : '失败' }}
                </NTag>
                <NTag size="small" round :bordered="false">{{ item.method }}</NTag>
                <span class="truncate text-sm font-medium text-[#111827]">
                  {{ item.username || '系统' }} · {{ formatRoutePath(item.path) }}
                </span>
              </div>
              <span class="text-sm text-[#64748B]">{{ formatDateTime(item.created_at) }}</span>
            </div>

            <div class="mt-3 flex flex-wrap items-center gap-4 text-sm text-[#64748B]">
              <span>状态码 {{ item.status_code }}</span>
              <span>耗时 {{ item.latency_ms }} ms</span>
            </div>
          </article>
        </div>

        <NEmpty v-else class="mt-4" description="还没有操作日志" />
      </NCard>

      <section class="grid gap-4">
        <NCard class="rounded-[24px]" :bordered="false" content-style="padding: 22px;">
          <div class="flex items-center justify-between gap-4">
            <div>
              <p class="text-sm font-semibold text-[#111827]">最近登录</p>
              <p class="mt-1 text-sm text-[#6B7280]">帮助你快速判断是否存在连续失败登录。</p>
            </div>
            <NTag round :bordered="false" type="warning">
              失败 {{ formatMetricValue(dashboard?.metrics.today_login_failed_total) }}
            </NTag>
          </div>

          <div v-if="recentLogins.length > 0" class="mt-4 space-y-3">
            <article
              v-for="item in recentLogins"
              :key="item.id"
              class="rounded-2xl bg-[#f8fafc] px-4 py-3"
            >
              <div class="flex items-center justify-between gap-3">
                <div class="min-w-0">
                  <div class="flex items-center gap-2">
                    <span class="text-sm font-semibold text-[#111827]">{{ item.username }}</span>
                    <NTag
                      :type="getLoginStatusTagType(item.status)"
                      size="small"
                      round
                      :bordered="false"
                    >
                      {{ getLoginStatusLabel(item.status) }}
                    </NTag>
                  </div>
                  <p class="mt-1 truncate text-sm text-[#64748B]">
                    {{ item.message || '登录状态已记录' }}
                  </p>
                </div>
                <span class="text-xs text-[#94A3B8]">{{ item.ip || '-' }}</span>
              </div>
              <p class="mt-3 text-xs text-[#94A3B8]">{{ formatDateTime(item.created_at) }}</p>
            </article>
          </div>

          <NEmpty v-else class="mt-4" description="还没有登录记录" />
        </NCard>

        <NCard class="rounded-[24px]" :bordered="false" content-style="padding: 22px;">
          <div class="flex items-center justify-between gap-4">
            <div>
              <p class="text-sm font-semibold text-[#111827]">最新公告</p>
              <p class="mt-1 text-sm text-[#6B7280]">展示启用中的公告，便于确认系统对外通知是否最新。</p>
            </div>
            <NTag round :bordered="false" type="success">
              {{ formatMetricValue(dashboard?.metrics.notice_total) }} 条
            </NTag>
          </div>

          <div v-if="latestNotices.length > 0" class="mt-4 space-y-3">
            <article
              v-for="item in latestNotices"
              :key="item.id"
              class="rounded-2xl border border-[#e5e7eb] px-4 py-3"
            >
              <p class="text-sm font-semibold text-[#111827]">{{ item.title }}</p>
              <p class="mt-2 text-xs uppercase tracking-[0.18em] text-[#94A3B8]">
                updated {{ formatDateTime(item.updated_at) }}
              </p>
            </article>
          </div>

          <NEmpty v-else class="mt-4" description="当前没有启用中的公告" />
        </NCard>
      </section>
    </section>
  </main>
</template>

<style scoped>
.overview-panel {
  background:
    radial-gradient(circle at top left, rgba(255, 255, 255, 0.82) 0%, rgba(255, 255, 255, 0) 25%),
    radial-gradient(circle at right center, rgba(22, 119, 255, 0.1) 0%, rgba(22, 119, 255, 0) 34%),
    linear-gradient(135deg, #f3f8ff 0%, #f8fbff 52%, #edf7f0 100%);
}

.quick-link-button {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  width: 100%;
  border: 1px solid #e2e8f0;
  border-radius: 18px;
  background: #f8fafc;
  padding: 14px 16px;
  text-align: left;
  transition:
    border-color 0.2s ease,
    background-color 0.2s ease,
    transform 0.2s ease,
    box-shadow 0.2s ease;
}

.quick-link-button:hover {
  border-color: #bfdbfe;
  background: #ffffff;
  box-shadow: 0 14px 28px rgba(15, 23, 42, 0.08);
  transform: translateY(-1px);
}

.metric-card {
  min-height: 152px;
}
</style>
