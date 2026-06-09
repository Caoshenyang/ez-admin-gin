// Package infra 实现操作日志的数据访问层。
package infra

import (
	"strings"

	operationlogdomain "ez-admin-gin/server/internal/modules/system/operationlog/domain"

	"gorm.io/gorm"
)

// Repository 封装操作日志表的数据访问操作。
type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// List 按用户名、方法、路径和成功状态分页查询操作日志。
func (r *Repository) List(query operationlogdomain.ListQuery, page int, pageSize int, success *bool) ([]operationlogdomain.Entity, int64, error) {
	queryDB := r.db.Model(&operationlogdomain.Entity{})

	username := strings.TrimSpace(query.Username)
	if username != "" {
		queryDB = queryDB.Where("username = ?", username)
	}
	method := strings.ToUpper(strings.TrimSpace(query.Method))
	if method != "" {
		queryDB = queryDB.Where("method = ?", method)
	}
	path := strings.TrimSpace(query.Path)
	if path != "" {
		queryDB = queryDB.Where("path LIKE ?", "%"+path+"%")
	}
	if success != nil {
		queryDB = queryDB.Where("success = ?", *success)
	}

	var total int64
	if err := queryDB.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var items []operationlogdomain.Entity
	if err := queryDB.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}
