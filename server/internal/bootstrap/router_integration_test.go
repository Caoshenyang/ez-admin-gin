package bootstrap

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	authnPlatform "ez-admin-gin/server/internal/platform/authn"
	authzPlatform "ez-admin-gin/server/internal/platform/authz"
	platformConfig "ez-admin-gin/server/internal/platform/config"
	"ez-admin-gin/server/internal/platform/model"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type apiResponse[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

type loginResponseData struct {
	UserID      uint   `json:"user_id"`
	Username    string `json:"username"`
	Nickname    string `json:"nickname"`
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresAt   string `json:"expires_at"`
}

type meResponseData struct {
	UserID       uint     `json:"user_id"`
	Username     string   `json:"username"`
	DepartmentID uint     `json:"department_id"`
	RoleCodes    []string `json:"role_codes"`
	IsSuperAdmin bool     `json:"is_super_admin"`
	DataScope    struct {
		AllowAll          bool `json:"allow_all"`
		RequireSelf       bool `json:"require_self"`
		IncludeDepartment bool `json:"include_department"`
		IncludeDeptTree   bool `json:"include_dept_tree"`
	} `json:"data_scope"`
}

type accountProfileResponseData struct {
	UserID         uint     `json:"user_id"`
	Username       string   `json:"username"`
	Nickname       string   `json:"nickname"`
	DepartmentID   uint     `json:"department_id"`
	DepartmentName string   `json:"department_name"`
	Status         int      `json:"status"`
	RoleCodes      []string `json:"role_codes"`
	IsSuperAdmin   bool     `json:"is_super_admin"`
	DataScope      struct {
		AllowAll            bool   `json:"allow_all"`
		RequireSelf         bool   `json:"require_self"`
		IncludeDepartment   bool   `json:"include_department"`
		IncludeDeptTree     bool   `json:"include_dept_tree"`
		CustomDepartmentIDs []uint `json:"custom_department_ids"`
	} `json:"data_scope"`
}

type userListItem struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
}

type userListResponseData struct {
	Items []userListItem `json:"items"`
	Total int64          `json:"total"`
}

type dictTypeResponseData struct {
	ID     uint   `json:"id"`
	Code   string `json:"code"`
	Name   string `json:"name"`
	Status int    `json:"status"`
}

type dictTypeListResponseData struct {
	Items []dictTypeResponseData `json:"items"`
	Total int64                  `json:"total"`
}

type dictItemResponseData struct {
	ID      uint   `json:"id"`
	TypeID  uint   `json:"type_id"`
	ItemKey string `json:"item_key"`
	Label   string `json:"label"`
	Value   string `json:"value"`
	Status  int    `json:"status"`
}

type dictItemListResponseData struct {
	Items []dictItemResponseData `json:"items"`
	Total int64                  `json:"total"`
}

type attachmentResponseData struct {
	ID           uint   `json:"id"`
	FileID       uint   `json:"file_id"`
	DisplayName  string `json:"display_name"`
	Category     string `json:"category"`
	BizType      string `json:"biz_type"`
	OriginalName string `json:"original_name"`
	FileName     string `json:"file_name"`
	Ext          string `json:"ext"`
	MimeType     string `json:"mime_type"`
	Size         int64  `json:"size"`
	URL          string `json:"url"`
	UploaderID   uint   `json:"uploader_id"`
	Status       int    `json:"status"`
	Remark       string `json:"remark"`
}

type attachmentListResponseData struct {
	Items []attachmentResponseData `json:"items"`
	Total int64                    `json:"total"`
}

type departmentResponseData struct {
	ID       uint                     `json:"id"`
	ParentID uint                     `json:"parent_id"`
	Name     string                   `json:"name"`
	Code     string                   `json:"code"`
	Status   int                      `json:"status"`
	Children []departmentResponseData `json:"children"`
}

type postResponseData struct {
	ID     uint   `json:"id"`
	Code   string `json:"code"`
	Name   string `json:"name"`
	Status int    `json:"status"`
}

type menuResponseData struct {
	ID        uint               `json:"id"`
	ParentID  uint               `json:"parent_id"`
	Type      model.MenuType     `json:"type"`
	Code      string             `json:"code"`
	Title     string             `json:"title"`
	Path      string             `json:"path"`
	Component string             `json:"component"`
	Icon      string             `json:"icon"`
	Sort      int                `json:"sort"`
	Children  []menuResponseData `json:"children"`
}

type testEnv struct {
	t      *testing.T
	db     *gorm.DB
	router *gin.Engine
	token  *authnPlatform.Manager
}

func TestSetupInitSuccessAndDuplicateReject(t *testing.T) {
	env := newTestEnv(t, func(db *gorm.DB) {
		mustCreateRole(t, db, model.Role{
			ID:        1,
			Code:      "super_admin",
			Name:      "超级管理员",
			Sort:      1,
			DataScope: "all",
			Status:    model.RoleStatusEnabled,
		})
	})

	body := map[string]any{
		"username": "admin",
		"password": "Admin12345",
		"nickname": "管理员",
	}

	first := env.doJSON(http.MethodPost, "/api/v1/setup/init", "", body)
	if first.Code != http.StatusOK {
		t.Fatalf("expected first init status 200, got %d: %s", first.Code, first.Body.String())
	}

	var firstBody map[string]any
	decodeJSON(t, first.Body.Bytes(), &firstBody)
	if firstBody["message"] != "管理员账号创建成功" {
		t.Fatalf("expected success message, got %#v", firstBody["message"])
	}

	second := env.doJSON(http.MethodPost, "/api/v1/setup/init", "", body)
	if second.Code != http.StatusConflict {
		t.Fatalf("expected second init status 409, got %d: %s", second.Code, second.Body.String())
	}

	var secondBody map[string]any
	decodeJSON(t, second.Body.Bytes(), &secondBody)
	if secondBody["error"] != "系统已初始化，不能重复执行" {
		t.Fatalf("expected duplicate init message, got %#v", secondBody["error"])
	}
}

func TestSetupInitRequiresSuperAdminSeedAndStaysAtomic(t *testing.T) {
	env := newTestEnv(t, nil)

	resp := env.doJSON(http.MethodPost, "/api/v1/setup/init", "", map[string]any{
		"username": "admin",
		"password": "Admin12345",
		"nickname": "管理员",
	})
	if resp.Code != http.StatusInternalServerError {
		t.Fatalf("expected setup init without super_admin seed status 500, got %d: %s", resp.Code, resp.Body.String())
	}

	var body map[string]any
	decodeJSON(t, resp.Body.Bytes(), &body)
	if body["error"] != "初始化角色不存在，请先检查种子数据" {
		t.Fatalf("expected missing super_admin seed message, got %#v", body["error"])
	}

	var count int64
	if err := env.db.Model(&model.User{}).Count(&count).Error; err != nil {
		t.Fatalf("count users after failed setup init: %v", err)
	}
	if count != 0 {
		t.Fatalf("expected failed setup init to stay atomic and create no users, got %d", count)
	}
}

func TestAuthLoginSuccessAndFailure(t *testing.T) {
	env := newTestEnv(t, func(db *gorm.DB) {
		mustCreateRole(t, db, model.Role{
			ID:        1,
			Code:      "super_admin",
			Name:      "超级管理员",
			Sort:      1,
			DataScope: "all",
			Status:    model.RoleStatusEnabled,
		})
		admin := mustCreateUser(t, db, seededUser{
			Username:     "admin",
			Password:     "Admin12345",
			Nickname:     "管理员",
			DepartmentID: 10,
			Status:       model.UserStatusEnabled,
		})
		mustCreateUserRole(t, db, admin.ID, 1)
	})

	success := env.doJSON(http.MethodPost, "/api/v1/auth/login", "", map[string]any{
		"username": "admin",
		"password": "Admin12345",
	})
	if success.Code != http.StatusOK {
		t.Fatalf("expected login success status 200, got %d: %s", success.Code, success.Body.String())
	}

	var loginResp apiResponse[loginResponseData]
	decodeJSON(t, success.Body.Bytes(), &loginResp)
	if loginResp.Code != 0 || loginResp.Data.AccessToken == "" || loginResp.Data.Username != "admin" {
		t.Fatalf("unexpected login success body: %+v", loginResp)
	}

	failure := env.doJSON(http.MethodPost, "/api/v1/auth/login", "", map[string]any{
		"username": "admin",
		"password": "wrong-password",
	})
	if failure.Code != http.StatusUnauthorized {
		t.Fatalf("expected login failure status 401, got %d: %s", failure.Code, failure.Body.String())
	}

	var failureResp apiResponse[map[string]any]
	decodeJSON(t, failure.Body.Bytes(), &failureResp)
	if failureResp.Code != 40100 || failureResp.Message != "用户名或密码错误" {
		t.Fatalf("unexpected login failure body: %+v", failureResp)
	}
}

func TestSetupInitSuperAdminCanSeeOrgMenusAndAccessOrgEndpoints(t *testing.T) {
	env := newTestEnv(t, func(db *gorm.DB) {
		mustCreateRole(t, db, model.Role{
			ID:        1,
			Code:      "super_admin",
			Name:      "超级管理员",
			Sort:      1,
			DataScope: "all",
			Status:    model.RoleStatusEnabled,
		})

		mustCreateMenu(t, db, model.Menu{
			ID:     100,
			Type:   model.MenuTypeDirectory,
			Code:   "system",
			Title:  "系统管理",
			Path:   "/system",
			Icon:   "setting",
			Sort:   10,
			Status: model.MenuStatusEnabled,
		})
		mustCreateMenu(t, db, model.Menu{
			ID:        211,
			ParentID:  100,
			Type:      model.MenuTypeMenu,
			Code:      "system:department",
			Title:     "部门管理",
			Path:      "/system/departments",
			Component: "system/DepartmentView",
			Icon:      "git-branch",
			Sort:      25,
			Status:    model.MenuStatusEnabled,
		})
		mustCreateMenu(t, db, model.Menu{
			ID:        212,
			ParentID:  100,
			Type:      model.MenuTypeMenu,
			Code:      "system:post",
			Title:     "岗位管理",
			Path:      "/system/posts",
			Component: "system/PostView",
			Icon:      "briefcase",
			Sort:      26,
			Status:    model.MenuStatusEnabled,
		})

		mustCreateRoleMenu(t, db, 1, 100)
		mustCreateRoleMenu(t, db, 1, 211)
		mustCreateRoleMenu(t, db, 1, 212)

		mustCreatePolicy(t, db, "super_admin", "/api/v1/system/departments", http.MethodGet)
		mustCreatePolicy(t, db, "super_admin", "/api/v1/system/posts", http.MethodGet)

		mustCreateDepartment(t, db, model.Department{
			ID:     10,
			Name:   "研发中心",
			Code:   "rd-center",
			Sort:   10,
			Status: model.DepartmentStatusEnabled,
		})
		if err := db.Create(&model.Post{
			ID:     20,
			Code:   "team-lead",
			Name:   "团队负责人",
			Sort:   5,
			Status: model.PostStatusEnabled,
		}).Error; err != nil {
			t.Fatalf("create post team-lead: %v", err)
		}
	})

	initResp := env.doJSON(http.MethodPost, "/api/v1/setup/init", "", map[string]any{
		"username": "admin",
		"password": "Admin12345",
		"nickname": "管理员",
	})
	if initResp.Code != http.StatusOK {
		t.Fatalf("expected setup init status 200, got %d: %s", initResp.Code, initResp.Body.String())
	}

	login := env.doJSON(http.MethodPost, "/api/v1/auth/login", "", map[string]any{
		"username": "admin",
		"password": "Admin12345",
	})
	if login.Code != http.StatusOK {
		t.Fatalf("expected login status 200, got %d: %s", login.Code, login.Body.String())
	}

	var loginResp apiResponse[loginResponseData]
	decodeJSON(t, login.Body.Bytes(), &loginResp)
	token := "Bearer " + loginResp.Data.AccessToken

	menus := env.doJSON(http.MethodGet, "/api/v1/auth/menus", token, nil)
	if menus.Code != http.StatusOK {
		t.Fatalf("expected menus status 200, got %d: %s", menus.Code, menus.Body.String())
	}

	var menuResp apiResponse[[]menuResponseData]
	decodeJSON(t, menus.Body.Bytes(), &menuResp)
	systemMenu := findMenuByCode(menuResp.Data, "system")
	if systemMenu == nil {
		t.Fatalf("expected system root menu in response, got %+v", menuResp.Data)
	}
	departmentMenu := findMenuByCode(systemMenu.Children, "system:department")
	if departmentMenu == nil || departmentMenu.Path != "/system/departments" || departmentMenu.Component != "system/DepartmentView" {
		t.Fatalf("expected department menu to be visible, got %+v", systemMenu.Children)
	}
	postMenu := findMenuByCode(systemMenu.Children, "system:post")
	if postMenu == nil || postMenu.Path != "/system/posts" || postMenu.Component != "system/PostView" {
		t.Fatalf("expected post menu to be visible, got %+v", systemMenu.Children)
	}

	listDepartments := env.doJSON(http.MethodGet, "/api/v1/system/departments", token, nil)
	if listDepartments.Code != http.StatusOK {
		t.Fatalf("expected list departments status 200, got %d: %s", listDepartments.Code, listDepartments.Body.String())
	}

	var departmentListResp apiResponse[[]departmentResponseData]
	decodeJSON(t, listDepartments.Body.Bytes(), &departmentListResp)
	if len(departmentListResp.Data) != 1 || departmentListResp.Data[0].Code != "rd-center" {
		t.Fatalf("unexpected department list response: %+v", departmentListResp.Data)
	}

	listPosts := env.doJSON(http.MethodGet, "/api/v1/system/posts", token, nil)
	if listPosts.Code != http.StatusOK {
		t.Fatalf("expected list posts status 200, got %d: %s", listPosts.Code, listPosts.Body.String())
	}

	var postListResp apiResponse[[]postResponseData]
	decodeJSON(t, listPosts.Body.Bytes(), &postListResp)
	if len(postListResp.Data) != 1 || postListResp.Data[0].Code != "team-lead" {
		t.Fatalf("unexpected post list response: %+v", postListResp.Data)
	}
}

func TestAuthMeReturnsActorSummary(t *testing.T) {
	var actor model.User
	env := newTestEnv(t, func(db *gorm.DB) {
		mustCreateRole(t, db, model.Role{
			ID:        2,
			Code:      "dept_admin",
			Name:      "部门管理员",
			Sort:      2,
			DataScope: "dept",
			Status:    model.RoleStatusEnabled,
		})
		actor = mustCreateUser(t, db, seededUser{
			Username:     "alice",
			Password:     "Alice12345",
			Nickname:     "Alice",
			DepartmentID: 10,
			Status:       model.UserStatusEnabled,
		})
		mustCreateUserRole(t, db, actor.ID, 2)
	})

	actorToken := env.mustIssueToken(actor.ID, actor.Username)
	resp := env.doJSON(http.MethodGet, "/api/v1/auth/me", actorToken, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("expected /auth/me status 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var meResp apiResponse[meResponseData]
	decodeJSON(t, resp.Body.Bytes(), &meResp)
	if meResp.Data.Username != "alice" {
		t.Fatalf("expected username alice, got %+v", meResp.Data)
	}
	if len(meResp.Data.RoleCodes) != 1 || meResp.Data.RoleCodes[0] != "dept_admin" {
		t.Fatalf("expected dept_admin role code, got %+v", meResp.Data.RoleCodes)
	}
	if !meResp.Data.DataScope.IncludeDepartment || meResp.Data.DataScope.AllowAll {
		t.Fatalf("expected department-only scope summary, got %+v", meResp.Data.DataScope)
	}
}

func TestSystemUsersRBACDenied(t *testing.T) {
	var actor model.User
	env := newTestEnv(t, func(db *gorm.DB) {
		mustCreateRole(t, db, model.Role{
			ID:        3,
			Code:      "auditor",
			Name:      "审计员",
			Sort:      3,
			DataScope: "all",
			Status:    model.RoleStatusEnabled,
		})
		actor = mustCreateUser(t, db, seededUser{
			Username:     "bob",
			Password:     "Bob123456",
			Nickname:     "Bob",
			DepartmentID: 20,
			Status:       model.UserStatusEnabled,
		})
		mustCreateUserRole(t, db, actor.ID, 3)
	})

	token := env.mustIssueToken(actor.ID, actor.Username)
	resp := env.doJSON(http.MethodGet, "/api/v1/system/users?page=1&page_size=10", token, nil)
	if resp.Code != http.StatusForbidden {
		t.Fatalf("expected RBAC denial status 403, got %d: %s", resp.Code, resp.Body.String())
	}

	var denied apiResponse[map[string]any]
	decodeJSON(t, resp.Body.Bytes(), &denied)
	if denied.Code != 40300 || denied.Message != "没有权限访问" {
		t.Fatalf("unexpected RBAC denial body: %+v", denied)
	}
}

func TestSystemUsersDataScopeFiltersByDepartment(t *testing.T) {
	var actor model.User
	env := newTestEnv(t, func(db *gorm.DB) {
		mustCreateRole(t, db, model.Role{
			ID:        4,
			Code:      "dept_manager",
			Name:      "部门经理",
			Sort:      4,
			DataScope: "dept",
			Status:    model.RoleStatusEnabled,
		})
		mustCreatePolicy(t, db, "dept_manager", "/api/v1/system/users", http.MethodGet)

		actor = mustCreateUser(t, db, seededUser{
			Username:     "carol",
			Password:     "Carol12345",
			Nickname:     "Carol",
			DepartmentID: 10,
			Status:       model.UserStatusEnabled,
		})
		sameDept := mustCreateUser(t, db, seededUser{
			Username:     "dave",
			Password:     "Dave123456",
			Nickname:     "Dave",
			DepartmentID: 10,
			Status:       model.UserStatusEnabled,
		})
		_ = sameDept
		mustCreateUser(t, db, seededUser{
			Username:     "erin",
			Password:     "Erin123456",
			Nickname:     "Erin",
			DepartmentID: 20,
			Status:       model.UserStatusEnabled,
		})
		mustCreateUserRole(t, db, actor.ID, 4)
	})

	token := env.mustIssueToken(actor.ID, actor.Username)
	resp := env.doJSON(http.MethodGet, "/api/v1/system/users?page=1&page_size=10", token, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("expected scoped list status 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var listResp apiResponse[userListResponseData]
	decodeJSON(t, resp.Body.Bytes(), &listResp)
	if listResp.Data.Total != 2 {
		t.Fatalf("expected 2 visible users in same department, got %+v", listResp.Data)
	}

	usernames := make(map[string]struct{}, len(listResp.Data.Items))
	for _, item := range listResp.Data.Items {
		usernames[item.Username] = struct{}{}
	}

	if _, ok := usernames["carol"]; !ok {
		t.Fatalf("expected actor carol to be visible, got %+v", listResp.Data.Items)
	}
	if _, ok := usernames["dave"]; !ok {
		t.Fatalf("expected same department user dave to be visible, got %+v", listResp.Data.Items)
	}
	if _, ok := usernames["erin"]; ok {
		t.Fatalf("did not expect other department user erin to be visible, got %+v", listResp.Data.Items)
	}
}

func TestSystemDepartmentAndPostHappyPath(t *testing.T) {
	var actor model.User
	env := newTestEnv(t, func(db *gorm.DB) {
		mustCreateRole(t, db, model.Role{
			ID:        12,
			Code:      "org_admin",
			Name:      "组织管理员",
			Sort:      12,
			DataScope: "all",
			Status:    model.RoleStatusEnabled,
		})
		mustCreatePolicy(t, db, "org_admin", "/api/v1/system/departments", http.MethodGet)
		mustCreatePolicy(t, db, "org_admin", "/api/v1/system/departments", http.MethodPost)
		mustCreatePolicy(t, db, "org_admin", "/api/v1/system/posts", http.MethodGet)
		mustCreatePolicy(t, db, "org_admin", "/api/v1/system/posts", http.MethodPost)

		actor = mustCreateUser(t, db, seededUser{
			Username:     "org-admin",
			Password:     "OrgAdmin123",
			Nickname:     "OrgAdmin",
			DepartmentID: 10,
			Status:       model.UserStatusEnabled,
		})
		mustCreateUserRole(t, db, actor.ID, 12)
	})

	token := env.mustIssueToken(actor.ID, actor.Username)

	createDepartment := env.doJSON(http.MethodPost, "/api/v1/system/departments", token, map[string]any{
		"parent_id":      0,
		"name":           "研发中心",
		"code":           "rd-center",
		"leader_user_id": actor.ID,
		"sort":           10,
		"status":         1,
		"remark":         "测试部门",
	})
	if createDepartment.Code != http.StatusOK {
		t.Fatalf("expected create department status 200, got %d: %s", createDepartment.Code, createDepartment.Body.String())
	}

	var createDepartmentResp apiResponse[departmentResponseData]
	decodeJSON(t, createDepartment.Body.Bytes(), &createDepartmentResp)
	if createDepartmentResp.Data.Code != "rd-center" {
		t.Fatalf("unexpected department create response: %+v", createDepartmentResp.Data)
	}

	createPost := env.doJSON(http.MethodPost, "/api/v1/system/posts", token, map[string]any{
		"code":   "team-lead",
		"name":   "团队负责人",
		"sort":   5,
		"status": 1,
		"remark": "测试岗位",
	})
	if createPost.Code != http.StatusOK {
		t.Fatalf("expected create post status 200, got %d: %s", createPost.Code, createPost.Body.String())
	}

	var createPostResp apiResponse[postResponseData]
	decodeJSON(t, createPost.Body.Bytes(), &createPostResp)
	if createPostResp.Data.Code != "team-lead" {
		t.Fatalf("unexpected post create response: %+v", createPostResp.Data)
	}

	listDepartments := env.doJSON(http.MethodGet, "/api/v1/system/departments", token, nil)
	if listDepartments.Code != http.StatusOK {
		t.Fatalf("expected list departments status 200, got %d: %s", listDepartments.Code, listDepartments.Body.String())
	}

	var departmentListResp apiResponse[[]departmentResponseData]
	decodeJSON(t, listDepartments.Body.Bytes(), &departmentListResp)
	if len(departmentListResp.Data) != 1 || departmentListResp.Data[0].Code != "rd-center" {
		t.Fatalf("unexpected department list response: %+v", departmentListResp.Data)
	}

	listPosts := env.doJSON(http.MethodGet, "/api/v1/system/posts", token, nil)
	if listPosts.Code != http.StatusOK {
		t.Fatalf("expected list posts status 200, got %d: %s", listPosts.Code, listPosts.Body.String())
	}

	var postListResp apiResponse[[]postResponseData]
	decodeJSON(t, listPosts.Body.Bytes(), &postListResp)
	if len(postListResp.Data) != 1 || postListResp.Data[0].Code != "team-lead" {
		t.Fatalf("unexpected post list response: %+v", postListResp.Data)
	}
}

func TestSystemDictHappyPath(t *testing.T) {
	var actor model.User
	env := newTestEnv(t, func(db *gorm.DB) {
		mustCreateRole(t, db, model.Role{
			ID:        5,
			Code:      "dict_admin",
			Name:      "字典管理员",
			Sort:      5,
			DataScope: "all",
			Status:    model.RoleStatusEnabled,
		})
		mustCreatePolicy(t, db, "dict_admin", "/api/v1/system/dict-types", http.MethodGet)
		mustCreatePolicy(t, db, "dict_admin", "/api/v1/system/dict-types", http.MethodPost)
		mustCreatePolicy(t, db, "dict_admin", "/api/v1/system/dict-items", http.MethodGet)
		mustCreatePolicy(t, db, "dict_admin", "/api/v1/system/dict-items", http.MethodPost)

		actor = mustCreateUser(t, db, seededUser{
			Username:     "dict-admin",
			Password:     "Dict123456",
			Nickname:     "DictAdmin",
			DepartmentID: 10,
			Status:       model.UserStatusEnabled,
		})
		mustCreateUserRole(t, db, actor.ID, 5)
	})

	token := env.mustIssueToken(actor.ID, actor.Username)

	createType := env.doJSON(http.MethodPost, "/api/v1/system/dict-types", token, map[string]any{
		"code":   "common:priority",
		"name":   "优先级",
		"sort":   10,
		"status": 1,
		"remark": "测试字典类型",
	})
	if createType.Code != http.StatusOK {
		t.Fatalf("expected create dict type status 200, got %d: %s", createType.Code, createType.Body.String())
	}

	var createTypeResp apiResponse[dictTypeResponseData]
	decodeJSON(t, createType.Body.Bytes(), &createTypeResp)
	if createTypeResp.Data.Code != "common:priority" {
		t.Fatalf("unexpected dict type create response: %+v", createTypeResp.Data)
	}

	createItem := env.doJSON(http.MethodPost, "/api/v1/system/dict-items", token, map[string]any{
		"type_id":  createTypeResp.Data.ID,
		"item_key": "high",
		"label":    "高优先级",
		"value":    "high",
		"tag_type": "error",
		"sort":     10,
		"status":   1,
		"remark":   "测试字典项",
	})
	if createItem.Code != http.StatusOK {
		t.Fatalf("expected create dict item status 200, got %d: %s", createItem.Code, createItem.Body.String())
	}

	var createItemResp apiResponse[dictItemResponseData]
	decodeJSON(t, createItem.Body.Bytes(), &createItemResp)
	if createItemResp.Data.TypeID != createTypeResp.Data.ID || createItemResp.Data.ItemKey != "high" {
		t.Fatalf("unexpected dict item create response: %+v", createItemResp.Data)
	}

	listTypes := env.doJSON(http.MethodGet, "/api/v1/system/dict-types?page=1&page_size=10", token, nil)
	if listTypes.Code != http.StatusOK {
		t.Fatalf("expected list dict types status 200, got %d: %s", listTypes.Code, listTypes.Body.String())
	}

	var listTypesResp apiResponse[dictTypeListResponseData]
	decodeJSON(t, listTypes.Body.Bytes(), &listTypesResp)
	if listTypesResp.Data.Total != 1 || len(listTypesResp.Data.Items) != 1 {
		t.Fatalf("unexpected dict type list response: %+v", listTypesResp.Data)
	}

	listItems := env.doJSON(
		http.MethodGet,
		fmt.Sprintf("/api/v1/system/dict-items?page=1&page_size=10&type_id=%d", createTypeResp.Data.ID),
		token,
		nil,
	)
	if listItems.Code != http.StatusOK {
		t.Fatalf("expected list dict items status 200, got %d: %s", listItems.Code, listItems.Body.String())
	}

	var listItemsResp apiResponse[dictItemListResponseData]
	decodeJSON(t, listItems.Body.Bytes(), &listItemsResp)
	if listItemsResp.Data.Total != 1 || len(listItemsResp.Data.Items) != 1 {
		t.Fatalf("unexpected dict item list response: %+v", listItemsResp.Data)
	}
	if listItemsResp.Data.Items[0].ItemKey != "high" {
		t.Fatalf("unexpected dict item list item: %+v", listItemsResp.Data.Items[0])
	}
}

func TestAuthAccountProfileAndPasswordUpdate(t *testing.T) {
	var actor model.User
	env := newTestEnv(t, func(db *gorm.DB) {
		mustCreateDepartment(t, db, model.Department{
			ID:       10,
			ParentID: 0,
			Name:     "研发中心",
			Code:     "rd-center",
			Sort:     10,
			Status:   model.DepartmentStatusEnabled,
		})
		mustCreateRole(t, db, model.Role{
			ID:        6,
			Code:      "account_user",
			Name:      "账户用户",
			Sort:      6,
			DataScope: "self",
			Status:    model.RoleStatusEnabled,
		})
		actor = mustCreateUser(t, db, seededUser{
			Username:     "frank",
			Password:     "Frank12345",
			Nickname:     "Frank",
			DepartmentID: 10,
			Status:       model.UserStatusEnabled,
		})
		mustCreateUserRole(t, db, actor.ID, 6)
	})

	token := env.mustIssueToken(actor.ID, actor.Username)

	profile := env.doJSON(http.MethodGet, "/api/v1/auth/account", token, nil)
	if profile.Code != http.StatusOK {
		t.Fatalf("expected account profile status 200, got %d: %s", profile.Code, profile.Body.String())
	}

	var profileResp apiResponse[accountProfileResponseData]
	decodeJSON(t, profile.Body.Bytes(), &profileResp)
	if profileResp.Data.Username != "frank" || profileResp.Data.DepartmentName != "研发中心" {
		t.Fatalf("unexpected account profile response: %+v", profileResp.Data)
	}

	updateProfile := env.doJSON(http.MethodPost, "/api/v1/auth/account/profile", token, map[string]any{
		"nickname": "Franky",
	})
	if updateProfile.Code != http.StatusOK {
		t.Fatalf("expected update account profile status 200, got %d: %s", updateProfile.Code, updateProfile.Body.String())
	}

	var updateProfileResp apiResponse[accountProfileResponseData]
	decodeJSON(t, updateProfile.Body.Bytes(), &updateProfileResp)
	if updateProfileResp.Data.Nickname != "Franky" {
		t.Fatalf("expected updated nickname Franky, got %+v", updateProfileResp.Data)
	}

	updatePassword := env.doJSON(http.MethodPost, "/api/v1/auth/account/password", token, map[string]any{
		"old_password": "Frank12345",
		"new_password": "Frank54321",
	})
	if updatePassword.Code != http.StatusOK {
		t.Fatalf("expected update account password status 200, got %d: %s", updatePassword.Code, updatePassword.Body.String())
	}

	loginOldPassword := env.doJSON(http.MethodPost, "/api/v1/auth/login", "", map[string]any{
		"username": "frank",
		"password": "Frank12345",
	})
	if loginOldPassword.Code != http.StatusUnauthorized {
		t.Fatalf("expected old password login status 401, got %d: %s", loginOldPassword.Code, loginOldPassword.Body.String())
	}

	loginNewPassword := env.doJSON(http.MethodPost, "/api/v1/auth/login", "", map[string]any{
		"username": "frank",
		"password": "Frank54321",
	})
	if loginNewPassword.Code != http.StatusOK {
		t.Fatalf("expected new password login status 200, got %d: %s", loginNewPassword.Code, loginNewPassword.Body.String())
	}

	var loginNewResp apiResponse[loginResponseData]
	decodeJSON(t, loginNewPassword.Body.Bytes(), &loginNewResp)
	if loginNewResp.Data.Nickname != "Franky" {
		t.Fatalf("expected updated nickname in login response, got %+v", loginNewResp.Data)
	}
}

func TestSystemAttachmentHappyPath(t *testing.T) {
	var actor model.User
	env := newTestEnv(t, func(db *gorm.DB) {
		mustCreateRole(t, db, model.Role{
			ID:        7,
			Code:      "attachment_admin",
			Name:      "附件管理员",
			Sort:      7,
			DataScope: "all",
			Status:    model.RoleStatusEnabled,
		})
		mustCreatePolicy(t, db, "attachment_admin", "/api/v1/system/attachments", http.MethodGet)
		mustCreatePolicy(t, db, "attachment_admin", "/api/v1/system/attachments", http.MethodPost)
		mustCreatePolicy(t, db, "attachment_admin", "/api/v1/system/attachments/:id/update", http.MethodPost)
		mustCreatePolicy(t, db, "attachment_admin", "/api/v1/system/attachments/:id/status", http.MethodPost)

		actor = mustCreateUser(t, db, seededUser{
			Username:     "attach-admin",
			Password:     "Attach12345",
			Nickname:     "AttachAdmin",
			DepartmentID: 10,
			Status:       model.UserStatusEnabled,
		})
		mustCreateUserRole(t, db, actor.ID, 7)
	})

	token := env.mustIssueToken(actor.ID, actor.Username)

	createResp := env.doMultipart(
		http.MethodPost,
		"/api/v1/system/attachments",
		token,
		map[string]string{
			"display_name": "合同模板",
			"category":     "contract",
			"biz_type":     "system-template",
			"status":       "1",
			"remark":       "测试附件上传",
		},
		"file",
		"contract-template.pdf",
		[]byte("attachment-center-test-content"),
	)
	if createResp.Code != http.StatusOK {
		t.Fatalf("expected create attachment status 200, got %d: %s", createResp.Code, createResp.Body.String())
	}

	var createBody apiResponse[attachmentResponseData]
	decodeJSON(t, createResp.Body.Bytes(), &createBody)
	if createBody.Data.DisplayName != "合同模板" || createBody.Data.Category != "contract" {
		t.Fatalf("unexpected attachment create response: %+v", createBody.Data)
	}

	listResp := env.doJSON(http.MethodGet, "/api/v1/system/attachments?page=1&page_size=10", token, nil)
	if listResp.Code != http.StatusOK {
		t.Fatalf("expected attachment list status 200, got %d: %s", listResp.Code, listResp.Body.String())
	}

	var listBody apiResponse[attachmentListResponseData]
	decodeJSON(t, listResp.Body.Bytes(), &listBody)
	if listBody.Data.Total != 1 || len(listBody.Data.Items) != 1 {
		t.Fatalf("unexpected attachment list response: %+v", listBody.Data)
	}

	updateResp := env.doJSON(
		http.MethodPost,
		fmt.Sprintf("/api/v1/system/attachments/%d/update", createBody.Data.ID),
		token,
		map[string]any{
			"display_name": "正式合同模板",
			"category":     "template",
			"biz_type":     "contract",
			"status":       1,
			"remark":       "更新后的备注",
		},
	)
	if updateResp.Code != http.StatusOK {
		t.Fatalf("expected update attachment status 200, got %d: %s", updateResp.Code, updateResp.Body.String())
	}

	var updateBody apiResponse[attachmentResponseData]
	decodeJSON(t, updateResp.Body.Bytes(), &updateBody)
	if updateBody.Data.DisplayName != "正式合同模板" || updateBody.Data.BizType != "contract" {
		t.Fatalf("unexpected attachment update response: %+v", updateBody.Data)
	}

	statusResp := env.doJSON(
		http.MethodPost,
		fmt.Sprintf("/api/v1/system/attachments/%d/status", createBody.Data.ID),
		token,
		map[string]any{"status": 2},
	)
	if statusResp.Code != http.StatusOK {
		t.Fatalf("expected update attachment status endpoint 200, got %d: %s", statusResp.Code, statusResp.Body.String())
	}
}

func newTestEnv(t *testing.T, seed func(db *gorm.DB)) *testEnv {
	t.Helper()
	gin.SetMode(gin.TestMode)

	dsn := fmt.Sprintf(
		"file:%s?mode=memory&cache=shared",
		strings.NewReplacer("/", "_", " ", "_").Replace(t.Name()),
	)
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}

	migrateTestSchema(t, db)
	if seed != nil {
		seed(db)
	}

	cfg := &platformConfig.Config{
		App: platformConfig.AppConfig{
			Name: "ez-admin-test",
			Env:  "test",
		},
		Auth: platformConfig.AuthConfig{
			JWTSecret:      strings.Repeat("t", 32),
			AccessTokenTTL: 3600,
			Issuer:         "ez-admin-test",
		},
		Upload: platformConfig.UploadConfig{
			Dir:        t.TempDir(),
			PublicPath: "/uploads",
			MaxSizeMB:  10,
		},
	}

	tokenManager, err := authnPlatform.NewManager(cfg.Auth)
	if err != nil {
		t.Fatalf("create token manager: %v", err)
	}

	enforcer, err := authzPlatform.NewEnforcer(db, filepath.Join("..", "..", "configs", "rbac_model.conf"))
	if err != nil {
		t.Fatalf("create permission enforcer: %v", err)
	}

	router := NewRouter(RouterOptions{
		Config:     cfg,
		Log:        zap.NewNop(),
		DB:         db,
		Token:      tokenManager,
		Permission: enforcer,
	})

	return &testEnv{
		t:      t,
		db:     db,
		router: router,
		token:  tokenManager,
	}
}

func migrateTestSchema(t *testing.T, db *gorm.DB) {
	t.Helper()
	err := db.AutoMigrate(
		&model.CasbinRule{},
		&model.Department{},
		&model.SystemDictItem{},
		&model.SystemDictType{},
		&model.SystemAttachment{},
		&model.LoginLog{},
		&model.Menu{},
		&model.Notice{},
		&model.OperationLog{},
		&model.Post{},
		&model.Role{},
		&model.RoleDataScope{},
		&model.RoleMenu{},
		&model.SystemConfig{},
		&model.SystemFile{},
		&model.User{},
		&model.UserPost{},
		&model.UserRole{},
	)
	if err != nil {
		t.Fatalf("auto migrate test schema: %v", err)
	}
}

func (e *testEnv) mustIssueToken(userID uint, username string) string {
	e.t.Helper()
	token, _, err := e.token.GenerateAccessToken(userID, username)
	if err != nil {
		e.t.Fatalf("generate access token: %v", err)
	}
	return "Bearer " + token
}

func (e *testEnv) doJSON(method string, path string, bearerToken string, body any) *httptest.ResponseRecorder {
	e.t.Helper()

	var reader *bytes.Reader
	if body == nil {
		reader = bytes.NewReader(nil)
	} else {
		raw, err := json.Marshal(body)
		if err != nil {
			e.t.Fatalf("marshal request body: %v", err)
		}
		reader = bytes.NewReader(raw)
	}

	req := httptest.NewRequest(method, path, reader)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if bearerToken != "" {
		req.Header.Set("Authorization", bearerToken)
	}

	recorder := httptest.NewRecorder()
	e.router.ServeHTTP(recorder, req)
	return recorder
}

func (e *testEnv) doMultipart(
	method string,
	path string,
	bearerToken string,
	fields map[string]string,
	fileField string,
	fileName string,
	fileContent []byte,
) *httptest.ResponseRecorder {
	e.t.Helper()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	for key, value := range fields {
		if err := writer.WriteField(key, value); err != nil {
			e.t.Fatalf("write multipart field %s: %v", key, err)
		}
	}

	part, err := writer.CreateFormFile(fileField, fileName)
	if err != nil {
		e.t.Fatalf("create multipart file %s: %v", fileName, err)
	}
	if _, err := io.Copy(part, bytes.NewReader(fileContent)); err != nil {
		e.t.Fatalf("write multipart file %s: %v", fileName, err)
	}
	if err := writer.Close(); err != nil {
		e.t.Fatalf("close multipart writer: %v", err)
	}

	req := httptest.NewRequest(method, path, &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if bearerToken != "" {
		req.Header.Set("Authorization", bearerToken)
	}

	recorder := httptest.NewRecorder()
	e.router.ServeHTTP(recorder, req)
	return recorder
}

type seededUser struct {
	Username     string
	Password     string
	Nickname     string
	DepartmentID uint
	Status       model.UserStatus
}

func mustCreateUser(t *testing.T, db *gorm.DB, input seededUser) model.User {
	t.Helper()
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}

	user := model.User{
		Username:     input.Username,
		PasswordHash: string(hash),
		Nickname:     input.Nickname,
		DepartmentID: input.DepartmentID,
		Status:       input.Status,
	}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("create user %s: %v", input.Username, err)
	}
	return user
}

func mustCreateRole(t *testing.T, db *gorm.DB, role model.Role) {
	t.Helper()
	if err := db.Create(&role).Error; err != nil {
		t.Fatalf("create role %s: %v", role.Code, err)
	}
}

func mustCreateDepartment(t *testing.T, db *gorm.DB, department model.Department) {
	t.Helper()
	if err := db.Create(&department).Error; err != nil {
		t.Fatalf("create department %s: %v", department.Code, err)
	}
}

func mustCreateMenu(t *testing.T, db *gorm.DB, menu model.Menu) {
	t.Helper()
	if err := db.Create(&menu).Error; err != nil {
		t.Fatalf("create menu %s: %v", menu.Code, err)
	}
}

func mustCreateUserRole(t *testing.T, db *gorm.DB, userID uint, roleID uint) {
	t.Helper()
	if err := db.Create(&model.UserRole{UserID: userID, RoleID: roleID}).Error; err != nil {
		t.Fatalf("create user role binding user=%d role=%d: %v", userID, roleID, err)
	}
}

func mustCreateRoleMenu(t *testing.T, db *gorm.DB, roleID uint, menuID uint) {
	t.Helper()
	if err := db.Create(&model.RoleMenu{RoleID: roleID, MenuID: menuID}).Error; err != nil {
		t.Fatalf("create role menu binding role=%d menu=%d: %v", roleID, menuID, err)
	}
}

func mustCreatePolicy(t *testing.T, db *gorm.DB, roleCode string, path string, method string) {
	t.Helper()
	row := model.CasbinRule{
		Ptype: "p",
		V0:    roleCode,
		V1:    path,
		V2:    method,
	}
	if err := db.Create(&row).Error; err != nil {
		t.Fatalf("create casbin policy %s %s %s: %v", roleCode, path, method, err)
	}
}

func decodeJSON(t *testing.T, raw []byte, target any) {
	t.Helper()
	if err := json.Unmarshal(raw, target); err != nil {
		t.Fatalf("decode json: %v\nbody=%s", err, string(raw))
	}
}

func findMenuByCode(menus []menuResponseData, code string) *menuResponseData {
	for i := range menus {
		if menus[i].Code == code {
			return &menus[i]
		}
	}
	return nil
}
