import type { DataTableColumns, FormRules } from 'naive-ui'
import { NButton, NPopconfirm, NSpace, NTag, useMessage } from 'naive-ui'
import { h, onMounted, reactive, ref } from 'vue'

import { useModalForm } from '@/composables/useModalForm'
import { usePermission } from '@/composables/usePermission'
import { useStatusToggle } from '@/composables/useStatusToggle'
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

export interface DictTypeFormModel {
  id: number
  code: string
  name: string
  sort: number
  status: DictStatus
  remark: string
}

export interface DictItemFormModel {
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

function toTagType(value: string) {
  if (value === 'success' || value === 'warning' || value === 'error' || value === 'info' || value === 'default') {
    return value
  }

  return 'default'
}

export function useDictPage() {
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

  const { handleToggleStatus: handleToggleTypeStatus } = useStatusToggle(updateDictTypeStatus, {
    onSuccess: async () => {
      successText.value = '字典类型状态已更新'
      await loadDictTypes()
    },
  })

  const { handleToggleStatus: handleToggleItemStatus } = useStatusToggle(updateDictItemStatus, {
    onSuccess: async () => {
      successText.value = '字典项状态已更新'
      await loadDictItems()
    },
  })

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

  async function submitType() {
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

  async function submitItem() {
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
