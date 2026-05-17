<script setup lang="ts">
import { NAlert, NButton } from 'naive-ui'

import PageHeader from '@/components/PageHeader.vue'
import ConfigFilterBar from '../components/ConfigFilterBar.vue'
import ConfigFormModal from '../components/ConfigFormModal.vue'
import ConfigTable from '../components/ConfigTable.vue'
import { useConfigPage } from '../composables/useConfigPage'

const {
  canUse,
  closeSuccess,
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
  submitForm,
  successText,
  total,
} = useConfigPage()
</script>

<template>
  <main class="admin-page">
    <section class="admin-page-section">
      <PageHeader title="配置管理" description="维护系统键值配置，按分组归类管理。">
        <template #actions>
          <NButton v-if="canUse('system:config:create')" type="primary" @click="openCreate">
            + 新增配置
          </NButton>
        </template>
      </PageHeader>

      <NAlert v-if="successText" type="success" :show-icon="true" closable class="mx-auto w-full max-w-[520px]" @close="closeSuccess">
        {{ successText }}
      </NAlert>

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
