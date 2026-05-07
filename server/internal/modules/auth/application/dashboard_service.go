package application

import (
	"errors"
	"time"

	authdomain "ez-admin-gin/server/internal/modules/auth/domain"
	authinfra "ez-admin-gin/server/internal/modules/auth/infra"
	errorsx "ez-admin-gin/server/internal/pkg/errorsx"
	platformConfig "ez-admin-gin/server/internal/platform/config"
	platformDatabase "ez-admin-gin/server/internal/platform/database"
	platformRedis "ez-admin-gin/server/internal/platform/redis"

	goredis "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type DashboardService struct {
	cfg   *platformConfig.Config
	db    *gorm.DB
	repo  *authinfra.Repository
	redis *goredis.Client
	log   *zap.Logger
}

func NewDashboardService(
	cfg *platformConfig.Config,
	db *gorm.DB,
	repo *authinfra.Repository,
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

func (s *DashboardService) Dashboard(userID uint, fallbackUsername string) (authdomain.DashboardResponse, error) {
	if err := platformDatabase.Ping(s.db); err != nil {
		return authdomain.DashboardResponse{}, errorsx.ServiceUnavailable("数据库不可用", err)
	}

	currentUser, err := s.loadCurrentUser(userID, fallbackUsername)
	if err != nil {
		return authdomain.DashboardResponse{}, err
	}

	metrics, err := s.loadMetrics()
	if err != nil {
		return authdomain.DashboardResponse{}, err
	}

	recentOperations, err := s.loadRecentOperations()
	if err != nil {
		return authdomain.DashboardResponse{}, err
	}

	recentLogins, err := s.loadRecentLogins()
	if err != nil {
		return authdomain.DashboardResponse{}, err
	}

	latestNotices, err := s.loadLatestNotices()
	if err != nil {
		return authdomain.DashboardResponse{}, err
	}

	health := authdomain.DashboardHealth{
		Env:      s.cfg.App.Env,
		Database: "ok",
		Redis:    "ok",
	}

	if err := platformRedis.Ping(s.redis); err != nil {
		health.Redis = "error"
		s.log.Warn("dashboard redis ping failed", zap.Error(err))
	}

	return authdomain.DashboardResponse{
		CurrentUser:      currentUser,
		Health:           health,
		Metrics:          metrics,
		RecentOperations: recentOperations,
		RecentLogins:     recentLogins,
		LatestNotices:    latestNotices,
	}, nil
}

func (s *DashboardService) loadCurrentUser(userID uint, fallbackUsername string) (authdomain.DashboardCurrentUser, error) {
	user, err := s.repo.FindUserProfileByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return authdomain.DashboardCurrentUser{}, errorsx.Unauthorized("登录状态无效，请重新登录")
		}
		return authdomain.DashboardCurrentUser{}, errorsx.Internal("查询当前用户失败", err)
	}

	if user.Username == "" {
		user.Username = fallbackUsername
	}

	return user, nil
}

func (s *DashboardService) loadMetrics() (authdomain.DashboardMetrics, error) {
	now := time.Now()
	dayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	userTotal, err := s.repo.CountUsers()
	if err != nil {
		return authdomain.DashboardMetrics{}, errorsx.Internal("查询用户总数失败", err)
	}
	enabledUsers, err := s.repo.CountEnabledUsers()
	if err != nil {
		return authdomain.DashboardMetrics{}, errorsx.Internal("查询启用用户总数失败", err)
	}
	enabledRoles, err := s.repo.CountEnabledRoles()
	if err != nil {
		return authdomain.DashboardMetrics{}, errorsx.Internal("查询启用角色总数失败", err)
	}
	configTotal, err := s.repo.CountEnabledConfigs()
	if err != nil {
		return authdomain.DashboardMetrics{}, errorsx.Internal("查询配置总数失败", err)
	}
	noticeTotal, err := s.repo.CountEnabledNotices()
	if err != nil {
		return authdomain.DashboardMetrics{}, errorsx.Internal("查询公告总数失败", err)
	}
	fileTotal, err := s.repo.CountFiles()
	if err != nil {
		return authdomain.DashboardMetrics{}, errorsx.Internal("查询文件总数失败", err)
	}
	todayOperations, err := s.repo.CountTodayOperations(dayStart)
	if err != nil {
		return authdomain.DashboardMetrics{}, errorsx.Internal("查询今日操作总数失败", err)
	}
	todayRiskOperations, err := s.repo.CountTodayRiskOperations(dayStart)
	if err != nil {
		return authdomain.DashboardMetrics{}, errorsx.Internal("查询今日失败操作总数失败", err)
	}
	todayLoginFailed, err := s.repo.CountTodayLoginFailures(dayStart)
	if err != nil {
		return authdomain.DashboardMetrics{}, errorsx.Internal("查询今日登录失败总数失败", err)
	}

	return authdomain.DashboardMetrics{
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

func (s *DashboardService) loadRecentOperations() ([]authdomain.DashboardOperationItem, error) {
	rows, err := s.repo.ListRecentOperations(6)
	if err != nil {
		return nil, errorsx.Internal("查询最近操作失败", err)
	}

	items := make([]authdomain.DashboardOperationItem, 0, len(rows))
	for _, item := range rows {
		items = append(items, authdomain.DashboardOperationItem{
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

func (s *DashboardService) loadRecentLogins() ([]authdomain.DashboardLoginItem, error) {
	rows, err := s.repo.ListRecentLogins(5)
	if err != nil {
		return nil, errorsx.Internal("查询最近登录记录失败", err)
	}

	items := make([]authdomain.DashboardLoginItem, 0, len(rows))
	for _, item := range rows {
		items = append(items, authdomain.DashboardLoginItem{
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

func (s *DashboardService) loadLatestNotices() ([]authdomain.DashboardNoticeItem, error) {
	rows, err := s.repo.ListLatestEnabledNotices(3)
	if err != nil {
		return nil, errorsx.Internal("查询最近公告失败", err)
	}

	items := make([]authdomain.DashboardNoticeItem, 0, len(rows))
	for _, item := range rows {
		items = append(items, authdomain.DashboardNoticeItem{
			ID:        item.ID,
			Title:     item.Title,
			Status:    item.Status,
			UpdatedAt: item.UpdatedAt,
		})
	}

	return items, nil
}
