package infra

import (
	"strings"

	filedomain "ez-admin-gin/server/internal/modules/system/file/domain"
	errorsx "ez-admin-gin/server/internal/pkg/errorsx"
	"ez-admin-gin/server/internal/platform/model"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) List(query filedomain.ListQuery, page int, pageSize int) ([]model.SystemFile, int64, error) {
	queryDB := r.db.Model(&model.SystemFile{})

	keyword := strings.TrimSpace(query.Keyword)
	if keyword != "" {
		like := "%" + keyword + "%"
		queryDB = queryDB.Where("original_name LIKE ? OR file_name LIKE ?", like, like)
	}

	ext := filedomain.NormalizeExt(query.Ext)
	if ext != "" {
		queryDB = queryDB.Where("ext = ?", ext)
	}

	if query.Status != 0 {
		status := model.SystemFileStatus(query.Status)
		if !filedomain.ValidStatus(status) {
			return nil, 0, errorsx.BadRequest("文件状态不正确")
		}
		queryDB = queryDB.Where("status = ?", status)
	}

	var total int64
	if err := queryDB.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var items []model.SystemFile
	if err := queryDB.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (r *Repository) Create(db *gorm.DB, item *model.SystemFile) error {
	return db.Create(item).Error
}

func (r *Repository) DeleteByID(db *gorm.DB, id uint) error {
	return db.Where("id = ?", id).Delete(&model.SystemFile{}).Error
}
