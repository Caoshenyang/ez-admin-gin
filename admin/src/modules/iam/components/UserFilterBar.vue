<script setup lang="ts">
import type { SelectOption } from 'naive-ui'
import { NButton, NInput, NSelect, NSpace } from 'naive-ui'

defineProps<{
  keyword: string
  roleId: number
  roleOptions: SelectOption[]
  status: number
  statusOptions: SelectOption[]
}>()

const emit = defineEmits<{
  'update:keyword': [value: string]
  'update:roleId': [value: number]
  'update:status': [value: number]
  search: []
  reset: []
}>()
</script>

<template>
  <div class="ez-toolbar">
    <NSpace align="center" :wrap="true" :size="12">
      <NInput
        :value="keyword"
        clearable
        placeholder="用户名 / 手机号"
        class="w-64"
        @update:value="(value) => emit('update:keyword', value)"
        @keyup.enter="emit('search')"
      />
      <NSelect
        :value="roleId"
        :options="roleOptions"
        class="w-40"
        @update:value="(value) => emit('update:roleId', Number(value ?? 0))"
      />
      <NSelect
        :value="status"
        :options="statusOptions"
        class="w-36"
        @update:value="(value) => emit('update:status', Number(value ?? 0))"
      />
      <NButton type="primary" @click="emit('search')">查询</NButton>
      <NButton @click="emit('reset')">重置</NButton>
    </NSpace>
  </div>
</template>
