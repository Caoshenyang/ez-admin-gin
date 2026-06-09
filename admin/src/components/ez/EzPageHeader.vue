<script setup lang="ts">
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
  <header class="ez-page-header-shell">
    <div>
      <div class="ez-page-header-shell__title-row">
        <h1>{{ title }}</h1>
        <slot name="title-extra" />
      </div>
      <p v-if="description">{{ description }}</p>
    </div>
    <div v-if="$slots.actions" class="ez-page-header-shell__actions">
      <slot name="actions" />
    </div>
  </header>
</template>

<style scoped>
.ez-page-header-shell {
  display: flex;
  min-height: 40px;
  align-items: flex-start;
  justify-content: space-between;
  gap: 14px;
}

.ez-page-header-shell h1 {
  margin: 0;
  font-size: 19px;
  font-weight: 700;
  color: var(--ez-text-main);
  line-height: 26px;
}

.ez-page-header-shell__title-row {
  display: flex;
  min-width: 0;
  align-items: center;
  gap: 6px;
}

.ez-page-header-shell p {
  margin: 3px 0 0;
  font-size: 12px;
  color: var(--ez-text-secondary);
  line-height: 19px;
}

.ez-page-header-shell__actions {
  display: flex;
  flex-shrink: 0;
  flex-wrap: wrap;
  align-items: center;
  gap: 8px;
  justify-content: flex-end;
}

@media (max-width: 720px) {
  .ez-page-header-shell {
    align-items: flex-start;
    flex-direction: column;
  }

  .ez-page-header-shell__actions {
    width: 100%;
    justify-content: flex-start;
  }
}
</style>
