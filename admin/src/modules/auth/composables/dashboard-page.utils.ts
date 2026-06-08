import {
  AlertCircleOutline,
  AnalyticsOutline,
  CheckmarkCircleOutline,
  FileTrayFullOutline,
  KeyOutline,
  LayersOutline,
  PeopleOutline,
  PulseOutline,
  ServerOutline,
  ShieldCheckmarkOutline,
  SpeedometerOutline,
  TimeOutline,
  WarningOutline,
} from '@vicons/ionicons5'

import { DashboardLoginStatus, type DashboardData } from '../types/dashboard'
import type {
  DashboardBarChartItem,
  DashboardChartSegment,
  DashboardChartStat,
  DashboardCommandItem,
  DashboardInsightCard,
  DashboardRingMetric,
  DashboardTrendPoint,
  DashboardResourceItem,
  DashboardTone,
  FocusFact,
  HealthSummaryItem,
  MetricCard,
  QuickLink,
} from '../types/dashboard-page'
import { MenuType, type AuthMenu } from '@/modules/iam/types/menu'

const quickLinkDescriptionMap: Record<string, string> = {
  系统状态: '检查当前环境、数据库和 Redis 状态',
  用户管理: '维护后台账号、状态和角色绑定',
  角色管理: '调整角色权限和菜单分配',
  菜单管理: '维护动态菜单、按钮权限和路由出口',
  系统配置: '查看系统参数和运行期配置项',
  文件管理: '检查上传文件和资源沉淀',
  操作日志: '回看后台操作链路和失败记录',
  登录日志: '检查最近登录结果和来源 IP',
  公告管理: '维护首页公告和系统通知内容',
}

export const dashboardHeroIcon = PulseOutline

interface DashboardCommandPaths {
  healthPath: string
  loginLogPath: string
  operationLogPath: string
}

export function dashboardCurrentDateLabel(date = new Date()) {
  return new Intl.DateTimeFormat('zh-CN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    weekday: 'long',
  }).format(date)
}

export function dashboardCurrentUserLabel(dashboard: DashboardData | null) {
  const user = dashboard?.current_user
  if (!user) {
    return '管理员'
  }

  return user.nickname || user.username
}

export function flattenPageMenus(menus: AuthMenu[]): AuthMenu[] {
  const result: AuthMenu[] = []

  for (const menu of menus) {
    if (menu.type === MenuType.Menu && menu.path) {
      result.push(menu)
    }

    result.push(...flattenPageMenus(menu.children ?? []))
  }

  return result
}

export function findMenuPathByTitle(menus: AuthMenu[], title: string) {
  return menus.find((menu) => menu.title === title)?.path || ''
}

export function dashboardIsHealthy(dashboard: DashboardData | null) {
  const health = dashboard?.health
  if (!health) {
    return false
  }

  return health.database === 'ok' && health.redis === 'ok'
}

export function dashboardHeroStatusText(dashboard: DashboardData | null, loading: boolean) {
  const health = dashboard?.health
  if (!health) {
    return loading ? '正在同步当前项目的运行快照...' : '等待首次同步项目运行数据'
  }

  if (health.database === 'ok' && health.redis === 'ok') {
    return '核心依赖在线，权限、操作和内容数据已经汇总到当前视图。'
  }

  if (health.database === 'ok') {
    return '数据库正常，Redis 有异常信号，建议优先检查缓存和登录态链路。'
  }

  return '核心依赖出现异常，请先处理环境问题，再继续做业务操作。'
}

export function formatMetricValue(value?: number) {
  return typeof value === 'number' ? new Intl.NumberFormat('zh-CN').format(value) : '--'
}

export function formatDashboardDateTime(value: string) {
  if (!value) {
    return '-'
  }

  const date = new Date(value)
  if (Number.isNaN(date.getTime())) {
    return '-'
  }

  return new Intl.DateTimeFormat('zh-CN', {
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
  }).format(date)
}

export function formatDashboardRoutePath(path: string) {
  return path.replace(/^\/api\/v1/, '') || path
}

export function dashboardHealthTagType(value?: string) {
  return value === 'ok' ? 'success' : value === 'pending' ? 'default' : 'error'
}

export function dashboardOperationStatusTagType(success: boolean) {
  return success ? 'success' : 'error'
}

export function dashboardLoginStatusTagType(status: number) {
  return status === DashboardLoginStatus.Success ? 'success' : 'error'
}

export function dashboardLoginStatusLabel(status: number) {
  return status === DashboardLoginStatus.Success ? '成功' : '失败'
}

export function dashboardHealthItems(dashboard: DashboardData | null): HealthSummaryItem[] {
  const health = dashboard?.health

  return [
    {
      label: '运行环境',
      value: health?.env || 'unknown',
      status: health ? 'ok' : 'pending',
    },
    {
      label: '数据库',
      value: health?.database || 'pending',
      status: health?.database || 'pending',
    },
    {
      label: 'Redis',
      value: health?.redis || 'pending',
      status: health?.redis || 'pending',
    },
  ]
}

export function dashboardMetricCards(
  dashboard: DashboardData | null,
  visiblePageTotal: number,
): MetricCard[] {
  const metrics = dashboard?.metrics

  return [
    {
      label: '启用账号',
      value: formatMetricValue(metrics?.enabled_user_total),
      hint: `总账号 ${formatMetricValue(metrics?.user_total)}`,
      iconClass: 'bg-[var(--ez-brand-soft)] text-[var(--ez-brand)]',
      panelClass: '',
      icon: ShieldCheckmarkOutline,
    },
    {
      label: '启用角色',
      value: formatMetricValue(metrics?.enabled_role_total),
      hint: `可访问页面 ${formatMetricValue(visiblePageTotal)}`,
      iconClass: 'bg-[var(--ez-accent-purple-soft)] text-[var(--ez-accent-purple)]',
      panelClass: '',
      icon: LayersOutline,
    },
    {
      label: '今日操作',
      value: formatMetricValue(metrics?.today_operation_total),
      hint: `失败 ${formatMetricValue(metrics?.today_risk_operation_total)}`,
      iconClass: 'bg-[var(--ez-warning-bg)] text-[var(--ez-warning-text)]',
      panelClass: '',
      icon: TimeOutline,
    },
    {
      label: '文件沉淀',
      value: formatMetricValue(metrics?.file_total),
      hint: `公告 ${formatMetricValue(metrics?.notice_total)} / 配置 ${formatMetricValue(metrics?.config_total)}`,
      iconClass: 'bg-[var(--ez-success-bg)] text-[var(--ez-success-text)]',
      panelClass: '',
      icon: PulseOutline,
    },
  ]
}

export function dashboardCommandItems(
  dashboard: DashboardData | null,
  paths: DashboardCommandPaths,
): DashboardCommandItem[] {
  const metrics = dashboard?.metrics
  const health = dashboard?.health
  const dependencyReady = dashboardIsHealthy(dashboard)
  const failedOperationTotal = metrics?.today_risk_operation_total ?? 0
  const todayOperationTotal = metrics?.today_operation_total ?? 0
  const failedLoginTotal = metrics?.today_login_failed_total ?? 0

  return [
    {
      actionText: '查看状态',
      description: health
        ? `数据库 ${health.database}，Redis ${health.redis}，环境 ${health.env || 'unknown'}。`
        : '健康检查尚未完成，刷新后可查看依赖状态。',
      icon: dependencyReady ? CheckmarkCircleOutline : ServerOutline,
      path: paths.healthPath,
      tagLabel: dependencyReady ? '稳定' : health ? '需处理' : '待同步',
      tagType: dependencyReady ? 'success' : health ? 'error' : 'default',
      title: '依赖健康',
      value: dependencyReady ? '在线' : health ? '异常' : '--',
    },
    {
      actionText: '看操作',
      description:
        failedOperationTotal > 0
          ? `今日 ${formatMetricValue(failedOperationTotal)} 次失败操作，建议先排查最近接口链路。`
          : `今日 ${formatMetricValue(todayOperationTotal)} 次操作已记录，最近行为可追踪。`,
      icon: failedOperationTotal > 0 ? AlertCircleOutline : AnalyticsOutline,
      path: paths.operationLogPath,
      tagLabel: failedOperationTotal > 0 ? '有风险' : '低风险',
      tagType: failedOperationTotal > 0 ? 'error' : 'success',
      title: '操作风险',
      value: formatMetricValue(failedOperationTotal),
    },
    {
      actionText: '查登录',
      description:
        failedLoginTotal > 0
          ? `今日 ${formatMetricValue(failedLoginTotal)} 次登录失败，需要关注账号或来源 IP。`
          : '今日暂无登录失败，最近登录记录保持可审计。',
      icon: failedLoginTotal > 0 ? WarningOutline : KeyOutline,
      path: paths.loginLogPath,
      tagLabel: failedLoginTotal > 0 ? '需关注' : '正常',
      tagType: failedLoginTotal > 0 ? 'warning' : 'success',
      title: '登录安全',
      value: formatMetricValue(failedLoginTotal),
    },
  ]
}

export function dashboardInsightCards(
  dashboard: DashboardData | null,
  visiblePageTotal: number,
): DashboardInsightCard[] {
  const metrics = dashboard?.metrics
  const userTotal = metrics?.user_total ?? 0
  const enabledUserTotal = metrics?.enabled_user_total ?? 0
  const accountPercent = formatPercent(enabledUserTotal, userTotal)

  return [
    {
      description: `启用率 ${accountPercent}，总账号 ${formatMetricValue(metrics?.user_total)}。`,
      icon: PeopleOutline,
      label: '账号治理',
      tone: 'info',
      value: formatMetricValue(metrics?.enabled_user_total),
    },
    {
      description: `${formatMetricValue(visiblePageTotal)} 个页面在当前角色下可访问。`,
      icon: ShieldCheckmarkOutline,
      label: '权限覆盖',
      tone: 'success',
      value: formatMetricValue(metrics?.enabled_role_total),
    },
    {
      description: `失败/风险 ${formatMetricValue(metrics?.today_risk_operation_total)}。`,
      icon: SpeedometerOutline,
      label: '今日操作',
      tone: (metrics?.today_risk_operation_total ?? 0) > 0 ? 'warning' : 'success',
      value: formatMetricValue(metrics?.today_operation_total),
    },
    {
      description: `公告 ${formatMetricValue(metrics?.notice_total)}，配置 ${formatMetricValue(metrics?.config_total)}。`,
      icon: FileTrayFullOutline,
      label: '内容资产',
      tone: 'default',
      value: formatMetricValue(metrics?.file_total),
    },
  ]
}

export function dashboardResourceItems(dashboard: DashboardData | null): DashboardResourceItem[] {
  const metrics = dashboard?.metrics
  const userTotal = metrics?.user_total ?? 0
  const enabledUserTotal = metrics?.enabled_user_total ?? 0
  const todayOperationTotal = metrics?.today_operation_total ?? 0
  const failedOperationTotal = metrics?.today_risk_operation_total ?? 0
  const successOperationTotal = Math.max(todayOperationTotal - failedOperationTotal, 0)
  const accountPercent = percentValue(enabledUserTotal, userTotal)
  const operationPercent = percentValue(successOperationTotal, todayOperationTotal)

  return [
    {
      detail: `${formatMetricValue(metrics?.enabled_user_total)} / ${formatMetricValue(metrics?.user_total)} 个账号启用`,
      label: '账号启用率',
      percent: accountPercent,
      tone: accountPercent >= 80 ? 'success' : accountPercent >= 50 ? 'warning' : 'error',
      value: userTotal > 0 ? `${accountPercent}%` : '--',
    },
    {
      detail: `${formatMetricValue(successOperationTotal)} / ${formatMetricValue(metrics?.today_operation_total)} 次操作成功`,
      label: '操作成功率',
      percent: operationPercent,
      tone: failedOperationTotal > 0 ? 'warning' : 'success',
      value: todayOperationTotal > 0 ? `${operationPercent}%` : '暂无操作',
    },
  ]
}

export function dashboardTrendPoints(dashboard: DashboardData | null): DashboardTrendPoint[] {
  const operations = [...(dashboard?.recent_operations ?? [])]
    .sort((left, right) => new Date(left.created_at).getTime() - new Date(right.created_at).getTime())
    .slice(-6)

  if (operations.length === 0) {
    return [
      {
        label: '--',
        latency: 0,
        name: '暂无数据',
        risk: 0,
        success: true,
      },
    ]
  }

  return operations.map((item) => ({
    label: formatDashboardTimeLabel(item.created_at),
    latency: item.latency_ms,
    name: formatCompactRouteLabel(item.path),
    risk: Math.min(Math.round(item.latency_ms / 18) + (item.success ? 8 : 28), 100),
    success: item.success,
  }))
}

export function dashboardRingMetrics(dashboard: DashboardData | null): DashboardRingMetric[] {
  const metrics = dashboard?.metrics
  const userTotal = metrics?.user_total ?? 0
  const enabledUserTotal = metrics?.enabled_user_total ?? 0
  const operationTotal = metrics?.today_operation_total ?? 0
  const failedOperationTotal = metrics?.today_risk_operation_total ?? 0
  const successfulOperationTotal = Math.max(operationTotal - failedOperationTotal, 0)
  const dependencyPercent = dashboardDependencyPercent(dashboard)
  const safetyScore = dashboardSafetyScore(dashboard)
  const accountPercent = percentValue(enabledUserTotal, userTotal)
  const operationPercent = percentValue(successfulOperationTotal, operationTotal)

  return [
    {
      detail: `失败操作 ${formatMetricValue(failedOperationTotal)} / 登录失败 ${formatMetricValue(metrics?.today_login_failed_total)}`,
      label: '安全评分',
      percent: safetyScore,
      tone: scoreTone(safetyScore),
      value: `${safetyScore}`,
    },
    {
      detail: `${formatMetricValue(enabledUserTotal)} / ${formatMetricValue(userTotal)} 个账号可用`,
      label: '账号启用率',
      percent: accountPercent,
      tone: percentTone(accountPercent),
      value: userTotal > 0 ? `${accountPercent}%` : '--',
    },
    {
      detail: `${formatMetricValue(successfulOperationTotal)} / ${formatMetricValue(operationTotal)} 次操作成功`,
      label: '操作成功率',
      percent: operationPercent,
      tone: failedOperationTotal > 0 ? 'warning' : 'success',
      value: operationTotal > 0 ? `${operationPercent}%` : '--',
    },
    {
      detail: dashboard
        ? `数据库 ${dashboard.health.database} / Redis ${dashboard.health.redis}`
        : '等待健康检查同步',
      label: '依赖在线率',
      percent: dependencyPercent,
      tone: percentTone(dependencyPercent),
      value: dashboard ? `${dependencyPercent}%` : '--',
    },
  ]
}

export function dashboardOperationSegments(dashboard: DashboardData | null): DashboardChartSegment[] {
  const metrics = dashboard?.metrics
  const total = metrics?.today_operation_total ?? 0
  const failed = metrics?.today_risk_operation_total ?? 0
  const success = Math.max(total - failed, 0)

  if (total <= 0) {
    return [{ label: '暂无操作', percent: 100, tone: 'default', value: '0' }]
  }

  return [
    {
      label: '成功',
      percent: percentValue(success, total),
      tone: 'success',
      value: formatMetricValue(success),
    },
    {
      label: '失败',
      percent: Math.max(percentValue(failed, total), failed > 0 ? 4 : 0),
      tone: failed > 0 ? 'error' : 'default',
      value: formatMetricValue(failed),
    },
  ]
}

export function dashboardLoginSegments(dashboard: DashboardData | null): DashboardChartSegment[] {
  const logins = dashboard?.recent_logins ?? []
  const success = logins.filter((item) => item.status === DashboardLoginStatus.Success).length
  const failed = logins.length - success

  if (logins.length <= 0) {
    return [{ label: '暂无记录', percent: 100, tone: 'default', value: '0' }]
  }

  return [
    {
      label: '成功',
      percent: percentValue(success, logins.length),
      tone: 'success',
      value: formatMetricValue(success),
    },
    {
      label: '失败',
      percent: Math.max(percentValue(failed, logins.length), failed > 0 ? 6 : 0),
      tone: failed > 0 ? 'warning' : 'default',
      value: formatMetricValue(failed),
    },
  ]
}

export function dashboardLatencyBars(dashboard: DashboardData | null): DashboardBarChartItem[] {
  const rows = [...(dashboard?.recent_operations ?? [])]
    .sort((left, right) => right.latency_ms - left.latency_ms)
    .slice(0, 6)
  const maxLatency = Math.max(...rows.map((item) => item.latency_ms), 1)

  return rows.map((item) => ({
    detail: `${item.method} · ${item.success ? '成功' : '失败'} · ${item.status_code}`,
    label: formatCompactRouteLabel(item.path),
    percent: Math.max(Math.round((item.latency_ms / maxLatency) * 100), 6),
    tone: dashboardLatencyTone(item.latency_ms),
    value: `${item.latency_ms} ms`,
  }))
}

export function dashboardModuleBars(dashboard: DashboardData | null): DashboardBarChartItem[] {
  const counts = new Map<string, number>()
  for (const item of dashboard?.recent_operations ?? []) {
    const label = dashboardModuleLabel(item.path)
    counts.set(label, (counts.get(label) ?? 0) + 1)
  }

  const rows = Array.from(counts.entries())
    .sort((left, right) => right[1] - left[1])
    .slice(0, 5)
  const max = Math.max(...rows.map(([, value]) => value), 1)

  return rows.map(([label, value]) => ({
    detail: '按最近操作记录聚合',
    label,
    percent: Math.max(Math.round((value / max) * 100), 10),
    tone: 'info',
    value: `${formatMetricValue(value)} 次`,
  }))
}

export function dashboardResourceBars(dashboard: DashboardData | null): DashboardBarChartItem[] {
  const metrics = dashboard?.metrics
  const rows = [
    {
      label: '文件',
      tone: 'info' as DashboardTone,
      value: metrics?.file_total ?? 0,
    },
    {
      label: '配置',
      tone: 'warning' as DashboardTone,
      value: metrics?.config_total ?? 0,
    },
    {
      label: '公告',
      tone: 'success' as DashboardTone,
      value: metrics?.notice_total ?? 0,
    },
    {
      label: '角色',
      tone: 'default' as DashboardTone,
      value: metrics?.enabled_role_total ?? 0,
    },
  ]
  const max = Math.max(...rows.map((item) => item.value), 1)

  return rows.map((item) => ({
    detail: '当前启用规模',
    label: item.label,
    percent: Math.max(Math.round((item.value / max) * 100), item.value > 0 ? 8 : 0),
    tone: item.tone,
    value: formatMetricValue(item.value),
  }))
}

export function dashboardChartStats(
  dashboard: DashboardData | null,
  visiblePageTotal: number,
): DashboardChartStat[] {
  const metrics = dashboard?.metrics

  return [
    {
      label: '今日操作',
      tone: (metrics?.today_risk_operation_total ?? 0) > 0 ? 'warning' : 'success',
      value: formatMetricValue(metrics?.today_operation_total),
    },
    {
      label: '风险操作',
      tone: (metrics?.today_risk_operation_total ?? 0) > 0 ? 'error' : 'success',
      value: formatMetricValue(metrics?.today_risk_operation_total),
    },
    {
      label: '登录失败',
      tone: (metrics?.today_login_failed_total ?? 0) > 0 ? 'warning' : 'success',
      value: formatMetricValue(metrics?.today_login_failed_total),
    },
    {
      label: '可访问页',
      tone: 'info',
      value: formatMetricValue(visiblePageTotal),
    },
  ]
}

export function dashboardQuickLinks(menus: AuthMenu[]): QuickLink[] {
  return menus
    .filter((menu) => menu.path && menu.path !== '/dashboard')
    .slice(0, 6)
    .map((menu) => ({
      title: menu.title,
      path: menu.path,
      description: quickLinkDescriptionMap[menu.title] || '进入对应系统页面继续处理业务',
    }))
}

export function dashboardFocusFacts(
  dashboard: DashboardData | null,
  visiblePageTotal: number,
  refreshedLabel: string,
): FocusFact[] {
  return [
    {
      label: '可访问页面',
      value: formatMetricValue(visiblePageTotal),
      hint: '按当前角色实时计算',
    },
    {
      label: '失败登录',
      value: formatMetricValue(dashboard?.metrics.today_login_failed_total),
      hint: '今日累计失败次数',
    },
    {
      label: '最近刷新',
      value: refreshedLabel,
      hint: '已同步当前项目快照',
    },
  ]
}

export function dashboardLatencyTone(latencyMs: number): DashboardTone {
  if (latencyMs >= 1000) {
    return 'warning'
  }

  if (latencyMs >= 300) {
    return 'info'
  }

  return 'success'
}

export function dashboardLatencyLabel(latencyMs: number) {
  if (latencyMs >= 1000) {
    return '偏慢'
  }

  if (latencyMs >= 300) {
    return '正常'
  }

  return '迅速'
}

export function dashboardErrorMessage(error: unknown) {
  if (typeof error === 'object' && error !== null) {
    const response = (error as { response?: { data?: { message?: string } } }).response
    if (typeof response?.data?.message === 'string' && response.data.message) {
      return response.data.message
    }
  }

  return '工作台数据获取失败，请稍后重试。'
}

function percentValue(value: number, total: number) {
  if (total <= 0) {
    return 0
  }

  return Math.min(Math.round((value / total) * 100), 100)
}

function formatPercent(value: number, total: number) {
  return total > 0 ? `${percentValue(value, total)}%` : '--'
}

function dashboardDependencyPercent(dashboard: DashboardData | null) {
  if (!dashboard) {
    return 0
  }

  const checks = [dashboard.health.database, dashboard.health.redis]
  const passed = checks.filter((item) => item === 'ok').length
  return percentValue(passed, checks.length)
}

function dashboardSafetyScore(dashboard: DashboardData | null) {
  const metrics = dashboard?.metrics
  const dependencyPenalty = 100 - dashboardDependencyPercent(dashboard)
  const failedOperationPenalty = Math.min((metrics?.today_risk_operation_total ?? 0) * 8, 36)
  const failedLoginPenalty = Math.min((metrics?.today_login_failed_total ?? 0) * 6, 30)

  return Math.max(100 - dependencyPenalty - failedOperationPenalty - failedLoginPenalty, 0)
}

function percentTone(percent: number): DashboardTone {
  if (percent >= 80) {
    return 'success'
  }

  if (percent >= 55) {
    return 'warning'
  }

  return 'error'
}

function scoreTone(score: number): DashboardTone {
  if (score >= 86) {
    return 'success'
  }

  if (score >= 70) {
    return 'warning'
  }

  return 'error'
}

function dashboardModuleLabel(path: string) {
  const parts = formatDashboardRoutePath(path).split('/').filter(Boolean)
  const moduleKey = parts[1] ?? parts[0] ?? 'other'
  const labelMap: Record<string, string> = {
    attachments: '附件',
    configs: '配置',
    departments: '部门',
    files: '文件',
    login: '登录',
    loginlog: '登录日志',
    logs: '日志',
    menus: '菜单',
    notices: '公告',
    operation: '操作日志',
    operationlogs: '操作日志',
    posts: '岗位',
    roles: '角色',
    users: '用户',
  }

  return labelMap[moduleKey.replace(/-/g, '').toLowerCase()] ?? moduleKey
}

function formatCompactRouteLabel(path: string) {
  const normalized = formatDashboardRoutePath(path)
  const parts = normalized.split('/').filter(Boolean)
  if (parts.length <= 3) {
    return normalized || path
  }

  return `/${parts.slice(0, 3).join('/')}`
}

function formatDashboardTimeLabel(value: string) {
  return new Intl.DateTimeFormat('zh-CN', {
    hour: '2-digit',
    minute: '2-digit',
  }).format(new Date(value))
}
