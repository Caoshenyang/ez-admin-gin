<script setup lang="ts">
import EzToolbar from './ez/EzToolbar.vue'

defineSlots<{
  actions?: () => unknown
  default?: () => unknown
}>()
</script>

<template>
  <EzToolbar>
    <div>
      <slot />
    </div>
    <template v-if="$slots.actions" #actions>
      <slot name="actions" />
    </template>
  </EzToolbar>
</template>
