package setup

import (
	"errors"
	"net/http"

	"ez-admin-gin/server/internal/platform/model"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var errAlreadyInitialized = errors.New("system already initialized")

type setupHandler struct {
	db  *gorm.DB
	log *zap.Logger
}

func newSetupHandler(db *gorm.DB, log *zap.Logger) *setupHandler {
	return &setupHandler{db: db, log: log}
}

type initRequest struct {
	Username string `json:"username" binding:"required,min=2,max=64"`
	Password string `json:"password" binding:"required,min=6,max=128"`
	Nickname string `json:"nickname" binding:"required,min=1,max=64"`
}

// Init godoc
// @Summary      系统初始化
// @Description  创建超级管理员账号，仅当系统中无用户时可执行。
// @Tags         系统
// @Accept       json
// @Produce      json
// @Param        body  body  initRequest  true  "初始化参数"
// @Success      200  {object}  map[string]any
// @Failure      400  {object}  map[string]any
// @Failure      409  {object}  map[string]any
// @Failure      500  {object}  map[string]any
// @Router       /setup/init [post]
func (h *setupHandler) Init(c *gin.Context) {
	var req initRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数无效"})
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		h.log.Error("hash password", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
		return
	}

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
		if err := tx.Where("code = ?", "super_admin").Where("status = ?", model.RoleStatusEnabled).First(&role).Error; err != nil {
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

		return tx.Create(&model.UserRole{UserID: user.ID, RoleID: role.ID}).Error
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
