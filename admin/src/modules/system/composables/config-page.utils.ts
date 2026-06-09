import {
  ConfigStatus,
  type CreateConfigPayload,
  type UpdateConfigPayload,
  type ConfigItem,
  type ConfigListQuery,
} from '../types/config'

export interface ConfigCategory {
  key: string
  group_code: string
  label: string
  description: string
}

export interface ConfigFormModel {
  id: number
  group_code: string
  key: string
  name: string
  value: string
  sort: number
  status: ConfigStatus
  remark: string
}

export const BUILTIN_CONFIG_CATEGORIES: ConfigCategory[] = [
  { key: 'all', group_code: '', label: '全部配置', description: '查看所有系统配置项' },
  {
    key: 'rate_limit',
    group_code: 'rate_limit',
    label: '限流配置',
    description: '登录 IP 限流、账号失败锁定',
  },
  {
    key: 'upload',
    group_code: 'upload',
    label: '上传配置',
    description: '上传大小、文件类型、存储策略',
  },
]

export function defaultConfigListQuery(): ConfigListQuery {
  return {
    page: 1,
    page_size: 10,
    keyword: '',
    group_code: '',
    status: 0,
  }
}

export function defaultConfigFormModel(): ConfigFormModel {
  return {
    id: 0,
    group_code: '',
    key: '',
    name: '',
    value: '',
    sort: 0,
    status: ConfigStatus.Enabled,
    remark: '',
  }
}

export function toConfigFormModel(config: ConfigItem): ConfigFormModel {
  return {
    id: config.id,
    group_code: config.group_code,
    key: config.key,
    name: config.name,
    value: config.value,
    sort: config.sort,
    status: config.status,
    remark: config.remark,
  }
}

// 配置更新不允许修改 key，所以单独拆出创建/更新 payload，避免页面层再手工挑字段。
export function buildConfigCreatePayload(formModel: ConfigFormModel): CreateConfigPayload {
  return {
    group_code: formModel.group_code,
    key: formModel.key,
    name: formModel.name,
    value: formModel.value,
    sort: formModel.sort,
    status: formModel.status,
    remark: formModel.remark,
  }
}

export function buildConfigUpdatePayload(formModel: ConfigFormModel): UpdateConfigPayload {
  return {
    group_code: formModel.group_code,
    name: formModel.name,
    value: formModel.value,
    sort: formModel.sort,
    status: formModel.status,
    remark: formModel.remark,
  }
}
