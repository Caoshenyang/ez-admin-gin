package post

import (
	"errors"
	"strings"

	"ez-admin-gin/server/internal/apperror"
	"ez-admin-gin/server/internal/model"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) List(query ListQuery) ([]model.Post, error) {
	queryDB := applyDataScope(r.db.Model(&model.Post{}))

	keyword := strings.TrimSpace(query.Keyword)
	if keyword != "" {
		like := "%" + keyword + "%"
		queryDB = queryDB.Where("code LIKE ? OR name LIKE ?", like, like)
	}

	if query.Status != 0 {
		status := model.PostStatus(query.Status)
		if !ValidStatus(status) {
			return nil, apperror.BadRequest("岗位状态不正确")
		}
		queryDB = queryDB.Where("status = ?", status)
	}

	var items []model.Post
	if err := queryDB.Order("sort ASC, id ASC").Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

func (r *Repository) FindByID(db *gorm.DB, postID uint) (model.Post, error) {
	var item model.Post
	err := db.First(&item, postID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.Post{}, apperror.NotFound("岗位不存在")
		}
		return model.Post{}, err
	}

	return item, nil
}

func (r *Repository) CodeExists(db *gorm.DB, code string, excludeID uint) (bool, error) {
	var item model.Post
	query := db.Unscoped().Where("code = ?", code)
	if excludeID != 0 {
		query = query.Where("id <> ?", excludeID)
	}

	err := query.First(&item).Error
	if err == nil {
		return true, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}

	return false, err
}

func (r *Repository) Create(db *gorm.DB, item *model.Post) error {
	return db.Create(item).Error
}

func (r *Repository) Update(db *gorm.DB, item *model.Post, code string, name string, sort int, status model.PostStatus, remark string) error {
	if err := db.Model(item).Updates(map[string]any{
		"code":   code,
		"name":   name,
		"sort":   sort,
		"status": status,
		"remark": remark,
	}).Error; err != nil {
		return err
	}

	item.Code = code
	item.Name = name
	item.Sort = sort
	item.Status = status
	item.Remark = remark
	return nil
}

func (r *Repository) UpdateStatus(db *gorm.DB, item *model.Post, status model.PostStatus) error {
	if err := db.Model(item).Update("status", status).Error; err != nil {
		return err
	}
	item.Status = status
	return nil
}
