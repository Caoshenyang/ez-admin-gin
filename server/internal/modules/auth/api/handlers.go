// Package api 提供 auth 模块的 HTTP 请求处理器与路由定义。
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

// Login godoc
// @Summary      用户登录
// @Tags         认证
// @Accept       json
// @Produce      json
// @Param        body  body  LoginRequest  true  "登录参数"
// @Success      200  {object}  httpx.Body{data=LoginResponse}
// @Failure      400  {object}  httpx.Body
// @Failure      500  {object}  httpx.Body
// @Router       /auth/login [post]
func (h *LoginHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.service.RecordLogin(c.Request.Context(), 0, "", 2, "用户名和密码不能为空", c.ClientIP(), c.Request.UserAgent())
		httpx.Error(c, errorsx.BadRequest("用户名和密码不能为空"), h.log)
		return
	}

	result, _, err := h.service.Login(c.Request.Context(), req, c.ClientIP(), c.Request.UserAgent())
	if err != nil {
		httpx.WriteError(c, err, "登录失败", h.log)
		return
	}

	httpx.Success(c, result)
}

// LoginWithRefresh 处理登录请求并签发双 token（access + refresh cookie）。
func (h *LoginHandler) LoginWithRefresh(env string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			h.service.RecordLogin(c.Request.Context(), 0, "", 2, "用户名和密码不能为空", c.ClientIP(), c.Request.UserAgent())
			httpx.Error(c, errorsx.BadRequest("用户名和密码不能为空"), h.log)
			return
		}

		result, refreshToken, err := h.service.Login(c.Request.Context(), req, c.ClientIP(), c.Request.UserAgent())
		if err != nil {
			httpx.WriteError(c, err, "登录失败", h.log)
			return
		}

		if refreshToken != "" {
			setRefreshTokenCookie(c, refreshToken, env)
		}
		httpx.Success(c, result)
	}
}

func setRefreshTokenCookie(c *gin.Context, token string, env string) {
	secure := env == "prod"
	maxAge := 7 * 86400
	c.SetCookie("ez_admin_refresh_token", token, maxAge, "/api/v1/auth", "", secure, true)
}

func clearRefreshTokenCookie(c *gin.Context) {
	c.SetCookie("ez_admin_refresh_token", "", -1, "/api/v1/auth", "", false, true)
}

func readRefreshTokenCookie(c *gin.Context) string {
	token, _ := c.Cookie("ez_admin_refresh_token")
	return token
}

type MeHandler struct {
	service *authapp.MeService
	log     *zap.Logger
}

func NewMeHandler(service *authapp.MeService, log *zap.Logger) *MeHandler {
	return &MeHandler{service: service, log: log}
}

// Me godoc
// @Summary      获取当前用户信息
// @Tags         认证
// @Accept       json
// @Produce      json
// @Success      200  {object}  httpx.Body{data=MeResponse}
// @Failure      401  {object}  httpx.Body
// @Security     BearerAuth
// @Router       /auth/me [get]
func (h *MeHandler) Me(c *gin.Context) {
	if actor, ok := middleware.CurrentActor(c); ok {
		httpx.Success(c, h.service.Build(actor))
		return
	}

	userID, ok := middleware.RequireUserID(c, h.log)
	if !ok {
		return
	}

	httpx.Success(c, authdomain.MeResponse{
		UserID:   userID,
		Username: middleware.Username(c),
	})
}

type AccountHandler struct {
	service *authapp.AccountService
	log     *zap.Logger
}

func NewAccountHandler(service *authapp.AccountService, log *zap.Logger) *AccountHandler {
	return &AccountHandler{service: service, log: log}
}

// Profile godoc
// @Summary      查询账户资料
// @Tags         认证
// @Accept       json
// @Produce      json
// @Success      200  {object}  httpx.Body{data=AccountProfileResponse}
// @Failure      401  {object}  httpx.Body
// @Security     BearerAuth
// @Router       /auth/account [get]
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

// UpdateProfile godoc
// @Summary      更新账户资料
// @Tags         认证
// @Accept       json
// @Produce      json
// @Param        body  body  UpdateAccountProfileRequest  true  "更新参数"
// @Success      200  {object}  httpx.Body{data=AccountProfileResponse}
// @Failure      400  {object}  httpx.Body
// @Failure      401  {object}  httpx.Body
// @Security     BearerAuth
// @Router       /auth/account/profile [post]
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

// UpdatePassword godoc
// @Summary      修改账户密码
// @Tags         认证
// @Accept       json
// @Produce      json
// @Param        body  body  UpdateAccountPasswordRequest  true  "密码参数"
// @Success      200  {object}  httpx.Body
// @Failure      400  {object}  httpx.Body
// @Failure      401  {object}  httpx.Body
// @Security     BearerAuth
// @Router       /auth/account/password [post]
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

// Menus godoc
// @Summary      查询当前用户菜单
// @Tags         认证
// @Accept       json
// @Produce      json
// @Success      200  {object}  httpx.Body{data=[]MenuResponse}
// @Failure      401  {object}  httpx.Body
// @Security     BearerAuth
// @Router       /auth/menus [get]
func (h *MenuHandler) Menus(c *gin.Context) {
	userID, ok := middleware.RequireUserID(c, h.log)
	if !ok {
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

// Dashboard godoc
// @Summary      查询工作台数据
// @Tags         认证
// @Accept       json
// @Produce      json
// @Success      200  {object}  httpx.Body{data=DashboardResponse}
// @Failure      401  {object}  httpx.Body
// @Security     BearerAuth
// @Router       /auth/dashboard [get]
func (h *DashboardHandler) Dashboard(c *gin.Context) {
	userID, ok := middleware.RequireUserID(c, h.log)
	if !ok {
		return
	}

	result, err := h.service.Dashboard(userID, middleware.Username(c))
	if err != nil {
		httpx.WriteError(c, err, "查询工作台失败", h.log)
		return
	}

	httpx.Success(c, result)
}
