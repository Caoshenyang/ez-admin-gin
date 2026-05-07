import { computed } from 'vue'
import { buttonPermissionCodes } from '@/router/dynamic-menu'

export function usePermission() {
  const codes = computed(() => buttonPermissionCodes.value)

  function canUse(code: string): boolean {
    return codes.value.includes(code)
  }

  return { canUse, buttonPermissionCodes: codes }
}
