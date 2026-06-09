export const MailStatus = {
  Enabled: 1,
  Disabled: 2,
} as const

export type MailStatus = (typeof MailStatus)[keyof typeof MailStatus]

export const MailLogStatus = {
  Success: 1,
  Failed: 2,
} as const

export type MailLogStatus = (typeof MailLogStatus)[keyof typeof MailLogStatus]

export const MailEncryption = {
  None: 'none',
  SSL: 'ssl',
  STARTTLS: 'starttls',
} as const

export type MailEncryption = (typeof MailEncryption)[keyof typeof MailEncryption]

export interface MailAccountItem {
  id: number
  name: string
  host: string
  port: number
  username: string
  from_email: string
  from_name: string
  encryption: MailEncryption
  is_default: boolean
  status: MailStatus
  last_test_at: string | null
  last_test_msg: string
  remark: string
  created_at: string
  updated_at: string
}

export interface MailAccountListQuery {
  page: number
  page_size: number
  keyword?: string
  status?: MailStatus | 0
}

export interface MailAccountListResponse {
  items: MailAccountItem[]
  total: number
  page: number
  page_size: number
}

export interface CreateMailAccountPayload {
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

export type UpdateMailAccountPayload = CreateMailAccountPayload

export interface UpdateMailStatusPayload {
  status: MailStatus
}

export interface TestMailAccountPayload {
  to: string[]
  subject: string
  content: string
}

export interface MailTemplateItem {
  id: number
  code: string
  name: string
  subject: string
  content: string
  is_html: boolean
  variables: string[]
  sort: number
  status: MailStatus
  remark: string
  created_at: string
  updated_at: string
}

export interface MailTemplateListQuery {
  page: number
  page_size: number
  keyword?: string
  status?: MailStatus | 0
}

export interface MailTemplateListResponse {
  items: MailTemplateItem[]
  total: number
  page: number
  page_size: number
}

export interface CreateMailTemplatePayload {
  code: string
  name: string
  subject: string
  content: string
  is_html: boolean
  variables: string[]
  sort: number
  status: MailStatus
  remark: string
}

export interface UpdateMailTemplatePayload {
  name: string
  subject: string
  content: string
  is_html: boolean
  variables: string[]
  sort: number
  status: MailStatus
  remark: string
}

export interface RenderMailTemplatePayload {
  variables: Record<string, string>
}

export interface RenderMailTemplateResponse {
  subject: string
  content: string
}

export interface SendMailPayload {
  account_id: number
  template_code: string
  to: string[]
  cc: string[]
  bcc: string[]
  subject: string
  content: string
  is_html: boolean
  variables: Record<string, string>
}

export interface SendMailResponse {
  log_id: number
  status: MailLogStatus
}

export interface MailLogItem {
  id: number
  account_id: number
  account_name: string
  template_id: number
  template_code: string
  subject: string
  from_email: string
  to_emails: string[]
  cc_emails: string[]
  bcc_emails: string[]
  status: MailLogStatus
  error_message: string
  created_at: string
}

export interface MailLogListQuery {
  page: number
  page_size: number
  keyword?: string
  status?: MailLogStatus | 0
  account_id?: number
  template_code?: string
}

export interface MailLogListResponse {
  items: MailLogItem[]
  total: number
  page: number
  page_size: number
}
