package loginlog

import (
	"ez-admin-gin/server/internal/model"

	"gorm.io/gorm"
)

// Repository 负责登录日志的查询读取。
type Repository struct {
	db *gorm.DB
}

// NewRepository 创建登录日志仓储。
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// List 返回登录日志分页结果和总数。
func (r *Repository) List(query ListQuery, page int, pageSize int, status *model.LoginLogStatus) ([]Entity, int64, error) {
	queryDB := r.db.Model(&Entity{})

	if username := NormalizeUsername(query.Username); username != "" {
		queryDB = queryDB.Where("username = ?", username)
	}

	if ip := NormalizeIP(query.IP); ip != "" {
		queryDB = queryDB.Where("ip = ?", ip)
	}

	if status != nil {
		queryDB = queryDB.Where("status = ?", *status)
	}

	var total int64
	if err := queryDB.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var items []Entity
	if err := queryDB.
		Order("id DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}
