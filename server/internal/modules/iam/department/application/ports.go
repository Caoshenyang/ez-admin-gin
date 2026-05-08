package application

import (
	departmentdomain "ez-admin-gin/server/internal/modules/iam/department/domain"
	"ez-admin-gin/server/internal/platform/database"
	"ez-admin-gin/server/internal/platform/datascope"
	"ez-admin-gin/server/internal/platform/model"

	"gorm.io/gorm"
)

// DepartmentTransactor 是部门模块使用的事务管理器类型别名。
type DepartmentTransactor = database.Transactor

// DepartmentRepository 定义部门聚合根的数据访问接口。
type DepartmentRepository interface {
	List(actor datascope.Actor, query departmentdomain.ListQuery) ([]model.Department, error)
	FindByIDInScope(db *gorm.DB, actor datascope.Actor, departmentID uint) (model.Department, error)
	FindByID(db *gorm.DB, departmentID uint) (model.Department, error)
	FindParent(db *gorm.DB, parentID uint) (model.Department, error)
	CodeExists(db *gorm.DB, code string, excludeID uint) (bool, error)
	LeaderUsable(db *gorm.DB, leaderUserID uint) error
	Create(db *gorm.DB, department *model.Department) error
	Update(db *gorm.DB, department *model.Department, parentID uint, ancestors string, name string, code string, leaderUserID uint, sort int, status model.DepartmentStatus, remark string) error
	UpdateStatus(db *gorm.DB, department *model.Department, status model.DepartmentStatus) error
	Subtree(db *gorm.DB, departmentID uint, fullPath string) ([]model.Department, error)
	UpdateAncestors(db *gorm.DB, departmentID uint, ancestors string) error
}
