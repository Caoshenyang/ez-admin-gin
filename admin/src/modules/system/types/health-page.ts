export type HealthServiceKey = 'database' | 'redis'

export interface HealthDependencyCard {
  key: HealthServiceKey
  label: string
  value?: string
  description: string
}

export interface HealthEndpointCard {
  title: string
  path: string
  description: string
}
