<!-- RolePermissionPanel 展示并编辑指定角色的 Casbin 权限策略列表。 -->
<script setup lang="ts">
import type { SelectOption, TreeOption } from 'naive-ui'
import { NButton, NCard, NCheckbox, NInput, NSelect, NTabPane, NTabs, NTag, NTree } from 'naive-ui'

import type { PermissionRow, PermissionTab } from '../types/role-page'
import type { RoleItem } from '../types/role'

defineProps<{
  canEditSelectedRole: boolean
  checkedButtonCount: number
  checkedMenuCount: number
  checkedTotal: number
  menuTreeOptions: TreeOption[]
  methodOptions: SelectOption[]
  selectedRole: RoleItem | null
  superAdminRoleCode: string
}>()

defineEmits<{
  addPermission: []
  checkAll: []
  clearAll: []
  removePermission: [id: number]
}>()

const activeTab = defineModel<PermissionTab>('activeTab', { required: true })
const checkedMenuIDs = defineModel<Array<string | number>>('checkedMenuIds', { required: true })
const permissionRows = defineModel<PermissionRow[]>('permissionRows', { required: true })
</script>

<template>
  <NCard class="ez-card min-h-0 rounded-[var(--ez-radius-sm)]" :bordered="false" content-class="ez-card-content-fill">
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
          <NTabPane name="base" tab="基础信息">
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

          <NTabPane name="menu" tab="菜单权限">
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

          <NTabPane name="button" tab="按钮权限">
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

          <NTabPane name="api" tab="接口权限">
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

          <NTabPane name="data" tab="数据权限">
            <div class="data-scope-panel">
              <button type="button" class="data-scope-card data-scope-card--active">
                <strong>全部数据</strong>
                <span>可查看所有组织与业务数据</span>
              </button>
              <button type="button" class="data-scope-card">
                <strong>本部门数据</strong>
                <span>仅查看当前归属部门数据</span>
              </button>
              <button type="button" class="data-scope-card">
                <strong>本部门及下级</strong>
                <span>适合部门负责人和区域管理员</span>
              </button>
              <button type="button" class="data-scope-card">
                <strong>仅本人数据</strong>
                <span>限制为当前登录用户创建或归属的数据</span>
              </button>
              <button type="button" class="data-scope-card">
                <strong>自定义部门</strong>
                <span>后续接入部门树后可精细授权</span>
              </button>
            </div>
          </NTabPane>

          <NTabPane name="users" tab="关联用户">
            <div class="related-users-empty">
              <strong>{{ selectedRole?.name ?? '当前角色' }}</strong>
              <span>关联用户列表后续可接入用户分页接口，这里先保留规范化入口。</span>
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

.role-basic-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 12px;
}

.role-basic-grid > div,
.related-users-empty,
.data-scope-card {
  border: 1px solid var(--ez-border);
  border-radius: 10px;
  background: var(--ez-page-bg);
  padding: 14px;
}

.role-basic-grid span,
.data-scope-card span,
.related-users-empty span {
  display: block;
  color: var(--ez-text-secondary);
  font-size: 12px;
  line-height: 1.6;
}

.role-basic-grid strong,
.data-scope-card strong,
.related-users-empty strong {
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

.data-scope-card--active,
.data-scope-card:hover {
  border-color: var(--ez-primary);
  background: var(--ez-primary-light);
}

.related-users-empty {
  min-height: 180px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  text-align: center;
}

@media (max-width: 900px) {
  .role-basic-grid,
  .data-scope-panel {
    grid-template-columns: 1fr;
  }
}
</style>
