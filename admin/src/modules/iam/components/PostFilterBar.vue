<script setup lang="ts">
import { NButton, NInput, NSelect, NSpace } from 'naive-ui'

import { STATUS_FILTER_OPTIONS } from '@/constants/status'
import type { PostPageQuery } from '../types/post-page'

defineEmits<{
  reset: []
  search: []
}>()

const query = defineModel<PostPageQuery>('query', { required: true })
</script>

<template>
  <div class="ez-toolbar">
    <NSpace align="center" :wrap="true">
      <NInput v-model:value="query.keyword" clearable placeholder="搜索岗位名称 / 编码" class="w-64" @keyup.enter="$emit('search')" />
      <NSelect v-model:value="query.status" :options="STATUS_FILTER_OPTIONS" class="w-36" />
      <NButton type="primary" @click="$emit('search')">查询</NButton>
      <NButton @click="$emit('reset')">重置</NButton>
    </NSpace>
  </div>
</template>
