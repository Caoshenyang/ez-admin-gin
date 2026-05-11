//go:build integration

package rbac

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"ez-admin-gin/server/tests/testutil"
)

// responseBody is the unified API response envelope.
type responseBody struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

type userListResponse struct {
	Items []struct {
		ID           uint   `json:"id"`
		Username     string `json:"username"`
		DepartmentID uint   `json:"department_id"`
	} `json:"items"`
	Total int64 `json:"total"`
}

// decodeUserList decodes the API envelope and extracts the user list from the data field.
func decodeUserList(t *testing.T, resp *http.Response) userListResponse {
	t.Helper()
	var envelope responseBody
	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		t.Fatalf("decode response envelope: %v", err)
	}
	if envelope.Code != 0 {
		t.Fatalf("API error: code=%d message=%s", envelope.Code, envelope.Message)
	}
	var list userListResponse
	if err := json.Unmarshal(envelope.Data, &list); err != nil {
		t.Fatalf("unmarshal user list: %v", err)
	}
	return list
}

// Common permission grants for datascope tests.
var usersListPerm = [][2]string{
	{"/api/v1/system/users", "GET"},
	{"/api/v1/system/departments", "GET"},
}

// --- API-level permission tests (previously implemented) ---

func TestUnauthenticatedAccessToSystemEndpoint(t *testing.T) {
	app := testutil.NewTestApp(t)
	defer app.Close(t)

	req, _ := http.NewRequest(http.MethodGet, app.URL("/api/v1/system/users"), nil)

	resp, err := app.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401 for unauthenticated request to system endpoint", resp.StatusCode)
	}
}

func TestPermissionDeniedWithoutRole(t *testing.T) {
	app := testutil.NewTestApp(t)
	defer app.Close(t)

	app.CleanupTestData(t)
	app.SeedAdmin(t, "admin", "admin123", "Admin")
	adminToken := app.LoginAs(t, "admin", "admin123")

	restrictedToken := app.SeedRestrictedUser(t, adminToken, "noperm_user", "password123", "No Perm User")

	req := app.AuthRequest(http.MethodGet, "/api/v1/system/users", restrictedToken, nil)

	resp, err := app.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusForbidden {
		var result responseBody
		json.NewDecoder(resp.Body).Decode(&result)
		t.Fatalf("status = %d, want 403 for user without permissions; code=%d msg=%s",
			resp.StatusCode, result.Code, result.Message)
	}
}

func TestPermissionGrantedWithRole(t *testing.T) {
	app := testutil.NewTestApp(t)
	defer app.Close(t)

	app.CleanupTestData(t)
	app.SeedAdmin(t, "admin", "admin123", "Admin")
	adminToken := app.LoginAs(t, "admin", "admin123")

	permittedToken := app.SeedUserWithPermissions(t, adminToken,
		"health_viewer", "viewer123", "Health Viewer",
		"test_health_viewer", "Test Health Viewer",
		[][2]string{
			{"/api/v1/system/health", "GET"},
		},
	)
	app.ReloadPolicies(t)

	req := app.AuthRequest(http.MethodGet, "/api/v1/system/health", permittedToken, nil)

	resp, err := app.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("status = %d, want 200 for permitted user; body: %s", resp.StatusCode, string(body))
	}
}

func TestAdminCanAccessAllEndpoints(t *testing.T) {
	app := testutil.NewTestApp(t)
	defer app.Close(t)

	app.CleanupTestData(t)
	app.SeedAdmin(t, "admin", "admin123", "Admin")
	adminToken := app.LoginAs(t, "admin", "admin123")

	req := app.AuthRequest(http.MethodGet, "/api/v1/system/users", adminToken, nil)

	resp, err := app.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200 for admin user", resp.StatusCode)
	}
}

func TestHTTPMethodPermissionDifferentiation(t *testing.T) {
	app := testutil.NewTestApp(t)
	defer app.Close(t)

	app.CleanupTestData(t)
	app.SeedAdmin(t, "admin", "admin123", "Admin")
	adminToken := app.LoginAs(t, "admin", "admin123")

	viewerToken := app.SeedUserWithPermissions(t, adminToken,
		"role_viewer", "viewer123", "Role Viewer",
		"test_role_getter", "Test Role Getter",
		[][2]string{
			{"/api/v1/system/roles", "GET"},
		},
	)
	app.ReloadPolicies(t)

	getReq := app.AuthRequest(http.MethodGet, "/api/v1/system/roles", viewerToken, nil)
	getResp, err := app.Do(getReq)
	if err != nil {
		t.Fatalf("GET request failed: %v", err)
	}
	defer getResp.Body.Close()

	if getResp.StatusCode != http.StatusOK {
		t.Fatalf("GET /system/roles status = %d, want 200", getResp.StatusCode)
	}

	postBody := `{"code":"should_fail","name":"Should Fail","sort":99,"data_scope":"self","status":1}`
	postReq := app.AuthRequest(http.MethodPost, "/api/v1/system/roles", viewerToken, strings.NewReader(postBody))
	postResp, err := app.Do(postReq)
	if err != nil {
		t.Fatalf("POST request failed: %v", err)
	}
	defer postResp.Body.Close()

	if postResp.StatusCode != http.StatusForbidden {
		t.Fatalf("POST /system/roles status = %d, want 403 (no POST permission)", postResp.StatusCode)
	}
}

// --- Data scope tests ---

// seedDepartmentTree creates a 3-level department tree for datascope tests:
//
//	hq (总部)
//	├── tech (技术部)
//	└── mkt (市场部)
//
// Returns the three department IDs.
func seedDepartmentTree(t *testing.T, app *testutil.TestApp, adminToken string) (hqID, techID, mktID uint) {
	t.Helper()
	hqID = app.SeedDepartment(t, adminToken, 0, "Test HQ", "test_hq")
	techID = app.SeedDepartment(t, adminToken, hqID, "Test Tech", "test_tech")
	mktID = app.SeedDepartment(t, adminToken, hqID, "Test Marketing", "test_mkt")
	return
}

// seedDepartmentTreeWithChild creates a deeper tree for dept_and_children tests:
//
//	hq (总部)
//	├── tech (技术部)
//	│   └── frontend (前端组)
//	└── mkt (市场部)
func seedDepartmentTreeWithChild(t *testing.T, app *testutil.TestApp, adminToken string) (hqID, techID, mktID, frontendID uint) {
	t.Helper()
	hqID, techID, mktID = seedDepartmentTree(t, app, adminToken)
	frontendID = app.SeedDepartment(t, adminToken, techID, "Test Frontend", "test_frontend")
	return
}

// TestDataScopeAll verifies that a user with data_scope=all can see all users
// across all departments.
func TestDataScopeAll(t *testing.T) {
	app := testutil.NewTestApp(t)
	defer app.Close(t)

	app.CleanupTestData(t)
	app.SeedAdmin(t, "admin", "admin123", "Admin")
	adminToken := app.LoginAs(t, "admin", "admin123")

	hqID, techID, mktID := seedDepartmentTree(t, app, adminToken)

	// Create a user in HQ with data_scope=all.
	allToken := app.SeedScopedUser(t, adminToken,
		"ds_all_user", "pass1234", "DS All",
		"test_ds_all", "Test DS All", "all", hqID, usersListPerm,
	)
	// Create users in different departments.
	app.SeedScopedUser(t, adminToken,
		"ds_tech_user", "pass1234", "Tech User",
		"test_ds_tech", "Test DS Tech", "self", techID, usersListPerm,
	)
	app.SeedScopedUser(t, adminToken,
		"ds_mkt_user", "pass1234", "Mkt User",
		"test_ds_mkt", "Test DS Mkt", "self", mktID, usersListPerm,
	)
	app.ReloadPolicies(t)

	// User with data_scope=all should see all users (admin + ds_all + tech + mkt = 4).
	req := app.AuthRequest(http.MethodGet, "/api/v1/system/users", allToken, nil)
	resp, err := app.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("status = %d, want 200; body: %s", resp.StatusCode, string(body))
	}

	var list userListResponse
	list = decodeUserList(t, resp)

	// Should see at least: admin (ID=1, dept=0), self (ds_all_user, hq), tech_user (tech), mkt_user (mkt)
	if list.Total < 4 {
		t.Errorf("data_scope=all returned %d users, want at least 4", list.Total)
	}
}

// TestDataScopeSelf verifies that a user with data_scope=self can only see
// their own user record.
func TestDataScopeSelf(t *testing.T) {
	app := testutil.NewTestApp(t)
	defer app.Close(t)

	app.CleanupTestData(t)
	app.SeedAdmin(t, "admin", "admin123", "Admin")
	adminToken := app.LoginAs(t, "admin", "admin123")

	hqID, techID, _ := seedDepartmentTree(t, app, adminToken)
	_ = hqID

	// Create a self-scoped user in tech department.
	selfToken := app.SeedScopedUser(t, adminToken,
		"ds_self_user", "pass1234", "DS Self",
		"test_ds_self", "Test DS Self", "self", techID, usersListPerm,
	)
	// Create another user in same department.
	app.SeedScopedUser(t, adminToken,
		"ds_tech_user2", "pass1234", "Tech User 2",
		"test_ds_tech2", "Test DS Tech 2", "self", techID, usersListPerm,
	)
	app.ReloadPolicies(t)

	req := app.AuthRequest(http.MethodGet, "/api/v1/system/users", selfToken, nil)
	resp, err := app.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("status = %d, want 200; body: %s", resp.StatusCode, string(body))
	}

	var list userListResponse
	list = decodeUserList(t, resp)

	// Should see only 1 user (themselves), not tech_user2 or admin (who is in dept 0).
	if list.Total != 1 {
	 usernames := make([]string, len(list.Items))
	 for i, u := range list.Items {
	  usernames[i] = u.Username
	 }
		t.Errorf("data_scope=self returned %d users (%v), want exactly 1", list.Total, usernames)
	}
}

// TestDataScopeDept verifies that a user with data_scope=dept can see users
// in their own department only (not sub-departments).
func TestDataScopeDept(t *testing.T) {
	app := testutil.NewTestApp(t)
	defer app.Close(t)

	app.CleanupTestData(t)
	app.SeedAdmin(t, "admin", "admin123", "Admin")
	adminToken := app.LoginAs(t, "admin", "admin123")

	_, techID, mktID := seedDepartmentTree(t, app, adminToken)

	// Create a dept-scoped user in tech department.
	deptToken := app.SeedScopedUser(t, adminToken,
		"ds_dept_user", "pass1234", "DS Dept",
		"test_ds_dept", "Test DS Dept", "dept", techID, usersListPerm,
	)
	// Create another user in tech department.
	app.SeedScopedUser(t, adminToken,
		"ds_tech_colleague", "pass1234", "Tech Colleague",
		"test_ds_tech_col", "Test DS Tech Col", "self", techID, usersListPerm,
	)
	// Create a user in marketing department.
	app.SeedScopedUser(t, adminToken,
		"ds_mkt_other", "pass1234", "Mkt Other",
		"test_ds_mkt_other", "Test DS Mkt Other", "self", mktID, usersListPerm,
	)
	app.ReloadPolicies(t)

	req := app.AuthRequest(http.MethodGet, "/api/v1/system/users", deptToken, nil)
	resp, err := app.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("status = %d, want 200; body: %s", resp.StatusCode, string(body))
	}

	var list userListResponse
	list = decodeUserList(t, resp)

	// Should see 2 users: self (ds_dept_user) + colleague (ds_tech_colleague).
	// Should NOT see mkt user or admin (dept=0).
	for _, u := range list.Items {
		if u.DepartmentID != techID {
			t.Errorf("data_scope=dept returned user %q in department %d, expected only department %d",
				u.Username, u.DepartmentID, techID)
		}
	}
	if list.Total != 2 {
		t.Errorf("data_scope=dept returned %d users, want 2", list.Total)
	}
}

// TestDataScopeDeptAndChildren verifies that a user with data_scope=dept_and_children
// can see users in their own department and all sub-departments.
func TestDataScopeDeptAndChildren(t *testing.T) {
	app := testutil.NewTestApp(t)
	defer app.Close(t)

	app.CleanupTestData(t)
	app.SeedAdmin(t, "admin", "admin123", "Admin")
	adminToken := app.LoginAs(t, "admin", "admin123")

	_, techID, mktID, frontendID := seedDepartmentTreeWithChild(t, app, adminToken)

	// Create a dept_and_children scoped user in tech department.
	treeToken := app.SeedScopedUser(t, adminToken,
		"ds_tree_user", "pass1234", "DS Tree",
		"test_ds_tree", "Test DS Tree", "dept_and_children", techID, usersListPerm,
	)
	// Create a user in tech department.
	app.SeedScopedUser(t, adminToken,
		"ds_tech_peer", "pass1234", "Tech Peer",
		"test_ds_tech_peer", "Test DS Tech Peer", "self", techID, usersListPerm,
	)
	// Create a user in frontend sub-department.
	app.SeedScopedUser(t, adminToken,
		"ds_fe_dev", "pass1234", "FE Dev",
		"test_ds_fe_dev", "Test DS FE Dev", "self", frontendID, usersListPerm,
	)
	// Create a user in marketing department (should NOT be visible).
	app.SeedScopedUser(t, adminToken,
		"ds_mkt_stranger", "pass1234", "Mkt Stranger",
		"test_ds_mkt_str", "Test DS Mkt Stranger", "self", mktID, usersListPerm,
	)
	app.ReloadPolicies(t)

	req := app.AuthRequest(http.MethodGet, "/api/v1/system/users", treeToken, nil)
	resp, err := app.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("status = %d, want 200; body: %s", resp.StatusCode, string(body))
	}

	list := decodeUserList(t, resp)

	// Should see 3 users: self (tech), tech_peer (tech), fe_dev (frontend sub-dept).
	// Should NOT see mkt_stranger or admin.
	visibleDepts := make(map[uint]int)
	for _, u := range list.Items {
		visibleDepts[u.DepartmentID]++
	}
	if visibleDepts[mktID] > 0 {
		t.Error("data_scope=dept_and_children should not see marketing department users")
	}
	if list.Total != 3 {
		t.Errorf("data_scope=dept_and_children returned %d users, want 3", list.Total)
	}
}

// TestDataScopeCustomDept verifies that a user with data_scope=custom_dept can
// only see users in the explicitly specified departments.
func TestDataScopeCustomDept(t *testing.T) {
	app := testutil.NewTestApp(t)
	defer app.Close(t)

	app.CleanupTestData(t)
	app.SeedAdmin(t, "admin", "admin123", "Admin")
	adminToken := app.LoginAs(t, "admin", "admin123")

	_, techID, mktID, frontendID := seedDepartmentTreeWithChild(t, app, adminToken)

	// Create a custom_dept user whose role grants access only to tech and frontend departments.
	// The user themselves sits in mkt department, but custom_dept overrides their own dept.
	customToken := app.SeedCustomDeptUser(t, adminToken,
		"ds_custom_user", "pass1234", "DS Custom",
		"test_ds_custom", "Test DS Custom",
		[]uint{techID, frontendID}, // explicitly grant access to tech + frontend
		mktID,                      // user's own department (should NOT be visible)
		usersListPerm,
	)
	// Create users in different departments.
	app.SeedScopedUser(t, adminToken,
		"ds_tech_target", "pass1234", "Tech Target",
		"test_ds_tech_tgt", "Test DS Tech Tgt", "self", techID, usersListPerm,
	)
	app.SeedScopedUser(t, adminToken,
		"ds_mkt_target", "pass1234", "Mkt Target",
		"test_ds_mkt_tgt", "Test DS Mkt Tgt", "self", mktID, usersListPerm,
	)
	app.SeedScopedUser(t, adminToken,
		"ds_fe_target", "pass1234", "FE Target",
		"test_ds_fe_tgt", "Test DS FE Tgt", "self", frontendID, usersListPerm,
	)
	app.ReloadPolicies(t)

	req := app.AuthRequest(http.MethodGet, "/api/v1/system/users", customToken, nil)
	resp, err := app.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("status = %d, want 200; body: %s", resp.StatusCode, string(body))
	}

	list := decodeUserList(t, resp)

	// Should see exactly 2 users: tech_target (tech) + fe_target (frontend).
	// Should NOT see mkt_target (mkt) or admin (dept=0) or self (mkt).
	visibleDepts := make(map[uint]int)
	for _, u := range list.Items {
		visibleDepts[u.DepartmentID]++
	}
	if visibleDepts[mktID] > 0 {
		t.Error("data_scope=custom_dept should not see marketing department users")
	}
	if list.Total != 2 {
		t.Errorf("data_scope=custom_dept returned %d users, want 2 (tech + frontend)", list.Total)
	}
}

// TestDataScopeDefaultDeny verifies that a user with no matching data scope
// conditions sees zero users (returns empty list, not all users).
func TestDataScopeDefaultDeny(t *testing.T) {
	app := testutil.NewTestApp(t)
	defer app.Close(t)

	app.CleanupTestData(t)
	app.SeedAdmin(t, "admin", "admin123", "Admin")
	adminToken := app.LoginAs(t, "admin", "admin123")

	_, techID, _ := seedDepartmentTree(t, app, adminToken)

	// Create a self-scoped user but assign them to department 0 (no department).
	// Their data_scope is "self", but since ownerColumn for user list is the user's
	// own ID, they should still see only themselves.
	// To test default deny, we need a user with data_scope="dept" but department_id=0.
	// In that case, actor.DepartmentID == 0, so the dept condition won't fire,
	// and since there's no "self" or other scope, we get 1=0 (deny all).
	noDeptToken := app.SeedScopedUser(t, adminToken,
		"ds_nodept_user", "pass1234", "DS No Dept",
		"test_ds_nodept", "Test DS No Dept", "dept", 0, usersListPerm,
	)
	// Create another user in tech.
	app.SeedScopedUser(t, adminToken,
		"ds_tech_someone", "pass1234", "Tech Someone",
		"test_ds_tech_some", "Test DS Tech Some", "self", techID, usersListPerm,
	)
	app.ReloadPolicies(t)

	req := app.AuthRequest(http.MethodGet, "/api/v1/system/users", noDeptToken, nil)
	resp, err := app.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("status = %d, want 200; body: %s", resp.StatusCode, string(body))
	}

	var list userListResponse
	list = decodeUserList(t, resp)

	// User with data_scope=dept but department_id=0 should see 0 users (default deny).
	if list.Total != 0 {
		t.Errorf("data_scope=dept with department_id=0 returned %d users, want 0 (default deny)", list.Total)
	}
}

// TestMultiRolePermissionUnion verifies that a user with multiple roles gets the
// union of all permissions across roles.
func TestMultiRolePermissionUnion(t *testing.T) {
	app := testutil.NewTestApp(t)
	defer app.Close(t)

	app.CleanupTestData(t)
	app.SeedAdmin(t, "admin", "admin123", "Admin")
	adminToken := app.LoginAs(t, "admin", "admin123")

	// Role A: only GET /system/roles
	roleAID := createRoleWithPerms(t, app, adminToken,
		"test_multi_role_a", "Test Multi Role A",
		[][2]string{{"/api/v1/system/roles", "GET"}},
	)
	// Role B: only GET /system/users
	roleBID := createRoleWithPerms(t, app, adminToken,
		"test_multi_role_b", "Test Multi Role B",
		[][2]string{{"/api/v1/system/users", "GET"}},
	)

	// Create a user with BOTH roles.
	userBody := fmt.Sprintf(
		`{"username":"multi_role_user","password":"pass1234","nickname":"Multi Role","department_id":0,"status":1,"role_ids":[%d,%d]}`,
		roleAID, roleBID,
	)
	userReq := app.AuthRequest(http.MethodPost, "/api/v1/system/users", adminToken, strings.NewReader(userBody))
	userResp, err := app.Do(userReq)
	if err != nil {
		t.Fatalf("create user request failed: %v", err)
	}
	defer userResp.Body.Close()
	if userResp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(userResp.Body)
		t.Fatalf("create user failed: status %d, body: %s", userResp.StatusCode, string(body))
	}

	multiToken := app.LoginAs(t, "multi_role_user", "pass1234")
	app.ReloadPolicies(t)

	// Should have GET /system/roles (from Role A).
	rolesReq := app.AuthRequest(http.MethodGet, "/api/v1/system/roles", multiToken, nil)
	rolesResp, err := app.Do(rolesReq)
	if err != nil {
		t.Fatalf("GET /system/roles request failed: %v", err)
	}
	defer rolesResp.Body.Close()
	if rolesResp.StatusCode != http.StatusOK {
		t.Fatalf("GET /system/roles status = %d, want 200 (permission from Role A)", rolesResp.StatusCode)
	}

	// Should have GET /system/users (from Role B).
	usersReq := app.AuthRequest(http.MethodGet, "/api/v1/system/users", multiToken, nil)
	usersResp, err := app.Do(usersReq)
	if err != nil {
		t.Fatalf("GET /system/users request failed: %v", err)
	}
	defer usersResp.Body.Close()
	if usersResp.StatusCode != http.StatusOK {
		t.Fatalf("GET /system/users status = %d, want 200 (permission from Role B)", usersResp.StatusCode)
	}

	// Should NOT have POST /system/roles (neither role grants it).
	postBody := `{"code":"should_fail","name":"Should Fail","sort":99,"data_scope":"self","status":1}`
	postReq := app.AuthRequest(http.MethodPost, "/api/v1/system/roles", multiToken, strings.NewReader(postBody))
	postResp, err := app.Do(postReq)
	if err != nil {
		t.Fatalf("POST /system/roles request failed: %v", err)
	}
	defer postResp.Body.Close()
	if postResp.StatusCode != http.StatusForbidden {
		t.Fatalf("POST /system/roles status = %d, want 403 (no role grants POST)", postResp.StatusCode)
	}
}

// createRoleWithPerms creates a role and assigns permissions, returning the role ID.
func createRoleWithPerms(t *testing.T, app *testutil.TestApp, adminToken, roleCode, roleName string, permissions [][2]string) uint {
	t.Helper()

	roleBody := fmt.Sprintf(
		`{"code":"%s","name":"%s","sort":95,"data_scope":"self","status":1}`,
		roleCode, roleName,
	)
	roleReq := app.AuthRequest(http.MethodPost, "/api/v1/system/roles", adminToken, strings.NewReader(roleBody))
	roleResp, err := app.Do(roleReq)
	if err != nil {
		t.Fatalf("create role %s request failed: %v", roleCode, err)
	}
	defer roleResp.Body.Close()
	if roleResp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(roleResp.Body)
		t.Fatalf("create role %s failed: status %d, body: %s", roleCode, roleResp.StatusCode, string(body))
	}

	var roleResult struct {
		Code int `json:"code"`
		Data struct {
			ID uint `json:"id"`
		} `json:"data"`
	}
	if err := json.NewDecoder(roleResp.Body).Decode(&roleResult); err != nil {
		t.Fatalf("decode role %s response: %v", roleCode, err)
	}
	if roleResult.Data.ID == 0 {
		t.Fatalf("created role %s has ID 0", roleCode)
	}
	roleID := roleResult.Data.ID

	if len(permissions) > 0 {
		permItems := make([]string, 0, len(permissions))
		for _, p := range permissions {
			permItems = append(permItems, fmt.Sprintf(`{"path":"%s","method":"%s"}`, p[0], p[1]))
		}
		permBody := fmt.Sprintf(`{"permissions":[%s]}`, strings.Join(permItems, ","))
		permReq := app.AuthRequest(http.MethodPost, fmt.Sprintf("/api/v1/system/roles/%d/permissions", roleID), adminToken, strings.NewReader(permBody))
		permResp, err := app.Do(permReq)
		if err != nil {
			t.Fatalf("assign permissions to role %s request failed: %v", roleCode, err)
		}
		defer permResp.Body.Close()
		if permResp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(permResp.Body)
			t.Fatalf("assign permissions to role %s failed: status %d, body: %s", roleCode, permResp.StatusCode, string(body))
		}
	}

	return roleID
}
