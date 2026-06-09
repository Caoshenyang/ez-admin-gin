package application

import (
	menudomain "ez-admin-gin/server/internal/modules/iam/menu/domain"
	"ez-admin-gin/server/internal/platform/database"
	"ez-admin-gin/server/internal/platform/model"

	"gorm.io/gorm"
)

// MenuTransactor 是菜单模块使用的事务管理器类型别名。
type MenuTransactor = database.Transactor

// MenuRepository 定义菜单聚合根的数据访问接口。
type MenuRepository interface {
	List() ([]model.Menu, error)
	FindByID(db *gorm.DB, menuID uint) (model.Menu, error)
	CodeExists(db *gorm.DB, code string) (bool, error)
	ParentUsable(db *gorm.DB, parentID uint, menuType model.MenuType, excludeID uint) error
	Create(db *gorm.DB, menu *model.Menu) error
	UpdateBase(db *gorm.DB, menu *model.Menu, req menudomain.UpdateRequest) error
	UpdateStatus(db *gorm.DB, menu *model.Menu, status model.MenuStatus) error
	CanDelete(db *gorm.DB, menuID uint) error
	Delete(db *gorm.DB, menu *model.Menu) error
}
