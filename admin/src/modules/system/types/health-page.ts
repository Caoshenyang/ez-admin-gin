import type { Component } from 'vue'

export type HealthServiceKey = 'database' | 'redis'
export type HealthTagType = 'default' | 'error' | 'info' | 'success' | 'warning'

export interface HealthDependencyCard {
  detail: string
  icon: Component
  key: HealthServiceKey
  label: string
  progress: number
  progressLabel: string
  statusLabel: string
  tagType: HealthTagType
  toneClass: string
  value?: string
  description: string
}

export interface HealthEndpointCard {
  badge: string
  description: string
  icon: Component
  method: string
  path: string
  title: string
  toneClass: string
}

export interface HealthOverview {
  description: string
  statusLabel: string
  tagType: HealthTagType
  title: string
  toneClass: string
}

export interface HealthSummaryCard {
  detail: string
  icon: Component
  key: string
  label: string
  tagLabel?: string
  tagType?: HealthTagType
  toneClass: string
  value: string
}

export interface HealthCheckItem {
  description: string
  key: string
  label: string
  statusLabel: string
  tagType: HealthTagType
  toneClass: string
}
