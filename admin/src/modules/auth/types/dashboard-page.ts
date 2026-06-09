import type { Component } from 'vue'

export type DashboardTone = 'success' | 'warning' | 'error' | 'info' | 'default'

export interface MetricCard {
  label: string
  value: string
  hint: string
  iconClass: string
  panelClass: string
  icon: Component
}

export interface QuickLink {
  title: string
  path: string
  description: string
}

export interface HealthSummaryItem {
  label: string
  value: string
  status: string
}

export interface FocusFact {
  label: string
  value: string
  hint: string
}

export interface DashboardCommandItem {
  actionText: string
  description: string
  icon: Component
  path: string
  tagLabel: string
  tagType: DashboardTone
  title: string
  value: string
}

export interface DashboardInsightCard {
  description: string
  icon: Component
  label: string
  tone: DashboardTone
  value: string
}

export interface DashboardResourceItem {
  detail: string
  label: string
  percent: number
  tone: DashboardTone
  value: string
}

export interface DashboardRingMetric {
  detail: string
  label: string
  percent: number
  tone: DashboardTone
  value: string
}

export interface DashboardChartSegment {
  label: string
  percent: number
  tone: DashboardTone
  value: string
}

export interface DashboardBarChartItem {
  detail: string
  label: string
  percent: number
  tone: DashboardTone
  value: string
}

export interface DashboardChartStat {
  label: string
  tone: DashboardTone
  value: string
}

export interface DashboardTrendPoint {
  label: string
  latency: number
  name: string
  risk: number
  success: boolean
}
