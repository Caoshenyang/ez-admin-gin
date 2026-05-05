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

	"ez-admin-gin/server/internal/config"
	"ez-admin-gin/server/internal/model"
	authnPlatform "ez-admin-gin/server/internal/platform/authn"
	authzPlatform "ez-admin-gin/server/internal/platform/authz"

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

type customerResponseData struct {
	ID             uint   `json:"id"`
	Name           string `json:"name"`
	ContactName    string `json:"contact_name"`
	Phone          string `json:"phone"`
	Level          string `json:"level"`
	Source         string `json:"source"`
	DepartmentID   uint   `json:"department_id"`
	DepartmentName string `json:"department_name"`
	OwnerUserID    uint   `json:"owner_user_id"`
	OwnerUsername  string `json:"owner_username"`
	OwnerNickname  string `json:"owner_nickname"`
	Status         int    `json:"status"`
	Remark         string `json:"remark"`
}

type customerListResponseData struct {
	Items []customerResponseData `json:"items"`
	Total int64                  `json:"total"`
}

type customerFollowUpResponseData struct {
	ID             uint    `json:"id"`
	CustomerID     uint    `json:"customer_id"`
	CustomerName   string  `json:"customer_name"`
	DepartmentID   uint    `json:"department_id"`
	DepartmentName string  `json:"department_name"`
	OwnerUserID    uint    `json:"owner_user_id"`
	OwnerUsername  string  `json:"owner_username"`
	OwnerNickname  string  `json:"owner_nickname"`
	FollowType     string  `json:"follow_type"`
	Subject        string  `json:"subject"`
	Content        string  `json:"content"`
	Result         string  `json:"result"`
	NextFollowAt   *string `json:"next_follow_at"`
	Status         int     `json:"status"`
}

type customerFollowUpListResponseData struct {
	Items []customerFollowUpResponseData `json:"items"`
	Total int64                          `json:"total"`
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

func TestCRMCustomerHappyPath(t *testing.T) {
	var actor model.User
	env := newTestEnv(t, func(db *gorm.DB) {
		mustCreateDepartment(t, db, model.Department{
			ID:       10,
			ParentID: 0,
			Name:     "销售一部",
			Code:     "sales-a",
			Sort:     10,
			Status:   model.DepartmentStatusEnabled,
		})
		mustCreateRole(t, db, model.Role{
			ID:        8,
			Code:      "customer_admin",
			Name:      "客户管理员",
			Sort:      8,
			DataScope: "all",
			Status:    model.RoleStatusEnabled,
		})
		mustCreatePolicy(t, db, "customer_admin", "/api/v1/crm/customers", http.MethodGet)
		mustCreatePolicy(t, db, "customer_admin", "/api/v1/crm/customers", http.MethodPost)
		mustCreatePolicy(t, db, "customer_admin", "/api/v1/crm/customers/:id/update", http.MethodPost)
		mustCreatePolicy(t, db, "customer_admin", "/api/v1/crm/customers/:id/status", http.MethodPost)

		actor = mustCreateUser(t, db, seededUser{
			Username:     "customer-admin",
			Password:     "Customer12345",
			Nickname:     "CustomerAdmin",
			DepartmentID: 10,
			Status:       model.UserStatusEnabled,
		})
		mustCreateUserRole(t, db, actor.ID, 8)
	})

	token := env.mustIssueToken(actor.ID, actor.Username)

	createResp := env.doJSON(http.MethodPost, "/api/v1/crm/customers", token, map[string]any{
		"name":         "星河工业设备有限公司",
		"contact_name": "陈经理",
		"phone":        "13800138000",
		"level":        "vip",
		"source":       "referral",
		"status":       1,
		"remark":       "重点客户",
	})
	if createResp.Code != http.StatusOK {
		t.Fatalf("expected create customer status 200, got %d: %s", createResp.Code, createResp.Body.String())
	}

	var createBody apiResponse[customerResponseData]
	decodeJSON(t, createResp.Body.Bytes(), &createBody)
	if createBody.Data.Name != "星河工业设备有限公司" || createBody.Data.OwnerUserID != actor.ID {
		t.Fatalf("unexpected customer create response: %+v", createBody.Data)
	}

	listResp := env.doJSON(http.MethodGet, "/api/v1/crm/customers?page=1&page_size=10", token, nil)
	if listResp.Code != http.StatusOK {
		t.Fatalf("expected customer list status 200, got %d: %s", listResp.Code, listResp.Body.String())
	}

	var listBody apiResponse[customerListResponseData]
	decodeJSON(t, listResp.Body.Bytes(), &listBody)
	if listBody.Data.Total != 1 || len(listBody.Data.Items) != 1 {
		t.Fatalf("unexpected customer list response: %+v", listBody.Data)
	}

	updateResp := env.doJSON(
		http.MethodPost,
		fmt.Sprintf("/api/v1/crm/customers/%d/update", createBody.Data.ID),
		token,
		map[string]any{
			"name":         "星河工业设备集团",
			"contact_name": "李总监",
			"phone":        "13900139000",
			"level":        "a",
			"source":       "offline",
			"status":       1,
			"remark":       "已进入方案评审",
		},
	)
	if updateResp.Code != http.StatusOK {
		t.Fatalf("expected update customer status 200, got %d: %s", updateResp.Code, updateResp.Body.String())
	}

	var updateBody apiResponse[customerResponseData]
	decodeJSON(t, updateResp.Body.Bytes(), &updateBody)
	if updateBody.Data.Name != "星河工业设备集团" || updateBody.Data.Level != "a" {
		t.Fatalf("unexpected customer update response: %+v", updateBody.Data)
	}

	statusResp := env.doJSON(
		http.MethodPost,
		fmt.Sprintf("/api/v1/crm/customers/%d/status", createBody.Data.ID),
		token,
		map[string]any{"status": 2},
	)
	if statusResp.Code != http.StatusOK {
		t.Fatalf("expected customer status update 200, got %d: %s", statusResp.Code, statusResp.Body.String())
	}
}

func TestCRMCustomerDataScopeFiltersByDepartment(t *testing.T) {
	var actor model.User
	env := newTestEnv(t, func(db *gorm.DB) {
		mustCreateDepartment(t, db, model.Department{
			ID:       10,
			ParentID: 0,
			Name:     "销售一部",
			Code:     "sales-a",
			Sort:     10,
			Status:   model.DepartmentStatusEnabled,
		})
		mustCreateDepartment(t, db, model.Department{
			ID:       20,
			ParentID: 0,
			Name:     "销售二部",
			Code:     "sales-b",
			Sort:     20,
			Status:   model.DepartmentStatusEnabled,
		})
		mustCreateRole(t, db, model.Role{
			ID:        9,
			Code:      "customer_manager",
			Name:      "客户经理",
			Sort:      9,
			DataScope: "dept",
			Status:    model.RoleStatusEnabled,
		})
		mustCreatePolicy(t, db, "customer_manager", "/api/v1/crm/customers", http.MethodGet)

		actor = mustCreateUser(t, db, seededUser{
			Username:     "iris",
			Password:     "Iris123456",
			Nickname:     "Iris",
			DepartmentID: 10,
			Status:       model.UserStatusEnabled,
		})
		sameDeptOwner := mustCreateUser(t, db, seededUser{
			Username:     "jack",
			Password:     "Jack123456",
			Nickname:     "Jack",
			DepartmentID: 10,
			Status:       model.UserStatusEnabled,
		})
		otherDeptOwner := mustCreateUser(t, db, seededUser{
			Username:     "kate",
			Password:     "Kate123456",
			Nickname:     "Kate",
			DepartmentID: 20,
			Status:       model.UserStatusEnabled,
		})
		mustCreateUserRole(t, db, actor.ID, 9)

		if err := db.Create(&model.Customer{
			Name:         "同部门客户",
			DepartmentID: 10,
			OwnerUserID:  sameDeptOwner.ID,
			Status:       model.CustomerStatusEnabled,
		}).Error; err != nil {
			t.Fatalf("create same department customer: %v", err)
		}
		if err := db.Create(&model.Customer{
			Name:         "跨部门客户",
			DepartmentID: 20,
			OwnerUserID:  otherDeptOwner.ID,
			Status:       model.CustomerStatusEnabled,
		}).Error; err != nil {
			t.Fatalf("create other department customer: %v", err)
		}
		if err := db.Create(&model.Customer{
			Name:         "当前负责人客户",
			DepartmentID: 10,
			OwnerUserID:  actor.ID,
			Status:       model.CustomerStatusEnabled,
		}).Error; err != nil {
			t.Fatalf("create actor customer: %v", err)
		}
	})

	token := env.mustIssueToken(actor.ID, actor.Username)
	resp := env.doJSON(http.MethodGet, "/api/v1/crm/customers?page=1&page_size=10", token, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("expected scoped customer list status 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var listResp apiResponse[customerListResponseData]
	decodeJSON(t, resp.Body.Bytes(), &listResp)
	if listResp.Data.Total != 2 {
		t.Fatalf("expected 2 visible customers in same department, got %+v", listResp.Data)
	}

	names := make(map[string]struct{}, len(listResp.Data.Items))
	for _, item := range listResp.Data.Items {
		names[item.Name] = struct{}{}
	}
	if _, ok := names["同部门客户"]; !ok {
		t.Fatalf("expected same department customer to be visible, got %+v", listResp.Data.Items)
	}
	if _, ok := names["当前负责人客户"]; !ok {
		t.Fatalf("expected actor customer to be visible, got %+v", listResp.Data.Items)
	}
	if _, ok := names["跨部门客户"]; ok {
		t.Fatalf("did not expect other department customer to be visible, got %+v", listResp.Data.Items)
	}
}

func TestCRMFollowUpHappyPath(t *testing.T) {
	var actor model.User
	env := newTestEnv(t, func(db *gorm.DB) {
		mustCreateDepartment(t, db, model.Department{
			ID:       10,
			ParentID: 0,
			Name:     "销售一部",
			Code:     "sales-a",
			Sort:     10,
			Status:   model.DepartmentStatusEnabled,
		})
		mustCreateRole(t, db, model.Role{
			ID:        10,
			Code:      "followup_admin",
			Name:      "跟进管理员",
			Sort:      10,
			DataScope: "all",
			Status:    model.RoleStatusEnabled,
		})
		mustCreatePolicy(t, db, "followup_admin", "/api/v1/crm/followups", http.MethodGet)
		mustCreatePolicy(t, db, "followup_admin", "/api/v1/crm/followups", http.MethodPost)
		mustCreatePolicy(t, db, "followup_admin", "/api/v1/crm/followups/:id/update", http.MethodPost)
		mustCreatePolicy(t, db, "followup_admin", "/api/v1/crm/followups/:id/status", http.MethodPost)

		actor = mustCreateUser(t, db, seededUser{
			Username:     "followup-admin",
			Password:     "Follow12345",
			Nickname:     "FollowUpAdmin",
			DepartmentID: 10,
			Status:       model.UserStatusEnabled,
		})
		mustCreateUserRole(t, db, actor.ID, 10)

		if err := db.Create(&model.Customer{
			ID:           1,
			Name:         "星河工业设备有限公司",
			ContactName:  "陈经理",
			Phone:        "13800138000",
			Level:        "vip",
			Source:       "referral",
			DepartmentID: 10,
			OwnerUserID:  actor.ID,
			Status:       model.CustomerStatusEnabled,
		}).Error; err != nil {
			t.Fatalf("create customer for followup: %v", err)
		}
	})

	token := env.mustIssueToken(actor.ID, actor.Username)
	nextFollowAt := "2026-05-10T09:30:00Z"

	createResp := env.doJSON(http.MethodPost, "/api/v1/crm/followups", token, map[string]any{
		"customer_id":    1,
		"follow_type":    "visit",
		"subject":        "首次上门沟通",
		"content":        "确认采购周期与预算窗口",
		"result":         "约定下周二继续方案评审",
		"next_follow_at": nextFollowAt,
		"status":         1,
	})
	if createResp.Code != http.StatusOK {
		t.Fatalf("expected create followup status 200, got %d: %s", createResp.Code, createResp.Body.String())
	}

	var createBody apiResponse[customerFollowUpResponseData]
	decodeJSON(t, createResp.Body.Bytes(), &createBody)
	if createBody.Data.CustomerName != "星河工业设备有限公司" || createBody.Data.FollowType != "visit" {
		t.Fatalf("unexpected followup create response: %+v", createBody.Data)
	}

	listResp := env.doJSON(http.MethodGet, "/api/v1/crm/followups?page=1&page_size=10", token, nil)
	if listResp.Code != http.StatusOK {
		t.Fatalf("expected followup list status 200, got %d: %s", listResp.Code, listResp.Body.String())
	}

	var listBody apiResponse[customerFollowUpListResponseData]
	decodeJSON(t, listResp.Body.Bytes(), &listBody)
	if listBody.Data.Total != 1 || len(listBody.Data.Items) != 1 {
		t.Fatalf("unexpected followup list response: %+v", listBody.Data)
	}

	updateResp := env.doJSON(
		http.MethodPost,
		fmt.Sprintf("/api/v1/crm/followups/%d/update", createBody.Data.ID),
		token,
		map[string]any{
			"follow_type":    "meeting",
			"subject":        "方案评审会",
			"content":        "确认实施范围和预算上限",
			"result":         "客户同意进入报价阶段",
			"next_follow_at": "2026-05-15T06:00:00Z",
			"status":         2,
		},
	)
	if updateResp.Code != http.StatusOK {
		t.Fatalf("expected update followup status 200, got %d: %s", updateResp.Code, updateResp.Body.String())
	}

	var updateBody apiResponse[customerFollowUpResponseData]
	decodeJSON(t, updateResp.Body.Bytes(), &updateBody)
	if updateBody.Data.Subject != "方案评审会" || updateBody.Data.Status != 2 {
		t.Fatalf("unexpected followup update response: %+v", updateBody.Data)
	}

	statusResp := env.doJSON(
		http.MethodPost,
		fmt.Sprintf("/api/v1/crm/followups/%d/status", createBody.Data.ID),
		token,
		map[string]any{"status": 3},
	)
	if statusResp.Code != http.StatusOK {
		t.Fatalf("expected followup status update 200, got %d: %s", statusResp.Code, statusResp.Body.String())
	}
}

func TestCRMFollowUpDataScopeFiltersByDepartment(t *testing.T) {
	var actor model.User
	env := newTestEnv(t, func(db *gorm.DB) {
		mustCreateDepartment(t, db, model.Department{
			ID:       10,
			ParentID: 0,
			Name:     "销售一部",
			Code:     "sales-a",
			Sort:     10,
			Status:   model.DepartmentStatusEnabled,
		})
		mustCreateDepartment(t, db, model.Department{
			ID:       20,
			ParentID: 0,
			Name:     "销售二部",
			Code:     "sales-b",
			Sort:     20,
			Status:   model.DepartmentStatusEnabled,
		})
		mustCreateRole(t, db, model.Role{
			ID:        11,
			Code:      "followup_manager",
			Name:      "跟进经理",
			Sort:      11,
			DataScope: "dept",
			Status:    model.RoleStatusEnabled,
		})
		mustCreatePolicy(t, db, "followup_manager", "/api/v1/crm/followups", http.MethodGet)

		actor = mustCreateUser(t, db, seededUser{
			Username:     "luna",
			Password:     "Luna123456",
			Nickname:     "Luna",
			DepartmentID: 10,
			Status:       model.UserStatusEnabled,
		})
		sameDeptOwner := mustCreateUser(t, db, seededUser{
			Username:     "mike",
			Password:     "Mike123456",
			Nickname:     "Mike",
			DepartmentID: 10,
			Status:       model.UserStatusEnabled,
		})
		otherDeptOwner := mustCreateUser(t, db, seededUser{
			Username:     "nina",
			Password:     "Nina123456",
			Nickname:     "Nina",
			DepartmentID: 20,
			Status:       model.UserStatusEnabled,
		})
		mustCreateUserRole(t, db, actor.ID, 11)

		if err := db.Create(&model.Customer{
			ID:           101,
			Name:         "同部门客户",
			DepartmentID: 10,
			OwnerUserID:  sameDeptOwner.ID,
			Status:       model.CustomerStatusEnabled,
		}).Error; err != nil {
			t.Fatalf("create same department customer for followup: %v", err)
		}
		if err := db.Create(&model.Customer{
			ID:           102,
			Name:         "跨部门客户",
			DepartmentID: 20,
			OwnerUserID:  otherDeptOwner.ID,
			Status:       model.CustomerStatusEnabled,
		}).Error; err != nil {
			t.Fatalf("create other department customer for followup: %v", err)
		}
		if err := db.Create(&model.Customer{
			ID:           103,
			Name:         "当前负责人客户",
			DepartmentID: 10,
			OwnerUserID:  actor.ID,
			Status:       model.CustomerStatusEnabled,
		}).Error; err != nil {
			t.Fatalf("create actor customer for followup: %v", err)
		}

		if err := db.Create(&model.CustomerFollowUp{
			CustomerID:   101,
			DepartmentID: 10,
			OwnerUserID:  sameDeptOwner.ID,
			FollowType:   "phone",
			Subject:      "同部门跟进",
			Content:      "同部门客户继续推进",
			Status:       model.CustomerFollowUpStatusPending,
		}).Error; err != nil {
			t.Fatalf("create same department followup: %v", err)
		}
		if err := db.Create(&model.CustomerFollowUp{
			CustomerID:   102,
			DepartmentID: 20,
			OwnerUserID:  otherDeptOwner.ID,
			FollowType:   "phone",
			Subject:      "跨部门跟进",
			Content:      "不应被当前部门看到",
			Status:       model.CustomerFollowUpStatusPending,
		}).Error; err != nil {
			t.Fatalf("create other department followup: %v", err)
		}
		if err := db.Create(&model.CustomerFollowUp{
			CustomerID:   103,
			DepartmentID: 10,
			OwnerUserID:  actor.ID,
			FollowType:   "visit",
			Subject:      "当前负责人跟进",
			Content:      "当前负责人自己的客户",
			Status:       model.CustomerFollowUpStatusPending,
		}).Error; err != nil {
			t.Fatalf("create actor followup: %v", err)
		}
	})

	token := env.mustIssueToken(actor.ID, actor.Username)
	resp := env.doJSON(http.MethodGet, "/api/v1/crm/followups?page=1&page_size=10", token, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("expected scoped followup list status 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var listResp apiResponse[customerFollowUpListResponseData]
	decodeJSON(t, resp.Body.Bytes(), &listResp)
	if listResp.Data.Total != 2 {
		t.Fatalf("expected 2 visible followups in same department, got %+v", listResp.Data)
	}

	subjects := make(map[string]struct{}, len(listResp.Data.Items))
	for _, item := range listResp.Data.Items {
		subjects[item.Subject] = struct{}{}
	}
	if _, ok := subjects["同部门跟进"]; !ok {
		t.Fatalf("expected same department followup to be visible, got %+v", listResp.Data.Items)
	}
	if _, ok := subjects["当前负责人跟进"]; !ok {
		t.Fatalf("expected actor followup to be visible, got %+v", listResp.Data.Items)
	}
	if _, ok := subjects["跨部门跟进"]; ok {
		t.Fatalf("did not expect other department followup to be visible, got %+v", listResp.Data.Items)
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

	cfg := &config.Config{
		App: config.AppConfig{
			Name: "ez-admin-test",
			Env:  "test",
		},
		Auth: config.AuthConfig{
			JWTSecret:      strings.Repeat("t", 32),
			AccessTokenTTL: 3600,
			Issuer:         "ez-admin-test",
		},
		Upload: config.UploadConfig{
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
		&model.Customer{},
		&model.CustomerFollowUp{},
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

func mustCreateUserRole(t *testing.T, db *gorm.DB, userID uint, roleID uint) {
	t.Helper()
	if err := db.Create(&model.UserRole{UserID: userID, RoleID: roleID}).Error; err != nil {
		t.Fatalf("create user role binding user=%d role=%d: %v", userID, roleID, err)
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
