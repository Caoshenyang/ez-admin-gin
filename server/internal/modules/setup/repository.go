package setup

import (
	"ez-admin-gin/server/internal/platform/model"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CountUsers(tx *gorm.DB) (int64, error) {
	var count int64
	err := r.dbOr(tx).Model(&model.User{}).Count(&count).Error
	return count, err
}

func (r *Repository) FindEnabledRoleByCode(tx *gorm.DB, code string) (model.Role, error) {
	var role model.Role
	err := r.dbOr(tx).
		Where("code = ?", code).
		Where("status = ?", model.RoleStatusEnabled).
		First(&role).Error
	return role, err
}

func (r *Repository) CreateUser(tx *gorm.DB, user *model.User) error {
	return r.dbOr(tx).Create(user).Error
}

func (r *Repository) CreateUserRole(tx *gorm.DB, binding *model.UserRole) error {
	return r.dbOr(tx).Create(binding).Error
}

func (r *Repository) dbOr(tx *gorm.DB) *gorm.DB {
	if tx != nil {
		return tx
	}
	return r.db
}
