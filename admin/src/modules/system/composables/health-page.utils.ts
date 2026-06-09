import {
  CheckmarkCircleOutline,
  CloudDoneOutline,
  GitNetworkOutline,
  HardwareChipOutline,
  PulseOutline,
  RadioOutline,
  ServerOutline,
  ShieldCheckmarkOutline,
  SpeedometerOutline,
  TerminalOutline,
  TimeOutline,
  WarningOutline,
  WifiOutline,
} from '@vicons/ionicons5'

import type { SystemHealthData } from '../types/health'
import type {
  HealthCheckItem,
  HealthDependencyCard,
  HealthEndpointCard,
  HealthOverview,
  HealthSummaryCard,
  HealthTagType,
} from '../types/health-page'

const dependencyTotal = 2

function isStatusOk(value?: string) {
  return value === 'ok'
}

function normalizeHealthEnv(value?: string) {
  return value?.trim() ?? ''
}

function isKnownHealthEnv(value?: string) {
  const env = normalizeHealthEnv(value)
  return env !== '' && env !== 'unknown'
}

function healthReadyCount(health: SystemHealthData | null) {
  if (!health) {
    return 0
  }

  return [health.database, health.redis].filter(isStatusOk).length
}

function healthToneClass(value?: string, loading = false) {
  if (isStatusOk(value)) {
    return 'health-tone--success'
  }

  if (value) {
    return 'health-tone--danger'
  }

  return loading ? 'health-tone--warning' : 'health-tone--neutral'
}

function healthProgress(value?: string, healthLoaded = false) {
  if (isStatusOk(value)) {
    return 100
  }

  if (value || healthLoaded) {
    return 28
  }

  return 0
}

function checkItemStatus(value?: string, loading = false) {
  if (isStatusOk(value)) {
    return {
      label: '通过',
      tagType: 'success' as const,
      toneClass: 'health-tone--success',
    }
  }

  if (value) {
    return {
      label: formatHealthStatusLabel(value),
      tagType: 'error' as const,
      toneClass: 'health-tone--danger',
    }
  }

  return {
    label: loading ? '检查中' : '待检查',
    tagType: loading ? ('warning' as const) : ('default' as const),
    toneClass: loading ? 'health-tone--warning' : 'health-tone--neutral',
  }
}

export function healthIsReady(health: SystemHealthData | null) {
  return healthReadyCount(health) === dependencyTotal
}

export function healthReadinessScore(health: SystemHealthData | null, loading: boolean) {
  if (!health) {
    return loading ? 36 : 0
  }

  const readyCount = healthReadyCount(health)
  if (readyCount === dependencyTotal) {
    return 100
  }

  return readyCount === 1 ? 62 : 28
}

export function healthOverview(
  health: SystemHealthData | null,
  loading: boolean,
  errorMessage: string,
): HealthOverview {
  if (errorMessage) {
    return {
      description: '保留最近一次快照，刷新后重新校验依赖连通性。',
      statusLabel: '失败',
      tagType: 'error',
      title: '状态检查失败',
      toneClass: 'health-overview--danger',
    }
  }

  if (healthIsReady(health)) {
    return {
      description: '数据库、Redis 和受保护接口均返回正常信号。',
      statusLabel: '稳定',
      tagType: 'success',
      title: '核心依赖运行稳定',
      toneClass: 'health-overview--success',
    }
  }

  if (health) {
    return {
      description: '已有依赖返回异常状态，请优先处理环境连接问题。',
      statusLabel: '需处理',
      tagType: 'error',
      title: '依赖存在异常',
      toneClass: 'health-overview--danger',
    }
  }

  return {
    description: loading ? '正在拉取后台运行快照。' : '刷新后同步运行环境与依赖状态。',
    statusLabel: loading ? '同步中' : '待检查',
    tagType: loading ? 'warning' : 'default',
    title: loading ? '正在同步运行快照' : '等待首次检查',
    toneClass: loading ? 'health-overview--warning' : 'health-overview--neutral',
  }
}

export function healthSummaryCards(
  health: SystemHealthData | null,
  lastCheckedLabel: string,
  loading: boolean,
): HealthSummaryCard[] {
  const readyCount = healthReadyCount(health)
  const isReady = readyCount === dependencyTotal
  const env = normalizeHealthEnv(health?.env)
  const hasKnownEnv = isKnownHealthEnv(health?.env)

  return [
    {
      detail: env === 'prod' ? '生产配置' : hasKnownEnv ? '开发 / 测试配置' : '等待同步',
      icon: TerminalOutline,
      key: 'env',
      label: '运行环境',
      tagLabel: env === 'prod' ? 'PROD' : hasKnownEnv ? 'DEV' : 'Pending',
      tagType: env === 'prod' ? 'success' : hasKnownEnv ? 'warning' : 'default',
      toneClass: hasKnownEnv ? 'health-tone--brand' : 'health-tone--neutral',
      value: hasKnownEnv ? env : '待同步',
    },
    {
      detail: isReady ? '全部依赖在线' : health ? '存在异常信号' : '未完成检查',
      icon: SpeedometerOutline,
      key: 'dependencies',
      label: '依赖在线',
      tagLabel: isReady ? 'Ready' : health ? 'Issue' : loading ? 'Syncing' : 'Pending',
      tagType: isReady ? 'success' : health ? 'error' : loading ? 'warning' : 'default',
      toneClass: isReady
        ? 'health-tone--success'
        : health
          ? 'health-tone--danger'
          : 'health-tone--neutral',
      value: `${readyCount}/${dependencyTotal}`,
    },
    {
      detail: loading ? '刷新中' : '最近一次检查',
      icon: TimeOutline,
      key: 'lastChecked',
      label: '刷新时间',
      tagLabel: loading ? 'Live' : 'Snapshot',
      tagType: loading ? 'warning' : 'info',
      toneClass: loading ? 'health-tone--warning' : 'health-tone--brand',
      value: lastCheckedLabel,
    },
    {
      detail: '/api/v1/system/health',
      icon: ShieldCheckmarkOutline,
      key: 'guard',
      label: '接口链路',
      tagLabel: health ? 'OK' : loading ? 'Syncing' : 'Pending',
      tagType: health ? 'success' : loading ? 'warning' : 'default',
      toneClass: health ? 'health-tone--success' : 'health-tone--neutral',
      value: health ? '受保护' : '待确认',
    },
  ]
}

export function healthDependencyCards(
  health: SystemHealthData | null,
  loading = false,
): HealthDependencyCard[] {
  const healthLoaded = Boolean(health)

  return [
    {
      description: '用户、权限、配置、日志等核心数据存储',
      detail: isStatusOk(health?.database) ? '连接可用' : '请检查数据库连接与迁移状态',
      icon: ServerOutline,
      key: 'database',
      label: '数据库',
      progress: healthProgress(health?.database, healthLoaded),
      progressLabel: isStatusOk(health?.database) ? 'ready' : health?.database || 'pending',
      statusLabel: formatHealthStatusLabel(health?.database),
      tagType: healthStatusTagType(health?.database, loading),
      toneClass: healthToneClass(health?.database, loading),
      value: health?.database,
    },
    {
      description: '登录态、限流和缓存依赖',
      detail: isStatusOk(health?.redis) ? '连接可用' : '请检查 Redis 服务和连接配置',
      icon: HardwareChipOutline,
      key: 'redis',
      label: 'Redis',
      progress: healthProgress(health?.redis, healthLoaded),
      progressLabel: isStatusOk(health?.redis) ? 'ready' : health?.redis || 'pending',
      statusLabel: formatHealthStatusLabel(health?.redis),
      tagType: healthStatusTagType(health?.redis, loading),
      toneClass: healthToneClass(health?.redis, loading),
      value: health?.redis,
    },
  ]
}

export const healthEndpointCards: HealthEndpointCard[] = [
  {
    badge: 'AUTH',
    description: '需要登录和权限，适合在管理台里确认认证链路与依赖状态。',
    icon: ShieldCheckmarkOutline,
    method: 'GET',
    path: '/api/v1/system/health',
    title: '后台接口',
    toneClass: 'health-tone--brand',
  },
  {
    badge: 'LIVE',
    description: '外部监控、Nginx 和容器健康检查使用的轻量存活探针。',
    icon: GitNetworkOutline,
    method: 'GET',
    path: '/health',
    title: '公开探针',
    toneClass: 'health-tone--success',
  },
]

export function healthCheckItems(
  health: SystemHealthData | null,
  loading: boolean,
): HealthCheckItem[] {
  const authStatus = health
    ? { label: '通过', tagType: 'success' as const, toneClass: 'health-tone--success' }
    : {
        label: loading ? '检查中' : '待检查',
        tagType: loading ? ('warning' as const) : ('default' as const),
        toneClass: loading ? 'health-tone--warning' : 'health-tone--neutral',
      }
  const databaseStatus = checkItemStatus(health?.database, loading)
  const redisStatus = checkItemStatus(health?.redis, loading)
  const env = normalizeHealthEnv(health?.env)
  const hasKnownEnv = isKnownHealthEnv(health?.env)

  return [
    {
      description: '当前页面的受保护接口返回成功响应',
      key: 'auth',
      label: '认证链路',
      statusLabel: authStatus.label,
      tagType: authStatus.tagType,
      toneClass: authStatus.toneClass,
    },
    {
      description: '主数据库连接可用',
      key: 'database',
      label: '数据库连接',
      statusLabel: databaseStatus.label,
      tagType: databaseStatus.tagType,
      toneClass: databaseStatus.toneClass,
    },
    {
      description: '缓存和会话依赖可用',
      key: 'redis',
      label: 'Redis 连接',
      statusLabel: redisStatus.label,
      tagType: redisStatus.tagType,
      toneClass: redisStatus.toneClass,
    },
    {
      description: '当前运行环境标识已返回',
      key: 'env',
      label: '环境标识',
      statusLabel: hasKnownEnv ? env : loading ? '检查中' : '待检查',
      tagType: hasKnownEnv ? healthEnvTagType(env) : loading ? 'warning' : 'default',
      toneClass: hasKnownEnv
        ? 'health-tone--brand'
        : loading
          ? 'health-tone--warning'
          : 'health-tone--neutral',
    },
  ]
}

export function healthSignalIcons(health: SystemHealthData | null, loading: boolean) {
  if (healthIsReady(health)) {
    return {
      icon: CheckmarkCircleOutline,
      softIcon: CloudDoneOutline,
    }
  }

  if (health) {
    return {
      icon: WarningOutline,
      softIcon: WifiOutline,
    }
  }

  return {
    icon: loading ? PulseOutline : WarningOutline,
    softIcon: RadioOutline,
  }
}

export function formatHealthStatusLabel(value?: string) {
  return value === 'ok' ? '正常' : value || '待检查'
}

export function healthStatusTagType(value?: string, loading = false): HealthTagType {
  if (value === 'ok') {
    return 'success'
  }

  if (value) {
    return 'error'
  }

  return loading ? 'warning' : 'default'
}

export function healthEnvTagType(value?: string): HealthTagType {
  if (value === 'prod') {
    return 'success'
  }

  return value ? 'warning' : 'default'
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
