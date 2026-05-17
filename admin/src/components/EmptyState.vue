<script setup lang="ts">
import { AddOutline, SearchOutline } from '@vicons/ionicons5'
import { NButton, NIcon } from 'naive-ui'
import { computed } from 'vue'

const props = withDefaults(defineProps<{
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

const iconComponent = computed(() => (props.kind === 'create' ? AddOutline : SearchOutline))
</script>

<template>
  <div class="empty-state">
    <div class="empty-state__icon">
      <NIcon :component="iconComponent" :size="28" />
    </div>
    <div class="empty-state__copy">
      <p>{{ title }}</p>
      <span>{{ description }}</span>
    </div>
    <NButton v-if="actionText" type="primary" ghost @click="$emit('action')">
      {{ actionText }}
    </NButton>
  </div>
</template>

<style scoped>
.empty-state {
  display: flex;
  width: 100%;
  min-height: var(--empty-state-min-height, 188px);
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 14px;
  border-radius: var(--ez-radius-lg);
  border: 1px dashed var(--ez-border);
  background:
    radial-gradient(circle at 50% 0%, var(--ez-brand-muted), transparent 34%),
    linear-gradient(180deg, var(--ez-card-bg), var(--ez-page-bg));
  padding: var(--empty-state-padding, 28px 18px);
  text-align: center;
}

.empty-state__icon {
  display: grid;
  height: 56px;
  width: 56px;
  place-items: center;
  border-radius: var(--ez-radius-2xl);
  background: var(--ez-brand-soft);
  color: var(--ez-brand);
}

.empty-state__copy p {
  margin: 0;
  color: var(--ez-text-main);
  font-size: var(--ez-text-md);
  font-weight: 700;
}

.empty-state__copy span {
  margin-top: 6px;
  display: block;
  max-width: 420px;
  color: var(--ez-text-sub);
  font-size: var(--ez-text-sm);
  line-height: 1.6;
}
</style>
