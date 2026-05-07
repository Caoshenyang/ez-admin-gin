package application

import (
	"context"

	userdomain "ez-admin-gin/server/internal/modules/iam/user/domain"
	errorsx "ez-admin-gin/server/internal/pkg/errorsx"
	"ez-admin-gin/server/internal/pkg/paging"
	"ez-admin-gin/server/internal/platform/datascope"
	"ez-admin-gin/server/internal/platform/model"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Service struct {
	tx   UserTransactor
	repo UserRepository
}

func NewService(tx UserTransactor, repo UserRepository) *Service {
	return &Service{
		tx:   tx,
		repo: repo,
	}
}

func (s *Service) List(actor datascope.Actor, query userdomain.ListQuery) (userdomain.ListResponse, error) {
	page, pageSize := paging.NormalizePage(query.Page, query.PageSize)

	users, total, err := s.repo.List(actor, query, page, pageSize)
	if err != nil {
		return userdomain.ListResponse{}, err
	}

	userIDs := make([]uint, 0, len(users))
	for _, item := range users {
		userIDs = append(userIDs, item.ID)
	}

	roleIDsByUser, err := s.repo.RoleIDsByUserIDs(userIDs)
	if err != nil {
		return userdomain.ListResponse{}, err
	}
	postIDsByUser, err := s.repo.PostIDsByUserIDs(userIDs)
	if err != nil {
		return userdomain.ListResponse{}, err
	}

	items := make([]userdomain.Response, 0, len(users))
	for _, item := range users {
		items = append(items, userdomain.BuildResponse(item, roleIDsByUser[item.ID], postIDsByUser[item.ID]))
	}

	return userdomain.ListResponse{
		Items:    items,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

func (s *Service) Create(actor datascope.Actor, req userdomain.CreateRequest) (userdomain.Response, error) {
	req, err := userdomain.NormalizeCreateRequest(req)
	if err != nil {
		return userdomain.Response{}, err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return userdomain.Response{}, errorsx.Internal("生成密码哈希失败", err)
	}

	var created userdomain.Entity
	err = s.tx.WithinTransaction(context.Background(), func(tx *gorm.DB) error {
		exists, err := s.repo.UsernameExists(tx, req.Username)
		if err != nil {
			return err
		}
		if exists {
			return errorsx.BadRequest("用户名已存在")
		}

		if err := s.repo.DepartmentUsable(tx, req.DepartmentID); err != nil {
			return err
		}
		if err := s.repo.RolesUsable(tx, req.RoleIDs); err != nil {
			return err
		}
		if err := s.repo.PostsUsable(tx, req.PostIDs); err != nil {
			return err
		}

		created = userdomain.Entity{
			Username:     req.Username,
			PasswordHash: string(passwordHash),
			Nickname:     req.Nickname,
			DepartmentID: req.DepartmentID,
			Status:       req.Status,
		}

		if err := s.repo.Create(tx, &created); err != nil {
			return err
		}
		if err := s.repo.ReplaceRoles(tx, created.ID, req.RoleIDs); err != nil {
			return err
		}

		return s.repo.ReplacePosts(tx, created.ID, req.PostIDs)
	})
	if err != nil {
		return userdomain.Response{}, err
	}

	_ = actor
	return userdomain.BuildResponse(created, req.RoleIDs, req.PostIDs), nil
}

func (s *Service) Update(actor datascope.Actor, userID uint, currentUserID uint, req userdomain.UpdateRequest) (userdomain.Response, error) {
	req, err := userdomain.NormalizeUpdateRequest(req)
	if err != nil {
		return userdomain.Response{}, err
	}
	if currentUserID == userID && req.Status == model.UserStatusDisabled {
		return userdomain.Response{}, errorsx.BadRequest("不能禁用当前登录用户")
	}

	var updated userdomain.Entity
	err = s.tx.WithinTransaction(context.Background(), func(tx *gorm.DB) error {
		user, err := s.repo.FindByIDInScope(tx, actor, userID)
		if err != nil {
			return err
		}
		if err := s.repo.DepartmentUsable(tx, req.DepartmentID); err != nil {
			return err
		}
		if err := s.repo.PostsUsable(tx, req.PostIDs); err != nil {
			return err
		}
		if err := s.repo.UpdateBase(tx, &user, req.Nickname, req.DepartmentID, req.Status); err != nil {
			return err
		}
		if err := s.repo.ReplacePosts(tx, user.ID, req.PostIDs); err != nil {
			return err
		}

		updated = user
		return nil
	})
	if err != nil {
		return userdomain.Response{}, err
	}

	roleIDsByUser, err := s.repo.RoleIDsByUserIDs([]uint{updated.ID})
	if err != nil {
		return userdomain.Response{}, err
	}
	postIDsByUser, err := s.repo.PostIDsByUserIDs([]uint{updated.ID})
	if err != nil {
		return userdomain.Response{}, err
	}

	return userdomain.BuildResponse(updated, roleIDsByUser[updated.ID], postIDsByUser[updated.ID]), nil
}

func (s *Service) UpdateStatus(actor datascope.Actor, userID uint, currentUserID uint, status uint) error {
	nextStatus := model.UserStatus(status)
	if !userdomain.ValidStatus(nextStatus) {
		return errorsx.BadRequest("用户状态不正确")
	}
	if currentUserID == userID && nextStatus == model.UserStatusDisabled {
		return errorsx.BadRequest("不能禁用当前登录用户")
	}

	return s.tx.WithinTransaction(context.Background(), func(tx *gorm.DB) error {
		user, err := s.repo.FindByIDInScope(tx, actor, userID)
		if err != nil {
			return err
		}

		return s.repo.UpdateStatus(tx, &user, nextStatus)
	})
}

func (s *Service) UpdateRoles(actor datascope.Actor, userID uint, currentUserID uint, roleIDs []uint) ([]uint, error) {
	if currentUserID == userID {
		return nil, errorsx.BadRequest("不能修改当前登录用户的角色")
	}

	normalizedRoleIDs, err := userdomain.NormalizeRoleIDs(roleIDs)
	if err != nil {
		return nil, err
	}

	err = s.tx.WithinTransaction(context.Background(), func(tx *gorm.DB) error {
		user, err := s.repo.FindByIDInScope(tx, actor, userID)
		if err != nil {
			return err
		}
		if err := s.repo.RolesUsable(tx, normalizedRoleIDs); err != nil {
			return err
		}
		return s.repo.ReplaceRoles(tx, user.ID, normalizedRoleIDs)
	})
	if err != nil {
		return nil, err
	}

	return normalizedRoleIDs, nil
}
