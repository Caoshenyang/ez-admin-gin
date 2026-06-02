<script setup lang="ts">
import type { FormInst, FormRules, SelectOption, TreeSelectOption } from 'naive-ui'
import { NButton, NForm, NFormItem, NInput, NModal, NSelect, NTreeSelect } from 'naive-ui'

import FormModalHeader from '@/components/FormModalHeader.vue'
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
    class="ez-modal-width-xl"
    @update:show="(value) => $emit('update:show', value)"
  >
    <template #header>
      <FormModalHeader
        :title="formMode === 'create' ? '新增用户' : '编辑用户'"
        :subtitle="formMode === 'create' ? '先创建账号主体，再补充岗位和角色归属。' : '编辑用户时仅调整资料、状态和岗位归属，角色可在列表中单独分配。'"
        @close="$emit('update:show', false)"
      />
    </template>

    <NForm
      ref="formRef"
      class="ez-modal-form"
      :model="formModel"
      :rules="rules"
      label-placement="left"
      label-width="76"
    >
      <section class="ez-modal-section ez-modal-section--soft">
        <div class="ez-modal-section__head">
          <h3>基础信息</h3>
          <p>先把账号主体信息补完整，这是本次弹窗的主要内容。</p>
        </div>

        <div v-if="formMode === 'create'" class="ez-form-grid ez-form-grid--2">
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

        <div v-else class="ez-form-grid ez-form-grid--2">
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

      <section class="ez-modal-section">
        <div class="ez-modal-section__head">
          <h3>岗位归属</h3>
          <p>岗位是用户归属的一部分，会直接影响后续通讯录、审批和业务协作能力的扩展空间。</p>
        </div>

        <NFormItem label="岗位" path="post_ids">
          <NSelect v-model:value="formModel.post_ids" multiple filterable :options="postOptions" placeholder="请选择岗位" />
        </NFormItem>
        <p class="ez-modal-section__tip">一个用户可以同时挂多个岗位，这里维护的是岗位归属，不会直接替代角色权限。</p>
      </section>

      <section v-if="formMode === 'create'" class="ez-modal-section">
        <div class="ez-modal-section__head">
          <h3>角色配置</h3>
          <p>这是补充信息，先给一个默认角色即可，后续仍可在列表中单独调整。</p>
        </div>

        <NFormItem label="角色" path="role_ids">
          <NSelect v-model:value="formModel.role_ids" multiple filterable :options="roleOptions" placeholder="请选择角色" />
        </NFormItem>
        <p class="ez-modal-section__tip">一个用户可以绑定多个角色，系统会自动合并其权限范围。</p>
      </section>
    </NForm>

    <template #footer>
      <div class="ez-modal-footer">
        <NButton @click="$emit('update:show', false)">取消</NButton>
        <NButton type="primary" :loading="saving" @click="$emit('submit')">保存</NButton>
      </div>
    </template>
  </NModal>
</template>
