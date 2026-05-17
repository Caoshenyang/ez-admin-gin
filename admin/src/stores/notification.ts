import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getNotifications, getUnreadCount, markRead, markAllRead } from '@/modules/system/api/notification'
import { useWebSocket } from '@/composables/useWebSocket'
import { getAccessToken } from '@/utils/auth'
import type { NotificationItem, NotificationListQuery, WSMessage } from '@/types/notification'

function buildWSUrl(): string {
  const protocol = location.protocol === 'https:' ? 'wss:' : 'ws:'
  const token = getAccessToken()
  if (!token) return ''
  return `${protocol}//${location.host}/api/v1/system/notifications/ws?token=${token}`
}

export const useNotificationStore = defineStore('notification', () => {
  const items = ref<NotificationItem[]>([])
  const total = ref(0)
  const unreadCount = ref(0)
  const loading = ref(false)
  const drawerVisible = ref(false)

  let wsHandle: ReturnType<typeof useWebSocket> | null = null

  async function loadNotifications(query?: Partial<NotificationListQuery>) {
    loading.value = true
    try {
      const data = await getNotifications({
        page: query?.page ?? 1,
        page_size: query?.page_size ?? 20,
        type: query?.type,
        is_read: query?.is_read,
      })
      items.value = data.items
      total.value = data.total
    } finally {
      loading.value = false
    }
  }

  async function fetchUnreadCount() {
    const data = await getUnreadCount()
    if (data) {
      unreadCount.value = data.count
    }
  }

  async function handleMarkRead(ids: number[]) {
    await markRead({ ids })
    await fetchUnreadCount()
  }

  async function handleMarkAllRead() {
    await markAllRead()
    await fetchUnreadCount()
  }

  function openDrawer() {
    drawerVisible.value = true
    void loadNotifications()
  }

  function closeDrawer() {
    drawerVisible.value = false
  }

  function handleWSMessage(data: unknown) {
    const msg = data as WSMessage

    if (msg.type === 'notification') {
      const item = msg.data as NotificationItem
      items.value.unshift(item)
      unreadCount.value++
    } else if (msg.type === 'unread_count') {
      const data = msg.data as { count: number }
      unreadCount.value = data.count
    }
  }

  function connectWS() {
    if (wsHandle) return

    wsHandle = useWebSocket(buildWSUrl, {
      onMessage: handleWSMessage,
    })
    wsHandle.connect()
  }

  function disconnectWS() {
    wsHandle?.disconnect()
    wsHandle = null
  }

  function wsSend(data: unknown) {
    wsHandle?.send(data)
  }

  function reset() {
    items.value = []
    total.value = 0
    unreadCount.value = 0
    loading.value = false
    drawerVisible.value = false
    disconnectWS()
  }

  return {
    items,
    total,
    unreadCount,
    loading,
    drawerVisible,
    loadNotifications,
    fetchUnreadCount,
    handleMarkRead,
    handleMarkAllRead,
    openDrawer,
    closeDrawer,
    connectWS,
    disconnectWS,
    wsSend,
    reset,
  }
})
