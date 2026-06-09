package api

import (
	"strings"

	authapp "ez-admin-gin/server/internal/modules/auth/application"
	httpx "ez-admin-gin/server/internal/pkg/httpx"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// LogoutHandler handles POST /auth/logout.
type LogoutHandler struct {
	service *authapp.LogoutService
	log     *zap.Logger
}

// NewLogoutHandler creates a new LogoutHandler.
func NewLogoutHandler(service *authapp.LogoutService, log *zap.Logger) *LogoutHandler {
	return &LogoutHandler{service: service, log: log}
}

// Logout godoc
// @Summary      退出登录
// @Tags         认证
// @Accept       json
// @Produce      json
// @Success      200  {object}  httpx.Body
// @Failure      401  {object}  httpx.Body
// @Security     BearerAuth
// @Router       /auth/logout [post]
func (h *LogoutHandler) Logout(c *gin.Context) {
	refreshToken := readRefreshTokenCookie(c)

	accessToken := ""
	if parts := strings.Fields(c.GetHeader("Authorization")); len(parts) == 2 {
		accessToken = parts[1]
	}

	_ = h.service.Logout(c.Request.Context(), refreshToken, accessToken)

	clearRefreshTokenCookie(c)
	httpx.Success(c, gin.H{"logged_out": true})
}
