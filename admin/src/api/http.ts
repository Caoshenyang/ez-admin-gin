import axios from 'axios'
import type { AxiosError, InternalAxiosRequestConfig } from 'axios'

import router from '../router'
import { clearAuthSession, getAuthorizationHeader, refreshAccessToken } from '../utils/auth'

const publicApiPaths = new Set(['/auth/login', '/setup/init', '/auth/refresh'])

const http = axios.create({
  baseURL: '/api/v1',
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

let isRedirectingToLogin = false

// Pending refresh promise — ensures only one refresh request at a time.
let pendingRefresh: Promise<string> | null = null

function getOrStartRefresh(): Promise<string> {
  if (pendingRefresh) {
    return pendingRefresh
  }
  pendingRefresh = refreshAccessToken().finally(() => {
    pendingRefresh = null
  })
  return pendingRefresh
}

// 响应拦截器：统一处理 401（含自动刷新）、403、网络错误和 5xx。
http.interceptors.response.use(
  (response) => response,
  async (error: AxiosError) => {
    if (!error.response) {
      showMessage('网络异常，请检查网络连接后重试')
      return Promise.reject(error)
    }

    const { status, config } = error.response

    if (status === 401) {
      // Login and refresh endpoints returning 401 should not trigger refresh.
      if (config.url === '/auth/login' || config.url === '/auth/refresh') {
        return Promise.reject(error)
      }

      // Try to refresh the access token.
      const newToken = await getOrStartRefresh()
      if (newToken) {
        // Retry the original request with the new token.
        const retryConfig: InternalAxiosRequestConfig = {
          ...config,
          headers: { ...config.headers, Authorization: `Bearer ${newToken}` },
        }
        return http.request(retryConfig)
      }

      // Refresh failed — clear session and redirect to login.
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
