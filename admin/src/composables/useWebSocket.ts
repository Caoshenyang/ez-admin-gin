import { ref, watch } from 'vue'
import { useThemeStore } from '@/stores/theme'

export interface WSOptions {
  onMessage: (data: unknown) => void
  onOpen?: () => void
  onClose?: () => void
}

export function useWebSocket(getUrl: () => string, opts: WSOptions) {
  const connected = ref(false)
  let ws: WebSocket | null = null
  let reconnectTimer: ReturnType<typeof setTimeout> | null = null
  let heartbeatTimer: ReturnType<typeof setInterval> | null = null
  let stopped = false

  const themeStore = useThemeStore()

  function connect() {
    if (stopped) return

    const url = getUrl()
    if (!url) return

    ws = new WebSocket(url)

    ws.onopen = () => {
      connected.value = true
      startHeartbeat()
      opts.onOpen?.()
    }

    ws.onmessage = (event) => {
      try {
        const msg = JSON.parse(event.data)
        if (msg.type === 'ping') {
          ws?.send(JSON.stringify({ type: 'pong' }))
          return
        }
        opts.onMessage(msg)
      } catch {
        // ignore non-JSON messages
      }
    }

    ws.onclose = () => {
      connected.value = false
      stopHeartbeat()
      opts.onClose?.()
      if (!stopped) {
        scheduleReconnect()
      }
    }

    ws.onerror = () => {
      ws?.close()
    }
  }

  function disconnect() {
    stopped = true
    stopHeartbeat()
    clearReconnect()
    if (ws) {
      ws.onclose = null
      ws.close()
      ws = null
    }
    connected.value = false
  }

  function send(data: unknown) {
    if (ws && ws.readyState === WebSocket.OPEN) {
      ws.send(JSON.stringify(data))
    }
  }

  function startHeartbeat() {
    stopHeartbeat()
    heartbeatTimer = setInterval(() => {
      // Server sends ping, we respond with pong (handled in onmessage)
    }, 25000)
  }

  function stopHeartbeat() {
    if (heartbeatTimer) {
      clearInterval(heartbeatTimer)
      heartbeatTimer = null
    }
  }

  function scheduleReconnect() {
    clearReconnect()
    reconnectTimer = setTimeout(() => {
      connect()
    }, 3000)
  }

  function clearReconnect() {
    if (reconnectTimer) {
      clearTimeout(reconnectTimer)
      reconnectTimer = null
    }
  }

  watch(() => themeStore.mode, () => {
    // Reconnect when theme changes to refresh state if needed
  })

  return {
    connected,
    connect,
    disconnect,
    send,
  }
}
