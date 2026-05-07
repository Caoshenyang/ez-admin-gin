<script setup lang="ts">
import type { DataTableColumns, FormInst, FormRules } from 'naive-ui'
import {
  NAlert,
  NButton,
  NCard,
  NDataTable,
  NEmpty,
  NForm,
  NFormItem,
  NInput,
  NInputNumber,
  NModal,
  NPagination,
  NPopconfirm,
  NSpace,
  NSelect,
  NTag,
  useMessage,
} from 'naive-ui'
import { h, onMounted, reactive, ref } from 'vue'

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
import { buttonPermissionCodes } from '@/router/dynamic-menu'
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
const typeLoading = ref(false)
const itemLoading = ref(false)
const typeSaving = ref(false)
const itemSaving = ref(false)

const dictTypes = ref<DictTypeItem[]>([])
const dictTypeTotal = ref(0)
const selectedTypeID = ref<number | null>(null)
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

const selectedType = ref<DictTypeItem | null>(null)

const typeFormRef = ref<FormInst | null>(null)
const typeFormVisible = ref(false)
const typeFormMode = ref<'create' | 'edit'>('create')
const typeFormModel = reactive<DictTypeFormModel>({
  id: 0,
  code: '',
  name: '',
  sort: 10,
  status: DictStatus.Enabled,
  remark: '',
})

const itemFormRef = ref<FormInst | null>(null)
const itemFormVisible = ref(false)
const itemFormMode = ref<'create' | 'edit'>('create')
const itemFormModel = reactive<DictItemFormModel>({
  id: 0,
  type_id: 0,
  item_key: '',
  label: '',
  value: '',
  tag_type: '',
  sort: 10,
  status: DictStatus.Enabled,
  remark: '',
})

const statusOptions = [
  { label: '状态：全部', value: 0 },
  { label: '启用', value: DictStatus.Enabled },
  { label: '禁用', value: DictStatus.Disabled },
]

const statusFormOptions = [
  { label: '启用', value: DictStatus.Enabled },
  { label: '禁用', value: DictStatus.Disabled },
]

const typeRules: FormRules = {
  code: [{ required: true, message: '请输入字典编码', trigger: 'blur' }],
  name: [{ required: true, message: '请输入字典名称', trigger: 'blur' }],
}

const itemRules: FormRules = {
  item_key: [{ required: true, message: '请输入字典项编码', trigger: 'blur' }],
  label: [{ required: true, message: '请输入字典项名称', trigger: 'blur' }],
  value: [{ required: true, message: '请输入字典项值', trigger: 'blur' }],
}

const typeColumns: DataTableColumns<DictTypeItem> = [
  {
    title: '编码',
    key: 'code',
    width: 180,
    ellipsis: { tooltip: true },
    render(row) {
      return h('span', { class: 'font-semibold text-[#111827]' }, row.code)
    },
  },
  {
    title: '名称',
    key: 'name',
    width: 140,
  },
  {
    title: '排序',
    key: 'sort',
    width: 80,
  },
  {
    title: '状态',
    key: 'status',
    width: 90,
    render(row) {
      return h(
        NTag,
        { type: row.status === DictStatus.Enabled ? 'success' : 'error', bordered: false },
        { default: () => (row.status === DictStatus.Enabled ? '启用' : '禁用') },
      )
    },
  },
  {
    title: '操作',
    key: 'actions',
    width: 180,
    fixed: 'right',
    render(row) {
      const nextStatus = row.status === DictStatus.Enabled ? DictStatus.Disabled : DictStatus.Enabled

      return h(
        NSpace,
        { size: 8 },
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
    title: '编码',
    key: 'item_key',
    width: 140,
  },
  {
    title: '名称',
    key: 'label',
    width: 140,
  },
  {
    title: '值',
    key: 'value',
    minWidth: 140,
    ellipsis: { tooltip: true },
  },
  {
    title: '标签样式',
    key: 'tag_type',
    width: 110,
    render(row) {
      if (!row.tag_type) {
        return h('span', { class: 'text-[#9CA3AF]' }, '-')
      }
      return h(NTag, { size: 'small', bordered: false, type: toTagType(row.tag_type) }, { default: () => row.tag_type })
    },
  },
  {
    title: '排序',
    key: 'sort',
    width: 80,
  },
  {
    title: '状态',
    key: 'status',
    width: 90,
    render(row) {
      return h(
        NTag,
        { type: row.status === DictStatus.Enabled ? 'success' : 'error', bordered: false },
        { default: () => (row.status === DictStatus.Enabled ? '启用' : '禁用') },
      )
    },
  },
  {
    title: '操作',
    key: 'actions',
    width: 180,
    fixed: 'right',
    render(row) {
      const nextStatus = row.status === DictStatus.Enabled ? DictStatus.Disabled : DictStatus.Enabled

      return h(
        NSpace,
        { size: 8 },
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

function canUse(code: string) {
  return buttonPermissionCodes.value.includes(code)
}

function toTagType(value: string) {
  if (value === 'success' || value === 'warning' || value === 'error' || value === 'info' || value === 'default') {
    return value
  }
  return 'default'
}

function formatTime(value: string) {
  if (!value) return '-'
  const d = new Date(value)
  const pad = (n: number) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}`
}

function typeRowProps(row: DictTypeItem) {
  return {
    class: row.id === selectedTypeID.value ? 'dict-type-row dict-type-row--active' : 'dict-type-row',
    onClick: () => selectType(row),
  }
}

function resetTypeForm() {
  Object.assign(typeFormModel, {
    id: 0,
    code: '',
    name: '',
    sort: 10,
    status: DictStatus.Enabled,
    remark: '',
  })
}

function resetItemForm() {
  Object.assign(itemFormModel, {
    id: 0,
    type_id: selectedTypeID.value ?? 0,
    item_key: '',
    label: '',
    value: '',
    tag_type: '',
    sort: 10,
    status: DictStatus.Enabled,
    remark: '',
  })
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

    if (data.items.length === 0) {
      selectedTypeID.value = null
      selectedType.value = null
      itemQuery.type_id = 0
      dictItems.value = []
      dictItemTotal.value = 0
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
  if (selectedTypeID.value === row.id && selectedType.value?.updated_at === row.updated_at) {
    return
  }

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

function handleTypePageChange(page: number) {
  typeQuery.page = page
  void loadDictTypes()
}

function handleTypePageSizeChange(pageSize: number) {
  typeQuery.page = 1
  typeQuery.page_size = pageSize
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
  typeFormMode.value = 'create'
  resetTypeForm()
  typeFormVisible.value = true
}

function openTypeEdit(row: DictTypeItem) {
  typeFormMode.value = 'edit'
  Object.assign(typeFormModel, {
    id: row.id,
    code: row.code,
    name: row.name,
    sort: row.sort,
    status: row.status,
    remark: row.remark,
  })
  typeFormVisible.value = true
}

function openItemCreate() {
  if (!selectedTypeID.value) {
    message.warning('请先选择一个字典类型')
    return
  }

  itemFormMode.value = 'create'
  resetItemForm()
  itemFormVisible.value = true
}

function openItemEdit(row: DictItem) {
  itemFormMode.value = 'edit'
  Object.assign(itemFormModel, {
    id: row.id,
    type_id: row.type_id,
    item_key: row.item_key,
    label: row.label,
    value: row.value,
    tag_type: row.tag_type,
    sort: row.sort,
    status: row.status,
    remark: row.remark,
  })
  itemFormVisible.value = true
}

async function handleTypeSubmit() {
  await typeFormRef.value?.validate()
  typeSaving.value = true
  try {
    if (typeFormMode.value === 'create') {
      const created = await createDictType({
        code: typeFormModel.code,
        name: typeFormModel.name,
        sort: typeFormModel.sort,
        status: typeFormModel.status,
        remark: typeFormModel.remark,
      })
      message.success(`已创建字典类型：${created.name}`)
    } else {
      const updated = await updateDictType(typeFormModel.id, {
        name: typeFormModel.name,
        sort: typeFormModel.sort,
        status: typeFormModel.status,
        remark: typeFormModel.remark,
      })
      message.success(`已更新字典类型：${updated.name}`)
    }

    typeFormVisible.value = false
    await loadDictTypes()
  } finally {
    typeSaving.value = false
  }
}

async function handleItemSubmit() {
  if (!selectedTypeID.value) {
    message.warning('请先选择一个字典类型')
    return
  }

  await itemFormRef.value?.validate()
  itemSaving.value = true
  try {
    if (itemFormMode.value === 'create') {
      const created = await createDictItem({
        type_id: selectedTypeID.value,
        item_key: itemFormModel.item_key,
        label: itemFormModel.label,
        value: itemFormModel.value,
        tag_type: itemFormModel.tag_type,
        sort: itemFormModel.sort,
        status: itemFormModel.status,
        remark: itemFormModel.remark,
      })
      message.success(`已创建字典项：${created.label}`)
    } else {
      const updated = await updateDictItem(itemFormModel.id, {
        label: itemFormModel.label,
        value: itemFormModel.value,
        tag_type: itemFormModel.tag_type,
        sort: itemFormModel.sort,
        status: itemFormModel.status,
        remark: itemFormModel.remark,
      })
      message.success(`已更新字典项：${updated.label}`)
    }

    itemFormVisible.value = false
    await loadDictItems()
  } finally {
    itemSaving.value = false
  }
}

async function handleToggleTypeStatus(row: DictTypeItem, status: DictStatus) {
  await updateDictTypeStatus(row.id, { status })
  message.success(`已${status === DictStatus.Disabled ? '禁用' : '启用'}字典类型：${row.name}`)
  await loadDictTypes()
}

async function handleToggleItemStatus(row: DictItem, status: DictStatus) {
  await updateDictItemStatus(row.id, { status })
  message.success(`已${status === DictStatus.Disabled ? '禁用' : '启用'}字典项：${row.label}`)
  await loadDictItems()
}

onMounted(async () => {
  await loadDictTypes()
})
</script>

<template>
  <div class="space-y-5">
    <NAlert type="info" :show-icon="false">
      数据字典当前按“系统级公共字典”处理，不进入部门或用户数据范围裁剪。
    </NAlert>

    <NCard title="字典类型" size="small" :bordered="false" class="shadow-sm">
      <template #header-extra>
        <NSpace align="center" size="small">
          <NInput
            v-model:value="typeQuery.keyword"
            clearable
            placeholder="搜索编码或名称"
            class="w-[220px]"
            @keydown.enter.prevent="handleTypeSearch"
          />
          <NSelect
            v-model:value="typeQuery.status"
            :options="statusOptions"
            class="w-[150px]"
          />
          <NButton type="primary" @click="handleTypeSearch">筛选</NButton>
          <NButton tertiary @click="handleTypeReset">重置</NButton>
          <NButton v-if="canUse('system:dict:type:create')" type="primary" @click="openTypeCreate">
            新建类型
          </NButton>
        </NSpace>
      </template>

      <NDataTable
        :columns="typeColumns"
        :data="dictTypes"
        :loading="typeLoading"
        :row-props="typeRowProps"
        :bordered="false"
        size="small"
        max-height="320"
      />

      <div class="mt-4 flex justify-end">
        <NPagination
          :page="typeQuery.page"
          :page-size="typeQuery.page_size"
          :item-count="dictTypeTotal"
          show-size-picker
          :page-sizes="[10, 20, 50]"
          @update:page="handleTypePageChange"
          @update:page-size="handleTypePageSizeChange"
        />
      </div>
    </NCard>

    <NCard size="small" :bordered="false" class="shadow-sm">
      <template #header>
        <div class="flex items-center gap-3">
          <span>字典项</span>
          <template v-if="selectedType">
            <NTag type="info" :bordered="false">{{ selectedType.name }}</NTag>
            <span class="text-sm text-[#6B7280]">{{ selectedType.code }}</span>
          </template>
        </div>
      </template>

      <template #header-extra>
        <NSpace v-if="selectedType" align="center" size="small">
          <NInput
            v-model:value="itemQuery.keyword"
            clearable
            placeholder="搜索编码、名称或值"
            class="w-[220px]"
            @keydown.enter.prevent="handleItemSearch"
          />
          <NSelect
            v-model:value="itemQuery.status"
            :options="statusOptions"
            class="w-[150px]"
          />
          <NButton type="primary" @click="handleItemSearch">筛选</NButton>
          <NButton tertiary @click="handleItemReset">重置</NButton>
          <NButton v-if="canUse('system:dict:item:create')" type="primary" @click="openItemCreate">
            新建字典项
          </NButton>
        </NSpace>
      </template>

      <template v-if="selectedType">
        <NDataTable
          :columns="itemColumns"
          :data="dictItems"
          :loading="itemLoading"
          :bordered="false"
          size="small"
          max-height="360"
        />

        <div class="mt-4 flex items-center justify-between gap-3">
          <div class="text-sm text-[#6B7280]">
            选中类型更新时间：{{ formatTime(selectedType.updated_at) }}
          </div>
          <NPagination
            :page="itemQuery.page"
            :page-size="itemQuery.page_size"
            :item-count="dictItemTotal"
            show-size-picker
            :page-sizes="[10, 20, 50]"
            @update:page="handleItemPageChange"
            @update:page-size="handleItemPageSizeChange"
          />
        </div>
      </template>
      <template v-else>
        <NEmpty description="先从上面的字典类型列表里选择一项，再查看和维护字典项。" />
      </template>
    </NCard>

    <NModal
      v-model:show="typeFormVisible"
      preset="card"
      :title="typeFormMode === 'create' ? '新建字典类型' : '编辑字典类型'"
      class="w-[620px]"
      :bordered="false"
      size="small"
    >
      <NForm ref="typeFormRef" :model="typeFormModel" :rules="typeRules" label-placement="top">
        <div class="grid gap-4 md:grid-cols-2">
          <NFormItem label="字典编码" path="code">
            <NInput
              v-model:value="typeFormModel.code"
              :disabled="typeFormMode === 'edit'"
              placeholder="例如 common:yes-no"
            />
          </NFormItem>
          <NFormItem label="字典名称" path="name">
            <NInput v-model:value="typeFormModel.name" placeholder="请输入字典名称" />
          </NFormItem>
          <NFormItem label="排序">
            <NInputNumber v-model:value="typeFormModel.sort" class="w-full" />
          </NFormItem>
          <NFormItem label="状态">
            <NSelect v-model:value="typeFormModel.status" :options="statusFormOptions" />
          </NFormItem>
        </div>
        <NFormItem label="备注">
          <NInput v-model:value="typeFormModel.remark" type="textarea" :rows="3" placeholder="可选备注" />
        </NFormItem>
      </NForm>

      <template #action>
        <div class="flex justify-end gap-3">
          <NButton @click="typeFormVisible = false">取消</NButton>
          <NButton type="primary" :loading="typeSaving" @click="handleTypeSubmit">保存</NButton>
        </div>
      </template>
    </NModal>

    <NModal
      v-model:show="itemFormVisible"
      preset="card"
      :title="itemFormMode === 'create' ? '新建字典项' : '编辑字典项'"
      class="w-[720px]"
      :bordered="false"
      size="small"
    >
      <NForm ref="itemFormRef" :model="itemFormModel" :rules="itemRules" label-placement="top">
        <div class="grid gap-4 md:grid-cols-2">
          <NFormItem label="字典项编码" path="item_key">
            <NInput
              v-model:value="itemFormModel.item_key"
              :disabled="itemFormMode === 'edit'"
              placeholder="例如 yes"
            />
          </NFormItem>
          <NFormItem label="字典项名称" path="label">
            <NInput v-model:value="itemFormModel.label" placeholder="请输入字典项名称" />
          </NFormItem>
          <NFormItem label="字典项值" path="value">
            <NInput v-model:value="itemFormModel.value" placeholder="例如 1 或 info" />
          </NFormItem>
          <NFormItem label="标签样式">
            <NInput v-model:value="itemFormModel.tag_type" placeholder="例如 success / warning / info" />
          </NFormItem>
          <NFormItem label="排序">
            <NInputNumber v-model:value="itemFormModel.sort" class="w-full" />
          </NFormItem>
          <NFormItem label="状态">
            <NSelect v-model:value="itemFormModel.status" :options="statusFormOptions" />
          </NFormItem>
        </div>
        <NFormItem label="备注">
          <NInput v-model:value="itemFormModel.remark" type="textarea" :rows="3" placeholder="可选备注" />
        </NFormItem>
      </NForm>

      <template #action>
        <div class="flex justify-end gap-3">
          <NButton @click="itemFormVisible = false">取消</NButton>
          <NButton type="primary" :loading="itemSaving" @click="handleItemSubmit">保存</NButton>
        </div>
      </template>
    </NModal>
  </div>
</template>

<style scoped>
:deep(.dict-type-row) {
  cursor: pointer;
}

:deep(.dict-type-row--active td) {
  background: #eff6ff;
}
</style>
