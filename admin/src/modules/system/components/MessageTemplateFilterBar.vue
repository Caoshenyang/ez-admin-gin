<script setup lang="ts">
import { NInput, NSelect } from 'naive-ui'

import EzActionButton from '@/components/ez/EzActionButton.vue'
import EzSearchPanel from '@/components/ez/EzSearchPanel.vue'
import { STATUS_FILTER_OPTIONS } from '@/constants/status'
import {
  MESSAGE_TEMPLATE_TYPE_FILTER_OPTIONS,
} from '../composables/message-page.utils'
import type { MessageTemplateListQuery } from '../types/message'

defineEmits<{
  reset: []
  search: []
}>()

const query = defineModel<MessageTemplateListQuery>('query', { required: true })
</script>

<template>
  <EzSearchPanel>
    <NInput
      v-model:value="query.keyword"
      clearable
      placeholder="编码 / 名称 / 标题"
      class="ez-search-field ez-search-field--primary"
      @keyup.enter="$emit('search')"
    />
    <NSelect
      v-model:value="query.type"
      :options="MESSAGE_TEMPLATE_TYPE_FILTER_OPTIONS"
      class="ez-search-field ez-search-field--sm"
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
