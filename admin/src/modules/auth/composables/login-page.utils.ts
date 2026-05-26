import axios from 'axios'
import type { FormRules } from 'naive-ui'

import type { LoginFormModel } from '../types/login-page'

const captchaAlphabet = 'ABCDEFGHJKLMNPQRSTUVWXYZ23456789'

export const loginProductFeatures = [
  'RBAC 权限体系',
  '动态菜单与按钮权限',
  '五级数据权限',
  '系统日志与公告通知',
  '高性能 WebSocket 推送',
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
    password: 'EzAdmin@123456',
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
