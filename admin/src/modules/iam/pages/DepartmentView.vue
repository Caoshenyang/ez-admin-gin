<script setup lang="ts">
import { NButton } from 'naive-ui'

import DepartmentFilterBar from '../components/DepartmentFilterBar.vue'
import DepartmentFormModal from '../components/DepartmentFormModal.vue'
import DepartmentTable from '../components/DepartmentTable.vue'
import { useDepartmentPage } from '../composables/useDepartmentPage'

const {
  canUse,
  departments,
  formMode,
  formModel,
  formRef,
  formStatusOptions,
  formVisible,
  handleReset,
  handleSearch,
  handleSubmit,
  handleToggleStatus,
  loading,
  openCreate,
  openEdit,
  parentOptions,
  query,
  rules,
  saving,
  statusOptions,
} = useDepartmentPage()
</script>

<template>
  <main class="admin-page">
    <section class="admin-page-section">
      <div class="flex items-center justify-between">
                  <div class="ez-page-header">
            <h1>部门管理</h1>
          <p>维护组织树结构，为用户归属与数据权限提供稳定边界。</p>
        </div>

        <NButton v-if="canUse('system:department:create')" type="primary" @click="openCreate">
          + 新增部门
        </NButton>
      </div>

      <DepartmentFilterBar
        :keyword="query.keyword"
        :status="query.status"
        :status-options="statusOptions"
        @update:keyword="query.keyword = $event"
        @update:status="query.status = $event"
        @reset="handleReset"
        @search="handleSearch"
      />

      <DepartmentTable
        :can-use="canUse"
        :departments="departments"
        :loading="loading"
        @edit="openEdit"
        @toggle-status="handleToggleStatus"
      />
    </section>

    <DepartmentFormModal
      v-model:show="formVisible"
      v-model:form-ref="formRef"
      v-model:model="formModel"
      :form-mode="formMode"
      :form-status-options="formStatusOptions"
      :parent-options="parentOptions"
      :rules="rules"
      :saving="saving"
      @submit="handleSubmit"
    />
  </main>
</template>
