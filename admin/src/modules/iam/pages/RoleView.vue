<script setup lang="ts">
import { NButton, NSpace } from 'naive-ui'

import PageHeader from '@/components/PageHeader.vue'
import RoleFormModal from '../components/RoleFormModal.vue'
import RoleListPanel from '../components/RoleListPanel.vue'
import RolePermissionPanel from '../components/RolePermissionPanel.vue'
import { useRolePage } from '../composables/useRolePage'

const {
  activeTab,
  addPermissionRow,
  canEditSelectedRole,
  canSavePermissionTab,
  canUse,
  checkedButtonCount,
  checkedMenuCount,
  checkedMenuIDs,
  checkedTotal,
  dataScopeDescription,
  dataScopeOptions,
  departmentNameMap,
  departmentTreeOptions,
  filteredRoles,
  formMode,
  formModel,
  formRef,
  formVisible,
  handleCheckAll,
  handleClearAll,
  handleDeleteRole,
  handleReset,
  handleSavePermissions,
  handleSearch,
  handleToggleRoleStatus,
  loadRelatedUsers,
  loading,
  menuTreeOptions,
  methodOptions,
  openCreate,
  openEdit,
  permissionRows,
  query,
  relatedUsers,
  relatedUsersLoading,
  relatedUsersTotal,
  removePermissionRow,
  rules,
  saving,
  selectRole,
  selectedRole,
  selectedRoleID,
  statusOptions,
  statusType,
  submitRole,
  superAdminRoleCode,
} = useRolePage()
</script>

<template>
  <main class="admin-page">
    <section class="admin-page-section">
      <PageHeader title="角色权限" description="维护角色本身，以及角色拥有的菜单、按钮和接口权限。">
        <template #actions>
          <NSpace>
            <NButton v-if="canUse('system:role:create')" type="primary" @click="openCreate">
              + 新增角色
            </NButton>
            <NButton type="primary" :loading="saving" :disabled="!canSavePermissionTab" @click="handleSavePermissions">
              保存权限
            </NButton>
          </NSpace>
        </template>
      </PageHeader>

      <div class="role-page-layout">
        <RoleListPanel
          v-model:query="query"
          :can-use="canUse"
          :loading="loading"
          :roles="filteredRoles"
          :selected-role-id="selectedRoleID"
          :status-options="statusOptions"
          :status-type="statusType"
          :super-admin-role-code="superAdminRoleCode"
          @edit="openEdit"
          @delete="handleDeleteRole"
          @reset="handleReset"
          @search="handleSearch"
          @select="selectRole"
          @toggle-status="handleToggleRoleStatus"
        />

        <RolePermissionPanel
          v-model:active-tab="activeTab"
          v-model:checked-menu-ids="checkedMenuIDs"
          v-model:permission-rows="permissionRows"
          :can-edit-selected-role="canEditSelectedRole"
          :checked-button-count="checkedButtonCount"
          :checked-menu-count="checkedMenuCount"
          :checked-total="checkedTotal"
          :data-scope-description="dataScopeDescription"
          :department-name-map="departmentNameMap"
          :menu-tree-options="menuTreeOptions"
          :method-options="methodOptions"
          :related-users="relatedUsers"
          :related-users-loading="relatedUsersLoading"
          :related-users-total="relatedUsersTotal"
          :selected-role="selectedRole"
          :super-admin-role-code="superAdminRoleCode"
          @add-permission="addPermissionRow"
          @check-all="handleCheckAll"
          @clear-all="handleClearAll"
          @refresh-related-users="loadRelatedUsers"
          @remove-permission="removePermissionRow"
        />
      </div>
    </section>

    <RoleFormModal
      v-model:show="formVisible"
      v-model:form-ref="formRef"
      :data-scope-options="dataScopeOptions"
      :department-tree-options="departmentTreeOptions"
      :form-mode="formMode"
      :model="formModel"
      :rules="rules"
      :saving="saving"
      :status-options="statusOptions.slice(1)"
      @submit="submitRole"
    />
  </main>
</template>

<style scoped>
.role-page-layout {
  display: grid;
  min-height: 0;
  flex: 1;
  gap: 14px;
}

@media (min-width: 1024px) {
  .role-page-layout {
    grid-template-columns: minmax(220px, 240px) minmax(0, 1fr);
    grid-template-rows: minmax(0, 1fr);
    overflow: hidden;
  }
}

@media (min-width: 1200px) {
  .role-page-layout {
    grid-template-columns: minmax(260px, 280px) minmax(0, 1fr);
  }
}

@media (max-width: 1023px) {
  .role-page-layout {
    grid-auto-rows: minmax(420px, auto);
  }
}
</style>
