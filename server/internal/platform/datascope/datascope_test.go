package datascope_test

import (
	"fmt"
	"sort"
	"strings"
	"testing"

	"ez-admin-gin/server/internal/platform/datascope"
	"ez-admin-gin/server/internal/platform/model"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func TestDepartmentQueryScopeIncludesDeptTree(t *testing.T) {
	db := newScopeTestDB(t)
	seedDepartments(t, db,
		model.Department{ID: 10, ParentID: 0, Ancestors: "0", Name: "销售中心", Code: "sales-root", Sort: 10, Status: model.DepartmentStatusEnabled},
		model.Department{ID: 11, ParentID: 10, Ancestors: "0,10", Name: "销售一部", Code: "sales-a", Sort: 10, Status: model.DepartmentStatusEnabled},
		model.Department{ID: 12, ParentID: 11, Ancestors: "0,10,11", Name: "KA 小组", Code: "sales-ka", Sort: 10, Status: model.DepartmentStatusEnabled},
		model.Department{ID: 20, ParentID: 0, Ancestors: "0", Name: "交付中心", Code: "delivery-root", Sort: 20, Status: model.DepartmentStatusEnabled},
		model.Department{ID: 21, ParentID: 20, Ancestors: "0,20", Name: "交付一部", Code: "delivery-a", Sort: 10, Status: model.DepartmentStatusEnabled},
	)

	actor := datascope.Actor{
		UserID:       100,
		DepartmentID: 10,
		Grants: []datascope.Grant{
			{Scope: datascope.ScopeDeptAndChildren},
		},
	}

	var items []model.Department
	err := db.Model(&model.Department{}).
		Scopes(datascope.DepartmentQueryScope(db, actor, "id")).
		Order("id ASC").
		Find(&items).Error
	if err != nil {
		t.Fatalf("query scoped departments: %v", err)
	}

	ids := make([]uint, 0, len(items))
	for _, item := range items {
		ids = append(ids, item.ID)
	}
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })

	expected := []uint{10, 11, 12}
	if len(ids) != len(expected) {
		t.Fatalf("expected visible departments %v, got %v", expected, ids)
	}
	for idx, item := range expected {
		if ids[idx] != item {
			t.Fatalf("expected visible departments %v, got %v", expected, ids)
		}
	}
}

func TestDepartmentQueryScopeSelfFallsBackToOwnDepartment(t *testing.T) {
	db := newScopeTestDB(t)
	seedDepartments(t, db,
		model.Department{ID: 10, ParentID: 0, Ancestors: "0", Name: "销售中心", Code: "sales-root", Sort: 10, Status: model.DepartmentStatusEnabled},
		model.Department{ID: 20, ParentID: 0, Ancestors: "0", Name: "交付中心", Code: "delivery-root", Sort: 20, Status: model.DepartmentStatusEnabled},
	)

	actor := datascope.Actor{
		UserID:       101,
		DepartmentID: 10,
		Grants: []datascope.Grant{
			{Scope: datascope.ScopeSelf},
		},
	}

	var items []model.Department
	err := db.Model(&model.Department{}).
		Scopes(datascope.DepartmentQueryScope(db, actor, "id")).
		Order("id ASC").
		Find(&items).Error
	if err != nil {
		t.Fatalf("query self scoped departments: %v", err)
	}

	if len(items) != 1 || items[0].ID != 10 {
		t.Fatalf("expected self scope to keep only current department 10, got %+v", items)
	}
}

func newScopeTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	dsn := fmt.Sprintf(
		"file:%s?mode=memory&cache=shared",
		strings.NewReplacer("/", "_", " ", "_").Replace(t.Name()),
	)
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite database: %v", err)
	}
	if err := db.AutoMigrate(&model.Department{}); err != nil {
		t.Fatalf("auto migrate department schema: %v", err)
	}
	return db
}

func seedDepartments(t *testing.T, db *gorm.DB, items ...model.Department) {
	t.Helper()

	for _, item := range items {
		if err := db.Create(&item).Error; err != nil {
			t.Fatalf("create department %s: %v", item.Code, err)
		}
	}
}
