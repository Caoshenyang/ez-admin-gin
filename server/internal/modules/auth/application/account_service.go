package application

import (
	"context"

	authdomain "ez-admin-gin/server/internal/modules/auth/domain"
	errorsx "ez-admin-gin/server/internal/pkg/errorsx"
	"ez-admin-gin/server/internal/platform/datascope"
	"ez-admin-gin/server/internal/platform/model"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AccountService 提供当前登录人的账户资料查询与修改服务。
type AccountService struct {
	tx   AuthTransactor
	repo AccountRepository
}

func NewAccountService(tx AuthTransactor, repo AccountRepository) *AccountService {
	return &AccountService{tx: tx, repo: repo}
}

// GetProfile 查询当前登录人的账户中心资料。
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

// UpdateProfile 更新当前登录人的昵称。
func (s *AccountService) UpdateProfile(
	actor datascope.Actor,
	req authdomain.UpdateAccountProfileRequest,
) (authdomain.AccountProfileResponse, error) {
	req, err := authdomain.NormalizeUpdateAccountProfileRequest(req)
	if err != nil {
		return authdomain.AccountProfileResponse{}, err
	}

	err = s.tx.WithinTransaction(context.Background(), func(tx *gorm.DB) error {
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

// UpdatePassword 校验旧密码后更新为新密码。
func (s *AccountService) UpdatePassword(actor datascope.Actor, req authdomain.UpdateAccountPasswordRequest) error {
	req, err := authdomain.NormalizeUpdateAccountPasswordRequest(req)
	if err != nil {
		return err
	}

	return s.tx.WithinTransaction(context.Background(), func(tx *gorm.DB) error {
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
