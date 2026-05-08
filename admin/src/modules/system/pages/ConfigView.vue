<script setup lang="ts">
import { NAlert, NButton } from 'naive-ui'

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
      <div class="flex items-center justify-between">
        <div>
          <h1 class="text-[26px] font-bold text-[#111827]">配置管理</h1>
          <p class="mt-1 text-sm text-[#6B7280]">维护系统键值配置，按分组归类管理。</p>
        </div>

        <NButton v-if="canUse('system:config:create')" type="primary" @click="openCreate">
          + 新增配置
        </NButton>
      </div>

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
