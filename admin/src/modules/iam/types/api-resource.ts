export const ApiStatus = {
  Enabled: 1,
  Disabled: 2,
} as const

export type ApiStatus = (typeof ApiStatus)[keyof typeof ApiStatus]

export interface AdminAPI {
  id: number
  code: string
  name: string
  module: string
  method: string
  path: string
  sort: number
  status: ApiStatus
  remark: string
  created_at: string
  updated_at: string
}
