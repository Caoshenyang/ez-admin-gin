<script setup lang="ts">
import { NButton } from 'naive-ui'

import PageHeader from '@/components/PageHeader.vue'
import RoleFormModal from '../components/RoleFormModal.vue'
import RoleListPanel from '../components/RoleListPanel.vue'
import RolePermissionPanel from '../components/RolePermissionPanel.vue'
import { useRolePage } from '../composables/useRolePage'

const {
  activeTab,
  addPermissionRow,
  canEditBaseRole,
  canDeleteSelectedRole,
  canEditPermissionTab,
  canEditSelectedRole,
  canSavePermissionTab,
  canToggleSelectedRoleStatus,
  canUse,
  checkedFeatureCount,
  checkedMenuIDs,
  dataScopeHelp,
  dataScopeLabel,
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
  permissionSaveLabel,
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
      <PageHeader title="角色管理" description="维护角色基础信息、功能权限、接口策略与数据范围。">
        <template #actions>
          <NButton secondary :loading="loading" @click="handleSearch">刷新</NButton>
          <NButton v-if="canUse('system:role:create')" type="primary" @click="openCreate">
            + 新增角色
          </NButton>
        </template>
      </PageHeader>

      <div class="role-page-layout">
        <RoleListPanel
          v-model:query="query"
          :loading="loading"
          :roles="filteredRoles"
          :selected-role-id="selectedRoleID"
          :status-options="statusOptions"
          :status-type="statusType"
          :super-admin-role-code="superAdminRoleCode"
          @reset="handleReset"
          @search="handleSearch"
          @select="selectRole"
        />

        <RolePermissionPanel
          v-model:active-tab="activeTab"
          v-model:checked-menu-ids="checkedMenuIDs"
          v-model:permission-rows="permissionRows"
          :can-edit-base-role="canEditBaseRole"
          :can-delete-selected-role="canDeleteSelectedRole"
          :can-edit-permission-tab="canEditPermissionTab"
          :can-edit-selected-role="canEditSelectedRole"
          :checked-feature-count="checkedFeatureCount"
          :can-save-permission-tab="canSavePermissionTab"
          :can-toggle-selected-role-status="canToggleSelectedRoleStatus"
          :data-scope-help="dataScopeHelp"
          :data-scope-label="dataScopeLabel"
          :department-name-map="departmentNameMap"
          :menu-tree-options="menuTreeOptions"
          :method-options="methodOptions"
          :permission-save-label="permissionSaveLabel"
          :related-users="relatedUsers"
          :related-users-loading="relatedUsersLoading"
          :related-users-total="relatedUsersTotal"
          :saving="saving"
          :selected-role="selectedRole"
          :super-admin-role-code="superAdminRoleCode"
          @add-permission="addPermissionRow"
          @check-all="handleCheckAll"
          @clear-all="handleClearAll"
          @delete-role="handleDeleteRole"
          @edit-role="openEdit"
          @refresh-related-users="loadRelatedUsers"
          @remove-permission="removePermissionRow"
          @save-permissions="handleSavePermissions"
          @toggle-status="handleToggleRoleStatus"
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
  gap: 12px;
}

@media (min-width: 1024px) {
  .role-page-layout {
    grid-template-columns: minmax(292px, 320px) minmax(0, 1fr);
    grid-template-rows: minmax(0, 1fr);
    overflow: hidden;
  }
}

@media (min-width: 1200px) {
  .role-page-layout {
    grid-template-columns: minmax(320px, 360px) minmax(0, 1fr);
  }
}

@media (max-width: 1023px) {
  .role-page-layout {
    grid-auto-rows: minmax(420px, auto);
    overflow: auto;
  }
}
</style>
