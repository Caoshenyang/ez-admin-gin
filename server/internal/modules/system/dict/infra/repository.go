// Package infra 实现字典的数据访问层。
package infra

import (
	"errors"
	"strings"

	dictdomain "ez-admin-gin/server/internal/modules/system/dict/domain"
	errorsx "ez-admin-gin/server/internal/pkg/errorsx"
	"ez-admin-gin/server/internal/platform/model"

	"gorm.io/gorm"
)

// Repository 封装字典类型和字典项表的数据访问操作。
type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// ListTypes 按关键词和状态分页查询字典类型列表。
func (r *Repository) ListTypes(query dictdomain.TypeListQuery, page int, pageSize int, status *model.SystemDictStatus) ([]dictdomain.DictTypeEntity, int64, error) {
	queryDB := r.db.Model(&dictdomain.DictTypeEntity{})

	keyword := strings.TrimSpace(query.Keyword)
	if keyword != "" {
		like := "%" + keyword + "%"
		queryDB = queryDB.Where("code LIKE ? OR name LIKE ?", like, like)
	}
	if status != nil {
		queryDB = queryDB.Where("status = ?", *status)
	}

	var total int64
	if err := queryDB.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var items []dictdomain.DictTypeEntity
	if err := queryDB.Order("sort ASC, id ASC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

// FindTypeByID 在指定事务中按主键查找字典类型，不存在时返回 NotFound 错误。
func (r *Repository) FindTypeByID(db *gorm.DB, typeID uint) (dictdomain.DictTypeEntity, error) {
	var item dictdomain.DictTypeEntity
	err := db.First(&item, typeID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dictdomain.DictTypeEntity{}, errorsx.NotFound("字典类型不存在")
		}
		return dictdomain.DictTypeEntity{}, err
	}

	return item, nil
}

// TypeCodeExists 检查字典类型编码是否已存在（包含已软删除的记录）。
func (r *Repository) TypeCodeExists(db *gorm.DB, code string) (bool, error) {
	var item dictdomain.DictTypeEntity
	err := db.Unscoped().Where("code = ?", code).First(&item).Error
	if err == nil {
		return true, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	return false, err
}

// CreateType 在指定事务中插入一条新的字典类型记录。
func (r *Repository) CreateType(db *gorm.DB, item *dictdomain.DictTypeEntity) error {
	return db.Create(item).Error
}

// UpdateTypeBase 更新字典类型的基本字段（名称、排序、状态、备注）。
func (r *Repository) UpdateTypeBase(db *gorm.DB, item *dictdomain.DictTypeEntity, req dictdomain.UpdateTypeRequest) error {
	if err := db.Model(item).Updates(map[string]any{
		"name":   req.Name,
		"sort":   req.Sort,
		"status": req.Status,
		"remark": req.Remark,
	}).Error; err != nil {
		return err
	}
	item.Name = req.Name
	item.Sort = req.Sort
	item.Status = req.Status
	item.Remark = req.Remark
	return nil
}

// UpdateTypeStatus 更新字典类型的状态字段。
func (r *Repository) UpdateTypeStatus(db *gorm.DB, item *dictdomain.DictTypeEntity, status model.SystemDictStatus) error {
	if err := db.Model(item).Update("status", status).Error; err != nil {
		return err
	}
	item.Status = status
	return nil
}

// CountItemsByType 统计指定字典类型下仍存在的字典项数量。
func (r *Repository) CountItemsByType(db *gorm.DB, typeID uint) (int64, error) {
	var count int64
	err := db.Model(&dictdomain.DictItemEntity{}).Where("type_id = ?", typeID).Count(&count).Error
	return count, err
}

// DeleteType 软删除字典类型记录。
func (r *Repository) DeleteType(db *gorm.DB, item *dictdomain.DictTypeEntity) error {
	return db.Delete(item).Error
}

// ListItems 按类型 ID、关键词和状态分页查询字典项列表。
func (r *Repository) ListItems(query dictdomain.ItemListQuery, page int, pageSize int, status *model.SystemDictStatus) ([]dictdomain.DictItemEntity, int64, error) {
	queryDB := r.db.Model(&dictdomain.DictItemEntity{}).Where("type_id = ?", query.TypeID)

	keyword := strings.TrimSpace(query.Keyword)
	if keyword != "" {
		like := "%" + keyword + "%"
		queryDB = queryDB.Where("item_key LIKE ? OR label LIKE ? OR value LIKE ?", like, like, like)
	}
	if status != nil {
		queryDB = queryDB.Where("status = ?", *status)
	}

	var total int64
	if err := queryDB.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var items []dictdomain.DictItemEntity
	if err := queryDB.Order("sort ASC, id ASC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

// FindItemByID 在指定事务中按主键查找字典项，不存在时返回 NotFound 错误。
func (r *Repository) FindItemByID(db *gorm.DB, itemID uint) (dictdomain.DictItemEntity, error) {
	var item dictdomain.DictItemEntity
	err := db.First(&item, itemID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dictdomain.DictItemEntity{}, errorsx.NotFound("字典项不存在")
		}
		return dictdomain.DictItemEntity{}, err
	}
	return item, nil
}

// ItemKeyExists 检查同一字典类型下的字典项编码是否已存在（包含已软删除的记录）。
func (r *Repository) ItemKeyExists(db *gorm.DB, typeID uint, itemKey string) (bool, error) {
	var item dictdomain.DictItemEntity
	err := db.Unscoped().Where("type_id = ? AND item_key = ?", typeID, itemKey).First(&item).Error
	if err == nil {
		return true, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	return false, err
}

// CreateItem 在指定事务中插入一条新的字典项记录。
func (r *Repository) CreateItem(db *gorm.DB, item *dictdomain.DictItemEntity) error {
	return db.Create(item).Error
}

// UpdateItemBase 更新字典项的基本字段（标签、值、标签类型、排序、状态、备注）。
func (r *Repository) UpdateItemBase(db *gorm.DB, item *dictdomain.DictItemEntity, req dictdomain.UpdateItemRequest) error {
	if err := db.Model(item).Updates(map[string]any{
		"label":    req.Label,
		"value":    req.Value,
		"tag_type": req.TagType,
		"sort":     req.Sort,
		"status":   req.Status,
		"remark":   req.Remark,
	}).Error; err != nil {
		return err
	}
	item.Label = req.Label
	item.Value = req.Value
	item.TagType = req.TagType
	item.Sort = req.Sort
	item.Status = req.Status
	item.Remark = req.Remark
	return nil
}

// UpdateItemStatus 更新字典项的状态字段。
func (r *Repository) UpdateItemStatus(db *gorm.DB, item *dictdomain.DictItemEntity, status model.SystemDictStatus) error {
	if err := db.Model(item).Update("status", status).Error; err != nil {
		return err
	}
	item.Status = status
	return nil
}

// DeleteItem 软删除字典项记录。
func (r *Repository) DeleteItem(db *gorm.DB, item *dictdomain.DictItemEntity) error {
	return db.Delete(item).Error
}
