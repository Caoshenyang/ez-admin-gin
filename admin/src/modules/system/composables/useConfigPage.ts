import type { DataTableColumns, FormRules } from 'naive-ui'
import { NButton, NPopconfirm, NSpace, NTag, useMessage } from 'naive-ui'
import { h, ref } from 'vue'

import { useModalForm } from '@/composables/useModalForm'
import { usePermission } from '@/composables/usePermission'
import { useRemotePagination } from '@/composables/useRemotePagination'
import { useStatusToggle } from '@/composables/useStatusToggle'
import { createConfig, getConfigs, updateConfig, updateConfigStatus } from '../api/config'
import { ConfigStatus, type ConfigItem, type ConfigListQuery } from '../types/config'

export interface ConfigFormModel {
  id: number
  group_code: string
  key: string
  name: string
  value: string
  sort: number
  status: ConfigStatus
  remark: string
}

function defaultFormModel(): ConfigFormModel {
  return { id: 0, group_code: '', key: '', name: '', value: '', sort: 0, status: ConfigStatus.Enabled, remark: '' }
}

export function useConfigPage() {
  const message = useMessage()
  const { canUse } = usePermission()
  const successText = ref('')

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
    page: 1,
    page_size: 10,
    keyword: '',
    group_code: '',
    status: 0,
  })

  const { formRef, formVisible, formMode, formModel, saving, rules, openCreate, openEdit, handleSubmit } =
    useModalForm<ConfigFormModel>(defaultFormModel, {
      rules: {
        group_code: [{ required: true, message: '请输入配置分组', trigger: 'blur' }],
        key: [{ required: true, message: '请输入配置键', trigger: 'blur' }],
        name: [{ required: true, message: '请输入配置名称', trigger: 'blur' }],
        value: [{ required: true, message: '请输入配置值', trigger: 'blur' }],
      } as FormRules,
    })

  const { handleToggleStatus } = useStatusToggle(updateConfigStatus, {
    onSuccess: async () => {
      successText.value = '配置状态已更新'
      await load()
    },
  })

  const columns: DataTableColumns<ConfigItem> = [
    {
      title: '分组',
      key: 'group_code',
      width: 140,
      render(row) {
        return h(NTag, { size: 'small', bordered: false, type: 'info' }, { default: () => row.group_code })
      },
    },
    {
      title: '键',
      key: 'key',
      width: 200,
      ellipsis: { tooltip: true },
    },
    {
      title: '名称',
      key: 'name',
      width: 160,
    },
    {
      title: '值',
      key: 'value',
      minWidth: 180,
      ellipsis: { tooltip: true },
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
                      { size: 'small', ghost: true, type: 'info', onClick: () => openEdit(row) },
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

  async function submitForm() {
    const payload = {
      group_code: formModel.group_code,
      name: formModel.name,
      value: formModel.value,
      sort: formModel.sort,
      status: formModel.status,
      remark: formModel.remark,
    }
    if (formMode.value === 'create') {
      await createConfig({ ...payload, key: formModel.key })
      successText.value = '配置创建成功'
      message.success('配置创建成功')
    } else {
      await updateConfig(formModel.id, payload)
      successText.value = '配置已更新'
      message.success('配置更新成功')
    }
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
    query,
    rules,
    saving,
    submitForm,
    successText,
    total,
  }
}
