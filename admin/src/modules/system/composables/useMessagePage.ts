import type { DataTableColumns, FormRules, SelectOption } from 'naive-ui'
import { NPopconfirm, NSpace, NTag, useMessage } from 'naive-ui'
import { computed, h, onMounted, ref } from 'vue'

import EzActionButton from '@/components/ez/EzActionButton.vue'
import { useModalForm } from '@/composables/useModalForm'
import { usePermission } from '@/composables/usePermission'
import { useRemotePagination } from '@/composables/useRemotePagination'
import { useStatusToggle } from '@/composables/useStatusToggle'
import { displayText, formatTime } from '@/utils/format'
import {
  createMessageReminder,
  createMessageTemplate,
  getMessageReminders,
  getMessageTemplates,
  updateMessageReminder,
  updateMessageReminderStatus,
  updateMessageTemplate,
  updateMessageTemplateStatus,
} from '../api/message'
import {
  MessageStatus,
  type MessageReminderItem,
  type MessageReminderListQuery,
  type MessageTemplateItem,
  type MessageTemplateListQuery,
} from '../types/message'
import {
  buildMessageReminderCreatePayload,
  buildMessageReminderUpdatePayload,
  buildMessageTemplateCreatePayload,
  buildMessageTemplateUpdatePayload,
  defaultMessageReminderFormModel,
  defaultMessageReminderListQuery,
  defaultMessageTemplateFormModel,
  defaultMessageTemplateListQuery,
  receiverTypeLabel,
  templateTypeLabel,
  toMessageReminderFormModel,
  toMessageTemplateFormModel,
} from './message-page.utils'

export function useMessagePage() {
  const message = useMessage()
  const { canUse } = usePermission()
  const activeTab = ref<'templates' | 'reminders'>('templates')
  const templateOptionItems = ref<MessageTemplateItem[]>([])

  const {
    items: templates,
    total: templateTotal,
    loading: templateLoading,
    query: templateQuery,
    load: loadTemplates,
    handleSearch: handleTemplateSearch,
    handleReset: handleTemplateReset,
    handlePageChange: handleTemplatePageChange,
    handlePageSizeChange: handleTemplatePageSizeChange,
  } = useRemotePagination<MessageTemplateItem, MessageTemplateListQuery>(getMessageTemplates, {
    ...defaultMessageTemplateListQuery(),
  })

  const {
    items: reminders,
    total: reminderTotal,
    loading: reminderLoading,
    query: reminderQuery,
    load: loadReminders,
    handleSearch: handleReminderSearch,
    handleReset: handleReminderReset,
    handlePageChange: handleReminderPageChange,
    handlePageSizeChange: handleReminderPageSizeChange,
  } = useRemotePagination<MessageReminderItem, MessageReminderListQuery>(getMessageReminders, {
    ...defaultMessageReminderListQuery(),
  })

  const {
    formRef: templateFormRef,
    formVisible: templateFormVisible,
    formMode: templateFormMode,
    formModel: templateFormModel,
    saving: templateSaving,
    rules: templateRules,
    openCreate: openTemplateCreateModal,
    openEdit: openTemplateEdit,
    handleSubmit: handleTemplateSubmit,
  } = useModalForm(defaultMessageTemplateFormModel, {
    rules: {
      code: [{ required: true, message: '请输入模板编码', trigger: 'blur' }],
      name: [{ required: true, message: '请输入模板名称', trigger: 'blur' }],
      title: [{ required: true, message: '请输入消息标题', trigger: 'blur' }],
      content: [{ required: true, message: '请输入消息内容', trigger: 'blur' }],
    } as FormRules,
  })

  const {
    formRef: reminderFormRef,
    formVisible: reminderFormVisible,
    formMode: reminderFormMode,
    formModel: reminderFormModel,
    saving: reminderSaving,
    rules: reminderRules,
    openCreate: openReminderCreateModal,
    openEdit: openReminderEdit,
    handleSubmit: handleReminderSubmit,
  } = useModalForm(defaultMessageReminderFormModel, {
    rules: {
      code: [{ required: true, message: '请输入提醒编码', trigger: 'blur' }],
      name: [{ required: true, message: '请输入提醒名称', trigger: 'blur' }],
      trigger_event: [{ required: true, message: '请输入触发事件', trigger: 'blur' }],
      template_id: [
        {
          required: true,
          type: 'number',
          message: '请选择消息模板',
          trigger: ['change', 'blur'],
        },
      ],
      channels: [{ required: true, message: '请输入提醒渠道', trigger: 'blur' }],
    } as FormRules,
  })

  const templateOptions = computed<SelectOption[]>(() =>
    templateOptionItems.value.map((item) => ({
      label: `${item.name}（${item.code}）`,
      value: item.id,
    })),
  )

  const { handleToggleStatus: handleToggleTemplateStatus } = useStatusToggle(
    updateMessageTemplateStatus,
    {
      onSuccess: async () => {
        await loadTemplates()
        await loadTemplateOptions()
      },
    },
  )

  const { handleToggleStatus: handleToggleReminderStatus } = useStatusToggle(
    updateMessageReminderStatus,
    {
      onSuccess: async () => {
        await loadReminders()
      },
    },
  )

  const templateColumns: DataTableColumns<MessageTemplateItem> = [
    {
      title: '模板',
      key: 'name',
      width: 210,
      ellipsis: { tooltip: true },
      render(row) {
        return h('div', { class: 'message-main-cell' }, [
          h('strong', displayText(row.name)),
          h('span', displayText(row.code)),
        ])
      },
    },
    {
      title: '标题',
      key: 'title',
      minWidth: 190,
      ellipsis: { tooltip: true },
      render(row) {
        return displayText(row.title)
      },
    },
    {
      title: '类型',
      key: 'type',
      width: 96,
      render(row) {
        return h(
          NTag,
          { size: 'small', bordered: false, type: templateTypeTagType(row.type) },
          { default: () => templateTypeLabel(row.type) },
        )
      },
    },
    {
      title: '系统级',
      key: 'is_system',
      width: 82,
      align: 'center',
      render(row) {
        return h(
          NTag,
          { size: 'small', bordered: false, type: row.is_system ? 'info' : 'default' },
          { default: () => (row.is_system ? '是' : '否') },
        )
      },
    },
    {
      title: '状态',
      key: 'status',
      width: 78,
      align: 'center',
      render(row) {
        return renderStatusTag(row.status)
      },
    },
    {
      title: '更新时间',
      key: 'updated_at',
      width: 150,
      render(row) {
        return h(
          'span',
          { class: 'tabular-nums text-[var(--ez-text-muted)]' },
          formatTime(row.updated_at),
        )
      },
    },
    {
      title: '操作',
      key: 'actions',
      width: 94,
      fixed: 'right',
      render(row) {
        const nextStatus =
          row.status === MessageStatus.Enabled ? MessageStatus.Disabled : MessageStatus.Enabled

        return h(
          NSpace,
          { class: 'ez-row-actions', size: 6, align: 'center' },
          {
            default: () =>
              [
                canUse('system:message:template:update')
                  ? h(EzActionButton, {
                      iconOnly: true,
                      kind: 'edit',
                      label: '编辑',
                      size: 'tiny',
                      secondary: true,
                      type: 'info',
                      onClick: () => openTemplateEdit(toMessageTemplateFormModel(row)),
                    })
                  : null,
                canUse('system:message:template:status')
                  ? h(
                      NPopconfirm,
                      { onPositiveClick: () => handleToggleTemplateStatus(row, nextStatus) },
                      {
                        trigger: () =>
                          h(EzActionButton, {
                            iconOnly: true,
                            kind: nextStatus === MessageStatus.Disabled ? 'disable' : 'enable',
                            label: nextStatus === MessageStatus.Disabled ? '禁用' : '启用',
                            size: 'tiny',
                            secondary: true,
                            tooltip: false,
                            type: nextStatus === MessageStatus.Disabled ? 'error' : 'success',
                          }),
                        default: () =>
                          `确认${nextStatus === MessageStatus.Disabled ? '禁用' : '启用'}该模板？`,
                      },
                    )
                  : null,
              ].filter(Boolean),
          },
        )
      },
    },
  ]

  const reminderColumns: DataTableColumns<MessageReminderItem> = [
    {
      title: '提醒规则',
      key: 'name',
      width: 220,
      ellipsis: { tooltip: true },
      render(row) {
        return h('div', { class: 'message-main-cell' }, [
          h('strong', displayText(row.name)),
          h('span', displayText(row.code)),
        ])
      },
    },
    {
      title: '触发事件',
      key: 'trigger_event',
      width: 160,
      ellipsis: { tooltip: true },
      render(row) {
        return displayText(row.trigger_event)
      },
    },
    {
      title: '模板',
      key: 'template_name',
      minWidth: 160,
      ellipsis: { tooltip: true },
      render(row) {
        return displayText(row.template_name || row.template_code)
      },
    },
    {
      title: '接收人',
      key: 'receiver_type',
      width: 110,
      render(row) {
        return h(
          NTag,
          { size: 'small', bordered: false, type: 'info' },
          { default: () => receiverTypeLabel(row.receiver_type) },
        )
      },
    },
    {
      title: '渠道',
      key: 'channels',
      width: 126,
      ellipsis: { tooltip: true },
      render(row) {
        return displayText(row.channels)
      },
    },
    {
      title: '提前',
      key: 'advance_minutes',
      width: 82,
      align: 'center',
      render(row) {
        return row.advance_minutes > 0 ? `${row.advance_minutes} 分钟` : '即时'
      },
    },
    {
      title: '状态',
      key: 'status',
      width: 78,
      align: 'center',
      render(row) {
        return renderStatusTag(row.status)
      },
    },
    {
      title: '操作',
      key: 'actions',
      width: 94,
      fixed: 'right',
      render(row) {
        const nextStatus =
          row.status === MessageStatus.Enabled ? MessageStatus.Disabled : MessageStatus.Enabled

        return h(
          NSpace,
          { class: 'ez-row-actions', size: 6, align: 'center' },
          {
            default: () =>
              [
                canUse('system:message:reminder:update')
                  ? h(EzActionButton, {
                      iconOnly: true,
                      kind: 'edit',
                      label: '编辑',
                      size: 'tiny',
                      secondary: true,
                      type: 'info',
                      onClick: () => openReminderEdit(toMessageReminderFormModel(row)),
                    })
                  : null,
                canUse('system:message:reminder:status')
                  ? h(
                      NPopconfirm,
                      { onPositiveClick: () => handleToggleReminderStatus(row, nextStatus) },
                      {
                        trigger: () =>
                          h(EzActionButton, {
                            iconOnly: true,
                            kind: nextStatus === MessageStatus.Disabled ? 'disable' : 'enable',
                            label: nextStatus === MessageStatus.Disabled ? '禁用' : '启用',
                            size: 'tiny',
                            secondary: true,
                            tooltip: false,
                            type: nextStatus === MessageStatus.Disabled ? 'error' : 'success',
                          }),
                        default: () =>
                          `确认${nextStatus === MessageStatus.Disabled ? '禁用' : '启用'}该提醒规则？`,
                      },
                    )
                  : null,
              ].filter(Boolean),
          },
        )
      },
    },
  ]

  async function submitTemplateForm() {
    if (templateFormMode.value === 'create') {
      await createMessageTemplate(buildMessageTemplateCreatePayload(templateFormModel))
      message.success('消息模板创建成功')
    } else {
      await updateMessageTemplate(
        templateFormModel.id,
        buildMessageTemplateUpdatePayload(templateFormModel),
      )
      message.success('消息模板更新成功')
    }
    await loadTemplates()
    await loadTemplateOptions()
  }

  async function submitReminderForm() {
    if (reminderFormMode.value === 'create') {
      await createMessageReminder(buildMessageReminderCreatePayload(reminderFormModel))
      message.success('提醒规则创建成功')
    } else {
      await updateMessageReminder(
        reminderFormModel.id,
        buildMessageReminderUpdatePayload(reminderFormModel),
      )
      message.success('提醒规则更新成功')
    }
    await loadReminders()
  }

  async function loadTemplateOptions() {
    const result = await getMessageTemplates({
      page: 1,
      page_size: 200,
      status: MessageStatus.Enabled,
    })
    templateOptionItems.value = result.items
  }

  function openCreateTemplate() {
    openTemplateCreateModal()
    templateFormModel.sort = templates.value.length * 10 + 10
  }

  function openCreateReminder() {
    openReminderCreateModal()
    reminderFormModel.sort = reminders.value.length * 10 + 10
    reminderFormModel.template_id = Number(templateOptions.value[0]?.value ?? 0) || null
  }

  onMounted(() => {
    void loadTemplateOptions()
  })

  return {
    activeTab,
    canUse,
    handleReminderPageChange,
    handleReminderPageSizeChange,
    handleReminderReset,
    handleReminderSearch,
    handleReminderSubmit,
    handleTemplatePageChange,
    handleTemplatePageSizeChange,
    handleTemplateReset,
    handleTemplateSearch,
    handleTemplateSubmit,
    loadReminders,
    loadTemplates,
    openCreateReminder,
    openCreateTemplate,
    reminderColumns,
    reminderFormMode,
    reminderFormModel,
    reminderFormRef,
    reminderFormVisible,
    reminderLoading,
    reminderQuery,
    reminderRules,
    reminders,
    reminderSaving,
    reminderTotal,
    submitReminderForm,
    submitTemplateForm,
    templateColumns,
    templateFormMode,
    templateFormModel,
    templateFormRef,
    templateFormVisible,
    templateLoading,
    templateOptions,
    templateQuery,
    templateRules,
    templates,
    templateSaving,
    templateTotal,
  }
}

function renderStatusTag(status: MessageStatus) {
  return h(
    NTag,
    { type: status === MessageStatus.Enabled ? 'success' : 'error', bordered: false },
    { default: () => (status === MessageStatus.Enabled ? '启用' : '禁用') },
  )
}

function templateTypeTagType(type: number) {
  if (type === 2) return 'warning'
  if (type === 3) return 'error'
  return 'info'
}
