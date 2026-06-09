import type { DataTableColumns, FormRules } from 'naive-ui'
import { NPopconfirm, NSpace, NTag, useMessage } from 'naive-ui'
import { computed, h, ref } from 'vue'

import EzActionButton from '@/components/ez/EzActionButton.vue'
import { useModalForm } from '@/composables/useModalForm'
import { usePermission } from '@/composables/usePermission'
import { useRemotePagination } from '@/composables/useRemotePagination'
import { useStatusToggle } from '@/composables/useStatusToggle'
import { displayText } from '@/utils/format'
import {
  createConfig,
  deleteConfig,
  getConfigs,
  updateConfig,
  updateConfigStatus,
} from '../api/config'
import { ConfigStatus, type ConfigItem, type ConfigListQuery } from '../types/config'
import {
  BUILTIN_CONFIG_CATEGORIES,
  buildConfigCreatePayload,
  buildConfigUpdatePayload,
  defaultConfigFormModel,
  defaultConfigListQuery,
  toConfigFormModel,
  type ConfigFormModel,
} from './config-page.utils'

const defaultConfigCategory = {
  key: 'all',
  group_code: '',
  label: '全部配置',
  description: '查看所有系统配置项',
}

export function useConfigPage() {
  const message = useMessage()
  const { canUse } = usePermission()
  const activeConfigGroup = ref('')

  const {
    items: configs,
    total,
    loading,
    query,
    load,
    handleSearch: searchConfigs,
    handleReset: resetConfigs,
    handlePageChange,
    handlePageSizeChange,
  } = useRemotePagination<ConfigItem, ConfigListQuery>(getConfigs, {
    ...defaultConfigListQuery(),
  })

  const {
    formRef,
    formVisible,
    formMode,
    formModel,
    saving,
    rules,
    openCreate: openCreateModal,
    openEdit,
    handleSubmit,
  } = useModalForm<ConfigFormModel>(defaultConfigFormModel, {
    rules: {
      group_code: [{ required: true, message: '请输入配置分组', trigger: 'blur' }],
      key: [{ required: true, message: '请输入配置键', trigger: 'blur' }],
      name: [{ required: true, message: '请输入配置名称', trigger: 'blur' }],
      value: [{ required: true, message: '请输入配置值', trigger: 'blur' }],
    } as FormRules,
  })

  const configCategories = computed(() => {
    const categoryMap = new Map(BUILTIN_CONFIG_CATEGORIES.map((item) => [item.group_code, item]))

    for (const config of configs.value) {
      if (!config.group_code || categoryMap.has(config.group_code)) {
        continue
      }

      categoryMap.set(config.group_code, {
        key: config.group_code,
        group_code: config.group_code,
        label: `${config.group_code} 分组`,
        description: '当前已有配置分组',
      })
    }

    return Array.from(categoryMap.values())
  })

  const activeConfigCategory = computed(() => {
    return (
      configCategories.value.find((item) => item.group_code === activeConfigGroup.value) ??
      BUILTIN_CONFIG_CATEGORIES[0] ??
      defaultConfigCategory
    )
  })

  const { handleToggleStatus } = useStatusToggle(updateConfigStatus, {
    onSuccess: async () => {
      await load()
    },
  })

  const handleSearch = () => {
    activeConfigGroup.value = query.group_code?.trim() ?? ''
    searchConfigs()
  }

  const handleReset = () => {
    activeConfigGroup.value = ''
    resetConfigs()
  }

  const selectConfigCategory = (groupCode: string) => {
    activeConfigGroup.value = groupCode
    query.group_code = groupCode
    query.page = 1
    void load()
  }

  const openCreateForm = (groupCode = activeConfigGroup.value) => {
    openCreateModal()

    if (groupCode) {
      formModel.group_code = groupCode
      formModel.sort = configs.value.length * 10 + 10
    }
  }

  const openEditForm = (row: ConfigItem) => {
    openEdit(toConfigFormModel(row))
  }

  const columns: DataTableColumns<ConfigItem> = [
    {
      title: '分组',
      key: 'group_code',
      width: 112,
      render(row) {
        return h(
          NTag,
          { size: 'small', bordered: false, type: 'info' },
          { default: () => displayText(row.group_code) },
        )
      },
    },
    {
      title: '键',
      key: 'key',
      width: 168,
      ellipsis: { tooltip: true },
      render(row) {
        return displayText(row.key)
      },
    },
    {
      title: '名称',
      key: 'name',
      width: 136,
      render(row) {
        return displayText(row.name)
      },
    },
    {
      title: '值',
      key: 'value',
      minWidth: 160,
      ellipsis: { tooltip: true },
      render(row) {
        return displayText(row.value)
      },
    },
    {
      title: '排序',
      key: 'sort',
      width: 66,
    },
    {
      title: '状态',
      key: 'status',
      width: 78,
      render(row) {
        return h(
          NTag,
          { type: row.status === ConfigStatus.Enabled ? 'success' : 'error', bordered: false },
          { default: () => (row.status === ConfigStatus.Enabled ? '启用' : '禁用') },
        )
      },
    },
    {
      title: '操作',
      key: 'actions',
      width: 128,
      fixed: 'right',
      render(row) {
        const nextStatus =
          row.status === ConfigStatus.Enabled ? ConfigStatus.Disabled : ConfigStatus.Enabled

        return h(
          NSpace,
          { class: 'ez-row-actions', size: 6, align: 'center' },
          {
            default: () =>
              [
                canUse('system:config:update')
                  ? h(EzActionButton, {
                      iconOnly: true,
                      kind: 'edit',
                      label: '编辑',
                      size: 'tiny',
                      secondary: true,
                      type: 'info',
                      onClick: () => openEditForm(row),
                    })
                  : null,
                canUse('system:config:status')
                  ? h(
                      NPopconfirm,
                      { onPositiveClick: () => handleToggleStatus(row, nextStatus) },
                      {
                        trigger: () =>
                          h(EzActionButton, {
                            iconOnly: true,
                            kind: nextStatus === ConfigStatus.Disabled ? 'disable' : 'enable',
                            label: nextStatus === ConfigStatus.Disabled ? '禁用' : '启用',
                            size: 'tiny',
                            secondary: true,
                            tooltip: false,
                            type: nextStatus === ConfigStatus.Disabled ? 'error' : 'success',
                          }),
                        default: () =>
                          `确认${nextStatus === ConfigStatus.Disabled ? '禁用' : '启用'}该配置？`,
                      },
                    )
                  : null,
                canUse('system:config:delete')
                  ? h(
                      NPopconfirm,
                      { onPositiveClick: () => handleDelete(row) },
                      {
                        trigger: () =>
                          h(EzActionButton, {
                            iconOnly: true,
                            kind: 'delete',
                            label: '删除',
                            size: 'tiny',
                            secondary: true,
                            tooltip: false,
                            type: 'error',
                          }),
                        default: () => '删除后该配置会从列表和缓存中移除，确认继续？',
                      },
                    )
                  : null,
              ].filter(Boolean),
          },
        )
      },
    },
  ]

  async function submitForm() {
    if (formMode.value === 'create') {
      await createConfig(buildConfigCreatePayload(formModel))
      message.success('配置创建成功')
    } else {
      await updateConfig(formModel.id, buildConfigUpdatePayload(formModel))
      message.success('配置更新成功')
    }
    await load()
  }

  async function handleDelete(row: ConfigItem) {
    await deleteConfig(row.id)
    message.success('配置已删除')
    await load()
  }

  return {
    activeConfigCategory,
    activeConfigGroup,
    canUse,
    configCategories,
    columns,
    configs,
    formMode,
    formModel,
    formRef,
    formVisible,
    handlePageChange,
    handlePageSizeChange,
    handleReset,
    handleSearch,
    handleSubmit,
    load,
    loading,
    openCreate: openCreateForm,
    openEdit: openEditForm,
    query,
    rules,
    saving,
    selectConfigCategory,
    submitForm,
    total,
  }
}
