package infra

import (
	"strings"

	attachmentdomain "ez-admin-gin/server/internal/modules/system/attachment/domain"
	"ez-admin-gin/server/internal/platform/model"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) List(query attachmentdomain.ListQuery, page int, pageSize int) ([]attachmentdomain.View, int64, error) {
	queryDB := r.listBase(query)

	var total int64
	if err := queryDB.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var items []attachmentdomain.View
	if err := queryDB.Select(`
			a.id, a.file_id, a.display_name, a.category, a.biz_type,
			f.original_name, f.file_name, f.ext, f.mime_type, f.size, f.url,
			a.uploader_id, a.status, a.remark,
			a.created_at, a.updated_at
		`).
		Order("a.id DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Scan(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (r *Repository) Create(tx *gorm.DB, item *attachmentdomain.Entity) error {
	return r.dbOr(tx).Create(item).Error
}

func (r *Repository) FindByID(tx *gorm.DB, id uint) (attachmentdomain.Entity, error) {
	var item attachmentdomain.Entity
	err := r.dbOr(tx).First(&item, id).Error
	return item, err
}

func (r *Repository) FindViewByID(id uint) (attachmentdomain.View, error) {
	var item attachmentdomain.View
	err := r.db.Table("sys_attachment AS a").
		Select(`
			a.id, a.file_id, a.display_name, a.category, a.biz_type,
			f.original_name, f.file_name, f.ext, f.mime_type, f.size, f.url,
			a.uploader_id, a.status, a.remark,
			a.created_at, a.updated_at
		`).
		Joins("JOIN sys_file AS f ON f.id = a.file_id").
		Where("a.id = ?", id).
		Where("a.deleted_at IS NULL").
		Where("f.deleted_at IS NULL").
		Scan(&item).Error
	if err != nil {
		return attachmentdomain.View{}, err
	}
	if item.ID == 0 {
		return attachmentdomain.View{}, gorm.ErrRecordNotFound
	}

	return item, nil
}

func (r *Repository) UpdateBase(tx *gorm.DB, item *attachmentdomain.Entity, req attachmentdomain.UpdateRequest) error {
	item.DisplayName = req.DisplayName
	item.Category = req.Category
	item.BizType = req.BizType
	item.Status = req.Status
	item.Remark = req.Remark

	return r.dbOr(tx).Model(item).Updates(map[string]any{
		"display_name": item.DisplayName,
		"category":     item.Category,
		"biz_type":     item.BizType,
		"status":       item.Status,
		"remark":       item.Remark,
	}).Error
}

func (r *Repository) UpdateStatus(tx *gorm.DB, item *attachmentdomain.Entity, status model.SystemAttachmentStatus) error {
	item.Status = status
	return r.dbOr(tx).Model(item).Update("status", status).Error
}

func (r *Repository) listBase(query attachmentdomain.ListQuery) *gorm.DB {
	queryDB := r.db.Table("sys_attachment AS a").
		Joins("JOIN sys_file AS f ON f.id = a.file_id").
		Where("a.deleted_at IS NULL").
		Where("f.deleted_at IS NULL")

	keyword := strings.TrimSpace(query.Keyword)
	if keyword != "" {
		like := "%" + keyword + "%"
		queryDB = queryDB.Where(
			"a.display_name LIKE ? OR f.original_name LIKE ? OR f.file_name LIKE ?",
			like, like, like,
		)
	}

	category := strings.TrimSpace(query.Category)
	if category != "" {
		queryDB = queryDB.Where("a.category = ?", category)
	}

	bizType := strings.TrimSpace(query.BizType)
	if bizType != "" {
		queryDB = queryDB.Where("a.biz_type = ?", bizType)
	}

	ext := strings.ToLower(strings.TrimSpace(query.Ext))
	if ext != "" {
		queryDB = queryDB.Where("f.ext = ?", ext)
	}

	statusFilter, err := attachmentdomain.NormalizeStatusFilter(query.Status)
	if err == nil && statusFilter != nil {
		queryDB = queryDB.Where("a.status = ?", *statusFilter)
	}

	return queryDB
}

func (r *Repository) dbOr(tx *gorm.DB) *gorm.DB {
	if tx != nil {
		return tx
	}
	return r.db
}
