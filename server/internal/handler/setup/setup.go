package handler

import (
	"errors"
	"net/http"

	"ez-admin-gin/server/internal/model"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// SetupHandler 处理管理员一次性初始化。
type SetupHandler struct {
	db  *gorm.DB
	log *zap.Logger
}

// NewSetupHandler 创建 SetupHandler 实例。
func NewSetupHandler(db *gorm.DB, log *zap.Logger) *SetupHandler {
	return &SetupHandler{db: db, log: log}
}

// InitRequest 是管理员初始化接口的请求体。
type InitRequest struct {
	Username string `json:"username" binding:"required,min=2,max=64"`
	Password string `json:"password" binding:"required,min=6,max=128"`
	Nickname string `json:"nickname" binding:"required,min=1,max=64"`
}

// Init 创建第一个管理员账号并绑定到 super_admin 角色。
// POST /api/v1/setup/init
func (h *SetupHandler) Init(c *gin.Context) {
	var req InitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数无效"})
		return
	}

	// bcrypt 加密密码
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		h.log.Error("hash password", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
		return
	}

	// 首次初始化必须在一个事务里完成，避免出现“用户已创建但未绑定管理员角色”的半成品状态。
	var user model.User
	if err := h.db.Transaction(func(tx *gorm.DB) error {
		var count int64
		if err := tx.Model(&model.User{}).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return errAlreadyInitialized
		}

		var role model.Role
		if err := tx.
			Where("code = ?", "super_admin").
			Where("status = ?", model.RoleStatusEnabled).
			First(&role).Error; err != nil {
			return err
		}

		user = model.User{
			Username:     req.Username,
			PasswordHash: string(passwordHash),
			Nickname:     req.Nickname,
			Status:       model.UserStatusEnabled,
		}
		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		return tx.Create(&model.UserRole{
			UserID: user.ID,
			RoleID: role.ID,
		}).Error
	}); err != nil {
		switch {
		case errors.Is(err, errAlreadyInitialized):
			c.JSON(http.StatusConflict, gin.H{"error": "系统已初始化，不能重复执行"})
		case errors.Is(err, gorm.ErrRecordNotFound):
			h.log.Error("super admin role missing for setup init")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "初始化角色不存在，请先检查种子数据"})
		default:
			h.log.Error("initialize admin user", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "初始化管理员失败"})
		}
		return
	}

	h.log.Info("admin user initialized", zap.String("username", req.Username))

	c.JSON(http.StatusOK, gin.H{
		"message":  "管理员账号创建成功",
		"user_id":  user.ID,
		"username": user.Username,
	})
}

var errAlreadyInitialized = errors.New("system already initialized")
