package application

import (
	"context"
	"testing"

	roledomain "ez-admin-gin/server/internal/modules/iam/role/domain"
	"ez-admin-gin/server/internal/platform/model"

	"gorm.io/gorm"
)

type roleTestTransactor struct {
	err error
}

func (t roleTestTransactor) WithinTransaction(_ context.Context, fn func(tx *gorm.DB) error) error {
	if t.err != nil {
		return t.err
	}
	return fn(nil)
}

type roleTestRepo struct {
	role        model.Role
	codeExists  bool
	findErr     error
	codeExistsV bool
}

func (r *roleTestRepo) List(query roledomain.ListQuery, page int, pageSize int) ([]model.Role, int64, error) {
	return nil, 0, nil
}
func (r *roleTestRepo) FindByID(db *gorm.DB, roleID uint) (model.Role, error) {
	return r.role, r.findErr
}
func (r *roleTestRepo) CodeExists(db *gorm.DB, code string) (bool, error) { return r.codeExistsV, nil }
func (r *roleTestRepo) DepartmentsUsable(db *gorm.DB, ids []uint) error   { return nil }
func (r *roleTestRepo) MenusUsable(db *gorm.DB, ids []uint) error         { return nil }
func (r *roleTestRepo) Create(db *gorm.DB, role *model.Role) error        { return nil }
func (r *roleTestRepo) UpdateBase(db *gorm.DB, role *model.Role, req roledomain.UpdateRequest) error {
	return nil
}
func (r *roleTestRepo) UpdateStatus(db *gorm.DB, role *model.Role, status model.RoleStatus) error {
	return nil
}
func (r *roleTestRepo) RolePermissions(codes []string) (map[string][]roledomain.PermissionItem, error) {
	return nil, nil
}
func (r *roleTestRepo) RoleMenuIDs(ids []uint) (map[uint][]uint, error)       { return nil, nil }
func (r *roleTestRepo) RoleCustomDepartmentIDs(ids []uint) (map[uint][]uint, error) { return nil, nil }
func (r *roleTestRepo) ReplacePermissions(db *gorm.DB, code string, perms []roledomain.PermissionItem) error {
	return nil
}
func (r *roleTestRepo) ReplaceMenus(db *gorm.DB, roleID uint, ids []uint) error         { return nil }
func (r *roleTestRepo) ReplaceCustomDepartments(db *gorm.DB, roleID uint, ids []uint) error { return nil }

func TestUpdateStatus_SuperAdminCannotBeDisabled(t *testing.T) {
	repo := &roleTestRepo{
		role: model.Role{ID: 1, Code: roledomain.SuperAdminRoleCode, Status: model.RoleStatusEnabled},
	}
	svc := NewService(roleTestTransactor{}, repo)

	err := svc.UpdateStatus(1, model.RoleStatusDisabled)
	if err == nil {
		t.Fatal("expected error when disabling super admin")
	}
}

func TestUpdate_SuperAdminScopeCannotChange(t *testing.T) {
	repo := &roleTestRepo{
		role: model.Role{ID: 1, Code: roledomain.SuperAdminRoleCode, Status: model.RoleStatusEnabled},
	}
	svc := NewService(roleTestTransactor{}, repo)

	_, err := svc.Update(1, roledomain.UpdateRequest{
		Name:       "Super Admin",
		Sort:       0,
		DataScope:  "custom",
		Status:     model.RoleStatusEnabled,
		Remark:     "",
	})
	if err == nil {
		t.Fatal("expected error when changing super admin data scope")
	}
}

func TestUpdatePermissions_SuperAdminBlocked(t *testing.T) {
	repo := &roleTestRepo{
		role: model.Role{ID: 1, Code: roledomain.SuperAdminRoleCode},
	}
	svc := NewService(roleTestTransactor{}, repo)

	_, _, err := svc.UpdatePermissions(1, []roledomain.PermissionItem{{Path: "/test", Method: "GET"}})
	if err == nil {
		t.Fatal("expected error when updating super admin permissions")
	}
}
