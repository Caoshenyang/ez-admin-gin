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

// TestCreateMenu verifies admin can create a directory menu item.
func TestCreateMenu(t *testing.T) {
	app := testutil.NewTestApp(t)
	defer app.Close(t)

	app.SeedAdmin(t, testAdminUser, testAdminPassword, testAdminNickname)
	app.CleanupTestData(t)
	token := app.LoginAs(t, testAdminUser, testAdminPassword)

	body := `{"parent_id":0,"type":1,"code":"test_dir","title":"Test Directory","path":"/test","icon":"folder","sort":1,"status":1}`
	req := app.AuthRequest(http.MethodPost, "/api/v1/system/menus", token, strings.NewReader(body))
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

	var menu struct {
		ID    uint   `json:"id"`
		Code  string `json:"code"`
		Title string `json:"title"`
		Type  int    `json:"type"`
	}
	if err := json.Unmarshal(result.Data, &menu); err != nil {
		t.Fatalf("unmarshal menu: %v", err)
	}
	if menu.ID == 0 {
		t.Error("menu ID should not be 0")
	}
	if menu.Code != "test_dir" {
		t.Errorf("code = %q, want %q", menu.Code, "test_dir")
	}
}

// TestListMenus verifies the menu list endpoint returns a tree structure.
func TestListMenus(t *testing.T) {
	app := testutil.NewTestApp(t)
	defer app.Close(t)

	app.SeedAdmin(t, testAdminUser, testAdminPassword, testAdminNickname)
	app.CleanupTestData(t)
	token := app.LoginAs(t, testAdminUser, testAdminPassword)

	req := app.AuthRequest(http.MethodGet, "/api/v1/system/menus", token, nil)
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

	// Menu list is an array (tree).
	var menus json.RawMessage
	if err := json.Unmarshal(result.Data, &menus); err != nil {
		t.Fatalf("unmarshal menus: %v", err)
	}
	// Verify it's a JSON array.
	if menus[0] != '[' {
		t.Error("menu data should be a JSON array")
	}
}

// TestUpdateMenu verifies updating a menu's title works.
func TestUpdateMenu(t *testing.T) {
	app := testutil.NewTestApp(t)
	defer app.Close(t)

	app.SeedAdmin(t, testAdminUser, testAdminPassword, testAdminNickname)
	app.CleanupTestData(t)
	token := app.LoginAs(t, testAdminUser, testAdminPassword)

	menuID := mustCreateMenu(t, app, token, "upd_menu", "Original Title")

	body := `{"parent_id":0,"type":1,"title":"Updated Title","path":"/test","sort":1,"status":1}`
	req := app.AuthRequest(http.MethodPost, fmt.Sprintf("/api/v1/system/menus/%d/update", menuID), token, strings.NewReader(body))
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

	var menu struct {
		Title string `json:"title"`
	}
	if err := json.Unmarshal(result.Data, &menu); err != nil {
		t.Fatalf("unmarshal menu: %v", err)
	}
	if menu.Title != "Updated Title" {
		t.Errorf("title = %q, want %q", menu.Title, "Updated Title")
	}
}

// TestDeleteMenu verifies deleting a menu item works.
func TestDeleteMenu(t *testing.T) {
	app := testutil.NewTestApp(t)
	defer app.Close(t)

	app.SeedAdmin(t, testAdminUser, testAdminPassword, testAdminNickname)
	app.CleanupTestData(t)
	token := app.LoginAs(t, testAdminUser, testAdminPassword)

	menuID := mustCreateMenu(t, app, token, "del_menu", "To Delete")

	req := app.AuthRequest(http.MethodPost, fmt.Sprintf("/api/v1/system/menus/%d/delete", menuID), token, nil)
	resp, err := app.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}
}

// TestCreateMenuMissingFields verifies the API rejects incomplete requests.
func TestCreateMenuMissingFields(t *testing.T) {
	app := testutil.NewTestApp(t)
	defer app.Close(t)

	app.SeedAdmin(t, testAdminUser, testAdminPassword, testAdminNickname)
	app.CleanupTestData(t)
	token := app.LoginAs(t, testAdminUser, testAdminPassword)

	// Missing title and code.
	body := `{"parent_id":0,"type":1,"path":"/test","status":1}`
	req := app.AuthRequest(http.MethodPost, "/api/v1/system/menus", token, strings.NewReader(body))
	resp, err := app.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		t.Fatal("missing required fields should not return 200")
	}
}

// mustCreateMenu is a test helper that creates a menu via API and returns its ID.
func mustCreateMenu(t *testing.T, app *testutil.TestApp, token, code, title string) uint {
	t.Helper()

	body := fmt.Sprintf(`{"parent_id":0,"type":1,"code":"%s","title":"%s","path":"/test","icon":"folder","sort":1,"status":1}`, code, title)
	req := app.AuthRequest(http.MethodPost, "/api/v1/system/menus", token, strings.NewReader(body))
	resp, err := app.Do(req)
	if err != nil {
		t.Fatalf("create menu request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("create menu failed: status %d", resp.StatusCode)
	}

	var result struct {
		Code int `json:"code"`
		Data struct {
			ID uint `json:"id"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("decode menu response: %v", err)
	}
	if result.Data.ID == 0 {
		t.Fatal("created menu has ID 0")
	}
	return result.Data.ID
}
