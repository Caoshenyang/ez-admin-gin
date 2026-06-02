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
  gap: 14px;
  border: 1px solid var(--ez-border);
  border-radius: var(--ez-radius-card);
  background: var(--ez-card-bg);
  padding: 14px 16px;
  box-shadow: var(--ez-shadow-sm);
}

.ez-search-panel__fields {
  display: flex;
  min-width: 0;
  flex: 1;
  flex-wrap: wrap;
  align-items: center;
  gap: 10px;
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
  width: 156px;
  min-width: 0;
  flex: 0 1 156px;
}

:slotted(.ez-search-field--xs) {
  width: 132px;
  flex-basis: 132px;
}

:slotted(.ez-search-field--sm) {
  width: 148px;
  flex-basis: 148px;
}

:slotted(.ez-search-field--md) {
  width: 172px;
  flex-basis: 172px;
}

:slotted(.ez-search-field--primary) {
  width: 268px;
  flex-basis: 268px;
}

:slotted(.ez-search-field--wide) {
  width: 320px;
  flex-basis: 320px;
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
