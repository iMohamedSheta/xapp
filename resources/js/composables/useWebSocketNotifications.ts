import { ref } from 'vue'

import { Notification } from '@/types'
import { formatRelativeTime } from './../lib/format'

export const WS_PATH = '/ws/connect'

const WS_PORT = import.meta.env.VITE_WS_PORT || '8081'
const MAX_RECONNECT_ATTEMPTS = 5

export type WsIncoming = Record<string, unknown>
export type WsMessageHandler = (data: WsIncoming) => void

// ─── Internal state ───────────────────────────────────────────────────────────

interface ConnectionState {
  ws: WebSocket | null
  reconnectAttempts: number
  reconnectTimeout: ReturnType<typeof setTimeout> | null
  /** channel → set of handlers registered for that channel */
  channelHandlers: Map<string, Set<WsMessageHandler>>
  /** channels we should re-subscribe to after reconnect */
  pendingSubscriptions: Set<string>
}

const state: ConnectionState = {
  ws: null,
  reconnectAttempts: 0,
  reconnectTimeout: null,
  channelHandlers: new Map(),
  pendingSubscriptions: new Set(),
}

// ─── Shared reactive refs ─────────────────────────────────────────────────────

const websocketNotifications = ref<Notification[]>([])
const isConnected = ref(false)

// ─── URL helper ───────────────────────────────────────────────────────────────

function wsUrl(): string {
  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
  // direct backend
  return `${protocol}//${window.location.hostname}:${WS_PORT}${WS_PATH}`
}

// ─── Default notification handler ────────────────────────────────────────────

function defaultNotificationHandler(data: WsIncoming) {
  if (data.type !== 'notification') return
  const payload = (data.data || {}) as Notification
  websocketNotifications.value.unshift({
    id: payload.id || Date.now(),
    title: payload.title || 'إشعار جديد',
    message: payload.message || '',
    time: formatRelativeTime(payload.time),
    link: payload.link,
    read: false,
    open: false,
    type: payload.type || 'notification',
  })
}

// ─── Message dispatch ─────────────────────────────────────────────────────────

function dispatchMessage(event: MessageEvent) {
  try {
    const data = JSON.parse(event.data) as WsIncoming
    const channel = data.channel as string | undefined

    if (!channel) return

    // Dispatch to all handlers registered for the exact channel
    const handlers = state.channelHandlers.get(channel)
    if (handlers) {
      handlers.forEach(h => h(data))
    }

    // Also dispatch to wildcard pattern handlers (e.g. "user_notifications.*")
    state.channelHandlers.forEach((handlerSet, registeredChannel) => {
      if (registeredChannel !== channel && matchesPattern(registeredChannel, channel)) {
        handlerSet.forEach(h => h(data))
      }
    })
  } catch (error) {
    console.error('Error parsing WebSocket message:', error)
  }
}

/** Match a wildcard pattern like "user_notifications.*" against a concrete channel name. */
function matchesPattern(pattern: string, channel: string): boolean {
  if (pattern.endsWith('.*')) {
    const prefix = pattern.slice(0, -2)
    return channel.startsWith(prefix + '.')
  }
  return pattern === channel
}

// ─── Socket lifecycle ─────────────────────────────────────────────────────────

function openSocket() {
  if (
    state.ws?.readyState === WebSocket.OPEN ||
    state.ws?.readyState === WebSocket.CONNECTING
  ) {
    return
  }

  try {
    const ws = new WebSocket(wsUrl())
    state.ws = ws

    ws.onopen = () => {
      state.reconnectAttempts = 0
      isConnected.value = true

      // Re-subscribe to all channels that were registered before/during reconnect
      state.pendingSubscriptions.forEach(channel => {
        ws.send(JSON.stringify({ type: 'subscribe', channel }))
      })
    }

    ws.onmessage = event => dispatchMessage(event)

    ws.onclose = () => {
      state.ws = null
      isConnected.value = false
      scheduleReconnect()
    }

    ws.onerror = () => {
      isConnected.value = false
    }
  } catch (error) {
    console.error('Error connecting to WebSocket:', error)
    isConnected.value = false
  }
}

function scheduleReconnect() {
  if (state.reconnectAttempts >= MAX_RECONNECT_ATTEMPTS) return
  const delay = Math.min(1000 * Math.pow(2, state.reconnectAttempts), 30000)
  state.reconnectAttempts++
  state.reconnectTimeout = setTimeout(openSocket, delay)
}

// ─── Public API ───────────────────────────────────────────────────────────────

export interface WsChannelOptions {
  channel: string
  onMessage?: WsMessageHandler
  /** @deprecated path is ignored — all channels share a single connection */
  path?: string
}



export function useWebSocketNotifications() {
  /**
   * Subscribe to a channel over the shared single connection.
   *
   * @example
   * // Notifications (string shorthand)
   * connect(`user_notifications.${userId}`)
   *
   * // Topology tests (options object)
   * connect({ channel: `topology.tests.${userId}`, onMessage: handler })
   */
  function connect(channelOrOptions: string | WsChannelOptions) {
    const options: WsChannelOptions =
      typeof channelOrOptions === 'string'
        ? { channel: channelOrOptions }
        : channelOrOptions

    const { channel } = options

    // Determine the handler to attach
    const handler: WsMessageHandler =
      options.onMessage ?? defaultNotificationHandler

    // Register the handler for this channel
    if (!state.channelHandlers.has(channel)) {
      state.channelHandlers.set(channel, new Set())
    }
    state.channelHandlers.get(channel)!.add(handler)

    // Mark channel for (re-)subscription
    state.pendingSubscriptions.add(channel)

    // If already connected, subscribe immediately
    if (state.ws?.readyState === WebSocket.OPEN) {
      state.ws.send(JSON.stringify({ type: 'subscribe', channel }))
    } else {
      // Connect (or let the onopen handler subscribe when it fires)
      openSocket()
    }
  }

  /**
   * Send a message on a specific channel.
   * Returns true if the socket was open and the message was sent.
   */
  function send(channel: string, data: object): boolean {
    if (state.ws?.readyState === WebSocket.OPEN) {
      state.ws.send(JSON.stringify({ type: 'message', channel, data }))
      return true
    }
    return false
  }

  /**
   * Unsubscribe a handler (or all handlers) from a channel.
   * If no handler is given, all handlers for that channel are removed.
   */
  function unsubscribe(channel: string, handler?: WsMessageHandler) {
    if (handler) {
      state.channelHandlers.get(channel)?.delete(handler)
    } else {
      state.channelHandlers.delete(channel)
    }

    if (!state.channelHandlers.get(channel)?.size) {
      state.channelHandlers.delete(channel)
      state.pendingSubscriptions.delete(channel)
      if (state.ws?.readyState === WebSocket.OPEN) {
        state.ws.send(JSON.stringify({ type: 'unsubscribe', channel }))
      }
    }
  }

  /** Close the connection entirely and clear all state. */
  function disconnect() {
    if (state.reconnectTimeout) clearTimeout(state.reconnectTimeout)
    state.reconnectAttempts = MAX_RECONNECT_ATTEMPTS // prevent auto-reconnect
    state.ws?.close()
    state.ws = null
    state.channelHandlers.clear()
    state.pendingSubscriptions.clear()
    isConnected.value = false
  }

  return {
    websocketNotifications,
    isConnected,
    connect,
    send,
    unsubscribe,
    disconnect
  }
}
