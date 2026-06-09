<script setup lang="ts">
defineSlots<{
  actions?: () => unknown
  default?: () => unknown
}>()
</script>

<template>
  <section class="ez-search-panel">
    <div class="ez-search-panel__fields">
      <slot />
    </div>
    <div v-if="$slots.actions" class="ez-search-panel__actions">
      <slot name="actions" />
    </div>
  </section>
</template>

<style scoped>
.ez-search-panel {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  overflow: hidden;
  border: 1px solid var(--ez-component-border);
  border-radius: var(--ez-radius-control);
  background: var(--ez-card-bg);
  padding: 11px 12px;
  box-shadow: var(--ez-component-shadow);
}

.ez-search-panel__fields {
  display: flex;
  min-width: 0;
  flex: 1;
  flex-wrap: wrap;
  align-items: center;
  gap: 8px;
}

.ez-search-panel__actions {
  display: flex;
  flex-shrink: 0;
  flex-wrap: wrap;
  align-items: center;
  justify-content: flex-end;
  gap: 8px;
}

:slotted(.ez-search-field) {
  width: 148px;
  min-width: 0;
  flex: 0 1 148px;
  --n-height: 32px;
}

:slotted(.ez-search-field--xs) {
  width: 124px;
  flex-basis: 124px;
}

:slotted(.ez-search-field--sm) {
  width: 140px;
  flex-basis: 140px;
}

:slotted(.ez-search-field--md) {
  width: 164px;
  flex-basis: 164px;
}

:slotted(.ez-search-field--primary) {
  width: 248px;
  flex-basis: 248px;
}

:slotted(.ez-search-field--wide) {
  width: 292px;
  flex-basis: 292px;
}

@media (max-width: 720px) {
  .ez-search-panel {
    align-items: stretch;
    flex-direction: column;
  }

  :slotted(.ez-search-field) {
    width: 100%;
    flex-basis: auto;
  }

  .ez-search-panel__actions {
    justify-content: flex-start;
  }
}
</style>
