// Package infra 实现部门的数据访问层，包含数据权限作用域过滤。
package infra

import (
	"errors"
	"strings"

	departmentdomain "ez-admin-gin/server/internal/modules/iam/department/domain"
	errorsx "ez-admin-gin/server/internal/pkg/errorsx"
	"ez-admin-gin/server/internal/platform/datascope"
	"ez-admin-gin/server/internal/platform/model"

	"gorm.io/gorm"
)

// Repository 实现部门的数据访问层。
type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// List 根据查询条件返回数据权限范围内的部门列表。
func (r *Repository) List(actor datascope.Actor, query departmentdomain.ListQuery) ([]model.Department, error) {
	queryDB := applyDataScope(r.db.Model(&model.Department{}), actor)

	keyword := strings.TrimSpace(query.Keyword)
	if keyword != "" {
		like := "%" + keyword + "%"
		queryDB = queryDB.Where("name LIKE ? OR code LIKE ?", like, like)
	}

	if query.Status != 0 {
		status := model.DepartmentStatus(query.Status)
		if !departmentdomain.ValidStatus(status) {
			return nil, errorsx.BadRequest("部门状态不正确")
		}
		queryDB = queryDB.Where("status = ?", status)
	}

	var items []model.Department
	if err := queryDB.Order("sort ASC, id ASC").Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

// FindByIDInScope 在数据权限范围内按 ID 查找部门。
func (r *Repository) FindByIDInScope(db *gorm.DB, actor datascope.Actor, departmentID uint) (model.Department, error) {
	var department model.Department
	err := applyDataScope(db, actor).First(&department, departmentID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.Department{}, errorsx.NotFound("部门不存在")
		}
		return model.Department{}, err
	}

	return department, nil
}

// FindByID 按 ID 查找部门。
func (r *Repository) FindByID(db *gorm.DB, departmentID uint) (model.Department, error) {
	var department model.Department
	err := db.First(&department, departmentID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.Department{}, errorsx.NotFound("部门不存在")
		}
		return model.Department{}, err
	}

	return department, nil
}

// FindParent 查找指定 ID 的父部门，parentID 为 0 时返回空值。
func (r *Repository) FindParent(db *gorm.DB, parentID uint) (model.Department, error) {
	if parentID == 0 {
		return model.Department{}, nil
	}

	return r.FindByID(db, parentID)
}

// CodeExists 检查部门编码是否已存在，可排除指定 ID。
func (r *Repository) CodeExists(db *gorm.DB, code string, excludeID uint) (bool, error) {
	var department model.Department
	query := db.Unscoped().Where("code = ?", code)
	if excludeID != 0 {
		query = query.Where("id <> ?", excludeID)
	}

	err := query.First(&department).Error
	if err == nil {
		return true, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}

	return false, err
}

// LeaderUsable 校验部门负责人是否存在且处于启用状态。
func (r *Repository) LeaderUsable(db *gorm.DB, leaderUserID uint) error {
	if leaderUserID == 0 {
		return nil
	}

	var count int64
	err := db.Model(&model.User{}).
		Where("id = ?", leaderUserID).
		Where("status = ?", model.UserStatusEnabled).
		Count(&count).Error
	if err != nil {
		return err
	}
	if count != 1 {
		return errorsx.BadRequest("负责人不存在或已禁用")
	}

	return nil
}

// Create 创建部门记录。
func (r *Repository) Create(db *gorm.DB, department *model.Department) error {
	return db.Create(department).Error
}

// Update 更新部门的所有可编辑字段。
func (r *Repository) Update(db *gorm.DB, department *model.Department, parentID uint, ancestors string, name string, code string, leaderUserID uint, sort int, status model.DepartmentStatus, remark string) error {
	if err := db.Model(department).Updates(map[string]any{
		"parent_id":      parentID,
		"ancestors":      ancestors,
		"name":           name,
		"code":           code,
		"leader_user_id": leaderUserID,
		"sort":           sort,
		"status":         status,
		"remark":         remark,
	}).Error; err != nil {
		return err
	}

	department.ParentID = parentID
	department.Ancestors = ancestors
	department.Name = name
	department.Code = code
	department.LeaderUserID = leaderUserID
	department.Sort = sort
	department.Status = status
	department.Remark = remark
	return nil
}

// UpdateStatus 更新部门的启用/禁用状态。
func (r *Repository) UpdateStatus(db *gorm.DB, department *model.Department, status model.DepartmentStatus) error {
	if err := db.Model(department).Update("status", status).Error; err != nil {
		return err
	}
	department.Status = status
	return nil
}

// Subtree 查询指定部门的整棵子树。
func (r *Repository) Subtree(db *gorm.DB, departmentID uint, fullPath string) ([]model.Department, error) {
	var items []model.Department
	if err := db.
		Where("ancestors = ? OR ancestors LIKE ?", fullPath, fullPath+",%").
		Order("id ASC").
		Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

// UpdateAncestors 更新指定部门的祖先路径。
func (r *Repository) UpdateAncestors(db *gorm.DB, departmentID uint, ancestors string) error {
	return db.Model(&model.Department{}).Where("id = ?", departmentID).Update("ancestors", ancestors).Error
}
