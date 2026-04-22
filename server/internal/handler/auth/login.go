package auth

import (
	"errors"
	"ez-admin-gin/server/internal/token"
	"strings"
	"time"

	"ez-admin-gin/server/internal/apperror"
	"ez-admin-gin/server/internal/model"
	"ez-admin-gin/server/internal/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// LoginHandler 负责登录相关接口。
type LoginHandler struct {
	db           *gorm.DB
	log          *zap.Logger
	tokenManager *token.Manager
}

// NewLoginHandler 创建登录 Handler。
func NewLoginHandler(db *gorm.DB, log *zap.Logger, tokenManager *token.Manager) *LoginHandler {
	return &LoginHandler{
		db:           db,
		log:          log,
		tokenManager: tokenManager,
	}
}

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type loginResponse struct {
	UserID      uint   `json:"user_id"`
	Username    string `json:"username"`
	Nickname    string `json:"nickname"`
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresAt   string `json:"expires_at"`
}

// Login 校验用户名和密码。
func (h *LoginHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.BadRequest("用户名和密码不能为空"), h.log)
		return
	}

	req.Username = strings.TrimSpace(req.Username)
	if req.Username == "" || req.Password == "" {
		response.Error(c, apperror.BadRequest("用户名和密码不能为空"), h.log)
		return
	}

	var user model.User
	// GORM 会自动过滤 deleted_at 不为空的记录。
	err := h.db.Where("username = ?", req.Username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(c, apperror.Unauthorized("用户名或密码错误"), h.log)
			return
		}

		h.log.Error("query login user failed", zap.Error(err))
		response.Error(c, apperror.Internal("登录失败", err), h.log)
		return
	}

	if user.Status != model.UserStatusEnabled {
		response.Error(c, apperror.Forbidden("用户已被禁用"), h.log)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		response.Error(c, apperror.Unauthorized("用户名或密码错误"), h.log)
		return
	}

	accessToken, expiresAt, err := h.tokenManager.GenerateAccessToken(user.ID, user.Username)
	if err != nil {
		response.Error(c, apperror.Internal("登录失败", err), h.log)
		return
	}

	response.Success(c, loginResponse{
		UserID:      user.ID,
		Username:    user.Username,
		Nickname:    user.Nickname,
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresAt:   expiresAt.UTC().Format(time.RFC3339),
	})
}
