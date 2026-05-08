<script setup lang="ts">
import { NButton, NCard, NInput, NSelect } from 'naive-ui'

import { methodOptions, successOptions } from '../composables/operation-log-page.utils'

const props = defineProps<{
  method: string
  path: string
  success: string
  username: string
}>()

const emit = defineEmits<{
  reset: []
  search: []
  'update:method': [value: string]
  'update:path': [value: string]
  'update:success': [value: string]
  'update:username': [value: string]
}>()
</script>

<template>
  <NCard :bordered="false" class="rounded-lg">
    <div class="grid gap-3 xl:grid-cols-[180px_150px_minmax(0,1fr)_150px_auto]">
      <NInput
        :value="props.username"
        clearable
        placeholder="操作人"
        @update:value="emit('update:username', $event ?? '')"
        @keyup.enter="emit('search')"
      />
      <NSelect :value="props.method" :options="methodOptions" @update:value="emit('update:method', $event ?? '')" />
      <NInput
        :value="props.path"
        clearable
        placeholder="请求路径"
        @update:value="emit('update:path', $event ?? '')"
        @keyup.enter="emit('search')"
      />
      <NSelect :value="props.success" :options="successOptions" @update:value="emit('update:success', $event ?? '')" />
      <div class="flex gap-3 xl:justify-end">
        <NButton type="primary" @click="emit('search')">查询</NButton>
        <NButton @click="emit('reset')">重置</NButton>
      </div>
    </div>
  </NCard>
</template>
