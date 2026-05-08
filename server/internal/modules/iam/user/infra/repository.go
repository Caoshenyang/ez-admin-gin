// Package infra 实现用户的数据访问层，包含数据权限作用域过滤。
package infra

import (
	"errors"
	"strings"

	userdomain "ez-admin-gin/server/internal/modules/iam/user/domain"
	errorsx "ez-admin-gin/server/internal/pkg/errorsx"
	"ez-admin-gin/server/internal/platform/datascope"
	"ez-admin-gin/server/internal/platform/model"

	"gorm.io/gorm"
)

// Repository 实现用户的数据访问层。
type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// List 根据查询条件分页返回数据权限范围内的用户列表。
func (r *Repository) List(actor datascope.Actor, query userdomain.ListQuery, page int, pageSize int) ([]model.User, int64, error) {
	queryDB := applyDataScope(r.db.Model(&model.User{}), actor)

	keyword := strings.TrimSpace(query.Keyword)
	if keyword != "" {
		like := "%" + keyword + "%"
		queryDB = queryDB.Where("username LIKE ? OR nickname LIKE ?", like, like)
	}

	if query.Status != 0 {
		status := model.UserStatus(query.Status)
		if !userdomain.ValidStatus(status) {
			return nil, 0, errorsx.BadRequest("用户状态不正确")
		}
		queryDB = queryDB.Where("status = ?", status)
	}

	var total int64
	if err := queryDB.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var users []model.User
	if err := queryDB.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// RoleIDsByUserIDs 批量查询指定用户关联的角色 ID。
func (r *Repository) RoleIDsByUserIDs(userIDs []uint) (map[uint][]uint, error) {
	result := make(map[uint][]uint, len(userIDs))
	if len(userIDs) == 0 {
		return result, nil
	}

	var rows []model.UserRole
	if err := r.db.Where("user_id IN ?", userIDs).Order("role_id ASC").Find(&rows).Error; err != nil {
		return nil, err
	}

	for _, row := range rows {
		result[row.UserID] = append(result[row.UserID], row.RoleID)
	}

	return result, nil
}

// PostIDsByUserIDs 批量查询指定用户关联的岗位 ID。
func (r *Repository) PostIDsByUserIDs(userIDs []uint) (map[uint][]uint, error) {
	result := make(map[uint][]uint, len(userIDs))
	if len(userIDs) == 0 {
		return result, nil
	}

	var rows []model.UserPost
	if err := r.db.Where("user_id IN ?", userIDs).Order("post_id ASC").Find(&rows).Error; err != nil {
		return nil, err
	}

	for _, row := range rows {
		result[row.UserID] = append(result[row.UserID], row.PostID)
	}

	return result, nil
}

// FindByIDInScope 在数据权限范围内按 ID 查找用户。
func (r *Repository) FindByIDInScope(db *gorm.DB, actor datascope.Actor, userID uint) (model.User, error) {
	var user model.User
	err := applyDataScope(db, actor).First(&user, userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.User{}, errorsx.NotFound("用户不存在")
		}
		return model.User{}, err
	}

	return user, nil
}

// UsernameExists 检查用户名是否已存在。
func (r *Repository) UsernameExists(db *gorm.DB, username string) (bool, error) {
	var user model.User
	err := db.Unscoped().Where("username = ?", username).First(&user).Error
	if err == nil {
		return true, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}

	return false, err
}

// DepartmentUsable 校验指定部门是否存在且启用。
func (r *Repository) DepartmentUsable(db *gorm.DB, departmentID uint) error {
	if departmentID == 0 {
		return nil
	}

	var count int64
	err := db.Model(&model.Department{}).
		Where("id = ?", departmentID).
		Where("status = ?", model.DepartmentStatusEnabled).
		Count(&count).Error
	if err != nil {
		return err
	}
	if count != 1 {
		return errorsx.BadRequest("部门不存在或已禁用")
	}

	return nil
}

// RolesUsable 校验给定的角色 ID 是否全部存在且启用。
func (r *Repository) RolesUsable(db *gorm.DB, roleIDs []uint) error {
	if len(roleIDs) == 0 {
		return nil
	}

	var count int64
	err := db.Model(&model.Role{}).
		Where("id IN ?", roleIDs).
		Where("status = ?", model.RoleStatusEnabled).
		Count(&count).Error
	if err != nil {
		return err
	}
	if count != int64(len(roleIDs)) {
		return errorsx.BadRequest("角色不存在或已禁用")
	}

	return nil
}

// PostsUsable 校验给定的岗位 ID 是否全部存在且启用。
func (r *Repository) PostsUsable(db *gorm.DB, postIDs []uint) error {
	if len(postIDs) == 0 {
		return nil
	}

	var count int64
	err := db.Model(&model.Post{}).
		Where("id IN ?", postIDs).
		Where("status = ?", model.PostStatusEnabled).
		Count(&count).Error
	if err != nil {
		return err
	}
	if count != int64(len(postIDs)) {
		return errorsx.BadRequest("岗位不存在或已禁用")
	}

	return nil
}

// Create 创建用户记录。
func (r *Repository) Create(db *gorm.DB, user *model.User) error {
	return db.Create(user).Error
}

// UpdateBase 更新用户的基本信息字段。
func (r *Repository) UpdateBase(db *gorm.DB, user *model.User, nickname string, departmentID uint, status model.UserStatus) error {
	if err := db.Model(user).Updates(map[string]any{
		"nickname":      nickname,
		"department_id": departmentID,
		"status":        status,
	}).Error; err != nil {
		return err
	}

	user.Nickname = nickname
	user.DepartmentID = departmentID
	user.Status = status
	return nil
}

// UpdateStatus 更新用户的启用/禁用状态。
func (r *Repository) UpdateStatus(db *gorm.DB, user *model.User, status model.UserStatus) error {
	if err := db.Model(user).Update("status", status).Error; err != nil {
		return err
	}

	user.Status = status
	return nil
}

// ReplaceRoles 替换指定用户的全部角色关联。
func (r *Repository) ReplaceRoles(db *gorm.DB, userID uint, roleIDs []uint) error {
	if err := db.Where("user_id = ?", userID).Delete(&model.UserRole{}).Error; err != nil {
		return err
	}

	if len(roleIDs) == 0 {
		return nil
	}

	rows := make([]model.UserRole, 0, len(roleIDs))
	for _, roleID := range roleIDs {
		rows = append(rows, model.UserRole{
			UserID: userID,
			RoleID: roleID,
		})
	}

	return db.Create(&rows).Error
}

// ReplacePosts 替换指定用户的全部岗位关联。
func (r *Repository) ReplacePosts(db *gorm.DB, userID uint, postIDs []uint) error {
	if err := db.Where("user_id = ?", userID).Delete(&model.UserPost{}).Error; err != nil {
		return err
	}

	if len(postIDs) == 0 {
		return nil
	}

	rows := make([]model.UserPost, 0, len(postIDs))
	for _, postID := range postIDs {
		rows = append(rows, model.UserPost{
			UserID: userID,
			PostID: postID,
		})
	}

	return db.Create(&rows).Error
}
