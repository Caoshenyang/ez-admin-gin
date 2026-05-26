<script setup lang="ts">
import { NButton, NInput, NSelect } from 'naive-ui'

import { STATUS_FILTER_OPTIONS } from '@/constants/status'
import EzSearchPanel from '@/components/ez/EzSearchPanel.vue'
import type { ConfigListQuery } from '../types/config'

defineEmits<{
  reset: []
  search: []
}>()

const query = defineModel<ConfigListQuery>('query', { required: true })
</script>

<template>
  <EzSearchPanel>
      <NInput v-model:value="query.keyword" clearable placeholder="键 / 名称" class="w-56" @keyup.enter="$emit('search')" />
      <NInput v-model:value="query.group_code" clearable placeholder="分组" class="w-44" @keyup.enter="$emit('search')" />
      <NSelect v-model:value="query.status" :options="STATUS_FILTER_OPTIONS" class="w-36" />
      <template #actions>
        <NButton type="primary" @click="$emit('search')">查询</NButton>
        <NButton @click="$emit('reset')">重置</NButton>
      </template>
  </EzSearchPanel>
</template>
