import type { DataTableColumns, FormRules } from 'naive-ui'
import { NButton, NPopconfirm, NSpace, NTag, useMessage } from 'naive-ui'
import { h } from 'vue'

import { useModalForm } from '@/composables/useModalForm'
import { usePermission } from '@/composables/usePermission'
import { useRemotePagination } from '@/composables/useRemotePagination'
import { useStatusToggle } from '@/composables/useStatusToggle'
import { displayText } from '@/utils/format'
import { createConfig, deleteConfig, getConfigs, updateConfig, updateConfigStatus } from '../api/config'
import { ConfigStatus, type ConfigItem, type ConfigListQuery } from '../types/config'
import {
  buildConfigCreatePayload,
  buildConfigUpdatePayload,
  defaultConfigFormModel,
  defaultConfigListQuery,
  toConfigFormModel,
  type ConfigFormModel,
} from './config-page.utils'

export function useConfigPage() {
  const message = useMessage()
  const { canUse } = usePermission()

  const {
    items: configs,
    total,
    loading,
    query,
    load,
    handleSearch,
    handleReset,
    handlePageChange,
    handlePageSizeChange,
  } = useRemotePagination<ConfigItem, ConfigListQuery>(getConfigs, {
    ...defaultConfigListQuery(),
  })

  const { formRef, formVisible, formMode, formModel, saving, rules, openCreate, openEdit, handleSubmit } =
    useModalForm<ConfigFormModel>(defaultConfigFormModel, {
      rules: {
        group_code: [{ required: true, message: '请输入配置分组', trigger: 'blur' }],
        key: [{ required: true, message: '请输入配置键', trigger: 'blur' }],
        name: [{ required: true, message: '请输入配置名称', trigger: 'blur' }],
        value: [{ required: true, message: '请输入配置值', trigger: 'blur' }],
      } as FormRules,
    })

  const { handleToggleStatus } = useStatusToggle(updateConfigStatus, {
    onSuccess: async () => {
      await load()
    },
  })

  const openEditForm = (row: ConfigItem) => {
    openEdit(toConfigFormModel(row))
  }

  const columns: DataTableColumns<ConfigItem> = [
    {
      title: '分组',
      key: 'group_code',
      width: 112,
      render(row) {
        return h(NTag, { size: 'small', bordered: false, type: 'info' }, { default: () => displayText(row.group_code) })
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
      width: 160,
      fixed: 'right',
      render(row) {
        const nextStatus = row.status === ConfigStatus.Enabled ? ConfigStatus.Disabled : ConfigStatus.Enabled

        return h(
          NSpace,
          { size: 8, align: 'center' },
          {
            default: () =>
              [
                canUse('system:config:update')
                  ? h(
                      NButton,
                      { size: 'small', ghost: true, type: 'info', onClick: () => openEditForm(row) },
                      { default: () => '编辑' },
                    )
                  : null,
                canUse('system:config:status')
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
                              type: nextStatus === ConfigStatus.Disabled ? 'error' : 'success',
                            },
                            { default: () => (nextStatus === ConfigStatus.Disabled ? '禁用' : '启用') },
                          ),
                        default: () => `确认${nextStatus === ConfigStatus.Disabled ? '禁用' : '启用'}该配置？`,
                      },
                    )
                  : null,
                canUse('system:config:delete')
                  ? h(
                      NPopconfirm,
                      { onPositiveClick: () => handleDelete(row) },
                      {
                        trigger: () =>
                          h(
                            NButton,
                            { size: 'small', ghost: true, type: 'error' },
                            { default: () => '删除' },
                          ),
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
    canUse,
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
    openCreate,
    openEdit: openEditForm,
    query,
    rules,
    saving,
    submitForm,
    total,
  }
}
