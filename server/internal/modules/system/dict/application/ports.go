package application

import (
	dictdomain "ez-admin-gin/server/internal/modules/system/dict/domain"
	"ez-admin-gin/server/internal/platform/database"
	"ez-admin-gin/server/internal/platform/model"

	"gorm.io/gorm"
)

type DictTransactor = database.Transactor

type DictRepository interface {
	ListTypes(query dictdomain.TypeListQuery, page int, pageSize int, status *model.SystemDictStatus) ([]dictdomain.DictTypeEntity, int64, error)
	FindTypeByID(db *gorm.DB, typeID uint) (dictdomain.DictTypeEntity, error)
	TypeCodeExists(db *gorm.DB, code string) (bool, error)
	CreateType(db *gorm.DB, item *dictdomain.DictTypeEntity) error
	UpdateTypeBase(db *gorm.DB, item *dictdomain.DictTypeEntity, req dictdomain.UpdateTypeRequest) error
	UpdateTypeStatus(db *gorm.DB, item *dictdomain.DictTypeEntity, status model.SystemDictStatus) error
	CountItemsByType(db *gorm.DB, typeID uint) (int64, error)
	DeleteType(db *gorm.DB, item *dictdomain.DictTypeEntity) error
	ListItems(query dictdomain.ItemListQuery, page int, pageSize int, status *model.SystemDictStatus) ([]dictdomain.DictItemEntity, int64, error)
	FindItemByID(db *gorm.DB, itemID uint) (dictdomain.DictItemEntity, error)
	ItemKeyExists(db *gorm.DB, typeID uint, itemKey string) (bool, error)
	CreateItem(db *gorm.DB, item *dictdomain.DictItemEntity) error
	UpdateItemBase(db *gorm.DB, item *dictdomain.DictItemEntity, req dictdomain.UpdateItemRequest) error
	UpdateItemStatus(db *gorm.DB, item *dictdomain.DictItemEntity, status model.SystemDictStatus) error
	DeleteItem(db *gorm.DB, item *dictdomain.DictItemEntity) error
}
