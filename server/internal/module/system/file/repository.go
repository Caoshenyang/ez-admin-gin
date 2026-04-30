package file

import (
	"strings"

	"ez-admin-gin/server/internal/apperror"
	"ez-admin-gin/server/internal/model"

	"gorm.io/gorm"
)

// Repository 负责文件模块的查询和持久化。
type Repository struct {
	db *gorm.DB
}

// NewRepository 创建文件仓储。
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// List 返回文件分页结果和总数。
func (r *Repository) List(query ListQuery, page int, pageSize int) ([]model.SystemFile, int64, error) {
	queryDB := r.db.Model(&model.SystemFile{})

	keyword := strings.TrimSpace(query.Keyword)
	if keyword != "" {
		like := "%" + keyword + "%"
		queryDB = queryDB.Where("original_name LIKE ? OR file_name LIKE ?", like, like)
	}

	ext := NormalizeExt(query.Ext)
	if ext != "" {
		queryDB = queryDB.Where("ext = ?", ext)
	}

	if query.Status != 0 {
		status := model.SystemFileStatus(query.Status)
		if !ValidStatus(status) {
			return nil, 0, apperror.BadRequest("文件状态不正确")
		}
		queryDB = queryDB.Where("status = ?", status)
	}

	var total int64
	if err := queryDB.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var items []model.SystemFile
	if err := queryDB.
		Order("id DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

// Create 创建文件记录。
func (r *Repository) Create(db *gorm.DB, item *model.SystemFile) error {
	return db.Create(item).Error
}
