<script setup lang="ts">
import { NInput, NSelect } from 'naive-ui'

import EzActionButton from '@/components/ez/EzActionButton.vue'
import EzSearchPanel from '@/components/ez/EzSearchPanel.vue'
import { STATUS_FILTER_OPTIONS } from '@/constants/status'
import type { MailTemplateListQuery } from '../types/mail'

defineEmits<{
  reset: []
  search: []
}>()

const query = defineModel<MailTemplateListQuery>('query', { required: true })
</script>

<template>
  <EzSearchPanel>
    <NInput
      v-model:value="query.keyword"
      clearable
      placeholder="模板编码 / 名称 / 主题"
      class="ez-search-field ez-search-field--primary"
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
