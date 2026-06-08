<!-- RolePermissionPanel 展示并编辑指定角色的 Casbin 权限策略列表。 -->
<script setup lang="ts">
import type { SelectOption, TreeOption } from 'naive-ui'
import { NButton, NCard, NCheckbox, NEmpty, NInput, NSelect, NSpin, NTabPane, NTabs, NTag, NTree } from 'naive-ui'

import { RoleDataScope } from '../types/role'
import type { PermissionRow, PermissionTab } from '../types/role-page'
import type { RoleItem } from '../types/role'
import type { UserItem } from '../types/user'
import { roleDataScopeDescriptions, roleDataScopeOptions } from '../composables/role-page.utils'

defineProps<{
  canEditSelectedRole: boolean
  checkedButtonCount: number
  checkedMenuCount: number
  checkedTotal: number
  dataScopeDescription: string
  departmentNameMap: Map<number, string>
  menuTreeOptions: TreeOption[]
  methodOptions: SelectOption[]
  relatedUsers: UserItem[]
  relatedUsersLoading: boolean
  relatedUsersTotal: number
  selectedRole: RoleItem | null
  superAdminRoleCode: string
}>()

defineEmits<{
  addPermission: []
  checkAll: []
  clearAll: []
  refreshRelatedUsers: []
  removePermission: [id: number]
}>()

const activeTab = defineModel<PermissionTab>('activeTab', { required: true })
const checkedMenuIDs = defineModel<Array<string | number>>('checkedMenuIds', { required: true })
const permissionRows = defineModel<PermissionRow[]>('permissionRows', { required: true })
</script>

<template>
  <NCard
    class="ez-card h-full min-h-0 overflow-hidden rounded-[var(--ez-radius-sm)]"
    :bordered="false"
    content-class="ez-card-content-fill"
  >
    <div class="flex h-full flex-col overflow-hidden">
      <div class="border-b border-[var(--ez-border)] px-5 py-5">
        <div class="flex items-start justify-between gap-4">
          <div>
            <h2 class="text-lg font-bold text-[var(--ez-text-main)]">菜单与按钮权限</h2>
            <p class="mt-2 text-sm text-[var(--ez-text-sub)]">
              当前角色：
              <span class="font-semibold text-[var(--ez-text-main)]">
                {{ selectedRole?.name ?? '未选择' }}
              </span>
              。半选状态表示部分子权限已授权。
            </p>
          </div>
          <NTag v-if="selectedRole?.code === superAdminRoleCode" type="warning" :bordered="false">受保护角色</NTag>
        </div>
      </div>

      <div class="min-h-0 flex-1 overflow-y-auto px-5 py-4">
        <NTabs v-model:value="activeTab" type="line" animated>
          <NTabPane name="base" tab="基础">
            <div class="role-basic-grid">
              <div>
                <span>角色名称</span>
                <strong>{{ selectedRole?.name ?? '未选择' }}</strong>
              </div>
              <div>
                <span>角色编码</span>
                <strong>{{ selectedRole?.code ?? '-' }}</strong>
              </div>
              <div>
                <span>角色状态</span>
                <NTag v-if="selectedRole" :type="selectedRole.status === 1 ? 'success' : 'error'" :bordered="false">
                  {{ selectedRole.status === 1 ? '启用' : '禁用' }}
                </NTag>
                <strong v-else>-</strong>
              </div>
              <div>
                <span>权限概览</span>
                <strong>{{ checkedTotal }} 项</strong>
              </div>
            </div>
          </NTabPane>

          <NTabPane name="menu" tab="菜单">
            <div class="permission-toolbar">
              <NCheckbox :checked="checkedTotal > 0" @update:checked="$emit('checkAll')">全选</NCheckbox>
              <NButton text type="primary" @click="$emit('checkAll')">全选可用节点</NButton>
              <NButton text type="primary" @click="$emit('clearAll')">清空全部</NButton>
            </div>

            <NTree
              v-model:checked-keys="checkedMenuIDs"
              checkable
              cascade
              block-line
              selectable
              :data="menuTreeOptions"
              :disabled="!canEditSelectedRole"
            />
          </NTabPane>

          <NTabPane name="button" tab="按钮">
            <div class="permission-toolbar">
              <NButton text type="primary" @click="$emit('checkAll')">全选可用节点</NButton>
              <NButton text type="primary" @click="$emit('clearAll')">清空全部</NButton>
            </div>

            <NTree
              v-model:checked-keys="checkedMenuIDs"
              checkable
              cascade
              block-line
              selectable
              :data="menuTreeOptions"
              :disabled="!canEditSelectedRole"
            />
          </NTabPane>

          <NTabPane name="api" tab="接口">
            <div class="mb-3 flex items-center justify-between">
              <p class="text-sm text-[var(--ez-text-sub)]">接口权限按请求路径和方法保存到 Casbin 策略表。</p>
              <NButton size="small" type="primary" ghost :disabled="!canEditSelectedRole" @click="$emit('addPermission')">
                + 添加接口
              </NButton>
            </div>

            <div class="space-y-3">
              <div v-for="row in permissionRows" :key="row.id" class="grid grid-cols-[130px_minmax(0,1fr)_80px] items-center gap-3 max-[720px]:grid-cols-1">
                <NSelect v-model:value="row.method" :options="methodOptions" :disabled="!canEditSelectedRole" />
                <NInput v-model:value="row.path" placeholder="/api/v1/system/users" :disabled="!canEditSelectedRole" />
                <NButton size="small" type="error" ghost :disabled="!canEditSelectedRole" @click="$emit('removePermission', row.id)">
                  删除
                </NButton>
              </div>
            </div>
          </NTabPane>

          <NTabPane name="data" tab="数据">
            <div class="data-scope-panel">
              <div
                v-for="option in roleDataScopeOptions"
                :key="String(option.value)"
                class="data-scope-card"
                :class="{ 'data-scope-card--active': selectedRole?.data_scope === option.value }"
              >
                <strong>{{ option.label }}</strong>
                <span>{{ roleDataScopeDescriptions.get(option.value as RoleDataScope) }}</span>
              </div>
            </div>

            <div class="data-scope-summary">
              <strong>{{ selectedRole ? dataScopeDescription : '未选择角色' }}</strong>
              <span v-if="selectedRole?.data_scope === RoleDataScope.CustomDept">
                已授权 {{ selectedRole.custom_department_ids.length }} 个自定义部门
              </span>
              <span v-else>通过“编辑角色”调整数据权限范围。</span>
              <div v-if="selectedRole?.data_scope === RoleDataScope.CustomDept" class="custom-department-tags">
                <NTag
                  v-for="departmentID in selectedRole.custom_department_ids"
                  :key="departmentID"
                  :bordered="false"
                  type="info"
                >
                  {{ departmentNameMap.get(departmentID) ?? `部门 ${departmentID}` }}
                </NTag>
                <NEmpty v-if="selectedRole.custom_department_ids.length === 0" size="small" description="未选择自定义部门" />
              </div>
            </div>
          </NTabPane>

          <NTabPane name="users" tab="用户">
            <div class="related-users-panel">
              <div class="related-users-panel__head">
                <strong>{{ selectedRole?.name ?? '当前角色' }}</strong>
                <div>
                  <span>{{ relatedUsersTotal }} 人</span>
                  <NButton size="tiny" text type="primary" @click="$emit('refreshRelatedUsers')">刷新</NButton>
                </div>
              </div>

              <NSpin :show="relatedUsersLoading">
                <NEmpty v-if="!selectedRole || relatedUsers.length === 0" description="暂无关联用户" />
                <div v-else class="related-user-list">
                  <div v-for="user in relatedUsers" :key="user.id" class="related-user-item">
                    <strong>{{ user.nickname || user.username }}</strong>
                    <span>
                      {{ user.username }} ·
                      {{ user.department_id === 0 ? '未分配部门' : departmentNameMap.get(user.department_id) ?? `部门 ${user.department_id}` }}
                    </span>
                  </div>
                </div>
              </NSpin>
            </div>
          </NTabPane>
        </NTabs>
      </div>

      <div class="permission-summary">
        <span>已授权菜单：{{ checkedMenuCount }}</span>
        <span>按钮权限：{{ checkedButtonCount }}</span>
        <span>接口权限：{{ permissionRows.length }}</span>
      </div>
    </div>
  </NCard>
</template>

<style scoped>
.permission-toolbar {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 16px;
  padding: 10px 12px;
  border-radius: var(--ez-radius-xs);
  background: var(--ez-page-bg);
}

.permission-summary {
  display: flex;
  gap: 32px;
  margin: 0 20px 20px;
  padding: 16px 18px;
  border-radius: var(--ez-radius-xs);
  background: var(--ez-brand-soft);
  color: var(--ez-brand);
  font-weight: 700;
}

:deep(.n-tabs-nav-scroll-wrapper) {
  overflow-x: auto;
}

:deep(.n-tabs-nav-scroll-content) {
  min-width: 0;
  flex-wrap: wrap;
  row-gap: 4px;
}

:deep(.n-tabs-tab) {
  padding-right: 10px;
  padding-left: 10px;
}

@media (max-width: 1200px) {
  :deep(.n-tabs-tab) {
    padding-right: 8px;
    padding-left: 8px;
    font-size: var(--ez-text-sm);
  }
}

.role-basic-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 12px;
}

.role-basic-grid > div,
.related-users-panel,
.related-user-item,
.data-scope-summary,
.data-scope-card {
  border: 1px solid var(--ez-border);
  border-radius: var(--ez-radius-xs);
  background: var(--ez-page-bg);
  padding: 14px;
}

.role-basic-grid span,
.data-scope-card span,
.data-scope-summary span,
.related-users-panel__head span,
.related-user-item span {
  display: block;
  color: var(--ez-text-secondary);
  font-size: 12px;
  line-height: 1.6;
}

.role-basic-grid strong,
.data-scope-card strong,
.data-scope-summary strong,
.related-users-panel__head strong,
.related-user-item strong {
  display: block;
  margin-top: 6px;
  color: var(--ez-text-main);
  font-size: 14px;
}

.data-scope-panel {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
}

.data-scope-card {
  text-align: left;
  transition:
    border-color 0.2s ease,
    background-color 0.2s ease;
}

.data-scope-card--active {
  border-color: var(--ez-primary);
  background: var(--ez-primary-light);
}

.data-scope-summary {
  margin-top: 12px;
}

.custom-department-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 12px;
}

.related-users-panel {
  min-height: 180px;
}

.related-users-panel__head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 12px;
}

.related-users-panel__head > div {
  display: flex;
  align-items: center;
  gap: 12px;
}

.related-user-list {
  display: grid;
  gap: 10px;
}

.related-user-item strong {
  margin-top: 0;
}

@media (max-width: 900px) {
  .role-basic-grid,
  .data-scope-panel {
    grid-template-columns: 1fr;
  }
}
</style>
