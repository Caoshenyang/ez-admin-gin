//go:build integration

package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"ez-admin-gin/server/tests/testutil"
)

const (
	testAdminUser     = "admin"
	testAdminPassword = "admin123"
	testAdminNickname = "Test Admin"
)

// responseBody is the unified API response envelope.
type responseBody struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

func TestLoginSuccess(t *testing.T) {
	app := testutil.NewTestApp(t)
	defer app.Close(t)

	app.SeedAdmin(t, testAdminUser, testAdminPassword, testAdminNickname)

	body := `{"username":"admin","password":"admin123"}`
	req, _ := http.NewRequest(http.MethodPost, app.URL("/api/v1/auth/login"), strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}

	var result responseBody
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if result.Code != 0 {
		t.Errorf("code = %d, want 0 (success)", result.Code)
	}
	if result.Message != "ok" {
		t.Errorf("message = %q, want %q", result.Message, "ok")
	}

	var loginData struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		UserID      uint   `json:"user_id"`
		Username    string `json:"username"`
	}
	if err := json.Unmarshal(result.Data, &loginData); err != nil {
		t.Fatalf("unmarshal login data: %v", err)
	}
	if loginData.AccessToken == "" {
		t.Error("access_token is empty")
	}
	if loginData.Username != "admin" {
		t.Errorf("username = %q, want %q", loginData.Username, "admin")
	}
}

func TestLoginWrongPassword(t *testing.T) {
	app := testutil.NewTestApp(t)
	defer app.Close(t)

	app.SeedAdmin(t, testAdminUser, testAdminPassword, testAdminNickname)

	body := `{"username":"admin","password":"wrong-password"}`
	req, _ := http.NewRequest(http.MethodPost, app.URL("/api/v1/auth/login"), strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		t.Fatal("login with wrong password should not return 200")
	}

	var result responseBody
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if result.Code == 0 {
		t.Error("code should not be 0 (success) for wrong password")
	}
}

func TestUnauthorizedAccessWithoutToken(t *testing.T) {
	app := testutil.NewTestApp(t)
	defer app.Close(t)

	req, _ := http.NewRequest(http.MethodGet, app.URL("/api/v1/auth/me"), nil)

	resp, err := app.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401 for unauthenticated request", resp.StatusCode)
	}

	var result responseBody
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if result.Code != 40100 {
		t.Errorf("code = %d, want 40100 (unauthorized)", result.Code)
	}
}

func TestUnauthorizedAccessWithInvalidToken(t *testing.T) {
	app := testutil.NewTestApp(t)
	defer app.Close(t)

	req, _ := http.NewRequest(http.MethodGet, app.URL("/api/v1/auth/me"), nil)
	req.Header.Set("Authorization", "Bearer invalid-token-here")

	resp, err := app.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401 for invalid token", resp.StatusCode)
	}
}

// TODO: TestLoginRateLimiting — verify login rate limiting kicks in after N failed attempts.
// Requires either Redis or careful request sequencing. Deferred to next phase.
