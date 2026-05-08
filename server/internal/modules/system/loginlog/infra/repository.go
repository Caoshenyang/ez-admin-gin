// Package infra 实现登录日志的数据访问层。
package infra

import (
	loginlogdomain "ez-admin-gin/server/internal/modules/system/loginlog/domain"
	"ez-admin-gin/server/internal/platform/model"

	"gorm.io/gorm"
)

// Repository 封装登录日志表的数据访问操作。
type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// List 按用户名、IP 和状态分页查询登录日志。
func (r *Repository) List(query loginlogdomain.ListQuery, page int, pageSize int, status *model.LoginLogStatus) ([]loginlogdomain.Entity, int64, error) {
	queryDB := r.db.Model(&loginlogdomain.Entity{})

	if username := loginlogdomain.NormalizeUsername(query.Username); username != "" {
		queryDB = queryDB.Where("username = ?", username)
	}
	if ip := loginlogdomain.NormalizeIP(query.IP); ip != "" {
		queryDB = queryDB.Where("ip = ?", ip)
	}
	if status != nil {
		queryDB = queryDB.Where("status = ?", *status)
	}

	var total int64
	if err := queryDB.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var items []loginlogdomain.Entity
	if err := queryDB.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}
