package application

import (
	roledomain "ez-admin-gin/server/internal/modules/iam/role/domain"
	errorsx "ez-admin-gin/server/internal/pkg/errorsx"
	"ez-admin-gin/server/internal/pkg/paging"
	"ez-admin-gin/server/internal/platform/model"

	"gorm.io/gorm"
)

type Service struct {
	tx   RoleTransactor
	repo RoleRepository
}

func NewService(tx RoleTransactor, repo RoleRepository) *Service {
	return &Service{tx: tx, repo: repo}
}

func (s *Service) List(query roledomain.ListQuery) (roledomain.ListResponse, error) {
	page, pageSize := paging.NormalizePage(query.Page, query.PageSize)

	roles, total, err := s.repo.List(query, page, pageSize)
	if err != nil {
		return roledomain.ListResponse{}, err
	}

	roleIDs := make([]uint, 0, len(roles))
	roleCodes := make([]string, 0, len(roles))
	for _, item := range roles {
		roleIDs = append(roleIDs, item.ID)
		roleCodes = append(roleCodes, item.Code)
	}

	customDepartmentIDsByRole, err := s.repo.RoleCustomDepartmentIDs(roleIDs)
	if err != nil {
		return roledomain.ListResponse{}, err
	}
	permissionsByRole, err := s.repo.RolePermissions(roleCodes)
	if err != nil {
		return roledomain.ListResponse{}, err
	}
	menuIDsByRole, err := s.repo.RoleMenuIDs(roleIDs)
	if err != nil {
		return roledomain.ListResponse{}, err
	}

	items := make([]roledomain.Response, 0, len(roles))
	for _, item := range roles {
		items = append(items, roledomain.BuildResponse(item, customDepartmentIDsByRole[item.ID], permissionsByRole[item.Code], menuIDsByRole[item.ID]))
	}

	return roledomain.ListResponse{Items: items, Total: total, Page: page, PageSize: pageSize}, nil
}

func (s *Service) Create(req roledomain.CreateRequest) (roledomain.Response, error) {
	req, err := roledomain.NormalizeCreateRequest(req)
	if err != nil {
		return roledomain.Response{}, err
	}

	var created roledomain.Entity
	err = s.tx.WithinTransaction(nil, func(tx *gorm.DB) error {
		exists, err := s.repo.CodeExists(tx, req.Code)
		if err != nil {
			return err
		}
		if exists {
			return errorsx.BadRequest("角色编码已存在")
		}
		if err := s.repo.DepartmentsUsable(tx, req.CustomDepartmentIDs); err != nil {
			return err
		}

		created = roledomain.Entity{
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
		return roledomain.Response{}, err
	}

	return roledomain.BuildResponse(created, req.CustomDepartmentIDs, nil, nil), nil
}

func (s *Service) Update(roleID uint, req roledomain.UpdateRequest) (roledomain.Response, error) {
	req, err := roledomain.NormalizeUpdateRequest(req)
	if err != nil {
		return roledomain.Response{}, err
	}

	var updated roledomain.Entity
	err = s.tx.WithinTransaction(nil, func(tx *gorm.DB) error {
		role, err := s.repo.FindByID(tx, roleID)
		if err != nil {
			return err
		}
		if role.Code == roledomain.SuperAdminRoleCode && req.Status == model.RoleStatusDisabled {
			return errorsx.BadRequest("不能禁用超级管理员角色")
		}
		if role.Code == roledomain.SuperAdminRoleCode && req.DataScope != role.DataScope {
			return errorsx.BadRequest("不能修改超级管理员角色的数据范围")
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
		return roledomain.Response{}, err
	}

	return roledomain.BuildResponse(updated, req.CustomDepartmentIDs, nil, nil), nil
}

func (s *Service) UpdateStatus(roleID uint, status model.RoleStatus) error {
	if !roledomain.ValidRoleStatus(status) {
		return errorsx.BadRequest("角色状态不正确")
	}

	return s.tx.WithinTransaction(nil, func(tx *gorm.DB) error {
		role, err := s.repo.FindByID(tx, roleID)
		if err != nil {
			return err
		}
		if role.Code == roledomain.SuperAdminRoleCode && status == model.RoleStatusDisabled {
			return errorsx.BadRequest("不能禁用超级管理员角色")
		}

		return s.repo.UpdateStatus(tx, &role, status)
	})
}

func (s *Service) UpdatePermissions(roleID uint, permissions []roledomain.PermissionItem) ([]roledomain.PermissionItem, string, error) {
	normalizedPermissions, err := roledomain.NormalizePermissions(permissions)
	if err != nil {
		return nil, "", err
	}

	var roleCode string
	err = s.tx.WithinTransaction(nil, func(tx *gorm.DB) error {
		role, err := s.repo.FindByID(tx, roleID)
		if err != nil {
			return err
		}
		if role.Code == roledomain.SuperAdminRoleCode {
			return errorsx.BadRequest("超级管理员角色权限不在这里修改")
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
	normalizedMenuIDs, err := roledomain.NormalizeIDs(menuIDs, "菜单 ID 不正确")
	if err != nil {
		return nil, err
	}

	err = s.tx.WithinTransaction(nil, func(tx *gorm.DB) error {
		role, err := s.repo.FindByID(tx, roleID)
		if err != nil {
			return err
		}
		if role.Code == roledomain.SuperAdminRoleCode {
			return errorsx.BadRequest("超级管理员菜单权限不在这里修改")
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
