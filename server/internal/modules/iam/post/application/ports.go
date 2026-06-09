package application

import (
	postdomain "ez-admin-gin/server/internal/modules/iam/post/domain"
	"ez-admin-gin/server/internal/platform/database"
	"ez-admin-gin/server/internal/platform/model"

	"gorm.io/gorm"
)

// PostTransactor 是岗位模块使用的事务管理器类型别名。
type PostTransactor = database.Transactor

// PostRepository 定义岗位聚合根的数据访问接口。
type PostRepository interface {
	List(query postdomain.ListQuery) ([]model.Post, error)
	FindByID(db *gorm.DB, postID uint) (model.Post, error)
	CodeExists(db *gorm.DB, code string, excludeID uint) (bool, error)
	Create(db *gorm.DB, item *model.Post) error
	Update(db *gorm.DB, item *model.Post, code string, name string, sort int, status model.PostStatus, remark string) error
	UpdateStatus(db *gorm.DB, item *model.Post, status model.PostStatus) error
	CountUsers(db *gorm.DB, postID uint) (int64, error)
	Delete(db *gorm.DB, item *model.Post) error
}
