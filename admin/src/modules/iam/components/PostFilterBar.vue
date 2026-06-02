<script setup lang="ts">
import { NButton, NInput, NSelect } from 'naive-ui'

import EzSearchPanel from '@/components/ez/EzSearchPanel.vue'
import { STATUS_FILTER_OPTIONS } from '@/constants/status'
import type { PostPageQuery } from '../types/post-page'

defineEmits<{
  reset: []
  search: []
}>()

const query = defineModel<PostPageQuery>('query', { required: true })
</script>

<template>
  <EzSearchPanel>
    <NInput v-model:value="query.keyword" clearable placeholder="搜索岗位名称 / 编码" class="ez-search-field ez-search-field--primary" @keyup.enter="$emit('search')" />
    <NSelect v-model:value="query.status" :options="STATUS_FILTER_OPTIONS" class="ez-search-field ez-search-field--sm" />

    <template #actions>
      <NButton type="primary" @click="$emit('search')">查询</NButton>
      <NButton @click="$emit('reset')">重置</NButton>
    </template>
  </EzSearchPanel>
</template>
