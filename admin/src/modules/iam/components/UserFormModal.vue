<script setup lang="ts">
import { CloseOutline } from '@vicons/ionicons5'
import type { FormInst, FormRules, SelectOption, TreeSelectOption } from 'naive-ui'
import { NButton, NForm, NFormItem, NIcon, NInput, NModal, NSelect, NTreeSelect } from 'naive-ui'

import type { UserFormModel } from '../types/user-page'

defineProps<{
  departmentTreeOptions: TreeSelectOption[]
  formMode: 'create' | 'edit'
  postOptions: SelectOption[]
  roleOptions: SelectOption[]
  rules: FormRules
  saving: boolean
  show: boolean
}>()

defineEmits<{
  'update:show': [value: boolean]
  submit: []
}>()

const formRef = defineModel<FormInst | null>('formRef')
const formModel = defineModel<UserFormModel>('model', { required: true })
</script>

<template>
  <NModal
    :show="show"
    preset="card"
    :closable="false"
    class="compact-user-modal"
    style="width: 640px; max-width: calc(100vw - 32px)"
    @update:show="(value) => $emit('update:show', value)"
  >
    <template #header>
      <div class="modal-header modal-header--hero">
        <h2 class="modal-header__title">{{ formMode === 'create' ? '新增用户' : '编辑用户' }}</h2>
        <p class="modal-header__hero-title">
          {{
            formMode === 'create'
              ? '先完成账号主体信息，再补充默认角色范围'
              : '这里仅维护展示资料，不修改登录凭证'
          }}
        </p>
        <p class="modal-header__hero-desc">
          {{
            formMode === 'create'
              ? '用户名和密码会作为首次登录凭证，角色支持后续在列表中继续微调。'
              : '编辑模式下不修改登录名和密码，避免影响已有账号的登录和追踪。'
          }}
        </p>
        <button type="button" class="modal-close" @click="$emit('update:show', false)">
          <NIcon :size="18">
            <CloseOutline />
          </NIcon>
        </button>
      </div>
    </template>

    <div class="user-modal-shell">
      <NForm
        ref="formRef"
        class="compact-user-form"
        :model="formModel"
        :rules="rules"
        label-placement="left"
        label-width="76"
      >
        <section class="form-section form-section--primary">
          <div class="form-section__head">
            <h3>基础信息</h3>
            <p>先把账号主体信息补完整，这是本次弹窗的主要内容。</p>
          </div>

          <div v-if="formMode === 'create'" class="form-section-grid">
            <NFormItem label="用户名" path="username">
              <NInput v-model:value="formModel.username" placeholder="请输入用户名" />
            </NFormItem>

            <NFormItem label="登录密码" path="password">
              <NInput v-model:value="formModel.password" type="password" show-password-on="click" placeholder="至少 8 位" />
            </NFormItem>

            <NFormItem label="昵称" path="nickname">
              <NInput v-model:value="formModel.nickname" placeholder="请输入昵称" />
            </NFormItem>

            <NFormItem label="部门" path="department_id">
              <NTreeSelect v-model:value="formModel.department_id" :options="departmentTreeOptions" placeholder="请选择部门" default-expand-all />
            </NFormItem>

            <NFormItem label="状态" path="status">
              <NSelect v-model:value="formModel.status" :options="[{ label: '启用', value: 1 }, { label: '禁用', value: 2 }]" />
            </NFormItem>
          </div>

          <div v-else class="form-section-grid">
            <NFormItem label="昵称" path="nickname">
              <NInput v-model:value="formModel.nickname" placeholder="请输入昵称" />
            </NFormItem>

            <NFormItem label="部门" path="department_id">
              <NTreeSelect v-model:value="formModel.department_id" :options="departmentTreeOptions" placeholder="请选择部门" default-expand-all />
            </NFormItem>

            <NFormItem label="状态" path="status">
              <NSelect v-model:value="formModel.status" :options="[{ label: '启用', value: 1 }, { label: '禁用', value: 2 }]" />
            </NFormItem>
          </div>
        </section>

        <section class="form-section form-section--muted">
          <div class="form-section__head">
            <h3>岗位归属</h3>
            <p>岗位是用户归属的一部分，会直接影响后续通讯录、审批和业务协作能力的扩展空间。</p>
          </div>

          <NFormItem label="岗位" path="post_ids" class="mb-0">
            <NSelect v-model:value="formModel.post_ids" multiple filterable :options="postOptions" placeholder="请选择岗位" />
          </NFormItem>
          <p class="form-section__tip">一个用户可以同时挂多个岗位，这里维护的是岗位归属，不会直接替代角色权限。</p>
        </section>

        <section v-if="formMode === 'create'" class="form-section form-section--muted">
          <div class="form-section__head">
            <h3>角色配置</h3>
            <p>这是补充信息，先给一个默认角色即可，后续仍可在列表中单独调整。</p>
          </div>

          <NFormItem label="角色" path="role_ids" class="mb-0">
            <NSelect v-model:value="formModel.role_ids" multiple filterable :options="roleOptions" placeholder="请选择角色" />
          </NFormItem>
          <p class="form-section__tip">一个用户可以绑定多个角色，系统会自动合并其权限范围。</p>
        </section>
      </NForm>
    </div>

    <template #footer>
      <div class="modal-footer-actions">
        <NButton quaternary class="modal-footer-button" @click="$emit('update:show', false)">取消</NButton>
        <NButton type="primary" class="modal-footer-button modal-footer-button--primary" :loading="saving" @click="$emit('submit')">
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

.compact-user-form :deep(.n-form-item) {
  margin-bottom: 16px;
}

.compact-user-form :deep(.n-form-item-label) {
  white-space: nowrap;
  align-items: center;
  padding-right: 14px;
  font-weight: 600;
  color: #374151;
}

.compact-user-form :deep(.n-form-item-blank) {
  min-height: 40px;
}

.compact-user-form :deep(.n-input-wrapper) {
  border-radius: 10px;
  background: #fbfcfe;
}

.compact-user-form :deep(.n-base-selection) {
  border-radius: 10px;
  background: #fbfcfe;
}

.compact-user-form :deep(.n-input),
.compact-user-form :deep(.n-base-selection) {
  box-shadow: none;
}

.compact-user-form {
  display: flex;
  flex-direction: column;
  gap: 18px;
}

.user-modal-shell {
  display: flex;
  flex-direction: column;
  gap: 16px;
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

.modal-header--hero::after {
  content: '';
  position: absolute;
  top: -18px;
  right: -10px;
  width: 118px;
  height: 118px;
  border-radius: 999px;
  background: radial-gradient(circle, rgba(34, 197, 94, 0.1) 0%, rgba(34, 197, 94, 0) 72%);
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

.modal-header__hero-desc {
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
  box-shadow: 0 10px 24px rgba(148, 163, 184, 0.12);
  backdrop-filter: blur(8px);
  cursor: pointer;
}

.form-section {
  border: 1px solid #e9eff6;
  border-radius: 14px;
  background: #ffffff;
  padding: 18px 18px 4px;
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.9);
}

.form-section--primary {
  border-color: #d9e7f8;
  background: linear-gradient(180deg, #ffffff 0%, #fcfdff 100%);
}

.form-section--muted {
  background: linear-gradient(180deg, #fcfdff 0%, #f9fbff 100%);
}

.form-section__head {
  margin-bottom: 12px;
}

.form-section__head h3 {
  font-size: 15px;
  font-weight: 600;
  color: #111827;
}

.form-section__head p,
.form-section__tip {
  margin-top: 4px;
  font-size: 12px;
  line-height: 1.6;
  color: #6b7280;
}

.form-section-grid {
  display: grid;
  grid-template-columns: minmax(0, 1fr) minmax(0, 1fr);
  column-gap: 20px;
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

@media (max-width: 720px) {
  .form-section-grid {
    grid-template-columns: minmax(0, 1fr);
  }
}
</style>
