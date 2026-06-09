// Package application 实现角色的业务逻辑：分页列表、CRUD、权限和菜单分配。
package application

import (
	"fmt"

	roledomain "ez-admin-gin/server/internal/modules/iam/role/domain"
	errorsx "ez-admin-gin/server/internal/pkg/errorsx"
	"ez-admin-gin/server/internal/pkg/paging"
	"ez-admin-gin/server/internal/platform/model"

	"gorm.io/gorm"
)

// Service 提供角色的业务操作服务。
type Service struct {
	tx       RoleTransactor
	repo     RoleRepository
	enforcer PolicyReloader
}

func NewService(tx RoleTransactor, repo RoleRepository, enforcer PolicyReloader) *Service {
	return &Service{tx: tx, repo: repo, enforcer: enforcer}
}

// List 分页查询角色列表，并附带每个角色的权限和菜单信息。
func (s *Service) List(query roledomain.ListQuery) (roledomain.ListResponse, error) {
	page, pageSize := paging.NormalizePage(query.Page, query.PageSize)

	roles, total, err := s.repo.List(query, page, pageSize)
	if err != nil {
		return roledomain.ListResponse{}, err
	}

	roleIDs := make([]uint, 0, len(roles))
	for _, item := range roles {
		roleIDs = append(roleIDs, item.ID)
	}

	customDepartmentIDsByRole, err := s.repo.RoleCustomDepartmentIDs(roleIDs)
	if err != nil {
		return roledomain.ListResponse{}, err
	}
	apiIDsByRole, err := s.repo.RoleAPIIDs(roleIDs)
	if err != nil {
		return roledomain.ListResponse{}, err
	}
	permissionsByRole, err := s.repo.RolePermissions(roleIDs)
	if err != nil {
		return roledomain.ListResponse{}, err
	}
	menuIDsByRole, err := s.repo.RoleMenuIDs(roleIDs)
	if err != nil {
		return roledomain.ListResponse{}, err
	}

	items := make([]roledomain.Response, 0, len(roles))
	for _, item := range roles {
		items = append(items, roledomain.BuildResponse(item, customDepartmentIDsByRole[item.ID], permissionsByRole[item.ID], apiIDsByRole[item.ID], menuIDsByRole[item.ID]))
	}

	return roledomain.ListResponse{Items: items, Total: total, Page: page, PageSize: pageSize}, nil
}

// Create 创建角色并关联自定义数据范围部门。
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

	return roledomain.BuildResponse(created, req.CustomDepartmentIDs, nil, nil, nil), nil
}

// Update 更新角色基本信息和自定义数据范围。
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
		// 超级管理员角色的关键属性不允许通过业务接口修改。
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

	return roledomain.BuildResponse(updated, req.CustomDepartmentIDs, nil, nil, nil), nil
}

// UpdateStatus 切换角色的启用/禁用状态。
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

// UpdatePermissions 更新角色的接口权限关联，并同步 Casbin 执行策略。
func (s *Service) UpdatePermissions(roleID uint, apiIDs []uint) ([]uint, string, error) {
	normalizedAPIIDs, err := roledomain.NormalizeIDs(apiIDs, "接口权限 ID 不正确")
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
		if err := s.repo.APIsUsable(tx, normalizedAPIIDs); err != nil {
			return err
		}
		if err := s.repo.ReplaceAPIs(tx, role.ID, normalizedAPIIDs); err != nil {
			return err
		}
		return s.repo.ReplacePoliciesByAPIs(tx, role.Code, normalizedAPIIDs)
	})
	if err != nil {
		return nil, "", err
	}

	if err := s.enforcer.ReloadPolicy(); err != nil {
		return nil, "", fmt.Errorf("reload casbin policy: %w", err)
	}

	return normalizedAPIIDs, roleCode, nil
}

// UpdateMenus 更新角色的菜单分配。
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

// Delete 删除角色，要求不是超级管理员且未分配给任何用户。
func (s *Service) Delete(roleID uint) error {
	var roleCode string

	err := s.tx.WithinTransaction(nil, func(tx *gorm.DB) error {
		role, err := s.repo.FindByID(tx, roleID)
		if err != nil {
			return err
		}
		if role.Code == roledomain.SuperAdminRoleCode {
			return errorsx.BadRequest("不能删除超级管理员角色")
		}

		userCount, err := s.repo.CountUsers(tx, roleID)
		if err != nil {
			return err
		}
		if userCount > 0 {
			return errorsx.BadRequest("角色已分配给用户，不能删除")
		}

		roleCode = role.Code
		if err := s.repo.ReplaceAPIs(tx, role.ID, nil); err != nil {
			return err
		}
		if err := s.repo.ReplacePoliciesByAPIs(tx, role.Code, nil); err != nil {
			return err
		}
		if err := s.repo.ReplaceMenus(tx, role.ID, nil); err != nil {
			return err
		}
		if err := s.repo.ReplaceCustomDepartments(tx, role.ID, nil); err != nil {
			return err
		}

		return s.repo.Delete(tx, &role)
	})
	if err != nil {
		return err
	}

	if roleCode != "" {
		if err := s.enforcer.ReloadPolicy(); err != nil {
			return fmt.Errorf("reload casbin policy: %w", err)
		}
	}

	return nil
}
