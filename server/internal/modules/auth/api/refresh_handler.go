package api

import (
	authapp "ez-admin-gin/server/internal/modules/auth/application"
	errorsx "ez-admin-gin/server/internal/pkg/errorsx"
	httpx "ez-admin-gin/server/internal/pkg/httpx"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RefreshHandler handles POST /auth/refresh.
type RefreshHandler struct {
	service *authapp.RefreshService
	env     string
	log     *zap.Logger
}

// NewRefreshHandler creates a new RefreshHandler.
func NewRefreshHandler(service *authapp.RefreshService, env string, log *zap.Logger) *RefreshHandler {
	return &RefreshHandler{service: service, env: env, log: log}
}

// Refresh godoc
// @Summary      刷新 Access Token
// @Tags         认证
// @Accept       json
// @Produce      json
// @Success      200  {object}  httpx.Body{data=authdomain.LoginResponse}
// @Failure      401  {object}  httpx.Body
// @Router       /auth/refresh [post]
func (h *RefreshHandler) Refresh(c *gin.Context) {
	refreshToken := readRefreshTokenCookie(c)
	if refreshToken == "" {
		httpx.Error(c, errorsx.Unauthorized("请先登录"), h.log)
		return
	}

	result, newRefreshToken, err := h.service.Refresh(c.Request.Context(), refreshToken)
	if err != nil {
		httpx.WriteError(c, err, "刷新 token 失败", h.log)
		return
	}

	if newRefreshToken != "" {
		setRefreshTokenCookie(c, newRefreshToken, h.env)
	}
	httpx.Success(c, result)
}
