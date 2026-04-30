package role

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

func (r *Repository) List(query ListQuery, page int, pageSize int) ([]model.Role, int64, error) {
	queryDB := r.db.Model(&model.Role{})

	keyword := strings.TrimSpace(query.Keyword)
	if keyword != "" {
		like := "%" + keyword + "%"
		queryDB = queryDB.Where("code LIKE ? OR name LIKE ?", like, like)
	}

	if query.Status != 0 {
		status := model.RoleStatus(query.Status)
		if !ValidRoleStatus(status) {
			return nil, 0, apperror.BadRequest("角色状态不正确")
		}
		queryDB = queryDB.Where("status = ?", status)
	}

	var total int64
	if err := queryDB.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var roles []model.Role
	if err := queryDB.Order("sort ASC, id ASC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&roles).Error; err != nil {
		return nil, 0, err
	}

	return roles, total, nil
}

func (r *Repository) FindByID(db *gorm.DB, roleID uint) (model.Role, error) {
	var role model.Role
	err := db.First(&role, roleID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.Role{}, apperror.NotFound("角色不存在")
		}
		return model.Role{}, err
	}

	return role, nil
}

func (r *Repository) CodeExists(db *gorm.DB, code string) (bool, error) {
	var role model.Role
	err := db.Unscoped().Where("code = ?", code).First(&role).Error
	if err == nil {
		return true, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}

	return false, err
}

func (r *Repository) DepartmentsUsable(db *gorm.DB, departmentIDs []uint) error {
	if len(departmentIDs) == 0 {
		return nil
	}

	var count int64
	err := db.Model(&model.Department{}).Where("id IN ?", departmentIDs).Where("status = ?", model.DepartmentStatusEnabled).Count(&count).Error
	if err != nil {
		return err
	}
	if count != int64(len(departmentIDs)) {
		return apperror.BadRequest("部门不存在或已禁用")
	}

	return nil
}

func (r *Repository) MenusUsable(db *gorm.DB, menuIDs []uint) error {
	if len(menuIDs) == 0 {
		return nil
	}

	var count int64
	err := db.Model(&model.Menu{}).Where("id IN ?", menuIDs).Where("status = ?", model.MenuStatusEnabled).Count(&count).Error
	if err != nil {
		return err
	}
	if count != int64(len(menuIDs)) {
		return apperror.BadRequest("菜单不存在或已禁用")
	}

	return nil
}

func (r *Repository) Create(db *gorm.DB, role *model.Role) error {
	return db.Create(role).Error
}

func (r *Repository) UpdateBase(db *gorm.DB, role *model.Role, req UpdateRequest) error {
	if err := db.Model(role).Updates(map[string]any{
		"name":       req.Name,
		"sort":       req.Sort,
		"data_scope": req.DataScope,
		"status":     req.Status,
		"remark":     req.Remark,
	}).Error; err != nil {
		return err
	}

	role.Name = req.Name
	role.Sort = req.Sort
	role.DataScope = req.DataScope
	role.Status = req.Status
	role.Remark = req.Remark
	return nil
}

func (r *Repository) UpdateStatus(db *gorm.DB, role *model.Role, status model.RoleStatus) error {
	if err := db.Model(role).Update("status", status).Error; err != nil {
		return err
	}
	role.Status = status
	return nil
}

func (r *Repository) RolePermissions(roleCodes []string) (map[string][]PermissionItem, error) {
	result := make(map[string][]PermissionItem, len(roleCodes))
	if len(roleCodes) == 0 {
		return result, nil
	}

	var rows []model.CasbinRule
	if err := r.db.Where("ptype = ?", "p").Where("v0 IN ?", roleCodes).Order("v1 ASC, v2 ASC").Find(&rows).Error; err != nil {
		return nil, err
	}

	for _, row := range rows {
		result[row.V0] = append(result[row.V0], PermissionItem{Path: row.V1, Method: row.V2})
	}

	return result, nil
}

func (r *Repository) RoleMenuIDs(roleIDs []uint) (map[uint][]uint, error) {
	result := make(map[uint][]uint, len(roleIDs))
	if len(roleIDs) == 0 {
		return result, nil
	}

	var rows []model.RoleMenu
	if err := r.db.Where("role_id IN ?", roleIDs).Order("menu_id ASC").Find(&rows).Error; err != nil {
		return nil, err
	}

	for _, row := range rows {
		result[row.RoleID] = append(result[row.RoleID], row.MenuID)
	}

	return result, nil
}

func (r *Repository) RoleCustomDepartmentIDs(roleIDs []uint) (map[uint][]uint, error) {
	result := make(map[uint][]uint, len(roleIDs))
	if len(roleIDs) == 0 {
		return result, nil
	}

	var rows []model.RoleDataScope
	if err := r.db.Where("role_id IN ?", roleIDs).Order("department_id ASC").Find(&rows).Error; err != nil {
		return nil, err
	}

	for _, row := range rows {
		result[row.RoleID] = append(result[row.RoleID], row.DepartmentID)
	}

	return result, nil
}

func (r *Repository) ReplacePermissions(db *gorm.DB, roleCode string, permissions []PermissionItem) error {
	if err := db.Where("ptype = ? AND v0 = ?", "p", roleCode).Delete(&model.CasbinRule{}).Error; err != nil {
		return err
	}
	if len(permissions) == 0 {
		return nil
	}

	rows := make([]model.CasbinRule, 0, len(permissions))
	for _, item := range permissions {
		rows = append(rows, model.CasbinRule{Ptype: "p", V0: roleCode, V1: item.Path, V2: item.Method})
	}

	return db.Create(&rows).Error
}

func (r *Repository) ReplaceMenus(db *gorm.DB, roleID uint, menuIDs []uint) error {
	if err := db.Where("role_id = ?", roleID).Delete(&model.RoleMenu{}).Error; err != nil {
		return err
	}
	if len(menuIDs) == 0 {
		return nil
	}

	rows := make([]model.RoleMenu, 0, len(menuIDs))
	for _, menuID := range menuIDs {
		rows = append(rows, model.RoleMenu{RoleID: roleID, MenuID: menuID})
	}

	return db.Create(&rows).Error
}

func (r *Repository) ReplaceCustomDepartments(db *gorm.DB, roleID uint, departmentIDs []uint) error {
	if err := db.Where("role_id = ?", roleID).Delete(&model.RoleDataScope{}).Error; err != nil {
		return err
	}
	if len(departmentIDs) == 0 {
		return nil
	}

	rows := make([]model.RoleDataScope, 0, len(departmentIDs))
	for _, departmentID := range departmentIDs {
		rows = append(rows, model.RoleDataScope{RoleID: roleID, DepartmentID: departmentID})
	}

	return db.Create(&rows).Error
}
