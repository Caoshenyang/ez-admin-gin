package ws

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"go.uber.org/zap"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

// Client 封装一个 WebSocket 连接。
type Client struct {
	conn       *websocket.Conn
	UserID     uint
	RemoteAddr string
	hub        *Hub
	log        *zap.Logger
	send       chan any
	cancel     context.CancelFunc
	wg         sync.WaitGroup
}

// NewClient 创建客户端实例。
func NewClient(conn *websocket.Conn, userID uint, hub *Hub, log *zap.Logger) *Client {
	_, cancel := context.WithCancel(context.Background())

	return &Client{
		conn:       conn,
		UserID:     userID,
		RemoteAddr: "ws-client",
		hub:        hub,
		log:        log,
		send:       make(chan any, 64),
		cancel:     cancel,
	}
}

// Run 启动读写循环。
func (c *Client) Run() {
	c.hub.Register(c)

	c.wg.Add(2)
	go c.writePump()
	go c.readPump()
	c.wg.Wait()
}

// Send 向客户端发送消息。
func (c *Client) Send(msg any) {
	select {
	case c.send <- msg:
	default:
		c.log.Warn("ws client send buffer full, dropping message",
			zap.Uint("user_id", c.UserID),
		)
	}
}

// Close 关闭客户端连接。
func (c *Client) Close() {
	c.cancel()
	c.conn.Close(websocket.StatusNormalClosure, "closing")
}

// writePump 从 send channel 读取消息并写入 WebSocket。
func (c *Client) writePump() {
	defer c.wg.Done()

	ctx := context.Background()
	for {
		select {
		case msg, ok := <-c.send:
			if !ok {
				return
			}
			writeCtx, writeCancel := context.WithTimeout(ctx, 10*time.Second)
			err := wsjson.Write(writeCtx, c.conn, msg)
			writeCancel()
			if err != nil {
				c.log.Debug("ws write error", zap.Error(err))
				c.cleanup()
				return
			}
		case <-time.After(30 * time.Second):
			// ping
			pingCtx, pingCancel := context.WithTimeout(ctx, 5*time.Second)
			err := c.conn.Ping(pingCtx)
			pingCancel()
			if err != nil {
				c.cleanup()
				return
			}
		}
	}
}

// readPump 从 WebSocket 读取客户端消息。
func (c *Client) readPump() {
	defer c.wg.Done()

	ctx := context.Background()
	for {
		readCtx, readCancel := context.WithTimeout(ctx, 60*time.Second)
		_, data, err := c.conn.Read(readCtx)
		readCancel()
		if err != nil {
			c.cleanup()
			return
		}

		c.handleMessage(data)
	}
}

type clientMessage struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data,omitempty"`
}

type markReadData struct {
	IDs []uint64 `json:"ids"`
}

func (c *Client) handleMessage(raw []byte) {
	var msg clientMessage
	if err := json.Unmarshal(raw, &msg); err != nil {
		return
	}

	switch msg.Type {
	case "pong":
		// heartbeat response, no-op
	case "mark_read":
		var data markReadData
		if err := json.Unmarshal(msg.Data, &data); err != nil {
			return
		}
		if err := c.hub.service.MarkRead(c.UserID, data.IDs); err != nil {
			c.log.Error("ws mark_read error", zap.Error(err))
		}
	case "mark_all_read":
		if err := c.hub.service.MarkAllRead(c.UserID); err != nil {
			c.log.Error("ws mark_all_read error", zap.Error(err))
		}
	}
}

func (c *Client) cleanup() {
	c.cancel()
	c.hub.Unregister(c)
	c.conn.Close(websocket.StatusNormalClosure, "closing")
}
