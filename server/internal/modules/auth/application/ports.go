package application

import (
	"context"
	"time"

	authdomain "ez-admin-gin/server/internal/modules/auth/domain"
	"ez-admin-gin/server/internal/platform/database"
	"ez-admin-gin/server/internal/platform/model"

	"gorm.io/gorm"
)

type AuthTransactor = database.Transactor

type TokenIssuer interface {
	GenerateAccessToken(userID uint, username string) (string, time.Time, error)
}

// SessionStore manages refresh token sessions (create, validate, revoke).
type SessionStore interface {
	Create(ctx context.Context, userID uint) (string, error)
	Validate(ctx context.Context, token string) (uint, error)
	Revoke(ctx context.Context, token string) error
	RevokeAllForUser(ctx context.Context, userID uint) error
}

type LoginRepository interface {
	FindUserByUsername(username string) (model.User, error)
	FindUserByIDSimple(userID uint) (model.User, error)
	CreateLoginLog(record *model.LoginLog) error
	IsNotFound(err error) bool
}

type AccountProfileRow struct {
	ID             uint
	Username       string
	Nickname       string
	DepartmentID   uint
	DepartmentName string
	Status         model.UserStatus
	UpdatedAt      time.Time
}

type AccountRepository interface {
	FindAccountProfileByID(userID uint) (AccountProfileRow, error)
	FindUserByID(tx *gorm.DB, userID uint) (model.User, error)
	UpdateAccountNickname(tx *gorm.DB, user *model.User, nickname string) error
	UpdateAccountPasswordHash(tx *gorm.DB, user *model.User, passwordHash string) error
	IsNotFound(err error) bool
}

type MenuRepository interface {
	ListMenusByUserID(userID uint) ([]model.Menu, error)
}

type DashboardRepository interface {
	FindUserProfileByID(userID uint) (authdomain.DashboardCurrentUser, error)
	CountUsers() (int64, error)
	CountEnabledUsers() (int64, error)
	CountEnabledRoles() (int64, error)
	CountEnabledConfigs() (int64, error)
	CountEnabledNotices() (int64, error)
	CountFiles() (int64, error)
	CountTodayOperations(dayStart time.Time) (int64, error)
	CountTodayRiskOperations(dayStart time.Time) (int64, error)
	CountTodayLoginFailures(dayStart time.Time) (int64, error)
	ListRecentOperations(limit int) ([]model.OperationLog, error)
	ListRecentLogins(limit int) ([]model.LoginLog, error)
	ListLatestEnabledNotices(limit int) ([]model.Notice, error)
}
