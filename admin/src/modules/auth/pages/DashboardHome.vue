<script setup lang="ts">
import { NAlert, NButton, NCard, NEmpty, NIcon, NTag } from 'naive-ui'
import OverviewPanel from '../components/OverviewPanel.vue'
import QuickEntry from '../components/QuickEntry.vue'
import StatCard from '../components/StatCard.vue'
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
      class="rounded-[14px]"
      :bordered="false"
    >
      {{ errorMessage }}
    </NAlert>

    <!-- 欢迎区 -->
    <section class="flex items-center justify-between gap-4">
      <div>
        <h1 class="text-[26px] font-semibold text-[#0F172A]">
          {{ currentUserLabel }}，{{ currentDateLabel }}
        </h1>
        <p class="mt-1 text-[14px] text-[#64748B]">
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
    <section class="grid gap-5 md:grid-cols-2 xl:grid-cols-4">
      <StatCard
        v-for="item in metricCards"
        :key="item.label"
        :label="item.label"
        :value="item.value"
        :hint="item.hint"
        :accent="item.accent"
        :icon-bg="item.iconBg"
        :icon="item.icon"
      />
    </section>

    <!-- 系统概览 + 快捷入口 -->
    <section class="grid gap-5 xl:grid-cols-[1fr_340px]">
      <NCard class="dash-card" :bordered="false" content-style="padding: 0;">
        <OverviewPanel
          :hero-icon="heroIcon"
          :health-items="healthItems"
          :focus-facts="focusFacts"
          :status-text="heroStatusText"
          :get-health-tag-type="getHealthTagType"
        />
      </NCard>

      <NCard class="dash-card" :bordered="false">
        <QuickEntry :links="quickLinks" @navigate="navigateTo" />
      </NCard>
    </section>

    <!-- 最近操作 + 登录/公告 -->
    <section class="grid gap-5 xl:grid-cols-[1fr_400px]">
      <NCard class="dash-card">
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
                <span class="truncate text-[13px] text-[#475569]">
                  {{ displayText(item.username, '系统') }} · {{ formatDashboardRoutePath(item.path) }}
                </span>
              </div>
              <span class="text-[12px] text-[#94A3B8]">
                {{ formatDashboardDateTime(item.created_at) }}
              </span>
            </div>
            <div class="mt-2 flex items-center gap-4 text-[12px] text-[#94A3B8]">
              <span>状态码 {{ item.status_code }}</span>
              <span>耗时 {{ item.latency_ms }} ms</span>
            </div>
          </article>
        </div>
        <NEmpty v-else description="还没有操作日志" />
      </NCard>

      <section class="grid gap-5">
        <NCard class="dash-card">
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
                    <span class="text-[13px] font-medium text-[#0F172A]">
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
                  <p class="mt-1 truncate text-[12px] text-[#94A3B8]">
                    {{ displayText(item.message, '登录状态已记录') }}
                  </p>
                </div>
                <span class="text-[12px] text-[#94A3B8]">{{ displayText(item.ip) }}</span>
              </div>
              <p class="mt-2 text-[12px] text-[#94A3B8]">
                {{ formatDashboardDateTime(item.created_at) }}
              </p>
            </article>
          </div>
          <NEmpty v-else description="还没有登录记录" />
        </NCard>

        <NCard class="dash-card">
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
              <p class="text-[13px] font-medium text-[#0F172A]">{{ displayText(item.title) }}</p>
              <p class="mt-1 text-[12px] text-[#94A3B8]">
                {{ formatDashboardDateTime(item.updated_at) }}
              </p>
            </article>
          </div>
          <NEmpty v-else description="当前没有启用中的公告" />
        </NCard>
      </section>
    </section>
  </main>
</template>

<style scoped>
.dash-gap {
  gap: 24px;
}

.dash-card {
  border: 1px solid #E5EAF3 !important;
  box-shadow: 0 1px 2px rgba(16, 24, 40, 0.04), 0 1px 3px rgba(16, 24, 40, 0.06) !important;
}

.dash-card-title {
  font-size: 15px;
  font-weight: 500;
  color: #0F172A;
}

.dash-log-row {
  border-radius: 12px;
  border: 1px solid #E5EAF3;
  padding: 12px 14px;
  transition: border-color 0.14s ease;
}

.dash-log-row:hover {
  border-color: #CBD5E1;
}

.dash-log-row--flat {
  background: #F6F8FB;
  border-color: transparent;
}

.dash-log-row--flat:hover {
  background: #EFF6FF;
}
</style>
