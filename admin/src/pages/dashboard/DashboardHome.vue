<script setup lang="ts">
import { NButton, NCard } from 'naive-ui'

const metrics = [
  { label: '总用户数', value: '12,846', hint: '较昨日 +128', color: '#18A058' },
  { label: '活跃角色', value: '18', hint: '权限策略 42 条', color: '#2080F0' },
  { label: '今日操作', value: '2,391', hint: '高风险 3 条', color: '#F0A020' },
  { label: '系统健康度', value: '99.98%', hint: 'API 平均 38ms', color: '#18A058' },
]

const trendBars = [
  { width: '78%', color: '#18A058' },
  { width: '92%', color: '#36AD6A' },
  { width: '70%', color: '#63E2B7' },
  { width: '94%', color: '#2080F0' },
]

const recentOperations = [
  { text: '09:18 林青 新增用户 chenwei', color: '#374151' },
  { text: '09:02 系统 同步菜单权限', color: '#374151' },
  { text: '08:41 周敏 导出操作日志', color: '#374151' },
  { text: '08:16 角色 devops 权限变更', color: '#D03050' },
]
</script>

<template>
  <main class="flex h-full flex-col gap-5 overflow-hidden">
    <section class="flex items-center justify-between">
      <div>
        <h1 class="text-[28px] font-bold text-[#111827]">工作台</h1>
        <p class="mt-1 text-sm text-[#6B7280]">今天是 2026-04-24，系统运行稳定。</p>
      </div>

      <NButton type="primary" color="#18A058"> 新建用户 </NButton>
    </section>

    <section class="grid gap-4 xl:grid-cols-4">
      <NCard
        v-for="item in metrics"
        :key="item.label"
        class="rounded-lg"
        :bordered="false"
        content-style="padding: 20px;"
      >
        <div class="flex flex-col gap-2">
          <p class="text-sm text-[#6B7280]">{{ item.label }}</p>
          <p class="text-3xl font-bold text-[#111827]">{{ item.value }}</p>
          <p class="text-[13px]" :style="{ color: item.color }">{{ item.hint }}</p>
        </div>
      </NCard>
    </section>

    <section class="grid min-h-0 flex-1 gap-5 xl:grid-cols-[minmax(0,1fr)_416px]">
      <NCard
        class="h-full rounded-lg"
        :bordered="false"
        content-style="height: 100%; padding: 22px;"
      >
        <section class="flex h-full flex-col gap-4.5">
          <h2 class="text-lg font-bold text-[#111827]">访问趋势</h2>

          <div class="grid flex-1 gap-3 rounded-md bg-[#F9FAFB] p-4.5">
            <div
              v-for="bar in trendBars"
              :key="`${bar.color}-${bar.width}`"
              class="h-6.5 rounded"
              :style="{ width: bar.width, background: bar.color }"
            />
          </div>

          <p class="text-sm text-[#6B7280]">近 7 日访问量稳定上升，异常请求占比 0.06%。</p>
        </section>
      </NCard>

      <NCard
        class="h-full rounded-lg"
        :bordered="false"
        content-style="height: 100%; padding: 22px;"
      >
        <section class="flex h-full flex-col gap-3.5">
          <h2 class="text-lg font-bold text-[#111827]">最近操作</h2>

          <p
            v-for="item in recentOperations"
            :key="item.text"
            class="text-sm"
            :style="{ color: item.color }"
          >
            {{ item.text }}
          </p>
        </section>
      </NCard>
    </section>

    <section>
      <NButton tertiary class="w-45"> 用户管理 </NButton>
    </section>
  </main>
</template>
