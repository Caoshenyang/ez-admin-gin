<script setup lang="ts">
import { NButton, NInput, NSelect } from 'naive-ui'

import EzSearchPanel from '@/components/ez/EzSearchPanel.vue'
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
  <EzSearchPanel>
    <NInput
      :value="query.username"
      clearable
      placeholder="用户名"
      class="ez-search-field ez-search-field--sm"
      @update:value="(value) => emit('update:username', value)"
      @keyup.enter="emit('search')"
    />
    <NInput
      :value="query.ip"
      clearable
      placeholder="IP 地址"
      class="ez-search-field ez-search-field--md"
      @update:value="(value) => emit('update:ip', value)"
      @keyup.enter="emit('search')"
    />
    <NSelect
      :value="query.status"
      :options="loginLogStatusOptions"
      class="ez-search-field ez-search-field--sm"
      @update:value="(value) => emit('update:status', Number(value ?? 0) as 0 | LoginLogStatus)"
    />

    <template #actions>
      <NButton type="primary" @click="emit('search')">查询</NButton>
      <NButton @click="emit('reset')">重置</NButton>
    </template>
  </EzSearchPanel>
</template>
