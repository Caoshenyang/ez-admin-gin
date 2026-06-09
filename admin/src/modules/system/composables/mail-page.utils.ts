import {
  MailEncryption,
  MailStatus,
  type MailAccountItem,
  type MailAccountListQuery,
  type MailLogListQuery,
  type MailTemplateItem,
  type MailTemplateListQuery,
} from '../types/mail'
import type {
  MailAccountFormModel,
  MailSendFormModel,
  MailTemplateFormModel,
  MailTestFormModel,
} from '../types/mail-page'

export function defaultMailAccountListQuery(): MailAccountListQuery {
  return {
    page: 1,
    page_size: 10,
    keyword: '',
    status: 0,
  }
}

export function defaultMailTemplateListQuery(): MailTemplateListQuery {
  return {
    page: 1,
    page_size: 10,
    keyword: '',
    status: 0,
  }
}

export function defaultMailLogListQuery(): MailLogListQuery {
  return {
    page: 1,
    page_size: 10,
    keyword: '',
    status: 0,
    account_id: undefined,
    template_code: '',
  }
}

export function defaultMailAccountFormModel(): MailAccountFormModel {
  return {
    id: 0,
    name: '',
    host: '',
    port: 465,
    username: '',
    password: '',
    from_email: '',
    from_name: '',
    encryption: MailEncryption.SSL,
    is_default: false,
    status: MailStatus.Enabled,
    remark: '',
  }
}

export function toMailAccountFormModel(account: MailAccountItem): MailAccountFormModel {
  return {
    id: account.id,
    name: account.name,
    host: account.host,
    port: account.port,
    username: account.username,
    password: '',
    from_email: account.from_email,
    from_name: account.from_name,
    encryption: account.encryption,
    is_default: account.is_default,
    status: account.status,
    remark: account.remark,
  }
}

export function buildMailAccountPayload(formModel: MailAccountFormModel) {
  return {
    name: formModel.name,
    host: formModel.host,
    port: formModel.port,
    username: formModel.username,
    password: formModel.password,
    from_email: formModel.from_email,
    from_name: formModel.from_name,
    encryption: formModel.encryption,
    is_default: formModel.is_default,
    status: formModel.status,
    remark: formModel.remark,
  }
}

export function defaultMailTemplateFormModel(): MailTemplateFormModel {
  return {
    id: 0,
    code: '',
    name: '',
    subject: '',
    content: '',
    is_html: true,
    variables_text: '',
    sort: 0,
    status: MailStatus.Enabled,
    remark: '',
  }
}

export function toMailTemplateFormModel(template: MailTemplateItem): MailTemplateFormModel {
  return {
    id: template.id,
    code: template.code,
    name: template.name,
    subject: template.subject,
    content: template.content,
    is_html: template.is_html,
    variables_text: template.variables.join(', '),
    sort: template.sort,
    status: template.status,
    remark: template.remark,
  }
}

export function buildCreateMailTemplatePayload(formModel: MailTemplateFormModel) {
  return {
    code: formModel.code,
    name: formModel.name,
    subject: formModel.subject,
    content: formModel.content,
    is_html: formModel.is_html,
    variables: parseTextList(formModel.variables_text),
    sort: formModel.sort,
    status: formModel.status,
    remark: formModel.remark,
  }
}

export function buildUpdateMailTemplatePayload(formModel: MailTemplateFormModel) {
  return {
    name: formModel.name,
    subject: formModel.subject,
    content: formModel.content,
    is_html: formModel.is_html,
    variables: parseTextList(formModel.variables_text),
    sort: formModel.sort,
    status: formModel.status,
    remark: formModel.remark,
  }
}

export function defaultMailTestFormModel(): MailTestFormModel {
  return {
    account_id: 0,
    to_text: '',
    subject: 'EZ Admin 邮箱测试',
    content: '这是一封来自 EZ Admin 的邮箱配置测试邮件。',
  }
}

export function defaultMailSendFormModel(): MailSendFormModel {
  return {
    account_id: 0,
    template_code: '',
    to_text: '',
    cc_text: '',
    bcc_text: '',
    subject: '',
    content: '',
    is_html: true,
    variables_text: '',
  }
}

export function buildMailSendPayload(formModel: MailSendFormModel) {
  return {
    account_id: formModel.account_id,
    template_code: formModel.template_code,
    to: parseTextList(formModel.to_text),
    cc: parseTextList(formModel.cc_text),
    bcc: parseTextList(formModel.bcc_text),
    subject: formModel.subject,
    content: formModel.content,
    is_html: formModel.is_html,
    variables: parseVariablesText(formModel.variables_text),
  }
}

export function parseTextList(value: string) {
  return value
    .split(/[\n,，;]/)
    .map((item) => item.trim())
    .filter(Boolean)
}

export function parseVariablesText(value: string) {
  const result: Record<string, string> = {}
  for (const line of value.split('\n')) {
    const trimmed = line.trim()
    if (!trimmed) {
      continue
    }
    const index = trimmed.indexOf('=')
    if (index <= 0) {
      continue
    }
    const key = trimmed.slice(0, index).trim()
    const itemValue = trimmed.slice(index + 1).trim()
    if (key) {
      result[key] = itemValue
    }
  }
  return result
}

export function buildVariablesText(keys: string[]) {
  return keys.map((key) => `${key}=`).join('\n')
}
