<script setup lang="ts">
import { NCard, NTag } from 'naive-ui'

import EmptyState from '@/components/EmptyState.vue'
import PageHeader from '@/components/PageHeader.vue'

defineProps<{
  title: string
  description: string
}>()
</script>

<template>
  <main class="admin-page">
    <section class="admin-page-section gap-6">
      <PageHeader :title="title" :description="description">
        <template #actions>
          <NTag type="info" :bordered="false">后续接入</NTag>
        </template>
      </PageHeader>

      <NCard class="ez-card min-h-0 flex-1 rounded-[var(--ez-radius-sm)]" :bordered="false" content-class="ez-card-content-fill">
        <div class="flex h-full items-center justify-center p-4">
          <EmptyState
            title="本页待接入"
            description="当前先验证后台布局、路由出口和工作标签，后续小节会继续补齐业务内容。"
          />
        </div>
      </NCard>
    </section>
  </main>
</template>
