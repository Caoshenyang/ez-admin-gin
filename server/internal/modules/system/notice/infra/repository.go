// Package infra 实现公告的数据访问层。
package infra

import (
	"errors"
	"strings"

	noticedomain "ez-admin-gin/server/internal/modules/system/notice/domain"
	errorsx "ez-admin-gin/server/internal/pkg/errorsx"
	"ez-admin-gin/server/internal/platform/model"

	"gorm.io/gorm"
)

// Repository 封装公告表的数据访问操作。
type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// List 按关键词和状态分页查询公告列表。
func (r *Repository) List(query noticedomain.ListQuery, page int, pageSize int, status *model.NoticeStatus) ([]noticedomain.Entity, int64, error) {
	queryDB := r.db.Model(&noticedomain.Entity{})

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

	var items []noticedomain.Entity
	if err := queryDB.Order("sort ASC, id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

// FindByID 在指定事务中按主键查找公告，不存在时返回 NotFound 错误。
func (r *Repository) FindByID(db *gorm.DB, noticeID uint) (noticedomain.Entity, error) {
	var item noticedomain.Entity
	err := db.First(&item, noticeID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return noticedomain.Entity{}, errorsx.NotFound("公告不存在")
		}
		return noticedomain.Entity{}, err
	}
	return item, nil
}

// Create 在指定事务中插入一条新的公告记录。
func (r *Repository) Create(db *gorm.DB, item *noticedomain.Entity) error {
	return db.Create(item).Error
}

// UpdateBase 更新公告的基本字段（标题、内容、排序、状态、备注）。
func (r *Repository) UpdateBase(db *gorm.DB, item *noticedomain.Entity, req noticedomain.UpdateRequest) error {
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

// UpdateStatus 更新公告的状态字段。
func (r *Repository) UpdateStatus(db *gorm.DB, item *noticedomain.Entity, status model.NoticeStatus) error {
	if err := db.Model(item).Update("status", status).Error; err != nil {
		return err
	}
	item.Status = status
	return nil
}
