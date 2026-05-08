import type { DataTableColumns, FormRules } from 'naive-ui'
import { NButton, NPopconfirm, NSpace, NTag, useMessage } from 'naive-ui'
import { h } from 'vue'

import { useModalForm } from '@/composables/useModalForm'
import { usePermission } from '@/composables/usePermission'
import { useRemotePagination } from '@/composables/useRemotePagination'
import { useSuccessFeedback } from '@/composables/useSuccessFeedback'
import { useStatusToggle } from '@/composables/useStatusToggle'
import { displayText } from '@/utils/format'
import { createConfig, getConfigs, updateConfig, updateConfigStatus } from '../api/config'
import { ConfigStatus, type ConfigItem, type ConfigListQuery } from '../types/config'
import {
  buildConfigCreatePayload,
  buildConfigUpdatePayload,
  defaultConfigFormModel,
  defaultConfigListQuery,
  toConfigFormModel,
  type ConfigFormModel,
} from './config-page.utils'

// 系统配置管理页面组合式函数，封装配置列表、创建、编辑、状态切换等逻辑
export function useConfigPage() {
  const message = useMessage()
  const { canUse } = usePermission()
  const { closeSuccess, showSuccess, successText } = useSuccessFeedback()

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
      showSuccess('配置状态已更新')
      await load()
    },
  })

  const openEditForm = (row: ConfigItem) => {
    openEdit(toConfigFormModel(row))
  }

  // 配置列表表格列定义
  const columns: DataTableColumns<ConfigItem> = [
    {
      title: '分组',
      key: 'group_code',
      width: 140,
      render(row) {
        return h(NTag, { size: 'small', bordered: false, type: 'info' }, { default: () => displayText(row.group_code) })
      },
    },
    {
      title: '键',
      key: 'key',
      width: 200,
      ellipsis: { tooltip: true },
      render(row) {
        return displayText(row.key)
      },
    },
    {
      title: '名称',
      key: 'name',
      width: 160,
      render(row) {
        return displayText(row.name)
      },
    },
    {
      title: '值',
      key: 'value',
      minWidth: 180,
      ellipsis: { tooltip: true },
      render(row) {
        return displayText(row.value)
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
          { type: row.status === ConfigStatus.Enabled ? 'success' : 'error', bordered: false },
          { default: () => (row.status === ConfigStatus.Enabled ? '启用' : '禁用') },
        )
      },
    },
    {
      title: '操作',
      key: 'actions',
      width: 180,
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
              ].filter(Boolean),
          },
        )
      },
    },
  ]

  // 提交配置表单（新建或更新）
  async function submitForm() {
    if (formMode.value === 'create') {
      await createConfig(buildConfigCreatePayload(formModel))
      showSuccess('配置创建成功')
      message.success('配置创建成功')
    } else {
      await updateConfig(formModel.id, buildConfigUpdatePayload(formModel))
      showSuccess('配置已更新')
      message.success('配置更新成功')
    }
    await load()
  }

  return {
    canUse,
    closeSuccess,
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
    successText,
    total,
  }
}
