package auth

import (
	"ez-admin-gin/server/internal/apperror"
	"ez-admin-gin/server/internal/middleware"
	"ez-admin-gin/server/internal/platform/datascope"
	"ez-admin-gin/server/internal/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// AccountHandler 负责当前登录人自助资料与密码接口。
type AccountHandler struct {
	service *AccountService
	log     *zap.Logger
}

// NewAccountHandler 创建账户中心 Handler。
func NewAccountHandler(service *AccountService, log *zap.Logger) *AccountHandler {
	return &AccountHandler{
		service: service,
		log:     log,
	}
}

// Profile 返回当前登录人的账户中心资料。
func (h *AccountHandler) Profile(c *gin.Context) {
	actor, ok := currentAuthActor(c, h.log)
	if !ok {
		return
	}

	result, err := h.service.GetProfile(actor)
	if err != nil {
		writeAuthError(c, err, "查询账户资料失败", h.log)
		return
	}

	response.Success(c, result)
}

// UpdateProfile 修改当前登录人的昵称。
func (h *AccountHandler) UpdateProfile(c *gin.Context) {
	actor, ok := currentAuthActor(c, h.log)
	if !ok {
		return
	}

	var req UpdateAccountProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.BadRequest("请求参数不正确"), h.log)
		return
	}

	result, err := h.service.UpdateProfile(actor, req)
	if err != nil {
		writeAuthError(c, err, "更新账户资料失败", h.log)
		return
	}

	response.Success(c, result)
}

// UpdatePassword 修改当前登录人的密码。
func (h *AccountHandler) UpdatePassword(c *gin.Context) {
	actor, ok := currentAuthActor(c, h.log)
	if !ok {
		return
	}

	var req UpdateAccountPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.BadRequest("请求参数不正确"), h.log)
		return
	}

	if err := h.service.UpdatePassword(actor, req); err != nil {
		writeAuthError(c, err, "修改账户密码失败", h.log)
		return
	}

	response.Success(c, gin.H{"updated": true})
}

func currentAuthActor(c *gin.Context, log *zap.Logger) (datascope.Actor, bool) {
	actor, ok := middleware.CurrentActor(c)
	if !ok {
		response.Error(c, apperror.Unauthorized("请先登录"), log)
		return datascope.Actor{}, false
	}

	return actor, true
}
