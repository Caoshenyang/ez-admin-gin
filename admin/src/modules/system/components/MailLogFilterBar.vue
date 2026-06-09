<script setup lang="ts">
import type { SelectOption } from 'naive-ui'
import { NInput, NSelect } from 'naive-ui'

import EzActionButton from '@/components/ez/EzActionButton.vue'
import EzSearchPanel from '@/components/ez/EzSearchPanel.vue'
import type { MailLogListQuery } from '../types/mail'

defineProps<{
  accountOptions: SelectOption[]
  statusOptions: SelectOption[]
  templateOptions: SelectOption[]
}>()

defineEmits<{
  reset: []
  search: []
}>()

const query = defineModel<MailLogListQuery>('query', { required: true })
</script>

<template>
  <EzSearchPanel>
    <NInput
      v-model:value="query.keyword"
      clearable
      placeholder="邮件主题 / 收件人 / 发件邮箱"
      class="ez-search-field ez-search-field--primary"
      @keyup.enter="$emit('search')"
    />
    <NSelect
      v-model:value="query.status"
      :options="statusOptions"
      class="ez-search-field ez-search-field--sm"
    />
    <NSelect
      v-model:value="query.account_id"
      clearable
      filterable
      placeholder="邮箱账号"
      :options="accountOptions"
      class="ez-search-field ez-search-field--md"
    />
    <NSelect
      v-model:value="query.template_code"
      clearable
      filterable
      placeholder="邮件模板"
      :options="templateOptions"
      class="ez-search-field ez-search-field--md"
    />

    <template #actions>
      <EzActionButton kind="search" label="查询" type="primary" @click="$emit('search')" />
      <EzActionButton kind="reset" label="重置" @click="$emit('reset')" />
    </template>
  </EzSearchPanel>
</template>
