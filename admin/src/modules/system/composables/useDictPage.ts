import type { DataTableColumns, FormRules } from 'naive-ui'
import { NButton, NPopconfirm, NSpace, NTag, useMessage } from 'naive-ui'
import { h, onMounted, reactive, ref } from 'vue'

import { useModalForm } from '@/composables/useModalForm'
import { usePermission } from '@/composables/usePermission'
import { useSuccessFeedback } from '@/composables/useSuccessFeedback'
import { useStatusToggle } from '@/composables/useStatusToggle'
import { displayText } from '@/utils/format'
import {
  buildDictItemCreatePayload,
  buildDictItemUpdatePayload,
  buildDictTypeCreatePayload,
  buildDictTypeUpdatePayload,
  defaultDictItemFormModel,
  defaultDictItemQuery,
  defaultDictTypeFormModel,
  defaultDictTypeQuery,
  normalizeDictItemQuery,
  normalizeDictTypeQuery,
  toDictItemFormModel,
  toDictTypeFormModel,
} from './dict-page.utils'
import {
  createDictItem as createDictItemRequest,
  createDictType as createDictTypeRequest,
  getDictItems as getDictItemsRequest,
  getDictTypes as getDictTypesRequest,
  updateDictItem as updateDictItemRequest,
  updateDictItemStatus,
  updateDictType as updateDictTypeRequest,
  updateDictTypeStatus,
} from '../api/dict'
import {
  DictStatus,
  type DictItem,
  type DictItemListQuery,
  type DictTypeItem,
  type DictTypeListQuery,
} from '../types/dict'
import type { DictItemFormModel, DictTypeFormModel } from '../types/dict-page'

function toTagType(value: string) {
  if (value === 'success' || value === 'warning' || value === 'error' || value === 'info' || value === 'default') {
    return value
  }

  return 'default'
}

export function useDictPage() {
  const message = useMessage()
  const { canUse } = usePermission()
  const { closeSuccess, showSuccess, successText } = useSuccessFeedback()

  const typeLoading = ref(false)
  const itemLoading = ref(false)

  const dictTypes = ref<DictTypeItem[]>([])
  const dictTypeTotal = ref(0)
  const selectedTypeID = ref<number | null>(null)
  const selectedType = ref<DictTypeItem | null>(null)

  const dictItems = ref<DictItem[]>([])
  const dictItemTotal = ref(0)

  const typeQuery = reactive<DictTypeListQuery>(defaultDictTypeQuery())

  const itemQuery = reactive<DictItemListQuery>(defaultDictItemQuery())

  function defaultItemFormModel(): DictItemFormModel {
    return defaultDictItemFormModel(selectedTypeID.value ?? 0)
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
  } = useModalForm<DictTypeFormModel>(defaultDictTypeFormModel, {
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

  async function loadDictItems() {
    if (!selectedTypeID.value) {
      dictItems.value = []
      dictItemTotal.value = 0
      return
    }

    itemLoading.value = true

    try {
      const data = await getDictItemsRequest(normalizeDictItemQuery(itemQuery, selectedTypeID.value))

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

  async function loadDictTypes() {
    typeLoading.value = true

    try {
      const data = await getDictTypesRequest(normalizeDictTypeQuery(typeQuery))

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

  const { handleToggleStatus: handleToggleTypeStatus } = useStatusToggle(updateDictTypeStatus, {
    onSuccess: async () => {
      showSuccess('字典类型状态已更新')
      await loadDictTypes()
    },
  })

  const { handleToggleStatus: handleToggleItemStatus } = useStatusToggle(updateDictItemStatus, {
    onSuccess: async () => {
      showSuccess('字典项状态已更新')
      await loadDictItems()
    },
  })

  function openTypeCreate() {
    openTypeCreateBase()
  }

  function openTypeEdit(row: DictTypeItem) {
    openTypeEditBase(toDictTypeFormModel(row))
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
    openItemEditBase(toDictItemFormModel(row))
  }

  async function submitType() {
    if (typeFormMode.value === 'create') {
      await createDictTypeRequest(buildDictTypeCreatePayload(typeFormModel))
      showSuccess('字典类型创建成功')
      message.success('字典类型创建成功')
    } else {
      await updateDictTypeRequest(typeFormModel.id, buildDictTypeUpdatePayload(typeFormModel))
      showSuccess('字典类型已更新')
      message.success('字典类型更新成功')
    }

    await loadDictTypes()
  }

  async function submitItem() {
    if (!selectedTypeID.value) {
      message.warning('请先选择一个字典类型')
      return
    }

    if (itemFormMode.value === 'create') {
      await createDictItemRequest(buildDictItemCreatePayload(selectedTypeID.value, itemFormModel))
      showSuccess('字典项创建成功')
      message.success('字典项创建成功')
    } else {
      await updateDictItemRequest(itemFormModel.id, buildDictItemUpdatePayload(itemFormModel))
      showSuccess('字典项已更新')
      message.success('字典项更新成功')
    }

    await loadDictItems()
  }

  function handleTypeSearch() {
    typeQuery.page = 1
    void loadDictTypes()
  }

  function handleTypeReset() {
    Object.assign(typeQuery, defaultDictTypeQuery())
    void loadDictTypes()
  }

  function handleItemSearch() {
    itemQuery.page = 1
    void loadDictItems()
  }

  function handleItemReset() {
    Object.assign(itemQuery, defaultDictItemQuery(selectedTypeID.value ?? 0))
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

  const typeColumns: DataTableColumns<DictTypeItem> = [
    {
      title: '字典类型',
      key: 'code',
      width: 168,
      render(row) {
        return h('div', { class: 'min-w-0 leading-5' }, [
          h('p', { class: 'truncate font-semibold text-[var(--ez-text-heading)]' }, displayText(row.name)),
          h('p', { class: 'truncate text-xs text-[var(--ez-text-muted)]' }, displayText(row.code)),
        ])
      },
    },
    {
      title: '排序',
      key: 'sort',
      width: 56,
      align: 'center',
    },
    {
      title: '状态',
      key: 'status',
      width: 76,
      align: 'center',
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
      width: 120,
      render(row) {
        const nextStatus = row.status === DictStatus.Enabled ? DictStatus.Disabled : DictStatus.Enabled

        return h(
          NSpace,
          { size: 6, align: 'center' },
          {
            default: () =>
              [
                canUse('system:dict:type:update')
                  ? h(
                      NButton,
                      { size: 'tiny', secondary: true, type: 'info', onClick: () => openTypeEdit(row) },
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
                              size: 'tiny',
                              secondary: true,
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
          h('p', { class: 'font-semibold text-[var(--ez-text-heading)]' }, displayText(row.label)),
          h('p', { class: 'text-xs text-[var(--ez-text-muted)]' }, `${displayText(row.item_key)} · ${displayText(row.value)}`),
        ])
      },
    },
    {
      title: '标签样式',
      key: 'tag_type',
      width: 120,
      render(row) {
        if (!row.tag_type) {
          return h('span', { class: 'text-[var(--ez-text-light)]' }, '-')
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

  function typeRowProps(row: DictTypeItem) {
    return {
      class: row.id === selectedTypeID.value ? 'dict-type-row dict-type-row--active' : 'dict-type-row',
      onClick: () => {
        void selectType(row)
      },
    }
  }

  onMounted(() => {
    void loadDictTypes()
  })

  return {
    canUse,
    closeSuccess,
    dictItemTotal,
    dictItems,
    dictTypeTotal,
    dictTypes,
    handleItemPageChange,
    handleItemPageSizeChange,
    handleItemReset,
    handleItemSearch,
    handleItemSubmit,
    handleTypePageChange,
    handleTypePageSizeChange,
    handleTypeReset,
    handleTypeRowProps: typeRowProps,
    handleTypeSearch,
    handleTypeSubmit,
    itemColumns,
    itemFormMode,
    itemFormModel,
    itemFormRef,
    itemFormVisible,
    itemLoading,
    itemQuery,
    itemRules,
    itemSaving,
    openItemCreate,
    openTypeCreate,
    selectedType,
    successText,
    submitItem,
    submitType,
    typeColumns,
    typeFormMode,
    typeFormModel,
    typeFormRef,
    typeFormVisible,
    typeLoading,
    typeQuery,
    typeRules,
    typeSaving,
  }
}
