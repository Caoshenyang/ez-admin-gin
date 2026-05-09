<script setup lang="ts">
import { NAlert, NButton, NCard, NTag } from 'naive-ui'
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
      <div class="flex items-center justify-between gap-4">
        <div>
          <h1 class="text-[28px] font-bold text-[#0F172A]">系统状态</h1>
          <p>登录后检查后台运行环境、数据库和 Redis 的连通性。</p>
        </div>

        <NButton type="primary" :loading="loading" @click="void loadHealth()">
          刷新状态
        </NButton>
      </div>

      <NAlert
        v-if="errorMessage"
        type="error"
        title="状态检查失败"
        class="rounded-lg"
        :bordered="false"
      >
        {{ errorMessage }}
      </NAlert>

      <div class="grid gap-4 xl:grid-cols-[minmax(0,1.4fr)_minmax(280px,0.9fr)]">
        <NCard class="rounded-lg" :bordered="false" content-class="p-6">
          <div class="flex h-full flex-col gap-5">
          <div class="flex items-start justify-between gap-4">
            <div>
              <p class="text-sm font-medium uppercase tracking-[0.24em] text-[#94A3B8]">
                Runtime Snapshot
              </p>
              <h2 class="mt-2 text-2xl font-bold text-[#0F172A]">
                {{ health ? '核心依赖全部在线' : '等待首次检查结果' }}
              </h2>
              <p class="mt-2 text-sm leading-6 text-[#64748B]">
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
              class="rounded-2xl border border-[#E6ECF3] bg-[#F8FAFC] px-5 py-4"
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
              <p class="mt-3 text-lg font-bold text-[#0F172A]">{{ item.value || 'pending' }}</p>
              <p class="mt-1 text-sm text-[#64748B]">{{ item.description }}</p>
            </article>
          </div>

          <div class="rounded-2xl bg-[#0F172A] px-5 py-4 text-white">
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
          <NCard class="rounded-lg" :bordered="false" content-class="p-5">
            <div class="flex flex-col gap-3">
              <div>
                <p class="text-sm font-semibold text-[#0F172A]">运行环境</p>
                <p>
                  当前后端 `app.env` 返回值，会随部署环境切换为 `dev` 或 `prod`。
                </p>
              </div>

              <div class="rounded-xl bg-[#F8FAFC] px-4 py-3">
                <p class="text-xs uppercase tracking-[0.18em] text-[#94A3B8]">Environment</p>
                <p class="mt-2 text-2xl font-bold text-[#0F172A]">{{ health?.env || 'unknown' }}</p>
              </div>
            </div>
          </NCard>

          <NCard class="rounded-lg" :bordered="false" content-class="p-5">
            <div class="flex flex-col gap-3">
              <div>
                <p class="text-sm font-semibold text-[#0F172A]">接口职责</p>
                <p>
                  同样是健康检查，公开探针和后台菜单入口分别服务于不同场景。
                </p>
              </div>

              <article
                v-for="endpoint in endpointCards"
                :key="endpoint.path"
                class="rounded-xl border border-[#E6ECF3] px-4 py-3"
              >
                <div class="flex items-center justify-between gap-3">
                  <span class="font-semibold text-[#0F172A]">{{ endpoint.title }}</span>
                  <code class="rounded bg-[#F8FAFC] px-2 py-1 text-xs text-[#475569]">
                    {{ endpoint.path }}
                  </code>
                </div>
                <p class="mt-2 text-sm leading-6 text-[#64748B]">{{ endpoint.description }}</p>
              </article>
            </div>
          </NCard>
        </section>
      </div>
    </section>
  </main>
</template>
