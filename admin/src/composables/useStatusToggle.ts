import { useMessage } from 'naive-ui'

interface StatusItem {
  id: number
}

interface StatusPayload<S extends number = number> {
  status: S
}

export function useStatusToggle<
  T extends StatusItem,
  P extends StatusPayload,
>(
  toggleFn: (id: number, payload: P) => Promise<unknown>,
  options?: { onSuccess?: () => Promise<void> | void },
) {
  const message = useMessage()

  async function handleToggleStatus(row: T, nextStatus: P['status']) {
    await toggleFn(row.id, { status: nextStatus } as P)
    message.success(nextStatus === 1 ? '已启用' : '已禁用')
    await options?.onSuccess?.()
  }

  return { handleToggleStatus }
}
