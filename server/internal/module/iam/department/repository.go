package department

import (
	"errors"
	"fmt"
	"strings"

	"ez-admin-gin/server/internal/apperror"
	"ez-admin-gin/server/internal/model"
	"ez-admin-gin/server/internal/platform/datascope"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) List(actor datascope.Actor, query ListQuery) ([]model.Department, error) {
	queryDB := applyDataScope(r.db.Model(&model.Department{}), actor)

	keyword := strings.TrimSpace(query.Keyword)
	if keyword != "" {
		like := "%" + keyword + "%"
		queryDB = queryDB.Where("name LIKE ? OR code LIKE ?", like, like)
	}

	if query.Status != 0 {
		status := model.DepartmentStatus(query.Status)
		if !ValidStatus(status) {
			return nil, apperror.BadRequest("部门状态不正确")
		}
		queryDB = queryDB.Where("status = ?", status)
	}

	var items []model.Department
	if err := queryDB.Order("sort ASC, id ASC").Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

func (r *Repository) FindByIDInScope(db *gorm.DB, actor datascope.Actor, departmentID uint) (model.Department, error) {
	var department model.Department
	err := applyDataScope(db, actor).First(&department, departmentID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.Department{}, apperror.NotFound("部门不存在")
		}
		return model.Department{}, err
	}

	return department, nil
}

func (r *Repository) FindByID(db *gorm.DB, departmentID uint) (model.Department, error) {
	var department model.Department
	err := db.First(&department, departmentID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.Department{}, apperror.NotFound("部门不存在")
		}
		return model.Department{}, err
	}

	return department, nil
}

func (r *Repository) FindParent(db *gorm.DB, parentID uint) (model.Department, error) {
	if parentID == 0 {
		return model.Department{}, nil
	}

	return r.FindByID(db, parentID)
}

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

func (r *Repository) LeaderUsable(db *gorm.DB, leaderUserID uint) error {
	if leaderUserID == 0 {
		return nil
	}

	var count int64
	err := db.Model(&model.User{}).Where("id = ?", leaderUserID).Where("status = ?", model.UserStatusEnabled).Count(&count).Error
	if err != nil {
		return err
	}
	if count != 1 {
		return apperror.BadRequest("负责人不存在或已禁用")
	}

	return nil
}

func (r *Repository) Create(db *gorm.DB, department *model.Department) error {
	return db.Create(department).Error
}

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

func (r *Repository) UpdateStatus(db *gorm.DB, department *model.Department, status model.DepartmentStatus) error {
	if err := db.Model(department).Update("status", status).Error; err != nil {
		return err
	}
	department.Status = status
	return nil
}

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

func (r *Repository) UpdateAncestors(db *gorm.DB, departmentID uint, ancestors string) error {
	return db.Model(&model.Department{}).Where("id = ?", departmentID).Update("ancestors", ancestors).Error
}

func BuildAncestors(parent model.Department) string {
	if parent.ID == 0 {
		return "0"
	}
	return fmt.Sprintf("%s,%d", parent.Ancestors, parent.ID)
}

func FullPath(item model.Department) string {
	if item.Ancestors == "" {
		return fmt.Sprintf("%d", item.ID)
	}
	return fmt.Sprintf("%s,%d", item.Ancestors, item.ID)
}

func IsDescendantPath(path string, target string) bool {
	return path == target || strings.HasPrefix(path, target+",")
}
