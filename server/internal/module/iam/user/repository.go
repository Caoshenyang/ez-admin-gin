package user

import (
	"errors"
	"strings"

	"ez-admin-gin/server/internal/apperror"
	"ez-admin-gin/server/internal/model"
	"ez-admin-gin/server/internal/platform/datascope"

	"gorm.io/gorm"
)

// Repository 负责用户模块的持久化和查询拼装。
type Repository struct {
	db *gorm.DB
}

// NewRepository 创建用户仓储。
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// List 返回当前数据范围内的用户分页结果和总数。
func (r *Repository) List(actor datascope.Actor, query ListQuery, page int, pageSize int) ([]model.User, int64, error) {
	queryDB := applyDataScope(r.db.Model(&model.User{}), actor)

	keyword := strings.TrimSpace(query.Keyword)
	if keyword != "" {
		like := "%" + keyword + "%"
		queryDB = queryDB.Where("username LIKE ? OR nickname LIKE ?", like, like)
	}

	if query.Status != 0 {
		status := model.UserStatus(query.Status)
		if !ValidStatus(status) {
			return nil, 0, apperror.BadRequest("用户状态不正确")
		}
		queryDB = queryDB.Where("status = ?", status)
	}

	var total int64
	if err := queryDB.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var users []model.User
	if err := queryDB.
		Order("id DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// RoleIDsByUserIDs 批量查询用户对应的角色 ID。
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

// PostIDsByUserIDs 批量查询用户对应的岗位 ID。
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

// FindByIDInScope 查询当前数据范围内的用户。
func (r *Repository) FindByIDInScope(db *gorm.DB, actor datascope.Actor, userID uint) (model.User, error) {
	var user model.User
	err := applyDataScope(db, actor).First(&user, userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.User{}, apperror.NotFound("用户不存在")
		}
		return model.User{}, err
	}

	return user, nil
}

// UsernameExists 判断用户名是否已存在。
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

// DepartmentUsable 校验部门是否存在且可用。
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
		return apperror.BadRequest("部门不存在或已禁用")
	}

	return nil
}

// RolesUsable 校验角色是否存在且可用。
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
		return apperror.BadRequest("角色不存在或已禁用")
	}

	return nil
}

// PostsUsable 校验岗位是否存在且可用。
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
		return apperror.BadRequest("岗位不存在或已禁用")
	}

	return nil
}

// Create 创建用户。
func (r *Repository) Create(db *gorm.DB, user *model.User) error {
	return db.Create(user).Error
}

// UpdateBase 更新用户基础字段。
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

// UpdateStatus 更新用户状态。
func (r *Repository) UpdateStatus(db *gorm.DB, user *model.User, status model.UserStatus) error {
	if err := db.Model(user).Update("status", status).Error; err != nil {
		return err
	}

	user.Status = status
	return nil
}

// ReplaceRoles 用整体替换的方式刷新用户角色集合。
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

// ReplacePosts 用整体替换的方式刷新用户岗位集合。
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
