package application

import (
	roledomain "ez-admin-gin/server/internal/modules/iam/role/domain"
	"ez-admin-gin/server/internal/platform/database"
	"ez-admin-gin/server/internal/platform/model"

	"gorm.io/gorm"
)

// RoleTransactor 是角色模块使用的事务管理器类型别名。
type RoleTransactor = database.Transactor

// PolicyReloader 在权限变更后刷新内存策略。
type PolicyReloader interface {
	ReloadPolicy() error
}

// RoleRepository 定义角色聚合根的数据访问接口。
type RoleRepository interface {
	List(query roledomain.ListQuery, page int, pageSize int) ([]model.Role, int64, error)
	FindByID(db *gorm.DB, roleID uint) (model.Role, error)
	CodeExists(db *gorm.DB, code string) (bool, error)
	DepartmentsUsable(db *gorm.DB, departmentIDs []uint) error
	MenusUsable(db *gorm.DB, menuIDs []uint) error
	APIsUsable(db *gorm.DB, apiIDs []uint) error
	Create(db *gorm.DB, role *model.Role) error
	UpdateBase(db *gorm.DB, role *model.Role, req roledomain.UpdateRequest) error
	UpdateStatus(db *gorm.DB, role *model.Role, status model.RoleStatus) error
	CountUsers(db *gorm.DB, roleID uint) (int64, error)
	RoleAPIIDs(roleIDs []uint) (map[uint][]uint, error)
	RolePermissions(roleIDs []uint) (map[uint][]roledomain.PermissionItem, error)
	RoleMenuIDs(roleIDs []uint) (map[uint][]uint, error)
	RoleCustomDepartmentIDs(roleIDs []uint) (map[uint][]uint, error)
	ReplaceAPIs(db *gorm.DB, roleID uint, apiIDs []uint) error
	ReplacePoliciesByAPIs(db *gorm.DB, roleCode string, apiIDs []uint) error
	ReplaceMenus(db *gorm.DB, roleID uint, menuIDs []uint) error
	ReplaceCustomDepartments(db *gorm.DB, roleID uint, departmentIDs []uint) error
	Delete(db *gorm.DB, role *model.Role) error
}
