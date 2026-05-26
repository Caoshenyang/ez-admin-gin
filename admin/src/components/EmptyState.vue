<script setup lang="ts">
import EzEmptyState from './ez/EzEmptyState.vue'

withDefaults(defineProps<{
  actionText?: string
  description?: string
  kind?: 'create' | 'search'
  title?: string
}>(), {
  actionText: '',
  description: '调整筛选条件后重试，或新建第一条业务数据。',
  kind: 'search',
  title: '当前没有数据',
})

defineEmits<{
  action: []
}>()

</script>

<template>
  <EzEmptyState
    :action-text="actionText"
    :description="description"
    :kind="kind"
    :title="title"
    @action="$emit('action')"
  />
</template>
