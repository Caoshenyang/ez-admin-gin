package auth

import (
	"errors"
	"time"

	"ez-admin-gin/server/internal/apperror"
	"ez-admin-gin/server/internal/config"
	"ez-admin-gin/server/internal/database"
	"ez-admin-gin/server/internal/middleware"
	"ez-admin-gin/server/internal/model"
	appRedis "ez-admin-gin/server/internal/redis"
	"ez-admin-gin/server/internal/response"

	"github.com/gin-gonic/gin"
	goredis "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// DashboardHandler 负责工作台概览接口。
type DashboardHandler struct {
	cfg         *config.Config
	db          *gorm.DB
	redisClient *goredis.Client
	log         *zap.Logger
}

// NewDashboardHandler 创建工作台 Handler。
func NewDashboardHandler(
	cfg *config.Config,
	db *gorm.DB,
	redisClient *goredis.Client,
	log *zap.Logger,
) *DashboardHandler {
	return &DashboardHandler{
		cfg:         cfg,
		db:          db,
		redisClient: redisClient,
		log:         log,
	}
}

type dashboardCurrentUser struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
}

type dashboardHealth struct {
	Env      string `json:"env"`
	Database string `json:"database"`
	Redis    string `json:"redis"`
}

type dashboardMetrics struct {
	UserTotal               int64 `json:"user_total"`
	EnabledUserTotal        int64 `json:"enabled_user_total"`
	EnabledRoleTotal        int64 `json:"enabled_role_total"`
	ConfigTotal             int64 `json:"config_total"`
	NoticeTotal             int64 `json:"notice_total"`
	FileTotal               int64 `json:"file_total"`
	TodayOperationTotal     int64 `json:"today_operation_total"`
	TodayRiskOperationTotal int64 `json:"today_risk_operation_total"`
	TodayLoginFailedTotal   int64 `json:"today_login_failed_total"`
}

type dashboardOperationItem struct {
	ID         uint      `json:"id"`
	Username   string    `json:"username"`
	Method     string    `json:"method"`
	Path       string    `json:"path"`
	StatusCode int       `json:"status_code"`
	Success    bool      `json:"success"`
	LatencyMs  int64     `json:"latency_ms"`
	CreatedAt  time.Time `json:"created_at"`
}

type dashboardLoginItem struct {
	ID        uint                 `json:"id"`
	Username  string               `json:"username"`
	Status    model.LoginLogStatus `json:"status"`
	Message   string               `json:"message"`
	IP        string               `json:"ip"`
	CreatedAt time.Time            `json:"created_at"`
}

type dashboardNoticeItem struct {
	ID        uint               `json:"id"`
	Title     string             `json:"title"`
	Status    model.NoticeStatus `json:"status"`
	UpdatedAt time.Time          `json:"updated_at"`
}

type dashboardResponse struct {
	CurrentUser      dashboardCurrentUser     `json:"current_user"`
	Health           dashboardHealth          `json:"health"`
	Metrics          dashboardMetrics         `json:"metrics"`
	RecentOperations []dashboardOperationItem `json:"recent_operations"`
	RecentLogins     []dashboardLoginItem     `json:"recent_logins"`
	LatestNotices    []dashboardNoticeItem    `json:"latest_notices"`
}

// Dashboard 返回工作台概览数据。
func (h *DashboardHandler) Dashboard(c *gin.Context) {
	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		response.Error(c, apperror.Unauthorized("请先登录"), h.log)
		return
	}

	username, _ := middleware.CurrentUsername(c)

	if err := database.Ping(h.db); err != nil {
		response.Error(c, apperror.ServiceUnavailable("数据库不可用", err), h.log)
		return
	}

	currentUser, err := h.loadCurrentUser(userID, username)
	if err != nil {
		response.Error(c, err, h.log)
		return
	}

	metrics, err := h.loadMetrics()
	if err != nil {
		response.Error(c, err, h.log)
		return
	}

	recentOperations, err := h.loadRecentOperations()
	if err != nil {
		response.Error(c, err, h.log)
		return
	}

	recentLogins, err := h.loadRecentLogins()
	if err != nil {
		response.Error(c, err, h.log)
		return
	}

	latestNotices, err := h.loadLatestNotices()
	if err != nil {
		response.Error(c, err, h.log)
		return
	}

	health := dashboardHealth{
		Env:      h.cfg.App.Env,
		Database: "ok",
		Redis:    "ok",
	}

	if err := appRedis.Ping(h.redisClient); err != nil {
		health.Redis = "error"
		h.log.Warn("dashboard redis ping failed", zap.Error(err))
	}

	response.Success(c, dashboardResponse{
		CurrentUser:      currentUser,
		Health:           health,
		Metrics:          metrics,
		RecentOperations: recentOperations,
		RecentLogins:     recentLogins,
		LatestNotices:    latestNotices,
	})
}

func (h *DashboardHandler) loadCurrentUser(userID uint, fallbackUsername string) (dashboardCurrentUser, error) {
	var user model.User
	if err := h.db.Select("id", "username", "nickname").First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dashboardCurrentUser{}, apperror.Unauthorized("登录状态无效，请重新登录")
		}
		return dashboardCurrentUser{}, apperror.Internal("查询当前用户失败", err)
	}

	if user.Username == "" {
		user.Username = fallbackUsername
	}

	return dashboardCurrentUser{
		UserID:   user.ID,
		Username: user.Username,
		Nickname: user.Nickname,
	}, nil
}

func (h *DashboardHandler) loadMetrics() (dashboardMetrics, error) {
	now := time.Now()
	dayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	metrics := dashboardMetrics{}

	if err := h.db.Model(&model.User{}).Count(&metrics.UserTotal).Error; err != nil {
		return dashboardMetrics{}, apperror.Internal("查询用户总数失败", err)
	}
	if err := h.db.Model(&model.User{}).
		Where("status = ?", model.UserStatusEnabled).
		Count(&metrics.EnabledUserTotal).Error; err != nil {
		return dashboardMetrics{}, apperror.Internal("查询启用用户总数失败", err)
	}
	if err := h.db.Model(&model.Role{}).
		Where("status = ?", model.RoleStatusEnabled).
		Count(&metrics.EnabledRoleTotal).Error; err != nil {
		return dashboardMetrics{}, apperror.Internal("查询启用角色总数失败", err)
	}
	if err := h.db.Model(&model.SystemConfig{}).
		Where("status = ?", model.SystemConfigStatusEnabled).
		Count(&metrics.ConfigTotal).Error; err != nil {
		return dashboardMetrics{}, apperror.Internal("查询配置总数失败", err)
	}
	if err := h.db.Model(&model.Notice{}).
		Where("status = ?", model.NoticeStatusEnabled).
		Count(&metrics.NoticeTotal).Error; err != nil {
		return dashboardMetrics{}, apperror.Internal("查询公告总数失败", err)
	}
	if err := h.db.Model(&model.SystemFile{}).Count(&metrics.FileTotal).Error; err != nil {
		return dashboardMetrics{}, apperror.Internal("查询文件总数失败", err)
	}
	if err := h.db.Model(&model.OperationLog{}).
		Where("created_at >= ?", dayStart).
		Count(&metrics.TodayOperationTotal).Error; err != nil {
		return dashboardMetrics{}, apperror.Internal("查询今日操作总数失败", err)
	}
	if err := h.db.Model(&model.OperationLog{}).
		Where("created_at >= ?", dayStart).
		Where("success = ?", false).
		Count(&metrics.TodayRiskOperationTotal).Error; err != nil {
		return dashboardMetrics{}, apperror.Internal("查询今日失败操作总数失败", err)
	}
	if err := h.db.Model(&model.LoginLog{}).
		Where("created_at >= ?", dayStart).
		Where("status = ?", model.LoginLogStatusFailed).
		Count(&metrics.TodayLoginFailedTotal).Error; err != nil {
		return dashboardMetrics{}, apperror.Internal("查询今日登录失败总数失败", err)
	}

	return metrics, nil
}

func (h *DashboardHandler) loadRecentOperations() ([]dashboardOperationItem, error) {
	var rows []model.OperationLog
	if err := h.db.Order("id DESC").Limit(6).Find(&rows).Error; err != nil {
		return nil, apperror.Internal("查询最近操作失败", err)
	}

	items := make([]dashboardOperationItem, 0, len(rows))
	for _, item := range rows {
		items = append(items, dashboardOperationItem{
			ID:         item.ID,
			Username:   item.Username,
			Method:     item.Method,
			Path:       item.Path,
			StatusCode: item.StatusCode,
			Success:    item.Success,
			LatencyMs:  item.LatencyMs,
			CreatedAt:  item.CreatedAt,
		})
	}

	return items, nil
}

func (h *DashboardHandler) loadRecentLogins() ([]dashboardLoginItem, error) {
	var rows []model.LoginLog
	if err := h.db.Order("id DESC").Limit(5).Find(&rows).Error; err != nil {
		return nil, apperror.Internal("查询最近登录记录失败", err)
	}

	items := make([]dashboardLoginItem, 0, len(rows))
	for _, item := range rows {
		items = append(items, dashboardLoginItem{
			ID:        item.ID,
			Username:  item.Username,
			Status:    item.Status,
			Message:   item.Message,
			IP:        item.IP,
			CreatedAt: item.CreatedAt,
		})
	}

	return items, nil
}

func (h *DashboardHandler) loadLatestNotices() ([]dashboardNoticeItem, error) {
	var rows []model.Notice
	if err := h.db.
		Where("status = ?", model.NoticeStatusEnabled).
		Order("updated_at DESC, id DESC").
		Limit(3).
		Find(&rows).Error; err != nil {
		return nil, apperror.Internal("查询最近公告失败", err)
	}

	items := make([]dashboardNoticeItem, 0, len(rows))
	for _, item := range rows {
		items = append(items, dashboardNoticeItem{
			ID:        item.ID,
			Title:     item.Title,
			Status:    item.Status,
			UpdatedAt: item.UpdatedAt,
		})
	}

	return items, nil
}
