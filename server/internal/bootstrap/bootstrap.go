package bootstrap

import (
	"errors"
	"fmt"

	"ez-admin-gin/server/internal/model"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	defaultAdminUsername = "admin"
	defaultAdminPassword = "EzAdmin@123456"
	defaultAdminRoleCode = "super_admin"
	defaultAdminRoleName = "超级管理员"
)

// Run 执行服务启动时必须完成的初始化动作。
func Run(db *gorm.DB, log *zap.Logger) error {
	admin, err := seedDefaultAdmin(db, log)
	if err != nil {
		return fmt.Errorf("seed default admin: %w", err)
	}

	role, err := seedSuperAdminRole(db, log)
	if err != nil {
		return fmt.Errorf("seed super admin role: %w", err)
	}

	if err := seedAdminRole(db, admin.ID, role.ID, log); err != nil {
		return fmt.Errorf("seed admin role: %w", err)
	}

	return nil
}

// seedDefaultAdmin 创建本地起步用的默认管理员。
func seedDefaultAdmin(db *gorm.DB, log *zap.Logger) (*model.User, error) {
	var user model.User
	// Unscoped 会把已逻辑删除记录也查出来，避免重复创建同名默认账号。
	err := db.Unscoped().Where("username = ?", defaultAdminUsername).First(&user).Error
	if err == nil {
		return &user, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(defaultAdminPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash default admin password: %w", err)
	}

	user = model.User{
		Username:     defaultAdminUsername,
		PasswordHash: string(passwordHash),
		Nickname:     "系统管理员",
		Status:       model.UserStatusEnabled,
	}

	if err := db.Create(&user).Error; err != nil {
		return nil, err
	}

	log.Info("default admin user created", zap.String("username", defaultAdminUsername))

	return &user, nil
}

// seedSuperAdminRole 创建超级管理员角色。
func seedSuperAdminRole(db *gorm.DB, log *zap.Logger) (*model.Role, error) {
	var role model.Role
	// 角色编码唯一，查询历史记录可以避免逻辑删除后重复创建同名编码。
	err := db.Unscoped().Where("code = ?", defaultAdminRoleCode).First(&role).Error
	if err == nil {
		return &role, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	role = model.Role{
		Code:   defaultAdminRoleCode,
		Name:   defaultAdminRoleName,
		Sort:   0,
		Status: model.RoleStatusEnabled,
		Remark: "系统内置角色",
	}

	if err := db.Create(&role).Error; err != nil {
		return nil, err
	}

	log.Info("default admin role created", zap.String("role_code", defaultAdminRoleCode))

	return &role, nil
}

// seedAdminRole 绑定默认管理员和超级管理员角色。
func seedAdminRole(db *gorm.DB, userID uint, roleID uint, log *zap.Logger) error {
	var userRole model.UserRole
	err := db.Where("user_id = ? AND role_id = ?", userID, roleID).First(&userRole).Error
	if err == nil {
		return nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	userRole = model.UserRole{
		UserID: userID,
		RoleID: roleID,
	}

	if err := db.Create(&userRole).Error; err != nil {
		return err
	}

	log.Info(
		"default admin role bound",
		zap.Uint("user_id", userID),
		zap.Uint("role_id", roleID),
	)

	return nil
}
