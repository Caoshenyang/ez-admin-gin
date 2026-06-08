<script setup lang="ts">
import EzActionButton from '@/components/ez/EzActionButton.vue'
import PageHeader from '@/components/PageHeader.vue'
import DepartmentFilterBar from '../components/DepartmentFilterBar.vue'
import DepartmentFormModal from '../components/DepartmentFormModal.vue'
import DepartmentTable from '../components/DepartmentTable.vue'
import { useDepartmentPage } from '../composables/useDepartmentPage'

const {
  canUse,
  checkedRowKeys,
  collapseAll,
  departments,
  expandAll,
  expandedRowKeys,
  formMode,
  formModel,
  formRef,
  formStatusOptions,
  formVisible,
  handleCheckedRowKeys,
  handleDeleteSelected,
  handleExpandedRowKeys,
  handleReset,
  handleSearch,
  handleSubmit,
  handleToggleStatus,
  leaderNameMap,
  loading,
  openCreate,
  openCreateChild,
  openEdit,
  parentOptions,
  query,
  rules,
  saving,
  selectedCount,
  statusOptions,
} = useDepartmentPage()
</script>

<template>
  <main class="admin-page">
    <section class="admin-page-section">
      <PageHeader title="部门管理" description="维护组织树结构，为用户归属与数据权限提供稳定边界。">
        <template #actions>
          <EzActionButton
            v-if="canUse('system:department:create')"
            kind="add"
            label="新增部门"
            type="primary"
            @click="openCreate"
          />
        </template>
      </PageHeader>

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
        :checked-row-keys="checkedRowKeys"
        :departments="departments"
        :expanded-row-keys="expandedRowKeys"
        :leader-name-map="leaderNameMap"
        :loading="loading"
        :selected-count="selectedCount"
        @checked-row-keys-change="handleCheckedRowKeys"
        @collapse-all="collapseAll"
        @create-child="openCreateChild"
        @delete-selected="handleDeleteSelected"
        @edit="openEdit"
        @expanded-row-keys-change="handleExpandedRowKeys"
        @expand-all="expandAll"
        @refresh="handleSearch"
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
