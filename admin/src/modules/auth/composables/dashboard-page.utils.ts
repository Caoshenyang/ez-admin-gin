import {
  LayersOutline,
  PulseOutline,
  ShieldCheckmarkOutline,
  TimeOutline,
} from '@vicons/ionicons5'

import { DashboardLoginStatus, type DashboardData } from '../types/dashboard'
import type {
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

export function dashboardCurrentDateLabel(date = new Date()) {
  return new Intl.DateTimeFormat('zh-CN', {
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
    return loading ? '正在拉取当前项目的实时概览...' : '等待首次同步项目运行数据'
  }

  if (health.database === 'ok' && health.redis === 'ok') {
    return '环境、数据库和缓存依赖都在线，可以直接作为管理员首页使用。'
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
  return value ? new Date(value).toLocaleString() : '-'
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
      accent: '#2563EB',
      iconBg: '#EFF6FF',
      panelClass: '',
      icon: ShieldCheckmarkOutline,
    },
    {
      label: '启用角色',
      value: formatMetricValue(metrics?.enabled_role_total),
      hint: `可访问页面 ${formatMetricValue(visiblePageTotal)}`,
      accent: '#8B5CF6',
      iconBg: '#F5F3FF',
      panelClass: '',
      icon: LayersOutline,
    },
    {
      label: '今日操作',
      value: formatMetricValue(metrics?.today_operation_total),
      hint: `失败 ${formatMetricValue(metrics?.today_risk_operation_total)}`,
      accent: '#F59E0B',
      iconBg: '#FFFBEB',
      panelClass: '',
      icon: TimeOutline,
    },
    {
      label: '文件沉淀',
      value: formatMetricValue(metrics?.file_total),
      hint: `公告 ${formatMetricValue(metrics?.notice_total)} / 配置 ${formatMetricValue(metrics?.config_total)}`,
      accent: '#22C55E',
      iconBg: '#F0FDF4',
      panelClass: '',
      icon: PulseOutline,
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

export function dashboardErrorMessage(error: unknown) {
  if (typeof error === 'object' && error !== null) {
    const response = (error as { response?: { data?: { message?: string } } }).response
    if (typeof response?.data?.message === 'string' && response.data.message) {
      return response.data.message
    }
  }

  return '工作台数据获取失败，请稍后重试。'
}
