<!-- UserRoleModal 弹窗分配用户角色，禁止修改当前登录用户自身的角色。 -->
<script setup lang="ts">
import { CloseOutline } from '@vicons/ionicons5'
import type { SelectOption } from 'naive-ui'
import { NButton, NForm, NFormItem, NIcon, NInput, NModal, NSelect } from 'naive-ui'

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
    class="compact-user-modal"
    style="width: 520px; max-width: calc(100vw - 32px)"
    @update:show="(value) => $emit('update:show', value)"
  >
    <template #header>
      <div class="modal-header modal-header--hero modal-header--role">
        <h2 class="modal-header__title">分配角色</h2>
        <p class="modal-header__hero-title">建议先给最小权限角色，再按职责逐步放开</p>
        <p class="modal-header__hero-desc">保存后立即生效，多角色场景下权限会按并集进行合并。</p>
        <button type="button" class="modal-close" @click="$emit('update:show', false)">
          <NIcon :size="18">
            <CloseOutline />
          </NIcon>
        </button>
      </div>
    </template>

    <div class="user-modal-shell">
      <NForm class="compact-user-form" label-placement="left" label-width="76">
        <section class="form-section form-section--muted">
          <div class="form-section__head">
            <h3>角色设置</h3>
            <p>为当前账号选择一个或多个角色，保存后立即生效。</p>
          </div>

          <NFormItem label="当前用户">
            <NInput :value="roleUser?.username ?? ''" disabled />
          </NFormItem>

          <NFormItem label="角色" class="mb-0">
            <NSelect
              :value="roleIds"
              multiple
              filterable
              :options="roleOptions"
              placeholder="请选择角色"
              @update:value="(value) => $emit('update:roleIds', value as number[])"
            />
          </NFormItem>
          <p class="form-section__tip">多角色场景下，菜单与按钮权限会按照并集生效。</p>
        </section>
      </NForm>
    </div>

    <template #footer>
      <div class="modal-footer-actions">
        <NButton quaternary class="modal-footer-button" @click="$emit('update:show', false)">取消</NButton>
        <NButton
          type="primary"
          class="modal-footer-button modal-footer-button--primary"
          :loading="roleSaving"
          @click="$emit('submit')"
        >
          保存
        </NButton>
      </div>
    </template>
  </NModal>
</template>

<style scoped>
.compact-user-modal :deep(.n-card) {
  overflow: hidden;
  border-radius: 32px;
  border: 1px solid #dfe9f5;
  background: #ffffff;
  box-shadow: 0 24px 72px rgba(15, 23, 42, 0.16);
}

.compact-user-modal :deep(.n-card-header) {
  padding: 0;
  border-bottom: 1px solid #dfe9f5;
  background: linear-gradient(135deg, #eff6ff 0%, #e8f2ff 58%, #f4f9ff 100%);
}

.compact-user-modal :deep(.n-card__content) {
  padding: 20px 28px 10px;
}

.compact-user-modal :deep(.n-card__footer) {
  padding: 16px 28px 24px;
  border-top: 1px solid #edf2f7;
  background: rgba(248, 250, 252, 0.85);
}

.modal-header {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.modal-header--hero {
  position: relative;
  overflow: hidden;
  min-height: 140px;
  padding: 26px 28px 22px;
  background:
    radial-gradient(circle at top right, rgba(34, 197, 94, 0.12), transparent 24%),
    linear-gradient(135deg, #eff6ff 0%, #e8f2ff 58%, #f4f9ff 100%);
}

.modal-header--role::after {
  content: '';
  position: absolute;
  top: -18px;
  right: -10px;
  width: 118px;
  height: 118px;
  border-radius: 999px;
  background: radial-gradient(circle, rgba(37, 99, 235, 0.08) 0%, rgba(37, 99, 235, 0) 72%);
}

.modal-header__title {
  position: relative;
  z-index: 1;
  font-size: 19px;
  font-weight: 600;
  line-height: 1.3;
  color: #111827;
}

.modal-header__hero-title {
  position: relative;
  z-index: 1;
  font-size: 14px;
  font-weight: 600;
  line-height: 1.6;
  color: #0f172a;
}

.modal-header__hero-desc,
.form-section__tip,
.form-section__head p {
  position: relative;
  z-index: 1;
  font-size: 12px;
  line-height: 1.6;
  color: #64748b;
}

.modal-close {
  position: absolute;
  top: 20px;
  right: 22px;
  z-index: 2;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 38px;
  height: 38px;
  border: none;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.76);
  color: #64748b;
  cursor: pointer;
}

.compact-user-form :deep(.n-form-item-label) {
  white-space: nowrap;
  align-items: center;
  padding-right: 14px;
  font-weight: 600;
  color: #374151;
}

.user-modal-shell {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.form-section {
  border: 1px solid #e9eff6;
  border-radius: 14px;
  background: linear-gradient(180deg, #fcfdff 0%, #f9fbff 100%);
  padding: 18px 18px 4px;
}

.form-section__head {
  margin-bottom: 12px;
}

.form-section__head h3 {
  font-size: 15px;
  font-weight: 600;
  color: #111827;
}

.modal-footer-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

.modal-footer-button {
  min-width: 92px;
  height: 40px;
  border-radius: 10px;
}

.modal-footer-button--primary {
  box-shadow: 0 10px 24px rgba(34, 197, 94, 0.18);
}

.mb-0 {
  margin-bottom: 0;
}
</style>
