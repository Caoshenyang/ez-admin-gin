import type { SystemHealthData } from '../types/health'
import type { HealthDependencyCard, HealthEndpointCard } from '../types/health-page'

export function healthDependencyCards(health: SystemHealthData | null): HealthDependencyCard[] {
  return [
    {
      key: 'database',
      label: '数据库',
      value: health?.database,
      description: '验证 PostgreSQL 连接是否可用',
    },
    {
      key: 'redis',
      label: 'Redis',
      value: health?.redis,
      description: '验证缓存和会话依赖是否可用',
    },
  ]
}

export const healthEndpointCards: HealthEndpointCard[] = [
  {
    title: '后台接口',
    path: '/api/v1/system/health',
    description: '需要登录和权限，适合在管理台里确认认证链路与依赖状态。',
  },
  {
    title: '公开探针',
    path: '/health',
    description: '给 Nginx、容器健康检查和外部监控使用，不依赖登录态。',
  },
]

export function formatHealthStatusLabel(value?: string) {
  return value === 'ok' ? '正常' : value || '待检查'
}

export function healthStatusTagType(value?: string) {
  return value === 'ok' ? 'success' : 'error'
}

export function healthStatusText(value: string | undefined, loading: boolean) {
  if (value === 'ok') {
    return '服务连通性正常'
  }

  if (loading) {
    return '正在刷新状态...'
  }

  return '请点击刷新重新检查'
}

export function healthErrorMessage(error: unknown) {
  if (typeof error === 'object' && error !== null) {
    const response = (error as { response?: { data?: { message?: string } } }).response
    if (typeof response?.data?.message === 'string' && response.data.message) {
      return response.data.message
    }
  }

  return '系统状态获取失败，请稍后重试。'
}
