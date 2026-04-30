package config

import (
	"errors"
	"strings"

	"ez-admin-gin/server/internal/apperror"
	"ez-admin-gin/server/internal/model"

	"gorm.io/gorm"
)

// Repository 负责系统配置的查询和持久化。
type Repository struct {
	db *gorm.DB
}

// NewRepository 创建配置仓储。
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// List 返回配置分页结果和总数。
func (r *Repository) List(query ListQuery, page int, pageSize int) ([]model.SystemConfig, int64, error) {
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
		if !ValidStatus(status) {
			return nil, 0, apperror.BadRequest("配置状态不正确")
		}
		queryDB = queryDB.Where("status = ?", status)
	}

	var total int64
	if err := queryDB.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var items []model.SystemConfig
	if err := queryDB.
		Order("group_code ASC, sort ASC, id ASC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

// FindByID 查询单个配置项。
func (r *Repository) FindByID(db *gorm.DB, configID uint) (model.SystemConfig, error) {
	var item model.SystemConfig
	err := db.First(&item, configID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.SystemConfig{}, apperror.NotFound("配置不存在")
		}
		return model.SystemConfig{}, err
	}

	return item, nil
}

// FindEnabledByKey 查询启用中的配置项。
func (r *Repository) FindEnabledByKey(key string) (model.SystemConfig, error) {
	var item model.SystemConfig
	err := r.db.Where("config_key = ?", key).
		Where("status = ?", model.SystemConfigStatusEnabled).
		First(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.SystemConfig{}, apperror.NotFound("配置不存在或已禁用")
		}
		return model.SystemConfig{}, err
	}

	return item, nil
}

// KeyExists 判断配置键是否已存在。
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

// Create 创建配置项。
func (r *Repository) Create(db *gorm.DB, item *model.SystemConfig) error {
	return db.Create(item).Error
}

// UpdateBase 更新配置基础字段。
func (r *Repository) UpdateBase(db *gorm.DB, item *model.SystemConfig, req UpdateRequest) error {
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

// UpdateStatus 单独更新配置状态。
func (r *Repository) UpdateStatus(db *gorm.DB, item *model.SystemConfig, status model.SystemConfigStatus) error {
	if err := db.Model(item).Update("status", status).Error; err != nil {
		return err
	}
	item.Status = status
	return nil
}
