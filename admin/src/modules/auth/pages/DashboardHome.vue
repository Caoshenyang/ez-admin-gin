<script setup lang="ts">
import {
  FolderOpenOutline,
  KeyOutline,
  ListOutline,
  NotificationsOutline,
  ServerOutline,
  SettingsOutline,
  ShieldCheckmarkOutline,
  TrendingUpOutline,
} from '@vicons/ionicons5'
import { NAlert, NButton, NEmpty, NIcon, NTag } from 'naive-ui'
import type { Component } from 'vue'
import { computed } from 'vue'

import { useDashboardHomePage } from '../composables/useDashboardHomePage'

interface MetricTile {
  change: string
  icon: Component
  label: string
  sparkPath: string
  tone: 'blue' | 'purple' | 'cyan' | 'green'
  value: string
}

interface ActivityItem {
  action: string
  avatar: string
  name: string
  time: string
}

interface CapabilityItem {
  description: string
  icon: Component
  label: string
  path: string
  tone: 'blue' | 'purple' | 'cyan' | 'green'
  value: string
}

interface DependencyItem {
  detail: string
  label: string
  statusTone: 'success' | 'warning' | 'error' | 'default'
  value: string
}

interface ResourceBarItem {
  detail: string
  label: string
  percent: number
  tone: 'blue' | 'purple' | 'cyan' | 'green' | 'orange'
  value: string
}

const {
  currentDateLabel,
  currentUserLabel,
  dashboard,
  errorMessage,
  formatDashboardDateTime,
  formatDashboardRoutePath,
  formatMetricValue,
  getLoginStatusLabel,
  getLoginStatusTagType,
  isHealthy,
  healthPath,
  latestNotices,
  loginLogPath,
  navigateTo,
  operationLogPath,
  recentLogins,
  recentOperations,
  userManagePath,
  visiblePageTotal,
} = useDashboardHomePage()

const metrics = computed(() => dashboard.value?.metrics)
const enabledAccountTotal = computed(() => metrics.value?.enabled_user_total ?? 0)
const accountTotal = computed(() => metrics.value?.user_total ?? 0)
const roleTotal = computed(() => metrics.value?.enabled_role_total ?? 0)
const configTotal = computed(() => metrics.value?.config_total ?? 0)
const todayOperationTotal = computed(() => metrics.value?.today_operation_total ?? 0)
const riskTotal = computed(() => metrics.value?.today_risk_operation_total ?? 0)
const failedLoginTotal = computed(() => metrics.value?.today_login_failed_total ?? 0)
const fileTotal = computed(() => metrics.value?.file_total ?? 0)
const noticeTotal = computed(() => metrics.value?.notice_total ?? 0)
const accountEnableRate = computed(() => {
  if (accountTotal.value <= 0) return '--'
  return `${Math.round((enabledAccountTotal.value / accountTotal.value) * 100)}%`
})
const supportAssetTotal = computed(() => configTotal.value + noticeTotal.value + fileTotal.value)
const operationSuccessRate = computed(() => {
  if (todayOperationTotal.value <= 0) return '--'
  const successTotal = Math.max(todayOperationTotal.value - riskTotal.value, 0)
  return `${Math.round((successTotal / todayOperationTotal.value) * 100)}%`
})
const securityScore = computed(() => {
  const dependencyPenalty = isHealthy.value ? 0 : 30
  const operationPenalty = Math.min(riskTotal.value * 8, 40)
  const loginPenalty = Math.min(failedLoginTotal.value * 6, 30)
  return Math.max(100 - dependencyPenalty - operationPenalty - loginPenalty, 0)
})
const securityScoreTone = computed(() => {
  if (securityScore.value >= 86) return 'success'
  if (securityScore.value >= 70) return 'warning'
  return 'error'
})
const greetingText = computed(() => {
  const hour = new Date().getHours()

  if (hour < 6) return '夜深了，注意休息'
  if (hour < 12) return '上午好'
  if (hour < 18) return '下午好'
  return '晚上好'
})

const metricTiles = computed<MetricTile[]>(() => [
  {
    change: `启用率 ${accountEnableRate.value}`,
    icon: ShieldCheckmarkOutline,
    label: '启用账号',
    sparkPath: 'M4 48 C24 44, 30 32, 48 35 S68 44, 78 24 S100 44, 112 12',
    tone: 'blue',
    value: formatMetricValue(enabledAccountTotal.value),
  },
  {
    change: `可访问页面 ${formatMetricValue(visiblePageTotal.value)}`,
    icon: KeyOutline,
    label: '启用角色',
    sparkPath: 'M4 50 C18 46, 18 22, 34 20 S54 58, 68 42 S82 8, 112 22',
    tone: 'purple',
    value: formatMetricValue(roleTotal.value),
  },
  {
    change: `风险 ${formatMetricValue(riskTotal.value)} / 失败 ${formatMetricValue(failedLoginTotal.value)}`,
    icon: TrendingUpOutline,
    label: '今日操作',
    sparkPath: 'M4 49 C20 45, 26 20, 42 24 S56 54, 70 42 S80 6, 112 22',
    tone: 'cyan',
    value: formatMetricValue(todayOperationTotal.value),
  },
  {
    change: `配置 ${formatMetricValue(configTotal.value)} / 公告 ${formatMetricValue(noticeTotal.value)}`,
    icon: FolderOpenOutline,
    label: '文件资产',
    sparkPath: 'M4 49 C20 45, 26 20, 42 24 S56 54, 70 42 S80 6, 112 22',
    tone: 'green',
    value: formatMetricValue(fileTotal.value),
  },
])

const capabilityItems = computed<CapabilityItem[]>(() => [
  {
    description: `启用账号 ${formatMetricValue(enabledAccountTotal.value)} / 总账号 ${formatMetricValue(accountTotal.value)}`,
    icon: ShieldCheckmarkOutline,
    label: '身份与权限',
    path: userManagePath.value || '/iam/users',
    tone: 'blue',
    value: `${formatMetricValue(roleTotal.value)} 个角色`,
  },
  {
    description: `系统参数 ${formatMetricValue(configTotal.value)} 项，支撑运行期开关和字典值。`,
    icon: SettingsOutline,
    label: '配置中心',
    path: '/system/configs',
    tone: 'purple',
    value: `${formatMetricValue(configTotal.value)} 项`,
  },
  {
    description: `今日操作 ${formatMetricValue(todayOperationTotal.value)} 次，成功率 ${operationSuccessRate.value}。`,
    icon: ListOutline,
    label: '审计日志',
    path: operationLogPath.value || '/audit/operation-logs',
    tone: 'cyan',
    value: `${formatMetricValue(riskTotal.value)} 个风险`,
  },
  {
    description: `公告 ${formatMetricValue(noticeTotal.value)} 条，文件 ${formatMetricValue(fileTotal.value)} 个。`,
    icon: NotificationsOutline,
    label: '内容与文件',
    path: '/system/notices',
    tone: 'green',
    value: `${formatMetricValue(supportAssetTotal.value)} 个资产`,
  },
])

const dependencyItems = computed<DependencyItem[]>(() => {
  const health = dashboard.value?.health

  return [
    {
      detail: '当前服务启动环境',
      label: '运行环境',
      statusTone: health ? 'default' : 'warning',
      value: health?.env || '等待同步',
    },
    {
      detail: '用户、权限、配置和日志的主存储',
      label: '数据库',
      statusTone: health?.database === 'ok' ? 'success' : 'error',
      value: health?.database === 'ok' ? '在线' : '异常',
    },
    {
      detail: '登录态、限流和缓存依赖',
      label: 'Redis',
      statusTone: health?.redis === 'ok' ? 'success' : 'error',
      value: health?.redis === 'ok' ? '在线' : '异常',
    },
  ]
})

const securityFactors = computed<ResourceBarItem[]>(() => [
  {
    detail: isHealthy.value ? '数据库和 Redis 均在线' : '依赖存在异常，请优先检查',
    label: '依赖状态',
    percent: isHealthy.value ? 100 : 50,
    tone: isHealthy.value ? 'green' : 'orange',
    value: isHealthy.value ? '正常' : '待检查',
  },
  {
    detail: `启用账号 ${formatMetricValue(enabledAccountTotal.value)} / ${formatMetricValue(accountTotal.value)}`,
    label: '账号启用',
    percent:
      accountTotal.value > 0
        ? Math.round((enabledAccountTotal.value / accountTotal.value) * 100)
        : 0,
    tone: 'blue',
    value: accountEnableRate.value,
  },
  {
    detail: `今日风险操作 ${formatMetricValue(riskTotal.value)} 次`,
    label: '操作成功',
    percent:
      todayOperationTotal.value > 0
        ? Math.round(
            ((todayOperationTotal.value - riskTotal.value) / todayOperationTotal.value) * 100,
          )
        : 100,
    tone: riskTotal.value > 0 ? 'orange' : 'green',
    value: operationSuccessRate.value,
  },
])

const activityItems = computed<ActivityItem[]>(() => {
  const rows = recentOperations.value.slice(0, 4).map((item) => ({
    action: `${item.success ? '访问了' : '触发异常'} ${formatDashboardRoutePath(item.path)}`,
    avatar: item.username.slice(0, 1).toUpperCase() || 'A',
    name: item.username || '系统用户',
    time: formatDashboardDateTime(item.created_at),
  }))

  if (rows.length > 0) {
    return rows
  }

  return [
    { action: '同步了工作台运行快照', avatar: '系', name: '系统服务', time: '刚刚' },
    { action: '完成了菜单权限加载', avatar: '权', name: '权限服务', time: '刚刚' },
    { action: '等待新的操作日志写入', avatar: '审', name: '审计服务', time: '刚刚' },
  ]
})
</script>

<template>
  <main class="admin-page admin-page-scroll dashboard-page">
    <NAlert
      v-if="errorMessage"
      type="error"
      title="工作台同步失败"
      class="dashboard-alert"
      :bordered="false"
    >
      {{ errorMessage }}
    </NAlert>

    <section class="welcome-block">
      <div class="welcome-copy">
        <span class="welcome-kicker">{{ greetingText }} · {{ currentDateLabel }}</span>
        <h1>欢迎回来，{{ currentUserLabel }}</h1>
        <p>权限规模、运行状态、登录与审计动态已汇总到当前工作台。</p>
      </div>
      <div class="welcome-actions">
        <span class="system-status" :class="{ 'system-status--warning': !isHealthy }">
          <i />
          {{ isHealthy ? '运行正常' : '待检查' }}
        </span>
      </div>
    </section>

    <section class="metric-grid">
      <article
        v-for="item in metricTiles"
        :key="item.label"
        class="metric-card"
        :class="`dashboard-tone--${item.tone}`"
      >
        <span class="metric-icon">
          <NIcon :component="item.icon" />
        </span>
        <div class="metric-copy">
          <span>{{ item.label }}</span>
          <strong>{{ item.value }}</strong>
          <small>{{ item.change }}</small>
        </div>
        <svg class="metric-spark" viewBox="0 0 116 62" aria-hidden="true">
          <path :d="item.sparkPath" />
        </svg>
      </article>
    </section>

    <section class="chart-grid">
      <article class="dashboard-card trend-card">
        <header class="panel-head">
          <div>
            <h2>常用系统模块</h2>
            <span>从当前工作台进入核心管理能力</span>
          </div>
          <NButton secondary size="small">当前</NButton>
        </header>
        <div class="capability-grid">
          <button
            v-for="item in capabilityItems"
            :key="item.label"
            type="button"
            class="capability-card"
            :class="`dashboard-tone--${item.tone}`"
            @click="navigateTo(item.path)"
          >
            <span class="capability-icon">
              <NIcon :component="item.icon" />
            </span>
            <span>
              <b>{{ item.label }}</b>
              <em>{{ item.value }}</em>
              <small>{{ item.description }}</small>
            </span>
          </button>
        </div>
      </article>

      <article class="dashboard-card bar-card">
        <header class="panel-head">
          <div>
            <h2>运行依赖</h2>
            <span>后端核心依赖当前状态</span>
          </div>
          <NButton secondary size="small" @click="navigateTo(healthPath)">详情</NButton>
        </header>
        <div class="dependency-list">
          <article v-for="item in dependencyItems" :key="item.label" class="dependency-row">
            <span class="dependency-icon">
              <NIcon :component="ServerOutline" />
            </span>
            <div>
              <strong>{{ item.label }}</strong>
              <small>{{ item.detail }}</small>
            </div>
            <NTag size="small" :type="item.statusTone" :bordered="false">{{ item.value }}</NTag>
          </article>
        </div>
      </article>

      <article class="dashboard-card security-card">
        <header class="panel-head">
          <div>
            <h2>安全概览</h2>
            <span>按依赖、登录和操作风险计算</span>
          </div>
        </header>
        <div class="security-panel">
          <div class="security-score">
            <span>安全评分</span>
            <strong>{{ securityScore }}</strong>
            <NTag size="small" :type="securityScoreTone" :bordered="false">
              {{ securityScore < 86 ? '需关注' : '稳定' }}
            </NTag>
          </div>
          <div class="security-factors">
            <article v-for="item in securityFactors" :key="item.label" class="security-factor">
              <div class="security-factor-head">
                <strong>{{ item.label }}</strong>
                <span>{{ item.value }}</span>
              </div>
              <div class="security-factor-body">
                <div class="resource-track">
                  <i
                    :class="`resource-fill--${item.tone}`"
                    :style="{ width: `${item.percent}%` }"
                  />
                </div>
                <small>{{ item.detail }}</small>
              </div>
            </article>
          </div>
        </div>
      </article>
    </section>

    <section class="list-grid">
      <article class="dashboard-card list-card">
        <header class="panel-head">
          <div>
            <h2>最近登录</h2>
            <span>账号访问与登录结果</span>
          </div>
          <button type="button" @click="navigateTo(loginLogPath || '/audit/login-logs')">
            查看全部
          </button>
        </header>
        <div v-if="recentLogins.length > 0" class="record-list">
          <div v-for="item in recentLogins.slice(0, 5)" :key="item.id" class="record-row">
            <span class="mini-avatar">{{ item.username.slice(0, 1).toUpperCase() || '用' }}</span>
            <div class="record-copy">
              <p>
                <strong>{{ item.username || '未知用户' }}</strong
                >{{ item.message || '登录状态已记录' }}
              </p>
              <span
                >{{ item.ip || '未知 IP' }} · {{ formatDashboardDateTime(item.created_at) }}</span
              >
            </div>
            <NTag size="small" :bordered="false" :type="getLoginStatusTagType(item.status)">
              {{ getLoginStatusLabel(item.status) }}
            </NTag>
          </div>
        </div>
        <NEmpty v-else class="compact-empty" size="small" description="暂无登录记录" />
      </article>

      <article class="dashboard-card list-card">
        <header class="panel-head">
          <div>
            <h2>最新公告</h2>
            <span>当前发布和更新的信息</span>
          </div>
          <button type="button" @click="navigateTo('/system/notices')">公告管理</button>
        </header>
        <div v-if="latestNotices.length > 0" class="record-list">
          <div v-for="item in latestNotices.slice(0, 5)" :key="item.id" class="record-row">
            <span class="notice-icon">
              <NIcon :component="NotificationsOutline" />
            </span>
            <div class="record-copy">
              <p>
                <strong>{{ item.title }}</strong>
              </p>
              <span>更新于 {{ formatDashboardDateTime(item.updated_at) }}</span>
            </div>
            <NTag size="small" :bordered="false" type="success">已发布</NTag>
          </div>
        </div>
        <NEmpty v-else class="compact-empty" size="small" description="暂无已发布公告" />
      </article>

      <article class="dashboard-card list-card">
        <header class="panel-head">
          <div>
            <h2>最近操作</h2>
            <span>系统审计和接口访问动态</span>
          </div>
          <button type="button" @click="navigateTo(operationLogPath)">查看全部</button>
        </header>
        <div class="activity-list">
          <div
            v-for="item in activityItems"
            :key="`${item.name}-${item.time}-${item.action}`"
            class="activity-row"
          >
            <span class="mini-avatar">{{ item.avatar }}</span>
            <div>
              <p>
                <strong>{{ item.name }}</strong
                >{{ item.action }}
              </p>
              <time>{{ item.time }}</time>
            </div>
          </div>
          <button
            class="create-entry"
            type="button"
            @click="navigateTo(operationLogPath || '/audit/operation-logs')"
          >
            <NIcon :component="ListOutline" />
            查看操作日志
          </button>
        </div>
      </article>
    </section>
  </main>
</template>

<style scoped>
.dashboard-page {
  --dashboard-page-bg: #f6f8fc;
  --dashboard-card-bg: #ffffff;
  --dashboard-card-bg-soft: #fbfcff;
  --dashboard-card-bg-hover: #ffffff;
  --dashboard-card-border: #e8edf5;
  --dashboard-card-border-soft: #eef2f7;
  --dashboard-card-shadow: 0 16px 36px rgba(31, 41, 55, 0.05);
  --dashboard-card-shadow-soft: 0 16px 36px rgba(31, 41, 55, 0.04);
  --dashboard-text-heading: #111827;
  --dashboard-text-strong: #273244;
  --dashboard-text-body: #4b5565;
  --dashboard-text-muted: #7b8798;
  --dashboard-action-text: #59667a;
  --dashboard-brand-text: #315dff;
  --dashboard-welcome-border: #dce8f7;
  --dashboard-welcome-bg:
    linear-gradient(120deg, rgba(49, 93, 255, 0.1), transparent 34%),
    linear-gradient(160deg, #ffffff 0%, #f6f9ff 52%, #effcf7 100%);
  --dashboard-welcome-accent: linear-gradient(90deg, #315dff, #18b883, transparent);
  --dashboard-welcome-shadow: 0 18px 42px rgba(31, 41, 55, 0.06);
  --dashboard-kicker-border: #dde8ff;
  --dashboard-kicker-bg: #f3f7ff;
  --dashboard-kicker-text: #315dff;
  --dashboard-status-border: #d7eee4;
  --dashboard-status-bg: #f3faf7;
  --dashboard-status-text: #269467;
  --dashboard-status-dot: #35b77d;
  --dashboard-status-dot-shadow: rgba(53, 183, 125, 0.12);
  --dashboard-warning-border: #f1dfb4;
  --dashboard-warning-bg: #fff9eb;
  --dashboard-warning-text: #b7791f;
  --dashboard-warning-dot: #e3a52f;
  --dashboard-warning-dot-shadow: rgba(227, 165, 47, 0.14);
  --dashboard-icon-soft: #edf4ff;
  --dashboard-icon-text: #2f6bff;
  --dashboard-security-bg: linear-gradient(135deg, #edf4ff, #f4f1ff);
  --dashboard-resource-track: #edf1f6;
  --dashboard-notice-bg: #eef3ff;
  --dashboard-avatar-bg: linear-gradient(135deg, #eaf1ff, #f3eaff);
  --dashboard-avatar-text: #2f55d4;
  --dashboard-create-bg: linear-gradient(135deg, #eef3ff, #f6f2ff);
  --dashboard-create-text: #3f46ff;

  gap: 18px;
  padding: 0;
  background: var(--dashboard-page-bg);
}

html.dark .dashboard-page {
  --dashboard-page-bg: var(--ez-page-bg);
  --dashboard-card-bg: var(--ez-card-bg);
  --dashboard-card-bg-soft: var(--ez-surface-subtle);
  --dashboard-card-bg-hover: var(--ez-surface-hover);
  --dashboard-card-border: var(--ez-component-border);
  --dashboard-card-border-soft: var(--ez-border-light);
  --dashboard-card-shadow: 0 12px 28px rgba(0, 0, 0, 0.2);
  --dashboard-card-shadow-soft: 0 10px 24px rgba(0, 0, 0, 0.18);
  --dashboard-text-heading: var(--ez-text-main);
  --dashboard-text-strong: var(--ez-text-heading);
  --dashboard-text-body: var(--ez-text-body);
  --dashboard-text-muted: var(--ez-text-muted);
  --dashboard-action-text: var(--ez-text-secondary);
  --dashboard-brand-text: #93c5fd;
  --dashboard-welcome-border: rgba(96, 165, 250, 0.24);
  --dashboard-welcome-bg:
    linear-gradient(120deg, rgba(59, 130, 246, 0.2), transparent 34%),
    linear-gradient(160deg, #111827 0%, #0f172a 54%, #082f2d 100%);
  --dashboard-welcome-accent: linear-gradient(90deg, #60a5fa, #34d399, transparent);
  --dashboard-welcome-shadow: 0 16px 36px rgba(0, 0, 0, 0.22);
  --dashboard-kicker-border: rgba(96, 165, 250, 0.26);
  --dashboard-kicker-bg: rgba(59, 130, 246, 0.16);
  --dashboard-kicker-text: #93c5fd;
  --dashboard-status-border: rgba(18, 185, 129, 0.3);
  --dashboard-status-bg: rgba(18, 185, 129, 0.14);
  --dashboard-status-text: #6ee7b7;
  --dashboard-status-dot: #34d399;
  --dashboard-status-dot-shadow: rgba(52, 211, 153, 0.18);
  --dashboard-warning-border: rgba(245, 158, 11, 0.32);
  --dashboard-warning-bg: rgba(245, 158, 11, 0.16);
  --dashboard-warning-text: #fcd34d;
  --dashboard-warning-dot: #f59e0b;
  --dashboard-warning-dot-shadow: rgba(245, 158, 11, 0.2);
  --dashboard-icon-soft: rgba(59, 130, 246, 0.16);
  --dashboard-icon-text: #93c5fd;
  --dashboard-security-bg: linear-gradient(135deg, rgba(59, 130, 246, 0.18), rgba(139, 92, 246, 0.16));
  --dashboard-resource-track: rgba(148, 163, 184, 0.18);
  --dashboard-notice-bg: rgba(59, 130, 246, 0.14);
  --dashboard-avatar-bg: linear-gradient(135deg, rgba(59, 130, 246, 0.22), rgba(139, 92, 246, 0.18));
  --dashboard-avatar-text: #bfdbfe;
  --dashboard-create-bg: linear-gradient(135deg, rgba(59, 130, 246, 0.18), rgba(139, 92, 246, 0.16));
  --dashboard-create-text: #93c5fd;
}

.dashboard-alert,
.dashboard-card,
.metric-card {
  border-radius: 8px;
}

.welcome-block,
.metric-card,
.panel-head,
.activity-row,
.create-entry {
  display: flex;
  align-items: center;
}

.panel-head button {
  border: 0;
  background: transparent;
  color: var(--dashboard-action-text);
}

.welcome-block {
  position: relative;
  overflow: hidden;
  justify-content: space-between;
  gap: 18px;
  border: 1px solid var(--dashboard-welcome-border);
  border-radius: 8px;
  background: var(--dashboard-welcome-bg);
  padding: 18px 20px;
  box-shadow: var(--dashboard-welcome-shadow);
}

.welcome-block::before {
  position: absolute;
  inset: 0 0 auto;
  height: 3px;
  background: var(--dashboard-welcome-accent);
  content: '';
}

.welcome-copy,
.welcome-actions {
  position: relative;
  z-index: 1;
}

.welcome-copy {
  min-width: 0;
}

.welcome-kicker {
  display: inline-flex;
  align-items: center;
  border: 1px solid var(--dashboard-kicker-border);
  border-radius: 999px;
  background: var(--dashboard-kicker-bg);
  padding: 4px 10px;
  color: var(--dashboard-kicker-text);
  font-size: 12px;
  font-weight: 800;
  line-height: 1.2;
  white-space: nowrap;
}

.welcome-block h1 {
  margin: 10px 0 0;
  color: var(--dashboard-text-heading);
  font-size: 24px;
  font-weight: 900;
  letter-spacing: 0;
  line-height: 1.3;
}

.welcome-block p,
.panel-head span,
.metric-copy span,
.metric-copy small,
.record-copy span,
.activity-row time {
  color: var(--dashboard-text-muted);
  font-size: 13px;
}

.welcome-block p {
  margin: 8px 0 0;
}

.welcome-actions {
  display: flex;
  flex-shrink: 0;
  flex-wrap: wrap;
  align-items: center;
  justify-content: flex-end;
  gap: 10px;
}

.system-status {
  display: inline-flex;
  height: 46px;
  align-items: center;
  gap: 7px;
  border: 1px solid var(--dashboard-status-border);
  border-radius: 8px;
  background: var(--dashboard-status-bg);
  padding: 0 14px;
  color: var(--dashboard-status-text);
  font-size: 13px;
  font-weight: 700;
  white-space: nowrap;
}

.system-status i {
  width: 7px;
  height: 7px;
  border-radius: 999px;
  background: var(--dashboard-status-dot);
  box-shadow: 0 0 0 3px var(--dashboard-status-dot-shadow);
}

.system-status--warning {
  border-color: var(--dashboard-warning-border);
  background: var(--dashboard-warning-bg);
  color: var(--dashboard-warning-text);
}

.system-status--warning i {
  background: var(--dashboard-warning-dot);
  box-shadow: 0 0 0 3px var(--dashboard-warning-dot-shadow);
}

.metric-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 18px;
}

.metric-card {
  min-width: 0;
  border: 1px solid var(--dashboard-card-border);
  background: var(--dashboard-card-bg);
  padding: 22px 20px;
  box-shadow: var(--dashboard-card-shadow);
}

.metric-icon {
  display: grid;
  width: 52px;
  height: 52px;
  flex-shrink: 0;
  place-items: center;
  border-radius: 14px;
  background: var(--tone-gradient);
  color: var(--ez-on-brand);
  font-size: 27px;
  box-shadow: var(--tone-shadow);
}

.metric-copy {
  min-width: 0;
  margin-left: 16px;
}

.metric-copy span,
.metric-copy strong,
.metric-copy small {
  display: block;
}

.metric-copy strong {
  margin-top: 6px;
  color: var(--dashboard-text-heading);
  font-size: 25px;
  font-weight: 900;
  line-height: 1.1;
}

.metric-copy small {
  margin-top: 8px;
  color: var(--dashboard-status-text);
  font-weight: 700;
  word-break: keep-all;
}

.metric-spark {
  width: 112px;
  height: 62px;
  margin-left: auto;
}

.metric-spark path {
  fill: none;
  stroke: var(--tone-line);
  stroke-linecap: round;
  stroke-linejoin: round;
  stroke-width: 3;
}

.dashboard-tone--blue {
  --tone-gradient: linear-gradient(135deg, #315dff, #7367ff);
  --tone-line: #315dff;
  --tone-shadow: 0 12px 24px rgba(49, 93, 255, 0.28);
}

.dashboard-tone--purple {
  --tone-gradient: linear-gradient(135deg, #6145ff, #8b35eb);
  --tone-line: #7a35e8;
  --tone-shadow: 0 12px 24px rgba(108, 65, 255, 0.28);
}

.dashboard-tone--cyan {
  --tone-gradient: linear-gradient(135deg, #1478ff, #26a3ff);
  --tone-line: #1478ff;
  --tone-shadow: 0 12px 24px rgba(20, 120, 255, 0.24);
}

.dashboard-tone--green {
  --tone-gradient: linear-gradient(135deg, #15be8b, #18d4b2);
  --tone-line: #18b883;
  --tone-shadow: 0 12px 24px rgba(21, 190, 139, 0.24);
}

.chart-grid {
  display: grid;
  grid-template-columns: minmax(0, 1.28fr) minmax(0, 1fr) minmax(320px, 0.9fr);
  gap: 18px;
}

.dashboard-card {
  min-width: 0;
  border: 1px solid var(--dashboard-card-border);
  background: var(--dashboard-card-bg);
  padding: 20px;
  box-shadow: var(--dashboard-card-shadow-soft);
}

.panel-head {
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 18px;
}

.panel-head h2 {
  margin: 0;
  color: var(--dashboard-text-heading);
  font-size: 17px;
  font-weight: 900;
  letter-spacing: 0;
}

.panel-head span {
  display: block;
  margin-top: 5px;
}

.panel-head button {
  flex-shrink: 0;
  color: var(--dashboard-create-text);
  font-size: 13px;
  font-weight: 800;
}

.capability-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 14px;
}

.capability-card {
  display: grid;
  min-height: 112px;
  grid-template-columns: 46px minmax(0, 1fr);
  gap: 12px;
  border: 1px solid var(--dashboard-card-border);
  border-radius: 8px;
  background: var(--dashboard-card-bg-soft);
  padding: 14px;
  color: inherit;
  text-align: left;
}

.capability-card:hover {
  border-color: color-mix(in srgb, var(--tone-line) 42%, var(--dashboard-card-border));
  background: var(--dashboard-card-bg-hover);
  box-shadow: var(--dashboard-card-shadow);
}

.capability-icon,
.dependency-icon {
  display: grid;
  flex-shrink: 0;
  place-items: center;
  border-radius: 12px;
}

.capability-icon {
  width: 46px;
  height: 46px;
  background: var(--tone-gradient);
  color: var(--ez-on-brand);
  font-size: 24px;
  box-shadow: var(--tone-shadow);
}

.capability-card b,
.capability-card em,
.capability-card small {
  display: block;
}

.capability-card b {
  color: var(--dashboard-text-strong);
  font-size: 14px;
  font-weight: 900;
}

.capability-card em {
  margin-top: 6px;
  color: var(--tone-line);
  font-style: normal;
  font-weight: 900;
}

.capability-card small {
  margin-top: 6px;
  color: var(--dashboard-text-muted);
  font-size: 12px;
  line-height: 1.45;
}

.dependency-list {
  display: grid;
  gap: 14px;
}

.dependency-row {
  display: grid;
  grid-template-columns: 42px minmax(0, 1fr) auto;
  align-items: center;
  gap: 12px;
  border: 1px solid var(--dashboard-card-border-soft);
  border-radius: 8px;
  background: var(--dashboard-card-bg-soft);
  padding: 13px;
}

.dependency-icon {
  width: 42px;
  height: 42px;
  background: var(--dashboard-icon-soft);
  color: var(--dashboard-icon-text);
  font-size: 22px;
}

.dependency-row strong,
.dependency-row small {
  display: block;
}

.dependency-row strong {
  color: var(--dashboard-text-strong);
  font-size: 14px;
  font-weight: 900;
}

.dependency-row small {
  margin-top: 4px;
  color: var(--dashboard-text-muted);
  font-size: 12px;
}

.security-panel {
  display: grid;
  gap: 18px;
}

.security-score {
  display: grid;
  min-height: 104px;
  place-items: center;
  border-radius: 8px;
  background: var(--dashboard-security-bg);
  text-align: center;
}

.security-score span {
  color: var(--dashboard-text-muted);
  font-size: 12px;
}

.security-score strong {
  color: var(--dashboard-text-heading);
  font-size: 36px;
  font-weight: 900;
  line-height: 1.1;
}

.security-factors {
  display: grid;
  gap: 15px;
}

.security-factor {
  display: grid;
  gap: 7px;
}

.security-factor-head,
.security-factor-body {
  display: flex;
  align-items: center;
}

.security-factor-head {
  justify-content: space-between;
  gap: 12px;
}

.security-factor-head strong {
  color: var(--dashboard-text-strong);
  font-size: 13px;
  font-weight: 900;
  white-space: nowrap;
}

.security-factor-head span {
  color: var(--dashboard-text-heading);
  font-size: 14px;
  font-weight: 900;
  white-space: nowrap;
}

.security-factor-body {
  gap: 12px;
}

.resource-track {
  overflow: hidden;
  height: 9px;
  flex: 1;
  border-radius: 999px;
  background: var(--dashboard-resource-track);
}

.resource-track i {
  display: block;
  height: 100%;
  border-radius: inherit;
}

.security-factor-body small {
  overflow: hidden;
  max-width: 150px;
  flex-shrink: 0;
  color: var(--dashboard-text-muted);
  font-size: 12px;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.resource-fill--blue {
  background: #2f6bff;
}

.resource-fill--purple {
  background: #7a35e8;
}

.resource-fill--cyan {
  background: #1478ff;
}

.resource-fill--green {
  background: #18b883;
}

.resource-fill--orange {
  background: var(--ez-warning);
}

.list-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 18px;
}

.list-card {
  padding-bottom: 16px;
}

.record-list,
.activity-list {
  display: grid;
  gap: 14px;
}

.record-row {
  display: flex;
  align-items: center;
  gap: 12px;
  min-width: 0;
}

.record-copy {
  min-width: 0;
  overflow: hidden;
  flex: 1;
}

.record-copy p {
  overflow: hidden;
  margin: 0;
  color: var(--dashboard-text-body);
  font-size: 13px;
  line-height: 1.5;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.record-copy strong {
  margin-right: 8px;
  color: var(--dashboard-text-strong);
  font-weight: 900;
}

.record-copy span {
  display: block;
  overflow: hidden;
  margin-top: 2px;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.record-row :deep(.n-tag) {
  flex-shrink: 0;
}

.notice-icon {
  display: grid;
  width: 28px;
  height: 28px;
  flex-shrink: 0;
  place-items: center;
  border-radius: 8px;
  background: var(--dashboard-notice-bg);
  color: var(--dashboard-brand-text);
  font-size: 15px;
}

.mini-avatar {
  display: grid;
  width: 28px;
  height: 28px;
  flex-shrink: 0;
  place-items: center;
  border-radius: 999px;
  background: var(--dashboard-avatar-bg);
  color: var(--dashboard-avatar-text);
  font-size: 12px;
  font-weight: 900;
}

.compact-empty {
  --n-icon-size: 32px;

  display: grid;
  min-height: 122px;
  place-content: center;
}

.activity-row {
  align-items: flex-start;
  gap: 12px;
}

.activity-row div {
  display: grid;
  min-width: 0;
  flex: 1;
  grid-template-columns: minmax(0, 1fr) auto;
  gap: 12px;
}

.activity-row p {
  overflow: hidden;
  margin: 0;
  color: var(--dashboard-text-body);
  font-size: 13px;
  line-height: 1.5;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.activity-row strong {
  margin-right: 8px;
  color: var(--dashboard-text-strong);
  font-weight: 900;
}

.activity-row time {
  white-space: nowrap;
}

.create-entry {
  justify-content: center;
  gap: 8px;
  margin-top: 2px;
  border: 0;
  border-radius: 8px;
  background: var(--dashboard-create-bg);
  padding: 11px 14px;
  color: var(--dashboard-create-text);
  font-weight: 900;
}

@media (max-width: 1320px) {
  .metric-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .chart-grid,
  .list-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 900px) {
  .welcome-block {
    align-items: flex-start;
    flex-direction: column;
  }

  .welcome-actions {
    justify-content: flex-start;
  }
}

@media (max-width: 720px) {
  .dashboard-page {
    gap: 14px;
  }

  .metric-grid {
    grid-template-columns: 1fr;
  }

  .welcome-block {
    padding: 16px;
  }

  .welcome-kicker {
    white-space: normal;
  }

  .welcome-actions,
  .system-status {
    width: 100%;
  }

  .metric-card {
    align-items: flex-start;
    flex-direction: column;
  }

  .metric-copy {
    margin: 12px 0 0;
  }

  .metric-spark {
    width: 100%;
    margin-left: 0;
  }

  .capability-grid {
    grid-template-columns: 1fr;
  }

  .activity-row div {
    grid-template-columns: 1fr;
    gap: 2px;
  }
}
</style>
