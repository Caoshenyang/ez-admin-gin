<script setup lang="ts">
import { NCard, NTag } from 'naive-ui'

import EzActionButton from '@/components/ez/EzActionButton.vue'
import PageHeader from '@/components/PageHeader.vue'
import ConfigFilterBar from '../components/ConfigFilterBar.vue'
import ConfigFormModal from '../components/ConfigFormModal.vue'
import ConfigTable from '../components/ConfigTable.vue'
import { useConfigPage } from '../composables/useConfigPage'

const {
  activeConfigCategory,
  activeConfigGroup,
  canUse,
  configCategories,
  columns,
  configs,
  formMode,
  formModel,
  formRef,
  formVisible,
  handlePageChange,
  handlePageSizeChange,
  handleReset,
  handleSearch,
  handleSubmit,
  load,
  loading,
  openCreate,
  query,
  rules,
  saving,
  selectConfigCategory,
  submitForm,
  total,
} = useConfigPage()
</script>

<template>
  <main class="admin-page">
    <section class="admin-page-section">
      <PageHeader title="配置管理" description="维护系统键值配置，按分组归类管理。">
        <template #actions>
          <EzActionButton
            v-if="canUse('system:config:create')"
            kind="add"
            :label="activeConfigGroup ? '新增本组配置' : '新增配置'"
            type="primary"
            @click="() => openCreate()"
          />
        </template>
      </PageHeader>

      <NCard class="config-group-card" :bordered="false">
        <div class="config-group-head">
          <div class="config-group-title">
            <span>配置分组</span>
            <strong>{{ activeConfigCategory.label }}</strong>
          </div>
          <div class="config-group-meta">
            <NTag size="small" :bordered="false" type="info">
              {{ activeConfigGroup || '全部' }}
            </NTag>
            <span>共 {{ total }} 条</span>
          </div>
        </div>

        <div class="config-category-strip">
          <button
            v-for="item in configCategories"
            :key="item.key"
            type="button"
            class="config-category"
            :class="{ 'config-category--active': item.group_code === activeConfigGroup }"
            @click="selectConfigCategory(item.group_code)"
          >
            <strong>{{ item.label }}</strong>
            <span>{{ item.description }}</span>
          </button>
        </div>
      </NCard>

      <ConfigFilterBar v-model:query="query" @search="handleSearch" @reset="handleReset" />

      <ConfigTable
        :columns="columns"
        :items="configs"
        :loading="loading"
        :query="query"
        :total="total"
        @page-change="handlePageChange"
        @page-size-change="handlePageSizeChange"
        @refresh="load"
      />
    </section>

    <ConfigFormModal
      v-model:show="formVisible"
      v-model:form-ref="formRef"
      :form-mode="formMode"
      :model="formModel"
      :rules="rules"
      :saving="saving"
      @submit="handleSubmit(submitForm)"
    />
  </main>
</template>

<style scoped>
.config-group-card {
  border: 1px solid var(--ez-border);
  border-radius: 8px;
  box-shadow: var(--ez-shadow-card);
}

.config-group-card :deep(.n-card__content) {
  padding: 14px 16px;
}

.config-group-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  margin-bottom: 12px;
}

.config-group-title {
  display: flex;
  min-width: 0;
  align-items: baseline;
  gap: 10px;
}

.config-group-title span {
  color: var(--ez-text-secondary);
  font-size: 13px;
}

.config-group-title strong {
  overflow: hidden;
  color: var(--ez-text-main);
  font-size: 16px;
  font-weight: 600;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.config-group-meta {
  display: flex;
  flex-shrink: 0;
  align-items: center;
  gap: 8px;
  color: var(--ez-text-secondary);
  font-size: 12px;
}

.config-category-strip {
  display: grid;
  grid-template-columns: repeat(6, minmax(0, 1fr));
  gap: 8px;
}

.config-category {
  display: flex;
  min-width: 0;
  min-height: 72px;
  flex-direction: column;
  justify-content: center;
  width: 100%;
  border: 1px solid transparent;
  border-radius: 8px;
  background: transparent;
  padding: 10px 12px;
  text-align: left;
  transition:
    border-color 0.2s ease,
    background-color 0.2s ease;
}

.config-category strong {
  display: block;
  overflow: hidden;
  color: var(--ez-text-main);
  font-size: 14px;
  font-weight: 600;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.config-category span {
  display: -webkit-box;
  overflow: hidden;
  min-height: 34px;
  margin-top: 4px;
  color: var(--ez-text-secondary);
  font-size: 12px;
  line-height: 1.5;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2;
}

.config-category--active,
.config-category:hover {
  border-color: var(--ez-primary);
  background: var(--ez-primary-light);
}

@media (max-width: 1024px) {
  .config-group-head {
    align-items: flex-start;
    flex-direction: column;
    gap: 8px;
  }

  .config-category-strip {
    display: flex;
    overflow-x: auto;
    padding-bottom: 2px;
    scrollbar-gutter: stable;
  }

  .config-category {
    width: 176px;
    flex: 0 0 176px;
  }
}

@media (max-height: 760px) and (min-width: 1025px) {
  .config-group-card :deep(.n-card__content) {
    padding: 12px 14px;
  }

  .config-group-head {
    margin-bottom: 10px;
  }

  .config-category {
    min-height: 64px;
    padding: 8px 10px;
  }

  .config-category span {
    margin-top: 2px;
    line-height: 1.35;
  }
}
</style>
