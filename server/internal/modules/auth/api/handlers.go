package api

import (
	authapp "ez-admin-gin/server/internal/modules/auth/application"
	authdomain "ez-admin-gin/server/internal/modules/auth/domain"
	errorsx "ez-admin-gin/server/internal/pkg/errorsx"
	httpx "ez-admin-gin/server/internal/pkg/httpx"
	"ez-admin-gin/server/internal/platform/middleware"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type LoginHandler struct {
	service *authapp.LoginService
	log     *zap.Logger
}

func NewLoginHandler(service *authapp.LoginService, log *zap.Logger) *LoginHandler {
	return &LoginHandler{service: service, log: log}
}

func (h *LoginHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.service.RecordLogin(c.Request.Context(), 0, "", 2, "用户名和密码不能为空", c.ClientIP(), c.Request.UserAgent())
		httpx.Error(c, errorsx.BadRequest("用户名和密码不能为空"), h.log)
		return
	}

	result, err := h.service.Login(c.Request.Context(), req, c.ClientIP(), c.Request.UserAgent())
	if err != nil {
		httpx.WriteError(c, err, "登录失败", h.log)
		return
	}

	httpx.Success(c, result)
}

type MeHandler struct {
	service *authapp.MeService
	log     *zap.Logger
}

func NewMeHandler(service *authapp.MeService, log *zap.Logger) *MeHandler {
	return &MeHandler{service: service, log: log}
}

func (h *MeHandler) Me(c *gin.Context) {
	if actor, ok := middleware.CurrentActor(c); ok {
		httpx.Success(c, h.service.Build(actor))
		return
	}

	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		httpx.Error(c, errorsx.Unauthorized("请先登录"), h.log)
		return
	}

	username, _ := middleware.CurrentUsername(c)
	httpx.Success(c, authdomain.MeResponse{
		UserID:   userID,
		Username: username,
	})
}

type AccountHandler struct {
	service *authapp.AccountService
	log     *zap.Logger
}

func NewAccountHandler(service *authapp.AccountService, log *zap.Logger) *AccountHandler {
	return &AccountHandler{service: service, log: log}
}

func (h *AccountHandler) Profile(c *gin.Context) {
	actor, ok := httpx.CurrentActor(c, h.log)
	if !ok {
		return
	}

	result, err := h.service.GetProfile(actor)
	if err != nil {
		httpx.WriteError(c, err, "查询账户资料失败", h.log)
		return
	}

	httpx.Success(c, result)
}

func (h *AccountHandler) UpdateProfile(c *gin.Context) {
	actor, ok := httpx.CurrentActor(c, h.log)
	if !ok {
		return
	}

	var req UpdateAccountProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, errorsx.BadRequest("请求参数不正确"), h.log)
		return
	}

	result, err := h.service.UpdateProfile(actor, req)
	if err != nil {
		httpx.WriteError(c, err, "更新账户资料失败", h.log)
		return
	}

	httpx.Success(c, result)
}

func (h *AccountHandler) UpdatePassword(c *gin.Context) {
	actor, ok := httpx.CurrentActor(c, h.log)
	if !ok {
		return
	}

	var req UpdateAccountPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, errorsx.BadRequest("请求参数不正确"), h.log)
		return
	}

	if err := h.service.UpdatePassword(actor, req); err != nil {
		httpx.WriteError(c, err, "修改账户密码失败", h.log)
		return
	}

	httpx.Success(c, gin.H{"updated": true})
}

type MenuHandler struct {
	service *authapp.MenuService
	log     *zap.Logger
}

func NewMenuHandler(service *authapp.MenuService, log *zap.Logger) *MenuHandler {
	return &MenuHandler{service: service, log: log}
}

func (h *MenuHandler) Menus(c *gin.Context) {
	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		httpx.Error(c, errorsx.Unauthorized("请先登录"), h.log)
		return
	}

	result, err := h.service.Menus(userID)
	if err != nil {
		httpx.WriteError(c, err, "查询菜单失败", h.log)
		return
	}

	httpx.Success(c, result)
}

type DashboardHandler struct {
	service *authapp.DashboardService
	log     *zap.Logger
}

func NewDashboardHandler(service *authapp.DashboardService, log *zap.Logger) *DashboardHandler {
	return &DashboardHandler{service: service, log: log}
}

func (h *DashboardHandler) Dashboard(c *gin.Context) {
	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		httpx.Error(c, errorsx.Unauthorized("请先登录"), h.log)
		return
	}

	username, _ := middleware.CurrentUsername(c)
	result, err := h.service.Dashboard(userID, username)
	if err != nil {
		httpx.WriteError(c, err, "查询工作台失败", h.log)
		return
	}

	httpx.Success(c, result)
}
