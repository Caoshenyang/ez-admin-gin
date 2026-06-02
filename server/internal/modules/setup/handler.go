package setup

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var errAlreadyInitialized = errors.New("system already initialized")

const (
	defaultAdminUsername = "admin"
	defaultAdminPassword = "EzAdmin@123456"
	defaultAdminNickname = "管理员"

	initSuccessMessage = "管理员账号创建成功，请及时修改默认密码"
)

type setupHandler struct {
	service *Service
	log     *zap.Logger
}

func newSetupHandler(service *Service, log *zap.Logger) *setupHandler {
	return &setupHandler{service: service, log: log}
}

type initResponse struct {
	Message  string `json:"message"`
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
}

// Init godoc
// @Summary      系统初始化
// @Description  使用默认账号创建超级管理员，仅当系统中无用户时可执行。
// @Tags         系统
// @Produce      json
// @Success      200  {object}  initResponse
// @Failure      409  {object}  map[string]any
// @Failure      500  {object}  map[string]any
// @Router       /setup/init [post]
func (h *setupHandler) Init(c *gin.Context) {
	req := InitRequest{
		Username: defaultAdminUsername,
		Password: defaultAdminPassword,
		Nickname: defaultAdminNickname,
	}
	user, err := h.service.Init(c.Request.Context(), req)
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
	c.JSON(http.StatusOK, initResponse{
		Message:  initSuccessMessage,
		UserID:   user.ID,
		Username: user.Username,
	})
}
