import { LoginLogStatus, type LoginLogListQuery } from '../types/login-log'

export const loginLogStatusOptions = [
  { label: '状态：全部', value: 0 },
  { label: '成功', value: LoginLogStatus.Success },
  { label: '失败', value: LoginLogStatus.Failed },
]

export function defaultLoginLogQuery(): LoginLogListQuery {
  return {
    page: 1,
    page_size: 10,
    username: '',
    ip: '',
    status: 0,
  }
}

export function normalizeLoginLogQuery(params: LoginLogListQuery) {
  return {
    ...params,
    username: params.username?.trim() || undefined,
    ip: params.ip?.trim() || undefined,
    status: params.status === 0 ? undefined : params.status,
  }
}
