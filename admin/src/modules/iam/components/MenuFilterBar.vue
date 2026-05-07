<script setup lang="ts">
import type { SelectOption } from 'naive-ui'
import { NButton, NCard, NInput, NSelect, NSpace } from 'naive-ui'

import { MenuStatus, MenuType } from '@/modules/iam/types/menu'

const props = defineProps<{
  keyword: string
  status: 0 | MenuStatus
  statusOptions: SelectOption[]
  type: 0 | MenuType
  typeOptions: SelectOption[]
}>()

const emit = defineEmits<{
  reset: []
  'update:keyword': [value: string]
  'update:status': [value: 0 | MenuStatus]
  'update:type': [value: 0 | MenuType]
}>()
</script>

<template>
  <NCard :bordered="false" class="rounded-lg">
    <NSpace align="center" :wrap="true">
      <NInput
        :value="props.keyword"
        clearable
        placeholder="菜单名称 / 路由 / 权限标识"
        class="w-56"
        @update:value="emit('update:keyword', $event ?? '')"
      />
      <NSelect :value="props.type" :options="typeOptions" class="w-40" @update:value="emit('update:type', $event as 0 | MenuType)" />
      <NSelect :value="props.status" :options="statusOptions" class="w-40" @update:value="emit('update:status', $event as 0 | MenuStatus)" />
      <NButton @click="emit('reset')">重置</NButton>
    </NSpace>
  </NCard>
</template>
