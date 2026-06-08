<!-- UserRoleModal 弹窗分配用户角色，禁止修改当前登录用户自身的角色。 -->
<script setup lang="ts">
import type { SelectOption } from 'naive-ui'
import { NButton, NForm, NFormItem, NInput, NModal, NSelect } from 'naive-ui'

import FormModalHeader from '@/components/FormModalHeader.vue'
import type { UserItem } from '../types/user'

defineProps<{
  roleOptions: SelectOption[]
  roleSaving: boolean
  roleUser: UserItem | null
  roleIds: number[]
  show: boolean
}>()

defineEmits<{
  'update:roleIds': [value: number[]]
  'update:show': [value: boolean]
  submit: []
}>()
</script>

<template>
  <NModal
    :show="show"
    preset="card"
    :closable="false"
    class="ez-modal-width-sm"
    @update:show="(value) => $emit('update:show', value)"
  >
    <template #header>
      <FormModalHeader
        title="分配角色"
        subtitle="建议先给最小权限角色，再按职责逐步放开；保存后立即生效，多角色会按并集合并权限。"
        @close="$emit('update:show', false)"
      />
    </template>

    <div class="ez-modal-shell">
      <NForm class="ez-modal-form" label-placement="left" label-width="76">
        <section class="ez-modal-section">
          <div class="ez-modal-section__head">
            <h3>角色设置</h3>
            <p>为当前账号选择一个或多个角色，保存后立即生效。</p>
          </div>

          <NFormItem label="当前用户">
            <NInput :value="roleUser?.username ?? ''" disabled />
          </NFormItem>

          <NFormItem label="角色">
            <NSelect
              :value="roleIds"
              multiple
              filterable
              :options="roleOptions"
              placeholder="请选择角色"
              @update:value="(value) => $emit('update:roleIds', value as number[])"
            />
          </NFormItem>
          <p class="ez-modal-section__tip">多角色场景下，菜单与按钮权限会按照并集生效。</p>
        </section>
      </NForm>
    </div>

    <template #footer>
      <div class="ez-modal-footer">
        <NButton quaternary class="min-w-[92px]" @click="$emit('update:show', false)">取消</NButton>
        <NButton type="primary" class="min-w-[92px]" :loading="roleSaving" @click="$emit('submit')">
          保存
        </NButton>
      </div>
    </template>
  </NModal>
</template>
