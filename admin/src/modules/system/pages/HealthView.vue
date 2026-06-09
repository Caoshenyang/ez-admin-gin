<script setup lang="ts">
import { NAlert, NIcon, NTag } from 'naive-ui'

import EzActionButton from '@/components/ez/EzActionButton.vue'
import PageHeader from '@/components/PageHeader.vue'
import { useHealthPage } from '../composables/useHealthPage'

const {
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
} = useHealthPage()
</script>

<template>
  <main class="admin-page admin-page-scroll health-page">
    <section class="admin-page-section health-page__section">
      <PageHeader title="系统状态" description="查看运行环境、核心依赖和健康探针。">
        <template #actions>
          <EzActionButton
            kind="refresh"
            label="刷新状态"
            type="primary"
            :loading="loading"
            @click="void loadHealth()"
          />
        </template>
      </PageHeader>

      <NAlert
        v-if="errorMessage"
        type="error"
        title="状态检查失败"
        class="health-alert"
        :bordered="false"
      >
        {{ errorMessage }}
      </NAlert>

      <section class="health-overview" :class="overview.toneClass">
        <div class="health-overview__body">
          <div class="health-overview__main">
            <span class="health-overview__icon" :class="{ 'is-loading': loading }">
              <NIcon :component="signalIcons.icon" />
            </span>

            <div class="health-overview__copy">
              <div class="health-overview__kicker">
                <span class="health-overview__dot" />
                <span>运行快照</span>
                <NTag size="small" :type="overview.tagType" :bordered="false">
                  {{ overview.statusLabel }}
                </NTag>
              </div>
              <h2>{{ overview.title }}</h2>
              <p>{{ overview.description }}</p>

              <div class="health-overview__meta">
                <span>
                  <NIcon :component="signalIcons.softIcon" />
                  {{ health?.env || 'unknown' }}
                </span>
                <span>最近刷新 {{ lastCheckedLabel }}</span>
                <span>{{ isHealthy ? '依赖全部在线' : '依赖待确认' }}</span>
              </div>
            </div>
          </div>

          <aside class="health-score-panel">
            <div class="health-score" :style="{ '--score': `${readinessScore}%` }">
              <div class="health-score__inner">
                <strong>{{ readinessScore }}</strong>
                <span>score</span>
              </div>
            </div>
            <div class="health-score-panel__copy">
              <p>就绪评分</p>
              <small>{{ isHealthy ? 'Ready' : loading ? 'Syncing' : 'Attention' }}</small>
            </div>
          </aside>
        </div>
      </section>

      <section class="health-summary-grid">
        <article
          v-for="item in summaryCards"
          :key="item.key"
          class="health-summary-card"
          :class="item.toneClass"
        >
          <span class="health-summary-card__icon">
            <NIcon :component="item.icon" />
          </span>
          <div class="health-summary-card__copy">
            <span>{{ item.label }}</span>
            <strong>{{ item.value }}</strong>
            <small>{{ item.detail }}</small>
          </div>
          <NTag v-if="item.tagLabel" size="small" :type="item.tagType" :bordered="false">
            {{ item.tagLabel }}
          </NTag>
        </article>
      </section>

      <section class="health-content-grid">
        <article class="health-panel health-panel--wide">
          <header class="health-panel__head">
            <div>
              <h3>核心依赖</h3>
              <p>数据库与缓存的实时连通性</p>
            </div>
            <NTag size="small" :type="isHealthy ? 'success' : 'warning'" :bordered="false">
              {{ isHealthy ? '全部在线' : '待处理' }}
            </NTag>
          </header>

          <div class="health-service-grid">
            <article
              v-for="item in dependencyCards"
              :key="item.key"
              class="health-service-card"
              :class="item.toneClass"
            >
              <header>
                <span class="health-service-card__icon">
                  <NIcon :component="item.icon" />
                </span>
                <div>
                  <span>{{ item.label }}</span>
                  <strong>{{ item.progressLabel }}</strong>
                </div>
                <NTag size="small" :type="item.tagType" :bordered="false">
                  {{ item.statusLabel }}
                </NTag>
              </header>

              <p>{{ item.description }}</p>

              <div class="health-progress">
                <i :style="{ width: `${item.progress}%` }" />
              </div>

              <footer>
                <span>{{ item.detail }}</span>
                <code>{{ item.value || 'pending' }}</code>
              </footer>
            </article>
          </div>
        </article>

        <article class="health-panel">
          <header class="health-panel__head">
            <div>
              <h3>健康探针</h3>
              <p>管理台与外部监控入口</p>
            </div>
          </header>

          <div class="health-endpoint-list">
            <article
              v-for="endpoint in endpointCards"
              :key="endpoint.path"
              class="health-endpoint-card"
              :class="endpoint.toneClass"
            >
              <span class="health-endpoint-card__icon">
                <NIcon :component="endpoint.icon" />
              </span>
              <div>
                <div class="health-endpoint-card__title">
                  <strong>{{ endpoint.title }}</strong>
                  <NTag size="tiny" :bordered="false">{{ endpoint.badge }}</NTag>
                </div>
                <code>{{ endpoint.method }} {{ endpoint.path }}</code>
                <p>{{ endpoint.description }}</p>
              </div>
            </article>
          </div>
        </article>
      </section>

      <article class="health-panel">
        <header class="health-panel__head">
          <div>
            <h3>检查清单</h3>
            <p>本次快照覆盖的关键项</p>
          </div>
        </header>

        <div class="health-check-list">
          <article
            v-for="item in checkItems"
            :key="item.key"
            class="health-check-row"
            :class="item.toneClass"
          >
            <span class="health-check-row__marker" />
            <div>
              <strong>{{ item.label }}</strong>
              <small>{{ item.description }}</small>
            </div>
            <NTag size="small" :type="item.tagType" :bordered="false">
              {{ item.statusLabel }}
            </NTag>
          </article>
        </div>
      </article>
    </section>
  </main>
</template>

<style scoped>
.health-page {
  gap: 16px;
  padding: 0;
  background: var(--ez-page-bg);
}

.health-page__section {
  gap: 16px;
}

.health-alert,
.health-overview,
.health-summary-card,
.health-panel,
.health-service-card,
.health-endpoint-card,
.health-check-row {
  border-radius: 8px;
}

.health-overview,
.health-summary-card,
.health-panel {
  border: 1px solid var(--ez-component-border);
  background: var(--ez-card-bg);
  box-shadow: var(--ez-component-shadow);
}

.health-overview {
  --overview-color: var(--ez-brand);
  --overview-soft: var(--ez-brand-soft);
  --overview-border: var(--ez-brand-border);
  overflow: hidden;
  padding: 20px;
}

.health-overview--success {
  --overview-color: var(--ez-success);
  --overview-soft: rgba(18, 185, 129, 0.1);
  --overview-border: rgba(18, 185, 129, 0.22);
}

.health-overview--warning {
  --overview-color: var(--ez-warning);
  --overview-soft: var(--ez-warning-bg);
  --overview-border: rgba(245, 158, 11, 0.24);
}

.health-overview--danger {
  --overview-color: var(--ez-danger);
  --overview-soft: var(--ez-danger-bg);
  --overview-border: rgba(239, 68, 68, 0.22);
}

.health-overview--neutral {
  --overview-color: var(--ez-text-secondary);
  --overview-soft: var(--ez-surface-subtle);
  --overview-border: var(--ez-component-border);
}

.health-overview__body,
.health-overview__main,
.health-overview__kicker,
.health-overview__meta,
.health-score-panel,
.health-summary-card,
.health-service-card header,
.health-service-card footer,
.health-endpoint-card,
.health-endpoint-card__title,
.health-check-row {
  display: flex;
  align-items: center;
}

.health-overview__body {
  max-width: 980px;
  align-items: stretch;
  gap: 28px;
}

.health-overview__main {
  align-items: flex-start;
  min-width: 0;
  flex: 1;
  gap: 16px;
}

.health-overview__icon {
  display: inline-flex;
  width: 54px;
  height: 54px;
  flex: 0 0 auto;
  align-items: center;
  justify-content: center;
  border: 1px solid var(--overview-border);
  border-radius: 8px;
  background: var(--overview-soft);
  color: var(--overview-color);
  font-size: 30px;
}

.health-overview__icon.is-loading {
  animation: health-pulse 1.4s ease-in-out infinite;
}

.health-overview__copy {
  min-width: 0;
}

.health-overview__kicker {
  flex-wrap: wrap;
  gap: 8px;
  color: var(--ez-text-secondary);
  font-size: 12px;
  font-weight: 700;
}

.health-overview__dot,
.health-check-row__marker {
  display: inline-block;
  width: 8px;
  height: 8px;
  flex: 0 0 auto;
  border-radius: 999px;
  background: var(--overview-color);
  box-shadow: 0 0 0 4px var(--overview-soft);
}

.health-overview h2 {
  margin: 9px 0 0;
  color: var(--ez-text-main);
  font-size: 24px;
  font-weight: 800;
  line-height: 1.2;
}

.health-overview p {
  margin: 8px 0 0;
  max-width: 620px;
  color: var(--ez-text-secondary);
  font-size: 13px;
  line-height: 1.7;
}

.health-overview__meta {
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 16px;
}

.health-overview__meta span {
  display: inline-flex;
  min-height: 30px;
  align-items: center;
  gap: 6px;
  border: 1px solid var(--ez-border-light);
  border-radius: 8px;
  background: var(--ez-surface-subtle);
  padding: 5px 10px;
  color: var(--ez-text-secondary);
  font-size: 12px;
  font-weight: 700;
}

.health-score-panel {
  width: 220px;
  flex: 0 0 auto;
  align-self: center;
  justify-content: flex-start;
  gap: 12px;
  border-left: 1px solid var(--ez-border-light);
  padding-left: 24px;
}

.health-score {
  position: relative;
  width: 88px;
  height: 88px;
  flex: 0 0 auto;
  border-radius: 50%;
  background: conic-gradient(var(--overview-color) var(--score, 0%), var(--ez-border-light) 0);
}

.health-score::before {
  position: absolute;
  inset: 8px;
  border-radius: inherit;
  background: var(--ez-card-bg);
  content: '';
}

.health-score__inner {
  position: absolute;
  inset: 8px;
  z-index: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
}

.health-score__inner strong,
.health-score__inner span {
  position: relative;
  display: block;
  text-align: center;
}

.health-score__inner strong {
  color: var(--ez-text-main);
  font-size: 23px;
  font-weight: 800;
  line-height: 0.95;
}

.health-score__inner span {
  margin-top: 6px;
  color: var(--ez-text-placeholder);
  font-size: 10px;
  font-weight: 700;
  line-height: 1;
}

.health-score-panel__copy {
  min-width: 0;
}

.health-score-panel p {
  margin: 0;
  color: var(--ez-text-main);
  font-size: 13px;
  font-weight: 800;
}

.health-score-panel small {
  display: block;
  margin-top: 4px;
  color: var(--ez-text-secondary);
  font-size: 12px;
  font-weight: 700;
}

.health-summary-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 12px;
}

.health-summary-card {
  --health-tone-color: var(--ez-brand);
  --health-tone-bg: var(--ez-brand-soft);
  --health-tone-border: var(--ez-brand-border);
  display: grid;
  grid-template-columns: 38px minmax(0, 1fr) auto;
  min-width: 0;
  align-items: center;
  gap: 12px;
  padding: 14px 12px;
}

.health-summary-card__icon,
.health-service-card__icon,
.health-endpoint-card__icon {
  display: inline-flex;
  width: 38px;
  height: 38px;
  flex: 0 0 auto;
  align-items: center;
  justify-content: center;
  border: 1px solid var(--health-tone-border);
  border-radius: 8px;
  background: var(--health-tone-bg);
  color: var(--health-tone-color);
  font-size: 22px;
}

.health-summary-card__copy {
  min-width: 0;
  flex: 1;
}

.health-summary-card__copy span,
.health-panel__head p,
.health-service-card p,
.health-service-card footer span,
.health-endpoint-card p,
.health-check-row small {
  color: var(--ez-text-secondary);
  font-size: 12px;
  line-height: 1.6;
}

.health-summary-card__copy strong {
  display: block;
  overflow: hidden;
  margin-top: 3px;
  color: var(--ez-text-main);
  font-size: 19px;
  font-weight: 800;
  line-height: 1.2;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.health-summary-card__copy small {
  display: block;
  overflow: hidden;
  margin-top: 3px;
  color: var(--ez-text-placeholder);
  font-size: 11px;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.health-content-grid {
  display: grid;
  grid-template-columns: minmax(0, 1.45fr) minmax(320px, 0.85fr);
  gap: 12px;
}

.health-panel {
  min-width: 0;
  padding: 16px;
}

.health-panel__head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 14px;
}

.health-panel__head h3 {
  margin: 0;
  color: var(--ez-text-main);
  font-size: 15px;
  font-weight: 800;
  line-height: 1.35;
}

.health-panel__head p {
  margin: 3px 0 0;
}

.health-service-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
}

.health-service-card,
.health-endpoint-card,
.health-check-row {
  --health-tone-color: var(--ez-brand);
  --health-tone-bg: var(--ez-brand-soft);
  --health-tone-border: var(--ez-brand-border);
}

.health-service-card {
  border: 1px solid var(--health-tone-border);
  background: var(--ez-surface-subtle);
  padding: 14px;
}

.health-service-card header {
  gap: 12px;
}

.health-service-card header > div {
  min-width: 0;
  flex: 1;
}

.health-service-card header span {
  color: var(--ez-text-secondary);
  font-size: 12px;
  font-weight: 700;
}

.health-service-card header strong {
  display: block;
  overflow: hidden;
  margin-top: 2px;
  color: var(--ez-text-main);
  font-size: 20px;
  font-weight: 800;
  line-height: 1.2;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.health-service-card p {
  margin: 14px 0 0;
}

.health-progress {
  overflow: hidden;
  height: 8px;
  margin-top: 14px;
  border-radius: 999px;
  background: var(--ez-border-light);
}

.health-progress i {
  display: block;
  height: 100%;
  border-radius: inherit;
  background: var(--health-tone-color);
  transition: width 0.2s ease;
}

.health-service-card footer {
  justify-content: space-between;
  gap: 12px;
  margin-top: 12px;
}

.health-service-card code,
.health-endpoint-card code {
  display: inline-flex;
  max-width: 100%;
  overflow: hidden;
  border-radius: 6px;
  background: var(--ez-card-bg);
  padding: 4px 7px;
  color: var(--ez-text-secondary);
  font-size: 11px;
  line-height: 1.4;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.health-endpoint-list,
.health-check-list {
  display: grid;
  gap: 10px;
}

.health-endpoint-card {
  align-items: flex-start;
  gap: 12px;
  border: 1px solid var(--ez-border-light);
  background: var(--ez-surface-subtle);
  padding: 12px;
}

.health-endpoint-card > div {
  min-width: 0;
}

.health-endpoint-card__title {
  justify-content: space-between;
  gap: 8px;
  margin-bottom: 7px;
}

.health-endpoint-card__title strong {
  color: var(--ez-text-main);
  font-size: 13px;
  font-weight: 800;
}

.health-endpoint-card p {
  margin: 8px 0 0;
}

.health-check-row {
  min-height: 58px;
  justify-content: space-between;
  gap: 12px;
  border: 1px solid var(--ez-border-light);
  background: var(--ez-surface-subtle);
  padding: 10px 12px;
}

.health-check-row__marker {
  background: var(--health-tone-color);
  box-shadow: 0 0 0 4px var(--health-tone-bg);
}

.health-check-row > div {
  min-width: 0;
  flex: 1;
}

.health-check-row strong {
  display: block;
  color: var(--ez-text-main);
  font-size: 13px;
  font-weight: 800;
}

.health-check-row small {
  display: block;
  margin-top: 2px;
}

.health-tone--brand {
  --health-tone-color: var(--ez-brand);
  --health-tone-bg: var(--ez-brand-soft);
  --health-tone-border: var(--ez-brand-border);
}

.health-tone--success {
  --health-tone-color: var(--ez-success);
  --health-tone-bg: rgba(18, 185, 129, 0.1);
  --health-tone-border: rgba(18, 185, 129, 0.22);
}

.health-tone--warning {
  --health-tone-color: var(--ez-warning);
  --health-tone-bg: var(--ez-warning-bg);
  --health-tone-border: rgba(245, 158, 11, 0.24);
}

.health-tone--danger {
  --health-tone-color: var(--ez-danger);
  --health-tone-bg: var(--ez-danger-bg);
  --health-tone-border: rgba(239, 68, 68, 0.22);
}

.health-tone--neutral {
  --health-tone-color: var(--ez-text-placeholder);
  --health-tone-bg: var(--ez-surface-subtle);
  --health-tone-border: var(--ez-border-light);
}

@keyframes health-pulse {
  0%,
  100% {
    transform: scale(1);
  }

  50% {
    transform: scale(1.04);
  }
}

@media (max-width: 1180px) {
  .health-summary-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .health-content-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 760px) {
  .health-service-grid,
  .health-summary-grid {
    grid-template-columns: 1fr;
  }

  .health-overview {
    padding: 16px;
  }

  .health-score-panel {
    width: 100%;
    justify-content: flex-start;
    border-top: 1px solid var(--ez-border-light);
    border-left: 0;
    padding-top: 16px;
    padding-left: 0;
  }

  .health-overview__main {
    flex-direction: column;
  }

  .health-overview__body {
    flex-direction: column;
    gap: 16px;
  }
}
</style>
