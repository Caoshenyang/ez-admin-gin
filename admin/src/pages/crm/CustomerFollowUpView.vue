<script setup lang="ts">
import type { DataTableColumns, FormInst, FormRules } from 'naive-ui'
import {
  NButton,
  NCard,
  NDataTable,
  NDatePicker,
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
  createCustomerFollowUp,
  getCustomerFollowUpCustomerOptions,
  getCustomerFollowUps,
  updateCustomerFollowUp,
  updateCustomerFollowUpStatus,
} from '../../api/followup'
import { buttonPermissionCodes } from '../../router/dynamic-menu'
import {
  CustomerFollowUpStatus,
  type CreateCustomerFollowUpPayload,
  type CustomerFollowUpCustomerOption,
  type CustomerFollowUpItem,
  type CustomerFollowUpListQuery,
} from '../../types/followup'

interface FollowUpFormModel extends CreateCustomerFollowUpPayload {
  id: number
  next_follow_at_ms: number | null
}

const message = useMessage()
const loading = ref(false)
const saving = ref(false)
const customerLoading = ref(false)
const followUps = ref<CustomerFollowUpItem[]>([])
const customerOptions = ref<CustomerFollowUpCustomerOption[]>([])
const total = ref(0)
const modalVisible = ref(false)
const modalMode = ref<'create' | 'edit'>('create')
const formRef = ref<FormInst | null>(null)

const query = reactive<CustomerFollowUpListQuery>({
  page: 1,
  page_size: 10,
  keyword: '',
  follow_type: '',
  customer_id: 0,
  status: 0,
})

const formModel = reactive<FollowUpFormModel>({
  id: 0,
  customer_id: 0,
  follow_type: 'phone',
  subject: '',
  content: '',
  result: '',
  next_follow_at: null,
  next_follow_at_ms: null,
  status: CustomerFollowUpStatus.Pending,
})

const statusOptions = [
  { label: '状态：全部', value: 0 },
  { label: '待跟进', value: CustomerFollowUpStatus.Pending },
  { label: '已完成', value: CustomerFollowUpStatus.Completed },
  { label: '已关闭', value: CustomerFollowUpStatus.Closed },
]

const statusFormOptions = [
  { label: '待跟进', value: CustomerFollowUpStatus.Pending },
  { label: '已完成', value: CustomerFollowUpStatus.Completed },
  { label: '已关闭', value: CustomerFollowUpStatus.Closed },
]

const followTypeOptions = [
  { label: '方式：全部', value: '' },
  { label: '电话', value: 'phone' },
  { label: '微信', value: 'wechat' },
  { label: '上门', value: 'visit' },
  { label: '会议', value: 'meeting' },
]

const followTypeFormOptions = followTypeOptions.filter((item) => item.value !== '')

const formRules: FormRules = {
  customer_id: [{ required: true, type: 'number', message: '请选择客户', trigger: ['change', 'blur'] }],
  follow_type: [{ required: true, message: '请选择跟进方式', trigger: ['change', 'blur'] }],
  subject: [{ required: true, message: '请输入跟进主题', trigger: ['blur', 'input'] }],
  content: [{ required: true, message: '请输入跟进内容', trigger: ['blur', 'input'] }],
  result: [{ max: 255, message: '跟进结果不能超过 255 个字符', trigger: ['blur', 'input'] }],
}

const hasRows = computed(() => followUps.value.length > 0)

const customerSelectOptions = computed(() =>
  customerOptions.value.map((item) => ({
    label: `${item.name} · ${item.owner_nickname || item.owner_username || `#${item.owner_user_id}`}`,
    value: item.id,
  })),
)

const columns: DataTableColumns<CustomerFollowUpItem> = [
  {
    title: '客户 / 主题',
    key: 'subject',
    minWidth: 260,
    render(row) {
      return h('div', { class: 'leading-5' }, [
        h('p', { class: 'font-medium text-[#111827]' }, row.subject),
        h('p', { class: 'text-xs text-[#6B7280]' }, row.customer_name),
      ])
    },
  },
  {
    title: '方式 / 内容',
    key: 'follow_type',
    minWidth: 220,
    render(row) {
      return h('div', { class: 'leading-5' }, [
        h('p', { class: 'text-sm font-medium text-[#111827]' }, formatFollowType(row.follow_type)),
        h('p', { class: 'line-clamp-2 text-xs text-[#6B7280]' }, row.content),
      ])
    },
  },
  {
    title: '负责人',
    key: 'owner_nickname',
    width: 180,
    render(row) {
      const owner = row.owner_nickname || row.owner_username || `#${row.owner_user_id}`
      return h('div', { class: 'leading-5' }, [
        h('p', { class: 'text-sm font-medium text-[#111827]' }, owner),
        h('p', { class: 'text-xs text-[#6B7280]' }, row.department_name || `部门 #${row.department_id}`),
      ])
    },
  },
  {
    title: '下次跟进',
    key: 'next_follow_at',
    width: 180,
    render(row) {
      return h('span', { class: 'text-sm text-[#334155]' }, formatTime(row.next_follow_at))
    },
  },
  {
    title: '状态',
    key: 'status',
    width: 100,
    render(row) {
      return h(
        NTag,
        { bordered: false, type: resolveStatusTagType(row.status) },
        { default: () => formatStatus(row.status) },
      )
    },
  },
  {
    title: '操作',
    key: 'actions',
    width: 200,
    fixed: 'right',
    render(row) {
      const nextStatus = resolveNextStatus(row.status)
      return h(
        NSpace,
        { size: 8 },
        {
          default: () =>
            [
              canUse('crm:followup:update')
                ? h(
                    NButton,
                    { size: 'small', ghost: true, type: 'primary', onClick: () => openEditModal(row) },
                    { default: () => '编辑' },
                  )
                : null,
              canUse('crm:followup:status') && nextStatus
                ? h(
                    NPopconfirm,
                    { onPositiveClick: () => handleToggleStatus(row, nextStatus.value) },
                    {
                      trigger: () =>
                        h(
                          NButton,
                          { size: 'small', ghost: true, type: nextStatus.type },
                          { default: () => nextStatus.label },
                        ),
                      default: () => `确认将该跟进改为“${nextStatus.label}”吗？`,
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

function resolveStatusTagType(status: number) {
  switch (status) {
    case CustomerFollowUpStatus.Completed:
      return 'success'
    case CustomerFollowUpStatus.Closed:
      return 'default'
    default:
      return 'warning'
  }
}

function resolveNextStatus(status: number) {
  switch (status) {
    case CustomerFollowUpStatus.Pending:
      return { value: CustomerFollowUpStatus.Completed, label: '标记完成', type: 'success' as const }
    case CustomerFollowUpStatus.Completed:
      return { value: CustomerFollowUpStatus.Closed, label: '关闭', type: 'default' as const }
    default:
      return null
  }
}

function formatStatus(status: number) {
  switch (status) {
    case CustomerFollowUpStatus.Completed:
      return '已完成'
    case CustomerFollowUpStatus.Closed:
      return '已关闭'
    default:
      return '待跟进'
  }
}

function formatFollowType(value: string) {
  switch (value) {
    case 'wechat':
      return '微信'
    case 'visit':
      return '上门'
    case 'meeting':
      return '会议'
    default:
      return '电话'
  }
}

function formatTime(value: string | null) {
  return value ? new Date(value).toLocaleString() : '未安排'
}

function resetForm() {
  formModel.id = 0
  formModel.customer_id = 0
  formModel.follow_type = 'phone'
  formModel.subject = ''
  formModel.content = ''
  formModel.result = ''
  formModel.next_follow_at = null
  formModel.next_follow_at_ms = null
  formModel.status = CustomerFollowUpStatus.Pending
}

function openCreateModal() {
  modalMode.value = 'create'
  modalVisible.value = true
  resetForm()
  void loadCustomerOptions()
}

function openEditModal(row: CustomerFollowUpItem) {
  modalMode.value = 'edit'
  modalVisible.value = true
  formModel.id = row.id
  formModel.customer_id = row.customer_id
  formModel.follow_type = row.follow_type
  formModel.subject = row.subject
  formModel.content = row.content
  formModel.result = row.result
  formModel.next_follow_at = row.next_follow_at
  formModel.next_follow_at_ms = row.next_follow_at ? new Date(row.next_follow_at).getTime() : null
  formModel.status = row.status
  void loadCustomerOptions()
}

function closeModal() {
  modalVisible.value = false
  formRef.value?.restoreValidation()
}

async function loadCustomerOptions(keyword = '') {
  customerLoading.value = true
  try {
    customerOptions.value = await getCustomerFollowUpCustomerOptions({ keyword, limit: 100 })
  }
  catch (error) {
    console.error(error)
    message.error('加载客户选项失败')
  }
  finally {
    customerLoading.value = false
  }
}

async function loadFollowUps() {
  loading.value = true
  try {
    const result = await getCustomerFollowUps(query)
    followUps.value = result.items
    total.value = result.total
    query.page = result.page
    query.page_size = result.page_size
  }
  catch (error) {
    console.error(error)
    message.error('加载客户跟进失败')
  }
  finally {
    loading.value = false
  }
}

async function handleSubmit() {
  try {
    saving.value = true
    await formRef.value?.validate()

    const createPayload: CreateCustomerFollowUpPayload = {
      customer_id: formModel.customer_id,
      follow_type: formModel.follow_type,
      subject: formModel.subject.trim(),
      content: formModel.content.trim(),
      result: formModel.result.trim(),
      next_follow_at: formModel.next_follow_at_ms ? new Date(formModel.next_follow_at_ms).toISOString() : null,
      status: formModel.status,
    }

    if (modalMode.value === 'create') {
      await createCustomerFollowUp(createPayload)
      message.success('客户跟进创建成功')
    }
    else {
      await updateCustomerFollowUp(formModel.id, {
        follow_type: createPayload.follow_type,
        subject: createPayload.subject,
        content: createPayload.content,
        result: createPayload.result,
        next_follow_at: createPayload.next_follow_at,
        status: createPayload.status,
      })
      message.success('客户跟进更新成功')
    }

    closeModal()
    await loadFollowUps()
  }
  catch (error) {
    if (error) {
      console.error(error)
      message.error(modalMode.value === 'create' ? '客户跟进创建失败' : '客户跟进更新失败')
    }
  }
  finally {
    saving.value = false
  }
}

async function handleToggleStatus(row: CustomerFollowUpItem, status: number) {
  try {
    await updateCustomerFollowUpStatus(row.id, status)
    message.success('客户跟进状态已更新')
    await loadFollowUps()
  }
  catch (error) {
    console.error(error)
    message.error('客户跟进状态更新失败')
  }
}

function handleSearch() {
  query.page = 1
  void loadFollowUps()
}

function handleReset() {
  query.keyword = ''
  query.follow_type = ''
  query.customer_id = 0
  query.status = 0
  query.page = 1
  void loadFollowUps()
}

function handlePageChange(page: number) {
  query.page = page
  void loadFollowUps()
}

function handlePageSizeChange(pageSize: number) {
  query.page_size = pageSize
  query.page = 1
  void loadFollowUps()
}

onMounted(() => {
  void loadCustomerOptions()
  void loadFollowUps()
})
</script>

<template>
  <section class="min-h-full bg-[#F6F8FB] p-4 md:p-6">
    <NCard
      title="客户跟进"
      :bordered="false"
      size="small"
      class="rounded-2xl shadow-[0_18px_48px_rgba(15,23,42,0.06)]"
    >
      <div class="space-y-4">
        <div class="flex flex-col gap-3 lg:flex-row lg:items-center lg:justify-between">
          <NSpace wrap>
            <NInput
              v-model:value="query.keyword"
              clearable
              placeholder="搜索主题、客户、负责人"
              class="w-full md:w-[240px]"
              @keyup.enter="handleSearch"
            />
            <NSelect
              v-model:value="query.follow_type"
              :options="followTypeOptions"
              class="w-full md:w-[160px]"
            />
            <NSelect
              v-model:value="query.customer_id"
              :options="[{ label: '客户：全部', value: 0 }, ...customerSelectOptions]"
              filterable
              clearable
              class="w-full md:w-[260px]"
              :loading="customerLoading"
            />
            <NSelect
              v-model:value="query.status"
              :options="statusOptions"
              class="w-full md:w-[160px]"
            />
          </NSpace>
          <NSpace>
            <NButton ghost @click="handleReset">
              重置
            </NButton>
            <NButton type="primary" @click="handleSearch">
              查询
            </NButton>
            <NButton v-if="canUse('crm:followup:create')" type="primary" secondary @click="openCreateModal">
              新建跟进
            </NButton>
          </NSpace>
        </div>

        <NDataTable
          :columns="columns"
          :data="followUps"
          :loading="loading"
          :pagination="false"
          :bordered="false"
          :single-line="false"
          size="small"
          class="rounded-2xl bg-white"
        />

        <div v-if="!loading && !hasRows" class="rounded-2xl border border-dashed border-[#D8E1EC] bg-[#FBFCFD] py-14">
          <NEmpty description="当前筛选条件下还没有客户跟进记录" />
        </div>

        <div class="flex justify-end">
          <NPagination
            :page="query.page"
            :page-size="query.page_size"
            :page-count="Math.max(1, Math.ceil(total / query.page_size))"
            show-size-picker
            :page-sizes="[10, 20, 50]"
            @update:page="handlePageChange"
            @update:page-size="handlePageSizeChange"
          />
        </div>
      </div>
    </NCard>

    <NModal v-model:show="modalVisible" preset="card" :title="modalMode === 'create' ? '新建客户跟进' : '编辑客户跟进'" class="max-w-[720px]">
      <NForm ref="formRef" :model="formModel" :rules="formRules" label-placement="top">
        <div class="grid gap-4 md:grid-cols-2">
          <NFormItem label="客户" path="customer_id">
            <NSelect
              v-model:value="formModel.customer_id"
              :options="customerSelectOptions"
              filterable
              :loading="customerLoading"
              :disabled="modalMode === 'edit'"
              placeholder="请选择客户"
            />
          </NFormItem>
          <NFormItem label="跟进方式" path="follow_type">
            <NSelect
              v-model:value="formModel.follow_type"
              :options="followTypeFormOptions"
              placeholder="请选择跟进方式"
            />
          </NFormItem>
          <NFormItem label="跟进主题" path="subject" class="md:col-span-2">
            <NInput v-model:value="formModel.subject" maxlength="128" placeholder="例如：方案评审会" />
          </NFormItem>
          <NFormItem label="跟进内容" path="content" class="md:col-span-2">
            <NInput
              v-model:value="formModel.content"
              type="textarea"
              :rows="4"
              maxlength="1000"
              placeholder="记录本次沟通内容、客户反馈、推进动作"
            />
          </NFormItem>
          <NFormItem label="跟进结果" path="result" class="md:col-span-2">
            <NInput
              v-model:value="formModel.result"
              type="textarea"
              :rows="2"
              maxlength="255"
              placeholder="例如：已同意进入报价阶段"
            />
          </NFormItem>
          <NFormItem label="下次跟进时间">
            <NDatePicker
              v-model:value="formModel.next_follow_at_ms"
              type="datetime"
              clearable
              class="w-full"
            />
          </NFormItem>
          <NFormItem label="状态" path="status">
            <NSelect
              v-model:value="formModel.status"
              :options="statusFormOptions"
              placeholder="请选择状态"
            />
          </NFormItem>
        </div>
      </NForm>

      <template #action>
        <NSpace justify="end">
          <NButton @click="closeModal">
            取消
          </NButton>
          <NButton type="primary" :loading="saving" @click="handleSubmit">
            {{ modalMode === 'create' ? '创建' : '保存' }}
          </NButton>
        </NSpace>
      </template>
    </NModal>
  </section>
</template>
