<script setup lang="ts">
import { NAlert, NButton } from 'naive-ui'

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
  closeSuccess,
  departmentNameMap,
  departmentTreeOptions,
  displayTotal,
  displayUsers,
  formMode,
  formModel,
  formRef,
  formVisible,
  handleCheckedRowKeys,
  handlePageChange,
  handlePageSizeChange,
  handleReset,
  handleSaveRoles,
  handleSearch,
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
  selectedRoleIDs,
  statusOptions,
  submitForm,
  successText,
} = useUserPage()
</script>

<template>
  <main class="admin-page">
    <section class="admin-page-section">
      <PageHeader title="用户管理" description="维护后台账号、启停状态和角色绑定。">
        <template #actions>
          <NButton v-if="canUse('system:user:create')" type="primary" @click="openCreate">
            + 新增用户
          </NButton>
        </template>
      </PageHeader>

      <NAlert
        v-if="successText"
        type="success"
        :show-icon="true"
        closable
        class="mx-auto w-full max-w-[520px]"
        @close="closeSuccess"
      >
        {{ successText }}
      </NAlert>

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
        @edit="openEdit"
        @page-change="handlePageChange"
        @page-size-change="handlePageSizeChange"
        @refresh="handleSearch"
        @role="openRole"
        @toggle-status="handleToggleStatus"
      />
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
