import axios from 'axios'
import type { FormRules } from 'naive-ui'

import type { LoginFormModel } from '../types/login-page'

const captchaAlphabet = 'ABCDEFGHJKLMNPQRSTUVWXYZ23456789'

export const loginProductFeatures = [
  '权限模型：用户 / 角色 / 菜单 / 按钮',
  '工作标签：多页面切换、刷新、关闭其他',
  '审计能力：登录日志、操作日志、风险等级',
  '蓝图友好：模块骨架清晰，业务页面可以快速扩展',
]

export function createCaptchaText() {
  return Array.from({ length: 4 }, () => {
    const index = Math.floor(Math.random() * captchaAlphabet.length)
    return captchaAlphabet[index]
  }).join('')
}

export function defaultLoginFormModel(): LoginFormModel {
  return {
    username: 'admin',
    password: 'Admin@123456',
    captcha: '',
    rememberLogin: true,
  }
}

export const loginFormRules: FormRules = {
  username: [
    {
      required: true,
      message: '请输入用户名',
      trigger: ['blur', 'input'],
    },
  ],
  password: [
    {
      required: true,
      message: '请输入密码',
      trigger: ['blur', 'input'],
    },
  ],
}

export function loginFooterText(year = new Date().getFullYear()) {
  return `© ${year} EZ Admin Gin · 高效 · 易用 · 企业级后台框架`
}

export function loginErrorMessage(error: unknown) {
  return axios.isAxiosError<{ message?: string }>(error)
    ? error.response?.data?.message ?? '登录失败，请稍后重试'
    : '登录失败，请稍后重试'
}
