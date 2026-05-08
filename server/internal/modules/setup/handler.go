package setup

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var errAlreadyInitialized = errors.New("system already initialized")

type setupHandler struct {
	service *Service
	log     *zap.Logger
}

func newSetupHandler(service *Service, log *zap.Logger) *setupHandler {
	return &setupHandler{service: service, log: log}
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

	user, err := h.service.Init(c.Request.Context(), InitRequest{
		Username: req.Username,
		Password: req.Password,
		Nickname: req.Nickname,
	})
	if err != nil {
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
