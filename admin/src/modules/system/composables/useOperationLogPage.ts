import { ref } from 'vue'

import { useRemotePagination } from '@/composables/useRemotePagination'
import { getOperationLogs } from '../api/operation-log'
import type { OperationLogItem } from '../types/operation-log'
import type { OperationLogPageQuery } from '../types/operation-log-page'
import { defaultOperationLogQuery, normalizeOperationLogQuery } from './operation-log-page.utils'

// 操作日志管理页面组合式函数，封装日志列表查询、详情查看等逻辑
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
  } = useRemotePagination<OperationLogItem, OperationLogPageQuery>(
    (params) => getOperationLogs(normalizeOperationLogQuery(params)),
    defaultOperationLogQuery(),
  )

  // 重置搜索条件并关闭详情弹窗
  function handleReset() {
    handleResetBase()
    closeDetail()
  }

  // 打开日志详情弹窗
  function openDetail(row: OperationLogItem) {
    detailRow.value = row
    detailVisible.value = true
  }

  // 抽屉开关统一从这里进出，关闭时顺手清掉旧详情数据。
  function handleDetailVisibleChange(visible: boolean) {
    if (visible) {
      detailVisible.value = true
      return
    }

    closeDetail()
  }

  // 关闭详情弹窗，并清空当前选中的日志行
  function closeDetail() {
    detailVisible.value = false
    detailRow.value = null
  }

  return {
    closeDetail,
    detailRow,
    detailVisible,
    handleDetailVisibleChange,
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
