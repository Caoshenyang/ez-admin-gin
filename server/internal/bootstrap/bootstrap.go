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
	defaultAdminUsername      = "admin"
	defaultAdminPassword      = "EzAdmin@123456"
	defaultAdminRoleCode      = "super_admin"
	defaultAdminRoleName      = "超级管理员"
	defaultPermissionPath     = "/api/v1/system/health"
	defaultPermissionMethod   = "GET"
	defaultSystemMenuCode     = "system"
	defaultHealthMenuCode     = "system:health"
	defaultHealthViewCode     = "system:health:view"
	defaultUserMenuCode       = "system:user"
	defaultUserListCode       = "system:user:list"
	defaultUserCreateCode     = "system:user:create"
	defaultUserUpdateCode     = "system:user:update"
	defaultUserStatusCode     = "system:user:status"
	defaultUserAssignRoleCode = "system:user:assign-role"
	defaultRoleMenuCode       = "system:role"
	defaultRoleListCode       = "system:role:list"
	defaultRoleCreateCode     = "system:role:create"
	defaultRoleUpdateCode     = "system:role:update"
	defaultRoleStatusCode     = "system:role:status"
	defaultRolePermissionCode = "system:role:permission"
	defaultRoleMenuAssignCode = "system:role:menu"
	defaultMenuManageCode     = "system:menu"
	defaultMenuListCode       = "system:menu:list"
	defaultMenuCreateCode     = "system:menu:create"
	defaultMenuUpdateCode     = "system:menu:update"
	defaultMenuStatusCode     = "system:menu:status"
	defaultMenuDeleteCode     = "system:menu:delete"
	defaultConfigMenuCode     = "system:config"
	defaultConfigListCode     = "system:config:list"
	defaultConfigCreateCode   = "system:config:create"
	defaultConfigUpdateCode   = "system:config:update"
	defaultConfigStatusCode   = "system:config:status"
	defaultConfigValueCode    = "system:config:value"
	defaultFileMenuCode       = "system:file"
	defaultFileListCode       = "system:file:list"
	defaultFileUploadCode     = "system:file:upload"
)

type defaultPermissionSeed struct {
	Path   string
	Method string
}

var defaultPermissionSeeds = []defaultPermissionSeed{
	{Path: "/api/v1/system/health", Method: "GET"},
	{Path: "/api/v1/system/users", Method: "GET"},
	{Path: "/api/v1/system/users", Method: "POST"},
	{Path: "/api/v1/system/users/:id", Method: "PUT"},
	{Path: "/api/v1/system/users/:id/status", Method: "PATCH"},
	{Path: "/api/v1/system/users/:id/roles", Method: "PUT"},
	{Path: "/api/v1/system/roles", Method: "GET"},
	{Path: "/api/v1/system/roles", Method: "POST"},
	{Path: "/api/v1/system/roles/:id", Method: "PUT"},
	{Path: "/api/v1/system/roles/:id/status", Method: "PATCH"},
	{Path: "/api/v1/system/roles/:id/permissions", Method: "PUT"},
	{Path: "/api/v1/system/roles/:id/menus", Method: "PUT"},
	{Path: "/api/v1/system/menus", Method: "GET"},
	{Path: "/api/v1/system/menus", Method: "POST"},
	{Path: "/api/v1/system/menus/:id", Method: "PUT"},
	{Path: "/api/v1/system/menus/:id/status", Method: "PATCH"},
	{Path: "/api/v1/system/menus/:id", Method: "DELETE"},
	{Path: "/api/v1/system/configs", Method: "GET"},
	{Path: "/api/v1/system/configs", Method: "POST"},
	{Path: "/api/v1/system/configs/:id", Method: "PUT"},
	{Path: "/api/v1/system/configs/:id/status", Method: "PATCH"},
	{Path: "/api/v1/system/configs/value/:key", Method: "GET"},
	{Path: "/api/v1/system/files", Method: "GET"},
	{Path: "/api/v1/system/files", Method: "POST"},
}

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

	if err := seedDefaultPermissions(db, log); err != nil {
		return fmt.Errorf("seed default permissions: %w", err)
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

	userMenu, err := seedMenu(db, model.Menu{
		ParentID:  systemMenu.ID,
		Type:      model.MenuTypeMenu,
		Code:      defaultUserMenuCode,
		Title:     "用户管理",
		Path:      "/system/users",
		Component: "system/UserView",
		Icon:      "user",
		Sort:      20,
		Status:    model.MenuStatusEnabled,
		Remark:    "系统内置菜单",
	}, log)
	if err != nil {
		return nil, err
	}

	userButtons := []model.Menu{
		{ParentID: userMenu.ID, Type: model.MenuTypeButton, Code: defaultUserListCode, Title: "查看用户", Sort: 10, Status: model.MenuStatusEnabled, Remark: "系统内置按钮"},
		{ParentID: userMenu.ID, Type: model.MenuTypeButton, Code: defaultUserCreateCode, Title: "创建用户", Sort: 20, Status: model.MenuStatusEnabled, Remark: "系统内置按钮"},
		{ParentID: userMenu.ID, Type: model.MenuTypeButton, Code: defaultUserUpdateCode, Title: "编辑用户", Sort: 30, Status: model.MenuStatusEnabled, Remark: "系统内置按钮"},
		{ParentID: userMenu.ID, Type: model.MenuTypeButton, Code: defaultUserStatusCode, Title: "修改用户状态", Sort: 40, Status: model.MenuStatusEnabled, Remark: "系统内置按钮"},
		{ParentID: userMenu.ID, Type: model.MenuTypeButton, Code: defaultUserAssignRoleCode, Title: "分配用户角色", Sort: 50, Status: model.MenuStatusEnabled, Remark: "系统内置按钮"},
	}

	menus := []model.Menu{*systemMenu, *healthMenu, *healthViewButton, *userMenu}
	for _, button := range userButtons {
		createdButton, err := seedMenu(db, button, log)
		if err != nil {
			return nil, err
		}
		menus = append(menus, *createdButton)
	}

	roleMenu, err := seedMenu(db, model.Menu{
		ParentID:  systemMenu.ID,
		Type:      model.MenuTypeMenu,
		Code:      defaultRoleMenuCode,
		Title:     "角色管理",
		Path:      "/system/roles",
		Component: "system/RoleView",
		Icon:      "team",
		Sort:      30,
		Status:    model.MenuStatusEnabled,
		Remark:    "系统内置菜单",
	}, log)
	if err != nil {
		return nil, err
	}

	roleButtons := []model.Menu{
		{ParentID: roleMenu.ID, Type: model.MenuTypeButton, Code: defaultRoleListCode, Title: "查看角色", Sort: 10, Status: model.MenuStatusEnabled, Remark: "系统内置按钮"},
		{ParentID: roleMenu.ID, Type: model.MenuTypeButton, Code: defaultRoleCreateCode, Title: "创建角色", Sort: 20, Status: model.MenuStatusEnabled, Remark: "系统内置按钮"},
		{ParentID: roleMenu.ID, Type: model.MenuTypeButton, Code: defaultRoleUpdateCode, Title: "编辑角色", Sort: 30, Status: model.MenuStatusEnabled, Remark: "系统内置按钮"},
		{ParentID: roleMenu.ID, Type: model.MenuTypeButton, Code: defaultRoleStatusCode, Title: "修改角色状态", Sort: 40, Status: model.MenuStatusEnabled, Remark: "系统内置按钮"},
		{ParentID: roleMenu.ID, Type: model.MenuTypeButton, Code: defaultRolePermissionCode, Title: "分配接口权限", Sort: 50, Status: model.MenuStatusEnabled, Remark: "系统内置按钮"},
		{ParentID: roleMenu.ID, Type: model.MenuTypeButton, Code: defaultRoleMenuAssignCode, Title: "分配菜单权限", Sort: 60, Status: model.MenuStatusEnabled, Remark: "系统内置按钮"},
	}

	menus = append(menus, *roleMenu)
	for _, button := range roleButtons {
		createdButton, err := seedMenu(db, button, log)
		if err != nil {
			return nil, err
		}
		menus = append(menus, *createdButton)
	}

	menuManage, err := seedMenu(db, model.Menu{
		ParentID:  systemMenu.ID,
		Type:      model.MenuTypeMenu,
		Code:      defaultMenuManageCode,
		Title:     "菜单管理",
		Path:      "/system/menus",
		Component: "system/MenuView",
		Icon:      "menu",
		Sort:      40,
		Status:    model.MenuStatusEnabled,
		Remark:    "系统内置菜单",
	}, log)
	if err != nil {
		return nil, err
	}

	menuButtons := []model.Menu{
		{ParentID: menuManage.ID, Type: model.MenuTypeButton, Code: defaultMenuListCode, Title: "查看菜单", Sort: 10, Status: model.MenuStatusEnabled, Remark: "系统内置按钮"},
		{ParentID: menuManage.ID, Type: model.MenuTypeButton, Code: defaultMenuCreateCode, Title: "创建菜单", Sort: 20, Status: model.MenuStatusEnabled, Remark: "系统内置按钮"},
		{ParentID: menuManage.ID, Type: model.MenuTypeButton, Code: defaultMenuUpdateCode, Title: "编辑菜单", Sort: 30, Status: model.MenuStatusEnabled, Remark: "系统内置按钮"},
		{ParentID: menuManage.ID, Type: model.MenuTypeButton, Code: defaultMenuStatusCode, Title: "修改菜单状态", Sort: 40, Status: model.MenuStatusEnabled, Remark: "系统内置按钮"},
		{ParentID: menuManage.ID, Type: model.MenuTypeButton, Code: defaultMenuDeleteCode, Title: "删除菜单", Sort: 50, Status: model.MenuStatusEnabled, Remark: "系统内置按钮"},
	}

	menus = append(menus, *menuManage)
	for _, button := range menuButtons {
		createdButton, err := seedMenu(db, button, log)
		if err != nil {
			return nil, err
		}
		menus = append(menus, *createdButton)
	}

	configMenu, err := seedMenu(db, model.Menu{
		ParentID:  systemMenu.ID,
		Type:      model.MenuTypeMenu,
		Code:      defaultConfigMenuCode,
		Title:     "系统配置",
		Path:      "/system/configs",
		Component: "system/ConfigView",
		Icon:      "tool",
		Sort:      50,
		Status:    model.MenuStatusEnabled,
		Remark:    "系统内置菜单",
	}, log)
	if err != nil {
		return nil, err
	}

	configButtons := []model.Menu{
		{ParentID: configMenu.ID, Type: model.MenuTypeButton, Code: defaultConfigListCode, Title: "查看配置", Sort: 10, Status: model.MenuStatusEnabled, Remark: "系统内置按钮"},
		{ParentID: configMenu.ID, Type: model.MenuTypeButton, Code: defaultConfigCreateCode, Title: "创建配置", Sort: 20, Status: model.MenuStatusEnabled, Remark: "系统内置按钮"},
		{ParentID: configMenu.ID, Type: model.MenuTypeButton, Code: defaultConfigUpdateCode, Title: "编辑配置", Sort: 30, Status: model.MenuStatusEnabled, Remark: "系统内置按钮"},
		{ParentID: configMenu.ID, Type: model.MenuTypeButton, Code: defaultConfigStatusCode, Title: "修改配置状态", Sort: 40, Status: model.MenuStatusEnabled, Remark: "系统内置按钮"},
		{ParentID: configMenu.ID, Type: model.MenuTypeButton, Code: defaultConfigValueCode, Title: "读取配置值", Sort: 50, Status: model.MenuStatusEnabled, Remark: "系统内置按钮"},
	}

	menus = append(menus, *configMenu)
	for _, button := range configButtons {
		createdButton, err := seedMenu(db, button, log)
		if err != nil {
			return nil, err
		}
		menus = append(menus, *createdButton)
	}

	fileMenu, err := seedMenu(db, model.Menu{
		ParentID:  systemMenu.ID,
		Type:      model.MenuTypeMenu,
		Code:      defaultFileMenuCode,
		Title:     "文件管理",
		Path:      "/system/files",
		Component: "system/FileView",
		Icon:      "folder",
		Sort:      60,
		Status:    model.MenuStatusEnabled,
		Remark:    "系统内置菜单",
	}, log)
	if err != nil {
		return nil, err
	}

	fileButtons := []model.Menu{
		{ParentID: fileMenu.ID, Type: model.MenuTypeButton, Code: defaultFileListCode, Title: "查看文件", Sort: 10, Status: model.MenuStatusEnabled, Remark: "系统内置按钮"},
		{ParentID: fileMenu.ID, Type: model.MenuTypeButton, Code: defaultFileUploadCode, Title: "上传文件", Sort: 20, Status: model.MenuStatusEnabled, Remark: "系统内置按钮"},
	}

	menus = append(menus, *fileMenu)
	for _, button := range fileButtons {
		createdButton, err := seedMenu(db, button, log)
		if err != nil {
			return nil, err
		}
		menus = append(menus, *createdButton)
	}

	return menus, nil

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

// seedDefaultPermissions 初始化超级管理员的默认接口权限。
func seedDefaultPermissions(db *gorm.DB, log *zap.Logger) error {
	for _, permission := range defaultPermissionSeeds {
		var rule model.CasbinRule
		err := db.Where(
			"ptype = ? AND v0 = ? AND v1 = ? AND v2 = ?",
			"p",
			defaultAdminRoleCode,
			permission.Path,
			permission.Method,
		).First(&rule).Error
		if err == nil {
			continue
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		rule = model.CasbinRule{
			Ptype: "p",
			V0:    defaultAdminRoleCode,
			V1:    permission.Path,
			V2:    permission.Method,
		}

		if err := db.Create(&rule).Error; err != nil {
			return err
		}

		log.Info(
			"default permission created",
			zap.String("role_code", defaultAdminRoleCode),
			zap.String("path", permission.Path),
			zap.String("method", permission.Method),
		)
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
