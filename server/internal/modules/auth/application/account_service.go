package application

import (
	authdomain "ez-admin-gin/server/internal/modules/auth/domain"
	authinfra "ez-admin-gin/server/internal/modules/auth/infra"
	errorsx "ez-admin-gin/server/internal/pkg/errorsx"
	"ez-admin-gin/server/internal/platform/datascope"
	"ez-admin-gin/server/internal/platform/model"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AccountService struct {
	db   *gorm.DB
	repo *authinfra.Repository
}

func NewAccountService(db *gorm.DB, repo *authinfra.Repository) *AccountService {
	return &AccountService{db: db, repo: repo}
}

func (s *AccountService) GetProfile(actor datascope.Actor) (authdomain.AccountProfileResponse, error) {
	row, err := s.repo.FindAccountProfileByID(actor.UserID)
	if err != nil {
		if s.repo.IsNotFound(err) {
			return authdomain.AccountProfileResponse{}, errorsx.Unauthorized("请先登录")
		}
		return authdomain.AccountProfileResponse{}, err
	}

	return authdomain.BuildAccountProfileResponse(
		actor,
		model.User{
			ID:           row.ID,
			Username:     row.Username,
			Nickname:     row.Nickname,
			DepartmentID: row.DepartmentID,
			Status:       row.Status,
			UpdatedAt:    row.UpdatedAt,
		},
		row.DepartmentName,
	), nil
}

func (s *AccountService) UpdateProfile(
	actor datascope.Actor,
	req authdomain.UpdateAccountProfileRequest,
) (authdomain.AccountProfileResponse, error) {
	req, err := authdomain.NormalizeUpdateAccountProfileRequest(req)
	if err != nil {
		return authdomain.AccountProfileResponse{}, err
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		user, err := s.repo.FindUserByID(tx, actor.UserID)
		if err != nil {
			return err
		}
		return s.repo.UpdateAccountNickname(tx, &user, req.Nickname)
	})
	if err != nil {
		return authdomain.AccountProfileResponse{}, err
	}

	return s.GetProfile(actor)
}

func (s *AccountService) UpdatePassword(actor datascope.Actor, req authdomain.UpdateAccountPasswordRequest) error {
	req, err := authdomain.NormalizeUpdateAccountPasswordRequest(req)
	if err != nil {
		return err
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		user, err := s.repo.FindUserByID(tx, actor.UserID)
		if err != nil {
			return err
		}
		if user.Status != model.UserStatusEnabled {
			return errorsx.Forbidden("用户已被禁用")
		}
		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.OldPassword)); err != nil {
			return errorsx.BadRequest("当前密码不正确")
		}

		passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			return errorsx.Internal("生成密码哈希失败", err)
		}

		return s.repo.UpdateAccountPasswordHash(tx, &user, string(passwordHash))
	})
}
