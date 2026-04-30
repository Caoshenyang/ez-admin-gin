package operationlog

import (
	"strings"

	"gorm.io/gorm"
)

// Repository 负责操作日志的查询和持久化读取。
type Repository struct {
	db *gorm.DB
}

// NewRepository 创建操作日志仓储。
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// List 返回操作日志分页结果和总数。
func (r *Repository) List(query ListQuery, page int, pageSize int, success *bool) ([]Entity, int64, error) {
	queryDB := r.db.Model(&Entity{})

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
