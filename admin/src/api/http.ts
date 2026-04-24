import axios from 'axios'

import { clearAuthSession, getAuthorizationHeader } from '../utils/auth'

const http = axios.create({
  // 通过 Vite 代理转发到本地后端。
  baseURL: '/api/v1',
  timeout: 10000,
})

http.interceptors.request.use((config) => {
  const authorization = getAuthorizationHeader()

  if (authorization) {
    config.headers.Authorization = authorization
  }

  return config
})

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
