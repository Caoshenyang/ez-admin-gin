import { ref } from 'vue'

import { useRemotePagination } from '@/composables/useRemotePagination'
import { getOperationLogs } from '../api/operation-log'
import type { OperationLogItem, OperationLogListQuery } from '../types/operation-log'

export type RiskLevel = 'high' | 'medium' | 'low'

export const riskMeta: Record<RiskLevel, { label: string; tagType: 'error' | 'warning' | 'success'; bg: string }> = {
  high: { label: '高风险', tagType: 'error', bg: '#fef2f2' },
  medium: { label: '中风险', tagType: 'warning', bg: '#fff7ed' },
  low: { label: '低风险', tagType: 'success', bg: '#f0fdf4' },
}

export const methodOptions = [
  { label: '方法：全部', value: '' },
  { label: 'GET', value: 'GET' },
  { label: 'POST', value: 'POST' },
  { label: 'PUT', value: 'PUT' },
  { label: 'DELETE', value: 'DELETE' },
]

export const successOptions = [
  { label: '结果：全部', value: '' },
  { label: '成功', value: 'true' },
  { label: '失败', value: 'false' },
]

export function getRiskLevel(row: OperationLogItem): RiskLevel {
  if (!row.success) return 'high'
  if (row.method === 'POST' || row.method === 'PUT' || row.method === 'DELETE') return 'medium'
  return 'low'
}

export function getModule(path: string): string {
  const segments = path.replace(/^\/api\/v\d+\//, '').split('/')
  const moduleKey = segments[1]
  const moduleMap: Record<string, string> = {
    users: '用户管理',
    roles: '角色权限',
    menus: '菜单管理',
    configs: '配置管理',
    'dict-types': '字典管理',
    'dict-items': '字典管理',
    files: '文件管理',
    attachments: '附件管理',
    notices: '公告管理',
    'operation-logs': '操作日志',
    'login-logs': '登录日志',
    health: '系统状态',
    auth: '认证',
  }

  return moduleKey ? (moduleMap[moduleKey] ?? moduleKey) : path
}

export function getAction(row: OperationLogItem): string {
  if (!row.success && row.error_message) return row.error_message

  const actionMap: Record<string, string> = {
    GET: '查询',
    POST: '提交',
    PUT: '更新',
    DELETE: '删除',
  }

  return actionMap[row.method] ?? row.method
}

export function formatTimeFull(value: string) {
  return value ? new Date(value).toLocaleString() : '-'
}

export function useOperationLogPage() {
  const detailVisible = ref(false)
  const detailRow = ref<OperationLogItem | null>(null)

  const {
    items: logs,
    total,
    loading,
    query,
    load,
    handleSearch,
    handleReset: handleResetBase,
    handlePageChange,
    handlePageSizeChange,
  } = useRemotePagination<OperationLogItem, OperationLogListQuery>(
    (params) =>
      getOperationLogs({
        ...params,
        username: params.username?.trim() || undefined,
        method: params.method || undefined,
        path: params.path?.trim() || undefined,
        success: params.success || undefined,
      }),
    {
      page: 1,
      page_size: 10,
      username: '',
      method: '',
      path: '',
      success: '',
    },
  )

  function handleReset() {
    handleResetBase()
    detailVisible.value = false
  }

  function openDetail(row: OperationLogItem) {
    detailRow.value = row
    detailVisible.value = true
  }

  return {
    detailRow,
    detailVisible,
    handlePageChange,
    handlePageSizeChange,
    handleReset,
    handleSearch,
    load,
    loading,
    logs,
    openDetail,
    query,
    total,
  }
}
