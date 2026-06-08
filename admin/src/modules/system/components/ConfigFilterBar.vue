<script setup lang="ts">
import { NInput, NSelect } from 'naive-ui'

import EzActionButton from '@/components/ez/EzActionButton.vue'
import { STATUS_FILTER_OPTIONS } from '@/constants/status'
import EzSearchPanel from '@/components/ez/EzSearchPanel.vue'
import type { ConfigListQuery } from '../types/config'

defineEmits<{
  reset: []
  search: []
}>()

const query = defineModel<ConfigListQuery>('query', { required: true })
</script>

<template>
  <EzSearchPanel>
    <NInput
      v-model:value="query.keyword"
      clearable
      placeholder="键 / 名称"
      class="ez-search-field ez-search-field--primary"
      @keyup.enter="$emit('search')"
    />
    <NInput
      v-model:value="query.group_code"
      clearable
      placeholder="分组"
      class="ez-search-field ez-search-field--sm"
      @keyup.enter="$emit('search')"
    />
    <NSelect
      v-model:value="query.status"
      :options="STATUS_FILTER_OPTIONS"
      class="ez-search-field ez-search-field--sm"
    />

    <template #actions>
      <EzActionButton kind="search" label="查询" type="primary" @click="$emit('search')" />
      <EzActionButton kind="reset" label="重置" @click="$emit('reset')" />
    </template>
  </EzSearchPanel>
</template>
