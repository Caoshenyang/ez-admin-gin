import axios from 'axios'
import type { AxiosError } from 'axios'

import router from '../router'
import { clearAuthSession, getAuthorizationHeader } from '../utils/auth'

// 不需要 Authorization 请求头的公开接口路径。
const publicApiPaths = new Set(['/auth/login', '/setup/init'])

// http 预配置的 Axios 实例，统一管理 baseURL、超时和拦截器。
const http = axios.create({
  // 通过 Vite 代理转发到本地后端。
  baseURL: '/api/v1',
  // timeout 请求超时时间（毫秒）。
  timeout: 10000,
})

// 请求拦截器：为非公开接口自动注入 Authorization 请求头。
http.interceptors.request.use((config) => {
  const requestPath = config.url ?? ''
  if (publicApiPaths.has(requestPath)) {
    return config
  }

  const authorization = getAuthorizationHeader()

  if (authorization) {
    config.headers.Authorization = authorization
  }

  return config
})

// 记录是否正在跳转登录页，防止多个 401 并发时重复跳转。
let isRedirectingToLogin = false

// 响应拦截器：统一处理 401、403、网络错误和 5xx。
http.interceptors.response.use(
  (response) => response,
  (error: AxiosError) => {
    if (!error.response) {
      // 网络异常（断网 / 超时）
      showMessage('网络异常，请检查网络连接后重试')
      return Promise.reject(error)
    }

    const { status, config } = error.response

    if (status === 401) {
      // 登录接口本身返回 401 时不跳转，只向上抛出。
      if (config.url === '/auth/login') {
        return Promise.reject(error)
      }

      clearAuthSession()

      if (!isRedirectingToLogin) {
        isRedirectingToLogin = true
        void router.push('/login').finally(() => {
          isRedirectingToLogin = false
        })
      }

      return Promise.reject(error)
    }

    if (status === 403) {
      showMessage('无权限访问该资源')
      return Promise.reject(error)
    }

    if (status >= 500) {
      showMessage('服务器错误，请稍后重试')
      return Promise.reject(error)
    }

    return Promise.reject(error)
  },
)

// 简单的消息展示机制，避免在 http 模块中直接依赖 Naive UI 实例。
// 由 main.ts 中通过 setMessageHandler 注入实际的 NMessage 调用。
let messageHandler: ((msg: string) => void) | null = null

export function setMessageHandler(handler: (msg: string) => void) {
  messageHandler = handler
}

function showMessage(msg: string) {
  if (messageHandler) {
    messageHandler(msg)
  }
}

export default http
