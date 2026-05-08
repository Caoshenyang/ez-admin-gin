import type { Component } from 'vue'

export interface MetricCard {
  label: string
  value: string
  hint: string
  accent: string
  iconBg: string
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
