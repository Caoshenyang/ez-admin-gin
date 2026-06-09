<script setup lang="ts">
import { BulbOutline, InformationCircleOutline, TrailSignOutline } from '@vicons/ionicons5'
import { NDrawer, NDrawerContent, NIcon } from 'naive-ui'

interface PageHelpSection {
  description?: string
  items: string[]
  title: string
}

defineProps<{
  flow: string[]
  sections: PageHelpSection[]
  show: boolean
  title: string
  warnings?: string[]
}>()

defineEmits<{
  'update:show': [value: boolean]
}>()
</script>

<template>
  <NDrawer
    :show="show"
    :width="460"
    placement="right"
    class="page-help-drawer"
    @update:show="(value) => $emit('update:show', value)"
  >
    <NDrawerContent closable body-content-class="page-help-drawer__body">
      <template #header>
        <div class="page-help-drawer__title">
          <span class="page-help-drawer__title-icon">
            <NIcon :component="InformationCircleOutline" :size="20" />
          </span>
          <span>{{ title }}</span>
        </div>
      </template>

      <div class="page-help-drawer__content">
        <section class="page-help-block">
          <div class="page-help-block__head">
            <span class="page-help-block__icon page-help-block__icon--flow">
              <NIcon :component="TrailSignOutline" :size="18" />
            </span>
            <h3>推荐流程</h3>
          </div>

          <ol class="page-help-flow">
            <li v-for="(item, index) in flow" :key="item">
              <span>{{ index + 1 }}</span>
              <p>{{ item }}</p>
            </li>
          </ol>
        </section>

        <section v-for="section in sections" :key="section.title" class="page-help-block">
          <div class="page-help-block__head">
            <span class="page-help-block__icon">
              <NIcon :component="BulbOutline" :size="18" />
            </span>
            <div>
              <h3>{{ section.title }}</h3>
              <p v-if="section.description">{{ section.description }}</p>
            </div>
          </div>

          <ul class="page-help-list">
            <li v-for="item in section.items" :key="item">{{ item }}</li>
          </ul>
        </section>

        <section v-if="warnings?.length" class="page-help-block page-help-block--warning">
          <div class="page-help-block__head">
            <span class="page-help-block__icon page-help-block__icon--warning">
              <NIcon :component="InformationCircleOutline" :size="18" />
            </span>
            <h3>注意事项</h3>
          </div>

          <ul class="page-help-list">
            <li v-for="item in warnings" :key="item">{{ item }}</li>
          </ul>
        </section>
      </div>
    </NDrawerContent>
  </NDrawer>
</template>

<style scoped>
.page-help-drawer__title {
  display: flex;
  min-width: 0;
  align-items: center;
  gap: 8px;
  color: var(--ez-text-heading);
  font-size: 16px;
  font-weight: 700;
}

.page-help-drawer__title-icon,
.page-help-block__icon {
  display: inline-flex;
  flex: 0 0 auto;
  align-items: center;
  justify-content: center;
  width: 30px;
  height: 30px;
  border-radius: 8px;
  background: var(--ez-info-bg);
  color: var(--ez-info-text);
}

.page-help-drawer__content {
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding-bottom: 10px;
}

.page-help-block {
  border: 1px solid var(--ez-component-border);
  border-radius: var(--ez-radius-sm);
  background: var(--ez-card-bg);
  padding: 14px;
}

.page-help-block--warning {
  border-color: rgba(245, 158, 11, 0.28);
  background: var(--ez-warning-bg);
}

.page-help-block__head {
  display: flex;
  align-items: flex-start;
  gap: 10px;
}

.page-help-block__head h3 {
  margin: 4px 0 0;
  color: var(--ez-text-heading);
  font-size: 14px;
  font-weight: 700;
  line-height: 20px;
}

.page-help-block__head p {
  margin: 2px 0 0;
  color: var(--ez-text-muted);
  font-size: 12px;
  line-height: 18px;
}

.page-help-block__icon {
  width: 28px;
  height: 28px;
  background: var(--ez-brand-soft);
  color: var(--ez-brand);
}

.page-help-block__icon--flow {
  background: rgba(18, 185, 129, 0.12);
  color: var(--ez-brand-green);
}

.page-help-block__icon--warning {
  background: rgba(245, 158, 11, 0.14);
  color: var(--ez-warning-text);
}

.page-help-flow {
  display: grid;
  gap: 9px;
  margin: 12px 0 0;
  padding: 0;
  list-style: none;
}

.page-help-flow li {
  display: grid;
  grid-template-columns: 24px minmax(0, 1fr);
  gap: 8px;
  align-items: flex-start;
}

.page-help-flow span {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 22px;
  height: 22px;
  border-radius: 999px;
  background: var(--ez-surface-subtle);
  color: var(--ez-text-muted);
  font-size: 12px;
  font-weight: 700;
}

.page-help-flow p,
.page-help-list li {
  margin: 0;
  color: var(--ez-text-body);
  font-size: 13px;
  line-height: 20px;
}

.page-help-list {
  display: grid;
  gap: 8px;
  margin: 12px 0 0;
  padding: 0 0 0 18px;
}

.page-help-list li::marker {
  color: var(--ez-brand);
}

:global(.page-help-drawer__body) {
  padding: 14px 16px 18px;
}

@media (max-width: 720px) {
  :global(.page-help-drawer) {
    width: min(100vw, 460px) !important;
  }
}
</style>
