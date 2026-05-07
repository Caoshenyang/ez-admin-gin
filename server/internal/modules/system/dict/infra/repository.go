package infra

import (
	"errors"
	"strings"

	dictdomain "ez-admin-gin/server/internal/modules/system/dict/domain"
	errorsx "ez-admin-gin/server/internal/pkg/errorsx"
	"ez-admin-gin/server/internal/platform/model"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

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

func (r *Repository) CreateType(db *gorm.DB, item *dictdomain.DictTypeEntity) error {
	return db.Create(item).Error
}

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

func (r *Repository) UpdateTypeStatus(db *gorm.DB, item *dictdomain.DictTypeEntity, status model.SystemDictStatus) error {
	if err := db.Model(item).Update("status", status).Error; err != nil {
		return err
	}
	item.Status = status
	return nil
}

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

func (r *Repository) CreateItem(db *gorm.DB, item *dictdomain.DictItemEntity) error {
	return db.Create(item).Error
}

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

func (r *Repository) UpdateItemStatus(db *gorm.DB, item *dictdomain.DictItemEntity, status model.SystemDictStatus) error {
	if err := db.Model(item).Update("status", status).Error; err != nil {
		return err
	}
	item.Status = status
	return nil
}
