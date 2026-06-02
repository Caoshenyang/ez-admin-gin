<script setup lang="ts">
import type { DataTableColumns } from 'naive-ui'
import { NButton, NDataTable } from 'naive-ui'

import EzTableCard from '@/components/ez/EzTableCard.vue'
import TableStatsBar from '@/components/TableStatsBar.vue'
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
  <EzTableCard>
    <TableStatsBar>
      <span />
      <template #actions>
        <NButton text type="primary" @click="$emit('refresh')">刷新</NButton>
      </template>
    </TableStatsBar>

    <NDataTable
      class="h-full"
      :columns="columns"
      :data="items"
      :loading="loading"
      :pagination="false"
      :bordered="false"
      flex-height
    />
  </EzTableCard>
</template>
