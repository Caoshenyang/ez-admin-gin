import type { FormRules } from 'naive-ui'

import type { AccountDataScopeSummary, AccountProfileResponse } from '../types/auth'
import type { PasswordFormModel, ProfileFormModel } from '../types/account-center-page'

export function defaultProfileFormModel(): ProfileFormModel {
  return {
    nickname: '',
  }
}

export function defaultPasswordFormModel(): PasswordFormModel {
  return {
    oldPassword: '',
    newPassword: '',
    confirmPassword: '',
  }
}

export const profileFormRules: FormRules = {
  nickname: [
    { required: true, message: '请输入昵称', trigger: ['blur', 'input'] },
    { max: 64, message: '昵称不能超过 64 个字符', trigger: ['blur', 'input'] },
  ],
}

export function passwordFormRules(formModel: PasswordFormModel): FormRules {
  return {
    oldPassword: [{ required: true, message: '请输入当前密码', trigger: ['blur', 'input'] }],
    newPassword: [
      { required: true, message: '请输入新密码', trigger: ['blur', 'input'] },
      { min: 8, message: '新密码至少 8 位', trigger: ['blur', 'input'] },
      { max: 72, message: '新密码不能超过 72 位', trigger: ['blur', 'input'] },
    ],
    confirmPassword: [
      { required: true, message: '请再次输入新密码', trigger: ['blur', 'input'] },
      {
        validator: () => {
          if (formModel.confirmPassword !== formModel.newPassword) {
            return new Error('两次输入的新密码不一致')
          }

          return true
        },
        trigger: ['blur', 'input'],
      },
    ],
  }
}

export function roleSummaryText(profile: AccountProfileResponse | null) {
  if (!profile || profile.role_codes.length === 0) {
    return '未绑定角色'
  }

  return profile.role_codes.join(' / ')
}

export function dataScopeSummaryText(summary?: AccountDataScopeSummary) {
  if (!summary) {
    return '加载中'
  }
  if (summary.allow_all) {
    return '全部数据'
  }
  if (summary.require_self) {
    return '仅本人数据'
  }
  if (summary.include_dept_tree) {
    return '本部门及下级部门数据'
  }
  if (summary.include_department) {
    return '本部门数据'
  }
  if (summary.custom_department_ids.length > 0) {
    return `自定义部门范围（${summary.custom_department_ids.length} 个部门）`
  }

  return '未配置数据范围'
}
