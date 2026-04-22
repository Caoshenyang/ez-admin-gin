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
	defaultAdminUsername    = "admin"
	defaultAdminPassword    = "EzAdmin@123456"
	defaultAdminRoleCode    = "super_admin"
	defaultAdminRoleName    = "超级管理员"
	defaultPermissionPath   = "/api/v1/system/health"
	defaultPermissionMethod = "GET"
	defaultSystemMenuCode   = "system"
	defaultHealthMenuCode   = "system:health"
	defaultHealthViewCode   = "system:health:view"
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

	menus, err := seedDefaultMenus(db, log)
	if err != nil {
		return fmt.Errorf("seed default menus: %w", err)
	}

	if err := seedDefaultPermission(db, log); err != nil {
		return fmt.Errorf("seed default permission: %w", err)
	}

	if err := seedAdminRole(db, admin.ID, role.ID, log); err != nil {
		return fmt.Errorf("seed admin role: %w", err)
	}

	if err := seedRoleMenus(db, role.ID, menus, log); err != nil {
		return fmt.Errorf("seed role menus: %w", err)
	}

	return nil
}

// seedDefaultMenus 初始化默认菜单和按钮。
func seedDefaultMenus(db *gorm.DB, log *zap.Logger) ([]model.Menu, error) {
	systemMenu, err := seedMenu(db, model.Menu{
		ParentID: 0,
		Type:     model.MenuTypeDirectory,
		Code:     defaultSystemMenuCode,
		Title:    "系统管理",
		Path:     "/system",
		Icon:     "setting",
		Sort:     10,
		Status:   model.MenuStatusEnabled,
		Remark:   "系统内置目录",
	}, log)
	if err != nil {
		return nil, err
	}

	healthMenu, err := seedMenu(db, model.Menu{
		ParentID:  systemMenu.ID,
		Type:      model.MenuTypeMenu,
		Code:      defaultHealthMenuCode,
		Title:     "系统状态",
		Path:      "/system/health",
		Component: "system/HealthView",
		Icon:      "monitor",
		Sort:      10,
		Status:    model.MenuStatusEnabled,
		Remark:    "系统内置菜单",
	}, log)
	if err != nil {
		return nil, err
	}

	healthViewButton, err := seedMenu(db, model.Menu{
		ParentID: healthMenu.ID,
		Type:     model.MenuTypeButton,
		Code:     defaultHealthViewCode,
		Title:    "查看系统状态",
		Sort:     10,
		Status:   model.MenuStatusEnabled,
		Remark:   "系统内置按钮",
	}, log)
	if err != nil {
		return nil, err
	}

	return []model.Menu{*systemMenu, *healthMenu, *healthViewButton}, nil
}

// seedMenu 按菜单编码创建默认菜单。
func seedMenu(db *gorm.DB, menu model.Menu, log *zap.Logger) (*model.Menu, error) {
	var exists model.Menu
	err := db.Unscoped().Where("code = ?", menu.Code).First(&exists).Error
	if err == nil {
		return &exists, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if err := db.Create(&menu).Error; err != nil {
		return nil, err
	}

	log.Info("default menu created", zap.String("menu_code", menu.Code))

	return &menu, nil
}

// seedRoleMenus 把默认菜单授权给指定角色。
func seedRoleMenus(db *gorm.DB, roleID uint, menus []model.Menu, log *zap.Logger) error {
	for _, menu := range menus {
		var roleMenu model.RoleMenu
		err := db.Where("role_id = ? AND menu_id = ?", roleID, menu.ID).First(&roleMenu).Error
		if err == nil {
			continue
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		roleMenu = model.RoleMenu{
			RoleID: roleID,
			MenuID: menu.ID,
		}

		if err := db.Create(&roleMenu).Error; err != nil {
			return err
		}

		log.Info(
			"default role menu bound",
			zap.Uint("role_id", roleID),
			zap.Uint("menu_id", menu.ID),
		)
	}

	return nil
}

// seedDefaultPermission 初始化超级管理员的默认接口权限。
func seedDefaultPermission(db *gorm.DB, log *zap.Logger) error {
	var rule model.CasbinRule
	err := db.Where(
		"ptype = ? AND v0 = ? AND v1 = ? AND v2 = ?",
		"p",
		defaultAdminRoleCode,
		defaultPermissionPath,
		defaultPermissionMethod,
	).First(&rule).Error
	if err == nil {
		return nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	rule = model.CasbinRule{
		Ptype: "p",
		V0:    defaultAdminRoleCode,
		V1:    defaultPermissionPath,
		V2:    defaultPermissionMethod,
	}

	if err := db.Create(&rule).Error; err != nil {
		return err
	}

	log.Info(
		"default permission created",
		zap.String("role_code", defaultAdminRoleCode),
		zap.String("path", defaultPermissionPath),
		zap.String("method", defaultPermissionMethod),
	)

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
