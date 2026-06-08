<script setup lang="ts">
import type { SelectOption } from 'naive-ui'
import { NInput, NSelect } from 'naive-ui'
import EzActionButton from '@/components/ez/EzActionButton.vue'
import EzSearchPanel from '@/components/ez/EzSearchPanel.vue'

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
  <EzSearchPanel>
    <NInput
      :value="keyword"
      clearable
      placeholder="用户名 / 手机号"
      class="ez-search-field ez-search-field--primary"
      @update:value="(value) => emit('update:keyword', value)"
      @keyup.enter="emit('search')"
    />
    <NSelect
      :value="roleId"
      :options="roleOptions"
      class="ez-search-field ez-search-field--md"
      @update:value="(value) => emit('update:roleId', Number(value ?? 0))"
    />
    <NSelect
      :value="status"
      :options="statusOptions"
      class="ez-search-field ez-search-field--sm"
      @update:value="(value) => emit('update:status', Number(value ?? 0))"
    />

    <template #actions>
      <EzActionButton kind="search" label="查询" type="primary" @click="emit('search')" />
      <EzActionButton kind="reset" label="重置" @click="emit('reset')" />
    </template>
  </EzSearchPanel>
</template>
