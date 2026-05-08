import type { LoginResponse } from '@/modules/auth/types/auth'

const ACCESS_TOKEN_KEY = 'ez-admin-access-token'
const TOKEN_TYPE_KEY = 'ez-admin-token-type'
const USER_INFO_KEY = 'ez-admin-user-info'
export const AUTH_USER_INFO_UPDATED_EVENT = 'ez-admin-auth-user-info-updated'

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

// readStorageValue 同时查找 localStorage 和 sessionStorage，支持"记住登录"和"不记住"两种存储。
function readStorageValue(key: string) {
  return localStorage.getItem(key) ?? sessionStorage.getItem(key) ?? ''
}

function currentUserInfoStorage() {
  if (localStorage.getItem(USER_INFO_KEY) !== null) {
    return localStorage
  }
  if (sessionStorage.getItem(USER_INFO_KEY) !== null) {
    return sessionStorage
  }

  return null
}

function notifyAuthUserInfoUpdated() {
  window.dispatchEvent(new Event(AUTH_USER_INFO_UPDATED_EVENT))
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

  notifyAuthUserInfoUpdated()
}

export function clearAuthSession() {
  localStorage.removeItem(ACCESS_TOKEN_KEY)
  localStorage.removeItem(TOKEN_TYPE_KEY)
  localStorage.removeItem(USER_INFO_KEY)

  sessionStorage.removeItem(ACCESS_TOKEN_KEY)
  sessionStorage.removeItem(TOKEN_TYPE_KEY)
  sessionStorage.removeItem(USER_INFO_KEY)

  notifyAuthUserInfoUpdated()
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

// updateAuthUserInfo 局部更新本地登录人摘要，用于昵称等自助资料修改后的同步。
export function updateAuthUserInfo(patch: Partial<AuthUserInfo>) {
  const storage = currentUserInfoStorage()
  const current = getAuthUserInfo()
  if (!storage || !current) {
    return
  }

  storage.setItem(
    USER_INFO_KEY,
    JSON.stringify({
      ...current,
      ...patch,
    } satisfies AuthUserInfo),
  )

  notifyAuthUserInfoUpdated()
}

// getAuthorizationHeader 统一拼接 Authorization 请求头。
export function getAuthorizationHeader() {
  const accessToken = getAccessToken()
  if (!accessToken) {
    return ''
  }

  return `${getTokenType()} ${accessToken}`
}
