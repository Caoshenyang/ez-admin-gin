// Package infra 实现 auth 模块的数据访问层。
package infra

import (
	"errors"
	"time"

	authapp "ez-admin-gin/server/internal/modules/auth/application"
	authdomain "ez-admin-gin/server/internal/modules/auth/domain"
	"ez-admin-gin/server/internal/platform/model"

	"gorm.io/gorm"
)

// Repository 封装 auth 模块所有数据库查询操作。
type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// FindUserByUsername 根据用户名查询用户。
func (r *Repository) FindUserByUsername(username string) (model.User, error) {
	var user model.User
	err := r.db.Where("username = ?", username).First(&user).Error
	return user, err
}

// FindUserByIDSimple 根据ID查询用户（无事务参数，用于 refresh service）。
func (r *Repository) FindUserByIDSimple(userID uint) (model.User, error) {
	var user model.User
	err := r.db.First(&user, userID).Error
	return user, err
}

// CreateLoginLog 写入一条登录日志记录。
func (r *Repository) CreateLoginLog(record *model.LoginLog) error {
	return r.db.Create(record).Error
}

// FindUserProfileByID 根据 ID 查询仪表盘所需的用户简要信息。
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

// FindAccountProfileByID 根据 ID 查询账户中心资料，含部门名称。
func (r *Repository) FindAccountProfileByID(userID uint) (authapp.AccountProfileRow, error) {
	var row authapp.AccountProfileRow
	err := r.db.
		Table("sys_user AS u").
		Select("u.id, u.username, u.nickname, u.department_id, u.status, u.updated_at, COALESCE(d.name, '') AS department_name").
		Joins("LEFT JOIN sys_department AS d ON d.id = u.department_id AND d.deleted_at IS NULL").
		Where("u.id = ?", userID).
		Where("u.deleted_at IS NULL").
		Scan(&row).Error
	if err != nil {
		return authapp.AccountProfileRow{}, err
	}
	if row.ID == 0 {
		return authapp.AccountProfileRow{}, gorm.ErrRecordNotFound
	}

	return row, nil
}

// FindUserByID 根据 ID 查询用户，支持在事务中执行。
func (r *Repository) FindUserByID(tx *gorm.DB, userID uint) (model.User, error) {
	var user model.User
	err := r.dbOr(tx).First(&user, userID).Error
	return user, err
}

// UpdateAccountNickname 更新用户昵称，支持在事务中执行。
func (r *Repository) UpdateAccountNickname(tx *gorm.DB, user *model.User, nickname string) error {
	user.Nickname = nickname
	return r.dbOr(tx).Model(user).Update("nickname", nickname).Error
}

// UpdateAccountPasswordHash 更新用户密码哈希，支持在事务中执行。
func (r *Repository) UpdateAccountPasswordHash(tx *gorm.DB, user *model.User, passwordHash string) error {
	user.PasswordHash = passwordHash
	return r.dbOr(tx).Model(user).Update("password_hash", passwordHash).Error
}

// ListMenusByUserID 查询指定用户通过角色关联的已启用菜单列表。
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

// CountUsers 统计用户总数。
func (r *Repository) CountUsers() (int64, error) {
	var total int64
	err := r.db.Model(&model.User{}).Count(&total).Error
	return total, err
}

// CountEnabledUsers 统计已启用用户总数。
func (r *Repository) CountEnabledUsers() (int64, error) {
	var total int64
	err := r.db.Model(&model.User{}).Where("status = ?", model.UserStatusEnabled).Count(&total).Error
	return total, err
}

// CountEnabledRoles 统计已启用角色总数。
func (r *Repository) CountEnabledRoles() (int64, error) {
	var total int64
	err := r.db.Model(&model.Role{}).Where("status = ?", model.RoleStatusEnabled).Count(&total).Error
	return total, err
}

// CountEnabledConfigs 统计已启用系统配置总数。
func (r *Repository) CountEnabledConfigs() (int64, error) {
	var total int64
	err := r.db.Model(&model.SystemConfig{}).Where("status = ?", model.SystemConfigStatusEnabled).Count(&total).Error
	return total, err
}

// CountEnabledNotices 统计已启用公告总数。
func (r *Repository) CountEnabledNotices() (int64, error) {
	var total int64
	err := r.db.Model(&model.Notice{}).Where("status = ?", model.NoticeStatusEnabled).Count(&total).Error
	return total, err
}

// CountFiles 统计文件总数。
func (r *Repository) CountFiles() (int64, error) {
	var total int64
	err := r.db.Model(&model.SystemFile{}).Count(&total).Error
	return total, err
}

// CountTodayOperations 统计从指定时间起的操作日志总数。
func (r *Repository) CountTodayOperations(dayStart time.Time) (int64, error) {
	var total int64
	err := r.db.Model(&model.OperationLog{}).Where("created_at >= ?", dayStart).Count(&total).Error
	return total, err
}

// CountTodayRiskOperations 统计从指定时间起的失败操作总数。
func (r *Repository) CountTodayRiskOperations(dayStart time.Time) (int64, error) {
	var total int64
	err := r.db.Model(&model.OperationLog{}).
		Where("created_at >= ?", dayStart).
		Where("success = ?", false).
		Count(&total).Error
	return total, err
}

// CountTodayLoginFailures 统计从指定时间起的登录失败总数。
func (r *Repository) CountTodayLoginFailures(dayStart time.Time) (int64, error) {
	var total int64
	err := r.db.Model(&model.LoginLog{}).
		Where("created_at >= ?", dayStart).
		Where("status = ?", model.LoginLogStatusFailed).
		Count(&total).Error
	return total, err
}

// ListRecentOperations 查询最近 N 条操作日志。
func (r *Repository) ListRecentOperations(limit int) ([]model.OperationLog, error) {
	var rows []model.OperationLog
	err := r.db.Order("id DESC").Limit(limit).Find(&rows).Error
	return rows, err
}

// ListRecentLogins 查询最近 N 条登录日志。
func (r *Repository) ListRecentLogins(limit int) ([]model.LoginLog, error) {
	var rows []model.LoginLog
	err := r.db.Order("id DESC").Limit(limit).Find(&rows).Error
	return rows, err
}

// ListLatestEnabledNotices 查询最近 N 条已启用公告。
func (r *Repository) ListLatestEnabledNotices(limit int) ([]model.Notice, error) {
	var rows []model.Notice
	err := r.db.
		Where("status = ?", model.NoticeStatusEnabled).
		Order("updated_at DESC, id DESC").
		Limit(limit).
		Find(&rows).Error
	return rows, err
}

// IsNotFound 判断错误是否为记录未找到。
func (r *Repository) IsNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

func (r *Repository) dbOr(tx *gorm.DB) *gorm.DB {
	if tx != nil {
		return tx
	}
	return r.db
}
