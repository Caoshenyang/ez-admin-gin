package auth

import (
	"errors"
	"time"

	"ez-admin-gin/server/internal/apperror"
	"ez-admin-gin/server/internal/config"
	"ez-admin-gin/server/internal/database"
	appRedis "ez-admin-gin/server/internal/redis"

	goredis "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// DashboardService 负责工作台概览数据组装。
type DashboardService struct {
	cfg   *config.Config
	db    *gorm.DB
	repo  *Repository
	redis *goredis.Client
	log   *zap.Logger
}

// NewDashboardService 创建工作台服务。
func NewDashboardService(
	cfg *config.Config,
	db *gorm.DB,
	repo *Repository,
	redis *goredis.Client,
	log *zap.Logger,
) *DashboardService {
	return &DashboardService{
		cfg:   cfg,
		db:    db,
		repo:  repo,
		redis: redis,
		log:   log,
	}
}

// Dashboard 返回工作台概览数据。
func (s *DashboardService) Dashboard(userID uint, fallbackUsername string) (DashboardResponse, error) {
	if err := database.Ping(s.db); err != nil {
		return DashboardResponse{}, apperror.ServiceUnavailable("数据库不可用", err)
	}

	currentUser, err := s.loadCurrentUser(userID, fallbackUsername)
	if err != nil {
		return DashboardResponse{}, err
	}

	metrics, err := s.loadMetrics()
	if err != nil {
		return DashboardResponse{}, err
	}

	recentOperations, err := s.loadRecentOperations()
	if err != nil {
		return DashboardResponse{}, err
	}

	recentLogins, err := s.loadRecentLogins()
	if err != nil {
		return DashboardResponse{}, err
	}

	latestNotices, err := s.loadLatestNotices()
	if err != nil {
		return DashboardResponse{}, err
	}

	health := DashboardHealth{
		Env:      s.cfg.App.Env,
		Database: "ok",
		Redis:    "ok",
	}

	if err := appRedis.Ping(s.redis); err != nil {
		health.Redis = "error"
		s.log.Warn("dashboard redis ping failed", zap.Error(err))
	}

	return DashboardResponse{
		CurrentUser:      currentUser,
		Health:           health,
		Metrics:          metrics,
		RecentOperations: recentOperations,
		RecentLogins:     recentLogins,
		LatestNotices:    latestNotices,
	}, nil
}

func (s *DashboardService) loadCurrentUser(userID uint, fallbackUsername string) (DashboardCurrentUser, error) {
	user, err := s.repo.FindUserProfileByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return DashboardCurrentUser{}, apperror.Unauthorized("登录状态无效，请重新登录")
		}
		return DashboardCurrentUser{}, apperror.Internal("查询当前用户失败", err)
	}

	if user.Username == "" {
		user.Username = fallbackUsername
	}

	return user, nil
}

func (s *DashboardService) loadMetrics() (DashboardMetrics, error) {
	now := time.Now()
	dayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	userTotal, err := s.repo.CountUsers()
	if err != nil {
		return DashboardMetrics{}, apperror.Internal("查询用户总数失败", err)
	}
	enabledUsers, err := s.repo.CountEnabledUsers()
	if err != nil {
		return DashboardMetrics{}, apperror.Internal("查询启用用户总数失败", err)
	}
	enabledRoles, err := s.repo.CountEnabledRoles()
	if err != nil {
		return DashboardMetrics{}, apperror.Internal("查询启用角色总数失败", err)
	}
	configTotal, err := s.repo.CountEnabledConfigs()
	if err != nil {
		return DashboardMetrics{}, apperror.Internal("查询配置总数失败", err)
	}
	noticeTotal, err := s.repo.CountEnabledNotices()
	if err != nil {
		return DashboardMetrics{}, apperror.Internal("查询公告总数失败", err)
	}
	fileTotal, err := s.repo.CountFiles()
	if err != nil {
		return DashboardMetrics{}, apperror.Internal("查询文件总数失败", err)
	}
	todayOperations, err := s.repo.CountTodayOperations(dayStart)
	if err != nil {
		return DashboardMetrics{}, apperror.Internal("查询今日操作总数失败", err)
	}
	todayRiskOperations, err := s.repo.CountTodayRiskOperations(dayStart)
	if err != nil {
		return DashboardMetrics{}, apperror.Internal("查询今日失败操作总数失败", err)
	}
	todayLoginFailed, err := s.repo.CountTodayLoginFailures(dayStart)
	if err != nil {
		return DashboardMetrics{}, apperror.Internal("查询今日登录失败总数失败", err)
	}

	return DashboardMetrics{
		UserTotal:               userTotal,
		EnabledUserTotal:        enabledUsers,
		EnabledRoleTotal:        enabledRoles,
		ConfigTotal:             configTotal,
		NoticeTotal:             noticeTotal,
		FileTotal:               fileTotal,
		TodayOperationTotal:     todayOperations,
		TodayRiskOperationTotal: todayRiskOperations,
		TodayLoginFailedTotal:   todayLoginFailed,
	}, nil
}

func (s *DashboardService) loadRecentOperations() ([]DashboardOperationItem, error) {
	rows, err := s.repo.ListRecentOperations(6)
	if err != nil {
		return nil, apperror.Internal("查询最近操作失败", err)
	}

	items := make([]DashboardOperationItem, 0, len(rows))
	for _, item := range rows {
		items = append(items, DashboardOperationItem{
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

func (s *DashboardService) loadRecentLogins() ([]DashboardLoginItem, error) {
	rows, err := s.repo.ListRecentLogins(5)
	if err != nil {
		return nil, apperror.Internal("查询最近登录记录失败", err)
	}

	items := make([]DashboardLoginItem, 0, len(rows))
	for _, item := range rows {
		items = append(items, DashboardLoginItem{
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

func (s *DashboardService) loadLatestNotices() ([]DashboardNoticeItem, error) {
	rows, err := s.repo.ListLatestEnabledNotices(3)
	if err != nil {
		return nil, apperror.Internal("查询最近公告失败", err)
	}

	items := make([]DashboardNoticeItem, 0, len(rows))
	for _, item := range rows {
		items = append(items, DashboardNoticeItem{
			ID:        item.ID,
			Title:     item.Title,
			Status:    item.Status,
			UpdatedAt: item.UpdatedAt,
		})
	}

	return items, nil
}
