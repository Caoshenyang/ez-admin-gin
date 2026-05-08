<script setup lang="ts">
import type { DataTableColumns } from 'naive-ui'
import { NButton, NCard, NDataTable } from 'naive-ui'

import type { PostItem } from '../types/post'

defineProps<{
  columns: DataTableColumns<PostItem>
  items: PostItem[]
  loading: boolean
}>()

defineEmits<{
  refresh: []
}>()
</script>

<template>
  <NCard class="min-h-0 flex-1 rounded-lg" :bordered="false" content-style="display: flex; height: 100%; min-height: 0; flex-direction: column; padding: 0;">
    <div class="flex items-center justify-end border-b border-[#E5E7EB] px-4 py-3">
      <NButton text type="primary" @click="$emit('refresh')">刷新</NButton>
    </div>

    <NDataTable
      class="h-full"
      style="height: 100%"
      :columns="columns"
      :data="items"
      :loading="loading"
      :pagination="false"
      :bordered="false"
      flex-height
    />
  </NCard>
</template>
