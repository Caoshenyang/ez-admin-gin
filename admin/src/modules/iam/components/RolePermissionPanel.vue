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
  <NCard class="min-h-0 rounded-lg" :bordered="false" content-style="display: flex; height: 100%; min-height: 0; flex-direction: column; padding: 0;">
    <div class="flex h-full flex-col overflow-hidden">
      <div class="border-b border-[#E5E7EB] px-5 py-5">
        <div class="flex items-start justify-between gap-4">
          <div>
            <h2 class="text-lg font-bold text-[#111827]">菜单与按钮权限</h2>
            <p class="mt-2 text-sm text-[#6B7280]">
              当前角色：
              <span class="font-semibold text-[#111827]">
                {{ selectedRole?.name ?? '未选择' }}
              </span>
              。半选状态表示部分子权限已授权。
            </p>
          </div>
          <NTag v-if="selectedRole?.code === superAdminRoleCode" type="warning" :bordered="false">受保护角色</NTag>
        </div>
      </div>

      <div class="min-h-0 flex-1 overflow-y-auto px-5 py-4">
        <NTabs v-model:value="activeTab" type="segment" animated>
          <NTabPane name="menu" tab="菜单权限">
            <div class="permission-toolbar">
              <NCheckbox :checked="checkedTotal > 0" @update:checked="$emit('checkAll')">全选</NCheckbox>
              <NButton text type="primary" @click="$emit('checkAll')">展开全部</NButton>
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
              <p class="text-sm text-[#6B7280]">接口权限按请求路径和方法保存到 Casbin 策略表。</p>
              <NButton size="small" type="primary" ghost :disabled="!canEditSelectedRole" @click="$emit('addPermission')">
                + 添加接口
              </NButton>
            </div>

            <div class="space-y-3">
              <div v-for="row in permissionRows" :key="row.id" class="grid grid-cols-[130px_minmax(0,1fr)_80px] items-center gap-3">
                <NSelect v-model:value="row.method" :options="methodOptions" :disabled="!canEditSelectedRole" />
                <NInput v-model:value="row.path" placeholder="/api/v1/system/users" :disabled="!canEditSelectedRole" />
                <NButton size="small" type="error" ghost :disabled="!canEditSelectedRole" @click="$emit('removePermission', row.id)">
                  删除
                </NButton>
              </div>
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
  border-radius: 6px;
  background: #f7fafc;
}

.permission-summary {
  display: flex;
  gap: 32px;
  margin: 0 20px 20px;
  padding: 16px 18px;
  border-radius: 6px;
  background: #e9fbf1;
  color: #18a058;
  font-weight: 700;
}
</style>
