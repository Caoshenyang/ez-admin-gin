// Package infra 实现岗位的数据访问层。
package infra

import (
	"errors"
	"strings"

	postdomain "ez-admin-gin/server/internal/modules/iam/post/domain"
	errorsx "ez-admin-gin/server/internal/pkg/errorsx"
	"ez-admin-gin/server/internal/platform/model"

	"gorm.io/gorm"
)

// Repository 实现岗位的数据访问层。
type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// List 根据查询条件返回岗位列表。
func (r *Repository) List(query postdomain.ListQuery) ([]model.Post, error) {
	queryDB := applyDataScope(r.db.Model(&model.Post{}))

	keyword := strings.TrimSpace(query.Keyword)
	if keyword != "" {
		like := "%" + keyword + "%"
		queryDB = queryDB.Where("code LIKE ? OR name LIKE ?", like, like)
	}

	if query.Status != 0 {
		status := model.PostStatus(query.Status)
		if !postdomain.ValidStatus(status) {
			return nil, errorsx.BadRequest("岗位状态不正确")
		}
		queryDB = queryDB.Where("status = ?", status)
	}

	var items []model.Post
	if err := queryDB.Order("sort ASC, id ASC").Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

// FindByID 按 ID 查找岗位。
func (r *Repository) FindByID(db *gorm.DB, postID uint) (model.Post, error) {
	var item model.Post
	err := db.First(&item, postID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.Post{}, errorsx.NotFound("岗位不存在")
		}
		return model.Post{}, err
	}

	return item, nil
}

// CodeExists 检查岗位编码是否已存在，可排除指定 ID。
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

// Create 创建岗位记录。
func (r *Repository) Create(db *gorm.DB, item *model.Post) error {
	return db.Create(item).Error
}

// Update 更新岗位的所有可编辑字段。
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

// UpdateStatus 更新岗位的启用/禁用状态。
func (r *Repository) UpdateStatus(db *gorm.DB, item *model.Post, status model.PostStatus) error {
	if err := db.Model(item).Update("status", status).Error; err != nil {
		return err
	}
	item.Status = status
	return nil
}
