package dict

import (
	"errors"
	"strings"

	"ez-admin-gin/server/internal/apperror"
	"ez-admin-gin/server/internal/model"

	"gorm.io/gorm"
)

// Repository 负责字典类型和字典项的查询与持久化。
type Repository struct {
	db *gorm.DB
}

// NewRepository 创建字典仓储。
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// ListTypes 返回字典类型分页结果和总数。
func (r *Repository) ListTypes(query TypeListQuery, page int, pageSize int, status *model.SystemDictStatus) ([]DictTypeEntity, int64, error) {
	queryDB := r.db.Model(&DictTypeEntity{})

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

	var items []DictTypeEntity
	if err := queryDB.
		Order("sort ASC, id ASC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

// FindTypeByID 查询单个字典类型。
func (r *Repository) FindTypeByID(db *gorm.DB, typeID uint) (DictTypeEntity, error) {
	var item DictTypeEntity
	err := db.First(&item, typeID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return DictTypeEntity{}, apperror.NotFound("字典类型不存在")
		}
		return DictTypeEntity{}, err
	}

	return item, nil
}

// TypeCodeExists 判断字典编码是否已存在。
func (r *Repository) TypeCodeExists(db *gorm.DB, code string) (bool, error) {
	var item DictTypeEntity
	err := db.Unscoped().Where("code = ?", code).First(&item).Error
	if err == nil {
		return true, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	return false, err
}

// CreateType 创建字典类型。
func (r *Repository) CreateType(db *gorm.DB, item *DictTypeEntity) error {
	return db.Create(item).Error
}

// UpdateTypeBase 更新字典类型基础字段。
func (r *Repository) UpdateTypeBase(db *gorm.DB, item *DictTypeEntity, req UpdateTypeRequest) error {
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

// UpdateTypeStatus 单独更新字典类型状态。
func (r *Repository) UpdateTypeStatus(db *gorm.DB, item *DictTypeEntity, status model.SystemDictStatus) error {
	if err := db.Model(item).Update("status", status).Error; err != nil {
		return err
	}
	item.Status = status
	return nil
}

// ListItems 返回字典项分页结果和总数。
func (r *Repository) ListItems(query ItemListQuery, page int, pageSize int, status *model.SystemDictStatus) ([]DictItemEntity, int64, error) {
	queryDB := r.db.Model(&DictItemEntity{}).Where("type_id = ?", query.TypeID)

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

	var items []DictItemEntity
	if err := queryDB.
		Order("sort ASC, id ASC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

// FindItemByID 查询单个字典项。
func (r *Repository) FindItemByID(db *gorm.DB, itemID uint) (DictItemEntity, error) {
	var item DictItemEntity
	err := db.First(&item, itemID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return DictItemEntity{}, apperror.NotFound("字典项不存在")
		}
		return DictItemEntity{}, err
	}

	return item, nil
}

// ItemKeyExists 判断字典项编码是否在类型内重复。
func (r *Repository) ItemKeyExists(db *gorm.DB, typeID uint, itemKey string) (bool, error) {
	var item DictItemEntity
	err := db.Unscoped().Where("type_id = ? AND item_key = ?", typeID, itemKey).First(&item).Error
	if err == nil {
		return true, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	return false, err
}

// CreateItem 创建字典项。
func (r *Repository) CreateItem(db *gorm.DB, item *DictItemEntity) error {
	return db.Create(item).Error
}

// UpdateItemBase 更新字典项基础字段。
func (r *Repository) UpdateItemBase(db *gorm.DB, item *DictItemEntity, req UpdateItemRequest) error {
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

// UpdateItemStatus 单独更新字典项状态。
func (r *Repository) UpdateItemStatus(db *gorm.DB, item *DictItemEntity, status model.SystemDictStatus) error {
	if err := db.Model(item).Update("status", status).Error; err != nil {
		return err
	}
	item.Status = status
	return nil
}
