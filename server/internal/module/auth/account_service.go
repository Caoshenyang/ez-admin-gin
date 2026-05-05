package auth

import (
	"ez-admin-gin/server/internal/apperror"
	"ez-admin-gin/server/internal/model"
	"ez-admin-gin/server/internal/platform/datascope"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AccountService 负责当前登录人自助资料与密码修改。
type AccountService struct {
	db   *gorm.DB
	repo *Repository
}

// NewAccountService 创建账户中心服务。
func NewAccountService(db *gorm.DB, repo *Repository) *AccountService {
	return &AccountService{
		db:   db,
		repo: repo,
	}
}

// GetProfile 返回当前登录人的账户中心资料。
func (s *AccountService) GetProfile(actor datascope.Actor) (AccountProfileResponse, error) {
	row, err := s.repo.FindAccountProfileByID(actor.UserID)
	if err != nil {
		if s.repo.IsNotFound(err) {
			return AccountProfileResponse{}, apperror.Unauthorized("请先登录")
		}
		return AccountProfileResponse{}, err
	}

	return BuildAccountProfileResponse(
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

// UpdateProfile 修改当前登录人的昵称。
func (s *AccountService) UpdateProfile(actor datascope.Actor, req UpdateAccountProfileRequest) (AccountProfileResponse, error) {
	req, err := NormalizeUpdateAccountProfileRequest(req)
	if err != nil {
		return AccountProfileResponse{}, err
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		user, err := s.repo.FindUserByID(tx, actor.UserID)
		if err != nil {
			return err
		}

		return s.repo.UpdateAccountNickname(tx, &user, req.Nickname)
	})
	if err != nil {
		return AccountProfileResponse{}, err
	}

	updatedActor := actor
	return s.GetProfile(updatedActor)
}

// UpdatePassword 修改当前登录人的登录密码。
func (s *AccountService) UpdatePassword(actor datascope.Actor, req UpdateAccountPasswordRequest) error {
	req, err := NormalizeUpdateAccountPasswordRequest(req)
	if err != nil {
		return err
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		user, err := s.repo.FindUserByID(tx, actor.UserID)
		if err != nil {
			return err
		}
		if user.Status != model.UserStatusEnabled {
			return apperror.Forbidden("用户已被禁用")
		}
		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.OldPassword)); err != nil {
			return apperror.BadRequest("当前密码不正确")
		}

		passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			return apperror.Internal("生成密码哈希失败", err)
		}

		return s.repo.UpdateAccountPasswordHash(tx, &user, string(passwordHash))
	})
}
