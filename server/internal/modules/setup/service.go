package setup

import (
	"context"
	"errors"

	errorsx "ez-admin-gin/server/internal/pkg/errorsx"
	"ez-admin-gin/server/internal/platform/model"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Service struct {
	tx   SetupTransactor
	repo RepositoryPort
}

func NewService(tx SetupTransactor, repo RepositoryPort) *Service {
	return &Service{tx: tx, repo: repo}
}

type InitRequest struct {
	Username string
	Password string
	Nickname string
}

func (s *Service) Init(ctx context.Context, req InitRequest) (model.User, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return model.User{}, errorsx.Internal("密码加密失败", err)
	}

	var user model.User
	err = s.tx.WithinTransaction(ctx, func(tx *gorm.DB) error {
		count, err := s.repo.CountUsers(tx)
		if err != nil {
			return err
		}
		if count > 0 {
			return errAlreadyInitialized
		}

		role, err := s.repo.FindEnabledRoleByCode(tx, "super_admin")
		if err != nil {
			return err
		}

		user = model.User{
			Username:     req.Username,
			PasswordHash: string(passwordHash),
			Nickname:     req.Nickname,
			Status:       model.UserStatusEnabled,
		}
		if err := s.repo.CreateUser(tx, &user); err != nil {
			return err
		}

		return s.repo.CreateUserRole(tx, &model.UserRole{UserID: user.ID, RoleID: role.ID})
	})
	if err != nil {
		if errors.Is(err, errAlreadyInitialized) || errors.Is(err, gorm.ErrRecordNotFound) {
			return model.User{}, err
		}
		return model.User{}, errorsx.Internal("初始化管理员失败", err)
	}

	return user, nil
}
