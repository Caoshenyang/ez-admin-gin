package setup

import (
	"ez-admin-gin/server/internal/platform/database"
	"ez-admin-gin/server/internal/platform/model"

	"gorm.io/gorm"
)

type SetupTransactor = database.Transactor

type RepositoryPort interface {
	CountUsers(tx *gorm.DB) (int64, error)
	FindEnabledRoleByCode(tx *gorm.DB, code string) (model.Role, error)
	CreateUser(tx *gorm.DB, user *model.User) error
	CreateUserRole(tx *gorm.DB, binding *model.UserRole) error
}
