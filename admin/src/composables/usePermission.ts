import { computed } from 'vue'
import { buttonPermissionCodes } from '@/router/dynamic-menu'

// usePermission 提供按钮级权限判断能力。
export function usePermission() {
  const codes = computed(() => buttonPermissionCodes.value)

  // canUse 判断当前用户是否拥有指定按钮权限码。
  function canUse(code: string): boolean {
    return codes.value.includes(code)
  }

  return { canUse, buttonPermissionCodes: codes }
}
