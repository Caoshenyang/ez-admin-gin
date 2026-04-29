package handler

import (
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
	// 检查是否已初始化（sys_user 是否有记录）
	var count int64
	if err := h.db.Model(&model.User{}).Count(&count).Error; err != nil {
		h.log.Error("check init status", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "检查初始化状态失败"})
		return
	}
	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "系统已初始化，不能重复执行"})
		return
	}

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

	// 创建管理员用户
	user := model.User{
		Username:     req.Username,
		PasswordHash: string(passwordHash),
		Nickname:     req.Nickname,
		Status:       model.UserStatusEnabled,
	}
	if err := h.db.Create(&user).Error; err != nil {
		h.log.Error("create admin user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建管理员失败"})
		return
	}

	// 绑定到 super_admin 角色（ID=1）
	userRole := model.UserRole{
		UserID: user.ID,
		RoleID: 1,
	}
	if err := h.db.Create(&userRole).Error; err != nil {
		h.log.Error("bind admin role", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "绑定管理员角色失败"})
		return
	}

	h.log.Info("admin user initialized", zap.String("username", req.Username))

	c.JSON(http.StatusOK, gin.H{
		"message":  "管理员账号创建成功",
		"user_id":  user.ID,
		"username": user.Username,
	})
}
