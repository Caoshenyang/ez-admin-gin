import axios from 'axios'

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

// 响应拦截器：正常响应直接返回；401 时清除本地登录态并向上抛出错误。
http.interceptors.response.use(
  (response) => response,
  (error) => {
    // 后面做完整登录态守卫前，先在 401 时清掉本地旧 Token。
    if (error.response?.status === 401) {
      clearAuthSession()
    }

    return Promise.reject(error)
  },
)

export default http
