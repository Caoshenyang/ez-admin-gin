<script setup lang="ts">
import { CloseOutline } from '@vicons/ionicons5'
import type { DataTableColumns, FormRules } from 'naive-ui'
import {
  NAlert,
  NButton,
  NCard,
  NDataTable,
  NEmpty,
  NForm,
  NFormItem,
  NIcon,
  NInput,
  NInputNumber,
  NModal,
  NPagination,
  NPopconfirm,
  NSelect,
  NSpace,
  NTag,
  useMessage,
} from 'naive-ui'
import { h, onMounted, reactive, ref } from 'vue'

import { useModalForm } from '@/composables/useModalForm'
import { usePermission } from '@/composables/usePermission'
import { useStatusToggle } from '@/composables/useStatusToggle'
import { STATUS_FILTER_OPTIONS, STATUS_FORM_OPTIONS } from '@/constants/status'
import { formatTime } from '@/utils/format'
import {
  createDictItem,
  createDictType,
  getDictItems,
  getDictTypes,
  updateDictItem,
  updateDictItemStatus,
  updateDictType,
  updateDictTypeStatus,
} from '../api/dict'
import {
  DictStatus,
  type DictItem,
  type DictItemListQuery,
  type DictTypeItem,
  type DictTypeListQuery,
} from '../types/dict'

interface DictTypeFormModel {
  id: number
  code: string
  name: string
  sort: number
  status: DictStatus
  remark: string
}

interface DictItemFormModel {
  id: number
  type_id: number
  item_key: string
  label: string
  value: string
  tag_type: string
  sort: number
  status: DictStatus
  remark: string
}

const message = useMessage()
const { canUse } = usePermission()
const successText = ref('')

const typeLoading = ref(false)
const itemLoading = ref(false)

const dictTypes = ref<DictTypeItem[]>([])
const dictTypeTotal = ref(0)
const selectedTypeID = ref<number | null>(null)
const selectedType = ref<DictTypeItem | null>(null)

const dictItems = ref<DictItem[]>([])
const dictItemTotal = ref(0)

const typeQuery = reactive<DictTypeListQuery>({
  page: 1,
  page_size: 10,
  keyword: '',
  status: 0,
})

const itemQuery = reactive<DictItemListQuery>({
  page: 1,
  page_size: 10,
  type_id: 0,
  keyword: '',
  status: 0,
})

function defaultTypeFormModel(): DictTypeFormModel {
  return {
    id: 0,
    code: '',
    name: '',
    sort: 10,
    status: DictStatus.Enabled,
    remark: '',
  }
}

function defaultItemFormModel(): DictItemFormModel {
  return {
    id: 0,
    type_id: selectedTypeID.value ?? 0,
    item_key: '',
    label: '',
    value: '',
    tag_type: '',
    sort: 10,
    status: DictStatus.Enabled,
    remark: '',
  }
}

const {
  formRef: typeFormRef,
  formVisible: typeFormVisible,
  formMode: typeFormMode,
  formModel: typeFormModel,
  saving: typeSaving,
  rules: typeRules,
  openCreate: openTypeCreateBase,
  openEdit: openTypeEditBase,
  handleSubmit: handleTypeSubmit,
} = useModalForm<DictTypeFormModel>(defaultTypeFormModel, {
  rules: {
    code: [{ required: true, message: '请输入字典编码', trigger: 'blur' }],
    name: [{ required: true, message: '请输入字典名称', trigger: 'blur' }],
  } as FormRules,
})

const {
  formRef: itemFormRef,
  formVisible: itemFormVisible,
  formMode: itemFormMode,
  formModel: itemFormModel,
  saving: itemSaving,
  rules: itemRules,
  openCreate: openItemCreateBase,
  openEdit: openItemEditBase,
  handleSubmit: handleItemSubmit,
} = useModalForm<DictItemFormModel>(defaultItemFormModel, {
  rules: {
    item_key: [{ required: true, message: '请输入字典项编码', trigger: 'blur' }],
    label: [{ required: true, message: '请输入字典项名称', trigger: 'blur' }],
    value: [{ required: true, message: '请输入字典项值', trigger: 'blur' }],
  } as FormRules,
})

const { handleToggleStatus: handleToggleTypeStatus } = useStatusToggle<DictTypeItem>(updateDictTypeStatus, {
  onSuccess: async () => {
    successText.value = '字典类型状态已更新'
    await loadDictTypes()
  },
})

const { handleToggleStatus: handleToggleItemStatus } = useStatusToggle<DictItem>(updateDictItemStatus, {
  onSuccess: async () => {
    successText.value = '字典项状态已更新'
    await loadDictItems()
  },
})

const typeColumns: DataTableColumns<DictTypeItem> = [
  {
    title: '字典类型',
    key: 'code',
    minWidth: 220,
    render(row) {
      return h('div', { class: 'leading-6' }, [
        h('p', { class: 'font-semibold text-[#111827]' }, row.name),
        h('p', { class: 'text-xs text-[#6B7280]' }, row.code),
      ])
    },
  },
  {
    title: '排序',
    key: 'sort',
    width: 76,
  },
  {
    title: '状态',
    key: 'status',
    width: 90,
    render(row) {
      return h(
        NTag,
        { bordered: false, type: row.status === DictStatus.Enabled ? 'success' : 'error' },
        { default: () => (row.status === DictStatus.Enabled ? '启用' : '禁用') },
      )
    },
  },
  {
    title: '操作',
    key: 'actions',
    width: 176,
    fixed: 'right',
    render(row) {
      const nextStatus = row.status === DictStatus.Enabled ? DictStatus.Disabled : DictStatus.Enabled

      return h(
        NSpace,
        { size: 8, align: 'center' },
        {
          default: () =>
            [
              canUse('system:dict:type:update')
                ? h(
                    NButton,
                    { size: 'small', ghost: true, type: 'info', onClick: () => openTypeEdit(row) },
                    { default: () => '编辑' },
                  )
                : null,
              canUse('system:dict:type:status')
                ? h(
                    NPopconfirm,
                    { onPositiveClick: () => handleToggleTypeStatus(row, nextStatus) },
                    {
                      trigger: () =>
                        h(
                          NButton,
                          {
                            size: 'small',
                            ghost: true,
                            type: nextStatus === DictStatus.Disabled ? 'error' : 'success',
                          },
                          { default: () => (nextStatus === DictStatus.Disabled ? '禁用' : '启用') },
                        ),
                      default: () => `确认${nextStatus === DictStatus.Disabled ? '禁用' : '启用'}该字典类型？`,
                    },
                  )
                : null,
            ].filter(Boolean),
        },
      )
    },
  },
]

const itemColumns: DataTableColumns<DictItem> = [
  {
    title: '字典项',
    key: 'item_key',
    minWidth: 220,
    render(row) {
      return h('div', { class: 'leading-6' }, [
        h('p', { class: 'font-semibold text-[#111827]' }, row.label),
        h('p', { class: 'text-xs text-[#6B7280]' }, `${row.item_key} · ${row.value}`),
      ])
    },
  },
  {
    title: '标签样式',
    key: 'tag_type',
    width: 120,
    render(row) {
      if (!row.tag_type) {
        return h('span', { class: 'text-[#9CA3AF]' }, '-')
      }

      return h(
        NTag,
        { size: 'small', bordered: false, type: toTagType(row.tag_type) },
        { default: () => row.tag_type },
      )
    },
  },
  {
    title: '排序',
    key: 'sort',
    width: 76,
  },
  {
    title: '状态',
    key: 'status',
    width: 90,
    render(row) {
      return h(
        NTag,
        { bordered: false, type: row.status === DictStatus.Enabled ? 'success' : 'error' },
        { default: () => (row.status === DictStatus.Enabled ? '启用' : '禁用') },
      )
    },
  },
  {
    title: '操作',
    key: 'actions',
    width: 176,
    fixed: 'right',
    render(row) {
      const nextStatus = row.status === DictStatus.Enabled ? DictStatus.Disabled : DictStatus.Enabled

      return h(
        NSpace,
        { size: 8, align: 'center' },
        {
          default: () =>
            [
              canUse('system:dict:item:update')
                ? h(
                    NButton,
                    { size: 'small', ghost: true, type: 'info', onClick: () => openItemEdit(row) },
                    { default: () => '编辑' },
                  )
                : null,
              canUse('system:dict:item:status')
                ? h(
                    NPopconfirm,
                    { onPositiveClick: () => handleToggleItemStatus(row, nextStatus) },
                    {
                      trigger: () =>
                        h(
                          NButton,
                          {
                            size: 'small',
                            ghost: true,
                            type: nextStatus === DictStatus.Disabled ? 'error' : 'success',
                          },
                          { default: () => (nextStatus === DictStatus.Disabled ? '禁用' : '启用') },
                        ),
                      default: () => `确认${nextStatus === DictStatus.Disabled ? '禁用' : '启用'}该字典项？`,
                    },
                  )
                : null,
            ].filter(Boolean),
        },
      )
    },
  },
]

function toTagType(value: string) {
  if (value === 'success' || value === 'warning' || value === 'error' || value === 'info' || value === 'default') {
    return value
  }

  return 'default'
}

function typeRowProps(row: DictTypeItem) {
  return {
    class: row.id === selectedTypeID.value ? 'dict-type-row dict-type-row--active' : 'dict-type-row',
    onClick: () => {
      void selectType(row)
    },
  }
}

async function loadDictTypes() {
  typeLoading.value = true

  try {
    const data = await getDictTypes({
      ...typeQuery,
      keyword: typeQuery.keyword?.trim() || undefined,
      status: typeQuery.status === 0 ? undefined : typeQuery.status,
    })

    dictTypes.value = data.items
    dictTypeTotal.value = data.total

    if (!data.items.length) {
      selectedTypeID.value = null
      selectedType.value = null
      dictItems.value = []
      dictItemTotal.value = 0
      itemQuery.type_id = 0
      return
    }

    const current = data.items.find((item) => item.id === selectedTypeID.value) ?? data.items[0] ?? null
    if (current) {
      await selectType(current)
    }
  } finally {
    typeLoading.value = false
  }
}

async function loadDictItems() {
  if (!selectedTypeID.value) {
    dictItems.value = []
    dictItemTotal.value = 0
    return
  }

  itemLoading.value = true

  try {
    const data = await getDictItems({
      ...itemQuery,
      type_id: selectedTypeID.value,
      keyword: itemQuery.keyword?.trim() || undefined,
      status: itemQuery.status === 0 ? undefined : itemQuery.status,
    })

    dictItems.value = data.items
    dictItemTotal.value = data.total
  } finally {
    itemLoading.value = false
  }
}

async function selectType(row: DictTypeItem) {
  selectedTypeID.value = row.id
  selectedType.value = row
  itemQuery.page = 1
  itemQuery.type_id = row.id
  await loadDictItems()
}

function handleTypeSearch() {
  typeQuery.page = 1
  void loadDictTypes()
}

function handleTypeReset() {
  typeQuery.page = 1
  typeQuery.page_size = 10
  typeQuery.keyword = ''
  typeQuery.status = 0
  void loadDictTypes()
}

function handleItemSearch() {
  itemQuery.page = 1
  void loadDictItems()
}

function handleItemReset() {
  itemQuery.page = 1
  itemQuery.page_size = 10
  itemQuery.keyword = ''
  itemQuery.status = 0
  void loadDictItems()
}

function handleTypePageChange(page: number) {
  typeQuery.page = page
  void loadDictTypes()
}

function handleTypePageSizeChange(pageSize: number) {
  typeQuery.page = 1
  typeQuery.page_size = pageSize
  void loadDictTypes()
}

function handleItemPageChange(page: number) {
  itemQuery.page = page
  void loadDictItems()
}

function handleItemPageSizeChange(pageSize: number) {
  itemQuery.page = 1
  itemQuery.page_size = pageSize
  void loadDictItems()
}

function openTypeCreate() {
  openTypeCreateBase()
}

function openTypeEdit(row: DictTypeItem) {
  openTypeEditBase(row)
}

function openItemCreate() {
  if (!selectedTypeID.value) {
    message.warning('请先选择一个字典类型')
    return
  }

  openItemCreateBase()
  itemFormModel.type_id = selectedTypeID.value
}

function openItemEdit(row: DictItem) {
  openItemEditBase(row)
}

async function onSubmitType() {
  if (typeFormMode.value === 'create') {
    await createDictType({
      code: typeFormModel.code.trim(),
      name: typeFormModel.name.trim(),
      sort: typeFormModel.sort,
      status: typeFormModel.status,
      remark: typeFormModel.remark.trim(),
    })
    successText.value = '字典类型创建成功'
    message.success('字典类型创建成功')
  } else {
    await updateDictType(typeFormModel.id, {
      name: typeFormModel.name.trim(),
      sort: typeFormModel.sort,
      status: typeFormModel.status,
      remark: typeFormModel.remark.trim(),
    })
    successText.value = '字典类型已更新'
    message.success('字典类型更新成功')
  }

  await loadDictTypes()
}

async function onSubmitItem() {
  if (!selectedTypeID.value) {
    message.warning('请先选择一个字典类型')
    return
  }

  if (itemFormMode.value === 'create') {
    await createDictItem({
      type_id: selectedTypeID.value,
      item_key: itemFormModel.item_key.trim(),
      label: itemFormModel.label.trim(),
      value: itemFormModel.value.trim(),
      tag_type: itemFormModel.tag_type.trim(),
      sort: itemFormModel.sort,
      status: itemFormModel.status,
      remark: itemFormModel.remark.trim(),
    })
    successText.value = '字典项创建成功'
    message.success('字典项创建成功')
  } else {
    await updateDictItem(itemFormModel.id, {
      label: itemFormModel.label.trim(),
      value: itemFormModel.value.trim(),
      tag_type: itemFormModel.tag_type.trim(),
      sort: itemFormModel.sort,
      status: itemFormModel.status,
      remark: itemFormModel.remark.trim(),
    })
    successText.value = '字典项已更新'
    message.success('字典项更新成功')
  }

  await loadDictItems()
}

onMounted(() => {
  void loadDictTypes()
})
</script>

<template>
  <main class="admin-page">
    <section class="admin-page-section">
      <div class="flex items-center justify-between gap-4">
        <div>
          <h1 class="text-[26px] font-bold text-[#111827]">字典管理</h1>
          <p class="mt-1 text-sm text-[#6B7280]">先维护字典类型，再按类型维护具体字典项，供全局表单和状态映射复用。</p>
        </div>

        <NButton v-if="canUse('system:dict:type:create')" type="primary" @click="openTypeCreate">
          + 新增字典类型
        </NButton>
      </div>

      <NAlert
        v-if="successText"
        type="success"
        :show-icon="true"
        closable
        class="mx-auto w-full max-w-[560px]"
        @close="successText = ''"
      >
        {{ successText }}
      </NAlert>

      <div class="grid min-h-0 flex-1 gap-4 xl:grid-cols-[420px_minmax(0,1fr)]">
        <NCard
          class="min-h-0 rounded-lg"
          :bordered="false"
          content-style="height: 100%; padding: 0;"
        >
          <div class="dict-card-shell">
            <div class="dict-card-shell__header">
              <div>
                <p class="dict-card-shell__eyebrow">Types</p>
                <h2 class="dict-card-shell__title">字典类型</h2>
              </div>
              <NButton
                v-if="canUse('system:dict:type:create')"
                size="small"
                type="primary"
                ghost
                @click="openTypeCreate"
              >
                新增
              </NButton>
            </div>

            <div class="dict-card-shell__filters">
              <NInput
                v-model:value="typeQuery.keyword"
                clearable
                placeholder="编码 / 名称"
                @keyup.enter="handleTypeSearch"
              />
              <NSelect v-model:value="typeQuery.status" :options="STATUS_FILTER_OPTIONS" />
              <div class="dict-filter-actions">
                <NButton type="primary" @click="handleTypeSearch">查询</NButton>
                <NButton @click="handleTypeReset">重置</NButton>
              </div>
            </div>

            <NDataTable
              remote
              class="dict-table flex-1"
              :columns="typeColumns"
              :data="dictTypes"
              :loading="typeLoading"
              :pagination="false"
              :row-key="(row: DictTypeItem) => row.id"
              :row-props="typeRowProps"
              :bordered="false"
              flex-height
            />

            <div class="dict-card-shell__footer">
              <span>共 {{ dictTypeTotal }} 条</span>
              <NPagination
                :page="typeQuery.page"
                :page-size="typeQuery.page_size"
                :item-count="dictTypeTotal"
                :page-sizes="[10, 20, 50]"
                show-size-picker
                @update:page="handleTypePageChange"
                @update:page-size="handleTypePageSizeChange"
              />
            </div>
          </div>
        </NCard>

        <NCard
          class="min-h-0 rounded-lg"
          :bordered="false"
          content-style="height: 100%; padding: 0;"
        >
          <div class="dict-card-shell">
            <div class="dict-card-shell__header">
              <div>
                <p class="dict-card-shell__eyebrow">Items</p>
                <div class="flex items-center gap-2">
                  <h2 class="dict-card-shell__title">字典项</h2>
                  <NTag v-if="selectedType" size="small" type="info" :bordered="false">
                    {{ selectedType.name }}
                  </NTag>
                </div>
                <p v-if="selectedType" class="mt-1 text-xs text-[#6B7280]">
                  {{ selectedType.code }} · 最近更新 {{ formatTime(selectedType.updated_at) }}
                </p>
              </div>
              <NButton
                v-if="canUse('system:dict:item:create')"
                size="small"
                type="primary"
                ghost
                :disabled="!selectedType"
                @click="openItemCreate"
              >
                新增
              </NButton>
            </div>

            <template v-if="selectedType">
              <div class="dict-card-shell__filters">
                <NInput
                  v-model:value="itemQuery.keyword"
                  clearable
                  placeholder="编码 / 名称 / 值"
                  @keyup.enter="handleItemSearch"
                />
                <NSelect v-model:value="itemQuery.status" :options="STATUS_FILTER_OPTIONS" />
                <div class="dict-filter-actions">
                  <NButton type="primary" @click="handleItemSearch">查询</NButton>
                  <NButton @click="handleItemReset">重置</NButton>
                </div>
              </div>

              <NDataTable
                remote
                class="dict-table flex-1"
                :columns="itemColumns"
                :data="dictItems"
                :loading="itemLoading"
                :pagination="false"
                :row-key="(row: DictItem) => row.id"
                :bordered="false"
                flex-height
              />

              <div class="dict-card-shell__footer">
                <span>共 {{ dictItemTotal }} 条</span>
                <NPagination
                  :page="itemQuery.page"
                  :page-size="itemQuery.page_size"
                  :item-count="dictItemTotal"
                  :page-sizes="[10, 20, 50]"
                  show-size-picker
                  @update:page="handleItemPageChange"
                  @update:page-size="handleItemPageSizeChange"
                />
              </div>
            </template>

            <div
              v-else
              class="flex flex-1 items-center justify-center rounded-2xl border border-dashed border-[#D9DEE8] bg-[#FAFBFC] m-4"
            >
              <NEmpty description="先从左侧选择一个字典类型，再维护它的字典项。" />
            </div>
          </div>
        </NCard>
      </div>
    </section>

    <NModal
      v-model:show="typeFormVisible"
      preset="card"
      :closable="false"
      class="compact-dict-modal"
      style="width: 620px; max-width: calc(100vw - 32px)"
    >
      <template #header>
        <div class="modal-header modal-header--hero">
          <h2 class="modal-header__title">
            {{ typeFormMode === 'create' ? '新增字典类型' : '编辑字典类型' }}
          </h2>
          <p class="modal-header__hero-title">
            {{ typeFormMode === 'create' ? '先定义稳定的字典编码，再让字典项围绕它展开。' : '修改字典展示信息，不影响已有字典项的主键归属。' }}
          </p>
          <button type="button" class="modal-close" @click="typeFormVisible = false">
            <NIcon :size="18">
              <CloseOutline />
            </NIcon>
          </button>
        </div>
      </template>

      <div class="dict-modal-shell">
        <NForm
          ref="typeFormRef"
          class="compact-dict-form"
          :model="typeFormModel"
          :rules="typeRules"
          label-placement="left"
          label-width="76"
        >
          <section class="form-section form-section--primary">
            <div class="form-section__head">
              <h3>基础信息</h3>
              <p>字典编码建议使用小写字母、数字、冒号、短横线和下划线，便于在前后端直接复用。</p>
            </div>

            <div class="form-section-grid">
              <NFormItem label="字典编码" path="code">
                <NInput
                  v-model:value="typeFormModel.code"
                  placeholder="例如 common:status"
                  :disabled="typeFormMode === 'edit'"
                />
              </NFormItem>

              <NFormItem label="字典名称" path="name">
                <NInput v-model:value="typeFormModel.name" placeholder="例如 通用状态" />
              </NFormItem>

              <NFormItem label="排序">
                <NInputNumber v-model:value="typeFormModel.sort" :min="0" class="w-full" />
              </NFormItem>

              <NFormItem label="状态">
                <NSelect v-model:value="typeFormModel.status" :options="STATUS_FORM_OPTIONS" />
              </NFormItem>
            </div>
          </section>

          <section class="form-section form-section--muted">
            <NFormItem label="备注" class="mb-0">
              <NInput
                v-model:value="typeFormModel.remark"
                type="textarea"
                :rows="3"
                placeholder="补充这个字典的适用场景或业务备注"
              />
            </NFormItem>
          </section>
        </NForm>
      </div>

      <template #footer>
        <div class="modal-footer-actions">
          <NButton quaternary class="modal-footer-button" @click="typeFormVisible = false">取消</NButton>
          <NButton
            type="primary"
            class="modal-footer-button modal-footer-button--primary"
            :loading="typeSaving"
            @click="handleTypeSubmit(onSubmitType)"
          >
            保存
          </NButton>
        </div>
      </template>
    </NModal>

    <NModal
      v-model:show="itemFormVisible"
      preset="card"
      :closable="false"
      class="compact-dict-modal"
      style="width: 680px; max-width: calc(100vw - 32px)"
    >
      <template #header>
        <div class="modal-header modal-header--hero modal-header--item">
          <h2 class="modal-header__title">
            {{ itemFormMode === 'create' ? '新增字典项' : '编辑字典项' }}
          </h2>
          <p class="modal-header__hero-title">
            {{ selectedType ? `当前归属：${selectedType.name}（${selectedType.code}）` : '请选择字典类型后再维护字典项。' }}
          </p>
          <button type="button" class="modal-close" @click="itemFormVisible = false">
            <NIcon :size="18">
              <CloseOutline />
            </NIcon>
          </button>
        </div>
      </template>

      <div class="dict-modal-shell">
        <NForm
          ref="itemFormRef"
          class="compact-dict-form"
          :model="itemFormModel"
          :rules="itemRules"
          label-placement="left"
          label-width="76"
        >
          <section class="form-section form-section--primary">
            <div class="form-section__head">
              <h3>基础信息</h3>
              <p>字典项编码和显示值建议稳定设计，这样状态色、下拉项和表格标签都能长期复用。</p>
            </div>

            <div class="form-section-grid">
              <NFormItem label="字典项编码" path="item_key">
                <NInput
                  v-model:value="itemFormModel.item_key"
                  placeholder="例如 enabled"
                  :disabled="itemFormMode === 'edit'"
                />
              </NFormItem>

              <NFormItem label="字典项名称" path="label">
                <NInput v-model:value="itemFormModel.label" placeholder="例如 启用" />
              </NFormItem>

              <NFormItem label="字典项值" path="value">
                <NInput v-model:value="itemFormModel.value" placeholder="例如 1" />
              </NFormItem>

              <NFormItem label="标签样式">
                <NInput v-model:value="itemFormModel.tag_type" placeholder="例如 success / warning / info" />
              </NFormItem>

              <NFormItem label="排序">
                <NInputNumber v-model:value="itemFormModel.sort" :min="0" class="w-full" />
              </NFormItem>

              <NFormItem label="状态">
                <NSelect v-model:value="itemFormModel.status" :options="STATUS_FORM_OPTIONS" />
              </NFormItem>
            </div>
          </section>

          <section class="form-section form-section--muted">
            <NFormItem label="备注" class="mb-0">
              <NInput
                v-model:value="itemFormModel.remark"
                type="textarea"
                :rows="3"
                placeholder="可填写这项配置值的展示说明或业务备注"
              />
            </NFormItem>
          </section>
        </NForm>
      </div>

      <template #footer>
        <div class="modal-footer-actions">
          <NButton quaternary class="modal-footer-button" @click="itemFormVisible = false">取消</NButton>
          <NButton
            type="primary"
            class="modal-footer-button modal-footer-button--primary"
            :loading="itemSaving"
            @click="handleItemSubmit(onSubmitItem)"
          >
            保存
          </NButton>
        </div>
      </template>
    </NModal>
  </main>
</template>

<style scoped>
.dict-card-shell {
  display: flex;
  height: 100%;
  flex-direction: column;
}

.dict-card-shell__header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 16px;
  border-bottom: 1px solid #e5e7eb;
  padding: 18px 20px 14px;
}

.dict-card-shell__eyebrow {
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 0.14em;
  text-transform: uppercase;
  color: #94a3b8;
}

.dict-card-shell__title {
  margin-top: 6px;
  font-size: 18px;
  font-weight: 700;
  color: #111827;
}

.dict-card-shell__filters {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 140px auto;
  gap: 12px;
  padding: 16px 20px;
}

.dict-filter-actions {
  display: flex;
  gap: 10px;
}

.dict-card-shell__footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  border-top: 1px solid #e5e7eb;
  padding: 14px 20px;
  font-size: 13px;
  color: #6b7280;
}

.dict-table :deep(.n-data-table-th) {
  font-weight: 700;
  color: #374151;
  background: #fff;
}

.dict-table :deep(.n-data-table-td) {
  color: #374151;
}

.dict-table :deep(.n-data-table-tr:hover .n-data-table-td) {
  background: #f8fbff;
}

:deep(.dict-type-row) {
  cursor: pointer;
}

:deep(.dict-type-row--active .n-data-table-td) {
  background: #eef6ff;
}

.compact-dict-modal :deep(.n-card) {
  overflow: hidden;
  border-radius: 32px;
  border: 1px solid #dfe9f5;
  background: #ffffff;
  box-shadow: 0 24px 72px rgba(15, 23, 42, 0.16);
}

.compact-dict-modal :deep(.n-card-header) {
  padding: 0;
  border-bottom: 1px solid #dfe9f5;
  background: linear-gradient(135deg, #eff6ff 0%, #e8f2ff 58%, #f4f9ff 100%);
}

.compact-dict-modal :deep(.n-card__content) {
  padding: 20px 28px 10px;
}

.compact-dict-modal :deep(.n-card__footer) {
  padding: 16px 28px 24px;
  border-top: 1px solid #edf2f7;
  background: rgba(248, 250, 252, 0.85);
}

.compact-dict-form :deep(.n-form-item) {
  margin-bottom: 16px;
}

.compact-dict-form :deep(.n-form-item-label) {
  white-space: nowrap;
  align-items: center;
  padding-right: 14px;
  font-weight: 600;
  color: #374151;
}

.compact-dict-form :deep(.n-form-item-blank) {
  min-height: 40px;
}

.compact-dict-form :deep(.n-input-wrapper),
.compact-dict-form :deep(.n-base-selection) {
  border-radius: 10px;
  background: #fbfcfe;
  box-shadow: none;
}

.compact-dict-form {
  display: flex;
  flex-direction: column;
  gap: 18px;
}

.dict-modal-shell {
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
  min-height: 124px;
  padding: 26px 28px 22px;
  background:
    radial-gradient(circle at top right, rgba(34, 197, 94, 0.12), transparent 24%),
    linear-gradient(135deg, #eff6ff 0%, #e8f2ff 58%, #f4f9ff 100%);
}

.modal-header--item {
  background:
    radial-gradient(circle at top right, rgba(59, 130, 246, 0.12), transparent 24%),
    linear-gradient(135deg, #eff6ff 0%, #e8f2ff 58%, #f4f9ff 100%);
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
  transition:
    background-color 0.2s ease,
    color 0.2s ease,
    box-shadow 0.2s ease,
    transform 0.2s ease;
}

.modal-close:hover {
  background: #ffffff;
  color: #111827;
  box-shadow: 0 14px 28px rgba(148, 163, 184, 0.18);
  transform: translateY(-1px);
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

.form-section__head p {
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

@media (max-width: 1280px) {
  .dict-card-shell__filters {
    grid-template-columns: minmax(0, 1fr);
  }

  .dict-filter-actions {
    justify-content: flex-end;
  }
}

@media (max-width: 720px) {
  .form-section-grid {
    grid-template-columns: minmax(0, 1fr);
  }

  .compact-dict-modal :deep(.n-card__content),
  .compact-dict-modal :deep(.n-card__footer) {
    padding-left: 20px;
    padding-right: 20px;
  }

  .modal-header--hero {
    min-height: 112px;
    padding: 22px 20px 18px;
  }

  .modal-close {
    top: 18px;
    right: 18px;
  }

  .dict-card-shell__header,
  .dict-card-shell__filters,
  .dict-card-shell__footer {
    padding-left: 16px;
    padding-right: 16px;
  }

  .dict-card-shell__footer {
    flex-direction: column;
    align-items: stretch;
  }
}
</style>
