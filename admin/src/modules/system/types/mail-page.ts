import type { MailEncryption, MailStatus } from './mail'

export type MailPanelKey = 'accounts' | 'templates' | 'logs'

export interface MailAccountFormModel {
  id: number
  name: string
  host: string
  port: number
  username: string
  password: string
  from_email: string
  from_name: string
  encryption: MailEncryption
  is_default: boolean
  status: MailStatus
  remark: string
}

export interface MailTemplateFormModel {
  id: number
  code: string
  name: string
  subject: string
  content: string
  is_html: boolean
  variables_text: string
  sort: number
  status: MailStatus
  remark: string
}

export interface MailTestFormModel {
  account_id: number
  to_text: string
  subject: string
  content: string
}

export interface MailSendFormModel {
  account_id: number
  template_code: string
  to_text: string
  cc_text: string
  bcc_text: string
  subject: string
  content: string
  is_html: boolean
  variables_text: string
}

export interface MailPreviewState {
  visible: boolean
  title: string
  subject: string
  content: string
}
