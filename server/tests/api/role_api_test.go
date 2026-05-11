//go:build integration

package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"ez-admin-gin/server/tests/testutil"
)

// TestCreateRole verifies admin can create a new role.
func TestCreateRole(t *testing.T) {
	app := testutil.NewTestApp(t)
	defer app.Close(t)

	app.SeedAdmin(t, testAdminUser, testAdminPassword, testAdminNickname)
	app.CleanupTestData(t)
	token := app.LoginAs(t, testAdminUser, testAdminPassword)

	body := `{"code":"test_editor","name":"Test Editor","sort":1,"data_scope":"self","status":1}`
	req := app.AuthRequest(http.MethodPost, "/api/v1/system/roles", token, strings.NewReader(body))
	resp, err := app.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}

	var result responseBody
	app.DecodeResponse(t, resp, &result)
	if result.Code != 0 {
		t.Fatalf("code = %d, want 0", result.Code)
	}

	var role struct {
		ID   uint   `json:"id"`
		Code string `json:"code"`
		Name string `json:"name"`
	}
	if err := json.Unmarshal(result.Data, &role); err != nil {
		t.Fatalf("unmarshal role: %v", err)
	}
	if role.ID == 0 {
		t.Error("role ID should not be 0")
	}
	if role.Code != "test_editor" {
		t.Errorf("code = %q, want %q", role.Code, "test_editor")
	}
}

// TestCreateRoleDuplicateCode verifies creating a role with an existing
// code returns an error.
func TestCreateRoleDuplicateCode(t *testing.T) {
	app := testutil.NewTestApp(t)
	defer app.Close(t)

	app.SeedAdmin(t, testAdminUser, testAdminPassword, testAdminNickname)
	app.CleanupTestData(t)
	token := app.LoginAs(t, testAdminUser, testAdminPassword)

	body := `{"code":"dup_role","name":"First Role","sort":1,"data_scope":"self","status":1}`
	req := app.AuthRequest(http.MethodPost, "/api/v1/system/roles", token, strings.NewReader(body))
	resp, err := app.Do(req)
	if err != nil {
		t.Fatalf("first create request failed: %v", err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("first create: status = %d, want 200", resp.StatusCode)
	}

	// Duplicate code.
	req2 := app.AuthRequest(http.MethodPost, "/api/v1/system/roles", token, strings.NewReader(body))
	resp2, err := app.Do(req2)
	if err != nil {
		t.Fatalf("duplicate request failed: %v", err)
	}
	defer resp2.Body.Close()

	if resp2.StatusCode == http.StatusOK {
		t.Fatal("duplicate role code should not return 200")
	}
}

// TestListRoles verifies the role list endpoint returns paginated results.
func TestListRoles(t *testing.T) {
	app := testutil.NewTestApp(t)
	defer app.Close(t)

	app.SeedAdmin(t, testAdminUser, testAdminPassword, testAdminNickname)
	app.CleanupTestData(t)
	token := app.LoginAs(t, testAdminUser, testAdminPassword)

	req := app.AuthRequest(http.MethodGet, "/api/v1/system/roles?page=1&page_size=10", token, nil)
	resp, err := app.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}

	var result responseBody
	app.DecodeResponse(t, resp, &result)
	if result.Code != 0 {
		t.Fatalf("code = %d, want 0", result.Code)
	}

	var list struct {
		Items []struct {
			ID   uint   `json:"id"`
			Code string `json:"code"`
		} `json:"items"`
		Total int64 `json:"total"`
	}
	if err := json.Unmarshal(result.Data, &list); err != nil {
		t.Fatalf("unmarshal list: %v", err)
	}
	if list.Total < 1 {
		t.Error("total should be at least 1 (super_admin)")
	}
}

// TestUpdateRole verifies updating a role's name works.
func TestUpdateRole(t *testing.T) {
	app := testutil.NewTestApp(t)
	defer app.Close(t)

	app.SeedAdmin(t, testAdminUser, testAdminPassword, testAdminNickname)
	app.CleanupTestData(t)
	token := app.LoginAs(t, testAdminUser, testAdminPassword)

	// Create a role.
	roleID := mustCreateRole(t, app, token, "upd_role", "Original Name")

	body := `{"name":"Updated Name","sort":1,"data_scope":"self","status":1}`
	req := app.AuthRequest(http.MethodPost, fmt.Sprintf("/api/v1/system/roles/%d/update", roleID), token, strings.NewReader(body))
	resp, err := app.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}

	var result responseBody
	app.DecodeResponse(t, resp, &result)
	if result.Code != 0 {
		t.Fatalf("code = %d, want 0", result.Code)
	}

	var role struct {
		Name string `json:"name"`
	}
	if err := json.Unmarshal(result.Data, &role); err != nil {
		t.Fatalf("unmarshal role: %v", err)
	}
	if role.Name != "Updated Name" {
		t.Errorf("name = %q, want %q", role.Name, "Updated Name")
	}
}

// TestUpdateRolePermissions verifies assigning permissions to a role.
func TestUpdateRolePermissions(t *testing.T) {
	app := testutil.NewTestApp(t)
	defer app.Close(t)

	app.SeedAdmin(t, testAdminUser, testAdminPassword, testAdminNickname)
	app.CleanupTestData(t)
	token := app.LoginAs(t, testAdminUser, testAdminPassword)

	roleID := mustCreateRole(t, app, token, "perm_role", "Perm Role")

	permBody := `{"permissions":[{"path":"/api/v1/system/users","method":"GET"},{"path":"/api/v1/system/users","method":"POST"}]}`
	req := app.AuthRequest(http.MethodPost, fmt.Sprintf("/api/v1/system/roles/%d/permissions", roleID), token, strings.NewReader(permBody))
	resp, err := app.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}

	var result responseBody
	app.DecodeResponse(t, resp, &result)
	if result.Code != 0 {
		t.Fatalf("code = %d, want 0", result.Code)
	}
}

// TestUpdateRoleStatus verifies toggling role status works.
func TestUpdateRoleStatus(t *testing.T) {
	app := testutil.NewTestApp(t)
	defer app.Close(t)

	app.SeedAdmin(t, testAdminUser, testAdminPassword, testAdminNickname)
	app.CleanupTestData(t)
	token := app.LoginAs(t, testAdminUser, testAdminPassword)

	roleID := mustCreateRole(t, app, token, "stat_role", "Status Role")

	body := `{"status":2}`
	req := app.AuthRequest(http.MethodPost, fmt.Sprintf("/api/v1/system/roles/%d/status", roleID), token, strings.NewReader(body))
	resp, err := app.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}
}

// mustCreateRole is a test helper that creates a role via API and returns its ID.
func mustCreateRole(t *testing.T, app *testutil.TestApp, token, code, name string) uint {
	t.Helper()

	body := fmt.Sprintf(`{"code":"%s","name":"%s","sort":1,"data_scope":"self","status":1}`, code, name)
	req := app.AuthRequest(http.MethodPost, "/api/v1/system/roles", token, strings.NewReader(body))
	resp, err := app.Do(req)
	if err != nil {
		t.Fatalf("create role request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		t.Fatalf("create role failed: status %d, body: %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		Code int `json:"code"`
		Data struct {
			ID uint `json:"id"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("decode role response: %v", err)
	}
	if result.Data.ID == 0 {
		t.Fatal("created role has ID 0")
	}
	return result.Data.ID
}
