package infra

import (
	"errors"
	"strings"

	configdomain "ez-admin-gin/server/internal/modules/system/config/domain"
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

func (r *Repository) List(query configdomain.ListQuery, page int, pageSize int) ([]model.SystemConfig, int64, error) {
	queryDB := r.db.Model(&model.SystemConfig{})

	keyword := strings.TrimSpace(query.Keyword)
	if keyword != "" {
		like := "%" + keyword + "%"
		queryDB = queryDB.Where("config_key LIKE ? OR name LIKE ?", like, like)
	}

	groupCode := strings.TrimSpace(query.GroupCode)
	if groupCode != "" {
		queryDB = queryDB.Where("group_code = ?", groupCode)
	}

	if query.Status != 0 {
		status := model.SystemConfigStatus(query.Status)
		if !configdomain.ValidStatus(status) {
			return nil, 0, errorsx.BadRequest("配置状态不正确")
		}
		queryDB = queryDB.Where("status = ?", status)
	}

	var total int64
	if err := queryDB.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var items []model.SystemConfig
	if err := queryDB.Order("group_code ASC, sort ASC, id ASC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (r *Repository) FindByID(db *gorm.DB, configID uint) (model.SystemConfig, error) {
	var item model.SystemConfig
	err := db.First(&item, configID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.SystemConfig{}, errorsx.NotFound("配置不存在")
		}
		return model.SystemConfig{}, err
	}

	return item, nil
}

func (r *Repository) FindEnabledByKey(key string) (model.SystemConfig, error) {
	var item model.SystemConfig
	err := r.db.Where("config_key = ?", key).Where("status = ?", model.SystemConfigStatusEnabled).First(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.SystemConfig{}, errorsx.NotFound("配置不存在或已禁用")
		}
		return model.SystemConfig{}, err
	}

	return item, nil
}

func (r *Repository) KeyExists(db *gorm.DB, key string) (bool, error) {
	var item model.SystemConfig
	err := db.Unscoped().Where("config_key = ?", key).First(&item).Error
	if err == nil {
		return true, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	return false, err
}

func (r *Repository) Create(db *gorm.DB, item *model.SystemConfig) error {
	return db.Create(item).Error
}

func (r *Repository) UpdateBase(db *gorm.DB, item *model.SystemConfig, req configdomain.UpdateRequest) error {
	if err := db.Model(item).Updates(map[string]any{
		"group_code": req.GroupCode,
		"name":       req.Name,
		"value":      req.Value,
		"sort":       req.Sort,
		"status":     req.Status,
		"remark":     req.Remark,
	}).Error; err != nil {
		return err
	}

	item.GroupCode = req.GroupCode
	item.Name = req.Name
	item.Value = req.Value
	item.Sort = req.Sort
	item.Status = req.Status
	item.Remark = req.Remark
	return nil
}

func (r *Repository) UpdateStatus(db *gorm.DB, item *model.SystemConfig, status model.SystemConfigStatus) error {
	if err := db.Model(item).Update("status", status).Error; err != nil {
		return err
	}
	item.Status = status
	return nil
}
