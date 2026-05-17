import { ref } from 'vue'

export interface WSOptions {
  onMessage: (data: unknown) => void
  onOpen?: () => void
  onClose?: () => void
}

export function useWebSocket(getUrl: () => string, opts: WSOptions) {
  const connected = ref(false)
  let ws: WebSocket | null = null
  let reconnectTimer: ReturnType<typeof setTimeout> | null = null
  let stopped = false

  function connect() {
    if (stopped) return

    const url = getUrl()
    if (!url) return

    ws = new WebSocket(url)

    ws.onopen = () => {
      connected.value = true
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

  return {
    connected,
    connect,
    disconnect,
    send,
  }
}
