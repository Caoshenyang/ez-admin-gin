package notice

import (
	"errors"
	"strings"

	"ez-admin-gin/server/internal/apperror"
	"ez-admin-gin/server/internal/model"

	"gorm.io/gorm"
)

// Repository 负责公告的查询和持久化。
type Repository struct {
	db *gorm.DB
}

// NewRepository 创建公告仓储。
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// List 返回公告分页结果和总数。
func (r *Repository) List(query ListQuery, page int, pageSize int, status *model.NoticeStatus) ([]Entity, int64, error) {
	queryDB := r.db.Model(&Entity{})

	keyword := strings.TrimSpace(query.Keyword)
	if keyword != "" {
		queryDB = queryDB.Where("title LIKE ?", "%"+keyword+"%")
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
		Order("sort ASC, id DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

// FindByID 查询单个公告。
func (r *Repository) FindByID(db *gorm.DB, noticeID uint) (Entity, error) {
	var item Entity
	err := db.First(&item, noticeID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return Entity{}, apperror.NotFound("公告不存在")
		}
		return Entity{}, err
	}

	return item, nil
}

// Create 创建公告。
func (r *Repository) Create(db *gorm.DB, item *Entity) error {
	return db.Create(item).Error
}

// UpdateBase 更新公告基础字段。
func (r *Repository) UpdateBase(db *gorm.DB, item *Entity, req UpdateRequest) error {
	if err := db.Model(item).Updates(map[string]any{
		"title":   req.Title,
		"content": req.Content,
		"sort":    req.Sort,
		"status":  req.Status,
		"remark":  req.Remark,
	}).Error; err != nil {
		return err
	}

	item.Title = req.Title
	item.Content = req.Content
	item.Sort = req.Sort
	item.Status = req.Status
	item.Remark = req.Remark
	return nil
}

// UpdateStatus 单独更新公告状态。
func (r *Repository) UpdateStatus(db *gorm.DB, item *Entity, status model.NoticeStatus) error {
	if err := db.Model(item).Update("status", status).Error; err != nil {
		return err
	}
	item.Status = status
	return nil
}
