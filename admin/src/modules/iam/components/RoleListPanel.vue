<script setup lang="ts">
import type { SelectOption } from 'naive-ui'
import { NCard, NEmpty, NInput, NSelect, NSpin, NTag } from 'naive-ui'

import EzActionButton from '@/components/ez/EzActionButton.vue'
import type { RoleItem, RoleListQuery, RoleStatus } from '../types/role'

defineProps<{
  loading: boolean
  roles: RoleItem[]
  selectedRoleId: number | null
  statusOptions: SelectOption[]
  statusType: (status: RoleStatus) => 'success' | 'error'
  superAdminRoleCode: string
}>()

defineEmits<{
  reset: []
  search: []
  select: [role: RoleItem]
}>()

const query = defineModel<RoleListQuery>('query', { required: true })
</script>

<template>
  <NCard
    class="ez-card role-rail-card h-full min-h-0 overflow-hidden rounded-[var(--ez-radius-sm)]"
    :bordered="false"
    content-class="ez-card-content-fill"
  >
    <div class="role-rail">
      <header class="role-rail-head">
        <div>
          <h2>角色列表</h2>
          <p>{{ roles.length }} 个角色</p>
        </div>
      </header>

      <div class="role-rail-search">
        <NInput
          v-model:value="query.keyword"
          clearable
          placeholder="搜索角色名称 / 编码"
          @keyup.enter="$emit('search')"
        />
        <div class="role-rail-filter">
          <NSelect v-model:value="query.status" :options="statusOptions" />
          <EzActionButton
            kind="search"
            label="查询"
            type="primary"
            secondary
            @click="$emit('search')"
          />
          <EzActionButton kind="reset" label="重置" secondary @click="$emit('reset')" />
        </div>
      </div>

      <NSpin :show="loading" class="min-h-0 flex-1">
        <div class="role-nav-list">
          <div
            v-for="role in roles"
            :key="role.id"
            class="role-nav-item"
            :class="{ 'role-nav-item--active': role.id === selectedRoleId }"
            role="button"
            tabindex="0"
            @click="$emit('select', role)"
            @keyup.enter="$emit('select', role)"
          >
            <span class="role-nav-content">
              <span class="role-nav-main">
                <strong>{{ role.name }}</strong>
                <span class="role-nav-tags">
                  <NTag :type="statusType(role.status)" :bordered="false" size="small">
                    {{ role.status === 1 ? '启用' : '禁用' }}
                  </NTag>
                  <NTag
                    v-if="role.code === superAdminRoleCode"
                    type="warning"
                    :bordered="false"
                    size="small"
                  >
                    保护
                  </NTag>
                </span>
              </span>
              <span class="role-nav-code">{{ role.code }}</span>
            </span>
          </div>

          <div v-if="roles.length === 0" class="role-rail-empty">
            <NEmpty size="small" description="暂无匹配角色" />
          </div>
        </div>
      </NSpin>
    </div>
  </NCard>
</template>

<style scoped>
.role-rail-card {
  background: #fff;
}

.role-rail {
  display: flex;
  height: 100%;
  min-height: 0;
  flex-direction: column;
  overflow: hidden;
}

.role-rail-head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
  border-bottom: 1px solid rgba(226, 232, 240, 0.86);
  padding: 14px 16px 12px;
}

.role-rail-head h2 {
  margin: 0;
  color: var(--ez-text-main);
  font-size: 16px;
  font-weight: 700;
  letter-spacing: 0;
}

.role-rail-head p {
  margin: 3px 0 0;
  color: var(--ez-text-secondary);
  font-size: 12px;
  line-height: 18px;
}

.role-rail-search {
  display: grid;
  gap: 10px;
  padding: 14px 14px 12px;
}

.role-rail-filter {
  display: grid;
  grid-template-columns: minmax(0, 1fr) auto auto;
  gap: 8px;
}

.role-nav-list {
  min-height: 0;
  height: 100%;
  overflow-y: auto;
  padding: 4px 12px 12px;
}

.role-nav-item {
  position: relative;
  display: block;
  border: 1px solid rgba(226, 232, 240, 0.92);
  border-radius: 8px;
  background: #fbfdff;
  padding: 11px 12px;
  cursor: pointer;
  outline: none;
  transition:
    background-color 0.18s ease,
    border-color 0.18s ease,
    box-shadow 0.18s ease,
    transform 0.18s ease;
}

.role-nav-item + .role-nav-item {
  margin-top: 8px;
}

.role-nav-item:hover {
  border-color: rgba(148, 163, 184, 0.68);
  background: #fff;
  box-shadow: 0 4px 12px rgba(15, 23, 42, 0.05);
  transform: translateY(-1px);
}

.role-nav-item--active {
  border-color: rgba(37, 99, 255, 0.38);
  background: rgba(37, 99, 255, 0.07);
  box-shadow:
    inset 3px 0 0 var(--ez-primary),
    0 4px 14px rgba(37, 99, 255, 0.08);
  transform: none;
}

.role-nav-content {
  min-width: 0;
}

.role-nav-main {
  display: flex;
  min-width: 0;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}

.role-nav-tags {
  display: inline-flex;
  flex-shrink: 0;
  align-items: center;
  gap: 4px;
}

.role-nav-main strong {
  overflow: hidden;
  color: var(--ez-text-main);
  font-size: 14px;
  font-weight: 700;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.role-nav-code {
  display: block;
  overflow: hidden;
  margin-top: 4px;
  color: var(--ez-text-light);
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, 'Liberation Mono', monospace;
  font-size: 12px;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.role-rail-empty {
  padding: 38px 12px;
}
</style>
