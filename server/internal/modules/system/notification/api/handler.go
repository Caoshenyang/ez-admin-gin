package api

import (
	"net/http"

	notiapp "ez-admin-gin/server/internal/modules/system/notification/application"
	notidomain "ez-admin-gin/server/internal/modules/system/notification/domain"
	notiws "ez-admin-gin/server/internal/modules/system/notification/ws"
	errorsx "ez-admin-gin/server/internal/pkg/errorsx"
	httpx "ez-admin-gin/server/internal/pkg/httpx"
	authnPlatform "ez-admin-gin/server/internal/platform/authn"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"nhooyr.io/websocket"
)

// Handler 处理通知的 HTTP 和 WebSocket 请求。
type Handler struct {
	service *notiapp.Service
	hub     *notiws.Hub
	log     *zap.Logger
}

func NewHandler(service *notiapp.Service, hub *notiws.Hub, log *zap.Logger) *Handler {
	return &Handler{service: service, hub: hub, log: log}
}

// WSHandler 处理 WebSocket 连接（单独的 handler，因为不走 Auth 中间件）。
type WSHandler struct {
	service *notiapp.Service
	hub     interface{ Run() }
	token   *authnPlatform.Manager
	log     *zap.Logger
}

func NewWSHandler(service *notiapp.Service, hub interface{ Run() }, token *authnPlatform.Manager, log *zap.Logger) *WSHandler {
	return &WSHandler{service: service, hub: hub, token: token, log: log}
}

// List godoc
// @Summary      查询通知列表
// @Tags         System / 通知管理
// @Accept       json
// @Produce      json
// @Param        page       query     int     false  "页码"
// @Param        page_size  query     int     false  "每页条数"
// @Param        type       query     int     false  "通知类型"
// @Param        is_read    query     int     false  "已读状态 0=全部 1=未读 2=已读"
// @Success      200  {object}  httpx.Body{data=notidomain.ListResponse}
// @Failure      400  {object}  httpx.Body
// @Failure      401  {object}  httpx.Body
// @Security     BearerAuth
// @Router       /system/notifications [get]
func (h *Handler) List(c *gin.Context) {
	actor, ok := httpx.CurrentActor(c, h.log)
	if !ok {
		return
	}

	var query ListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		httpx.Error(c, errorsx.BadRequest("查询参数不正确"), h.log)
		return
	}

	result, err := h.service.List(actor.UserID, notidomain.ListQuery{
		Page:     query.Page,
		PageSize: query.PageSize,
		Type:     query.Type,
		IsRead:   query.IsRead,
	})
	if err != nil {
		httpx.WriteError(c, err, "查询通知列表失败", h.log)
		return
	}
	httpx.Success(c, result)
}

// UnreadCount godoc
// @Summary      获取未读通知数
// @Tags         System / 通知管理
// @Produce      json
// @Success      200  {object}  httpx.Body{data=notidomain.UnreadCountResponse}
// @Failure      401  {object}  httpx.Body
// @Security     BearerAuth
// @Router       /system/notifications/unread-count [get]
func (h *Handler) UnreadCount(c *gin.Context) {
	actor, ok := httpx.CurrentActor(c, h.log)
	if !ok {
		return
	}

	count, err := h.service.UnreadCount(actor.UserID)
	if err != nil {
		httpx.WriteError(c, err, "获取未读数失败", h.log)
		return
	}
	httpx.Success(c, notidomain.UnreadCountResponse{Count: count})
}

// MarkRead godoc
// @Summary      标记通知已读
// @Tags         System / 通知管理
// @Accept       json
// @Produce      json
// @Param        body  body  notidomain.MarkReadRequest  true  "通知 ID 列表"
// @Success      200  {object}  httpx.Body
// @Failure      400  {object}  httpx.Body
// @Failure      401  {object}  httpx.Body
// @Security     BearerAuth
// @Router       /system/notifications/mark-read [post]
func (h *Handler) MarkRead(c *gin.Context) {
	actor, ok := httpx.CurrentActor(c, h.log)
	if !ok {
		return
	}

	var req MarkReadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, errorsx.BadRequest("请求参数不正确"), h.log)
		return
	}

	if err := h.service.MarkRead(actor.UserID, req.IDs); err != nil {
		httpx.WriteError(c, err, "标记已读失败", h.log)
		return
	}
	httpx.Success(c, nil)
}

// MarkAllRead godoc
// @Summary      全部标记已读
// @Tags         System / 通知管理
// @Produce      json
// @Success      200  {object}  httpx.Body
// @Failure      401  {object}  httpx.Body
// @Security     BearerAuth
// @Router       /system/notifications/mark-all-read [post]
func (h *Handler) MarkAllRead(c *gin.Context) {
	actor, ok := httpx.CurrentActor(c, h.log)
	if !ok {
		return
	}

	if err := h.service.MarkAllRead(actor.UserID); err != nil {
		httpx.WriteError(c, err, "全部标记已读失败", h.log)
		return
	}
	httpx.Success(c, nil)
}

// ServeWebSocket 处理 WebSocket 连接升级。
// 浏览器 WebSocket API 不支持自定义 Header，因此从 query param ?token=xxx 解析 access token。
func (h *WSHandler) ServeWebSocket(c *gin.Context) {
	tokenStr := c.Query("token")
	if tokenStr == "" {
		c.JSON(http.StatusUnauthorized, httpx.Body{
			Code:    errorsx.CodeUnauthorized,
			Message: "缺少 token 参数",
		})
		return
	}

	claims, err := h.token.ParseAccessToken(tokenStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, httpx.Body{
			Code:    errorsx.CodeUnauthorized,
			Message: "token 无效或已过期",
		})
		return
	}

	wsConn, err := websocket.Accept(c.Writer, c.Request, &websocket.AcceptOptions{
		OriginPatterns: []string{"*"},
	})
	if err != nil {
		h.log.Error("websocket accept error", zap.Error(err))
		return
	}

	hub, ok := h.hub.(*notiws.Hub)
	if !ok {
		wsConn.Close(websocket.StatusInternalError, "hub not available")
		return
	}

	client := notiws.NewClient(wsConn, claims.UserID, hub, h.log)
	go client.Run()

	// 推送当前未读数
	count, err := h.service.UnreadCount(claims.UserID)
	if err == nil {
		client.Send(map[string]any{
			"type": "unread_count",
			"data": map[string]any{"count": count},
		})
	}
}

var _ = notidomain.PermissionList
