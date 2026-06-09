// Package infra 实现角色的数据访问层，含 Casbin 策略写入。
package infra

import (
	"errors"
	"strings"

	roledomain "ez-admin-gin/server/internal/modules/iam/role/domain"
	errorsx "ez-admin-gin/server/internal/pkg/errorsx"
	"ez-admin-gin/server/internal/platform/model"

	"gorm.io/gorm"
)

// Repository 实现角色的数据访问层。
type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// List 根据查询条件分页返回角色列表。
func (r *Repository) List(query roledomain.ListQuery, page int, pageSize int) ([]model.Role, int64, error) {
	queryDB := r.db.Model(&model.Role{})

	keyword := strings.TrimSpace(query.Keyword)
	if keyword != "" {
		like := "%" + keyword + "%"
		queryDB = queryDB.Where("code LIKE ? OR name LIKE ?", like, like)
	}

	if query.Status != 0 {
		status := model.RoleStatus(query.Status)
		if !roledomain.ValidRoleStatus(status) {
			return nil, 0, errorsx.BadRequest("角色状态不正确")
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

// FindByID 按 ID 查找角色。
func (r *Repository) FindByID(db *gorm.DB, roleID uint) (model.Role, error) {
	var role model.Role
	err := db.First(&role, roleID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.Role{}, errorsx.NotFound("角色不存在")
		}
		return model.Role{}, err
	}

	return role, nil
}

// CodeExists 检查角色编码是否已存在。
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

// DepartmentsUsable 校验给定的部门 ID 是否全部存在且启用。
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
		return errorsx.BadRequest("部门不存在或已禁用")
	}

	return nil
}

// MenusUsable 校验给定的菜单 ID 是否全部存在且启用。
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
		return errorsx.BadRequest("菜单不存在或已禁用")
	}

	return nil
}

// APIsUsable 校验给定的接口权限 ID 是否全部存在且启用。
func (r *Repository) APIsUsable(db *gorm.DB, apiIDs []uint) error {
	if len(apiIDs) == 0 {
		return nil
	}

	var count int64
	err := db.Model(&model.API{}).Where("id IN ?", apiIDs).Where("status = ?", model.APIStatusEnabled).Count(&count).Error
	if err != nil {
		return err
	}
	if count != int64(len(apiIDs)) {
		return errorsx.BadRequest("接口权限不存在或已禁用")
	}

	return nil
}

// Create 创建角色记录。
func (r *Repository) Create(db *gorm.DB, role *model.Role) error {
	return db.Create(role).Error
}

// UpdateBase 更新角色的基本信息字段。
func (r *Repository) UpdateBase(db *gorm.DB, role *model.Role, req roledomain.UpdateRequest) error {
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

// UpdateStatus 更新角色的启用/禁用状态。
func (r *Repository) UpdateStatus(db *gorm.DB, role *model.Role, status model.RoleStatus) error {
	if err := db.Model(role).Update("status", status).Error; err != nil {
		return err
	}
	role.Status = status
	return nil
}

// CountUsers 统计绑定到该角色的用户数量。
func (r *Repository) CountUsers(db *gorm.DB, roleID uint) (int64, error) {
	var count int64
	err := db.Model(&model.UserRole{}).Where("role_id = ?", roleID).Count(&count).Error
	return count, err
}

// RoleAPIIDs 批量查询指定角色关联的接口权限 ID。
func (r *Repository) RoleAPIIDs(roleIDs []uint) (map[uint][]uint, error) {
	result := make(map[uint][]uint, len(roleIDs))
	if len(roleIDs) == 0 {
		return result, nil
	}

	var rows []model.RoleAPI
	if err := r.db.Where("role_id IN ?", roleIDs).Order("api_id ASC").Find(&rows).Error; err != nil {
		return nil, err
	}

	for _, row := range rows {
		result[row.RoleID] = append(result[row.RoleID], row.APIID)
	}

	return result, nil
}

// RolePermissions 批量查询指定角色关联的接口权限元数据快照。
func (r *Repository) RolePermissions(roleIDs []uint) (map[uint][]roledomain.PermissionItem, error) {
	result := make(map[uint][]roledomain.PermissionItem, len(roleIDs))
	if len(roleIDs) == 0 {
		return result, nil
	}

	type row struct {
		RoleID uint
		ID     uint
		Code   string
		Name   string
		Module string
		Method string
		Path   string
		Status model.APIStatus
	}

	var rows []row
	err := r.db.
		Table("sys_role_api AS ra").
		Select("ra.role_id, a.id, a.code, a.name, a.module, a.method, a.path, a.status").
		Joins("JOIN sys_api AS a ON a.id = ra.api_id").
		Where("ra.role_id IN ?", roleIDs).
		Where("a.deleted_at IS NULL").
		Order("a.module ASC, a.sort ASC, a.id ASC").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	for _, row := range rows {
		result[row.RoleID] = append(result[row.RoleID], roledomain.PermissionItem{
			ID:     row.ID,
			Code:   row.Code,
			Name:   row.Name,
			Module: row.Module,
			Method: row.Method,
			Path:   row.Path,
			Status: row.Status,
		})
	}

	return result, nil
}

// RoleMenuIDs 批量查询指定角色关联的菜单 ID。
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

// RoleCustomDepartmentIDs 批量查询指定角色的自定义数据范围部门 ID。
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

// ReplaceAPIs 替换指定角色的全部接口权限关联。
func (r *Repository) ReplaceAPIs(db *gorm.DB, roleID uint, apiIDs []uint) error {
	if err := db.Where("role_id = ?", roleID).Delete(&model.RoleAPI{}).Error; err != nil {
		return err
	}
	if len(apiIDs) == 0 {
		return nil
	}

	rows := make([]model.RoleAPI, 0, len(apiIDs))
	for _, apiID := range apiIDs {
		rows = append(rows, model.RoleAPI{RoleID: roleID, APIID: apiID})
	}

	return db.Create(&rows).Error
}

// ReplacePoliciesByAPIs 用角色接口关联重建指定角色的 Casbin 执行策略。
func (r *Repository) ReplacePoliciesByAPIs(db *gorm.DB, roleCode string, apiIDs []uint) error {
	if err := db.Where("ptype = ? AND v0 = ?", "p", roleCode).Delete(&model.CasbinRule{}).Error; err != nil {
		return err
	}
	if len(apiIDs) == 0 {
		return nil
	}

	var apis []model.API
	if err := db.Where("id IN ?", apiIDs).Where("status = ?", model.APIStatusEnabled).Order("id ASC").Find(&apis).Error; err != nil {
		return err
	}
	if len(apis) == 0 {
		return nil
	}

	rows := make([]model.CasbinRule, 0, len(apis))
	for _, item := range apis {
		rows = append(rows, model.CasbinRule{Ptype: "p", V0: roleCode, V1: item.Path, V2: item.Method})
	}

	return db.Create(&rows).Error
}

// ReplaceMenus 替换指定角色的全部菜单关联。
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

// ReplaceCustomDepartments 替换指定角色的自定义数据范围部门。
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

// Delete 软删除角色记录。
func (r *Repository) Delete(db *gorm.DB, role *model.Role) error {
	return db.Delete(role).Error
}
