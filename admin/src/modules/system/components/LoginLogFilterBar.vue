<script setup lang="ts">
import { NButton, NInput, NSelect, NSpace } from 'naive-ui'

import type { LoginLogListQuery, LoginLogStatus } from '../types/login-log'
import { loginLogStatusOptions } from '../composables/login-log-page.utils'

defineProps<{
  query: LoginLogListQuery
}>()

const emit = defineEmits<{
  'update:ip': [value: string]
  'update:status': [value: 0 | LoginLogStatus]
  'update:username': [value: string]
  reset: []
  search: []
}>()
</script>

<template>
  <div class="ez-toolbar">
    <NSpace align="center" :wrap="true" :size="12">
      <NInput
        :value="query.username"
        clearable
        placeholder="用户名"
        class="w-40"
        @update:value="(value) => emit('update:username', value)"
        @keyup.enter="emit('search')"
      />
      <NInput
        :value="query.ip"
        clearable
        placeholder="IP 地址"
        class="w-44"
        @update:value="(value) => emit('update:ip', value)"
        @keyup.enter="emit('search')"
      />
      <NSelect
        :value="query.status"
        :options="loginLogStatusOptions"
        class="w-36"
        @update:value="(value) => emit('update:status', Number(value ?? 0) as 0 | LoginLogStatus)"
      />
      <NButton type="primary" @click="emit('search')">查询</NButton>
      <NButton @click="emit('reset')">重置</NButton>
    </NSpace>
  </div>
</template>
