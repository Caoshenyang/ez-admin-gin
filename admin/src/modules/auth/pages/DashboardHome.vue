<script setup lang="ts">
import { NAlert, NButton, NCard, NTag } from 'naive-ui'
import EzMetricCard from '@/components/ez/EzMetricCard.vue'
import EmptyState from '@/components/EmptyState.vue'
import OverviewPanel from '../components/OverviewPanel.vue'
import QuickEntry from '../components/QuickEntry.vue'
import { useDashboardHomePage } from '../composables/useDashboardHomePage'

const {
  currentDateLabel,
  currentUserLabel,
  dashboard,
  displayText,
  errorMessage,
  focusFacts,
  formatDashboardDateTime,
  formatDashboardRoutePath,
  formatMetricValue,
  getHealthTagType,
  getLoginStatusLabel,
  getLoginStatusTagType,
  getStatusTagType,
  healthItems,
  healthPath,
  heroIcon,
  heroStatusText,
  isHealthy,
  latestNotices,
  loading,
  loadDashboard,
  metricCards,
  navigateTo,
  quickLinks,
  recentLogins,
  recentOperations,
  userManagePath,
} = useDashboardHomePage()
</script>

<template>
  <main class="admin-page admin-page-scroll dash-gap">
    <NAlert
      v-if="errorMessage"
      type="error"
      title="工作台同步失败"
      class="rounded-[var(--ez-radius-lg)]"
      :bordered="false"
    >
      {{ errorMessage }}
    </NAlert>

    <!-- 欢迎区 -->
    <section class="flex items-center justify-between gap-4">
      <div>
        <h1 class="text-[var(--ez-text-3xl)] font-bold text-[var(--ez-text-main)]">
          {{ currentUserLabel }}，{{ currentDateLabel }}
        </h1>
        <p class="mt-1.5 text-[var(--ez-text-base)] text-[var(--ez-text-sub)]">
          {{ heroStatusText }}
        </p>
      </div>
      <div class="flex items-center gap-2.5">
        <NTag
          :type="isHealthy ? 'success' : 'warning'"
          size="medium"
          round
          :bordered="false"
        >
          {{ isHealthy ? '运行稳定' : '存在待检查项' }}
        </NTag>
        <NButton
          v-if="userManagePath"
          type="primary"
          @click="navigateTo(userManagePath)"
        >
          用户管理
        </NButton>
        <NButton v-else-if="healthPath" secondary @click="navigateTo(healthPath)">
          系统状态
        </NButton>
        <NButton secondary :loading="loading" @click="void loadDashboard()">刷新</NButton>
      </div>
    </section>

    <!-- 统计卡片 -->
    <section class="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
      <EzMetricCard
        v-for="item in metricCards"
        :key="item.label"
        :title="item.label"
        :value="item.value"
        :trend="item.hint"
        trend-type="up"
        :icon="item.icon"
      />
    </section>

    <!-- 系统概览 + 快捷入口 -->
    <section class="grid gap-5 xl:grid-cols-[minmax(0,1fr)_minmax(360px,420px)]">
      <NCard class="ez-card-elevated" :bordered="false" content-class="ez-card-content-reset">
        <OverviewPanel
          :hero-icon="heroIcon"
          :health-items="healthItems"
          :focus-facts="focusFacts"
          :status-text="heroStatusText"
          :get-health-tag-type="getHealthTagType"
        />
      </NCard>

      <NCard class="ez-card-elevated" :bordered="false">
        <QuickEntry :links="quickLinks" @navigate="navigateTo" />
      </NCard>
    </section>

    <section class="grid gap-5 xl:grid-cols-[minmax(0,1fr)_360px]">
      <NCard class="ez-card-elevated">
        <template #header>
          <span class="dash-card-title">登录趋势（近 7 天）</span>
        </template>
        <div class="dash-trend-chart">
          <svg viewBox="0 0 640 220" role="img" aria-label="登录趋势折线图">
            <defs>
              <linearGradient id="loginTrendFill" x1="0" x2="0" y1="0" y2="1">
                <stop offset="0%" stop-color="#2563FF" stop-opacity="0.22" />
                <stop offset="100%" stop-color="#2563FF" stop-opacity="0" />
              </linearGradient>
            </defs>
            <path class="dash-trend-grid" d="M40 36H610M40 84H610M40 132H610M40 180H610" />
            <path class="dash-trend-area" d="M44 166L126 130L208 150L290 84L372 124L454 72L536 94L610 58V190H44Z" />
            <path class="dash-trend-line dash-trend-line--primary" d="M44 166L126 130L208 150L290 84L372 124L454 72L536 94L610 58" />
            <path class="dash-trend-line dash-trend-line--success" d="M44 176L126 152L208 164L290 128L372 148L454 116L536 126L610 102" />
          </svg>
          <div class="dash-trend-legend">
            <span><i class="dash-dot dash-dot--primary" />登录次数</span>
            <span><i class="dash-dot dash-dot--success" />独立用户</span>
          </div>
        </div>
      </NCard>

      <NCard class="ez-card-elevated">
        <template #header>
          <span class="dash-card-title">系统运行状态</span>
        </template>
        <div class="dash-health-grid">
          <div v-for="item in healthItems" :key="item.label" class="dash-health-item">
            <span>{{ item.label }}</span>
            <NTag :type="getHealthTagType(item.status)" size="small" :bordered="false">{{ item.value }}</NTag>
          </div>
        </div>
      </NCard>
    </section>

    <!-- 最近操作 + 登录/公告 -->
    <section class="grid gap-5 xl:grid-cols-[1fr_400px]">
      <NCard class="ez-card-elevated">
        <template #header>
          <span class="dash-card-title">最近操作</span>
        </template>
        <template #header-extra>
          <NTag round :bordered="false" type="info" size="small">
            {{ recentOperations.length }} 条
          </NTag>
        </template>

        <div v-if="recentOperations.length > 0" class="space-y-3">
          <article
            v-for="item in recentOperations"
            :key="item.id"
            class="dash-log-row"
          >
            <div class="flex flex-wrap items-center justify-between gap-2">
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
                <span class="truncate text-[var(--ez-text-sm)] text-[var(--ez-text-sub)]">
                  {{ displayText(item.username, '系统') }} · {{ formatDashboardRoutePath(item.path) }}
                </span>
              </div>
              <span class="text-[var(--ez-text-xs)] text-[var(--ez-text-light)]">
                {{ formatDashboardDateTime(item.created_at) }}
              </span>
            </div>
            <div class="mt-2 flex items-center gap-4 text-[var(--ez-text-xs)] text-[var(--ez-text-light)]">
              <span>状态码 {{ item.status_code }}</span>
              <span>耗时 {{ item.latency_ms }} ms</span>
            </div>
          </article>
        </div>
        <EmptyState v-else class="dash-empty-state" title="还没有操作日志" description="系统记录到操作后会显示在这里。" />
      </NCard>

      <section class="grid gap-5">
        <NCard class="ez-card-elevated">
          <template #header>
            <span class="dash-card-title">最近登录</span>
          </template>
          <template #header-extra>
            <NTag round :bordered="false" type="warning" size="small">
              失败 {{ formatMetricValue(dashboard?.metrics.today_login_failed_total) }}
            </NTag>
          </template>

          <div v-if="recentLogins.length > 0" class="space-y-3">
            <article
              v-for="item in recentLogins"
              :key="item.id"
              class="dash-log-row dash-log-row--flat"
            >
              <div class="flex items-center justify-between gap-2">
                <div class="min-w-0">
                  <div class="flex items-center gap-2">
                    <span class="text-[var(--ez-text-sm)] font-medium text-[var(--ez-text-main)]">
                      {{ displayText(item.username) }}
                    </span>
                    <NTag
                      :type="getLoginStatusTagType(item.status)"
                      size="small"
                      round
                      :bordered="false"
                    >
                      {{ getLoginStatusLabel(item.status) }}
                    </NTag>
                  </div>
                  <p class="mt-1 truncate text-[var(--ez-text-xs)] text-[var(--ez-text-light)]">
                    {{ displayText(item.message, '登录状态已记录') }}
                  </p>
                </div>
                <span class="text-[var(--ez-text-xs)] text-[var(--ez-text-light)]">{{ displayText(item.ip) }}</span>
              </div>
              <p class="mt-2 text-[var(--ez-text-xs)] text-[var(--ez-text-light)]">
                {{ formatDashboardDateTime(item.created_at) }}
              </p>
            </article>
          </div>
          <EmptyState v-else class="dash-empty-state" title="还没有登录记录" description="用户登录后会显示最近状态。" />
        </NCard>

        <NCard class="ez-card-elevated">
          <template #header>
            <span class="dash-card-title">最新公告</span>
          </template>
          <template #header-extra>
            <NTag round :bordered="false" type="success" size="small">
              {{ formatMetricValue(dashboard?.metrics.notice_total) }} 条
            </NTag>
          </template>

          <div v-if="latestNotices.length > 0" class="space-y-3">
            <article
              v-for="item in latestNotices"
              :key="item.id"
              class="dash-log-row"
            >
              <p class="text-[var(--ez-text-sm)] font-medium text-[var(--ez-text-main)]">{{ displayText(item.title) }}</p>
              <p class="mt-1 text-[var(--ez-text-xs)] text-[var(--ez-text-light)]">
                {{ formatDashboardDateTime(item.updated_at) }}
              </p>
            </article>
          </div>
          <EmptyState v-else class="dash-empty-state" title="当前没有启用中的公告" description="发布并启用公告后会显示在这里。" />
        </NCard>
      </section>
    </section>
  </main>
</template>

<style scoped>
.dash-gap {
  gap: 24px;
}

.dash-card-title {
  font-size: var(--ez-text-lg);
  font-weight: 600;
  color: var(--ez-text-main);
}

.dash-trend-chart {
  position: relative;
}

.dash-trend-chart svg {
  width: 100%;
  height: 240px;
}

.dash-trend-grid {
  stroke: var(--ez-border-light);
  stroke-width: 1;
}

.dash-trend-area {
  fill: url(#loginTrendFill);
}

.dash-trend-line {
  fill: none;
  stroke-linecap: round;
  stroke-linejoin: round;
  stroke-width: 4;
}

.dash-trend-line--primary {
  stroke: var(--ez-primary);
}

.dash-trend-line--success {
  stroke: var(--ez-success);
}

.dash-trend-legend {
  display: flex;
  justify-content: center;
  gap: 18px;
  color: var(--ez-text-secondary);
  font-size: 12px;
}

.dash-dot {
  display: inline-block;
  width: 8px;
  height: 8px;
  margin-right: 6px;
  border-radius: 50%;
}

.dash-dot--primary {
  background: var(--ez-primary);
}

.dash-dot--success {
  background: var(--ez-success);
}

.dash-health-grid {
  display: grid;
  gap: 12px;
}

.dash-health-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  border: 1px solid var(--ez-border-light);
  border-radius: 10px;
  background: var(--ez-page-bg);
  padding: 12px 14px;
  color: var(--ez-text-main);
  font-size: 13px;
}

.dash-empty-state {
  --empty-state-min-height: 144px;
  --empty-state-padding: 22px 16px;
}

.dash-log-row {
  border-radius: var(--ez-radius-md);
  border: 1px solid var(--ez-border);
  padding: 12px 14px;
  transition: all 0.2s ease;
}

.dash-log-row:hover {
  border-color: var(--ez-border);
  background: var(--ez-brand-soft);
}

.dash-log-row--flat {
  background: var(--ez-page-bg);
  border-color: transparent;
}

.dash-log-row--flat:hover {
  background: var(--ez-brand-soft);
}
</style>
