<!-- RolePermissionPanel 承接角色权限主流程：功能权限与接口权限。 -->
<script setup lang="ts">
import type { SelectOption, TreeOption } from 'naive-ui'
import {
  NButton,
  NCard,
  NDrawer,
  NDrawerContent,
  NEmpty,
  NInput,
  NPopconfirm,
  NSelect,
  NSpin,
  NTag,
  NTree,
} from 'naive-ui'
import { ref } from 'vue'

import { formatTime } from '@/utils/format'
import type { RoleItem } from '../types/role'
import type { PermissionRow, PermissionTab } from '../types/role-page'
import type { UserItem } from '../types/user'

defineProps<{
  canDeleteSelectedRole: boolean
  canEditBaseRole: boolean
  canEditPermissionTab: boolean
  canEditSelectedRole: boolean
  canSavePermissionTab: boolean
  canToggleSelectedRoleStatus: boolean
  checkedFeatureCount: number
  dataScopeHelp: string
  dataScopeLabel: string
  departmentNameMap: Map<number, string>
  menuTreeOptions: TreeOption[]
  methodOptions: SelectOption[]
  permissionSaveLabel: string
  relatedUsers: UserItem[]
  relatedUsersLoading: boolean
  relatedUsersTotal: number
  saving: boolean
  selectedRole: RoleItem | null
  superAdminRoleCode: string
}>()

const emit = defineEmits<{
  addPermission: []
  checkAll: []
  clearAll: []
  deleteRole: [role: RoleItem]
  editRole: [role: RoleItem]
  refreshRelatedUsers: []
  removePermission: [id: number]
  savePermissions: []
  toggleStatus: [role: RoleItem]
}>()

const activeTab = defineModel<PermissionTab>('activeTab', { required: true })
const checkedMenuIDs = defineModel<Array<string | number>>('checkedMenuIds', { required: true })
const permissionRows = defineModel<PermissionRow[]>('permissionRows', { required: true })

const relatedUsersVisible = ref(false)
const customApiVisible = ref(false)

const permissionTabs: Array<{ label: string; value: PermissionTab }> = [
  { label: '功能权限', value: 'feature' },
  { label: '接口权限', value: 'api' },
]

function openRelatedUsers() {
  relatedUsersVisible.value = true
  emit('refreshRelatedUsers')
}
</script>

<template>
  <NCard
    class="ez-card permission-console-card h-full min-h-0 overflow-hidden rounded-[var(--ez-radius-sm)]"
    :bordered="false"
    content-class="ez-card-content-fill"
  >
    <section class="permission-console">
      <header class="permission-command">
        <div class="permission-identity">
          <div class="permission-identity__main">
            <h2>{{ selectedRole?.name ?? '选择一个角色' }}</h2>
            <NTag v-if="selectedRole" :type="selectedRole.status === 1 ? 'success' : 'error'" :bordered="false">
              {{ selectedRole.status === 1 ? '启用' : '禁用' }}
            </NTag>
            <NTag v-if="selectedRole?.code === superAdminRoleCode" type="warning" :bordered="false">
              受保护
            </NTag>
          </div>
          <p>{{ selectedRole ? `编码：${selectedRole.code}` : '从左侧角色列表选择后开始配置权限' }}</p>
        </div>

        <div class="permission-command__actions">
          <NButton
            v-if="selectedRole"
            secondary
            :disabled="!canEditBaseRole"
            @click="emit('editRole', selectedRole)"
          >
            编辑角色
          </NButton>
          <NPopconfirm
            v-if="selectedRole"
            :disabled="!canToggleSelectedRoleStatus"
            @positive-click="emit('toggleStatus', selectedRole)"
          >
            <template #trigger>
              <NButton
                secondary
                :type="selectedRole.status === 1 ? 'warning' : 'success'"
                :disabled="!canToggleSelectedRoleStatus"
              >
                {{ selectedRole.status === 1 ? '禁用角色' : '启用角色' }}
              </NButton>
            </template>
            确认{{ selectedRole.status === 1 ? '禁用' : '启用' }}该角色？
          </NPopconfirm>
          <NPopconfirm
            v-if="selectedRole"
            :disabled="!canDeleteSelectedRole"
            @positive-click="emit('deleteRole', selectedRole)"
          >
            <template #trigger>
              <NButton secondary type="error" :disabled="!canDeleteSelectedRole">删除角色</NButton>
            </template>
            删除前请确认该角色没有分配给任何用户。
          </NPopconfirm>
          <NButton
            type="primary"
            :loading="saving"
            :disabled="!canSavePermissionTab"
            @click="emit('savePermissions')"
          >
            {{ permissionSaveLabel }}
          </NButton>
        </div>
      </header>

      <div class="permission-shell">
        <section v-if="!selectedRole" class="permission-empty-state">
          <NEmpty description="请先从左侧选择一个角色">
            <template #extra>
              <span>选择后可维护菜单、按钮、接口策略、数据范围和关联用户。</span>
            </template>
          </NEmpty>
        </section>

        <template v-else>
        <section class="permission-editor">
          <div class="role-detail-strip">
            <span>排序 {{ selectedRole.sort }}</span>
            <span>创建 {{ formatTime(selectedRole.created_at) }}</span>
            <span>更新 {{ formatTime(selectedRole.updated_at) }}</span>
          </div>

          <div class="permission-editor-tabs">
            <button
              v-for="tab in permissionTabs"
              :key="tab.value"
              type="button"
              :class="{ 'permission-editor-tab--active': activeTab === tab.value }"
              class="permission-editor-tab"
              @click="activeTab = tab.value"
            >
              {{ tab.label }}
            </button>
          </div>

          <div v-if="activeTab === 'feature'" class="permission-editor-panel">
            <div v-if="!canEditSelectedRole" class="permission-warning">
              <span>!</span>
              <strong>当前角色为内置保护角色，基础信息和权限配置保持只读。</strong>
            </div>
            <div v-else-if="!canEditPermissionTab" class="permission-warning">
              <span>!</span>
              <strong>当前账号没有分配功能权限的操作权限。</strong>
            </div>

            <div class="tree-workspace">
              <div class="tree-workspace__head">
                <div>
                  <strong>{{ checkedFeatureCount }} 项已选</strong>
                </div>
                <div class="tree-workspace__actions">
                  <NButton
                    size="small"
                    secondary
                    type="primary"
                    :disabled="!canEditPermissionTab"
                    @click="emit('checkAll')"
                  >
                    全选
                  </NButton>
                  <NButton size="small" secondary :disabled="!canEditPermissionTab" @click="emit('clearAll')">
                    清空
                  </NButton>
                </div>
              </div>
              <NEmpty
                v-if="menuTreeOptions.length === 0"
                size="small"
                description="暂无功能权限节点"
              />
              <NTree
                v-else
                v-model:checked-keys="checkedMenuIDs"
                checkable
                cascade
                block-line
                selectable
                :data="menuTreeOptions"
                :disabled="!canEditPermissionTab"
              />
            </div>
          </div>

          <div v-else class="permission-editor-panel">
            <div class="permission-warning">
              <span>!</span>
              <strong>当前后端暂未提供接口元数据列表，先维护已保存的 METHOD + PATH 策略。</strong>
            </div>
            <div v-if="!canEditSelectedRole" class="permission-warning">
              <span>!</span>
              <strong>当前角色为内置保护角色，接口策略保持只读。</strong>
            </div>
            <div v-else-if="!canEditPermissionTab" class="permission-warning">
              <span>!</span>
              <strong>当前账号没有分配接口权限的操作权限。</strong>
            </div>

            <div class="api-toolbar">
              <div>
                <strong>当前接口策略</strong>
                <span>{{ permissionRows.length }} 条策略</span>
              </div>
              <NButton
                size="small"
                secondary
                type="primary"
                :disabled="!canEditPermissionTab"
                @click="customApiVisible = !customApiVisible"
              >
                自定义策略
              </NButton>
            </div>

            <div v-if="permissionRows.length === 0" class="api-empty">
              <NEmpty size="small" description="暂无接口策略" />
            </div>

            <div v-else class="api-policy-table">
              <div class="api-policy-row api-policy-row--head">
                <span>方法</span>
                <span>路径</span>
                <span>操作</span>
              </div>
              <div v-for="row in permissionRows" :key="row.id" class="api-policy-row">
                <NSelect v-model:value="row.method" :options="methodOptions" :disabled="!canEditPermissionTab" />
                <NInput v-model:value="row.path" placeholder="/api/v1/system/users" :disabled="!canEditPermissionTab" />
                <NButton
                  size="small"
                  type="error"
                  secondary
                  :disabled="!canEditPermissionTab"
                  @click="emit('removePermission', row.id)"
                >
                  删除
                </NButton>
              </div>
            </div>

            <div v-if="customApiVisible" class="custom-api-entry">
              <NButton
                size="small"
                type="primary"
                ghost
                :disabled="!canEditPermissionTab"
                @click="emit('addPermission')"
              >
                添加一条自定义接口策略
              </NButton>
              <span>仅用于临时补录；建议后续从接口元数据中勾选。</span>
            </div>
          </div>
        </section>

        <aside class="permission-inspector">
          <div class="inspector-block">
            <h3>权限摘要</h3>
            <dl>
              <div>
                <dt>功能</dt>
                <dd>{{ checkedFeatureCount }}</dd>
              </div>
              <div>
                <dt>接口</dt>
                <dd>{{ permissionRows.length }}</dd>
              </div>
              <div>
                <dt>用户</dt>
                <dd>{{ relatedUsersTotal }}</dd>
              </div>
            </dl>
            <NButton size="small" secondary block :disabled="!selectedRole" @click="openRelatedUsers">
              查看关联用户
            </NButton>
          </div>

          <div class="inspector-block">
            <h3>数据范围</h3>
            <strong>{{ selectedRole ? dataScopeLabel : '-' }}</strong>
            <p v-if="selectedRole && dataScopeHelp">{{ dataScopeHelp }}</p>
          </div>

          <div v-if="selectedRole.custom_department_ids.length > 0" class="inspector-block">
            <h3>授权部门</h3>
            <div class="department-chip-list">
              <NTag
                v-for="departmentID in selectedRole.custom_department_ids"
                :key="departmentID"
                size="small"
                :bordered="false"
              >
                {{ departmentNameMap.get(departmentID) ?? `部门 ${departmentID}` }}
              </NTag>
            </div>
          </div>

          <div class="inspector-block">
            <h3>备注</h3>
            <p>{{ selectedRole.remark || '暂无备注' }}</p>
          </div>
        </aside>
        </template>
      </div>
    </section>

    <NDrawer v-model:show="relatedUsersVisible" width="420">
      <NDrawerContent title="关联用户">
        <div class="related-users-drawer-head">
          <div>
            <strong>{{ selectedRole?.name ?? '当前角色' }}</strong>
            <span>{{ relatedUsersTotal }} 人</span>
          </div>
          <NButton size="small" secondary @click="emit('refreshRelatedUsers')">刷新</NButton>
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
      </NDrawerContent>
    </NDrawer>
  </NCard>
</template>

<style scoped>
.permission-console-card {
  background: #fff;
}

.permission-console {
  display: flex;
  height: 100%;
  min-height: 0;
  flex-direction: column;
  overflow: hidden;
}

.permission-command {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 20px;
  border-bottom: 1px solid rgba(226, 232, 240, 0.9);
  background: #fff;
  padding: 15px 18px;
}

.permission-identity__main {
  display: flex;
  min-width: 0;
  flex-wrap: wrap;
  align-items: center;
  gap: 8px;
}

.permission-identity h2 {
  margin: 0;
  color: var(--ez-text-main);
  font-size: 18px;
  font-weight: 700;
  line-height: 1.18;
}

.permission-identity p {
  overflow: hidden;
  margin-top: 6px;
  color: var(--ez-text-sub);
  font-size: 12px;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.permission-command__actions {
  display: flex;
  flex-shrink: 0;
  align-items: center;
  gap: 10px;
}

.permission-shell {
  display: grid;
  min-height: 0;
  flex: 1;
  grid-template-columns: minmax(0, 1fr) 220px;
  background: #f8fafc;
}

.permission-empty-state {
  display: flex;
  grid-column: 1 / -1;
  min-height: 320px;
  align-items: center;
  justify-content: center;
  background: #fff;
}

.permission-empty-state span {
  color: var(--ez-text-secondary);
  font-size: 12px;
}

.permission-editor {
  display: flex;
  flex-direction: column;
  min-width: 0;
  min-height: 0;
  border-right: 1px solid rgba(226, 232, 240, 0.86);
  background: #fff;
  overflow: hidden;
}

.role-detail-strip {
  display: flex;
  flex-wrap: wrap;
  gap: 14px;
  border-bottom: 1px solid rgba(226, 232, 240, 0.82);
  background: #fff;
  padding: 10px 16px;
}

.role-detail-strip span,
.role-detail-strip button {
  border: 0;
  background: transparent;
  color: var(--ez-text-secondary);
  font-size: 12px;
  line-height: 20px;
}

.role-detail-strip button {
  color: var(--ez-primary);
}

.permission-editor-tabs {
  display: flex;
  gap: 4px;
  border-bottom: 1px solid rgba(226, 232, 240, 0.82);
  padding: 12px 16px 0;
}

.permission-editor-tab {
  border: 0;
  border-bottom: 2px solid transparent;
  background: transparent;
  padding: 0 4px 12px;
  color: var(--ez-text-sub);
  font-size: 14px;
  font-weight: 700;
}

.permission-editor-tab + .permission-editor-tab {
  margin-left: 18px;
}

.permission-editor-tab:hover {
  color: var(--ez-text-main);
}

.permission-editor-tab--active {
  border-bottom-color: var(--ez-primary);
  color: var(--ez-primary);
}

.permission-editor-panel {
  display: flex;
  flex: 1;
  flex-direction: column;
  gap: 8px;
  min-height: 0;
  overflow: auto;
  padding: 14px;
}

.permission-warning {
  display: flex;
  flex-shrink: 0;
  min-height: 24px;
  align-items: center;
  gap: 5px;
  border-radius: 6px;
  background: rgba(245, 158, 11, 0.08);
  padding: 3px 8px;
  color: #92400e;
}

.permission-warning span {
  display: inline-flex;
  width: 14px;
  height: 14px;
  flex-shrink: 0;
  align-items: center;
  justify-content: center;
  border-radius: 999px;
  background: #f59e0b;
  color: #fff;
  font-size: 10px;
  font-weight: 800;
  line-height: 1;
}

.permission-warning strong {
  min-width: 0;
  overflow: hidden;
  color: var(--ez-text-main);
  font-size: 11px;
  font-weight: 600;
  line-height: 14px;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.api-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 14px;
  border: 1px solid rgba(226, 232, 240, 0.9);
  border-radius: 8px;
  background: #fff;
  padding: 11px 12px;
}

.api-toolbar strong {
  display: block;
  color: var(--ez-text-main);
  font-size: 14px;
  font-weight: 700;
}

.api-toolbar span,
.custom-api-entry span {
  display: block;
  margin-top: 3px;
  color: var(--ez-text-sub);
  font-size: 12px;
}

.tree-workspace {
  flex: 1;
  min-height: 0;
  overflow: hidden;
  border: 1px solid rgba(226, 232, 240, 0.9);
  border-radius: 8px;
  background: #fff;
}

.tree-workspace__head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  border-bottom: 1px solid rgba(226, 232, 240, 0.74);
  background: #f8fafc;
  padding: 10px 13px;
}

.tree-workspace__head strong {
  display: block;
  color: var(--ez-primary);
  font-size: 12px;
  font-weight: 700;
}

.tree-workspace__actions {
  display: flex;
  flex-shrink: 0;
  align-items: center;
  gap: 8px;
}

:deep(.tree-workspace .n-tree) {
  padding: 12px;
}

:deep(.tree-workspace .n-tree-node-content) {
  min-height: 32px;
  border-radius: 7px;
}

:deep(.tree-workspace .n-tree-node-content:hover) {
  background: rgba(37, 99, 255, 0.055);
}

.api-empty,
.custom-api-entry,
.related-user-item {
  border: 1px solid rgba(226, 232, 240, 0.9);
  border-radius: 8px;
  background: #fff;
  padding: 14px;
}

.custom-api-entry {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 14px;
}

.api-policy-table {
  overflow: auto;
  border: 1px solid rgba(226, 232, 240, 0.9);
  border-radius: 8px;
  background: #fff;
}

.api-policy-row {
  display: grid;
  grid-template-columns: 120px minmax(240px, 1fr) 76px;
  align-items: center;
  gap: 12px;
  border-bottom: 1px solid rgba(226, 232, 240, 0.72);
  padding: 10px 12px;
}

.api-policy-row:last-child {
  border-bottom: 0;
}

.api-policy-row--head {
  background: #f8fafc;
  color: var(--ez-text-sub);
  font-size: 12px;
  font-weight: 900;
}

.permission-inspector {
  display: grid;
  align-content: start;
  gap: 10px;
  min-width: 0;
  padding: 14px;
}

.inspector-block {
  border: 1px solid rgba(226, 232, 240, 0.9);
  border-radius: 8px;
  background: #fff;
  padding: 12px;
  text-align: left;
}

.inspector-block h3 {
  margin: 0 0 10px;
  color: var(--ez-text-main);
  font-size: 12px;
  font-weight: 700;
}

.inspector-block strong {
  display: block;
  color: var(--ez-text-main);
  font-size: 13px;
  font-weight: 700;
  line-height: 1.5;
}

.inspector-block p {
  margin: 6px 0 0;
  color: var(--ez-text-sub);
  font-size: 12px;
  line-height: 1.55;
}

.inspector-block dl {
  display: grid;
  gap: 8px;
  margin: 0 0 12px;
}

.inspector-block dl > div {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.inspector-block dt {
  color: var(--ez-text-secondary);
  font-size: 12px;
}

.inspector-block dd {
  margin: 0;
  color: var(--ez-text-main);
  font-size: 13px;
  font-weight: 700;
}

.department-chip-list {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  margin-top: 10px;
}

.related-users-drawer-head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 14px;
}

.related-users-drawer-head strong,
.related-user-item strong {
  display: block;
  color: var(--ez-text-main);
  font-size: 14px;
  font-weight: 900;
}

.related-users-drawer-head span,
.related-user-item span {
  display: block;
  margin-top: 4px;
  color: var(--ez-text-secondary);
  font-size: 12px;
  line-height: 1.6;
}

.related-user-list {
  display: grid;
  gap: 10px;
}

@media (max-width: 1200px) {
  .permission-shell {
    grid-template-columns: minmax(0, 1fr);
  }

  .permission-editor {
    border-right: 0;
  }

  .permission-inspector {
    border-top: 1px solid rgba(226, 232, 240, 0.86);
  }
}

@media (max-width: 900px) {
  .permission-command,
  .api-toolbar,
  .custom-api-entry {
    align-items: stretch;
    flex-direction: column;
  }

  .permission-command__actions,
  .tree-workspace__actions {
    justify-content: flex-start;
  }

  .permission-inspector,
  .api-policy-row {
    grid-template-columns: 1fr;
  }
}
</style>
