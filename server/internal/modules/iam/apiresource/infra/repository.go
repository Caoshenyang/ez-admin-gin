// Package infra 实现接口权限元数据的数据访问层。
package infra

import (
	"ez-admin-gin/server/internal/platform/model"

	"gorm.io/gorm"
)

// Repository 实现接口权限元数据的数据访问层。
type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// List 返回所有接口权限元数据记录。
func (r *Repository) List() ([]model.API, error) {
	var items []model.API
	if err := r.db.Order("module ASC, sort ASC, id ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}
