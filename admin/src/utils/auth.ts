import type { LoginResponse } from '../types/auth'

const ACCESS_TOKEN_KEY = 'ez-admin-access-token'
const TOKEN_TYPE_KEY = 'ez-admin-token-type'
const USER_INFO_KEY = 'ez-admin-user-info'

type StorageMode = 'local' | 'session'

export interface AuthUserInfo {
  userId: number
  username: string
  nickname: string
  expiresAt: string
}

function getStorage(mode: StorageMode) {
  return mode === 'local' ? localStorage : sessionStorage
}

function readStorageValue(key: string) {
  return localStorage.getItem(key) ?? sessionStorage.getItem(key) ?? ''
}

// setAuthSession 在登录成功后保存本地登录态。
export function setAuthSession(payload: LoginResponse, rememberLogin: boolean) {
  clearAuthSession()

  const storage = getStorage(rememberLogin ? 'local' : 'session')

  storage.setItem(ACCESS_TOKEN_KEY, payload.access_token)
  storage.setItem(TOKEN_TYPE_KEY, payload.token_type)
  storage.setItem(
    USER_INFO_KEY,
    JSON.stringify({
      userId: payload.user_id,
      username: payload.username,
      nickname: payload.nickname,
      expiresAt: payload.expires_at,
    } satisfies AuthUserInfo),
  )
}

export function clearAuthSession() {
  localStorage.removeItem(ACCESS_TOKEN_KEY)
  localStorage.removeItem(TOKEN_TYPE_KEY)
  localStorage.removeItem(USER_INFO_KEY)

  sessionStorage.removeItem(ACCESS_TOKEN_KEY)
  sessionStorage.removeItem(TOKEN_TYPE_KEY)
  sessionStorage.removeItem(USER_INFO_KEY)
}

export function getAccessToken() {
  return readStorageValue(ACCESS_TOKEN_KEY)
}

export function getTokenType() {
  return readStorageValue(TOKEN_TYPE_KEY) || 'Bearer'
}

export function hasAccessToken() {
  return getAccessToken() !== ''
}

export function getAuthUserInfo() {
  const raw = readStorageValue(USER_INFO_KEY)
  if (!raw) {
    return null
  }

  try {
    return JSON.parse(raw) as AuthUserInfo
  } catch {
    clearAuthSession()
    return null
  }
}

// getAuthorizationHeader 统一拼接 Authorization 请求头。
export function getAuthorizationHeader() {
  const accessToken = getAccessToken()
  if (!accessToken) {
    return ''
  }

  return `${getTokenType()} ${accessToken}`
}
