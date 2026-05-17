<script setup lang="ts">
import type { SelectOption } from 'naive-ui'
import { NButton, NCard, NInput, NPopconfirm, NSelect, NSpace, NTag } from 'naive-ui'

import type { RoleListQuery, RoleItem, RoleStatus } from '../types/role'

defineProps<{
  canUse: (code: string) => boolean
  loading: boolean
  roles: RoleItem[]
  selectedRoleId: number | null
  statusOptions: SelectOption[]
  statusType: (status: RoleStatus) => 'success' | 'error'
  superAdminRoleCode: string
}>()

defineEmits<{
  edit: [role: RoleItem]
  reset: []
  search: []
  select: [role: RoleItem]
  toggleStatus: [role: RoleItem]
}>()

const query = defineModel<RoleListQuery>('query', { required: true })
</script>

<template>
  <NCard class="ez-card min-h-0 rounded-[var(--ez-radius-sm)]" :bordered="false" content-class="ez-card-content-fill">
    <div class="flex h-full flex-col overflow-hidden">
      <div class="mb-4">
        <h2 class="text-lg font-bold text-[var(--ez-text-main)]">角色列表</h2>
        <p class="mt-1 text-xs text-[var(--ez-text-sub)]">点击左侧角色后，在右侧维护权限。</p>
      </div>

      <NSpace vertical :size="10" class="mb-4">
        <NInput v-model:value="query.keyword" clearable placeholder="角色编码 / 名称" @keyup.enter="$emit('search')" />
        <div class="grid grid-cols-[1fr_auto] gap-2">
          <NSelect v-model:value="query.status" :options="statusOptions" />
          <NButton @click="$emit('reset')">重置</NButton>
        </div>
      </NSpace>

      <div class="min-h-0 flex-1 space-y-3 overflow-y-auto pr-1">
        <button
          v-for="role in roles"
          :key="role.id"
          type="button"
          class="role-card"
          :class="{ 'role-card--active': role.id === selectedRoleId }"
          @click="$emit('select', role)"
        >
          <span class="flex items-center justify-between gap-2">
            <span class="min-w-0 truncate text-base font-bold text-[var(--ez-text-main)]">
              {{ role.name }}
            </span>
            <NTag :type="statusType(role.status)" :bordered="false" size="small">
              {{ role.status === 1 ? '启用' : '禁用' }}
            </NTag>
          </span>
          <span class="mt-1 block text-left text-xs text-[var(--ez-text-sub)]">
            {{ role.code }} · 菜单 {{ (role.menu_ids ?? []).length }} · 接口 {{ (role.permissions ?? []).length }}
          </span>
          <span class="mt-2 flex items-center gap-2">
            <NButton v-if="canUse('system:role:update')" size="tiny" @click.stop="$emit('edit', role)">编辑</NButton>
            <NPopconfirm v-if="canUse('system:role:status') && role.code !== superAdminRoleCode" @positive-click="$emit('toggleStatus', role)">
              <template #trigger>
                <NButton size="tiny" :type="role.status === 1 ? 'error' : 'success'" ghost @click.stop>
                  {{ role.status === 1 ? '禁用' : '启用' }}
                </NButton>
              </template>
              确认{{ role.status === 1 ? '禁用' : '启用' }}该角色？
            </NPopconfirm>
          </span>
        </button>
      </div>
    </div>
  </NCard>
</template>

<style scoped>
.role-card {
  width: 100%;
  border: 1px solid var(--ez-border);
  border-radius: var(--ez-radius-sm);
  background: var(--ez-card-bg);
  padding: 14px 12px;
  text-align: left;
  transition:
    border-color 0.2s ease,
    background-color 0.2s ease,
    box-shadow 0.2s ease;
}

.role-card:hover {
  border-color: var(--ez-brand);
  box-shadow: var(--ez-shadow-md);
}

.role-card--active {
  border-color: var(--ez-brand);
  background: var(--ez-brand-soft);
}
</style>
