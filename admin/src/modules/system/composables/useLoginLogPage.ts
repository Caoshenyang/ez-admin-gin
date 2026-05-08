import { useRemotePagination } from '@/composables/useRemotePagination'
import { getLoginLogs } from '../api/login-log'
import type { LoginLogItem, LoginLogListQuery } from '../types/login-log'
import { defaultLoginLogQuery, normalizeLoginLogQuery } from './login-log-page.utils'

// 登录日志页面组合式函数，统一封装查询、分页和刷新逻辑。
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
