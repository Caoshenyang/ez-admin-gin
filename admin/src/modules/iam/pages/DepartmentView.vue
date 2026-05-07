<script setup lang="ts">
import type {
  DataTableColumns,
  FormInst,
  FormRules,
  TreeSelectOption,
} from 'naive-ui'
import {
  NButton,
  NCard,
  NDataTable,
  NForm,
  NFormItem,
  NInput,
  NInputNumber,
  NModal,
  NPopconfirm,
  NSelect,
  NSpace,
  NTag,
  NTreeSelect,
  useMessage,
} from 'naive-ui'
import { computed, h, onMounted, reactive, ref } from 'vue'

import {
  createDepartment,
  getDepartments,
  updateDepartment,
  updateDepartmentStatus,
} from '../api/department'
import { buttonPermissionCodes } from '@/router/dynamic-menu'
import {
  DepartmentStatus,
  type CreateDepartmentPayload,
  type DepartmentItem,
  type DepartmentListQuery,
} from '../types/department'

interface DepartmentFormModel extends CreateDepartmentPayload {
  id: number
}

const message = useMessage()
const loading = ref(false)
const saving = ref(false)
const departments = ref<DepartmentItem[]>([])
const formRef = ref<FormInst | null>(null)
const formVisible = ref(false)
const formMode = ref<'create' | 'edit'>('create')

const query = reactive<DepartmentListQuery>({
  keyword: '',
  status: 0,
})

const formModel = reactive<DepartmentFormModel>({
  id: 0,
  parent_id: 0,
  name: '',
  code: '',
  leader_user_id: 0,
  sort: 0,
  status: DepartmentStatus.Enabled,
  remark: '',
})

const statusOptions = [
  { label: '状态：全部', value: 0 },
  { label: '启用', value: DepartmentStatus.Enabled },
  { label: '禁用', value: DepartmentStatus.Disabled },
]

const formStatusOptions = statusOptions.slice(1)

const parentOptions = computed<TreeSelectOption[]>(() => {
  return [
    {
      label: '作为根部门',
      value: 0,
    },
    ...buildDepartmentTreeOptions(departments.value),
  ]
})

const rules: FormRules = {
  name: [{ required: true, message: '请输入部门名称', trigger: ['blur', 'input'] }],
  code: [{ required: true, message: '请输入部门编码', trigger: ['blur', 'input'] }],
}

const columns: DataTableColumns<DepartmentItem> = [
  {
    title: '部门',
    key: 'name',
    minWidth: 260,
    render(row) {
      return h('div', { class: 'leading-6' }, [
        h('p', { class: 'font-semibold text-[#111827]' }, row.name),
        h('p', { class: 'text-xs text-[#6B7280]' }, row.code),
      ])
    },
  },
  {
    title: '负责人',
    key: 'leader_user_id',
    width: 120,
    render(row) {
      return row.leader_user_id === 0 ? '未设置' : `用户 ${row.leader_user_id}`
    },
  },
  {
    title: '排序',
    key: 'sort',
    width: 90,
  },
  {
    title: '状态',
    key: 'status',
    width: 100,
    render(row) {
      return h(
        NTag,
        {
          bordered: false,
          type: row.status === DepartmentStatus.Enabled ? 'success' : 'error',
        },
        { default: () => (row.status === DepartmentStatus.Enabled ? '启用' : '禁用') },
      )
    },
  },
  {
    title: '更新时间',
    key: 'updated_at',
    width: 180,
    render(row) {
      return formatTime(row.updated_at)
    },
  },
  {
    title: '操作',
    key: 'actions',
    width: 180,
    render(row) {
      const nextStatus =
        row.status === DepartmentStatus.Enabled
          ? DepartmentStatus.Disabled
          : DepartmentStatus.Enabled

      return h(
        NSpace,
        { size: 8 },
        {
          default: () => [
            canUse('system:department:update')
              ? h(
                  NButton,
                  {
                    size: 'small',
                    ghost: true,
                    type: 'info',
                    onClick: () => openEdit(row),
                  },
                  { default: () => '编辑' },
                )
              : null,
            canUse('system:department:status')
              ? h(
                  NPopconfirm,
                  { onPositiveClick: () => handleToggleStatus(row, nextStatus) },
                  {
                    trigger: () =>
                      h(
                        NButton,
                        {
                          size: 'small',
                          ghost: true,
                          type: nextStatus === DepartmentStatus.Disabled ? 'error' : 'success',
                        },
                        { default: () => (nextStatus === DepartmentStatus.Disabled ? '禁用' : '启用') },
                      ),
                    default: () =>
                      `确认${nextStatus === DepartmentStatus.Disabled ? '禁用' : '启用'}该部门？`,
                  },
                )
              : null,
          ],
        },
      )
    },
  },
]

function canUse(code: string) {
  return buttonPermissionCodes.value.includes(code)
}

function formatTime(value: string) {
  return value ? new Date(value).toLocaleString() : '-'
}

function resetForm() {
  Object.assign(formModel, {
    id: 0,
    parent_id: 0,
    name: '',
    code: '',
    leader_user_id: 0,
    sort: 0,
    status: DepartmentStatus.Enabled,
    remark: '',
  })
}

async function loadDepartments() {
  loading.value = true
  try {
    departments.value = await getDepartments({
      keyword: query.keyword?.trim() || undefined,
      status: query.status === 0 ? undefined : query.status,
    })
  } finally {
    loading.value = false
  }
}

function handleSearch() {
  void loadDepartments()
}

function handleReset() {
  query.keyword = ''
  query.status = 0
  void loadDepartments()
}

function openCreate() {
  formMode.value = 'create'
  resetForm()
  formVisible.value = true
}

function openEdit(row: DepartmentItem) {
  formMode.value = 'edit'
  Object.assign(formModel, {
    id: row.id,
    parent_id: row.parent_id,
    name: row.name,
    code: row.code,
    leader_user_id: row.leader_user_id,
    sort: row.sort,
    status: row.status,
    remark: row.remark,
  })
  formVisible.value = true
}

async function handleSubmit() {
  await formRef.value?.validate()
  saving.value = true

  try {
    if (formMode.value === 'create') {
      await createDepartment({
        parent_id: formModel.parent_id,
        name: formModel.name,
        code: formModel.code,
        leader_user_id: formModel.leader_user_id,
        sort: formModel.sort,
        status: formModel.status,
        remark: formModel.remark,
      })
      message.success('部门创建成功')
    } else {
      await updateDepartment(formModel.id, {
        parent_id: formModel.parent_id,
        name: formModel.name,
        code: formModel.code,
        leader_user_id: formModel.leader_user_id,
        sort: formModel.sort,
        status: formModel.status,
        remark: formModel.remark,
      })
      message.success('部门更新成功')
    }

    formVisible.value = false
    await loadDepartments()
  } catch {
    message.error(formMode.value === 'create' ? '部门创建失败' : '部门更新失败')
  } finally {
    saving.value = false
  }
}

async function handleToggleStatus(row: DepartmentItem, status: DepartmentStatus) {
  try {
    await updateDepartmentStatus(row.id, { status })
    message.success(status === DepartmentStatus.Enabled ? '部门已启用' : '部门已禁用')
    await loadDepartments()
  } catch {
    message.error('部门状态更新失败')
  }
}

function buildDepartmentTreeOptions(items: DepartmentItem[]): TreeSelectOption[] {
  return items.map((item) => ({
    label: `${item.name}（${item.code}）`,
    value: item.id,
    children: item.children?.length ? buildDepartmentTreeOptions(item.children) : undefined,
  }))
}

onMounted(() => {
  void loadDepartments()
})
</script>

<template>
  <main class="admin-page">
    <section class="admin-page-section">
      <div class="flex items-center justify-between">
        <div>
          <h1 class="text-[26px] font-bold text-[#111827]">部门管理</h1>
          <p class="mt-1 text-sm text-[#6B7280]">维护组织树结构，为用户归属与数据权限提供稳定边界。</p>
        </div>

        <NButton v-if="canUse('system:department:create')" type="primary" @click="openCreate">
          + 新增部门
        </NButton>
      </div>

      <NCard :bordered="false" class="rounded-lg">
        <NSpace align="center" :wrap="true">
          <NInput
            v-model:value="query.keyword"
            clearable
            placeholder="搜索部门名称 / 编码"
            class="w-64"
            @keyup.enter="handleSearch"
          />
          <NSelect v-model:value="query.status" :options="statusOptions" class="w-36" />
          <NButton type="primary" @click="handleSearch">查询</NButton>
          <NButton @click="handleReset">重置</NButton>
        </NSpace>
      </NCard>

      <NCard class="min-h-0 flex-1 rounded-lg" :bordered="false" content-style="height: 100%; padding: 0;">
        <NDataTable
          class="h-full"
          style="height: 100%"
          :columns="columns"
          :data="departments"
          :loading="loading"
          :pagination="false"
          :bordered="false"
          children-key="children"
          flex-height
        />
      </NCard>
    </section>

    <NModal v-model:show="formVisible" preset="card" class="w-[680px]" :title="formMode === 'create' ? '新增部门' : '编辑部门'">
      <NForm ref="formRef" :model="formModel" :rules="rules" label-placement="top" class="grid gap-4 md:grid-cols-2">
        <NFormItem label="上级部门" path="parent_id">
          <NTreeSelect
            v-model:value="formModel.parent_id"
            :options="parentOptions"
            default-expand-all
            placeholder="请选择上级部门"
          />
        </NFormItem>

        <NFormItem label="负责人用户 ID" path="leader_user_id">
          <NInputNumber v-model:value="formModel.leader_user_id" :min="0" class="w-full" />
        </NFormItem>

        <NFormItem label="部门名称" path="name">
          <NInput v-model:value="formModel.name" placeholder="请输入部门名称" />
        </NFormItem>

        <NFormItem label="部门编码" path="code">
          <NInput v-model:value="formModel.code" placeholder="请输入部门编码" />
        </NFormItem>

        <NFormItem label="排序" path="sort">
          <NInputNumber v-model:value="formModel.sort" :min="0" class="w-full" />
        </NFormItem>

        <NFormItem label="状态" path="status">
          <NSelect v-model:value="formModel.status" :options="formStatusOptions" />
        </NFormItem>

        <NFormItem label="备注" path="remark" class="md:col-span-2">
          <NInput v-model:value="formModel.remark" type="textarea" :rows="4" placeholder="补充记录这个部门的职责、边界或特殊说明" />
        </NFormItem>
      </NForm>

      <template #footer>
        <div class="flex justify-end gap-3">
          <NButton quaternary @click="formVisible = false">取消</NButton>
          <NButton type="primary" :loading="saving" @click="handleSubmit">保存</NButton>
        </div>
      </template>
    </NModal>
  </main>
</template>
