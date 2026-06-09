<script setup lang="ts">
import { NInput, NSelect } from 'naive-ui'

import EzActionButton from '@/components/ez/EzActionButton.vue'
import EzSearchPanel from '@/components/ez/EzSearchPanel.vue'
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
  <EzSearchPanel>
    <NInput
      :value="props.username"
      clearable
      placeholder="操作人"
      class="ez-search-field ez-search-field--sm"
      @update:value="emit('update:username', $event ?? '')"
      @keyup.enter="emit('search')"
    />
    <NSelect
      :value="props.method"
      :options="methodOptions"
      class="ez-search-field ez-search-field--xs"
      @update:value="emit('update:method', $event ?? '')"
    />
    <NInput
      :value="props.path"
      clearable
      placeholder="请求路径"
      class="ez-search-field ez-search-field--wide"
      @update:value="emit('update:path', $event ?? '')"
      @keyup.enter="emit('search')"
    />
    <NSelect
      :value="props.success"
      :options="successOptions"
      class="ez-search-field ez-search-field--sm"
      @update:value="emit('update:success', $event ?? '')"
    />

    <template #actions>
      <EzActionButton kind="search" label="查询" type="primary" @click="emit('search')" />
      <EzActionButton kind="reset" label="重置" @click="emit('reset')" />
    </template>
  </EzSearchPanel>
</template>
