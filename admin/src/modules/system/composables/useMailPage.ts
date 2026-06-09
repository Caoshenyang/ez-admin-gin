import type { DataTableColumns, FormRules, SelectOption } from 'naive-ui'
import { NPopconfirm, NSpace, NTag, useMessage } from 'naive-ui'
import { computed, h, reactive, ref, watch } from 'vue'

import EzActionButton from '@/components/ez/EzActionButton.vue'
import { useModalForm } from '@/composables/useModalForm'
import { usePermission } from '@/composables/usePermission'
import { useRemotePagination } from '@/composables/useRemotePagination'
import { useStatusToggle } from '@/composables/useStatusToggle'
import { displayText, formatTime } from '@/utils/format'
import {
  createMailAccount,
  createMailTemplate,
  deleteMailAccount,
  deleteMailTemplate,
  getMailAccounts,
  getMailLogs,
  getMailTemplates,
  renderMailTemplate,
  sendMail,
  testMailAccount,
  updateMailAccount,
  updateMailAccountStatus,
  updateMailTemplate,
  updateMailTemplateStatus,
} from '../api/mail'
import {
  MailEncryption,
  MailLogStatus,
  MailStatus,
  type MailAccountItem,
  type MailAccountListQuery,
  type MailLogItem,
  type MailLogListQuery,
  type MailTemplateItem,
  type MailTemplateListQuery,
} from '../types/mail'
import type {
  MailAccountFormModel,
  MailPanelKey,
  MailPreviewState,
  MailSendFormModel,
  MailTemplateFormModel,
  MailTestFormModel,
} from '../types/mail-page'
import {
  buildCreateMailTemplatePayload,
  buildMailAccountPayload,
  buildMailSendPayload,
  buildUpdateMailTemplatePayload,
  buildVariablesText,
  defaultMailAccountFormModel,
  defaultMailAccountListQuery,
  defaultMailLogListQuery,
  defaultMailSendFormModel,
  defaultMailTemplateFormModel,
  defaultMailTemplateListQuery,
  defaultMailTestFormModel,
  parseTextList,
  parseVariablesText,
  toMailAccountFormModel,
  toMailTemplateFormModel,
} from './mail-page.utils'

export function useMailPage() {
  const message = useMessage()
  const { canUse } = usePermission()
  const activePanel = ref<MailPanelKey>('accounts')

  const accountsPage = useRemotePagination<MailAccountItem, MailAccountListQuery>(
    getMailAccounts,
    defaultMailAccountListQuery(),
  )
  const templatesPage = useRemotePagination<MailTemplateItem, MailTemplateListQuery>(
    getMailTemplates,
    defaultMailTemplateListQuery(),
  )
  const logsPage = useRemotePagination<MailLogItem, MailLogListQuery>(
    getMailLogs,
    defaultMailLogListQuery(),
  )

  const accountForm = useModalForm<MailAccountFormModel>(defaultMailAccountFormModel, {
    rules: {
      name: [{ required: true, message: '请输入邮箱名称', trigger: 'blur' }],
      host: [{ required: true, message: '请输入 SMTP 主机', trigger: 'blur' }],
      from_email: [{ required: true, message: '请输入发件邮箱', trigger: 'blur' }],
    } as FormRules,
  })

  const templateForm = useModalForm<MailTemplateFormModel>(defaultMailTemplateFormModel, {
    rules: {
      code: [{ required: true, message: '请输入模板编码', trigger: 'blur' }],
      name: [{ required: true, message: '请输入模板名称', trigger: 'blur' }],
      subject: [{ required: true, message: '请输入邮件主题', trigger: 'blur' }],
      content: [{ required: true, message: '请输入邮件正文', trigger: 'blur' }],
    } as FormRules,
  })

  const sendForm = useModalForm<MailSendFormModel>(defaultMailSendFormModel, {
    rules: {
      to_text: [{ required: true, message: '请输入收件人', trigger: 'blur' }],
    } as FormRules,
  })

  const testForm = useModalForm<MailTestFormModel>(defaultMailTestFormModel, {
    rules: {
      to_text: [{ required: true, message: '请输入收件人', trigger: 'blur' }],
    } as FormRules,
  })

  const preview = reactive<MailPreviewState>({
    visible: false,
    title: '',
    subject: '',
    content: '',
  })

  const accountOptions = computed<SelectOption[]>(() =>
    accountsPage.items.value.map((account) => ({
      label: `${account.name} (${account.from_email})`,
      value: account.id,
    })),
  )

  const templateOptions = computed<SelectOption[]>(() =>
    templatesPage.items.value.map((template) => ({
      label: `${template.name} (${template.code})`,
      value: template.code,
    })),
  )

  const { handleToggleStatus: toggleAccountStatus } = useStatusToggle(updateMailAccountStatus, {
    onSuccess: async () => {
      await accountsPage.load()
    },
  })

  const { handleToggleStatus: toggleTemplateStatus } = useStatusToggle(updateMailTemplateStatus, {
    onSuccess: async () => {
      await templatesPage.load()
    },
  })

  const accountColumns: DataTableColumns<MailAccountItem> = [
    {
      title: '邮箱账号',
      key: 'name',
      minWidth: 180,
      render(row) {
        return h('div', { class: 'mail-cell-main' }, [
          h('span', { class: 'mail-cell-title' }, displayText(row.name)),
          h('span', { class: 'mail-cell-subtitle' }, displayText(row.from_email)),
        ])
      },
    },
    {
      title: 'SMTP',
      key: 'smtp',
      minWidth: 180,
      render(row) {
        return h('span', { class: 'tabular-nums' }, `${row.host}:${row.port}`)
      },
    },
    {
      title: '加密',
      key: 'encryption',
      width: 96,
      align: 'center',
      render(row) {
        return h(NTag, { bordered: false, type: 'info' }, { default: () => encryptionText(row.encryption) })
      },
    },
    {
      title: '默认',
      key: 'is_default',
      width: 78,
      align: 'center',
      render(row) {
        return row.is_default
          ? h(NTag, { bordered: false, type: 'success' }, { default: () => '默认' })
          : h('span', { class: 'text-[var(--ez-text-muted)]' }, '-')
      },
    },
    {
      title: '状态',
      key: 'status',
      width: 78,
      align: 'center',
      render(row) {
        return h(
          NTag,
          { type: row.status === MailStatus.Enabled ? 'success' : 'error', bordered: false },
          { default: () => (row.status === MailStatus.Enabled ? '启用' : '禁用') },
        )
      },
    },
    {
      title: '最近测试',
      key: 'last_test_at',
      width: 156,
      render(row) {
        return h('span', { class: 'text-[var(--ez-text-muted)]' }, formatTime(row.last_test_at ?? ''))
      },
    },
    {
      title: '操作',
      key: 'actions',
      width: 178,
      fixed: 'right',
      render(row) {
        const nextStatus =
          row.status === MailStatus.Enabled ? MailStatus.Disabled : MailStatus.Enabled
        return h(NSpace, { class: 'ez-row-actions', size: 6, align: 'center' }, {
          default: () =>
            [
              canUse('system:mail:account:update')
                ? h(EzActionButton, {
                    iconOnly: true,
                    kind: 'edit',
                    label: '编辑',
                    size: 'tiny',
                    secondary: true,
                    type: 'info',
                    onClick: () => openEditAccount(row),
                  })
                : null,
              canUse('system:mail:account:test')
                ? h(EzActionButton, {
                    iconOnly: true,
                    kind: 'view',
                    label: '测试发送',
                    size: 'tiny',
                    secondary: true,
                    type: 'primary',
                    onClick: () => openTestAccount(row),
                  })
                : null,
              canUse('system:mail:account:status')
                ? h(
                    NPopconfirm,
                    { onPositiveClick: () => toggleAccountStatus(row, nextStatus) },
                    {
                      trigger: () =>
                        h(EzActionButton, {
                          iconOnly: true,
                          kind: nextStatus === MailStatus.Disabled ? 'disable' : 'enable',
                          label: nextStatus === MailStatus.Disabled ? '禁用' : '启用',
                          size: 'tiny',
                          secondary: true,
                          tooltip: false,
                          type: nextStatus === MailStatus.Disabled ? 'error' : 'success',
                        }),
                      default: () =>
                        `确认${nextStatus === MailStatus.Disabled ? '禁用' : '启用'}该邮箱账号？`,
                    },
                  )
                : null,
              canUse('system:mail:account:delete')
                ? h(
                    NPopconfirm,
                    { onPositiveClick: () => handleDeleteAccount(row) },
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
                      default: () => '确认删除该邮箱账号？',
                    },
                  )
                : null,
            ].filter(Boolean),
        })
      },
    },
  ]

  const templateColumns: DataTableColumns<MailTemplateItem> = [
    {
      title: '模板',
      key: 'name',
      minWidth: 200,
      render(row) {
        return h('div', { class: 'mail-cell-main' }, [
          h('span', { class: 'mail-cell-title' }, displayText(row.name)),
          h('span', { class: 'mail-cell-subtitle' }, row.code),
        ])
      },
    },
    {
      title: '主题',
      key: 'subject',
      minWidth: 220,
      ellipsis: { tooltip: true },
    },
    {
      title: '变量',
      key: 'variables',
      minWidth: 160,
      render(row) {
        return row.variables.length > 0
          ? h(
              NSpace,
              { size: 4 },
              { default: () => row.variables.map((item) => h(NTag, { size: 'small', bordered: false }, { default: () => item })) },
            )
          : h('span', { class: 'text-[var(--ez-text-muted)]' }, '-')
      },
    },
    {
      title: '类型',
      key: 'is_html',
      width: 82,
      align: 'center',
      render(row) {
        return h(NTag, { bordered: false, type: row.is_html ? 'info' : 'default' }, { default: () => (row.is_html ? 'HTML' : '文本') })
      },
    },
    {
      title: '状态',
      key: 'status',
      width: 78,
      align: 'center',
      render(row) {
        return h(
          NTag,
          { type: row.status === MailStatus.Enabled ? 'success' : 'error', bordered: false },
          { default: () => (row.status === MailStatus.Enabled ? '启用' : '禁用') },
        )
      },
    },
    {
      title: '操作',
      key: 'actions',
      width: 178,
      fixed: 'right',
      render(row) {
        const nextStatus =
          row.status === MailStatus.Enabled ? MailStatus.Disabled : MailStatus.Enabled
        return h(NSpace, { class: 'ez-row-actions', size: 6, align: 'center' }, {
          default: () =>
            [
              canUse('system:mail:template:update')
                ? h(EzActionButton, {
                    iconOnly: true,
                    kind: 'edit',
                    label: '编辑',
                    size: 'tiny',
                    secondary: true,
                    type: 'info',
                    onClick: () => openEditTemplate(row),
                  })
                : null,
              canUse('system:mail:template:render')
                ? h(EzActionButton, {
                    iconOnly: true,
                    kind: 'view',
                    label: '预览',
                    size: 'tiny',
                    secondary: true,
                    type: 'primary',
                    onClick: () => handlePreviewTemplate(row),
                  })
                : null,
              canUse('system:mail:send')
                ? h(EzActionButton, {
                    iconOnly: true,
                    kind: 'save',
                    label: '按模板发送',
                    size: 'tiny',
                    secondary: true,
                    type: 'success',
                    onClick: () => openSendMail(row),
                  })
                : null,
              canUse('system:mail:template:status')
                ? h(
                    NPopconfirm,
                    { onPositiveClick: () => toggleTemplateStatus(row, nextStatus) },
                    {
                      trigger: () =>
                        h(EzActionButton, {
                          iconOnly: true,
                          kind: nextStatus === MailStatus.Disabled ? 'disable' : 'enable',
                          label: nextStatus === MailStatus.Disabled ? '禁用' : '启用',
                          size: 'tiny',
                          secondary: true,
                          tooltip: false,
                          type: nextStatus === MailStatus.Disabled ? 'error' : 'success',
                        }),
                      default: () =>
                        `确认${nextStatus === MailStatus.Disabled ? '禁用' : '启用'}该邮件模板？`,
                    },
                  )
                : null,
              canUse('system:mail:template:delete')
                ? h(
                    NPopconfirm,
                    { onPositiveClick: () => handleDeleteTemplate(row) },
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
                      default: () => '确认删除该邮件模板？',
                    },
                  )
                : null,
            ].filter(Boolean),
        })
      },
    },
  ]

  const logColumns: DataTableColumns<MailLogItem> = [
    {
      title: '主题',
      key: 'subject',
      minWidth: 220,
      ellipsis: { tooltip: true },
    },
    {
      title: '收件人',
      key: 'to_emails',
      minWidth: 200,
      ellipsis: { tooltip: true },
      render(row) {
        return row.to_emails.join(', ')
      },
    },
    {
      title: '邮箱账号',
      key: 'account_name',
      width: 140,
      ellipsis: { tooltip: true },
    },
    {
      title: '模板',
      key: 'template_code',
      width: 140,
      ellipsis: { tooltip: true },
      render(row) {
        return row.template_code || '-'
      },
    },
    {
      title: '状态',
      key: 'status',
      width: 78,
      align: 'center',
      render(row) {
        return h(
          NTag,
          { type: row.status === MailLogStatus.Success ? 'success' : 'error', bordered: false },
          { default: () => (row.status === MailLogStatus.Success ? '成功' : '失败') },
        )
      },
    },
    {
      title: '错误',
      key: 'error_message',
      minWidth: 180,
      ellipsis: { tooltip: true },
      render(row) {
        return displayText(row.error_message)
      },
    },
    {
      title: '发送时间',
      key: 'created_at',
      width: 156,
      render(row) {
        return h('span', { class: 'tabular-nums text-[var(--ez-text-muted)]' }, formatTime(row.created_at))
      },
    },
  ]

  const encryptionOptions: SelectOption[] = [
    { label: 'SSL/TLS', value: MailEncryption.SSL },
    { label: 'STARTTLS', value: MailEncryption.STARTTLS },
    { label: '无加密', value: MailEncryption.None },
  ]

  const logStatusOptions: SelectOption[] = [
    { label: '状态：全部', value: 0 },
    { label: '成功', value: MailLogStatus.Success },
    { label: '失败', value: MailLogStatus.Failed },
  ]

  watch(
    () => sendForm.formModel.template_code,
    (code) => {
      const template = templatesPage.items.value.find((item) => item.code === code)
      if (!template) {
        return
      }
      sendForm.formModel.subject = ''
      sendForm.formModel.content = ''
      sendForm.formModel.is_html = template.is_html
      sendForm.formModel.variables_text = buildVariablesText(template.variables)
    },
  )

  function openCreateAccount() {
    accountForm.openCreate()
  }

  function openEditAccount(row: MailAccountItem) {
    accountForm.openEdit(toMailAccountFormModel(row))
  }

  function openCreateTemplate() {
    templateForm.openCreate()
  }

  function openEditTemplate(row: MailTemplateItem) {
    templateForm.openEdit(toMailTemplateFormModel(row))
  }

  function openSendMail(template?: MailTemplateItem) {
    sendForm.openCreate()
    if (template) {
      sendForm.formModel.template_code = template.code
      sendForm.formModel.variables_text = buildVariablesText(template.variables)
    }
  }

  function openTestAccount(row: MailAccountItem) {
    testForm.openCreate()
    testForm.formModel.account_id = row.id
  }

  async function submitAccountForm() {
    const payload = buildMailAccountPayload(accountForm.formModel)
    if (accountForm.formMode.value === 'create') {
      await createMailAccount(payload)
      message.success('邮箱账号创建成功')
    } else {
      await updateMailAccount(accountForm.formModel.id, payload)
      message.success('邮箱账号更新成功')
    }
    await accountsPage.load()
  }

  async function submitTemplateForm() {
    if (templateForm.formMode.value === 'create') {
      await createMailTemplate(buildCreateMailTemplatePayload(templateForm.formModel))
      message.success('邮件模板创建成功')
    } else {
      await updateMailTemplate(
        templateForm.formModel.id,
        buildUpdateMailTemplatePayload(templateForm.formModel),
      )
      message.success('邮件模板更新成功')
    }
    await templatesPage.load()
  }

  async function submitSendForm() {
    const result = await sendMail(buildMailSendPayload(sendForm.formModel))
    if (result.status === MailLogStatus.Success) {
      message.success('邮件已发送')
    } else {
      message.error('邮件发送失败，已记录发送日志')
    }
    await logsPage.load()
  }

  async function submitTestForm() {
    const result = await testMailAccount(testForm.formModel.account_id, {
      to: parseTextList(testForm.formModel.to_text),
      subject: testForm.formModel.subject,
      content: testForm.formModel.content,
    })
    if (result.status === MailLogStatus.Success) {
      message.success('测试邮件已发送')
    } else {
      message.error('测试邮件发送失败，已记录发送日志')
    }
    await accountsPage.load()
    await logsPage.load()
  }

  async function handleDeleteAccount(row: MailAccountItem) {
    await deleteMailAccount(row.id)
    message.success('邮箱账号已删除')
    await accountsPage.load()
  }

  async function handleDeleteTemplate(row: MailTemplateItem) {
    await deleteMailTemplate(row.id)
    message.success('邮件模板已删除')
    await templatesPage.load()
  }

  async function handlePreviewTemplate(row: MailTemplateItem) {
    const rendered = await renderMailTemplate(row.id, {
      variables: parseVariablesText(buildVariablesText(row.variables)),
    })
    preview.title = row.name
    preview.subject = rendered.subject
    preview.content = rendered.content
    preview.visible = true
  }

  function encryptionText(value: MailEncryption) {
    switch (value) {
      case MailEncryption.SSL:
        return 'SSL'
      case MailEncryption.STARTTLS:
        return 'STARTTLS'
      default:
        return '无'
    }
  }

  return {
    accountColumns,
    accountForm,
    accountOptions,
    accountsPage,
    activePanel,
    canUse,
    encryptionOptions,
    logColumns,
    logStatusOptions,
    logsPage,
    openCreateAccount,
    openCreateTemplate,
    openSendMail,
    preview,
    sendForm,
    submitAccountForm,
    submitSendForm,
    submitTemplateForm,
    submitTestForm,
    templateColumns,
    templateForm,
    templateOptions,
    templatesPage,
    testForm,
  }
}
