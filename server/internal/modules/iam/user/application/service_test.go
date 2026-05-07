package application

import (
	"context"
	"testing"

	userdomain "ez-admin-gin/server/internal/modules/iam/user/domain"
	"ez-admin-gin/server/internal/platform/datascope"
	"ez-admin-gin/server/internal/platform/model"

	"gorm.io/gorm"
)

type userTestTransactor struct{}

func (userTestTransactor) WithinTransaction(_ context.Context, fn func(tx *gorm.DB) error) error {
	return fn(nil)
}

type userTestRepo struct {
	replacedRoleIDs []uint
}

func (r *userTestRepo) List(actor datascope.Actor, query userdomain.ListQuery, page int, pageSize int) ([]model.User, int64, error) {
	return nil, 0, nil
}

func (r *userTestRepo) RoleIDsByUserIDs(userIDs []uint) (map[uint][]uint, error) {
	return map[uint][]uint{}, nil
}
func (r *userTestRepo) PostIDsByUserIDs(userIDs []uint) (map[uint][]uint, error) {
	return map[uint][]uint{}, nil
}
func (r *userTestRepo) FindByIDInScope(db *gorm.DB, actor datascope.Actor, userID uint) (model.User, error) {
	return model.User{ID: userID}, nil
}
func (r *userTestRepo) UsernameExists(db *gorm.DB, username string) (bool, error) { return false, nil }
func (r *userTestRepo) DepartmentUsable(db *gorm.DB, departmentID uint) error     { return nil }
func (r *userTestRepo) RolesUsable(db *gorm.DB, roleIDs []uint) error             { return nil }
func (r *userTestRepo) PostsUsable(db *gorm.DB, postIDs []uint) error             { return nil }
func (r *userTestRepo) Create(db *gorm.DB, user *model.User) error                { return nil }
func (r *userTestRepo) UpdateBase(db *gorm.DB, user *model.User, nickname string, departmentID uint, status model.UserStatus) error {
	return nil
}
func (r *userTestRepo) UpdateStatus(db *gorm.DB, user *model.User, status model.UserStatus) error {
	return nil
}
func (r *userTestRepo) ReplaceRoles(db *gorm.DB, userID uint, roleIDs []uint) error {
	r.replacedRoleIDs = append([]uint(nil), roleIDs...)
	return nil
}
func (r *userTestRepo) ReplacePosts(db *gorm.DB, userID uint, postIDs []uint) error { return nil }

func TestUpdateStatusRejectsDisablingCurrentUser(t *testing.T) {
	repo := &userTestRepo{}
	service := NewService(userTestTransactor{}, repo)

	err := service.UpdateStatus(datascope.Actor{}, 7, 7, uint(model.UserStatusDisabled))
	if err == nil || err.Error() != "不能禁用当前登录用户" {
		t.Fatalf("expected current-user disable rejection, got %v", err)
	}
}

func TestUpdateRolesNormalizesAndPersistsRoleIDs(t *testing.T) {
	repo := &userTestRepo{}
	service := NewService(userTestTransactor{}, repo)

	roleIDs, err := service.UpdateRoles(datascope.Actor{}, 9, 1, []uint{2, 2, 1})
	if err != nil {
		t.Fatalf("UpdateRoles returned error: %v", err)
	}

	expected := []uint{2, 1}
	if len(roleIDs) != len(expected) {
		t.Fatalf("expected %v, got %v", expected, roleIDs)
	}
	for index, item := range expected {
		if roleIDs[index] != item {
			t.Fatalf("expected %v, got %v", expected, roleIDs)
		}
		if repo.replacedRoleIDs[index] != item {
			t.Fatalf("expected persisted role IDs %v, got %v", expected, repo.replacedRoleIDs)
		}
	}
}
