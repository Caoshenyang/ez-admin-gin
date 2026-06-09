import { useRemotePagination } from '@/composables/useRemotePagination'
import { getLoginLogs } from '../api/login-log'
import type { LoginLogItem, LoginLogListQuery } from '../types/login-log'
import { defaultLoginLogQuery, normalizeLoginLogQuery } from './login-log-page.utils'

export function useLoginLogPage() {
  const {
    items: logs,
    total,
    loading,
    query,
    load,
    handleSearch,
    handleReset,
    handlePageChange,
    handlePageSizeChange,
  } = useRemotePagination<LoginLogItem, LoginLogListQuery>(
    (params) => getLoginLogs(normalizeLoginLogQuery(params)),
    defaultLoginLogQuery(),
  )

  return {
    handlePageChange,
    handlePageSizeChange,
    handleReset,
    handleSearch,
    load,
    loading,
    logs,
    query,
    total,
  }
}
