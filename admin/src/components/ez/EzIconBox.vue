<script setup lang="ts">
import { NIcon } from 'naive-ui'
import type { Component } from 'vue'

withDefaults(
  defineProps<{
    color?: string
    icon: Component
    softColor?: string
  }>(),
  {
    color: 'var(--ez-primary)',
    softColor: 'var(--ez-primary-light)',
  },
)
</script>

<template>
  <span class="ez-icon-box" :style="{ color, background: softColor }">
    <NIcon :component="icon" :size="22" />
  </span>
</template>

<style scoped>
.ez-icon-box {
  display: inline-flex;
  width: 38px;
  height: 38px;
  flex-shrink: 0;
  align-items: center;
  justify-content: center;
  border: 1px solid var(--ez-component-border);
  border-radius: var(--ez-radius-control);
  box-shadow: 0 4px 10px rgba(15, 23, 42, 0.04);
}
</style>
