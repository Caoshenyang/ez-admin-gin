<script setup lang="ts">
import type { SelectOption } from 'naive-ui'
import { NButton, NCard, NInput, NSelect, NSpace } from 'naive-ui'

import { DepartmentStatus } from '@/modules/iam/types/department'

const props = defineProps<{
  keyword: string
  status: 0 | DepartmentStatus
  statusOptions: SelectOption[]
}>()

const emit = defineEmits<{
  reset: []
  search: []
  'update:keyword': [value: string]
  'update:status': [value: 0 | DepartmentStatus]
}>()
</script>

<template>
  <NCard :bordered="false" class="rounded-lg">
    <NSpace align="center" :wrap="true">
      <NInput
        :value="props.keyword"
        clearable
        placeholder="搜索部门名称 / 编码"
        class="w-64"
        @update:value="emit('update:keyword', $event ?? '')"
        @keyup.enter="emit('search')"
      />
      <NSelect :value="props.status" :options="statusOptions" class="w-36" @update:value="emit('update:status', $event as 0 | DepartmentStatus)" />
      <NButton type="primary" @click="emit('search')">查询</NButton>
      <NButton @click="emit('reset')">重置</NButton>
    </NSpace>
  </NCard>
</template>
