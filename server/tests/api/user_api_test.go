//go:build integration

package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"ez-admin-gin/server/tests/testutil"
)

// mustCreateUser creates a user via API and returns the user ID.
func mustCreateUser(t *testing.T, app *testutil.TestApp, token, username, password, nickname string) uint {
	t.Helper()
	body := fmt.Sprintf(
		`{"username":"%s","password":"%s","nickname":"%s","department_id":0,"status":1}`,
		username, password, nickname,
	)
	req := app.AuthRequest(http.MethodPost, "/api/v1/system/users", token, strings.NewReader(body))
	resp, err := app.Do(req)
	if err != nil {
		t.Fatalf("create user request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("create user failed: status %d", resp.StatusCode)
	}
	var result struct {
		Code int `json:"code"`
		Data struct {
			ID uint `json:"id"`
		} `json:"data"`
	}
	app.DecodeResponse(t, resp, &result)
	if result.Data.ID == 0 {
		t.Fatal("created user has ID 0")
	}
	return result.Data.ID
}

// TestCreateUser verifies admin can create a new user and the response
// contains the expected fields.
func TestCreateUser(t *testing.T) {
	app := testutil.NewTestApp(t)
	defer app.Close(t)

	app.SeedAdmin(t, testAdminUser, testAdminPassword, testAdminNickname)
	app.CleanupTestData(t)
	token := app.LoginAs(t, testAdminUser, testAdminPassword)

	body := `{"username":"testuser1","password":"pass1234","nickname":"Test User","department_id":0,"status":1}`
	req := app.AuthRequest(http.MethodPost, "/api/v1/system/users", token, strings.NewReader(body))
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

	var user struct {
		ID       uint   `json:"id"`
		Username string `json:"username"`
		Nickname string `json:"nickname"`
	}
	if err := json.Unmarshal(result.Data, &user); err != nil {
		t.Fatalf("unmarshal user: %v", err)
	}
	if user.ID == 0 {
		t.Error("user ID should not be 0")
	}
	if user.Username != "testuser1" {
		t.Errorf("username = %q, want %q", user.Username, "testuser1")
	}
}

// TestCreateUserDuplicateUsername verifies creating a user with an existing
// username returns an error.
func TestCreateUserDuplicateUsername(t *testing.T) {
	app := testutil.NewTestApp(t)
	defer app.Close(t)

	app.SeedAdmin(t, testAdminUser, testAdminPassword, testAdminNickname)
	app.CleanupTestData(t)
	token := app.LoginAs(t, testAdminUser, testAdminPassword)

	// Create first user.
	body := `{"username":"dupuser","password":"pass1234","nickname":"First","department_id":0,"status":1}`
	req := app.AuthRequest(http.MethodPost, "/api/v1/system/users", token, strings.NewReader(body))
	resp, err := app.Do(req)
	if err != nil {
		t.Fatalf("first create request failed: %v", err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("first create: status = %d, want 200", resp.StatusCode)
	}

	// Try duplicate username.
	req2 := app.AuthRequest(http.MethodPost, "/api/v1/system/users", token, strings.NewReader(body))
	resp2, err := app.Do(req2)
	if err != nil {
		t.Fatalf("duplicate request failed: %v", err)
	}
	defer resp2.Body.Close()

	if resp2.StatusCode == http.StatusOK {
		t.Fatal("duplicate username should not return 200")
	}
}

// TestListUsers verifies the user list endpoint returns paginated results.
func TestListUsers(t *testing.T) {
	app := testutil.NewTestApp(t)
	defer app.Close(t)

	app.SeedAdmin(t, testAdminUser, testAdminPassword, testAdminNickname)
	app.CleanupTestData(t)
	token := app.LoginAs(t, testAdminUser, testAdminPassword)

	req := app.AuthRequest(http.MethodGet, "/api/v1/system/users?page=1&page_size=10", token, nil)
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
			ID       uint   `json:"id"`
			Username string `json:"username"`
		} `json:"items"`
		Total int64 `json:"total"`
	}
	if err := json.Unmarshal(result.Data, &list); err != nil {
		t.Fatalf("unmarshal list: %v", err)
	}
	if list.Total < 1 {
		t.Error("total should be at least 1 (admin user)")
	}
}

// TestUpdateUser verifies updating a user's nickname works.
func TestUpdateUser(t *testing.T) {
	app := testutil.NewTestApp(t)
	defer app.Close(t)

	app.SeedAdmin(t, testAdminUser, testAdminPassword, testAdminNickname)
	app.CleanupTestData(t)
	token := app.LoginAs(t, testAdminUser, testAdminPassword)

	// Create a user first.
	userID := mustCreateUser(t, app, token, "updateme", "pass1234", "Before Update")

	body := `{"nickname":"After Update","department_id":0,"status":1}`
	req := app.AuthRequest(http.MethodPost, fmt.Sprintf("/api/v1/system/users/%d/update", userID), token, strings.NewReader(body))
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

	var user struct {
		Nickname string `json:"nickname"`
	}
	if err := json.Unmarshal(result.Data, &user); err != nil {
		t.Fatalf("unmarshal user: %v", err)
	}
	if user.Nickname != "After Update" {
		t.Errorf("nickname = %q, want %q", user.Nickname, "After Update")
	}

	_ = userID // suppress unused warning
}

// TestUpdateUserStatus verifies toggling user status works.
func TestUpdateUserStatus(t *testing.T) {
	app := testutil.NewTestApp(t)
	defer app.Close(t)

	app.SeedAdmin(t, testAdminUser, testAdminPassword, testAdminNickname)
	app.CleanupTestData(t)
	token := app.LoginAs(t, testAdminUser, testAdminPassword)

	// Create a user and capture the ID.
	userID := mustCreateUser(t, app, token, "statustest", "pass1234", "Status Test")

	// Disable user (status=2).
	body := `{"status":2}`
	req := app.AuthRequest(http.MethodPost, fmt.Sprintf("/api/v1/system/users/%d/status", userID), token, strings.NewReader(body))
	resp, err := app.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}
}

// TestCreateUserMissingFields verifies the API rejects incomplete requests.
func TestCreateUserMissingFields(t *testing.T) {
	app := testutil.NewTestApp(t)
	defer app.Close(t)

	app.SeedAdmin(t, testAdminUser, testAdminPassword, testAdminNickname)
	app.CleanupTestData(t)
	token := app.LoginAs(t, testAdminUser, testAdminPassword)

	// Missing password.
	body := `{"username":"nopass","nickname":"No Password","department_id":0,"status":1}`
	req := app.AuthRequest(http.MethodPost, "/api/v1/system/users", token, strings.NewReader(body))
	resp, err := app.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		t.Fatal("missing password should not return 200")
	}
}
