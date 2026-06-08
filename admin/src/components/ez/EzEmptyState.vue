<script setup lang="ts">
import { AddOutline, SearchOutline } from '@vicons/ionicons5'
import { NButton, NIcon } from 'naive-ui'
import { computed } from 'vue'

const props = withDefaults(
  defineProps<{
    actionText?: string
    description?: string
    kind?: 'create' | 'search'
    title?: string
  }>(),
  {
    actionText: '',
    description: '调整筛选条件后重试，或新建第一条业务数据。',
    kind: 'search',
    title: '当前没有数据',
  },
)

defineEmits<{
  action: []
}>()

const iconComponent = computed(() => (props.kind === 'create' ? AddOutline : SearchOutline))
</script>

<template>
  <div class="ez-empty-state">
    <div class="ez-empty-state__icon">
      <NIcon :component="iconComponent" :size="28" />
    </div>
    <div class="ez-empty-state__copy">
      <p>{{ title }}</p>
      <span>{{ description }}</span>
    </div>
    <NButton v-if="actionText" type="primary" ghost @click="$emit('action')">
      {{ actionText }}
    </NButton>
  </div>
</template>

<style scoped>
.ez-empty-state {
  display: flex;
  width: 100%;
  min-height: var(--empty-state-min-height, 168px);
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 12px;
  overflow: hidden;
  border: 1px dashed var(--ez-component-border);
  border-radius: var(--ez-radius-control);
  background: var(--ez-surface-subtle);
  padding: var(--empty-state-padding, 24px 16px);
  text-align: center;
}

.ez-empty-state__icon {
  display: grid;
  width: 48px;
  height: 48px;
  place-items: center;
  border: 1px solid var(--ez-component-border);
  border-radius: var(--ez-radius-control);
  background: var(--ez-card-bg);
  color: var(--ez-primary);
}

.ez-empty-state__copy p {
  margin: 0;
  color: var(--ez-text-main);
  font-size: 14px;
  font-weight: 600;
}

.ez-empty-state__copy span {
  display: block;
  max-width: 420px;
  margin-top: 5px;
  color: var(--ez-text-secondary);
  font-size: 12px;
  line-height: 1.6;
}
</style>
