package infra

import (
	"errors"
	"time"

	authdomain "ez-admin-gin/server/internal/modules/auth/domain"
	"ez-admin-gin/server/internal/platform/model"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
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

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindUserByUsername(username string) (model.User, error) {
	var user model.User
	err := r.db.Where("username = ?", username).First(&user).Error
	return user, err
}

func (r *Repository) CreateLoginLog(record *model.LoginLog) error {
	return r.db.Create(record).Error
}

func (r *Repository) FindUserProfileByID(userID uint) (authdomain.DashboardCurrentUser, error) {
	var user model.User
	err := r.db.Select("id", "username", "nickname").First(&user, userID).Error
	if err != nil {
		return authdomain.DashboardCurrentUser{}, err
	}

	return authdomain.DashboardCurrentUser{
		UserID:   user.ID,
		Username: user.Username,
		Nickname: user.Nickname,
	}, nil
}

func (r *Repository) FindAccountProfileByID(userID uint) (AccountProfileRow, error) {
	var row AccountProfileRow
	err := r.db.
		Table("sys_user AS u").
		Select("u.id, u.username, u.nickname, u.department_id, u.status, u.updated_at, COALESCE(d.name, '') AS department_name").
		Joins("LEFT JOIN sys_department AS d ON d.id = u.department_id AND d.deleted_at IS NULL").
		Where("u.id = ?", userID).
		Where("u.deleted_at IS NULL").
		Scan(&row).Error
	if err != nil {
		return AccountProfileRow{}, err
	}
	if row.ID == 0 {
		return AccountProfileRow{}, gorm.ErrRecordNotFound
	}

	return row, nil
}

func (r *Repository) FindUserByID(tx *gorm.DB, userID uint) (model.User, error) {
	var user model.User
	err := r.dbOr(tx).First(&user, userID).Error
	return user, err
}

func (r *Repository) UpdateAccountNickname(tx *gorm.DB, user *model.User, nickname string) error {
	user.Nickname = nickname
	return r.dbOr(tx).Model(user).Update("nickname", nickname).Error
}

func (r *Repository) UpdateAccountPasswordHash(tx *gorm.DB, user *model.User, passwordHash string) error {
	user.PasswordHash = passwordHash
	return r.dbOr(tx).Model(user).Update("password_hash", passwordHash).Error
}

func (r *Repository) ListMenusByUserID(userID uint) ([]model.Menu, error) {
	var menus []model.Menu
	err := r.db.
		Table("sys_menu AS m").
		Select("DISTINCT m.*").
		Joins("JOIN sys_role_menu AS rm ON rm.menu_id = m.id").
		Joins("JOIN sys_user_role AS ur ON ur.role_id = rm.role_id").
		Joins("JOIN sys_role AS r ON r.id = ur.role_id").
		Where("ur.user_id = ?", userID).
		Where("m.status = ?", model.MenuStatusEnabled).
		Where("r.status = ?", model.RoleStatusEnabled).
		Where("m.deleted_at IS NULL").
		Where("r.deleted_at IS NULL").
		Order("m.sort ASC, m.id ASC").
		Find(&menus).Error
	return menus, err
}

func (r *Repository) CountUsers() (int64, error) {
	var total int64
	err := r.db.Model(&model.User{}).Count(&total).Error
	return total, err
}

func (r *Repository) CountEnabledUsers() (int64, error) {
	var total int64
	err := r.db.Model(&model.User{}).Where("status = ?", model.UserStatusEnabled).Count(&total).Error
	return total, err
}

func (r *Repository) CountEnabledRoles() (int64, error) {
	var total int64
	err := r.db.Model(&model.Role{}).Where("status = ?", model.RoleStatusEnabled).Count(&total).Error
	return total, err
}

func (r *Repository) CountEnabledConfigs() (int64, error) {
	var total int64
	err := r.db.Model(&model.SystemConfig{}).Where("status = ?", model.SystemConfigStatusEnabled).Count(&total).Error
	return total, err
}

func (r *Repository) CountEnabledNotices() (int64, error) {
	var total int64
	err := r.db.Model(&model.Notice{}).Where("status = ?", model.NoticeStatusEnabled).Count(&total).Error
	return total, err
}

func (r *Repository) CountFiles() (int64, error) {
	var total int64
	err := r.db.Model(&model.SystemFile{}).Count(&total).Error
	return total, err
}

func (r *Repository) CountTodayOperations(dayStart time.Time) (int64, error) {
	var total int64
	err := r.db.Model(&model.OperationLog{}).Where("created_at >= ?", dayStart).Count(&total).Error
	return total, err
}

func (r *Repository) CountTodayRiskOperations(dayStart time.Time) (int64, error) {
	var total int64
	err := r.db.Model(&model.OperationLog{}).
		Where("created_at >= ?", dayStart).
		Where("success = ?", false).
		Count(&total).Error
	return total, err
}

func (r *Repository) CountTodayLoginFailures(dayStart time.Time) (int64, error) {
	var total int64
	err := r.db.Model(&model.LoginLog{}).
		Where("created_at >= ?", dayStart).
		Where("status = ?", model.LoginLogStatusFailed).
		Count(&total).Error
	return total, err
}

func (r *Repository) ListRecentOperations(limit int) ([]model.OperationLog, error) {
	var rows []model.OperationLog
	err := r.db.Order("id DESC").Limit(limit).Find(&rows).Error
	return rows, err
}

func (r *Repository) ListRecentLogins(limit int) ([]model.LoginLog, error) {
	var rows []model.LoginLog
	err := r.db.Order("id DESC").Limit(limit).Find(&rows).Error
	return rows, err
}

func (r *Repository) ListLatestEnabledNotices(limit int) ([]model.Notice, error) {
	var rows []model.Notice
	err := r.db.
		Where("status = ?", model.NoticeStatusEnabled).
		Order("updated_at DESC, id DESC").
		Limit(limit).
		Find(&rows).Error
	return rows, err
}

func (r *Repository) IsNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

func (r *Repository) dbOr(tx *gorm.DB) *gorm.DB {
	if tx != nil {
		return tx
	}
	return r.db
}
