<script setup lang="ts">
import { NAlert, NCard, NTag } from 'naive-ui'

import EzActionButton from '@/components/ez/EzActionButton.vue'
import PageHeader from '@/components/PageHeader.vue'
import { useHealthPage } from '../composables/useHealthPage'

const {
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
} = useHealthPage()
</script>

<template>
  <main class="admin-page admin-page-scroll">
    <section class="admin-page-section">
      <PageHeader title="系统状态" description="登录后检查后台运行环境、数据库和 Redis 的连通性。">
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
        class="rounded-[var(--ez-radius-sm)]"
        :bordered="false"
      >
        {{ errorMessage }}
      </NAlert>

      <div class="grid gap-4 xl:grid-cols-[minmax(0,1.4fr)_minmax(280px,0.9fr)]">
        <NCard class="rounded-[var(--ez-radius-sm)]" :bordered="false" content-class="p-6">
          <div class="flex h-full flex-col gap-5">
            <div class="flex items-start justify-between gap-4">
              <div>
                <p
                  class="text-sm font-medium uppercase tracking-[0.24em] text-[var(--ez-text-light)]"
                >
                  Runtime Snapshot
                </p>
                <h2 class="mt-2 text-2xl font-bold text-[var(--ez-text-main)]">
                  {{ health ? '核心依赖全部在线' : '等待首次检查结果' }}
                </h2>
                <p class="mt-2 text-sm leading-6 text-[var(--ez-text-sub)]">
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
                class="rounded-[var(--ez-radius-2xl)] border border-[var(--ez-border)] bg-[var(--ez-page-bg)] px-5 py-4"
              >
                <div class="flex items-center justify-between gap-3">
                  <span class="text-sm font-semibold text-[var(--ez-text-main)]">{{
                    item.label
                  }}</span>
                  <NTag :type="getStatusTagType(item.value)" size="small" round :bordered="false">
                    {{ formatStatusLabel(item.value) }}
                  </NTag>
                </div>
                <p class="mt-3 text-lg font-bold text-[var(--ez-text-main)]">
                  {{ item.value || 'pending' }}
                </p>
                <p class="mt-1 text-sm text-[var(--ez-text-sub)]">{{ item.description }}</p>
              </article>
            </div>

            <div
              class="rounded-[var(--ez-radius-2xl)] bg-[var(--ez-panel-dark)] px-5 py-4 text-[var(--ez-on-dark)]"
            >
              <div class="flex items-center justify-between gap-4">
                <div>
                  <p class="text-xs uppercase tracking-[0.2em] text-[var(--ez-on-dark-muted)]">
                    Last Check
                  </p>
                  <p class="mt-2 text-base font-semibold">{{ lastCheckedLabel }}</p>
                </div>
                <p class="text-sm text-[var(--ez-on-dark-sub)]">
                  {{ getStatusText(health?.database) }}
                </p>
              </div>
            </div>
          </div>
        </NCard>

        <section class="grid gap-4">
          <NCard class="rounded-[var(--ez-radius-sm)]" :bordered="false" content-class="p-5">
            <div class="flex flex-col gap-3">
              <div>
                <p class="text-sm font-semibold text-[var(--ez-text-main)]">运行环境</p>
                <p>当前后端 `app.env` 返回值，会随部署环境切换为 `dev` 或 `prod`。</p>
              </div>

              <div class="rounded-[var(--ez-radius-md)] bg-[var(--ez-page-bg)] px-4 py-3">
                <p class="text-xs uppercase tracking-[0.18em] text-[var(--ez-text-light)]">
                  Environment
                </p>
                <p class="mt-2 text-2xl font-bold text-[var(--ez-text-main)]">
                  {{ health?.env || 'unknown' }}
                </p>
              </div>
            </div>
          </NCard>

          <NCard class="rounded-[var(--ez-radius-sm)]" :bordered="false" content-class="p-5">
            <div class="flex flex-col gap-3">
              <div>
                <p class="text-sm font-semibold text-[var(--ez-text-main)]">接口职责</p>
                <p>同样是健康检查，公开探针和后台菜单入口分别服务于不同场景。</p>
              </div>

              <article
                v-for="endpoint in endpointCards"
                :key="endpoint.path"
                class="rounded-[var(--ez-radius-md)] border border-[var(--ez-border)] px-4 py-3"
              >
                <div class="flex items-center justify-between gap-3">
                  <span class="font-semibold text-[var(--ez-text-main)]">{{ endpoint.title }}</span>
                  <code
                    class="rounded bg-[var(--ez-page-bg)] px-2 py-1 text-xs text-[var(--ez-text-sub)]"
                  >
                    {{ endpoint.path }}
                  </code>
                </div>
                <p class="mt-2 text-sm leading-6 text-[var(--ez-text-sub)]">
                  {{ endpoint.description }}
                </p>
              </article>
            </div>
          </NCard>
        </section>
      </div>
    </section>
  </main>
</template>
