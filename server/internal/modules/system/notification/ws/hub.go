package ws

import (
	"sync"

	notiapp "ez-admin-gin/server/internal/modules/system/notification/application"
	notidomain "ez-admin-gin/server/internal/modules/system/notification/domain"

	"go.uber.org/zap"
)

// Hub 管理所有 WebSocket 连接，按 userID 分组。
type Hub struct {
	mu       sync.RWMutex
	clients  map[uint][]*Client // userID → clients
	service  *notiapp.Service
	transport Transport
	log      *zap.Logger
}

// NewHub 创建 Hub 实例。
func NewHub(service *notiapp.Service, transport Transport, log *zap.Logger) *Hub {
	return &Hub{
		clients:   make(map[uint][]*Client),
		service:   service,
		transport: transport,
		log:       log,
	}
}

// Run 启动 Hub 的消息分发循环。
func (h *Hub) Run() {
	if h.transport == nil {
		return
	}
	ch := h.transport.Subscribe()
	go func() {
		for msg := range ch {
			h.dispatch(msg)
		}
	}()
}

// Register 注册客户端连接。
func (h *Hub) Register(client *Client) {
	h.mu.Lock()
	h.clients[client.UserID] = append(h.clients[client.UserID], client)
	h.mu.Unlock()

	h.log.Debug("ws client registered",
		zap.Uint("user_id", client.UserID),
		zap.String("remote", client.RemoteAddr),
	)
}

// Unregister 注销客户端连接。
func (h *Hub) Unregister(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	clients := h.clients[client.UserID]
	for i, c := range clients {
		if c == client {
			h.clients[client.UserID] = append(clients[:i], clients[i+1:]...)
			break
		}
	}
	if len(h.clients[client.UserID]) == 0 {
		delete(h.clients, client.UserID)
	}

	h.log.Debug("ws client unregistered",
		zap.Uint("user_id", client.UserID),
		zap.String("remote", client.RemoteAddr),
	)
}

// SendToUser 向指定用户的所有在线连接推送消息。
func (h *Hub) SendToUser(userID uint, msg any) {
	h.mu.RLock()
	clients := h.clients[userID]
	h.mu.RUnlock()

	for _, c := range clients {
		c.Send(msg)
	}
}

// SendNotificationAndCount 推送新通知并更新未读数。
func (h *Hub) SendNotificationAndCount(userID uint, resp notidomain.Response, count int64) {
	h.SendToUser(userID, map[string]any{
		"type": "notification",
		"data": resp,
	})
	h.SendToUser(userID, map[string]any{
		"type": "unread_count",
		"data": map[string]any{"count": count},
	})
}

// dispatch 分发从 Redis 收到的消息。
func (h *Hub) dispatch(msg TransportMessage) {
	switch msg.Type {
	case "notification":
		h.SendToUser(msg.UserID, map[string]any{
			"type": "notification",
			"data": msg.Data,
		})
	case "unread_count":
		h.SendToUser(msg.UserID, map[string]any{
			"type": "unread_count",
			"data": msg.Data,
		})
	}
}

// OnlineUsers 返回当前在线用户数。
func (h *Hub) OnlineUsers() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}
