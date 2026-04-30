package role

import (
	"ez-admin-gin/server/internal/apperror"
	"ez-admin-gin/server/internal/model"

	"gorm.io/gorm"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) List(query ListQuery) (ListResponse, error) {
	page, pageSize := NormalizePage(query.Page, query.PageSize)

	roles, total, err := s.repo.List(query, page, pageSize)
	if err != nil {
		return ListResponse{}, err
	}

	roleIDs := make([]uint, 0, len(roles))
	roleCodes := make([]string, 0, len(roles))
	for _, item := range roles {
		roleIDs = append(roleIDs, item.ID)
		roleCodes = append(roleCodes, item.Code)
	}

	customDepartmentIDsByRole, err := s.repo.RoleCustomDepartmentIDs(roleIDs)
	if err != nil {
		return ListResponse{}, err
	}
	permissionsByRole, err := s.repo.RolePermissions(roleCodes)
	if err != nil {
		return ListResponse{}, err
	}
	menuIDsByRole, err := s.repo.RoleMenuIDs(roleIDs)
	if err != nil {
		return ListResponse{}, err
	}

	items := make([]Response, 0, len(roles))
	for _, item := range roles {
		items = append(items, BuildResponse(item, customDepartmentIDsByRole[item.ID], permissionsByRole[item.Code], menuIDsByRole[item.ID]))
	}

	return ListResponse{Items: items, Total: total, Page: page, PageSize: pageSize}, nil
}

func (s *Service) Create(req CreateRequest) (Response, error) {
	req, err := NormalizeCreateRequest(req)
	if err != nil {
		return Response{}, err
	}

	var created Entity
	err = s.repo.db.Transaction(func(tx *gorm.DB) error {
		exists, err := s.repo.CodeExists(tx, req.Code)
		if err != nil {
			return err
		}
		if exists {
			return apperror.BadRequest("角色编码已存在")
		}
		if err := s.repo.DepartmentsUsable(tx, req.CustomDepartmentIDs); err != nil {
			return err
		}

		created = Entity{
			Code:      req.Code,
			Name:      req.Name,
			Sort:      req.Sort,
			DataScope: req.DataScope,
			Status:    req.Status,
			Remark:    req.Remark,
		}
		if err := s.repo.Create(tx, &created); err != nil {
			return err
		}

		return s.repo.ReplaceCustomDepartments(tx, created.ID, req.CustomDepartmentIDs)
	})
	if err != nil {
		return Response{}, err
	}

	return BuildResponse(created, req.CustomDepartmentIDs, nil, nil), nil
}

func (s *Service) Update(roleID uint, req UpdateRequest) (Response, error) {
	req, err := NormalizeUpdateRequest(req)
	if err != nil {
		return Response{}, err
	}

	var updated Entity
	err = s.repo.db.Transaction(func(tx *gorm.DB) error {
		role, err := s.repo.FindByID(tx, roleID)
		if err != nil {
			return err
		}
		if role.Code == superAdminRoleCode {
			if req.Status == model.RoleStatusDisabled {
				return apperror.BadRequest("不能禁用超级管理员角色")
			}
			if req.DataScope != role.DataScope {
				return apperror.BadRequest("不能修改超级管理员角色的数据范围")
			}
		}
		if err := s.repo.DepartmentsUsable(tx, req.CustomDepartmentIDs); err != nil {
			return err
		}
		if err := s.repo.UpdateBase(tx, &role, req); err != nil {
			return err
		}
		if err := s.repo.ReplaceCustomDepartments(tx, role.ID, req.CustomDepartmentIDs); err != nil {
			return err
		}

		updated = role
		return nil
	})
	if err != nil {
		return Response{}, err
	}

	return BuildResponse(updated, req.CustomDepartmentIDs, nil, nil), nil
}

func (s *Service) UpdateStatus(roleID uint, status model.RoleStatus) error {
	if !ValidRoleStatus(status) {
		return apperror.BadRequest("角色状态不正确")
	}

	return s.repo.db.Transaction(func(tx *gorm.DB) error {
		role, err := s.repo.FindByID(tx, roleID)
		if err != nil {
			return err
		}
		if role.Code == superAdminRoleCode && status == model.RoleStatusDisabled {
			return apperror.BadRequest("不能禁用超级管理员角色")
		}

		return s.repo.UpdateStatus(tx, &role, status)
	})
}

func (s *Service) UpdatePermissions(roleID uint, permissions []PermissionItem) ([]PermissionItem, string, error) {
	normalizedPermissions, err := NormalizePermissions(permissions)
	if err != nil {
		return nil, "", err
	}

	var roleCode string
	err = s.repo.db.Transaction(func(tx *gorm.DB) error {
		role, err := s.repo.FindByID(tx, roleID)
		if err != nil {
			return err
		}
		if role.Code == superAdminRoleCode {
			return apperror.BadRequest("超级管理员角色权限不在这里修改")
		}

		roleCode = role.Code
		return s.repo.ReplacePermissions(tx, role.Code, normalizedPermissions)
	})
	if err != nil {
		return nil, "", err
	}

	return normalizedPermissions, roleCode, nil
}

func (s *Service) UpdateMenus(roleID uint, menuIDs []uint) ([]uint, error) {
	normalizedMenuIDs, err := NormalizeIDs(menuIDs, "菜单 ID 不正确")
	if err != nil {
		return nil, err
	}

	err = s.repo.db.Transaction(func(tx *gorm.DB) error {
		role, err := s.repo.FindByID(tx, roleID)
		if err != nil {
			return err
		}
		if role.Code == superAdminRoleCode {
			return apperror.BadRequest("超级管理员菜单权限不在这里修改")
		}
		if err := s.repo.MenusUsable(tx, normalizedMenuIDs); err != nil {
			return err
		}

		return s.repo.ReplaceMenus(tx, roleID, normalizedMenuIDs)
	})
	if err != nil {
		return nil, err
	}

	return normalizedMenuIDs, nil
}
