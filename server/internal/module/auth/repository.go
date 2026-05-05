package auth

import (
	"errors"
	"time"

	"ez-admin-gin/server/internal/model"

	"gorm.io/gorm"
)

// Repository 负责认证模块的查询与持久化。
type Repository struct {
	db *gorm.DB
}

type accountProfileRow struct {
	ID             uint
	Username       string
	Nickname       string
	DepartmentID   uint
	DepartmentName string
	Status         model.UserStatus
	UpdatedAt      time.Time
}

// NewRepository 创建认证模块仓储。
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// FindUserByUsername 查询登录用户。
func (r *Repository) FindUserByUsername(username string) (model.User, error) {
	var user model.User
	err := r.db.Where("username = ?", username).First(&user).Error
	return user, err
}

// CreateLoginLog 写入登录日志。
func (r *Repository) CreateLoginLog(record *model.LoginLog) error {
	return r.db.Create(record).Error
}

// FindUserProfileByID 查询当前用户摘要。
func (r *Repository) FindUserProfileByID(userID uint) (DashboardCurrentUser, error) {
	var user model.User
	err := r.db.Select("id", "username", "nickname").First(&user, userID).Error
	if err != nil {
		return DashboardCurrentUser{}, err
	}

	return DashboardCurrentUser{
		UserID:   user.ID,
		Username: user.Username,
		Nickname: user.Nickname,
	}, nil
}

// FindAccountProfileByID 查询当前登录人账户中心资料。
func (r *Repository) FindAccountProfileByID(userID uint) (accountProfileRow, error) {
	var row accountProfileRow
	err := r.db.
		Table("sys_user AS u").
		Select(
			"u.id, u.username, u.nickname, u.department_id, u.status, u.updated_at, COALESCE(d.name, '') AS department_name",
		).
		Joins("LEFT JOIN sys_department AS d ON d.id = u.department_id AND d.deleted_at IS NULL").
		Where("u.id = ?", userID).
		Where("u.deleted_at IS NULL").
		Scan(&row).Error
	if err != nil {
		return accountProfileRow{}, err
	}
	if row.ID == 0 {
		return accountProfileRow{}, gorm.ErrRecordNotFound
	}

	return row, nil
}

// FindUserByID 查询指定用户。
func (r *Repository) FindUserByID(tx *gorm.DB, userID uint) (model.User, error) {
	var user model.User
	err := r.dbOr(tx).First(&user, userID).Error
	return user, err
}

// UpdateAccountNickname 更新当前登录人昵称。
func (r *Repository) UpdateAccountNickname(tx *gorm.DB, user *model.User, nickname string) error {
	user.Nickname = nickname
	return r.dbOr(tx).Model(user).Update("nickname", nickname).Error
}

// UpdateAccountPasswordHash 更新当前登录人密码哈希。
func (r *Repository) UpdateAccountPasswordHash(tx *gorm.DB, user *model.User, passwordHash string) error {
	user.PasswordHash = passwordHash
	return r.dbOr(tx).Model(user).Update("password_hash", passwordHash).Error
}

// ListMenusByUserID 查询当前用户可见菜单。
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

// CountUsers 返回用户总数。
func (r *Repository) CountUsers() (int64, error) {
	var total int64
	err := r.db.Model(&model.User{}).Count(&total).Error
	return total, err
}

// CountEnabledUsers 返回启用用户总数。
func (r *Repository) CountEnabledUsers() (int64, error) {
	var total int64
	err := r.db.Model(&model.User{}).
		Where("status = ?", model.UserStatusEnabled).
		Count(&total).Error
	return total, err
}

// CountEnabledRoles 返回启用角色总数。
func (r *Repository) CountEnabledRoles() (int64, error) {
	var total int64
	err := r.db.Model(&model.Role{}).
		Where("status = ?", model.RoleStatusEnabled).
		Count(&total).Error
	return total, err
}

// CountEnabledConfigs 返回启用配置总数。
func (r *Repository) CountEnabledConfigs() (int64, error) {
	var total int64
	err := r.db.Model(&model.SystemConfig{}).
		Where("status = ?", model.SystemConfigStatusEnabled).
		Count(&total).Error
	return total, err
}

// CountEnabledNotices 返回启用公告总数。
func (r *Repository) CountEnabledNotices() (int64, error) {
	var total int64
	err := r.db.Model(&model.Notice{}).
		Where("status = ?", model.NoticeStatusEnabled).
		Count(&total).Error
	return total, err
}

// CountFiles 返回文件总数。
func (r *Repository) CountFiles() (int64, error) {
	var total int64
	err := r.db.Model(&model.SystemFile{}).Count(&total).Error
	return total, err
}

// CountTodayOperations 返回今日操作总数。
func (r *Repository) CountTodayOperations(dayStart time.Time) (int64, error) {
	var total int64
	err := r.db.Model(&model.OperationLog{}).
		Where("created_at >= ?", dayStart).
		Count(&total).Error
	return total, err
}

// CountTodayRiskOperations 返回今日失败操作总数。
func (r *Repository) CountTodayRiskOperations(dayStart time.Time) (int64, error) {
	var total int64
	err := r.db.Model(&model.OperationLog{}).
		Where("created_at >= ?", dayStart).
		Where("success = ?", false).
		Count(&total).Error
	return total, err
}

// CountTodayLoginFailures 返回今日登录失败总数。
func (r *Repository) CountTodayLoginFailures(dayStart time.Time) (int64, error) {
	var total int64
	err := r.db.Model(&model.LoginLog{}).
		Where("created_at >= ?", dayStart).
		Where("status = ?", model.LoginLogStatusFailed).
		Count(&total).Error
	return total, err
}

// ListRecentOperations 返回最近操作记录。
func (r *Repository) ListRecentOperations(limit int) ([]model.OperationLog, error) {
	var rows []model.OperationLog
	err := r.db.Order("id DESC").Limit(limit).Find(&rows).Error
	return rows, err
}

// ListRecentLogins 返回最近登录记录。
func (r *Repository) ListRecentLogins(limit int) ([]model.LoginLog, error) {
	var rows []model.LoginLog
	err := r.db.Order("id DESC").Limit(limit).Find(&rows).Error
	return rows, err
}

// ListLatestEnabledNotices 返回最近启用公告。
func (r *Repository) ListLatestEnabledNotices(limit int) ([]model.Notice, error) {
	var rows []model.Notice
	err := r.db.
		Where("status = ?", model.NoticeStatusEnabled).
		Order("updated_at DESC, id DESC").
		Limit(limit).
		Find(&rows).Error
	return rows, err
}

// IsNotFound 判断错误是否为记录不存在。
func (r *Repository) IsNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

func (r *Repository) dbOr(tx *gorm.DB) *gorm.DB {
	if tx != nil {
		return tx
	}

	return r.db
}
