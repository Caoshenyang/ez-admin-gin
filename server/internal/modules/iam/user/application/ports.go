package application

import (
	"ez-admin-gin/server/internal/modules/iam/user/domain"
	"ez-admin-gin/server/internal/platform/database"
	"ez-admin-gin/server/internal/platform/datascope"
	"ez-admin-gin/server/internal/platform/model"

	"gorm.io/gorm"
)

// UserRepository 定义用户聚合根的数据访问接口。
type UserRepository interface {
	List(actor datascope.Actor, query domain.ListQuery, page int, pageSize int) ([]model.User, int64, error)
	RoleIDsByUserIDs(userIDs []uint) (map[uint][]uint, error)
	PostIDsByUserIDs(userIDs []uint) (map[uint][]uint, error)
	FindByIDInScope(db *gorm.DB, actor datascope.Actor, userID uint) (model.User, error)
	UsernameExists(db *gorm.DB, username string) (bool, error)
	DepartmentUsable(db *gorm.DB, departmentID uint) error
	RolesUsable(db *gorm.DB, roleIDs []uint) error
	PostsUsable(db *gorm.DB, postIDs []uint) error
	Create(db *gorm.DB, user *model.User) error
	UpdateBase(db *gorm.DB, user *model.User, nickname string, departmentID uint, status model.UserStatus) error
	UpdateStatus(db *gorm.DB, user *model.User, status model.UserStatus) error
	ReplaceRoles(db *gorm.DB, userID uint, roleIDs []uint) error
	ReplacePosts(db *gorm.DB, userID uint, postIDs []uint) error
}

// UserTransactor 是用户模块使用的事务管理器类型别名。
type UserTransactor = database.Transactor
