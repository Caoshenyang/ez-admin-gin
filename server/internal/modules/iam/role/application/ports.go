package application

import (
	roledomain "ez-admin-gin/server/internal/modules/iam/role/domain"
	"ez-admin-gin/server/internal/platform/database"
	"ez-admin-gin/server/internal/platform/model"

	"gorm.io/gorm"
)

type RoleTransactor = database.Transactor

type RoleRepository interface {
	List(query roledomain.ListQuery, page int, pageSize int) ([]model.Role, int64, error)
	FindByID(db *gorm.DB, roleID uint) (model.Role, error)
	CodeExists(db *gorm.DB, code string) (bool, error)
	DepartmentsUsable(db *gorm.DB, departmentIDs []uint) error
	MenusUsable(db *gorm.DB, menuIDs []uint) error
	Create(db *gorm.DB, role *model.Role) error
	UpdateBase(db *gorm.DB, role *model.Role, req roledomain.UpdateRequest) error
	UpdateStatus(db *gorm.DB, role *model.Role, status model.RoleStatus) error
	RolePermissions(roleCodes []string) (map[string][]roledomain.PermissionItem, error)
	RoleMenuIDs(roleIDs []uint) (map[uint][]uint, error)
	RoleCustomDepartmentIDs(roleIDs []uint) (map[uint][]uint, error)
	ReplacePermissions(db *gorm.DB, roleCode string, permissions []roledomain.PermissionItem) error
	ReplaceMenus(db *gorm.DB, roleID uint, menuIDs []uint) error
	ReplaceCustomDepartments(db *gorm.DB, roleID uint, departmentIDs []uint) error
}
