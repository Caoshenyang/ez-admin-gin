package user

import (
	"ez-admin-gin/server/internal/apperror"
	"ez-admin-gin/server/internal/model"
	"ez-admin-gin/server/internal/platform/datascope"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Service 负责用户模块的业务规则、事务边界和跨仓储协作。
type Service struct {
	db   *gorm.DB
	repo *Repository
}

// NewService 创建用户服务。
func NewService(db *gorm.DB, repo *Repository) *Service {
	return &Service{
		db:   db,
		repo: repo,
	}
}

// List 返回当前数据范围内的用户分页结果。
func (s *Service) List(actor datascope.Actor, query ListQuery) (ListResponse, error) {
	page, pageSize := NormalizePage(query.Page, query.PageSize)

	users, total, err := s.repo.List(actor, query, page, pageSize)
	if err != nil {
		return ListResponse{}, err
	}

	userIDs := make([]uint, 0, len(users))
	for _, item := range users {
		userIDs = append(userIDs, item.ID)
	}

	roleIDsByUser, err := s.repo.RoleIDsByUserIDs(userIDs)
	if err != nil {
		return ListResponse{}, err
	}
	postIDsByUser, err := s.repo.PostIDsByUserIDs(userIDs)
	if err != nil {
		return ListResponse{}, err
	}

	items := make([]Response, 0, len(users))
	for _, item := range users {
		items = append(items, BuildResponse(item, roleIDsByUser[item.ID], postIDsByUser[item.ID]))
	}

	return ListResponse{
		Items:    items,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

// Create 创建后台用户。
func (s *Service) Create(actor datascope.Actor, req CreateRequest) (Response, error) {
	req, err := NormalizeCreateRequest(req)
	if err != nil {
		return Response{}, err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return Response{}, apperror.Internal("生成密码哈希失败", err)
	}

	var created Entity
	err = s.db.Transaction(func(tx *gorm.DB) error {
		exists, err := s.repo.UsernameExists(tx, req.Username)
		if err != nil {
			return err
		}
		if exists {
			return apperror.BadRequest("用户名已存在")
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

		created = Entity{
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
		return Response{}, err
	}

	return BuildResponse(created, req.RoleIDs, req.PostIDs), nil
}

// Update 编辑用户基础信息。
func (s *Service) Update(actor datascope.Actor, userID uint, currentUserID uint, req UpdateRequest) (Response, error) {
	req, err := NormalizeUpdateRequest(req)
	if err != nil {
		return Response{}, err
	}
	if currentUserID == userID && req.Status == 2 {
		return Response{}, apperror.BadRequest("不能禁用当前登录用户")
	}

	var updated Entity
	err = s.db.Transaction(func(tx *gorm.DB) error {
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
		return Response{}, err
	}

	roleIDsByUser, err := s.repo.RoleIDsByUserIDs([]uint{updated.ID})
	if err != nil {
		return Response{}, err
	}
	postIDsByUser, err := s.repo.PostIDsByUserIDs([]uint{updated.ID})
	if err != nil {
		return Response{}, err
	}

	return BuildResponse(updated, roleIDsByUser[updated.ID], postIDsByUser[updated.ID]), nil
}

// UpdateStatus 单独修改用户状态。
func (s *Service) UpdateStatus(actor datascope.Actor, userID uint, currentUserID uint, status uint) error {
	nextStatus := model.UserStatus(status)
	if !ValidStatus(nextStatus) {
		return apperror.BadRequest("用户状态不正确")
	}
	if currentUserID == userID && nextStatus == model.UserStatusDisabled {
		return apperror.BadRequest("不能禁用当前登录用户")
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		user, err := s.repo.FindByIDInScope(tx, actor, userID)
		if err != nil {
			return err
		}

		return s.repo.UpdateStatus(tx, &user, nextStatus)
	})
}

// UpdateRoles 更新用户角色集合。
func (s *Service) UpdateRoles(actor datascope.Actor, userID uint, currentUserID uint, roleIDs []uint) ([]uint, error) {
	if currentUserID == userID {
		return nil, apperror.BadRequest("不能修改当前登录用户的角色")
	}

	normalizedRoleIDs, err := NormalizeRoleIDs(roleIDs)
	if err != nil {
		return nil, err
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
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
