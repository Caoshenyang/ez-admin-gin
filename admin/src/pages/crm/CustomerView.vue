<script setup lang="ts">
import type { DataTableColumns, FormInst, FormRules } from 'naive-ui'
import {
  NButton,
  NCard,
  NDataTable,
  NEmpty,
  NForm,
  NFormItem,
  NInput,
  NModal,
  NPagination,
  NPopconfirm,
  NSelect,
  NSpace,
  NTag,
  useMessage,
} from 'naive-ui'
import { computed, h, onMounted, reactive, ref } from 'vue'

import {
  createCustomer,
  getCustomers,
  updateCustomer,
  updateCustomerStatus,
} from '../../api/customer'
import { buttonPermissionCodes } from '../../router/dynamic-menu'
import {
  CustomerStatus,
  type CreateCustomerPayload,
  type CustomerItem,
  type CustomerListQuery,
} from '../../types/customer'

interface CustomerFormModel extends CreateCustomerPayload {
  id: number
}

const message = useMessage()
const loading = ref(false)
const saving = ref(false)
const customers = ref<CustomerItem[]>([])
const total = ref(0)
const modalVisible = ref(false)
const modalMode = ref<'create' | 'edit'>('create')
const formRef = ref<FormInst | null>(null)

const query = reactive<CustomerListQuery>({
  page: 1,
  page_size: 10,
  keyword: '',
  level: '',
  source: '',
  status: 0,
})

const formModel = reactive<CustomerFormModel>({
  id: 0,
  name: '',
  contact_name: '',
  phone: '',
  level: '',
  source: '',
  status: CustomerStatus.Enabled,
  remark: '',
})

const statusOptions = [
  { label: '状态：全部', value: 0 },
  { label: '启用', value: CustomerStatus.Enabled },
  { label: '停用', value: CustomerStatus.Disabled },
]

const statusFormOptions = [
  { label: '启用', value: CustomerStatus.Enabled },
  { label: '停用', value: CustomerStatus.Disabled },
]

const levelOptions = [
  { label: '等级：全部', value: '' },
  { label: 'A 类', value: 'a' },
  { label: 'B 类', value: 'b' },
  { label: 'VIP', value: 'vip' },
]

const sourceOptions = [
  { label: '来源：全部', value: '' },
  { label: '转介绍', value: 'referral' },
  { label: '广告投放', value: 'ads' },
  { label: '线下拜访', value: 'offline' },
]

const levelFormOptions = [
  { label: 'A 类', value: 'a' },
  { label: 'B 类', value: 'b' },
  { label: 'VIP', value: 'vip' },
]

const sourceFormOptions = [
  { label: '转介绍', value: 'referral' },
  { label: '广告投放', value: 'ads' },
  { label: '线下拜访', value: 'offline' },
]

const formRules: FormRules = {
  name: [{ required: true, message: '请输入客户名称', trigger: ['blur', 'input'] }],
  contact_name: [{ max: 64, message: '联系人不能超过 64 个字符', trigger: ['blur', 'input'] }],
  phone: [{ max: 32, message: '联系电话不能超过 32 个字符', trigger: ['blur', 'input'] }],
  level: [{ max: 32, message: '客户等级不能超过 32 个字符', trigger: ['blur', 'change'] }],
  source: [{ max: 32, message: '客户来源不能超过 32 个字符', trigger: ['blur', 'change'] }],
  remark: [{ max: 255, message: '备注不能超过 255 个字符', trigger: ['blur', 'input'] }],
}

const hasRows = computed(() => customers.value.length > 0)

const columns: DataTableColumns<CustomerItem> = [
  {
    title: '客户',
    key: 'name',
    minWidth: 240,
    render(row) {
      return h('div', { class: 'leading-5' }, [
        h('p', { class: 'font-medium text-[#111827]' }, row.name),
        h('p', { class: 'text-xs text-[#6B7280]' }, row.contact_name || '未填写联系人'),
      ])
    },
  },
  {
    title: '联系方式',
    key: 'phone',
    width: 150,
    render(row) {
      return h('span', { class: 'text-sm text-[#334155]' }, row.phone || '-')
    },
  },
  {
    title: '等级 / 来源',
    key: 'level',
    width: 160,
    render(row) {
      const text = [formatLevel(row.level), formatSource(row.source)].join(' / ')
      return h('span', { class: 'text-sm text-[#334155]' }, text)
    },
  },
  {
    title: '归属',
    key: 'owner_nickname',
    width: 190,
    render(row) {
      const owner = row.owner_nickname || row.owner_username || `#${row.owner_user_id}`
      return h('div', { class: 'leading-5' }, [
        h('p', { class: 'text-sm font-medium text-[#111827]' }, owner),
        h('p', { class: 'text-xs text-[#6B7280]' }, row.department_name || `部门 #${row.department_id}`),
      ])
    },
  },
  {
    title: '状态',
    key: 'status',
    width: 90,
    render(row) {
      return h(
        NTag,
        { bordered: false, type: row.status === CustomerStatus.Enabled ? 'success' : 'error' },
        { default: () => (row.status === CustomerStatus.Enabled ? '启用' : '停用') },
      )
    },
  },
  {
    title: '创建时间',
    key: 'created_at',
    width: 180,
    render(row) {
      return formatTime(row.created_at)
    },
  },
  {
    title: '操作',
    key: 'actions',
    width: 180,
    fixed: 'right',
    render(row) {
      const nextStatus = row.status === CustomerStatus.Enabled ? CustomerStatus.Disabled : CustomerStatus.Enabled

      return h(
        NSpace,
        { size: 8 },
        {
          default: () =>
            [
              canUse('crm:customer:update')
                ? h(
                    NButton,
                    { size: 'small', ghost: true, type: 'primary', onClick: () => openEditModal(row) },
                    { default: () => '编辑' },
                  )
                : null,
              canUse('crm:customer:status')
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
                            type: nextStatus === CustomerStatus.Disabled ? 'error' : 'success',
                          },
                          { default: () => (nextStatus === CustomerStatus.Disabled ? '停用' : '启用') },
                        ),
                      default: () => `确认${nextStatus === CustomerStatus.Disabled ? '停用' : '启用'}该客户？`,
                    },
                  )
                : null,
            ].filter(Boolean),
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

function formatLevel(value: string) {
  switch (value) {
    case 'a':
      return 'A 类'
    case 'b':
      return 'B 类'
    case 'vip':
      return 'VIP'
    default:
      return value || '未分级'
  }
}

function formatSource(value: string) {
  switch (value) {
    case 'referral':
      return '转介绍'
    case 'ads':
      return '广告投放'
    case 'offline':
      return '线下拜访'
    default:
      return value || '未知来源'
  }
}

async function fetchCustomers() {
  loading.value = true
  try {
    const result = await getCustomers(query)
    customers.value = result.items
    total.value = result.total
  } catch (error) {
    console.error(error)
    message.error('客户列表加载失败')
  } finally {
    loading.value = false
  }
}

function handleSearch() {
  query.page = 1
  void fetchCustomers()
}

function resetFilters() {
  query.keyword = ''
  query.level = ''
  query.source = ''
  query.status = 0
  query.page = 1
  void fetchCustomers()
}

function openCreateModal() {
  modalMode.value = 'create'
  modalVisible.value = true
  resetForm()
}

function openEditModal(item: CustomerItem) {
  modalMode.value = 'edit'
  modalVisible.value = true
  formModel.id = item.id
  formModel.name = item.name
  formModel.contact_name = item.contact_name
  formModel.phone = item.phone
  formModel.level = item.level
  formModel.source = item.source
  formModel.status = item.status
  formModel.remark = item.remark
}

function resetForm() {
  formModel.id = 0
  formModel.name = ''
  formModel.contact_name = ''
  formModel.phone = ''
  formModel.level = ''
  formModel.source = ''
  formModel.status = CustomerStatus.Enabled
  formModel.remark = ''
  formRef.value?.restoreValidation()
}

function closeModal() {
  modalVisible.value = false
  resetForm()
}

async function handleSubmit() {
  try {
    await formRef.value?.validate()
  } catch {
    return
  }

  saving.value = true
  try {
    if (modalMode.value === 'create') {
      await createCustomer({
        name: formModel.name,
        contact_name: formModel.contact_name,
        phone: formModel.phone,
        level: formModel.level,
        source: formModel.source,
        status: formModel.status,
        remark: formModel.remark,
      })
      message.success('客户创建成功')
    } else {
      await updateCustomer(formModel.id, {
        name: formModel.name,
        contact_name: formModel.contact_name,
        phone: formModel.phone,
        level: formModel.level,
        source: formModel.source,
        status: formModel.status,
        remark: formModel.remark,
      })
      message.success('客户更新成功')
    }

    closeModal()
    await fetchCustomers()
  } catch (error) {
    console.error(error)
    message.error(modalMode.value === 'create' ? '客户创建失败' : '客户更新失败')
  } finally {
    saving.value = false
  }
}

async function handleToggleStatus(item: CustomerItem, status: CustomerStatus) {
  try {
    await updateCustomerStatus(item.id, status)
    message.success(status === CustomerStatus.Enabled ? '客户已启用' : '客户已停用')
    await fetchCustomers()
  } catch (error) {
    console.error(error)
    message.error('客户状态更新失败')
  }
}

onMounted(() => {
  void fetchCustomers()
})
</script>

<template>
  <div class="flex flex-col gap-5">
    <NCard :bordered="false" class="rounded-3xl shadow-[0_24px_80px_rgba(15,23,42,0.06)]">
      <div class="flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between">
        <div>
          <p class="text-sm font-medium uppercase tracking-[0.24em] text-[#EA580C]">CRM / Customers</p>
          <h1 class="mt-2 text-2xl font-semibold text-[#0F172A]">客户档案</h1>
          <p class="mt-2 max-w-2xl text-sm leading-6 text-[#475569]">
            这页是一个真实的非 system 分组业务模块示例。它把客户资源接进了菜单、按钮、接口权限和数据权限链路，
            用来证明第 8 章的模块接入规范已经可以稳定落地在业务域。
          </p>
        </div>

        <NSpace align="center" :wrap="true" size="small">
          <NInput
            v-model:value="query.keyword"
            clearable
            placeholder="搜索客户 / 联系人 / 电话 / 负责人"
            class="w-[260px]"
            @keyup.enter="handleSearch"
          />
          <NSelect
            v-model:value="query.level"
            :options="levelOptions"
            class="w-[140px]"
          />
          <NSelect
            v-model:value="query.source"
            :options="sourceOptions"
            class="w-[150px]"
          />
          <NSelect
            v-model:value="query.status"
            :options="statusOptions"
            class="w-[130px]"
          />
          <NButton type="primary" @click="handleSearch">搜索</NButton>
          <NButton quaternary @click="resetFilters">重置</NButton>
          <NButton v-if="canUse('crm:customer:create')" type="warning" @click="openCreateModal">
            新建客户
          </NButton>
        </NSpace>
      </div>
    </NCard>

    <NCard :bordered="false" class="rounded-3xl shadow-[0_24px_80px_rgba(15,23,42,0.06)]">
      <div v-if="hasRows" class="flex flex-col gap-4">
        <NDataTable
          :columns="columns"
          :data="customers"
          :loading="loading"
          :pagination="false"
          size="small"
          scroll-x="1120"
        />

        <div class="flex justify-end">
          <NPagination
            v-model:page="query.page"
            v-model:page-size="query.page_size"
            :item-count="total"
            :page-sizes="[10, 20, 50]"
            show-size-picker
            @update:page="fetchCustomers"
            @update:page-size="fetchCustomers"
          />
        </div>
      </div>

      <NEmpty
        v-else
        description="当前筛选条件下还没有客户记录。"
        class="rounded-3xl border border-dashed border-[#E2E8F0] py-14"
      >
        <template v-if="canUse('crm:customer:create')" #extra>
          <NButton type="warning" @click="openCreateModal">创建第一条客户</NButton>
        </template>
      </NEmpty>
    </NCard>

    <NModal v-model:show="modalVisible" preset="card" class="w-[640px]" :title="modalMode === 'create' ? '新建客户' : '编辑客户'">
      <NForm ref="formRef" :model="formModel" :rules="formRules" label-placement="top" class="grid grid-cols-1 gap-x-4 md:grid-cols-2">
        <NFormItem label="客户名称" path="name" class="md:col-span-2">
          <NInput v-model:value="formModel.name" placeholder="例如：星河工业设备有限公司" />
        </NFormItem>
        <NFormItem label="联系人" path="contact_name">
          <NInput v-model:value="formModel.contact_name" placeholder="联系人姓名" />
        </NFormItem>
        <NFormItem label="联系电话" path="phone">
          <NInput v-model:value="formModel.phone" placeholder="手机号 / 座机" />
        </NFormItem>
        <NFormItem label="客户等级" path="level">
          <NSelect v-model:value="formModel.level" :options="levelFormOptions" clearable placeholder="选择客户等级" />
        </NFormItem>
        <NFormItem label="客户来源" path="source">
          <NSelect v-model:value="formModel.source" :options="sourceFormOptions" clearable placeholder="选择客户来源" />
        </NFormItem>
        <NFormItem label="状态" path="status">
          <NSelect v-model:value="formModel.status" :options="statusFormOptions" />
        </NFormItem>
        <NFormItem label="备注" path="remark" class="md:col-span-2">
          <NInput v-model:value="formModel.remark" type="textarea" :rows="4" placeholder="记录客户补充说明、交付节奏或跟进重点" />
        </NFormItem>
      </NForm>

      <template #footer>
        <div class="flex justify-end gap-3">
          <NButton quaternary @click="closeModal">取消</NButton>
          <NButton type="primary" :loading="saving" @click="handleSubmit">
            {{ modalMode === 'create' ? '创建客户' : '保存修改' }}
          </NButton>
        </div>
      </template>
    </NModal>
  </div>
</template>
