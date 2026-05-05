package attachment

import (
	"strings"

	"ez-admin-gin/server/internal/model"

	"gorm.io/gorm"
)

// Repository 负责附件中心的查询和持久化。
type Repository struct {
	db *gorm.DB
}

// NewRepository 创建附件中心仓储。
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// List 返回附件分页结果和总数。
func (r *Repository) List(query ListQuery, page int, pageSize int) ([]View, int64, error) {
	queryDB := r.listBase(query)

	var total int64
	if err := queryDB.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var items []View
	if err := queryDB.
		Select(`
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

// Create 创建附件记录。
func (r *Repository) Create(tx *gorm.DB, item *model.SystemAttachment) error {
	return r.dbOr(tx).Create(item).Error
}

// FindByID 查询附件实体。
func (r *Repository) FindByID(tx *gorm.DB, id uint) (Entity, error) {
	var item Entity
	err := r.dbOr(tx).First(&item, id).Error
	return item, err
}

// FindViewByID 查询附件联合视图。
func (r *Repository) FindViewByID(id uint) (View, error) {
	var item View
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
		return View{}, err
	}
	if item.ID == 0 {
		return View{}, gorm.ErrRecordNotFound
	}

	return item, nil
}

// UpdateBase 更新附件元数据。
func (r *Repository) UpdateBase(tx *gorm.DB, item *Entity, req UpdateRequest) error {
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

// UpdateStatus 单独更新附件状态。
func (r *Repository) UpdateStatus(tx *gorm.DB, item *Entity, status model.SystemAttachmentStatus) error {
	item.Status = status
	return r.dbOr(tx).Model(item).Update("status", status).Error
}

func (r *Repository) listBase(query ListQuery) *gorm.DB {
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

	statusFilter, err := NormalizeStatusFilter(query.Status)
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
