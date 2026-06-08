<script setup lang="ts">
import { NCard, NInput, NTree } from 'naive-ui'

import EzActionButton from '@/components/ez/EzActionButton.vue'
import PageHeader from '@/components/PageHeader.vue'
import UserFilterBar from '../components/UserFilterBar.vue'
import UserFormModal from '../components/UserFormModal.vue'
import UserRoleModal from '../components/UserRoleModal.vue'
import UserTable from '../components/UserTable.vue'
import { useUserPage } from '../composables/useUserPage'
import type { UserStatus } from '../types/user'

const {
  canUse,
  checkedRowKeys,
  departmentCount,
  departmentKeyword,
  departmentNameMap,
  departmentTreeOptions,
  displayTotal,
  displayUsers,
  filteredDepartmentTreeOptions,
  formMode,
  formModel,
  formRef,
  formVisible,
  handleCheckedRowKeys,
  handleClearDepartment,
  handleDelete,
  handlePageChange,
  handlePageSizeChange,
  handleReset,
  handleSaveRoles,
  handleSearch,
  handleSelectDepartment,
  handleToggleStatus,
  loading,
  openCreate,
  openEdit,
  openRole,
  postNameMap,
  postOptions,
  query,
  roleFilterOptions,
  roleNameMap,
  roleOptions,
  roleSaving,
  roleUser,
  roleVisible,
  rules,
  saving,
  selectedCount,
  selectedDepartmentKeys,
  selectedRoleIDs,
  statusOptions,
  submitForm,
} = useUserPage()
</script>

<template>
  <main class="admin-page">
    <section class="admin-page-section">
      <PageHeader title="用户管理" description="维护后台账号、启停状态和角色绑定。">
        <template #actions>
          <EzActionButton
            v-if="canUse('system:user:create')"
            kind="add"
            label="新增用户"
            type="primary"
            @click="openCreate"
          />
        </template>
      </PageHeader>

      <div class="user-page-layout">
        <NCard class="user-dept-card" :bordered="false" content-class="ez-card-content-fill">
          <div class="user-dept-card__head">
            <strong>部门树</strong>
            <EzActionButton
              kind="reset"
              label="全部部门"
              size="tiny"
              text
              type="primary"
              @click="handleClearDepartment"
            />
          </div>
          <NInput v-model:value="departmentKeyword" size="small" placeholder="搜索部门" clearable />
          <NTree
            class="mt-3"
            block-line
            default-expand-all
            selectable
            :data="filteredDepartmentTreeOptions"
            key-field="value"
            label-field="label"
            :selected-keys="selectedDepartmentKeys"
            @update:selected-keys="handleSelectDepartment"
          />
          <div class="user-dept-card__foot">{{ departmentCount }} 个部门</div>
        </NCard>

        <section class="user-table-panel">
          <UserFilterBar
            :keyword="query.keyword ?? ''"
            :role-id="query.role_id ?? 0"
            :role-options="roleFilterOptions"
            :status="query.status ?? 0"
            :status-options="statusOptions"
            @update:keyword="query.keyword = $event"
            @update:role-id="query.role_id = $event"
            @update:status="query.status = $event as 0 | UserStatus"
            @search="handleSearch"
            @reset="handleReset"
          />

          <UserTable
            :checked-row-keys="checkedRowKeys"
            :department-name-map="departmentNameMap"
            :display-total="displayTotal"
            :loading="loading"
            :page="query.page"
            :page-size="query.page_size"
            :post-name-map="postNameMap"
            :role-name-map="roleNameMap"
            :selected-count="selectedCount"
            :users="displayUsers"
            :can-use="canUse"
            @checked-row-keys-change="handleCheckedRowKeys"
            @delete="handleDelete"
            @edit="openEdit"
            @page-change="handlePageChange"
            @page-size-change="handlePageSizeChange"
            @refresh="handleSearch"
            @role="openRole"
            @toggle-status="handleToggleStatus"
          />
        </section>
      </div>
    </section>

    <UserFormModal
      v-model:show="formVisible"
      v-model:form-ref="formRef"
      :department-tree-options="departmentTreeOptions"
      :form-mode="formMode"
      :model="formModel"
      :post-options="postOptions"
      :role-options="roleOptions"
      :rules="rules"
      :saving="saving"
      @submit="submitForm"
    />

    <UserRoleModal
      v-model:show="roleVisible"
      :role-options="roleOptions"
      :role-saving="roleSaving"
      :role-user="roleUser"
      :role-ids="selectedRoleIDs"
      @update:role-ids="selectedRoleIDs = $event"
      @submit="handleSaveRoles"
    />
  </main>
</template>

<style scoped>
.user-page-layout {
  display: grid;
  min-height: 0;
  flex: 1;
  grid-template-columns: 260px minmax(0, 1fr);
  grid-template-rows: minmax(0, 1fr);
  gap: 14px;
  overflow: hidden;
}

.user-dept-card {
  min-height: 0;
  overflow: hidden;
  border: 1px solid var(--ez-component-border);
  border-radius: var(--ez-radius-control);
  background: var(--ez-card-bg);
  box-shadow: var(--ez-component-shadow);
}

.user-dept-card__head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 9px;
  color: var(--ez-text-main);
  font-size: var(--ez-text-sm);
}

.user-dept-card__foot {
  margin-top: 10px;
  color: var(--ez-text-secondary);
  font-size: 12px;
}

.user-table-panel {
  display: flex;
  min-width: 0;
  min-height: 0;
  flex-direction: column;
  gap: 12px;
}

@media (max-width: 1024px) {
  .user-page-layout {
    grid-template-columns: 1fr;
    overflow: auto;
  }
}
</style>
