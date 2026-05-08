<script setup lang="ts">
import { NButton, NCard, NInput, NSelect, NSpace } from 'naive-ui'

import type { LoginLogListQuery } from '../types/login-log'
import { loginLogStatusOptions } from '../composables/login-log-page.utils'

defineProps<{
  query: LoginLogListQuery
}>()

const emit = defineEmits<{
  reset: []
  search: []
}>()
</script>

<template>
  <NCard :bordered="false" class="rounded-lg">
    <NSpace align="center" :wrap="true" :size="12">
      <NInput
        v-model:value="query.username"
        clearable
        placeholder="用户名"
        class="w-40"
        @keyup.enter="emit('search')"
      />
      <NInput
        v-model:value="query.ip"
        clearable
        placeholder="IP 地址"
        class="w-44"
        @keyup.enter="emit('search')"
      />
      <NSelect v-model:value="query.status" :options="loginLogStatusOptions" class="w-36" />
      <NButton type="primary" @click="emit('search')">查询</NButton>
      <NButton @click="emit('reset')">重置</NButton>
    </NSpace>
  </NCard>
</template>
