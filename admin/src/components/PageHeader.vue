<script setup lang="ts">
import EzPageHeader from './ez/EzPageHeader.vue'

defineProps<{
  description?: string
  title: string
}>()

defineSlots<{
  actions?: () => unknown
  'title-extra'?: () => unknown
}>()
</script>

<template>
  <EzPageHeader :title="title" :description="description">
    <template v-if="$slots['title-extra']" #title-extra>
      <slot name="title-extra" />
    </template>

    <template v-if="$slots.actions" #actions>
      <slot name="actions" />
    </template>
  </EzPageHeader>
</template>
