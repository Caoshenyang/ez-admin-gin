import { useMessage } from 'naive-ui'

interface StatusItem {
  id: number
  status: number
}

export function useStatusToggle<T extends StatusItem>(
  toggleFn: (id: number, payload: { status: number }) => Promise<unknown>,
  options?: { onSuccess?: () => Promise<void> | void },
) {
  const message = useMessage()

  async function handleToggleStatus(row: T, nextStatus: number) {
    await toggleFn(row.id, { status: nextStatus })
    message.success(nextStatus === 1 ? '已启用' : '已禁用')
    await options?.onSuccess?.()
  }

  return { handleToggleStatus }
}
